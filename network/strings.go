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

package network

import (
	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// STRING_N returns the value a STRING(n) keeps of s: its first n
// characters.
func STRING_N(s iec.STRING, n iec.INT) iec.STRING {
	if LEN(s) <= n {
		return s
	}
	return LEFT(s, n)
}
