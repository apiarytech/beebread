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
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

func TestDates(t *testing.T) {
	dt := iec.DT(time.Date(2024, 7, 1, 12, 30, 15, 0, time.UTC))
	if DT_TO_DWORD(dt) != 1719837015 {
		t.Errorf("DT_TO_DWORD = %v", DT_TO_DWORD(dt))
	}
	if DT_TO_DWORD(iec.DT{}) != 0 || DATE_TO_DWORD(iec.DATE{}) != 0 {
		t.Error("the zero DT and DATE are not 0")
	}
	if DATE_TO_DWORD(DT_TO_DATE(dt)) != 1719792000 {
		t.Errorf("DT_TO_DATE = %v", time.Time(DT_TO_DATE(dt)))
	}
	if TOD_TO_DWORD(DT_TO_TOD(dt)) != 45015000 {
		t.Errorf("DT_TO_TOD = %v", TOD_TO_DWORD(DT_TO_TOD(dt)))
	}
	if DT_TO_DWORD(DATE_TO_DT(DT_TO_DATE(dt))) != 1719792000 {
		t.Error("DATE_TO_DT")
	}
	if DWORD_TO_DT(1719837015) != dt {
		t.Error("DWORD_TO_DT")
	}
	if TIME_TO_DWORD(iec.TIME(1500*time.Millisecond)) != 1500 || DWORD_TO_TIME(250) != iec.TIME(250*time.Millisecond) {
		t.Error("TIME_TO_DWORD, DWORD_TO_TIME")
	}
	if TOD_TO_TIME(DWORD_TO_TOD(90000000)) != iec.TIME(time.Hour) {
		t.Error("DWORD_TO_TOD wraps at midnight")
	}
	if REAL_TO_TIME(1.5) != iec.TIME(2*time.Millisecond) || TIME_TO_REAL(iec.TIME(time.Second)) != 1000 {
		t.Error("REAL_TO_TIME, TIME_TO_REAL")
	}
}

func TestClock(t *testing.T) {
	start := time.Unix(1000, 0)
	SetEpoch(start)
	SetClock(func() time.Time { return start.Add(1500 * time.Millisecond) })
	defer SetClock(time.Now)
	if T_PLC_MS() != 1500 || T_PLC_US() != 1500000 || PLC_MS(start.Add(-time.Millisecond)) != 0xFFFFFFFF {
		t.Errorf("T_PLC_MS = %v, T_PLC_US = %v", T_PLC_MS(), T_PLC_US())
	}
	if !Now().Equal(start.Add(1500 * time.Millisecond)) {
		t.Error("Now")
	}
}

func TestNumbers(t *testing.T) {
	if REAL_TO_INT(2.5) != 2 || REAL_TO_INT(3.5) != 4 || REAL_TO_DINT(-1.6) != -2 || TRUNC(-1.6) != -1 {
		t.Error("REAL_TO_INT, REAL_TO_DINT, TRUNC")
	}
	if big := uint64(1e10); REAL_TO_DWORD(1e10) != iec.DWORD(big) {
		t.Error("REAL_TO_DWORD cuts to 32 bits")
	}
	if LIMIT(0, 5, 3) != 3 || LIMIT(0, -1, 3) != 0 || SEL(true, 1, 2) != 2 || BOOL_TO_INT(true) != 1 {
		t.Error("LIMIT, SEL, BOOL_TO_INT")
	}
	if SHL(iec.BYTE(1), 8) != 0 || SHR(iec.DWORD(8), -1) != 0 || ROL(iec.BYTE(0x81), 1) != 0x03 ||
		ROR(iec.WORD(1), 1) != 0x8000 || !BIT(iec.DWORD(4), 2) {
		t.Error("SHL, SHR, ROL, ROR, BIT")
	}
	if SQRT(16) != 4 || EXPT(2, 10) != 1024 || ABS(-2) != 2 || LOG(100) != 2 {
		t.Error("SQRT, EXPT, ABS, LOG")
	}
}

func TestStrings(t *testing.T) {
	s := iec.STRING("Grüße")
	checks := []struct {
		name      string
		got, want iec.STRING
	}{
		{"LEFT", LEFT(s, 3), "Grü"},
		{"RIGHT", RIGHT(s, 2), "ße"},
		{"MID", MID(s, 2, 3), "üß"},
		{"MID past end", MID(s, 10, 4), "ße"},
		{"INSERT", INSERT("abc", "XY", 1), "aXYbc"},
		{"DELETE", DELETE("abcdef", 2, 3), "abef"},
		{"REPLACE", REPLACE("abcdef", "XY", 3, 2), "aXYef"},
		{"CONCAT", CONCAT("a", "b", "c"), "abc"},
		{"STR", STR([]iec.BYTE{'a', 228, 0, 'b'}), "aä"},
		{"INT_TO_STRING", INT_TO_STRING(-42), "-42"},
		{"DWORD_TO_STRING", DWORD_TO_STRING(4000000000), "4000000000"},
		{"REAL_TO_STRING", REAL_TO_STRING(2.5), "2.5"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
	if LEN(s) != 5 || FIND(s, "ße") != 4 || FIND(s, "x") != 0 {
		t.Error("LEN, FIND")
	}
	if c := CHARS("aä€"); len(c) != 3 || c[1] != 228 || c[2] != '?' {
		t.Errorf("CHARS = %v", c)
	}
	if STRING_TO_INT(" -12abc") != -12 || STRING_TO_DINT("x") != 0 || STRING_TO_REAL("3.5kg") != 3.5 {
		t.Error("STRING_TO_INT, STRING_TO_DINT, STRING_TO_REAL")
	}
}

func TestGlobals(t *testing.T) {
	if SETUP.CHARNAMES[0][:10] != ";\"&quot;&&" {
		t.Errorf("CHARNAMES[0] = %q", SETUP.CHARNAMES[0][:10])
	}
	if LANGUAGE.MONTHS[1][2] != "März" || MATH.FACTS[12] != 479001600 || LOCATION.LANGUAGE[2] != 3 {
		t.Error("constants")
	}
}
