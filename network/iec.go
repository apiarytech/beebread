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

// The IEC 61131-3 conversions between numbers and strings the library uses
// that basic does not have. A number is written in decimal; a string is
// read as STRING_TO_DINT reads it, a decimal number after blanks, or 0, and
// the result is truncated to the type.

import (
	"strconv"
	"strings"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

func dec(v uint64) iec.STRING { return iec.STRING(strconv.FormatUint(v, 10)) }

// BYTE_TO_STRING returns b in decimal.
func BYTE_TO_STRING(b iec.BYTE) iec.STRING { return dec(uint64(b)) }

// WORD_TO_STRING returns w in decimal.
func WORD_TO_STRING(w iec.WORD) iec.STRING { return dec(uint64(w)) }

// UINT_TO_STRING returns u in decimal.
func UINT_TO_STRING(u iec.UINT) iec.STRING { return dec(uint64(u)) }

// UDINT_TO_STRING returns u in decimal.
func UDINT_TO_STRING(u iec.UDINT) iec.STRING { return dec(uint64(u)) }

// USINT_TO_STRING returns u in decimal.
func USINT_TO_STRING(u iec.USINT) iec.STRING { return dec(uint64(u)) }

// stringToUnsigned reads the decimal number at the start of s, after
// blanks, or 0; it may be larger than a DINT.
func stringToUnsigned(s iec.STRING) uint64 {
	t := strings.TrimLeft(string(s), " ")
	if strings.HasPrefix(t, "-") {
		return uint64(int64(STRING_TO_DINT(s)))
	}
	t = strings.TrimPrefix(t, "+")
	n := 0
	for n < len(t) && t[n] >= '0' && t[n] <= '9' {
		n++
	}
	v, _ := strconv.ParseUint(t[:n], 10, 64)
	return v
}

// STRING_TO_BYTE reads a BYTE from s.
func STRING_TO_BYTE(s iec.STRING) iec.BYTE { return iec.BYTE(stringToUnsigned(s)) }

// STRING_TO_WORD reads a WORD from s.
func STRING_TO_WORD(s iec.STRING) iec.WORD { return iec.WORD(stringToUnsigned(s)) }

// STRING_TO_DWORD reads a DWORD from s.
func STRING_TO_DWORD(s iec.STRING) iec.DWORD { return iec.DWORD(stringToUnsigned(s)) }

// STRING_TO_UINT reads a UINT from s.
func STRING_TO_UINT(s iec.STRING) iec.UINT { return iec.UINT(stringToUnsigned(s)) }

// STRING_TO_UDINT reads a UDINT from s.
func STRING_TO_UDINT(s iec.STRING) iec.UDINT { return iec.UDINT(stringToUnsigned(s)) }

// BYTES returns the characters of s as ISO 8859-1 codes, followed by a 0
// and padded with 0 to n bytes at least, as a STRING is in memory.
func BYTES(s iec.STRING, n int) []iec.BYTE {
	b := append(CHARS(s), 0)
	for len(b) < n {
		b = append(b, 0)
	}
	return b
}
