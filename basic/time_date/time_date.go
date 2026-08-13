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

package basic

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
)

// CalendarCalc calculates all calendar data based on UTC time and location data.
// It corresponds to the CALENDAR_CALC function block in OSCAT.
func CalendarCalc(utc time.Time, locationNo int, dstEnable bool, languageNo int, longitude, latitude float32) Calendar {
	var cal Calendar
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
	cal.Offset = int16(UtcToLtimeOffset(utc, loc, &cal.DstOn))
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
		cal.SunRise, cal.SunSet, cal.SunMidday, cal.SunHeigth = SunTime(cal.LocalDate, longitude, latitude)
	}

	// Calculate current sun position
	cal.SunHor, cal.SunVer = SunPos(cal.LocalDT, longitude, latitude)

	// Calculate night and day
	cal.Night = cal.LocalTOD < cal.SunRise || cal.LocalTOD > cal.SunSet

	// Calculate holiday
	cal.Holiday, cal.HolyName = Holiday(cal.LocalDate, loc)

	// Calculate work week
	cal.WorkWeek = int16(WorkWeek(cal.LocalDate))

	return cal
}

// DateAdd adds a time duration to a date.
func DateAdd(da time.Time, t time.Duration) time.Time {
	return da.Add(t)
}

// DayOfDate returns the day of a date (1-31).
// Note: The original ST code `DATE_TO_DINT(DI) MOD 31 + 1` is incorrect for getting the day of the month.
// A correct implementation is used here.
func DayOfDate(di time.Time) int {
	return di.Day()
}

// DayOfMonth returns the day of a month (1-31).
// This is an alias for DayOfDate, as the ST original is identical and incorrect.
func DayOfMonth(di time.Time) int {
	return di.Day()
}

// DayOfWeek returns the day of the week (Monday=1, ..., Sunday=7).
// Note: The original ST code `(DATE_TO_DINT(DI) / 31) MOD 7 + 1` is incorrect.
// A correct implementation is used here.
func DayOfWeek(di time.Time) int {
	wd := di.Weekday()
	if wd == time.Sunday {
		return 7
	}
	return int(wd)
}

// DayOfYear returns the day of the year (1-366).
func DayOfYear(di time.Time) int {
	return di.YearDay()
}

// DayToTime converts a number of days to a time.Duration.
func DayToTime(d int) time.Duration {
	return time.Duration(d) * 24 * time.Hour
}

// DaysDelta calculates the number of days between two dates.
func DaysDelta(d1, d2 time.Time) int {
	// Truncate to the beginning of the day to get whole days
	d1 = d1.Truncate(24 * time.Hour)
	d2 = d2.Truncate(24 * time.Hour)
	return int(d2.Sub(d1).Hours() / 24)
}

// DaysInMonth returns the number of days in a given month of a given year.
func DaysInMonth(m int, y int) int {
	if m == 2 {
		if LeapDay(y) {
			return 29
		}
		return 28
	}
	// The original ST logic `30 + (M + (M/8)) MOD 2` is a clever bit-twiddle for the 30/31 day pattern.
	// A more readable approach is used here.
	if m == 4 || m == 6 || m == 9 || m == 11 {
		return 30
	}
	return 31
}

// DaysInYear returns the number of days in a year.
func DaysInYear(y int) int {
	if LeapDay(y) {
		return 366
	}
	return 365
}

// LeapDay checks if a year is a leap year.
func LeapDay(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// Dst checks if a given date is within the European daylight saving time period.
func Dst(di time.Time) bool {
	y := di.Year()
	m := int(di.Month())
	d := di.Day()
	w := DayOfWeek(di) // Monday = 1

	if m > 3 && m < 10 {
		return true
	}
	if m == 3 && d-w > 24 {
		return true
	}
	if m == 10 && d-w < 25 {
		return true
	}
	return false
}

// Dt2ToSdt converts a time.Time (DT) into a structured data type Sdt.
func Dt2ToSdt(dtIn time.Time) Sdt {
	var temp Sdt
	temp.Year = int16(dtIn.Year())
	temp.Month = int16(dtIn.Month())
	temp.Day = int16(dtIn.Day())
	temp.Weekday = int16(DayOfWeek(dtIn))
	temp.Hour = int16(dtIn.Hour())
	temp.Minute = int16(dtIn.Minute())
	temp.Second = int16(dtIn.Second())
	return temp
}

// DtToSdt is an alias for Dt2ToSdt as their ST implementations are identical.
func DtToSdt(dtIn time.Time) Sdt {
	return Dt2ToSdt(dtIn)
}

// Easter calculates the date of Easter Sunday for a given year using the Meeus/Jones/Butcher algorithm.
func Easter(y int) time.Time {
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

// Events is a 16 channel timer switch.
// In Go, this is better managed with channels and goroutines rather than a polling function block.
// This is a direct translation for completeness.
func Events(e bool, dtIn time.Time, events []TimerEvent) [16]bool {
	var q [16]bool
	if e {
		td := time.Duration(dtIn.Hour())*time.Hour + time.Duration(dtIn.Minute())*time.Minute + time.Duration(dtIn.Second())*time.Second
		di := dtIn.Truncate(24 * time.Hour)
		dow := DayOfWeek(di)

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
				isHoliday, _ := Holiday(di, int(events[i].Land))
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

// Holiday checks if a given date is a holiday for a specific location.
// This is a simplified placeholder. A full implementation would require the holiday data arrays.
func Holiday(di time.Time, l int) (bool, string) {
	// Placeholder for complex holiday logic.
	// A full implementation would need the h_de, h_at, etc. arrays and logic from the ST code.
	wd := DayOfWeek(di)
	if wd == 6 || wd == 7 {
		return true, "Weekend"
	}
	return false, ""
}

// Hour returns the hour part of a time.Duration representing time of day.
func Hour(t time.Duration) int {
	return int(t.Hours())
}

// HourOfDt returns the hour of a date_and_time.
func HourOfDt(di time.Time) int {
	return di.Hour()
}

// HourToTime converts an integer number of hours to a time.Duration.
func HourToTime(h int) time.Duration {
	return time.Duration(h) * time.Hour
}

// HourToTod is an alias for HourToTime.
func HourToTod(h int) time.Duration {
	return HourToTime(h)
}
