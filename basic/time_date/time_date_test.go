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

package time_date

import (
	gomath "math"
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

func date(y, m, d int) iec.DATE {
	return iec.DATE(time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC))
}

func dt(y, m, d, hh, mm, ss int) iec.DT {
	return iec.DT(time.Date(y, time.Month(m), d, hh, mm, ss, 0, time.UTC))
}

func tod(hh, mm, ss int) iec.TOD {
	return DWORD_TO_TOD(iec.DWORD(hh*3600000 + mm*60000 + ss*1000))
}

// TestDatesAgainstGo checks the date functions against Go's time package for
// every day OSCAT supports.
func TestDatesAgainstGo(t *testing.T) {
	for d := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC); d.Year() < 2100; d = d.AddDate(0, 0, 1) {
		id := iec.DATE(d)
		wd := int(d.Weekday())
		if wd == 0 {
			wd = 7
		}
		leap := d.Year()%4 == 0
		daysInMonth := time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		isoYear, isoWeek := d.ISOWeek()
		_ = isoYear
		checks := []struct {
			name      string
			got, want int
		}{
			{"YEAR_OF_DATE", int(YEAR_OF_DATE(id)), d.Year()},
			{"MONTH_OF_DATE", int(MONTH_OF_DATE(id)), int(d.Month())},
			{"DAY_OF_MONTH", int(DAY_OF_MONTH(id)), d.Day()},
			{"DAY_OF_WEEK", int(DAY_OF_WEEK(id)), wd},
			{"DAY_OF_YEAR", int(DAY_OF_YEAR(id)), d.YearDay()},
			{"DAYS_IN_MONTH", int(DAYS_IN_MONTH(id)), daysInMonth},
			{"SET_DATE", int(DATE_TO_DWORD(SET_DATE(iec.INT(d.Year()), iec.INT(d.Month()), iec.INT(d.Day())))), int(d.Unix())},
		}
		if d.Year() > 1970 && d.Year() < 2099 {
			checks = append(checks, struct {
				name      string
				got, want int
			}{"WORK_WEEK", int(WORK_WEEK(id)), isoWeek})
		}
		for _, c := range checks {
			if c.got != c.want {
				t.Fatalf("%s(%s) = %d, want %d", c.name, d.Format("2006-01-02"), c.got, c.want)
			}
		}
		if bool(LEAP_OF_DATE(id)) != leap {
			t.Fatalf("LEAP_OF_DATE(%s) = %v", d.Format("2006-01-02"), !leap)
		}
		if bool(LEAP_DAY(id)) != (d.Month() == 2 && d.Day() == 29) {
			t.Fatalf("LEAP_DAY(%s)", d.Format("2006-01-02"))
		}
		if d.Day() == 1 && d.Month() == 1 {
			if DATE_TO_DWORD(YEAR_BEGIN(iec.INT(d.Year()))) != iec.DWORD(d.Unix()) {
				t.Fatalf("YEAR_BEGIN(%d)", d.Year())
			}
			if DATE_TO_DWORD(YEAR_END(iec.INT(d.Year()))) != iec.DWORD(d.AddDate(1, 0, -1).Unix()) {
				t.Fatalf("YEAR_END(%d)", d.Year())
			}
			if DAYS_IN_YEAR(id) != iec.INT(SEL[int](iec.BOOL(leap), 365, 366)) {
				t.Fatalf("DAYS_IN_YEAR(%d)", d.Year())
			}
		}
		if d.Day() == 15 {
			if DATE_TO_DWORD(MONTH_BEGIN(id)) != iec.DWORD(d.AddDate(0, 0, -14).Unix()) {
				t.Fatalf("MONTH_BEGIN(%s)", d.Format("2006-01-02"))
			}
			if d.Year() < 2099 && DATE_TO_DWORD(MONTH_END(id)) != iec.DWORD(time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, time.UTC).Unix()) {
				t.Fatalf("MONTH_END(%s)", d.Format("2006-01-02"))
			}
		}
	}
}

func TestEaster(t *testing.T) {
	for y, want := range map[iec.INT]iec.DATE{
		2000: date(2000, 4, 23), 2019: date(2019, 4, 21), 2024: date(2024, 3, 31),
		2025: date(2025, 4, 20), 2038: date(2038, 4, 25), 1981: date(1981, 4, 19),
	} {
		if got := EASTER(y); DATE_TO_DWORD(got) != DATE_TO_DWORD(want) {
			t.Errorf("EASTER(%d) = %v, want %v", y, time.Time(got), time.Time(want))
		}
	}
}

func TestDST(t *testing.T) {
	tests := []struct {
		utc  iec.DT
		want iec.BOOL
	}{
		{dt(2024, 3, 31, 0, 59, 59), false},
		{dt(2024, 3, 31, 1, 0, 0), true},
		{dt(2024, 10, 27, 0, 59, 59), true},
		{dt(2024, 10, 27, 1, 0, 0), false},
		{dt(2025, 3, 30, 1, 0, 0), true},
		{dt(2025, 7, 1, 12, 0, 0), true},
		{dt(2025, 12, 1, 12, 0, 0), false},
	}
	for _, tt := range tests {
		if got := DST(tt.utc); got != tt.want {
			t.Errorf("DST(%v) = %v, want %v", time.Time(tt.utc), got, tt.want)
		}
	}
}

func TestTimeFunctions(t *testing.T) {
	x := dt(2024, 7, 1, 12, 34, 56)
	if got := UTC_TO_LTIME(x, true, 60); DT_TO_DWORD(got) != DT_TO_DWORD(dt(2024, 7, 1, 14, 34, 56)) {
		t.Errorf("UTC_TO_LTIME = %v", time.Time(got))
	}
	if got := UTC_TO_LTIME(x, false, -300); DT_TO_DWORD(got) != DT_TO_DWORD(dt(2024, 7, 1, 7, 34, 56)) {
		t.Errorf("UTC_TO_LTIME west = %v", time.Time(got))
	}
	if got := LTIME_TO_UTC(dt(2024, 7, 1, 14, 34, 56), true, 60); DT_TO_DWORD(got) != DT_TO_DWORD(x) {
		t.Errorf("LTIME_TO_UTC = %v", time.Time(got))
	}
	if HOUR_OF_DT(x) != 12 || MINUTE_OF_DT(x) != 34 || SECOND_OF_DT(x) != 56 {
		t.Error("HOUR_OF_DT, MINUTE_OF_DT, SECOND_OF_DT")
	}
	td := DT_TO_TOD(x)
	if HOUR(td) != 12 || MINUTE(td) != 34 || SECOND(td) != 56 {
		t.Error("HOUR, MINUTE, SECOND")
	}
	if got := SET_DT(2024, 7, 1, 12, 34, 56); DT_TO_DWORD(got) != DT_TO_DWORD(x) {
		t.Errorf("SET_DT = %v", time.Time(got))
	}
	s := DT_TO_SDT(x)
	if s.YEAR != 2024 || s.MONTH != 7 || s.DAY != 1 || s.WEEKDAY != 1 || s.HOUR != 12 || s.MINUTE != 34 || s.SECOND != 56 {
		t.Errorf("DT_TO_SDT = %+v", s)
	}
	if DT_TO_DWORD(SDT_TO_DT(s)) != DT_TO_DWORD(x) || DATE_TO_DWORD(SDT_TO_DATE(s)) != DATE_TO_DWORD(date(2024, 7, 1)) {
		t.Error("SDT_TO_DT, SDT_TO_DATE")
	}
	s2 := DT2_TO_SDT(date(2024, 7, 1), SET_TOD(12, 34, 56.789))
	if s2.MS != 789 || s2.SECOND != 56 || TOD_TO_DWORD(SDT_TO_TOD(s2)) != TOD_TO_DWORD(SET_TOD(12, 34, 56.789)) {
		t.Errorf("DT2_TO_SDT = %+v", s2)
	}
	if got := DATE_ADD(date(2024, 3, 15), 1, 1, 1, 1); DATE_TO_DWORD(got) != DATE_TO_DWORD(date(2025, 4, 23)) {
		t.Errorf("DATE_ADD = %v", time.Time(got))
	}
	if got := DATE_ADD(date(2024, 1, 15), -20, 0, 0, 0); DATE_TO_DWORD(got) != DATE_TO_DWORD(date(2023, 12, 26)) {
		t.Errorf("DATE_ADD back = %v", time.Time(got))
	}
	if DAYS_DELTA(date(2024, 1, 1), date(2024, 3, 1)) != 60 || DAYS_DELTA(date(2024, 3, 1), date(2024, 1, 1)) != -60 {
		t.Error("DAYS_DELTA")
	}
	if DAY_OF_DATE(date(1970, 1, 11)) != 10 {
		t.Error("DAY_OF_DATE")
	}
	if DAY_TO_TIME(1.5) != iec.TIME(36*time.Hour) || HOUR_TO_TIME(0.5) != iec.TIME(30*time.Minute) ||
		MINUTE_TO_TIME(1.5) != iec.TIME(90*time.Second) || SECOND_TO_TIME(0.25) != iec.TIME(250*time.Millisecond) {
		t.Error("DAY_TO_TIME, HOUR_TO_TIME, MINUTE_TO_TIME, SECOND_TO_TIME")
	}
	if TOD_TO_DWORD(HOUR_TO_TOD(12.5)) != 45000000 {
		t.Error("HOUR_TO_TOD")
	}
	if MULTIME(iec.TIME(time.Second), 2.5) != iec.TIME(2500*time.Millisecond) {
		t.Error("MULTIME")
	}
	if JD2000(dt(2000, 1, 2, 12, 0, 0)) != 1 {
		t.Error("JD2000")
	}
	if !LEAP_YEAR(2024) || LEAP_YEAR(2023) {
		t.Error("LEAP_YEAR")
	}
	if !TIMECHECK(tod(23, 0, 0), tod(22, 0, 0), tod(6, 0, 0)) || TIMECHECK(tod(12, 0, 0), tod(22, 0, 0), tod(6, 0, 0)) ||
		!TIMECHECK(tod(8, 0, 0), tod(8, 0, 0), tod(9, 0, 0)) || TIMECHECK(tod(9, 0, 0), tod(8, 0, 0), tod(9, 0, 0)) {
		t.Error("TIMECHECK")
	}
	if !PERIOD(date(2020, 12, 1), date(2024, 1, 5), date(2020, 1, 10)) || PERIOD(date(2020, 3, 1), date(2023, 2, 28), date(2020, 3, 31)) {
		t.Error("PERIOD")
	}
	var dp [4][2]iec.DATE
	dp[2] = [2]iec.DATE{date(2024, 5, 1), date(2024, 5, 31)}
	if !PERIOD2(dp, date(2024, 5, 10)) || PERIOD2(dp, date(2024, 6, 1)) {
		t.Error("PERIOD2")
	}
	if r := REFRACTION(0); gomath.Abs(float64(r)-0.48) > 0.02 {
		t.Errorf("REFRACTION(0) = %v", r)
	}
}

func TestSun(t *testing.T) {
	var s SUN_TIME
	s.INIT()
	s.LATITUDE, s.LONGITUDE = 52.52, 13.405
	s.UTC = date(2024, 6, 21)
	s.Execute(time.Time{})
	near := func(got iec.TOD, h, m int) bool {
		d := int(TOD_TO_DWORD(got)/60000) - (h*60 + m)
		return d >= -10 && d <= 10
	}
	if !near(s.SUN_RISE, 2, 43) || !near(s.SUN_SET, 19, 33) || !near(s.MIDDAY, 11, 8) {
		t.Errorf("SUN_TIME Berlin midsummer: rise %v set %v midday %v",
			TOD_TO_DWORD(s.SUN_RISE)/60000, TOD_TO_DWORD(s.SUN_SET)/60000, TOD_TO_DWORD(s.MIDDAY)/60000)
	}
	if gomath.Abs(float64(s.SUN_DECLINATION)-60.9) > 1 {
		t.Errorf("SUN_TIME declination = %v", s.SUN_DECLINATION)
	}

	var p SUN_POS
	p.LATITUDE, p.LONGITUDE = 52.52, 13.405
	p.UTC = dt(2024, 6, 21, 11, 8, 0)
	p.Execute(time.Time{})
	if gomath.Abs(float64(p.H)-60.9) > 1 || gomath.Abs(float64(p.B)-180) > 3 || p.HR <= p.H {
		t.Errorf("SUN_POS = B %v H %v HR %v", p.B, p.H, p.HR)
	}
}

func TestHoliday(t *testing.T) {
	list := HOLIDAY_DE
	list[15].USE = -3 // Buss und Bettag, the wednesday before 23.11.
	var h HOLIDAY
	h.HOLIDAYS = &list
	h.SATURDAY = true
	tests := []struct {
		d    iec.DATE
		y    iec.BOOL
		name iec.STRING
	}{
		{date(2024, 3, 29), true, "Karfreitag"},
		{date(2024, 5, 9), true, "Christi Himmelfahrt"},
		{date(2024, 10, 3), true, "Tag der Deutschen Einheit"},
		{date(2024, 11, 20), true, "Buss und Bettag"},
		{date(2024, 6, 22), true, "Samstag"},
		{date(2024, 6, 20), false, ""},
		{date(2024, 12, 25), true, "1. Weihnachtstag"},
	}
	for _, tt := range tests {
		h.DATE_IN = tt.d
		h.Execute(time.Time{})
		if h.Y != tt.y || h.NAME != tt.name {
			t.Errorf("HOLIDAY(%v) = %v %q, want %v %q", time.Time(tt.d).Format("2006-01-02"), h.Y, h.NAME, tt.y, tt.name)
		}
	}

	var elist [50]HOLIDAY_DATA
	elist[3] = HOLIDAY_DATA{NAME: "Oktoberfest", DAY: 21, MONTH: 9, USE: 16}
	var e EVENTS
	e.ELIST = &elist
	e.ENA = true
	e.DATE_IN = date(2024, 10, 1)
	e.Execute(time.Time{})
	if !e.Y || e.NAME != "Oktoberfest" {
		t.Errorf("EVENTS = %v %q", e.Y, e.NAME)
	}
	e.DATE_IN = date(2024, 10, 7)
	e.Execute(time.Time{})
	if e.Y {
		t.Error("EVENTS after the event")
	}
}

func TestCalendarCalc(t *testing.T) {
	cal := CALENDAR{UTC: dt(2024, 12, 25, 10, 0, 0), OFFSET: 60, DST_EN: true, LATITUDE: 52.52, LONGITUDE: 13.405}
	list := HOLIDAY_DE
	var c CALENDAR_CALC
	c.XCAL, c.HOLIDAYS = &cal, &list
	c.INIT()
	c.SPE = true
	c.Execute(time.Time{})
	if cal.YEAR != 2024 || cal.MONTH != 12 || cal.DAY != 25 || cal.WEEKDAY != 3 || !cal.HOLIDAY ||
		cal.HOLY_NAME != "1. Weihnachtstag" || cal.WORK_WEEK != 52 || cal.DST_ON {
		t.Errorf("CALENDAR_CALC = %+v", cal)
	}
	if DT_TO_DWORD(cal.LOCAL_DT) != DT_TO_DWORD(dt(2024, 12, 25, 11, 0, 0)) {
		t.Errorf("LOCAL_DT = %v", time.Time(cal.LOCAL_DT))
	}
	// Sun rise in Berlin on christmas is about 08:15 local time.
	if rise := TOD_TO_DWORD(cal.SUN_RISE) / 60000; rise < 8*60 || rise > 8*60+30 {
		t.Errorf("SUN_RISE = %v min", rise)
	}
	if cal.SUN_VER < 5 || cal.SUN_VER > 20 {
		t.Errorf("SUN_VER = %v", cal.SUN_VER)
	}
}

func TestRTC(t *testing.T) {
	now := time.Unix(5000, 0)
	var r RTC_2
	r.SDT = dt(2024, 7, 1, 12, 0, 0)
	r.DEN = true
	r.OFS = 60
	r.Execute(now)
	for i := 0; i < 25; i++ {
		now = now.Add(100 * time.Millisecond)
		r.Execute(now)
	}
	if DT_TO_DWORD(r.UDT) != DT_TO_DWORD(dt(2024, 7, 1, 12, 0, 2)) || r.XMS != 500 || !r.DSO {
		t.Errorf("RTC_2 = %v %v %v", time.Time(r.UDT), r.XMS, r.DSO)
	}
	if DT_TO_DWORD(r.LOCAL_DT) != DT_TO_DWORD(dt(2024, 7, 1, 14, 0, 2)) {
		t.Errorf("RTC_2 LOCAL_DT = %v", time.Time(r.LOCAL_DT))
	}
}

// dcfMinute returns the 59 bits DCF77 sends for the minute of the local
// time mez, in summer time.
func dcfMinute(mez time.Time) [59]bool {
	var b [59]bool
	bcd := func(from int, v int, weights ...int) {
		for i := len(weights) - 1; i >= 0; i-- {
			if v >= weights[i] {
				b[from+i] = true
				v -= weights[i]
			}
		}
	}
	par := func(from, to int) {
		p := false
		for i := from; i < to; i++ {
			p = p != b[i]
		}
		b[to] = p
	}
	b[17], b[20] = true, true
	bcd(21, mez.Minute(), 1, 2, 4, 8, 10, 20, 40)
	par(21, 28)
	bcd(29, mez.Hour(), 1, 2, 4, 8, 10, 20)
	par(29, 35)
	bcd(36, mez.Day(), 1, 2, 4, 8, 10, 20)
	wd := int(mez.Weekday())
	if wd == 0 {
		wd = 7
	}
	bcd(42, wd, 1, 2, 4)
	bcd(45, int(mez.Month()), 1, 2, 4, 8, 10)
	bcd(50, mez.Year()%100, 1, 2, 4, 8, 10, 20, 40, 80)
	par(36, 58)
	return b
}

func TestDCF77(t *testing.T) {
	var d DCF77
	d.INIT()
	now := time.Unix(10000, 0)
	scan := func(rec bool, dur time.Duration) {
		d.REC = iec.BOOL(rec)
		d.Execute(now)
		now = now.Add(dur)
	}
	// Each second the signal drops for 100 ms (0) or 200 ms (1); the 59th
	// second has no drop, which marks the minute.
	send := func(bits [59]bool) {
		for i := 0; i < 59; i++ {
			low := 100 * time.Millisecond
			if bits[i] {
				low = 200 * time.Millisecond
			}
			scan(false, low)
			scan(true, time.Second-low)
		}
		now = now.Add(time.Second)
	}
	// The minute before the transmission starts.
	scan(true, 1900*time.Millisecond)
	mez := time.Date(2024, 7, 1, 14, 30, 0, 0, time.UTC)
	send(dcfMinute(mez))
	send(dcfMinute(mez.Add(time.Minute)))
	scan(false, 0)
	if !d.TP || d.ERROR {
		t.Fatalf("DCF77 after two minutes: TP %v ERROR %v", d.TP, d.ERROR)
	}
	want := DT_TO_DWORD(dt(2024, 7, 1, 12, 31, 0))
	if DT_TO_DWORD(d.RTC) != want || !d.DS || !d.SYNC {
		t.Fatalf("DCF77 RTC = %v DS %v SYNC %v", time.Time(d.RTC), d.DS, d.SYNC)
	}
	// RTC1 is one hour ahead plus summer time.
	if DT_TO_DWORD(d.RTC1) != want+7200 {
		t.Fatalf("DCF77 RTC1 = %v", time.Time(d.RTC1))
	}
	now = now.Add(1500 * time.Millisecond)
	d.REC = false
	d.Execute(now)
	if DT_TO_DWORD(d.RTC) != want+1 || d.WDAY != 1 {
		t.Fatalf("DCF77 free running RTC = %v WDAY %v", time.Time(d.RTC), d.WDAY)
	}
}
