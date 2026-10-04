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

// Package parser is the port of the OSCAT NETWORK parsers of CSV and INI
// data, in a buffer or in a file read with FILE_BLOCK, and of XML.
//
// OSCAT stops a parse after 1 ms of a scan, with a watchdog timer, and goes
// on in the next scan. The scan time of Execute does not move during a
// scan, so here a parse of a buffer ends in one scan; a parse of a file
// still waits for the reads of its FILE_SERVER.
package parser

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/file"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// The results of the parsers.
const (
	RESULT_ELEMENT     iec.BYTE = 1  // CSV: an element; INI: a section
	RESULT_LINE        iec.BYTE = 2  // CSV: an element at the end of its line; INI: a key
	RESULT_BUSY        iec.BYTE = 5  // the parse goes on in the next scan
	RESULT_END         iec.BYTE = 10 // the end of the data
	RESULT_SECTION_END iec.BYTE = 11 // INI: the key was not in the section
)

// text is a string a parser writes, a character at a time, up to
// STRING_LENGTH characters, as OSCAT does through a pointer to it.
type text struct{ b []iec.BYTE }

func (t *text) reset()             { t.b = t.b[:0] }
func (t *text) empty() bool        { return len(t.b) == 0 }
func (t *text) string() iec.STRING { return STR(t.b) }
func (t *text) add(c iec.BYTE) {
	if c > 0 && len(t.b) < int(STRING_LENGTH) {
		t.b = append(t.b, c)
	}
}

// source is where a parser reads its characters: next returns the
// character at i and whether it is the last, or busy if it must wait for
// the next scan, or end at the end of the data.
type source interface {
	next(now time.Time, i iec.UDINT) (c iec.BYTE, eof, busy, end bool)
}

// bufSource is the data of a NETWORK_BUFFER.
type bufSource struct{ pt *network.NETWORK_BUFFER }

func (b bufSource) next(now time.Time, i iec.UDINT) (iec.BYTE, bool, bool, bool) {
	if iec.INT(i) > iec.INT(b.pt.SIZE)-1 {
		return 0, false, false, true
	}
	return b.pt.BUFFER[i], i+1 == iec.UDINT(b.pt.SIZE), false, false
}

// fileSource is the data of a file, read with a FILE_BLOCK.
type fileSource struct {
	fb       *file.FILE_BLOCK
	mode     *iec.BYTE
	filename *iec.STRING
	fsd      *network.FILE_SERVER_DATA
	pt       *network.NETWORK_BUFFER
}

func (f fileSource) next(now time.Time, i iec.UDINT) (iec.BYTE, bool, bool, bool) {
	*f.mode = 1
	f.fb.MODE, f.fb.FSD, f.fb.PT, f.fb.FILENAME, f.fb.POS = f.mode, f.fsd, f.pt, f.filename, i
	f.fb.Execute(now)
	eof := i+1 == f.fsd.FILE_SIZE
	switch {
	case *f.mode > 0:
		return 0, eof, true, false
	case f.fb.ERROR > 0:
		// an error, or the end of the file
		return 0, eof, false, true
	}
	return f.fb.DATA, eof, false, false
}

// csv is CSV_PARSER_BUF and CSV_PARSER_FILE.
type csv struct {
	step     iec.INT
	c        iec.BYTE
	vAdd     iec.BYTE
	i        iec.UDINT
	state    iec.BYTE
	eof      iec.BOOL
	value    text
	watchdog timers.TON
	x        iec.BOOL
}

// run runs a scan of the parse of src.
func (p *csv) run(now time.Time, src source, SEP iec.BYTE, RUN *iec.BYTE, OFFSET *iec.UDINT, VALUE *iec.STRING, RESULT *iec.BYTE) {
	if *RUN == 0 {
		return
	}
	if p.state == 0 {
		// a new parse
		p.value.reset()
		*VALUE = ""
		p.i = *OFFSET
		*RESULT = 0
		p.state = 5
		p.step = 0
	}
	p.watchdog.IN, p.watchdog.PT = false, iec.TIME(time.Millisecond)
	p.watchdog.Execute(now)
	p.watchdog.IN = true

	for p.state == 5 {
		var busy, end bool
		p.c, p.eof, busy, end = srcNext(src, now, p.i)
		p.watchdog.Execute(now)
		if busy || bool(p.watchdog.Q) {
			break
		}
		if end {
			p.state = 10 // the end of the data
			break
		}
		p.vAdd = 0
		switch {
		case p.c == SEP: // an element
			p.state = 1
			p.x = true
		case p.c >= 32:
			p.vAdd = p.c
			if p.eof {
				p.state = 2 // an element and the end of the line
			} else {
				p.x = true
			}
		case bool(p.x):
			p.state = 2 // an element and the end of the line
			p.x = false
		}
		p.value.add(p.vAdd)
		*VALUE = p.value.string()
		p.i++
	}

	switch p.state {
	case 1, 2, 10:
		*RESULT = p.state
		*OFFSET = p.i
		*RUN = 0
		p.state = 0
	case 5:
		*RESULT = 5
	}
}

// srcNext reads the character at i, with eof as an iec.BOOL.
func srcNext(src source, now time.Time, i iec.UDINT) (iec.BYTE, iec.BOOL, bool, bool) {
	c, eof, busy, end := src.next(now, i)
	return c, iec.BOOL(eof), busy, end
}

// CSV_PARSER_BUF reads the next element of the CSV data in PT, separated by
// SEP, from OFFSET into VALUE: RUN set to non-zero starts it and is 0 when
// done, with RESULT 1 for an element, 2 for one at the end of its line and
// 10 at the end of the data, and OFFSET past it.
type CSV_PARSER_BUF struct {
	SEP    iec.BYTE
	RUN    *iec.BYTE
	OFFSET *iec.UDINT
	VALUE  *iec.STRING // STRING(STRING_LENGTH)
	PT     *network.NETWORK_BUFFER
	RESULT iec.BYTE

	csv
}

// INIT resets the block.
func (p *CSV_PARSER_BUF) INIT() {
	*p = CSV_PARSER_BUF{RUN: p.RUN, OFFSET: p.OFFSET, VALUE: p.VALUE, PT: p.PT}
}

// Execute runs the block once.
func (p *CSV_PARSER_BUF) Execute(now time.Time) {
	if p.RUN == nil || p.OFFSET == nil || p.VALUE == nil || p.PT == nil {
		return
	}
	p.run(now, bufSource{p.PT}, p.SEP, p.RUN, p.OFFSET, p.VALUE, &p.RESULT)
}

// CSV_PARSER_FILE is CSV_PARSER_BUF for the data of the file FILENAME, read
// with the FILE_SERVER of FSD and PT; RESULT is 5 while it waits for a read.
type CSV_PARSER_FILE struct {
	SEP      iec.BYTE
	FILENAME *iec.STRING
	FSD      *network.FILE_SERVER_DATA
	RUN      *iec.BYTE
	OFFSET   *iec.UDINT
	VALUE    *iec.STRING // STRING(STRING_LENGTH)
	PT       *network.NETWORK_BUFFER
	RESULT   iec.BYTE

	csv
	fb   file.FILE_BLOCK
	mode iec.BYTE
}

// INIT resets the block.
func (p *CSV_PARSER_FILE) INIT() {
	*p = CSV_PARSER_FILE{FILENAME: p.FILENAME, FSD: p.FSD, RUN: p.RUN, OFFSET: p.OFFSET, VALUE: p.VALUE, PT: p.PT}
}

// Execute runs the block once.
func (p *CSV_PARSER_FILE) Execute(now time.Time) {
	if p.FILENAME == nil || p.FSD == nil || p.RUN == nil || p.OFFSET == nil || p.VALUE == nil || p.PT == nil {
		return
	}
	p.run(now, fileSource{&p.fb, &p.mode, p.FILENAME, p.FSD, p.PT}, p.SEP, p.RUN, p.OFFSET, p.VALUE, &p.RESULT)
}

// ini is INI_PARSER_BUF and INI_PARSER_FILE.
type ini struct {
	step     iec.INT
	c        iec.BYTE
	kAdd     iec.BYTE
	vAdd     iec.BYTE
	i        iec.UDINT
	state    iec.BYTE
	eof      iec.BOOL
	key      text
	value    text
	watchdog timers.TON
}

// run runs a scan of the parse of src.
func (p *ini) run(now time.Time, src source, STR *iec.STRING, RUN *iec.BYTE, OFFSET *iec.UDINT, KEY, VALUE *iec.STRING, RESULT *iec.BYTE) {
	if *RUN == 0 {
		return
	}
	if p.state == 0 {
		// a new parse
		p.value.reset()
		p.key.reset()
		*VALUE, *KEY = "", ""
		p.i = *OFFSET
		*RESULT = 0
		p.state = 5
		p.step = 0
	}
	p.watchdog.IN, p.watchdog.PT = false, iec.TIME(time.Millisecond)
	p.watchdog.Execute(now)
	p.watchdog.IN = true

	for p.state == 5 {
		var busy, end bool
		p.c, p.eof, busy, end = srcNext(src, now, p.i)
		p.watchdog.Execute(now)
		if busy || bool(p.watchdog.Q) {
			break
		}
		if end {
			p.state = 10 // the end of the data
			break
		}
		p.kAdd, p.vAdd = 0, 0
		c := p.c
		switch p.step {
		case 0: // what the line is
			switch {
			case c == 91: // [ a section
				p.step = 200
			case c == 59 || c == 35: // ; or # a comment
				p.step = 100
			case c > 32: // a key
				p.kAdd = c
				p.step = 300
			}
		case 100: // a comment: find the end of the line
			if c < 32 {
				p.state = 3
			}
		case 200: // a section: find the ]
			switch {
			case c == 93:
				if !p.key.empty() {
					p.state = 1
				} else {
					// an empty section: go on
					p.key.reset()
					p.step = 0
				}
			case c < 32 || bool(p.eof):
				// a broken section line: go on
				p.key.reset()
				p.step = 0
			default:
				p.kAdd = c
			}
		case 300: // the name of a key
			switch {
			case c < 32:
				// a broken key line: go on
				p.key.reset()
				p.step = 0
			case c == 61: // =
				p.step = 330
			default:
				p.kAdd = c
			}
		case 330: // the start of the value
			switch {
			case c < 32:
				p.state = 2 // a key without a value
			case bool(p.eof):
				p.state = 2
				p.vAdd = c
			default:
				p.vAdd = c
				p.step = 350
			}
		case 350: // the end of the value
			switch {
			case c < 32:
				p.state = 2
			case bool(p.eof):
				p.vAdd = c
				p.state = 2
			default:
				p.vAdd = c
			}
		}
		p.value.add(p.vAdd)
		p.key.add(p.kAdd)
		*VALUE, *KEY = p.value.string(), p.key.string()
		p.i++
	}

	// A section or key that is not the one asked for leaves the state as it
	// is, in OSCAT: the parse does not go on.
	switch p.state {
	case 1: // a section
		if *RUN == 1 || *RUN == 3 {
			if *KEY == *STR || LEN(*STR) == 0 {
				*RESULT = 1
				*OFFSET = p.i
				*RUN = 0
				p.state = 0
			}
		} else {
			p.key.reset()
			*KEY = ""
			*RESULT = 11 // the key is not in the section; OFFSET stays
			*RUN = 0
			p.state = 0
		}
	case 2: // a key
		if (*RUN == 2 || *RUN == 3) && (*KEY == *STR || LEN(*STR) == 0) {
			*RESULT = 2
			*OFFSET = p.i
			*RUN = 0
			p.state = 0
		}
	case 3: // a comment
		*OFFSET = p.i
		p.state = 0
	case 5:
		*RESULT = 5
	case 10:
		*RESULT = 10
		*OFFSET = p.i
		*RUN = 0
		p.state = 0
	}
}

// INI_PARSER_BUF reads the INI data in PT from OFFSET: RUN set to 1 finds
// the section STR, 2 the key STR, 3 the next section or key, STR ” for
// any; RUN is 0 when done, with RESULT 1 for a section, its name in KEY, 2
// for a key, its value in VALUE, 10 at the end of the data and 11 for a key
// not found before the next section, and OFFSET past it.
type INI_PARSER_BUF struct {
	STR    *iec.STRING // STRING(STRING_LENGTH)
	RUN    *iec.BYTE
	OFFSET *iec.UDINT
	KEY    *iec.STRING // STRING(STRING_LENGTH)
	VALUE  *iec.STRING // STRING(STRING_LENGTH)
	PT     *network.NETWORK_BUFFER
	RESULT iec.BYTE

	ini
}

// INIT resets the block.
func (p *INI_PARSER_BUF) INIT() {
	*p = INI_PARSER_BUF{STR: p.STR, RUN: p.RUN, OFFSET: p.OFFSET, KEY: p.KEY, VALUE: p.VALUE, PT: p.PT}
}

// Execute runs the block once.
func (p *INI_PARSER_BUF) Execute(now time.Time) {
	if p.STR == nil || p.RUN == nil || p.OFFSET == nil || p.KEY == nil || p.VALUE == nil || p.PT == nil {
		return
	}
	p.run(now, bufSource{p.PT}, p.STR, p.RUN, p.OFFSET, p.KEY, p.VALUE, &p.RESULT)
}

// INI_PARSER_FILE is INI_PARSER_BUF for the data of the file FILENAME,
// read with the FILE_SERVER of FSD and PT; RESULT is 5 while it waits for a
// read.
type INI_PARSER_FILE struct {
	STR      *iec.STRING // STRING(STRING_LENGTH)
	FILENAME *iec.STRING
	FSD      *network.FILE_SERVER_DATA
	RUN      *iec.BYTE
	OFFSET   *iec.UDINT
	KEY      *iec.STRING // STRING(STRING_LENGTH)
	VALUE    *iec.STRING // STRING(STRING_LENGTH)
	PT       *network.NETWORK_BUFFER
	RESULT   iec.BYTE

	ini
	fb   file.FILE_BLOCK
	mode iec.BYTE
}

// INIT resets the block.
func (p *INI_PARSER_FILE) INIT() {
	*p = INI_PARSER_FILE{STR: p.STR, FILENAME: p.FILENAME, FSD: p.FSD, RUN: p.RUN, OFFSET: p.OFFSET, KEY: p.KEY, VALUE: p.VALUE, PT: p.PT}
}

// Execute runs the block once.
func (p *INI_PARSER_FILE) Execute(now time.Time) {
	if p.STR == nil || p.FILENAME == nil || p.FSD == nil || p.RUN == nil || p.OFFSET == nil || p.KEY == nil || p.VALUE == nil || p.PT == nil {
		return
	}
	p.run(now, fileSource{&p.fb, &p.mode, p.FILENAME, p.FSD, p.PT}, p.STR, p.RUN, p.OFFSET, p.KEY, p.VALUE, &p.RESULT)
}
