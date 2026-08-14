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

// TPlcUs returns the system time as a high-resolution duration from a fixed point, in microseconds.
// This is a Go-native way to get a monotonic clock reading.
func TPlcUs() int64 {
	// time.Now().UnixNano() provides a monotonic clock reading suitable for measuring intervals.
	return time.Now().UnixNano() / 1000
}

// TCS delivers the time since it was last called on the output TC in seconds.
type TCS struct {
	TC float64

	// internal state
	init bool
	last int64
}

// Update executes the cycle time measurement logic.
func (t *TCS) Update() {
	tx := TPlcUs()

	if !t.init {
		t.init = true
		t.TC = 0.0
	} else {
		t.TC = float64(tx-t.last) * 1.0e-6
	}
	t.last = tx
}
