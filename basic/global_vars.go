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

const (
	// StringLength defines a default maximum string length, equivalent to STRING_LENGTH in OSCAT.
	StringLength = 250
	// ListLength defines a default maximum list (as a string) length, equivalent to LIST_LENGTH in OSCAT.
	ListLength = 250
)

// Global variables holding constant values from the OSCAT library.
var (
	// Initialize Math constants
	Math = CONSTANTS_MATH{
		Pi:     3.141592653589793,
		Pi2:    6.283185307179586,
		Pi4:    12.566370614359172,
		Pi05:   1.5707963267949,
		Pi025:  0.785398163397448,
		Pi_inv: 0.318309886183791,
		E:      2.718281828459045,
		E_inv:  0.367879441171442,
		Sq2:    1.4142135623731,
		Facts: [13]int32{
			1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800, 39916800, 479001600,
		},
	}

	// Initialize Physics constants
	Phys = CONSTANTS_PHYS{
		C:  299792458.0,
		E:  1.60217653e-19,
		G:  9.80665,
		T0: -273.15,
		Ru: 8.314472,
		Pn: 101325.0,
	}
	Language CONSTANTS_LANGUAGE
	Setup    CONSTANTS_SETUP
	Location CONSTANTS_LOCATION
)

// init initializes the global constant structures with their default values from the OSCAT library.
func init() {

	// Initialize Language constants
	Language = CONSTANTS_LANGUAGE{
		Default: 1,
		Lmax:    3,
		Weekdays: [3][7]string{
			{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
			{"Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag", "Sonntag"},
			{"Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi", "Dimanche"},
		},
		Weekdays2: [3][7]string{
			{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"},
			{"Mo", "Di", "Mi", "Do", "Fr", "Sa", "So"},
			{"Lu", "Ma", "Me", "Je", "Ve", "Sa", "Di"},
		},
		Months: [3][12]string{
			{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
			{"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"},
			{"Janvier", "Février", "mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Decembre"},
		},
		Months3: [3][12]string{
			{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
			{"Jan", "Feb", "Mrz", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"},
			{"Jan", "Fev", "Mar", "Avr", "Mai", "Jun", "Jul", "Aou", "Sep", "Oct", "Nov", "Dec"},
		},
		Dirs: [3][16]string{
			{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"},
			{"N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"},
			{"N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"},
		},
	}

	// Initialize Location constants
	Location = CONSTANTS_LOCATION{
		Default:  1,
		Lmax:     5,
		Language: [5]int16{2, 2, 3, 2, 2},
	}

	// Initialize Setup constants
	Setup = CONSTANTS_SETUP{
		ExtendedASCII: true,
		Charnames: [4]string{
			`"&quot;&&amp;<&lt;>&gt;&euro;&nbsp;&iexcl;&cent;&pound;&curren;&yen;&brvbar;&sect;&uml;&copy;&ordf;&laquo;&not;&shy;&reg;&macr;&deg;&plusmn;&sup2;&sup3;&acute;&micro;&para;&middot;&cedil;&sup1;&ordm;&raquo;&frac14;&Ucirc;`,
			`&frac34;&iquest;&Agrave;&Aacute;&Acirc;&Atilde;&Auml;&Aring;&AElig;&Ccedil;&Egrave;&Eacute;&Ecirc;&Euml;&Igrave;&Iacute;&Icirc;&Iuml;&ETH;&Ntilde;&Ograve;&Oacute;&Ocirc;&Otilde;&Ouml;&times;&Oslash;&Ugrave;&Uacute;&frac12;`,
			`&Uuml;&Yacute;&THORN;&szlig;&agrave;&aacute;&acirc;&atilde;&auml;&aring;&aelig;&ccedil;&egrave;&eacute;&ecirc;&euml;&igrave;&iacute;&icirc;&iuml;&eth;&ntilde;&ograve;&oacute;&ocirc;&otilde;&ouml;&divide;&oslash;&ugrave;`,
			`&uacute;&ucirc;&uuml;&yacute;&thorn;&yuml;`,
		},
		MthOfs: [12]int16{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334},
		Decades: [9]float32{
			1.0, 10.0, 100.0, 1000.0, 10000.0, 100000.0, 1000000.0, 10000000.0, 100000000.0,
		},
	}
}

// OSCAT_VERSION returns the library version number or release date.
func OSCAT_VERSION(in bool) uint32 {
	if in {
		// Corresponds to DATE_TO_DWORD(D#2024-07-16)
		// This is a placeholder. A real implementation would calculate this.
		return 19736 // Days since 1970-01-01 for 2024-07-16
	}
	return 335
}
