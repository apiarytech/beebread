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

// Package time_date is the port of the OSCAT BASIC time and date functions.
// As in OSCAT, the calculations work on dates from 1970 to 2099.
package time_date

import (
	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// A DATE, TOD or DT compares by its seconds or milliseconds, as in OSCAT.
func dw(d iec.DATE) iec.DWORD { return DATE_TO_DWORD(d) }
func dtw(d iec.DT) iec.DWORD  { return DT_TO_DWORD(d) }
func tw(t iec.TOD) iec.DWORD  { return TOD_TO_DWORD(t) }

// intToDword converts an INT to a DWORD as INT_TO_DWORD does: a negative
// value is extended with its sign.
func intToDword(i iec.INT) iec.DWORD { return iec.DWORD(int32(i)) }

// DATE_ADD adds days d, weeks w, months m and years y to a date. Negative
// values subtract.
func DATE_ADD(idate iec.DATE, d, w, m, y iec.INT) iec.DATE {
	out := DWORD_TO_DATE(dw(idate) + intToDword(d+w*7)*86400)
	yr := y + YEAR_OF_DATE(out)
	mo := m + MONTH_OF_DATE(out)
	dm := DAY_OF_MONTH(out)
	for mo > 12 {
		mo -= 12
		yr++
	}
	for mo < 1 {
		mo += 12
		yr--
	}
	return SET_DATE(yr, mo, dm)
}

// DAY_OF_DATE returns the days since 1970-01-01.
func DAY_OF_DATE(idate iec.DATE) iec.DINT {
	return iec.DINT(dw(idate) / 86400)
}

// DAY_OF_MONTH returns the day of the month of a date.
func DAY_OF_MONTH(idate iec.DATE) iec.INT {
	ofs := func(month int) iec.INT { return SETUP.MTH_OFS[month-1] }
	day := DAY_OF_YEAR(idate)
	leap := BOOL_TO_INT(LEAP_OF_DATE(idate))
	day -= leap
	switch {
	case day > ofs(9):
		switch {
		case day > ofs(12):
			day -= ofs(12)
		case day > ofs(11):
			day -= ofs(11)
		case day > ofs(10):
			day -= ofs(10)
		default:
			day -= ofs(9)
		}
	case day > ofs(5):
		switch {
		case day > ofs(8):
			day -= ofs(8)
		case day > ofs(7):
			day -= ofs(7)
		case day > ofs(6):
			day -= ofs(6)
		default:
			day -= ofs(5)
		}
	case day > ofs(3):
		if day > ofs(4) {
			day -= ofs(4)
		} else {
			day -= ofs(3)
		}
	default:
		// January or February: the leap day counts again.
		day += leap
		if day > ofs(2) {
			day -= ofs(2)
		}
	}
	return day
}

// DAY_OF_WEEK returns the weekday of a date according to ISO 8601,
// monday = 1 ... sunday = 7.
func DAY_OF_WEEK(idate iec.DATE) iec.INT {
	return iec.INT((dw(idate)/86400+3)%7) + 1
}

// DAY_OF_YEAR returns the day of the year of a date, 1..366.
func DAY_OF_YEAR(idate iec.DATE) iec.INT {
	return iec.INT((dw(idate)-dw(YEAR_BEGIN(YEAR_OF_DATE(idate))))/86400) + 1
}

// DAY_TO_TIME converts a number of days to a TIME.
func DAY_TO_TIME(in iec.REAL) iec.TIME {
	return DWORD_TO_TIME(REAL_TO_DWORD(in * 86400000.0))
}

// DAYS_DELTA returns the days from date1 to date2.
func DAYS_DELTA(date1, date2 iec.DATE) iec.DINT {
	if dw(date1) > dw(date2) {
		return -iec.DINT((dw(date1) - dw(date2)) / 86400)
	}
	return iec.DINT((dw(date2) - dw(date1)) / 86400)
}

// DAYS_IN_MONTH returns the number of days of the month of a date.
func DAYS_IN_MONTH(idate iec.DATE) iec.INT {
	d := DAY_OF_YEAR(idate)
	in := func(lo, hi iec.INT) bool { return d >= lo && d <= hi }
	if LEAP_OF_DATE(idate) {
		switch {
		case in(32, 60):
			return 29
		case in(92, 121), in(153, 182), in(245, 274), in(306, 335):
			return 30
		}
		return 31
	}
	switch {
	case in(32, 59):
		return 28
	case in(91, 120), in(152, 181), in(244, 273), in(305, 334):
		return 30
	}
	return 31
}

// DAYS_IN_YEAR returns 366 for a date in a leap year and 365 otherwise.
func DAYS_IN_YEAR(idate iec.DATE) iec.INT {
	if LEAP_OF_DATE(idate) {
		return 366
	}
	return 365
}

// DST reports whether European daylight saving time is on at the time utc.
// It starts at 01:00 UTC on the last sunday of march and ends at 01:00 UTC
// on the last sunday of october.
func DST(utc iec.DT) iec.BOOL {
	yr := YEAR_OF_DATE(DT_TO_DATE(utc))
	ltc := dtw(utc)
	idate := dtw(SET_DT(yr, 3, 31, 1, 0, 0))
	yr4 := (5*intToDword(yr))>>2 + 1
	return idate-((yr4+3)%7)*86400 <= ltc && idate+(214-yr4%7)*86400 > ltc
}

// DT2_TO_SDT converts a date and a time of day to a structured date and
// time.
func DT2_TO_SDT(di iec.DATE, ti iec.TOD) SDT {
	t := tw(ti)
	return SDT{
		YEAR:    YEAR_OF_DATE(di),
		MONTH:   MONTH_OF_DATE(di),
		DAY:     DAY_OF_MONTH(di),
		WEEKDAY: DAY_OF_WEEK(di),
		MS:      iec.INT(t % 1000),
		SECOND:  iec.INT(t / 1000 % 60),
		MINUTE:  iec.INT(t / 60000 % 60),
		HOUR:    iec.INT(t / 3600000),
	}
}

// DT_TO_SDT converts a date and time to a structured date and time.
func DT_TO_SDT(dti iec.DT) SDT {
	tmp := DT_TO_DATE(dti)
	tdt := dtw(dti) - dw(tmp)
	return SDT{
		YEAR:    YEAR_OF_DATE(tmp),
		MONTH:   MONTH_OF_DATE(tmp),
		DAY:     DAY_OF_MONTH(tmp),
		WEEKDAY: DAY_OF_WEEK(tmp),
		SECOND:  iec.INT(tdt % 60),
		MINUTE:  iec.INT(tdt / 60 % 60),
		HOUR:    iec.INT(tdt / 3600),
	}
}

// EASTER returns the date of easter sunday of a year.
func EASTER(year iec.INT) iec.DATE {
	b := (204 - 11*(year%19)) % 30
	if b > 27 {
		b--
	}
	c := (year + iec.INT(uint16(year)>>2) + b - 13) % 7
	oday := 28 + b - c
	if oday > 33 {
		return SET_DATE(year, 4, oday-31)
	}
	return SET_DATE(year, 3, oday)
}

// HOUR returns the hour of a time of day.
func HOUR(itod iec.TOD) iec.INT {
	return iec.INT(tw(itod) / 3600000)
}

// HOUR_OF_DT returns the hour of a date and time.
func HOUR_OF_DT(xdt iec.DT) iec.INT {
	return iec.INT(dtw(xdt) % 86400 / 3600)
}

// HOUR_TO_TIME converts a number of hours to a TIME.
func HOUR_TO_TIME(in iec.REAL) iec.TIME {
	return DWORD_TO_TIME(REAL_TO_DWORD(in * 3600000))
}

// HOUR_TO_TOD converts a number of hours to a time of day.
func HOUR_TO_TOD(in iec.REAL) iec.TOD {
	return DWORD_TO_TOD(REAL_TO_DWORD(in * 3600000))
}

// JD2000 returns the days since 2000-01-01 12:00, the astronomic Julian
// date of 2000.
func JD2000(dti iec.DT) iec.REAL {
	return iec.REAL(dtw(dti)-946728000) / 86400.0
}

// LEAP_DAY reports whether a date is a leap day, the 29th of february.
func LEAP_DAY(idate iec.DATE) iec.BOOL {
	return dw(idate)%126230400 == 68169600
}

// LEAP_OF_DATE reports whether a date is in a leap year.
func LEAP_OF_DATE(idate iec.DATE) iec.BOOL {
	return ((dw(idate)+43200)/31557600)<<30 == 0x80000000
}

// LEAP_YEAR reports whether a year is a leap year.
func LEAP_YEAR(yr iec.INT) iec.BOOL {
	return uint16(yr)<<14 == 0
}

// LTIME_TO_UTC calculates the world time UTC from the local time ltime of a
// time zone offset in minutes, and one hour less if dst is true.
func LTIME_TO_UTC(ltime iec.DT, dst iec.BOOL, timeZoneOffset iec.INT) iec.DT {
	out := dtw(ltime) - intToDword(timeZoneOffset)*60
	if dst {
		out -= 3600
	}
	return DWORD_TO_DT(out)
}

// MINUTE returns the minutes of a time of day.
func MINUTE(itod iec.TOD) iec.INT {
	return iec.INT(tw(itod)/60000 - tw(itod)/3600000*60)
}

// MINUTE_OF_DT returns the minute of the hour of a date and time.
func MINUTE_OF_DT(xdt iec.DT) iec.INT {
	return iec.INT(dtw(xdt)%3600) / 60
}

// MINUTE_TO_TIME converts a number of minutes to a TIME.
func MINUTE_TO_TIME(in iec.REAL) iec.TIME {
	return DWORD_TO_TIME(REAL_TO_DWORD(in * 60000.0))
}

// MONTH_BEGIN returns the first day of the month of a date.
func MONTH_BEGIN(idate iec.DATE) iec.DATE {
	return DWORD_TO_DATE(dw(idate) - intToDword(DAY_OF_MONTH(idate)-1)*86400)
}

// MONTH_END returns the last day of the month of a date.
func MONTH_END(idate iec.DATE) iec.DATE {
	return DWORD_TO_DATE(dw(SET_DATE(YEAR_OF_DATE(idate), MONTH_OF_DATE(idate)+1, 1)) - 86400)
}

// MONTH_OF_DATE returns the month of a date.
func MONTH_OF_DATE(idate iec.DATE) iec.INT {
	m := DAY_OF_YEAR(idate)
	switch {
	case m < 32:
		return 1
	case bool(LEAP_OF_DATE(idate)):
		return (m*53 + 1668) / 1623
	}
	return (m*53 + 1700) / 1620
}

// MULTIME multiplies a TIME by a real number.
func MULTIME(t iec.TIME, m iec.REAL) iec.TIME {
	return DWORD_TO_TIME(REAL_TO_DWORD(iec.REAL(TIME_TO_DWORD(t)) * m))
}

// PERIOD reports whether the date dx is in the period from d1 to d2,
// whatever the years of the dates. A period may span the new year.
func PERIOD(d1, dx, d2 iec.DATE) iec.BOOL {
	day1 := DAY_OF_YEAR(d1)
	day2 := DAY_OF_YEAR(d2)
	dayx := DAY_OF_YEAR(dx)
	if !LEAP_OF_DATE(dx) && dayx > 58 {
		dayx++
	}
	if !LEAP_OF_DATE(d1) && day1 > 58 {
		day1++
	}
	if !LEAP_OF_DATE(d2) && day2 > 58 {
		day2++
	}
	if day2 < day1 {
		return dayx <= day2 || dayx >= day1
	}
	return dayx >= day1 && dayx <= day2
}

// PERIOD2 reports whether the date dx is in one of 4 periods, from dp[i][0]
// to dp[i][1].
func PERIOD2(dp [4][2]iec.DATE, dx iec.DATE) iec.BOOL {
	x := dw(dx)
	for _, p := range dp {
		if x >= dw(p[0]) && x <= dw(p[1]) {
			return true
		}
	}
	return false
}

// SDT_TO_DATE converts a structured date and time to its date.
func SDT_TO_DATE(dti SDT) iec.DATE {
	return SET_DATE(dti.YEAR, dti.MONTH, dti.DAY)
}

// SDT_TO_DT converts a structured date and time to a date and time.
func SDT_TO_DT(dti SDT) iec.DT {
	return SET_DT(dti.YEAR, dti.MONTH, dti.DAY, dti.HOUR, dti.MINUTE, dti.SECOND)
}

// SDT_TO_TOD converts a structured date and time to its time of day.
func SDT_TO_TOD(dti SDT) iec.TOD {
	return DWORD_TO_TOD(intToDword(dti.HOUR)*3600000 + intToDword(dti.MINUTE)*60000 +
		intToDword(dti.SECOND)*1000 + intToDword(dti.MS))
}

// SECOND returns the seconds and milliseconds of a time of day.
func SECOND(itod iec.TOD) iec.REAL {
	return iec.REAL(tw(itod)-tw(itod)/60000*60000) / 1000.0
}

// SECOND_OF_DT returns the second of the minute of a date and time.
func SECOND_OF_DT(xdt iec.DT) iec.INT {
	return iec.INT(dtw(xdt) % 60)
}

// SECOND_TO_TIME converts a number of seconds to a TIME.
func SECOND_TO_TIME(in iec.REAL) iec.TIME {
	return DWORD_TO_TIME(REAL_TO_DWORD(in * 1000.0))
}

// setDateOfs are the days of the year before each month.
var setDateOfs = [12]iec.INT{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}

// SET_DATE creates a date from a year, month and day. A month outside
// 1..12, which OSCAT reads outside its table, counts on into the years
// before or after: month 13 is january of the next year, as MONTH_END
// needs.
func SET_DATE(year, month, day iec.INT) iec.DATE {
	for month > 12 {
		month -= 12
		year++
	}
	for month < 1 {
		month += 12
		year--
	}
	ofs := setDateOfs[month-1]
	base := (intToDword(year)*1461 - 2878169) >> 2
	if month > 2 && uint16(year)<<14 == 0 {
		// One more day in a leap year.
		return DWORD_TO_DATE((intToDword(ofs+day) + base) * 86400)
	}
	return DWORD_TO_DATE((intToDword(ofs+day-1) + base) * 86400)
}

// SET_DT creates a date and time from a year, month, day, hour, minute and
// second.
func SET_DT(year, month, day, hour, minute, second iec.INT) iec.DT {
	return DWORD_TO_DT(dw(SET_DATE(year, month, day)) + intToDword(second) +
		intToDword(minute)*60 + intToDword(hour)*3600)
}

// SET_TOD creates a time of day from an hour, minute and second.
func SET_TOD(hour, minute iec.INT, second iec.REAL) iec.TOD {
	return DWORD_TO_TOD(REAL_TO_DWORD(second*1000.0) + intToDword(minute)*60000 + intToDword(hour)*3600000)
}

// TIMECHECK reports whether the time of day td is from start to before
// stop. A period over midnight has a start later than its stop.
func TIMECHECK(td, start, stop iec.TOD) iec.BOOL {
	t, a, b := tw(td), tw(start), tw(stop)
	if b < a {
		return a <= t || t < b
	}
	return a <= t && t < b
}

// UTC_TO_LTIME calculates the local time from the world time utc for a time
// zone offset in minutes, and one hour more during daylight saving time if
// dstEnable is true.
func UTC_TO_LTIME(utc iec.DT, dstEnable iec.BOOL, timeZoneOffset iec.INT) iec.DT {
	tmp := iec.DINT(timeZoneOffset)*60 + iec.DINT(BOOL_TO_INT(dstEnable && DST(utc)))*3600
	return DWORD_TO_DT(dtw(utc) + iec.DWORD(tmp))
}

// WORK_WEEK returns the calendar week of a date according to ISO 8601.
func WORK_WEEK(idate iec.DATE) iec.INT {
	yr := YEAR_OF_DATE(idate)
	d1 := YEAR_BEGIN(yr)
	w1 := DAY_OF_WEEK(d1)
	// The monday of the last week of the year before: if january 1st is
	// after thursday the week starts on the monday before it, otherwise
	// two mondays before.
	var ds iec.DWORD
	if w1 < 5 {
		ds = dw(d1) - intToDword(w1+6)*86400
	} else {
		ds = dw(d1) - intToDword(w1-1)*86400
	}
	week := iec.INT((dw(idate) - ds) / 604800)
	if week == 0 {
		// The last week of the year before is 53 if its january 1st or
		// december 31st was a thursday.
		var w31, w01 iec.INT
		if w1 > 1 {
			w31 = w1 - 1
		} else {
			w31 = 7
		}
		if LEAP_YEAR(yr-1) && w31 > 1 {
			w01 = w31 - 1
		}
		return 52 + BOOL_TO_INT(w31 == 4 || w01 == 4)
	}
	var w31 iec.INT
	if LEAP_YEAR(yr) {
		if w1 < 7 {
			w31 = w1 + 1
		} else {
			w31 = 1
		}
	} else {
		w31 = w1
	}
	wm := 52 + BOOL_TO_INT(w31 == 4 || w1 == 4)
	if week > wm {
		return 1
	}
	return week
}

// YEAR_BEGIN returns january 1st of a year.
func YEAR_BEGIN(y iec.INT) iec.DATE {
	return DWORD_TO_DATE((intToDword(y)*1461 - 2878169) >> 2 * 86400)
}

// YEAR_END returns december 31st of a year.
func YEAR_END(y iec.INT) iec.DATE {
	return DWORD_TO_DATE((intToDword(y)*1461 - 2876712) >> 2 * 86400)
}

// YEAR_OF_DATE returns the year of a date.
func YEAR_OF_DATE(idate iec.DATE) iec.INT {
	return iec.INT((dw(idate)+43200)/31557600 + 1970)
}
