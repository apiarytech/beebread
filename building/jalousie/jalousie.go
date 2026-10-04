/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - EPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

// Package jalousie is the port of the OSCAT BUILDING blind (Jalousie)
// control: actuators, inputs, and the blocks that set a blind's position
// and slat angle, chained from the input to the actuator. Each block of the
// chain passes the motor commands UP and DN (both true for automatic
// operation), the position PI and angle AI, and the status S_IN on to QU,
// QD, PO, AO and STATUS, unless it acts itself.
package jalousie

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/engineering"
	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/beebread/building/electrical"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// ms returns a TIME in milliseconds, as the PLC timer counts.
func ms(t iec.TIME) iec.DWORD { return TIME_TO_DWORD(t) }

// BLIND_ACTUATOR drives the motor of a blind, QU up and QD down, never both
// and with T_LOCKOUT between a change of direction, and simulates its
// position POS, which takes T_UD from bottom to top, and slat angle ANG,
// which takes T_ANGLE. STATUS is 121 up, 122 down, 1 for both, or S_IN.
type BLIND_ACTUATOR struct {
	UP, DN    iec.BOOL
	S_IN      iec.BYTE
	T_UD      iec.TIME // default T#10s
	T_ANGLE   iec.TIME // default T#3s
	T_LOCKOUT iec.TIME // default T#100ms
	POS, ANG  iec.BYTE
	QU, QD    iec.BOOL
	STATUS    iec.BYTE

	position, angle engineering.RMP_B
	lock            engineering.INTERLOCK
	initialized     bool
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BLIND_ACTUATOR) INIT() {
	*b = BLIND_ACTUATOR{
		T_UD: iec.TIME(10 * time.Second), T_ANGLE: iec.TIME(3 * time.Second),
		T_LOCKOUT: iec.TIME(100 * time.Millisecond), initialized: true,
	}
	b.position.INIT()
	b.angle.INIT()
}

// Execute runs the block once.
func (b *BLIND_ACTUATOR) Execute(now time.Time) {
	if !b.initialized {
		b.initialized = true
		b.position.INIT()
		b.angle.INIT()
	}
	// One motor at a time.
	b.lock.I1, b.lock.I2, b.lock.TL = b.UP, b.DN, b.T_LOCKOUT
	b.lock.Execute(now)

	// Simulate the slat angle, then the position.
	b.angle.E, b.angle.UP, b.angle.PT = b.lock.Q1 || b.lock.Q2, b.lock.Q1, b.T_ANGLE
	b.angle.Execute(now)
	b.position.E = b.lock.Q1 && b.angle.HIGH || b.lock.Q2 && b.angle.LOW
	b.position.UP, b.position.PT = b.lock.Q1, b.T_UD
	b.position.Execute(now)

	b.POS = b.position.OUT
	b.ANG = b.angle.OUT
	b.QU = b.lock.Q1
	b.QD = b.lock.Q2

	switch {
	case bool(b.UP) && bool(b.DN):
		b.STATUS = 1 // up and down together are an error
	case bool(b.UP):
		b.STATUS = 121
	case bool(b.DN):
		b.STATUS = 122
	default:
		b.STATUS = b.S_IN
	}
}

// BLIND_CONTROL controls a blind with a BLIND_ACTUATOR: UP or DN alone move
// it, and both together move it to the position PI and then the angle AI,
// within SENS. POS and ANG are the simulated position and angle, MU and MD
// the motor outputs.
type BLIND_CONTROL struct {
	UP, DN        iec.BOOL
	S_IN          iec.BYTE
	PI            iec.BYTE
	AI            iec.BYTE
	T_UD, T_ANGLE iec.TIME
	SENS          iec.BYTE // default 5
	T_LOCKOUT     iec.TIME // default T#100ms
	POS, ANG      iec.BYTE
	MU, MD        iec.BOOL
	STATUS        iec.BYTE

	act         BLIND_ACTUATOR
	delta       iec.BYTE
	bTimeTest   iec.BOOL
	iPos        iec.BYTE
	iAngel      iec.BYTE
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BLIND_CONTROL) INIT() {
	*b = BLIND_CONTROL{SENS: 5, T_LOCKOUT: iec.TIME(100 * time.Millisecond), initialized: true}
	b.act.INIT()
}

// runAct runs the actuator with the block's times.
func (b *BLIND_CONTROL) runAct(now time.Time) {
	b.act.T_UD, b.act.T_ANGLE, b.act.T_LOCKOUT = b.T_UD, b.T_ANGLE, b.T_LOCKOUT
	b.act.Execute(now)
}

// Execute runs the block once.
func (b *BLIND_CONTROL) Execute(now time.Time) {
	if !b.initialized {
		b.initialized = true
		b.act.INIT()
	}
	// A test of the run times, off in OSCAT.
	if b.bTimeTest {
		b.PI = b.iPos
		b.AI = b.iAngel
		b.UP = true
		b.DN = true
	}

	// The position.
	b.runAct(now)

	if b.UP && b.DN {
		// automatic mode: the position first, then the angle
		pos, ang := iec.INT(b.act.POS), iec.INT(b.act.ANG)
		pi, ai, delta := iec.INT(b.PI), iec.INT(b.AI), iec.INT(b.delta)
		angled := b.T_ANGLE > iec.TIME(100*time.Millisecond)
		switch {
		case pos < pi-delta:
			b.act.UP, b.act.DN = true, false
			b.delta = 0
			b.STATUS = 121
		case pos > pi+delta:
			b.act.UP, b.act.DN = false, true
			b.delta = 0
			b.STATUS = 122
		case ang < ai-delta && angled:
			b.act.UP, b.act.DN = true, false
			b.delta = b.SENS / 2
			b.STATUS = 123
		case ang > ai+delta && angled:
			b.act.UP, b.act.DN = false, true
			b.delta = b.SENS / 2
			b.STATUS = 124
		default:
			// the position is reached
			b.act.UP, b.act.DN = false, false
			b.delta = b.SENS
			b.STATUS = b.S_IN
		}
	} else {
		b.act.UP = b.UP
		b.act.DN = b.DN
		b.STATUS = b.S_IN
	}

	b.runAct(now)
	b.POS = b.act.POS
	b.ANG = b.act.ANG
	b.MU = b.act.QU
	b.MD = b.act.QD
	b.STATUS = b.act.STATUS
}

// BLIND_CONTROL_S controls a roller shutter, without slats: UP or DN alone
// move it, both together move it to the position PI, simulated from the
// times T_UP and T_DN, with T_LOCKOUT between changes of direction. At the
// top or bottom (within EXT_TRIG) it runs on for T_EXT to calibrate, and it
// calibrates at power up. RU and RD revert it from the bottom to R_POS_BOT
// or the top to R_POS_TOP. STATUS is 121 up, 122 down, 123 positioning, 124
// reverting, 127 lockout, 128 calibrating, 129 extending, or S_IN.
type BLIND_CONTROL_S struct {
	UP, DN    iec.BOOL
	S_IN      iec.BYTE // default 125
	PI        iec.BYTE
	T_UP      iec.TIME
	T_DN      iec.TIME
	RU        iec.BOOL
	RD        iec.BOOL
	T_LOCKOUT iec.TIME // default T#100ms
	T_EXT     iec.TIME
	EXT_TRIG  iec.BYTE // default 5
	R_POS_TOP iec.BYTE // default 255
	R_POS_BOT iec.BYTE
	POS       iec.BYTE
	MU, MD    iec.BOOL
	STATUS    iec.BYTE

	rmp         engineering.RMP_NEXT_
	last        iec.DWORD
	piLast      iec.BYTE
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BLIND_CONTROL_S) INIT() {
	*b = BLIND_CONTROL_S{
		S_IN: 125, T_LOCKOUT: iec.TIME(100 * time.Millisecond), EXT_TRIG: 5, R_POS_TOP: 255,
		initialized: true,
	}
	b.rmp.INIT()
}

// Execute runs the block once.
func (b *BLIND_CONTROL_S) Execute(now time.Time) {
	if !b.initialized {
		b.initialized = true
		b.rmp.INIT()
	}
	tx := PLC_MS(now)

	// The inputs.
	switch {
	case bool(b.UP) && !bool(b.DN):
		// manual up
		b.rmp.IN = 255
		b.STATUS = 121
	case bool(b.DN) && !bool(b.UP):
		// manual down
		b.rmp.IN = 0
		b.STATUS = 122
	case bool(!(b.UP || b.DN)):
		// manual standby
		b.rmp.IN = b.PI
		b.STATUS = b.S_IN
	}

	// Simulate the position.
	b.rmp.E, b.rmp.TR, b.rmp.TF, b.rmp.TL = b.UP || b.DN, b.T_UP, b.T_DN, b.T_LOCKOUT
	b.rmp.OUT = &b.POS
	b.rmp.Execute(now)

	switch b.STATUS {
	case 0: // power up
		b.last = tx
		b.piLast = b.PI ^ 255
		b.STATUS = 128 // calibrate

	case 121: // manual up
		b.MU, b.MD = true, false
		if b.POS >= 255-b.EXT_TRIG {
			// extend at the top
			b.POS = 255
			b.last = tx
			b.STATUS = 129
		}

	case 122: // manual down
		b.MD, b.MU = true, false
		if b.POS <= b.EXT_TRIG {
			// extend at the bottom
			b.POS = 0
			b.last = tx
			b.STATUS = 129
		}

	case 123: // automatic positioning
		b.MD, b.MU = b.rmp.DN, b.rmp.UP
		if !(b.rmp.DN || b.rmp.UP) {
			// the position is reached
			switch {
			case b.POS <= b.EXT_TRIG:
				b.MD = true
				b.last = tx
				b.STATUS = 129
			case b.POS >= 255-b.EXT_TRIG:
				b.MU = true
				b.last = tx
				b.STATUS = 129
			default:
				b.STATUS = b.S_IN
			}
		}

	case 124: // revert from the top or bottom
		b.MD, b.MU = b.rmp.DN, b.rmp.UP
		if !(b.rmp.DN || b.rmp.UP) {
			b.piLast = b.PI
			b.STATUS = b.S_IN
		}

	case 127: // lockout
		if tx-b.last >= ms(b.T_LOCKOUT) {
			b.STATUS = b.S_IN
		}

	case 128: // calibration
		b.MU, b.MD = true, false
		b.rmp.IN = 255
		if tx-b.last >= ms(b.T_UP)+ms(b.T_EXT) {
			b.MU = false
			b.last = tx
			b.STATUS = 127
		}

	case 129: // extend
		if tx-b.last >= ms(b.T_EXT) {
			b.MU, b.MD = false, false
			b.last = tx
			b.STATUS = 127
		}

	default:
		b.MU, b.MD = false, false
		switch {
		case b.PI != b.piLast:
			b.piLast = b.PI
			b.rmp.IN = b.PI
			b.STATUS = 123 // automatic positioning
		case b.POS == 0 && bool(b.RU):
			b.rmp.IN = b.R_POS_BOT
			b.STATUS = 124 // revert
		case b.POS == 255 && bool(b.RD):
			b.rmp.IN = b.R_POS_TOP
			b.STATUS = 124 // revert
		default:
			b.STATUS = b.S_IN
		}
	}
}

// BLIND_INPUT reads the switches of a blind, S1 up and S2 down (or S1 alone
// if SINGLE_SWITCH), debounced for DEBOUNCE_TIME: holding a switch moves
// the blind while held, a click (if CLICK_EN) moves it for MAX_RUNTIME or
// stops it, and a double click moves it to DBL_POS1/DBL_ANG1 or
// DBL_POS2/DBL_ANG2 if DBL_CLK1 or DBL_CLK2, or toggles D1 or D2. Manual
// operation ends after MANUAL_TIMEOUT; holding both switches stops the
// automatic operation until they are released. IN forces the position PI
// and angle AI. In automatic operation QU and QD are both true and PO and
// AO follow POS and ANG, unless MASTER_MODE.
type BLIND_INPUT struct {
	POS, ANG       iec.BYTE
	S1, S2         iec.BOOL
	IN             iec.BOOL
	PI, AI         iec.BYTE
	SINGLE_SWITCH  iec.BOOL
	CLICK_EN       iec.BOOL // default TRUE
	CLICK_TIME     iec.TIME // default T#500ms
	MAX_RUNTIME    iec.TIME // default T#60s
	MANUAL_TIMEOUT iec.TIME // default T#1h
	DEBOUNCE_TIME  iec.TIME // default T#20ms
	DBL_CLK1       iec.BOOL
	DBL_POS1       iec.BYTE
	DBL_ANG1       iec.BYTE
	DBL_CLK2       iec.BOOL
	DBL_POS2       iec.BYTE // default 255
	DBL_ANG2       iec.BYTE // default 255
	D1_TOGGLE      iec.BOOL // default TRUE
	D2_TOGGLE      iec.BOOL // default TRUE
	MASTER_MODE    iec.BOOL
	QU             iec.BOOL // default TRUE
	QD             iec.BOOL // default TRUE
	STATUS         iec.BYTE
	PO             iec.BYTE // default 255
	AO             iec.BYTE // default 255
	D1, D2         iec.BOOL

	s1e, s2e    timers.TOF
	s1d, s2d    electrical.CLICK_MODE
	dir         iec.BOOL
	last        iec.DWORD
	initialized bool
}

// INIT resets the block and sets its inputs and outputs to their initial
// values.
func (b *BLIND_INPUT) INIT() {
	*b = BLIND_INPUT{
		CLICK_EN: true, CLICK_TIME: iec.TIME(500 * time.Millisecond), MAX_RUNTIME: iec.TIME(60 * time.Second),
		MANUAL_TIMEOUT: iec.TIME(time.Hour), DEBOUNCE_TIME: iec.TIME(20 * time.Millisecond),
		DBL_POS2: 255, DBL_ANG2: 255, D1_TOGGLE: true, D2_TOGGLE: true,
		QU: true, QD: true, PO: 255, AO: 255,
		initialized: true,
	}
	b.s1d.INIT()
	b.s2d.INIT()
}

// Execute runs the block once.
func (b *BLIND_INPUT) Execute(now time.Time) {
	if !b.initialized {
		b.initialized = true
		b.QU, b.QD, b.PO, b.AO = true, true, 255, 255
		b.s1d.INIT()
		b.s2d.INIT()
	}
	tx := PLC_MS(now)

	// Debounce S1 and S2 and decode their clicks.
	b.s1e.IN, b.s1e.PT = b.S1, b.DEBOUNCE_TIME
	b.s1e.Execute(now)
	b.s2e.IN, b.s2e.PT = b.S2, b.DEBOUNCE_TIME
	b.s2e.Execute(now)
	b.s1d.IN = b.s1e.Q && !b.SINGLE_SWITCH || b.s1e.Q && b.SINGLE_SWITCH && b.dir
	b.s1d.T_LONG = b.CLICK_TIME
	b.s1d.Execute(now)
	b.s2d.IN = b.s2e.Q && !b.SINGLE_SWITCH || b.s1e.Q && b.SINGLE_SWITCH && !b.dir
	b.s2d.T_LONG = b.CLICK_TIME
	b.s2d.Execute(now)

	// D1 and D2 are pulses unless they toggle.
	if !b.D1_TOGGLE {
		b.D1 = false
	}
	if !b.D2_TOGGLE {
		b.D2 = false
	}

	// The actions.
	running := b.QU != b.QD
	switch {
	case bool(b.s1d.LONG && b.s2d.LONG || b.STATUS == 139):
		b.STATUS = 139
		if !(b.s1d.LONG || b.s2d.LONG) {
			b.STATUS = 130
		}
	case bool(b.s1d.TP_LONG):
		b.STATUS = 132
	case bool(b.s2d.TP_LONG):
		b.STATUS = 133
	case bool(b.s1d.SINGLE_):
		if b.CLICK_EN {
			// a click stops a running blind, or starts it
			if running {
				b.STATUS = 131
			} else {
				b.STATUS = 134
				b.last = tx
				b.dir = !b.dir
			}
		}
	case bool(b.s2d.SINGLE_):
		if b.CLICK_EN {
			if running {
				b.STATUS = 131
			} else {
				b.STATUS = 135
				b.last = tx
				b.dir = !b.dir
			}
		}
	case bool(b.IN):
		b.STATUS = 136
		b.last = tx
	case bool(b.s1d.DOUBLE):
		if b.DBL_CLK1 {
			b.STATUS = 137
			b.last = tx
		} else {
			b.D1 = !b.D1
		}
	case bool(b.s2d.DOUBLE):
		if b.DBL_CLK2 {
			b.STATUS = 138
			b.last = tx
		} else if b.SINGLE_SWITCH {
			// a single switch toggles D1
			b.D1 = !b.D1
		} else {
			b.D2 = !b.D2
		}
	}

	// The state machine.
	switch b.STATUS {
	case 0: // power up
		b.STATUS = 130
	case 130: // automatic operation
		if !b.MASTER_MODE {
			b.PO, b.AO = b.POS, b.ANG
		}
		b.QU, b.QD = true, true
	case 131: // manual standby, until MANUAL_TIMEOUT
		b.QU, b.QD = false, false
		b.PO, b.AO = b.POS, b.ANG
		if tx-b.last >= ms(b.MANUAL_TIMEOUT) {
			b.STATUS = 130
		}
	case 132: // manual up
		b.QU, b.QD = true, false
		b.PO, b.AO = b.POS, b.ANG
		b.last = tx
		if !b.s1d.LONG {
			b.STATUS = 131
			b.dir = !b.dir
		}
	case 133: // manual down
		b.QU, b.QD = false, true
		b.PO, b.AO = b.POS, b.ANG
		b.last = tx
		if !b.s2d.LONG {
			b.STATUS = 131
			b.dir = !b.dir
		}
	case 134: // a click up
		b.QU, b.QD = true, false
		b.PO, b.AO = b.POS, b.ANG
		if tx-b.last >= ms(b.MAX_RUNTIME) {
			b.STATUS = 131
		}
	case 135: // a click down
		b.QU, b.QD = false, true
		b.PO, b.AO = b.POS, b.ANG
		if tx-b.last >= ms(b.MAX_RUNTIME) {
			b.STATUS = 131
		}
	case 136: // forced by IN
		b.QU, b.QD = true, true
		b.PO, b.AO = b.PI, b.AI
		if tx-b.last >= ms(b.MAX_RUNTIME) {
			b.STATUS = 130
		}
	case 137: // the position of a double click of S1
		b.QU, b.QD = true, true
		b.PO, b.AO = b.DBL_POS1, b.DBL_ANG1
		if tx-b.last >= ms(b.MAX_RUNTIME) {
			b.STATUS = 131
		}
	case 138: // the position of a double click of S2
		b.QU, b.QD = true, true
		b.PO, b.AO = b.DBL_POS2, b.DBL_ANG2
		if tx-b.last >= ms(b.MAX_RUNTIME) {
			b.STATUS = 131
		}
	case 139: // both switches held: standby
		b.QU, b.QD = false, false
		b.PO, b.AO = b.POS, b.ANG
	}
}

// BLIND_NIGHT closes a blind at night, to NIGHT_POSITION and NIGHT_ANGLE,
// from SUNSET + SUNSET_OFFSET (if E_NIGHT) to SUNRISE + SUNRISE_OFFSET (if
// E_DAY), at the date and time DTIN, in automatic mode. Manual operation
// at night ends the night. STATUS is 141 at night.
type BLIND_NIGHT struct {
	UP, DN          iec.BOOL
	S_IN            iec.BYTE
	PI, AI          iec.BYTE
	E_NIGHT         iec.BOOL // default TRUE
	E_DAY           iec.BOOL // default TRUE
	DTIN            iec.DT
	SUNRISE, SUNSET iec.TOD
	SUNRISE_OFFSET  iec.TIME
	SUNSET_OFFSET   iec.TIME
	NIGHT_POSITION  iec.BYTE
	NIGHT_ANGLE     iec.BYTE
	QU, QD          iec.BOOL
	STATUS          iec.BYTE
	PO, AO          iec.BYTE

	night     iec.BOOL
	lastNight iec.DWORD
	lastDay   iec.DWORD
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BLIND_NIGHT) INIT() { *b = BLIND_NIGHT{E_NIGHT: true, E_DAY: true} }

// Execute runs the block once.
func (b *BLIND_NIGHT) Execute(now time.Time) {
	tod := TOD_TO_DWORD(DT_TO_TOD(b.DTIN))
	date := DATE_TO_DWORD(DT_TO_DATE(b.DTIN))
	switch {
	case !bool(b.UP && b.DN) && bool(b.night):
		// manual operation at night ends the night
		b.night = false
	case tod > TOD_TO_DWORD(b.SUNSET)+ms(b.SUNSET_OFFSET) && b.lastNight < date && !bool(b.night) && bool(b.E_NIGHT):
		// the night begins
		b.night = true
		b.lastNight = date
	case tod > TOD_TO_DWORD(b.SUNRISE)+ms(b.SUNRISE_OFFSET) && b.lastDay < date && bool(b.night) && bool(b.E_DAY) && b.lastNight < date:
		// the night ends
		b.night = false
		b.lastDay = date
	}

	// Close the blind at night, in automatic mode.
	if b.UP && b.DN && b.night {
		b.STATUS = 141
		b.PO, b.AO = b.NIGHT_POSITION, b.NIGHT_ANGLE
	} else {
		b.QU, b.QD = b.UP, b.DN
		b.PO, b.AO = b.PI, b.AI
		b.STATUS = b.S_IN
	}
}

// BLIND_SCENE stores 16 scenes of a blind's position and angle: SWRITE
// stores PI and AI in the scene SCENE (0..15), enabled if ENABLE, and in
// automatic mode an enabled scene sets PO and AO, with STATUS 160 + SCENE.
type BLIND_SCENE struct {
	UP, DN iec.BOOL
	S_IN   iec.BYTE
	PI, AI iec.BYTE
	ENABLE iec.BOOL
	SWRITE iec.BOOL
	SCENE  iec.BYTE
	QU, QD iec.BOOL
	STATUS iec.BYTE
	PO, AO iec.BYTE

	// sx is RETAIN in OSCAT: the position, angle and enable of each scene.
	sx [16][3]iec.BYTE
}

// INIT resets the block.
func (b *BLIND_SCENE) INIT() { *b = BLIND_SCENE{} }

// Execute runs the block once.
func (b *BLIND_SCENE) Execute(now time.Time) {
	// The lower 4 bits of SCENE.
	x := b.SCENE & 0x0F
	if b.ENABLE && b.sx[x][2] > 0 && b.UP && b.DN {
		b.PO, b.AO = b.sx[x][0], b.sx[x][1]
		b.STATUS = 160 + x // 160 to 175 for the 16 scenes
		b.QU, b.QD = true, true
	} else {
		b.QU, b.QD = b.UP, b.DN
		b.STATUS = b.S_IN
		b.PO, b.AO = b.PI, b.AI
	}

	// Store a scene.
	if b.SWRITE {
		b.STATUS = 176
		b.sx[x][0], b.sx[x][1] = b.PI, b.AI
		b.sx[x][2] = SEL[iec.BYTE](b.ENABLE, 0, 1)
	}
}

// BLIND_SECURITY moves a blind for safety: up on FIRE, by WIND (up if
// WIND_UP, else down), on ALARM (up if ALARM_UP), on an open DOOR, and by
// RAIN (up if RAIN_UP) unless moved by hand, with STATUS 111..115.
type BLIND_SECURITY struct {
	UP, DN   iec.BOOL
	S_IN     iec.BYTE
	PI, AI   iec.BYTE
	FIRE     iec.BOOL
	WIND     iec.BOOL
	ALARM    iec.BOOL
	DOOR     iec.BOOL
	RAIN     iec.BOOL
	ALARM_UP iec.BOOL // default TRUE
	WIND_UP  iec.BOOL // default TRUE
	RAIN_UP  iec.BOOL
	QU, QD   iec.BOOL
	STATUS   iec.BYTE
	PO, AO   iec.BYTE
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BLIND_SECURITY) INIT() { *b = BLIND_SECURITY{ALARM_UP: true, WIND_UP: true} }

// Execute runs the block once.
func (b *BLIND_SECURITY) Execute(now time.Time) {
	switch {
	case bool(b.FIRE):
		b.QU, b.QD, b.STATUS = true, false, 111
	case bool(b.WIND):
		b.QU, b.QD, b.STATUS = b.WIND_UP, !b.WIND_UP, 112
	case bool(b.ALARM):
		b.QU, b.QD, b.STATUS = b.ALARM_UP, !b.ALARM_UP, 113
	case bool(b.DOOR):
		b.QU, b.QD, b.STATUS = true, false, 114
	case bool(b.RAIN) && b.UP == b.DN:
		b.QU, b.QD, b.STATUS = b.RAIN_UP, !b.RAIN_UP, 115
	default:
		b.QU, b.QD = b.UP, b.DN
		b.STATUS = b.S_IN
		b.PO, b.AO = b.PI, b.AI
	}
}

// BLIND_SET forces a blind to the position PX and angle AX while IN is
// true, in automatic mode or if OVERRIDE_MANUAL, with STATUS 178; then, if
// RESTORE_POSITION, it returns to the position before (STATUS 179), until
// PI and AI reach it or RESTORE_TIME passes.
type BLIND_SET struct {
	UP, DN           iec.BOOL
	S_IN             iec.BYTE
	PI, AI           iec.BYTE
	IN               iec.BOOL
	PX, AX           iec.BYTE
	OVERRIDE_MANUAL  iec.BOOL
	RESTORE_POSITION iec.BOOL
	RESTORE_TIME     iec.TIME // default T#60s
	QU, QD           iec.BOOL
	STATUS           iec.BYTE
	PO, AO           iec.BYTE

	ps, as iec.BYTE
	last   iec.DWORD
}

// INIT resets the block and sets RESTORE_TIME to its initial value.
func (b *BLIND_SET) INIT() { *b = BLIND_SET{RESTORE_TIME: iec.TIME(60 * time.Second)} }

// Execute runs the block once.
func (b *BLIND_SET) Execute(now time.Time) {
	tx := PLC_MS(now)
	if b.IN && (b.OVERRIDE_MANUAL || b.UP && b.DN) {
		b.STATUS = 178
	}

	switch b.STATUS {
	case 0: // power on
		b.STATUS = b.S_IN
	case 178: // forced to PX and AX
		b.PO, b.AO = b.PX, b.AX
		b.QU, b.QD = true, true
		if !b.IN {
			b.STATUS = SEL(b.RESTORE_POSITION, b.S_IN, 179)
			b.last = tx
		}
	case 179: // restore the position before
		b.PO, b.AO = b.ps, b.as
		if b.PO == b.PI && b.AO == b.AI || tx-b.last >= ms(b.RESTORE_TIME) {
			b.STATUS = b.S_IN
		}
	default:
		// pass the inputs on
		b.PO, b.ps = b.PI, b.PI
		b.AO, b.as = b.AI, b.AI
		b.STATUS = b.S_IN
		b.QU, b.QD = b.UP, b.DN
	}
}

// BLIND_SHADE shades a blind from the sun, in automatic mode while ENABLE
// and SUN (held on for SHADE_DELAY): when the sun, by the calendar CX,
// stands within ANGLE_OFFSET of the DIRECTION the window faces, between
// SUNRISE_OFFSET after sunrise and SUNSET_PRESET before sunset, it moves
// the blind to SHADE_POS and sets the slat angle, from SLAT_WIDTH and
// SLAT_SPACING, to keep the sun out. STATUS is 151 while shading.
type BLIND_SHADE struct {
	UP, DN         iec.BOOL
	S_IN           iec.BYTE
	PI, AI         iec.BYTE
	ENABLE         iec.BOOL
	SUN            iec.BOOL
	CX             *CALENDAR
	SUNRISE_OFFSET iec.TIME // default T#1h
	SUNSET_PRESET  iec.TIME // default T#1h
	DIRECTION      iec.REAL // default 180.0
	ANGLE_OFFSET   iec.REAL // default 80.0
	SLAT_WIDTH     iec.REAL // default 80.0
	SLAT_SPACING   iec.REAL // default 60.0
	SHADE_DELAY    iec.TIME // default T#60s
	SHADE_POS      iec.BYTE
	QU, QD         iec.BOOL
	STATUS         iec.BYTE
	PO, AO         iec.BYTE

	angle    iec.REAL
	sunDelay timers.TOF
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BLIND_SHADE) INIT() {
	*b = BLIND_SHADE{
		CX: b.CX, SUNRISE_OFFSET: iec.TIME(time.Hour), SUNSET_PRESET: iec.TIME(time.Hour),
		DIRECTION: 180, ANGLE_OFFSET: 80, SLAT_WIDTH: 80, SLAT_SPACING: 60,
		SHADE_DELAY: iec.TIME(60 * time.Second),
	}
}

// Execute runs the block once.
func (b *BLIND_SHADE) Execute(now time.Time) {
	if b.CX == nil {
		return
	}
	cx := b.CX
	// SUN, held on for SHADE_DELAY.
	b.sunDelay.IN, b.sunDelay.PT = b.SUN, b.SHADE_DELAY
	b.sunDelay.Execute(now)

	utc := TOD_TO_DWORD(DT_TO_TOD(cx.UTC))
	if b.UP && b.DN && b.ENABLE && b.sunDelay.Q &&
		cx.SUN_HOR > b.DIRECTION-b.ANGLE_OFFSET && cx.SUN_HOR < b.DIRECTION+b.ANGLE_OFFSET &&
		utc > TOD_TO_DWORD(cx.SUN_RISE)+ms(b.SUNRISE_OFFSET) && utc < TOD_TO_DWORD(cx.SUN_SET)-ms(b.SUNSET_PRESET) {
		b.STATUS = 151
		b.QU, b.QD = b.UP, b.DN
		b.PO = b.SHADE_POS
		// The steepest angle of the slats.
		b.angle = math.DEG(ATAN(b.SLAT_SPACING / b.SLAT_WIDTH))
		if cx.SUN_VER > 0.0 && cx.SUN_VER < b.angle {
			b.angle = cx.SUN_VER + math.DEG(ACOS(COS(math.RAD(cx.SUN_VER))*b.SLAT_SPACING/b.SLAT_WIDTH))
			b.AO = iec.BYTE(LIMIT(0, TRUNC(b.angle*2.833333333), 255))
		} else {
			b.AO = 255
		}
	} else {
		b.QU, b.QD = b.UP, b.DN
		b.PO, b.AO = b.PI, b.AI
		b.STATUS = b.S_IN
	}
}

// BLIND_SHADE_S shades a roller shutter from the sun, in automatic mode
// while ENABLE and SUN (held on for SHADE_DELAY): when the sun, by the
// calendar CX, stands between HORZ1 and HORZ2 and below VERT, between
// SUNRISE_OFFSET after sunrise and SUNSET_PRESET before sunset, it moves
// the shutter down to SHADE_POS at least, with STATUS 151. ALERT moves it
// up, with STATUS 152.
type BLIND_SHADE_S struct {
	UP, DN         iec.BOOL
	S_IN           iec.BYTE
	PI             iec.BYTE
	ENABLE         iec.BOOL
	SUN            iec.BOOL
	HORZ1          iec.REAL // default 100.0
	HORZ2          iec.REAL // default 260.0
	VERT           iec.REAL // default 90.0
	ALERT          iec.BOOL
	CX             *CALENDAR
	SUNRISE_OFFSET iec.TIME // default T#1h
	SUNSET_PRESET  iec.TIME // default T#1h
	SHADE_DELAY    iec.TIME // default T#60s
	SHADE_POS      iec.BYTE
	QU, QD         iec.BOOL
	STATUS         iec.BYTE
	PO             iec.BYTE

	sunDelay timers.TOF
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BLIND_SHADE_S) INIT() {
	*b = BLIND_SHADE_S{
		CX: b.CX, HORZ1: 100, HORZ2: 260, VERT: 90,
		SUNRISE_OFFSET: iec.TIME(time.Hour), SUNSET_PRESET: iec.TIME(time.Hour),
		SHADE_DELAY: iec.TIME(60 * time.Second),
	}
}

// Execute runs the block once.
func (b *BLIND_SHADE_S) Execute(now time.Time) {
	if b.CX == nil {
		return
	}
	cx := b.CX
	b.sunDelay.IN, b.sunDelay.PT = b.SUN, b.SHADE_DELAY
	b.sunDelay.Execute(now)

	utc := TOD_TO_DWORD(DT_TO_TOD(cx.UTC))
	switch {
	case bool(b.ALERT):
		b.QU, b.QD = true, false
		b.STATUS = 152
	case bool(b.UP && b.DN && b.ENABLE && b.sunDelay.Q) &&
		cx.SUN_HOR > b.HORZ1 && cx.SUN_HOR < b.HORZ2 && cx.SUN_VER < b.VERT &&
		utc > TOD_TO_DWORD(cx.SUN_RISE)+ms(b.SUNRISE_OFFSET) && utc < TOD_TO_DWORD(cx.SUN_SET)-ms(b.SUNSET_PRESET):
		b.QU, b.QD = b.UP, b.DN
		b.STATUS = 151
		b.PO = min(b.PI, b.SHADE_POS)
	default:
		b.QU, b.QD = b.UP, b.DN
		b.PO = b.PI
		b.STATUS = b.S_IN
	}
}
