/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package engineering

import (
	"math"
	"time"

	"beebread/basic"
	"beebread/basic/logic"
	"beebread/basic/time_date"

	beeMath "beebread/basic/math"
)

// ALARM_2 checks two pairs of limits and signals when the input is above or below a set limit.
type ALARM_2 struct {
	Q1Lo bool
	Q1Hi bool
	Q2Lo bool
	Q2Hi bool
}

// Update executes the alarm logic.
func (a *ALARM_2) Update(x, lo1, hi1, lo2, hi2, hys float64) {
	tmp := x - hys*0.5
	if tmp > lo1 {
		a.Q1Lo = false
	}
	if tmp > lo2 {
		a.Q2Lo = false
	}
	if tmp > hi1 {
		a.Q1Hi = true
	}
	if tmp > hi2 {
		a.Q2Hi = true
	}

	tmp += hys
	if tmp < lo1 {
		a.Q1Lo = true
	}
	if tmp < lo2 {
		a.Q2Lo = true
	}
	if tmp < hi1 {
		a.Q1Hi = false
	}
	if tmp < hi2 {
		a.Q2Hi = false
	}
}

// BAR_GRAPH is a multi-window comparator that displays an analog signal on 8 digital outputs.
type BAR_GRAPH struct {
	Low    bool
	Q1     bool
	Q2     bool
	Q3     bool
	Q4     bool
	Q5     bool
	Q6     bool
	High   bool
	Alarm  bool
	Status byte

	// internal state
	init bool
	t    [5]float64
}

// Update executes the bar graph logic.
func (b *BAR_GRAPH) Update(x, triggerLow, triggerHigh float64, rst, alarmLow, alarmHigh, logScale bool) {
	if !b.init {
		b.init = true
		if logScale {
			temp := math.Exp(math.Log(triggerHigh/triggerLow) * 0.16666666666666666)
			b.t[0] = triggerLow * temp
			for i := 1; i < 5; i++ {
				b.t[i] = b.t[i-1] * temp
			}
		} else {
			temp := (triggerHigh - triggerLow) * 0.142857142
			b.t[0] = triggerLow + temp
			for i := 1; i < 5; i++ {
				b.t[i] = b.t[i-1] + temp
			}
		}
	}

	b.Q1, b.Q2, b.Q3, b.Q4, b.Q5, b.Q6 = false, false, false, false, false, false
	b.Status = 110

	if !alarmLow {
		b.Low = false
	}
	if !alarmHigh {
		b.High = false
	}
	if rst {
		b.Alarm, b.Low, b.High = false, false, false
	}

	if x < triggerLow {
		b.Low = true
		b.Status = 111
		if alarmLow {
			b.Alarm = true
			b.Status = 1
		}
	} else if x < b.t[0] {
		b.Q1 = true
	} else if x < b.t[1] {
		b.Q2 = true
	} else if x < b.t[2] {
		b.Q3 = true
	} else if x < b.t[3] {
		b.Q4 = true
	} else if x < b.t[4] {
		b.Q5 = true
	} else if x < triggerHigh {
		b.Q6 = true
	} else {
		b.High = true
		b.Status = 112
		if alarmHigh {
			b.Alarm = true
			b.Status = 2
		}
	}
}

// CALIBRATE allows for offset and scale calibration of an analog input.
type CALIBRATE struct {
	Y float64
	// RETAIN fields
	Offset float64
	Scale  float64
}

// NewCALIBRATE creates a CALIBRATE struct with default scale.
func NewCALIBRATE() *CALIBRATE {
	return &CALIBRATE{Scale: 1.0}
}

// Update executes the calibration logic.
func (c *CALIBRATE) Update(x float64, co, cs bool, yOffset, yScale float64) {
	if co {
		c.Offset = yOffset - x
	} else if cs {
		c.Scale = yScale / (x + c.Offset)
	}
	c.Y = (x + c.Offset) * c.Scale
}

// CYCLE_TIME measures PLC cycle time statistics.
type CYCLE_TIME struct {
	CtMin   time.Duration
	CtMax   time.Duration
	CtLast  time.Duration
	SysTime time.Duration
	SysDays int
	Cycles  uint32

	// internal state
	lastCycle time.Duration
	init      bool
}

// Update executes the cycle time measurement.
func (c *CYCLE_TIME) Update(rst bool) {
	tx := time.Duration(logic.T_PLC_US()) * time.Microsecond
	elapsed := tx - c.lastCycle

	if rst {
		c.CtMin = 10 * time.Hour // A large value
		c.CtMax = 0
		c.Cycles = 0
	} else if c.lastCycle > 0 {
		if elapsed < c.CtMin {
			c.CtMin = elapsed
		}
		if elapsed > c.CtMax {
			c.CtMax = elapsed
		}
		c.CtLast = elapsed
	}

	if c.init {
		c.SysTime += elapsed
		if c.SysTime >= 24*time.Hour {
			c.SysTime -= 24 * time.Hour
			c.SysDays++
		}
	}
	c.init = true
	c.lastCycle = tx
	c.Cycles++
}

// DT_SIMU simulates a real-time clock with adjustable speed.
type DT_SIMU struct {
	Dts time.Time

	// internal state
	init bool
	last int64
}

// Update executes the simulation logic.
func (d *DT_SIMU) Update(start time.Time, speed float64) {
	tx := logic.T_PLC_US()
	if !d.init {
		d.init = true
		d.Dts = start
		d.last = tx
	}

	elapsedMicroseconds := tx - d.last
	d.last = tx

	if speed == 0.0 {
		// In ST, this was an increment by 1 tick. In Go, we can just hold the time.
	} else {
		scaledDuration := time.Duration(float64(elapsedMicroseconds) * speed)
		d.Dts = d.Dts.Add(scaledDuration * time.Microsecond)
	}
}

// METER_STAT calculates statistics for a metered value (daily, weekly, monthly, yearly).
type METER_STAT struct {
	LastDay      float64
	CurrentDay   float64
	LastWeek     float64
	CurrentWeek  float64
	LastMonth    float64
	CurrentMonth float64
	LastYear     float64
	CurrentYear  float64

	// RETAIN fields
	yearStart  float64
	monthStart float64
	weekStart  float64
	dayStart   float64
	lastRun    time.Time
}

// Update executes the statistics logic.
func (m *METER_STAT) Update(in float64, di time.Time, rst bool) {
	if rst {
		m.LastDay, m.CurrentDay = 0.0, 0.0
		m.dayStart = in
		m.LastWeek, m.CurrentWeek = 0.0, 0.0
		m.weekStart = in
		m.LastMonth, m.CurrentMonth = 0.0, 0.0
		m.monthStart = in
		m.LastYear, m.CurrentYear = 0.0, 0.0
		m.yearStart = in
	} else {
		m.CurrentDay = in - m.dayStart
		m.CurrentWeek = in - m.weekStart
		m.CurrentMonth = in - m.monthStart
		m.CurrentYear = in - m.yearStart
	}

	if !m.lastRun.IsZero() {
		if di.Year() > m.lastRun.Year() {
			m.LastYear = m.CurrentYear
			m.yearStart = in
			m.LastMonth = m.CurrentMonth
			m.monthStart = in
			m.LastWeek = m.CurrentWeek // Week can also span year-end
			m.weekStart = in
			m.LastDay = m.CurrentDay
			m.dayStart = in
		} else if di.Month() > m.lastRun.Month() {
			m.LastMonth = m.CurrentMonth
			m.monthStart = in
			m.LastDay = m.CurrentDay
			m.dayStart = in
		} else if di.YearDay() > m.lastRun.YearDay() {
			m.LastDay = m.CurrentDay
			m.dayStart = in
		}

		_, w_di := di.ISOWeek()
		w_last := time_date.WORK_WEEK(m.lastRun)
		if w_di != w_last {
			m.LastWeek = m.CurrentWeek
			m.weekStart = in
		}
	}

	m.lastRun = di
}

// ONTIME measures the total on-time and number of cycles for a boolean signal.
type ONTIME struct {
	Seconds uint32
	Cycles  uint32

	// internal state
	last int64
	edge bool
	init bool
	ms   int64
}

// Update executes the on-time measurement logic.
func (o *ONTIME) Update(in, rst bool) {
	tx := logic.T_PLC_US()

	if !o.init {
		o.init = true
		o.last = tx
	}

	if rst {
		o.Seconds = 0
		o.Cycles = 0
		o.ms = 0
	} else if in {
		o.ms += (tx - o.last) / 1000 // Add elapsed milliseconds
		if o.ms >= 1000 {
			o.Seconds += uint32(o.ms / 1000)
			o.ms %= 1000
		}
		if !o.edge {
			o.Cycles++
		}
	}

	o.last = tx
	o.edge = in
}

// FLOW_METER measures flow according to gated time or pulses.
type FLOW_METER struct {
	F float64
	X float64 // VAR_IN_OUT
	Y uint64  // VAR_IN_OUT

	// internal state
	tl    time.Time
	int1  Integrate
	init  bool
	eLast bool
	xLast float64
	yLast uint64
}

// Update executes the flow meter logic.
func (fm *FLOW_METER) Update(vx float64, e, rst, pulseMode bool, updateTime time.Duration) {
	tx := time.Now()

	if !fm.init {
		fm.init = true
		fm.tl = tx
		fm.xLast = fm.X
		fm.yLast = fm.Y
	}

	// Gated operation
	fm.int1.Update(vx, 2.7777777777777777e-4, !rst && !pulseMode && e, &fm.X)

	if rst {
		fm.X = 0.0
		fm.Y = 0
		fm.tl = tx
		fm.xLast = 0.0
		fm.yLast = 0
	} else if e && pulseMode {
		if !fm.eLast {
			fm.X += vx
		}
	}
	fm.eLast = e

	// Reduce X to be less than 1 and increase Y
	if fm.X > 1.0 {
		tmp := math.Floor(fm.X)
		fm.Y += uint64(tmp)
		fm.X -= tmp
	}

	// Calculate current flow
	if tx.Sub(fm.tl) >= updateTime && updateTime > 0 {
		fm.F = (float64(fm.Y-fm.yLast) + fm.X - fm.xLast) / tx.Sub(fm.tl).Seconds() * 3600.0
		fm.yLast = fm.Y
		fm.xLast = fm.X
		fm.tl = tx
	}
}

// M_D measures the time between a rising edge on Start and a rising edge on Stop.
type M_D struct {
	PT  time.Duration
	ET  time.Duration
	Run bool

	// internal state
	edge    bool
	t0      time.Time
	startup bool
}

// Update executes the measurement logic.
func (m *M_D) Update(start, stop, rst bool, tmax time.Duration) {
	if rst || m.ET >= tmax {
		m.PT = 0
		m.ET = 0
		m.startup = false
		m.Run = false
	}

	if !m.startup {
		m.edge = start
		m.startup = true
	}

	tx := time.Now()

	if start && !m.edge && !stop {
		m.t0 = tx
		m.Run = true
		m.PT = 0
	} else if stop && m.Run {
		m.PT = m.ET
		m.Run = false
	}
	m.edge = start

	if m.Run {
		m.ET = tx.Sub(m.t0)
	}
}

// M_T measures the width of a high pulse.
type M_T struct {
	PT time.Duration
	ET time.Duration

	// internal state
	edge  bool
	start time.Time
}

// Update executes the measurement logic.
func (m *M_T) Update(in, rst bool, tmax time.Duration) {
	tx := time.Now()

	if rst || m.ET >= tmax {
		m.PT = 0
		m.ET = 0
	} else if in {
		if !m.edge {
			m.start = tx
		}
		m.ET = tx.Sub(m.start)
	} else {
		m.PT = m.ET
	}
	m.edge = in
}

// M_TX measures the timing of a signal (High time, Low time, Duty Cycle, Frequency).
type M_TX struct {
	TH time.Duration
	TL time.Duration
	DC float64
	F  float64
	ET time.Duration

	// internal state
	edge    bool
	start   time.Time
	stop    time.Time
	rise    bool
	fall    bool
	startup bool
}

// Update executes the measurement logic.
func (m *M_TX) Update(in, rst bool, tmax time.Duration) {
	if rst || (m.ET >= tmax) {
		m.rise, m.fall, m.startup = false, false, false
		m.TH, m.TL, m.ET = 0, 0, 0
		m.DC, m.F = 0.0, 0.0
	}

	if !m.startup {
		m.edge = in
		m.startup = true
	}

	tx := time.Now()

	if in != m.edge {
		m.edge = in
		if in { // Rising edge
			m.start = tx
			m.rise = true
			if m.fall {
				m.TL = m.start.Sub(m.stop)
			}
		} else { // Falling edge
			m.stop = tx
			m.fall = true
			if m.rise {
				m.TH = m.stop.Sub(m.start)
			}
		}

		if m.TH > 0 && m.TL > 0 {
			total := float64(m.TH + m.TL)
			m.DC = float64(m.TH) / total
			m.F = 1.0 / (total / float64(time.Second))
		}
	}

	if m.rise {
		m.ET = tx.Sub(m.start)
	}
}

// METER measures usage of power or similar values over time.
type METER struct {
	Mx float64 // VAR_IN_OUT

	// internal state
	mr   basic.REAL2
	last int64
	init bool
}

// Update executes the meter logic.
func (m *METER) Update(m1, m2, d float64, i1, i2, rst bool) {
	tx := logic.T_PLC_US()
	if !m.init {
		m.init = true
		m.last = tx
		m.mr.Rx = float32(m.Mx)
		m.mr.R1 = 0.0
	}

	if tx == m.last {
		return
	}

	tc := float64(tx-m.last) * 0.001 // Cycle time in seconds
	m.last = tx

	if rst {
		m.mr = basic.REAL2{}
	} else {
		var mx1, mx2 float64
		if i1 {
			mx1 = m1
		}
		if i2 {
			mx2 = m2
		}
		if d != 0.0 {
			m.mr = beeMath.R2_ADD(m.mr, float32(((mx1+mx2)/d)*tc))
		}
	}
	m.Mx = float64(m.mr.Rx)
}

// TC_MS measures the cycle time in milliseconds.
type TC_MS struct {
	TC uint32
	// internal state
	init bool
	last int64
}

// Update executes the logic.
func (t *TC_MS) Update() {
	tx := logic.T_PLC_US() / 1000 // to milliseconds
	if !t.init {
		t.init = true
		t.TC = 0
	} else {
		t.TC = uint32(tx - t.last)
	}
	t.last = tx
}

// TC_US measures the cycle time in microseconds.
type TC_US struct {
	TC uint32
	// internal state
	init bool
	last int64
}

// Update executes the logic.
func (t *TC_US) Update() {
	tx := logic.T_PLC_US()
	if !t.init {
		t.init = true
		t.TC = 0
	} else {
		t.TC = uint32(tx - t.last)
	}
	t.last = tx
}
