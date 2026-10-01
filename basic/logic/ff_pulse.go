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

package logic

import (
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// LTCH is a transparent latch with asynchronous reset: while L is true, Q
// follows D.
type LTCH struct {
	D, L, RST iec.BOOL
	Q         iec.BOOL
}

// INIT resets the block.
func (l *LTCH) INIT() { *l = LTCH{} }

// Execute runs the block once.
func (l *LTCH) Execute(now time.Time) {
	if l.RST {
		l.Q = false
	} else if l.L {
		l.Q = l.D
	}
}

// LTCH_4 is a quad transparent latch with asynchronous reset: while L is
// true, Q0..Q3 follow D0..D3.
type LTCH_4 struct {
	D0, D1, D2, D3, L, RST iec.BOOL
	Q0, Q1, Q2, Q3         iec.BOOL
}

// INIT resets the block.
func (l *LTCH_4) INIT() { *l = LTCH_4{} }

// Execute runs the block once.
func (l *LTCH_4) Execute(now time.Time) {
	if l.RST {
		l.Q0, l.Q1, l.Q2, l.Q3 = false, false, false, false
	} else if l.L {
		l.Q0, l.Q1, l.Q2, l.Q3 = l.D0, l.D1, l.D2, l.D3
	}
}

// STORE_8 stores 8 inputs: an output is set when its input is true and
// stays set until RST. SET sets all outputs, and a rising edge of CLR clears
// the lowest output that is set.
type STORE_8 struct {
	SET                            iec.BOOL
	D0, D1, D2, D3, D4, D5, D6, D7 iec.BOOL
	CLR, RST                       iec.BOOL
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7 iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (s *STORE_8) INIT() { *s = STORE_8{} }

// Execute runs the block once.
func (s *STORE_8) Execute(now time.Time) {
	q := [8]*iec.BOOL{&s.Q0, &s.Q1, &s.Q2, &s.Q3, &s.Q4, &s.Q5, &s.Q6, &s.Q7}
	if s.RST || s.SET {
		for _, p := range q {
			*p = !s.RST
		}
		return
	}
	d := [8]iec.BOOL{s.D0, s.D1, s.D2, s.D3, s.D4, s.D5, s.D6, s.D7}
	for i, p := range q {
		if d[i] {
			*p = true
		}
	}
	if s.CLR && !s.edge {
		i := 0
		for i < 7 && !*q[i] {
			i++
		}
		*q[i] = false
	}
	s.edge = s.CLR
}
