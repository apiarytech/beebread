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

package string

import (
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

func TestStrings(t *testing.T) {
	tests := []struct {
		name      string
		got, want iec.STRING
	}{
		{"BYTE_TO_STRB", BYTE_TO_STRB(0xA5), "10100101"},
		{"BYTE_TO_STRH", BYTE_TO_STRH(0x3F), "3F"},
		{"DWORD_TO_STRB", DWORD_TO_STRB(5), "00000000000000000000000000000101"},
		{"DWORD_TO_STRH", DWORD_TO_STRH(0x12AB), "000012AB"},
		{"DWORD_TO_STRF pad", DWORD_TO_STRF(123, 4), "0123"},
		{"DWORD_TO_STRF cut", DWORD_TO_STRF(123, 2), "23"},
		{"CAPITALIZE", CAPITALIZE("hello big world"), "Hello Big World"},
		{"CHARNAME euro", CHARNAME(128), "euro"},
		{"CHARNAME ä", CHARNAME(228), "auml"},
		{"CHARNAME &", CHARNAME('&'), "amp"},
		{"CHARNAME Û", CHARNAME(219), "Ucirc"},
		{"CHARNAME none", CHARNAME('a'), "a"},
		{"CHR_TO_STRING", CHR_TO_STRING(228), "ä"},
		{"CLEAN", CLEAN("a1b2c3", "abc"), "abc"},
		{"DEL_CHARS", DEL_CHARS("a1b2c3", "abc"), "123"},
		{"EXEC +", EXEC("1.5 + 2"), "3.5"},
		{"EXEC *", EXEC("3*4"), "12.0"},
		{"EXEC /0", EXEC("3/0"), "ERROR"},
		{"EXEC sqrt", EXEC("sqrt 16"), "4.0"},
		{"EXEC number", EXEC("42"), "42"},
		{"FILL", FILL('x', 3), "xxx"},
		{"FIX cut end", FIX("abcdef", 3, ' ', 0), "abc"},
		{"FIX cut start", FIX("abcdef", 3, ' ', 1), "def"},
		{"FIX fill start", FIX("ab", 5, '*', 1), "***ab"},
		{"FIX fill both", FIX("ab", 5, '*', 2), "*ab**"},
		{"FIX fill end", FIX("ab", 4, '*', 0), "ab**"},
		{"LOWERCASE", LOWERCASE("ÄBC Def"), "äbc def"},
		{"UPPERCASE", UPPERCASE("äbc Def"), "ÄBC DEF"},
		{"MIRROR", MIRROR("abc"), "cba"},
		{"MONTH_TO_STRING", MONTH_TO_STRING(3, 2, 0), "März"},
		{"MONTH_TO_STRING 3", MONTH_TO_STRING(8, 3, 3), "Aou"},
		{"MONTH_TO_STRING out of range", MONTH_TO_STRING(13, 1, 0), ""},
		{"WEEKDAY_TO_STRING", WEEKDAY_TO_STRING(7, 1, 0), "Sunday"},
		{"WEEKDAY_TO_STRING 2", WEEKDAY_TO_STRING(1, 2, 2), "Mo"},
		{"REAL_TO_STRF", REAL_TO_STRF(3.14159, 2, "."), "3.14"},
		{"REAL_TO_STRF small", REAL_TO_STRF(-0.05, 3, ","), "-0,050"},
		{"REAL_TO_STRF 0 digits", REAL_TO_STRF(2.6, 0, "."), "3"},
		{"REPLACE_ALL", REPLACE_ALL("a-b-c", "-", "--"), "a--b--c"},
		{"REPLACE_CHARS", REPLACE_CHARS("hello", "lo", "01"), "he001"},
		{"REPLACE_UML", REPLACE_UML("Müßig Äpfel"), "Muessig Aepfel"},
		{"TO_UML", TO_UML(223), "ss"},
		{"TRIM", TRIM(" a b  c "), "abc"},
		{"TRIM1", TRIM1("  a   b c  "), "a b c"},
		{"TRIME", TRIME("  a b  "), "a b"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestStringNumbers(t *testing.T) {
	tests := []struct {
		name      string
		got, want int64
	}{
		{"BIN_TO_BYTE", int64(BIN_TO_BYTE("1010")), 10},
		{"BIN_TO_DWORD", int64(BIN_TO_DWORD("1 0000 0000")), 256},
		{"DEC_TO_BYTE", int64(DEC_TO_BYTE("2x5")), 25},
		{"DEC_TO_DWORD", int64(DEC_TO_DWORD("4000000000")), 4000000000},
		{"DEC_TO_INT", int64(DEC_TO_INT("-123")), -123},
		{"HEX_TO_BYTE", int64(HEX_TO_BYTE("fF")), 255},
		{"HEX_TO_DWORD", int64(HEX_TO_DWORD("DEADbeef")), 0xDEADBEEF},
		{"OCT_TO_BYTE", int64(OCT_TO_BYTE("17")), 15},
		{"OCT_TO_DWORD", int64(OCT_TO_DWORD("777")), 511},
		{"FSTRING_TO_BYTE bin", int64(FSTRING_TO_BYTE("2#1111")), 15},
		{"FSTRING_TO_BYTE oct", int64(FSTRING_TO_BYTE("8#17")), 15},
		{"FSTRING_TO_BYTE hex", int64(FSTRING_TO_BYTE("16#0F")), 15},
		{"FSTRING_TO_BYTE dec", int64(FSTRING_TO_BYTE("15")), 15},
		{"FSTRING_TO_DWORD", int64(FSTRING_TO_DWORD("16#FFFF")), 65535},
		{"FSTRING_TO_MONTH name", int64(FSTRING_TO_MONTH(" march ", 1)), 3},
		{"FSTRING_TO_MONTH short", int64(FSTRING_TO_MONTH("okt", 2)), 10},
		{"FSTRING_TO_MONTH number", int64(FSTRING_TO_MONTH("11", 0)), 11},
		{"FSTRING_TO_WEEKDAY", int64(FSTRING_TO_WEEKDAY("fr", 1)), 5},
		{"FSTRING_TO_WEEKDAY number", int64(FSTRING_TO_WEEKDAY("6", 1)), 6},
		{"FSTRING_TO_WEEK", int64(FSTRING_TO_WEEK("mo,we,su", 1)), 0x40 | 0x10 | 0x01},
		{"CHARCODE euro", int64(CHARCODE("euro")), 128},
		{"CHARCODE auml", int64(CHARCODE("auml")), 228},
		{"CHARCODE yuml", int64(CHARCODE("yuml")), 255},
		{"CHARCODE char", int64(CHARCODE("A")), 65},
		{"CHARCODE unknown", int64(CHARCODE("nope")), 0},
		{"CODE", int64(CODE("abc", 2)), 'b'},
		{"CODE out of range", int64(CODE("abc", 4)), 0},
		{"COUNT_CHAR", int64(COUNT_CHAR("banana", 'a')), 3},
		{"COUNT_SUBSTRING", int64(COUNT_SUBSTRING("an", "banana")), 2},
		{"FIND_CHAR", int64(FIND_CHAR("\t\n ab", 1)), 3},
		{"FIND_CTRL", int64(FIND_CTRL("ab\tc", 1)), 3},
		{"FIND_NONUM", int64(FIND_NONUM("12.5kg", 1)), 5},
		{"FIND_NUM", int64(FIND_NUM("abc3", 1)), 4},
		{"FINDB", int64(FINDB("abcabc", "bc")), 5},
		{"FINDB_NONUM", int64(FINDB_NONUM("ab12")), 2},
		{"FINDB_NUM", int64(FINDB_NUM("1a2bc")), 3},
		{"FINDP", int64(FINDP("abcabc", "bc", 3)), 5},
		{"FINDP none", int64(FINDP("abc", "x", 1)), 0},
		{"TO_LOWER", int64(TO_LOWER('A')), 'a'},
		{"TO_UPPER ä", int64(TO_UPPER(228)), 196},
		{"TO_UPPER ÿ", int64(TO_UPPER(255)), 255},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.want)
		}
	}
	for _, tt := range []struct {
		in   iec.STRING
		want iec.REAL
	}{{"3.25", 3.25}, {"-1,5", -1.5}, {"2.5E3", 2500}, {"12e-1", 1.2}, {"x7", 7}} {
		if got := FLOAT_TO_REAL(tt.in); got < tt.want-1e-5 || got > tt.want+1e-5 {
			t.Errorf("FLOAT_TO_REAL(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
	bools := []struct {
		name      string
		got, want iec.BOOL
	}{
		{"IS_ALNUM", IS_ALNUM("abc123"), true},
		{"IS_ALNUM blank", IS_ALNUM("abc 123"), false},
		{"IS_ALPHA", IS_ALPHA("Grüße"), true},
		{"IS_CC", IS_CC("aab", "ab"), true},
		{"IS_CC empty", IS_CC("", "ab"), false},
		{"IS_CTRL", IS_CTRL("\t\n"), true},
		{"IS_HEX", IS_HEX("1aF"), true},
		{"IS_HEX g", IS_HEX("1g"), false},
		{"IS_LOWER", IS_LOWER("abcä"), true},
		{"IS_NCC", IS_NCC("abc", "xyz"), true},
		{"IS_NUM", IS_NUM("0123"), true},
		{"IS_NUM empty", IS_NUM(""), false},
		{"IS_UPPER", IS_UPPER("ABCÄ"), true},
		{"ISC_ALPHA ×", ISC_ALPHA(215), false},
	}
	for _, tt := range bools {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestDateStrings(t *testing.T) {
	x := iec.DT(time.Date(2024, 3, 5, 14, 7, 9, 0, time.UTC))
	tests := []struct {
		fmt  iec.STRING
		lang iec.INT
		want iec.STRING
	}{
		{"#A-#D-#H #N:#R:#T.#V", 0, "2024-03-05 14:07:09.042"},
		{"#B/#C/#G #M:#Q:#S #U", 0, "24/3/5 14:7:9 42"},
		{"#K #F #W, #O #L", 2, "Dienstag März  5, 2 PM"},
		{"#J #E #X #P #I", 1, "Tu Mar  3 02 2"},
	}
	for _, tt := range tests {
		if got := DT_TO_STRF(x, 42, tt.fmt, tt.lang); got != tt.want {
			t.Errorf("DT_TO_STRF(%q) = %q, want %q", tt.fmt, got, tt.want)
		}
	}
	got := FSTRING_TO_DT("2024-03-05 14:07:09", "#Y-#M-#D #h:#m:#s")
	if DT_TO_DWORD(got) != DT_TO_DWORD(x) {
		t.Errorf("FSTRING_TO_DT = %v", time.Time(got))
	}
	got = FSTRING_TO_DT("5. March 24", "#D. #N #Y")
	if DT_TO_DWORD(got) != DT_TO_DWORD(iec.DT(time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC))) {
		t.Errorf("FSTRING_TO_DT with month name = %v", time.Time(got))
	}
}

func TestMessages(t *testing.T) {
	now := time.Unix(100, 0)
	var m MESSAGE_4R
	m.INIT()
	m.M0, m.M1, m.M2, m.M3 = "zero", "one", "two", "three"
	m.CLK = false
	m.Execute(now)
	m.CLK = true
	m.Execute(now)
	if m.MX != "one" || !m.TR {
		t.Fatalf("MESSAGE_4R edge = %q", m.MX)
	}
	now = now.Add(3 * time.Second)
	m.Execute(now)
	m.Execute(now)
	if m.MX != "two" {
		t.Fatalf("MESSAGE_4R after T1 = %q", m.MX)
	}
	m.ENQ = false
	m.Execute(now)
	if m.MX != "" || m.MN != 0 {
		t.Fatalf("MESSAGE_4R off = %q", m.MX)
	}

	var m8 MESSAGE_8
	m8.S2, m8.S5 = "two", "five"
	m8.IN5, m8.IN2 = true, true
	m8.Execute(now)
	if m8.M != "two" {
		t.Fatalf("MESSAGE_8 = %q", m8.M)
	}

	text := iec.STRING("abcdef")
	var tk TICKER
	tk.TEXT = &text
	tk.N = 3
	tk.PT = iec.TIME(time.Second)
	tk.Execute(now)
	if tk.DISPLAY != "abc" {
		t.Fatalf("TICKER = %q", tk.DISPLAY)
	}
	for i := 0; i < 3; i++ {
		now = now.Add(600 * time.Millisecond)
		tk.Execute(now)
	}
	if tk.DISPLAY != "bcd" {
		t.Fatalf("TICKER after a second = %q", tk.DISPLAY)
	}
}
