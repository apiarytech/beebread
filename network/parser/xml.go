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

package parser

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/basic/logic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// The types of the items XML_READER finds, in CTRL.TYP.
const (
	XML_UNKNOWN   iec.INT = 0
	XML_ELEMENT   iec.INT = 1  // an element
	XML_CLOSE     iec.INT = 2  // the end of an element
	XML_TEXT      iec.INT = 3  // the text of an element
	XML_ATTRIBUTE iec.INT = 4  // an attribute
	XML_PI        iec.INT = 5  // a processing instruction
	XML_CDATA     iec.INT = 12 // <![CDATA[...]]>
	XML_COMMENT   iec.INT = 13 // <!--...-->
	XML_DTD       iec.INT = 14 // a document type declaration
	XML_WATCHDOG  iec.INT = 98 // the watchdog stopped the parse
	XML_END       iec.INT = 99 // the end of the data
)

// XML_READER parses the XML in BUF, from CTRL.START_POS to CTRL.STOP_POS,
// an item a run: a CTRL.COMMAND with bit 15 set starts it, and its bits 1
// to 14 select the types of the items it stops at, CTRL.TYP; it stops at
// the end of the data with XML_END. CTRL has the element, its path from
// the root, its level, the attribute and the value, and BLOCK1 and BLOCK2
// the positions of the name and value in BUF.
type XML_READER struct {
	CTRL *network.XML_CONTROL
	BUF  *network.NW_BUF_LONG

	index        iec.INT
	index2       iec.INT
	stop         iec.INT
	mode         iec.INT
	command      iec.WORD
	c            iec.BYTE
	s1           iec.INT
	e1           iec.INT
	pathOverflow iec.BOOL
	emptyTag     iec.BOOL
	sCode        iec.STRING // STRING(10)
	eCode        iec.STRING // STRING(10)
	watchdog     timers.TON
}

// INIT resets the block.
func (x *XML_READER) INIT() { *x = XML_READER{CTRL: x.CTRL, BUF: x.BUF} }

// at returns the byte at i of BUF, 0 outside it.
func (x *XML_READER) at(i iec.INT) iec.BYTE {
	if i < 0 || int(i) >= len(x.BUF) {
		return 0
	}
	return x.BUF[i]
}

// text returns the text of BUF from start to stop.
func (x *XML_READER) text(start, stop iec.UINT) iec.STRING {
	return buffer.BUFFER_TO_STRING(x.BUF[:], iec.UINT(x.stop+1), start, stop)
}

// wants reports whether the command asks to stop at the type typ.
func (x *XML_READER) wants(typ iec.INT) bool {
	return bool(logic.BIT_OF_DWORD(iec.DWORD(x.command), typ))
}

// Execute runs the block once.
func (x *XML_READER) Execute(now time.Time) {
	if x.CTRL == nil || x.BUF == nil {
		return
	}
	ctrl := x.CTRL
	if logic.BIT_OF_DWORD(iec.DWORD(ctrl.COMMAND), 15) {
		x.command = ctrl.COMMAND
		ctrl.COMMAND = 0
		x.index = iec.INT(ctrl.START_POS)
		x.stop = iec.INT(ctrl.STOP_POS)
		x.mode = 100
		x.pathOverflow = false
		ctrl.TYP = 0
		ctrl.COUNT, ctrl.LEVEL = 0, 0
		ctrl.ATTRIBUTE, ctrl.ELEMENT, ctrl.PATH, ctrl.VALUE = "", "", "", ""
		x.watchdog.PT = ctrl.WATCHDOG
	}
	if x.index < 0 {
		return
	}
	ctrl.BLOCK1_START, ctrl.BLOCK1_STOP, ctrl.BLOCK2_START, ctrl.BLOCK2_STOP = 0, 0, 0, 0

	x.watchdog.IN = false
	x.watchdog.Execute(now)

	for {
		x.watchdog.IN = ctrl.WATCHDOG > 0
		x.watchdog.Execute(now)
		if x.watchdog.Q {
			ctrl.TYP = XML_WATCHDOG
			return
		}
		if x.index > x.stop {
			// the end of the data
			ctrl.COUNT++
			ctrl.TYP = XML_END
			x.mode = 0
			x.index = 0
			return
		}

		switch x.mode {
		case 100: // an element
			x.s1, x.e1 = 0, 0
			// The start of the element.
			for x.index <= x.stop {
				if x.at(x.index) == 60 { // <
					x.index++
					x.c = x.at(x.index)
					switch x.c {
					case 47: // /
						x.mode = 300
					case 33: // !
						x.index++
						x.mode = 500
					default:
						x.s1 = x.index
					}
					break
				}
				x.index++
			}

			// The element.
			ctrl.TYP = 0
			if x.s1 > 0 {
			element:
				for x.index <= x.stop && x.e1 == 0 {
					x.c = x.at(x.index)
					if x.c <= 32 {
						// an element with attributes, or an empty one
						x.index2 = x.index
						for x.index2 <= x.stop && x.at(x.index2) <= 32 {
							x.index2++
						}
						if x.at(x.index2) == 47 { // /
							x.index = x.index2
							x.e1 = x.index - 2
							ctrl.TYP = XML_ELEMENT
							x.mode = 300
							break element
						}
						x.e1 = x.index - 1
						ctrl.TYP = XML_ELEMENT
						x.mode = 400
					} else if x.c == 62 { // >
						x.index2 = x.index - 1
						if x.at(x.index2) == 47 { // an empty element <x/>
							x.e1 = x.index2 - 1
							ctrl.TYP = XML_ELEMENT
							x.emptyTag = true
							x.mode = 200
							break element
						}
						x.e1 = x.index - 1
						x.index2 = x.index + 1
						// What follows.
						for x.index2 <= x.stop {
							x.c = x.at(x.index2)
							if x.c == 60 { // <
								if x.at(x.index2+1) == 47 { // /
									x.emptyTag = true
									x.mode = 200
								} else {
									x.mode = 100
								}
								ctrl.TYP = XML_ELEMENT
								break
							} else if x.c > 32 {
								ctrl.TYP = XML_ELEMENT // an element with text
								x.mode = 200
								break
							}
							x.index2++
						}
					}
					x.index++
				}
			}

			if ctrl.TYP > 0 {
				// a processing instruction
				if x.at(x.s1) == 63 { // ?
					x.s1++
					ctrl.TYP = XML_PI
				}
				ctrl.BLOCK1_START, ctrl.BLOCK1_STOP = iec.UINT(x.s1), iec.UINT(x.e1)
				ctrl.ELEMENT = x.text(ctrl.BLOCK1_START, ctrl.BLOCK1_STOP)
				ctrl.COUNT++
				ctrl.LEVEL++
				if !x.pathOverflow {
					if LEN(ctrl.PATH)+LEN(ctrl.ELEMENT)+1 > 250 {
						x.pathOverflow = true
						ctrl.PATH = "OVERFLOW"
					} else {
						ctrl.PATH = CONCAT(ctrl.PATH, "/", ctrl.ELEMENT)
					}
				}
				if x.wants(ctrl.TYP) {
					return
				}
			}

		case 200: // the text
			ctrl.VALUE = ""
			if !x.emptyTag {
				x.s1 = x.index
				for x.index <= x.stop && x.at(x.index) != 60 {
					x.index++
				}
				x.e1 = x.index - 1
				ctrl.BLOCK1_START, ctrl.BLOCK1_STOP = iec.UINT(x.s1), iec.UINT(x.e1)
				ctrl.VALUE = x.text(ctrl.BLOCK1_START, ctrl.BLOCK1_STOP)
			}
			x.emptyTag = false
			x.mode = 300
			ctrl.COUNT++
			ctrl.TYP = XML_TEXT
			if x.wants(ctrl.TYP) {
				return
			}

		case 300: // the end of an element
			for x.index <= x.stop {
				x.c = x.at(x.index)
				if x.c == 62 { // >
					x.index++
					x.s1 = str.FINDB(ctrl.PATH, "/")
					if !x.pathOverflow {
						if x.s1 > 1 {
							ctrl.ELEMENT = RIGHT(ctrl.PATH, LEN(ctrl.PATH)-x.s1)
							ctrl.PATH = LEFT(ctrl.PATH, x.s1-1)
						} else {
							ctrl.ELEMENT = RIGHT(ctrl.PATH, LEN(ctrl.PATH)-1)
							ctrl.PATH = ""
						}
					}
					ctrl.LEVEL--
					ctrl.COUNT++
					ctrl.TYP = XML_CLOSE
					x.mode = 100
					if x.wants(ctrl.TYP) {
						return
					}
					break
				}
				x.index++
			}

		case 400: // an attribute
			for x.index <= x.stop && x.at(x.index) <= 32 {
				x.index++
			}
			// The name, up to the =.
			x.e1 = 0
			x.s1 = x.index
			for x.index <= x.stop {
				if x.at(x.index) == 61 {
					x.e1 = x.index - 1
					break
				}
				x.index++
			}
			if x.e1 > 0 {
				ctrl.BLOCK1_START, ctrl.BLOCK1_STOP = iec.UINT(x.s1), iec.UINT(x.e1)
				ctrl.ATTRIBUTE = x.text(ctrl.BLOCK1_START, ctrl.BLOCK1_STOP)
				// The value, between quotes.
				x.index += 2
				x.e1 = 0
				x.s1 = x.index
				for x.index <= x.stop && x.e1 == 0 {
					x.c = x.at(x.index)
					if x.c == 34 || x.c == 39 { // " or '
						x.e1 = x.index - 1
					}
					x.index++
				}
				if x.e1 > 0 {
					ctrl.BLOCK2_START, ctrl.BLOCK2_STOP = iec.UINT(x.s1), iec.UINT(x.e1)
					// An empty value has a start after its stop: ''.
					ctrl.VALUE = x.text(ctrl.BLOCK2_START, ctrl.BLOCK2_STOP)
					for x.index <= x.stop && x.at(x.index) <= 32 {
						x.index++
					}
					x.c = x.at(x.index)
					switch {
					case x.c == 62: // >
						x.index++
						x.index2 = x.index
						for x.index <= x.stop && x.at(x.index) <= 32 {
							x.index++
						}
						if x.at(x.index) == 60 { // <
							x.mode = 100
						} else {
							x.index = x.index2
							x.mode = 200
						}
					case x.c == 47 || x.c == 63: // / or ?
						x.mode = 300
					default:
						x.mode = 400
					}
					ctrl.COUNT++
					ctrl.TYP = XML_ATTRIBUTE
					if x.wants(ctrl.TYP) {
						return
					}
				}
			}

		case 500: // <!...>: CDATA, a comment or a declaration
			ctrl.TYP = 0
			x.s1 = 0
			x.index2 = 11
			for x.s1 == 0 {
				x.index2++
				ctrl.TYP = x.index2
				if x.index2 == 14 {
					// any other declaration
					x.eCode = ""
					x.s1 = x.index
					break
				}
				switch x.index2 {
				case 12:
					x.sCode, x.eCode = "[CDATA[", "]]>"
				case 13:
					x.sCode, x.eCode = "--", "-->"
				}
				x.e1 = LEN(x.sCode)
				if buffer.BUFFER_SEARCH(x.BUF[:], x.index+x.e1, x.sCode, x.index, false) == x.index {
					x.s1 = x.index + x.e1
				}
			}
			ctrl.BLOCK1_START = iec.UINT(x.s1)
			if LEN(x.eCode) > 0 {
				x.e1 = buffer.BUFFER_SEARCH(x.BUF[:], x.stop+1, x.eCode, x.index, false)
				ctrl.BLOCK1_STOP = iec.UINT(x.e1 - 1)
				x.index = x.e1 + LEN(x.eCode)
			} else {
				// The >, counting the <> pairs and skipping quotes.
				x.index2 = 1
				for x.index <= x.stop && x.index2 > 0 {
					x.c = x.at(x.index)
					switch x.c {
					case 60:
						x.index2++
					case 62:
						x.index2--
					case 34, 39:
						x.index++
						for x.index <= x.stop {
							x.c = x.at(x.index)
							if x.c == 34 || x.c == 39 {
								break
							}
							x.index++
						}
					}
					x.index++
				}
				ctrl.BLOCK1_STOP = iec.UINT(x.index - 2)
			}
			ctrl.VALUE = x.text(ctrl.BLOCK1_START, ctrl.BLOCK1_STOP)
			ctrl.COUNT++
			x.mode = 100
			if x.wants(ctrl.TYP) {
				return
			}
		}
	}
}
