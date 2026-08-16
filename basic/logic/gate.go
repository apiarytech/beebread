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

package logic

import (
	"math"
	"math/bits"
)

// BCDC_TO_INT converts a two-digit BCD number into an integer.
func BCDC_TO_INT(in byte) int {
	return int((in & 0x0F) + ((in>>4)&0x0F)*10)
}

// BIT_COUNT counts the number of set bits (1s) in a dword.
func BIT_COUNT(in uint32) int {
	return bits.OnesCount32(in)
}

// BIT_LOAD_B sets or clears a bit at a specific position in a byte.
func BIT_LOAD_B(in byte, val bool, pos uint) byte {
	if pos >= 8 {
		return in
	}
	if val {
		return in | (1 << pos)
	}
	return in &^ (1 << pos)
}

// BIT_LOAD_B2 sets or clears N bits starting at position P in a byte.
func BIT_LOAD_B2(i byte, d bool, p, n uint) byte {
	if p >= 8 || n == 0 {
		return i
	}
	if p+n > 8 {
		n = 8 - p
	}
	mask := byte(((1 << n) - 1) << p)
	if d {
		return i | mask
	}
	return i &^ mask
}

// BIT_LOAD_DW sets or clears a bit at a specific position in a dword.
func BIT_LOAD_DW(in uint32, val bool, pos uint) uint32 {
	if pos >= 32 {
		return in
	}
	if val {
		return in | (1 << pos)
	}
	return in &^ (1 << pos)
}

// BIT_LOAD_DW2 sets or clears N bits starting at position P in a dword.
func BIT_LOAD_DW2(i uint32, d bool, p, n uint) uint32 {
	if p >= 32 || n == 0 {
		return i
	}
	if p+n > 32 {
		n = 32 - p
	}
	mask := uint32(((1 << n) - 1) << p)
	if d {
		return i | mask
	}
	return i &^ mask
}

// BIT_LOAD_W sets or clears a bit at a specific position in a word.
func BIT_LOAD_W(in uint16, val bool, pos uint) uint16 {
	if pos >= 16 {
		return in
	}
	if val {
		return in | (1 << pos)
	}
	return in &^ (1 << pos)
}

// BIT_LOAD_W2 sets or clears N bits starting at position P in a word.
func BIT_LOAD_W2(i uint16, d bool, p, n uint) uint16 {
	if p >= 16 || n == 0 {
		return i
	}
	if p+n > 16 {
		n = 16 - p
	}
	mask := uint16(((1 << n) - 1) << p)
	if d {
		return i | mask
	}
	return i &^ mask
}

// BIT_OF_DWORD extracts a single bit from the Nth position of a dword.
func BIT_OF_DWORD(in uint32, n uint) bool {
	if n >= 32 {
		return false
	}
	return (in>>n)&1 == 1
}

// BIT_TOGGLE_B toggles a bit of a byte at a specific position.
func BIT_TOGGLE_B(in byte, pos uint) byte {
	if pos >= 8 {
		return in
	}
	return in ^ (1 << pos)
}

// BIT_TOGGLE_DW toggles a bit of a dword at a specific position.
func BIT_TOGGLE_DW(in uint32, pos uint) uint32 {
	if pos >= 32 {
		return in
	}
	return in ^ (1 << pos)
}

// BIT_TOGGLE_W toggles a bit of a word at a specific position.
func BIT_TOGGLE_W(in uint16, pos uint) uint16 {
	if pos >= 16 {
		return in
	}
	return in ^ (1 << pos)
}

// BYTE_OF_BIT creates a byte from 8 individual bits.
func BYTE_OF_BIT(b7, b6, b5, b4, b3, b2, b1, b0 bool) byte {
	var res byte
	if b0 {
		res |= 1 << 0
	}
	if b1 {
		res |= 1 << 1
	}
	if b2 {
		res |= 1 << 2
	}
	if b3 {
		res |= 1 << 3
	}
	if b4 {
		res |= 1 << 4
	}
	if b5 {
		res |= 1 << 5
	}
	if b6 {
		res |= 1 << 6
	}
	if b7 {
		res |= 1 << 7
	}
	return res
}

// BYTE_OF_DWORD extracts the Nth byte from a dword. N=0 is the least significant byte.
func BYTE_OF_DWORD(in uint32, n uint) byte {
	if n >= 4 {
		return 0
	}
	return byte(in >> (n * 8))
}

// BYTE_TO_BITS extracts the 8 bits from a byte.
func BYTE_TO_BITS(in byte) (b0, b1, b2, b3, b4, b5, b6, b7 bool) {
	return (in>>0)&1 == 1, (in>>1)&1 == 1, (in>>2)&1 == 1, (in>>3)&1 == 1,
		(in>>4)&1 == 1, (in>>5)&1 == 1, (in>>6)&1 == 1, (in>>7)&1 == 1
}

// BYTE_TO_GRAY converts a binary byte to Gray code.
func BYTE_TO_GRAY(in byte) byte {
	return in ^ (in >> 1)
}

// DW_TO_REAL converts a dword to a real (float32) bitwise.
func DW_TO_REAL(x uint32) float32 {
	return math.Float32frombits(x)
}

// DWORD_OF_BYTE creates a dword from 4 individual bytes.
func DWORD_OF_BYTE(b3, b2, b1, b0 byte) uint32 {
	return uint32(b3)<<24 | uint32(b2)<<16 | uint32(b1)<<8 | uint32(b0)
}

// DWORD_OF_WORD creates a dword from 2 individual words.
func DWORD_OF_WORD(w1, w0 uint16) uint32 {
	return uint32(w1)<<16 | uint32(w0)
}

// GRAY_TO_BYTE converts a Gray code byte to binary.
func GRAY_TO_BYTE(in byte) byte {
	in ^= in >> 4
	in ^= in >> 2
	in ^= in >> 1
	return in
}

// INT_TO_BCDC converts an integer (0-99) into a two-digit BCD number.
func INT_TO_BCDC(in int) byte {
	if in < 0 || in > 99 {
		return 0
	}
	return byte(((in / 10) << 4) | (in % 10))
}

// REAL_TO_DW converts a real (float32) to a dword bitwise.
func REAL_TO_DW(x float32) uint32 {
	return math.Float32bits(x)
}

// Reflect reverses the order of the L least significant bits in a dword.
func Reflect(d uint32, l uint) uint32 {
	if l == 0 || l > 32 {
		return d
	}
	var res uint32
	mask := (uint32(1) << l) - 1
	val := d & mask
	for i := uint(0); i < l; i++ {
		if (val & (1 << i)) != 0 {
			res |= 1 << (l - 1 - i)
		}
	}
	// Combine with the un-reflected part of the original dword
	return (d &^ mask) | res
}

// REVERSE reverses the bits of a byte.
func REVERSE(in byte) byte {
	return bits.Reverse8(in)
}

// SHL1 shifts N bits to the left, filling the new bits with 1s.
func SHL1(in uint32, n uint) uint32 {
	if n >= 32 {
		return 0xFFFFFFFF
	}
	return (in << n) | ((1 << n) - 1)
}

// SHR1 shifts N bits to the right, filling the new bits with 1s (arithmetic-like shift).
func SHR1(in uint32, n uint) uint32 {
	if n >= 32 {
		return 0xFFFFFFFF
	}
	return (in >> n) | (0xFFFFFFFF << (32 - n))
}

// SWAP_BYTE swaps the high and low byte of a word.
func SWAP_BYTE(in uint16) uint16 {
	return bits.ReverseBytes16(in)
}

// SWAP_BYTE2 reverses the byte order in a dword (e.g., from little-endian to big-endian).
func SWAP_BYTE2(in uint32) uint32 {
	return bits.ReverseBytes32(in)
}

// WORD_OF_BYTE creates a word from 2 individual bytes.
func WORD_OF_BYTE(b1, b0 byte) uint16 {
	return uint16(b1)<<8 | uint16(b0)
}

// WORD_OF_DWORD extracts the Nth word from a dword. N=0 is the least significant word.
func WORD_OF_DWORD(in uint32, n uint) uint16 {
	if n >= 2 {
		return 0
	}
	return uint16(in >> (n * 16))
}

// CHECK_PARITY checks for even parity for a dword and a parity bit.
func CHECK_PARITY(in uint32, p bool) bool {
	// The original ST code has a complex loop. A more direct way is to count
	// the bits and check the parity of the count.
	// Even parity means the total number of 1s (including the parity bit) is even.
	count := bits.OnesCount32(in)
	if p {
		count++
	}
	return count%2 == 0
}

// CHK_REAL checks a float32 for NaN and infinity.
// Returns: 0 for normal, 20 for +Inf, 40 for -Inf, 80 for NaN.
func CHK_REAL(x float32) byte {
	bits := math.Float32bits(x)
	// Exponent is all 1s (0xFF)
	if (bits>>23)&0xFF == 0xFF {
		// Mantissa is 0 -> Infinity
		if bits&0x7FFFFF == 0 {
			// Sign bit is 0 -> +Inf
			if bits>>31 == 0 {
				return 20
			}
			// Sign bit is 1 -> -Inf
			return 40
		}
		// Mantissa is non-zero -> NaN
		return 80
	}
	// Normal or denormalized number
	return 0
}

// DEC2 is a 1-to-2 decoder.
// If D is true, Q0 is true if A is false, and Q1 is true if A is true.
func DEC2(d, a bool) (q0, q1 bool) {
	return d && !a, d && a
}

// DEC4 is a 2-to-4 decoder.
func DEC4(d, a0, a1 bool) (q0, q1, q2, q3 bool) {
	return d && !a0 && !a1, d && a0 && !a1, d && !a0 && a1, d && a0 && a1
}

// DEC8 is a 3-to-8 decoder.
func DEC8(d, a0, a1, a2 bool) [8]bool {
	var q [8]bool
	if d {
		var sel byte
		if a0 {
			sel |= 1
		}
		if a1 {
			sel |= 2
		}
		if a2 {
			sel |= 4
		}
		if sel < 8 {
			q[sel] = true
		}
	}
	return q
}

// MUX2 is a 2-to-1 multiplexer.
func MUX2(d0, d1, a0 bool) bool {
	if a0 {
		return d1
	}
	return d0
}

// MUX4 is a 4-to-1 multiplexer.
func MUX4(d0, d1, d2, d3, a0, a1 bool) bool {
	if a1 {
		return MUX2(d2, d3, a0)
	}
	return MUX2(d0, d1, a0)
}

// PARITY calculates the odd parity of a dword.
// Returns true if the number of set bits is odd.
func PARITY(in uint32) bool {
	return bits.OnesCount32(in)%2 != 0
}
