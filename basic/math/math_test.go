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
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

func near(a, b, eps iec.REAL) bool {
	return gomath.Abs(float64(a-b)) <= float64(eps)
}

func TestRealFunctions(t *testing.T) {
	pi := MATH.PI
	tests := []struct {
		name      string
		got, want iec.REAL
		eps       iec.REAL
	}{
		{"ACOSH(2)", ACOSH(2), 1.3169579, 1e-5},
		{"ACOTH(2)", ACOTH(2), 0.5493061, 1e-5},
		{"ASINH(1)", ASINH(1), 0.8813736, 1e-5},
		{"ATANH(0.5)", ATANH(0.5), 0.5493061, 1e-5},
		{"ATAN2(1,-1)", ATAN2(1, -1), 2.3561945, 1e-5},
		{"ATAN2(-1,-1)", ATAN2(-1, -1), -2.3561945, 1e-5},
		{"ATAN2(1,0)", ATAN2(1, 0), pi / 2, 1e-6},
		{"ATAN2(0,0)", ATAN2(0, 0), 0, 0},
		{"CAUCHY(0,0,1)", CAUCHY(0, 0, 1), 0.3183099, 1e-5},
		{"CAUCHYCD(0,0,1)", CAUCHYCD(0, 0, 1), 0.5, 1e-6},
		{"COSH(1)", COSH(1), 1.5430806, 1e-5},
		{"COTH(1)", COTH(1), 1.3130353, 1e-5},
		{"COTH(30)", COTH(30), 1, 0},
		{"DEG(pi)", DEG(pi), 180, 1e-3},
		{"DEG(-pi/2)", DEG(-pi / 2), 270, 1e-3},
		{"ERF(1)", ERF(1), 0.8427, 1e-3},
		{"ERF(-1)", ERF(-1), -0.8427, 1e-3},
		{"ERFC(1)", ERFC(1), 0.1573, 1e-3},
		{"EXP10(3)", EXP10(3), 1000, 1e-2},
		{"EXPN(2,10)", EXPN(2, 10), 1024, 0},
		{"EXPN(2,-2)", EXPN(2, -2), 0.25, 0},
		{"EXPN(3,0)", EXPN(3, 0), 1, 0},
		{"FRACT(3.25)", FRACT(3.25), 0.25, 1e-6},
		{"FRACT(-3.25)", FRACT(-3.25), -0.25, 1e-6},
		{"GAMMA(5)", GAMMA(5), 24, 0.01},
		{"BETA(2,3)", BETA(2, 3), 1.0 / 12, 1e-3},
		{"GAUSS(0,0,1)", GAUSS(0, 0, 1), 0.3989423, 1e-5},
		{"GAUSSCD(0,0,1)", GAUSSCD(0, 0, 1), 0.5, 1e-6},
		{"GDF(1)", GDF(1), 0.8657694, 1e-5},
		{"AGDF(GDF(1))", AGDF(GDF(1)), 1, 1e-4},
		{"GOLD(1)", GOLD(1), 1.618034, 1e-5},
		{"HYPOT(3,4)", HYPOT(3, 4), 5, 0},
		{"INV(4)", INV(4), 0.25, 0},
		{"INV(0)", INV(0), 0, 0},
		{"LAMBERT_W(e)", LAMBERT_W(MATH.E), 1, 1e-5},
		{"LAMBERT_W(1)", LAMBERT_W(1), 0.5671433, 1e-5},
		{"LAMBERT_W(-1)", LAMBERT_W(-1), -1000, 0},
		{"LANGEVIN(1)", LANGEVIN(1), 0.3130353, 1e-5},
		{"MAX3", MAX3(1, 3, 2), 3, 0},
		{"MID3(3,1,2)", MID3(3, 1, 2), 2, 0},
		{"MID3(1,2,3)", MID3(1, 2, 3), 2, 0},
		{"MID3(2,3,1)", MID3(2, 3, 1), 2, 0},
		{"MIN3", MIN3(2, 1, 3), 1, 0},
		{"MODR(5.5,2.5)", MODR(5.5, 2.5), 0.5, 1e-6},
		{"MODR(-1,3)", MODR(-1, 3), 2, 1e-6},
		{"MODR(1,0)", MODR(1, 0), 0, 0},
		{"MUL_ADD", MUL_ADD(2, 3, 4), 10, 0},
		{"NEGX", NEGX(2), -2, 0},
		{"RAD(180)", RAD(180), pi, 1e-5},
		{"RND(3.1415,2)", RND(3.1415, 2), 3.1, 1e-6},
		{"RND(1234.5,2)", RND(1234.5, 2), 1200, 1e-3},
		{"ROUND(3.14159,2)", ROUND(3.14159, 2), 3.14, 1e-6},
		{"SIGMOID(0)", SIGMOID(0), 0.5, 0},
		{"SINC(0)", SINC(0), 1, 0},
		{"SINH(1)", SINH(1), 1.1752012, 1e-5},
		{"SINH(0.001)", SINH(0.001), 0.001, 0},
		{"SQRTN(27,3)", SQRTN(27, 3), 3, 1e-5},
		{"SQRTN(27,0)", SQRTN(27, 0), 0, 0},
		{"TANC(0)", TANC(0), 1, 0},
		{"TANH(1)", TANH(1), 0.7615942, 1e-5},
		{"TANH(-30)", TANH(-30), -1, 0},
		{"CIRCLE_A(1,360)", CIRCLE_A(1, 360), pi, 1e-4},
		{"CIRCLE_C(1,360)", CIRCLE_C(1, 360), 2 * pi, 1e-4},
		{"CIRCLE_SEG(1,1)", CIRCLE_SEG(1, 1), pi / 2, 1e-5},
		{"CIRCLE_SEG(1,2)", CIRCLE_SEG(1, 2), pi, 1e-5},
		{"CONE_V(1,3)", CONE_V(1, 3), pi, 1e-5},
		{"ELLIPSE_A(1,1)", ELLIPSE_A(1, 1), pi, 0},
		{"ELLIPSE_C(1,1)", ELLIPSE_C(1, 1), 2 * pi, 1e-5},
		{"SPHERE_V(1)", SPHERE_V(1), 4.18879, 1e-5},
		{"TRIANGLE_A(3,0,4,5)", TRIANGLE_A(3, 0, 4, 5), 6, 1e-5},
		{"TRIANGLE_A(3,90,4,0)", TRIANGLE_A(3, 90, 4, 0), 6, 1e-5},
		{"F_LIN", F_LIN(2, 3, 1), 7, 0},
		{"F_LIN2", F_LIN2(5, 0, 0, 10, 20), 10, 0},
		{"F_POLY", F_POLY(2, [8]iec.REAL{1, 1, 1}), 7, 0},
		{"F_POWER", F_POWER(2, 3, 2), 18, 1e-5},
		{"F_QUAD", F_QUAD(2, 1, 1, 1), 7, 0},
	}
	for _, tt := range tests {
		if !near(tt.got, tt.want, tt.eps) {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestIntegerFunctions(t *testing.T) {
	tests := []struct {
		name      string
		got, want int64
	}{
		{"BINOM(5,2)", int64(BINOM(5, 2)), 10},
		{"BINOM(10,3)", int64(BINOM(10, 3)), 120},
		{"BINOM(49,6)", int64(BINOM(49, 6)), 13983816},
		{"BINOM(5,0)", int64(BINOM(5, 0)), 1},
		{"BINOM(5,5)", int64(BINOM(5, 5)), 1},
		{"CEIL(3.14)", int64(CEIL(3.14)), 4},
		{"CEIL(-3.14)", int64(CEIL(-3.14)), -3},
		{"CEIL(3)", int64(CEIL(3)), 3},
		{"CEIL2(2.5)", int64(CEIL2(2.5)), 3},
		{"FLOOR(3.14)", int64(FLOOR(3.14)), 3},
		{"FLOOR(-3.14)", int64(FLOOR(-3.14)), -4},
		{"FLOOR2(-2.5)", int64(FLOOR2(-2.5)), -3},
		{"D_TRUNC(1.5)", int64(D_TRUNC(1.5)), 1},
		{"D_TRUNC(-1.5)", int64(D_TRUNC(-1.5)), -1},
		{"D_TRUNC(2.7)", int64(D_TRUNC(2.7)), 2},
		{"DEC1(0,3)", int64(DEC1(0, 3)), 2},
		{"DEC1(2,3)", int64(DEC1(2, 3)), 1},
		{"FACT(5)", int64(FACT(5)), 120},
		{"FACT(12)", int64(FACT(12)), 479001600},
		{"FACT(13)", int64(FACT(13)), -1},
		{"FIB(0)", int64(FIB(0)), 0},
		{"FIB(1)", int64(FIB(1)), 1},
		{"FIB(3)", int64(FIB(3)), 2},
		{"FIB(4)", int64(FIB(4)), 3},
		{"FIB(10)", int64(FIB(10)), 55},
		{"FIB(46)", int64(FIB(46)), 1836311903},
		{"FIB(47)", int64(FIB(47)), -1},
		{"GCD(12,18)", int64(GCD(12, 18)), 6},
		{"GCD(0,-5)", int64(GCD(0, -5)), 5},
		{"GCD(-48,36)", int64(GCD(-48, 36)), 12},
		{"GCD(17,5)", int64(GCD(17, 5)), 1},
		{"INC(9,1,9)", int64(INC(9, 1, 9)), 0},
		{"INC(3,2,9)", int64(INC(3, 2, 9)), 5},
		{"INC(0,-1,9)", int64(INC(0, -1, 9)), 9},
		{"INC1(2,3)", int64(INC1(2, 3)), 0},
		{"INC1(1,3)", int64(INC1(1, 3)), 2},
		{"INC2(5,1,3,5)", int64(INC2(5, 1, 3, 5)), 3},
		{"INC2(3,-1,3,5)", int64(INC2(3, -1, 3, 5)), 5},
		{"SGN(-2)", int64(SGN(-2)), -1},
		{"SGN(0)", int64(SGN(0)), 0},
		{"SGN(2)", int64(SGN(2)), 1},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
	bools := []struct {
		name      string
		got, want iec.BOOL
	}{
		{"CMP 6 digits", CMP(3.141516, 3.141517, 6), true},
		{"CMP 7 digits", CMP(3.141516, 3.141518, 7), false},
		{"DIFFER", DIFFER(1, 2, 0.5), true},
		{"EVEN(4)", EVEN(4), true},
		{"EVEN(-3)", EVEN(-3), false},
		{"SIGN_I(-1)", SIGN_I(-1), true},
		{"SIGN_R(1)", SIGN_R(1), false},
		{"WINDOW(1,1,2)", WINDOW(1, 1, 2), false},
		{"WINDOW2(1,1,2)", WINDOW2(1, 1, 2), true},
	}
	for _, tt := range bools {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestREAL_TO_FRAC(t *testing.T) {
	tests := []struct {
		x        iec.REAL
		n        iec.INT
		num, den iec.INT
	}{
		{0.75, 100, 3, 4},
		{3.14159265, 1000, 355, 113},
		{-0.5, 10, -1, 2},
		{2, 10, 2, 1},
		{0.333333, 10, 1, 3},
	}
	for _, tt := range tests {
		got := REAL_TO_FRAC(tt.x, tt.n)
		if got.NUMERATOR != tt.num || got.DENOMINATOR != tt.den {
			t.Errorf("REAL_TO_FRAC(%v, %v) = %v/%v, want %v/%v", tt.x, tt.n, got.NUMERATOR, got.DENOMINATOR, tt.num, tt.den)
		}
	}
	if got := REAL_TO_FRAC(iec.REAL(gomath.NaN()), 10); got.DENOMINATOR != 1 {
		t.Errorf("REAL_TO_FRAC(NaN) = %v, want 0/1", got)
	}
}

func TestRandom(t *testing.T) {
	for i := 0; i < 100; i++ {
		if r := RDM(iec.REAL(i) / 100); r < 0 || r >= 1 {
			t.Fatalf("RDM = %v, want 0 <= RDM < 1", r)
		}
		if r := RDM2(iec.INT(i), 3, 7); r < 3 || r > 7 {
			t.Fatalf("RDM2 = %v, want 3..7", r)
		}
	}
	_ = RDMDW(12345)
}

func TestComplex(t *testing.T) {
	c := func(re, im iec.REAL) COMPLEX { return COMPLEX{RE: re, IM: im} }
	eq := func(a, b COMPLEX) bool { return near(a.RE, b.RE, 1e-5) && near(a.IM, b.IM, 1e-5) }
	pi := MATH.PI
	tests := []struct {
		name      string
		got, want COMPLEX
	}{
		{"CADD", CADD(c(1, 2), c(3, 4)), c(4, 6)},
		{"CSUB", CSUB(c(1, 2), c(3, 4)), c(-2, -2)},
		{"CMUL", CMUL(c(1, 2), c(3, 4)), c(-5, 10)},
		{"CDIV", CDIV(c(-5, 10), c(3, 4)), c(1, 2)},
		{"CINV", CINV(c(0, 2)), c(0, -0.5)},
		{"CCON", CCON(c(1, 2)), c(1, -2)},
		{"CSQRT", CSQRT(c(3, 4)), c(2, 1)},
		{"CEXP", CEXP(c(0, pi)), c(-1, 0)},
		{"CLOG", CLOG(c(-1, 0)), c(0, pi)},
		{"CPOW", CPOW(c(0, 1), c(2, 0)), c(-1, 0)},
		{"CPOL", CPOL(2, pi/2), c(0, 2)},
		{"CSET", CSET(1, 2), c(1, 2)},
		{"CSIN", CSIN(c(1, 1)), c(1.2984576, 0.6349639)},
		{"CCOS", CCOS(c(1, 1)), c(0.8337300, -0.9888977)},
		{"CTAN", CTAN(c(1, 1)), c(0.2717526, 1.0839234)},
		{"CSINH", CSINH(c(1, 1)), c(0.6349639, 1.2984576)},
		{"CCOSH", CCOSH(c(1, 1)), c(0.8337300, 0.9888977)},
		{"CTANH", CTANH(c(1, 1)), c(1.0839234, 0.2717526)},
		{"CASIN", CASIN(c(1, 1)), c(0.6662394, 1.0612751)},
		{"CACOS", CACOS(c(1, 1)), c(0.9045569, -1.0612751)},
		{"CATAN", CATAN(c(0.5, 0.5)), c(0.5535744, 0.4023595)},
		{"CASINH", CASINH(c(1, 1)), c(1.0612751, 0.6662394)},
		{"CACOSH", CACOSH(c(1, 1)), c(1.0612751, 0.9045569)},
		{"CATANH", CATANH(c(0.5, 0.5)), c(0.4023595, 0.5535744)},
	}
	for _, tt := range tests {
		if !eq(tt.got, tt.want) {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
	if got := CABS(c(3, 4)); got != 5 {
		t.Errorf("CABS = %v, want 5", got)
	}
	if got := CARG(c(0, 1)); !near(got, pi/2, 1e-6) {
		t.Errorf("CARG = %v, want pi/2", got)
	}
}

func TestDouble(t *testing.T) {
	x := R2_SET(1e7)
	for i := 0; i < 10; i++ {
		x = R2_ADD(x, 0.1)
	}
	if got := float64(x.RX) + float64(x.R1); gomath.Abs(got-10000001) > 1e-3 {
		t.Errorf("R2_ADD sum = %v, want 10000001", got)
	}
	if got := R2_ABS(REAL2{RX: -2, R1: -0.5}); got.RX != 2 || got.R1 != 0.5 {
		t.Errorf("R2_ABS = %v", got)
	}
	if got := R2_MUL(REAL2{RX: 2, R1: 0.5}, 2); got.RX != 4 || got.R1 != 1 {
		t.Errorf("R2_MUL = %v", got)
	}
	if got := R2_ADD2(REAL2{RX: 2, R1: 0.5}, REAL2{RX: 1, R1: 0.25}); got.RX != 3 || got.R1 != 0.75 {
		t.Errorf("R2_ADD2 = %v", got)
	}
}

func TestVector(t *testing.T) {
	a := VECTOR_3{X: 1}
	b := VECTOR_3{Y: 2}
	if got := V3_XPRO(a, b); got != (VECTOR_3{Z: 2}) {
		t.Errorf("V3_XPRO = %v", got)
	}
	if got := V3_ANG(a, b); !near(got, MATH.PI05, 1e-6) {
		t.Errorf("V3_ANG = %v", got)
	}
	if got := V3_NORM(VECTOR_3{X: 3, Y: 4}); !near(got.X, 0.6, 1e-6) || !near(got.Y, 0.8, 1e-6) {
		t.Errorf("V3_NORM = %v", got)
	}
	if got := V3_ABS(VECTOR_3{X: 2, Y: 3, Z: 6}); got != 7 {
		t.Errorf("V3_ABS = %v", got)
	}
	if !V3_PAR(a, V3_SMUL(a, 3)) || V3_PAR(a, b) {
		t.Error("V3_PAR")
	}
	if !V3_NUL(V3_SUB(a, a)) || V3_NUL(a) {
		t.Error("V3_NUL")
	}
	if V3_DPRO(V3_ADD(a, b), V3_REV(b)) != -4 {
		t.Error("V3_DPRO")
	}
	if !near(V3_XANG(b), MATH.PI05, 1e-6) || V3_YANG(b) != 0 || !near(V3_ZANG(a), MATH.PI05, 1e-6) {
		t.Error("V3_XANG, V3_YANG, V3_ZANG")
	}
}

func TestArrays(t *testing.T) {
	size := func(a []iec.REAL) iec.UINT { return iec.UINT(len(a) * 4) }

	a := []iec.REAL{12, 0, 4, 7, 1}
	if got := ARRAY_MEDIAN_(a, size(a)); got != 4 {
		t.Errorf("ARRAY_MEDIAN_ = %v, want 4", got)
	}
	if !IS_SORTED(a, size(a)) || a[0] != 0 || a[4] != 12 {
		t.Errorf("array after ARRAY_MEDIAN_ = %v, want it sorted", a)
	}
	if got := ARRAY_MEDIAN_([]iec.REAL{4, 1, 3, 2}, 16); got != 2.5 {
		t.Errorf("ARRAY_MEDIAN_ of 4 = %v, want 2.5", got)
	}
	if IS_SORTED([]iec.REAL{1, 3, 2}, 12) {
		t.Error("IS_SORTED([1 3 2]) = true")
	}

	trend := []iec.REAL{0, 1, 4, 5, 3, 4, 6, 3}
	if got := ARRAY_TREND(trend, size(trend)); got != 1.5 {
		t.Errorf("ARRAY_TREND = %v, want 1.5", got)
	}

	b := []iec.REAL{1, 2, 3, 4}
	tests := []struct {
		name      string
		got, want iec.REAL
	}{
		{"ARRAY_AVG", ARRAY_AVG(b, 16), 2.5},
		{"ARRAY_SUM", ARRAY_SUM(b, 16), 10},
		{"ARRAY_MAX", ARRAY_MAX(b, 16), 4},
		{"ARRAY_MIN", ARRAY_MIN(b, 16), 1},
		{"ARRAY_SPR", ARRAY_SPR(b, 16), 3},
		{"ARRAY_VAR", ARRAY_VAR(b, 16), 5.0 / 3},
		{"ARRAY_SDV", ARRAY_SDV(b, 16), SQRT(5.0 / 3)},
		{"ARRAY_GAV", ARRAY_GAV([]iec.REAL{1, 4}, 8), 2},
		{"ARRAY_HAV", ARRAY_HAV([]iec.REAL{1, 2}, 8), 4.0 / 3},
		{"ARRAY_GAV with 0", ARRAY_GAV([]iec.REAL{1, 0}, 8), 0},
		{"ARRAY_SUM of 2", ARRAY_SUM(b, 8), 3},
	}
	for _, tt := range tests {
		if !near(tt.got, tt.want, 1e-5) {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}

	ARRAY_ADD_(b, 16, 1)
	ARRAY_MUL_(b, 16, 2)
	if b[0] != 4 || b[3] != 10 {
		t.Errorf("ARRAY_ADD_, ARRAY_MUL_ = %v", b)
	}
	ARRAY_INIT_(b, 8, -1)
	ARRAY_ABS_(b, 16)
	if b[0] != 1 || b[1] != 1 || b[2] != 8 {
		t.Errorf("ARRAY_INIT_, ARRAY_ABS_ = %v", b)
	}
	c := []iec.REAL{1, 2, 3, 4, 5, 6}
	ARRAY_SHUFFLE_(c, size(c))
	ARRAY_SORT_(c, size(c))
	for i, v := range c {
		if v != iec.REAL(i+1) {
			t.Fatalf("ARRAY_SHUFFLE_ lost elements: %v", c)
		}
	}
}

func TestInterpolation(t *testing.T) {
	var xy [20][2]iec.REAL
	xy[0], xy[1], xy[2] = [2]iec.REAL{0, 0}, [2]iec.REAL{10, 100}, [2]iec.REAL{20, 0}
	for _, tt := range []struct{ x, want iec.REAL }{{5, 50}, {15, 50}, {25, -50}, {-5, -50}} {
		if got := LINEAR_INT(tt.x, xy, 3); got != tt.want {
			t.Errorf("LINEAR_INT(%v) = %v, want %v", tt.x, got, tt.want)
		}
	}
	var p [5][2]iec.REAL
	p[0], p[1], p[2] = [2]iec.REAL{0, 0}, [2]iec.REAL{1, 1}, [2]iec.REAL{2, 4}
	if got := POLYNOM_INT(3, p, 3); !near(got, 9, 1e-5) {
		t.Errorf("POLYNOM_INT(3) = %v, want 9", got)
	}
}

func TestFRMP_B(t *testing.T) {
	ms := func(n int) iec.TIME { return iec.TIME(time.Duration(n) * time.Millisecond) }
	tests := []struct {
		start  iec.BYTE
		dir    iec.BOOL
		td, tr iec.TIME
		want   iec.BYTE
	}{
		{0, true, ms(500), ms(1000), 128},
		{200, true, ms(500), ms(1000), 255},
		{100, false, ms(500), ms(1000), 0},
		{200, false, ms(250), ms(1000), 136},
		{7, true, ms(1000), ms(1000), 255},
		{7, false, ms(2000), ms(1000), 0},
	}
	for _, tt := range tests {
		if got := FRMP_B(tt.start, tt.dir, tt.td, tt.tr); got != tt.want {
			t.Errorf("FRMP_B(%v, %v, %v, %v) = %v, want %v", tt.start, tt.dir, tt.td, tt.tr, got, tt.want)
		}
	}
}

func TestFT_RMP(t *testing.T) {
	now := time.Unix(1000, 0)
	var r FT_RMP
	r.INIT()
	r.IN, r.KR, r.KF = 10, 1, 2
	r.Execute(now)
	if r.OUT != 10 || r.BUSY {
		t.Fatalf("first scan: OUT = %v, BUSY = %v, want 10, false", r.OUT, r.BUSY)
	}
	r.IN = 20
	now = now.Add(2 * time.Second)
	r.Execute(now)
	if !near(r.OUT, 12, 1e-4) || !bool(r.BUSY && r.UD) {
		t.Fatalf("ramp up: OUT = %v, BUSY = %v, UD = %v, want 12, true, true", r.OUT, r.BUSY, r.UD)
	}
	r.IN = 0
	now = now.Add(time.Second)
	r.Execute(now)
	if !near(r.OUT, 10, 1e-4) || !bool(r.BUSY && !r.UD) {
		t.Fatalf("ramp down: OUT = %v, BUSY = %v, UD = %v, want 10, true, false", r.OUT, r.BUSY, r.UD)
	}
	now = now.Add(time.Minute)
	r.Execute(now)
	if r.OUT != 0 || r.BUSY {
		t.Fatalf("end: OUT = %v, BUSY = %v, want 0, false", r.OUT, r.BUSY)
	}
	r.RMP = false
	r.IN = 5
	r.Execute(now)
	if r.OUT != 5 {
		t.Fatalf("RMP false: OUT = %v, want 5", r.OUT)
	}
}

func TestFT_AVG(t *testing.T) {
	var f FT_AVG
	f.INIT()
	f.N = 4
	f.Execute(time.Time{})
	if f.AVG != 0 || !f.E {
		t.Fatalf("first scan: AVG = %v, E = %v", f.AVG, f.E)
	}
	f.IN = 4
	want := []iec.REAL{1, 2, 3, 4, 4}
	for i, w := range want {
		f.Execute(time.Time{})
		if !near(f.AVG, w, 1e-5) {
			t.Fatalf("scan %d: AVG = %v, want %v", i, f.AVG, w)
		}
	}
	f.RST, f.IN = true, 8
	f.Execute(time.Time{})
	if f.AVG != 8 {
		t.Fatalf("RST: AVG = %v, want 8", f.AVG)
	}
}

func TestFT_MIN_MAX(t *testing.T) {
	var f FT_MIN_MAX
	for _, in := range []iec.REAL{5, 3, 9, 7} {
		f.IN = in
		f.Execute(time.Time{})
	}
	if f.MN != 3 || f.MX != 9 {
		t.Fatalf("MN, MX = %v, %v, want 3, 9", f.MN, f.MX)
	}
	f.RST, f.IN = true, 6
	f.Execute(time.Time{})
	if f.MN != 6 || f.MX != 6 {
		t.Fatalf("after RST MN, MX = %v, %v, want 6, 6", f.MN, f.MX)
	}
}
