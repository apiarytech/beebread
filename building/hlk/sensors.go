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
	"github.com/apiarytech/beebread/basic/math"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/beebread/building/actuators"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// T_AVG24 averages the outside temperature TS, in 0.1 °C, filtered over
// T_FILTER, over 24 hours, sampled every 30 minutes at the date and time
// DTI: TA is the filtered temperature, T24 the average, T24_MAX and T24_MIN
// the extremes, all scaled by SCALE after the offset OFS. TP is true for one
// scan when they change. At the start, or on RST, the samples are T24, or
// TS if T24 is -1000.
type T_AVG24 struct {
	TS       iec.INT
	DTI      iec.DT
	RST      iec.BOOL
	T_FILTER iec.TIME // default T#10m
	SCALE    iec.REAL // default 1.0
	OFS      iec.REAL
	TA       iec.REAL
	TP       iec.BOOL
	T24      *iec.REAL
	T24_MAX  *iec.REAL
	T24_MIN  *iec.REAL

	samples [48]iec.INT
	pos     iec.INT
	init    iec.BOOL
	ft1     engineering.FILTER_I
	sum     iec.DINT
	last    iec.DT
}

// INIT resets the block and sets its inputs to their initial values.
func (t *T_AVG24) INIT() {
	*t = T_AVG24{
		T24: t.T24, T24_MAX: t.T24_MAX, T24_MIN: t.T24_MIN,
		T_FILTER: iec.TIME(10 * time.Minute), SCALE: 1,
	}
}

// Execute runs the block once.
func (t *T_AVG24) Execute(now time.Time) {
	if t.T24 == nil || t.T24_MAX == nil || t.T24_MIN == nil {
		return
	}
	// Filter the sensor.
	t.ft1.X, t.ft1.T = t.TS, t.T_FILTER
	t.ft1.Execute(now)
	scaled := func(x iec.REAL) iec.REAL { return (x + t.OFS) * t.SCALE }

	dti := DT_TO_DWORD(t.DTI)
	switch {
	case bool(t.RST || !t.init):
		t.init = true
		if *t.T24 == -1000.0 {
			*t.T24 = iec.REAL(t.ft1.Y) * 0.1
		}
		for i := range t.samples {
			t.samples[i] = REAL_TO_INT(*t.T24 * 10.0)
		}
		t.pos = 0
		t.sum = iec.DINT(t.samples[0]) * 48
		t.TA = scaled(iec.REAL(t.ft1.Y) * 0.1)
		*t.T24 = scaled(iec.REAL(t.sum) * 0.00208333333333)
		t.TP = true
	case dti/60%30 == 0 && dti > DT_TO_DWORD(t.last):
		// A sample every 30 minutes, once.
		t.last = DWORD_TO_DT(dti + 60)
		// sum is the sum of the samples, a ring buffer.
		t.sum -= iec.DINT(t.samples[t.pos])
		t.samples[t.pos] = t.ft1.Y
		t.sum += iec.DINT(t.samples[t.pos])
		t.pos = math.INC1(t.pos, 48)
		t.TA = scaled(iec.REAL(t.ft1.Y) * 0.1)
		*t.T24 = scaled(iec.REAL(t.sum) * 0.00208333333333)
		// The extremes of the last 24 hours.
		tmpMax, tmpMin := iec.INT(-32000), iec.INT(32000)
		for _, s := range t.samples {
			tmpMax = max(tmpMax, s)
			tmpMin = min(tmpMin, s)
		}
		*t.T24_MAX = scaled(iec.REAL(tmpMax) * 0.1)
		*t.T24_MIN = scaled(iec.REAL(tmpMin) * 0.1)
		t.TP = true
	default:
		t.TP = false
	}
}

// TANK_LEVEL fills a tank: the VALVE opens while the LEVEL switch, delayed
// by LEVEL_DELAY_TIME, reports a low level. A LEAK, or the valve open
// longer than MAX_VALVE_TIME if that is not 0, closes it with an ALARM
// (STATUS 1 or 2) until ACLR.
type TANK_LEVEL struct {
	LEVEL            iec.BOOL
	LEAK             iec.BOOL
	ACLR             iec.BOOL
	MAX_VALVE_TIME   iec.TIME
	LEVEL_DELAY_TIME iec.TIME
	VALVE            iec.BOOL
	ALARM            iec.BOOL
	STATUS           iec.BYTE

	cx          actuators.ACTUATOR_COIL
	tn          timers.TON
	tl          logic.TONOF
	open        iec.BOOL
	initialized bool
}

// INIT resets the block.
func (t *TANK_LEVEL) INIT() {
	*t = TANK_LEVEL{initialized: true}
	t.cx.INIT()
}

// coil runs the valve's coil with the input in.
func (t *TANK_LEVEL) coil(in iec.BOOL, now time.Time) {
	t.cx.IN = in
	t.cx.Execute(now)
}

// Execute runs the block once.
func (t *TANK_LEVEL) Execute(now time.Time) {
	if !t.initialized {
		t.initialized = true
		t.cx.INIT()
	}
	// The level, delayed.
	t.tl.IN, t.tl.T_ON, t.tl.T_OFF = t.LEVEL, t.LEVEL_DELAY_TIME, t.LEVEL_DELAY_TIME
	t.tl.Execute(now)
	t.open = t.tl.Q

	switch {
	case bool(t.ALARM):
		// An alarm waits for ACLR.
		if t.ACLR {
			t.ALARM = false
			t.STATUS = 101 // cleared
			t.coil(false, now)
		}
		return
	case bool(t.LEAK):
		t.coil(false, now)
		t.ALARM = true
		t.STATUS = 1 // leak
	case bool(t.open):
		t.coil(true, now)
		t.STATUS = 102 // open by a low level
	default:
		t.coil(false, now)
		t.STATUS = 100 // idle
	}

	// The valve open too long is an alarm.
	t.tn.IN, t.tn.PT = t.cx.OUT && t.MAX_VALVE_TIME > 0, t.MAX_VALVE_TIME
	t.tn.Execute(now)
	if t.tn.Q {
		t.ALARM = true
		t.STATUS = 2 // open too long
		t.coil(false, now)
	}
	t.VALVE = t.cx.OUT
}

// TEMP_EXT gives the outside temperature T_EXT from up to three sensors,
// as MULTI_IN selects by T_EXT_CONFIG within T_EXT_MIN..T_EXT_MAX, and
// decides on heating and cooling: HEAT within the heating period, from
// HEAT_PERIOD_START to HEAT_PERIOD_STOP (only the month and day count),
// when T_EXT falls to HEAT_START_TEMP_DAY by day (START_DAY to
// START_NIGHT) or HEAT_START_TEMP_NIGHT by night, until HEAT_STOP_TEMP;
// COOL likewise. It runs every CYCLE_TIME at the date and time DT_IN.
type TEMP_EXT struct {
	T_EXT1                iec.REAL
	T_EXT2                iec.REAL
	T_EXT3                iec.REAL
	T_EXT_CONFIG          iec.BYTE
	DT_IN                 iec.DT
	T_EXT_MIN             iec.REAL // default -40.0
	T_EXT_MAX             iec.REAL // default 60.0
	T_EXT_DEFAULT         iec.REAL // default -10.0
	HEAT_PERIOD_START     iec.DATE // default D#1970-09-01
	HEAT_PERIOD_STOP      iec.DATE // default D#1970-04-30
	COOL_PERIOD_START     iec.DATE // default D#1970-04-01
	COOL_PERIOD_STOP      iec.DATE // default D#1970-09-30
	HEAT_START_TEMP_DAY   iec.REAL // default 15.0
	HEAT_START_TEMP_NIGHT iec.REAL // default 10.0
	HEAT_STOP_TEMP        iec.REAL // default 18.0
	COOL_START_TEMP_DAY   iec.REAL // default 26.0
	COOL_START_TEMP_NIGHT iec.REAL // default 26.0
	COOL_STOP_TEMP        iec.REAL // default 24.0
	START_DAY             iec.TOD  // default TOD#09:00:00
	START_NIGHT           iec.TOD  // default TOD#21:00:00
	CYCLE_TIME            iec.TIME // default T#10m
	T_EXT                 iec.REAL
	HEAT                  iec.BOOL
	COOL                  iec.BOOL

	lastRun   iec.DWORD
	init      iec.BOOL
	coolStart iec.DWORD
	coolStop  iec.DWORD
	heatStart iec.DWORD
	heatStop  iec.DWORD
}

// INIT resets the block and sets its inputs to their initial values.
func (t *TEMP_EXT) INIT() {
	*t = TEMP_EXT{
		T_EXT_MIN: -40, T_EXT_MAX: 60, T_EXT_DEFAULT: -10,
		HEAT_PERIOD_START: td.SET_DATE(1970, 9, 1), HEAT_PERIOD_STOP: td.SET_DATE(1970, 4, 30),
		COOL_PERIOD_START: td.SET_DATE(1970, 4, 1), COOL_PERIOD_STOP: td.SET_DATE(1970, 9, 30),
		HEAT_START_TEMP_DAY: 15, HEAT_START_TEMP_NIGHT: 10, HEAT_STOP_TEMP: 18,
		COOL_START_TEMP_DAY: 26, COOL_START_TEMP_NIGHT: 26, COOL_STOP_TEMP: 24,
		START_DAY: DWORD_TO_TOD(9 * 3600000), START_NIGHT: DWORD_TO_TOD(21 * 3600000),
		CYCLE_TIME: iec.TIME(10 * time.Minute),
	}
}

// in1972 returns the day of the date in 1972, a leap year, so that dates
// compare by their month and day.
func in1972(d iec.DATE) iec.DWORD {
	return DATE_TO_DWORD(td.SET_DATE(1972, td.MONTH_OF_DATE(d), td.DAY_OF_MONTH(d)))
}

// Execute runs the block once.
func (t *TEMP_EXT) Execute(now time.Time) {
	if !t.init {
		t.init = true
		t.heatStart = in1972(t.HEAT_PERIOD_START)
		t.heatStop = in1972(t.HEAT_PERIOD_STOP)
		t.coolStart = in1972(t.COOL_PERIOD_START)
		t.coolStop = in1972(t.COOL_PERIOD_STOP)
	}

	// The block runs every CYCLE_TIME.
	tx := PLC_MS(now)
	if tx-t.lastRun < ms(t.CYCLE_TIME) {
		return
	}

	xdate := in1972(DT_TO_DATE(t.DT_IN))
	tod := TOD_TO_DWORD(DT_TO_TOD(t.DT_IN))
	day := tod >= TOD_TO_DWORD(t.START_DAY) && tod < TOD_TO_DWORD(t.START_NIGHT)

	t.T_EXT = engineering.MULTI_IN(t.T_EXT1, t.T_EXT2, t.T_EXT3, t.T_EXT_DEFAULT, t.T_EXT_MIN, t.T_EXT_MAX, t.T_EXT_CONFIG)

	// within reports whether xdate is in the period from start to stop.
	within := func(start, stop iec.DWORD) bool {
		return start <= stop && xdate >= start && xdate <= stop || start > stop && (xdate >= start || xdate <= stop)
	}

	// Heating.
	if within(t.heatStart, t.heatStop) {
		switch {
		case day && t.T_EXT <= t.HEAT_START_TEMP_DAY:
			t.HEAT = true
		case !day && t.T_EXT <= t.HEAT_START_TEMP_NIGHT:
			t.HEAT = true
		case t.T_EXT >= t.HEAT_STOP_TEMP:
			t.HEAT = false
		}
	} else {
		t.HEAT = false
	}

	// Cooling.
	if within(t.coolStart, t.coolStop) {
		switch {
		case day && t.T_EXT >= t.COOL_START_TEMP_DAY:
			t.COOL = true
		case !day && t.T_EXT >= t.COOL_START_TEMP_NIGHT:
			t.COOL = true
		case t.T_EXT <= t.COOL_STOP_TEMP:
			t.COOL = false
		}
	} else {
		t.COOL = false
	}
	t.lastRun = tx
}
