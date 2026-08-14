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

import (
	"github.com/apiarytech/beebread/basic/math"
)

// COUNT_BR is a byte counter with independent up and down inputs.
// The counter counts from 0 to MX and wraps around.
// A step input sets the counter's stepping width.
type COUNT_BR struct {
	Cnt byte
	// internal state
	lastUp bool
	lastDn bool
}

// Update executes the counter logic for one cycle.
func (c *COUNT_BR) Update(set bool, in byte, up, dn bool, step, mx byte, rst bool) {
	if rst {
		c.Cnt = 0
	} else if set {
		c.Cnt = math.LimitB(0, in, mx)
	} else if up && !c.lastUp {
		c.Cnt = byte(math.Inc(int(c.Cnt), int(step), int(mx)))
	} else if dn && !c.lastDn {
		c.Cnt = byte(math.Inc(int(c.Cnt), -int(step), int(mx)))
	}
	c.lastUp = up
	c.lastDn = dn
}

// COUNT_DR is a DWORD counter with independent up and down inputs.
// The counter counts from 0 to MX and wraps around.
// A step input sets the counter's stepping width.
type COUNT_DR struct {
	Cnt uint32
	// internal state
	lastUp bool
	lastDn bool
}

// Update executes the counter logic for one cycle.
func (c *COUNT_DR) Update(set bool, in uint32, up, dn bool, step, mx uint32, rst bool) {
	if rst {
		c.Cnt = 0
	} else if set {
		c.Cnt = math.LimitDW(0, in, mx)
	} else if up && !c.lastUp {
		if step > mx-c.Cnt {
			c.Cnt = c.Cnt - mx + step - 1
		} else {
			c.Cnt += step
		}
	} else if dn && !c.lastDn {
		if step > c.Cnt {
			c.Cnt = c.Cnt - step + mx + 1
		} else {
			c.Cnt -= step
		}
	}
	c.lastUp = up
	c.lastDn = dn
}

// FF_D2E is a dual D-type flip-flop with reset and rising clock trigger.
type FF_D2E struct {
	Q0 bool
	Q1 bool
	// internal state
	edge bool
}

// Update executes the flip-flop logic for one cycle.
func (f *FF_D2E) Update(d0, d1, clk, rst bool) {
	if rst {
		f.Q0 = false
		f.Q1 = false
	} else if clk && !f.edge {
		f.Q0 = d0
		f.Q1 = d1
	}
	f.edge = clk
}

// FF_D4E is a quad D-type flip-flop with reset and rising clock trigger.
type FF_D4E struct {
	Q0 bool
	Q1 bool
	Q2 bool
	Q3 bool
	// internal state
	edge bool
}

// Update executes the flip-flop logic for one cycle.
func (f *FF_D4E) Update(d0, d1, d2, d3, clk, rst bool) {
	if rst {
		f.Q0 = false
		f.Q1 = false
		f.Q2 = false
		f.Q3 = false
	} else if clk && !f.edge {
		f.Q0 = d0
		f.Q1 = d1
		f.Q2 = d2
		f.Q3 = d3
	}
	f.edge = clk
}

// FF_DRE is a D-type flip-flop with set, reset, and rising clock trigger.
type FF_DRE struct {
	Q bool
	// internal state
	edge bool
}

// Update executes the flip-flop logic for one cycle.
func (f *FF_DRE) Update(set, d, clk, rst bool) {
	if rst {
		f.Q = false
	} else if set {
		f.Q = true
	} else if clk && !f.edge {
		f.Q = d
	}
	f.edge = clk
}

// FF_RSE is a rising edge-triggered RS flip-flop.
// A rising edge on CS sets Q to true, and a rising edge on CR sets Q to false.
// CR has priority over CS.
type FF_RSE struct {
	Q bool
	// internal state
	csEdge, crEdge bool
}

// Update executes the flip-flop logic for one cycle.
func (f *FF_RSE) Update(cs, cr, rst bool) {
	if rst {
		f.Q = false
	} else if cr && !f.crEdge {
		f.Q = false
	} else if cs && !f.csEdge {
		f.Q = true
	}
	f.csEdge = cs
	f.crEdge = cr
}

// TOGGLE is a flip-flop where the output changes state with every rising edge of CLK.
type TOGGLE struct {
	Q bool
	// internal state
	edge bool
}

// Update executes the toggle logic for one cycle.
func (t *TOGGLE) Update(clk, rst bool) {
	if rst {
		t.Q = false
	} else if clk && !t.edge {
		t.Q = !t.Q
	}
	t.edge = clk
}
