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

// V3_ABS calculates the absolute value (length) of a vector.
func V3_ABS(v VECTOR_3) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// V3_ADD adds two vectors.
func V3_ADD(v1, v2 VECTOR_3) VECTOR_3 {
	return VECTOR_3{
		X: v1.X + v2.X, // ST: V3_ADD.X := A.X + B.X;
		Y: v1.Y + v2.Y, // ST: V3_ADD.Y := A.Y + B.Y;
		Z: v1.Z + v2.Z, // ST: V3_ADD.Z := A.Z + B.Z;
	}
}

// V3_XPRO calculates the cross product of two vectors.
func V3_XPRO(v1, v2 VECTOR_3) VECTOR_3 {
	return VECTOR_3{
		X: v1.Y*v2.Z - v1.Z*v2.Y, // ST: V3_XPRO.X := A.Y * B.Z - A.Z * B.Y;
		Y: v1.Z*v2.X - v1.X*v2.Z, // ST: V3_XPRO.Y := A.Z * B.X - A.X * B.Z;
		Z: v1.X*v2.Y - v1.Y*v2.X, // ST: V3_XPRO.Z := A.X * B.Y - A.Y * B.X;
	}
}

// V3_DPRO calculates the dot product of two vectors.
func V3_DPRO(v1, v2 VECTOR_3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// V3_SMUL multiplies a vector with a scalar.
func V3_SMUL(v VECTOR_3, s float32) VECTOR_3 {
	return VECTOR_3{
		X: v.X * s, // ST: V3_SMUL.X := A.X * M;
		Y: v.Y * s, // ST: V3_SMUL.Y := A.Y * M;
		Z: v.Z * s, // ST: V3_SMUL.Z := A.Z * M;
	}
}

// V3_SUB subtracts two vectors.
func V3_SUB(v1, v2 VECTOR_3) VECTOR_3 {
	return VECTOR_3{
		X: v1.X - v2.X, // ST: V3_SUB.X := A.X - B.X;
		Y: v1.Y - v2.Y, // ST: V3_SUB.Y := A.Y - B.Y;
		Z: v1.Z - v2.Z, // ST: V3_SUB.Z := A.Z - B.Z;
	}
}

// V3_ANG calculates the angle between two vectors in a 3-dimensional space.
func V3_ANG(a, b VECTOR_3) float32 {
	d := V3_ABS(a) * V3_ABS(b)
	if d > 0 {
		// ST: V3_ANG := ACOS(LIMIT(-1.0, V3_DPRO(A, B) / d,1.0));
		return float32(math.Acos(LIMIT(-1.0, float64(V3_DPRO(a, b)/d), 1.0)))
	}
	return 0.0 // Angle is 0 if one or both vectors are zero.
}

// V3_NORM generates a unit vector (length 1) from a vector.
func V3_NORM(a VECTOR_3) VECTOR_3 {
	la := V3_ABS(a)
	if la > 0.0 {
		// ST: V3_NORM := V3_SMUL(A, 1.0 / la);
		return V3_SMUL(a, float32(1.0/la))
	}
	return VECTOR_3{} // Return zero vector if input is zero vector.
}

// V3_NUL checks if a vector is a null vector (all components are zero).
func V3_NUL(a VECTOR_3) bool {
	// ST: V3_NUL := A.X = 0.0 AND A.Y = 0.0 AND A.Z = 0.0;
	return a.X == 0.0 && a.Y == 0.0 && a.Z == 0.0
}

// V3_PAR checks if two vectors are parallel.
func V3_PAR(a, b VECTOR_3) bool {
	// Two vectors are parallel if their cross product is a zero vector.
	// ST: V3_PAR := V3_ABS(V3_XPRO(A, B)) = 0.0;
	return V3_ABS(V3_XPRO(a, b)) == 0.0
}

// V3_REV reverses a vector (multiplies all components by -1).
func V3_REV(a VECTOR_3) VECTOR_3 {
	return VECTOR_3{
		X: -a.X, // ST: V3_REV.X := -A.X;
		Y: -a.Y, // ST: V3_REV.Y := -A.Y;
		Z: -a.Z, // ST: V3_REV.Z := -A.Z;
	}
}

// V3_XANG calculates the angle between the X-axis and a vector.
func V3_XANG(a VECTOR_3) float32 {
	la := V3_ABS(a)
	if la > 0.0 {
		// ST: V3_XANG := ACOS(A.X / la);
		return float32(math.Acos(float64(a.X / la)))
	}
	return 0.0 // Undefined for zero vector, return 0.
}

// V3_YANG calculates the angle between the Y-axis and a vector.
func V3_YANG(a VECTOR_3) float32 {
	la := V3_ABS(a)
	if la > 0.0 {
		// ST: V3_YANG := ACOS(A.Y / la);
		return float32(math.Acos(float64(a.Y / la)))
	}
	return 0.0
}

// V3_ZANG calculates the angle between the Z-axis and a vector.
func V3_ZANG(a VECTOR_3) float32 {
	la := V3_ABS(a)
	if la > 0.0 {
		// ST: V3_ZANG := ACOS(A.Z / la);
		return float32(math.Acos(float64(a.Z / la)))
	}
	return 0.0
}
