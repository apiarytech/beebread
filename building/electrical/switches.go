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

// Package electrical is the port of the OSCAT BUILDING electrical
// installation: switch and push button inputs, dimmers, lamps and timers.
package electrical

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// ms returns a TIME in milliseconds, as the PLC timer counts.
func ms(t iec.TIME) iec.DWORD { return TIME_TO_DWORD(t) }

// CLICK decodes the clicks of a push button IN, debounced and configured by
// a SW_RECONFIG: after a pause of T_PAUSE (or a press of T_SHORT) it sets
// SINGLE_, DOUBLE or TRIPLE for one, two or three clicks, until the button
// is released. Q is the debounced input. OSCAT names SINGLE_ SINGLE, a
// keyword of IEC 61131-3.
type CLICK struct {
	IN         iec.BOOL
	Q          iec.BOOL
	SINGLE_    iec.BOOL
	DOUBLE     iec.BOOL
	TRIPLE     iec.BOOL
	STATUS     iec.BYTE
	T_DEBOUNCE iec.TIME // default T#10ms
	T_SHORT    iec.TIME // default T#200ms
	T_PAUSE    iec.TIME // default T#500ms
	T_RECONFIG iec.TIME // default T#1m

	sIn   SW_RECONFIG
	state iec.INT
	last  iec.DWORD
}

// INIT resets the block and sets its inputs to their initial values.
func (c *CLICK) INIT() {
	*c = CLICK{
		T_DEBOUNCE: iec.TIME(10 * time.Millisecond), T_SHORT: iec.TIME(200 * time.Millisecond),
		T_PAUSE: iec.TIME(500 * time.Millisecond), T_RECONFIG: iec.TIME(time.Minute),
	}
}

// Execute runs the block once.
func (c *CLICK) Execute(now time.Time) {
	// Reconfiguration and debouncing.
	c.sIn.IN, c.sIn.TD, c.sIn.TR = c.IN, c.T_DEBOUNCE, c.T_RECONFIG
	c.sIn.Execute(now)

	tx := PLC_MS(now)

	// Count the clicks.
	if c.sIn.Q != c.Q {
		c.last = tx
		if c.sIn.Q {
			c.state++
		}
	}
	c.Q = c.sIn.Q

	// Set the outputs.
	if c.state > 0 {
		if c.Q && tx-c.last > ms(c.T_SHORT) || !c.Q && tx-c.last > ms(c.T_PAUSE) {
			switch c.state {
			case 1:
				c.SINGLE_ = true
				c.STATUS = 111
			case 2:
				c.DOUBLE = true
				c.STATUS = 112
			case 3:
				c.TRIPLE = true
				c.STATUS = 113
			}
			c.state = 0
		}
	} else if !c.Q {
		c.SINGLE_, c.DOUBLE, c.TRIPLE = false, false, false
		c.STATUS = 110
		c.last = tx
		c.state = 0
	}
}

// CLICK_MODE decodes a push button IN: SINGLE_ or DOUBLE is true for one
// scan after one or two short clicks within T_LONG, LONG is true while the
// button is held longer than T_LONG, and TP_LONG for the first scan of
// LONG. OSCAT names SINGLE_ SINGLE, a keyword of IEC 61131-3.
type CLICK_MODE struct {
	IN      iec.BOOL
	SINGLE_ iec.BOOL
	DOUBLE  iec.BOOL
	LONG    iec.BOOL
	TP_LONG iec.BOOL
	T_LONG  iec.TIME // default T#500ms

	timer timers.TP
	cnt   iec.INT
	last  iec.BOOL
}

// INIT resets the block and sets T_LONG to its initial value.
func (c *CLICK_MODE) INIT() { *c = CLICK_MODE{T_LONG: iec.TIME(500 * time.Millisecond)} }

// Execute runs the block once.
func (c *CLICK_MODE) Execute(now time.Time) {
	// A rising edge of IN starts the timer that decodes the clicks.
	c.timer.IN, c.timer.PT = c.IN, c.T_LONG
	c.timer.Execute(now)
	c.SINGLE_, c.DOUBLE = false, false

	if c.timer.Q {
		// Count the clicks while the timer runs.
		if !c.IN && c.last {
			c.cnt++
		}
	} else {
		switch c.cnt {
		case 1:
			c.SINGLE_ = true
		case 2:
			c.DOUBLE = true
		}
		c.cnt = 0
	}
	c.last = c.IN
	c.TP_LONG = !c.timer.Q && !c.LONG && c.IN
	c.LONG = !c.timer.Q && c.IN
}

// DEBOUNCE debounces IN: a rising edge sets Q at once, and Q stays on for TD
// after IN goes off, unless PM (pulse mode) is true, where Q is a pulse of
// one scan.
type DEBOUNCE struct {
	IN iec.BOOL
	TD iec.TIME
	PM iec.BOOL
	Q  iec.BOOL

	deb timers.TOF
}

// INIT resets the block.
func (d *DEBOUNCE) INIT() { *d = DEBOUNCE{} }

// Execute runs the block once.
func (d *DEBOUNCE) Execute(now time.Time) {
	switch {
	case !bool(d.deb.Q) && bool(d.IN):
		// a rising edge of the input
		d.Q = true
	case bool(!d.PM):
		d.Q = d.deb.Q
	default:
		d.Q = false
	}
	d.deb.IN, d.deb.PT = d.IN, d.TD
	d.deb.Execute(now)
}

// PULSE_LENGTH measures the pulses of IN: when IN goes off it sets SHORT for
// a pulse shorter than T_SHORT, LONG for one longer than T_LONG and MIDDLE
// otherwise, for one scan. LONG is also set while IN stays on past T_LONG.
type PULSE_LENGTH struct {
	IN      iec.BOOL
	SHORT   iec.BOOL
	MIDDLE  iec.BOOL
	LONG    iec.BOOL
	T_SHORT iec.TIME // default T#100ms
	T_LONG  iec.TIME // default T#1s

	tn   iec.DWORD
	edge iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (p *PULSE_LENGTH) INIT() {
	*p = PULSE_LENGTH{T_SHORT: iec.TIME(100 * time.Millisecond), T_LONG: iec.TIME(time.Second)}
}

// Execute runs the block once.
func (p *PULSE_LENGTH) Execute(now time.Time) {
	tx := PLC_MS(now)
	// The outputs are on for one scan.
	p.SHORT, p.MIDDLE, p.LONG = false, false, false
	switch {
	case bool(p.IN) && !bool(p.edge):
		// a rising edge
		p.edge = true
		p.tn = tx
	case !bool(p.IN) && bool(p.edge):
		// a falling edge
		p.edge = false
		p.tn = tx - p.tn
		switch {
		case p.tn < ms(p.T_SHORT):
			p.SHORT = true
		case p.tn > ms(p.T_LONG):
			p.LONG = true
		default:
			p.MIDDLE = true
		}
	case bool(p.IN) && tx-p.tn > ms(p.T_LONG):
		// a long pulse as soon as T_LONG is reached
		p.LONG = true
	}
}

// PULSE_T is a pulse output toggled by a push button: a rising edge of IN
// sets Q, a falling edge clears it if IN was on longer than T1, and Q goes
// off after T2. RST clears Q.
type PULSE_T struct {
	IN  iec.BOOL
	T1  iec.TIME
	T2  iec.TIME
	RST iec.BOOL
	Q   iec.BOOL

	init iec.BOOL
	last iec.DWORD
	edge iec.BOOL
}

// INIT resets the block.
func (p *PULSE_T) INIT() { *p = PULSE_T{} }

// Execute runs the block once.
func (p *PULSE_T) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(!p.init):
		p.init = true
		p.last = tx
	case bool(p.RST):
		// asynchronous reset
		p.Q = false
	case bool(p.IN) && !bool(p.edge) && !bool(p.Q):
		// a rising edge starts a pulse
		p.last = tx
		p.Q = true
	case !bool(p.IN) && bool(p.edge) && tx-p.last > ms(p.T1):
		// a falling edge ends it if IN was on longer than T1
		p.Q = false
	case tx-p.last >= ms(p.T2):
		// the pulse times out
		p.Q = false
	}
	p.edge = p.IN
}

// SW_RECONFIG debounces a switch IN for TD and configures itself to a
// normally open or normally closed contact: if Q stays on for TR, it is
// inverted. TR = 0 turns the reconfiguration off.
type SW_RECONFIG struct {
	IN iec.BOOL
	TD iec.TIME
	TR iec.TIME
	Q  iec.BOOL

	t1, t2 timers.TON
	inv    iec.BOOL
}

// INIT resets the block.
func (s *SW_RECONFIG) INIT() { *s = SW_RECONFIG{} }

// Execute runs the block once.
func (s *SW_RECONFIG) Execute(now time.Time) {
	// debounce
	s.t1.IN, s.t1.PT = s.IN, s.TD
	s.t1.Execute(now)

	if s.TR > 0 {
		s.Q = s.t1.Q != s.inv
		// A Q that stays on for TR inverts the input.
		s.t2.IN, s.t2.PT = s.Q, s.TR
		s.t2.Execute(now)
		if s.t2.Q {
			s.inv = !s.inv
		}
	} else {
		s.Q = s.t1.Q
	}
}

// SWITCH_I is a switch input that toggles Q on each change of IN, debounced
// for T_DEBOUNCE; after T_RECONFIG it reacts to the other edge, so that it
// works with switches and push buttons alike. SET and RST set and clear Q,
// and Q goes off after T_ON_MAX if that is not 0.
type SWITCH_I struct {
	SET        iec.BOOL
	IN         iec.BOOL
	RST        iec.BOOL
	T_DEBOUNCE iec.TIME // default T#10ms
	T_RECONFIG iec.TIME // default T#1s
	T_ON_MAX   iec.TIME
	Q          iec.BOOL

	state iec.BYTE
	edge  iec.BOOL
	rEdge iec.BOOL
	tOn   iec.DWORD
}

// INIT resets the block and sets its inputs to their initial values.
func (s *SWITCH_I) INIT() {
	*s = SWITCH_I{T_DEBOUNCE: iec.TIME(10 * time.Millisecond), T_RECONFIG: iec.TIME(time.Second)}
}

// Execute runs the block once.
func (s *SWITCH_I) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(s.SET) && !bool(s.RST):
		// asynchronous set and reset first
		s.Q = true
		s.tOn = tx
	case bool(s.RST):
		s.Q = false
	case s.IN != s.edge && s.state != 1:
		// an edge of the input starts the debounce time
		s.state = 1
		s.tOn = tx
	case s.state == 1 && tx-ms(s.T_DEBOUNCE) >= s.tOn:
		// the debounce time is over: the edge given by rEdge toggles Q
		s.state = 2
		if s.rEdge != s.IN {
			s.Q = !s.Q
		}
	case s.state == 2 && tx-ms(s.T_RECONFIG) >= s.tOn:
		// after T_RECONFIG, rEdge follows the input
		s.rEdge = s.IN
	}
	if s.Q && s.T_ON_MAX > 0 && tx >= s.tOn+ms(s.T_ON_MAX) {
		s.Q = false
	}
	s.edge = s.IN
}

// SWITCH_X decodes six push buttons, debounced for T_DEBOUNCE (at least
// 50 ms): IN1 or IN2 pressed alone sets Q1 or Q2 when released, IN3..IN6
// alone set Q3..Q6, and IN3..IN6 pressed while IN1 or IN2 is held set
// Q31..Q61 or Q32..Q62. The outputs are on for one scan.
type SWITCH_X struct {
	IN1, IN2, IN3, IN4, IN5, IN6 iec.BOOL
	Q1, Q2, Q3, Q4, Q5, Q6       iec.BOOL
	Q31, Q41, Q51, Q61           iec.BOOL
	Q32, Q42, Q52, Q62           iec.BOOL
	T_DEBOUNCE                   iec.TIME // default T#50ms

	init   iec.BOOL
	t      [6]timers.TOF
	x1, x2 iec.BOOL
	e1, e2 iec.BOOL
}

// INIT resets the block and sets T_DEBOUNCE to its initial value.
func (s *SWITCH_X) INIT() { *s = SWITCH_X{T_DEBOUNCE: iec.TIME(50 * time.Millisecond)} }

// Execute runs the block once.
func (s *SWITCH_X) Execute(now time.Time) {
	if !s.init {
		s.init = true
		tx := s.T_DEBOUNCE
		if tx < iec.TIME(50*time.Millisecond) {
			tx = iec.TIME(50 * time.Millisecond)
		}
		for i := range s.t {
			s.t[i].PT = tx
			s.t[i].Execute(now)
		}
	} else {
		s.Q1, s.Q2, s.Q3, s.Q4, s.Q5, s.Q6 = false, false, false, false, false, false
		s.Q31, s.Q41, s.Q51, s.Q61 = false, false, false, false
		s.Q32, s.Q42, s.Q52, s.Q62 = false, false, false, false
	}

	// Read and debounce the inputs.
	for i, in := range [6]iec.BOOL{s.IN1, s.IN2, s.IN3, s.IN4, s.IN5, s.IN6} {
		s.t[i].IN = in
		s.t[i].Execute(now)
	}
	t1, t2, t3, t4, t5, t6 := s.t[0].Q, s.t[1].Q, s.t[2].Q, s.t[3].Q, s.t[4].Q, s.t[5].Q

	// The rising edges of IN1 and IN2.
	if t1 && !s.e1 {
		s.x1 = true
	}
	if t2 && !s.e2 {
		s.x2 = true
	}

	switch {
	case bool(t1):
		switch {
		case bool(t3):
			s.Q31, s.x1 = true, false
		case bool(t4):
			s.Q41, s.x1 = true, false
		case bool(t5):
			s.Q51, s.x1 = true, false
		case bool(t6):
			s.Q61, s.x1 = true, false
		}
	case bool(t2):
		switch {
		case bool(t3):
			s.Q32, s.x2 = true, false
		case bool(t4):
			s.Q42, s.x2 = true, false
		case bool(t5):
			s.Q52, s.x2 = true, false
		case bool(t6):
			s.Q62, s.x2 = true, false
		}
	case !bool(t1) && bool(s.e1) && bool(s.x1):
		// IN1 was pressed alone
		s.Q1, s.x1 = true, false
	case !bool(t2) && bool(s.e2) && bool(s.x2):
		s.Q2, s.x2 = true, false
	case bool(t3):
		s.Q3 = true
	case bool(t4):
		s.Q4 = true
	case bool(t5):
		s.Q5 = true
	case bool(t6):
		s.Q6 = true
	}

	// The state of IN1 and IN2.
	s.e1, s.e2 = t1, t2
}
