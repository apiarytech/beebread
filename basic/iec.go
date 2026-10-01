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

package basic

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// This file holds the IEC 61131-3 standard conversions and string functions
// the OSCAT sources rely on, with the behaviour OSCAT was written for: a
// date is a number of seconds since 1970-01-01, a time of day and a TIME a
// number of milliseconds, and string positions start at 1.

// Integer is any integer type.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Float is any floating point type.
type Float interface {
	~float32 | ~float64
}

// LIMIT returns in limited to the range mn..mx.
func LIMIT[T Integer | Float](mn, in, mx T) T {
	if in < mn {
		return mn
	}
	if in > mx {
		return mx
	}
	return in
}

// SEL returns in1 if g is true and in0 otherwise.
func SEL[T any](g iec.BOOL, in0, in1 T) T {
	if g {
		return in1
	}
	return in0
}

// BOOL_TO_INT returns 1 for true and 0 for false.
func BOOL_TO_INT(b iec.BOOL) iec.INT {
	if b {
		return 1
	}
	return 0
}

// ROUND_TO converts a real value to the integer type T, rounding to the
// nearest integer, halves to the even one. A value outside the 64 bit
// integers saturates before it is cut to T.
func ROUND_TO[T Integer, F Float](x F) T {
	return floatTo[T](math.RoundToEven(float64(x)))
}

// TRUNC_TO converts a real value to the integer type T, dropping the
// fraction.
func TRUNC_TO[T Integer, F Float](x F) T {
	return floatTo[T](math.Trunc(float64(x)))
}

func floatTo[T Integer](v float64) T {
	var i int64
	switch {
	case v != v:
		i = 0
	case v >= math.MaxInt64:
		i = math.MaxInt64
	case v <= math.MinInt64:
		i = math.MinInt64
	default:
		i = int64(v)
	}
	return T(i)
}

// REAL_TO_INT converts a REAL to an INT, rounding to the nearest integer.
func REAL_TO_INT(x iec.REAL) iec.INT { return ROUND_TO[iec.INT](x) }

// REAL_TO_DINT converts a REAL to a DINT, rounding to the nearest integer.
func REAL_TO_DINT(x iec.REAL) iec.DINT { return ROUND_TO[iec.DINT](x) }

// REAL_TO_DWORD converts a REAL to a DWORD, rounding to the nearest integer.
func REAL_TO_DWORD(x iec.REAL) iec.DWORD { return ROUND_TO[iec.DWORD](x) }

// TRUNC returns the integer part of a REAL.
func TRUNC(x iec.REAL) iec.DINT { return TRUNC_TO[iec.DINT](x) }

// The standard numerical functions on REAL, calculated in 64 bits.

func SQRT(x iec.REAL) iec.REAL { return iec.REAL(math.Sqrt(float64(x))) }
func SIN(x iec.REAL) iec.REAL  { return iec.REAL(math.Sin(float64(x))) }
func COS(x iec.REAL) iec.REAL  { return iec.REAL(math.Cos(float64(x))) }
func TAN(x iec.REAL) iec.REAL  { return iec.REAL(math.Tan(float64(x))) }
func ASIN(x iec.REAL) iec.REAL { return iec.REAL(math.Asin(float64(x))) }
func ACOS(x iec.REAL) iec.REAL { return iec.REAL(math.Acos(float64(x))) }
func ATAN(x iec.REAL) iec.REAL { return iec.REAL(math.Atan(float64(x))) }
func EXP(x iec.REAL) iec.REAL  { return iec.REAL(math.Exp(float64(x))) }
func LN(x iec.REAL) iec.REAL   { return iec.REAL(math.Log(float64(x))) }
func LOG(x iec.REAL) iec.REAL  { return iec.REAL(math.Log10(float64(x))) }
func ABS(x iec.REAL) iec.REAL  { return iec.REAL(math.Abs(float64(x))) }
func EXPT(x, y iec.REAL) iec.REAL {
	return iec.REAL(math.Pow(float64(x), float64(y)))
}

// Dates and times.

// DT_TO_DWORD returns the seconds from 1970-01-01 00:00:00 to dt. The zero
// DT, which a variable that was never set holds, is 0.
func DT_TO_DWORD(dt iec.DT) iec.DWORD {
	t := time.Time(dt)
	if t.IsZero() {
		return 0
	}
	return iec.DWORD(t.Unix())
}

// DWORD_TO_DT returns the date and time d seconds after 1970-01-01 00:00:00.
func DWORD_TO_DT(d iec.DWORD) iec.DT {
	return iec.DT(time.Unix(int64(d), 0).UTC())
}

// DATE_TO_DWORD returns the seconds from 1970-01-01 to the start of the date.
func DATE_TO_DWORD(d iec.DATE) iec.DWORD {
	t := time.Time(d)
	if t.IsZero() {
		return 0
	}
	s := iec.DWORD(t.Unix())
	return s - s%86400
}

// DWORD_TO_DATE returns the date of the day d seconds after 1970-01-01.
func DWORD_TO_DATE(d iec.DWORD) iec.DATE {
	return iec.DATE(time.Unix(int64(d-d%86400), 0).UTC())
}

// DT_TO_DATE returns the date of dt.
func DT_TO_DATE(dt iec.DT) iec.DATE { return DWORD_TO_DATE(DT_TO_DWORD(dt)) }

// DATE_TO_DT returns the start of the date d.
func DATE_TO_DT(d iec.DATE) iec.DT { return DWORD_TO_DT(DATE_TO_DWORD(d)) }

// TOD_TO_DWORD returns the milliseconds from midnight to the time of day.
func TOD_TO_DWORD(tod iec.TOD) iec.DWORD {
	t := time.Time(tod)
	return iec.DWORD(t.Hour()*3600000 + t.Minute()*60000 + t.Second()*1000 + t.Nanosecond()/1000000)
}

// DWORD_TO_TOD returns the time of day d milliseconds after midnight.
func DWORD_TO_TOD(d iec.DWORD) iec.TOD {
	return iec.TOD(time.Time{}.Add(time.Duration(d%86400000) * time.Millisecond))
}

// DT_TO_TOD returns the time of day of dt.
func DT_TO_TOD(dt iec.DT) iec.TOD {
	return DWORD_TO_TOD(DT_TO_DWORD(dt) % 86400 * 1000)
}

// TOD_TO_TIME returns the time from midnight to the time of day.
func TOD_TO_TIME(tod iec.TOD) iec.TIME { return DWORD_TO_TIME(TOD_TO_DWORD(tod)) }

// TIME_TO_DWORD returns a TIME in milliseconds.
func TIME_TO_DWORD(t iec.TIME) iec.DWORD {
	return iec.DWORD(time.Duration(t).Milliseconds())
}

// DWORD_TO_TIME returns the TIME of d milliseconds.
func DWORD_TO_TIME(d iec.DWORD) iec.TIME {
	return iec.TIME(time.Duration(d) * time.Millisecond)
}

// TIME_TO_REAL returns a TIME in milliseconds.
func TIME_TO_REAL(t iec.TIME) iec.REAL { return iec.REAL(TIME_TO_DWORD(t)) }

// REAL_TO_TIME returns the TIME of x milliseconds.
func REAL_TO_TIME(x iec.REAL) iec.TIME { return DWORD_TO_TIME(REAL_TO_DWORD(x)) }

// Conversions between numbers and strings.

// INT_TO_STRING returns the decimal string of an INT.
func INT_TO_STRING(i iec.INT) iec.STRING { return iec.STRING(strconv.Itoa(int(i))) }

// DINT_TO_STRING returns the decimal string of a DINT.
func DINT_TO_STRING(i iec.DINT) iec.STRING { return iec.STRING(strconv.Itoa(int(i))) }

// DWORD_TO_STRING returns the decimal string of a DWORD.
func DWORD_TO_STRING(d iec.DWORD) iec.STRING {
	return iec.STRING(strconv.FormatUint(uint64(d), 10))
}

// REAL_TO_STRING returns the shortest decimal string of a REAL.
func REAL_TO_STRING(x iec.REAL) iec.STRING {
	return iec.STRING(strconv.FormatFloat(float64(x), 'g', -1, 32))
}

// STRING_TO_DINT reads the integer at the start of s, after blanks, with an
// optional sign. It is 0 if there is none.
func STRING_TO_DINT(s iec.STRING) iec.DINT {
	t := strings.TrimLeft(string(s), " ")
	n := 0
	if n < len(t) && (t[n] == '-' || t[n] == '+') {
		n++
	}
	for n < len(t) && t[n] >= '0' && t[n] <= '9' {
		n++
	}
	v, _ := strconv.ParseInt(t[:n], 10, 64)
	return iec.DINT(v)
}

// STRING_TO_INT reads the integer at the start of s, as STRING_TO_DINT.
func STRING_TO_INT(s iec.STRING) iec.INT { return iec.INT(STRING_TO_DINT(s)) }

// STRING_TO_REAL reads the number at the start of s, after blanks. It is 0
// if there is none.
func STRING_TO_REAL(s iec.STRING) iec.REAL {
	t := strings.TrimLeft(string(s), " ")
	// The longest start of t that is a number.
	for n := len(t); n > 0; n-- {
		if v, err := strconv.ParseFloat(t[:n], 32); err == nil {
			return iec.REAL(v)
		}
	}
	return 0
}

// Strings. A string is a sequence of characters; positions start at 1.

// LEN returns the number of characters of s.
func LEN(s iec.STRING) iec.INT {
	n := 0
	for range string(s) {
		n++
	}
	return iec.INT(n)
}

// LEFT returns the first n characters of s.
func LEFT(s iec.STRING, n iec.INT) iec.STRING {
	r := []rune(string(s))
	if n <= 0 {
		return ""
	}
	if int(n) >= len(r) {
		return s
	}
	return iec.STRING(r[:n])
}

// RIGHT returns the last n characters of s.
func RIGHT(s iec.STRING, n iec.INT) iec.STRING {
	r := []rune(string(s))
	if n <= 0 {
		return ""
	}
	if int(n) >= len(r) {
		return s
	}
	return iec.STRING(r[len(r)-int(n):])
}

// MID returns n characters of s, starting at position pos.
func MID(s iec.STRING, n, pos iec.INT) iec.STRING {
	r := []rune(string(s))
	if n <= 0 || pos < 1 || int(pos) > len(r) {
		return ""
	}
	end := int(pos) - 1 + int(n)
	if end > len(r) {
		end = len(r)
	}
	return iec.STRING(r[pos-1 : end])
}

// CONCAT returns the strings joined.
func CONCAT(s ...iec.STRING) iec.STRING {
	var out iec.STRING
	for _, v := range s {
		out += v
	}
	return out
}

// INSERT returns s1 with s2 inserted after position pos.
func INSERT(s1, s2 iec.STRING, pos iec.INT) iec.STRING {
	r := []rune(string(s1))
	if pos < 0 {
		pos = 0
	}
	if int(pos) > len(r) {
		pos = iec.INT(len(r))
	}
	return iec.STRING(r[:pos]) + s2 + iec.STRING(r[pos:])
}

// DELETE returns s without the n characters starting at position pos.
func DELETE(s iec.STRING, n, pos iec.INT) iec.STRING {
	r := []rune(string(s))
	if n <= 0 || pos < 1 || int(pos) > len(r) {
		return s
	}
	end := int(pos) - 1 + int(n)
	if end > len(r) {
		end = len(r)
	}
	return iec.STRING(r[:pos-1]) + iec.STRING(r[end:])
}

// REPLACE returns s1 with the n characters starting at position pos
// replaced by s2.
func REPLACE(s1, s2 iec.STRING, n, pos iec.INT) iec.STRING {
	r := []rune(string(s1))
	if pos < 1 || int(pos) > len(r)+1 {
		return s1
	}
	if n < 0 {
		n = 0
	}
	end := int(pos) - 1 + int(n)
	if end > len(r) {
		end = len(r)
	}
	return iec.STRING(r[:pos-1]) + s2 + iec.STRING(r[end:])
}

// FIND returns the position of the first s2 in s1, or 0 if there is none.
func FIND(s1, s2 iec.STRING) iec.INT {
	if s2 == "" {
		return 0
	}
	a, b := []rune(string(s1)), []rune(string(s2))
	for i := 0; i+len(b) <= len(a); i++ {
		j := 0
		for j < len(b) && a[i+j] == b[j] {
			j++
		}
		if j == len(b) {
			return iec.INT(i + 1)
		}
	}
	return 0
}

// CHARS returns the character codes of s.
func CHARS(s iec.STRING) []iec.BYTE {
	out := make([]iec.BYTE, 0, len(s))
	for _, c := range string(s) {
		out = append(out, charCode(c))
	}
	return out
}

// STR returns the string of the character codes c, up to the first 0.
func STR(c []iec.BYTE) iec.STRING {
	r := make([]rune, 0, len(c))
	for _, b := range c {
		if b == 0 {
			break
		}
		r = append(r, rune(b))
	}
	return iec.STRING(r)
}

// charCode returns a character's code; a character that OSCAT's single byte
// strings cannot hold is a question mark.
func charCode(c rune) iec.BYTE {
	if c > 255 {
		return '?'
	}
	return iec.BYTE(c)
}
