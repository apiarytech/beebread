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

	"beebread/basic/logic"
	beeMath "beebread/basic/math"
)

// DRIVER_1 is a multi-purpose driver.
// A rising edge on IN sets the output high if ToggleMode is false.
// If ToggleMode is true, a rising edge on IN toggles the output Q.
// If a Timeout is specified, the output Q will be reset to false automatically after the timeout has elapsed.
// An asynchronous reset and set will force the output high or low respectively.
type DRIVER_1 struct {
	ToggleMode bool
	Timeout    time.Duration

	// internal state
	q        bool
	edge     bool
	offTimer logic.TON
}

// Update executes the driver logic for one cycle.
func (d *DRIVER_1) Update(set, in, rst bool) bool {
	if d.offTimer.Q {
		d.q = false
	}

	if rst {
		d.q = false
	} else if set {
		d.q = true
	} else if in && !d.edge {
		if d.ToggleMode {
			d.q = !d.q
		} else {
			d.q = true
		}
	}
	d.edge = in

	if d.Timeout > 0 {
		d.offTimer.Update(d.q, d.Timeout)
	}

	return d.q
}

// DRIVER_4 is a 4-channel multi-purpose driver.
type DRIVER_4 struct {
	d0, d1, d2, d3 DRIVER_1
}

// Update executes the logic for all four drivers.
func (d *DRIVER_4) Update(set, rst bool, in [4]bool, toggleMode bool, timeout time.Duration) [4]bool {
	d.d0.ToggleMode = toggleMode
	d.d0.Timeout = timeout
	d.d1.ToggleMode = toggleMode
	d.d1.Timeout = timeout
	d.d2.ToggleMode = toggleMode
	d.d2.Timeout = timeout
	d.d3.ToggleMode = toggleMode
	d.d3.Timeout = timeout

	var q [4]bool
	q[0] = d.d0.Update(set, in[0], rst)
	q[1] = d.d1.Update(set, in[1], rst)
	q[2] = d.d2.Update(set, in[2], rst)
	q[3] = d.d3.Update(set, in[3], rst)

	return q
}

// DRIVER_4C is a multi-purpose cycling driver.
// A rising edge on IN switches from one state to the next.
// The state of the outputs in any state is configurable with SX.
type DRIVER_4C struct {
	Timeout time.Duration
	SX      [7]byte

	// internal state
	SN       int
	edge     bool
	offTimer logic.TON
}

// Update executes the driver logic for one cycle.
func (d *DRIVER_4C) Update(in, rst bool) (int, [4]bool) {
	var q [4]bool

	d.offTimer.Update(d.SN > 0, d.Timeout)

	if rst || d.offTimer.Q {
		d.SN = 0
	} else if in && !d.edge {
		d.SN++
		if d.SN > 7 || (d.SN > 0 && d.SN <= 7 && d.SX[d.SN-1] == 0) {
			d.SN = 0
		}
	}
	d.edge = in

	if d.SN > 0 && d.SN <= 7 {
		sxVal := d.SX[d.SN-1]
		q[0] = (sxVal & 0x01) != 0
		q[1] = (sxVal & 0x02) != 0
		q[2] = (sxVal & 0x04) != 0
		q[3] = (sxVal & 0x08) != 0
	}

	return d.SN, q
}

// FLOW_CONTROL switches a valve depending on the input IN.
// It also limits the maximum on-time of the valve and controls pressure on the output side.
type FLOW_CONTROL struct {
	TAuto  time.Duration
	TDelay time.Duration

	// internal state
	q      bool
	status byte
	timer  logic.TP1D
}

// Update executes the flow control logic for one cycle.
func (fc *FLOW_CONTROL) Update(in, req, enq, rst bool) (bool, byte) {
	fc.status = 100
	if rst {
		fc.q = false
		fc.timer.RST = true
		fc.status = 103
	} else if enq {
		if in {
			fc.status = 101
		}
		if req {
			fc.timer.PT1 = fc.TAuto
			fc.timer.PTD = fc.TDelay
			fc.timer.IN = true
			fc.status = 102
		}
	}

	fc.timer.Update()
	fc.timer.IN = false // Pulse behavior
	fc.timer.RST = false

	fc.q = (in && enq) || fc.timer.Q
	return fc.q, fc.status
}

// FT_PROFILE generates an output signal defined by values over a time scale.
type FT_PROFILE struct {
	// Configuration
	Value0  float32
	Time1   time.Duration
	Value1  float32
	Time2   time.Duration
	Value2  float32
	Time3   time.Duration
	Value3  float32
	Time10  time.Duration
	Value10 float32
	Time11  time.Duration
	Value11 float32
	Time12  time.Duration
	Value12 float32
	Time13  time.Duration
	Value13 float32

	// Internal State
	Y     float32
	Run   bool
	ET    time.Duration
	edge  bool
	state byte
	ta    time.Time
	tb    time.Duration
	t0    time.Time
	temp  float32
	va    float32
	vb    float32
}

// Update executes the profile generation logic.
func (p *FT_PROFILE) Update(k, o, m float32, e bool) {
	tx := time.Now()

	if e && !p.edge {
		p.Run = true
		p.ET = 0
		p.t0 = tx
		p.ta = tx
		p.tb = time.Duration(float64(p.Time1) * float64(m))
		p.va = p.Value0
		p.vb = p.Value1
		p.temp = p.Value0
		p.state = 1
	}
	p.edge = e

	if p.Run {
		updateState := func(nextState byte, nextTime, prevTime time.Duration, nextVal, prevVal float32) {
			p.ta = p.ta.Add(p.tb)
			p.tb = time.Duration(float64(nextTime-prevTime) * float64(m))
			p.va = prevVal
			p.vb = nextVal
			p.temp = prevVal
			p.state = nextState
		}

		interpolate := func() {
			if p.tb > 0 {
				p.temp = (p.vb-p.va)*float32(tx.Sub(p.ta))/float32(p.tb) + p.va
			}
		}

		switch p.state {
		case 1:
			if tx.Sub(p.ta) >= p.tb {
				updateState(2, p.Time2, p.Time1, p.Value2, p.Value1)
			} else {
				interpolate()
			}
		case 2:
			if tx.Sub(p.ta) >= p.tb {
				updateState(3, p.Time3, p.Time2, p.Value3, p.Value2)
			} else {
				interpolate()
			}
		case 3:
			if tx.Sub(p.ta) >= p.tb {
				updateState(4, p.Time10, p.Time3, p.Value10, p.Value3)
			} else {
				interpolate()
			}
		case 4:
			if tx.Sub(p.ta) >= p.tb {
				updateState(5, p.Time11, p.Time10, p.Value11, p.Value10)
				if !e {
					p.state = 6
				}
			} else {
				interpolate()
			}
		case 5: // Extend while E is true
			if e {
				p.ta = tx
			} else {
				p.state = 6
			}
		case 6:
			if tx.Sub(p.ta) >= p.tb {
				updateState(7, p.Time12, p.Time11, p.Value12, p.Value11)
			} else {
				interpolate()
			}
		case 7:
			if tx.Sub(p.ta) >= p.tb {
				updateState(8, p.Time13, p.Time12, p.Value13, p.Value12)
			} else {
				interpolate()
			}
		case 8:
			if tx.Sub(p.ta) >= p.tb {
				p.temp = p.Value13
				p.Run = false
			} else {
				interpolate()
			}
		}
		p.Y = p.temp*k + o
		p.ET = tx.Sub(p.t0)
	}
}

// INC_DEC is an incremental decoder with quadruple accuracy.
type INC_DEC struct {
	Dir bool
	Cnt int
	// internal state
	edgeA, edgeB bool
}

// Update executes the decoder logic.
func (id *INC_DEC) Update(cha, chb, rst bool) {
	axb := cha != chb
	clka := cha != id.edgeA
	id.edgeA = cha
	clkb := chb != id.edgeB
	id.edgeB = chb
	clk := clka || clkb

	if axb && clka {
		id.Dir = true
	}
	if axb && clkb {
		id.Dir = false
	}

	if clk {
		if id.Dir {
			id.Cnt++
		} else {
			id.Cnt--
		}
	}

	if rst {
		id.Cnt = 0
	}
}

// INTERLOCK has two inputs I1 and I2 which drive the corresponding outputs Q1 and Q2.
// The input signals lock each other out.
type INTERLOCK struct {
	t1, t2 logic.TOF
}

// Update executes the interlock logic.
func (il *INTERLOCK) Update(i1, i2 bool, tl time.Duration) (bool, bool) {
	q1 := i1 && !il.t2.Update(i2, tl)
	q2 := i2 && !il.t1.Update(i1, tl)
	return q1, q2
}

// Interlock4 detects one of 4 switches and delivers the number of the switch pressed.
type INTERLOCK_4 struct {
	Out byte
	TP  bool
	// internal state
	last, old, in byte
	lmode         int
}

// Update executes the interlock logic.
func (il *INTERLOCK_4) Update(i0, i1, i2, i3, e bool, mode int) {
	if e {
		if mode != il.lmode {
			il.Out, il.last, il.old, il.lmode = 0, 0, 0, mode
		}

		il.in = 0
		if i0 {
			il.in |= 1
		}
		if i1 {
			il.in |= 2
		}
		if i2 {
			il.in |= 4
		}
		if i3 {
			il.in |= 8
		}

		if il.in != il.last {
			switch mode {
			case 0:
				il.Out = il.in
			case 1:
				if (il.in & 8) != 0 {
					il.Out = 8
				} else if (il.in & 4) != 0 {
					il.Out = 4
				} else if (il.in & 2) != 0 {
					il.Out = 2
				} else {
					il.Out = il.in
				}
			case 2:
				il.last = (il.in ^ il.last) & il.in
				if (il.last & 8) != 0 {
					il.Out = 8
				} else if (il.last & 4) != 0 {
					il.Out = 4
				} else if (il.last & 2) != 0 {
					il.Out = 2
				} else {
					il.Out = il.last
				}
			case 3:
				if (il.Out & il.in) == 0 {
					if (il.in & 8) != 0 {
						il.Out = 8
					} else if (il.in & 4) != 0 {
						il.Out = 4
					} else if (il.in & 2) != 0 {
						il.Out = 2
					} else {
						il.Out = il.in
					}
				}
			}
			il.last = il.in
		}
		il.TP = il.Out != il.old
		il.old = il.Out
	} else {
		il.Out, il.last, il.old, il.lmode, il.TP = 0, 0, 0, 0, false
	}
}

// MANUAL is a manual override for digital signals.
func MANUAL(in, on, off bool) bool {
	return !off && (in || on)
}

// MANUAL_1 is a manual override for digital signals.
type MANUAL_1 struct {
	Q      bool
	Status byte
	// internal state
	sEdge, rEdge, edge bool
}

// Update executes the manual override logic.
func (m *MANUAL_1) Update(in, man, mI, set, rst bool) {
	if !man {
		m.Q = in
		m.Status = 100
		m.edge = false
	} else if !m.sEdge && set {
		m.Q = true
		m.edge = true
		m.Status = 101
	} else if !m.rEdge && rst {
		m.Q = false
		m.edge = true
		m.Status = 102
	} else if !m.edge {
		m.Q = mI
		m.Status = 103
	}
	m.sEdge = set
	m.rEdge = rst
}

// MANUAL_2 is a manual override for boolean signals.
func MANUAL_2(in, ena, on, off, man bool) (bool, byte) {
	if ena {
		if !on && !off {
			return in, 100
		} else if on && !off {
			return true, 101
		} else if !on && off {
			return false, 102
		} else {
			return man, 103
		}
	}
	return false, 104
}

// MANUAL_4 is a manual override for 4 digital signals.
type MANUAL_4 struct {
	Q      [4]bool
	Status byte
	// internal state
	edge bool
	pos  int
	tog  bool
}

// Update executes the logic.
func (m *MANUAL_4) Update(i [4]bool, man, stp bool, mIn [4]bool) {
	if man {
		if !m.tog {
			m.Q = mIn
			m.Status = 101
		}
		if stp && !m.edge {
			m.tog = true
			m.Q = [4]bool{} // Reset all
			m.Q[m.pos] = true
			m.Status = 110 + byte(m.pos)
			m.pos = (m.pos + 1) % 4
		}
	} else {
		m.Q = i
		m.Status = 100
		m.tog = false
		m.pos = 0
	}
	m.edge = stp
}

// PARSET selects one of 4 parameter sets addressed by A0 and A1.
// If TC is specified, the change of the outputs is ramped by the time TC.
type PARSET struct {
	P1, P2, P3, P4 float64

	// internal state
	x     [4][4]float64
	s     [4]float64
	last  time.Time
	start bool
	set   byte
	init  bool
}

// Update executes the parameter set selection logic.
func (p *PARSET) Update(a0, a1 bool, tc time.Duration, params [4][4]float64) {
	tx := time.Now()

	if !p.init {
		p.init = true
		p.x = params
		p.set = 0
		if a0 {
			p.set |= 1
		}
		if a1 {
			p.set |= 2
		}
		p.P1 = p.x[p.set][0]
		p.P2 = p.x[p.set][1]
		p.P3 = p.x[p.set][2]
		p.P4 = p.x[p.set][3]
	}

	newSet := byte(0)
	if a0 {
		newSet |= 1
	}
	if a1 {
		newSet |= 2
	}

	if newSet != p.set {
		p.set = newSet
		if tc > 0 {
			p.start = true
			p.last = tx
			tcSec := tc.Seconds()
			p.s[0] = (p.x[p.set][0] - p.P1) / tcSec
			p.s[1] = (p.x[p.set][1] - p.P2) / tcSec
			p.s[2] = (p.x[p.set][2] - p.P3) / tcSec
			p.s[3] = (p.x[p.set][3] - p.P4) / tcSec
		}
	} else if p.start && time.Since(p.last) < tc {
		elapsed := time.Since(p.last).Seconds()
		p.P1 += p.s[0] * elapsed
		p.P2 += p.s[1] * elapsed
		p.P3 += p.s[2] * elapsed
		p.P4 += p.s[3] * elapsed
	} else {
		p.start = false
		p.P1 = p.x[p.set][0]
		p.P2 = p.x[p.set][1]
		p.P3 = p.x[p.set][2]
		p.P4 = p.x[p.set][3]
	}
	p.last = tx
}

// PARSET2 selects one of 4 parameter sets depending on the value of X.
type PARSET2 struct {
	P1, P2, P3, P4 float64
	pset           PARSET
}

// Update executes the logic.
func (p *PARSET2) Update(x, l1, l2, l3 float64, tc time.Duration, params [4][4]float64) {
	var a0, a1 bool
	absX := math.Abs(x)
	if absX < l1 {
		a0, a1 = false, false
	} else if absX < l2 {
		a0, a1 = true, false
	} else if absX < l3 {
		a0, a1 = false, true
	} else {
		a0, a1 = true, true
	}
	p.pset.Update(a0, a1, tc, params)
	p.P1 = p.pset.P1
	p.P2 = p.pset.P2
	p.P3 = p.pset.P3
	p.P4 = p.pset.P4
}

// SIGNAL generates an output signal according to a bit pattern SIG.
type SIGNAL struct {
	Q bool
}

// Update executes the signal generation logic.
func (s *SIGNAL) Update(in bool, sig byte, ts time.Duration) {
	if in {
		tx := time.Now().UnixMilli()
		var step byte
		if ts > 0 {
			step = byte(tx/int64(ts.Milliseconds())) & 0x07
		} else {
			step = byte(tx>>7) & 0x07
		}
		step = 1 << step
		s.Q = (step & sig) > 0
	} else {
		s.Q = false
	}
}

// SIGNAL_4 generates one out of 4 signals specified by bit patterns S1..S4.
type SIGNAL_4 struct {
	Q   bool
	sig SIGNAL
}

// Update executes the logic.
func (s *SIGNAL_4) Update(in1, in2, in3, in4 bool, ts time.Duration, s1, s2, s3, s4 byte) {
	var sigIn bool
	var sigPattern byte

	if in1 {
		sigIn = true
		sigPattern = s1
	} else if in2 {
		sigIn = true
		sigPattern = s2
	} else if in3 {
		sigIn = true
		sigPattern = s3
	} else if in4 {
		sigIn = true
		sigPattern = s4
	} else {
		sigIn = false
	}

	s.sig.Update(sigIn, sigPattern, ts)
	s.Q = s.sig.Q
}

// SRAMP generates an output signal that is slew rate and acceleration controlled.
type SRAMP struct {
	Y float64
	V float64

	// internal state
	cycleTime logic.TCS
	init      bool
}

// Update executes the S-Ramp logic.
func (s *SRAMP) Update(x, aUp, aDn, vuMax, vdMax, limitHigh, limitLow float64, rst bool) {
	s.cycleTime.Update()
	tc := s.cycleTime.TC

	aUp = math.Max(0.0, aUp)
	aDn = math.Min(0.0, aDn)
	vuMax = math.Max(0.0, vuMax)
	vdMax = math.Min(0.0, vdMax)

	if rst || !s.init {
		s.init = true
		s.Y = 0.0
		s.V = 0.0
	} else if x == s.Y {
		s.V = 0.0
	} else if x > s.Y { // Ramp up
		s.V = math.Min(s.V+aUp*tc, vuMax)
		s.V = math.Min(math.Sqrt((s.Y-x)*2.0*aDn), s.V)
		s.Y = beeMath.Limit(limitLow, s.Y+math.Min(s.V*tc, x-s.Y), limitHigh)
	} else { // Ramp down
		s.V = math.Max(s.V+aDn*tc, vdMax)
		s.V = math.Max(-math.Sqrt((s.Y-x)*2.0*aUp), s.V)
		s.Y = beeMath.Limit(limitLow, s.Y+math.Max(s.V*tc, x-s.Y), limitHigh)
	}
}

// TUNE generates an output signal which is set by input switches.
type TUNE struct {
	Y float64

	// internal state
	start, start2 time.Time
	state         int
	step          float64
	speed         float64
	yStart        float64
	yStart2       float64
}

// Update executes the tuning logic.
func (t *TUNE) Update(set, su, sd, rst bool, ss, limitL, limitH, rstVal, setVal, s1, s2 float64, t1, t2 time.Duration) {
	tx := time.Now()

	if rst {
		t.Y = rstVal
		t.state = 0
	} else if set {
		t.Y = setVal
		t.state = 0
	} else if t.state > 0 {
		in := (t.state == 1 && su) || (t.state == 2 && sd)

		if !in && tx.Sub(t.start) <= t1 {
			t.Y = t.yStart + t.step
			t.state = 0
		} else if in && tx.Sub(t.start) >= t2 {
			t.Y = t.yStart2 + tx.Sub(t.start2).Seconds()*s2/t.speed
		} else if in && tx.Sub(t.start) >= t1 {
			t.Y = t.yStart + (tx.Sub(t.start)-t1).Seconds()*s1/t.speed
			t.start2 = tx
			t.yStart2 = t.Y
		} else if !in {
			t.state = 0
		}
	} else if su {
		t.state = 1
		t.start = tx
		t.step = ss
		t.speed = 1.0
		t.yStart = t.Y
	} else if sd {
		t.state = 2
		t.start = tx
		t.step = -ss
		t.speed = -1.0
		t.yStart = t.Y
	}

	t.Y = beeMath.Limit(limitL, t.Y, limitH)
}

// PWM_DC generates a square wave signal specified by frequency and duty cycle.
type PWM_DC struct {
	Q bool

	// internal state
	clk   CLK_PRG
	pulse TP_X
}

// Update executes the PWM logic.
func (p *PWM_DC) Update(f, dc float64) {
	if f > 0.0 {
		tmp := 1000.0 / f
		p.clk.Update(time.Duration(tmp) * time.Millisecond)
		p.pulse.Update(p.clk.Q, time.Duration(tmp*dc*float64(time.Millisecond)))
		p.Q = p.pulse.Q
	} else {
		p.Q = false
	}
}

// CLK_PRG uses the system time to generate a clock with a programmable period time.
// A pulse is generated for one cycle only.
type CLK_PRG struct {
	Q bool

	// internal state
	init bool
	last time.Time
}

// Update executes the clock logic.
func (c *CLK_PRG) Update(pt time.Duration) {
	tx := time.Now()

	if !c.init {
		c.init = true
		c.last = tx.Add(-pt) // Ensure first pulse is generated
	}

	if tx.Sub(c.last) >= pt {
		c.Q = true
		c.last = tx
	} else {
		c.Q = false
	}
}

// TP_X is a retriggerable, edge-triggered pulse timer.
type TP_X struct {
	Q  bool
	ET time.Duration

	// internal state
	edge      bool
	startTime time.Time
}

// Update executes the pulse timer logic.
func (t *TP_X) Update(in bool, pt time.Duration) {
	tx := time.Now()

	// Rising edge trigger
	if in && !t.edge {
		t.startTime = tx
		t.Q = pt > 0
	} else if t.Q {
		t.ET = tx.Sub(t.startTime)
		if t.ET >= pt {
			t.Q = false
			t.ET = 0
		}
	}
	t.edge = in
}

// PWM_PW generates a square wave signal specified by frequency and pulse width.
type PWM_PW struct {
	Q bool

	// internal state
	clk   CLK_PRG
	pulse TP_X
}

// Update executes the PWM logic.
func (p *PWM_PW) Update(f float64, pw time.Duration) {
	if f > 0.0 {
		p.clk.Update(time.Duration(1000.0/f) * time.Millisecond)
		p.pulse.Update(p.clk.Q, pw)
		p.Q = p.pulse.Q
	} else {
		p.Q = false
	}
}
