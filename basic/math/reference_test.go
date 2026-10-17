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
	gomath "math"
	"math/big"
	"testing"

	"github.com/apiarytech/royaljelly/iec"
)

// The tests below check the REAL functions against Go's math package over
// their domains, not at a few points: a wrong coefficient or a wrong branch
// of the cleaned OSCAT source shows as an error far above the 32 bits a
// REAL is calculated in.

// points are the arguments a function is checked at in each domain.
const points = 2000

// relErr is the error of got against want, relative to want, but absolute
// near 0.
func relErr(got iec.REAL, want float64) float64 {
	g := float64(got)
	switch {
	case gomath.IsNaN(want):
		if gomath.IsNaN(g) {
			return 0
		}
		return gomath.Inf(1)
	case gomath.IsInf(want, 0):
		if g == want {
			return 0
		}
		return gomath.Inf(1)
	}
	return gomath.Abs(g-want) / gomath.Max(1, gomath.Abs(want))
}

// lambertW is the principal branch of the Lambert W function, by Newton's
// method in 64 bits.
func lambertW(x float64) float64 {
	w := gomath.Log1p(x)
	if x > gomath.E {
		w = gomath.Log(x) - gomath.Log(gomath.Log(x))
	}
	for range 100 {
		e := gomath.Exp(w)
		d := (w*e - x) / (e*(w+1) - (w+2)*(w*e-x)/(2*w+2))
		w -= d
		if gomath.Abs(d) < 1e-15*(1+gomath.Abs(w)) {
			break
		}
	}
	return w
}

func modulo(x, m float64) float64 { return x - m*gomath.Floor(x/m) }

func TestRealFunctionsAgainstGo(t *testing.T) {
	tests := []struct {
		name   string
		f      func(iec.REAL) iec.REAL
		ref    func(float64) float64
		lo, hi float64
		tol    float64
	}{
		{"ACOSH", ACOSH, gomath.Acosh, 1, 1e4, 1e-5},
		{"ASINH", ASINH, gomath.Asinh, 0, 1e4, 1e-5},
		// OSCAT's LN(x + SQRT(x²+1)) cancels for negative x: a REAL has only
		// a few digits of x + SQRT(x²+1) left. See TestASINHNegative.
		{"ASINH (x < 0)", ASINH, gomath.Asinh, -10, 0, 1e-5},
		{"ATANH", ATANH, gomath.Atanh, -0.999, 0.999, 1e-5},
		{"ACOTH", ACOTH, func(x float64) float64 { return gomath.Atanh(1 / x) }, 1.001, 1e3, 1e-5},
		{"COSH", COSH, gomath.Cosh, -80, 80, 1e-5},
		{"SINH", SINH, gomath.Sinh, -80, 80, 1e-5},
		{"TANH", TANH, gomath.Tanh, -20, 20, 1e-5},
		{"COTH", COTH, func(x float64) float64 { return 1 / gomath.Tanh(x) }, 0.01, 20, 1e-5},
		// OSCAT's GAMMA is an approximation of Stirling's kind, good for a large
		// x and less so towards 0; see TestGAMMASmall.
		{"GAMMA", GAMMA, gomath.Gamma, 2, 34, 5e-5},
		{"GAMMA (1..2)", GAMMA, gomath.Gamma, 1, 2, 1e-3},
		{"GAMMA (0.5..1)", GAMMA, gomath.Gamma, 0.5, 1, 1e-2},
		{"ERF", ERF, gomath.Erf, -5, 5, 2e-4},
		{"ERFC", ERFC, gomath.Erfc, -5, 5, 2e-4},
		{"GDF", GDF, func(x float64) float64 { return 2 * gomath.Atan(gomath.Tanh(x/2)) }, -20, 20, 1e-5},
		{"AGDF", AGDF, func(x float64) float64 { return gomath.Atanh(gomath.Sin(x)) }, -1.5, 1.5, 1e-4},
		{"EXP10", EXP10, func(x float64) float64 { return gomath.Pow(10, x) }, -30, 30, 1e-5},
		{"SIGMOID", SIGMOID, func(x float64) float64 { return 1 / (1 + gomath.Exp(-x)) }, -80, 80, 1e-5},
		{"SINC", SINC, func(x float64) float64 { return gomath.Sin(x) / x }, 1e-3, 100, 1e-5},
		{"TANC", TANC, func(x float64) float64 { return gomath.Tan(x) / x }, 1e-3, 1.5, 1e-4},
		{"GOLD", GOLD, func(x float64) float64 { return (x + gomath.Sqrt(x*x+4)) / 2 }, 0, 1e4, 1e-5},
		// coth(x) - 1/x cancels for small x as well.
		{"LANGEVIN", LANGEVIN, func(x float64) float64 { return 1/gomath.Tanh(x) - 1/x }, 0.5, 50, 1e-5},
		// At the branch point -1/e itself the iteration is off by 0.0024.
		{"LAMBERT_W", LAMBERT_W, lambertW, -0.36, 1e4, 1e-5},
		{"DEG", DEG, func(x float64) float64 { return modulo(x*180/gomath.Pi, 360) }, -100, 100, 1e-3},
		{"RAD", RAD, func(x float64) float64 { return modulo(x, 360) * gomath.Pi / 180 }, -1000, 1000, 1e-4},
		{"FRACT", FRACT, func(x float64) float64 { return x - gomath.Trunc(x) }, -1000, 1000, 1e-4},
	}
	for _, tt := range tests {
		worst, at := 0.0, 0.0
		for i := range points + 1 {
			x := iec.REAL(tt.lo + (tt.hi-tt.lo)*float64(i)/points)
			// The reference gets the argument as the REAL the port gets.
			if e := relErr(tt.f(x), tt.ref(float64(x))); e > worst {
				worst, at = e, float64(x)
			}
		}
		t.Logf("%-14s worst error %.2g at %g", tt.name, worst, at)
		if worst > tt.tol {
			t.Errorf("%s(%g) is off by %.2g, more than %g", tt.name, at, worst, tt.tol)
		}
	}
}

// TestASINHNegative documents how far OSCAT's ASINH is off for negative
// arguments, where its formula cancels: by 0.3 % at -1000, and from about
// -4100 x + SQRT(x²+1) is 0 in 32 bits and ASINH is -Inf, on a PLC too.
func TestASINHNegative(t *testing.T) {
	for _, x := range []float64{-100, -1000} {
		if e := relErr(ASINH(iec.REAL(x)), gomath.Asinh(x)); e > 0.01 {
			t.Errorf("ASINH(%g) is off by %.2g", x, e)
		}
	}
	if got := ASINH(-10000); !gomath.IsInf(float64(got), -1) {
		t.Errorf("ASINH(-10000) = %v; OSCAT's formula gives -Inf", got)
	}
}

// TestGAMMASmall documents OSCAT's GAMMA towards 0: below x = 0.0913 the
// base of its power, 12x - 0.1/x, is negative and GAMMA is NaN, on a PLC
// too; above it, it is off by up to 13 % at 0.12.
func TestGAMMASmall(t *testing.T) {
	if got := GAMMA(0.05); !gomath.IsNaN(float64(got)) {
		t.Errorf("GAMMA(0.05) = %v; OSCAT's formula gives NaN", got)
	}
	for x := 0.1; x < 0.5; x += 0.01 {
		if e := relErr(GAMMA(iec.REAL(x)), gomath.Gamma(x)); e > 0.15 {
			t.Errorf("GAMMA(%g) is off by %.2g", x, e)
		}
	}
}

func TestTwoArgumentFunctionsAgainstGo(t *testing.T) {
	for i := range points {
		x := iec.REAL(-50 + 100*float64(i)/points)
		for _, y := range []iec.REAL{-7.5, -1, -0.25, 0.25, 1, 3, 12.5} {
			if e := relErr(HYPOT(x, y), gomath.Hypot(float64(x), float64(y))); e > 1e-6 {
				t.Fatalf("HYPOT(%g, %g) is off by %.2g", x, y, e)
			}
			if e := relErr(ATAN2(y, x), gomath.Atan2(float64(y), float64(x))); e > 1e-5 {
				t.Fatalf("ATAN2(%g, %g) is off by %.2g", y, x, e)
			}
			if e := relErr(MODR(x, y), modulo(float64(x), float64(y))); e > 1e-4 {
				t.Fatalf("MODR(%g, %g) is off by %.2g", x, y, e)
			}
		}
	}
	for n := iec.INT(-12); n <= 12; n++ {
		for _, x := range []iec.REAL{-3, -1.5, -1, -0.5, 0.5, 1, 1.5, 2, 3} {
			if e := relErr(EXPN(x, n), gomath.Pow(float64(x), float64(n))); e > 1e-5 {
				t.Errorf("EXPN(%g, %d) is off by %.2g", x, n, e)
			}
		}
		if n > 0 {
			for _, x := range []iec.REAL{0.5, 2, 27, 1000} {
				if e := relErr(SQRTN(x, n), gomath.Pow(float64(x), 1/float64(n))); e > 1e-5 {
					t.Errorf("SQRTN(%g, %d) is off by %.2g", x, n, e)
				}
			}
		}
	}
	// BETA is GAMMA(x) * GAMMA(y) / GAMMA(x + y), with GAMMA's error.
	for _, x := range []float64{1, 2, 3.5, 10} {
		for _, y := range []float64{1, 2, 4.25, 10} {
			want := gomath.Gamma(x) * gomath.Gamma(y) / gomath.Gamma(x+y)
			if e := relErr(BETA(iec.REAL(x), iec.REAL(y)), want); e > 1e-3 {
				t.Errorf("BETA(%g, %g) is off by %.2g", x, y, e)
			}
		}
	}
}

func TestIntegerFunctionsAgainstGo(t *testing.T) {
	fact, a, b := int64(1), int64(0), int64(1) // n!, FIB(n), FIB(n+1)
	for n := iec.INT(0); n <= 46; n++ {
		if n > 0 {
			fact *= int64(n)
		}
		if n <= 12 && int64(FACT(n)) != fact {
			t.Errorf("FACT(%d) = %d, want %d", n, FACT(n), fact)
		}
		if int64(FIB(n)) != a {
			t.Errorf("FIB(%d) = %d, want %d", n, FIB(n), a)
		}
		a, b = b, a+b
	}
	// From n = 30 BINOM's product overflows a DINT before it is divided, on
	// a PLC too.
	for n := iec.INT(0); n <= 29; n++ {
		for k := iec.INT(0); k <= n; k++ {
			want := new(big.Int).Binomial(int64(n), int64(k)).Int64()
			if int64(BINOM(n, k)) != want {
				t.Errorf("BINOM(%d, %d) = %d, want %d", n, k, BINOM(n, k), want)
			}
		}
	}
	gcd := func(a, b int64) int64 {
		a, b = max(a, -a), max(b, -b)
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	for a := iec.DINT(-60); a <= 60; a++ {
		for b := iec.DINT(-60); b <= 60; b++ {
			if int64(GCD(a, b)) != gcd(int64(a), int64(b)) {
				t.Fatalf("GCD(%d, %d) = %d, want %d", a, b, GCD(a, b), gcd(int64(a), int64(b)))
			}
		}
	}
	// GCD is an INT, so the pairs have a divisor below 32768.
	for _, p := range [][2]iec.DINT{{1 << 20, 3 << 14}, {1_000_000, 7_500}, {2_147_483_647, 2}, {32760, 65520}, {-49_152, 1 << 30}} {
		if int64(GCD(p[0], p[1])) != gcd(int64(p[0]), int64(p[1])) {
			t.Errorf("GCD(%d, %d) = %d", p[0], p[1], GCD(p[0], p[1]))
		}
	}
}
