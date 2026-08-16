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

import "time"

// CALENDAR corresponds to the CALENDAR struct in OSCAT, holding date, time,
// and astronomical data.
type CALENDAR struct {
	UTC       time.Time     // world time UTC
	LocalDT   time.Time     // local time
	LocalDate time.Time     // local date
	LocalTOD  time.Duration // local time of day
	Year      int16         // year of LocalDate
	Month     int16         // month of LocalDate
	Day       int16         // day of LocalDate
	Weekday   int16         // weekday of LocalDate
	Offset    int16         // Time Zone Offset for Local time in minutes
	DstEn     bool          // daylight savings time enable
	DstOn     bool          // true when daylight savings time is on
	Name      string        // name of time zone (originally STRING(5))
	Language  int16         // location number pls see location setup
	Longitude float32       // longitude of current location
	Latitude  float32       // latitude of current location
	SunRise   time.Duration // sun_rise for current location
	SunSet    time.Duration // sun_set for current location
	SunMidday time.Duration // worldtime when sun stands at south position
	SunHeigth float32       // suns heigth at midday, south position
	SunHor    float32       // sun angle horizontal 0 = north in degrees
	SunVer    float32       // sun angle vertical above horizon in degrees
	Night     bool          // true between sun_set and sun_rise
	Holiday   bool          // true when holiday
	HolyName  string        // name of holiday (originally STRING(30))
	WorkWeek  int16         // current work week
}

// COMPLEX corresponds to the COMPLEX struct in OSCAT for representing complex numbers.
type COMPLEX struct {
	Re float32
	Im float32
}

// CONSTANTS_LANGUAGE corresponds to the CONSTANTS_LANGUAGE struct in OSCAT.
type CONSTANTS_LANGUAGE struct {
	// Language Setup
	Default   int16 // 1=english, 2=german 3=french
	Lmax      int16
	Weekdays  [3][7]string  // Corresponds to ARRAY[1..3, 1..7] OF STRING(10)
	Weekdays2 [3][7]string  // Corresponds to ARRAY[1..3, 1..7] OF STRING(2)
	Months    [3][12]string // Corresponds to ARRAY[1..3, 1..12] OF STRING(10)
	Months3   [3][12]string // Corresponds to ARRAY[1..3, 1..12] OF STRING(3)
	Dirs      [3][16]string // Corresponds to ARRAY[1..3, 0..15] OF STRING(3)
}

// CONSTANTS_LOCATION corresponds to the CONSTANTS_LOCATION struct in OSCAT.
type CONSTANTS_LOCATION struct {
	// location setup
	Default int16 // 1=germany, 2=austria 3=france 4=belgium-german 5= italien-Sdtirol
	Lmax    int16

	// language spoken in the location
	Language [5]int16 // Corresponds to ARRAY[1..5] OF INT
}

// CONSTANTS_MATH corresponds to the CONSTANTS_MATH struct in OSCAT.
type CONSTANTS_MATH struct {
	Pi     float64   // Kreiszahl PI
	Pi2    float64   // PI * 2
	Pi4    float64   // PI * 4
	Pi05   float64   // PI / 2
	Pi025  float64   // PI / 4
	Pi_inv float64   // 1 / PI
	E      float64   // Euler constant e
	E_inv  float64   // 1 / e
	Sq2    float64   // Wurzel von 2
	Facts  [13]int32 // Corresponds to ARRAY[0..12] OF DINT
}

// CONSTANTS_PHYS corresponds to the CONSTANTS_PHYS struct in OSCAT.
type CONSTANTS_PHYS struct {
	C  float64 // Lichtgeschwindigkeit in m/s
	E  float64 // elementarladung in Colomb = A * s
	G  float64 // Erdbeschleunigung in m / s
	T0 float32 // absoluter Nullpunkt in C
	Ru float32 // Universelle Gaskonstante in J / (mol  K)
	Pn float32 // NormalDruck in Pa
}

// CONSTANTS_SETUP corresponds to the CONSTANTS_SETUP struct in OSCAT.
type CONSTANTS_SETUP struct {
	// setup Parameters
	ExtendedASCII bool
	Charnames     [4]string  // Corresponds to ARRAY[1..4] OF STRING(253)
	MthOfs        [12]int16  // Corresponds to ARRAY[1..12] OF INT
	Decades       [9]float32 // Corresponds to ARRAY[0..8] OF REAL
}

// ESR_DATA corresponds to the ESR_DATA struct in OSCAT.
type ESR_DATA struct {
	Typ    byte
	Adress string // Originally STRING(10)
	Ds     time.Time
	Ts     time.Duration
	Data   [8]byte // Corresponds to ARRAY[0..7] OF BYTE
}

// FRACTION corresponds to the FRACTION struct in OSCAT.
type FRACTION struct {
	Numerator   int16
	Denominator int16
}

// HOLIDAY_DATA corresponds to the HOLIDAY_DATA struct in OSCAT.
type HOLIDAY_DATA struct {
	Name  string // Originally STRING(30)
	Day   int8
	Month int8
	Use   int8
}

// REAL2 corresponds to the REAL2 struct in OSCAT for double-precision emulation.
type REAL2 struct {
	R1 float32 // small value
	Rx float32 // big value
}

// Sdt corresponds to the SDT (Structured Date Time) struct in OSCAT.
type SDT struct {
	Year    int16
	Month   int16
	Day     int16
	Weekday int16
	Hour    int16
	Minute  int16
	Second  int16
	Ms      int16
}

// TIMER_EVENT corresponds to the TIMER_EVENT struct in OSCAT.
type TIMER_EVENT struct {
	Typ      byte
	Channel  byte
	Day      byte
	Start    time.Duration
	Duration time.Duration
	Land     byte
	Lor      byte
	Last     time.Time
}

// VECTOR_3 corresponds to the VECTOR_3 struct in OSCAT.
type VECTOR_3 struct {
	X float32
	Y float32
	Z float32
}
