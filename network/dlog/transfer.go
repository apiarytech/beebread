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
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/inet"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// transfer is what DLOG_FILE_TO_FTP and DLOG_FILE_TO_SMTP share: the queue
// of the files a store closed, and the retries of their transfer.
type transfer struct {
	X           *network.DLOG_DATA
	FILE_DELETE iec.BOOL
	TIMEOUT     iec.TIME // default T#30s
	RETRY       iec.INT
	RETRY_TIME  iec.TIME // default T#30s
	DNS_IP4     iec.DWORD
	DONE        iec.BOOL
	BUSY        iec.BOOL
	ERROR_C     iec.DWORD
	ERROR_T     iec.BYTE

	ucbd     network.UNI_CIRCULAR_BUFFER_DATA
	ucb      UNI_CIRCULAR_BUFFER
	wt       timers.TON
	ftrigOld iec.BOOL
	file     iec.STRING
	step     iec.INT
	cBusy    iec.BOOL
	cDone    iec.BOOL
	cnt      iec.INT
}

// newTransfer returns a transfer of the data x, with the initial values.
func newTransfer(x *network.DLOG_DATA) transfer {
	return transfer{X: x, TIMEOUT: iec.TIME(30 * time.Second), RETRY_TIME: iec.TIME(30 * time.Second)}
}

// queue runs the queue, before the client: it returns the next file to
// send, if any, which a new job starts with.
func (t *transfer) queue(now time.Time) (iec.STRING, bool) {
	x := t.X
	t.ucb.DATA = &t.ucbd
	if x.NEW_FILE_RTRIG && !t.ftrigOld { // a file closed
		t.ucbd.D_HEAD, t.ucbd.D_MODE, t.ucbd.D_STRING = 1, ucbAdd, x.NEW_FILE
		t.ucb.Execute(now)
	}
	t.ftrigOld = x.NEW_FILE_RTRIG

	next, ok := iec.STRING(""), false
	switch t.step {
	case 0:
		if t.ucbd.BUF_COUNT > 0 && !t.cBusy {
			t.ucbd.D_MODE = ucbReadRemove
			t.ucb.Execute(now)
			next, ok = t.ucbd.D_STRING, true
			t.cnt = t.RETRY
			t.step = 10
		}
	case 10:
		if !t.cBusy {
			t.step = SEL[iec.INT](t.cDone, 20, 0) // a pause after an error
		}
	case 20:
		if t.wt.Q {
			t.cnt--
			// RETRY 0 retries for ever.
			t.step = SEL[iec.INT](t.RETRY == 0 || t.cnt >= 0, 0, 10)
		}
	}
	t.BUSY = t.step == 10
	t.DONE = t.step == 0
	t.wt.IN, t.wt.PT = t.step == 20, t.RETRY_TIME
	t.wt.Execute(now)
	return next, ok
}

// result takes the outputs of the client, and sets the errors.
func (t *transfer) result(done, busy iec.BOOL, errC iec.DWORD, errT iec.BYTE) {
	t.cDone, t.cBusy = done, busy
	switch {
	case errT != 0:
		t.ERROR_C, t.ERROR_T = errC, errT
	case t.ucbd.BUF_DATA_LOST > 0:
		t.ERROR_C, t.ERROR_T = 2, 6 // the queue overflowed
	case t.ucbd.BUF_USED_MAX > 90:
		t.ERROR_C, t.ERROR_T = 1, 6 // the queue is over 90% full
	default:
		t.ERROR_C, t.ERROR_T = 0, 0
	}
}

// DLOG_FILE_TO_FTP uploads each file a store of X closes to the FTP server
// FTP_URL, active if FTP_ACTIV, and deletes it then if FILE_DELETE. After
// a failure it tries again after RETRY_TIME, RETRY times or, if 0, for
// ever. BUSY while it sends, DONE when it waits. ERROR_C and ERROR_T: those
// of FTP_CLIENT, or 6 the queue (1 over 90% full, 2 overflow).
type DLOG_FILE_TO_FTP struct {
	transfer
	FTP_URL   iec.STRING // STRING(STRING_LENGTH)
	FTP_ACTIV iec.BOOL
	PLC_IP4   iec.DWORD

	ftp inet.FTP_CLIENT
}

// INIT resets the block and sets its inputs to their initial values.
func (d *DLOG_FILE_TO_FTP) INIT() { *d = DLOG_FILE_TO_FTP{transfer: newTransfer(d.X)} }

// Execute runs the block once.
func (d *DLOG_FILE_TO_FTP) Execute(now time.Time) {
	if d.X == nil {
		return
	}
	if f, ok := d.queue(now); ok {
		d.file = f
	}
	c := &d.ftp
	c.ACTIVATE, c.FILENAME, c.FTP_URL, c.FTP_DOWNLOAD = d.BUSY, d.file, d.FTP_URL, false
	c.FTP_ACTIV, c.FILE_DELETE, c.TIMEOUT, c.DNS_IP4, c.PLC_IP4 = d.FTP_ACTIV, d.FILE_DELETE, d.TIMEOUT, d.DNS_IP4, d.PLC_IP4
	c.Execute(now)
	d.result(c.DONE, c.BUSY, c.ERROR_C, c.ERROR_T)
}

// DLOG_FILE_TO_SMTP mails each file a store of X closes, as an attachment
// from MAILFROM to MAILTO through the SMTP server SERVER, with SUBJECT and
// BODY and the date DTI in the time zone DTI_OFFSET, and deletes it then
// if FILE_DELETE; see DLOG_FILE_TO_FTP.
type DLOG_FILE_TO_SMTP struct {
	transfer
	SERVER     *iec.STRING
	MAILFROM   *iec.STRING
	MAILTO     *iec.STRING // STRING(STRING_LENGTH)
	SUBJECT    *iec.STRING
	BODY       *iec.STRING // STRING(STRING_LENGTH)
	DTI        iec.DT
	DTI_OFFSET iec.INT

	smtp inet.SMTP_CLIENT
}

// INIT resets the block and sets its inputs to their initial values.
func (d *DLOG_FILE_TO_SMTP) INIT() {
	*d = DLOG_FILE_TO_SMTP{transfer: newTransfer(d.X),
		SERVER: d.SERVER, MAILFROM: d.MAILFROM, MAILTO: d.MAILTO, SUBJECT: d.SUBJECT, BODY: d.BODY}
}

// Execute runs the block once.
func (d *DLOG_FILE_TO_SMTP) Execute(now time.Time) {
	if d.X == nil || d.SERVER == nil || d.MAILFROM == nil || d.MAILTO == nil || d.SUBJECT == nil || d.BODY == nil {
		return
	}
	if f, ok := d.queue(now); ok {
		d.file = f
		if d.FILE_DELETE {
			d.file = CONCAT(d.file, ";#DEL#")
		}
	}
	c := &d.smtp
	c.ACTIVATE, c.TIMEOUT, c.DTI, c.DTI_OFFSET, c.DNS_IP4 = d.BUSY, d.TIMEOUT, d.DTI, d.DTI_OFFSET, d.DNS_IP4
	c.SERVER, c.MAILFROM, c.MAILTO, c.SUBJECT, c.BODY, c.FILES = d.SERVER, d.MAILFROM, d.MAILTO, d.SUBJECT, d.BODY, &d.file
	c.Execute(now)
	d.result(c.DONE, c.BUSY, c.ERROR_C, c.ERROR_T)
}
