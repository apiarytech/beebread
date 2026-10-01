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

package engineering

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/math"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

// ALARM_2 checks X against two pairs of limits with the hysteresis HYS:
// Qn_LO is true below LO_n and Qn_HI above HI_n.
type ALARM_2 struct {
	X                          iec.REAL
	LO_1, HI_1, LO_2, HI_2     iec.REAL
	HYS                        iec.REAL
	Q1_LO, Q1_HI, Q2_LO, Q2_HI iec.BOOL
}

// INIT resets the block.
func (a *ALARM_2) INIT() { *a = ALARM_2{} }

// Execute runs the block once.
func (a *ALARM_2) Execute(now time.Time) {
	tmp := a.X - a.HYS*0.5
	if tmp > a.LO_1 {
		a.Q1_LO = false
	}
	if tmp > a.LO_2 {
		a.Q2_LO = false
	}
	if tmp > a.HI_1 {
		a.Q1_HI = true
	}
	if tmp > a.HI_2 {
		a.Q2_HI = true
	}
	tmp = tmp + a.HYS
	if tmp < a.LO_1 {
		a.Q1_LO = true
	}
	if tmp < a.LO_2 {
		a.Q2_LO = true
	}
	if tmp < a.HI_1 {
		a.Q1_HI = false
	}
	if tmp < a.HI_2 {
		a.Q2_HI = false
	}
}

// BAR_GRAPH shows X on one of 8 outputs: LOW below TRIGGER_LOW, Q1..Q6 for
// 6 equal steps, linear or logarithmic if LOG_SCALE is true, up to
// TRIGGER_HIGH, and HIGH above. ALARM_LOW and ALARM_HIGH make LOW and HIGH
// latch and set ALARM until RST. STATUS is 110 normal, 111 low, 112 high, 1
// low alarm and 2 high alarm.
type BAR_GRAPH struct {
	X                                 iec.REAL
	RST                               iec.BOOL
	TRIGGER_LOW, TRIGGER_HIGH         iec.REAL
	ALARM_LOW, ALARM_HIGH             iec.BOOL
	LOG_SCALE                         iec.BOOL
	LOW, Q1, Q2, Q3, Q4, Q5, Q6, HIGH iec.BOOL
	ALARM                             iec.BOOL
	STATUS                            iec.BYTE

	init iec.BOOL
	t    [5]iec.REAL
}

// INIT resets the block.
func (b *BAR_GRAPH) INIT() { *b = BAR_GRAPH{} }

// Execute runs the block once.
func (b *BAR_GRAPH) Execute(now time.Time) {
	if !b.init {
		b.init = true
		if b.LOG_SCALE {
			temp := EXP(LN(b.TRIGGER_HIGH/b.TRIGGER_LOW) * 0.166666666666666666666)
			b.t[0] = b.TRIGGER_LOW * temp
			for i := 1; i < 5; i++ {
				b.t[i] = b.t[i-1] * temp
			}
		} else {
			temp := (b.TRIGGER_HIGH - b.TRIGGER_LOW) * 0.142857142
			b.t[0] = b.TRIGGER_LOW + temp
			for i := 1; i < 5; i++ {
				b.t[i] = b.t[i-1] + temp
			}
		}
	}
	b.Q1, b.Q2, b.Q3, b.Q4, b.Q5, b.Q6 = false, false, false, false, false, false
	b.STATUS = 110
	if !b.ALARM_LOW {
		b.LOW = false
	}
	if !b.ALARM_HIGH {
		b.HIGH = false
	}
	if b.RST {
		b.ALARM, b.LOW, b.HIGH = false, false, false
	}
	x := b.X
	switch {
	case x < b.TRIGGER_LOW:
		b.LOW = true
		b.STATUS = 111
		if b.ALARM_LOW {
			b.ALARM = true
			b.STATUS = 1
		}
	case x < b.t[0]:
		b.Q1 = true
	case x < b.t[1]:
		b.Q2 = true
	case x < b.t[2]:
		b.Q3 = true
	case x < b.t[3]:
		b.Q4 = true
	case x < b.t[4]:
		b.Q5 = true
	case x < b.TRIGGER_HIGH:
		b.Q6 = true
	default:
		b.HIGH = true
		b.STATUS = 112
		if b.ALARM_HIGH {
			b.ALARM = true
			b.STATUS = 2
		}
	}
}

// CALIBRATE calibrates an analog value X to Y: while CO is true the offset
// is set so Y is Y_OFFSET, and then while CS is true the scale is set so Y
// is Y_SCALE.
type CALIBRATE struct {
	X        iec.REAL
	CO, CS   iec.BOOL
	Y_OFFSET iec.REAL
	Y_SCALE  iec.REAL
	Y        iec.REAL

	offset      iec.REAL
	scale       iec.REAL
	initialized bool
}

// INIT resets the block.
func (c *CALIBRATE) INIT() { *c = CALIBRATE{scale: 1, initialized: true} }

// Execute runs the block once.
func (c *CALIBRATE) Execute(now time.Time) {
	if !c.initialized {
		c.initialized = true
		c.scale = 1
	}
	if c.CO {
		c.offset = c.Y_OFFSET - c.X
	} else if c.CS {
		c.scale = c.Y_SCALE / (c.X + c.offset)
	}
	c.Y = (c.X + c.offset) * c.scale
}

// CYCLE_TIME measures the time between its runs: the minimum CT_MIN, the
// maximum CT_MAX and the last CT_LAST, the time running SYSTIME and SYSDAYS,
// and the number of runs CYCLES. RST clears them.
type CYCLE_TIME struct {
	RST     iec.BOOL
	CT_MIN  iec.TIME
	CT_MAX  iec.TIME
	CT_LAST iec.TIME
	SYSTIME iec.TIME
	SYSDAYS iec.INT
	CYCLES  iec.DWORD

	lastCycle iec.DWORD
	init      iec.BOOL
}

// INIT resets the block.
func (c *CYCLE_TIME) INIT() { *c = CYCLE_TIME{} }

// Execute runs the block once.
func (c *CYCLE_TIME) Execute(now time.Time) {
	tx := PLC_MS(now) - c.lastCycle
	t := DWORD_TO_TIME(tx)
	switch {
	case bool(c.RST):
		c.CT_MIN = iec.TIME(10 * time.Hour)
		c.CT_MAX = 0
		c.CYCLES = 0
	case c.lastCycle > 0:
		if t < c.CT_MIN {
			c.CT_MIN = t
		} else if t > c.CT_MAX {
			c.CT_MAX = t
		}
		c.CT_LAST = t
	case c.CT_MIN == 0:
		// The largest TIME, as t#0s - t#1ms is.
		c.CT_MIN = DWORD_TO_TIME(0xFFFFFFFF)
	}
	if c.init {
		c.SYSTIME += t
		if c.SYSTIME >= iec.TIME(24*time.Hour) {
			c.SYSTIME -= iec.TIME(24 * time.Hour)
			c.SYSDAYS++
		}
	}
	c.init = true
	c.lastCycle += tx
	c.CYCLES++
}

// DT_SIMU simulates a clock DTS that starts at START and runs SPEED times as
// fast as real time, or one second each run if SPEED is 0.
type DT_SIMU struct {
	START iec.DT
	SPEED iec.REAL // default 1.0
	DTS   iec.DT

	init iec.BOOL
	last iec.DWORD
}

// INIT resets the block and sets SPEED to its initial value.
func (d *DT_SIMU) INIT() { *d = DT_SIMU{SPEED: 1} }

// Execute runs the block once.
func (d *DT_SIMU) Execute(now time.Time) {
	tx := PLC_MS(now)
	tc := REAL_TO_DWORD(iec.REAL(tx-d.last) * d.SPEED)
	switch {
	case !bool(d.init):
		d.init = true
		d.DTS = d.START
		d.last = tx
	case d.SPEED == 0.0:
		d.DTS = DWORD_TO_DT(DT_TO_DWORD(d.DTS) + 1)
	case tc >= 1000:
		t := tc / 1000 * 1000
		d.DTS = DWORD_TO_DT(DT_TO_DWORD(d.DTS) + t/1000)
		d.last += REAL_TO_DWORD(iec.REAL(t) / d.SPEED)
	}
}

// FLOW_METER measures a flow: while E is true it integrates the flow VX, in
// units per hour, into X and whole units into Y, or with PULSE_MODE it adds
// VX to X on each rising edge of E. F is the flow in units per hour,
// updated every UPDATE_TIME.
type FLOW_METER struct {
	VX          iec.REAL
	E           iec.BOOL
	RST         iec.BOOL
	PULSE_MODE  iec.BOOL
	UPDATE_TIME iec.TIME // default T#1s
	F           iec.REAL
	X           *iec.REAL
	Y           *iec.UDINT

	tl    iec.DWORD
	int1  INTEGRATE
	init  iec.BOOL
	eLast iec.BOOL
	xLast iec.REAL
	yLast iec.UDINT
}

// INIT resets the block and sets UPDATE_TIME to its initial value.
func (f *FLOW_METER) INIT() { *f = FLOW_METER{X: f.X, Y: f.Y, UPDATE_TIME: iec.TIME(time.Second)} }

// Execute runs the block once.
func (f *FLOW_METER) Execute(now time.Time) {
	if f.X == nil || f.Y == nil {
		return
	}
	tx := PLC_MS(now)
	if !f.init {
		f.init = true
		// OSCAT takes the time before it has read it, which is 0.
		f.tl = 0
		f.xLast = *f.X
		f.yLast = *f.Y
		f.int1.INIT()
		f.int1.K = 2.7777777777777777e-4
	}
	// Gated operation.
	f.int1.E = !(f.RST || f.PULSE_MODE) && f.E
	f.int1.X = f.VX
	f.int1.Y = f.X
	f.int1.Execute(now)
	if f.RST {
		*f.X = 0
		*f.Y = 0
		f.tl = tx
		f.xLast = 0
		f.yLast = 0
	} else if f.E && f.PULSE_MODE && !f.eLast {
		*f.X += f.VX
	}
	f.eLast = f.E
	if *f.X > 1.0 {
		tmp := math.FLOOR(*f.X)
		*f.Y += iec.UDINT(tmp)
		*f.X -= iec.REAL(tmp)
	}
	if tx-f.tl >= ms(f.UPDATE_TIME) && f.UPDATE_TIME > 0 {
		f.F = (iec.REAL(*f.Y-f.yLast) + *f.X - f.xLast) / iec.REAL(tx-f.tl) * 3.6e6
		f.yLast = *f.Y
		f.xLast = *f.X
		f.tl = tx
	}
}

// M_D measures the time from a rising edge of START to STOP: ET counts
// while it runs and PT is the last time measured. A time over TMAX, or RST,
// clears them.
type M_D struct {
	START, STOP iec.BOOL
	TMAX        iec.TIME // default T#10d
	RST         iec.BOOL
	PT          iec.TIME
	ET          iec.TIME
	RUN         iec.BOOL

	edge    iec.BOOL
	t0      iec.DWORD
	startup iec.BOOL
}

// INIT resets the block and sets TMAX to its initial value.
func (m *M_D) INIT() { *m = M_D{TMAX: iec.TIME(240 * time.Hour)} }

// Execute runs the block once.
func (m *M_D) Execute(now time.Time) {
	if m.RST || m.ET >= m.TMAX {
		m.PT, m.ET = 0, 0
		m.startup = false
		m.RUN = false
	}
	if !m.startup {
		m.edge = m.START
		m.startup = true
	}
	tx := PLC_MS(now)
	if m.START && !m.edge && !m.STOP {
		m.t0 = tx
		m.RUN = true
		m.PT = 0
	} else if m.STOP && m.RUN {
		m.PT = m.ET
		m.RUN = false
	}
	m.edge = m.START
	if m.RUN {
		m.ET = DWORD_TO_TIME(tx - m.t0)
	}
}

// M_T measures the width of a pulse of IN: ET counts while IN is true and
// PT is the last width. A pulse over TMAX, or RST, clears them.
type M_T struct {
	IN   iec.BOOL
	TMAX iec.TIME // default T#10d
	RST  iec.BOOL
	PT   iec.TIME
	ET   iec.TIME

	edge  iec.BOOL
	start iec.DWORD
}

// INIT resets the block and sets TMAX to its initial value.
func (m *M_T) INIT() { *m = M_T{TMAX: iec.TIME(240 * time.Hour)} }

// Execute runs the block once.
func (m *M_T) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(m.RST) || m.ET >= m.TMAX:
		m.PT, m.ET = 0, 0
	case bool(m.IN):
		if !m.edge {
			m.start = tx
		}
		m.ET = DWORD_TO_TIME(tx - m.start)
	default:
		m.PT = m.ET
	}
	m.edge = m.IN
}

// M_TX measures a signal IN: the high time TH, the low time TL, the duty
// cycle DC, the frequency F in Hz, and the time since the last rising edge
// ET. A period over TMAX, or RST, clears them.
type M_TX struct {
	IN   iec.BOOL
	TMAX iec.TIME // default T#10d
	RST  iec.BOOL
	TH   iec.TIME
	TL   iec.TIME
	DC   iec.REAL
	F    iec.REAL
	ET   iec.TIME

	edge        iec.BOOL
	start, stop iec.DWORD
	rise, fall  iec.BOOL
	startup     iec.BOOL
}

// INIT resets the block and sets TMAX to its initial value.
func (m *M_TX) INIT() { *m = M_TX{TMAX: iec.TIME(240 * time.Hour)} }

// Execute runs the block once.
func (m *M_TX) Execute(now time.Time) {
	if m.RST || m.ET >= m.TMAX {
		m.rise, m.fall, m.startup = false, false, false
		m.TH, m.TL, m.DC, m.F, m.ET = 0, 0, 0, 0, 0
	}
	if !m.startup {
		m.edge = m.IN
		m.startup = true
	}
	tx := PLC_MS(now)
	ratio := func() {
		if m.TH > 0 && m.TL > 0 {
			m.DC = TIME_TO_REAL(m.TH) / TIME_TO_REAL(m.TH+m.TL)
			m.F = 1000.0 / TIME_TO_REAL(m.TH+m.TL)
		}
	}
	if m.IN != m.edge {
		m.edge = m.IN
		if m.IN {
			m.start = tx
			m.rise = true
			if m.fall {
				m.TL = DWORD_TO_TIME(m.start - m.stop)
			}
		} else {
			m.stop = tx
			m.fall = true
			if m.rise {
				m.TH = DWORD_TO_TIME(m.stop - m.start)
			}
		}
		ratio()
	}
	if m.rise {
		m.ET = DWORD_TO_TIME(tx - m.start)
	}
}

// METER sums a consumption over time in MX: M1 while I1 is true plus M2
// while I2 is true, per second, divided by D. It sums with double precision.
type METER struct {
	M1, M2 iec.REAL
	I1, I2 iec.BOOL
	D      iec.REAL // default 1.0
	RST    iec.BOOL
	MX     *iec.REAL

	mr   REAL2
	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets D to its initial value.
func (m *METER) INIT() { *m = METER{MX: m.MX, D: 1} }

// Execute runs the block once.
func (m *METER) Execute(now time.Time) {
	if m.MX == nil {
		return
	}
	tx := PLC_MS(now)
	var tc iec.REAL
	if !m.init {
		m.init = true
		m.last = tx
		m.mr.RX = *m.MX
		m.mr.R1 = 0.0
	} else if tx == m.last {
		return
	} else {
		tc = iec.REAL(tx-m.last) * 0.001
	}
	m.last = tx
	if m.RST {
		m.mr = REAL2{}
		return
	}
	var mx1, mx2 iec.REAL
	if m.I1 {
		mx1 = m.M1
	}
	if m.I2 {
		mx2 = m.M2
	}
	m.mr = math.R2_ADD(m.mr, (mx1+mx2)/m.D*tc)
	*m.MX = m.mr.RX
}

// METER_STAT keeps the statistics of a meter reading IN on the date DI:
// the consumption of the current and last day, week, month and year.
type METER_STAT struct {
	IN            iec.REAL
	DI            iec.DATE
	RST           iec.BOOL
	LAST_DAY      *iec.REAL
	CURRENT_DAY   *iec.REAL
	LAST_WEEK     *iec.REAL
	CURRENT_WEEK  *iec.REAL
	LAST_MONTH    *iec.REAL
	CURRENT_MONTH *iec.REAL
	LAST_YEAR     *iec.REAL
	CURRENT_YEAR  *iec.REAL

	yearStart, monthStart, weekStart, dayStart iec.REAL
	lastRun                                    iec.DATE
}

// INIT resets the block.
func (m *METER_STAT) INIT() {
	*m = METER_STAT{LAST_DAY: m.LAST_DAY, CURRENT_DAY: m.CURRENT_DAY, LAST_WEEK: m.LAST_WEEK,
		CURRENT_WEEK: m.CURRENT_WEEK, LAST_MONTH: m.LAST_MONTH, CURRENT_MONTH: m.CURRENT_MONTH,
		LAST_YEAR: m.LAST_YEAR, CURRENT_YEAR: m.CURRENT_YEAR}
}

// Execute runs the block once.
func (m *METER_STAT) Execute(now time.Time) {
	for _, p := range []*iec.REAL{m.LAST_DAY, m.CURRENT_DAY, m.LAST_WEEK, m.CURRENT_WEEK,
		m.LAST_MONTH, m.CURRENT_MONTH, m.LAST_YEAR, m.CURRENT_YEAR} {
		if p == nil {
			return
		}
	}
	if m.RST {
		*m.LAST_DAY, *m.CURRENT_DAY, m.dayStart = 0, 0, m.IN
		*m.LAST_WEEK, *m.CURRENT_WEEK, m.weekStart = 0, 0, m.IN
		*m.LAST_MONTH, *m.CURRENT_MONTH, m.monthStart = 0, 0, m.IN
		*m.LAST_YEAR, *m.CURRENT_YEAR, m.yearStart = 0, 0, m.IN
	} else {
		*m.CURRENT_DAY = m.IN - m.dayStart
		*m.CURRENT_WEEK = m.IN - m.weekStart
		*m.CURRENT_MONTH = m.IN - m.monthStart
		*m.CURRENT_YEAR = m.IN - m.yearStart
	}
	newDay := func() {
		*m.LAST_DAY, *m.CURRENT_DAY, m.dayStart = *m.CURRENT_DAY, 0, m.IN
	}
	newMonth := func() {
		*m.LAST_MONTH, *m.CURRENT_MONTH, m.monthStart = *m.CURRENT_MONTH, 0, m.IN
		newDay()
	}
	switch {
	case td.YEAR_OF_DATE(m.DI) > td.YEAR_OF_DATE(m.lastRun):
		*m.LAST_YEAR, *m.CURRENT_YEAR, m.yearStart = *m.CURRENT_YEAR, 0, m.IN
		newMonth()
	case td.MONTH_OF_DATE(m.DI) > td.MONTH_OF_DATE(m.lastRun):
		newMonth()
	case td.DAY_OF_YEAR(m.DI) > td.DAY_OF_YEAR(m.lastRun):
		newDay()
	}
	if td.DAY_OF_WEEK(m.DI) < td.DAY_OF_WEEK(m.lastRun) {
		*m.LAST_WEEK, *m.CURRENT_WEEK, m.weekStart = *m.CURRENT_WEEK, 0, m.IN
	}
	m.lastRun = m.DI
}

// ONTIME measures the time IN is true, in SECONDS, and counts its rising
// edges in CYCLES.
type ONTIME struct {
	IN      iec.BOOL
	RST     iec.BOOL
	SECONDS *iec.UDINT
	CYCLES  *iec.UDINT

	last iec.DWORD
	edge iec.BOOL
	init iec.BOOL
	ms   iec.DWORD
}

// INIT resets the block.
func (o *ONTIME) INIT() { *o = ONTIME{SECONDS: o.SECONDS, CYCLES: o.CYCLES} }

// Execute runs the block once.
func (o *ONTIME) Execute(now time.Time) {
	if o.SECONDS == nil || o.CYCLES == nil {
		return
	}
	tx := PLC_MS(now)
	if !o.init {
		o.init = true
		o.last = tx
		o.ms = 0
	}
	if o.RST {
		*o.SECONDS = 0
		*o.CYCLES = 0
		o.last = tx
		o.ms = 0
	} else if o.IN {
		o.ms += tx - o.last
		if o.ms >= 1000 {
			*o.SECONDS++
			o.ms -= 1000
		}
		if !o.edge {
			*o.CYCLES++
		}
	}
	o.last = tx
	o.edge = o.IN
}

// TC_MS gives the time since its last run in TC, in milliseconds.
type TC_MS struct {
	TC iec.DWORD

	init iec.BOOL
	last iec.DWORD
}

// INIT resets the block.
func (t *TC_MS) INIT() { *t = TC_MS{} }

// Execute runs the block once.
func (t *TC_MS) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !t.init {
		t.init = true
		t.TC = 0
	} else {
		t.TC = tx - t.last
	}
	t.last = tx
}

// TC_S gives the time since its last run in TC, in seconds.
type TC_S struct {
	TC iec.REAL

	init iec.BOOL
	last iec.DWORD
}

// INIT resets the block.
func (t *TC_S) INIT() { *t = TC_S{} }

// Execute runs the block once.
func (t *TC_S) Execute(now time.Time) {
	tx := PLC_US(now)
	if !t.init {
		t.init = true
		t.TC = 0
	} else {
		t.TC = iec.REAL(tx-t.last) * 1.0e-6
	}
	t.last = tx
}

// TC_US gives the time since its last run in TC, in microseconds.
type TC_US struct {
	TC iec.DWORD

	init iec.BOOL
	last iec.DWORD
}

// INIT resets the block.
func (t *TC_US) INIT() { *t = TC_US{} }

// Execute runs the block once.
func (t *TC_US) Execute(now time.Time) {
	tx := PLC_US(now)
	if !t.init {
		t.init = true
		t.TC = 0
	} else {
		t.TC = tx - t.last
	}
	t.last = tx
}
