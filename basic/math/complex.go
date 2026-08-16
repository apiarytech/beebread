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

// toC128 converts the library's COMPLEX struct to Go's native COMPLEX128.
func toC128(c COMPLEX) complex128 {
	return complex(float64(c.Re), float64(c.Im))
}

// fromC128 converts Go's native COMPLEX128 to the library's COMPLEX struct.
func fromC128(c complex128) COMPLEX {
	return COMPLEX{
		Re: float32(real(c)),
		Im: float32(imag(c)),
	}
}

// C_ABS calculates the absolute value (magnitude) of a COMPLEX number.
func C_ABS(x COMPLEX) float32 {
	return float32(cmplx.Abs(toC128(x)))
}

// C_ACOS calculates the inverse cosine of a COMPLEX number.
func C_ACOS(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Acos(toC128(x)))
}

// C_ADD adds two COMPLEX numbers.
func C_ADD(x1, x2 COMPLEX) COMPLEX {
	return COMPLEX{
		Re: x1.Re + x2.Re,
		Im: x1.Im + x2.Im,
	}
}

// C_ASIN calculates the inverse sine of a COMPLEX number.
func C_ASIN(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Asin(toC128(x)))
}

// C_ATAN calculates the inverse tangent of a COMPLEX number.
func C_ATAN(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Atan(toC128(x)))
}

// C_COS calculates the cosine of a COMPLEX number.
func C_COS(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Cos(toC128(x)))
}

// C_DIV divides two COMPLEX numbers (x1 / x2).
func C_DIV(x1, x2 COMPLEX) COMPLEX {
	return fromC128(toC128(x1) / toC128(x2))
}

// C_EXP calculates the base-e exponential of a COMPLEX number.
func C_EXP(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Exp(toC128(x)))
}

// C_LOG calculates the natural logarithm of a COMPLEX number.
func C_LOG(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Log(toC128(x)))
}

// C_MUL multiplies two COMPLEX numbers.
func C_MUL(x1, x2 COMPLEX) COMPLEX {
	return COMPLEX{
		Re: x1.Re*x2.Re - x1.Im*x2.Im,
		Im: x1.Re*x2.Im + x1.Im*x2.Re,
	}
}

// C_SET creates a COMPLEX number from two real values.
func C_SET(re, im float32) COMPLEX {
	return COMPLEX{Re: re, Im: im}
}

// HYPOT calculates the hypotenuse Sqrt(p*p + q*q).
func HYPOT(p, q float64) float64 {
	return math.Hypot(p, q)
}

// C_ARG calculates the phase angle (argument) of a COMPLEX number.
func C_ARG(x COMPLEX) float32 {
	return float32(cmplx.Phase(toC128(x)))
}

// C_CON calculates the conjugation of a COMPLEX number.
func C_CON(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Conj(toC128(x)))
}

// C_INV calculates the inverse of a COMPLEX number (1 / x).
func C_INV(x COMPLEX) COMPLEX {
	return fromC128(1 / (toC128(x)))
}

// C_POL creates a COMPLEX number from polar coordinates (length and angle).
func C_POL(l, a float64) COMPLEX {
	return fromC128(cmplx.Rect(l, a))
}

// C_POW calculates the COMPLEX power function x to the power of y.
func C_POW(x, y COMPLEX) COMPLEX {
	return fromC128(cmplx.Pow(toC128(x), toC128(y)))
}

// C_SIN calculates the sine of a COMPLEX number.
func C_SIN(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Sin(toC128(x)))
}

// C_SINH calculates the hyperbolic sine of a COMPLEX number.
func C_SINH(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Sinh(toC128(x)))
}

// C_ACOSH calculates the hyperbolic arccosine of a COMPLEX number.
func C_ACOSH(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Acosh(toC128(x)))
}

// C_ASINH calculates the hyperbolic arcsine of a COMPLEX number.
func C_ASINH(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Asinh(toC128(x)))
}

// C_ATANH calculates the hyperbolic arctangent of a COMPLEX number.
func C_ATANH(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Atanh(toC128(x)))
}

// C_COSH calculates the hyperbolic cosine of a COMPLEX number.
func C_COSH(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Cosh(toC128(x)))
}

// C_SQRT calculates the square root of a COMPLEX number.
func C_SQRT(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Sqrt(toC128(x)))
}

// C_SUB subtracts two COMPLEX numbers (x1 - x2).
func C_SUB(x1, x2 COMPLEX) COMPLEX {
	return COMPLEX{
		Re: x1.Re - x2.Re,
		Im: x1.Im - x2.Im,
	}
}

// C_TAN calculates the tangent of a COMPLEX number.
func C_TAN(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Tan(toC128(x)))
}

// C_TANH calculates the hyperbolic tangent of a COMPLEX number.
func C_TANH(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Tanh(toC128(x)))
}

// C_COT calculates the cotangent of a COMPLEX number.
func C_COT(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Cot(toC128(x)))
}

// C_COTH calculates the hyperbolic cotangent of a COMPLEX number.
func C_COTH(x COMPLEX) COMPLEX {
	return fromC128(cmplx.Cosh(toC128(x)) / cmplx.Sinh(toC128(x))) // No direct coth in cmplx, calculate from cosh/sinh
}
