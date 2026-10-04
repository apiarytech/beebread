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

// Package actuators is the port of the OSCAT BUILDING actuators: valves,
// coils, pumps and up/down drives, with self activation to keep them from
// seizing.
package actuators

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/engineering"
	"github.com/apiarytech/beebread/basic/logic"
	"github.com/apiarytech/royaljelly/iec"
)

// ms returns a TIME in milliseconds, as the PLC timer counts.
func ms(t iec.TIME) iec.DWORD { return TIME_TO_DWORD(t) }

// ACTUATOR_2P drives a valve with a two point output OUT, pulse width
// modulated by IN (0..255) over CYCLE_TIME. IN below SENS is off and above
// 255 - SENS is on. An AUTORUN moves the valve every SELF_ACT_TIME, or when
// TEST is true, for SELF_ACT_CYCLES pulses of SELF_ACT_PULSE; ARX, shared
// by the actuators, lets one self activate at a time.
type ACTUATOR_2P struct {
	IN              iec.BYTE
	TEST            iec.BOOL
	ARE             iec.BOOL // default TRUE
	ARX             *iec.BOOL
	CYCLE_TIME      iec.TIME
	SENS            iec.BYTE
	SELF_ACT_TIME   iec.TIME
	SELF_ACT_PULSE  iec.TIME
	SELF_ACT_CYCLES iec.INT // default 1
	OUT             iec.BOOL
	ARO             iec.BOOL

	timer       AUTORUN
	pwgen       engineering.GEN_PULSE
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (a *ACTUATOR_2P) INIT() {
	*a = ACTUATOR_2P{ARX: a.ARX, ARE: true, SELF_ACT_CYCLES: 1, initialized: true}
	a.timer.INIT()
	a.pwgen.INIT()
}

// Execute runs the block once.
func (a *ACTUATOR_2P) Execute(now time.Time) {
	if !a.initialized {
		a.initialized = true
		a.timer.INIT()
		a.pwgen.INIT()
	}
	if a.ARX == nil {
		return
	}
	a.timer.TRUN = iec.TIME(time.Duration(a.SELF_ACT_PULSE) * time.Duration(2*a.SELF_ACT_CYCLES))
	a.timer.TOFF = a.SELF_ACT_TIME
	a.timer.TEST = a.TEST
	a.timer.ARE = a.ARE
	a.timer.ARX = a.ARX
	a.timer.Execute(now)
	a.ARO = a.timer.ARO

	switch {
	case bool(a.ARO):
		a.pwgen.PTL, a.pwgen.PTH = a.SELF_ACT_PULSE, a.SELF_ACT_PULSE
		a.pwgen.Execute(now)
		a.OUT = a.pwgen.Q
	case a.IN < a.SENS:
		a.OUT = false
	case a.IN > 255-a.SENS:
		a.OUT = true
	default:
		a.pwgen.PTH = DWORD_TO_TIME(ms(a.CYCLE_TIME) * iec.DWORD(engineering.BAND_B(a.IN, a.SENS)) / 255)
		a.pwgen.PTL = a.CYCLE_TIME - a.pwgen.PTH
		a.pwgen.Execute(now)
		a.OUT = a.pwgen.Q
	}
}

// ACTUATOR_3P drives a valve with a three point output, OUT1 to open and
// OUT2 to close, to the position IN (0..255), simulating its position POS
// from the runtime T_RUN. It calibrates every T_CAL by running to an end,
// for T_EXT longer than T_RUN, and runs a diagnosis every T_DIAG, or when
// TEST is true, that sets ERROR if the end switch END_POS (with
// SWITCH_AVAIL) is not reached or the runtimes up and down differ by more
// than 10%. ARX, shared by the actuators, lets one self activate at a time.
type ACTUATOR_3P struct {
	IN           iec.BYTE
	TEST         iec.BOOL
	ARE          iec.BOOL // default TRUE
	END_POS      iec.BOOL
	T_RUN        iec.TIME // default T#60s
	T_EXT        iec.TIME // default T#10s
	T_CAL        iec.TIME // default T#600s
	T_DIAG       iec.TIME // default T#10d
	SWITCH_AVAIL iec.BOOL
	ARX          *iec.BOOL
	OUT1         iec.BOOL
	OUT2         iec.BOOL
	POS          iec.BYTE
	ERROR        iec.BOOL
	STATUS       iec.BYTE

	ramp        engineering.RMP_NEXT_
	nextCal     iec.DWORD
	nextDiag    iec.DWORD
	last        iec.DWORD
	start       iec.DWORD
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (a *ACTUATOR_3P) INIT() {
	*a = ACTUATOR_3P{
		ARX: a.ARX, ARE: true,
		T_RUN: iec.TIME(60 * time.Second), T_EXT: iec.TIME(10 * time.Second),
		T_CAL: iec.TIME(600 * time.Second), T_DIAG: iec.TIME(10 * 24 * time.Hour),
		initialized: true,
	}
	a.ramp.INIT()
}

// Execute runs the block once.
func (a *ACTUATOR_3P) Execute(now time.Time) {
	if !a.initialized {
		a.initialized = true
		a.ramp.INIT()
	}
	if a.ARX == nil {
		return
	}
	tx := PLC_MS(now)

	// The test input.
	if a.TEST {
		a.STATUS = 103
		a.start = tx
		*a.ARX = true
	}

	switch a.STATUS {
	case 0: // power on setup
		if a.ARE && !*a.ARX {
			a.STATUS = 103
			a.start = tx
			*a.ARX = true
		}

	case 100: // normal operation
		switch {
		case bool(a.T_DIAG > 0 && tx > a.nextDiag && a.ARE && !*a.ARX):
			// auto diagnostics
			a.STATUS = 103
			a.start = tx
			*a.ARX = true
		case bool(a.T_CAL > 0 && tx > a.nextCal && a.ARE && !*a.ARX):
			// auto calibration
			if a.POS > 127 {
				a.OUT1, a.OUT2 = true, false
				a.ramp.IN = 255
			} else {
				a.OUT1, a.OUT2 = false, true
				a.ramp.IN = 0
			}
			*a.ARX = true
			a.STATUS = 101
			a.start = tx
		default:
			// Calibration waits while the valve does not move.
			if !(a.OUT1 || a.OUT2) {
				a.nextCal += tx - a.last
			}
			a.ramp.IN = a.IN
		}

	case 101: // calibrate
		switch {
		case tx-a.start < ms(a.T_EXT):
			a.nextCal = tx + ms(a.T_CAL)
		case bool(a.SWITCH_AVAIL && a.END_POS):
			a.STATUS = 100
			*a.ARX = false
		case tx-a.start > ms(a.T_EXT)+ms(a.T_RUN):
			a.ERROR = a.SWITCH_AVAIL
			*a.ARX = false
		}

	case 103: // diagnostics up, for T_EXT
		switch {
		case tx-a.start < ms(a.T_EXT):
			a.ERROR = false
			a.ramp.TR = a.T_RUN
			a.ramp.TF = a.T_RUN
			a.OUT1, a.OUT2 = true, false
			a.ramp.IN = 255
		case bool(a.SWITCH_AVAIL && a.END_POS):
			a.ramp.TR = DWORD_TO_TIME(tx - a.start)
			a.STATUS = 104
		case tx-a.start > ms(a.T_EXT)+ms(a.T_RUN):
			a.ERROR = a.SWITCH_AVAIL
			a.STATUS = 104
			a.start = tx
		}

	case 104: // diagnostics down
		switch {
		case tx-a.start < ms(a.T_EXT):
			a.OUT1, a.OUT2 = false, true
			a.ramp.IN = 0
			a.nextDiag = tx + ms(a.T_DIAG)
		case bool(a.SWITCH_AVAIL && a.END_POS):
			a.ramp.TR = DWORD_TO_TIME(tx - a.start)
			// The runtimes up and down differ by more than 10%.
			diff := iec.DINT(ms(a.ramp.TR)) - iec.DINT(ms(a.ramp.TF))
			if diff < 0 {
				diff = -diff
			}
			if iec.DWORD(diff*10) > ms(a.T_RUN) {
				a.ERROR = true
			}
			a.STATUS = 100
			*a.ARX = false
			a.nextCal = tx + ms(a.T_CAL)
		case tx-a.start > ms(a.T_EXT)+ms(a.T_RUN):
			if a.SWITCH_AVAIL {
				a.ERROR = true
			}
			a.STATUS = 100
			*a.ARX = false
			a.nextCal = tx + ms(a.T_CAL)
		}
	}

	// The simulated position of the flap, and the outputs.
	a.ramp.OUT = &a.POS
	a.ramp.Execute(now)
	if a.STATUS == 100 {
		a.OUT1 = a.ramp.UP
		a.OUT2 = a.ramp.DN
	}

	// An end switch sets the position.
	if a.SWITCH_AVAIL && a.END_POS {
		a.POS = SEL[iec.BYTE](a.POS > 127, 0, 255)
		a.nextCal = tx + ms(a.T_CAL)
	}
	a.last = tx
}

// ACTUATOR_A drives an analog actuator: Y goes from OUT_MIN to OUT_MAX as
// I1 (or I2 if IS is true) goes from 0 to 255, reversed if RV is true. A
// rising edge of DX starts a self activation that drives Y to OUT_MIN and
// then OUT_MAX for RUNTIME each; it repeats every SELF_ACT_TIME if that is
// not 0.
type ACTUATOR_A struct {
	I1            iec.BYTE
	IS            iec.BOOL
	I2            iec.BYTE
	RV            iec.BOOL
	DX            iec.BOOL
	RUNTIME       iec.TIME
	SELF_ACT_TIME iec.TIME
	OUT_MIN       iec.WORD
	OUT_MAX       iec.WORD
	Y             iec.WORD

	timer       logic.CYCLE_4
	dxEdge      iec.BOOL
	initialized bool
}

// INIT resets the block.
func (a *ACTUATOR_A) INIT() {
	*a = ACTUATOR_A{initialized: true}
	a.timer.INIT()
}

// Execute runs the block once.
func (a *ACTUATOR_A) Execute(now time.Time) {
	if !a.initialized {
		a.initialized = true
		a.timer.INIT()
	}
	// The cycle stays in state 0 if SELF_ACT_TIME is 0.
	a.timer.T0, a.timer.T1, a.timer.T3 = a.RUNTIME, a.RUNTIME, a.SELF_ACT_TIME
	a.timer.SL = a.DX && !a.dxEdge
	a.timer.SX = 0
	a.timer.S0 = a.SELF_ACT_TIME > 0
	a.timer.Execute(now)
	a.dxEdge = a.DX

	in := int64(SEL(a.IS, a.I1, a.I2))
	lo, hi := int64(a.OUT_MIN), int64(a.OUT_MAX)
	switch a.timer.STATE {
	case 0: // self activation, minimum
		a.Y = a.OUT_MIN
	case 1: // self activation, maximum
		a.Y = a.OUT_MAX
	case 3: // normal operation
		if a.RV {
			a.Y = iec.WORD(hi - (hi-lo)*in/255)
		} else {
			a.Y = iec.WORD((hi-lo)*in/255 + lo)
		}
	}
}

// ACTUATOR_COIL drives a coil: OUT follows IN, and while IN is false OUT
// goes on for SELF_ACT_TIME every SELF_ACT_CYCLE. STATUS is 100 when idle,
// 101 when on by IN and 102 when self activated.
type ACTUATOR_COIL struct {
	IN             iec.BOOL
	SELF_ACT_CYCLE iec.TIME // default T#10d
	SELF_ACT_TIME  iec.TIME // default T#1s
	OUT            iec.BOOL
	STATUS         iec.BYTE

	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (a *ACTUATOR_COIL) INIT() {
	*a = ACTUATOR_COIL{SELF_ACT_CYCLE: iec.TIME(10 * 24 * time.Hour), SELF_ACT_TIME: iec.TIME(time.Second)}
}

// Execute runs the block once.
func (a *ACTUATOR_COIL) Execute(now time.Time) {
	tn := PLC_MS(now)
	switch {
	case bool(!a.init):
		a.last = tn
		a.init = true
	case bool(a.IN):
		a.OUT = true
		a.STATUS = 101 // activated by the input
		a.last = tn
	default:
		a.OUT = false
		a.STATUS = 100 // disabled
		// Self activation.
		tx := tn - a.last
		if a.SELF_ACT_CYCLE > 0 && tx >= ms(a.SELF_ACT_CYCLE) {
			a.OUT = true
			a.STATUS = 102 // self activated
			if tx >= ms(a.SELF_ACT_CYCLE)+ms(a.SELF_ACT_TIME) {
				a.last = tn
				a.OUT = false
				a.STATUS = 100 // idle
			}
		}
	}
}

// ACTUATOR_PUMP drives a pump: PUMP follows IN, or MANUAL, but stays on for
// MIN_ONTIME and off for MIN_OFFTIME, and goes on after RUN_EVERY off if
// that is not 0. It counts the pump's RUNTIME in seconds and its CYCLES,
// which RST clears.
type ACTUATOR_PUMP struct {
	IN          iec.BOOL
	MANUAL      iec.BOOL
	RST         iec.BOOL
	MIN_ONTIME  iec.TIME // default T#10s
	MIN_OFFTIME iec.TIME // default T#10s
	RUN_EVERY   iec.TIME // default T#10000m
	PUMP        iec.BOOL
	RUNTIME     *iec.UDINT
	CYCLES      *iec.UDINT

	lastChange iec.DWORD
	meter      engineering.ONTIME
	oldMan     iec.BOOL
	init       iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (a *ACTUATOR_PUMP) INIT() {
	*a = ACTUATOR_PUMP{
		RUNTIME: a.RUNTIME, CYCLES: a.CYCLES,
		MIN_ONTIME: iec.TIME(10 * time.Second), MIN_OFFTIME: iec.TIME(10 * time.Second),
		RUN_EVERY: iec.TIME(10000 * time.Minute),
	}
}

// Execute runs the block once.
func (a *ACTUATOR_PUMP) Execute(now time.Time) {
	if a.RUNTIME == nil || a.CYCLES == nil {
		return
	}
	tx := PLC_MS(now)
	switch {
	case bool(!a.init):
		a.init = true
		a.lastChange = tx
	case bool(a.RST):
		a.RST = false
		*a.RUNTIME = 0
		*a.CYCLES = 0
	case bool(a.MANUAL && !a.PUMP && !a.oldMan):
		a.lastChange = tx
		a.PUMP = true
	case bool(!a.MANUAL && a.oldMan && a.PUMP && !a.IN):
		a.lastChange = tx
		a.PUMP = false
	case bool(a.IN && !a.PUMP && tx-a.lastChange >= ms(a.MIN_OFFTIME)):
		a.lastChange = tx
		a.PUMP = true
	case bool(a.PUMP && !a.IN && !a.MANUAL && tx-a.lastChange >= ms(a.MIN_ONTIME)):
		a.lastChange = tx
		a.PUMP = false
	case bool(!a.PUMP && tx-a.lastChange >= ms(a.RUN_EVERY) && a.RUN_EVERY > 0):
		a.lastChange = tx
		a.PUMP = true
	}
	a.meter.IN = a.PUMP
	a.meter.SECONDS, a.meter.CYCLES = a.RUNTIME, a.CYCLES
	a.meter.Execute(now)
	a.oldMan = a.MANUAL
}

// ACTUATOR_UD drives a motor up (YUP) or down (YDN): automatically, by UD
// while ON is true, or by UP and DN while MANUAL is true. OFF turns it off.
// An output stays on for TON and off for TOFF at least, and on a change of
// direction both go off first. If OUT_RETURN is true, the outputs fed back
// in YUP_IN and YDN_IN keep both from being on.
type ACTUATOR_UD struct {
	UD         iec.BOOL
	ON         iec.BOOL
	MANUAL     iec.BOOL
	UP         iec.BOOL
	DN         iec.BOOL
	OFF        iec.BOOL
	YUP_IN     iec.BOOL
	YDN_IN     iec.BOOL
	TON        iec.TIME
	TOFF       iec.TIME
	OUT_RETURN iec.BOOL
	YUP        iec.BOOL
	YDN        iec.BOOL
	STATUS     iec.BYTE

	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block.
func (a *ACTUATOR_UD) INIT() { *a = ACTUATOR_UD{} }

// Execute runs the block once.
func (a *ACTUATOR_UD) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(!a.init):
		a.last = tx
		a.init = true
	case bool(a.OFF):
		// emergency shut off
		a.YUP, a.YDN = false, false
		a.last = tx
		a.STATUS = 101
	case bool((a.YUP || a.YDN) && tx-a.last < ms(a.TON)):
		// the minimum time on
		return
	case bool(!a.YUP && !a.YDN && tx-a.last < ms(a.TOFF)):
		// the minimum time off
		return
	case bool(a.MANUAL):
		a.STATUS = 102
		switch {
		case bool(a.YUP && !a.UP || a.YDN && !a.DN):
			// a change of direction turns both outputs off first
			a.YDN, a.YUP = false, false
			a.last = tx
		case bool(a.UP && !a.DN && !a.OFF):
			a.YDN, a.YUP = false, true
			a.last = tx
			a.STATUS = 103
		case bool(a.DN && !a.UP && !a.OFF):
			a.YUP, a.YDN = false, true
			a.last = tx
			a.STATUS = 104
		default:
			if a.YUP || a.YDN {
				a.last = tx
			}
			a.YUP, a.YDN = false, false
		}
	default:
		// automatic operation
		switch {
		case bool(a.YUP && !a.UD || a.YDN && a.UD):
			a.YUP, a.YDN = false, false
			a.last = tx
		case bool(a.UD && a.ON && !a.OFF):
			a.YDN, a.YUP = false, true
			a.last = tx
			a.STATUS = 111
		case bool(!a.UD && a.ON && !a.OFF):
			a.YUP, a.YDN = false, true
			a.last = tx
			a.STATUS = 112
		default:
			if a.YUP || a.YDN {
				a.last = tx
			}
			a.YUP, a.YDN = false, false
			a.STATUS = 110
		}
	}

	// Yup and Ydn are never on together.
	if a.YDN && a.YUP_IN && a.OUT_RETURN {
		a.YDN = false
		a.STATUS = 1
	}
	if a.YUP && a.YDN_IN && a.OUT_RETURN {
		a.YUP = false
		a.STATUS = 2
	}
}

// AUTORUN runs a self activation: when its timer, which counts TOFF while
// OUT is false and TRUN while it is true, runs out, or when TEST is true, ARO
// goes on until the timer runs out again, if ARE is true and no other
// actuator sharing ARX is self activating. OUT is IN or ARO.
type AUTORUN struct {
	IN   iec.BOOL
	TEST iec.BOOL
	ARE  iec.BOOL // default TRUE
	ARX  *iec.BOOL
	TRUN iec.TIME
	TOFF iec.TIME
	OUT  iec.BOOL
	ARO  iec.BOOL

	timer       engineering.RMP_B_
	val         iec.BYTE
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (a *AUTORUN) INIT() {
	*a = AUTORUN{ARX: a.ARX, ARE: true, initialized: true}
	a.timer.INIT()
}

// Execute runs the block once.
func (a *AUTORUN) Execute(now time.Time) {
	if !a.initialized {
		a.initialized = true
		a.timer.INIT()
	}
	if a.ARX == nil {
		return
	}
	// The timer runs while TOFF is not 0.
	if a.TOFF > 0 {
		a.timer.DIR = a.OUT
		a.timer.TR = SEL(a.OUT, a.TOFF, a.TRUN)
		a.timer.RMP = &a.val
		a.timer.Execute(now)
	} else {
		a.val = 255
	}

	// When the timer is 0, or on a test, self activate until it is 255.
	if !*a.ARX && a.ARE && a.val == 0 || a.TEST {
		a.val = 0
		a.ARO = true
		*a.ARX = true
	} else if a.val == 255 && a.ARO {
		a.ARO = false
		*a.ARX = false
	}
	a.OUT = a.IN || a.ARO
}
