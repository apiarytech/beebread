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

package other

import (
	"time"

	. "beebread/basic"
	"beebread/basic/logic"
	"beebread/basic/math"
)

// ESR_COLLECT collects ESR data from up to 8 ESR_MON modules and stores them in an output array.
type ESR_COLLECT struct {
	EsrOut [32]ESR_DATA
	pos    int
	cnt    int
}

// Update executes the ESR collection logic.
// The variadic esrIn parameter allows passing multiple slices of EsrData.
func (e *ESR_COLLECT) Update(rst bool, esrIn ...[]ESR_DATA) int {
	if rst {
		e.pos = -1
		e.cnt = 0
	} else if e.cnt >= 0 {
		for _, esrArray := range esrIn {
			for _, esrItem := range esrArray {
				if esrItem.Typ > 0 {
					e.pos = math.INC1(e.pos, 32)
					e.EsrOut[e.pos] = esrItem
				}
			}
		}
	}
	return e.pos
}

// ESR_MON_R4 monitors up to 4 real inputs and reports changes with a timestamp and address label.
type ESR_MON_R4 struct {
	EsrFlag bool
	EsrOut  [4]ESR_DATA

	// internal state
	lastState [4]float32
}

// Update executes the monitoring logic.
func (e *ESR_MON_R4) Update(dtIn time.Time, r [4]float32, a [4]string, s [4]float32) {
	e.EsrFlag = false
	// Clear previous output
	e.EsrOut = [4]ESR_DATA{}
	cnt := 0

	for i := 0; i < 4 && cnt < 4; i++ {
		if math.DIFFER(float64(r[i]), float64(e.lastState[i]), float64(s[i])) {
			e.EsrOut[cnt].Typ = 20
			e.EsrOut[cnt].Adress = a[i]
			e.EsrOut[cnt].Ds = dtIn
			e.EsrOut[cnt].Ts = time.Duration(time.Now().UnixNano())
			// Store the float32 bits in the data array
			bits := logic.REAL_TO_DW(r[i])
			e.EsrOut[cnt].Data[0] = logic.BYTE_OF_DWORD(bits, 0)
			e.EsrOut[cnt].Data[1] = logic.BYTE_OF_DWORD(bits, 1)
			e.EsrOut[cnt].Data[2] = logic.BYTE_OF_DWORD(bits, 2)
			e.EsrOut[cnt].Data[3] = logic.BYTE_OF_DWORD(bits, 3)
			e.lastState[i] = r[i]
			cnt++
			e.EsrFlag = true
		}
	}
}

// ESR_MON_X8 monitors up to 8 status inputs (bytes) and reports changes.
type ESR_MON_X8 struct {
	EsrFlag bool
	EsrOut  [4]ESR_DATA

	// internal state
	lastState [8]byte
}

// Update executes the monitoring logic.
func (e *ESR_MON_X8) Update(dtIn time.Time, s [8]byte, a [8]string, mode byte) {
	e.EsrFlag = false
	// Clear previous output
	e.EsrOut = [4]ESR_DATA{}
	cnt := 0

	for i := 0; i < 8 && cnt < 4; i++ {
		if s[i] != e.lastState[i] {
			// Check mode: 1=error only, 2=error+status, 3=error+status+debug
			if (s[i] < 100) || (s[i] >= 100 && s[i] < 200 && mode >= 2) || (s[i] >= 200 && mode == 3) {
				e.EsrOut[cnt] = STATUS_TO_ESR(s[i], a[i], dtIn, time.Duration(time.Now().UnixNano()))
				e.lastState[i] = s[i]
				cnt++
				e.EsrFlag = true
			}
		}
	}
}

// STATUS_TO_ESR creates ESR data from a status byte.
func STATUS_TO_ESR(status byte, address string, dtIn time.Time, ts time.Duration) ESR_DATA {
	var esr ESR_DATA
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

// RDM calculates a pseudo-random number between 0.0 and 1.0.
// To use Rdm more than once per cycle, it needs to be called with different seed values for `last`.
func RDM(last float64) float64 {
	tn := uint32(logic.T_PLC_US())
	tc := logic.BIT_COUNT(tn)

	// Scramble bits based on original ST logic
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 2), 31)
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 5), 30)
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 4), 29)
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 1), 28)
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 0), 27)
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 7), 26)
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 6), 25)
	tn = logic.BIT_LOAD_DW(tn, logic.BIT_OF_DWORD(tn, 3), 24)

	tn = (tn << uint(tc)) | (tn >> (32 - uint(tc))) // ROL
	tn |= 0x80000001
	tn = tn%71474513 + uint32(tc+77)

	return math.FRACT(float64(tn) / 10000000.0 * (Math.E - math.LIMIT(0.0, last, 1.0)))
}

// RDM2 calculates an integer pseudo-random number in a given range.
func RDM2(last, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return int(RDM(math.FRACT(float64(last)*Math.Pi))*(float64(high-low+1))) + low
}

// RDMDW calculates a DWORD pseudo-random number.
func RDMDW(last uint32) uint32 {
	m := float64(logic.BIT_COUNT(last))
	rx1 := RDM(math.FRACT(m * Math.Pi))
	rdm1 := uint32(rx1 * 65535)

	rx2 := RDM(math.FRACT(m * Math.E))
	rdm2 := uint32(rx2 * 65535)

	return (rdm1 << 16) | (rdm2 & 0x0000FFFF)
}
