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

package logic

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// CRC_GEN calculates the CRC of the first size bytes of pt with the
// polynomial pn of length pl, whose highest bit is left out: x4 + x + 1 is
// 2#0011 with length 4. init is the start value and xorOut is XORed with the
// result. revIn reflects each input byte and revOut the result. A message
// shorter than 4 bytes is filled with 0s at the beginning.
func CRC_GEN(pt []iec.BYTE, size, pl iec.INT, pn, init iec.DWORD, revIn, revOut iec.BOOL, xorOut iec.DWORD) iec.DWORD {
	at := func(i iec.INT) iec.BYTE {
		if int(i) < len(pt) && i < size {
			return pt[i]
		}
		return 0
	}
	shift := 32 - pl
	pn = SHL(pn, shift)
	var crc iec.DWORD
	for pos := iec.INT(0); pos <= 3; pos++ {
		if revIn {
			crc = crc<<8 | iec.DWORD(REVERSE(at(pos)))
		} else {
			crc = crc<<8 | iec.DWORD(at(pos))
		}
	}
	crc = crc ^ SHL(init, shift)
	for pos := iec.INT(4); pos < size; pos++ {
		dx := at(pos)
		if revIn {
			dx = REVERSE(dx)
		}
		for bit := 0; bit < 8; bit++ {
			in := iec.DWORD(dx >> 7)
			if crc&0x80000000 != 0 {
				crc = (crc<<1 | in) ^ pn
			} else {
				crc = crc<<1 | in
			}
			dx <<= 1
		}
	}
	// Finish the register's 32 bits.
	for bit := 0; bit < 32; bit++ {
		if crc&0x80000000 != 0 {
			crc = crc<<1 ^ pn
		} else {
			crc = crc << 1
		}
	}
	crc = SHR(crc, shift) ^ xorOut
	if revOut {
		crc = REFLECT(crc, pl)
	}
	return crc
}

// MATRIX encodes a keyboard matrix of 4 rows and up to 5 columns. It drives
// one row, Y1..Y4, each scan and reads the columns on X1..X5. When a key is
// pressed, TP is true for one scan and CODE holds the column 1..5 in bits
// 0..2, the row in bits 4..6 and 1 in bit 7. When RELEASE is true, the
// release of a key sends its code with bit 7 = 0.
type MATRIX struct {
	X1, X2, X3, X4, X5 iec.BOOL
	RELEASE            iec.BOOL
	CODE               iec.BYTE
	TP                 iec.BOOL
	Y1                 iec.BOOL // default TRUE
	Y2, Y3, Y4         iec.BOOL

	line        iec.BYTE
	x           [4]iec.BYTE // scan line inputs
	l           [4]iec.BYTE // scan line status
	initialized bool
}

// INIT resets the block and sets Y1 to its initial value.
func (m *MATRIX) INIT() { *m = MATRIX{Y1: true, initialized: true} }

// Execute runs the block once.
func (m *MATRIX) Execute(now time.Time) {
	if !m.initialized {
		m.initialized = true
		m.Y1 = true
	}
	m.TP = false
	m.CODE = 0
	m.x[m.line] = BYTE_OF_BIT(m.X1, m.X2, m.X3, m.X4, m.X5, false, false, false) |
		m.x[m.line]&0xE0
	for i := 0; i <= 3; i++ {
		if m.x[i] == m.l[i] {
			continue
		}
		// The scan line has changed: find and send the code.
		temp := m.x[i] ^ m.l[i]
		for bit := iec.BYTE(0); bit <= 4; bit++ {
			if temp>>bit&1 != 0 {
				mask := iec.BYTE(1) << bit
				m.CODE = bit + 1 | m.x[i]&mask>>bit<<7
				m.l[i] = m.l[i]&^mask | m.x[i]&mask
				break
			}
		}
		m.TP = true
		m.CODE = m.CODE&0x8F | (m.line&7)<<4
		if !m.RELEASE && m.CODE < 127 {
			m.CODE = 0
			m.TP = false
		}
		break
	}
	m.line = (m.line + 1) & 3
	temp := iec.BYTE(1) << m.line
	m.Y1, m.Y2, m.Y3, m.Y4 = temp&1 != 0, temp&2 != 0, temp&4 != 0, temp&8 != 0
}

// PIN_CODE sets TP for one scan when the codes CB, taken each scan E is
// true, spell PIN.
type PIN_CODE struct {
	CB  iec.BYTE
	E   iec.BOOL
	PIN iec.STRING // STRING(8)
	TP  iec.BOOL

	pos         iec.INT
	initialized bool
}

// INIT resets the block.
func (p *PIN_CODE) INIT() { *p = PIN_CODE{pos: 1, initialized: true} }

// Execute runs the block once.
func (p *PIN_CODE) Execute(now time.Time) {
	if !p.initialized {
		p.initialized = true
		p.pos = 1
	}
	p.TP = false
	if !p.E {
		return
	}
	pin := CHARS(p.PIN)
	if int(p.pos) <= len(pin) && p.CB == pin[p.pos-1] {
		p.pos++
		if int(p.pos) > len(pin) {
			p.TP = true
			p.pos = 1
		}
	} else {
		p.pos = 1
	}
}
