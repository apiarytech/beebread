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
	beeMath "beebread/basic/math"
)

// BAND_B limits a byte value X to a band B.
// If X < B, the result is 0. If X > 255-B, the result is 255. Otherwise, it's X.
func BAND_B(x, b byte) byte {
	if x < b {
		return 0
	}
	if x > 255-b {
		return 255
	}
	return x
}

// CONTROL_SET1 calculates controller parameters for P, PI, and PID controllers
// based on the Ziegler-Nichols method using Kt and Tt.
type CONTROL_SET1 struct {
	// Constants
	PK    float64
	PIK   float64
	PITN  float64
	PIDK  float64
	PIDTN float64
	PIDTV float64

	// Outputs
	KP float64
	TN float64
	TV float64
	KI float64
	KD float64
}

// NewCONTROL_SET1 creates a CONTROL_SET1 with default values.
func NewCONTROL_SET1() *CONTROL_SET1 {
	return &CONTROL_SET1{
		PK:    0.5,
		PIK:   0.45,
		PITN:  0.83,
		PIDK:  0.6,
		PIDTN: 0.5,
		PIDTV: 0.125,
	}
}

// Update calculates the controller parameters.
func (c *CONTROL_SET1) Update(kt, tt float64, pi, pid bool) {
	if pi && pid {
		c.KP, c.TN, c.TV = 0.0, 0.0, 0.0
	} else if pid {
		c.KP = c.PIDK * kt
		c.TN = c.PIDTN * tt
		c.TV = c.PIDTV * tt
	} else if pi {
		c.KP = c.PIK * kt
		c.TN = c.PITN * tt
	} else {
		c.KP = c.PK * kt
		c.TV = 0.0 // P controller has no derivative time
	}

	if c.TN > 0.0 {
		c.KI = c.KP / c.TN
	} else {
		c.KI = 0.0
	}
	c.KD = c.KP * c.TV
}

// CONTROL_SET2 calculates controller parameters for P, PI, and PID controllers
// based on the Ziegler-Nichols method using KS, TU, and TG.
type CONTROL_SET2 struct {
	// Constants
	PK    float64
	PIK   float64
	PITN  float64
	PIDK  float64
	PIDTN float64
	PIDTV float64

	// Outputs
	KP float64
	TN float64
	TV float64
	KI float64
	KD float64
}

// NewCONTROL_SET2 creates a CONTROL_SET2 with default values.
func NewCONTROL_SET2() *CONTROL_SET2 {
	return &CONTROL_SET2{
		PK:    1.0,
		PIK:   0.9,
		PITN:  3.33,
		PIDK:  1.2,
		PIDTN: 2.0,
		PIDTV: 0.5,
	}
}

// Update calculates the controller parameters.
func (c *CONTROL_SET2) Update(ks, tu, tg float64, pi, pid bool) {
	var tx float64
	if tu > 0.0 && ks > 0.0 {
		tx = tg / tu / ks
	}

	if pi && pid {
		c.KP, c.TN, c.TV = 0.0, 0.0, 0.0
	} else if pid {
		c.KP = c.PIDK * tx
		c.TN = c.PIDTN * tu
		c.TV = c.PIDTV * tu
	} else if pi {
		c.KP = c.PIK * tx
		c.TN = c.PITN * tu
	} else {
		c.KP = c.PK * tx
		c.TV = 0.0
	}

	if c.TN > 0.0 {
		c.KI = c.KP / c.TN
	} else {
		c.KI = 0.0
	}
	c.KD = c.KP * c.TV
}

// CTRL_IN calculates the process error (difference) with a dead zone for noise.
func CTRL_IN(setPoint, actual, noise float64) float64 {
	return DEAD_ZONE(setPoint-actual, noise)
}

// CTRL_OUT handles manual override and output limiting for a controller.
type CTRL_OUT struct {
	Y   float64
	Lim bool
}

// Update executes the logic.
func (c *CTRL_OUT) Update(ci, offset, manIn, limL, limH float64, manual bool) {
	if manual {
		c.Y = manIn + offset
	} else {
		c.Y = ci + offset
	}

	if c.Y > limL && c.Y < limH {
		c.Lim = false
	} else {
		c.Y = beeMath.LIMIT(limL, c.Y, limH)
		c.Lim = true
	}
}

// CTRL_PI is a PI controller with manual functionality.
type CTRL_PI struct {
	Y    float64
	Diff float64
	Lim  bool

	// internal state
	pi FT_PIWL
	co CTRL_OUT
}

// Update executes the PI controller logic.
func (c *CTRL_PI) Update(act, set, sup, ofs, mI, kp, ki, ll, lh float64, man, rst bool) {
	c.Diff = CTRL_IN(set, act, sup)
	c.pi.Update(c.Diff, kp, ki, ll, lh, rst)
	c.co.Update(c.pi.Y, ofs, mI, ll, lh, man)
	c.Y = c.co.Y
	c.Lim = c.co.Lim
}

// CTRL_PID is a PID controller with manual functionality.
type CTRL_PID struct {
	Y    float64
	Diff float64
	Lim  bool

	// internal state
	pid FT_PIDWL
	co  CTRL_OUT
}

// Update executes the PID controller logic.
func (c *CTRL_PID) Update(act, set, sup, ofs, mI, kp, tn, tv, ll, lh float64, man, rst bool) {
	c.Diff = CTRL_IN(set, act, sup)
	c.pid.Update(c.Diff, kp, tn, tv, ll, lh, rst)
	c.co.Update(c.pid.Y, ofs, mI, ll, lh, man)
	c.Y = c.co.Y
	c.Lim = c.co.Lim
}

// CTRL_PWM converts a controller output to a PWM signal.
type CTRL_PWM struct {
	Q  bool
	pw PWM_DC
}

// Update executes the PWM logic.
func (c *CTRL_PWM) Update(ci, manIn, f float64, manual bool) {
	var dc float64
	if manual {
		dc = manIn
	} else {
		dc = ci
	}
	c.pw.Update(f, dc)
	c.Q = c.pw.Q
}

// DEAD_BAND is a linear transfer function with a dead band.
// Y = X - L for X > L
// Y = X + L for X < -L
// Y = 0 for |X| <= L
func DEAD_BAND(x, l float64) float64 {
	if x > l {
		return x - l
	}
	if x < -l {
		return x + l
	}
	return 0.0
}

// DEAD_BAND_A is a dead band function with automatic width calculation.
type DEAD_BAND_A struct {
	Y float64
	L float64

	// internal state
	tp1, tp2 FT_PT1
}

// Update executes the logic.
func (d *DEAD_BAND_A) Update(x, kl, lm float64, t time.Duration) {
	d.tp1.Update(x, t, 1.0)
	d.tp2.Update(math.Abs(d.tp1.Out-x), t*4, 1.0)
	d.L = math.Min(kl*d.tp2.Out, lm)

	if x > d.L {
		d.Y = x - d.L
	} else if x < -d.L {
		d.Y = x + d.L
	} else {
		d.Y = 0.0
	}
}

// DEAD_ZONE is a linear transfer function where Y=X if |X| > L, otherwise Y=0.
func DEAD_ZONE(x, l float64) float64 {
	if math.Abs(x) > l {
		return x
	}
	return 0.0
}

// DEAD_ZONE2 is a dead zone with hysteresis.
type DEAD_ZONE2 struct {
	Y float64
}

// Update executes the logic.
func (d *DEAD_ZONE2) Update(x, l float64) {
	if math.Abs(x) > l {
		d.Y = x
	} else if d.Y > 0.0 {
		d.Y = l
	} else {
		d.Y = -l
	}
}

// FT_DERIV calculates the derivative of a signal.
type FT_DERIV struct {
	Out float64

	// internal state
	old  float64
	last int64
	init bool
}

// Update executes the derivative calculation.
func (f *FT_DERIV) Update(in, k float64, run bool) {
	tx := logic.T_PLC_US()
	tc := float64(tx - f.last)
	f.last = tx

	if !f.init {
		f.init = true
		f.old = in
	} else if run && tc > 0.0 {
		f.Out = (in - f.old) / tc * 1000000.0 * k
		f.old = in
	} else {
		f.Out = 0.0
	}
}

// FT_IMP is an impulse filter (high-pass).
type FT_IMP struct {
	Out float64
	t1  FT_PT1
}

// Update executes the filter logic.
func (f *FT_IMP) Update(in, k float64, t time.Duration) {
	f.t1.Update(in, t, 1.0)
	f.Out = (in - f.t1.Out) * k
}

// FT_INT is an integrator with limits.
type FT_INT struct {
	Out float64
	Lim bool

	// internal state
	integ Integrate
}

// Update executes the integration logic.
func (f *FT_INT) Update(in, k, outMin, outMax float64, run, rst bool) {
	if rst {
		f.Out = 0.0
	} else {
		// The original uses Y as IN_OUT, so we pass f.Out as the initial value.
		f.integ.Update(in, k, run, &f.Out)
	}

	if f.Out >= outMax {
		f.Out = outMax
		f.Lim = true
	} else if f.Out <= outMin {
		f.Out = outMin
		f.Lim = true
	} else {
		f.Lim = false
	}
}

// FT_INT2 is a double-precision integrator.
type FT_INT2 struct {
	Out float64
	Lim bool

	// internal state
	integ Integrate
	ix    float64
	val   basic.REAL2
}

// Update executes the integration logic.
func (f *FT_INT2) Update(in, k, outMin, outMax float64, run, rst bool) {
	if rst {
		f.val = beeMath.R2_SET(0.0)
		f.Out = 0.0 // ST: out := 0.0;
	} else {
		f.ix = 0.0 // Reset temporary integrator value
		f.integ.Update(in, k, run, &f.ix)
		f.val = beeMath.R2_ADD(f.val, float32(f.ix))
		f.Out = float64(f.val.Rx)
	}

	if f.Out > outMin && f.Out < outMax {
		f.Lim = false
	} else {
		f.Out = beeMath.LIMIT(outMin, f.Out, outMax) // ST: OUT := LIMIT(OUT_MIN, OUT, OUT_MAX);
		f.val = beeMath.R2_SET(float32(f.Out))
		f.Lim = true
	}
}

// FT_PD is a PD controller.
type FT_PD struct {
	Y    float64
	diff FT_DERIV
}

// Update executes the PD controller logic.
func (f *FT_PD) Update(in, kp, tv float64) {
	f.diff.Update(in, tv, true)
	f.Y = kp * (f.diff.Out + in)
}

// FT_PDT1 is a PD controller with a first-order lag on the derivative part.
type FT_PDT1 struct {
	Y    float64
	diff FT_DERIV
	tp   FT_PT1
}

// Update executes the PDT1 controller logic.
func (f *FT_PDT1) Update(in, kp, tv, t1 float64) {
	f.diff.Update(in, tv, true)
	f.tp.Update(f.diff.Out, time.Duration(t1*float64(time.Millisecond)), 1.0)
	f.Y = kp * (f.tp.Out + in)
}

// FT_PI is a PI controller.
type FT_PI struct {
	Y     float64
	Lim   bool
	integ FT_INT
}

// Update executes the PI controller logic.
func (f *FT_PI) Update(in, kp, ki, ilimL, ilimH float64, ien, rst bool) {
	f.integ.Update(in, ki, ilimL, ilimH, ien, rst)
	f.Lim = f.integ.Lim
	f.Y = kp*in + f.integ.Out
}

// FT_PID is a PID controller.
type FT_PID struct {
	Y     float64
	Lim   bool
	integ FT_INT
	diff  FT_DERIV
}

// Update executes the PID controller logic.
func (f *FT_PID) Update(in, kp, tn, tv, ilimL, ilimH float64, ien, rst bool) {
	var ki float64
	if tn > 0.0 {
		ki = 1.0 / tn
	}
	f.integ.Update(in, ki, ilimL, ilimH, ien, rst)
	f.diff.Update(in, tv, true)
	f.Y = kp * (f.integ.Out + f.diff.Out + in)
	f.Lim = f.integ.Lim
}

// FT_PIDW is a PID controller with anti-windup.
type FT_PIDW struct {
	Y     float64
	Lim   bool
	integ Integrate
	diff  FT_DERIV
	yi    float64
}

// Update executes the PIDW controller logic.
func (f *FT_PIDW) Update(in, kp, tn, tv, limL, limH float64, rst bool) {
	if tn == 0.0 || rst {
		f.yi = 0.0
		f.integ.Update(0, 0, false, &f.yi) // Reset integrator
	} else {
		f.integ.Update(in, 1.0/tn, !f.Lim, &f.yi)
	}

	f.Y = kp * (in + f.yi)
	f.diff.Update(in, tv, true)

	// Set Lim before adding derivative part
	f.Lim = f.Y <= limL || f.Y >= limH

	f.Y = beeMath.LIMIT(limL, f.Y+kp*f.diff.Out, limH)
}

// FT_PIDWL is a PID controller with anti-windup and output limiting.
type FT_PIDWL struct {
	Y    float64
	Lim  bool
	piwl FT_PIWL
	diff FT_DERIV
}

// Update executes the PIDWL logic.
func (f *FT_PIDWL) Update(in, kp, tn, tv, limL, limH float64, rst bool) {
	if rst {
		f.piwl.Update(0, 0, 0, 0, 0, true)
		f.Y = 0
		f.Lim = false
		return
	}

	var ki float64
	if tn > 0.0 {
		ki = 1.0 / tn
	}

	f.piwl.Update(in*kp, 1.0, ki, limL, limH, false)
	f.diff.Update(in, kp*tv, true)
	f.Y = f.piwl.Y + f.diff.Out

	if f.Y < limL {
		f.Lim = true
		f.Y = limL
	} else if f.Y > limH {
		f.Lim = true
		f.Y = limH
	} else {
		f.Lim = false
	}
}

// FT_PIW is a PI controller with anti-windup.
type FT_PIW struct {
	Y     float64
	Lim   bool
	integ FT_INT
}

// Update executes the PIW logic.
func (f *FT_PIW) Update(in, kp, ki, limL, limH float64, rst bool) {
	f.integ.Update(in, ki, -1e38, 1e38, !f.Lim, rst) // Integrator limits are not used here
	f.Y = kp*in + f.integ.Out

	if f.Y < limL {
		f.Y = limL
		f.Lim = true
	} else if f.Y > limH {
		f.Y = limH
		f.Lim = true
	} else {
		f.Lim = false
	}
}

// FT_PIWL is a PI controller with anti-windup and output limiting.
type FT_PIWL struct {
	Y   float64
	Lim bool

	// internal state
	init   bool
	tLast  int64
	inLast float64
	i      float64
}

// Update executes the PIWL logic.
func (f *FT_PIWL) Update(in, kp, ki, limL, limH float64, rst bool) {
	if !f.init || rst {
		f.init = true
		f.inLast = in
		f.tLast = logic.T_PLC_US()
		f.i = 0.0
		f.Y = 0.0
		f.Lim = false
		return
	}

	tx := logic.T_PLC_US()
	tc := float64(tx - f.tLast)
	f.tLast = tx

	p := kp * in
	f.i += (in + f.inLast) * 5.0e-7 * ki * tc
	f.inLast = in

	f.Y = p + f.i

	if f.Y >= limH {
		f.Y = limH
		if ki != 0.0 {
			f.i = limH - p
		} else {
			f.i = 0.0
		}
		f.Lim = true
	} else if f.Y <= limL {
		f.Y = limL
		if ki != 0.0 {
			f.i = limL - p
		} else {
			f.i = 0.0
		}
		f.Lim = true
	} else {
		f.Lim = false
	}
}

// FT_PT1 is a first-order low-pass filter.
type FT_PT1 struct {
	Out float64

	// internal state
	last int64
	init bool
}

// Update executes the filter logic.
func (f *FT_PT1) Update(in float64, t time.Duration, k float64) {
	tx := logic.T_PLC_US()

	if !f.init || t == 0 {
		f.init = true
		f.Out = k * in
	} else {
		tReal := float64(t.Microseconds())
		f.Out += (in*k - f.Out) * float64(tx-f.last) / tReal
		if math.Abs(f.Out) < 1.0e-20 {
			f.Out = 0.0
		}
	}
	f.last = tx
}

// FT_PT2 is a second-order low-pass filter.
type FT_PT2 struct {
	Out float64

	// internal state
	init   bool
	int1   Integrate
	int2   Integrate
	i1, i2 float64
}

// Update executes the filter logic.
func (f *FT_PT2) Update(in float64, t time.Duration, d, k float64) {
	if !f.init || t == 0 {
		f.init = true
		f.Out = k * in
		f.i2 = f.Out
		return
	}

	tn := t.Seconds()
	tn2 := tn * tn

	// The original uses Y as IN_OUT, so we pass pointers.
	f.int1.Update(in*k/tn2-f.i1*0.5*d/tn-f.i2/tn2, 1.0, true, &f.i1)
	f.int2.Update(f.i1, 1.0, true, &f.i2)
	f.Out = f.i2
}

// Integrate is a plain integrator with an IN_OUT parameter for the output.
type Integrate struct {
	// internal state
	xLast float64
	init  bool
	last  int64
}

// Update executes the integration logic. Y is a pointer to the integrated value.
func (i *Integrate) Update(x, k float64, e bool, y *float64) {
	tx := logic.T_PLC_US()

	if !i.init {
		i.init = true
		i.xLast = x
	} else if e {
		*y += (x + i.xLast) * 0.5e-6 * float64(tx-i.last) * k
		i.xLast = x
	}
	i.last = tx
}

// FT_TN8 is an 8-sample signal delay line.
// It samples the input IN at intervals of T/8.
type FT_TN8 struct {
	Out  float64
	Trig bool

	// internal state
	x    [8]float64
	cnt  int
	last time.Time
	init bool
}

// Update executes the delay logic.
func (d *FT_TN8) Update(in float64, t time.Duration) {
	tx := time.Now()
	d.Trig = false

	if !d.init {
		d.init = true
		d.x[d.cnt] = in
		d.last = tx
		return
	}

	if t > 0 && tx.Sub(d.last) >= t/8 {
		d.cnt = (d.cnt + 1) % 8
		d.Out = d.x[d.cnt]
		d.x[d.cnt] = in
		d.last = tx
		d.Trig = true
	}
}

// FT_TN16 is a 16-sample signal delay line.
// It samples the input IN at intervals of T/16.
type FT_TN16 struct {
	Out  float64
	Trig bool

	// internal state
	x    [16]float64
	cnt  int
	last time.Time
	init bool
}

// Update executes the delay logic.
func (d *FT_TN16) Update(in float64, t time.Duration) {
	tx := time.Now()
	d.Trig = false

	if !d.init {
		d.init = true
		d.x[d.cnt] = in
		d.last = tx
		return
	}

	if t > 0 && tx.Sub(d.last) >= t/16 {
		d.cnt = (d.cnt + 1) % 16
		d.Out = d.x[d.cnt]
		d.x[d.cnt] = in
		d.last = tx
		d.Trig = true
	}
}

// FT_TN64 is a 64-sample signal delay line.
// It samples the input IN at intervals of T/64.
type FT_TN64 struct {
	Out  float64
	Trig bool

	// internal state
	x    [64]float64
	cnt  int
	last time.Time
	init bool
}

// Update executes the delay logic.
func (d *FT_TN64) Update(in float64, t time.Duration) {
	tx := time.Now()
	d.Trig = false

	if !d.init {
		d.init = true
		d.x[d.cnt] = in
		d.last = tx
		return
	}

	if t > 0 && tx.Sub(d.last) >= t/64 {
		d.cnt = (d.cnt + 1) % 64
		d.Out = d.x[d.cnt]
		d.x[d.cnt] = in
		d.last = tx
		d.Trig = true
	}
}
