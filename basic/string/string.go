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

// Package string is the port of the OSCAT BASIC string functions. A
// character's code is its ISO 8859-1 code; see the package basic.
package string

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/math"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// mapChars returns str with each character of a single byte code replaced
// by f of its code. Other characters are kept.
func mapChars(str iec.STRING, f func(iec.BYTE) iec.BYTE) iec.STRING {
	r := []rune(string(str))
	for i, c := range r {
		if c <= 255 {
			r[i] = rune(f(iec.BYTE(c)))
		}
	}
	return iec.STRING(r)
}

// all reports whether str is not empty and each of its characters is f.
func all(str iec.STRING, f func(iec.BYTE) iec.BOOL) iec.BOOL {
	c := CHARS(str)
	for _, x := range c {
		if !f(x) {
			return false
		}
	}
	return len(c) > 0
}

func hexChar(v iec.BYTE) iec.BYTE {
	if v <= 9 {
		return v + 48
	}
	return v + 55
}

// BIN_TO_BYTE converts a binary string to a byte. Characters other than 0
// and 1 are ignored.
func BIN_TO_BYTE(bin iec.STRING) iec.BYTE {
	var out iec.BYTE
	for _, x := range CHARS(bin) {
		switch x {
		case '0':
			out <<= 1
		case '1':
			out = out<<1 | 1
		}
	}
	return out
}

// BIN_TO_DWORD converts a binary string to a DWORD. Characters other than 0
// and 1 are ignored.
func BIN_TO_DWORD(bin iec.STRING) iec.DWORD {
	var out iec.DWORD
	for _, x := range CHARS(bin) {
		switch x {
		case '0':
			out <<= 1
		case '1':
			out = out<<1 | 1
		}
	}
	return out
}

// BYTE_TO_STRB converts a byte to a string of 8 bits, the highest first.
func BYTE_TO_STRB(in iec.BYTE) iec.STRING {
	b := make([]iec.BYTE, 8)
	for i := range b {
		b[i] = in>>(7-i)&1 + 48
	}
	return STR(b)
}

// BYTE_TO_STRH converts a byte to a string of 2 hexadecimal digits.
func BYTE_TO_STRH(in iec.BYTE) iec.STRING {
	return STR([]iec.BYTE{hexChar(in >> 4), hexChar(in & 0x0F)})
}

// CAPITALIZE returns str with each first letter after a blank, and at the
// start, in upper case.
func CAPITALIZE(str iec.STRING) iec.STRING {
	first := true
	return mapChars(str, func(c iec.BYTE) iec.BYTE {
		if first {
			c = TO_UPPER(c)
		}
		first = c == 32
		return c
	})
}

// CHARCODE returns the code of the character whose HTML name is str, such as
// "euro". A string of one character gives that character's code. It is 0
// for an unknown name.
func CHARCODE(str iec.STRING) iec.BYTE {
	if LEN(str) == 1 {
		return CODE(str, 1)
	}
	if str == "" {
		return 0
	}
	search := "&" + str + ";"
	for _, names := range SETUP.CHARNAMES {
		if pos := FIND(names, search); pos > 0 {
			return CODE(MID(names, 1, pos-1), 1)
		}
	}
	return 0
}

// CHARNAME returns the HTML name of the character with the code c, such as
// "euro", or the character itself if it has no name.
func CHARNAME(c iec.BYTE) iec.STRING {
	if c == 0 {
		return ""
	}
	search := ";" + CHR_TO_STRING(c) + "&"
	for _, names := range SETUP.CHARNAMES {
		if pos := FIND(names, search); pos > 0 {
			name := MID(names, 10, pos+3)
			return LEFT(name, FIND(name, ";")-1)
		}
	}
	return CHR_TO_STRING(c)
}

// CHR_TO_STRING returns the string of the one character with the code c.
func CHR_TO_STRING(c iec.BYTE) iec.STRING {
	return iec.STRING(rune(c))
}

// CLEAN returns in without the characters that are not in cx.
func CLEAN(in, cx iec.STRING) iec.STRING {
	var out []rune
	for _, c := range string(in) {
		if FIND(cx, iec.STRING(c)) > 0 {
			out = append(out, c)
		}
	}
	return iec.STRING(out)
}

// CODE returns the code of the character at position pos of str, or 0 if
// there is none.
func CODE(str iec.STRING, pos iec.INT) iec.BYTE {
	c := CHARS(str)
	if pos < 1 || int(pos) > len(c) {
		return 0
	}
	return c[pos-1]
}

// COUNT_CHAR counts the characters with the code chr in str.
func COUNT_CHAR(str iec.STRING, chr iec.BYTE) iec.INT {
	var n iec.INT
	for _, c := range CHARS(str) {
		if c == chr {
			n++
		}
	}
	return n
}

// COUNT_SUBSTRING counts the occurrences of search in str.
func COUNT_SUBSTRING(search, str iec.STRING) iec.INT {
	var n iec.INT
	size := LEN(search)
	for {
		pos := FIND(str, search)
		if pos == 0 {
			return n
		}
		str = REPLACE(str, "", size, pos)
		n++
	}
}

// DEC_TO_BYTE converts a decimal string to a byte. Characters other than
// digits are ignored.
func DEC_TO_BYTE(dec iec.STRING) iec.BYTE {
	var out iec.BYTE
	for _, x := range CHARS(dec) {
		if x > 47 && x < 58 {
			out = out*10 + x - 48
		}
	}
	return out
}

// DEC_TO_DWORD converts a decimal string to a DWORD. Characters other than
// digits are ignored.
func DEC_TO_DWORD(dec iec.STRING) iec.DWORD {
	var out iec.DWORD
	for _, x := range CHARS(dec) {
		if x > 47 && x < 58 {
			out = out*10 + iec.DWORD(x) - 48
		}
	}
	return out
}

// DEC_TO_INT converts a decimal string to an INT. A minus before the first
// digit makes it negative; other characters are ignored.
func DEC_TO_INT(dec iec.STRING) iec.INT {
	var out iec.INT
	sign := false
	for _, x := range CHARS(dec) {
		if x > 47 && x < 58 {
			out = out*10 + iec.INT(x) - 48
		} else if x == 45 && out == 0 {
			sign = true
		}
	}
	if sign {
		return -out
	}
	return out
}

// DEL_CHARS returns in without the characters that are in cx.
func DEL_CHARS(in, cx iec.STRING) iec.STRING {
	var out []rune
	for _, c := range string(in) {
		if FIND(cx, iec.STRING(c)) == 0 {
			out = append(out, c)
		}
	}
	return iec.STRING(out)
}

// DT_TO_STRF formats a date and time with the milliseconds ms. Each #X in
// fmt is replaced:
//
//	#A year, 4 digits          #M hour 0..23, 1 or 2 digits
//	#B year, 2 digits          #N hour 0..23, 2 digits
//	#C month, 1 or 2 digits    #O hour 1..12, 1 or 2 digits
//	#D month, 2 digits         #P hour 1..12, 2 digits
//	#E month, 3 letters        #Q minute, 1 or 2 digits
//	#F month, full name        #R minute, 2 digits
//	#G day, 1 or 2 digits      #S second, 1 or 2 digits
//	#H day, 2 digits           #T second, 2 digits
//	#I weekday 1..7            #U milliseconds, 1 to 3 digits
//	#J weekday, 2 letters      #V milliseconds, 3 digits
//	#K weekday, full name      #W day, 2 characters, blank first
//	#L AM or PM                #X month, 2 characters, blank first
//
// Names are in the language lang, or the default language if lang < 1.
func DT_TO_STRF(dti iec.DT, ms iec.INT, fmt iec.STRING, lang iec.INT) iec.STRING {
	ly := LANGUAGE.DEFAULT
	if lang >= 1 {
		ly = min(LANGUAGE.LMAX, lang)
	}
	dx := DT_TO_DATE(dti)
	tod := DT_TO_TOD(dti)
	pad := func(s, fill iec.STRING) iec.STRING {
		if LEN(s) < 2 {
			return fill + s
		}
		return s
	}
	hour12 := func() iec.INT {
		h := td.HOUR(tod) % 12
		if h == 0 {
			h = 12
		}
		return h
	}
	out := fmt
	pos := FIND(out, "#")
	for pos > 0 {
		var fs iec.STRING
		switch CODE(out, pos+1) {
		case 'A':
			fs = INT_TO_STRING(td.YEAR_OF_DATE(dx))
		case 'B':
			fs = RIGHT(INT_TO_STRING(td.YEAR_OF_DATE(dx)), 2)
		case 'C':
			fs = INT_TO_STRING(td.MONTH_OF_DATE(dx))
		case 'D':
			fs = pad(INT_TO_STRING(td.MONTH_OF_DATE(dx)), "0")
		case 'E':
			fs = MONTH_TO_STRING(td.MONTH_OF_DATE(dx), ly, 3)
		case 'F':
			fs = MONTH_TO_STRING(td.MONTH_OF_DATE(dx), ly, 0)
		case 'G':
			fs = INT_TO_STRING(td.DAY_OF_MONTH(dx))
		case 'H':
			fs = pad(INT_TO_STRING(td.DAY_OF_MONTH(dx)), "0")
		case 'I':
			fs = INT_TO_STRING(td.DAY_OF_WEEK(dx))
		case 'J':
			fs = WEEKDAY_TO_STRING(td.DAY_OF_WEEK(dx), ly, 2)
		case 'K':
			fs = WEEKDAY_TO_STRING(td.DAY_OF_WEEK(dx), ly, 0)
		case 'L':
			fs = SEL[iec.STRING](TOD_TO_DWORD(tod) >= 12*3600000, "AM", "PM")
		case 'M':
			fs = INT_TO_STRING(td.HOUR(tod))
		case 'N':
			fs = pad(INT_TO_STRING(td.HOUR(tod)), "0")
		case 'O':
			fs = INT_TO_STRING(hour12())
		case 'P':
			fs = pad(INT_TO_STRING(hour12()), "0")
		case 'Q':
			fs = INT_TO_STRING(td.MINUTE(tod))
		case 'R':
			fs = pad(INT_TO_STRING(td.MINUTE(tod)), "0")
		case 'S':
			fs = INT_TO_STRING(REAL_TO_INT(td.SECOND(tod)))
		case 'T':
			fs = pad(INT_TO_STRING(REAL_TO_INT(td.SECOND(tod))), "0")
		case 'U':
			fs = INT_TO_STRING(ms)
		case 'V':
			fs = RIGHT("00"+INT_TO_STRING(ms), 3)
		case 'W':
			fs = pad(INT_TO_STRING(td.DAY_OF_MONTH(dx)), " ")
		case 'X':
			fs = pad(INT_TO_STRING(td.MONTH_OF_DATE(dx)), " ")
		}
		out = REPLACE(out, fs, 2, pos)
		pos = FIND(out, "#")
	}
	return out
}

// DWORD_TO_STRB converts a DWORD to a string of 32 bits, the highest first.
func DWORD_TO_STRB(in iec.DWORD) iec.STRING {
	b := make([]iec.BYTE, 32)
	for i := range b {
		b[i] = iec.BYTE(in>>(31-i)&1) + 48
	}
	return STR(b)
}

// DWORD_TO_STRF converts a DWORD to a string of n digits, 0..20, filled
// with leading zeros or cut to its lowest digits: DWORD_TO_STRF(123, 4) is
// "0123" and DWORD_TO_STRF(123, 2) is "23".
func DWORD_TO_STRF(in iec.DWORD, n iec.INT) iec.STRING {
	return FIX(DWORD_TO_STRING(in), LIMIT(0, n, 20), 48, 1)
}

// DWORD_TO_STRH converts a DWORD to a string of 8 hexadecimal digits.
func DWORD_TO_STRH(in iec.DWORD) iec.STRING {
	b := make([]iec.BYTE, 8)
	for i := range b {
		b[i] = hexChar(iec.BYTE(in >> (28 - 4*i) & 0xF))
	}
	return STR(b)
}

// EXEC calculates a simple term of two numbers and one operator, +, -, *, /
// or ^, or a function sqrt, sin, cos or tan of a number, and returns the
// result, or "ERROR".
func EXEC(str iec.STRING) iec.STRING {
	var r1, r2 iec.REAL
	exec := UPPERCASE(TRIM(str))
	pos := FINDB_NONUM(exec)
	if pos > 1 {
		r1 = STRING_TO_REAL(LEFT(exec, pos-1))
	}
	r2 = STRING_TO_REAL(RIGHT(exec, LEN(exec)-pos))
	exec = LEFT(exec, pos)
	pos = FINDB_NUM(exec)
	operator := RIGHT(exec, LEN(exec)-pos)
	switch {
	case operator == "" && LEN(str) == 0:
		return ""
	case operator == "":
		return str
	}
	switch operator {
	case "^":
		exec = REAL_TO_STRING(EXPT(r1, r2))
	case "SQRT":
		exec = REAL_TO_STRING(SQRT(r2))
	case "SIN":
		exec = REAL_TO_STRING(SIN(r2))
	case "COS":
		exec = REAL_TO_STRING(COS(r2))
	case "TAN":
		exec = REAL_TO_STRING(TAN(r2))
	case "*":
		exec = REAL_TO_STRING(r1 * r2)
	case "/":
		if r2 != 0 {
			exec = REAL_TO_STRING(r1 / r2)
		} else {
			exec = "ERROR"
		}
	case "+":
		exec = REAL_TO_STRING(r1 + r2)
	case "-":
		exec = REAL_TO_STRING(r1 - r2)
	default:
		exec = "ERROR"
	}
	switch {
	case exec == "ERROR":
	case FIND(exec, ".") == 0:
		// Some systems give an integer instead of a real.
		exec += ".0"
	case RIGHT(exec, 1) == ".":
		exec += "0"
	}
	return exec
}

// FILL returns a string of l characters with the code c, up to
// STRING_LENGTH.
func FILL(c iec.BYTE, l iec.INT) iec.STRING {
	l = LIMIT(0, l, STRING_LENGTH)
	r := make([]rune, l)
	for i := range r {
		r[i] = rune(c)
	}
	return iec.STRING(r)
}

// FIND_CHAR returns the position of the first character of str, from pos
// on, that is not a control character, or 0 if there is none.
func FIND_CHAR(str iec.STRING, pos iec.INT) iec.INT {
	c := CHARS(str)
	for i := max(pos, 1); int(i) <= len(c); i++ {
		x := c[i-1]
		if x > 31 && ((SETUP.EXTENDED_ASCII && x != 127) || (!SETUP.EXTENDED_ASCII && x < 127)) {
			return i
		}
	}
	return 0
}

// FIND_CTRL returns the position of the first control character of str,
// from pos on, or 0 if there is none.
func FIND_CTRL(str iec.STRING, pos iec.INT) iec.INT {
	c := CHARS(str)
	for i := max(pos, 1); int(i) <= len(c); i++ {
		if x := c[i-1]; x < 32 || x == 127 {
			return i
		}
	}
	return 0
}

// isNum reports whether a character belongs to a number: 0..9 or ".".
func isNum(x iec.BYTE) bool { return x > 47 && x < 58 || x == 46 }

// FIND_NONUM returns the position of the first character of str, from pos
// on, that is not 0..9 or ".", or 0 if there is none.
func FIND_NONUM(str iec.STRING, pos iec.INT) iec.INT {
	c := CHARS(str)
	for i := max(pos, 1); int(i) <= len(c); i++ {
		if !isNum(c[i-1]) {
			return i
		}
	}
	return 0
}

// FIND_NUM returns the position of the first character of str, from pos
// on, that is 0..9 or ".", or 0 if there is none.
func FIND_NUM(str iec.STRING, pos iec.INT) iec.INT {
	c := CHARS(str)
	for i := max(pos, 1); int(i) <= len(c); i++ {
		if isNum(c[i-1]) {
			return i
		}
	}
	return 0
}

// FINDB returns the position of the last str2 in str1, or 0 if there is
// none.
func FINDB(str1, str2 iec.STRING) iec.INT {
	length := LEN(str2)
	for pos := LEN(str1) - length + 1; pos >= 1; pos-- {
		if MID(str1, length, pos) == str2 {
			return pos
		}
	}
	return 0
}

// FINDB_NONUM returns the position of the last character of str that is not
// 0..9 or ".", or 0 if there is none.
func FINDB_NONUM(str iec.STRING) iec.INT {
	c := CHARS(str)
	for pos := len(c); pos >= 1; pos-- {
		if !isNum(c[pos-1]) {
			return iec.INT(pos)
		}
	}
	return 0
}

// FINDB_NUM returns the position of the last character of str that is 0..9
// or ".", or 0 if there is none.
func FINDB_NUM(str iec.STRING) iec.INT {
	c := CHARS(str)
	for pos := len(c); pos >= 1; pos-- {
		if isNum(c[pos-1]) {
			return iec.INT(pos)
		}
	}
	return 0
}

// FINDP returns the position of the first src in str from the position pos
// on, or 0 if there is none.
func FINDP(str, src iec.STRING, pos iec.INT) iec.INT {
	ls := LEN(str)
	lx := LEN(src)
	if ls < lx || lx == 0 {
		return 0
	}
	for i := max(pos, 1); i <= ls-lx+1; i++ {
		if MID(str, lx, i) == src {
			return i
		}
	}
	return 0
}

// FIX returns str cut or filled with the character c to the length l: m = 1
// cuts or fills at the start, m = 2 fills on both sides, with one more at
// the end for an odd number, and any other m cuts or fills at the end.
func FIX(str iec.STRING, l iec.INT, c iec.BYTE, m iec.INT) iec.STRING {
	n := LIMIT(0, l, STRING_LENGTH) - LEN(str)
	switch {
	case n <= 0:
		if m == 1 {
			return RIGHT(str, l)
		}
		return LEFT(str, l)
	case m == 1:
		return FILL(c, n) + str
	case m == 2:
		sx := FILL(c, (n+1)>>1)
		return LEFT(sx, n>>1) + str + sx
	}
	return str + FILL(c, n)
}

// FLOAT_TO_REAL converts a string to a REAL. The decimal separator can be
// "," or ".", the exponent starts with "E" or "e", and other characters are
// ignored.
func FLOAT_TO_REAL(flt iec.STRING) iec.REAL {
	pt := CHARS(flt)
	stop := len(pt)
	at := func(i int) iec.BYTE { return pt[i-1] }
	var sign iec.DINT = 1
	var tmp iec.DINT
	var d iec.INT
	var x iec.BYTE
	// The sign, up to the first digit or dot.
	i := 1
	for ; i <= stop; i++ {
		x = at(i)
		if x > 47 && x < 58 || x == 46 {
			break
		} else if x == 45 {
			sign = -1
		}
	}
	// The digits up to a separator or the exponent.
	for ; i <= stop; i++ {
		x = at(i)
		if x == 44 || x == 46 || x == 69 || x == 101 {
			break
		} else if x > 47 && x < 58 {
			tmp = tmp*10 + iec.DINT(x) - 48
		}
	}
	// The digits after the separator.
	if x == 44 || x == 46 {
		for i = i + 1; i <= stop; i++ {
			x = at(i)
			if x == 69 || x == 101 {
				break
			} else if x > 47 && x < 58 {
				tmp = tmp*10 + iec.DINT(x) - 48
				d--
			}
		}
	}
	if x == 69 || x == 101 {
		d += DEC_TO_INT(RIGHT(flt, iec.INT(stop-i)))
	}
	return math.EXPN(10, d) * iec.REAL(tmp*sign)
}

// FSTRING_TO_BYTE converts a string of the form 2#0101, 8#17, 16#2A or 123
// to a byte.
func FSTRING_TO_BYTE(in iec.STRING) iec.BYTE {
	switch {
	case LEFT(in, 2) == "2#":
		return BIN_TO_BYTE(RIGHT(in, LEN(in)-2))
	case LEFT(in, 2) == "8#":
		return OCT_TO_BYTE(RIGHT(in, LEN(in)-2))
	case LEFT(in, 3) == "16#":
		return HEX_TO_BYTE(RIGHT(in, LEN(in)-3))
	}
	return DEC_TO_BYTE(CLEAN(in, "0123456789"))
}

// FSTRING_TO_DT reads a date and time from sdt with the format fmt: #Y is
// the year, #M the month, #N the month's name, #D the day, #h the hour, #m
// the minute and #s the second; * skips a character, and other characters
// must match.
func FSTRING_TO_DT(sdt, fmt iec.STRING) iec.DT {
	var dy, dm, dd iec.INT = 1970, 1, 1
	var th, tm, ts iec.INT
	for fmt != "" {
		c := LEFT(fmt, 1)
		switch {
		case c == "*":
			fmt = DELETE(fmt, 1, 1)
			sdt = DELETE(sdt, 1, 1)
		case c == "#":
			c = MID(fmt, 1, 2)
			fmt = DELETE(fmt, 2, 1)
			var tmp iec.STRING
			if fmt == "" {
				tmp = sdt
			} else {
				end := FIND(sdt, LEFT(fmt, 1)) - 1
				tmp = LEFT(sdt, end)
				sdt = DELETE(sdt, end, 1)
			}
			switch c {
			case "Y":
				dy = STRING_TO_INT(tmp)
				if dy < 100 {
					dy += 2000
				}
			case "M":
				dm = STRING_TO_INT(tmp)
			case "N":
				dm = FSTRING_TO_MONTH(tmp, 0)
			case "D":
				dd = STRING_TO_INT(tmp)
			case "h":
				th = STRING_TO_INT(tmp)
			case "m":
				tm = STRING_TO_INT(tmp)
			case "s":
				ts = STRING_TO_INT(tmp)
			}
		case c == LEFT(sdt, 1):
			fmt = DELETE(fmt, 1, 1)
			sdt = DELETE(sdt, 1, 1)
		default:
			return iec.DT{}
		}
	}
	return td.SET_DT(dy, dm, dd, th, tm, ts)
}

// FSTRING_TO_DWORD converts a string of the form 2#0101, 8#17, 16#2A or 123
// to a DWORD.
func FSTRING_TO_DWORD(in iec.STRING) iec.DWORD {
	switch {
	case LEFT(in, 2) == "2#":
		return BIN_TO_DWORD(RIGHT(in, LEN(in)-2))
	case LEFT(in, 2) == "8#":
		return OCT_TO_DWORD(RIGHT(in, LEN(in)-2))
	case LEFT(in, 3) == "16#":
		return HEX_TO_DWORD(RIGHT(in, LEN(in)-3))
	}
	return DEC_TO_DWORD(CLEAN(in, "0123456789"))
}

// FSTRING_TO_MONTH converts a month's name or number to its number 1..12,
// in the language lang, or the default language if lang is 0.
func FSTRING_TO_MONTH(mth iec.STRING, lang iec.INT) iec.INT {
	lx := LANGUAGE.DEFAULT
	if lang != 0 {
		lx = min(lang, LANGUAGE.LMAX)
	}
	mth = CAPITALIZE(LOWERCASE(TRIM(mth)))
	for m := iec.INT(1); m <= 12; m++ {
		if mth == LANGUAGE.MONTHS[lx-1][m-1] || mth == LANGUAGE.MONTHS3[lx-1][m-1] {
			return m
		}
	}
	return STRING_TO_INT(mth)
}

// FSTRING_TO_WEEK converts a comma separated list of weekdays to a byte with
// bit 6 for monday ... bit 0 for sunday.
func FSTRING_TO_WEEK(week iec.STRING, lang iec.INT) iec.BYTE {
	var out iec.BYTE
	pos := FIND(week, ",")
	for pos > 0 {
		out |= SHR(iec.BYTE(128), FSTRING_TO_WEEKDAY(MID(week, pos-1, 1), lang))
		week = RIGHT(week, LEN(week)-pos)
		pos = FIND(week, ",")
	}
	return (out | SHR(iec.BYTE(128), FSTRING_TO_WEEKDAY(week, lang))) & 127
}

// FSTRING_TO_WEEKDAY converts a weekday's two letter name or number to its
// number 1..7, in the language lang, or the default language if lang is 0.
func FSTRING_TO_WEEKDAY(wday iec.STRING, lang iec.INT) iec.INT {
	ly := LANGUAGE.DEFAULT
	if lang != 0 {
		ly = min(lang, LANGUAGE.LMAX)
	}
	tmp := LEFT(CAPITALIZE(LOWERCASE(TRIM(wday))), 2)
	for i := iec.INT(1); i <= 7; i++ {
		if LANGUAGE.WEEKDAYS2[ly-1][i-1] == tmp {
			return i
		}
	}
	return STRING_TO_INT(wday)
}

// hexDigit returns the value of a hexadecimal digit and whether it is one.
func hexDigit(x iec.BYTE) (iec.BYTE, bool) {
	switch {
	case x > 47 && x < 58:
		return x - 48, true
	case x > 64 && x < 71:
		return x - 55, true
	case x > 96 && x < 103:
		return x - 87, true
	}
	return 0, false
}

// HEX_TO_BYTE converts a hexadecimal string to a byte. Other characters are
// ignored.
func HEX_TO_BYTE(hex iec.STRING) iec.BYTE {
	var out iec.BYTE
	for _, x := range CHARS(hex) {
		if v, ok := hexDigit(x); ok {
			out = out<<4 + v
		}
	}
	return out
}

// HEX_TO_DWORD converts a hexadecimal string to a DWORD. Other characters
// are ignored.
func HEX_TO_DWORD(hex iec.STRING) iec.DWORD {
	var out iec.DWORD
	for _, x := range CHARS(hex) {
		if v, ok := hexDigit(x); ok {
			out = out<<4 + iec.DWORD(v)
		}
	}
	return out
}

// IS_ALNUM reports whether str has only letters and digits.
func IS_ALNUM(str iec.STRING) iec.BOOL {
	return all(str, func(c iec.BYTE) iec.BOOL { return ISC_ALPHA(c) || ISC_NUM(c) })
}

// IS_ALPHA reports whether str has only letters.
func IS_ALPHA(str iec.STRING) iec.BOOL { return all(str, ISC_ALPHA) }

// IS_CC reports whether str has only characters of cmp.
func IS_CC(str, cmp iec.STRING) iec.BOOL {
	for _, c := range string(str) {
		if FIND(cmp, iec.STRING(c)) == 0 {
			return false
		}
	}
	return str != ""
}

// IS_CTRL reports whether str has only control characters.
func IS_CTRL(str iec.STRING) iec.BOOL { return all(str, ISC_CTRL) }

// IS_HEX reports whether str has only hexadecimal digits.
func IS_HEX(str iec.STRING) iec.BOOL { return all(str, ISC_HEX) }

// IS_LOWER reports whether str has only lower case letters.
func IS_LOWER(str iec.STRING) iec.BOOL { return all(str, ISC_LOWER) }

// IS_NCC reports whether str has none of the characters of cmp.
func IS_NCC(str, cmp iec.STRING) iec.BOOL {
	for _, c := range string(str) {
		if FIND(cmp, iec.STRING(c)) > 0 {
			return false
		}
	}
	return true
}

// IS_NUM reports whether str has only digits.
func IS_NUM(str iec.STRING) iec.BOOL { return all(str, ISC_NUM) }

// IS_UPPER reports whether str has only upper case letters.
func IS_UPPER(str iec.STRING) iec.BOOL { return all(str, ISC_UPPER) }

// ISC_ALPHA reports whether a character is a letter, with the ISO 8859-1
// letters if SETUP.EXTENDED_ASCII is true.
func ISC_ALPHA(in iec.BYTE) iec.BOOL {
	if SETUP.EXTENDED_ASCII {
		return in > 64 && in < 91 || in > 191 && in != 215 && in != 247 || in > 96 && in < 123
	}
	return in > 64 && in < 91 || in > 96 && in < 123
}

// ISC_CTRL reports whether a character is a control character.
func ISC_CTRL(in iec.BYTE) iec.BOOL { return in < 32 || in == 127 }

// ISC_HEX reports whether a character is a hexadecimal digit.
func ISC_HEX(in iec.BYTE) iec.BOOL {
	_, ok := hexDigit(in)
	return iec.BOOL(ok)
}

// ISC_LOWER reports whether a character is a lower case letter.
func ISC_LOWER(in iec.BYTE) iec.BOOL {
	if SETUP.EXTENDED_ASCII {
		return in > 96 && in < 123 || in > 222 && in != 247
	}
	return in > 96 && in < 123
}

// ISC_NUM reports whether a character is a digit.
func ISC_NUM(in iec.BYTE) iec.BOOL { return in > 47 && in < 58 }

// ISC_UPPER reports whether a character is an upper case letter.
func ISC_UPPER(in iec.BYTE) iec.BOOL {
	if SETUP.EXTENDED_ASCII {
		return in > 64 && in < 91 || in > 191 && in < 223 && in != 215
	}
	return in > 64 && in < 91
}

// LOWERCASE returns str in lower case.
func LOWERCASE(str iec.STRING) iec.STRING { return mapChars(str, TO_LOWER) }

// MESSAGE_4R shows the messages M0..MM in turn on MX, the next one on each
// rising edge of CLK or after T1 while CLK stays true. MN is the number of
// the message and TR is true for one scan when it changes. MX is empty
// while ENQ is false.
type MESSAGE_4R struct {
	M0, M1, M2, M3 iec.STRING
	MM             iec.INT  // default 3
	ENQ            iec.BOOL // default TRUE
	CLK            iec.BOOL // default TRUE
	T1             iec.TIME // default T#3s
	MX             iec.STRING
	MN             iec.INT
	TR             iec.BOOL

	timer timers.TON
	edge  iec.BOOL
}

// INIT resets the block and sets MM, ENQ, CLK and T1 to their initial
// values.
func (m *MESSAGE_4R) INIT() {
	*m = MESSAGE_4R{MM: 3, ENQ: true, CLK: true, T1: iec.TIME(3 * time.Second)}
}

// Execute runs the block once.
func (m *MESSAGE_4R) Execute(now time.Time) {
	m.TR = false
	if !m.ENQ {
		m.MX = ""
		m.MN = 0
		return
	}
	if (!m.edge && m.CLK) || m.timer.Q {
		m.MN = math.INC1(m.MN, m.MM)
		m.TR = true
		m.timer.IN = false
		m.timer.Execute(now)
		switch m.MN {
		case 0:
			m.MX = m.M0
		case 1:
			m.MX = m.M1
		case 2:
			m.MX = m.M2
		case 3:
			m.MX = m.M3
		}
	}
	m.edge = m.CLK
	m.timer.IN = m.CLK
	m.timer.PT = m.T1
	m.timer.Execute(now)
}

// MESSAGE_8 shows on M the message S1..S8 of the first input IN1..IN8 that
// is true, or an empty string.
type MESSAGE_8 struct {
	IN1, IN2, IN3, IN4, IN5, IN6, IN7, IN8 iec.BOOL
	S1, S2, S3, S4, S5, S6, S7, S8         iec.STRING
	M                                      iec.STRING
}

// INIT resets the block.
func (m *MESSAGE_8) INIT() { *m = MESSAGE_8{} }

// Execute runs the block once.
func (m *MESSAGE_8) Execute(now time.Time) {
	in := [8]iec.BOOL{m.IN1, m.IN2, m.IN3, m.IN4, m.IN5, m.IN6, m.IN7, m.IN8}
	s := [8]iec.STRING{m.S1, m.S2, m.S3, m.S4, m.S5, m.S6, m.S7, m.S8}
	m.M = ""
	for i, b := range in {
		if b {
			m.M = s[i]
			return
		}
	}
}

// MIRROR returns str reversed.
func MIRROR(str iec.STRING) iec.STRING {
	r := []rune(string(str))
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return iec.STRING(r)
}

// MONTH_TO_STRING returns the name of the month mth, 1..12, in the language
// lang, or the default language if lang <= 0: the full name if lx is 0 and
// three letters if lx is 3.
func MONTH_TO_STRING(mth, lang, lx iec.INT) iec.STRING {
	ly := LANGUAGE.DEFAULT
	if lang > 0 {
		ly = min(lang, LANGUAGE.LMAX)
	}
	switch {
	case mth < 1 || mth > 12:
		return ""
	case lx == 0:
		return LANGUAGE.MONTHS[ly-1][mth-1]
	case lx == 3:
		return LANGUAGE.MONTHS3[ly-1][mth-1]
	}
	return ""
}

// OCT_TO_BYTE converts an octal string to a byte. Other characters are
// ignored.
func OCT_TO_BYTE(oct iec.STRING) iec.BYTE {
	var out iec.BYTE
	for _, x := range CHARS(oct) {
		if x > 47 && x < 56 {
			out = out<<3 + x - 48
		}
	}
	return out
}

// OCT_TO_DWORD converts an octal string to a DWORD. Other characters are
// ignored.
func OCT_TO_DWORD(oct iec.STRING) iec.DWORD {
	var out iec.DWORD
	for _, x := range CHARS(oct) {
		if x > 47 && x < 56 {
			out = out<<3 + iec.DWORD(x) - 48
		}
	}
	return out
}

// REAL_TO_STRF converts a REAL to a string with n, 0..7, digits after the
// decimal separator d.
func REAL_TO_STRF(in iec.REAL, n iec.INT, d iec.STRING) iec.STRING {
	n = LIMIT(0, n, 7)
	o := ABS(in) * math.EXP10(iec.REAL(n))
	out := DINT_TO_STRING(REAL_TO_DINT(o))
	for i := LEN(out); i <= n; i++ {
		out = "0" + out
	}
	if n > 0 {
		out = INSERT(out, d, LEN(out)-n)
	}
	if in < 0 {
		out = "-" + out
	}
	return out
}

// REPLACE_ALL returns str with each src replaced by rep.
func REPLACE_ALL(str, src, rep iec.STRING) iec.STRING {
	lx := LEN(src)
	lp := LEN(rep)
	pos := FINDP(str, src, 1)
	for pos > 0 {
		str = REPLACE(str, rep, lx, pos)
		pos = FINDP(str, src, pos+lp)
	}
	return str
}

// REPLACE_CHARS returns str with each character of src replaced by the
// character at the same position of rep.
func REPLACE_CHARS(str, src, rep iec.STRING) iec.STRING {
	a, b := LEN(src), LEN(rep)
	if a < b {
		rep = LEFT(rep, a)
	} else if b < a {
		src = LEFT(src, b)
	}
	r := []rune(string(str))
	for i, c := range r {
		if p := FIND(src, iec.STRING(c)); p > 0 {
			r[i] = []rune(string(rep))[p-1]
		}
	}
	return iec.STRING(r)
}

// REPLACE_UML returns str with the umlauts Ä, ä, Ö, ö, Ü, ü and ß replaced
// by Ae, ae, Oe, oe, Ue, ue and ss, up to STRING_LENGTH characters.
func REPLACE_UML(str iec.STRING) iec.STRING {
	out := make([]rune, 0, len(str))
	for i, c := range []rune(string(str)) {
		if i >= int(STRING_LENGTH) {
			break
		}
		if c < 127 {
			out = append(out, c)
			continue
		}
		su := []rune(string(TO_UML(charCode(c))))
		if c > 255 {
			su = []rune{c}
		}
		out = append(out, su[0])
		if len(out) < int(STRING_LENGTH) && len(su) > 1 {
			out = append(out, su[1])
		}
	}
	return iec.STRING(out)
}

// charCode returns a character's code, or ? for a character outside
// ISO 8859-1.
func charCode(c rune) iec.BYTE {
	if c > 255 {
		return '?'
	}
	return iec.BYTE(c)
}

// TICKER shows N characters of TEXT on DISPLAY, one character further every
// PT, so that the text moves across. A text of N characters or less is shown
// whole.
type TICKER struct {
	N       iec.INT
	PT      iec.TIME
	TEXT    *iec.STRING
	DISPLAY iec.STRING

	delay timers.TP
	step  iec.INT
}

// INIT resets the block.
func (t *TICKER) INIT() { *t = TICKER{TEXT: t.TEXT} }

// Execute runs the block once.
func (t *TICKER) Execute(now time.Time) {
	if t.TEXT == nil {
		return
	}
	text := *t.TEXT
	if t.N >= LEN(text) {
		t.DISPLAY = text
		return
	}
	if !t.delay.Q {
		t.step++
		if t.step > LEN(text) {
			t.step = 1
		}
		t.DISPLAY = MID(text, t.N, t.step)
		t.delay.IN = true
		t.delay.PT = t.PT
	} else {
		t.delay.IN = false
	}
	t.delay.Execute(now)
}

// TO_LOWER returns the lower case of a character.
func TO_LOWER(in iec.BYTE) iec.BYTE {
	switch {
	case in > 64 && in < 91:
		return in | 0x20
	case in > 191 && in < 223 && in != 215 && bool(SETUP.EXTENDED_ASCII):
		return in | 0x20
	}
	return in
}

// TO_UML returns the two letters of an umlaut, Ä is "Ae", or the character
// itself.
func TO_UML(in iec.BYTE) iec.STRING {
	switch in {
	case 196:
		return "Ae"
	case 214:
		return "Oe"
	case 220:
		return "Ue"
	case 223:
		return "ss"
	case 228:
		return "ae"
	case 246:
		return "oe"
	case 252:
		return "ue"
	}
	return CHR_TO_STRING(in)
}

// TO_UPPER returns the upper case of a character.
func TO_UPPER(in iec.BYTE) iec.BYTE {
	switch {
	case in > 96 && in < 123:
		return in & 0xDF
	case in > 223 && in != 247 && in != 255 && bool(SETUP.EXTENDED_ASCII):
		return in & 0xDF
	}
	return in
}

// TRIM returns str without any blanks.
func TRIM(str iec.STRING) iec.STRING {
	return DEL_CHARS(str, " ")
}

// TRIM1 returns str with each run of blanks replaced by one blank and
// without leading and trailing blanks.
func TRIM1(str iec.STRING) iec.STRING {
	for pos := FIND(str, "  "); pos > 0; pos = FIND(str, "  ") {
		str = REPLACE(str, " ", 2, pos)
	}
	if LEFT(str, 1) == " " {
		str = DELETE(str, 1, 1)
	}
	if RIGHT(str, 1) == " " {
		str = DELETE(str, 1, LEN(str))
	}
	return str
}

// TRIME returns str without leading and trailing blanks.
func TRIME(str iec.STRING) iec.STRING {
	for LEFT(str, 1) == " " {
		str = DELETE(str, 1, 1)
	}
	for RIGHT(str, 1) == " " {
		str = DELETE(str, 1, LEN(str))
	}
	return str
}

// UPPERCASE returns str in upper case.
func UPPERCASE(str iec.STRING) iec.STRING { return mapChars(str, TO_UPPER) }

// WEEKDAY_TO_STRING returns the name of the weekday wday, 1..7, in the
// language lang, or the default language if lang is 0: the full name if lx
// is 0 and two letters if lx is 2.
func WEEKDAY_TO_STRING(wday, lang, lx iec.INT) iec.STRING {
	ly := LANGUAGE.DEFAULT
	if lang != 0 {
		ly = min(lang, LANGUAGE.LMAX)
	}
	switch {
	case wday < 1 || wday > 7:
		return ""
	case lx == 0:
		return LANGUAGE.WEEKDAYS[ly-1][wday-1]
	case lx == 2:
		return LANGUAGE.WEEKDAYS2[ly-1][wday-1]
	}
	return ""
}
