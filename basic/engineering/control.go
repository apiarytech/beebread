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

// A block that holds blocks with initial input values runs their INIT the
// first time it runs, as OSCAT's instances start with those values.

// BAND_B returns 0 for X below B, 255 for X above 255-B and X otherwise.
func BAND_B(x, b iec.BYTE) iec.BYTE {
	switch {
	case x < b:
		return 0
	case x > 255-b:
		return 255
	}
	return x
}

// controlSet sets the parameters of CONTROL_SET1 and CONTROL_SET2.
func controlSet(pi, pid iec.BOOL, k, t, pK, piK, piTn, pidK, pidTn, pidTv iec.REAL, kp, tn, tv, ki, kd *iec.REAL) {
	switch {
	case bool(pi && pid):
		*kp, *tn, *tv = 0, 0, 0
	case bool(pid):
		*kp, *tn, *tv = pidK*k, pidTn*t, pidTv*t
	case bool(pi):
		*kp, *tn = piK*k, piTn*t
	default:
		*kp = pK * k
	}
	if *tn > 0.0 {
		*ki = *kp / *tn
	} else {
		*ki = 0
	}
	*kd = *kp * *tv
}

// CONTROL_SET1 calculates the parameters of a P, PI or PID controller with
// the Ziegler-Nichols method from the critical gain KT and period TT.
type CONTROL_SET1 struct {
	KT, TT                iec.REAL
	PI, PID               iec.BOOL
	P_K, PI_K, PI_TN      iec.REAL // default 0.5, 0.45, 0.83
	PID_K, PID_TN, PID_TV iec.REAL // default 0.6, 0.5, 0.125
	KP, TN, TV, KI, KD    iec.REAL
}

// INIT resets the block and sets its constants to their initial values.
func (c *CONTROL_SET1) INIT() {
	*c = CONTROL_SET1{P_K: 0.5, PI_K: 0.45, PI_TN: 0.83, PID_K: 0.6, PID_TN: 0.5, PID_TV: 0.125}
}

// Execute runs the block once.
func (c *CONTROL_SET1) Execute(now time.Time) {
	controlSet(c.PI, c.PID, c.KT, c.TT, c.P_K, c.PI_K, c.PI_TN, c.PID_K, c.PID_TN, c.PID_TV,
		&c.KP, &c.TN, &c.TV, &c.KI, &c.KD)
}

// CONTROL_SET2 calculates the parameters of a P, PI or PID controller with
// the Ziegler-Nichols method from the step response: the gain KS, the
// delay TU and the rise time TG.
type CONTROL_SET2 struct {
	KS, TU, TG            iec.REAL
	PI, PID               iec.BOOL
	P_K, PI_K, PI_TN      iec.REAL // default 1.0, 0.9, 3.33
	PID_K, PID_TN, PID_TV iec.REAL // default 1.2, 2.0, 0.5
	KP, TN, TV, KI, KD    iec.REAL

	tx iec.REAL
}

// INIT resets the block and sets its constants to their initial values.
func (c *CONTROL_SET2) INIT() {
	*c = CONTROL_SET2{P_K: 1.0, PI_K: 0.9, PI_TN: 3.33, PID_K: 1.2, PID_TN: 2.0, PID_TV: 0.5}
}

// Execute runs the block once.
func (c *CONTROL_SET2) Execute(now time.Time) {
	if c.TU > 0.0 && c.KS > 0.0 {
		c.tx = c.TG / c.TU / c.KS
	}
	// The times scale with TU.
	pidTn, pidTv := c.PID_TN, c.PID_TV
	switch {
	case bool(c.PI && c.PID):
		c.KP, c.TN, c.TV = 0, 0, 0
	case bool(c.PID):
		c.KP, c.TN, c.TV = c.PID_K*c.tx, pidTn*c.TU, pidTv*c.TU
	case bool(c.PI):
		c.KP, c.TN = c.PI_K*c.tx, c.PI_TN*c.TU
	default:
		c.KP = c.P_K * c.tx
	}
	if c.TN > 0.0 {
		c.KI = c.KP / c.TN
	} else {
		c.KI = 0
	}
	c.KD = c.KP * c.TV
}

// CTRL_IN returns the control error SET_POINT - ACTUAL, 0 within NOISE.
func CTRL_IN(setPoint, actual, noise iec.REAL) iec.REAL {
	return DEAD_ZONE(setPoint-actual, noise)
}

// CTRL_OUT is the output stage of a controller: Y is CI, or MAN_IN while
// MANUAL is true, plus OFFSET, limited to LIM_L..LIM_H. LIM is true at the
// limits.
type CTRL_OUT struct {
	CI, OFFSET, MAN_IN, LIM_L, LIM_H iec.REAL
	MANUAL                           iec.BOOL
	Y                                iec.REAL
	LIM                              iec.BOOL
}

// INIT resets the block.
func (c *CTRL_OUT) INIT() { *c = CTRL_OUT{} }

// Execute runs the block once.
func (c *CTRL_OUT) Execute(now time.Time) {
	c.Y = SEL(c.MANUAL, c.CI, c.MAN_IN) + c.OFFSET
	if c.Y > c.LIM_L && c.Y < c.LIM_H {
		c.LIM = false
	} else {
		c.Y = LIMIT(c.LIM_L, c.Y, c.LIM_H)
		c.LIM = true
	}
}

// CTRL_PI is a PI controller for the error SET - ACT, ignored within SUP,
// with the output Y limited to LL..LH, the offset OFS, and M_I while MAN is
// true.
type CTRL_PI struct {
	ACT, SET, SUP, OFS, M_I iec.REAL
	MAN                     iec.BOOL
	RST                     iec.BOOL
	KP                      iec.REAL // default 1.0
	KI                      iec.REAL // default 1.0
	LL                      iec.REAL // default -1000.0
	LH                      iec.REAL // default 1000.0
	Y                       iec.REAL
	DIFF                    iec.REAL
	LIM                     iec.BOOL

	pi          FT_PIWL
	co          CTRL_OUT
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (c *CTRL_PI) INIT() {
	*c = CTRL_PI{KP: 1, KI: 1, LL: -1000, LH: 1000, initialized: true}
	c.pi.INIT()
}

// Execute runs the block once.
func (c *CTRL_PI) Execute(now time.Time) {
	if !c.initialized {
		c.initialized = true
		c.pi.INIT()
	}
	c.DIFF = CTRL_IN(c.SET, c.ACT, c.SUP)
	c.pi.IN, c.pi.KP, c.pi.KI, c.pi.LIM_L, c.pi.LIM_H, c.pi.RST = c.DIFF, c.KP, c.KI, c.LL, c.LH, c.RST
	c.pi.Execute(now)
	c.co.CI, c.co.OFFSET, c.co.MAN_IN, c.co.LIM_L, c.co.LIM_H, c.co.MANUAL = c.pi.Y, c.OFS, c.M_I, c.LL, c.LH, c.MAN
	c.co.Execute(now)
	c.Y = c.co.Y
	c.LIM = c.co.LIM
}

// CTRL_PID is a PID controller for the error SET - ACT, ignored within SUP,
// with the output Y limited to LL..LH, the offset OFS, and M_I while MAN is
// true.
type CTRL_PID struct {
	ACT, SET, SUP, OFS, M_I iec.REAL
	MAN                     iec.BOOL
	RST                     iec.BOOL
	KP                      iec.REAL // default 1.0
	TN                      iec.REAL // default 1.0
	TV                      iec.REAL // default 1.0
	LL                      iec.REAL // default -1000.0
	LH                      iec.REAL // default 1000.0
	Y                       iec.REAL
	DIFF                    iec.REAL
	LIM                     iec.BOOL

	pid         FT_PIDWL
	co          CTRL_OUT
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (c *CTRL_PID) INIT() {
	*c = CTRL_PID{KP: 1, TN: 1, TV: 1, LL: -1000, LH: 1000, initialized: true}
	c.pid.INIT()
}

// Execute runs the block once.
func (c *CTRL_PID) Execute(now time.Time) {
	if !c.initialized {
		c.initialized = true
		c.pid.INIT()
	}
	c.DIFF = CTRL_IN(c.SET, c.ACT, c.SUP)
	p := &c.pid
	p.IN, p.KP, p.TN, p.TV, p.LIM_L, p.LIM_H, p.RST = c.DIFF, c.KP, c.TN, c.TV, c.LL, c.LH, c.RST
	p.Execute(now)
	c.co.CI, c.co.OFFSET, c.co.MAN_IN, c.co.LIM_L, c.co.LIM_H, c.co.MANUAL = p.Y, c.OFS, c.M_I, c.LL, c.LH, c.MAN
	c.co.Execute(now)
	c.Y = c.co.Y
	c.LIM = c.co.LIM
}

// CTRL_PWM is a PWM output stage of a controller with the frequency F and
// the duty cycle CI, or MAN_IN while MANUAL is true.
type CTRL_PWM struct {
	CI, MAN_IN iec.REAL
	MANUAL     iec.BOOL
	F          iec.REAL
	Q          iec.BOOL

	pw PWM_DC
}

// INIT resets the block.
func (c *CTRL_PWM) INIT() { *c = CTRL_PWM{} }

// Execute runs the block once.
func (c *CTRL_PWM) Execute(now time.Time) {
	c.pw.F = c.F
	c.pw.DC = SEL(c.MANUAL, c.CI, c.MAN_IN)
	c.pw.Execute(now)
	c.Q = c.pw.Q
}

// DEAD_BAND is a linear function with a dead band: X - L above L, X + L
// below -L and 0 between.
func DEAD_BAND(x, l iec.REAL) iec.REAL {
	switch {
	case x > l:
		return x - l
	case x < -l:
		return x + l
	}
	return 0.0
}

// DEAD_BAND_A is a DEAD_BAND whose band L is KL times the noise of X,
// measured with the time T, up to LM.
type DEAD_BAND_A struct {
	X  iec.REAL
	T  iec.TIME
	KL iec.REAL // default 1.0
	LM iec.REAL
	Y  iec.REAL
	L  iec.REAL

	tp1, tp2    FT_PT1
	initialized bool
}

// INIT resets the block and sets KL to its initial value.
func (d *DEAD_BAND_A) INIT() {
	*d = DEAD_BAND_A{KL: 1, initialized: true}
	d.tp1.INIT()
	d.tp2.INIT()
}

// Execute runs the block once.
func (d *DEAD_BAND_A) Execute(now time.Time) {
	if !d.initialized {
		d.initialized = true
		d.tp1.INIT()
		d.tp2.INIT()
	}
	d.tp1.IN, d.tp1.T = d.X, d.T
	d.tp1.Execute(now)
	d.tp2.IN, d.tp2.T = ABS(d.tp1.OUT-d.X), td.MULTIME(d.T, 4.0)
	d.tp2.Execute(now)
	d.L = min(d.KL*d.tp2.OUT, d.LM)
	d.Y = DEAD_BAND(d.X, d.L)
}

// DEAD_ZONE returns X, or 0 while |X| <= L.
func DEAD_ZONE(x, l iec.REAL) iec.REAL {
	if ABS(x) > l {
		return x
	}
	return 0.0
}

// DEAD_ZONE2 follows X with Y while |X| > L, and holds Y at L or -L within
// it.
type DEAD_ZONE2 struct {
	X, L iec.REAL
	Y    iec.REAL
}

// INIT resets the block.
func (d *DEAD_ZONE2) INIT() { *d = DEAD_ZONE2{} }

// Execute runs the block once.
func (d *DEAD_ZONE2) Execute(now time.Time) {
	switch {
	case ABS(d.X) > d.L:
		d.Y = d.X
	case d.Y > 0.0:
		d.Y = d.L
	default:
		d.Y = -d.L
	}
}

// FT_DERIV calculates the derivative of IN per second times K while RUN is
// true.
type FT_DERIV struct {
	IN  iec.REAL
	K   iec.REAL // default 1.0
	RUN iec.BOOL // default TRUE
	OUT iec.REAL

	old  iec.REAL
	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets K and RUN to their initial values.
func (f *FT_DERIV) INIT() { *f = FT_DERIV{K: 1, RUN: true} }

// Execute runs the block once.
func (f *FT_DERIV) Execute(now time.Time) {
	tx := PLC_US(now)
	tc := iec.REAL(tx - f.last)
	f.last = tx
	switch {
	case !bool(f.init):
		f.init = true
		f.old = f.IN
	case bool(f.RUN) && tc > 0.0:
		f.OUT = (f.IN - f.old) / tc * 1000000.0 * f.K
		f.old = f.IN
	default:
		f.OUT = 0.0
	}
}

// FT_IMP is a high pass filter with the time T and the factor K.
type FT_IMP struct {
	IN  iec.REAL
	T   iec.TIME
	K   iec.REAL // default 1.0
	OUT iec.REAL

	t1          FT_PT1
	initialized bool
}

// INIT resets the block and sets K to its initial value.
func (f *FT_IMP) INIT() {
	*f = FT_IMP{K: 1, initialized: true}
	f.t1.INIT()
}

// Execute runs the block once.
func (f *FT_IMP) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.t1.INIT()
	}
	f.t1.IN, f.t1.T = f.IN, f.T
	f.t1.Execute(now)
	f.OUT = (f.IN - f.t1.OUT) * f.K
}

// FT_INT integrates IN times K per second while RUN is true, within
// OUT_MIN..OUT_MAX. LIM is true at the limits and RST clears it.
type FT_INT struct {
	IN      iec.REAL
	K       iec.REAL // default 1
	RUN     iec.BOOL // default TRUE
	RST     iec.BOOL
	OUT_MIN iec.REAL // default -1E37
	OUT_MAX iec.REAL // default 1E37
	OUT     iec.REAL
	LIM     iec.BOOL

	integ INTEGRATE
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_INT) INIT() {
	*f = FT_INT{K: 1, RUN: true, OUT_MIN: -1e37, OUT_MAX: 1e37}
	f.integ.INIT()
}

// Execute runs the block once.
func (f *FT_INT) Execute(now time.Time) {
	if f.RST {
		f.OUT = 0
	} else {
		f.integ.X, f.integ.E, f.integ.K, f.integ.Y = f.IN, f.RUN, f.K, &f.OUT
		f.integ.Execute(now)
	}
	switch {
	case f.OUT >= f.OUT_MAX:
		f.OUT = f.OUT_MAX
		f.LIM = true
	case f.OUT <= f.OUT_MIN:
		f.OUT = f.OUT_MIN
		f.LIM = true
	default:
		f.LIM = false
	}
}

// FT_INT2 is FT_INT with double precision.
type FT_INT2 struct {
	IN      iec.REAL
	K       iec.REAL // default 1.0
	RUN     iec.BOOL // default TRUE
	RST     iec.BOOL
	OUT_MIN iec.REAL // default -1.0E38
	OUT_MAX iec.REAL // default 1.0E38
	OUT     iec.REAL
	LIM     iec.BOOL

	integ INTEGRATE
	ix    iec.REAL
	val   REAL2
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_INT2) INIT() {
	*f = FT_INT2{K: 1, RUN: true, OUT_MIN: -1e38, OUT_MAX: 1e38}
	f.integ.INIT()
}

// Execute runs the block once.
func (f *FT_INT2) Execute(now time.Time) {
	if f.RST {
		f.val = math.R2_SET(0)
		f.OUT = 0
	} else {
		f.integ.X, f.integ.E, f.integ.K, f.integ.Y = f.IN, f.RUN, f.K, &f.ix
		f.integ.Execute(now)
		f.val = math.R2_ADD(f.val, f.ix)
		f.ix = 0
		f.OUT = f.val.RX
	}
	if f.OUT > f.OUT_MIN && f.OUT < f.OUT_MAX {
		f.LIM = false
	} else {
		f.OUT = LIMIT(f.OUT_MIN, f.OUT, f.OUT_MAX)
		f.val = math.R2_SET(f.OUT)
		f.LIM = true
	}
}

// FT_PD is a PD controller: Y = KP * (IN + TV * d IN / dt).
type FT_PD struct {
	IN iec.REAL
	KP iec.REAL // default 1.0
	TV iec.REAL // default 1.0
	Y  iec.REAL

	diff        FT_DERIV
	initialized bool
}

// INIT resets the block and sets KP and TV to their initial values.
func (f *FT_PD) INIT() {
	*f = FT_PD{KP: 1, TV: 1, initialized: true}
	f.diff.INIT()
}

// Execute runs the block once.
func (f *FT_PD) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.diff.INIT()
	}
	f.diff.IN, f.diff.K = f.IN, f.TV
	f.diff.Execute(now)
	f.Y = f.KP * (f.diff.OUT + f.IN)
}

// FT_PDT1 is a PD controller whose derivative is filtered with the time T1
// in seconds.
type FT_PDT1 struct {
	IN iec.REAL
	KP iec.REAL // default 1.0
	TV iec.REAL // default 1.0
	T1 iec.REAL // default 1.0
	Y  iec.REAL

	diff        FT_DERIV
	tp          FT_PT1
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_PDT1) INIT() {
	*f = FT_PDT1{KP: 1, TV: 1, T1: 1, initialized: true}
	f.diff.INIT()
	f.tp.INIT()
}

// Execute runs the block once.
func (f *FT_PDT1) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.diff.INIT()
		f.tp.INIT()
	}
	f.diff.IN, f.diff.K = f.IN, f.TV
	f.diff.Execute(now)
	// OSCAT converts the seconds T1 to a TIME of T1 milliseconds.
	f.tp.IN, f.tp.T = f.diff.OUT, REAL_TO_TIME(f.T1)
	f.tp.Execute(now)
	f.Y = f.KP * (f.tp.OUT + f.IN)
}

// FT_PI is a PI controller: Y = KP * IN + KI * integral of IN, with the
// integral limited to ILIM_L..ILIM_H and running while IEN is true.
type FT_PI struct {
	IN     iec.REAL
	KP     iec.REAL // default 1.0
	KI     iec.REAL // default 1.0
	ILIM_L iec.REAL // default -1E38
	ILIM_H iec.REAL // default 1E38
	IEN    iec.BOOL // default TRUE
	RST    iec.BOOL
	Y      iec.REAL
	LIM    iec.BOOL

	integ FT_INT
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_PI) INIT() {
	*f = FT_PI{KP: 1, KI: 1, ILIM_L: -1e38, ILIM_H: 1e38, IEN: true}
	f.integ.INIT()
}

// Execute runs the block once.
func (f *FT_PI) Execute(now time.Time) {
	i := &f.integ
	i.IN, i.K, i.RUN, i.RST, i.OUT_MIN, i.OUT_MAX = f.IN, f.KI, f.IEN, f.RST, f.ILIM_L, f.ILIM_H
	i.Execute(now)
	f.LIM = i.LIM
	f.Y = f.KP*f.IN + i.OUT
}

// FT_PID is a PID controller: Y = KP * (IN + integral of IN / TN + TV *
// d IN / dt), with the integral limited to ILIM_L..ILIM_H and running while
// IEN is true.
type FT_PID struct {
	IN     iec.REAL
	KP     iec.REAL // default 1.0
	TN     iec.REAL // default 1.0
	TV     iec.REAL // default 1.0
	ILIM_L iec.REAL // default -1.0E38
	ILIM_H iec.REAL // default 1.0E38
	IEN    iec.BOOL // default TRUE
	RST    iec.BOOL
	Y      iec.REAL
	LIM    iec.BOOL

	integ       FT_INT
	diff        FT_DERIV
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_PID) INIT() {
	*f = FT_PID{KP: 1, TN: 1, TV: 1, ILIM_L: -1e38, ILIM_H: 1e38, IEN: true, initialized: true}
	f.integ.INIT()
	f.diff.INIT()
}

// Execute runs the block once.
func (f *FT_PID) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.integ.INIT()
		f.diff.INIT()
	}
	i := &f.integ
	if f.TN > 0.0 {
		i.IN, i.K, i.RUN, i.RST, i.OUT_MIN, i.OUT_MAX = f.IN, 1.0/f.TN, f.IEN, f.RST, f.ILIM_L, f.ILIM_H
	} else {
		i.RST = false
	}
	i.Execute(now)
	f.diff.IN, f.diff.K = f.IN, f.TV
	f.diff.Execute(now)
	f.Y = f.KP * (i.OUT + f.diff.OUT + f.IN)
	f.LIM = i.LIM
}

// FT_PIDW is a PID controller whose output is limited to LIM_L..LIM_H,
// where the integral stops.
type FT_PIDW struct {
	IN    iec.REAL
	KP    iec.REAL // default 1.0
	TN    iec.REAL // default 1.0
	TV    iec.REAL // default 1.0
	LIM_L iec.REAL // default -1.0E38
	LIM_H iec.REAL // default 1.0E38
	RST   iec.BOOL
	Y     iec.REAL
	LIM   iec.BOOL

	integ       INTEGRATE
	diff        FT_DERIV
	yi          iec.REAL
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_PIDW) INIT() {
	*f = FT_PIDW{KP: 1, TN: 1, TV: 1, LIM_L: -1e38, LIM_H: 1e38, initialized: true}
	f.integ.INIT()
	f.diff.INIT()
}

// Execute runs the block once.
func (f *FT_PIDW) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.integ.INIT()
		f.diff.INIT()
	}
	f.integ.Y = &f.yi
	if f.TN == 0.0 || f.RST {
		f.integ.E = false
		f.integ.Execute(now)
		f.yi = 0
	} else {
		f.integ.X, f.integ.K, f.integ.E = f.IN, 1.0/f.TN, !f.LIM
		f.integ.Execute(now)
	}
	f.Y = f.KP * (f.IN + f.yi)
	f.diff.IN, f.diff.K = f.IN, f.TV
	f.diff.Execute(now)
	f.LIM = !(f.Y > f.LIM_L && f.Y < f.LIM_H)
	f.Y = LIMIT(f.LIM_L, f.Y+f.KP*f.diff.OUT, f.LIM_H)
}

// FT_PIDWL is a PID controller whose output is limited to LIM_L..LIM_H,
// with the anti wind-up of FT_PIWL.
type FT_PIDWL struct {
	IN    iec.REAL
	KP    iec.REAL // default 1.0
	TN    iec.REAL // default 1.0
	TV    iec.REAL // default 1.0
	LIM_L iec.REAL // default -1.0E38
	LIM_H iec.REAL // default 1.0E38
	RST   iec.BOOL
	Y     iec.REAL
	LIM   iec.BOOL

	piwl        FT_PIWL
	diff        FT_DERIV
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_PIDWL) INIT() {
	*f = FT_PIDWL{KP: 1, TN: 1, TV: 1, LIM_L: -1e38, LIM_H: 1e38, initialized: true}
	f.piwl.INIT()
	f.diff.INIT()
}

// Execute runs the block once.
func (f *FT_PIDWL) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.piwl.INIT()
		f.diff.INIT()
	}
	p := &f.piwl
	if f.RST {
		p.RST = true
		p.Execute(now)
		p.RST = false
		return
	}
	var ki iec.REAL
	if f.TN != 0.0 {
		ki = 1.0 / f.TN
	}
	p.IN, p.KP, p.KI, p.LIM_L, p.LIM_H = f.IN*f.KP, 1.0, ki, f.LIM_L, f.LIM_H
	p.Execute(now)
	f.diff.IN, f.diff.K = f.IN, f.KP*f.TV
	f.diff.Execute(now)
	f.Y = p.Y + f.diff.OUT
	switch {
	case f.Y < f.LIM_L:
		f.LIM = true
		f.Y = f.LIM_L
	case f.Y > f.LIM_H:
		f.LIM = true
		f.Y = f.LIM_H
	default:
		f.LIM = false
	}
}

// FT_PIW is a PI controller whose output is limited to LIM_L..LIM_H, where
// the integral stops.
type FT_PIW struct {
	IN    iec.REAL
	KP    iec.REAL // default 1.0
	KI    iec.REAL // default 1.0
	LIM_L iec.REAL // default -1E38
	LIM_H iec.REAL // default 1E38
	RST   iec.BOOL
	Y     iec.REAL
	LIM   iec.BOOL

	integ       FT_INT
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_PIW) INIT() {
	*f = FT_PIW{KP: 1, KI: 1, LIM_L: -1e38, LIM_H: 1e38, initialized: true}
	f.integ.INIT()
}

// Execute runs the block once.
func (f *FT_PIW) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.integ.INIT()
	}
	f.integ.IN, f.integ.K, f.integ.RUN, f.integ.RST = f.IN, f.KI, !f.LIM, f.RST
	f.integ.Execute(now)
	f.Y = f.KP*f.IN + f.integ.OUT
	switch {
	case f.Y < f.LIM_L:
		f.Y = f.LIM_L
		f.LIM = true
	case f.Y > f.LIM_H:
		f.Y = f.LIM_H
		f.LIM = true
	default:
		f.LIM = false
	}
}

// FT_PIWL is a PI controller whose output is limited to LIM_L..LIM_H, where
// the integral is set so the output stays at the limit.
type FT_PIWL struct {
	IN    iec.REAL
	KP    iec.REAL // default 1.0
	KI    iec.REAL // default 1.0
	LIM_L iec.REAL // default -1.0E38
	LIM_H iec.REAL // default 1.0E38
	RST   iec.BOOL
	Y     iec.REAL
	LIM   iec.BOOL

	init         iec.BOOL
	tLast        iec.DWORD
	inLast, i, p iec.REAL
}

// INIT resets the block and sets its inputs to their initial values.
func (f *FT_PIWL) INIT() { *f = FT_PIWL{KP: 1, KI: 1, LIM_L: -1e38, LIM_H: 1e38} }

// Execute runs the block once.
func (f *FT_PIWL) Execute(now time.Time) {
	if !f.init || f.RST {
		f.init = true
		f.inLast = f.IN
		f.tLast = PLC_US(now)
		f.i = 0
		return
	}
	tx := PLC_US(now)
	tc := iec.REAL(tx - f.tLast)
	f.tLast = tx
	f.p = f.KP * f.IN
	f.i = (f.IN+f.inLast)*5.0e-7*f.KI*tc + f.i
	f.inLast = f.IN
	f.Y = f.p + f.i
	switch {
	case f.Y >= f.LIM_H:
		f.Y = f.LIM_H
		f.i = SEL[iec.REAL](f.KI != 0, 0, f.LIM_H-f.p)
		f.LIM = true
	case f.Y <= f.LIM_L:
		f.Y = f.LIM_L
		f.i = SEL[iec.REAL](f.KI != 0, 0, f.LIM_L-f.p)
		f.LIM = true
	default:
		f.LIM = false
	}
}

// FT_PT1 is a low pass filter of first order with the time T and the factor
// K.
type FT_PT1 struct {
	IN  iec.REAL
	T   iec.TIME
	K   iec.REAL // default 1.0
	OUT iec.REAL

	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets K to its initial value.
func (f *FT_PT1) INIT() { *f = FT_PT1{K: 1} }

// Execute runs the block once.
func (f *FT_PT1) Execute(now time.Time) {
	tx := PLC_US(now)
	if !f.init || f.T == 0 {
		f.init = true
		f.OUT = f.K * f.IN
	} else {
		f.OUT = f.OUT + (f.IN*f.K-f.OUT)*iec.REAL(tx-f.last)/TIME_TO_REAL(f.T)*1.0e-3
		if ABS(f.OUT) < 1.0e-20 {
			f.OUT = 0
		}
	}
	f.last = tx
}

// FT_PT2 is a low pass filter of second order with the time T, the damping
// D and the factor K.
type FT_PT2 struct {
	IN  iec.REAL
	T   iec.TIME
	D   iec.REAL
	K   iec.REAL // default 1.0
	OUT iec.REAL

	init        iec.BOOL
	int1, int2  INTEGRATE
	i1, i2      iec.REAL
	initialized bool
}

// INIT resets the block and sets K to its initial value.
func (f *FT_PT2) INIT() {
	*f = FT_PT2{K: 1, initialized: true}
	f.int1.INIT()
	f.int2.INIT()
}

// Execute runs the block once.
func (f *FT_PT2) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.int1.INIT()
		f.int2.INIT()
	}
	if !f.init || f.T == 0 {
		f.init = true
		f.OUT = f.K * f.IN
		f.i2 = f.OUT
		return
	}
	tn := TIME_TO_REAL(f.T) * 1.0e-3
	tn2 := tn * tn
	f.int1.X, f.int1.Y = f.IN*f.K/tn2-f.i1*0.5*f.D/tn-f.i2/tn2, &f.i1
	f.int1.Execute(now)
	f.int2.X, f.int2.Y = f.i1, &f.i2
	f.int2.Execute(now)
	f.OUT = f.i2
}

// ftTn delays IN by T with a buffer of len(x) values; see FT_TN8.
type ftTn struct {
	IN   iec.REAL
	T    iec.TIME
	OUT  iec.REAL
	TRIG iec.BOOL

	cnt  int
	last iec.DWORD
	init iec.BOOL
}

func (f *ftTn) run(now time.Time, x []iec.REAL) {
	tx := PLC_MS(now)
	f.TRIG = false
	if !f.init {
		x[f.cnt] = f.IN
		f.init = true
		f.last = tx
	} else if tx-f.last >= ms(f.T)/iec.DWORD(len(x)) {
		f.cnt = (f.cnt + 1) % len(x)
		f.OUT = x[f.cnt]
		x[f.cnt] = f.IN
		f.last = tx
		f.TRIG = true
	}
}

// FT_TN8 delays IN by the time T, storing 8 values in the time. TRIG is
// true for one scan when a value is stored.
type FT_TN8 struct {
	ftTn
	x [8]iec.REAL
}

// INIT resets the block.
func (f *FT_TN8) INIT() { *f = FT_TN8{} }

// Execute runs the block once.
func (f *FT_TN8) Execute(now time.Time) { f.run(now, f.x[:]) }

// FT_TN16 delays IN by the time T, storing 16 values in the time; see
// FT_TN8.
type FT_TN16 struct {
	ftTn
	x [16]iec.REAL
}

// INIT resets the block.
func (f *FT_TN16) INIT() { *f = FT_TN16{} }

// Execute runs the block once.
func (f *FT_TN16) Execute(now time.Time) { f.run(now, f.x[:]) }

// FT_TN64 delays IN by the time T, storing 64 values in the time; see
// FT_TN8.
type FT_TN64 struct {
	ftTn
	x [64]iec.REAL
}

// INIT resets the block.
func (f *FT_TN64) INIT() { *f = FT_TN64{} }

// Execute runs the block once.
func (f *FT_TN64) Execute(now time.Time) { f.run(now, f.x[:]) }

// HYST is a hysteresis: if ON >= OFF, Q switches on above ON and off below
// OFF; if ON < OFF, Q switches on below ON and off above OFF. WIN is true
// between the two.
type HYST struct {
	IN, ON, OFF iec.REAL
	Q, WIN      iec.BOOL
}

// INIT resets the block.
func (h *HYST) INIT() { *h = HYST{} }

// Execute runs the block once.
func (h *HYST) Execute(now time.Time) {
	if h.ON >= h.OFF {
		switch {
		case h.IN < h.OFF:
			h.Q, h.WIN = false, false
		case h.IN > h.ON:
			h.Q, h.WIN = true, false
		default:
			h.WIN = true
		}
		return
	}
	switch {
	case h.IN > h.OFF:
		h.Q, h.WIN = false, false
	case h.IN < h.ON:
		h.Q, h.WIN = true, false
	default:
		h.WIN = true
	}
}

// HYST_1 is a hysteresis: Q switches on above HIGH and off below LOW. WIN
// is true between the two.
type HYST_1 struct {
	IN, HIGH, LOW iec.REAL
	Q, WIN        iec.BOOL
}

// INIT resets the block.
func (h *HYST_1) INIT() { *h = HYST_1{} }

// Execute runs the block once.
func (h *HYST_1) Execute(now time.Time) {
	switch {
	case h.IN < h.LOW:
		h.Q, h.WIN = false, false
	case h.IN > h.HIGH:
		h.Q, h.WIN = true, false
	default:
		h.WIN = true
	}
}

// HYST_2 is a hysteresis around VAL: Q switches on above VAL + HYS/2 and
// off below VAL - HYS/2. WIN is true between the two.
type HYST_2 struct {
	IN, VAL, HYS iec.REAL
	Q, WIN       iec.BOOL
}

// INIT resets the block.
func (h *HYST_2) INIT() { *h = HYST_2{} }

// Execute runs the block once.
func (h *HYST_2) Execute(now time.Time) {
	tmp := h.VAL - h.HYS*0.5
	switch {
	case h.IN < tmp:
		h.Q, h.WIN = false, false
	case h.IN > tmp+h.HYS:
		h.Q, h.WIN = true, false
	default:
		h.WIN = true
	}
}

// HYST_3 is a double hysteresis: Q1 switches on below VAL1 and Q2 on above
// VAL2, each with the hysteresis HYST.
type HYST_3 struct {
	IN, HYST, VAL1, VAL2 iec.REAL
	Q1, Q2               iec.BOOL
}

// INIT resets the block.
func (h *HYST_3) INIT() { *h = HYST_3{} }

// Execute runs the block once.
func (h *HYST_3) Execute(now time.Time) {
	x := h.HYST * 0.5
	if h.IN < h.VAL1-x {
		h.Q1 = true
	} else if h.IN > h.VAL1+x {
		h.Q1 = false
	}
	if h.IN < h.VAL2-x {
		h.Q2 = false
	} else if h.IN > h.VAL2+x {
		h.Q2 = true
	}
}

// INTEGRATE integrates X times K per second into Y while E is true.
type INTEGRATE struct {
	E iec.BOOL // default TRUE
	X iec.REAL
	K iec.REAL // default 1.0
	Y *iec.REAL

	xLast iec.REAL
	init  iec.BOOL
	last  iec.DWORD
}

// INIT resets the block and sets E and K to their initial values.
func (i *INTEGRATE) INIT() { *i = INTEGRATE{Y: i.Y, E: true, K: 1} }

// Execute runs the block once.
func (i *INTEGRATE) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !i.init {
		i.init = true
		i.xLast = i.X
	} else if i.E && i.Y != nil {
		*i.Y = (i.X+i.xLast)*0.5e-3*iec.REAL(tx-i.last)*i.K + *i.Y
		i.xLast = i.X
	}
	i.last = tx
}
