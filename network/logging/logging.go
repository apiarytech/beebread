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

// Package logging is the port of the OSCAT NETWORK log: LOG_MSG adds
// messages to the global log network.LOG_CL, formatted by PRINT_SF, and
// LOG_VIEWPORT pages through a log.
package logging

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// PRINT_SF replaces ~1..~9 in STR with the arguments PRINTF_DATA[1..9],
// '..' for one that would make STR longer than LOG_SIZE. ~ followed by
// anything else is removed with it.
type PRINT_SF struct {
	PRINTF_DATA *network.PRINTF_DATA
	STR         *iec.STRING // STRING(LOG_SIZE)

	pos iec.INT
	c   iec.INT
	src iec.STRING // default '~'
	run iec.BOOL
}

// INIT resets the block.
func (p *PRINT_SF) INIT() { *p = PRINT_SF{PRINTF_DATA: p.PRINTF_DATA, STR: p.STR, src: "~"} }

// Execute runs the block once.
func (p *PRINT_SF) Execute(now time.Time) {
	if p.PRINTF_DATA == nil || p.STR == nil {
		return
	}
	p.src = "~"
	args := p.PRINTF_DATA
	if LEN(*p.STR) == 0 {
		return
	}
	p.run = true
	for p.run {
		p.pos = FIND(*p.STR, p.src)
		args[9] = ""
		if p.pos == 0 {
			p.run = false
			continue
		}
		p.c = iec.INT(str.CODE(*p.STR, p.pos+1))
		if str.ISC_NUM(iec.BYTE(p.c)) {
			p.c -= 48
			// ~1..~9 are the arguments 1..9; OSCAT reads outside the
			// arguments for ~0, which is '' here.
			if p.c >= 1 {
				args[9] = args[p.c-1]
			}
		}
		if LEN(args[9])+LEN(*p.STR) > network.LOG_SIZE {
			args[9] = ".."
		}
		args[10] = *p.STR
		*p.STR = network.STRING_N(REPLACE(args[10], args[9], 2, p.pos), network.LOG_SIZE)
	}
}

// LOG_MSG adds the message LOG_CL.NEW_MSG, with its arguments
// LOG_CL.PRINTF, to the log network.LOG_CL, a ring of LOG_CL.SIZE
// messages, and clears it. LOG_CL.RESET clears the log.
type LOG_MSG struct {
	fbPrintSF PRINT_SF
	idx       iec.INT
}

// INIT resets the block.
func (l *LOG_MSG) INIT() { *l = LOG_MSG{} }

// Execute runs the block once.
func (l *LOG_MSG) Execute(now time.Time) {
	lc := &network.LOG_CL
	// OSCAT shifts the option left, not right, so the level is always 0
	// and every message is logged.
	if iec.BYTE(SHL(lc.NEW_MSG_OPTION, 16)) > lc.LEVEL {
		return
	}
	l.fbPrintSF.PRINTF_DATA, l.fbPrintSF.STR = &lc.PRINTF, &lc.NEW_MSG
	l.fbPrintSF.Execute(now)

	if lc.RESET {
		lc.RESET = false
		lc.RING_MODE = false
		lc.IDX = 0
		lc.UPDATE_COUNT++
	}
	if lc.SIZE > 0 && LEN(lc.NEW_MSG) > 0 {
		if lc.IDX >= lc.SIZE {
			// the log wraps around
			lc.RING_MODE = true
			lc.IDX = 0
		}
		lc.UPDATE_COUNT++
		lc.IDX++
		l.idx = lc.IDX
		if l.idx >= 0 && int(l.idx) < len(lc.MSG) {
			lc.MSG[l.idx] = lc.NEW_MSG
			lc.MSG_OPTION[l.idx] = lc.NEW_MSG_OPTION
		}
		lc.NEW_MSG = ""
	}
}

// LOG_VIEWPORT pages through the log LC: LV.LINE_ARRAY[1..LV.COUNT] are the
// indexes in LC.MSG of the LV.COUNT messages to show, 0 for none, and
// LV.UPDATE is set when they change. LV.MOVE_TO_X moves the page: 30000 to
// the oldest message, 30001 to the newest, 30002 and 30003 a page down
// (newer) and up (older), and any other value by that many messages.
type LOG_VIEWPORT struct {
	LC *network.LOG_CONTROL
	LV *network.US_LOG_VIEWPORT

	pos       iec.INT
	count     iec.INT
	idx       iec.INT
	base      iec.INT
	updatePos iec.BOOL
}

// INIT resets the block.
func (l *LOG_VIEWPORT) INIT() { *l = LOG_VIEWPORT{LC: l.LC, LV: l.LV} }

// Execute runs the block once.
func (l *LOG_VIEWPORT) Execute(now time.Time) {
	if l.LC == nil || l.LV == nil {
		return
	}
	lc, lv := l.LC, l.LV
	if lv.MOVE_TO_X != 0 {
		switch lv.MOVE_TO_X {
		case 30000: // the oldest message
			l.pos = 1
		case 30001: // the newest message
			l.pos = lc.SIZE
		case 30002: // a page down, to newer messages
			l.pos += lv.COUNT
		case 30003: // a page up, to older messages
			l.pos -= lv.COUNT
		default:
			l.pos += lv.MOVE_TO_X
		}
		l.updatePos = true
		lv.MOVE_TO_X = 0
	}

	if lv.UPDATE_COUNT != lc.UPDATE_COUNT || l.updatePos {
		lv.UPDATE = true
		lv.UPDATE_COUNT = lc.UPDATE_COUNT
		l.updatePos = false

		switch {
		case bool(lc.RING_MODE):
			l.pos = LIMIT(1, l.pos, lc.SIZE-lv.COUNT+1)
		case lc.IDX > lv.COUNT:
			l.pos = LIMIT(1, l.pos, lc.IDX-lv.COUNT+1)
		default:
			l.pos = 1
		}

		l.base = lc.IDX - l.pos - lv.COUNT + 1
		for l.count = 1; l.count <= lv.COUNT && int(l.count) <= len(lv.LINE_ARRAY); l.count++ {
			l.idx = l.base + l.count
			if l.idx < 1 {
				if lc.RING_MODE {
					l.idx += lc.SIZE
				} else {
					// an empty line
					l.idx = 0
				}
			}
			lv.LINE_ARRAY[l.count-1] = l.idx
		}
	}
}
