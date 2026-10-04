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

package main

// This file reads a CODESYS 2.3 library (.lib), a binary file, into the
// lines of a CODESYS export, which clean then cleans as any export.
//
// The library keeps the text of each POU as length prefixed strings: a
// 4-byte little endian length and the declaration, then 0x12, a length and
// the body. A TYPE or global variable list has a declaration alone. The
// library may also come as text an editor made of it: read as Windows-1252
// and saved as UTF-8, with its NUL bytes written as spaces and a lone CR or
// LF written as CR LF. Its bytes are recovered, but a 0x20 in a length may
// then be a NUL and a CR LF a lone CR or LF, so each reading of a length is
// tried, and the one that ends where a record ends is taken.

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// isLib reports whether data is a CODESYS library rather than an export.
func isLib(data []byte) bool { return bytes.HasPrefix(data, []byte("CoDeSys")) }

// cp1252 maps the characters of Windows-1252 that are not ISO 8859-1 to
// their bytes.
var cp1252 = map[rune]byte{
	'€': 0x80, '‚': 0x82, 'ƒ': 0x83, '„': 0x84, '…': 0x85, '†': 0x86, '‡': 0x87, 'ˆ': 0x88,
	'‰': 0x89, 'Š': 0x8A, '‹': 0x8B, 'Œ': 0x8C, 'Ž': 0x8E, '‘': 0x91, '’': 0x92, '“': 0x93,
	'”': 0x94, '•': 0x95, '–': 0x96, '—': 0x97, '˜': 0x98, '™': 0x99, 'š': 0x9A, '›': 0x9B,
	'œ': 0x9C, 'ž': 0x9E, 'Ÿ': 0x9F,
}

// libBytes returns the bytes of the library: data itself, or the bytes an
// editor read as Windows-1252 and saved as UTF-8.
func libBytes(data []byte) []byte {
	if !utf8.Valid(data) {
		return data
	}
	raw := make([]byte, 0, len(data))
	for len(data) > 0 {
		r, n := utf8.DecodeRune(data)
		data = data[n:]
		if b, ok := cp1252[r]; ok {
			raw = append(raw, b)
		} else {
			raw = append(raw, byte(r))
		}
	}
	return raw
}

// latin1 returns the text of bytes in Windows-1252 as UTF-8.
func latin1(b []byte) string {
	var s strings.Builder
	for _, c := range b {
		r := rune(c)
		for k, v := range cp1252 {
			if v == c {
				r = k
				break
			}
		}
		s.WriteRune(r)
	}
	return s.String()
}

// control reports whether c ends a text: a control character other than a
// tab or an end of line.
func control(c byte) bool { return c < 0x20 && c != '\t' && c != '\r' && c != '\n' || c == 0x7f }

// lengthField is a reading of a length: its value and where it ends.
type lengthField struct{ val, end int }

// readLength returns the readings of the length at p. A NUL may be written
// as a space, and a lone CR or LF as CR LF.
func readLength(raw []byte, p int) []lengthField {
	var out []lengthField
	var walk func(q, i, val int)
	walk = func(q, i, val int) {
		if i == 4 {
			out = append(out, lengthField{val, q})
			return
		}
		if q < 0 || q >= len(raw) {
			return
		}
		c := raw[q]
		if i >= 2 {
			// The texts are shorter than 64K.
			if c == 0 || c == 0x20 {
				walk(q+1, i+1, val)
			}
			return
		}
		walk(q+1, i+1, val|int(c)<<(8*i))
		if c == 0x20 {
			walk(q+1, i+1, val)
		}
		if c == '\r' && q+1 < len(raw) && raw[q+1] == '\n' {
			walk(q+2, i+1, val|'\n'<<(8*i))
			walk(q+2, i+1, val|'\r'<<(8*i))
		}
	}
	walk(p, 0, 0)
	return out
}

// isText reports whether raw[a:b] is text.
func isText(raw []byte, a, b int) bool {
	for _, c := range raw[a:b] {
		if control(c) {
			return false
		}
	}
	return true
}

// libRecord is a POU, TYPE or global variable list of a library.
type libRecord struct {
	kind, name string
	decl, body string
}

var reLibHead = regexp.MustCompile(`^(FUNCTION_BLOCK|FUNCTION|PROGRAM|TYPE|VAR_GLOBAL)\b\s*([A-Za-z_0-9]*)`)

// libRecords returns the records of the library raw, in their order.
func libRecords(raw []byte) []libRecord {
	var out []libRecord
	for s := 4; s < len(raw); s++ {
		m := reLibHead.FindSubmatch(raw[s:min(s+80, len(raw))])
		if m == nil {
			continue
		}
		kind := string(m[1])
		var found []libRecord
		var ends []int
		for _, q := range []int{s - 4, s - 5} {
			for _, f := range readLength(raw, q) {
				e := s + f.val
				if f.end != s || f.val < len(m[0]) || e+1 > len(raw) || !isText(raw, s, e) {
					continue
				}
				decl := string(raw[s:e])
				if kind == "TYPE" || kind == "VAR_GLOBAL" {
					end := map[string]string{"TYPE": "END_TYPE", "VAR_GLOBAL": "END_VAR"}[kind]
					if (control(raw[e]) || raw[e] == 0 || raw[e] == 0x20) && strings.HasSuffix(strings.TrimSpace(decl), end) {
						found = append(found, libRecord{kind: kind, name: string(m[2]), decl: decl})
						ends = append(ends, e)
					}
					continue
				}
				if raw[e] != 0x12 {
					continue
				}
				for _, g := range readLength(raw, e+1) {
					b := g.end + g.val
					if g.val > 0 && b <= len(raw) && isText(raw, g.end, b) && (b == len(raw) || control(raw[b]) || raw[b] == 0) {
						found = append(found, libRecord{kind, string(m[2]), decl, string(raw[g.end:b])})
						ends = append(ends, b)
					}
				}
			}
		}
		if len(found) == 1 {
			out = append(out, found[0])
			s = ends[0] - 1
		}
	}
	return out
}

// libExport returns the lines of a CODESYS export of the library data: its
// functions and function blocks, then its types and global variables. Its
// programs, the library's demos, are left out, and they are returned.
func libExport(data []byte) (lines, programs []string) {
	var pous, types, globals []string
	add := func(to *[]string, text string) {
		text = strings.ReplaceAll(latin1([]byte(text)), "\r\n", "\n")
		*to = append(*to, strings.Split(strings.TrimRight(text, "\n "), "\n")...)
		*to = append(*to, "")
	}
	for _, r := range libRecords(libBytes(data)) {
		switch r.kind {
		case "PROGRAM":
			programs = append(programs, r.name)
		case "TYPE":
			add(&types, r.decl)
		case "VAR_GLOBAL":
			add(&globals, r.decl)
		default:
			add(&pous, fmt.Sprintf("%s\n(* @END_DECLARATION := '0' *)\n%s\nEND_%s", strings.TrimRight(r.decl, "\r\n "), r.body, r.kind))
		}
	}
	lines = append(append(pous, types...), globals...)
	return lines, programs
}
