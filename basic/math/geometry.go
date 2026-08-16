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
	. "beebread/basic"
	"math"
)

// CIRCLE_AREA calculates the area of a circle with radius R. This is a helper function.
func CIRCLE_AREA(r float64) float64 {
	return r * r * Math.Pi
}

// CIRCLE_C calculates the circumference of a circle sector with radius R and angle A.
func CIRCLE_C(r, a float64) float64 {
	return r * Math.Pi2
}

// CIRCLE_A calculates the area of a circle sector with radius R and angle A.
func CIRCLE_A(r, a float64) float64 {
	return r * r * Math.Pi * a / 360.0
}

// CONE_V calculates the volume of a cone with radius R and height H.
func CONE_V(r, h float64) float64 {
	return r * r * Math.Pi * h / 3.0
}

// SPHERE_V calculates the volume of a sphere with radius R.
func SPHERE_V(r float64) float64 {
	return r * r * r * Math.Pi4 / 3.0
}

// ELLIPSE_A calculates the area of an ellipse based on its two radii.
func ELLIPSE_A(r1, r2 float64) float64 {
	return Math.Pi * r1 * r2
}

// ELLIPSE_C calculates an approximation of the circumference of an ellipse.
// This uses Ramanujan's second approximation.
func ELLIPSE_C(r1, r2 float64) float64 {
	return Math.Pi * (3.0*(r1+r2) - math.Sqrt((3.0*r1+r2)*(r1+3.0*r2)))
}

// TRIANGLE_A calculates the area of a triangle.
// If angle A is 0, it uses Heron's formula with three sides (s1, s2, s3).
// Otherwise, it uses the formula with two sides and the included angle (s1, s2, A).
func TRIANGLE_A(s1, a, s2, s3 float64) float64 {
	if a == 0.0 {
		// Heron's formula for area from three sides
		s := (s1 + s2 + s3) / 2.0
		return math.Sqrt(s * (s - s1) * (s - s2) * (s - s3))
	}

	// Area from two sides and the included angle
	// The angle is converted from degrees to radians.
	return 0.5 * s1 * s2 * math.Sin(a/math.Pi)
}
