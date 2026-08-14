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

import . "beebread/basic"

// CircleA calculates the area of a circle with radius R.
func CircleA(r float64) float64 {
	return r * r * Math.Pi
}

// CircleC calculates the circumference of a circle with radius R.
func CircleC(r float64) float64 {
	return r * Math.Pi2
}

// CircleSeg calculates the area of a circle segment with radius R and angle A.
func CircleSeg(r, a float64) float64 {
	return r * r * Math.Pi * a / 360.0
}

// ConeV calculates the volume of a cone with radius R and height H.
func ConeV(r, h float64) float64 {
	return r * r * Math.Pi * h / 3.0
}

// SphereV calculates the volume of a sphere with radius R.
func SphereV(r float64) float64 {
	return r * r * r * Math.Pi4 / 3.0
}
