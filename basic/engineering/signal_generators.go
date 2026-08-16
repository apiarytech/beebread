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
	"beebread/basic"
	"beebread/basic/other"
	"math"
	"time"

	"beebread/basic/logic"
	beeMath "beebread/basic/math"
)

// RMP_B generates a ramp on an external var of type byte.
type RMP_B struct {
	// internal state
	init     bool
	tl       time.Time
	tn       time.Duration
	lastDir  bool
	startVal byte
}

// Update executes the ramp logic. Rmp is a pointer to the value being ramped.
func (r *RMP_B) Update(rmp *byte, e, dir bool, tr time.Duration) {
	tx := time.Now()

	var sel byte
	if dir {
		sel = 255
	}

	if e && r.init && (dir == r.lastDir) && (*rmp != sel) && tr == r.tn {
		*rmp = beeMath.FRMP_B(r.startVal, dir, tx.Sub(r.tl), tr)
	} else {
		r.init = true
		r.tl = tx
		r.tn = tr
		r.startVal = *rmp
	}
	r.lastDir = dir
}

// RMP_NEXT generates a ramp output following the input IN.
type RMP_NEXT struct {
	Dir bool // upwards = TRUE
	Up  bool
	Dn  bool

	// internal state
	rmx   RMP_B
	dirx  TREND_DW
	tLock logic.TP
	xen   bool
	xdir  bool
}

// Update executes the logic. Out is a pointer to the value being ramped.
func (r *RMP_NEXT) Update(out *byte, e bool, in byte, tr, tf, tl time.Duration) {
	r.dirx.Update(uint32(in))
	r.tLock.Update(false, tl)

	if r.dirx.TU && (*out < in) {
		if !r.xdir && r.xen {
			r.tLock.Update(true, tl)
		}
		r.xen = true
		r.xdir = true
	} else if r.dirx.TD && (*out > in) {
		if r.xdir && r.xen {
			r.tLock.Update(true, tl)
		}
		r.xen = true
		r.xdir = false
	} else if r.xen {
		if (r.xdir && (*out >= in)) || (!r.xdir && (*out <= in)) {
			r.xen = false
			if tl > 0 {
				r.tLock.Update(true, tl)
			}
		}
	}

	if !r.tLock.Q && r.xen {
		r.Up = r.xdir
		r.Dir = r.xdir
		r.Dn = !r.xdir
	} else {
		r.Up = false
		r.Dn = false
	}

	var rampTime time.Duration
	if r.Dir {
		rampTime = tr
	} else {
		rampTime = tf
	}
	r.rmx.Update(out, e && (r.Up || r.Dn), r.Dir, rampTime)
}

// RMP_W generates a word-wide ramp with set/reset and direction control.
type RMP_W struct {
	Out  uint16
	Busy bool
	High bool
	Low  bool

	// internal state
	rmp _RmpW
}

// Update executes the ramp logic.
func (r *RMP_W) Update(set, rst, e, up bool, pt time.Duration) {
	r.rmp.Update(&r.Out, e, up, pt)
	if rst {
		r.Out = 0
	} else if set {
		r.Out = 65535
	}
	r.Low = r.Out == 0
	r.High = r.Out == 65535
	r.Busy = !(r.Low || r.High) && e
}

// _RmpW is the internal implementation for a word ramp.
type _RmpW struct {
	// internal state
	init    bool
	tl      int64
	lastDir bool
}

// Update executes the ramp logic. Rmp is a pointer to the value being ramped.
func (r *_RmpW) Update(rmp *uint16, e, dir bool, tr time.Duration) {
	tx := logic.T_PLC_US() / 1000 // T_PLC_MS

	if e && r.init {
		if dir != r.lastDir {
			r.tl = tx
			r.lastDir = dir
		}

		var step int64
		if tr > 0 {
			step = ((tx - r.tl) << 16) / tr.Milliseconds()
		} else {
			step = 65535
		}

		if step > 0 {
			r.tl = tx
			if !dir {
				step = -step
			}
			newVal := beeMath.LIMIT_I64(0, int64(*rmp)+step, 65535)
			*rmp = uint16(newVal)
		}
	} else {
		r.tl = tx
		r.init = true
	}
}

// FADE cross-fades between two inputs.
type FADE struct {
	Y   float64
	rmx RMP_W
}

// Update executes the fade logic.
func (f *FADE) Update(in1, in2 float64, fade bool, tf time.Duration, rst bool) {
	f.rmx.Update(rst && fade, rst && !fade, true, fade, tf)
	f.Y = (in2-in1)/65535.0*float64(f.rmx.Out) + in1
}

// GEN_PULSE generates a continuous output waveform with programmable high and low times.
type GEN_PULSE struct {
	Q bool

	// internal state
	init bool
	tn   time.Time
}

// Update executes the pulse generator logic.
func (g *GEN_PULSE) Update(enq bool, pth, ptl time.Duration) {
	if enq {
		tx := time.Now()
		if !g.init {
			g.init = true
			g.tn = tx
		}

		var period time.Duration
		if g.Q {
			period = pth
		} else {
			period = ptl
		}

		if tx.Sub(g.tn) >= period {
			g.tn = g.tn.Add(period)
			g.Q = !g.Q
		}
	} else {
		g.Q = false
		g.init = false
	}
}

// GEN_PW2 generates a pulse with two selectable timings.
type GEN_PW2 struct {
	Q  bool
	TH time.Duration
	TL time.Duration

	// internal state
	init  bool
	start time.Time
}

// Update executes the logic.
func (g *GEN_PW2) Update(enq bool, th1, tl1, th2, tl2 time.Duration, ts bool) {
	tx := time.Now()

	if !g.init {
		g.start = tx
		g.init = true
		g.TH, g.TL = 0, 0
	}

	tHigh := th1
	tLow := tl1
	if ts {
		tHigh = th2
		tLow = tl2
	}

	if enq {
		et := tx.Sub(g.start)
		if !g.Q {
			if et >= tLow {
				g.Q = true
				g.start = tx
				g.TL = 0
			} else {
				g.TL = et
			}
		} else {
			if et >= tHigh {
				g.Q = false
				g.start = tx
				g.TH = 0
			} else {
				g.TH = et
			}
		}
	} else {
		g.Q = false
		g.TH, g.TL = 0, 0
		g.start = tx
	}
}

// GEN_RDM generates a random signal at periodic intervals.
type GEN_RDM struct {
	Q   bool
	Out float64

	// internal state
	init bool
	last time.Time
}

// Update executes the random signal generation logic.
func (g *GEN_RDM) Update(pt time.Duration, am, os float64) {
	tx := time.Now()

	if !g.init {
		g.init = true
		g.last = tx
	}

	if tx.Sub(g.last) >= pt {
		g.last = g.last.Add(pt)
		g.Out = am*(other.RDM(0)-0.5) + os
		g.Q = true
	} else {
		g.Q = false
	}
}

// GEN_RDT generates a defined pulse with a random period.
type GEN_RDT struct {
	XQ bool

	// internal state
	tonRDMTimer logic.TON
	tofXQ       logic.TOF
	tRDMTime    time.Duration
	rRDMTime    float64
}

// Update executes the logic.
func (g *GEN_RDT) Update(enable bool, minTime, maxTime, pulseTime time.Duration) {
	g.tonRDMTimer.Update(enable, g.tRDMTime)
	g.tofXQ.Update(g.tonRDMTimer.Q, pulseTime)
	g.XQ = g.tofXQ.Q

	if g.tonRDMTimer.Q {
		g.XQ = true
		g.rRDMTime = other.RDM(g.rRDMTime)
		randomRange := float64(maxTime - minTime)
		g.tRDMTime = time.Duration(g.rRDMTime*randomRange) + minTime
		g.tonRDMTimer.Update(false, g.tRDMTime) // Reset the timer
	}
}

// GEN_RMP generates a ramp wave output.
type GEN_RMP struct {
	Q   bool
	Out float64

	// internal state
	init  bool
	last  time.Time
	temp  float64
	ltemp float64
}

// Update executes the ramp generation logic.
func (g *GEN_RMP) Update(pt time.Duration, am, os, dl float64) {
	tx := time.Now()

	dl = beeMath.MODR(dl, 1.0)
	if dl < 0.0 {
		dl += 1.0
	}

	if !g.init {
		g.init = true
		g.last = tx
	}

	elapsed := tx.Sub(g.last)
	if elapsed >= pt {
		g.last = g.last.Add(pt)
		elapsed -= pt
	}

	g.ltemp = g.temp
	if pt > 0 {
		totalElapsed := float64(elapsed) + float64(pt)*dl
		g.temp = beeMath.FRACT(totalElapsed / float64(pt))
	}
	g.Out = am*g.temp + os
	g.Q = g.temp < g.ltemp
}

// TREND_DW analyses the trend of a dword input signal.
// Q is true if the input is rising or falling.
// TU is true for one cycle on a rising trend.
// TD is true for one cycle on a falling trend.
// D is the delta of the change.
type TREND_DW struct {
	Q  bool
	TU bool
	TD bool
	D  uint32
	// internal state
	lastX uint32
}

// Update executes the trend analysis logic.
func (t *TREND_DW) Update(x uint32) {
	if x > t.lastX {
		t.TU = true
		t.TD = false
		t.D = x - t.lastX
		t.Q = true
	} else if x < t.lastX {
		t.TD = true
		t.TU = false
		t.D = t.lastX - x
		t.Q = false
	} else {
		t.TU = false
		t.TD = false
		t.D = 0
	}
	t.lastX = x
}

// FRMP_B calculates a ramp for a byte value and limits the output to 0-255.
// It avoids overflow issues during calculation.
func FRMP_B(start byte, dir bool, td, tr time.Duration) byte {
	if td < tr && tr > 0 {
		val := byte((uint64(td) * 256) / uint64(tr))
		if dir { // Ramp up_
			return byte(beeMath.LIMIT(0, float64(start)+float64(val), 255))
		}
		// Ramp down
		return byte(beeMath.LIMIT(0, float64(start)-float64(val), 255))
	} else if dir {
		return 255
	}
	return 0
}

// GEN_SIN generates a sine wave output.
type GEN_SIN struct {
	Q   bool
	Out float64

	// internal state
	init bool
	last time.Time
	temp float64
}

// Update executes the sine wave generation logic.
func (g *GEN_SIN) Update(pt time.Duration, am, os, dl float64) {
	tx := time.Now()
	if !g.init {
		g.init = true
		g.last = tx
	}

	dl = beeMath.MODR(dl, 1.0)
	if dl < 0.0 {
		dl += 1.0
	}

	elapsed := tx.Sub(g.last)
	if pt > 0 && elapsed >= pt {
		g.last = g.last.Add(pt)
		elapsed -= pt
	}

	if pt > 0 {
		totalElapsed := float64(elapsed) + float64(pt)*dl
		g.temp = math.Sin(basic.Math.Pi2 * totalElapsed / float64(pt))
	}

	g.Out = am*0.5*g.temp + os
	g.Q = g.temp >= 0.0
}

// GEN_SQR generates a square wave output.
type GEN_SQR struct {
	Q   bool
	Out float64

	// internal state
	init bool
	last time.Time
}

// Update executes the square wave generation logic.
func (g *GEN_SQR) Update(pt time.Duration, am, os, dc, dl float64) {
	if dc <= 0.0 {
		g.Out = -am*0.5 + os
		g.Q = false
		return
	}
	if dc >= 1.0 {
		g.Out = am*0.5 + os
		g.Q = true
		return
	}

	tx := time.Now()
	if !g.init {
		g.init = true
		g.last = tx
	}

	dl = beeMath.MODR(dl, 1.0)
	if dl < 0.0 {
		dl += 1.0
	}

	elapsed := tx.Sub(g.last)
	if pt > 0 && elapsed >= pt {
		g.last = g.last.Add(pt)
		elapsed -= pt
	}

	var temp float64
	if pt > 0 {
		temp = beeMath.FRACT(float64(elapsed)/float64(pt) + dl)
	}

	if temp < dc {
		g.Q = true
		g.Out = am*0.5 + os
	} else {
		g.Q = false
		g.Out = -am*0.5 + os
	}
}

// RMP_SOFT generates a soft on/off ramp for a byte value.
type RMP_SOFT struct {
	Out byte

	// internal state
	rmp RMP_B
}

// Update executes the soft ramp logic.
func (r *RMP_SOFT) Update(in bool, val byte, ptOn, ptOff time.Duration) {
	var target byte
	if in {
		target = val
	}

	if target > r.Out {
		r.rmp.Update(&r.Out, true, true, ptOn)
		if r.Out > target {
			r.Out = target
		}
	} else if target < r.Out {
		r.rmp.Update(&r.Out, true, false, ptOff)
		if r.Out < target {
			r.Out = target
		}
	} else {
		r.rmp.Update(&r.Out, false, false, 0)
	}
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
