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

	. "beebread/basic"
)

// R2Abs calculates the absolute value of a double-precision real number.
func R2Abs(x Real2) Real2 {
	return Real2{
		Rx: float32(math.Abs(float64(x.Rx))),
		R1: float32(math.Abs(float64(x.R1))),
	}
}

// R2Add adds a real to a double-precision real, maintaining higher precision.
// This implements the TwoSum algorithm.
func R2Add(x1 Real2, x2 float32) Real2 {
	var res Real2
	t1 := x1.Rx + x2
	t2 := t1 - x1.Rx
	res.R1 = (x1.Rx - (t1 - t2)) + (x2 - t2) + x1.R1
	res.Rx = t1 + res.R1
	res.R1 = (t1 - res.Rx) + res.R1
	return res
}

// R2Add2 adds two double-precision real numbers.
func R2Add2(x1, x2 Real2) Real2 {
	res := R2Add(x1, x2.Rx)
	res = R2Add(res, x2.R1)
	return res
}

// R2Mul multiplies a double-precision real with a standard real.
func R2Mul(x1 Real2, x2 float32) Real2 {
	return Real2{
		Rx: x1.Rx * x2,
		R1: x1.R1 * x2,
	}
}

// R2Set converts a standard real (float32) to a double-precision real.
func R2Set(x float32) Real2 {
	return Real2{
		Rx: x,
		R1: 0.0,
	}
}
