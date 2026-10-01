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
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// OSCAT measures time with the PLC timer, a DWORD of milliseconds or
// microseconds since the PLC started that wraps around. The port derives
// that timer from a point in time: the scan time a function block is given,
// or the clock for the functions T_PLC_MS and T_PLC_US.

var (
	// epoch is the time the PLC timer counts from.
	epoch = time.Now()
	// clock gives the current time to the functions that read the timer.
	clock = time.Now
)

// SetEpoch sets the time the PLC timer counts from. It is the time the
// program started unless set.
func SetEpoch(t time.Time) { epoch = t }

// SetClock sets the clock T_PLC_MS and T_PLC_US read, which is time.Now
// unless set. A program that runs on simulated time sets it to that time.
func SetClock(now func() time.Time) { clock = now }

// Now returns the current time of the clock T_PLC_MS and T_PLC_US read.
func Now() time.Time { return clock() }

// PLC_MS returns the PLC timer in milliseconds at the time now.
func PLC_MS(now time.Time) iec.DWORD {
	return iec.DWORD(now.Sub(epoch).Milliseconds())
}

// PLC_US returns the PLC timer in microseconds at the time now.
func PLC_US(now time.Time) iec.DWORD {
	return iec.DWORD(now.Sub(epoch).Microseconds())
}

// T_PLC_MS reads the PLC timer in milliseconds.
func T_PLC_MS() iec.DWORD { return PLC_MS(clock()) }

// T_PLC_US reads the PLC timer in microseconds.
func T_PLC_US() iec.DWORD { return PLC_US(clock()) }
