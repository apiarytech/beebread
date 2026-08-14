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
	. "beebread/basic"
	"math"
	"time"
)

// CrcGen generates a CRC checksum from a block of data.
// The CRC Polynom is specified with PN and the length of the Polynom is specified by PL.
// A Polynom x^4 + x + 1 is represented by 0x3 with length 4.
func CrcGen(data []byte, pl int, pn, init, xorOut uint32, revIn, revOut bool) uint32 {
	if pl <= 0 || pl > 32 {
		return 0
	}

	var crc uint32
	size := len(data)
	shift := uint(32 - pl)
	poly := pn << shift

	// Load first bytes into register
	for i := 0; i < 4 && i < size; i++ {
		var d byte
		if revIn {
			d = Reverse(data[i])
		} else {
			d = data[i]
		}
		crc = (crc << 8) | uint32(d)
	}

	// XOR with init value
	crc ^= (init << shift)

	// Process remaining bytes
	for i := 4; i < size; i++ {
		var d byte
		if revIn {
			d = Reverse(data[i])
		} else {
			d = data[i]
		}
		for j := 0; j < 8; j++ {
			if (crc & 0x80000000) != 0 {
				crc = (crc << 1) | uint32((d>>7)&1) ^ poly
			} else {
				crc = (crc << 1) | uint32((d>>7)&1)
			}
			d <<= 1
		}
	}

	// Finish the register's 32 bits
	for i := 0; i < 32; i++ {
		if (crc & 0x80000000) != 0 {
			crc = (crc << 1) ^ poly
		} else {
			crc <<= 1
		}
	}

	// Final XOR and shift
	crc = (crc >> shift) ^ xorOut

	// Reverse output if necessary
	if revOut {
		crc = Reflect(crc, uint(pl))
	}

	return crc
}

// Matrix is a matrix keyboard encoder for 4 rows and up to 5 columns.
type Matrix struct {
	Code byte
	TP   bool
	Y    [4]bool

	// internal state
	line byte
	x    [4]byte // scan line inputs
	l    [4]byte // scan line status
}

// Update executes the matrix scan logic for one cycle.
func (m *Matrix) Update(x1, x2, x3, x4, x5, release bool) {
	m.TP = false
	m.Code = 0

	// Read scan lines
	var currentX byte
	if x1 {
		currentX |= 1 << 0
	}
	if x2 {
		currentX |= 1 << 1
	}
	if x3 {
		currentX |= 1 << 2
	}
	if x4 {
		currentX |= 1 << 3
	}
	if x5 {
		currentX |= 1 << 4
	}
	m.x[m.line] = currentX

	// Compare for change
	for i := 0; i < 4; i++ {
		if m.x[i] != m.l[i] {
			diff := m.x[i] ^ m.l[i]
			var col byte
			// Find which bit changed
			for j := 0; j < 5; j++ {
				if (diff>>j)&1 != 0 {
					col = byte(j + 1)
					break
				}
			}

			if col > 0 {
				m.Code = col
				isPressed := (m.x[i]>>(col-1))&1 != 0
				m.Code = BitLoadB(m.Code, isPressed, 7)
				m.l[i] = BitLoadB(m.l[i], isPressed, uint(col-1))

				m.TP = true
				m.Code = BitLoadB(m.Code, (byte(i)>>0)&1 != 0, 4)
				m.Code = BitLoadB(m.Code, (byte(i)>>1)&1 != 0, 5)
				m.Code = BitLoadB(m.Code, (byte(i)>>2)&1 != 0, 6)

				if !release && !isPressed {
					m.Code = 0
					m.TP = false
				}
				goto end_loop // Exit after finding the first change
			}
		}
	}
end_loop:

	// Increment scan line for the next cycle
	m.line = (m.line + 1) & 0x03
	temp := byte(1 << m.line)
	m.Y[0] = (temp & 0x01) != 0
	m.Y[1] = (temp & 0x02) != 0
	m.Y[2] = (temp & 0x04) != 0
	m.Y[3] = (temp & 0x08) != 0
}

// PinCode scans the input of a keypad (Matrix) for a sequence of characters.
type PinCode struct {
	TP bool
	// internal state
	pos int
}

// Update executes the pin code checking logic.
func (p *PinCode) Update(cb byte, e bool, pin string) {
	p.TP = false
	if e {
		if p.pos < len(pin) && cb == pin[p.pos] {
			p.pos++
			if p.pos == len(pin) {
				p.TP = true
				p.pos = 0 // Reset for next attempt
			}
		} else {
			p.pos = 0 // Reset on wrong character
		}
	}
}

// EsrCollect collects ESR data from up to 8 ESR_MON modules and stores them in an output array.
type EsrCollect struct {
	EsrOut [32]EsrData
	pos    int
	cnt    int
}

// Update executes the ESR collection logic.
func (e *EsrCollect) Update(rst bool, esrIn ...[]EsrData) int {
	if rst || e.cnt < 0 {
		e.pos = -1
		e.cnt = 0 // Set to 0 to allow processing
	} else {
		for _, esrArray := range esrIn {
			for _, esrItem := range esrArray {
				if esrItem.Typ > 0 {
					e.pos = (e.pos + 1) % 32 // Wrap around using modulo
					e.EsrOut[e.pos] = esrItem
				}
			}
		}
	}
	return e.pos
}

// EsrMonB8 monitors up to 8 binary inputs and reports changes with a timestamp and address label.
type EsrMonB8 struct {
	EsrFlag bool
	EsrOut  [4]EsrData

	// internal state
	lastState [8]bool
}

// Update executes the monitoring logic.
func (e *EsrMonB8) Update(dtIn time.Time, s [8]bool, a [8]string) {
	e.EsrFlag = false
	// Clear previous output
	e.EsrOut = [4]EsrData{}
	cnt := 0

	for i := 0; i < 8 && cnt < 4; i++ {
		if s[i] != e.lastState[i] {
			e.EsrOut[cnt].Typ = 10 + boolToByte(s[i])
			e.EsrOut[cnt].Adress = a[i]
			e.EsrOut[cnt].Ds = dtIn
			e.EsrOut[cnt].Ts = time.Duration(time.Now().UnixNano())
			e.lastState[i] = s[i]
			cnt++
			e.EsrFlag = true
		}
	}
}

// EsrMonR4 monitors up to 4 real inputs and reports changes with a timestamp and address label.
type EsrMonR4 struct {
	EsrFlag bool
	EsrOut  [4]EsrData

	// internal state
	lastState [4]float32
}

// Update executes the monitoring logic.
func (e *EsrMonR4) Update(dtIn time.Time, r [4]float32, a [4]string, s [4]float32) {
	e.EsrFlag = false
	// Clear previous output
	e.EsrOut = [4]EsrData{}
	cnt := 0

	for i := 0; i < 4 && cnt < 4; i++ {
		if math.Abs(float64(r[i]-e.lastState[i])) > float64(s[i]) {
			e.EsrOut[cnt].Typ = 20
			e.EsrOut[cnt].Adress = a[i]
			e.EsrOut[cnt].Ds = dtIn
			e.EsrOut[cnt].Ts = time.Duration(time.Now().UnixNano())
			// Store the float32 bits in the data array
			bits := math.Float32bits(r[i])
			e.EsrOut[cnt].Data[0] = byte(bits)
			e.EsrOut[cnt].Data[1] = byte(bits >> 8)
			e.EsrOut[cnt].Data[2] = byte(bits >> 16)
			e.EsrOut[cnt].Data[3] = byte(bits >> 24)
			e.lastState[i] = r[i]
			cnt++
			e.EsrFlag = true
		}
	}
}

// EsrMonX8 monitors up to 8 status inputs (bytes) and reports changes.
type EsrMonX8 struct {
	EsrFlag bool
	EsrOut  [4]EsrData

	// internal state
	lastState [8]byte
}

// Update executes the monitoring logic.
func (e *EsrMonX8) Update(dtIn time.Time, s [8]byte, a [8]string, mode byte) {
	e.EsrFlag = false
	// Clear previous output
	e.EsrOut = [4]EsrData{}
	cnt := 0

	for i := 0; i < 8 && cnt < 4; i++ {
		if s[i] != e.lastState[i] {
			// Check mode: 1=error only, 2=error+status, 3=error+status+debug
			if (s[i] < 100) || (s[i] >= 100 && s[i] < 200 && mode >= 2) || (s[i] >= 200 && mode == 3) {
				e.EsrOut[cnt] = StatusToEsr(s[i], a[i], dtIn, time.Duration(time.Now().UnixNano()))
				e.lastState[i] = s[i]
				cnt++
				e.EsrFlag = true
			}
		}
	}
}

// OscatVersion returns the library version number or release date.
func OscatVersion(in bool) uint32 {
	if in {
		// Corresponds to DATE_TO_DWORD(D#2024-07-16)
		// This is a placeholder. A real implementation would calculate this.
		return 19736 // Days since 1970-01-01 for 2024-07-16
	}
	return 335
}

// StatusToEsr creates ESR data from a status byte.
func StatusToEsr(status byte, address string, dtIn time.Time, ts time.Duration) EsrData {
	var esr EsrData
	if status < 100 {
		esr.Typ = 1
	} else if status < 200 {
		esr.Typ = 2
	} else {
		esr.Typ = 3
	}
	esr.Adress = address
	esr.Ds = dtIn
	esr.Ts = ts
	esr.Data[0] = status
	return esr
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
