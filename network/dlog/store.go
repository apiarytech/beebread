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
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/file"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// The heads of the records a file store puts in the ring for itself.
const (
	headOpen   iec.WORD = 0xF201 // the file name: open the file and write
	headCreate iec.WORD = 0xF301 // the file name: create the file and write
	headStop   iec.WORD = 0xEE00 // the log stops
)

// The formats of the file stores.
const (
	formatCSV = iota
	formatHTML
	formatXML
)

// dtNever is the date and time a store has seen before its first, the
// largest DT, DT#2106-02-07-06:28:15 as CODESYS reads OSCAT's
// DT#2070-02-06-06:28:15.
const dtNever iec.DWORD = 0xFFFF_FFFF

// fileStore is DLOG_STORE_FILE_CSV, _HTML and _XML: the trigger, the
// records of the ring, and the file they go to.
type fileStore struct {
	X          *network.DLOG_DATA
	SAVE_DATA  *network.DLOG_SAVE
	ENABLE     iec.BOOL
	TRIG_M     iec.BOOL
	TRIG_T     iec.TIME
	FILENAME   iec.STRING
	DTI        iec.DT
	AUTO_CLOSE iec.TIME // default T#15s
	ERROR_C    iec.DWORD
	ERROR_T    iec.BYTE

	fs          file.FILE_SERVER
	fsd         network.FILE_SERVER_DATA
	pt          network.NETWORK_BUFFER
	ucb         UNI_CIRCULAR_BUFFER
	trigAuto    iec.BOOL
	trigStored  iec.BOOL
	trigTmp     iec.BOOL
	trigMLast   iec.BOOL
	dtiLast     iec.DWORD
	fnLast      iec.STRING
	fn          iec.STRING
	step1       iec.INT
	step2       iec.INT
	idx         iec.INT
	logState    iec.USINT
	logStop     iec.BOOL
	n           iec.INT
	wdTon       timers.TON
	wdTime      iec.TIME
	awTon       timers.TON
	awEnable    iec.BOOL
	awAktiv     iec.BOOL
	totalBytes  iec.UDINT
	initialized bool
}

// defaults sets the initial values of the internal state.
func (s *fileStore) defaults() {
	s.initialized = true
	s.dtiLast = dtNever
}

// put writes text into the file buffer at idx, and moves idx past it.
func (s *fileStore) put(text iec.STRING) {
	buffer.STRING_TO_BUFFER_(text, s.idx, s.pt.BUFFER[:], iec.UINT(len(s.pt.BUFFER)))
	s.idx += LEN(text)
}

// ring runs the ring buffer with the mode m.
func (s *fileStore) ring(now time.Time, m iec.INT) {
	s.X.UCB.D_MODE = m
	s.ucb.DATA = &s.X.UCB
	s.ucb.Execute(now)
}

// openFile starts writing the file of the record read, a create or an open.
func (s *fileStore) openFile() {
	x, fsd := s.X, &s.fsd
	fsd.MODE = iec.BYTE(SHR(SHL(x.UCB.D_HEAD, 4), 12)) // 3 create, 2 open
	fsd.OFFSET = 0
	fsd.FILENAME = x.UCB.D_STRING
	s.pt.SIZE = 0
	s.SAVE_DATA.FN_REM = fsd.FILENAME
	x.NEW_FILE = fsd.FILENAME
	x.NEW_FILE_SIZE = 0
	x.NEW_FILE_RTRIG = false
	s.step2 = 20
}

// run runs a scan of the store of the format format, with the separator
// sep of CSV and the parts html of HTML: caption, table, and the rows of
// the head, even and odd.
func (s *fileStore) run(now time.Time, format int, sep iec.BYTE, html [5]iec.STRING) {
	if !s.initialized {
		s.defaults()
	}
	if s.X == nil || s.SAVE_DATA == nil {
		return
	}
	x, save, fsd := s.X, s.SAVE_DATA, &s.fsd
	s.fs.FSD, s.fs.PT = fsd, &s.pt

	// The trigger, a second at a time.
	x.DTI = s.DTI
	x.CLOCK_TRIG = DT_TO_DWORD(s.DTI) != s.dtiLast
	s.dtiLast = DT_TO_DWORD(s.DTI)
	s.trigAuto = false
	if x.CLOCK_TRIG {
		s.fn = str.DT_TO_STRF(s.DTI, 0, s.FILENAME, 0) // the file name of the time
		if s.TRIG_T >= iec.TIME(time.Second) {
			s.trigAuto = iec.UDINT(DT_TO_DWORD(s.DTI))%(iec.UDINT(TIME_TO_DWORD(s.TRIG_T))/1000) == 0
		}
	}

	switch s.step1 {
	case 0:
		if s.ENABLE {
			s.awEnable = s.AUTO_CLOSE >= iec.TIME(5*time.Second)
			s.wdTime = SEL(x.LOAD_TIME_MAX > 0, iec.TIME(5*time.Millisecond), x.LOAD_TIME_MAX)
			s.fnLast = s.fn
			if x.ID_MAX == 0 { // count the value blocks, once
				x.ADD_COM = addInfo
				x.STORE_TYPE = []iec.BYTE{2, 4, 3}[format]
			}
			s.step1 = 10
		}
	case 10: // the file
		x.UCB.D_STRING = s.fnLast
		if s.fnLast == save.FN_REM {
			x.UCB.D_HEAD = headOpen
			x.ADD_COM = 0
		} else {
			x.UCB.D_HEAD = headCreate
			x.ADD_COM = addHeader // a new file: the column names first
		}
		s.ring(now, ucbAdd)
		s.step1 = 30
	case 30: // the trigger
		s.trigTmp = s.TRIG_M && !s.trigMLast || s.trigAuto || x.ADD_DATA_REQ
		if !s.ENABLE || s.fn != s.fnLast {
			// stop the log, for good or for a new file
			s.trigStored = s.trigTmp
			x.UCB.D_HEAD = headStop
			s.ring(now, ucbAdd)
			x.ADD_COM = 0
			s.step1 = 0
		} else if s.trigTmp || s.trigStored {
			s.trigStored = false
			x.ADD_COM = addData
			save.TRIG_CNT++
			save.TRIG_CNT_TOTAL++
		} else {
			x.ADD_COM = addDataReq
		}
	}
	x.ADD_DATA_REQ = false

	switch s.step2 {
	case 0:
		s.awTon.IN, s.awTon.PT = s.idx > 0 && s.awEnable, s.AUTO_CLOSE
		s.awTon.Execute(now)
		if s.awTon.Q {
			// write after AUTO_CLOSE
			s.awAktiv = true
			s.step2 = 10
			break
		}
		s.wdTon.IN, s.wdTon.PT = false, s.wdTime
		s.wdTon.Execute(now)
		s.wdTon.IN = true
		s.wdTon.Execute(now)
		for x.UCB.BUF_COUNT > 0 && s.step2 == 0 && !s.wdTon.Q {
			count, idx := x.UCB.BUF_COUNT, s.idx
			s.record(now, format, sep, html)
			s.wdTon.Execute(now)
			// OSCAT waits for the watchdog on a record that does not fit;
			// the time of a scan does not move here.
			if x.UCB.BUF_COUNT == count && s.idx == idx && s.step2 == 0 {
				break
			}
		}
	case 10: // write the buffer at the end of the file
		if s.idx > 0 {
			fsd.MODE = 2
			fsd.OFFSET = 0xFFFF_FFFE
			s.pt.SIZE = iec.UINT(s.idx)
			s.idx = 0
		}
		s.step2 = 20
	case 20:
		if fsd.MODE == 0 {
			s.totalBytes = fsd.FILE_SIZE
			x.NEW_FILE_SIZE = s.totalBytes + iec.UDINT(s.idx)
			if s.logStop || s.awAktiv {
				s.awAktiv = false
				fsd.MODE = 5 // close the file
			}
			s.step2 = 30
		}
	case 30:
		if fsd.MODE == 0 {
			if s.logStop {
				x.NEW_FILE_RTRIG = true
				s.logStop = false
				save.FN_REM = ""
			}
			s.step2 = 0
		}
	}

	s.fs.Execute(now)
	switch {
	case fsd.ERROR != 0:
		s.ERROR_C, s.ERROR_T = iec.DWORD(fsd.ERROR), 1
	case x.UCB.BUF_DATA_LOST > 0:
		s.ERROR_C, s.ERROR_T = 2, 6 // the ring overflowed
	case x.UCB.BUF_USED_MAX > 90:
		s.ERROR_C, s.ERROR_T = 1, 6 // the ring is over 90% full
	default:
		s.ERROR_C, s.ERROR_T = 0, 0
	}
	s.trigMLast = s.TRIG_M
}

// record reads the oldest record of the ring and writes it into the file
// buffer, if it fits.
func (s *fileStore) record(now time.Time, format int, sep iec.BYTE, html [5]iec.STRING) {
	x, save := s.X, s.SAVE_DATA
	size := iec.INT(len(s.pt.BUFFER))
	s.ring(now, ucbRead)
	switch {
	case x.UCB.D_HEAD > 0xF000: // a file to open or create
		if format == formatCSV {
			s.ring(now, ucbRemove)
			s.openFile()
			if s.fsd.MODE == 3 {
				save.TRIG_CNT = 0
			}
			return
		}
		head := iec.STRING(`<?xml version="1.0" encoding="UTF-8"?><table>`)
		if format == formatHTML {
			head = CONCAT("<html><body><table ", html[1], "><caption>", html[0], "</caption>")
		}
		if s.idx+LEN(head) < size {
			s.openFile()
			if s.fsd.MODE == 3 { // a new file: its head
				save.TRIG_CNT = 0
				if format == formatHTML {
					save.COLOR, save.HEAD = false, false
				}
				s.put(head)
				x.NEW_FILE_SIZE = s.totalBytes + iec.UDINT(s.idx)
			}
			s.ring(now, ucbRemove)
		}
	case x.UCB.D_HEAD == headStop:
		tail := map[int]iec.STRING{formatCSV: "", formatHTML: "</table></body></html>", formatXML: "</table>"}[format]
		if s.idx+LEN(tail) < size {
			s.put(tail)
			s.logStop = true
			s.step2 = 10
			s.ring(now, ucbRemove)
		}
	case format == formatCSV && s.idx+x.UCB.D_SIZE+2 < size:
		s.ring(now, ucbRemove)
		s.put(x.UCB.D_STRING)
		s.n++
		if s.n == iec.INT(x.ID_MAX) { // the last element of a line
			s.n = 0
			s.pt.BUFFER[s.idx], s.pt.BUFFER[s.idx+1] = 0x0D, 0x0A
			s.idx++
		} else {
			s.pt.BUFFER[s.idx] = sep
		}
		s.idx++
		x.NEW_FILE_SIZE = s.totalBytes + iec.UDINT(s.idx)
	case format != formatCSV && s.idx+LEN(x.UCB.D_STRING)+120 < size:
		if format == formatHTML {
			if s.logState == 0 {
				row := html[4] // odd
				switch {
				case bool(!save.HEAD):
					row = html[2]
					save.HEAD = true
				case bool(save.COLOR):
					row = html[3]
				}
				s.put(CONCAT("<TR ", row, "><TD>"))
			} else {
				s.put("<TD>")
			}
		} else {
			s.put(SEL[iec.STRING](s.logState == 0, "<entry>", "<row><entry>"))
		}
		s.logState++
		s.put(x.UCB.D_STRING)
		s.n++
		if s.n >= iec.INT(x.ID_MAX) { // the last element of a row
			s.n = 0
			if format == formatHTML {
				s.put("</TD></TR>")
				save.COLOR = !save.COLOR
			} else {
				s.put("</entry></row>")
			}
			s.logState = 0
		} else {
			s.put(SEL[iec.STRING](format == formatHTML, "</entry>", "</TD>"))
		}
		x.NEW_FILE_SIZE = s.totalBytes + iec.UDINT(s.idx)
		s.ring(now, ucbRemove)
	default: // no room: write the buffer
		s.step2 = 10
	}
}

// DLOG_STORE_FILE_CSV writes the data log X to the CSV file FILENAME, a
// format of DT_TO_STRF of the time DTI, so that a new time makes a new
// file: a line of the column names, then a line, separated by SEP, at each
// rising edge of TRIG_M, every TRIG_T (1 s at least) and when a value
// changes by its DELTA, while ENABLE. It writes the file when its buffer
// is full and after AUTO_CLOSE (5 s at least) and closes it then.
// SAVE_DATA keeps the file and counts across a restart. ERROR_C and
// ERROR_T: 1 the file, 6 the ring buffer (1 over 90% full, 2 overflow).
type DLOG_STORE_FILE_CSV struct {
	fileStore
	SEP iec.BYTE
}

// INIT resets the block and sets AUTO_CLOSE to its initial value.
func (d *DLOG_STORE_FILE_CSV) INIT() {
	*d = DLOG_STORE_FILE_CSV{fileStore: fileStore{X: d.X, SAVE_DATA: d.SAVE_DATA, AUTO_CLOSE: iec.TIME(15 * time.Second)}}
	d.defaults()
}

// Execute runs the block once.
func (d *DLOG_STORE_FILE_CSV) Execute(now time.Time) { d.run(now, formatCSV, d.SEP, [5]iec.STRING{}) }

// DLOG_STORE_FILE_HTML writes the data log X to an HTML table in the file
// FILENAME, with the caption HTML_CAPTION, the attributes HTML_TABLE of the
// table and HTML_TR_HEAD, HTML_TR_EVEN and HTML_TR_ODD of the rows; see
// DLOG_STORE_FILE_CSV.
type DLOG_STORE_FILE_HTML struct {
	fileStore
	HTML_CAPTION iec.STRING
	HTML_TABLE   iec.STRING
	HTML_TR_HEAD iec.STRING
	HTML_TR_EVEN iec.STRING
	HTML_TR_ODD  iec.STRING
}

// INIT resets the block and sets AUTO_CLOSE to its initial value.
func (d *DLOG_STORE_FILE_HTML) INIT() {
	*d = DLOG_STORE_FILE_HTML{fileStore: fileStore{X: d.X, SAVE_DATA: d.SAVE_DATA, AUTO_CLOSE: iec.TIME(15 * time.Second)}}
	d.defaults()
}

// Execute runs the block once.
func (d *DLOG_STORE_FILE_HTML) Execute(now time.Time) {
	d.run(now, formatHTML, 0, [5]iec.STRING{d.HTML_CAPTION, d.HTML_TABLE, d.HTML_TR_HEAD, d.HTML_TR_EVEN, d.HTML_TR_ODD})
}

// DLOG_STORE_FILE_XML writes the data log X to the XML file FILENAME,
// <table><row><entry>...; see DLOG_STORE_FILE_CSV.
type DLOG_STORE_FILE_XML struct {
	fileStore
}

// INIT resets the block and sets AUTO_CLOSE to its initial value.
func (d *DLOG_STORE_FILE_XML) INIT() {
	*d = DLOG_STORE_FILE_XML{fileStore: fileStore{X: d.X, SAVE_DATA: d.SAVE_DATA, AUTO_CLOSE: iec.TIME(15 * time.Second)}}
	d.defaults()
}

// Execute runs the block once.
func (d *DLOG_STORE_FILE_XML) Execute(now time.Time) { d.run(now, formatXML, 0, [5]iec.STRING{}) }
