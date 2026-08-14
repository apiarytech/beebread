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

package logic

import "time"

// TON is a Timer On-Delay function block.
// It delays a rising edge of IN by the duration PT.
type TON struct {
	Q  bool
	ET time.Duration

	// internal state
	startTime time.Time
	IN        bool
}

// Update executes the TON logic for one cycle.
func (t *TON) Update(in bool, pt time.Duration) {
	if !t.IN && in { // Rising edge
		t.startTime = time.Now()
	}
	t.IN = in

	if t.IN {
		t.ET = time.Since(t.startTime)
		if t.ET >= pt {
			t.Q = true
			t.ET = pt
		}
	} else {
		t.Q = false
		t.ET = 0
	}
}

// TOF is a Timer Off-Delay function block.
// It delays a falling edge of IN by the duration PT.
type TOF struct {
	Q  bool
	ET time.Duration

	// internal state
	startTime time.Time
	IN        bool
}

// Update executes the TOF logic for one cycle.
func (t *TOF) Update(in bool, pt time.Duration) bool {
	if t.IN && !in { // Falling edge
		t.startTime = time.Now()
	}
	t.IN = in

	if t.IN {
		t.Q = true
		t.ET = 0
	} else {
		t.ET = time.Since(t.startTime)
		if t.ET >= pt {
			t.Q = false
			t.ET = pt
		}
	}
	return t.Q
}

// TP is a Pulse Timer function block.
// It generates a pulse of duration PT on a rising edge of IN.
type TP struct {
	Q  bool
	ET time.Duration

	// internal state
	ton TON
}

// Update executes the TP logic for one cycle.
func (t *TP) Update(in bool, pt time.Duration) {
	t.ton.Update(in, pt)
	t.Q = in && !t.ton.Q
	t.ET = t.ton.ET
}

// TP1D is a delayed pulse timer.
// It generates a pulse of duration PT1 on a rising edge of IN.
// After the pulse, it waits for a delay of PTD before it can be triggered again.
type TP1D struct {
	Q   bool // Pulse output
	W   bool // Wait state (locked out)
	IN  bool // Trigger input (will be reset internally)
	RST bool // Reset input (will be reset internally)
	PT1 time.Duration
	PTD time.Duration

	// internal state
	startTime time.Time
	edge      bool
}

// Update executes the TP1D logic for one cycle.
func (t *TP1D) Update() {
	tx := time.Now()

	if t.RST {
		t.Q = false
		t.W = false
		t.RST = false // Reset is self-clearing
		return
	}

	if t.W {
		// We are in the lockout/delay period
		if tx.Sub(t.startTime) >= t.PTD {
			t.W = false
		}
	} else if t.IN && !t.edge {
		// Rising edge on IN, and not in lockout period
		t.Q = true
		t.startTime = tx
	} else if t.Q {
		// Pulse is currently active, check if it's time to end it
		if tx.Sub(t.startTime) >= t.PT1 {
			t.Q = false
			t.W = true // Start the lockout period
			t.startTime = tx
		}
	}

	t.edge = t.IN
	t.IN = false // The IN signal is a pulse, so it's reset after being processed.
}
