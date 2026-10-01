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

// Package list is the port of the OSCAT BASIC list processing functions. A
// list is a string of elements that each start with the separator SEP, such
// as ",a,b,c". A list holds up to LIST_LENGTH characters.
package list

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/royaljelly/iec"
)

// mem returns the characters of a list as OSCAT sees its memory: position
// 1 is the first character, and the list is followed by 0s.
func mem(list iec.STRING) []iec.BYTE {
	m := make([]iec.BYTE, int(LIST_LENGTH)+2)
	copy(m[1:], CHARS(list))
	return m
}

// text returns the string in memory m, up to its first 0.
func text(m []iec.BYTE) iec.STRING {
	return STR(m[1:])
}

// LIST_ADD adds the element ins at the end of the list. It returns false,
// and leaves the list, if the element does not fit.
func LIST_ADD(sep iec.BYTE, ins iec.STRING, list *iec.STRING) iec.BOOL {
	ins = str.CHR_TO_STRING(sep) + ins
	if LEN(*list)+LEN(ins) > LIST_LENGTH {
		return false
	}
	*list += ins
	return true
}

// LIST_CLEAN removes the empty elements of the list.
func LIST_CLEAN(sep iec.BYTE, list *iec.STRING) iec.BOOL {
	pt := mem(*list)
	write := 1
	var last iec.BYTE
	for read := 1; read <= int(LIST_LENGTH); read++ {
		c := pt[read]
		if c == 0 {
			break
		}
		// Copy the character unless it is a second separator.
		if c != sep || sep != last {
			pt[write] = c
			write++
		}
		last = c
	}
	// A separator at the end is an empty element.
	if last == sep {
		write--
	}
	if write <= int(STRING_LENGTH) {
		pt[write] = 0
	}
	*list = text(pt)
	return true
}

// LIST_GET returns the element pos of the list; the first element is 1.
func LIST_GET(sep iec.BYTE, pos iec.INT, list *iec.STRING) iec.STRING {
	pt := mem(*list)
	var out []iec.BYTE
	var cnt iec.INT
	for i := 1; ; {
		c := pt[i]
		if cnt == pos {
			if c == sep {
				break
			}
			out = append(out, c)
		} else if c == sep {
			cnt++
		}
		i++
		if i == int(LIST_LENGTH) || c == 0 {
			break
		}
	}
	return STR(out)
}

// LIST_INSERT inserts the element ins at the position pos of the list,
// adding empty elements if the list is shorter. It returns false, and
// leaves the list, if the element does not fit.
func LIST_INSERT(sep iec.BYTE, pos iec.INT, ins iec.STRING, list *iec.STRING) iec.BOOL {
	if LEN(ins)+1+LEN(*list) > LIST_LENGTH {
		return false
	}
	pt := mem(*list)
	read := 1
	var cnt iec.INT = 1
	for read < int(LIST_LENGTH) {
		if cnt >= pos {
			*list = INSERT(text(pt), str.CHR_TO_STRING(sep)+ins, iec.INT(read-1))
			return true
		}
		if pt[read] == 0 {
			// The list is too short: add an empty element.
			pt[read] = sep
			pt[read+1] = 0
		}
		read++
		if pt[read] == sep || pt[read] == 0 {
			cnt++
		}
	}
	*list = text(pt)
	return true
}

// LIST_LEN returns the number of elements of the list.
func LIST_LEN(sep iec.BYTE, list *iec.STRING) iec.INT {
	pt := mem(*list)
	var n iec.INT
	for pos := 1; ; pos++ {
		c := pt[pos]
		if c == sep {
			n++
		}
		if c == 0 || pos+1 > int(LIST_LENGTH) {
			return n
		}
	}
}

// LIST_NEXT returns the next element of LIST on LEL each time it runs,
// starting with the first element after RST or its first run. At the end of
// the list LEL is empty and NUL is true.
type LIST_NEXT struct {
	SEP  iec.BYTE
	RST  iec.BOOL
	LIST *iec.STRING
	LEL  iec.STRING
	NUL  iec.BOOL

	pos         int
	initialized bool
}

// INIT resets the block.
func (l *LIST_NEXT) INIT() { *l = LIST_NEXT{LIST: l.LIST, pos: 1, initialized: true} }

// Execute runs the block once.
func (l *LIST_NEXT) Execute(now time.Time) {
	if !l.initialized {
		l.initialized = true
		l.pos = 1
	}
	if l.LIST == nil {
		return
	}
	pt := mem(*l.LIST)
	if l.RST {
		l.pos = 1
	}
	if l.pos >= len(pt) || pt[l.pos] == 0 || l.pos == int(LIST_LENGTH) {
		l.LEL = ""
		l.NUL = true
		return
	}
	l.NUL = false
	var out []iec.BYTE
	for l.pos = l.pos + 1; l.pos <= int(LIST_LENGTH); l.pos++ {
		c := pt[l.pos]
		if c == 0 || c == l.SEP {
			break
		}
		out = append(out, c)
	}
	l.LEL = STR(out)
}

// LIST_RETRIEVE returns the element pos of the list and removes it from the
// list.
func LIST_RETRIEVE(sep iec.BYTE, pos iec.INT, list *iec.STRING) iec.STRING {
	if pos <= 0 {
		return ""
	}
	pt := mem(*list)
	var out []iec.BYTE
	w := 1
	var cnt iec.INT
	for i := 1; i <= int(LIST_LENGTH); i++ {
		c := pt[i]
		switch {
		case c == 0:
			if cnt < pos {
				pt[w+1] = 0
			} else {
				pt[w] = 0
			}
			*list = text(pt)
			return STR(out)
		case cnt == pos && c != sep:
			// The element.
			out = append(out, pt[i])
		case cnt >= pos:
			pt[w] = c
			w++
		default:
			w = i
		}
		if c == sep {
			cnt++
		}
	}
	*list = text(pt)
	return STR(out)
}

// LIST_RETRIEVE_LAST returns the last element of the list and removes it
// from the list.
func LIST_RETRIEVE_LAST(sep iec.BYTE, list *iec.STRING) iec.STRING {
	pt := mem(*list)
	last := 1
	for i := 1; i <= int(LIST_LENGTH); i++ {
		c := pt[i]
		if c == 0 {
			break
		} else if c == sep {
			last = i
		}
	}
	out := MID(*list, LIST_LENGTH, iec.INT(last+1))
	pt[last] = 0
	*list = text(pt)
	return out
}
