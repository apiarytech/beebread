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

// Package dlog is the port of the OSCAT NETWORK data logger. A store block,
// DLOG_STORE_FILE_CSV, _HTML or _XML, DLOG_STORE_MYSQL or DLOG_STORE_RRD,
// and the value blocks DLOG_BOOL, DLOG_DINT, DLOG_DT, DLOG_REAL,
// DLOG_REAL_ARRAY and DLOG_STRING share a DLOG_DATA. The store runs first
// in each scan and sets X.ADD_COM: 1 asks the value blocks to count
// themselves, 2 for their column names, 3 for their values and 4 whether
// one changed by its DELTA. The value blocks put their answers in the ring
// buffer X.UCB, a UNI_CIRCULAR_BUFFER, where the store reads them.
// DLOG_FILE_TO_FTP and DLOG_FILE_TO_SMTP send the files the file stores
// finish, and DLOG_CRON_TAB triggers by a cron table.
package dlog

import (
	"encoding/binary"
	gomath "math"
	"time"

	. "github.com/apiarytech/beebread/basic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// The commands of a store to the value blocks, in X.ADD_COM.
const (
	addInfo    iec.INT = 1 // count the blocks
	addHeader  iec.INT = 2 // the column names
	addData    iec.INT = 3 // the values
	addDataReq iec.INT = 4 // ask for values if one changed
)

// The modes of UNI_CIRCULAR_BUFFER, in D_MODE.
const (
	ucbAdd        iec.INT = 1  // add a value
	ucbRead       iec.INT = 10 // read the oldest value
	ucbRemove     iec.INT = 11 // remove the value just read
	ucbReadRemove iec.INT = 12 // read and remove the oldest value
	ucbReset      iec.INT = 99 // empty the buffer
)

// headerSize is the size of the header of a value in the ring.
const headerSize = 4

// UNI_CIRCULAR_BUFFER is a ring buffer of values in DATA.BUF, each a header
// of 4 bytes, the type D_HEAD and the size, high byte first, and the
// value: D_STRING (D_HEAD 1), D_REAL (2) or D_DWORD (3), little endian. The
// mode D_MODE asks it to add a value (1), read the oldest (10), remove the
// one just read (11), read and remove the oldest (12) or empty itself
// (99); it is 0 when done. A value that does not fit is counted in
// BUF_DATA_LOST.
type UNI_CIRCULAR_BUFFER struct {
	DATA *network.UNI_CIRCULAR_BUFFER_DATA

	cnt      iec.UINT
	modeLast iec.INT
}

// INIT resets the block.
func (u *UNI_CIRCULAR_BUFFER) INIT() { *u = UNI_CIRCULAR_BUFFER{DATA: u.DATA} }

// Execute runs the block once.
func (u *UNI_CIRCULAR_BUFFER) Execute(now time.Time) {
	if u.DATA == nil {
		return
	}
	d := u.DATA
	total := iec.UINT(len(d.BUF))
	switch d.D_MODE {
	case ucbAdd:
		var data []iec.BYTE
		switch iec.BYTE(d.D_HEAD) { // the low byte
		case 1:
			data = CHARS(d.D_STRING)
		case 2:
			data = le(gomath.Float32bits(float32(d.D_REAL)))
		case 3:
			data = le(uint32(d.D_DWORD))
		}
		u.cnt = iec.UINT(len(data))
		d.BUF_DATA_CNT++
		if u.cnt+headerSize > total-d.BUF_SIZE {
			// no room
			d.BUF_DATA_LOST++
			return
		}
		d.D_SIZE = iec.INT(u.cnt)
		u.modeLast = d.D_MODE
		head := SHL(iec.DWORD(d.D_HEAD), 16) | iec.DWORD(u.cnt)
		for n := 0; n < headerSize; n++ {
			head = ROL(head, 8)
			d.BUF[d.LAST_] = iec.BYTE(head)
			d.LAST_ = (d.LAST_ + 1) % total
		}
		for _, b := range data {
			d.BUF[d.LAST_] = b
			d.LAST_ = (d.LAST_ + 1) % total
		}
		d.BUF_SIZE += u.cnt + headerSize
		d.BUF_COUNT++
		d.BUF_USED = iec.USINT(100 * iec.UDINT(d.BUF_SIZE) / iec.UDINT(total))
		d.D_MODE = 0

	case ucbRead, ucbRemove, ucbReadRemove:
		if d.D_MODE != ucbRemove {
			if d.BUF_COUNT == 0 {
				return
			}
			d.GETEND_ = d.FIRST_ // the oldest value
			u.modeLast = d.D_MODE
			var head iec.DWORD
			for n := 0; n < headerSize; n++ {
				head = ROL(head, 8) | iec.DWORD(d.BUF[d.GETEND_])
				d.GETEND_ = (d.GETEND_ + 1) % total
			}
			d.GETSTART_ = d.GETEND_
			u.cnt = iec.UINT(head & 0xFFFF)
			d.D_SIZE = iec.INT(u.cnt)
			d.D_HEAD = iec.WORD(SHR(head, 16))
			data := make([]iec.BYTE, u.cnt)
			for i := range data {
				data[i] = d.BUF[d.GETEND_]
				d.GETEND_ = (d.GETEND_ + 1) % total
			}
			switch iec.BYTE(d.D_HEAD) {
			case 1:
				d.D_STRING = STR(data)
			case 2:
				if len(data) >= 4 {
					d.D_REAL = iec.REAL(gomath.Float32frombits(binary.LittleEndian.Uint32(bytesOf(data))))
				}
			case 3:
				if len(data) >= 4 {
					d.D_DWORD = iec.DWORD(binary.LittleEndian.Uint32(bytesOf(data)))
				}
			}
			d.D_MODE = 0
		}
		// Remove only right after a read.
		if u.modeLast == ucbRead && d.D_MODE == ucbRemove || u.modeLast == ucbReadRemove && d.D_MODE == 0 {
			u.modeLast = d.D_MODE
			d.BUF_SIZE -= u.cnt + headerSize
			d.BUF_COUNT--
			d.FIRST_ = d.GETEND_
			d.BUF_USED = iec.USINT(100 * iec.UDINT(d.BUF_SIZE) / iec.UDINT(total))
			d.D_MODE = 0
		}

	case ucbReset:
		d.FIRST_, d.LAST_, d.GETSTART_, d.GETEND_ = 0, 0, 0, 0
		d.BUF_COUNT, d.BUF_SIZE, d.BUF_USED = 0, 0, 0
		d.D_STRING, d.D_DWORD, d.D_REAL, d.D_HEAD = "", 0, 0, 0
		d.D_MODE, d.D_SIZE = 0, 0
		u.modeLast = 0
	}
	if d.BUF_USED > d.BUF_USED_MAX {
		d.BUF_USED_MAX = d.BUF_USED
	}
}

// le returns the bytes of d, little endian.
func le(d uint32) []iec.BYTE {
	return []iec.BYTE{iec.BYTE(d), iec.BYTE(d >> 8), iec.BYTE(d >> 16), iec.BYTE(d >> 24)}
}

func bytesOf(b []iec.BYTE) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = byte(b[i])
	}
	return out
}

// value is what the value blocks share: their type and their buffer.
type value struct {
	X *network.DLOG_DATA

	ucb UNI_CIRCULAR_BUFFER
	id  iec.WORD
}

// add puts the string s, of the block's type, in the ring.
func (v *value) add(now time.Time, s iec.STRING) {
	v.X.UCB.D_STRING = network.STRING_N(s, STRING_LENGTH)
	v.X.UCB.D_HEAD = v.id
	v.X.UCB.D_MODE = ucbAdd
	v.ucb.DATA = &v.X.UCB
	v.ucb.Execute(now)
}

// info counts the block, of the type id: its source type in the high byte,
// STRING (1) its value's.
func (v *value) info(id iec.WORD) {
	v.X.ID_MAX++
	v.id = id
}

// DLOG_BOOL logs STATE, as ON or OFF ('1' or '0' for MySQL), in the column
// COLUMN, and asks for a log when it changes if DELTA.
type DLOG_BOOL struct {
	value
	STATE  iec.BOOL
	OFF    iec.STRING
	ON     iec.STRING
	COLUMN iec.STRING // STRING(40)
	DELTA  iec.BOOL

	lastDelta iec.BOOL
	stateLast iec.BOOL
}

// INIT resets the block.
func (d *DLOG_BOOL) INIT() { *d = DLOG_BOOL{value: value{X: d.X}} }

// Execute runs the block once.
func (d *DLOG_BOOL) Execute(now time.Time) {
	if d.X == nil {
		return
	}
	switch d.X.ADD_COM {
	case addInfo:
		d.info(0x0601) // BOOL
	case addHeader:
		d.add(now, d.COLUMN)
	case addData:
		if d.X.STORE_TYPE == 5 { // MySQL
			d.add(now, SEL[iec.STRING](d.stateLast, "0", "1"))
		} else {
			d.add(now, SEL(d.stateLast, d.OFF, d.ON))
		}
		d.lastDelta = d.STATE
	case addDataReq:
		if d.DELTA && d.STATE != d.lastDelta {
			d.X.ADD_DATA_REQ = true
			d.lastDelta = d.STATE
		}
	}
	d.stateLast = d.STATE
}

// DLOG_DINT logs VALUE in the column COLUMN, and asks for a log when it
// changes by DELTA if that is not 0.
type DLOG_DINT struct {
	value
	VALUE  iec.DINT
	COLUMN iec.STRING // STRING(40)
	DELTA  iec.DINT

	deltaLast iec.DINT
	valueLast iec.DINT
}

// INIT resets the block.
func (d *DLOG_DINT) INIT() { *d = DLOG_DINT{value: value{X: d.X}} }

// Execute runs the block once.
func (d *DLOG_DINT) Execute(now time.Time) {
	if d.X == nil {
		return
	}
	switch d.X.ADD_COM {
	case addInfo:
		d.info(0x0801) // DINT
	case addHeader:
		d.add(now, d.COLUMN)
	case addData:
		d.add(now, DINT_TO_STRING(d.valueLast))
		d.deltaLast = d.valueLast
	case addDataReq:
		if d.DELTA != 0 && (d.VALUE <= d.deltaLast-d.DELTA || d.VALUE >= d.deltaLast+d.DELTA) {
			d.X.ADD_DATA_REQ = true
			d.deltaLast = d.VALUE
		}
	}
	d.valueLast = d.VALUE
}

// DLOG_DT logs the time of the log, X.DTI, in the format FMT of DT_TO_STRF,
// in the column COLUMN, and asks for a log every DELTA seconds if that is
// not 0.
type DLOG_DT struct {
	value
	FMT    iec.STRING // default '#A-#D-#H #N:#R:#T'
	COLUMN iec.STRING // STRING(40)
	DELTA  iec.UDINT

	deltaLast iec.DT
}

// INIT resets the block and sets FMT to its initial value.
func (d *DLOG_DT) INIT() { *d = DLOG_DT{value: value{X: d.X}, FMT: "#A-#D-#H #N:#R:#T"} }

// Execute runs the block once.
func (d *DLOG_DT) Execute(now time.Time) {
	if d.X == nil {
		return
	}
	switch d.X.ADD_COM {
	case addInfo:
		d.info(0x0D01) // DT
	case addHeader:
		d.add(now, d.COLUMN)
	case addData:
		d.add(now, str.DT_TO_STRF(d.X.DTI, 0, d.FMT, 0))
		d.deltaLast = d.X.DTI
	case addDataReq:
		if d.DELTA != 0 && iec.UDINT(DT_TO_DWORD(d.X.DTI)) >= iec.UDINT(DT_TO_DWORD(d.deltaLast))+d.DELTA {
			d.X.ADD_DATA_REQ = true
			d.deltaLast = d.X.DTI
		}
	}
}

// DLOG_REAL logs VALUE with N decimals and the decimal point D in the column
// COLUMN, and asks for a log when it changes by DELTA if that is not 0.
type DLOG_REAL struct {
	value
	VALUE  iec.REAL
	N      iec.INT
	D      iec.STRING // STRING(1), default ','
	COLUMN iec.STRING // STRING(40)
	DELTA  iec.REAL

	deltaLast iec.REAL
	valueLast iec.REAL
}

// INIT resets the block and sets D to its initial value.
func (d *DLOG_REAL) INIT() { *d = DLOG_REAL{value: value{X: d.X}, D: ","} }

// Execute runs the block once.
func (d *DLOG_REAL) Execute(now time.Time) {
	if d.X == nil {
		return
	}
	switch d.X.ADD_COM {
	case addInfo:
		d.info(0x0201) // REAL
	case addHeader:
		d.add(now, d.COLUMN)
	case addData:
		d.add(now, str.REAL_TO_STRF(d.valueLast, d.N, d.D))
		d.deltaLast = d.valueLast
	case addDataReq:
		if d.DELTA != 0 && (d.VALUE <= d.deltaLast-d.DELTA || d.VALUE >= d.deltaLast+d.DELTA) {
			d.X.ADD_DATA_REQ = true
			d.deltaLast = d.VALUE
		}
	}
	d.valueLast = d.VALUE
}

// DLOG_STRING logs STR in the column COLUMN.
type DLOG_STRING struct {
	value
	STR    iec.STRING
	COLUMN iec.STRING // STRING(40)

	strLast iec.STRING
}

// INIT resets the block.
func (d *DLOG_STRING) INIT() { *d = DLOG_STRING{value: value{X: d.X}} }

// Execute runs the block once.
func (d *DLOG_STRING) Execute(now time.Time) {
	if d.X == nil {
		return
	}
	switch d.X.ADD_COM {
	case addInfo:
		d.info(0x0101) // STRING
	case addHeader:
		d.add(now, d.COLUMN)
	case addData:
		d.add(now, d.strLast)
	}
	d.strLast = d.STR
}

// DLOG_REAL_ARRAY logs the first R_COUNT values of R_ARRAY, each a column,
// with N decimals and the decimal point D; a value without a column name
// is 'Spalte_n'.
type DLOG_REAL_ARRAY struct {
	value
	R_ARRAY *network.DLOG_REAL_ARRAY_DATA
	N       iec.INT
	D       iec.STRING // STRING(1)
	R_COUNT iec.INT
}

// INIT resets the block.
func (d *DLOG_REAL_ARRAY) INIT() { *d = DLOG_REAL_ARRAY{value: value{X: d.X}, R_ARRAY: d.R_ARRAY} }

// Execute runs the block once.
func (d *DLOG_REAL_ARRAY) Execute(now time.Time) {
	if d.X == nil || d.R_ARRAY == nil {
		return
	}
	for i := 0; i < int(d.R_COUNT) && i < len(d.R_ARRAY); i++ {
		e := &d.R_ARRAY[i]
		switch d.X.ADD_COM {
		case addInfo:
			d.info(0x0201) // REAL
			if LEN(e.COLUMN) == 0 {
				e.COLUMN = CONCAT("Spalte_", network.USINT_TO_STRING(d.X.ID_MAX))
			}
		case addHeader:
			d.add(now, e.COLUMN)
		case addData:
			d.add(now, str.REAL_TO_STRF(e.VALUE_LAST_, d.N, d.D))
			e.DELTA_LAST_ = e.VALUE_LAST_
		case addDataReq:
			if e.DELTA != 0 && (e.VALUE <= e.DELTA_LAST_-e.DELTA || e.VALUE >= e.DELTA_LAST_+e.DELTA) {
				d.X.ADD_DATA_REQ = true
				e.DELTA_LAST_ = e.VALUE
			}
		}
		e.VALUE_LAST_ = e.VALUE
	}
}
