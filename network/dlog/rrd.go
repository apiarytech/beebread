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
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/inet"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/iec"
)

// DLOG_STORE_RRD sends the data log X to the web server of URL, for a
// round robin database: an HTTP GET of URL with the values, separated by
// SEP, added to its query, at each rising edge of TRIG_M, every TRIG_T (1 s
// at least) and when a value changes by its DELTA, while ENABLE. The server
// answers 0 for OK. ERROR_C and ERROR_T: 1 DNS, 2 HTTP, 3 not OK, 4 the
// values do not fit in the URL, 6 the ring buffer (1 over 90% full, 2
// overflow). After an error it waits for ENABLE to fall.
type DLOG_STORE_RRD struct {
	X       *network.DLOG_DATA
	ENABLE  iec.BOOL
	TRIG_M  iec.BOOL
	TRIG_T  iec.TIME
	URL     iec.STRING
	DTI     iec.DT
	SEP     iec.BYTE
	DNS_IP4 iec.DWORD
	TIMEOUT iec.TIME
	ERROR_C iec.DWORD
	ERROR_T iec.BYTE

	ucb         UNI_CIRCULAR_BUFFER
	ipC         network.IP_C
	sBuf        network.NETWORK_BUFFER
	rBuf        network.NETWORK_BUFFER
	urlData     network.URL
	dns         inet.DNS_CLIENT
	http        inet.HTTP_GET
	ipc         ip.IP_CONTROL
	trigAuto    iec.BOOL
	trigStored  iec.BOOL
	trigMLast   iec.BOOL
	dtiLast     iec.DWORD
	step1       iec.INT
	step2       iec.INT
	sepChar     iec.STRING
	urlQuery    iec.STRING
	initialized bool
}

// INIT resets the block.
func (d *DLOG_STORE_RRD) INIT() {
	*d = DLOG_STORE_RRD{X: d.X}
	d.defaults()
}

// defaults sets the initial values of the internal state.
func (d *DLOG_STORE_RRD) defaults() {
	d.initialized = true
	d.dtiLast = dtNever
	d.http.INIT()
	d.dns.INIT()
	d.ipc.INIT()
}

// ring runs the ring buffer with the mode m.
func (d *DLOG_STORE_RRD) ring(now time.Time, m iec.INT) {
	d.X.UCB.D_MODE = m
	d.ucb.DATA = &d.X.UCB
	d.ucb.Execute(now)
}

// Execute runs the block once.
func (d *DLOG_STORE_RRD) Execute(now time.Time) {
	if !d.initialized {
		d.defaults()
	}
	if d.X == nil {
		return
	}
	x := d.X

	// The trigger, a second at a time.
	x.DTI = d.DTI
	x.CLOCK_TRIG = DT_TO_DWORD(d.DTI) != d.dtiLast
	d.dtiLast = DT_TO_DWORD(d.DTI)
	d.trigAuto = false
	if x.CLOCK_TRIG && d.TRIG_T >= iec.TIME(time.Second) {
		d.trigAuto = iec.UDINT(DT_TO_DWORD(d.DTI))%(iec.UDINT(TIME_TO_DWORD(d.TRIG_T))/1000) == 0
	}

	switch d.step1 {
	case 0:
		if d.ENABLE {
			if x.ID_MAX == 0 { // count the value blocks, once
				x.ADD_COM = addInfo
				x.STORE_TYPE = 1 // RRD
			}
			d.ERROR_T, d.ERROR_C = 0, 0
			d.ring(now, ucbReset)
			d.urlData = encoding.STRING_TO_URL(d.URL, "", "")
			d.sepChar = str.CHR_TO_STRING(d.SEP)
			d.urlQuery = d.urlData.QUERY
			d.step1 = 30
		}
	case 30: // the trigger
		trig := d.TRIG_M && !d.trigMLast || d.trigAuto || x.ADD_DATA_REQ
		switch {
		case bool(!d.ENABLE):
			x.ADD_COM = 0
			d.step1 = 0
		case bool(trig || d.trigStored):
			d.trigStored = false
			x.ADD_COM = addData
		default:
			x.ADD_COM = addDataReq
		}
	}
	x.ADD_DATA_REQ = false

	switch d.step2 {
	case 0:
		if x.UCB.BUF_DATA_LOST > 0 {
			d.ERROR_C, d.ERROR_T = 2, 6 // the ring overflowed
			d.step2 = 100
		} else if x.UCB.BUF_USED_MAX > 90 {
			d.ERROR_C, d.ERROR_T = 1, 6 // the ring is over 90% full
		}
		if iec.INT(x.UCB.BUF_COUNT) >= iec.INT(x.ID_MAX) { // a row: the query
			d.urlData.QUERY = d.urlQuery
			for n := iec.INT(1); n <= iec.INT(x.ID_MAX); n++ {
				d.ring(now, ucbReadRemove)
				if n < iec.INT(x.ID_MAX) {
					x.UCB.D_STRING = CONCAT(x.UCB.D_STRING, d.sepChar)
				}
				// OSCAT goes on after a value that does not fit, and a
				// later one that fits sends the query.
				if LEN(x.UCB.D_STRING)+LEN(d.urlData.QUERY) <= STRING_LENGTH {
					d.urlData.QUERY = CONCAT(d.urlData.QUERY, x.UCB.D_STRING)
					d.step2 = 40
				} else {
					d.ERROR_C, d.ERROR_T = 1, 4
					d.step2 = 100
				}
			}
		}
	case 40:
		if d.dns.DONE {
			d.step2 = 60
		} else if d.dns.ERROR > 0 {
			d.ERROR_C, d.ERROR_T = d.dns.ERROR, 1
			d.step2 = 100
		}
	case 60:
		if d.http.DONE {
			if d.rBuf.BUFFER[d.http.BODY_START] != '0' {
				d.ERROR_C, d.ERROR_T = 1, 3
			}
			d.step2 = 100
		} else if d.http.ERROR > 0 {
			d.ERROR_C, d.ERROR_T = d.http.ERROR, 2
			d.step2 = 100
		}
	case 100:
		if !d.http.DONE { // the buffer is free again
			d.step2 = SEL[iec.INT](d.ERROR_T == 0, 120, 0)
		}
	case 120:
		if !d.ENABLE {
			d.step2 = 0
		}
	}

	h := &d.http
	h.IP_C, h.S_BUF, h.R_BUF = &d.ipC, &d.sBuf, &d.rBuf
	h.IP4, h.GET, h.MODE, h.UNLOCK_BUF, h.URL_DATA = d.dns.IP4, d.step2 == 60, 2, d.step2 == 100, &d.urlData
	h.Execute(now)
	n := &d.dns
	n.IP_C, n.S_BUF, n.R_BUF = &d.ipC, &d.sBuf, &d.rBuf
	n.DOMAIN, n.IP4_DNS, n.ACTIVATE = d.urlData.DOMAIN, d.DNS_IP4, d.step2 == 40
	n.Execute(now)
	c := &d.ipc
	c.IP, c.PORT, c.TIME_OUT, c.IP_C, c.S_BUF, c.R_BUF = 0, 0, d.TIMEOUT, &d.ipC, &d.sBuf, &d.rBuf
	c.Execute(now)

	d.trigMLast = d.TRIG_M
}
