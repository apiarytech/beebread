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

package dlog

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/inet"
	"github.com/apiarytech/royaljelly/iec"
)

// The MySQL commands DLOG_STORE_MYSQL sends.
const (
	comInitDB iec.BYTE = 2 // COM_INIT_DB
	comQuery  iec.BYTE = 3 // COM_QUERY
)

// DLOG_STORE_MYSQL writes the data log X to the table TB_NAME of the
// database DB_NAME of the MySQL server URL, creating both as needed, with
// a column of the time and one per value block: a row at each rising edge
// of TRIG_M, every TRIG_T (1 s at least) and when a value changes by its
// DELTA, while ENABLE. Every AUTO_DELETE seconds, if not 0, it deletes
// the rows older than that. SQL_INFO is the state of the connection.
// ERROR_C and ERROR_T: those of MYSQL_CONTROL, or 6 the ring buffer (1 over
// 90% full, 2 overflow, 3 a command too long for the send buffer).
type DLOG_STORE_MYSQL struct {
	SQL_INFO    *network.MYSQL_INFO
	X           *network.DLOG_DATA
	ENABLE      iec.BOOL
	TRIG_M      iec.BOOL
	TRIG_T      iec.TIME
	DTI         iec.DT
	URL         iec.STRING
	DB_NAME     iec.STRING // STRING(64)
	TB_NAME     iec.STRING // STRING(64)
	AUTO_DELETE iec.UDINT
	TIMEOUT     iec.TIME
	DNS_IP4     iec.DWORD
	ERROR_C     iec.DWORD
	ERROR_T     iec.BYTE

	ucb         UNI_CIRCULAR_BUFFER
	mc          inet.MYSQL_CONTROL
	y           network.MYSQL_COM
	trigAuto    iec.BOOL
	trigMLast   iec.BOOL
	dtiLast     iec.DWORD
	enableLast  iec.BOOL
	step        iec.INT
	sndText     iec.STRING
	sndSize     iec.INT
	sqlCommand  iec.BYTE
	nextStep    iec.INT
	delDT       iec.UDINT
	delStart    iec.BOOL
	initialized bool
}

// INIT resets the block.
func (d *DLOG_STORE_MYSQL) INIT() {
	*d = DLOG_STORE_MYSQL{SQL_INFO: d.SQL_INFO, X: d.X}
	d.defaults()
}

// defaults sets the initial values of the internal state.
func (d *DLOG_STORE_MYSQL) defaults() {
	d.initialized = true
	d.dtiLast = dtNever
}

// ring runs the ring buffer with the mode m.
func (d *DLOG_STORE_MYSQL) ring(now time.Time, m iec.INT) {
	d.X.UCB.D_MODE = m
	d.ucb.DATA = &d.X.UCB
	d.ucb.Execute(now)
}

// query builds a query in the send buffer from the head head, a part per
// value block from the ring, read by part, and ')'.
func (d *DLOG_STORE_MYSQL) query(now time.Time, head iec.STRING, part func() iec.STRING) {
	idx := iec.INT(5)
	put := func(s iec.STRING) {
		buffer.STRING_TO_BUFFER_(s, idx, d.y.S_BUF.BUFFER[:], iec.UINT(len(d.y.S_BUF.BUFFER)))
		idx += LEN(s)
	}
	put(head)
	for range int(d.X.ID_MAX) {
		d.ring(now, ucbReadRemove)
		put(part())
	}
	put(")")
	d.sqlCommand = comQuery
	d.sndSize = idx
	d.nextStep = 40
}

// Execute runs the block once.
func (d *DLOG_STORE_MYSQL) Execute(now time.Time) {
	if !d.initialized {
		d.defaults()
	}
	if d.X == nil || d.SQL_INFO == nil {
		return
	}
	x, y := d.X, &d.y

	// The trigger, a second at a time.
	x.DTI = d.DTI
	x.CLOCK_TRIG = DT_TO_DWORD(d.DTI) != d.dtiLast
	d.dtiLast = DT_TO_DWORD(d.DTI)
	d.trigAuto = false
	if x.CLOCK_TRIG {
		d.delStart = iec.UDINT(DT_TO_DWORD(d.DTI))-d.delDT > d.AUTO_DELETE && d.AUTO_DELETE > 0
		if d.TRIG_T >= iec.TIME(time.Second) {
			d.trigAuto = iec.UDINT(DT_TO_DWORD(d.DTI))%(iec.UDINT(TIME_TO_DWORD(d.TRIG_T))/1000) == 0
		}
	}
	x.ADD_COM = 0

	switch d.step {
	case 0:
		if d.ENABLE && !d.enableLast {
			y.SQL_URL, y.TIMEOUT, y.DNS_IP4, y.SQL_CON = d.URL, d.TIMEOUT, d.DNS_IP4, true
			d.ring(now, ucbReset)
			if x.ID_MAX == 0 { // count the value blocks, once
				x.ADD_COM = addInfo
				x.STORE_TYPE = 5 // MySQL
			}
			d.step = 5
		}
	case 5:
		x.ADD_COM = addHeader
		d.step = 10
	case 10:
		if d.SQL_INFO.SQL_CONNECTED {
			d.sqlCommand = comQuery
			d.sndText = CONCAT("CREATE DATABASE IF NOT EXISTS `", d.DB_NAME, "`")
			d.nextStep = 20
		}
	case 20:
		d.sqlCommand = comInitDB
		d.sndText = d.DB_NAME
		d.nextStep = 30
	case 30: // the table, of the column names and types
		d.query(now, CONCAT("CREATE TABLE IF NOT EXISTS `", d.TB_NAME,
			"` (`ID` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY ,",
			" store_timestamp TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP"), func() iec.STRING {
			kind := map[iec.WORD]iec.STRING{2: "` FLOAT", 6: "` BOOL", 8: "` INT", 13: "` DATETIME"}[SHR(x.UCB.D_HEAD, 8)]
			if kind == "" {
				kind = "` VARCHAR(80)"
			}
			return CONCAT(", `", x.UCB.D_STRING, kind)
		})
	case 40:
		if d.delStart {
			d.delStart = false
			d.sqlCommand = comQuery
			d.sndText = CONCAT("DELETE FROM `", d.TB_NAME,
				"` WHERE TIMESTAMPDIFF(SECOND, store_timestamp, NOW()) > ", network.UDINT_TO_STRING(d.AUTO_DELETE))
			d.nextStep = 40
			d.delDT = iec.UDINT(DT_TO_DWORD(d.DTI))
		} else if iec.INT(x.UCB.BUF_COUNT) >= iec.INT(x.ID_MAX) { // a row
			d.query(now, CONCAT("INSERT INTO `", d.TB_NAME, "` VALUES (NULL, CURRENT_TIMESTAMP"), func() iec.STRING {
				return CONCAT(`, "`, x.UCB.D_STRING, `"`)
			})
		}
		if d.ENABLE {
			if d.TRIG_M && !d.trigMLast || d.trigAuto || x.ADD_DATA_REQ {
				x.ADD_COM = addData
			} else {
				x.ADD_COM = addDataReq
			}
		} else {
			y.SQL_CON = false
			d.step = 0
		}
		x.ADD_DATA_REQ = false
	case 199: // the answer
		switch y.SQL_RCV_STATE {
		case 1: // OK
			d.step = d.nextStep
		case 2:
			d.step = 0
		}
	}

	// Send the command.
	if d.sqlCommand > 0 {
		if d.sndSize > 0 {
			y.S_BUF.SIZE = iec.UINT(d.sndSize)
			d.sndSize = 0
		} else {
			buffer.STRING_TO_BUFFER_(d.sndText, 5, y.S_BUF.BUFFER[:], iec.UINT(len(y.S_BUF.BUFFER)))
			y.S_BUF.SIZE = iec.UINT(5 + LEN(d.sndText))
			d.sndText = ""
		}
		if int(y.S_BUF.SIZE) > len(y.S_BUF.BUFFER) { // too long: give up
			y.S_BUF.SIZE = 0
			d.sqlCommand = 0
			y.SQL_CON = false
			d.step = 999
		} else {
			for i := 5; i < int(y.S_BUF.SIZE); i++ { // " becomes '
				if y.S_BUF.BUFFER[i] == '"' {
					y.S_BUF.BUFFER[i] = '\''
				}
			}
			y.S_BUF.BUFFER[4] = d.sqlCommand
			y.SQL_PACKET_NO = 255
			d.sqlCommand = 0
			d.step = 199
		}
	}

	d.mc.COM, d.mc.INFO = y, d.SQL_INFO
	d.mc.Execute(now)

	switch {
	case d.step == 999:
		d.ERROR_C, d.ERROR_T = 3, 6 // the send buffer overflowed
	case y.ERROR_T > 0:
		d.ERROR_C, d.ERROR_T = y.ERROR_C, y.ERROR_T
		d.step = 0
	case x.UCB.BUF_DATA_LOST > 0:
		d.ERROR_C, d.ERROR_T = 2, 6
	case x.UCB.BUF_USED_MAX > 90:
		d.ERROR_C, d.ERROR_T = 1, 6
	default:
		d.ERROR_C, d.ERROR_T = 0, 0
	}
	d.trigMLast = d.TRIG_M
	d.enableLast = d.ENABLE
}
