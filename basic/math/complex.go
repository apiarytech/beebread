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
	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// CABS returns the absolute value of a complex number.
func CABS(x COMPLEX) iec.REAL {
	return HYPOT(x.RE, x.IM)
}

// CACOS calculates the arcus cosinus of a complex number.
func CACOS(x COMPLEX) COMPLEX {
	y := CACOSH(x)
	return COMPLEX{RE: y.IM, IM: -y.RE}
}

// CACOSH calculates the hyperbolic arcus cosinus of a complex number.
func CACOSH(x COMPLEX) COMPLEX {
	y := CSQRT(COMPLEX{
		RE: (x.RE-x.IM)*(x.RE+x.IM) - 1.0,
		IM: 2.0 * x.RE * x.IM,
	})
	y.RE = y.RE + x.RE
	y.IM = y.IM + x.IM
	return CLOG(y)
}

// CADD adds two complex numbers.
func CADD(x, y COMPLEX) COMPLEX {
	return COMPLEX{RE: x.RE + y.RE, IM: x.IM + y.IM}
}

// CARG calculates the phase angle (argument) of a complex number.
func CARG(x COMPLEX) iec.REAL {
	return ATAN2(x.IM, x.RE)
}

// CASIN calculates the arcus sinus of a complex number.
func CASIN(x COMPLEX) COMPLEX {
	y := CASINH(COMPLEX{RE: -x.IM, IM: x.RE})
	return COMPLEX{RE: y.IM, IM: -y.RE}
}

// CASINH calculates the hyperbolic arcus sinus of a complex number.
func CASINH(x COMPLEX) COMPLEX {
	y := CSQRT(COMPLEX{
		RE: (x.RE-x.IM)*(x.RE+x.IM) + 1.0,
		IM: 2.0 * x.RE * x.IM,
	})
	y.RE = y.RE + x.RE
	y.IM = y.IM + x.IM
	return CLOG(y)
}

// CATAN calculates the arcus tangens of a complex number.
func CATAN(x COMPLEX) COMPLEX {
	var out COMPLEX
	r2 := x.RE * x.RE
	den := 1.0 - r2 - x.IM*x.IM
	out.RE = 0.5 * ATAN(2.0*x.RE/den)
	num := x.IM + 1.0
	num = r2 + num*num
	den = x.IM - 1.0
	den = r2 + den*den
	out.IM = 0.25 * (LN(num) - LN(den))
	return out
}

// CATANH calculates the hyperbolic arcus tangens of a complex number.
func CATANH(x COMPLEX) COMPLEX {
	var out COMPLEX
	i2 := x.IM * x.IM
	num := 1.0 + x.RE
	num = i2 + num*num
	den := 1.0 - x.RE
	den = i2 + den*den
	out.RE = 0.25 * (LN(num) - LN(den))
	den = 1 - x.RE*x.RE - i2
	out.IM = 0.5 * ATAN(2.0*x.IM/den)
	return out
}

// CCON returns the conjugate of a complex number.
func CCON(x COMPLEX) COMPLEX {
	return COMPLEX{RE: x.RE, IM: -x.IM}
}

// CCOS calculates the cosinus of a complex number.
func CCOS(x COMPLEX) COMPLEX {
	return CCOSH(CSET(-x.IM, x.RE))
}

// CCOSH calculates the hyperbolic cosinus of a complex number.
func CCOSH(x COMPLEX) COMPLEX {
	return COMPLEX{RE: COSH(x.RE) * COS(x.IM), IM: SINH(x.RE) * SIN(x.IM)}
}

// CDIV divides two complex numbers.
func CDIV(x, y COMPLEX) COMPLEX {
	temp := y.RE*y.RE + y.IM*y.IM
	return COMPLEX{
		RE: (x.RE*y.RE + x.IM*y.IM) / temp,
		IM: (x.IM*y.RE - x.RE*y.IM) / temp,
	}
}

// CEXP calculates e to the power of a complex number.
func CEXP(x COMPLEX) COMPLEX {
	temp := EXP(x.RE)
	return COMPLEX{RE: temp * COS(x.IM), IM: temp * SIN(x.IM)}
}

// CINV calculates the inverse of a complex number, 1 / x.
func CINV(x COMPLEX) COMPLEX {
	temp := x.RE*x.RE + x.IM*x.IM
	return COMPLEX{RE: x.RE / temp, IM: -x.IM / temp}
}

// CLOG calculates the natural logarithm of a complex number.
func CLOG(x COMPLEX) COMPLEX {
	return COMPLEX{RE: LN(HYPOT(x.RE, x.IM)), IM: ATAN2(x.IM, x.RE)}
}

// CMUL multiplies two complex numbers.
func CMUL(x, y COMPLEX) COMPLEX {
	return COMPLEX{
		RE: x.RE*y.RE - x.IM*y.IM,
		IM: x.RE*y.IM + x.IM*y.RE,
	}
}

// CPOL creates a complex number from the polar form, the length l and the
// angle a.
func CPOL(l, a iec.REAL) COMPLEX {
	return COMPLEX{RE: l * COS(a), IM: l * SIN(a)}
}

// CPOW calculates x to the power of y.
func CPOW(x, y COMPLEX) COMPLEX {
	return CEXP(CMUL(y, CLOG(x)))
}

// CSET creates a complex number from its real and imaginary parts.
func CSET(re, im iec.REAL) COMPLEX {
	return COMPLEX{RE: re, IM: im}
}

// CSIN calculates the sinus of a complex number.
func CSIN(x COMPLEX) COMPLEX {
	return COMPLEX{RE: COSH(x.IM) * SIN(x.RE), IM: SINH(x.IM) * COS(x.RE)}
}

// CSINH calculates the hyperbolic sinus of a complex number.
func CSINH(x COMPLEX) COMPLEX {
	return COMPLEX{RE: SINH(x.RE) * COS(x.IM), IM: COSH(x.RE) * SIN(x.IM)}
}

// CSQRT calculates the square root of a complex number.
func CSQRT(x COMPLEX) COMPLEX {
	temp := HYPOT(x.RE, x.IM)
	return COMPLEX{
		RE: SQRT(0.5 * (temp + x.RE)),
		IM: iec.REAL(SGN(x.IM)) * SQRT(0.5*(temp-x.RE)),
	}
}

// CSUB subtracts two complex numbers.
func CSUB(x, y COMPLEX) COMPLEX {
	return COMPLEX{RE: x.RE - y.RE, IM: x.IM - y.IM}
}

// CTAN calculates the tangens of a complex number.
func CTAN(x COMPLEX) COMPLEX {
	xi2 := 2.0 * x.IM
	xr2 := 2.0 * x.RE
	temp := 1.0 / (COS(xr2) + COSH(xi2))
	return COMPLEX{RE: temp * SIN(xr2), IM: temp * SINH(xi2)}
}

// CTANH calculates the hyperbolic tangens of a complex number.
func CTANH(x COMPLEX) COMPLEX {
	xi2 := 2.0 * x.IM
	xr2 := 2.0 * x.RE
	temp := 1.0 / (COSH(xr2) + COS(xi2))
	return COMPLEX{RE: temp * SINH(xr2), IM: temp * SIN(xi2)}
}
