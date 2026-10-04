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

package jalousie

import (
	"testing"

	"github.com/apiarytech/beebread/building/internal/inittest"
)

// initBlocks are the function blocks of the package.
var initBlocks = map[string]func() inittest.Initer{
	"BLIND_ACTUATOR":  func() inittest.Initer { return new(BLIND_ACTUATOR) },
	"BLIND_CONTROL":   func() inittest.Initer { return new(BLIND_CONTROL) },
	"BLIND_CONTROL_S": func() inittest.Initer { return new(BLIND_CONTROL_S) },
	"BLIND_INPUT":     func() inittest.Initer { return new(BLIND_INPUT) },
	"BLIND_NIGHT":     func() inittest.Initer { return new(BLIND_NIGHT) },
	"BLIND_SCENE":     func() inittest.Initer { return new(BLIND_SCENE) },
	"BLIND_SECURITY":  func() inittest.Initer { return new(BLIND_SECURITY) },
	"BLIND_SET":       func() inittest.Initer { return new(BLIND_SET) },
	"BLIND_SHADE":     func() inittest.Initer { return new(BLIND_SHADE) },
	"BLIND_SHADE_S":   func() inittest.Initer { return new(BLIND_SHADE_S) },
}

// TestINIT checks that INIT restores every block.
func TestINIT(t *testing.T) { inittest.CheckINIT(t, initBlocks) }

// TestINITDefaults checks INIT against the OSCAT initial values.
func TestINITDefaults(t *testing.T) { inittest.CheckDefaults(t, initBlocks) }
