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
	"math"
	"math/cmplx"

	. "beebread/basic"
)

// toC128 converts the library's Complex struct to Go's native complex128.
func toC128(c Complex) complex128 {
	return complex(float64(c.Re), float64(c.Im))
}

// fromC128 converts Go's native complex128 to the library's Complex struct.
func fromC128(c complex128) Complex {
	return Complex{
		Re: float32(real(c)),
		Im: float32(imag(c)),
	}
}

// CAbs calculates the absolute value (magnitude) of a complex number.
func CAbs(x Complex) float32 {
	return float32(cmplx.Abs(toC128(x)))
}

// CAcos calculates the inverse cosine of a complex number.
func CAcos(x Complex) Complex {
	return fromC128(cmplx.Acos(toC128(x)))
}

// CAdd adds two complex numbers.
func CAdd(x1, x2 Complex) Complex {
	return Complex{
		Re: x1.Re + x2.Re,
		Im: x1.Im + x2.Im,
	}
}

// CAsin calculates the inverse sine of a complex number.
func CAsin(x Complex) Complex {
	return fromC128(cmplx.Asin(toC128(x)))
}

// CAtan calculates the inverse tangent of a complex number.
func CAtan(x Complex) Complex {
	return fromC128(cmplx.Atan(toC128(x)))
}

// CCos calculates the cosine of a complex number.
func CCos(x Complex) Complex {
	return fromC128(cmplx.Cos(toC128(x)))
}

// CDiv divides two complex numbers (x1 / x2).
func CDiv(x1, x2 Complex) Complex {
	return fromC128(toC128(x1) / toC128(x2))
}

// CExp calculates the base-e exponential of a complex number.
func CExp(x Complex) Complex {
	return fromC128(cmplx.Exp(toC128(x)))
}

// CLn calculates the natural logarithm of a complex number.
func CLn(x Complex) Complex {
	return fromC128(cmplx.Log(toC128(x)))
}

// CMul multiplies two complex numbers.
func CMul(x1, x2 Complex) Complex {
	return Complex{
		Re: x1.Re*x2.Re - x1.Im*x2.Im,
		Im: x1.Re*x2.Im + x1.Im*x2.Re,
	}
}

// CSet creates a complex number from two real values.
func CSet(re, im float32) Complex {
	return Complex{Re: re, Im: im}
}

// Hypot calculates the hypotenuse Sqrt(p*p + q*q).
func Hypot(p, q float64) float64 {
	return math.Hypot(p, q)
}

// Sgn returns the sign of a real number.
// Returns -1 for negative, 1 for positive, and 0 for zero.
func Sgn(x float32) float32 {
	if x > 0 {
		return 1.0
	}
	if x < 0 {
		return -1.0
	}
	return 0.0
}

// CArg calculates the phase angle (argument) of a complex number.
func CArg(x Complex) float32 {
	return float32(cmplx.Phase(toC128(x)))
}

// CCon calculates the conjugation of a complex number.
func CCon(x Complex) Complex {
	return fromC128(cmplx.Conj(toC128(x)))
}

// CInv calculates the inverse of a complex number (1 / x).
func CInv(x Complex) Complex {
	return fromC128(1 / (toC128(x)))
}

// CPol creates a complex number from polar coordinates (length and angle).
func CPol(l, a float32) Complex {
	return fromC128(cmplx.Rect(float64(l), float64(a)))
}

// CPow calculates the complex power function x to the power of y.
func CPow(x, y Complex) Complex {
	return fromC128(cmplx.Pow(toC128(x), toC128(y)))
}

// CSin calculates the sine of a complex number.
func CSin(x Complex) Complex {
	return fromC128(cmplx.Sin(toC128(x)))
}

// CSinh calculates the hyperbolic sine of a complex number.
func CSinh(x Complex) Complex {
	return fromC128(cmplx.Sinh(toC128(x)))
}

// CAcosh calculates the hyperbolic arccosine of a complex number.
func CAcosh(x Complex) Complex {
	return fromC128(cmplx.Acosh(toC128(x)))
}

// CAsinh calculates the hyperbolic arcsine of a complex number.
func CAsinh(x Complex) Complex {
	return fromC128(cmplx.Asinh(toC128(x)))
}

// CAtanh calculates the hyperbolic arctangent of a complex number.
func CAtanh(x Complex) Complex {
	return fromC128(cmplx.Atanh(toC128(x)))
}

// CCosh calculates the hyperbolic cosine of a complex number.
func CCosh(x Complex) Complex {
	return fromC128(cmplx.Cosh(toC128(x)))
}

// CSqrt calculates the square root of a complex number.
func CSqrt(x Complex) Complex {
	return fromC128(cmplx.Sqrt(toC128(x)))
}

// CSub subtracts two complex numbers (x1 - x2).
func CSub(x1, x2 Complex) Complex {
	return Complex{
		Re: x1.Re - x2.Re,
		Im: x1.Im - x2.Im,
	}
}

// CTan calculates the tangent of a complex number.
func CTan(x Complex) Complex {
	return fromC128(cmplx.Tan(toC128(x)))
}

// CTanh calculates the hyperbolic tangent of a complex number.
func CTanh(x Complex) Complex {
	return fromC128(cmplx.Tanh(toC128(x)))
}
