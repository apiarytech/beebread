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

// BcdcToInt converts a two-digit BCD number into an integer.
func BcdcToInt(in byte) int {
	return int((in & 0x0F) + ((in>>4)&0x0F)*10)
}

// BitCount counts the number of set bits (1s) in a dword.
func BitCount(in uint32) int {
	return bits.OnesCount32(in)
}

// BitLoadB sets or clears a bit at a specific position in a byte.
func BitLoadB(in byte, val bool, pos uint) byte {
	if pos >= 8 {
		return in
	}
	if val {
		return in | (1 << pos)
	}
	return in &^ (1 << pos)
}

// BitLoadB2 sets or clears N bits starting at position P in a byte.
func BitLoadB2(i byte, d bool, p, n uint) byte {
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

// BitLoadDw sets or clears a bit at a specific position in a dword.
func BitLoadDw(in uint32, val bool, pos uint) uint32 {
	if pos >= 32 {
		return in
	}
	if val {
		return in | (1 << pos)
	}
	return in &^ (1 << pos)
}

// BitLoadDw2 sets or clears N bits starting at position P in a dword.
func BitLoadDw2(i uint32, d bool, p, n uint) uint32 {
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

// BitLoadW sets or clears a bit at a specific position in a word.
func BitLoadW(in uint16, val bool, pos uint) uint16 {
	if pos >= 16 {
		return in
	}
	if val {
		return in | (1 << pos)
	}
	return in &^ (1 << pos)
}

// BitLoadW2 sets or clears N bits starting at position P in a word.
func BitLoadW2(i uint16, d bool, p, n uint) uint16 {
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

// BitOfDword extracts a single bit from the Nth position of a dword.
func BitOfDword(in uint32, n uint) bool {
	if n >= 32 {
		return false
	}
	return (in>>n)&1 == 1
}

// BitToggleB toggles a bit of a byte at a specific position.
func BitToggleB(in byte, pos uint) byte {
	if pos >= 8 {
		return in
	}
	return in ^ (1 << pos)
}

// BitToggleDw toggles a bit of a dword at a specific position.
func BitToggleDw(in uint32, pos uint) uint32 {
	if pos >= 32 {
		return in
	}
	return in ^ (1 << pos)
}

// BitToggleW toggles a bit of a word at a specific position.
func BitToggleW(in uint16, pos uint) uint16 {
	if pos >= 16 {
		return in
	}
	return in ^ (1 << pos)
}

// ByteOfBit creates a byte from 8 individual bits.
func ByteOfBit(b7, b6, b5, b4, b3, b2, b1, b0 bool) byte {
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

// ByteOfDword extracts the Nth byte from a dword. N=0 is the least significant byte.
func ByteOfDword(in uint32, n uint) byte {
	if n >= 4 {
		return 0
	}
	return byte(in >> (n * 8))
}

// ByteToBits extracts the 8 bits from a byte.
func ByteToBits(in byte) (b0, b1, b2, b3, b4, b5, b6, b7 bool) {
	return (in>>0)&1 == 1, (in>>1)&1 == 1, (in>>2)&1 == 1, (in>>3)&1 == 1,
		(in>>4)&1 == 1, (in>>5)&1 == 1, (in>>6)&1 == 1, (in>>7)&1 == 1
}

// ByteToGray converts a binary byte to Gray code.
func ByteToGray(in byte) byte {
	return in ^ (in >> 1)
}

// DwToReal converts a dword to a real (float32) bitwise.
func DwToReal(x uint32) float32 {
	return math.Float32frombits(x)
}

// DwordOfByte creates a dword from 4 individual bytes.
func DwordOfByte(b3, b2, b1, b0 byte) uint32 {
	return uint32(b3)<<24 | uint32(b2)<<16 | uint32(b1)<<8 | uint32(b0)
}

// DwordOfWord creates a dword from 2 individual words.
func DwordOfWord(w1, w0 uint16) uint32 {
	return uint32(w1)<<16 | uint32(w0)
}

// GrayToByte converts a Gray code byte to binary.
func GrayToByte(in byte) byte {
	in ^= in >> 4
	in ^= in >> 2
	in ^= in >> 1
	return in
}

// IntToBcdc converts an integer (0-99) into a two-digit BCD number.
func IntToBcdc(in int) byte {
	if in < 0 || in > 99 {
		return 0
	}
	return byte(((in / 10) << 4) | (in % 10))
}

// RealToDw converts a real (float32) to a dword bitwise.
func RealToDw(x float32) uint32 {
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

// Reverse reverses the bits of a byte.
func Reverse(in byte) byte {
	return bits.Reverse8(in)
}

// Shl1 shifts N bits to the left, filling the new bits with 1s.
func Shl1(in uint32, n uint) uint32 {
	if n >= 32 {
		return 0xFFFFFFFF
	}
	return (in << n) | ((1 << n) - 1)
}

// Shr1 shifts N bits to the right, filling the new bits with 1s (arithmetic-like shift).
func Shr1(in uint32, n uint) uint32 {
	if n >= 32 {
		return 0xFFFFFFFF
	}
	return (in >> n) | (0xFFFFFFFF << (32 - n))
}

// SwapByte swaps the high and low byte of a word.
func SwapByte(in uint16) uint16 {
	return bits.ReverseBytes16(in)
}

// SwapByte2 reverses the byte order in a dword (e.g., from little-endian to big-endian).
func SwapByte2(in uint32) uint32 {
	return bits.ReverseBytes32(in)
}

// WordOfByte creates a word from 2 individual bytes.
func WordOfByte(b1, b0 byte) uint16 {
	return uint16(b1)<<8 | uint16(b0)
}

// WordOfDword extracts the Nth word from a dword. N=0 is the least significant word.
func WordOfDword(in uint32, n uint) uint16 {
	if n >= 2 {
		return 0
	}
	return uint16(in >> (n * 16))
}

// CheckParity checks for even parity for a dword and a parity bit.
func CheckParity(in uint32, p bool) bool {
	// The original ST code has a complex loop. A more direct way is to count
	// the bits and check the parity of the count.
	// Even parity means the total number of 1s (including the parity bit) is even.
	count := bits.OnesCount32(in)
	if p {
		count++
	}
	return count%2 == 0
}

// ChkReal checks a float32 for NaN and infinity.
// Returns: 0 for normal, 20 for +Inf, 40 for -Inf, 80 for NaN.
func ChkReal(x float32) byte {
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

// Dec2 is a 1-to-2 decoder.
// If D is true, Q0 is true if A is false, and Q1 is true if A is true.
func Dec2(d, a bool) (q0, q1 bool) {
	return d && !a, d && a
}

// Dec4 is a 2-to-4 decoder.
func Dec4(d, a0, a1 bool) (q0, q1, q2, q3 bool) {
	return d && !a0 && !a1, d && a0 && !a1, d && !a0 && a1, d && a0 && a1
}

// Dec8 is a 3-to-8 decoder.
func Dec8(d, a0, a1, a2 bool) [8]bool {
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

// Mux2 is a 2-to-1 multiplexer.
func Mux2(d0, d1, a0 bool) bool {
	if a0 {
		return d1
	}
	return d0
}

// Mux4 is a 4-to-1 multiplexer.
func Mux4(d0, d1, d2, d3, a0, a1 bool) bool {
	if a1 {
		return Mux2(d2, d3, a0)
	}
	return Mux2(d0, d1, a0)
}

// Parity calculates the odd parity of a dword.
// Returns true if the number of set bits is odd.
func Parity(in uint32) bool {
	return bits.OnesCount32(in)%2 != 0
}
