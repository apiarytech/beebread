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

import "github.com/apiarytech/royaljelly/iec"

// CALENDAR holds date, time and astronomical data for a location.
type CALENDAR struct {
	UTC        iec.DT     // world time UTC
	LOCAL_DT   iec.DT     // local time
	LOCAL_DATE iec.DATE   // local date
	LOCAL_TOD  iec.TOD    // local time of day
	YEAR       iec.INT    // year of LOCAL_DATE
	MONTH      iec.INT    // month of LOCAL_DATE
	DAY        iec.INT    // day of LOCAL_DATE
	WEEKDAY    iec.INT    // weekday of LOCAL_DATE
	OFFSET     iec.INT    // time zone offset for local time in minutes
	DST_EN     iec.BOOL   // daylight savings time enable
	DST_ON     iec.BOOL   // true when daylight savings time is on
	NAME       iec.STRING // name of time zone, STRING(5)
	LANGUAGE   iec.INT    // language number, see the language setup
	LONGITUDE  iec.REAL   // longitude of current location
	LATITUDE   iec.REAL   // latitude of current location
	SUN_RISE   iec.TOD    // sun rise for current location
	SUN_SET    iec.TOD    // sun set for current location
	SUN_MIDDAY iec.TOD    // world time when the sun stands at south position
	SUN_HEIGTH iec.REAL   // sun's height at midday, south position
	SUN_HOR    iec.REAL   // sun angle horizontal, 0 = north, in degrees
	SUN_VER    iec.REAL   // sun angle vertical above horizon in degrees
	NIGHT      iec.BOOL   // true between sun set and sun rise
	HOLIDAY    iec.BOOL   // true when holiday
	HOLY_NAME  iec.STRING // name of holiday, STRING(30)
	WORK_WEEK  iec.INT    // current work week
}

// COMPLEX is a complex number.
type COMPLEX struct {
	RE iec.REAL
	IM iec.REAL
}

// CONSTANTS_LANGUAGE is the language setup. The arrays' first index is the
// language, 1=english, 2=german, 3=french, at Go index language-1.
type CONSTANTS_LANGUAGE struct {
	DEFAULT   iec.INT
	LMAX      iec.INT
	WEEKDAYS  [3][7]iec.STRING  // ARRAY[1..3, 1..7] OF STRING(10)
	WEEKDAYS2 [3][7]iec.STRING  // ARRAY[1..3, 1..7] OF STRING(2)
	MONTHS    [3][12]iec.STRING // ARRAY[1..3, 1..12] OF STRING(10)
	MONTHS3   [3][12]iec.STRING // ARRAY[1..3, 1..12] OF STRING(3)
	DIRS      [3][16]iec.STRING // ARRAY[1..3, 0..15] OF STRING(3)
}

// CONSTANTS_LOCATION is the location setup: 1=germany, 2=austria, 3=france,
// 4=belgium-german, 5=italy-south tyrol.
type CONSTANTS_LOCATION struct {
	DEFAULT  iec.INT
	LMAX     iec.INT
	LANGUAGE [5]iec.INT // ARRAY[1..5] OF INT, the language spoken in the location
}

// CONSTANTS_MATH holds mathematical constants.
type CONSTANTS_MATH struct {
	PI     iec.REAL
	PI2    iec.REAL // PI * 2
	PI4    iec.REAL // PI * 4
	PI05   iec.REAL // PI / 2
	PI025  iec.REAL // PI / 4
	PI_INV iec.REAL // 1 / PI
	E      iec.REAL // Euler constant e
	E_INV  iec.REAL // 1 / e
	SQ2    iec.REAL // square root of 2
	FACTS  [13]iec.DINT
}

// CONSTANTS_PHYS holds physical constants.
type CONSTANTS_PHYS struct {
	C  iec.REAL // speed of light in m/s
	E  iec.REAL // elementary charge in Coulomb = A * s
	G  iec.REAL // acceleration of gravity in m/s²
	T0 iec.REAL // absolute zero in °C
	RU iec.REAL // universal gas constant in J / (mol * K)
	PN iec.REAL // standard pressure in Pa
}

// CONSTANTS_SETUP holds the library's setup parameters.
type CONSTANTS_SETUP struct {
	EXTENDED_ASCII iec.BOOL
	CHARNAMES      [4]iec.STRING // ARRAY[1..4] OF STRING(253)
	MTH_OFS        [12]iec.INT   // ARRAY[1..12] OF INT
	DECADES        [9]iec.REAL   // ARRAY[0..8] OF REAL
}

// ESR_DATA is an event, status or error record.
type ESR_DATA struct {
	TYP    iec.BYTE
	ADRESS iec.STRING // STRING(10)
	DS     iec.DT
	TS     iec.TIME
	DATA   [8]iec.BYTE
}

// FRACTION is a fraction of two integers.
type FRACTION struct {
	NUMERATOR   iec.INT
	DENOMINATOR iec.INT
}

// HOLIDAY_DATA describes a holiday. If MONTH is 0, DAY is the offset in days
// from easter. USE is 0 for not used, 1 for used, and -1..-7 for the weekday
// before the date, for example -3 is the wednesday before DAY.MONTH.
type HOLIDAY_DATA struct {
	NAME  iec.STRING // STRING(30)
	DAY   iec.SINT
	MONTH iec.SINT
	USE   iec.SINT
}

// REAL2 emulates a double precision value with two REALs.
type REAL2 struct {
	R1 iec.REAL // small value
	RX iec.REAL // big value
}

// SDT is a structured date and time.
type SDT struct {
	YEAR    iec.INT
	MONTH   iec.INT
	DAY     iec.INT
	WEEKDAY iec.INT
	HOUR    iec.INT
	MINUTE  iec.INT
	SECOND  iec.INT
	MS      iec.INT
}

// TIMER_EVENT describes one event of a timer switch.
type TIMER_EVENT struct {
	TYP      iec.BYTE
	CHANNEL  iec.BYTE
	DAY      iec.BYTE
	START    iec.TOD
	DURATION iec.TIME
	LAND     iec.BYTE
	LOR      iec.BYTE
	LAST     iec.DT
}

// VECTOR_3 is a vector in three dimensional space.
type VECTOR_3 struct {
	X iec.REAL
	Y iec.REAL
	Z iec.REAL
}
