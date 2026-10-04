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

// Package crypto is the port of the OSCAT NETWORK hashes and ciphers: MD5,
// SHA1, RC4 and the CRAM-MD5 authentication.
//
// The stream blocks MD5_STREAM, SHA1_STREAM and RC4_CRYPT_STREAM work on
// data 64 bytes at a time: the caller sets MODE to 1 and SIZE to the size
// of the data; the block sets MODE to 2 and SIZE to the size of the first
// chunk, which the caller copies into BUF, from the position POS, before
// each run; MODE is 3 when the data are done.
package crypto

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/basic/logic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/royaljelly/iec"
)

// MD5_AUX is a step of round N (1..4) of MD5.
func MD5_AUX(N iec.INT, A, B, C, D, X iec.DWORD, U iec.INT, T iec.DWORD) iec.DWORD {
	var w iec.DWORD
	switch N {
	case 1:
		w = B&C | ^B&D
	case 2:
		w = B&D | C&^D
	case 3:
		w = B ^ C ^ D
	case 4:
		w = C ^ (B | ^D)
	}
	return B + ROL(A+w+X+T, U)
}

// The steps of MD5: the word of the block, the shift and the constant.
var md5Steps = [4][16]struct {
	x int
	u iec.INT
	t iec.DWORD
}{
	{{0, 7, 0xD76AA478}, {1, 12, 0xE8C7B756}, {2, 17, 0x242070DB}, {3, 22, 0xC1BDCEEE},
		{4, 7, 0xF57C0FAF}, {5, 12, 0x4787C62A}, {6, 17, 0xA8304613}, {7, 22, 0xFD469501},
		{8, 7, 0x698098D8}, {9, 12, 0x8B44F7AF}, {10, 17, 0xFFFF5BB1}, {11, 22, 0x895CD7BE},
		{12, 7, 0x6B901122}, {13, 12, 0xFD987193}, {14, 17, 0xA679438E}, {15, 22, 0x49B40821}},
	{{1, 5, 0xF61E2562}, {6, 9, 0xC040B340}, {11, 14, 0x265E5A51}, {0, 20, 0xE9B6C7AA},
		{5, 5, 0xD62F105D}, {10, 9, 0x02441453}, {15, 14, 0xD8A1E681}, {4, 20, 0xE7D3FBC8},
		{9, 5, 0x21E1CDE6}, {14, 9, 0xC33707D6}, {3, 14, 0xF4D50D87}, {8, 20, 0x455A14ED},
		{13, 5, 0xA9E3E905}, {2, 9, 0xFCEFA3F8}, {7, 14, 0x676F02D9}, {12, 20, 0x8D2A4C8A}},
	{{5, 4, 0xFFFA3942}, {8, 11, 0x8771F681}, {11, 16, 0x6D9D6122}, {14, 23, 0xFDE5380C},
		{1, 4, 0xA4BEEA44}, {4, 11, 0x4BDECFA9}, {7, 16, 0xF6BB4B60}, {10, 23, 0xBEBFBC70},
		{13, 4, 0x289B7EC6}, {0, 11, 0xEAA127FA}, {3, 16, 0xD4EF3085}, {6, 23, 0x04881D05},
		{9, 4, 0xD9D4D039}, {12, 11, 0xE6DB99E5}, {15, 16, 0x1FA27CF8}, {2, 23, 0xC4AC5665}},
	{{0, 6, 0xF4292244}, {7, 10, 0x432AFF97}, {14, 15, 0xAB9423A7}, {5, 21, 0xFC93A039},
		{12, 6, 0x655B59C3}, {3, 10, 0x8F0CCC92}, {10, 15, 0xFFEFF47D}, {1, 21, 0x85845DD1},
		{8, 6, 0x6FA87E4F}, {15, 10, 0xFE2CE6E0}, {6, 15, 0xA3014314}, {13, 21, 0x4E0811A1},
		{4, 6, 0xF7537E82}, {11, 10, 0xBD3AF235}, {2, 15, 0x2AD7D2BB}, {9, 21, 0xEB86D391}},
}

// stream is the state the hash streams share: the blocks left, the size of
// the data and the padding.
type stream struct {
	end   iec.UDINT
	block iec.UDINT
	pad1  iec.BOOL
}

// start sets up a hash of size bytes and returns the size of the first
// chunk.
func (s *stream) start(size iec.UDINT) iec.UDINT {
	s.block = SHR(size, 6) + iec.UDINT(BOOL_TO_INT(iec.BYTE(size)&0x3F > 55))
	s.pad1 = false
	s.end = size
	return min(64, s.end)
}

// pad pads the chunk of size bytes in buf, and moves pos and size to the
// next chunk; it reports whether the chunk is the last, which holds the
// length.
func (s *stream) pad(buf *[64]iec.BYTE, size, pos *iec.UDINT) bool {
	for n := *size; n <= 63; n++ {
		buf[n] = 0
	}
	if *size < 64 && !s.pad1 {
		buf[*size] = 0x80
		s.pad1 = true
	}
	*pos += *size
	*size = min(64, s.end-*pos)
	if s.block == 0 {
		return true
	}
	s.block--
	return false
}

// MD5_STREAM computes the MD5 hash MD5 of data given 64 bytes at a time in
// BUF; see the package documentation.
type MD5_STREAM struct {
	MODE *iec.INT
	BUF  *[64]iec.BYTE
	MD5  *[16]iec.BYTE
	SIZE *iec.UDINT
	POS  iec.UDINT

	stream
	hash [4]iec.DWORD
	x    [16]iec.DWORD
}

// INIT resets the block.
func (m *MD5_STREAM) INIT() { *m = MD5_STREAM{MODE: m.MODE, BUF: m.BUF, MD5: m.MD5, SIZE: m.SIZE} }

// Execute runs the block once.
func (m *MD5_STREAM) Execute(now time.Time) {
	if m.MODE == nil || m.BUF == nil || m.MD5 == nil || m.SIZE == nil {
		return
	}
	switch *m.MODE {
	case 1:
		m.hash = [4]iec.DWORD{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476}
		m.POS = 0
		*m.SIZE = m.start(*m.SIZE)
		*m.MODE = 2
	case 2:
		last := m.pad(m.BUF, m.SIZE, &m.POS)
		buf := m.BUF
		for n1 := 0; n1 < 16; n1++ {
			n := 4 * n1
			// little endian
			m.x[n1] = logic.DWORD_OF_BYTE(buf[n+3], buf[n+2], buf[n+1], buf[n])
		}
		if last {
			// The length of the data in bits, little endian.
			m.x[14] = SHL(iec.DWORD(m.end), 3)
			*m.MODE = 3
		}
		v := m.hash
		for round, steps := range md5Steps {
			for i, s := range steps {
				// a, d, c, b in turn.
				j := (4 - i%4) % 4
				v[j] = MD5_AUX(iec.INT(round+1), v[j], v[(j+1)%4], v[(j+2)%4], v[(j+3)%4], m.x[s.x], s.u, s.t)
			}
		}
		for i := range m.hash {
			m.hash[i] += v[i]
		}
	}
	if *m.MODE == 3 {
		z := 0
		for n := range m.hash {
			for n1 := 0; n1 < 4; n1++ {
				m.MD5[z] = iec.BYTE(m.hash[n])
				m.hash[n] = ROR(m.hash[n], 8)
				z++
			}
		}
	}
}

// MD5_STR computes the MD5 hash MD5 of the string STR: a rising edge of
// RUN starts it, 64 characters a scan, and DONE is set at the end.
type MD5_STR struct {
	RUN  iec.BOOL
	DONE iec.BOOL
	STR  *iec.STRING // STRING(STRING_LENGTH)
	MD5  *[16]iec.BYTE

	runLast iec.BOOL
	stream  MD5_STREAM
	buf     [64]iec.BYTE
	mode    iec.INT
	size    iec.UDINT
	pos     iec.UDINT
}

// INIT resets the block.
func (m *MD5_STR) INIT() { *m = MD5_STR{STR: m.STR, MD5: m.MD5} }

// Execute runs the block once.
func (m *MD5_STR) Execute(now time.Time) {
	if m.STR == nil || m.MD5 == nil {
		return
	}
	switch m.mode {
	case 0: // wait for a start
		if m.RUN && !m.runLast {
			m.DONE = false
			m.mode = 1
			m.size = iec.UDINT(LEN(*m.STR))
		}
	case 2: // copy the data
		if m.size > 0 {
			buffer.STRING_TO_BUFFER_(MID(*m.STR, iec.INT(m.size), iec.INT(m.pos)+1), 0, m.buf[:], 64)
		}
	case 3:
		m.DONE = true
		m.mode = 0
	}
	if m.mode > 0 {
		m.stream.SIZE, m.stream.MODE, m.stream.BUF, m.stream.MD5 = &m.size, &m.mode, &m.buf, m.MD5
		m.stream.Execute(now)
		m.pos = m.stream.POS
	}
	m.runLast = m.RUN
}

// hexString writes the bytes h in hexadecimal, lower case.
func hexString(h []iec.BYTE) iec.STRING {
	out := make([]iec.BYTE, 0, 2*len(h))
	for _, b := range h {
		for _, tmp := range []iec.BYTE{SHR(b, 4), b & 0x0F} {
			out = append(out, tmp+SEL[iec.BYTE](tmp <= 9, 87, 48))
		}
	}
	return STR(out)
}

// MD5_TO_STRH writes the MD5 hash MD5 in hexadecimal, lower case.
func MD5_TO_STRH(MD5 [16]iec.BYTE) iec.STRING { return hexString(MD5[:]) }

// SHA1_STREAM computes the SHA1 hash SHA1 of data given 64 bytes at a time
// in BUF; see the package documentation.
type SHA1_STREAM struct {
	MODE *iec.INT
	BUF  *[64]iec.BYTE
	SHA1 *[20]iec.BYTE
	SIZE *iec.UDINT
	POS  iec.UDINT

	stream
	hash [5]iec.DWORD
	w    [80]iec.DWORD
}

// INIT resets the block.
func (s *SHA1_STREAM) INIT() {
	*s = SHA1_STREAM{MODE: s.MODE, BUF: s.BUF, SHA1: s.SHA1, SIZE: s.SIZE}
}

// Execute runs the block once.
func (s *SHA1_STREAM) Execute(now time.Time) {
	if s.MODE == nil || s.BUF == nil || s.SHA1 == nil || s.SIZE == nil {
		return
	}
	switch *s.MODE {
	case 1:
		s.hash = [5]iec.DWORD{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476, 0xC3D2E1F0}
		for n := 16; n < 80; n++ {
			s.w[n] = 0
		}
		s.POS = 0
		*s.SIZE = s.start(*s.SIZE)
		*s.MODE = 2
	case 2:
		last := s.pad(s.BUF, s.SIZE, &s.POS)
		buf := s.BUF
		for n1 := 0; n1 < 16; n1++ {
			n := 4 * n1
			s.w[n1] = logic.DWORD_OF_BYTE(buf[n], buf[n+1], buf[n+2], buf[n+3])
		}
		if last {
			// The length of the data in bits.
			s.w[15] = SHL(iec.DWORD(s.end), 3)
			*s.MODE = 3
		}
		for n := 16; n < 80; n++ {
			s.w[n] = ROL(s.w[n-3]^s.w[n-8]^s.w[n-14]^s.w[n-16], 1)
		}
		a, b, c, d, e := s.hash[0], s.hash[1], s.hash[2], s.hash[3], s.hash[4]
		for n := 0; n < 80; n++ {
			var f, k iec.DWORD
			switch {
			case n <= 19:
				f, k = b&c|^b&d, 0x5A827999
			case n <= 39:
				f, k = b^c^d, 0x6ED9EBA1
			case n <= 59:
				f, k = b&c|b&d|c&d, 0x8F1BBCDC
			default:
				f, k = b^c^d, 0xCA62C1D6
			}
			x := ROL(a, 5) + f + e + k + s.w[n]
			e, d, c, b, a = d, c, ROL(b, 30), a, x
		}
		s.hash[0] += a
		s.hash[1] += b
		s.hash[2] += c
		s.hash[3] += d
		s.hash[4] += e
	}
	if *s.MODE == 3 {
		z := 0
		for n := range s.hash {
			for n1 := 0; n1 < 4; n1++ {
				s.hash[n] = ROL(s.hash[n], 8)
				s.SHA1[z] = iec.BYTE(s.hash[n])
				z++
			}
		}
	}
}

// SHA1_STR computes the SHA1 hash SHA1 of the string STR: a rising edge of
// RUN starts it, 64 characters a scan, and DONE is set at the end.
type SHA1_STR struct {
	RUN  iec.BOOL
	DONE iec.BOOL
	STR  *iec.STRING // STRING(STRING_LENGTH)
	SHA1 *[20]iec.BYTE

	runLast iec.BOOL
	stream  SHA1_STREAM
	buf     [64]iec.BYTE
	mode    iec.INT
	size    iec.UDINT
	pos     iec.UDINT
}

// INIT resets the block.
func (s *SHA1_STR) INIT() { *s = SHA1_STR{STR: s.STR, SHA1: s.SHA1} }

// Execute runs the block once.
func (s *SHA1_STR) Execute(now time.Time) {
	if s.STR == nil || s.SHA1 == nil {
		return
	}
	switch s.mode {
	case 0: // wait for a start
		if s.RUN && !s.runLast {
			s.DONE = false
			s.mode = 1
			s.size = iec.UDINT(LEN(*s.STR))
		}
	case 2: // copy the data
		if s.size > 0 {
			buffer.STRING_TO_BUFFER_(MID(*s.STR, iec.INT(s.size), iec.INT(s.pos)+1), 0, s.buf[:], 64)
		}
	case 3:
		s.DONE = true
		s.mode = 0
	}
	if s.mode > 0 {
		s.stream.SIZE, s.stream.MODE, s.stream.BUF, s.stream.SHA1 = &s.size, &s.mode, &s.buf, s.SHA1
		s.stream.Execute(now)
		s.pos = s.stream.POS
	}
	s.runLast = s.RUN
}

// SHA1_TO_STRH writes the SHA1 hash SHA1 in hexadecimal, lower case.
func SHA1_TO_STRH(SHA1 [20]iec.BYTE) iec.STRING { return hexString(SHA1[:]) }

// RC4_CRYPT_STREAM encrypts or decrypts, in place, data given 64 bytes at a
// time in BUF with RC4 and the key KEY; see the package documentation. An
// empty key or no data end it at once.
type RC4_CRYPT_STREAM struct {
	MODE *iec.INT
	KEY  *iec.STRING // STRING(40)
	BUF  *[64]iec.BYTE
	SIZE *iec.UDINT
	POS  iec.UDINT

	sbox [256]iec.BYTE
	skey [256]iec.BYTE
	d, e iec.USINT
	end  iec.UDINT
}

// INIT resets the block.
func (r *RC4_CRYPT_STREAM) INIT() {
	*r = RC4_CRYPT_STREAM{MODE: r.MODE, KEY: r.KEY, BUF: r.BUF, SIZE: r.SIZE}
}

// Execute runs the block once.
func (r *RC4_CRYPT_STREAM) Execute(now time.Time) {
	if r.MODE == nil || r.KEY == nil || r.BUF == nil || r.SIZE == nil {
		return
	}
	switch *r.MODE {
	case 1:
		key := network.STRING_N(*r.KEY, 40)
		b := int(LEN(key)) - 1
		if b < 0 || *r.SIZE < 1 {
			*r.MODE = 3
			return
		}
		for a := 0; a <= b; a++ {
			r.skey[a] = str.CODE(key, iec.INT(a+1))
			r.sbox[a] = iec.BYTE(a)
		}
		c := 0
		for a := b + 1; a <= 255; a++ {
			r.skey[a] = r.skey[c]
			r.sbox[a] = iec.BYTE(a)
			c++
			if c > b {
				c = 0
			}
		}
		r.d = 0
		for a := 0; a <= 255; a++ {
			r.d += iec.USINT(r.sbox[a]) + iec.USINT(r.skey[a])
			r.sbox[a], r.sbox[r.d] = r.sbox[r.d], r.sbox[a]
		}
		r.POS = 0
		r.end = *r.SIZE
		*r.SIZE = min(64, r.end)
		r.d, r.e = 0, 0
		*r.MODE = 2
	case 2:
		for a := 0; a < int(*r.SIZE) && a < 64; a++ {
			r.d++
			r.e += iec.USINT(r.sbox[r.d])
			r.sbox[r.d], r.sbox[r.e] = r.sbox[r.e], r.sbox[r.d]
			r.BUF[a] ^= r.sbox[r.sbox[r.d]+r.sbox[r.e]]
		}
		r.POS += *r.SIZE
		*r.SIZE = min(64, r.end-r.POS)
		if *r.SIZE == 0 {
			*r.MODE = 3
		}
	}
}

// MD5_CRAM_AUTH computes the answer AUTH_KEY to a CRAM-MD5 challenge, the
// Base64 time stamp B64_TS of a server: Base64 of USERNAME, a space and
// the HMAC-MD5 of the time stamp with the key PASSWORD (64 characters at
// most). RUN starts it and is cleared at the end.
type MD5_CRAM_AUTH struct {
	RUN      *iec.BOOL
	USERNAME *iec.STRING // STRING(64)
	PASSWORD *iec.STRING // STRING(64)
	B64_TS   *iec.STRING // STRING(64)
	AUTH_KEY *iec.STRING // STRING(192)

	decode    encoding.BASE64_DECODE_STR
	encode    encoding.BASE64_ENCODE_STR
	md5Stream MD5_STREAM
	step      iec.INT
	md5Mode   iec.INT
	buf       [64]iec.BYTE
	md5       [16]iec.BYTE
	md5First  [16]iec.BYTE
	md5Str1   iec.STRING // STRING(64)
	md5Str2   iec.STRING // STRING(64)
	size      iec.UDINT
	pos       iec.UDINT
	b64dRun   iec.BOOL
	b64dDone  iec.BOOL
	b64eRun   iec.BOOL
	b64eDone  iec.BOOL
	str192    iec.STRING // STRING(192)
	str144    iec.STRING // STRING(144)
	xpad      iec.BYTE
}

// INIT resets the block.
func (m *MD5_CRAM_AUTH) INIT() {
	*m = MD5_CRAM_AUTH{RUN: m.RUN, USERNAME: m.USERNAME, PASSWORD: m.PASSWORD, B64_TS: m.B64_TS, AUTH_KEY: m.AUTH_KEY}
}

// Execute runs the block once.
func (m *MD5_CRAM_AUTH) Execute(now time.Time) {
	if m.RUN == nil || m.USERNAME == nil || m.PASSWORD == nil || m.B64_TS == nil || m.AUTH_KEY == nil {
		return
	}
	switch m.step {
	case 0:
		if *m.RUN {
			m.step = 5
			m.str192 = *m.B64_TS
			m.b64dRun = true
		}
	case 5: // the inner pad
		if m.b64dDone {
			m.b64dRun = false
			m.xpad = 0x36
			m.md5Str1 = network.STRING_N(*m.PASSWORD, 64)
			m.md5Str2 = network.STRING_N(m.str144, 64)
			m.size = 64 + iec.UDINT(LEN(m.md5Str2))
			m.md5Mode = 1
			m.step = 10
		}
	case 10: // the outer pad
		if m.md5Mode == 3 {
			m.md5First = m.md5
			m.xpad = 0x5C
			m.size = 64 + 16
			m.md5Mode = 1
			m.step = 20
		}
	case 20:
		if m.md5Mode == 3 {
			m.md5Mode = 0
			m.str144 = network.STRING_N(CONCAT(*m.USERNAME, " ", MD5_TO_STRH(m.md5)), 144)
			m.b64eRun = true
			m.step = 30
		}
	case 30:
		if m.b64eDone {
			*m.AUTH_KEY = m.str192
			m.b64eRun = false
			*m.RUN = false
			m.step = 0
		}
	}

	switch m.md5Mode {
	case 1:
		// The key, XORed with the pad.
		key := CHARS(m.md5Str1)
		for n1 := range m.buf {
			if n1 < len(key) {
				m.buf[n1] = key[n1] ^ m.xpad
			} else {
				m.buf[n1] = m.xpad
			}
		}
	case 2:
		if m.size > 0 && m.pos >= 64 {
			if m.step == 20 {
				copy(m.buf[:16], m.md5First[:])
			} else {
				copy(m.buf[:], CHARS(m.md5Str2))
			}
		}
	}
	if m.md5Mode > 0 {
		m.md5Stream.SIZE, m.md5Stream.MODE, m.md5Stream.BUF, m.md5Stream.MD5 = &m.size, &m.md5Mode, &m.buf, &m.md5
		m.md5Stream.Execute(now)
		m.pos = m.md5Stream.POS
	}

	m.decode.RUN, m.decode.STR1, m.decode.STR2 = m.b64dRun, &m.str192, &m.str144
	m.decode.Execute(now)
	m.b64dDone = m.decode.DONE
	m.encode.RUN, m.encode.STR1, m.encode.STR2 = m.b64eRun, &m.str144, &m.str192
	m.encode.Execute(now)
	m.b64eDone = m.encode.DONE
}
