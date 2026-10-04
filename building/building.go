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

// Package building is the port of the OSCAT BUILDING 1.00 library, built on
// the port of OSCAT BASIC, from doc/beedance_building.st, the library
// cleaned for IEC 61131-3 by tools/stclean. Its POUs are in a package for
// each subject: actuators, electrical, hlk (heating, ventilation and air
// conditioning) and jalousie (blinds); this package has BUILDING_VERSION.
//
// The ports follow beebread's conventions: royaljelly's iec types, OSCAT's
// names in upper case, a function block as a struct with INIT and
// Execute(now), and a VAR_IN_OUT as a pointer field.
package building

import (
	. "github.com/apiarytech/beebread/basic"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

// BUILDING_VERSION returns the version of the library, 100 for 1.00, or its
// release date as DATE_TO_DWORD gives it if IN is true.
func BUILDING_VERSION(IN iec.BOOL) iec.DWORD {
	if IN {
		return DATE_TO_DWORD(td.SET_DATE(2011, 2, 3))
	}
	return 100
}
