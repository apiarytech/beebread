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
	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/royaljelly/iec"
)

// MULTI_IN reads a value from up to 3 sensors, ignoring a sensor whose value
// is not between inMin and inMax. mode 0 averages the sensors, 1..3 take
// in1..in3, 4 takes def, 5 the lowest, 6 the highest and 7 the middle one,
// or the lower of two. Without a valid sensor it is def, and 0 for a mode
// above 7.
func MULTI_IN(in1, in2, in3, def, inMin, inMax iec.REAL, mode iec.BYTE) iec.REAL {
	f1 := in1 > inMin && in1 < inMax
	f2 := in2 > inMin && in2 < inMax
	f3 := in3 > inMin && in3 < inMax
	var out iec.REAL
	switch mode {
	case 0:
		count := 0
		for _, v := range []struct {
			ok bool
			x  iec.REAL
		}{{f1, in1}, {f2, in2}, {f3, in3}} {
			if v.ok {
				count++
				out += v.x
			}
		}
		if count == 0 {
			return def
		}
		return out / iec.REAL(count)
	case 1:
		return SEL(iec.BOOL(f1), def, in1)
	case 2:
		return SEL(iec.BOOL(f2), def, in2)
	case 3:
		return SEL(iec.BOOL(f3), def, in3)
	case 4:
		return def
	case 5:
		out = SEL(iec.BOOL(f1), inMax, in1)
		if f2 && in2 < out {
			out = in2
		}
		if f3 && in3 < out {
			out = in3
		}
		if out == inMax {
			out = def
		}
		return out
	case 6:
		out = SEL(iec.BOOL(f1), inMin, in1)
		if f2 && in2 > out {
			out = in2
		}
		if f3 && in3 > out {
			out = in3
		}
		if out == inMin {
			out = def
		}
		return out
	case 7:
		switch {
		case f1 && f2 && f3:
			return math.MID3(in1, in2, in3)
		case f1 && f2:
			return min(in1, in2)
		case f1 && f3:
			return min(in1, in3)
		case f2 && f3:
			return min(in2, in3)
		case f1:
			return in1
		case f2:
			return in2
		case f3:
			return in3
		}
		return def
	}
	return 0
}

// RES_NI returns the resistance of a nickel sensor with the resistance r0
// at 0 °C at the temperature t, -60..+180 °C.
func RES_NI(t, r0 iec.REAL) iec.REAL {
	const a, b, c = 0.5485, 0.665e-3, 2.805e-9
	t2 := t * t
	return r0 + a*t + b*t2 + c*t2*t2
}

// RES_NTC returns the resistance of an NTC sensor with the resistance rn at
// 25 °C and the constant b at the temperature t in °C.
func RES_NTC(t, rn, b iec.REAL) iec.REAL {
	return rn * EXP(b*(1.0/(t+273.15)-0.00335401643468053))
}

// RES_PT returns the resistance of a platinum sensor with the resistance r0
// at 0 °C at the temperature t, -200..+850 °C.
func RES_PT(t, r0 iec.REAL) iec.REAL {
	const a, b, c = 3.90802e-3, -5.802e-7, -4.27350e-12
	t2 := t * t
	if t >= 0.0 {
		return r0 * (1.0 + a*t + b*t2)
	}
	return r0 * (1.0 + a*t + b*t2 + c*(t-100.0)*t2*t)
}

// RES_SI returns the resistance of a silicon sensor with the resistance rs
// at the temperature ts at the temperature t, -50..+150 °C.
func RES_SI(t, rs, ts iec.REAL) iec.REAL {
	const a, b = 7.64e-3, 1.66e-5
	tx := t - ts
	return rs * (1.0 + a*tx + b*tx*tx)
}

// SENSOR_INT returns the resistance of a sensor measured with the voltage
// and current, with the parasitic resistance rp parallel to it and rs in
// series.
func SENSOR_INT(voltage, current, rp, rs iec.REAL) iec.REAL {
	if current != 0.0 {
		rg := voltage / current
		return rp * (rg - rs) / (rp + rs - rg)
	}
	return 0.0
}

// TEMP_NI returns the temperature of a nickel sensor with the resistance r0
// at 0 °C and the resistance res, -60..+180 °C.
func TEMP_NI(res, r0 iec.REAL) iec.REAL {
	return (SQRT(0.30085225-2.66e-3*(r0-res)) - 0.5485) * 751.8796992
}

// TEMP_NTC returns the temperature of an NTC sensor with the resistance rn
// at 25 °C and the constant b and the resistance res, 0..85 °C.
func TEMP_NTC(res, rn, b iec.REAL) iec.REAL {
	if res > 0.0 {
		return b*298.15/(b+LN(res/rn)*298.15) - 273.15
	}
	return 0.0
}

// TEMP_PT returns the temperature of a platinum sensor with the resistance
// r0 at 0 °C and the resistance res, -200..+850 °C. It is 10000 if there is
// none.
func TEMP_PT(res, r0 iec.REAL) iec.REAL {
	const a, b, accuracy = 3.9083e-3, -5.775e-7, 0.01
	x := a * r0
	y := b * r0
	if res >= r0 {
		t1 := x*x - 4.0*y*(r0-res)
		if t1 < 0.0 {
			return 10000.0
		}
		return (-x + SQRT(t1)) / (2.0 * y)
	}
	var t, step iec.REAL = -100.0, 50.0
	for step > accuracy {
		if RES_PT(t, r0) < res {
			t += step
		} else {
			t -= step
		}
		step /= 2
	}
	return t
}

// TEMP_SI returns the temperature of a silicon sensor with the resistance
// rs at the temperature ts and the resistance res, -60..+180 °C.
func TEMP_SI(res, rs, ts iec.REAL) iec.REAL {
	return (-7.64e-3+SQRT(res/rs*6.64e-5-0.803e-5))*30120.48193 + ts
}
