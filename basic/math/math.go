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

// Package math is the port of the OSCAT BASIC mathematical functions.
package math

import (
	gomath "math"
	"math/bits"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// ACOSH calculates the arcus cosinus hyperbolicus.
func ACOSH(x iec.REAL) iec.REAL {
	return LN(SQRT(x*x-1.0) + x)
}

// ACOTH calculates the arcus cotangens hyperbolicus.
func ACOTH(x iec.REAL) iec.REAL {
	return LN((x+1.0)/(x-1.0)) * 0.5
}

// AGDF calculates the inverse Gudermannian function. It is accurate only up
// to x = 1.57.
func AGDF(x iec.REAL) iec.REAL {
	return LN((1.0 + SIN(x)) / COS(x))
}

// ASINH calculates the arcus sinus hyperbolicus.
func ASINH(x iec.REAL) iec.REAL {
	return LN(SQRT(x*x+1.0) + x)
}

// ATAN2 calculates the angle of the point x, y in a coordinate system, in
// rad.
func ATAN2(y, x iec.REAL) iec.REAL {
	switch {
	case x > 0.0:
		return ATAN(y / x)
	case x < 0.0:
		if y >= 0.0 {
			return ATAN(y/x) + MATH.PI
		}
		return ATAN(y/x) - MATH.PI
	case y > 0.0:
		return MATH.PI05
	case y < 0.0:
		return -MATH.PI05
	}
	return 0.0
}

// ATANH calculates the arcus tangens hyperbolicus.
func ATANH(x iec.REAL) iec.REAL {
	return LN((1.0+x)/(1.0-x)) * 0.5
}

// BETA calculates the beta function for real numbers > 0.
func BETA(x, y iec.REAL) iec.REAL {
	return GAMMA(x) * GAMMA(y) / GAMMA(x+y)
}

// BINOM calculates the binomial coefficient, n over k.
func BINOM(n, k iec.INT) iec.DINT {
	if 2*k > n {
		k = n - k
	}
	switch {
	case k > n:
		return 0
	case k == 0 || k == n:
		return 1
	case k == 1:
		return iec.DINT(n)
	}
	binom := iec.DINT(n)
	n++
	for i := iec.INT(2); i <= k; i++ {
		binom = binom * iec.DINT(n-i) / iec.DINT(i)
	}
	return binom
}

// CAUCHY calculates the Cauchy density function.
func CAUCHY(x, t, u iec.REAL) iec.REAL {
	tmp := x - t
	return MATH.PI_INV * u / (u*u + tmp*tmp)
}

// CAUCHYCD calculates the Cauchy cumulative distribution function.
func CAUCHYCD(x, t, u iec.REAL) iec.REAL {
	return 0.5 + MATH.PI_INV*ATAN((x-t)/u)
}

// CEIL returns the smallest integer that is greater than or equal to x:
// CEIL(3.14) = 4, CEIL(-3.14) = -3.
func CEIL(x iec.REAL) iec.INT {
	ceil := REAL_TO_INT(x)
	if iec.REAL(ceil) < x {
		ceil++
	}
	return ceil
}

// CEIL2 returns the smallest integer that is greater than or equal to x, as
// a DINT.
func CEIL2(x iec.REAL) iec.DINT {
	ceil := REAL_TO_DINT(x)
	if iec.REAL(ceil) < x {
		ceil++
	}
	return ceil
}

// CMP reports whether x and y are identical in their first n digits:
// CMP(3.141516, 3.141517, 6) is true.
func CMP(x, y iec.REAL, n iec.INT) iec.BOOL {
	tmp := ABS(x)
	if tmp > 0.0 {
		tmp = EXP10(iec.REAL(FLOOR(LOG(tmp)) - n + 1))
	} else {
		tmp = EXP10(tmp)
	}
	return ABS(x-y) < tmp
}

// COSH calculates the cosinus hyperbolicus.
func COSH(x iec.REAL) iec.REAL {
	t := EXP(x)
	return (1.0/t + t) * 0.5
}

// COTH calculates the cotangens hyperbolicus.
func COTH(x iec.REAL) iec.REAL {
	switch {
	case x > 20.0:
		return 1.0
	case x < -20.0:
		return -1.0
	}
	return 1.0 + 2.0/(EXP(x*2.0)-1.0)
}

// D_TRUNC truncates a REAL to a DINT: 1.5 is 1 and -1.5 is -1.
func D_TRUNC(x iec.REAL) iec.DINT {
	t := REAL_TO_DINT(x)
	if x > 0.0 {
		if iec.REAL(t) > x {
			t--
		}
	} else if iec.REAL(t) < x {
		t++
	}
	return t
}

// DEC1 decrements x by 1 and starts again at n-1 after 0: DEC1(x, 3)
// generates 2, 1, 0, 2, ...
func DEC1(x, n iec.INT) iec.INT {
	if x == 0 {
		return n - 1
	}
	return x - 1
}

// DEG converts an angle in rad to degrees, 0..360.
func DEG(rad iec.REAL) iec.REAL {
	return MODR(57.29577951308232*rad, 360.0)
}

// DIFFER reports whether in1 differs by more than x from in2.
func DIFFER(in1, in2, x iec.REAL) iec.BOOL {
	return ABS(in1-in2) > x
}

// ERF calculates the error function.
func ERF(x iec.REAL) iec.REAL {
	x2 := x * x
	ax2 := 0.147*x2 + 1.0
	return SQRT(1.0-EXP(-x2*((0.27323954473516+ax2)/ax2))) * iec.REAL(SGN(x))
}

// ERFC calculates the complementary error function.
func ERFC(x iec.REAL) iec.REAL {
	return 1.0 - ERF(x)
}

// EVEN reports whether in is even.
func EVEN(in iec.DINT) iec.BOOL {
	return in&1 == 0
}

// EXP10 calculates 10 to the power of x.
func EXP10(x iec.REAL) iec.REAL {
	return EXP(x * 2.30258509299405)
}

// EXPN calculates x to the power of the integer n.
func EXPN(x iec.REAL, n iec.INT) iec.REAL {
	sign := n < 0
	// The absolute value, as an unsigned number so -32768 is kept.
	m := uint16(n)
	if sign {
		m = -m
	}
	var expn iec.REAL = 1.0
	if m&1 != 0 {
		expn = x
	}
	m >>= 1
	for m > 0 {
		x = x * x
		if m&1 != 0 {
			expn = expn * x
		}
		m >>= 1
	}
	if sign {
		expn = 1.0 / expn
	}
	return expn
}

// FACT calculates the factorial of x. It is -1 if x is negative or higher
// than 12.
func FACT(x iec.INT) iec.DINT {
	if x >= 0 && x <= 12 {
		return MATH.FACTS[x]
	}
	return -1
}

// FIB calculates the Fibonacci number of x. It is -1 if x is negative or
// higher than 46.
func FIB(x iec.INT) iec.DINT {
	if x < 0 || x > 46 {
		return -1
	}
	if x < 2 {
		return iec.DINT(x)
	}
	var t1, t2 iec.DINT = 0, 1
	for x > 3 {
		x -= 2
		t1 = t1 + t2
		t2 = t1 + t2
	}
	if x > 2 {
		t1 = t1 + t2
	}
	return t1 + t2
}

// FLOOR returns the biggest integer that is less than or equal to x:
// FLOOR(3.14) = 3, FLOOR(-3.14) = -4.
func FLOOR(x iec.REAL) iec.INT {
	floor := REAL_TO_INT(x)
	if iec.REAL(floor) > x {
		floor--
	}
	return floor
}

// FLOOR2 returns the biggest integer that is less than or equal to x, as a
// DINT.
func FLOOR2(x iec.REAL) iec.DINT {
	floor := REAL_TO_DINT(x)
	if iec.REAL(floor) > x {
		floor--
	}
	return floor
}

// FRACT returns the fraction of a real number: FRACT(3.14) = 0.14.
func FRACT(x iec.REAL) iec.REAL {
	if ABS(x) < 2.0e9 {
		return x - iec.REAL(D_TRUNC(x))
	}
	return 0.0
}

// GAMMA calculates an approximation of the gamma function for x > 0, with
// the Stirling formula.
func GAMMA(x iec.REAL) iec.REAL {
	if x > 0.0 {
		return SQRT(MATH.PI2/x) * EXPT(MATH.E_INV*(x+1.0/(12.0*x-0.1/x)), x)
	}
	return 0.0
}

// GAUSS calculates the Gaussian density function.
func GAUSS(x, u, si iec.REAL) iec.REAL {
	temp := x - u
	siInv := 1.0 / si
	return EXP(temp*temp*siInv*siInv*-0.5) * 0.39894228 * siInv
}

// GAUSSCD calculates the Gaussian cumulative distribution function.
func GAUSSCD(x, u, si iec.REAL) iec.REAL {
	return (ERF((x-u)/(si*1.414213562)) + 1.0) * 0.5
}

// GCD calculates the greatest common divisor of a and b.
func GCD(a, b iec.DINT) iec.INT {
	abs := func(v iec.DINT) iec.DINT {
		if v < 0 {
			return -v
		}
		return v
	}
	if a == 0 {
		return iec.INT(abs(b))
	}
	if b == 0 {
		return iec.INT(abs(a))
	}
	a, b = abs(a), abs(b)
	var gcd iec.INT = 1
	for a&1 == 0 && b&1 == 0 {
		a >>= 1
		b >>= 1
		gcd <<= 1
	}
	for a > 0 {
		switch {
		case a&1 == 0:
			a >>= 1
		case b&1 == 0:
			b >>= 1
		default:
			t := abs(a-b) >> 1
			if a < b {
				b = t
			} else {
				a = t
			}
		}
	}
	return gcd * iec.INT(b)
}

// GDF calculates the Gudermannian function.
func GDF(x iec.REAL) iec.REAL {
	switch {
	case x == 0.0:
		return 0.0
	case x > 15.0:
		return MATH.PI05
	case x < -15.0:
		return -MATH.PI05
	}
	return ATAN(EXP(x))*2.0 - MATH.PI05
}

// GOLD calculates the golden function.
func GOLD(x iec.REAL) iec.REAL {
	return (x + SQRT(x*x+4.0)) * 0.5
}

// HYPOT calculates the Pythagorean function, the length of the hypotenuse.
func HYPOT(x, y iec.REAL) iec.REAL {
	return SQRT(x*x + y*y)
}

// INC increments x by d and continues from 0 when the result exceeds m.
func INC(x, d, m iec.INT) iec.INT {
	return (x + d + m + 1) % (m + 1)
}

// INC1 increments x by 1 and starts again at 0 when n is reached: INC1(x, 3)
// generates 0, 1, 2, 0, 1, ...
func INC1(x, n iec.INT) iec.INT {
	if x >= n-1 {
		return 0
	}
	return x + 1
}

// INC2 increments x by d and continues from l when the result exceeds u.
func INC2(x, d, l, u iec.INT) iec.INT {
	tmp := u - l + 1
	return (x+d-l+tmp)%tmp + l
}

// INV calculates 1 / x. It is 0 for x = 0.
func INV(x iec.REAL) iec.REAL {
	if x != 0.0 {
		return 1.0 / x
	}
	return 0.0
}

// LAMBERT_W calculates the Lambert W function. It is -1000 if x is below
// -1/e.
func LAMBERT_W(x iec.REAL) iec.REAL {
	var w iec.REAL
	switch {
	case x < -0.367879441171442:
		return -1000.0
	case x == 0.0:
		return 0.0
	case x <= 500.0:
		// First an estimate is calculated.
		w = LN(x + 1.0)
		w = 0.665*(1.0+0.0195*w)*w + 0.04
	default:
		w = LN(x-4.0) - (1.0-1.0/LN(x))*LN(LN(x))
	}
	// Use the estimate to calculate the exact result.
	for i := 0; i <= 5; i++ {
		ewx := EXP(w)
		we := w*ewx - x
		w1e := (w + 1.0) * ewx
		last := gomath.Float32bits(float32(w)) & 0xFFFF_FFFC
		w = w - (we / (w1e - (w+2.0)*we/(2.0*w+2.0)))
		if gomath.Float32bits(float32(w))&0xFFFF_FFFC == last {
			break
		}
	}
	return w
}

// LANGEVIN calculates the Langevin function.
func LANGEVIN(x iec.REAL) iec.REAL {
	if x == 0.0 {
		return 0.0
	}
	return COTH(x) - 1.0/x
}

// MAX3 returns the highest of three values.
func MAX3(in1, in2, in3 iec.REAL) iec.REAL {
	return max(in1, in2, in3)
}

// MID3 returns the middle of three values.
func MID3(in1, in2, in3 iec.REAL) iec.REAL {
	if in1 > in2 {
		in1, in2 = in2, in1
	}
	if in2 > in3 {
		in2 = in3
	}
	return SEL(in1 > in2, in2, in1)
}

// MIN3 returns the lowest of three values.
func MIN3(in1, in2, in3 iec.REAL) iec.REAL {
	return min(in1, in2, in3)
}

// MODR is the modulo function for real numbers: MODR(5.5, 2.5) = 0.5. It
// is 0 for a divisor of 0.
func MODR(in, divi iec.REAL) iec.REAL {
	if divi == 0.0 {
		return 0.0
	}
	return in - iec.REAL(FLOOR2(in/divi))*divi
}

// MUL_ADD multiplies x with k and adds the offset o.
func MUL_ADD(x, k, o iec.REAL) iec.REAL {
	return x*k + o
}

// NEGX returns -x.
func NEGX(x iec.REAL) iec.REAL {
	return -x
}

// RAD converts an angle in degrees to rad, 0..2π.
func RAD(deg iec.REAL) iec.REAL {
	return MODR(0.0174532925199433*deg, MATH.PI2)
}

// RDM calculates a pseudo random number between 0 and 1 from the PLC timer.
// To use it more than once in a scan, call it with different seeds last.
func RDM(last iec.REAL) iec.REAL {
	tn := uint32(PLC_MS(Now()))
	tc := bits.OnesCount32(tn)
	// Bits 31..24 are bits 2, 5, 4, 1, 0, 7, 6 and 3.
	for i, src := range [8]uint{2, 5, 4, 1, 0, 7, 6, 3} {
		dst := uint(31 - i)
		tn = tn&^(1<<dst) | (tn>>src&1)<<dst
	}
	tn = bits.RotateLeft32(tn, bits.OnesCount32(tn)) | 0x80000001
	tn = tn%71474513 + uint32(tc+77)
	return FRACT(iec.REAL(tn) / 10000000.0 * (MATH.E - LIMIT(0.0, last, 1.0)))
}

// RDM2 calculates an integer pseudo random number in the range low..high.
func RDM2(last, low, high iec.INT) iec.INT {
	return iec.INT(TRUNC(RDM(FRACT(iec.REAL(last)*MATH.PI))*iec.REAL(high-low+1))) + low
}

// RDMDW calculates a DWORD pseudo random number.
func RDMDW(last iec.DWORD) iec.DWORD {
	m := iec.REAL(bits.OnesCount32(uint32(last)))
	rx := RDM(FRACT(m * MATH.PI))
	out := REAL_TO_DWORD(rx*65535) << 16
	rx = RDM(FRACT(m * MATH.E))
	return out | REAL_TO_DWORD(rx*65535)&0x0000FFFF
}

// REAL_TO_FRAC calculates the closest fraction for a real number, with a
// denominator of at most n.
func REAL_TO_FRAC(x iec.REAL, n iec.INT) FRACTION {
	abs := func(v iec.DINT) iec.DINT {
		if v < 0 {
			return -v
		}
		return v
	}
	var numerator, denominator iec.DINT = 1, 0
	var numeratorOld, denominatorOld iec.DINT = 0, 1

	sign := x < 0.0
	if sign {
		x = ABS(x)
	}
	// x != x is true for a value that is not a number, which has no
	// fraction and would never end the loop.
	for x == x {
		rounded := REAL_TO_DINT(x)

		// Extend the numerator and the denominator by the integer part.
		numeratorOld, numerator = numerator, numerator*rounded+numeratorOld
		denominatorOld, denominator = denominator, denominator*rounded+denominatorOld

		whole := iec.REAL(rounded)
		if x == whole {
			// The fraction has no remainder.
			if abs(denominator) <= iec.DINT(n) {
				numeratorOld = numerator
				denominatorOld = denominator
			}
			break
		}
		// The reciprocal of the remainder is the next fraction.
		x = 1.0 / (x - whole)
		if abs(denominator) > iec.DINT(n) {
			break
		}
	}

	num := iec.INT(numeratorOld)
	if num < 0 {
		num = -num
	}
	den := iec.INT(denominatorOld)
	if den < 0 {
		den = -den
	}
	if sign {
		num = -num
	}
	return FRACTION{NUMERATOR: num, DENOMINATOR: den}
}

// RND rounds a real to n digits in total: RND(3.1415, 2) = 3.1.
func RND(x iec.REAL, n iec.INT) iec.REAL {
	if x == 0.0 {
		return 0.0
	}
	m := EXPN(10.0, n-CEIL(LOG(ABS(x))))
	return iec.REAL(REAL_TO_DINT(x*m)) / m
}

// ROUND rounds a real to n digits behind the decimal point, 0..8.
func ROUND(in iec.REAL, n iec.INT) iec.REAL {
	x := SETUP.DECADES[LIMIT(0, n, 8)]
	return iec.REAL(REAL_TO_DINT(in*x)) / x
}

// SGN returns 0 when x = 0, -1 when x < 0 and +1 when x > 0.
func SGN(x iec.REAL) iec.INT {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	}
	return 0
}

// SIGMOID calculates the sigmoid function.
func SIGMOID(x iec.REAL) iec.REAL {
	switch {
	case x > 20.0:
		return 1.0
	case x < -85.0:
		return 0.0
	}
	return 1.0 / (1.0 + EXP(-x))
}

// SIGN_I reports whether an integer is negative.
func SIGN_I(in iec.DINT) iec.BOOL {
	return in < 0
}

// SIGN_R reports whether a real is negative.
func SIGN_R(in iec.REAL) iec.BOOL {
	return in < 0.0
}

// SINC calculates the sinc function, sin(x) / x.
func SINC(x iec.REAL) iec.REAL {
	if x == 0.0 {
		return 1.0
	}
	return SIN(x) / x
}

// SINH calculates the sinus hyperbolicus.
func SINH(x iec.REAL) iec.REAL {
	if ABS(x) < 2e-3 {
		return x
	}
	return (EXP(x) - EXP(-x)) * 0.5
}

// SQRTN calculates the nth root of x, x^(1/n). It is 0 for n <= 0.
func SQRTN(x iec.REAL, n iec.INT) iec.REAL {
	if n > 0 {
		return EXP(LN(x) / iec.REAL(n))
	}
	return 0.0
}

// TANC calculates the tanc function, tan(x) / x.
func TANC(x iec.REAL) iec.REAL {
	if x == 0.0 {
		return 1.0
	}
	return TAN(x) / x
}

// TANH calculates the tangens hyperbolicus.
func TANH(x iec.REAL) iec.REAL {
	switch {
	case x > 20.0:
		return 1.0
	case x < -20.0:
		return -1.0
	}
	return 1.0 - 2.0/(EXP(2.0*x)+1.0)
}

// WINDOW reports whether in is between low and high.
func WINDOW(low, in, high iec.REAL) iec.BOOL {
	return in < high && in > low
}

// WINDOW2 reports whether in is between low and high, including the limits.
func WINDOW2(low, in, high iec.REAL) iec.BOOL {
	return in >= low && in <= high
}
