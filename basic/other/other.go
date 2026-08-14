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

// EsrCollect collects ESR data from up to 8 ESR_MON modules and stores them in an output array.
type EsrCollect struct {
	EsrOut [32]EsrData
	pos    int
	cnt    int
}

// Update executes the ESR collection logic.
// The variadic esrIn parameter allows passing multiple slices of EsrData.
func (e *EsrCollect) Update(rst bool, esrIn ...[]EsrData) int {
	if rst {
		e.pos = -1
		e.cnt = 0
	} else if e.cnt >= 0 {
		for _, esrArray := range esrIn {
			for _, esrItem := range esrArray {
				if esrItem.Typ > 0 {
					e.pos = math.Inc1(e.pos, 32)
					e.EsrOut[e.pos] = esrItem
				}
			}
		}
	}
	return e.pos
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
		if math.Differ(float64(r[i]), float64(e.lastState[i]), float64(s[i])) {
			e.EsrOut[cnt].Typ = 20
			e.EsrOut[cnt].Adress = a[i]
			e.EsrOut[cnt].Ds = dtIn
			e.EsrOut[cnt].Ts = time.Duration(time.Now().UnixNano())
			// Store the float32 bits in the data array
			bits := logic.RealToDw(r[i])
			e.EsrOut[cnt].Data[0] = logic.ByteOfDword(bits, 0)
			e.EsrOut[cnt].Data[1] = logic.ByteOfDword(bits, 1)
			e.EsrOut[cnt].Data[2] = logic.ByteOfDword(bits, 2)
			e.EsrOut[cnt].Data[3] = logic.ByteOfDword(bits, 3)
			e.lastState[i] = r[i]
			cnt++
			e.EsrFlag = true
		}
	}
}

// OscatVersion returns the version of the OSCAT library.
func OscatVersion() string {
	return "BASIC 3.35"
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
