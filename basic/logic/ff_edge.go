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

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/royaljelly/fb/triggers"
	"github.com/apiarytech/royaljelly/iec"
)

// COUNT_BR is a byte counter with separate up and down inputs. It counts by
// STEP from 0 to MX and continues at 0 after MX.
type COUNT_BR struct {
	SET  iec.BOOL
	IN   iec.BYTE
	UP   iec.BOOL
	DN   iec.BOOL
	STEP iec.BYTE // default 1
	MX   iec.BYTE // default 255
	RST  iec.BOOL
	CNT  iec.BYTE

	lastUp, lastDn iec.BOOL
}

// INIT resets the block and sets STEP and MX to their initial values.
func (c *COUNT_BR) INIT() { *c = COUNT_BR{STEP: 1, MX: 255} }

// Execute runs the block once.
func (c *COUNT_BR) Execute(now time.Time) {
	switch {
	case bool(c.RST):
		c.CNT = 0
	case bool(c.SET):
		c.CNT = LIMIT(0, c.IN, c.MX)
	case bool(c.UP && !c.lastUp):
		c.CNT = iec.BYTE(math.INC(iec.INT(c.CNT), iec.INT(c.STEP), iec.INT(c.MX)))
	case bool(c.DN && !c.lastDn):
		c.CNT = iec.BYTE(math.INC(iec.INT(c.CNT), -iec.INT(c.STEP), iec.INT(c.MX)))
	}
	c.lastUp = c.UP
	c.lastDn = c.DN
}

// COUNT_DR is a DWORD counter with separate up and down inputs. It counts
// by STEP from 0 to MX and continues at 0 after MX.
type COUNT_DR struct {
	SET  iec.BOOL
	IN   iec.DWORD
	UP   iec.BOOL
	DN   iec.BOOL
	STEP iec.DWORD // default 1
	MX   iec.DWORD // default 16#FFFFFFFF
	RST  iec.BOOL
	CNT  iec.DWORD

	lastUp, lastDn iec.BOOL
}

// INIT resets the block and sets STEP and MX to their initial values.
func (c *COUNT_DR) INIT() { *c = COUNT_DR{STEP: 1, MX: 0xFFFFFFFF} }

// Execute runs the block once.
func (c *COUNT_DR) Execute(now time.Time) {
	switch {
	case bool(c.RST):
		c.CNT = 0
	case bool(c.SET):
		c.CNT = LIMIT(0, c.IN, c.MX)
	case bool(c.UP && !c.lastUp):
		if c.STEP > c.MX-c.CNT {
			c.CNT = c.CNT - c.MX + c.STEP - 1
		} else {
			c.CNT = c.CNT + c.STEP
		}
	case bool(c.DN && !c.lastDn):
		if c.STEP > c.CNT {
			c.CNT = c.CNT - c.STEP + c.MX + 1
		} else {
			c.CNT = c.CNT - c.STEP
		}
	}
	c.lastUp = c.UP
	c.lastDn = c.DN
}

// FF_D2E is a dual D flip-flop with reset, triggered by a rising edge of CLK.
type FF_D2E struct {
	D0, D1, CLK, RST iec.BOOL
	Q0, Q1           iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (f *FF_D2E) INIT() { *f = FF_D2E{} }

// Execute runs the block once.
func (f *FF_D2E) Execute(now time.Time) {
	if f.RST {
		f.Q0, f.Q1 = false, false
	} else if f.CLK && !f.edge {
		f.Q0, f.Q1 = f.D0, f.D1
	}
	f.edge = f.CLK
}

// FF_D4E is a quad D flip-flop with reset, triggered by a rising edge of
// CLK.
type FF_D4E struct {
	D0, D1, D2, D3, CLK, RST iec.BOOL
	Q0, Q1, Q2, Q3           iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (f *FF_D4E) INIT() { *f = FF_D4E{} }

// Execute runs the block once.
func (f *FF_D4E) Execute(now time.Time) {
	if f.RST {
		f.Q0, f.Q1, f.Q2, f.Q3 = false, false, false, false
	} else if f.CLK && !f.edge {
		f.Q0, f.Q1, f.Q2, f.Q3 = f.D0, f.D1, f.D2, f.D3
	}
	f.edge = f.CLK
}

// FF_DRE is a D flip-flop with set and reset, triggered by a rising edge of
// CLK.
type FF_DRE struct {
	SET, D, CLK, RST iec.BOOL
	Q                iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (f *FF_DRE) INIT() { *f = FF_DRE{} }

// Execute runs the block once.
func (f *FF_DRE) Execute(now time.Time) {
	if f.RST || f.SET {
		f.Q = !f.RST
	} else if f.CLK && !f.edge {
		f.Q = f.D
	}
	f.edge = f.CLK
}

// FF_JKE is a JK flip-flop with asynchronous set and reset, triggered by a
// rising edge of CLK: J and K false keep Q, J sets it, K clears it, and J
// and K together toggle it.
type FF_JKE struct {
	SET, J, CLK, K, RST iec.BOOL
	Q                   iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (f *FF_JKE) INIT() { *f = FF_JKE{} }

// Execute runs the block once.
func (f *FF_JKE) Execute(now time.Time) {
	if f.RST || f.SET {
		f.Q = !f.RST
	} else if f.CLK && !f.edge {
		if f.J != f.K {
			f.Q = f.J
		} else {
			f.Q = f.K != f.Q
		}
	}
	f.edge = f.CLK
}

// FF_RSE is an RS flip-flop triggered by edges: a rising edge of CS sets Q
// and a rising edge of CR clears it. CR has priority.
type FF_RSE struct {
	CS, CR, RST iec.BOOL
	Q           iec.BOOL

	es, er iec.BOOL
}

// INIT resets the block.
func (f *FF_RSE) INIT() { *f = FF_RSE{} }

// Execute runs the block once.
func (f *FF_RSE) Execute(now time.Time) {
	switch {
	case bool(f.RST):
		f.Q = false
	case bool(f.CR && !f.er):
		f.Q = false
	case bool(f.CS && !f.es):
		f.Q = true
	}
	f.es = f.CS
	f.er = f.CR
}

// SELECT_8 selects one of 8 outputs, stepping with rising edges of UP and
// DN. All outputs are off while E is false.
type SELECT_8 struct {
	E, SET                         iec.BOOL
	IN                             iec.BYTE
	UP, DN, RST                    iec.BOOL
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7 iec.BOOL
	STATE                          iec.INT

	lastUp, lastDn iec.BOOL
}

// INIT resets the block.
func (s *SELECT_8) INIT() { *s = SELECT_8{} }

// Execute runs the block once.
func (s *SELECT_8) Execute(now time.Time) {
	switch {
	case bool(s.RST):
		s.STATE = 0
	case bool(s.SET):
		s.STATE = iec.INT(s.IN)
	case bool(s.UP && !s.lastUp):
		s.STATE = math.INC(s.STATE, 1, 7)
	case bool(s.DN && !s.lastDn):
		s.STATE = math.INC(s.STATE, -1, 7)
	}
	s.lastUp = s.UP
	s.lastDn = s.DN
	q := [8]*iec.BOOL{&s.Q0, &s.Q1, &s.Q2, &s.Q3, &s.Q4, &s.Q5, &s.Q6, &s.Q7}
	for i, p := range q {
		*p = s.E && iec.INT(i) == s.STATE
	}
}

// SHR_4E is a 4 bit shift register with set and reset that shifts D0 in on
// a rising edge of CLK.
type SHR_4E struct {
	SET, D0, CLK, RST iec.BOOL
	Q0, Q1, Q2, Q3    iec.BOOL

	trig triggers.R_TRIG
}

// INIT resets the block.
func (s *SHR_4E) INIT() { *s = SHR_4E{} }

// Execute runs the block once.
func (s *SHR_4E) Execute(now time.Time) {
	s.trig.CLK = s.CLK
	s.trig.R_TRIG()
	if s.SET || s.RST {
		s.Q0 = !s.RST
		s.Q1, s.Q2, s.Q3 = s.Q0, s.Q0, s.Q0
	} else if s.trig.Q {
		s.Q3, s.Q2, s.Q1, s.Q0 = s.Q2, s.Q1, s.Q0, s.D0
	}
}

// SHR_4UDE is a 4 bit shift register with set and reset that shifts on a
// rising edge of CLK: D0 in up, or D3 in down if DN is true.
type SHR_4UDE struct {
	SET, D0, D3, CLK, DN, RST iec.BOOL
	Q0, Q1, Q2, Q3            iec.BOOL

	trig triggers.R_TRIG
}

// INIT resets the block.
func (s *SHR_4UDE) INIT() { *s = SHR_4UDE{} }

// Execute runs the block once.
func (s *SHR_4UDE) Execute(now time.Time) {
	s.trig.CLK = s.CLK
	s.trig.R_TRIG()
	if s.SET || s.RST {
		s.Q0 = !s.RST
		s.Q1, s.Q2, s.Q3 = s.Q0, s.Q0, s.Q0
	} else if s.trig.Q {
		if s.DN {
			s.Q0, s.Q1, s.Q2, s.Q3 = s.Q1, s.Q2, s.Q3, s.D3
		} else {
			s.Q3, s.Q2, s.Q1, s.Q0 = s.Q2, s.Q1, s.Q0, s.D0
		}
	}
}

// SHR_8PLE is an 8 bit shift register with serial input DIN, parallel load
// of DLOAD and reset. On a rising edge of CLK it shifts up, DIN into bit 0
// and bit 7 out on DOUT, or down if UP is false, DIN into bit 7 and bit 0
// out. If LOAD is true, DLOAD is loaded after the shift.
type SHR_8PLE struct {
	DIN   iec.BOOL
	DLOAD iec.BYTE
	CLK   iec.BOOL
	UP    iec.BOOL // default TRUE
	LOAD  iec.BOOL
	RST   iec.BOOL
	DOUT  iec.BOOL

	edge        iec.BOOL
	register    iec.BYTE
	initialized bool
}

// INIT resets the block and sets UP to its initial value.
func (s *SHR_8PLE) INIT() { *s = SHR_8PLE{UP: true, edge: true, initialized: true} }

// Execute runs the block once.
func (s *SHR_8PLE) Execute(now time.Time) {
	if !s.initialized {
		s.initialized = true
		s.edge = true
	}
	if s.CLK && s.edge && !s.RST {
		s.edge = false
		if s.UP {
			s.register = s.register<<1 | boolByte(s.DIN)
			s.DOUT = s.register&0x80 != 0
		} else {
			s.register = s.register>>1 | boolByte(s.DIN)<<7
			s.DOUT = s.register&1 != 0
		}
		if s.LOAD {
			s.register = s.DLOAD
			if s.UP {
				s.DOUT = s.register&0x80 != 0
			} else {
				s.DOUT = s.register&1 != 0
			}
		}
	}
	if !s.CLK {
		s.edge = true
	}
	if s.RST {
		s.register = 0
		s.DOUT = false
	}
}

// SHR_8UDE is an 8 bit shift register with set and reset that shifts on a
// rising edge of CLK: D0 in up, or D7 in down if DN is true.
type SHR_8UDE struct {
	SET, D0, D7, CLK, DN, RST      iec.BOOL
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7 iec.BOOL

	trig triggers.R_TRIG
}

// INIT resets the block.
func (s *SHR_8UDE) INIT() { *s = SHR_8UDE{} }

// Execute runs the block once.
func (s *SHR_8UDE) Execute(now time.Time) {
	s.trig.CLK = s.CLK
	s.trig.R_TRIG()
	q := [8]*iec.BOOL{&s.Q0, &s.Q1, &s.Q2, &s.Q3, &s.Q4, &s.Q5, &s.Q6, &s.Q7}
	if s.SET || s.RST {
		s.Q0 = !s.RST
		for _, p := range q[1:] {
			*p = s.Q0
		}
	} else if s.trig.Q {
		if s.DN {
			for i := 0; i < 7; i++ {
				*q[i] = *q[i+1]
			}
			s.Q7 = s.D7
		} else {
			for i := 7; i > 0; i-- {
				*q[i] = *q[i-1]
			}
			s.Q0 = s.D0
		}
	}
}

// TOGGLE is a toggle flip-flop: Q changes with every rising edge of CLK.
type TOGGLE struct {
	CLK, RST iec.BOOL
	Q        iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (t *TOGGLE) INIT() { *t = TOGGLE{} }

// Execute runs the block once.
func (t *TOGGLE) Execute(now time.Time) {
	if t.RST {
		t.Q = false
	} else if t.CLK && !t.edge {
		t.Q = !t.Q
	}
	t.edge = t.CLK
}
