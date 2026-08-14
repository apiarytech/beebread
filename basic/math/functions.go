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
	"time"
)

// FLin calculates the linear equation f_lin = a*x + b.
func FLin(x, a, b float64) float64 {
	return a*x + b
}

// FLin2 calculates the linear equation f_lin = a*x + b given by two points (x1, y1) and (x2, y2).
// It handles the vertical line case by returning y1.
func FLin2(x, x1, y1, x2, y2 float64) float64 {
	if x2-x1 == 0.0 {
		return y1
	}
	return (y2-y1)/(x2-x1)*(x-x1) + y1
}

// FPoly calculates a polynomial using Horner's method for efficiency.
// It evaluates C[0] + C[1]*X^1 + C[2]*X^2 + ... + C[7]*X^7.
func FPoly(x float64, c [8]float64) float64 {
	res := c[7]
	for i := 6; i >= 0; i-- {
		res = res*x + c[i]
	}
	return res
}

// FPower calculates the power equation f_power = a * x^n.
func FPower(a, x, n float64) float64 {
	return a * math.Pow(x, n)
}

// FQuad calculates the quadratic equation f_quad = a*x^2 + b*x + c.
func FQuad(x, a, b, c float64) float64 {
	return (a*x+b)*x + c
}

// FRmpB calculates a ramp for a byte value and limits the output to 0-255.
// It avoids overflow issues during calculation.
func FRmpB(start byte, dir bool, td, tr time.Duration) byte {
	if td < tr && tr > 0 {
		// Calculate the ramped value.
		// The calculation is done using larger integer types to prevent overflow before scaling down.
		val := byte((uint64(td) * 256) / uint64(tr))

		if dir { // Ramp up
			// Prevent overflow when adding
			if int(start)+int(val) > 255 {
				return 255
			}
			return start + val
		} else { // Ramp down
			// Prevent underflow when subtracting
			if int(val) > int(start) {
				return 0
			}
			return start - val
		}
	} else if dir {
		return 255
	}
	return 0
}

// LinearInt calculates an output based on a linear interpolation of up to 20 coordinates.
// The input coordinates in xy must be sorted by ascending X values.
func LinearInt(x float64, xy [][2]float64, pts int) float64 {
	// Ensure pts is within the array bounds
	numPts := len(xy)
	if pts > numPts {
		pts = numPts
	}
	if pts < 2 {
		return 0.0 // Not enough points to interpolate
	}

	// Find the correct segment for interpolation
	i := 1
	for i < pts-1 && xy[i][0] < x {
		i++
	}

	// Calculate the output value on the corresponding segment
	return FLin2(x, xy[i-1][0], xy[i-1][1], xy[i][0], xy[i][1])
}

// PolynomInt calculates an output based on a Newton polynomial interpolation of up to 5 coordinates.
// The input coordinates in xy must be sorted by ascending X values.
func PolynomInt(x float64, xy [][2]float64, pts int) float64 {
	// Ensure pts is within the array bounds
	numPts := len(xy)
	if pts > numPts {
		pts = numPts
	}
	if pts == 0 {
		return 0.0
	}

	// Create a mutable copy to perform the divided differences calculation
	xyCopy := make([][2]float64, pts)
	copy(xyCopy, xy)

	// Calculate divided differences
	for i := 1; i < pts; i++ {
		for j := pts - 1; j >= i; j-- {
			if xyCopy[j][0]-xyCopy[j-i][0] != 0 {
				xyCopy[j][1] = (xyCopy[j][1] - xyCopy[j-1][1]) / (xyCopy[j][0] - xyCopy[j-i][0])
			} else {
				xyCopy[j][1] = 0 // Avoid division by zero
			}
		}
	}

	// Evaluate the polynomial using Horner's method
	res := xyCopy[pts-1][1]
	for i := pts - 2; i >= 0; i-- {
		res = res*(x-xyCopy[i][0]) + xyCopy[i][1]
	}

	return res
}

// FtAvg is a moving average filter over N samples.
type FtAvg struct {
	Avg float64

	// internal state
	buff []float64
	init bool
}

// Update calculates the moving average.
func (f *FtAvg) Update(in float64, e bool, n int, rst bool) {
	if n <= 0 {
		n = 1
	}
	if n > 32 {
		n = 32
	}

	if !f.init || rst {
		f.buff = make([]float64, n)
		for i := 0; i < n; i++ {
			f.buff[i] = in
		}
		f.Avg = in
		f.init = true
	} else if e {
		// The original ST code uses a DELAY block which is a circular buffer.
		// A more efficient way to calculate moving average is to subtract the
		// oldest value and add the new one.
		f.Avg = f.Avg + (in-f.buff[0])/float64(n)

		// Shift buffer
		copy(f.buff, f.buff[1:])
		f.buff[n-1] = in
	}
}

// FtMinMax stores the minimum and maximum value of an input signal.
type FtMinMax struct {
	Max float64
	Min float64

	// internal state
	init bool
}

// Update executes the min/max tracking logic.
func (f *FtMinMax) Update(in float64, rst bool) {
	if rst || !f.init {
		f.Min = in
		f.Max = in
		f.init = true
	} else if in < f.Min {
		f.Min = in
	} else if in > f.Max {
		f.Max = in
	}
}

// FtRmp is a ramp function that follows an input signal with a linear ramp.
type FtRmp struct {
	Out  float64
	Busy bool
	UD   bool // Up/Down direction

	// internal state
	last time.Time
	init bool
}

// Update executes the ramp logic for one cycle.
func (f *FtRmp) Update(rmp bool, in, kr, kf float64) {
	tx := time.Now()

	if !f.init {
		f.init = true
		f.last = tx
		f.Out = in
	}

	elapsed := tx.Sub(f.last).Seconds()

	if !rmp {
		f.Out = in
		f.Busy = false
	} else if f.Out > in {
		// Ramp down
		f.Out -= elapsed * kf
		f.Out = math.Max(in, f.Out)
	} else if f.Out < in {
		// Ramp up
		f.Out += elapsed * kr
		f.Out = math.Min(in, f.Out)
	}

	// Set busy and direction flags
	if f.Out < in {
		f.Busy = true
		f.UD = true
	} else if f.Out > in {
		f.Busy = true
		f.UD = false
	} else {
		f.Busy = false
	}
	f.last = tx
}

// Limit restricts a value to a given range.
func Limit(min, val, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}
