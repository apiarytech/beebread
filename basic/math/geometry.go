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

// CIRCLE_A calculates the area of a circle segment with the radius rx and
// the angle ax in degrees; ax = 360 is the whole circle.
func CIRCLE_A(rx, ax iec.REAL) iec.REAL {
	return rx * rx * 8.726646e-3 * ax
}

// CIRCLE_C calculates the length of the arc of a circle with the radius rx
// over the angle ax in degrees; ax = 360 is the whole circumference.
func CIRCLE_C(rx, ax iec.REAL) iec.REAL {
	return 1.7453293e-2 * rx * ax
}

// CIRCLE_SEG calculates the area of the circle segment between a secant
// line and the circumference, for the radius rx and the segment's height hx.
func CIRCLE_SEG(rx, hx iec.REAL) iec.REAL {
	if rx > 0.0 {
		seg := 2.0 * ACOS(1.0-LIMIT(0.0, hx/rx, 2.0))
		return (seg - SIN(seg)) * rx * rx / 2.0
	}
	return 0.0
}

// CONE_V calculates the volume of a cone with the radius rx and height hx.
func CONE_V(rx, hx iec.REAL) iec.REAL {
	return 1.047197551 * rx * rx * hx
}

// ELLIPSE_A calculates the area of an ellipse with the radii r1 and r2.
func ELLIPSE_A(r1, r2 iec.REAL) iec.REAL {
	return MATH.PI * r1 * r2
}

// ELLIPSE_C calculates the circumference of an ellipse with the radii r1
// and r2.
func ELLIPSE_C(r1, r2 iec.REAL) iec.REAL {
	return MATH.PI * (3.0*(r1+r2) - SQRT((3.0*r1+r2)*(3.0*r2+r1)))
}

// SPHERE_V calculates the volume of a sphere with the radius rx.
func SPHERE_V(rx iec.REAL) iec.REAL {
	return 4.188790205 * rx * rx * rx
}

// TRIANGLE_A calculates the area of a triangle from two sides s1 and s2 and
// the angle a between them in degrees, or, if a is 0, from the three sides.
func TRIANGLE_A(s1, a, s2, s3 iec.REAL) iec.REAL {
	if a == 0.0 {
		return SQRT((s1+s2+s3)*(s1+s2-s3)*(s2+s3-s1)*(s3+s1-s2)) * 0.25
	}
	return s1 * s2 * SIN(RAD(a)) * 0.5
}
