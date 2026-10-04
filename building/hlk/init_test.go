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

package hlk

import (
	"testing"

	"github.com/apiarytech/beebread/building/internal/inittest"
)

// initBlocks are the function blocks of the package.
var initBlocks = map[string]func() inittest.Initer{
	"BOILER":     func() inittest.Initer { return new(BOILER) },
	"BURNER":     func() inittest.Initer { return new(BURNER) },
	"HEAT_METER": func() inittest.Initer { return new(HEAT_METER) },
	"HEAT_TEMP":  func() inittest.Initer { return new(HEAT_TEMP) },
	"LEGIONELLA": func() inittest.Initer { return new(LEGIONELLA) },
	"TANK_LEVEL": func() inittest.Initer { return new(TANK_LEVEL) },
	"TEMP_EXT":   func() inittest.Initer { return new(TEMP_EXT) },
	"T_AVG24":    func() inittest.Initer { return new(T_AVG24) },
}

// TestINIT checks that INIT restores every block.
func TestINIT(t *testing.T) { inittest.CheckINIT(t, initBlocks) }

// TestINITDefaults checks INIT against the OSCAT initial values.
func TestINITDefaults(t *testing.T) { inittest.CheckDefaults(t, initBlocks) }
