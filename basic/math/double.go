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

// R2_ABS returns the absolute value of a double precision real.
func R2_ABS(x REAL2) REAL2 {
	if x.RX >= 0.0 {
		return REAL2{RX: x.RX, R1: x.R1}
	}
	return REAL2{RX: -x.RX, R1: -x.R1}
}

// R2_ADD adds a real to a double precision real, which extends the accuracy
// of a real to twice as many digits.
func R2_ADD(x REAL2, y iec.REAL) REAL2 {
	var out REAL2
	temp := x.RX
	out.RX = y + x.R1 + x.RX
	out.R1 = temp - out.RX + y + x.R1
	return out
}

// R2_ADD2 adds two double precision reals.
func R2_ADD2(x, y REAL2) REAL2 {
	return REAL2{R1: x.R1 + y.R1, RX: x.RX + y.RX}
}

// R2_MUL multiplies a double precision real with a real.
func R2_MUL(x REAL2, y iec.REAL) REAL2 {
	return REAL2{RX: x.RX * y, R1: x.R1 * y}
}

// R2_SET sets a double precision real to a real value.
func R2_SET(x iec.REAL) REAL2 {
	return REAL2{RX: x, R1: 0.0}
}
