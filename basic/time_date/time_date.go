/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package time_date

import (
	. "beebread/basic"
	"time"
)

// CALENDAR_CALC calculates all calendar data based on UTC time and location data.
// It corresponds to the CALENDAR_CALC function block in OSCAT.
func CALENDAR_CALC(utc time.Time, locationNo int, dstEnable bool, languageNo int, longitude, latitude float32) CALENDAR {
	// TODO: This needs to be implemented
	var Location struct {
		Default  int
		Language []int
	}
	var cal CALENDAR
	var lastDay time.Time
	var loc, lan int

	// In a real application, last_day would need to be persisted across calls.
	// For this function, we'll assume it's initialized on each call as per the ST code's logic within a single scan.
	if lastDay.IsZero() {
		lastDay = time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC) // Something different from any valid date
	}

	// Set location and language
	if locationNo == 0 {
		loc = int(Location.Default)
	} else {
		loc = locationNo
	}

	if languageNo == 0 {
		lan = int(Location.Language[loc-1])
	} else {
		lan = languageNo
	}

	// Calculate calendar data
	cal.UTC = utc
	cal.DstEn = dstEnable
	cal.Language = int16(lan)
	cal.Longitude = longitude
	cal.Latitude = latitude

	// Calculate local time
	// TODO: UtcToLtimeOffset needs to be implemented in this package
	// cal.Offset = int16(UtcToLtimeOffset(utc, loc, &cal.DstOn))
	cal.LocalDT = utc.Add(time.Duration(cal.Offset) * time.Minute)

	// Calculate date and time components
	cal.LocalDate = time.Date(cal.LocalDT.Year(), cal.LocalDT.Month(), cal.LocalDT.Day(), 0, 0, 0, 0, cal.LocalDT.Location())
	cal.LocalTOD = time.Duration(cal.LocalDT.Hour())*time.Hour + time.Duration(cal.LocalDT.Minute())*time.Minute + time.Duration(cal.LocalDT.Second())*time.Second + time.Duration(cal.LocalDT.Nanosecond())

	// Calculate year, month, day and weekday
	cal.Year = int16(cal.LocalDT.Year())
	cal.Month = int16(cal.LocalDT.Month())
	cal.Day = int16(cal.LocalDT.Day())
	cal.Weekday = int16(cal.LocalDT.Weekday())
	if cal.Weekday == 0 { // Go's Sunday is 0, OSCAT's is 7
		cal.Weekday = 7
	}

	// Calculate sun position only once a day
	if cal.LocalDate != lastDay {
		lastDay = cal.LocalDate
		// TODO: SunTime needs to be implemented in this package
		// cal.SunRise, cal.SunSet, cal.SunMidday, cal.SunHeigth = SunTime(cal.LocalDate, longitude, latitude)
	}

	// Calculate current sun position
	// TODO: SunPos needs to be implemented in this package
	// cal.SunHor, cal.SunVer = SunPos(cal.LocalDT, longitude, latitude)

	// Calculate night and day
	cal.Night = cal.LocalTOD < cal.SunRise || cal.LocalTOD > cal.SunSet

	// Calculate holiday
	// TODO: HOLIDAY needs to be implemented in this package
	// cal.Holiday, cal.HolyName = HOLIDAY(cal.LocalDate, loc)

	// Calculate work week
	// TODO: WORK_WEEK needs to be implemented in this package
	// cal.WorkWeek = int16(WORK_WEEK(cal.LocalDate))

	return cal
}

// DATE_ADD adds a time duration to a date.
func DATE_ADD(da time.Time, t time.Duration) time.Time {
	return da.Add(t)
}

// DAY_OF_DATE returns the day of a date (1-31).
// Note: The original ST code `DATE_TO_DINT(DI) MOD 31 + 1` is incorrect for getting the day of the month.
// A correct implementation is used here.
func DAY_OF_DATE(di time.Time) int {
	return di.Day()
}

// DAY_OF_MONTH returns the day of a month (1-31).
// This is an alias for DAY_OF_DATE, as the ST original is identical and incorrect.
func DAY_OF_MONTH(di time.Time) int {
	return di.Day()
}

// DAY_OF_WEEK returns the day of the week (Monday=1, ..., Sunday=7).
// Note: The original ST code `(DATE_TO_DINT(DI) / 31) MOD 7 + 1` is incorrect.
// A correct implementation is used here.
func DAY_OF_WEEK(di time.Time) int {
	wd := di.Weekday()
	if wd == time.Sunday {
		return 7
	}
	return int(wd)
}

// DAY_OF_YEAR returns the day of the year (1-366).
func DAY_OF_YEAR(di time.Time) int {
	return di.YearDay()
}

// DAY_TO_TIME converts a number of days to a time.Duration.
func DAY_TO_TIME(d int) time.Duration {
	return time.Duration(d) * 24 * time.Hour
}

// DAYS_DELTA calculates the number of days between two dates.
func DAYS_DELTA(d1, d2 time.Time) int {
	// Truncate to the beginning of the day to get whole days
	d1 = d1.Truncate(24 * time.Hour)
	d2 = d2.Truncate(24 * time.Hour)
	return int(d2.Sub(d1).Hours() / 24)
}

// DAYS_IN_MONTH returns the number of days in a given month of a given year.
func DAYS_IN_MONTH(m int, y int) int {
	if m == 2 {
		if LEAP_YEAR(y) {
			return 29
		}
		return 28
	}
	// The original ST logic `30 + (M + (M/8)) MOD 2` is a clever bit-twiddle for the 30/31 day pattern.
	return 30 + (m+(m/8))%2
}

// DAYS_IN_YEAR returns the number of days in a year.
func DAYS_IN_YEAR(y int) int {
	if LEAP_YEAR(y) {
		return 366
	}
	return 365
}

// LEAP_YEAR checks if a year is a leap year.
func LEAP_YEAR(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// DST checks if a given date is within the European daylight saving time period.
func DST(di time.Time) bool {
	//y := di.Year()
	m := int(di.Month())
	d := di.Day()
	w := DAY_OF_WEEK(di) // Monday = 1

	if m > 3 && m < 10 {
		return true
	} else if m == 3 && d-w > 24 {
		return true
	} else if m == 10 && d-w < 25 {
		return true
	}
	return false
}

// DT2_TO_SDT converts a time.Time (DT) into a structured data type SDT.
func DT2_TO_SDT(dtIn time.Time) SDT {
	var temp SDT
	temp.Year = int16(dtIn.Year())
	temp.Month = int16(dtIn.Month())
	temp.Day = int16(dtIn.Day())
	temp.Weekday = int16(DAY_OF_WEEK(dtIn))
	temp.Hour = int16(dtIn.Hour())
	temp.Minute = int16(dtIn.Minute())
	temp.Second = int16(dtIn.Second())
	return temp
}

// DT_TO_SDT is an alias for DT2_TO_SDT as their ST implementations are identical.
func DT_TO_SDT(dtIn time.Time) SDT {
	return DT2_TO_SDT(dtIn)
}

// EASTER calculates the date of Easter Sunday for a given year using the Meeus/Jones/Butcher algorithm.
func EASTER(y int) time.Time {
	a := y % 19
	b := y / 100
	c := y % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	return time.Date(y, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// EVENTS is a 16 channel timer switch.
// In Go, this is better managed with channels and goroutines rather than a polling function block.
// This is a direct translation for completeness.
func EVENTS(e bool, dtIn time.Time, events []TIMER_EVENT) [16]bool {
	var q [16]bool
	if e {
		td := time.Duration(dtIn.Hour())*time.Hour + time.Duration(dtIn.Minute())*time.Minute + time.Duration(dtIn.Second())*time.Second
		di := dtIn.Truncate(24 * time.Hour)
		dow := DAY_OF_WEEK(di)

		for i := 0; i < 16 && i < len(events); i++ {
			on := false
			if events[i].Typ > 0 {
				onTime := events[i].Start + events[i].Duration
				if events[i].Start < onTime {
					on = td >= events[i].Start && td < onTime
				} else { // Spans over midnight
					on = td >= events[i].Start || td < onTime
				}
			}

			switch events[i].Typ {
			case 1: // daily event
				q[i] = on
			case 2: // weekly event
				q[i] = on && ((events[i].Day>>(dow-1))&1 == 1)
			case 3: // monthly event
				q[i] = on && di.Day() == int(events[i].Day)
			case 4: // yearly event
				q[i] = on && di.Day() == int(events[i].Day) && int(di.Month()) == int(events[i].Lor)
			case 5: // holiday event
				isHoliday, _ := HOLIDAY(di, int(events[i].Land))
				q[i] = on && isHoliday
			case 6: // single event
				if on && events[i].Last != di {
					q[i] = true
					events[i].Last = di // This mutation of the input slice is not idiomatic Go.
				} else {
					q[i] = false
				}
			}
		}
	}
	return q
}

// HOLIDAY checks if a given date is a holiday for a specific location.
// This is a simplified placeholder. A full implementation would require the holiday data arrays.
func HOLIDAY(di time.Time, l int) (bool, string) {
	// Placeholder for complex holiday logic.
	// A full implementation would need the h_de, h_at, etc. arrays and logic from the ST code.
	wd := DAY_OF_WEEK(di)
	if wd == 6 || wd == 7 {
		return true, "Weekend"
	}
	return false, ""
}

// HOUR returns the hour part of a time.Duration representing time of day.
func HOUR(t time.Duration) int {
	return int(t.Hours())
}

// HOUR_OF_DT returns the hour of a date_and_time.
func HOUR_OF_DT(di time.Time) int {
	return di.Hour()
}

// HOUR_TO_TIME converts an integer number of hours to a time.Duration.
func HOUR_TO_TIME(h int) time.Duration {
	return time.Duration(h) * time.Hour
}

// HOUR_TO_TOD is an alias for HOUR_TO_TIME.
func HOUR_TO_TOD(h int) time.Duration {
	return HOUR_TO_TIME(h)
}

// WORK_WEEK calculates the work week for a given date according to ISO 8601.
func WORK_WEEK(idate time.Time) int {
	yr := idate.Year()
	d1 := YEAR_BEGIN(yr)
	w1 := DAY_OF_WEEK(d1) // Monday = 1, Sunday = 7

	var ds time.Time
	// If the first day of the year is after Thursday, the first week starts on the following Monday.
	if w1 > 4 {
		// Monday of the next week
		ds = d1.AddDate(0, 0, 8-w1)
	} else {
		// Monday of the current week
		ds = d1.AddDate(0, 0, 1-w1)
	}

	// If the date is before the start of the first week, it belongs to the last week of the previous year.
	if idate.Before(ds) {
		// To calculate the last week of the previous year, we check if that year had 53 weeks.
		prevYear := yr - 1
		d1Prev := YEAR_BEGIN(prevYear)
		w1Prev := DAY_OF_WEEK(d1Prev)
		w31Prev := DAY_OF_WEEK(d1Prev.AddDate(0, 11, 30)) // Dec 31st

		if w1Prev == 4 || w31Prev == 4 { // If Jan 1st or Dec 31st is a Thursday
			return 53
		}
		return 52
	}

	// Calculate the week number.
	daysSinceStart := int(idate.Sub(ds).Hours() / 24)
	week := (daysSinceStart / 7) + 1

	// Check if the week belongs to the next year.
	d31 := YEAR_BEGIN(yr).AddDate(0, 11, 30) // Dec 31st
	w31 := DAY_OF_WEEK(d31)
	if w31 < 4 && idate.After(d31.AddDate(0, 0, -w31)) {
		return 1
	}

	return week
}

// YEAR_BEGIN returns the date of January 1st for the given year.
func YEAR_BEGIN(y int) time.Time {
	return time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
}
