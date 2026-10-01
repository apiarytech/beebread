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

// V3_ABS calculates the length of a vector.
func V3_ABS(a VECTOR_3) iec.REAL {
	return SQRT(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// V3_ADD adds two vectors.
func V3_ADD(a, b VECTOR_3) VECTOR_3 {
	return VECTOR_3{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

// V3_ANG calculates the angle between two vectors, in rad.
func V3_ANG(a, b VECTOR_3) iec.REAL {
	d := V3_ABS(a) * V3_ABS(b)
	if d > 0 {
		return ACOS(LIMIT(-1.0, V3_DPRO(a, b)/d, 1.0))
	}
	return 0.0
}

// V3_DPRO calculates the dot product of two vectors.
func V3_DPRO(a, b VECTOR_3) iec.REAL {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// V3_NORM returns the vector of length 1 with the direction of a.
func V3_NORM(a VECTOR_3) VECTOR_3 {
	la := V3_ABS(a)
	if la > 0.0 {
		return V3_SMUL(a, 1.0/la)
	}
	return VECTOR_3{}
}

// V3_NUL reports whether a vector is the null vector.
func V3_NUL(a VECTOR_3) iec.BOOL {
	return a.X == 0.0 && a.Y == 0.0 && a.Z == 0.0
}

// V3_PAR reports whether two vectors are parallel.
func V3_PAR(a, b VECTOR_3) iec.BOOL {
	return V3_ABS(V3_XPRO(a, b)) == 0.0
}

// V3_REV reverses a vector.
func V3_REV(a VECTOR_3) VECTOR_3 {
	return VECTOR_3{X: -a.X, Y: -a.Y, Z: -a.Z}
}

// V3_SMUL multiplies a vector with the scalar m.
func V3_SMUL(a VECTOR_3, m iec.REAL) VECTOR_3 {
	return VECTOR_3{X: a.X * m, Y: a.Y * m, Z: a.Z * m}
}

// V3_SUB subtracts the vector b from a.
func V3_SUB(a, b VECTOR_3) VECTOR_3 {
	return VECTOR_3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

// V3_XANG calculates the angle between the X axis and a vector, in rad.
func V3_XANG(a VECTOR_3) iec.REAL {
	la := V3_ABS(a)
	if la > 0.0 {
		return ACOS(a.X / la)
	}
	return 0.0
}

// V3_XPRO calculates the cross product of two vectors.
func V3_XPRO(a, b VECTOR_3) VECTOR_3 {
	return VECTOR_3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

// V3_YANG calculates the angle between the Y axis and a vector, in rad.
func V3_YANG(a VECTOR_3) iec.REAL {
	la := V3_ABS(a)
	if la > 0.0 {
		return ACOS(a.Y / la)
	}
	return 0.0
}

// V3_ZANG calculates the angle between the Z axis and a vector, in rad.
func V3_ZANG(a VECTOR_3) iec.REAL {
	la := V3_ABS(a)
	if la > 0.0 {
		return ACOS(a.Z / la)
	}
	return 0.0
}
