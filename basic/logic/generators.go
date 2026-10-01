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
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// The timers of these blocks count in milliseconds of the PLC timer, a
// DWORD that wraps around, as OSCAT's TIME arithmetic does.

func ms(t iec.TIME) iec.DWORD { return TIME_TO_DWORD(t) }

// A_TRIG triggers Q for one scan when IN has changed by more than RES since
// the last trigger. D is the change.
type A_TRIG struct {
	IN, RES iec.REAL
	Q       iec.BOOL
	D       iec.REAL

	lastIn iec.REAL
}

// INIT resets the block.
func (a *A_TRIG) INIT() { *a = A_TRIG{} }

// Execute runs the block once.
func (a *A_TRIG) Execute(now time.Time) {
	a.D = a.IN - a.lastIn
	a.Q = ABS(a.D) > a.RES
	if a.Q {
		a.lastIn = a.IN
	}
	a.D = a.IN - a.lastIn
}

// B_TRIG triggers Q for one scan on a rising or falling edge of CLK.
type B_TRIG struct {
	CLK iec.BOOL
	Q   iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (b *B_TRIG) INIT() { *b = B_TRIG{} }

// Execute runs the block once.
func (b *B_TRIG) Execute(now time.Time) {
	b.Q = b.CLK != b.edge
	b.edge = b.CLK
}

// CLICK_CNT sets Q for one scan when IN has N pulses within the time TC
// after its first rising edge.
type CLICK_CNT struct {
	IN iec.BOOL
	N  iec.INT
	TC iec.TIME
	Q  iec.BOOL

	tx          timers.TP
	edge        iec.BOOL
	cnt         iec.INT
	initialized bool
}

// INIT resets the block.
func (c *CLICK_CNT) INIT() { *c = CLICK_CNT{cnt: -1, initialized: true} }

// Execute runs the block once.
func (c *CLICK_CNT) Execute(now time.Time) {
	if !c.initialized {
		c.initialized = true
		c.cnt = -1
	}
	c.Q = false
	switch {
	case bool(c.IN && !c.edge && !c.tx.Q):
		// A rising edge starts the count.
		c.cnt = 0
	case bool(c.tx.Q && !c.IN && c.edge):
		// Count the falling edges within TC.
		c.cnt++
	case !bool(c.tx.Q):
		c.Q = c.cnt == c.N
		c.cnt = -1
	}
	c.edge = c.IN
	c.tx.IN = c.IN
	c.tx.PT = c.TC
	c.tx.Execute(now)
}

// CLICK_DEC decodes the number of pulses of IN within the time TC after its
// first rising edge: Q0 for a single rising edge, Q1 for one pulse and so
// on up to Q3. The output stays on until IN is false.
type CLICK_DEC struct {
	IN             iec.BOOL
	TC             iec.TIME
	Q0, Q1, Q2, Q3 iec.BOOL

	tx          timers.TP
	edge        iec.BOOL
	cnt         iec.INT
	initialized bool
}

// INIT resets the block.
func (c *CLICK_DEC) INIT() { *c = CLICK_DEC{cnt: -1, initialized: true} }

// Execute runs the block once.
func (c *CLICK_DEC) Execute(now time.Time) {
	if !c.initialized {
		c.initialized = true
		c.cnt = -1
	}
	if !c.IN {
		c.Q0, c.Q1, c.Q2, c.Q3 = false, false, false, false
	}
	switch {
	case bool(c.IN && !c.edge && !c.tx.Q):
		c.cnt = 0
	case bool(c.tx.Q && !c.IN && c.edge):
		c.cnt++
	case !bool(c.tx.Q):
		switch c.cnt {
		case 0:
			c.Q0 = true
		case 1:
			c.Q1 = true
		case 2:
			c.Q2 = true
		case 3:
			c.Q3 = true
		}
		c.cnt = -1
	}
	c.edge = c.IN
	c.tx.IN = c.IN
	c.tx.PT = c.TC
	c.tx.Execute(now)
}

// CLK_DIV divides a clock: each scan CLK is true counts, and Q0 is bit 0 of
// the count, Q1 bit 1 and so on.
type CLK_DIV struct {
	CLK, RST                       iec.BOOL
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7 iec.BOOL

	cnt iec.BYTE
}

// INIT resets the block.
func (c *CLK_DIV) INIT() { *c = CLK_DIV{} }

// Execute runs the block once.
func (c *CLK_DIV) Execute(now time.Time) {
	if c.RST {
		c.cnt = 0
	} else if c.CLK {
		c.cnt++
	} else {
		return
	}
	q := [8]*iec.BOOL{&c.Q0, &c.Q1, &c.Q2, &c.Q3, &c.Q4, &c.Q5, &c.Q6, &c.Q7}
	for i, p := range q {
		*p = c.cnt>>i&1 != 0
	}
}

// CLK_N generates a pulse of one scan every 2^N milliseconds.
type CLK_N struct {
	N iec.INT
	Q iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (c *CLK_N) INIT() { *c = CLK_N{} }

// Execute runs the block once.
func (c *CLK_N) Execute(now time.Time) {
	clk := iec.BOOL(SHR(PLC_MS(now), c.N)&1 != 0)
	c.Q = clk != c.edge
	c.edge = clk
}

// CLK_PRG generates a pulse of one scan every PT, starting with a pulse.
type CLK_PRG struct {
	PT iec.TIME // default T#10ms
	Q  iec.BOOL

	init iec.BOOL
	last iec.DWORD
}

// INIT resets the block and sets PT to its initial value.
func (c *CLK_PRG) INIT() {
	*c = CLK_PRG{PT: iec.TIME(10 * time.Millisecond)}
}

// Execute runs the block once.
func (c *CLK_PRG) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !c.init {
		c.init = true
		c.last = tx - ms(c.PT)
	}
	c.Q = tx-c.last >= ms(c.PT)
	if c.Q {
		c.last = tx
	}
}

// CLK_PULSE generates N pulses of one scan every PT, or pulses without end
// if N is 0. CNT counts the pulses and RUN is true while they run.
type CLK_PULSE struct {
	PT  iec.TIME
	N   iec.INT
	RST iec.BOOL
	Q   iec.BOOL
	CNT iec.INT
	RUN iec.BOOL

	tn   iec.DWORD
	init iec.BOOL
}

// INIT resets the block.
func (c *CLK_PULSE) INIT() { *c = CLK_PULSE{} }

// Execute runs the block once.
func (c *CLK_PULSE) Execute(now time.Time) {
	tx := PLC_MS(now)
	c.Q = false
	c.RUN = c.CNT < c.N
	if !c.init || c.RST {
		c.init = true
		c.CNT = 0
		c.tn = tx - ms(c.PT)
		c.RUN = false
	} else if (c.CNT < c.N || c.N == 0) && tx-c.tn >= ms(c.PT) {
		c.CNT++
		c.Q = true
		c.tn = c.tn + ms(c.PT)
	}
}

// CYCLE_4 runs through the states 0..3, staying T0..T3 in each. After state
// 3 it starts again at 0 if S0 is true. SL loads the state SX. While E is
// false the state is 0.
type CYCLE_4 struct {
	E              iec.BOOL // default TRUE
	T0, T1, T2, T3 iec.TIME
	S0             iec.BOOL
	SX             iec.INT
	SL             iec.BOOL
	STATE          iec.INT

	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets E to its initial value.
func (c *CYCLE_4) INIT() { *c = CYCLE_4{E: true} }

// Execute runs the block once.
func (c *CYCLE_4) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !c.init {
		c.init = true
		c.last = tx
	}
	if !c.E {
		c.STATE = 0
		c.last = tx
		return
	}
	if c.SL {
		c.STATE = LIMIT(0, c.SX, 3)
		c.last = tx
		c.SL = false
		return
	}
	t := [4]iec.TIME{c.T0, c.T1, c.T2, c.T3}
	if c.STATE >= 0 && c.STATE <= 3 && tx-c.last >= ms(t[c.STATE]) {
		switch {
		case c.STATE < 3:
			c.STATE++
		case bool(c.S0):
			// If S0 is false, the sequence stops at state 3.
			c.STATE = 0
		}
		c.last = tx
	}
}

// D_TRIG triggers Q for one scan when IN has changed. X is the change.
type D_TRIG struct {
	IN iec.DWORD
	Q  iec.BOOL
	X  iec.DWORD

	lastIn iec.DWORD
}

// INIT resets the block.
func (d *D_TRIG) INIT() { *d = D_TRIG{} }

// Execute runs the block once.
func (d *D_TRIG) Execute(now time.Time) {
	d.Q = d.IN != d.lastIn
	d.X = d.IN - d.lastIn
	d.lastIn = d.IN
}

// GEN_BIT is a 4 bit pattern generator: each scan CLK is true it shifts the
// next bit of IN0..IN3 to Q0..Q3, for STEPS bits, and repeats the pattern
// REP times, or without end if REP is 0.
type GEN_BIT struct {
	IN0, IN1, IN2, IN3 iec.DWORD
	CLK                iec.BOOL
	STEPS              iec.INT
	REP                iec.INT
	RST                iec.BOOL
	Q0, Q1, Q2, Q3     iec.BOOL
	CNT                iec.INT
	RUN                iec.BOOL

	r0, r1, r2, r3 iec.DWORD
	rx             iec.INT
	initialized    bool
}

// INIT resets the block.
func (g *GEN_BIT) INIT() { *g = GEN_BIT{rx: 1, initialized: true} }

// Execute runs the block once.
func (g *GEN_BIT) Execute(now time.Time) {
	if !g.initialized {
		g.initialized = true
		g.rx = 1
	}
	if g.CLK && !g.RST {
		g.RUN = g.REP == 0 || g.rx <= g.REP
		if !g.RUN {
			return
		}
		if g.CNT == g.STEPS {
			g.CNT = 0
		}
		if g.CNT == 0 {
			g.r0, g.r1, g.r2, g.r3 = g.IN0, g.IN1, g.IN2, g.IN3
		}
		if g.CNT < g.STEPS {
			g.Q0, g.Q1, g.Q2, g.Q3 = g.r0&1 != 0, g.r1&1 != 0, g.r2&1 != 0, g.r3&1 != 0
			g.r0, g.r1, g.r2, g.r3 = g.r0>>1, g.r1>>1, g.r2>>1, g.r3>>1
		}
		g.CNT++
		if g.CNT == g.STEPS && g.REP != 0 {
			g.rx++
		}
		if g.rx > g.REP && g.REP != 0 {
			g.RUN = false
		}
	} else if g.RST {
		g.RUN = false
		g.Q0, g.Q1, g.Q2, g.Q3 = false, false, false, false
		g.r0, g.r1, g.r2, g.r3 = 0, 0, 0, 0
		g.CNT = 0
		g.rx = 1
	}
}

// GEN_SQ generates a square wave with the period PT.
type GEN_SQ struct {
	PT iec.TIME
	Q  iec.BOOL

	tn   iec.DWORD
	init iec.BOOL
}

// INIT resets the block.
func (g *GEN_SQ) INIT() { *g = GEN_SQ{} }

// Execute runs the block once.
func (g *GEN_SQ) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !g.init {
		g.init = true
		g.tn = tx
		g.Q = true
	} else if tx-g.tn >= ms(g.PT)>>1 {
		g.Q = !g.Q
		g.tn = g.tn + ms(g.PT)>>1
	}
}

// SCHEDULER sets Qn for one scan every Tn while En is true. It checks one of
// the four channels each scan, in turn.
type SCHEDULER struct {
	E0, E1, E2, E3 iec.BOOL
	T0, T1, T2, T3 iec.TIME
	Q0, Q1, Q2, Q3 iec.BOOL

	init iec.BOOL
	s    [4]iec.DWORD
	c    iec.INT
}

// INIT resets the block.
func (s *SCHEDULER) INIT() { *s = SCHEDULER{} }

// Execute runs the block once.
func (s *SCHEDULER) Execute(now time.Time) {
	tx := PLC_MS(now)
	t := [4]iec.TIME{s.T0, s.T1, s.T2, s.T3}
	e := [4]iec.BOOL{s.E0, s.E1, s.E2, s.E3}
	q := [4]*iec.BOOL{&s.Q0, &s.Q1, &s.Q2, &s.Q3}
	if !s.init {
		s.init = true
		for i := range s.s {
			s.s[i] = tx - ms(t[i])
		}
	}
	s.Q0, s.Q1, s.Q2, s.Q3 = false, false, false, false
	if c := s.c; c >= 0 && c <= 3 {
		if tx-s.s[c] >= ms(t[c]) {
			*q[c] = e[c]
			s.s[c] = tx
		}
		s.c = (c + 1) % 4
	}
}

// SCHEDULER_2 sets Qn every Cn scans while En is true, with the offset On
// scans.
type SCHEDULER_2 struct {
	E0, E1, E2, E3 iec.BOOL
	C0, C1, C2, C3 iec.UINT
	O0, O1, O2, O3 iec.UINT
	Q0, Q1, Q2, Q3 iec.BOOL

	sx iec.UINT
}

// INIT resets the block.
func (s *SCHEDULER_2) INIT() { *s = SCHEDULER_2{} }

// Execute runs the block once.
func (s *SCHEDULER_2) Execute(now time.Time) {
	due := func(c, o iec.UINT) iec.BOOL {
		// OSCAT divides by C, which must not be 0.
		if c == 0 {
			return o == 0
		}
		return s.sx%c-o == 0
	}
	s.Q0 = s.E0 && due(s.C0, s.O0)
	s.Q1 = s.E1 && due(s.C1, s.O1)
	s.Q2 = s.E2 && due(s.C2, s.O2)
	s.Q3 = s.E3 && due(s.C3, s.O3)
	s.sx++
}

// sequence runs the steps of SEQUENCE_4 and SEQUENCE_8.
type sequence struct {
	last iec.DWORD
	edge iec.BOOL
	init iec.BOOL
}

// run runs a sequencer of len(q) steps once. Each step waits up to wait[k]
// for in[k], sets its output, which clears the previous one, and holds it
// for delay[k].
func (s *sequence) run(tx iec.DWORD, in []iec.BOOL, wait, delay []iec.TIME, q []*iec.BOOL,
	start, rst, stopOnError iec.BOOL, qx, run *iec.BOOL, step *iec.INT, status *iec.BYTE) {
	n := len(q)
	if !s.init {
		s.last = tx
		s.init = true
		*status = 110
	}
	if rst {
		*step = -1
		for _, p := range q {
			*p = false
		}
		*status = 110
		*run = false
	} else if start && !s.edge {
		// A rising edge on start restarts the sequencer.
		*step = 0
		s.last = tx
		*status = 111
		for _, p := range q {
			*p = false
		}
		*run = true
	}
	s.edge = start
	if *status > 0 && *status < 100 && stopOnError {
		return
	}
	for k := 0; k < n; k++ {
		if !*run || *step != iec.INT(k) {
			continue
		}
		switch {
		case !bool(*q[k]) && bool(in[k]) && tx-s.last <= ms(wait[k]):
			if k > 0 {
				*q[k-1] = false
			}
			*q[k] = true
			s.last = tx
		case !bool(*q[k]) && tx-s.last > ms(wait[k]):
			*status = iec.BYTE(k + 1)
			if k > 0 {
				*q[k-1] = false
			}
			*run = false
		case bool(*q[k]) && tx-s.last >= ms(delay[k]):
			if k < n-1 {
				*step = iec.INT(k + 1)
				s.last = tx
			} else {
				*step = -1
				*q[k] = false
				*run = false
				*status = 110
			}
		}
	}
	*qx = false
	for _, p := range q {
		*qx = *qx || *p
	}
}

// SEQUENCE_4 is a sequencer of 4 steps. A rising edge of START starts it at
// step 0: it waits up to WAITn for INn, then sets Qn for DELAYn and moves on
// to the next step. If INn does not come, STATUS is the error n+1 and the
// sequencer stops. STATUS is 110 while waiting and 111 while running.
type SEQUENCE_4 struct {
	IN0, IN1, IN2, IN3 iec.BOOL // default TRUE
	START, RST         iec.BOOL
	WAIT0, DELAY0      iec.TIME
	WAIT1, DELAY1      iec.TIME
	WAIT2, DELAY2      iec.TIME
	WAIT3, DELAY3      iec.TIME
	STOP_ON_ERROR      iec.BOOL
	Q0, Q1, Q2, Q3     iec.BOOL
	QX                 iec.BOOL
	RUN                iec.BOOL
	STEP               iec.INT // default -1
	STATUS             iec.BYTE

	seq         sequence
	initialized bool
}

// INIT resets the block and sets IN0..IN3 and STEP to their initial values.
func (s *SEQUENCE_4) INIT() {
	*s = SEQUENCE_4{IN0: true, IN1: true, IN2: true, IN3: true, STEP: -1, initialized: true}
}

// Execute runs the block once.
func (s *SEQUENCE_4) Execute(now time.Time) {
	if !s.initialized {
		s.initialized = true
		s.STEP = -1
	}
	s.seq.run(PLC_MS(now),
		[]iec.BOOL{s.IN0, s.IN1, s.IN2, s.IN3},
		[]iec.TIME{s.WAIT0, s.WAIT1, s.WAIT2, s.WAIT3},
		[]iec.TIME{s.DELAY0, s.DELAY1, s.DELAY2, s.DELAY3},
		[]*iec.BOOL{&s.Q0, &s.Q1, &s.Q2, &s.Q3},
		s.START, s.RST, s.STOP_ON_ERROR, &s.QX, &s.RUN, &s.STEP, &s.STATUS)
}

// SEQUENCE_8 is a sequencer of 8 steps; see SEQUENCE_4.
type SEQUENCE_8 struct {
	IN0, IN1, IN2, IN3, IN4, IN5, IN6, IN7 iec.BOOL // default TRUE
	START, RST                             iec.BOOL
	WAIT0, DELAY0                          iec.TIME
	WAIT1, DELAY1                          iec.TIME
	WAIT2, DELAY2                          iec.TIME
	WAIT3, DELAY3                          iec.TIME
	WAIT4, DELAY4                          iec.TIME
	WAIT5, DELAY5                          iec.TIME
	WAIT6, DELAY6                          iec.TIME
	WAIT7, DELAY7                          iec.TIME
	STOP_ON_ERROR                          iec.BOOL
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7         iec.BOOL
	QX                                     iec.BOOL
	RUN                                    iec.BOOL
	STEP                                   iec.INT // default -1
	STATUS                                 iec.BYTE

	seq         sequence
	initialized bool
}

// INIT resets the block and sets IN0..IN7 and STEP to their initial values.
func (s *SEQUENCE_8) INIT() {
	*s = SEQUENCE_8{STEP: -1, initialized: true}
	s.IN0, s.IN1, s.IN2, s.IN3, s.IN4, s.IN5, s.IN6, s.IN7 = true, true, true, true, true, true, true, true
}

// Execute runs the block once.
func (s *SEQUENCE_8) Execute(now time.Time) {
	if !s.initialized {
		s.initialized = true
		s.STEP = -1
	}
	s.seq.run(PLC_MS(now),
		[]iec.BOOL{s.IN0, s.IN1, s.IN2, s.IN3, s.IN4, s.IN5, s.IN6, s.IN7},
		[]iec.TIME{s.WAIT0, s.WAIT1, s.WAIT2, s.WAIT3, s.WAIT4, s.WAIT5, s.WAIT6, s.WAIT7},
		[]iec.TIME{s.DELAY0, s.DELAY1, s.DELAY2, s.DELAY3, s.DELAY4, s.DELAY5, s.DELAY6, s.DELAY7},
		[]*iec.BOOL{&s.Q0, &s.Q1, &s.Q2, &s.Q3, &s.Q4, &s.Q5, &s.Q6, &s.Q7},
		s.START, s.RST, s.STOP_ON_ERROR, &s.QX, &s.RUN, &s.STEP, &s.STATUS)
}

// SEQUENCE_64 steps through the states 0..SMAX, staying PROG[state] in
// each, and then ends in state -1. A rising edge of START starts it; TRIG is
// true for one scan at each change of state.
type SEQUENCE_64 struct {
	START iec.BOOL
	SMAX  iec.INT
	PROG  [64]iec.TIME
	RST   iec.BOOL
	STATE iec.INT // default -1
	TRIG  iec.BOOL

	edge        iec.BOOL
	last        iec.DWORD
	initialized bool
}

// INIT resets the block and sets STATE to its initial value.
func (s *SEQUENCE_64) INIT() { *s = SEQUENCE_64{STATE: -1, initialized: true} }

// Execute runs the block once.
func (s *SEQUENCE_64) Execute(now time.Time) {
	if !s.initialized {
		s.initialized = true
		s.STATE = -1
	}
	tx := PLC_MS(now)
	s.TRIG = false
	switch {
	case bool(s.RST):
		s.STATE = -1
	case bool(s.START && !s.edge):
		s.STATE = 0
		s.last = tx
		s.TRIG = true
	case s.STATE >= 0:
		if tx-s.last >= ms(s.PROG[min(s.STATE, 63)]) {
			s.STATE = math.INC2(s.STATE, 1, -1, s.SMAX)
			s.last = tx
			s.TRIG = true
		}
	}
	s.edge = s.START
}

// TMAX follows IN with Q, but for at most PT. Z is true for one scan when
// the time runs out.
type TMAX struct {
	IN iec.BOOL
	PT iec.TIME
	Q  iec.BOOL
	Z  iec.BOOL

	start  iec.DWORD
	lastIn iec.BOOL
}

// INIT resets the block.
func (t *TMAX) INIT() { *t = TMAX{} }

// Execute runs the block once.
func (t *TMAX) Execute(now time.Time) {
	tx := PLC_MS(now)
	t.Z = false
	switch {
	case !bool(t.IN):
		t.Q = false
	case !bool(t.lastIn):
		t.Q = true
		t.start = tx
	case tx-t.start >= ms(t.PT) && bool(t.Q):
		t.Q = false
		t.Z = true
	}
	t.lastIn = t.IN
}

// TMIN follows IN with Q, but for at least PT.
type TMIN struct {
	IN iec.BOOL
	PT iec.TIME
	Q  iec.BOOL

	pm timers.TP
}

// INIT resets the block.
func (t *TMIN) INIT() { *t = TMIN{} }

// Execute runs the block once.
func (t *TMIN) Execute(now time.Time) {
	t.pm.IN = t.IN
	t.pm.PT = t.PT
	t.pm.Execute(now)
	t.Q = t.IN || t.pm.Q
}

// TOF_1 is an off delay: Q is true while IN is and for PT after, unless RST
// clears it.
type TOF_1 struct {
	IN  iec.BOOL
	PT  iec.TIME
	RST iec.BOOL
	Q   iec.BOOL

	start iec.DWORD
}

// INIT resets the block.
func (t *TOF_1) INIT() { *t = TOF_1{} }

// Execute runs the block once.
func (t *TOF_1) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(t.RST):
		t.Q = false
	case bool(t.IN):
		t.Q = true
		t.start = tx
	case tx-t.start >= ms(t.PT):
		t.Q = false
	}
}

// TONOF delays the rising edge of IN by T_ON and the falling edge by T_OFF.
type TONOF struct {
	IN          iec.BOOL
	T_ON, T_OFF iec.TIME
	Q           iec.BOOL

	x    timers.TON
	old  iec.BOOL
	mode iec.BOOL
}

// INIT resets the block.
func (t *TONOF) INIT() { *t = TONOF{} }

// Execute runs the block once.
func (t *TONOF) Execute(now time.Time) {
	if t.IN != t.old {
		t.x.IN = false
		t.x.PT = SEL(t.IN, t.T_OFF, t.T_ON)
		t.x.Execute(now)
		t.mode = t.IN
		t.old = t.IN
	}
	t.x.IN = true
	t.x.Execute(now)
	if t.x.Q {
		t.Q = t.mode
	}
}

// TP_1 is a pulse of PT on each rising edge of IN; a new edge restarts it.
// RST clears it.
type TP_1 struct {
	IN  iec.BOOL
	PT  iec.TIME
	RST iec.BOOL
	Q   iec.BOOL

	start iec.DWORD
	ix    iec.BOOL
}

// INIT resets the block.
func (t *TP_1) INIT() { *t = TP_1{} }

// Execute runs the block once.
func (t *TP_1) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(t.RST):
		t.Q = false
	case bool(t.IN && !t.ix):
		t.Q = true
		t.start = tx
	case tx-t.start >= ms(t.PT):
		t.Q = false
	}
	t.ix = t.IN
}

// TP_1D is a pulse of PT1 on a rising edge of IN, after which it waits PTD,
// with W true, before it takes a new edge. It clears IN and RST itself.
type TP_1D struct {
	IN  iec.BOOL
	PT1 iec.TIME
	PTD iec.TIME
	RST iec.BOOL
	Q   iec.BOOL
	W   iec.BOOL

	start iec.DWORD
	ix    iec.BOOL
}

// INIT resets the block.
func (t *TP_1D) INIT() { *t = TP_1D{} }

// Execute runs the block once.
func (t *TP_1D) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(t.RST):
		t.Q = false
		t.RST = false
		t.W = false
	case bool(t.W):
		if tx-t.start >= ms(t.PTD) {
			t.W = false
		}
	case bool(t.IN && !t.ix):
		t.Q = true
		t.start = tx
		t.IN = false
	case tx-t.start >= ms(t.PT1):
		t.Q = false
		t.W = true
		t.start = tx
	}
	t.ix = t.IN
}

// TP_X is a pulse of PT on each rising edge of IN; a new edge restarts it.
// ET is the time since the edge. Q stays false if PT is 0.
type TP_X struct {
	IN iec.BOOL
	PT iec.TIME
	Q  iec.BOOL
	ET iec.TIME

	edge  iec.BOOL
	start iec.DWORD
}

// INIT resets the block.
func (t *TP_X) INIT() { *t = TP_X{} }

// Execute runs the block once.
func (t *TP_X) Execute(now time.Time) {
	tx := PLC_MS(now)
	if t.IN && !t.edge {
		t.start = tx
		t.Q = t.PT > 0
	} else if t.Q {
		t.ET = DWORD_TO_TIME(tx - t.start)
		if t.ET >= t.PT {
			t.Q = false
			t.ET = 0
		}
	}
	t.edge = t.IN
}
