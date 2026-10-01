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

package math

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// F_LIN calculates the linear equation a*x + b.
func F_LIN(x, a, b iec.REAL) iec.REAL {
	return a*x + b
}

// F_LIN2 calculates the linear equation through the points x1/y1 and x2/y2.
func F_LIN2(x, x1, y1, x2, y2 iec.REAL) iec.REAL {
	return (y2-y1)/(x2-x1)*(x-x1) + y1
}

// F_POLY calculates the polynomial c[0] + c[1]*x + c[2]*x^2 + ... + c[7]*x^7.
func F_POLY(x iec.REAL, c [8]iec.REAL) iec.REAL {
	return ((((((c[7]*x+c[6])*x+c[5])*x+c[4])*x+c[3])*x+c[2])*x+c[1])*x + c[0]
}

// F_POWER calculates the power equation a*x^n.
func F_POWER(a, x, n iec.REAL) iec.REAL {
	return a * EXPT(x, n)
}

// F_QUAD calculates the quadratic equation a*x^2 + b*x + c.
func F_QUAD(x, a, b, c iec.REAL) iec.REAL {
	return (a*x+b)*x + c
}

// FRMP_B calculates a ramp from start, up if dir is true and down if not,
// that has run for td of the time tr a full ramp of 255 takes. The result is
// limited to 0..255.
func FRMP_B(start iec.BYTE, dir iec.BOOL, td, tr iec.TIME) iec.BYTE {
	if td < tr {
		step := iec.BYTE((TIME_TO_DWORD(td) << 8) / TIME_TO_DWORD(tr))
		step = min(step, SEL(dir, start, 255-start))
		if dir {
			return start + step
		}
		return start - step
	}
	return SEL[iec.BYTE](dir, 0, 255)
}

// delay is the DELAY function block of the signal processing functions,
// which FT_AVG uses.
type delay struct {
	in   iec.REAL
	n    iec.INT
	out  iec.REAL
	buf  [32]iec.REAL
	i    iec.INT
	init iec.BOOL
}

func (d *delay) run() {
	stop := LIMIT(0, d.n, 32) - 1
	switch {
	case !bool(d.init):
		d.init = true
		for i := iec.INT(0); i <= stop; i++ {
			d.buf[i] = d.in
		}
		d.out = d.in
		d.i = 0
	case stop < 0:
		d.out = d.in
	default:
		d.out = d.buf[d.i]
		d.buf[d.i] = d.in
		d.i = INC1(d.i, d.n)
	}
}

// FT_AVG calculates the moving average of the last N samples of IN, up to
// 32. A sample is taken each scan E is true. RST loads the buffer with IN.
type FT_AVG struct {
	IN  iec.REAL
	E   iec.BOOL // default TRUE
	N   iec.INT  // default 32
	RST iec.BOOL
	AVG iec.REAL

	buff delay
	init iec.BOOL
}

// INIT resets the block and sets E and N to their initial values.
func (f *FT_AVG) INIT() {
	*f = FT_AVG{E: true, N: 32}
}

// Execute runs the block once.
func (f *FT_AVG) Execute(now time.Time) {
	f.buff.n = LIMIT(0, f.N, 32)
	if !f.init || f.RST {
		for i := iec.INT(1); i <= f.N; i++ {
			f.buff.in = f.IN
			f.buff.run()
		}
		f.AVG = f.IN
		f.init = true
	} else if f.E {
		f.buff.in = f.IN
		f.buff.run()
		f.AVG = f.AVG + (f.IN-f.buff.out)/iec.REAL(f.N)
	}
}

// FT_MIN_MAX holds the minimum and maximum of IN since the first scan or the
// last RST.
type FT_MIN_MAX struct {
	IN  iec.REAL
	RST iec.BOOL
	MX  iec.REAL
	MN  iec.REAL

	init iec.BOOL
}

// INIT resets the block.
func (f *FT_MIN_MAX) INIT() { *f = FT_MIN_MAX{} }

// Execute runs the block once.
func (f *FT_MIN_MAX) Execute(now time.Time) {
	switch {
	case bool(f.RST || !f.init):
		f.MN = f.IN
		f.MX = f.IN
		f.init = true
	case f.IN < f.MN:
		f.MN = f.IN
	case f.IN > f.MX:
		f.MX = f.IN
	}
}

// FT_RMP follows IN with a ramp that rises KR and falls KF units per second.
// If RMP is false OUT follows IN directly. BUSY is true while the ramp runs
// and UD is true when it runs up.
type FT_RMP struct {
	RMP  iec.BOOL // default TRUE
	IN   iec.REAL
	KR   iec.REAL
	KF   iec.REAL
	OUT  iec.REAL
	BUSY iec.BOOL
	UD   iec.BOOL

	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets RMP to its initial value.
func (f *FT_RMP) INIT() { *f = FT_RMP{RMP: true} }

// Execute runs the block once.
func (f *FT_RMP) Execute(now time.Time) {
	tx := PLC_MS(now) - f.last
	if !f.init {
		f.init = true
		f.last = tx
		tx = 0
		f.OUT = f.IN
	}
	switch {
	case !bool(f.RMP):
		f.OUT = f.IN
		f.BUSY = false
	case f.OUT > f.IN:
		f.OUT = f.OUT - iec.REAL(tx)*0.001*f.KF
		f.OUT = max(f.IN, f.OUT)
	case f.OUT < f.IN:
		f.OUT = f.OUT + iec.REAL(tx)*0.001*f.KR
		f.OUT = min(f.IN, f.OUT)
	}
	switch {
	case f.OUT < f.IN:
		f.BUSY = true
		f.UD = true
	case f.OUT > f.IN:
		f.BUSY = true
		f.UD = false
	default:
		f.BUSY = false
	}
	f.last = f.last + tx
}

// LINEAR_INT interpolates linearly between up to 20 points xy, sorted by
// ascending x. Below and above the points the first and last segments are
// extended. xy[i] is the point i+1 of OSCAT's ARRAY[1..20, 0..1].
func LINEAR_INT(x iec.REAL, xy [20][2]iec.REAL, pts iec.INT) iec.REAL {
	pts = min(pts, 20)
	// OSCAT's index i is Go's i-1.
	i := iec.INT(2)
	for i < pts && xy[i-1][0] < x {
		i++
	}
	a, b := xy[i-2], xy[i-1]
	return ((b[1]-a[1])*x - b[1]*a[0] + a[1]*b[0]) / (b[0] - a[0])
}

// POLYNOM_INT interpolates with a polynomial through up to 5 points xy,
// sorted by ascending x. xy[i] is the point i+1 of OSCAT's
// ARRAY[1..5, 0..1].
func POLYNOM_INT(x iec.REAL, xy [5][2]iec.REAL, pts iec.INT) iec.REAL {
	pts = min(pts, 5)
	// Newton's divided differences, stored in the y values. OSCAT's
	// index i is Go's i-1.
	for i := iec.INT(1); i <= pts; i++ {
		for j := pts; j >= i+1; j-- {
			xy[j-1][1] = (xy[j-1][1] - xy[j-2][1]) / (xy[j-1][0] - xy[j-i-1][0])
		}
	}
	var out iec.REAL
	for i := pts; i >= 1; i-- {
		out = out*(x-xy[i-1][0]) + xy[i-1][1]
	}
	return out
}
