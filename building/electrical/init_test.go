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

package electrical

import (
	"testing"

	"github.com/apiarytech/beebread/building/internal/inittest"
)

// initBlocks are the function blocks of the package.
var initBlocks = map[string]func() inittest.Initer{
	"CLICK":        func() inittest.Initer { return new(CLICK) },
	"CLICK_MODE":   func() inittest.Initer { return new(CLICK_MODE) },
	"DEBOUNCE":     func() inittest.Initer { return new(DEBOUNCE) },
	"DIMM_2":       func() inittest.Initer { return new(DIMM_2) },
	"DIMM_I":       func() inittest.Initer { return new(DIMM_I) },
	"F_LAMP":       func() inittest.Initer { return new(F_LAMP) },
	"PULSE_LENGTH": func() inittest.Initer { return new(PULSE_LENGTH) },
	"PULSE_T":      func() inittest.Initer { return new(PULSE_T) },
	"SWITCH_I":     func() inittest.Initer { return new(SWITCH_I) },
	"SWITCH_X":     func() inittest.Initer { return new(SWITCH_X) },
	"SW_RECONFIG":  func() inittest.Initer { return new(SW_RECONFIG) },
	"TIMER_1":      func() inittest.Initer { return new(TIMER_1) },
	"TIMER_2":      func() inittest.Initer { return new(TIMER_2) },
	"TIMER_EXT":    func() inittest.Initer { return new(TIMER_EXT) },
	"TIMER_P4":     func() inittest.Initer { return new(TIMER_P4) },
}

// TestINIT checks that INIT restores every block.
func TestINIT(t *testing.T) { inittest.CheckINIT(t, initBlocks) }

// TestINITDefaults checks INIT against the OSCAT initial values.
func TestINITDefaults(t *testing.T) { inittest.CheckDefaults(t, initBlocks) }
