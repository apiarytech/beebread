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

// Package logic is the port of the OSCAT BASIC logic functions: gate logic,
// flip-flops, generators and memory.
package logic

import (
	gomath "math"
	"math/bits"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

func boolByte(b iec.BOOL) iec.BYTE {
	if b {
		return 1
	}
	return 0
}

// BCDC_TO_INT converts a two digit BCD number to an integer.
func BCDC_TO_INT(in iec.BYTE) iec.INT {
	return iec.INT(in&0x0F) + iec.INT(in>>4)*10
}

// BIT_COUNT counts the bits of in that are set.
func BIT_COUNT(in iec.DWORD) iec.INT {
	return iec.INT(bits.OnesCount32(uint32(in)))
}

// BIT_LOAD_B sets bit pos of in to val.
func BIT_LOAD_B(in iec.BYTE, val iec.BOOL, pos iec.INT) iec.BYTE {
	if val {
		return in | SHL(iec.BYTE(1), pos)
	}
	return in &^ SHL(iec.BYTE(1), pos)
}

// BIT_LOAD_B2 sets n bits of i, starting at bit p, to d.
func BIT_LOAD_B2(i iec.BYTE, d iec.BOOL, p, n iec.INT) iec.BYTE {
	if d {
		return ROL(SHR(iec.BYTE(255), 8-n)|ROR(i, p), p)
	}
	return ROL(SHL(iec.BYTE(255), n)&ROR(i, p), p)
}

// BIT_LOAD_DW sets bit pos of in to val.
func BIT_LOAD_DW(in iec.DWORD, val iec.BOOL, pos iec.INT) iec.DWORD {
	if val {
		return in | SHL(iec.DWORD(1), pos)
	}
	return in &^ SHL(iec.DWORD(1), pos)
}

// BIT_LOAD_DW2 sets n bits of i, starting at bit p, to d.
func BIT_LOAD_DW2(i iec.DWORD, d iec.BOOL, p, n iec.INT) iec.DWORD {
	if d {
		return ROL(SHR(iec.DWORD(0xFFFFFFFF), 32-n)|ROR(i, p), p)
	}
	return ROL(SHL(iec.DWORD(0xFFFFFFFF), n)&ROR(i, p), p)
}

// BIT_LOAD_W sets bit pos of in to val.
func BIT_LOAD_W(in iec.WORD, val iec.BOOL, pos iec.INT) iec.WORD {
	if val {
		return in | SHL(iec.WORD(1), pos)
	}
	return in &^ SHL(iec.WORD(1), pos)
}

// BIT_LOAD_W2 sets n bits of i, starting at bit p, to d.
func BIT_LOAD_W2(i iec.WORD, d iec.BOOL, p, n iec.INT) iec.WORD {
	if d {
		return ROL(SHR(iec.WORD(0xFFFF), 16-n)|ROR(i, p), p)
	}
	return ROL(SHL(iec.WORD(0xFFFF), n)&ROR(i, p), p)
}

// BIT_OF_DWORD returns bit n of in; bit 0 is the lowest.
func BIT_OF_DWORD(in iec.DWORD, n iec.INT) iec.BOOL {
	return iec.BOOL(BIT(in, n))
}

// BIT_TOGGLE_B toggles bit pos of in.
func BIT_TOGGLE_B(in iec.BYTE, pos iec.INT) iec.BYTE {
	return SHL(iec.BYTE(1), pos) ^ in
}

// BIT_TOGGLE_DW toggles bit pos of in.
func BIT_TOGGLE_DW(in iec.DWORD, pos iec.INT) iec.DWORD {
	return SHL(iec.DWORD(1), pos) ^ in
}

// BIT_TOGGLE_W toggles bit pos of in.
func BIT_TOGGLE_W(in iec.WORD, pos iec.INT) iec.WORD {
	return SHL(iec.WORD(1), pos) ^ in
}

// BYTE_OF_BIT creates a byte from 8 bits.
func BYTE_OF_BIT(b0, b1, b2, b3, b4, b5, b6, b7 iec.BOOL) iec.BYTE {
	return boolByte(b7)<<7 | boolByte(b6)<<6 | boolByte(b5)<<5 | boolByte(b4)<<4 |
		boolByte(b3)<<3 | boolByte(b2)<<2 | boolByte(b1)<<1 | boolByte(b0)
}

// BYTE_OF_DWORD returns byte n of in; byte 0 is the lowest.
func BYTE_OF_DWORD(in iec.DWORD, n iec.BYTE) iec.BYTE {
	return iec.BYTE(SHR(in, SHL(n, 3)))
}

// BYTE_TO_BITS splits a byte into its 8 bits.
type BYTE_TO_BITS struct {
	IN                             iec.BYTE
	B0, B1, B2, B3, B4, B5, B6, B7 iec.BOOL
}

// INIT resets the block.
func (b *BYTE_TO_BITS) INIT() { *b = BYTE_TO_BITS{} }

// Execute runs the block once.
func (b *BYTE_TO_BITS) Execute(now time.Time) {
	bit := func(n int) iec.BOOL { return b.IN>>n&1 != 0 }
	b.B0, b.B1, b.B2, b.B3 = bit(0), bit(1), bit(2), bit(3)
	b.B4, b.B5, b.B6, b.B7 = bit(4), bit(5), bit(6), bit(7)
}

// BYTE_TO_GRAY converts a binary number to Gray code.
func BYTE_TO_GRAY(in iec.BYTE) iec.BYTE {
	return in ^ in>>1
}

// CHECK_PARITY reports whether in and the parity bit p have even parity.
func CHECK_PARITY(in iec.DWORD, p iec.BOOL) iec.BOOL {
	return iec.BOOL(bits.OnesCount32(uint32(in))%2 == 1) == p
}

// CHK_REAL checks a REAL: 16#00 is a number, 16#20 is +infinity, 16#40 is
// -infinity and 16#80 is not a number.
func CHK_REAL(x iec.REAL) iec.BYTE {
	tmp := bits.RotateLeft32(gomath.Float32bits(float32(x)), 1)
	switch {
	case tmp < 0xFF000000:
		return 0x00
	case tmp == 0xFF000000:
		return 0x20
	case tmp == 0xFF000001:
		return 0x40
	}
	return 0x80
}

// DEC_2 switches D to Q0 if A is false and to Q1 if A is true.
type DEC_2 struct {
	D, A   iec.BOOL
	Q0, Q1 iec.BOOL
}

// INIT resets the block.
func (d *DEC_2) INIT() { *d = DEC_2{} }

// Execute runs the block once.
func (d *DEC_2) Execute(now time.Time) {
	d.Q0 = d.D && !d.A
	d.Q1 = d.D && d.A
}

// DEC_4 switches D to the output Q0..Q3 the address A1, A0 selects.
type DEC_4 struct {
	D, A0, A1      iec.BOOL
	Q0, Q1, Q2, Q3 iec.BOOL
}

// INIT resets the block.
func (d *DEC_4) INIT() { *d = DEC_4{} }

// Execute runs the block once.
func (d *DEC_4) Execute(now time.Time) {
	d.Q0 = d.D && !d.A0 && !d.A1
	d.Q1 = d.D && d.A0 && !d.A1
	d.Q2 = d.D && !d.A0 && d.A1
	d.Q3 = d.D && d.A0 && d.A1
}

// DEC_8 switches D to the output Q0..Q7 the address A2, A1, A0 selects.
type DEC_8 struct {
	D, A0, A1, A2                  iec.BOOL
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7 iec.BOOL
}

// INIT resets the block.
func (d *DEC_8) INIT() { *d = DEC_8{} }

// Execute runs the block once.
func (d *DEC_8) Execute(now time.Time) {
	x := boolByte(d.A0) | boolByte(d.A1)<<1 | boolByte(d.A2)<<2
	q := [8]*iec.BOOL{&d.Q0, &d.Q1, &d.Q2, &d.Q3, &d.Q4, &d.Q5, &d.Q6, &d.Q7}
	for i, p := range q {
		*p = iec.BYTE(i) == x && d.D
	}
}

// DW_TO_REAL returns the REAL whose bits are x.
func DW_TO_REAL(x iec.DWORD) iec.REAL {
	return iec.REAL(gomath.Float32frombits(uint32(x)))
}

// DWORD_OF_BYTE creates a DWORD from 4 bytes; b3 is the highest.
func DWORD_OF_BYTE(b3, b2, b1, b0 iec.BYTE) iec.DWORD {
	return iec.DWORD(b3)<<24 | iec.DWORD(b2)<<16 | iec.DWORD(b1)<<8 | iec.DWORD(b0)
}

// DWORD_OF_WORD creates a DWORD from 2 words; w1 is the high word.
func DWORD_OF_WORD(w1, w0 iec.WORD) iec.DWORD {
	return iec.DWORD(w1)<<16 | iec.DWORD(w0)
}

// GRAY_TO_BYTE converts a Gray code to a binary number.
func GRAY_TO_BYTE(in iec.BYTE) iec.BYTE {
	out := in>>4 ^ in
	out = out>>2 ^ out
	return out>>1 ^ out
}

// INT_TO_BCDC converts an integer 0..99 to a two digit BCD number.
func INT_TO_BCDC(in iec.INT) iec.BYTE {
	return iec.BYTE(in/10)<<4 | iec.BYTE(in%10)
}

// MUX_2 returns D0 if A0 is false and D1 if it is true.
func MUX_2(d0, d1, a0 iec.BOOL) iec.BOOL {
	return SEL(a0, d0, d1)
}

// MUX_4 returns the input D0..D3 the address A1, A0 selects.
func MUX_4(d0, d1, d2, d3, a0, a1 iec.BOOL) iec.BOOL {
	if a1 {
		return SEL(a0, d2, d3)
	}
	return SEL(a0, d0, d1)
}

// PARITY reports whether the number of bits of in that are set is odd.
func PARITY(in iec.DWORD) iec.BOOL {
	return bits.OnesCount32(uint32(in))%2 == 1
}

// REAL_TO_DW returns the bits of the REAL x.
func REAL_TO_DW(x iec.REAL) iec.DWORD {
	return iec.DWORD(gomath.Float32bits(float32(x)))
}

// REFLECT reverses the order of the lowest l bits of d.
func REFLECT(d iec.DWORD, l iec.INT) iec.DWORD {
	var out iec.DWORD
	for i := iec.INT(1); i <= l; i++ {
		out = out<<1 | d&1
		d >>= 1
	}
	return out | SHL(d, l)
}

// REVERSE reverses the order of the bits of a byte.
func REVERSE(in iec.BYTE) iec.BYTE {
	return iec.BYTE(bits.Reverse8(uint8(in)))
}

// SHL1 shifts in left by n bits, filling in 1s.
func SHL1(in iec.DWORD, n iec.INT) iec.DWORD {
	return SHR(iec.DWORD(0xFFFFFFFF), 32-n) | SHL(in, n)
}

// SHR1 shifts in right by n bits, filling in 1s.
func SHR1(in iec.DWORD, n iec.INT) iec.DWORD {
	return SHL(iec.DWORD(0xFFFFFFFF), 32-n) | SHR(in, n)
}

// SWAP_BYTE swaps the high and low byte of a word.
func SWAP_BYTE(in iec.WORD) iec.WORD {
	return ROL(in, 8)
}

// SWAP_BYTE2 reverses the order of the bytes of a DWORD.
func SWAP_BYTE2(in iec.DWORD) iec.DWORD {
	return iec.DWORD(bits.ReverseBytes32(uint32(in)))
}

// WORD_OF_BYTE creates a word from 2 bytes; b1 is the high byte.
func WORD_OF_BYTE(b1, b0 iec.BYTE) iec.WORD {
	return iec.WORD(b1)<<8 | iec.WORD(b0)
}

// WORD_OF_DWORD returns word n of in; word 0 is the low word.
func WORD_OF_DWORD(in iec.DWORD, n iec.BYTE) iec.WORD {
	return iec.WORD(SHR(in, SHL(n, 4)))
}
