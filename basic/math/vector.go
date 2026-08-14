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

// V3Abs calculates the absolute value (length) of a vector.
func V3Abs(v Vector3) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// V3Add adds two vectors.
func V3Add(v1, v2 Vector3) Vector3 {
	return Vector3{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

// V3Cross calculates the cross product of two vectors.
func V3Cross(v1, v2 Vector3) Vector3 {
	return Vector3{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

// V3Dot calculates the dot product of two vectors.
func V3Dot(v1, v2 Vector3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// V3Mul multiplies a vector with a scalar.
func V3Mul(v Vector3, s float32) Vector3 {
	return Vector3{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

// V3Sub subtracts two vectors.
func V3Sub(v1, v2 Vector3) Vector3 {
	return Vector3{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}
