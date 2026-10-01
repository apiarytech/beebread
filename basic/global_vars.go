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
	"strings"

	"github.com/apiarytech/royaljelly/iec"
)

const (
	// STRING_LENGTH is the length of the library's strings.
	STRING_LENGTH iec.INT = 250
	// LIST_LENGTH is the length of the library's lists.
	LIST_LENGTH iec.INT = 250
)

// The library's global constants. They are variables, as in OSCAT, so an
// application can change the setup.
var (
	MATH     = newConstantsMath()
	PHYS     = newConstantsPhys()
	LANGUAGE = newConstantsLanguage()
	SETUP    = newConstantsSetup()
	LOCATION = newConstantsLocation()
)

func newConstantsMath() CONSTANTS_MATH {
	return CONSTANTS_MATH{
		PI:     3.14159265358979323846264338327950288,
		PI2:    6.28318530717958647692528676655900576,
		PI4:    12.56637061435917295385057353311801152,
		PI05:   1.5707963267949,
		PI025:  0.785398163397448,
		PI_INV: 0.318309886183791,
		E:      2.71828182845904523536028747135266249,
		E_INV:  0.367879441171442,
		SQ2:    1.4142135623731,
		FACTS:  [13]iec.DINT{1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800, 39916800, 479001600},
	}
}

func newConstantsPhys() CONSTANTS_PHYS {
	return CONSTANTS_PHYS{
		C:  299792458.0,
		E:  1.60217653e-19,
		G:  9.80665,
		T0: -273.15,
		RU: 8.314472,
		PN: 101325.0,
	}
}

func newConstantsLanguage() CONSTANTS_LANGUAGE {
	return CONSTANTS_LANGUAGE{
		DEFAULT: 1,
		LMAX:    3,
		WEEKDAYS: [3][7]iec.STRING{
			{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
			{"Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag", "Sonntag"},
			{"Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi", "Dimanche"},
		},
		WEEKDAYS2: [3][7]iec.STRING{
			{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"},
			{"Mo", "Di", "Mi", "Do", "Fr", "Sa", "So"},
			{"Lu", "Ma", "Me", "Je", "Ve", "Sa", "Di"},
		},
		MONTHS: [3][12]iec.STRING{
			{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
			{"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"},
			{"Janvier", "Février", "mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Decembre"},
		},
		MONTHS3: [3][12]iec.STRING{
			{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
			{"Jan", "Feb", "Mrz", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"},
			{"Jan", "Fev", "Mar", "Avr", "Mai", "Jun", "Jul", "Aou", "Sep", "Oct", "Nov", "Dec"},
		},
		DIRS: [3][16]iec.STRING{
			{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"},
			{"N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"},
			{"N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"},
		},
	}
}

func newConstantsLocation() CONSTANTS_LOCATION {
	return CONSTANTS_LOCATION{
		DEFAULT:  1,
		LMAX:     5,
		LANGUAGE: [5]iec.INT{2, 2, 3, 2, 2},
	}
}

// charnames lists the characters of each SETUP.CHARNAMES string with their
// HTML entity names. Each character of a list is followed by "&name;". The
// euro sign has the code it has in Windows-1252, 128.
var charnames = [4][]string{
	{
		"\"", "quot", "&", "amp", "<", "lt", ">", "gt", "\u0080", "euro", " ", "nbsp", "¡", "iexcl",
		"¢", "cent", "£", "pound", "¤", "curren", "¥", "yen", "¦", "brvbar", "§", "sect", "¨", "uml",
		"©", "copy", "ª", "ordf", "«", "laquo", "¬", "not", "­", "shy", "®", "reg", "¯", "macr",
		"°", "deg", "±", "plusmn", "²", "sup2", "³", "sup3", "´", "acute", "µ", "micro", "¶", "para",
		"·", "middot", "¸", "cedil", "¹", "sup1", "º", "ordm", "»", "raquo", "¼", "frac14", "Û", "Ucirc",
	},
	{
		"¾", "frac34", "¿", "iquest", "À", "Agrave", "Á", "Aacute", "Â", "Acirc", "Ã", "Atilde",
		"Ä", "Auml", "Å", "Aring", "Æ", "AElig", "Ç", "Ccedil", "È", "Egrave", "É", "Eacute",
		"Ê", "Ecirc", "Ë", "Euml", "Ì", "Igrave", "Í", "Iacute", "Î", "Icirc", "Ï", "Iuml",
		"Ð", "ETH", "Ñ", "Ntilde", "Ò", "Ograve", "Ó", "Oacute", "Ô", "Ocirc", "Õ", "Otilde",
		"Ö", "Ouml", "×", "times", "Ø", "Oslash", "Ù", "Ugrave", "Ú", "Uacute", "½", "frac12",
	},
	{
		"Ü", "Uuml", "Ý", "Yacute", "Þ", "THORN", "ß", "szlig", "à", "agrave", "á", "aacute",
		"â", "acirc", "ã", "atilde", "ä", "auml", "å", "aring", "æ", "aelig", "ç", "ccedil",
		"è", "egrave", "é", "eacute", "ê", "ecirc", "ë", "euml", "ì", "igrave", "í", "iacute",
		"î", "icirc", "ï", "iuml", "ð", "eth", "ñ", "ntilde", "ò", "ograve", "ó", "oacute",
		"ô", "ocirc", "õ", "otilde", "ö", "ouml", "÷", "divide", "ø", "oslash", "ù", "ugrave",
	},
	{
		"ú", "uacute", "û", "ucirc", "ü", "uuml", "ý", "yacute", "þ", "thorn", "ÿ", "yuml",
	},
}

func newConstantsSetup() CONSTANTS_SETUP {
	s := CONSTANTS_SETUP{
		EXTENDED_ASCII: true,
		MTH_OFS:        [12]iec.INT{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334},
		DECADES:        [9]iec.REAL{1.0, 10.0, 100.0, 1000.0, 10000.0, 100000.0, 1000000.0, 10000000.0, 100000000.0},
	}
	// Each list is ";", then character "&" name ";" for each entry, as in
	// OSCAT.
	for i, list := range charnames {
		var b strings.Builder
		b.WriteString(";")
		for j := 0; j < len(list); j += 2 {
			b.WriteString(list[j] + "&" + list[j+1] + ";")
		}
		s.CHARNAMES[i] = iec.STRING(b.String())
	}
	return s
}
