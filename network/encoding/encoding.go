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

// Package encoding is the port of the OSCAT NETWORK encodings and string
// functions: Base64, HTML and URL encoding, URLs, lists of elements, file
// paths and IPv4 addresses.
package encoding

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// base64 is the alphabet of Base64.
const base64 iec.STRING = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// BASE64_ENCODE_STREAM encodes the SIZE1 bytes of BUF1 (48 at most) in
// Base64 into the SIZE2 bytes of BUF2, padded with '='.
type BASE64_ENCODE_STREAM struct {
	BUF1  *[48]iec.BYTE
	BUF2  *[64]iec.BYTE
	SIZE1 iec.INT
	SIZE2 iec.INT

	a, i, i2, i3, c iec.INT
	b               iec.WORD
}

// INIT resets the block.
func (e *BASE64_ENCODE_STREAM) INIT() {
	*e = BASE64_ENCODE_STREAM{BUF1: e.BUF1, BUF2: e.BUF2}
}

// Execute runs the block once.
func (e *BASE64_ENCODE_STREAM) Execute(now time.Time) {
	if e.BUF1 == nil || e.BUF2 == nil {
		return
	}
	if e.SIZE1 <= 0 {
		e.SIZE2 = 0
		return
	}
	e.c = e.SIZE1 - 1
	e.i2, e.a, e.b = 0, 0, 0
	for e.i = 0; e.i <= e.c && int(e.i) < len(e.BUF1); e.i++ {
		e.b = SHL(e.b, 8) | iec.WORD(e.BUF1[e.i])
		e.a += 8
		for e.a >= 6 {
			e.a -= 6
			e.BUF2[e.i2] = str.CODE(base64, iec.INT(SHR(e.b, e.a))+1)
			e.i2++
			e.b &^= SHL(iec.WORD(0xFFFF), e.a)
		}
	}
	if e.a > 0 {
		e.BUF2[e.i2] = str.CODE(base64, iec.INT(SHL(e.b, 6-e.a))+1)
		e.i2++
		// Pad with '=' to a multiple of 4.
		e.i3 = e.i2 + iec.INT(^iec.WORD(e.i2)&3)
		for e.i = e.i2; e.i <= e.i3 && int(e.i) < len(e.BUF2); e.i++ {
			e.BUF2[e.i] = 61
		}
		e.SIZE2 = e.i3 + 1
	} else {
		e.SIZE2 = e.i2
	}
}

// BASE64_DECODE_STREAM decodes the SIZE1 Base64 bytes of BUF1 (64 at most)
// into the SIZE2 bytes of BUF2. A '=' ends the data.
type BASE64_DECODE_STREAM struct {
	BUF1  *[64]iec.BYTE
	BUF2  *[48]iec.BYTE
	SIZE1 iec.INT
	SIZE2 iec.INT
}

// INIT resets the block.
func (d *BASE64_DECODE_STREAM) INIT() {
	*d = BASE64_DECODE_STREAM{BUF1: d.BUF1, BUF2: d.BUF2}
}

// Execute runs the block once.
func (d *BASE64_DECODE_STREAM) Execute(now time.Time) {
	if d.BUF1 == nil || d.BUF2 == nil {
		return
	}
	if d.SIZE1 <= 0 {
		d.SIZE2 = 0
		return
	}
	var a, off, i2 iec.INT
	var b iec.WORD
	for i1 := iec.INT(0); i1 <= d.SIZE1-1 && int(i1) < len(d.BUF1); i1++ {
		o := iec.INT(d.BUF1[i1])
		switch {
		case o == 61: // =
			d.SIZE2 = i2
			return
		case o > 96: // a-z
			off = 71
		case o > 64: // A-Z
			off = 65
		case o > 47: // 0-9
			off = -4
		case o == 43: // +
			off = -19
		case o == 47: // /
			off = -16
		}
		b = SHL(b, 6) | iec.WORD(iec.BYTE(o-off))
		a += 6
		if a >= 8 {
			a -= 8
			if int(i2) < len(d.BUF2) {
				d.BUF2[i2] = iec.BYTE(SHR(b, a))
			}
			i2++
			b &^= SHL(iec.WORD(0xFFFF), a)
		}
	}
	d.SIZE2 = i2
}

// base64Str is the conversion of BASE64_ENCODE_STR and BASE64_DECODE_STR: a
// rising edge of RUN converts STR1 into STR2, 48 characters a scan, and
// DONE is set at the end.
type base64Str struct {
	RUN  iec.BOOL
	DONE iec.BOOL

	runLast iec.BOOL
	mode    iec.INT
	size1   iec.INT
	size2   iec.INT
	end     iec.INT
	pos     iec.INT
}

// step runs a scan, converting a chunk with convert, which writes into
// buf2 the size2 bytes of the size1 bytes of buf1.
func (b *base64Str) step(str1, str2 *iec.STRING, buf1, buf2 []iec.BYTE, convert func(size1 iec.INT) iec.INT) {
	switch b.mode {
	case 0:
		if b.RUN && !b.runLast {
			b.DONE = false
			b.mode = 1
			b.end = LEN(*str1)
			b.pos = 0
			*str2 = ""
		}
	case 1:
		b.size1 = min(48, b.end-b.pos)
		if b.size1 > 0 {
			buffer.STRING_TO_BUFFER_(MID(*str1, b.size1, b.pos+1), 0, buf1, iec.UINT(len(buf1)))
			b.size2 = convert(b.size1)
			*str2 = CONCAT(*str2, buffer.BUFFER_TO_STRING(buf2, iec.UINT(len(buf2)), 0, iec.UINT(b.size2-1)))
		} else {
			b.mode = 0
			b.DONE = true
		}
		b.pos += b.size1
	}
	b.runLast = b.RUN
}

// BASE64_ENCODE_STR encodes the string STR1 (144 characters at most) in
// Base64 into STR2: a rising edge of RUN starts it, 48 characters a scan,
// and DONE is set at the end.
type BASE64_ENCODE_STR struct {
	base64Str
	STR1 *iec.STRING // STRING(144)
	STR2 *iec.STRING // STRING(192)

	stream BASE64_ENCODE_STREAM
	buf1   [48]iec.BYTE
	buf2   [64]iec.BYTE
}

// INIT resets the block.
func (e *BASE64_ENCODE_STR) INIT() { *e = BASE64_ENCODE_STR{STR1: e.STR1, STR2: e.STR2} }

// Execute runs the block once.
func (e *BASE64_ENCODE_STR) Execute(now time.Time) {
	if e.STR1 == nil || e.STR2 == nil {
		return
	}
	e.step(e.STR1, e.STR2, e.buf1[:], e.buf2[:], func(size1 iec.INT) iec.INT {
		e.stream.BUF1, e.stream.BUF2, e.stream.SIZE1 = &e.buf1, &e.buf2, size1
		e.stream.Execute(now)
		return e.stream.SIZE2
	})
}

// BASE64_DECODE_STR decodes the Base64 string STR1 (192 characters at
// most) into STR2: a rising edge of RUN starts it, 48 characters a scan,
// and DONE is set at the end.
type BASE64_DECODE_STR struct {
	base64Str
	STR1 *iec.STRING // STRING(192)
	STR2 *iec.STRING // STRING(144)

	stream BASE64_DECODE_STREAM
	buf1   [64]iec.BYTE
	buf2   [48]iec.BYTE
}

// INIT resets the block.
func (d *BASE64_DECODE_STR) INIT() { *d = BASE64_DECODE_STR{STR1: d.STR1, STR2: d.STR2} }

// Execute runs the block once.
func (d *BASE64_DECODE_STR) Execute(now time.Time) {
	if d.STR1 == nil || d.STR2 == nil {
		return
	}
	d.step(d.STR1, d.STR2, d.buf1[:], d.buf2[:], func(size1 iec.INT) iec.INT {
		d.stream.BUF1, d.stream.BUF2, d.stream.SIZE1 = &d.buf1, &d.buf2, size1
		d.stream.Execute(now)
		return d.stream.SIZE2
	})
}

// HTML_DECODE decodes the character references of IN: &#x..; in
// hexadecimal, &#..; in decimal and the named ones.
func HTML_DECODE(IN iec.STRING) iec.STRING {
	out := IN
	pos := FIND(out, "&")
	for pos > 0 {
		tmp := MID(out, 2, pos+1)
		switch {
		case tmp == "#x" || tmp == "#X":
			// hexadecimal
			tmp = MID(out, 10, pos+3)
			end := FIND(tmp, ";")
			code := str.CHR_TO_STRING(iec.BYTE(str.HEX_TO_DWORD(LEFT(tmp, end-1))))
			out = REPLACE(out, code, end+3, pos)
		case LEFT(tmp, 1) == "#":
			// decimal
			tmp = MID(out, 10, pos+2)
			end := FIND(tmp, ";")
			code := str.CHR_TO_STRING(iec.BYTE(STRING_TO_INT(LEFT(tmp, end-1))))
			out = REPLACE(out, code, end+2, pos)
		default:
			// named
			tmp = MID(out, 10, pos+1)
			end := FIND(tmp, ";")
			code := str.CHR_TO_STRING(str.CHARCODE(LEFT(tmp, end-1)))
			out = REPLACE(out, code, end+1, pos)
		}
		pos = str.FINDP(out, "&", pos+1)
	}
	return network.STRING_N(out, STRING_LENGTH)
}

// writer writes the characters of a string of 250 characters at most, as
// OSCAT does through a pointer to it.
type writer struct {
	out [257]iec.BYTE // ARRAY[1..256], at out[1..256]
}

// string returns the string written, up to the first 0.
func (w *writer) string() iec.STRING { return STR(w.out[1:]) }

// HTML_ENCODE encodes the characters ", &, < and > of IN as character
// references, and if M is true the characters above 127 by their names.
func HTML_ENCODE(IN iec.STRING, M iec.BOOL) iec.STRING {
	in := network.BYTES(IN, 256)
	var w writer
	put := func(pos *iec.INT, s string) {
		for i := 0; i < len(s); i++ {
			if i > 0 {
				*pos++
			}
			w.out[*pos] = iec.BYTE(s[i])
		}
	}
	posIn := iec.INT(1)
	stop := LEN(IN)
	var posOut iec.INT
	for posOut = 1; posOut <= 250; posOut++ {
		if posIn > stop {
			break
		}
		b := in[posIn-1]
		switch b {
		case 34: // "
			if posOut > STRING_LENGTH-6 {
				goto done
			}
			put(&posOut, "&quot;")
		case 38: // &
			if posOut > STRING_LENGTH-5 {
				goto done
			}
			put(&posOut, "&amp;")
		case 60: // <
			if posOut > STRING_LENGTH-4 {
				goto done
			}
			put(&posOut, "&lt;")
		case 62: // >
			if posOut > STRING_LENGTH-4 {
				goto done
			}
			put(&posOut, "&gt;")
		default:
			if M && b > 127 {
				tmp := str.CHARNAME(b)
				// Leave it out if it does not fit.
				if posOut+LEN(tmp)+2 <= 250 {
					w.out[posOut] = 38
					posOut++
					for i := iec.INT(1); i <= LEN(tmp); i++ {
						w.out[posOut] = str.CODE(tmp, i)
						posOut++
					}
					w.out[posOut] = 59
				}
			} else {
				w.out[posOut] = b
			}
		}
		posIn++
	}
done:
	// Terminate the string.
	w.out[posOut] = 0
	return w.string()
}

// URL_DECODE decodes the %XX sequences of IN.
func URL_DECODE(IN iec.STRING) iec.STRING {
	out := IN
	pos := FIND(out, "%")
	for pos > 0 {
		seq := MID(out, 2, pos+1)
		out = REPLACE(out, str.CHR_TO_STRING(iec.BYTE(str.HEX_TO_DWORD(seq))), 3, pos)
		pos = FIND(out, "%")
	}
	return out
}

// URL_ENCODE encodes the characters of IN that may not be in a URL as %XX.
func URL_ENCODE(IN iec.STRING) iec.STRING {
	in := network.BYTES(IN, 256)
	var w writer
	hex := func(tb iec.BYTE) iec.BYTE {
		if tb > 9 {
			return tb + 55
		}
		return tb + 48
	}
	posIn := iec.INT(1)
	stop := LEN(IN)
	var posOut iec.INT
	for posOut = 1; posOut <= 250; posOut++ {
		if posIn > stop {
			break
		}
		c := in[posIn-1]
		if IS_URLCHR(c) {
			w.out[posOut] = c
			posIn++
			continue
		}
		// Stop if the 3 characters do not fit.
		if posOut > 248 {
			break
		}
		w.out[posOut] = 37 // %
		posOut++
		w.out[posOut] = hex(SHR(c, 4))
		posOut++
		w.out[posOut] = hex(c & 0x0F)
		posIn++
	}
	w.out[posOut] = 0
	return w.string()
}

// IS_URLCHR reports whether the character IN may be in a URL as it is:
// letters, digits and ~ _ - .
func IS_URLCHR(IN iec.BYTE) iec.BOOL {
	return IN > 47 && IN < 58 || IN > 64 && IN < 91 || IN > 96 && IN < 123 || IN == 126 || IN == 95 || IN == 45 || IN == 46
}

// URL_TO_STRING writes the URL IN as a string.
func URL_TO_STRING(IN network.URL) iec.STRING {
	var s iec.STRING
	if IN.PROTOCOL != "" {
		s = CONCAT(IN.PROTOCOL, "://")
	}
	if IN.USER != "" {
		s = CONCAT(s, IN.USER)
		if IN.PASSWORD != "" {
			s = CONCAT(s, ":", IN.PASSWORD)
		}
		s = CONCAT(s, "@")
	}
	s = CONCAT(s, IN.DOMAIN)
	if IN.PORT > 0 {
		s = CONCAT(s, ":", network.WORD_TO_STRING(IN.PORT))
	}
	s = CONCAT(s, IN.PATH)
	if IN.QUERY != "" {
		s = CONCAT(s, "?", IN.QUERY)
	}
	if IN.ANCHOR != "" {
		s = CONCAT(s, "#", IN.ANCHOR)
	}
	return network.STRING_N(s, STRING_LENGTH)
}

// STRING_TO_URL splits the URL STR into its parts, with DEFAULT_PROTOCOL
// and DEFAULT_PATH where it has none:
// protocol://user:password@domain:port/path?query#anchor. OSCAT leaves the
// last character before the @ out of a user without a password.
func STRING_TO_URL(STR, DEFAULT_PROTOCOL, DEFAULT_PATH iec.STRING) network.URL {
	var u network.URL
	if STR == "" {
		return u
	}
	x := STR
	if pos := FIND(STR, "://"); pos > 0 {
		u.PROTOCOL = LEFT(STR, pos-1)
		x = RIGHT(STR, LEN(STR)-pos-2)
	} else {
		u.PROTOCOL = DEFAULT_PROTOCOL
	}

	// The user and password end with the last @.
	if pos := str.FINDB(x, "@"); pos > 0 {
		pos2 := FIND(x, ":") + 1
		if pos2 > 1 && pos2 <= pos {
			u.USER = LEFT(x, pos2-2)
			u.PASSWORD = MID(x, pos-pos2, pos2)
		} else {
			u.USER = LEFT(x, pos-2)
		}
		x = RIGHT(x, LEN(x)-pos)
	}

	if pos := FIND(x, "#"); pos > 0 {
		u.ANCHOR = RIGHT(x, LEN(x)-pos)
		x = LEFT(x, pos-1)
	}
	if pos := FIND(x, "?"); pos > 0 {
		u.QUERY = RIGHT(x, LEN(x)-pos)
		x = LEFT(x, pos-1)
	}
	if pos := FIND(x, "/"); pos > 0 {
		u.PATH = RIGHT(x, LEN(x)-pos+1)
		x = LEFT(x, pos-1)
	} else {
		u.PATH = DEFAULT_PATH
	}
	if pos := FIND(x, ":"); pos > 0 {
		u.PORT = network.STRING_TO_WORD(RIGHT(x, LEN(x)-pos))
		u.DOMAIN = LEFT(x, pos-1)
	} else {
		u.DOMAIN = x
	}
	u.PROTOCOL = network.STRING_N(u.PROTOCOL, 10)
	u.USER = network.STRING_N(u.USER, 32)
	u.PASSWORD = network.STRING_N(u.PASSWORD, 32)
	u.DOMAIN = network.STRING_N(u.DOMAIN, 80)
	u.PATH = network.STRING_N(u.PATH, 80)
	u.ANCHOR = network.STRING_N(u.ANCHOR, 40)
	return u
}

// ELEMENT_COUNT returns the number of elements of the list ELEMENT,
// separated by the character SEP: 0 for an empty list.
func ELEMENT_COUNT(SEP iec.BYTE, ELEMENT *iec.STRING) iec.INT {
	if ELEMENT == nil || *ELEMENT == "" {
		return 0
	}
	n := iec.INT(1)
	for _, c := range CHARS(*ELEMENT) {
		if c == SEP {
			n++
		}
	}
	return n
}

// ELEMENT_GET returns the element POS, from 0, of the list ELEMENT,
// separated by the character SEP.
func ELEMENT_GET(SEP iec.BYTE, POS iec.INT, ELEMENT *iec.STRING) iec.STRING {
	if ELEMENT == nil {
		return ""
	}
	const length = int(network.ELEMENT_LENGTH)
	pt := network.BYTES(*ELEMENT, length+2) // ARRAY[1..ELEMENT_LENGTH], at pt[i-1]
	var out []iec.BYTE
	i, o, cnt := 1, 1, iec.INT(0)
	// The POS-th separator.
	for cnt < POS && i < length && pt[i-1] > 0 {
		if pt[i-1] == SEP {
			cnt++
		}
		i++
	}
	// The element up to the next separator.
	var c iec.BYTE
	if i < length {
		c = pt[i-1]
	}
	for c != SEP && c > 0 && o < length && i < length {
		out = append(out, pt[i-1])
		o++
		i++
		c = pt[i-1]
	}
	return STR(out)
}

// FILE_PATH_SPLIT splits the path FILENAME into X: the drive, such as
// 'C:', the directory, with its / or \ at both ends, and the file name. It
// is false for an empty path.
func FILE_PATH_SPLIT(FILENAME iec.STRING, X *network.FILE_PATH_DATA) iec.BOOL {
	if X == nil {
		return false
	}
	X.DRIVE, X.DIRECTORY, X.FILENAME = "", "", ""
	c := LEN(FILENAME)
	if c == 0 {
		return false
	}
	var p1, p2 iec.INT
	for b := iec.INT(1); b <= c; b++ {
		switch MID(FILENAME, 1, b) {
		case ":":
			p1 = b
		case "/", "\\":
			p2 = b
		}
	}
	if p1 == 2 {
		X.DRIVE = LEFT(FILENAME, p1)
	}
	if p2 > 0 && p2 > p1 {
		X.DIRECTORY = MID(FILENAME, p2-p1, p1+1)
	}
	X.FILENAME = RIGHT(FILENAME, c-max(p1, p2))
	return true
}

// IP4_CHECK reports whether the IPv4 addresses NIP and LIP are in the same
// network of the subnet mask SM.
func IP4_CHECK(NIP, LIP, SM iec.DWORD) iec.BOOL { return NIP&SM == LIP&SM }

// IP4_DECODE returns the IPv4 address written as STR, such as
// '192.168.1.1', as a DWORD.
func IP4_DECODE(STR iec.STRING) iec.DWORD {
	s := network.STRING_N(STR, 15)
	var ip iec.DWORD
	for pos := FIND(s, "."); pos > 0; pos = FIND(s, ".") {
		ip = SHL(ip, 8) | network.STRING_TO_DWORD(LEFT(s, pos-1))
		s = DELETE(s, pos, 1)
	}
	return SHL(ip, 8) | network.STRING_TO_DWORD(s)
}

// IS_IP4 reports whether STR has the three dots of an IPv4 address.
func IS_IP4(STR iec.STRING) iec.BOOL { return str.COUNT_CHAR(STR, 46) == 3 }

// IP4_TO_STRING writes the IPv4 address IP4 as a string, such as
// '192.168.1.1'.
func IP4_TO_STRING(IP4 iec.DWORD) iec.STRING {
	s := CONCAT("...", network.BYTE_TO_STRING(iec.BYTE(IP4)))
	s = INSERT(s, network.BYTE_TO_STRING(iec.BYTE(SHR(IP4, 8))), 2)
	s = INSERT(s, network.BYTE_TO_STRING(iec.BYTE(SHR(IP4, 16))), 1)
	return CONCAT(network.BYTE_TO_STRING(iec.BYTE(SHR(IP4, 24))), s)
}
