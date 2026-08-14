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

package math

import (
	"beebread/basic"
	. "beebread/basic"
	"beebread/basic/logic"
	"math"
)

// Acosh calculates the hyperbolic arccos of X.
func Acosh(x float64) float64 {
	return math.Log(x + math.Sqrt(x*x-1.0))
}

// Acoth calculates the hyperbolic arccot of X.
func Acoth(x float64) float64 {
	return 0.5 * math.Log((x+1.0)/(x-1.0))
}

// Agdf calculates the inverse Gudermann function of X.
func Agdf(x float64) float64 {
	return math.Log(math.Tan(Math.Pi025 + 0.5*x))
}

// Asinh calculates the hyperbolic arcsin of X.
func Asinh(x float64) float64 {
	return math.Log(x + math.Sqrt(x*x+1.0))
}

// Atan2 calculates the arctan of Y/X.
func Atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

// Atanh calculates the hyperbolic arctan of X.
func Atanh(x float64) float64 {
	return 0.5 * math.Log((1.0+x)/(1.0-x))
}

// Beta calculates the beta function of X and Y.
func Beta(x, y float64) float64 {
	// math.Gamma is the standard library implementation for the Gamma function.
	return (math.Gamma(x) * math.Gamma(y)) / math.Gamma(x+y)
}

// Binom calculates the binomial coefficient of N over K.
func Binom(n, k int) int64 {
	if k >= 0 && k <= n && n < 13 {
		return int64(Math.Facts[n] / (Math.Facts[k] * Math.Facts[n-k]))
	}
	return 0 // Or handle error appropriately
}

// Cauchy calculates the Cauchy distribution.
func Cauchy(x, m, g float64) float64 {
	return Math.Pi_inv * g / ((x-m)*(x-m) + g*g)
}

// CauchyCD calculates the cumulative Cauchy distribution.
func CauchyCD(x, m, g float64) float64 {
	return 0.5 + Math.Pi_inv*math.Atan((x-m)/g)
}

// Ceil calculates the next higher integer of X.
func Ceil(x float64) int {
	return int(math.Ceil(x))
}

// Ceil2 calculates the next higher integer of X/Y * Y.
func Ceil2(x, y float64) float64 {
	if y != 0.0 {
		return math.Ceil(x/y) * y
	}
	return 0.0
}

// Cmp compares two real values with a given deviation D.
// Returns 0 if |X1-X2| <= D.
// Returns 1 if X1 > X2.
// Returns -1 if X1 < X2.
func Cmp(x1, x2, d float64) int {
	if math.Abs(x1-x2) <= d {
		return 0
	} else if x1 > x2 {
		return 1
	}
	return -1
}

// Cosh calculates the cosinus hyperbolicus of X.
func Cosh(x float64) float64 {
	return math.Cosh(x)
}

// Coth calculates the cotangens hyperbolicus of X.
func Coth(x float64) float64 {
	return (math.Exp(x) + math.Exp(-x)) / (math.Exp(x) - math.Exp(-x))
}

// Differ checks if in1 differs more than x from in2.
// It returns true if abs(in1 - in2) > x.
func Differ(in1, in2, x float64) bool {
	return math.Abs(in1-in2) > x
}

// DTrunc truncates a real value to a dint (int64).
func DTrunc(x float64) int64 {
	return int64(math.Trunc(x))
}

// LimitB limits a byte value to a given range.
func LimitB(min, val, max byte) byte {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// LimitDW limits a dword (uint32) value to a given range.
func LimitDW(min, val, max uint32) uint32 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// LimitInt64 limits an int64 value to a given range.
func LimitInt64(min, val, max int64) int64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// Inc implements a modular increment.
func Inc(x, d, m int) int {
	// Ensure the modulo is positive
	mod := m + 1
	if mod <= 0 {
		return x // Or handle error
	}

	res := (x + d) % mod
	if res < 0 {
		res += mod
	}
	return res
}

// Inc1 increments X by 1 and wraps around to 0 if N is reached.
// It generates a sequence: 0, 1, 2, ..., N-1, 0, ...
func Inc1(x, n int) int {
	if x >= n-1 {
		return 0
	}
	return x + 1
}

// Mid3 returns the middle value of three real inputs.
func Mid3(in1, in2, in3 float64) float64 {
	// This is a more readable Go implementation of the original ST logic.
	if (in1 <= in2 && in2 <= in3) || (in3 <= in2 && in2 <= in1) {
		return in2
	}
	if (in2 <= in1 && in1 <= in3) || (in3 <= in1 && in1 <= in2) {
		return in1
	}
	return in3
}

// ModR is a modulo function for real numbers.
// e.g., ModR(5.5, 2.5) = 0.5
func ModR(in, divi float64) float64 {
	if divi == 0.0 {
		return 0.0
	}
	return in - math.Floor(in/divi)*divi
}

// Fract returns the fractional part of a real number.
// e.g., Fract(3.14) = 0.14
func Fract(x float64) float64 {
	if math.Abs(x) < 2.0e9 {
		return x - math.Trunc(x)
	}
	return 0.0
}

// Rdm calculates a pseudo-random number between 0.0 and 1.0.
// To use Rdm more than once per cycle, it needs to be called with different seed values for `last`.
func Rdm(last float64) float64 {
	tn := uint32(logic.TPlcUs())
	tc := logic.BitCount(tn)

	// Scramble bits based on original ST logic
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 2), 31)
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 5), 30)
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 4), 29)
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 1), 28)
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 0), 27)
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 7), 26)
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 6), 25)
	tn = logic.BitLoadDw(tn, logic.BitOfDword(tn, 3), 24)

	tn = (tn << uint(tc)) | (tn >> (32 - uint(tc))) // ROL
	tn |= 0x80000001
	tn = tn%71474513 + uint32(tc+77)

	return Fract(float64(tn) / 10000000.0 * (basic.Math.E - Limit(0.0, last, 1.0)))
}
