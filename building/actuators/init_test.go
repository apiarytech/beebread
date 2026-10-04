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

package actuators

import (
	"testing"

	"github.com/apiarytech/beebread/building/internal/inittest"
)

// initBlocks are the function blocks of the package.
var initBlocks = map[string]func() inittest.Initer{
	"ACTUATOR_2P":   func() inittest.Initer { return new(ACTUATOR_2P) },
	"ACTUATOR_3P":   func() inittest.Initer { return new(ACTUATOR_3P) },
	"ACTUATOR_A":    func() inittest.Initer { return new(ACTUATOR_A) },
	"ACTUATOR_COIL": func() inittest.Initer { return new(ACTUATOR_COIL) },
	"ACTUATOR_PUMP": func() inittest.Initer { return new(ACTUATOR_PUMP) },
	"ACTUATOR_UD":   func() inittest.Initer { return new(ACTUATOR_UD) },
	"AUTORUN":       func() inittest.Initer { return new(AUTORUN) },
}

// TestINIT checks that INIT restores every block.
func TestINIT(t *testing.T) { inittest.CheckINIT(t, initBlocks) }

// TestINITDefaults checks INIT against the OSCAT initial values.
func TestINITDefaults(t *testing.T) { inittest.CheckDefaults(t, initBlocks) }
