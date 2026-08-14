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

// Ltch is a transparent latch with an asynchronous reset.
// As long as L (Latch) is true, the output Q follows the input D.
// When L goes false, Q holds its last value.
// RST asynchronously forces Q to false.
type Ltch struct {
	Q bool
}

// Update executes the latch logic for one cycle.
func (l *Ltch) Update(d, latch, rst bool) {
	if rst {
		l.Q = false
	} else if latch {
		l.Q = d
	}
	// If neither rst nor latch is true, Q retains its state.
}

// Ltch4 is a quad transparent latch with a common asynchronous reset and latch input.
type Ltch4 struct {
	Q0 bool
	Q1 bool
	Q2 bool
	Q3 bool
}

// Update executes the latch logic for one cycle.
func (l *Ltch4) Update(d0, d1, d2, d3, latch, rst bool) {
	if rst {
		l.Q0 = false
		l.Q1 = false
		l.Q2 = false
		l.Q3 = false
	} else if latch {
		l.Q0 = d0
		l.Q1 = d1
		l.Q2 = d2
		l.Q3 = d3
	}
	// If neither rst nor latch is true, outputs retain their state.
}

// Store8 stores up to 8 boolean inputs until a reset clears the outputs.
// The respective output is set with a true value at the respective input and stays true until a reset.
// A Set input sets all outputs true simultaneously.
// A rising edge on Clr resets the lowest priority output (Q0 first, then Q1, etc.).
type Store8 struct {
	Q [8]bool

	// internal state
	clrEdge bool
}

// Update executes the storage logic for one cycle.
func (s *Store8) Update(set bool, d [8]bool, clr, rst bool) {
	if rst {
		s.Q = [8]bool{} // Zeros all elements
	} else if set {
		for i := range s.Q {
			s.Q[i] = true
		}
	} else {
		// Set individual bits based on data inputs
		for i, val := range d {
			if val {
				s.Q[i] = true
			}
		}

		// Rising edge on Clr
		if clr && !s.clrEdge {
			// Find the first set bit (from Q0 to Q7) and clear it.
			for i := range s.Q {
				if s.Q[i] {
					s.Q[i] = false
					break // Exit after clearing the first one found
				}
			}
		}
	}
	s.clrEdge = clr
}
