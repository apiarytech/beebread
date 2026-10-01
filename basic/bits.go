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

package basic

import "unsafe"

// Unsigned is a bit string type: BYTE, WORD, DWORD or LWORD.
type Unsigned interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

func width[T Unsigned]() uint {
	var x T
	return uint(unsafe.Sizeof(x)) * 8
}

// SHL shifts x left by n bits. Bits shifted in are 0; a shift by the width
// of x or more gives 0, and so does a negative shift.
func SHL[T Unsigned, N Integer](x T, n N) T {
	if n < 0 || uint64(n) >= uint64(width[T]()) {
		return 0
	}
	return x << uint(n)
}

// SHR shifts x right by n bits. Bits shifted in are 0; a shift by the
// width of x or more gives 0, and so does a negative shift.
func SHR[T Unsigned, N Integer](x T, n N) T {
	if n < 0 || uint64(n) >= uint64(width[T]()) {
		return 0
	}
	return x >> uint(n)
}

// ROL rotates x left by n bits.
func ROL[T Unsigned, N Integer](x T, n N) T {
	w := int64(width[T]())
	s := uint((int64(n)%w + w) % w)
	if s == 0 {
		return x
	}
	return x<<s | x>>(uint(w)-s)
}

// ROR rotates x right by n bits.
func ROR[T Unsigned, N Integer](x T, n N) T {
	return ROL(x, -int64(n))
}

// BIT reports whether bit n of x is set, as x.n does in ST.
func BIT[T Unsigned, N Integer](x T, n N) bool {
	return SHR(x, n)&1 != 0
}
