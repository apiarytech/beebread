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
	"testing"

	"github.com/apiarytech/royaljelly/iec"
)

// The expected values below follow from the OSCAT definitions: INTEGRATE and
// FT_INT integrate with the trapezoid rule, FT_DERIV divides the change by
// the time since its last run, and the first run of each only initializes.

func TestFT_PD(t *testing.T) {
	c := newClock()
	var f FT_PD
	f.INIT()
	f.KP, f.TV = 2, 0.5
	f.Execute(c.now) // the derivative initializes
	if f.Y != 0 {
		t.Fatalf("Y = %v at start", f.Y)
	}
	f.IN = 1
	f.Execute(c.stepS(1)) // derivative 1/s * TV 0.5 = 0.5; Y = 2*(1+0.5)
	if !near(f.Y, 3, 1e-4) {
		t.Fatalf("Y = %v after the step, want 3", f.Y)
	}
	f.Execute(c.stepS(1)) // no change: only the P part
	if !near(f.Y, 2, 1e-4) {
		t.Fatalf("Y = %v when steady, want 2", f.Y)
	}

	// Without INIT the derivative still starts with its initial values.
	var g FT_PD
	g.KP = 1
	g.Execute(c.now)
	g.IN = 1
	g.Execute(c.stepS(1))
	if !near(g.Y, 1, 1e-4) { // TV is 0 without INIT: no derivative
		t.Fatalf("FT_PD without INIT: Y = %v", g.Y)
	}
}

func TestFT_PDT1(t *testing.T) {
	c := newClock()
	var f FT_PDT1
	f.INIT()
	f.T1 = 1000 // OSCAT reads T1 as milliseconds: a 1 s filter
	f.Execute(c.now)
	f.IN = 1
	// The derivative jumps to 1/0.01 s = 100; the filter takes 1% of it.
	f.Execute(c.stepS(0.01))
	if !near(f.Y, 2, 1e-3) {
		t.Fatalf("Y = %v after the step, want 1 + 1", f.Y)
	}
	f.Execute(c.stepS(0.01)) // the filtered derivative decays by 1%
	if !near(f.Y, 1.99, 1e-3) {
		t.Fatalf("Y = %v one scan later, want 1.99", f.Y)
	}
	for range 2000 {
		f.Execute(c.stepS(0.01))
	}
	if !near(f.Y, 1, 1e-3) {
		t.Fatalf("Y = %v when steady, want KP*IN = 1", f.Y)
	}
}

func TestFT_PI(t *testing.T) {
	c := newClock()
	var f FT_PI
	f.INIT()
	f.KP, f.KI, f.IN = 2, 0.5, 1
	f.Execute(c.now)
	if !near(f.Y, 2, 1e-4) {
		t.Fatalf("Y = %v at start, want KP*IN = 2", f.Y)
	}
	f.Execute(c.stepS(2)) // integral: 1 * 2 s * KI 0.5 = 1
	if !near(f.Y, 3, 1e-4) || f.LIM {
		t.Fatalf("Y = %v LIM = %v after 2 s, want 3", f.Y, f.LIM)
	}

	f.IEN = false // the integral holds
	f.Execute(c.stepS(2))
	if !near(f.Y, 3, 1e-4) {
		t.Fatalf("Y = %v with IEN false, want the integral held at 1", f.Y)
	}

	f.IEN, f.ILIM_H = true, 1.5
	for range 5 {
		f.Execute(c.stepS(1))
	}
	if !near(f.Y, 3.5, 1e-4) || !f.LIM {
		t.Fatalf("Y = %v LIM = %v, want the integral limited to 1.5", f.Y, f.LIM)
	}

	f.RST = true
	f.Execute(c.stepS(1))
	if !near(f.Y, 2, 1e-4) {
		t.Fatalf("Y = %v after RST, want the integral cleared", f.Y)
	}
}

func TestFT_PIW(t *testing.T) {
	c := newClock()
	var f FT_PIW
	f.INIT()
	f.LIM_H, f.LIM_L, f.IN = 3, -5, 1
	f.Execute(c.now)
	if !near(f.Y, 1, 1e-4) || f.LIM {
		t.Fatalf("Y = %v at start", f.Y)
	}
	for range 10 {
		f.Execute(c.stepS(1))
	}
	if f.Y != 3 || !f.LIM {
		t.Fatalf("Y = %v LIM = %v, want held at LIM_H 3", f.Y, f.LIM)
	}
	// Anti wind-up: the integral stopped near the limit instead of growing
	// for ten seconds, so a negative error brings Y down at once.
	f.IN = -1
	f.Execute(c.stepS(1))
	if f.LIM || f.Y > 2.5 {
		t.Fatalf("Y = %v LIM = %v after the error turned, want below the limit", f.Y, f.LIM)
	}
	f.IN = -20
	f.Execute(c.stepS(1))
	if f.Y != -5 || !f.LIM {
		t.Fatalf("Y = %v LIM = %v, want held at LIM_L -5", f.Y, f.LIM)
	}
	f.RST, f.IN = true, 1
	f.Execute(c.stepS(1))
	if !near(f.Y, 1, 1e-4) {
		t.Fatalf("Y = %v after RST, want KP*IN", f.Y)
	}

	var g FT_PIW // without INIT: the integrator still has its limits
	g.KP, g.LIM_H, g.IN = 1, 10, 1
	g.Execute(c.now)
	g.Execute(c.stepS(1))
	if !near(g.Y, 1, 1e-4) { // KI is 0 without INIT
		t.Fatalf("FT_PIW without INIT: Y = %v", g.Y)
	}
}

func TestFT_PIDW(t *testing.T) {
	c := newClock()
	var f FT_PIDW
	f.INIT()
	f.TV, f.LIM_H, f.IN = 0, 2, 1
	f.Execute(c.now)
	if !near(f.Y, 1, 1e-4) || f.LIM {
		t.Fatalf("Y = %v at start", f.Y)
	}
	f.Execute(c.stepS(1)) // integral 1 / TN 1 s = 1; Y = 1 + 1 = LIM_H
	if f.Y != 2 || !f.LIM {
		t.Fatalf("Y = %v LIM = %v after 1 s, want 2 at the limit", f.Y, f.LIM)
	}
	f.Execute(c.stepS(1)) // at the limit the integral stops
	if f.Y != 2 || f.yi != 1 {
		t.Fatalf("Y = %v integral = %v, want the integral stopped at 1", f.Y, f.yi)
	}
	f.TN = 0 // no integral
	f.Execute(c.stepS(1))
	if !near(f.Y, 1, 1e-4) {
		t.Fatalf("Y = %v with TN 0, want KP*IN", f.Y)
	}

	// The derivative: with TN 0, a step of 1 in 1 s adds KP * TV * 1.
	var d FT_PIDW
	d.INIT()
	d.TN = 0
	d.Execute(c.now)
	d.IN = 1
	d.Execute(c.stepS(1))
	if !near(d.Y, 2, 1e-4) {
		t.Fatalf("FT_PIDW derivative: Y = %v, want 1 + 1", d.Y)
	}
	d.RST = true
	d.Execute(c.stepS(1))
	if !near(d.Y, 1, 1e-4) {
		t.Fatalf("FT_PIDW after RST: Y = %v", d.Y)
	}
}

func TestFT_PIDWL(t *testing.T) {
	c := newClock()
	var f FT_PIDWL
	f.INIT()
	f.KP, f.TV, f.LIM_H, f.IN = 2, 0, 5, 1
	f.Execute(c.now) // the PI part initializes
	if f.Y != 0 {
		t.Fatalf("Y = %v at start", f.Y)
	}
	f.Execute(c.stepS(1)) // P = 2; integral (2+2)/2 * 1 s = 2
	if !near(f.Y, 4, 1e-4) {
		t.Fatalf("Y = %v after 1 s, want 4", f.Y)
	}
	f.Execute(c.stepS(1)) // 6 is above LIM_H: the integral is set to 5 - 2
	if f.Y != 5 {
		t.Fatalf("Y = %v, want held at 5", f.Y)
	}
	f.IN = 0 // integral 3 + (0+2)/2 = 4, P = 0
	f.Execute(c.stepS(1))
	if !near(f.Y, 4, 1e-4) {
		t.Fatalf("Y = %v after the error dropped, want 4 (anti wind-up)", f.Y)
	}

	f.LIM_L, f.IN = -1, -10
	f.Execute(c.stepS(1))
	if f.Y != -1 {
		t.Fatalf("Y = %v, want held at LIM_L -1", f.Y)
	}

	f.RST = true // resets the PI part; Y keeps its value this scan
	f.Execute(c.stepS(1))
	if f.Y != -1 {
		t.Fatalf("Y = %v during RST", f.Y)
	}

	// The derivative, with no integral (TN 0): Y = P + KP*TV*dIN/dt.
	var d FT_PIDWL
	d.INIT()
	d.TN = 0
	d.Execute(c.now)
	d.IN = 1
	d.Execute(c.stepS(1))
	if !near(d.Y, 2, 1e-4) {
		t.Fatalf("FT_PIDWL derivative: Y = %v, want 1 + 1", d.Y)
	}
	// LIM reports the outer limit only: with the PI part held exactly at
	// LIM_H and no derivative, Y equals LIM_H and LIM stays false, as in
	// OSCAT. A rising input pushes the sum past it.
	d.LIM_H = 0.5
	d.Execute(c.stepS(1))
	if d.Y != 0.5 || d.LIM {
		t.Fatalf("FT_PIDWL at the PI limit: Y = %v LIM = %v", d.Y, d.LIM)
	}
	d.IN = 2
	d.Execute(c.stepS(1))
	if d.Y != 0.5 || !d.LIM {
		t.Fatalf("FT_PIDWL past the limit: Y = %v LIM = %v", d.Y, d.LIM)
	}
}

func TestCTRL_PID(t *testing.T) {
	c := newClock()
	var p CTRL_PID
	p.INIT()
	p.SET, p.ACT = 10, 8
	p.Execute(c.now)
	if p.DIFF != 2 || p.Y != 0 {
		t.Fatalf("DIFF = %v Y = %v at start", p.DIFF, p.Y)
	}
	p.Execute(c.stepS(1)) // P 2 + integral (2+2)/2 * 1 s / TN 1 s
	if !near(p.Y, 4, 1e-3) || p.LIM {
		t.Fatalf("Y = %v after 1 s, want 4", p.Y)
	}

	p.OFS = 1
	p.Execute(c.stepS(0)) // the offset is added after the controller
	if !near(p.Y, 5, 1e-3) {
		t.Fatalf("Y = %v with OFS 1, want 5", p.Y)
	}

	p.MAN, p.M_I = true, 7 // manual: M_I + OFS
	p.Execute(c.stepS(1))
	if !near(p.Y, 8, 1e-3) {
		t.Fatalf("Y = %v in manual, want 8", p.Y)
	}
	p.MAN = false

	p.LH = 3
	p.Execute(c.stepS(1))
	if p.Y != 3 || !p.LIM {
		t.Fatalf("Y = %v LIM = %v, want held at LH 3", p.Y, p.LIM)
	}

	p.SUP = 5 // the error is within the noise band: DIFF 0
	p.Execute(c.stepS(1))
	if p.DIFF != 0 {
		t.Fatalf("DIFF = %v within SUP, want 0", p.DIFF)
	}

	var q CTRL_PID // without INIT: the PID inside starts with its initial values
	q.KP, q.TN, q.LL, q.LH, q.SET = 1, 1, -100, 100, 1
	q.Execute(c.now)
	q.Execute(c.stepS(1))
	if !near(q.Y, 2, 1e-3) {
		t.Fatalf("CTRL_PID without INIT: Y = %v, want 2", q.Y)
	}
}

func TestCTRL_PWM(t *testing.T) {
	c := newClock()
	var p CTRL_PWM
	p.INIT()
	p.F, p.CI = 1, 0.25 // 1 Hz, on for 25% of each second
	high := 0
	for range 1000 {
		p.Execute(c.stepS(0.001))
		if p.Q {
			high++
		}
	}
	if high < 240 || high > 260 {
		t.Fatalf("on for %d ms of 1000, want 250", high)
	}

	p.MANUAL, p.MAN_IN = true, 0.75
	high = 0
	for range 1000 {
		p.Execute(c.stepS(0.001))
		if p.Q {
			high++
		}
	}
	if high < 740 || high > 760 {
		t.Fatalf("manual: on for %d ms of 1000, want 750", high)
	}
}

func TestDEAD_BAND_A(t *testing.T) {
	c := newClock()
	var d DEAD_BAND_A
	d.INIT()
	d.T, d.LM, d.X = secs(0.1), 0.3, 5
	for range 100 {
		d.Execute(c.stepS(0.01))
	}
	if d.L != 0 || d.Y != 5 {
		t.Fatalf("steady X: L = %v Y = %v, want no band", d.L, d.Y)
	}
	// Noise of ±1 around 0: the band grows to its limit LM.
	for i := range 2000 {
		d.X = iec.REAL(1 - 2*(i%2))
		d.Execute(c.stepS(0.01))
	}
	if d.L != 0.3 {
		t.Fatalf("L = %v with noise, want LM 0.3", d.L)
	}
	if !near(d.Y, DEAD_BAND(d.X, 0.3), 1e-6) {
		t.Fatalf("Y = %v, want DEAD_BAND(X, L)", d.Y)
	}

	var e DEAD_BAND_A // without INIT: KL is 0, so no band
	e.X, e.LM = 1, 1
	e.Execute(c.now)
	if e.Y != 1 {
		t.Fatalf("DEAD_BAND_A without INIT: Y = %v", e.Y)
	}
}

func TestFT_IMP(t *testing.T) {
	c := newClock()
	var f FT_IMP
	f.INIT()
	f.T, f.K = secs(1), 2
	f.Execute(c.now)
	f.IN = 1
	f.Execute(c.stepS(0.001))
	if !near(f.OUT, 2, 0.01) {
		t.Fatalf("OUT = %v right after the step, want about K", f.OUT)
	}
	for range 999 {
		f.Execute(c.stepS(0.001))
	}
	if !near(f.OUT, 2*0.368, 0.01) { // the low pass reached 63%
		t.Fatalf("OUT = %v after T, want K * 0.368", f.OUT)
	}
	for range 10000 {
		f.Execute(c.stepS(0.001))
	}
	if !near(f.OUT, 0, 0.001) {
		t.Fatalf("OUT = %v when steady, want 0", f.OUT)
	}
}

func TestFT_INT2(t *testing.T) {
	c := newClock()
	var f FT_INT2
	f.INIT()
	f.IN = 1
	f.Execute(c.now)
	for range 3 {
		f.Execute(c.stepS(1))
	}
	if !near(f.OUT, 3, 1e-4) || f.LIM {
		t.Fatalf("OUT = %v after 3 s, want 3", f.OUT)
	}
	f.OUT_MAX = 2
	f.Execute(c.stepS(1))
	if f.OUT != 2 || !f.LIM {
		t.Fatalf("OUT = %v LIM = %v, want held at 2", f.OUT, f.LIM)
	}
	f.RST = true
	f.Execute(c.stepS(1))
	if f.OUT != 0 {
		t.Fatalf("OUT = %v after RST", f.OUT)
	}

	// Double precision: small increments on a large value are not lost, as
	// they would be in a REAL (1.5E7 + 0.01 == 1.5E7 in single precision).
	var g FT_INT2
	g.INIT()
	g.IN = 1e7
	g.Execute(c.now)
	g.Execute(c.stepS(1)) // 1E7
	g.IN = 0.01
	g.Execute(c.stepS(1)) // + (1E7 + 0.01) / 2
	for range 99 {
		g.Execute(c.stepS(1)) // + 0.01 each
	}
	got := float64(g.val.RX) + float64(g.val.R1)
	if want := 15000000.0 + 0.005 + 0.99; got < want-0.02 || got > want+0.02 {
		t.Fatalf("double precision total = %.3f, want %.3f", got, want)
	}
}

func TestFT_PT2(t *testing.T) {
	c := newClock()
	run := func(d iec.REAL) (final, peak iec.REAL) {
		var f FT_PT2
		f.INIT()
		f.T, f.D = secs(1), d
		f.Execute(c.now)
		f.IN = 1
		for range 10000 {
			f.Execute(c.stepS(0.001))
			peak = max(peak, f.OUT)
		}
		return f.OUT, peak
	}
	// OSCAT's damping D is 4 times the damping ratio: D 4 is critical.
	if final, peak := run(4); !near(final, 1, 0.01) || peak > 1.001 {
		t.Fatalf("critical damping: final %v peak %v, want 1 without overshoot", final, peak)
	}
	if final, peak := run(0.4); peak < 1.5 || !near(final, 1, 0.5) {
		t.Fatalf("low damping: final %v peak %v, want a large overshoot", final, peak)
	}

	var z FT_PT2 // T 0: the output follows K * IN
	z.INIT()
	z.K, z.IN = 3, 2
	z.Execute(c.now)
	z.Execute(c.stepS(1))
	if z.OUT != 6 {
		t.Fatalf("T 0: OUT = %v, want 6", z.OUT)
	}
}

func TestFT_TN16AndTN64(t *testing.T) {
	c := newClock()
	var t16 FT_TN16
	t16.INIT()
	t16.T = secs(16)
	for i := range 40 {
		t16.IN = iec.REAL(i)
		t16.Execute(c.now)
		c.stepS(1)
	}
	if t16.OUT != 39-16 {
		t.Fatalf("FT_TN16 OUT = %v, want the value of 16 s ago (23)", t16.OUT)
	}
	t16.Execute(c.now) // a second later: a value is stored
	if !t16.TRIG {
		t.Fatal("FT_TN16 TRIG not set when a value was stored")
	}
	t16.Execute(c.now) // the same time again: nothing stored
	if t16.TRIG {
		t.Fatal("FT_TN16 TRIG set without a store")
	}

	var t64 FT_TN64
	t64.INIT()
	t64.T = secs(64)
	for i := range 100 {
		t64.IN = iec.REAL(i)
		t64.Execute(c.now)
		c.stepS(1)
	}
	if t64.OUT != 99-64 {
		t.Fatalf("FT_TN64 OUT = %v, want the value of 64 s ago (35)", t64.OUT)
	}
}
