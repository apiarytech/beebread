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

package engineering

import (
	"math"

	beeMath "beebread/basic/math"
)

// MULTI_IN is a signal conditioning function for multiple redundant sensors.
func MULTI_IN(in1, in2, in3, def, inMin, inMax float64, mode byte) float64 {
	f1 := in1 > inMin && in1 < inMax
	f2 := in2 > inMin && in2 < inMax
	f3 := in3 > inMin && in3 < inMax

	switch mode {
	case 0: // Average of valid inputs
		var sum float64
		count := 0
		if f1 {
			sum += in1
			count++
		}
		if f2 {
			sum += in2
			count++
		}
		if f3 {
			sum += in3
			count++
		}
		if count == 0 {
			return def
		}
		return sum / float64(count)

	case 1: // Input 1 or default
		if f1 {
			return in1
		}
		return def

	case 2: // Input 2 or default
		if f2 {
			return in2
		}
		return def

	case 3: // Input 3 or default
		if f3 {
			return in3
		}
		return def

	case 4: // Default
		return def

	case 5: // Lowest of valid inputs
		res := inMax
		if f1 {
			res = in1
		}
		if f2 && in2 < res {
			res = in2
		}
		if f3 && in3 < res {
			res = in3
		}
		if res == inMax {
			return def
		}
		return res

	case 6: // Highest of valid inputs
		res := inMin
		if f1 {
			res = in1
		}
		if f2 && in2 > res {
			res = in2
		}
		if f3 && in3 > res {
			res = in3
		}
		if res == inMin {
			return def
		}
		return res

	case 7: // Middle or average of 2
		if f1 && f2 && f3 {
			return beeMath.MID3(in1, in2, in3)
		} else if f1 && f2 {
			return math.Min(in1, in2)
		} else if f1 && f3 {
			return math.Min(in1, in3)
		} else if f2 && f3 {
			return math.Min(in2, in3)
		} else if f1 {
			return in1
		} else if f2 {
			return in2
		} else if f3 {
			return in3
		}
		return def
	}

	return 0.0
}

// RES_NI returns the resistance for a nickel sensor for a given temperature.
func RES_NI(t, r0 float64) float64 {
	const (
		a = 0.5485
		b = 0.665e-3
		c = 2.805e-9
	)
	t2 := t * t
	return r0 + a*t + b*t2 + c*t2*t2
}

// RES_NTC returns the resistance for an NTC sensor for a given temperature.
func RES_NTC(t, rn, b float64) float64 {
	return rn * math.Exp(b*(1.0/(t+273.15)-0.00335401643468053))
}

// RES_PT returns the resistance for a platinum sensor for a given temperature.
func RES_PT(t, r0 float64) float64 {
	const (
		a = 3.90802e-3
		b = -5.802e-7
		c = -4.27350e-12
	)
	t2 := t * t
	if t >= 0.0 {
		return r0 * (1.0 + a*t + b*t2)
	}
	return r0 * (1.0 + a*t + b*t2 + c*(t-100.0)*t2*t)
}

// RES_SI returns the resistance for a silicon sensor for a given temperature.
func RES_SI(t, rs, ts float64) float64 {
	const (
		a = 7.64e-3
		b = 1.66e-5
	)
	tx := t - ts
	return rs * (1.0 + a*tx + b*tx*tx)
}

// SENSOR_INT calculates the true resistance of a sensor considering parallel and series parasitic resistances.
func SENSOR_INT(voltage, current, rp, rs float64) float64 {
	if current != 0.0 {
		rg := voltage / current
		if (rp + rs - rg) != 0.0 {
			return rp * (rg - rs) / (rp + rs - rg)
		}
	}
	return 0.0
}

// TEMP_NI returns the temperature for a nickel sensor given its resistance.
func TEMP_NI(res, r0 float64) float64 {
	return (math.Sqrt(0.30085225-2.66e-3*(r0-res)) - 0.5485) * 751.8796992
}

// TEMP_NTC returns the temperature for an NTC sensor given its resistance.
func TEMP_NTC(res, rn, b float64) float64 {
	if res > 0.0 {
		return b*298.15/(b+math.Log(res/rn)*298.15) - 273.15
	}
	return 0.0
}

// TEMP_PT returns the temperature for a platinum sensor given its resistance.
func TEMP_PT(res, r0 float64) float64 {
	const (
		a        = 3.9083e-3
		b        = -5.775e-7
		accuracy = 0.01
	)
	x := a * r0
	y := b * r0

	if res >= r0 {
		t1 := x*x - 4.0*y*(r0-res)
		if t1 < 0.0 {
			return 10000.0
		}
		return (-x + math.Sqrt(t1)) / (2.0 * y)
	}

	// Successive approximation for temperatures below 0°C
	tempPt := -100.0
	step := 50.0
	for step > accuracy {
		if RES_PT(tempPt, r0) < res {
			tempPt += step
		} else {
			tempPt -= step
		}
		step *= 0.5
	}
	return tempPt
}

// TEMP_SI returns the temperature for a silicon sensor given its resistance.
func TEMP_SI(res, rs, ts float64) float64 {
	val := res/rs*6.64e-5 - 0.803e-5
	if val < 0 {
		return ts // Return base temperature if sqrt is invalid
	}
	return (-7.64e-3+math.Sqrt(val))*30120.48193 + ts
}
