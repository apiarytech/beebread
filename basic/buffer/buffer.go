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

// Package buffer is the port of the OSCAT BASIC buffer management
// functions. A buffer is an array of bytes, which OSCAT passes with
// ADR(array) and SIZEOF(array) and the port passes as a slice and its size.
// A size larger than the slice is limited to the slice.
package buffer

import (
	. "github.com/apiarytech/beebread/basic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/royaljelly/iec"
)

// buf gives access to the bytes of a buffer that are within its slice.
type buf []iec.BYTE

func (b buf) at(i int) iec.BYTE {
	if i < 0 || i >= len(b) {
		return 0
	}
	return b[i]
}

func (b buf) set(i int, v iec.BYTE) {
	if i >= 0 && i < len(b) {
		b[i] = v
	}
}

// BUFFER_CLEAR_ sets the first size bytes of pt to 0.
func BUFFER_CLEAR_(pt []iec.BYTE, size iec.UINT) iec.BOOL {
	return BUFFER_INIT_(pt, size, 0)
}

// BUFFER_INIT_ sets the first size bytes of pt to init.
func BUFFER_INIT_(pt []iec.BYTE, size iec.UINT, init iec.BYTE) iec.BOOL {
	n := min(int(size), len(pt))
	for i := 0; i < n; i++ {
		pt[i] = init
	}
	return true
}

// BUFFER_INSERT_ inserts s at the position pos of the buffer of size bytes,
// moving the bytes after it towards the end, and returns the position after
// it.
func BUFFER_INSERT_(s iec.STRING, pos iec.INT, pt []iec.BYTE, size iec.UINT) iec.INT {
	b := buf(pt)
	lx := int(LEN(s))
	end := int(pos) + lx
	for i := int(size) - 1; i >= end; i-- {
		b.set(i, b.at(i-lx))
	}
	return STRING_TO_BUFFER_(s, pos, pt, size)
}

// BUFFER_UPPERCASE_ converts the first size bytes of pt to upper case.
func BUFFER_UPPERCASE_(pt []iec.BYTE, size iec.INT) iec.BOOL {
	n := min(int(size), len(pt))
	for i := 0; i < n; i++ {
		pt[i] = str.TO_UPPER(pt[i])
	}
	return true
}

// STRING_TO_BUFFER_ copies s into the buffer of size bytes at the position
// pos and returns the position after it.
func STRING_TO_BUFFER_(s iec.STRING, pos iec.INT, pt []iec.BYTE, size iec.UINT) iec.INT {
	b := buf(pt)
	c := CHARS(s)
	end := min(int(pos)+len(c), int(size))
	if end > 0 {
		end--
	}
	i := int(pos)
	for ; i <= end; i++ {
		// Past the string OSCAT copies its closing 0.
		k := i - int(pos)
		var v iec.BYTE
		if k < len(c) {
			v = c[k]
		}
		b.set(i, v)
	}
	return iec.INT(i)
}

// BUFFER_COMP returns the first position, from start on, where the buffer
// pt2 of size2 bytes is found in the buffer pt1 of size1 bytes, or -1.
func BUFFER_COMP(pt1 []iec.BYTE, size1 iec.INT, pt2 []iec.BYTE, size2 iec.INT, start iec.INT) iec.INT {
	b1, b2 := buf(pt1), buf(pt2)
	if size2 <= size1 {
		end := int(size1 - size2)
		first := b2.at(0)
		for i := int(start); i <= end; i++ {
			if b1.at(i) != first {
				continue
			}
			j := 1
			for j < int(size2) && b2.at(j) == b1.at(j+i) {
				j++
			}
			if j == int(size2) {
				return iec.INT(i)
			}
		}
	}
	return -1
}

// BUFFER_SEARCH returns the first position, from pos on, of s in the buffer
// of size bytes, or -1. If ign is true the buffer is compared in upper case,
// so s must be in upper case.
func BUFFER_SEARCH(pt []iec.BYTE, size iec.INT, s iec.STRING, pos iec.INT, ign iec.BOOL) iec.INT {
	b := buf(pt)
	c := CHARS(s)
	lx := len(c)
	end := min(int(size)-lx, int(size)-1)
	for i := int(pos); i <= end; i++ {
		k := 0
		for ; k < lx; k++ {
			chx := b.at(i + k)
			if ign {
				chx = str.TO_UPPER(chx)
			}
			if c[k] != chx {
				break
			}
		}
		if k == lx {
			return iec.INT(i)
		}
	}
	return -1
}

// BUFFER_TO_STRING returns the bytes start to stop of the buffer of size
// bytes as a string, up to STRING_LENGTH characters and the first 0.
func BUFFER_TO_STRING(pt []iec.BYTE, size, start, stop iec.UINT) iec.STRING {
	if size == 0 {
		return ""
	}
	b := buf(pt)
	sta := min(start, size-1)
	stp := min(stop, size-1)
	if int(stp)-int(sta)+1 >= int(STRING_LENGTH) {
		stp = sta + iec.UINT(STRING_LENGTH) - 1
	}
	var out []iec.BYTE
	for i := int(sta); i <= int(stp); i++ {
		out = append(out, b.at(i))
	}
	return STR(out)
}
