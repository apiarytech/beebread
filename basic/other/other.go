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

// Package other is the port of the OSCAT BASIC functions for event, status
// and error reports (ESR), and the library's version.
package other

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/logic"
	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/royaljelly/iec"
)

// ESR_COLLECT collects the reports of up to 8 ESR_MON blocks in ESR_OUT. POS
// is the position of the last report, or -1 if there is none; whoever reads
// the reports sets it back to -1. When ESR_OUT is full, it starts again at
// position 0.
//
// OSCAT 3.35 never leaves its reset state, because its counter starts at -1
// and only the collecting code changes it; the port starts collecting after
// the first run.
type ESR_COLLECT struct {
	ESR_0, ESR_1, ESR_2, ESR_3, ESR_4, ESR_5, ESR_6, ESR_7 [4]ESR_DATA
	RST                                                    iec.BOOL
	POS                                                    *iec.INT
	ESR_OUT                                                [32]ESR_DATA

	init iec.BOOL
}

// INIT resets the block.
func (e *ESR_COLLECT) INIT() { *e = ESR_COLLECT{POS: e.POS} }

// Execute runs the block once.
func (e *ESR_COLLECT) Execute(now time.Time) {
	if e.POS == nil {
		return
	}
	if e.RST || !e.init {
		e.init = true
		*e.POS = -1
		return
	}
	in := [8]*[4]ESR_DATA{&e.ESR_0, &e.ESR_1, &e.ESR_2, &e.ESR_3, &e.ESR_4, &e.ESR_5, &e.ESR_6, &e.ESR_7}
	for cnt := 0; cnt <= 3; cnt++ {
		for _, esr := range in {
			if esr[cnt].TYP > 0 {
				*e.POS = math.INC1(*e.POS, 32)
				e.ESR_OUT[*e.POS] = esr[cnt]
			}
		}
	}
}

// esrMon is the part the ESR_MON blocks share: ESR_OUT, cleared each run,
// takes up to 4 reports, and ESR_FLAG is true when there are any.
type esrMon struct {
	cnt int
}

func (m *esrMon) start(out *[4]ESR_DATA, flag *iec.BOOL) {
	*flag = false
	for i := range out {
		out[i].TYP = 0
	}
	m.cnt = 0
}

// ESR_MON_B8 reports the changes of up to 8 binary inputs S0..S7, labelled
// A0..A7, with the time DT_IN, in ESR_OUT: type 10 for a change to false
// and 11 for a change to true. It reports up to 4 changes in a run; the
// changes of S4..S7 beyond that wait for the next run.
type ESR_MON_B8 struct {
	S0, S1, S2, S3, S4, S5, S6, S7 iec.BOOL
	DT_IN                          iec.DT
	A0, A1, A2, A3, A4, A5, A6, A7 iec.STRING // STRING(10)
	ESR_FLAG                       iec.BOOL
	ESR_OUT                        *[4]ESR_DATA

	x   [8]iec.BOOL
	mon esrMon
}

// INIT resets the block.
func (e *ESR_MON_B8) INIT() { *e = ESR_MON_B8{ESR_OUT: e.ESR_OUT} }

// Execute runs the block once.
func (e *ESR_MON_B8) Execute(now time.Time) {
	if e.ESR_OUT == nil {
		return
	}
	tx := DWORD_TO_TIME(PLC_MS(now))
	e.mon.start(e.ESR_OUT, &e.ESR_FLAG)
	s := [8]iec.BOOL{e.S0, e.S1, e.S2, e.S3, e.S4, e.S5, e.S6, e.S7}
	a := [8]iec.STRING{e.A0, e.A1, e.A2, e.A3, e.A4, e.A5, e.A6, e.A7}
	for i := range s {
		// The first 4 inputs always fit; the others only while there is
		// room.
		if s[i] == e.x[i] || (i >= 4 && e.mon.cnt >= 4) {
			continue
		}
		out := &e.ESR_OUT[e.mon.cnt]
		out.TYP = 10 + iec.BYTE(BOOL_TO_INT(s[i]))
		out.ADRESS = a[i]
		out.DS = e.DT_IN
		out.TS = tx
		e.x[i] = s[i]
		e.mon.cnt++
		e.ESR_FLAG = true
	}
}

// ESR_MON_R4 reports the changes of up to 4 REAL inputs R0..R3, labelled
// A0..A3, of more than S0..S3, with the time DT_IN, in ESR_OUT: type 20,
// with the value's 4 bytes in DATA.
type ESR_MON_R4 struct {
	R0, R1, R2, R3 iec.REAL
	DT_IN          iec.DT
	A0, A1, A2, A3 iec.STRING // STRING(10)
	S0, S1, S2, S3 iec.REAL
	ESR_FLAG       iec.BOOL
	ESR_OUT        *[4]ESR_DATA

	x   [4]iec.REAL
	mon esrMon
}

// INIT resets the block.
func (e *ESR_MON_R4) INIT() { *e = ESR_MON_R4{ESR_OUT: e.ESR_OUT} }

// Execute runs the block once.
func (e *ESR_MON_R4) Execute(now time.Time) {
	if e.ESR_OUT == nil {
		return
	}
	tx := DWORD_TO_TIME(PLC_MS(now))
	e.mon.start(e.ESR_OUT, &e.ESR_FLAG)
	r := [4]iec.REAL{e.R0, e.R1, e.R2, e.R3}
	a := [4]iec.STRING{e.A0, e.A1, e.A2, e.A3}
	s := [4]iec.REAL{e.S0, e.S1, e.S2, e.S3}
	for i := range r {
		if !math.DIFFER(r[i], e.x[i], s[i]) {
			continue
		}
		out := &e.ESR_OUT[e.mon.cnt]
		out.TYP = 20
		out.ADRESS = a[i]
		out.DS = e.DT_IN
		out.TS = tx
		bits := logic.REAL_TO_DW(r[i])
		for b := iec.BYTE(0); b < 4; b++ {
			out.DATA[b] = logic.BYTE_OF_DWORD(bits, b)
		}
		e.x[i] = r[i]
		e.mon.cnt++
		e.ESR_FLAG = true
	}
}

// ESR_MON_X8 reports the changes of up to 8 status inputs S0..S7, labelled
// A0..A7, with the time DT_IN, in ESR_OUT; see STATUS_TO_ESR. MODE 1
// reports errors, below 100, 2 also status, below 200, and 3 also debug
// messages. It reports up to 4 changes in a run.
type ESR_MON_X8 struct {
	S0, S1, S2, S3, S4, S5, S6, S7 iec.BYTE
	DT_IN                          iec.DT
	MODE                           iec.BYTE // default 3
	A0, A1, A2, A3, A4, A5, A6, A7 iec.STRING
	ESR_FLAG                       iec.BOOL
	ESR_OUT                        *[4]ESR_DATA

	x   [8]iec.BYTE
	mon esrMon
}

// INIT resets the block and sets MODE to its initial value.
func (e *ESR_MON_X8) INIT() { *e = ESR_MON_X8{ESR_OUT: e.ESR_OUT, MODE: 3} }

// Execute runs the block once.
func (e *ESR_MON_X8) Execute(now time.Time) {
	if e.ESR_OUT == nil {
		return
	}
	tx := DWORD_TO_TIME(PLC_MS(now))
	e.mon.start(e.ESR_OUT, &e.ESR_FLAG)
	s := [8]iec.BYTE{e.S0, e.S1, e.S2, e.S3, e.S4, e.S5, e.S6, e.S7}
	a := [8]iec.STRING{e.A0, e.A1, e.A2, e.A3, e.A4, e.A5, e.A6, e.A7}
	for i, v := range s {
		if i >= 4 && e.mon.cnt >= 4 {
			continue
		}
		report := v < 100 || v > 99 && v < 200 && e.MODE >= 2 || v > 199 && e.MODE == 3
		if v == e.x[i] || !report {
			continue
		}
		e.ESR_OUT[e.mon.cnt] = STATUS_TO_ESR(v, a[i], e.DT_IN, tx)
		e.x[i] = v
		e.mon.cnt++
		e.ESR_FLAG = true
	}
}

// OSCAT_VERSION returns the library's version, 335 for 3.35, or if in is
// true its release date as DATE_TO_DWORD(D#2024-07-16).
func OSCAT_VERSION(in iec.BOOL) iec.DWORD {
	if in {
		return DATE_TO_DWORD(iec.DATE(time.Date(2024, 7, 16, 0, 0, 0, 0, time.UTC)))
	}
	return 335
}

// STATUS_TO_ESR creates a report of a status byte: type 1 for an error,
// below 100, 2 for a status, below 200, and 3 for a debug message.
func STATUS_TO_ESR(status iec.BYTE, adress iec.STRING, dtIn iec.DT, ts iec.TIME) ESR_DATA {
	var out ESR_DATA
	switch {
	case status < 100:
		out.TYP = 1
	case status < 200:
		out.TYP = 2
	default:
		out.TYP = 3
	}
	out.ADRESS = adress
	out.DS = dtIn
	out.TS = ts
	out.DATA[0] = status
	return out
}
