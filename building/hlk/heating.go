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

package hlk

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/engineering"
	"github.com/apiarytech/beebread/basic/logic"
	"github.com/apiarytech/beebread/building/electrical"
	"github.com/apiarytech/royaljelly/iec"
)

// ms returns a TIME in milliseconds, as the PLC timer counts.
func ms(t iec.TIME) iec.DWORD { return TIME_TO_DWORD(t) }

// BOILER controls the heating of a hot water boiler from the temperatures
// T_UPPER and T_LOWER (if T_LOWER_ENABLE): while ENABLE it keeps T_UPPER
// between T_UPPER_MIN and T_UPPER_MAX, REQ_1 and REQ_2 request T_REQUEST_1
// and T_REQUEST_2 (plus T_REQUEST_HYS to stop), and a rising edge of BOOST
// heats it once to T_UPPER_MAX. Temperatures outside T_PROTECT_LOW and
// T_PROTECT_HIGH, or no PRESSURE, are errors: STATUS 1..5.
type BOILER struct {
	T_UPPER        iec.REAL
	T_LOWER        iec.REAL
	PRESSURE       iec.BOOL // default TRUE
	ENABLE         iec.BOOL // default TRUE
	REQ_1          iec.BOOL
	REQ_2          iec.BOOL
	BOOST          iec.BOOL
	T_UPPER_MIN    iec.REAL // default 50.0
	T_UPPER_MAX    iec.REAL // default 60.0
	T_LOWER_ENABLE iec.BOOL
	T_LOWER_MAX    iec.REAL // default 60.0
	T_REQUEST_1    iec.REAL // default 70.0
	T_REQUEST_2    iec.REAL // default 50.0
	T_REQUEST_HYS  iec.REAL // default 5.0
	T_PROTECT_HIGH iec.REAL // default 80.0
	T_PROTECT_LOW  iec.REAL // default 10.0
	HEAT           iec.BOOL
	ERROR          iec.BOOL
	STATUS         iec.BYTE

	edge      iec.BOOL
	boostMode iec.BOOL
	flag0     iec.BOOL
	flag1     iec.BOOL
	flag2     iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BOILER) INIT() {
	*b = BOILER{
		PRESSURE: true, ENABLE: true,
		T_UPPER_MIN: 50, T_UPPER_MAX: 60, T_LOWER_MAX: 60,
		T_REQUEST_1: 70, T_REQUEST_2: 50, T_REQUEST_HYS: 5,
		T_PROTECT_HIGH: 80, T_PROTECT_LOW: 10,
	}
}

// Execute runs the block once.
func (b *BOILER) Execute(now time.Time) {
	// Check the sensors.
	switch {
	case b.T_UPPER > b.T_PROTECT_HIGH:
		b.STATUS, b.HEAT, b.ERROR = 1, false, true
	case b.T_UPPER < b.T_PROTECT_LOW:
		b.STATUS, b.HEAT, b.ERROR = 2, true, true
	case b.T_LOWER > b.T_PROTECT_HIGH && bool(b.T_LOWER_ENABLE):
		b.STATUS, b.HEAT, b.ERROR = 3, false, true
	case b.T_LOWER < b.T_PROTECT_LOW && bool(b.T_LOWER_ENABLE):
		b.STATUS, b.HEAT, b.ERROR = 4, true, true
	case bool(!b.PRESSURE):
		b.STATUS, b.HEAT, b.ERROR = 5, false, true
	case bool(b.REQ_1 || b.REQ_2 || b.ENABLE || b.BOOST):
		b.ERROR = false

		// Whether to heat.
		switch {
		case bool(b.BOOST) && !bool(b.edge) && b.T_UPPER < b.T_UPPER_MAX:
			b.STATUS, b.HEAT = 101, true
			b.boostMode = true
		case bool(b.ENABLE) && b.T_UPPER < b.T_UPPER_MIN:
			b.STATUS, b.HEAT = 102, true
		case bool(b.REQ_1) && b.T_UPPER < b.T_REQUEST_1:
			b.STATUS, b.HEAT = 103, true
		case bool(b.REQ_2) && b.T_UPPER < b.T_REQUEST_2:
			b.STATUS, b.HEAT = 104, true
		}

		// When to stop heating.
		if b.HEAT {
			if b.ENABLE || b.boostMode {
				b.flag0 = true
				if b.T_LOWER_ENABLE && b.T_LOWER > b.T_LOWER_MAX {
					b.boostMode = false
					b.flag0 = b.boostMode
				} else if !b.T_LOWER_ENABLE && b.T_UPPER > b.T_UPPER_MAX {
					b.boostMode = false
					b.flag0 = b.boostMode
				}
			} else {
				b.flag0 = false
			}
			b.flag1 = b.REQ_1 && b.T_UPPER > b.T_REQUEST_1+b.T_REQUEST_HYS
			b.flag2 = b.REQ_2 && b.T_UPPER > b.T_REQUEST_2+b.T_REQUEST_HYS

			b.HEAT = b.flag0 || b.flag1 || b.flag2
			if !b.HEAT {
				b.STATUS = 100
			}
		}
	default:
		b.STATUS, b.HEAT, b.ERROR = 100, false, false
	}
	b.edge = b.BOOST
}

// BURNER controls an oil burner with two stages: IN starts the sequence of
// oil pre heating (until OIL_TEMP, for PRE_HEAT_TIME at most), the motor
// for PRE_VENT_TIME, ignition from PRE_IGNITE_TIME before the oil valve
// COIL1 opens until POST_IGNITE_TIME after, and COIL2 for STAGE2 after
// STAGE2_DELAY. No FLAME within SAFETY_TIME, a flame too early, or
// OVER_TEMP lock it out with FAIL and STATUS 1..9 until RST, after
// LOCKOUT_TIME. It counts the RUNTIME1 and RUNTIME2 of the stages in
// seconds, the CYCLES, and the energy KWH from the powers KW1 and KW2.
// RST_TIMER clears the counters.
type BURNER struct {
	IN                iec.BOOL
	STAGE2            iec.BOOL
	OVER_TEMP         iec.BOOL
	OIL_TEMP          iec.BOOL // default TRUE
	FLAME             iec.BOOL
	RST               iec.BOOL
	RST_TIMER         iec.BOOL
	PRE_HEAT_TIME     iec.TIME // default T#5s
	PRE_VENT_TIME     iec.TIME // default T#15s
	PRE_IGNITE_TIME   iec.TIME // default T#5s
	POST_IGNITE_TIME  iec.TIME // default T#25s
	STAGE2_DELAY      iec.TIME // default T#10s
	SAFETY_TIME       iec.TIME // default T#5s
	LOCKOUT_TIME      iec.TIME // default T#10s
	MULTIPLE_IGNITION iec.BOOL // default TRUE
	KW1               iec.REAL
	KW2               iec.REAL
	MOTOR             iec.BOOL
	COIL1             iec.BOOL
	COIL2             iec.BOOL
	PRE_HEAT          iec.BOOL
	IGNITE            iec.BOOL
	KWH               iec.REAL
	STATUS            iec.BYTE
	FAIL              iec.BOOL
	RUNTIME1          *iec.UDINT
	RUNTIME2          *iec.UDINT
	CYCLES            *iec.UDINT

	state      iec.BYTE
	last       iec.DWORD
	lastChange iec.DWORD
	timer1     engineering.ONTIME
	timer2     engineering.ONTIME
	// oilTempLast is never set in OSCAT.
	oilTempLast iec.BOOL
	cycles2     iec.UDINT
}

// INIT resets the block and sets its inputs to their initial values.
func (b *BURNER) INIT() {
	*b = BURNER{
		RUNTIME1: b.RUNTIME1, RUNTIME2: b.RUNTIME2, CYCLES: b.CYCLES,
		OIL_TEMP:      true,
		PRE_HEAT_TIME: iec.TIME(5 * time.Second), PRE_VENT_TIME: iec.TIME(15 * time.Second),
		PRE_IGNITE_TIME: iec.TIME(5 * time.Second), POST_IGNITE_TIME: iec.TIME(25 * time.Second),
		STAGE2_DELAY: iec.TIME(10 * time.Second), SAFETY_TIME: iec.TIME(5 * time.Second),
		LOCKOUT_TIME: iec.TIME(10 * time.Second), MULTIPLE_IGNITION: true,
	}
}

// off turns all outputs off.
func (b *BURNER) off() {
	b.MOTOR, b.COIL1, b.COIL2, b.IGNITE, b.PRE_HEAT = false, false, false, false, false
}

// Execute runs the block once.
func (b *BURNER) Execute(now time.Time) {
	if b.RUNTIME1 == nil || b.RUNTIME2 == nil || b.CYCLES == nil {
		return
	}
	tx := PLC_MS(now)

	// Reset and over temperature.
	if b.RST || b.OVER_TEMP || b.state == 0 {
		if b.STATUS > 0 && tx-b.lastChange >= ms(b.LOCKOUT_TIME) && bool(b.RST) {
			b.STATUS = 110
			b.FAIL = false
			b.state = 1
		} else {
			b.off()
			if b.OVER_TEMP {
				b.STATUS = 9
				b.FAIL = true
			}
			b.lastChange = tx
			b.last = tx
			b.state = 1
		}
	}

	// Reset the counters.
	if b.RST_TIMER {
		*b.RUNTIME1, *b.RUNTIME2, *b.CYCLES = 0, 0, 0
		b.cycles2 = 0
	}

	// Stop on an error.
	if b.STATUS > 0 && b.STATUS < 100 || bool(b.RST) {
		return
	}

	// lockout stops the burner with an error.
	lockout := func(status iec.BYTE) {
		b.state = 7
		b.STATUS = status
		b.lastChange = tx
	}
	switch b.state {
	case 1: // IN starts the oil pre heating
		if b.IN && b.FLAME {
			b.PRE_HEAT = false
			lockout(2)
		} else if b.IN {
			b.PRE_HEAT = true
			b.state = 2
			b.lastChange = tx
		}

	case 2: // after pre heating start the motor
		switch {
		case tx-b.lastChange >= ms(b.PRE_HEAT_TIME) && bool(b.OIL_TEMP) || bool(b.OIL_TEMP) && !bool(b.oilTempLast):
			b.MOTOR = true
			b.state = 3
			b.lastChange = tx
		case tx-b.lastChange >= ms(b.PRE_HEAT_TIME) && !bool(b.OIL_TEMP):
			// the oil did not heat in time
			b.PRE_HEAT = false
			lockout(1)
		case bool(b.FLAME):
			// a flame before ignition
			b.PRE_HEAT = false
			lockout(2)
		}

	case 3: // wait for the ignition
		switch {
		case tx-b.lastChange >= ms(b.PRE_VENT_TIME)-ms(b.PRE_IGNITE_TIME):
			b.IGNITE = true
			b.state = 4
			b.lastChange = tx
		case bool(b.FLAME):
			b.PRE_HEAT = false
			b.MOTOR = false
			lockout(3)
		}

	case 4: // wait to open the oil valve
		if tx-b.lastChange >= ms(b.PRE_IGNITE_TIME) {
			b.COIL1 = true
			b.state = 5
			b.lastChange = tx
		}

	case 5: // wait for the flame
		if tx-b.lastChange >= ms(b.SAFETY_TIME) || bool(b.FLAME) {
			if !b.FLAME {
				// no flame: emergency stop
				b.MOTOR, b.COIL1, b.PRE_HEAT, b.IGNITE = false, false, false, false
				lockout(4)
			} else {
				b.state = 6
				b.lastChange = tx
			}
		}

	case 6: // the burner runs: watch the flame and end the ignition
		switch {
		case !bool(b.FLAME) && !bool(b.MULTIPLE_IGNITION):
			// the flame went out: emergency stop
			b.off()
			lockout(5)
		case !bool(b.FLAME) && bool(b.MULTIPLE_IGNITION):
			b.IGNITE = true
			b.state = 5
			b.COIL2 = false
			b.lastChange = tx
		default:
			if tx-b.lastChange >= ms(b.POST_IGNITE_TIME) {
				b.IGNITE = false
			}
			b.COIL2 = tx-b.lastChange >= ms(b.STAGE2_DELAY) && b.STAGE2
		}
	}

	// No input stops the burner.
	if !b.IN {
		b.state = 1
		b.off()
		b.lastChange = tx
	}

	// The runtime counters.
	b.timer1.IN = b.FLAME && b.IN && b.MOTOR && b.COIL1 && !b.COIL2
	b.timer1.SECONDS, b.timer1.CYCLES = b.RUNTIME1, b.CYCLES
	b.timer1.Execute(now)
	b.timer2.IN = b.FLAME && b.IN && b.MOTOR && b.COIL1 && b.COIL2
	b.timer2.SECONDS, b.timer2.CYCLES = b.RUNTIME2, &b.cycles2
	b.timer2.Execute(now)
	b.KWH = iec.REAL(*b.RUNTIME1)*b.KW1/3600.0 + iec.REAL(*b.RUNTIME2)*b.KW2/3600.0

	b.last = tx

	// FAIL on an error, and the status of normal operation.
	if b.STATUS > 0 && b.STATUS < 100 {
		b.FAIL = true
		return
	}
	b.FAIL = false
	switch {
	case bool(!b.IN):
		b.STATUS = 110
	case bool(b.FLAME && b.IN && b.MOTOR && b.COIL2 && b.COIL1):
		b.STATUS = 113
	case bool(b.FLAME && b.IN && b.MOTOR && b.COIL1):
		b.STATUS = 112
	default:
		b.STATUS = 111
	}
}

// HEAT_METER measures heat from the flow temperature TF, the return
// temperature TR and the flow LPH in l/h, while E is true: Y is the heat in
// J (or, if PULSE_MODE, X·LPH for each rising edge of E, LPH being the
// liters of a pulse), and C the power over AVG_TIME in J/h. CP, DENSITY and
// CONTENT describe a heat transfer medium mixed in, by the fraction
// CONTENT. RETURN_METER tells the meter is on the return line. RST clears
// the counters.
type HEAT_METER struct {
	TF, TR       iec.REAL
	LPH          iec.REAL
	E            iec.BOOL
	RST          iec.BOOL
	CP           iec.REAL
	DENSITY      iec.REAL
	CONTENT      iec.REAL
	PULSE_MODE   iec.BOOL
	RETURN_METER iec.BOOL
	AVG_TIME     iec.TIME // default T#5s
	C            iec.REAL
	Y            *iec.REAL

	tx, last    iec.DWORD
	int1        engineering.FT_INT2
	edge        iec.BOOL
	x           iec.REAL
	init        iec.BOOL
	yLast       iec.REAL
	initialized bool
}

// INIT resets the block and sets AVG_TIME to its initial value.
func (h *HEAT_METER) INIT() {
	*h = HEAT_METER{Y: h.Y, AVG_TIME: iec.TIME(5 * time.Second), initialized: true}
	h.int1.INIT()
}

// Execute runs the block once.
func (h *HEAT_METER) Execute(now time.Time) {
	if !h.initialized {
		h.initialized = true
		h.int1.INIT()
	}
	if h.Y == nil {
		return
	}
	if h.RST {
		h.int1.RST = true
		h.int1.Execute(now)
		h.int1.RST = false
		h.C = 0.0
		*h.Y = 0.0
	} else if h.E {
		h.x = WATER_DENSITY(SEL(h.RETURN_METER, h.TF, h.TR), false)*(WATER_ENTHALPY(h.TF)-WATER_ENTHALPY(h.TR))*(1.0-h.CONTENT) +
			h.CP*h.DENSITY*h.CONTENT*(h.TF-h.TR)
	}

	// Integrate, or add the pulses.
	h.int1.RUN = !h.PULSE_MODE && h.E
	h.int1.IN = h.x * h.LPH * 2.77777777777e-4
	h.int1.Execute(now)
	if h.PULSE_MODE {
		if !h.edge && h.E {
			*h.Y += h.x * h.LPH
		}
	} else {
		*h.Y = h.int1.OUT
	}
	h.edge = h.E

	h.tx = PLC_MS(now)
	if !h.init {
		h.init = true
		h.last = h.tx
	}

	// The power over AVG_TIME.
	if h.tx-h.last >= ms(h.AVG_TIME) && h.AVG_TIME > 0 {
		h.last = h.tx
		h.C = (*h.Y - h.yLast) * 3.6e6 / iec.REAL(ms(h.AVG_TIME))
		h.yLast = *h.Y
	}
}

// HEAT_TEMP is a heating curve: TY is the flow temperature for the inside
// temperature T_INT + OFFSET and the outside temperature T_EXT, for
// radiators of the exponent C designed for TY_CONFIG at T_INT_CONFIG and
// T_EXT_CONFIG, with the spread T_DIFF, limited to TY_MIN..TY_MAX; 0 when
// T_EXT is within H of the inside temperature. TY is at least T_REQ, and
// HEAT is TY > 0.
type HEAT_TEMP struct {
	T_EXT        iec.REAL
	T_INT        iec.REAL
	OFFSET       iec.REAL
	T_REQ        iec.REAL
	TY_MAX       iec.REAL // default 70.0
	TY_MIN       iec.REAL // default 25.0
	TY_CONFIG    iec.REAL // default 70.0
	T_INT_CONFIG iec.REAL // default 20.0
	T_EXT_CONFIG iec.REAL // default -15.0
	T_DIFF       iec.REAL // default 10.0
	C            iec.REAL // default 1.33
	H            iec.REAL // default 3.0
	TY           iec.REAL
	HEAT         iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (h *HEAT_TEMP) INIT() {
	*h = HEAT_TEMP{
		TY_MAX: 70, TY_MIN: 25, TY_CONFIG: 70, T_INT_CONFIG: 20, T_EXT_CONFIG: -15,
		T_DIFF: 10, C: 1.33, H: 3,
	}
}

// Execute runs the block once.
func (h *HEAT_TEMP) Execute(now time.Time) {
	tr := h.T_INT + h.OFFSET
	tx := (tr - h.T_EXT) / (h.T_INT_CONFIG - h.T_EXT_CONFIG)
	if h.T_EXT+h.H > tr {
		h.TY = 0.0
	} else {
		h.TY = LIMIT(h.TY_MIN, tr+h.T_DIFF*0.5*tx+(h.TY_CONFIG-h.T_DIFF*0.5-tr)*EXPT(tx, 1.0/h.C), h.TY_MAX)
	}
	h.TY = max(h.TY, h.T_REQ)
	h.HEAT = h.TY > 0.0
}

// LEGIONELLA runs the thermal disinfection of a hot water system against
// legionella, every DAY of the week (1 Monday) at T_START, or while MANUAL:
// it heats the boiler to TEMP_SET + TEMP_OFFSET (HEAT, with the hysteresis
// TEMP_HYS), for T_MAX_HEAT at most, then runs the PUMP and opens the
// valves VALVE0..VALVE7 in turn for TP_0..TP_7 each, waiting T_MAX_RET at
// most for TEMP_RETURN to reach TEMP_SET. RUN is true while it runs, and
// STATUS is its SEQUENCE_8's. RST resets it.
type LEGIONELLA struct {
	MANUAL      iec.BOOL
	TEMP_BOILER iec.REAL
	TEMP_RETURN iec.REAL // default 100.0
	DT_IN       iec.DT
	RST         iec.BOOL
	T_START     iec.TOD
	DAY         iec.INT  // default 7
	TEMP_SET    iec.REAL // default 70.0
	TEMP_OFFSET iec.REAL // default 10.0
	TEMP_HYS    iec.REAL // default 5.0
	T_MAX_HEAT  iec.TIME // default T#10m
	T_MAX_RET   iec.TIME // default T#10m
	TP_0        iec.TIME // default T#5m
	TP_1        iec.TIME // default T#5m
	TP_2        iec.TIME // default T#5m
	TP_3        iec.TIME // default T#5m
	TP_4        iec.TIME // default T#5m
	TP_5        iec.TIME // default T#5m
	TP_6        iec.TIME // default T#5m
	TP_7        iec.TIME // default T#5m
	HEAT        iec.BOOL
	PUMP        iec.BOOL
	VALVE0      iec.BOOL
	VALVE1      iec.BOOL
	VALVE2      iec.BOOL
	VALVE3      iec.BOOL
	VALVE4      iec.BOOL
	VALVE5      iec.BOOL
	VALVE6      iec.BOOL
	VALVE7      iec.BOOL
	RUN         iec.BOOL
	STATUS      iec.BYTE

	x1          electrical.TIMER_1
	x2          logic.SEQUENCE_8
	x3          engineering.HYST_1
	init        iec.BOOL
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (l *LEGIONELLA) INIT() {
	tp := iec.TIME(5 * time.Minute)
	*l = LEGIONELLA{
		TEMP_RETURN: 100, T_START: DWORD_TO_TOD(3 * 3600000), DAY: 7,
		TEMP_SET: 70, TEMP_OFFSET: 10, TEMP_HYS: 5,
		T_MAX_HEAT: iec.TIME(10 * time.Minute), T_MAX_RET: iec.TIME(10 * time.Minute),
		TP_0: tp, TP_1: tp, TP_2: tp, TP_3: tp, TP_4: tp, TP_5: tp, TP_6: tp, TP_7: tp,
		initialized: true,
	}
	l.x1.INIT()
	l.x2.INIT()
}

// Execute runs the block once.
func (l *LEGIONELLA) Execute(now time.Time) {
	if !l.initialized {
		l.initialized = true
		l.x1.INIT()
		l.x2.INIT()
	}
	// Set up the timer, the hysteresis and the sequence.
	if !l.init {
		l.init = true
		l.x1.DAY = SHR(iec.BYTE(128), l.DAY)
		l.x1.START = l.T_START
		l.x3.LOW = l.TEMP_OFFSET + l.TEMP_SET
		l.x3.HIGH = l.TEMP_HYS + l.x3.LOW
		x2 := &l.x2
		x2.WAIT0 = l.T_MAX_HEAT
		x2.DELAY0, x2.DELAY1, x2.DELAY2, x2.DELAY3 = l.TP_0, l.TP_1, l.TP_2, l.TP_3
		x2.DELAY4, x2.DELAY5, x2.DELAY6, x2.DELAY7 = l.TP_4, l.TP_5, l.TP_6, l.TP_7
		x2.WAIT1, x2.WAIT2, x2.WAIT3, x2.WAIT4 = l.T_MAX_RET, l.T_MAX_RET, l.T_MAX_RET, l.T_MAX_RET
		x2.WAIT5, x2.WAIT6, x2.WAIT7 = l.T_MAX_RET, l.T_MAX_RET, l.T_MAX_RET
		x2.Execute(now)
	}

	l.x1.DTI = l.DT_IN
	l.x1.Execute(now)
	if l.x1.Q || l.MANUAL || l.x2.RUN {
		l.x3.IN = l.TEMP_BOILER
		l.x3.Execute(now)
		x2 := &l.x2
		x2.IN0 = l.x3.Q || l.x3.WIN
		x2.IN1 = l.TEMP_RETURN >= l.TEMP_SET
		x2.IN2, x2.IN3, x2.IN4, x2.IN5, x2.IN6, x2.IN7 = x2.IN1, x2.IN1, x2.IN1, x2.IN1, x2.IN1, x2.IN1
		x2.RST = l.RST
		x2.START = l.x1.Q || l.MANUAL
		x2.Execute(now)
		l.RUN = x2.RUN
		l.HEAT = !l.x3.Q && x2.RUN
		l.VALVE0, l.VALVE1, l.VALVE2, l.VALVE3 = x2.Q0, x2.Q1, x2.Q2, x2.Q3
		l.VALVE4, l.VALVE5, l.VALVE6, l.VALVE7 = x2.Q4, x2.Q5, x2.Q6, x2.Q7
		l.PUMP = x2.QX
		l.STATUS = x2.STATUS
	} else {
		l.x2.START = false
		l.x2.Execute(now)
		l.STATUS = l.x2.STATUS
	}
}
