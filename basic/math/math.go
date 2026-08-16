/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package math

import (
	. "beebread/basic"
	"math"
)

// ACOSH calculates the hyperbolic arccos of X.
func ACOSH(x float64) float64 {
	return math.Log(x + math.Sqrt(x*x-1.0))
}

// ACOTH calculates the hyperbolic arccot of X.
func ACOTH(x float64) float64 {
	return 0.5 * math.Log((x+1.0)/(x-1.0))
}

// AGDF calculates the inverse Gudermann function of X.
func AGDF(x float64) float64 {
	return math.Log(math.Tan(Math.Pi025 + 0.5*x))
}

// ASINH calculates the hyperbolic arcsin of X.
func ASINH(x float64) float64 {
	return math.Log(x + math.Sqrt(x*x+1.0))
}

// ATAN2 calculates the arctan of Y/X.
func ATAN2(y, x float64) float64 {
	return math.Atan2(y, x)
}

// ATANH calculates the hyperbolic arctan of X.
func ATANH(x float64) float64 {
	return 0.5 * math.Log((1.0+x)/(1.0-x))
}

// BETA calculates the beta function of X and Y.
func BETA(x, y float64) float64 {
	// math.Gamma is the standard library implementation for the Gamma function.
	return (math.Gamma(x) * math.Gamma(y)) / math.Gamma(x+y)
}

// BINOM calculates the binomial coefficient of N over K.
func BINOM(n, k int) int64 {
	if k >= 0 && k <= n && n < 13 {
		return int64(Math.Facts[n] / (Math.Facts[k] * Math.Facts[n-k]))
	}
	return 0 // Or handle error appropriately
}

// CAUCHY calculates the Cauchy distribution.
func CAUCHY(x, m, g float64) float64 {
	return Math.Pi_inv * g / ((x-m)*(x-m) + g*g)
}

// CAUCHY_CD calculates the cumulative Cauchy distribution.
func CAUCHY_CD(x, m, g float64) float64 {
	return 0.5 + Math.Pi_inv*math.Atan((x-m)/g)
}

// CEIL calculates the next higher integer of X.
func CEIL(x float64) int {
	return int(math.Ceil(x))
}

// CEIL2 calculates the next higher integer of X/Y * Y.
func CEIL2(x, y float64) float64 {
	if y != 0.0 {
		return math.Ceil(x/y) * y
	}
	return 0.0
}

// CMP compares two real values with a given deviation D.
// Returns 0 if |X1-X2| <= D.
// Returns 1 if X1 > X2.
// Returns -1 if X1 < X2.
func CMP(x1, x2, d float64) int {
	if math.Abs(x1-x2) <= d {
		return 0
	} else if x1 > x2 {
		return 1
	}
	return -1
}

// COSH calculates the cosinus hyperbolicus of X.
func COSH(x float64) float64 {
	return math.Cosh(x)
}

// COTH calculates the cotangens hyperbolicus of X.
func COTH(x float64) float64 {
	return (math.Exp(x) + math.Exp(-x)) / (math.Exp(x) - math.Exp(-x))
}

// DIFFER checks if in1 differs more than x from in2.
// It returns true if abs(in1 - in2) > x.
func DIFFER(in1, in2, x float64) bool {
	return math.Abs(in1-in2) > x
}

// D_TRUNC truncates a real value to a dint (int64).
func D_TRUNC(x float64) int64 {
	return int64(math.Trunc(x))
}

// ERF calculates the error function.
func ERF(x float64) float64 {
	return math.Erf(x)
}

// ERFC calculates the complementary error function.
func ERFC(x float64) float64 {
	return math.Erfc(x)
}

// LIMIT restricts a value to a given range.
func LIMIT(min, val, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}

// LIMIT_B limits a byte value to a given range.
func LIMIT_B(min, val, max byte) byte {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// LIMIT_DW limits a dword (uint32) value to a given range.
func LIMIT_DW(min, val, max uint32) uint32 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// LIMIT_I64 limits an int64 value to a given range.
func LIMIT_I64(min, val, max int64) int64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// EVEN checks if an integer is even.
func EVEN(in int64) bool {
	return in%2 == 0
}

// EXP10 calculates the base-10 exponential of X (10^X).
func EXP10(x float64) float64 {
	// This is the direct Go equivalent of the original ST logic EXP(X * 2.302585...)
	return math.Pow(10.0, x)
}

// EXPN calculates X to the power of N, where N is an integer.
// This is an implementation of exponentiation by squaring, which is more
// efficient for integer exponents than the generic math.Pow().
func EXPN(x float64, n int) float64 {
	if n < 0 {
		x = 1.0 / x
		n = -n
	}
	if n == 0 {
		return 1.0
	}

	res := 1.0
	for n > 0 {
		// If n is odd, multiply x with the result
		if n%2 == 1 {
			res *= x
		}
		// n must be even now, and x is squared
		n /= 2
		x *= x
	}
	return res
}

// FACT calculates the factorial of x.
// Returns -1 if the input is negative or > 12.
func FACT(x int) int64 {
	if x >= 0 && x < len(Math.Facts) {
		return int64(Math.Facts[x])
	}
	return -1
}

// FIB calculates the Fibonacci sequence number for X.
// Returns -1 for X < 0 or X > 46.
func FIB(x int) int64 {
	if x < 0 || x > 46 {
		return -1
	}
	if x < 2 {
		return int64(x)
	}
	var t1, t2 int64 = 0, 1
	for x > 2 {
		t1, t2 = t2, t1+t2
		x--
	}
	return t1 + t2
}

// FLOOR returns the greatest integer value less than or equal to x.
func FLOOR(x float64) int {
	return int(math.Floor(x))
}

// FLOOR2 returns the greatest integer value less than or equal to x as a 64-bit integer.
func FLOOR2(x float64) int64 {
	return int64(math.Floor(x))
}

// GAMMA calculates the Gamma function of x.
func GAMMA(x float64) float64 {
	return math.Gamma(x)
}

// GAUSS calculates the Gaussian density function.
func GAUSS(x, u, si float64) float64 {
	if si == 0 {
		return 0
	}
	temp := x - u
	si_inv := 1.0 / si
	return math.Exp(temp*temp*si_inv*si_inv*-0.5) * 0.39894228 * si_inv
}

// GAUSSCD calculates the Gaussian cumulative distribution function.
func GAUSSCD(x, u, si float64) float64 {
	if si == 0 {
		if x < u {
			return 0
		}
		return 1
	}
	return (math.Erf((x-u)/(si*1.414213562)) + 1.0) * 0.5
}

// GCD calculates the greatest common divisor of two numbers.
func GCD(a, b int64) int {
	if a == 0 {
		return int(math.Abs(float64(b)))
	}
	if b == 0 {
		return int(math.Abs(float64(a)))
	}
	a = int64(math.Abs(float64(a)))
	b = int64(math.Abs(float64(b)))
	gcd := 1
	for a&1 == 0 && b&1 == 0 {
		a >>= 1
		b >>= 1
		gcd <<= 1
	}
	for a > 0 {
		if a&1 == 0 {
			a >>= 1
		} else if b&1 == 0 {
			b >>= 1
		} else {
			t := int64(math.Abs(float64(a-b))) >> 1
			if a < b {
				b = t
			} else {
				a = t
			}
		}
	}
	return gcd * int(b)
}

// INC implements a modular increment.
func INC(x, d, m int) int {
	if m <= 0 {
		return x
	}
	return (x + d) % (m + 1)
}

// INC1 increments X by 1 and wraps around to 0 if N is reached.
// It generates a sequence: 0, 1, 2, ..., N-1, 0, ...
func INC1(x, n int) int {
	if x >= n-1 {
		return 0
	}
	return x + 1
}

// INC2 increments X by D and wraps around between L and U.
func INC2(x, d, l, u int) int {
	tmp := u - l + 1
	if tmp <= 0 {
		return x
	}
	res := (x + d - l) % tmp
	if res < 0 {
		res += tmp
	}
	return res + l
}

// GDF calculates the Gudermannian function.
func GDF(x float64) float64 {
	if x == 0.0 {
		return 0.0
	}
	if x > 15.0 {
		return Math.Pi05
	}
	if x < -15.0 {
		return -Math.Pi05
	}
	return math.Atan(math.Exp(x))*2.0 - Math.Pi05
}

// GOLD calculates the golden function.
func GOLD(x float64) float64 {
	return (x + math.Sqrt(x*x+4.0)) * 0.5
}

// INV calculates the inverse of a number (1/x).
func INV(x float64) float64 {
	if x != 0.0 {
		return 1.0 / x
	}
	return 0.0
}

// LAMBERT_W calculates the Lambert W function.
func LAMBERT_W(x float64) float64 {
	if x < -0.367879441171442 {
		return -1000.0
	}
	if x == 0.0 {
		return 0.0
	}

	var w float64
	if x <= 500.0 {
		w = math.Log(x + 1.0)
		w = 0.665*(1.0+0.0195*w)*w + 0.04
	} else {
		logX := math.Log(x)
		w = logX - math.Log(logX)
	}

	for i := 0; i < 10; i++ {
		ew := math.Exp(w)
		we := w*ew - x
		if we == 0 {
			break
		}
		w1e := (w + 1.0) * ew
		w -= we / (w1e - (w+2.0)*we/(2.0*w+2.0))
	}
	return w
}

// LANGEVIN calculates the Langevin function.
func LANGEVIN(x float64) float64 {
	if x == 0.0 {
		return 0.0
	}
	return COTH(x) - 1.0/x
}

// MAX3 returns the maximum of three values.
func MAX3(in1, in2, in3 float64) float64 {
	return math.Max(math.Max(in1, in2), in3)
}

// MIN3 returns the minimum of three values.
func MIN3(in1, in2, in3 float64) float64 {
	return math.Min(math.Min(in1, in2), in3)
}

// MID3 returns the middle value of three real inputs.
func MID3(in1, in2, in3 float64) float64 {
	// This is a more readable Go implementation of the original ST
	if (in1 <= in2 && in2 <= in3) || (in3 <= in2 && in2 <= in1) {
		return in2
	}
	if (in2 <= in1 && in1 <= in3) || (in3 <= in1 && in1 <= in2) {
		return in1
	}
	return in3
}

// MODR is a modulo function for real numbers.
// e.g., ModR(5.5, 2.5) = 0.5
func MODR(in, divi float64) float64 {
	if divi == 0.0 {
		return 0.0
	}
	return in - math.Floor(in/divi)*divi
}

// NEGX returns the negated value of X.
func NEGX(x float64) float64 {
	return -x
}

// FRACT returns the fractional part of a real number.
// e.g., Fract(3.14) = 0.14
func FRACT(x float64) float64 {
	if math.Abs(x) < 2.0e9 {
		return x - math.Trunc(x)
	}
	return 0.0
}

func RAD(DEG float64) float64 {
	if DEG == 0.0 {
		return 0.0
	}
	return MODR(0.0174532925199433*DEG, Math.Pi2)
}

// REAL_TO_FRAC calculates the closest fraction for a real number.
func REAL_TO_FRAC(x float64, n int) (num, den int) {
	sign := x < 0
	if sign {
		x = -x
	}

	var num_old, den_old int64
	var num_cur, den_cur int64 = 1, 0

	for {
		x_rounded := int64(x)

		num_old, num_cur = num_cur, num_cur*x_rounded+num_old
		den_old, den_cur = den_cur, den_cur*x_rounded+den_old

		x_no_frac := float64(x_rounded)
		if x == x_no_frac {
			if den_cur <= int64(n) {
				num_old, den_old = num_cur, den_cur
			}
			break
		}
		x = 1.0 / (x - x_no_frac)
		if den_cur > int64(n) {
			break
		}
	}

	if sign {
		return -int(num_old), int(den_old)
	}
	return int(num_old), int(den_old)
}

// RND rounds a real number to N total digits.
func RND(x float64, n int) float64 {
	if x == 0.0 {
		return 0.0
	}
	m := EXPN(10.0, n-CEIL(math.Log10(math.Abs(x))))
	return float64(int64(x*m+0.5)) / m
}

// ROUND rounds a real number to N decimal places.
func ROUND(in float64, n int) float64 {
	n = int(LIMIT(0, float64(n), 8))
	m := float64(Setup.Decades[n])
	return math.Round(in*m) / m
}

// SGN returns the sign of a real number (-1, 0, 1).
func SGN(x float64) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// SIGMOID calculates the sigmoid function.
func SIGMOID(x float64) float64 {
	if x > 20.0 {
		return 1.0
	}
	if x < -85.0 {
		return 0.0
	}
	return 1.0 / (1.0 + math.Exp(-x))
}

// SIGN_I checks if an integer is negative.
func SIGN_I(in int64) bool {
	return in < 0
}

// SIGN_R checks if a real number is negative.
func SIGN_R(in float64) bool {
	return in < 0.0
}

// SINC calculates the normalized sinc function.
func SINC(x float64) float64 {
	if x == 0.0 {
		return 1.0
	}
	return math.Sin(x) / x
}

// SQRTN calculates the Nth root of X.
func SQRTN(x float64, n int) float64 {
	if n > 0 && x >= 0 {
		return math.Pow(x, 1.0/float64(n))
	}
	return 0.0
}

// TANC calculates the tanc function.
func TANC(x float64) float64 {
	if x == 0.0 {
		return 1.0
	}
	return math.Tan(x) / x
}

// WINDOW checks if a value is within a range (exclusive).
func WINDOW(low, in, high float64) bool {
	return in > low && in < high
}

// WINDOW2 checks if a value is within a range (inclusive).
func WINDOW2(low, in, high float64) bool {
	return in >= low && in <= high
}
