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

// Package engineering is the port of the OSCAT BASIC engineering functions:
// automation, control, conversion, measurements, sensors, signal generators
// and signal processing.
package engineering

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/logic"
	"github.com/apiarytech/beebread/basic/math"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// ms returns a TIME in milliseconds, as the PLC timer counts.
func ms(t iec.TIME) iec.DWORD { return TIME_TO_DWORD(t) }

// DRIVER_1 is a driver for an output: a rising edge of IN sets Q, or
// toggles it if TOGGLE_MODE is true. SET and RST set and clear Q, and if
// TIMEOUT is not 0 Q goes off after TIMEOUT.
type DRIVER_1 struct {
	TOGGLE_MODE iec.BOOL
	TIMEOUT     iec.TIME
	SET         iec.BOOL
	IN          iec.BOOL
	RST         iec.BOOL
	Q           iec.BOOL

	off  timers.TON
	edge iec.BOOL
}

// INIT resets the block.
func (d *DRIVER_1) INIT() { *d = DRIVER_1{} }

// Execute runs the block once.
func (d *DRIVER_1) Execute(now time.Time) {
	if d.off.Q {
		d.Q = false
	}
	switch {
	case bool(d.RST):
		d.Q = false
	case bool(d.SET):
		d.Q = true
	case bool(d.IN && !d.edge):
		if d.TOGGLE_MODE {
			d.Q = !d.Q
		} else {
			d.Q = true
		}
	}
	d.edge = d.IN
	if d.TIMEOUT > 0 {
		d.off.IN = d.Q
		d.off.PT = d.TIMEOUT
		d.off.Execute(now)
	}
}

// DRIVER_4 is 4 DRIVER_1 with common SET, RST, TOGGLE_MODE and TIMEOUT.
type DRIVER_4 struct {
	TOGGLE_MODE        iec.BOOL
	TIMEOUT            iec.TIME
	SET                iec.BOOL
	IN0, IN1, IN2, IN3 iec.BOOL
	RST                iec.BOOL
	Q0, Q1, Q2, Q3     iec.BOOL

	d [4]DRIVER_1
}

// INIT resets the block.
func (d *DRIVER_4) INIT() { *d = DRIVER_4{} }

// Execute runs the block once.
func (d *DRIVER_4) Execute(now time.Time) {
	in := [4]iec.BOOL{d.IN0, d.IN1, d.IN2, d.IN3}
	q := [4]*iec.BOOL{&d.Q0, &d.Q1, &d.Q2, &d.Q3}
	for i := range d.d {
		x := &d.d[i]
		x.SET, x.IN, x.RST = d.SET, in[i], d.RST
		x.TOGGLE_MODE, x.TIMEOUT = d.TOGGLE_MODE, d.TIMEOUT
		x.Execute(now)
		*q[i] = x.Q
	}
}

// DRIVER_4C steps through the states SX on rising edges of IN: state 0 has
// all outputs off, and in state n the bits 0..3 of SX[n] are Q0..Q3. The
// sequence ends at the first SX that is 0, or after 7 states. RST, or
// TIMEOUT if it is not 0, return to state 0.
type DRIVER_4C struct {
	IN             iec.BOOL
	RST            iec.BOOL
	TIMEOUT        iec.TIME
	SX             [7]iec.BYTE // ARRAY[1..7], default 1, 3, 7, 15
	SN             iec.INT
	Q0, Q1, Q2, Q3 iec.BOOL

	off  timers.TON
	edge iec.BOOL
}

// INIT resets the block and sets SX to its initial value.
func (d *DRIVER_4C) INIT() { *d = DRIVER_4C{SX: [7]iec.BYTE{1, 3, 7, 15}} }

// Execute runs the block once.
func (d *DRIVER_4C) Execute(now time.Time) {
	if d.RST || d.off.Q {
		d.SN = 0
	} else if d.IN && !d.edge {
		d.SN++
		if d.SN > 7 || d.SX[d.SN-1] == 0 {
			d.SN = 0
		}
	}
	d.edge = d.IN
	var s iec.BYTE
	if d.SN > 0 {
		s = d.SX[d.SN-1]
	}
	d.Q0, d.Q1, d.Q2, d.Q3 = s&1 != 0, s&2 != 0, s&4 != 0, s&8 != 0
	if d.TIMEOUT > 0 {
		d.off.IN = d.SN > 0
		d.off.PT = d.TIMEOUT
		d.off.Execute(now)
	}
}

// FLOW_CONTROL switches a valve Q: on while IN and ENQ are true, and for
// T_AUTO after a request REQ, after which requests wait for T_DELAY. STATUS
// is 100 idle, 101 on by IN, 102 on by request and 103 reset.
type FLOW_CONTROL struct {
	IN      iec.BOOL
	REQ     iec.BOOL
	ENQ     iec.BOOL
	RST     iec.BOOL
	T_AUTO  iec.TIME // default T#1h
	T_DELAY iec.TIME // default T#23h
	Q       iec.BOOL
	STATUS  iec.BYTE

	timer logic.TP_1D
}

// INIT resets the block and sets T_AUTO and T_DELAY to their initial values.
func (f *FLOW_CONTROL) INIT() {
	*f = FLOW_CONTROL{T_AUTO: iec.TIME(time.Hour), T_DELAY: iec.TIME(23 * time.Hour)}
}

// Execute runs the block once.
func (f *FLOW_CONTROL) Execute(now time.Time) {
	f.STATUS = 100
	if f.RST {
		f.Q = false
		f.timer.RST = true
		f.timer.Execute(now)
		f.timer.RST = false
		f.STATUS = 103
	} else if f.ENQ {
		if f.IN {
			f.STATUS = 101
		}
		if f.REQ {
			f.timer.PT1 = f.T_AUTO
			f.timer.PTD = f.T_DELAY
			f.timer.IN = true
			f.STATUS = 102
		}
	}
	f.timer.Execute(now)
	f.timer.IN = false
	f.Q = f.IN && f.ENQ || f.timer.Q
}

// FT_PROFILE generates a profile over time: from VALUE_0 it ramps to
// VALUE_1 at TIME_1, to VALUE_2 at TIME_2 and VALUE_3 at TIME_3, then to
// VALUE_10 at TIME_10; while E stays true the profile holds there, and
// then it ramps to VALUE_11, VALUE_12 and VALUE_13 over the times from
// TIME_10 to TIME_11, TIME_12 and TIME_13. A rising edge of E starts it. M
// scales the times, and the output Y is the profile * K + O. RUN is true
// while it runs and ET is the time since the start.
type FT_PROFILE struct {
	K        iec.REAL // default 1.0
	O        iec.REAL
	M        iec.REAL // default 1.0
	E        iec.BOOL
	VALUE_0  iec.REAL
	TIME_1   iec.TIME
	VALUE_1  iec.REAL
	TIME_2   iec.TIME
	VALUE_2  iec.REAL
	TIME_3   iec.TIME
	VALUE_3  iec.REAL
	TIME_10  iec.TIME
	VALUE_10 iec.REAL
	TIME_11  iec.TIME
	VALUE_11 iec.REAL
	TIME_12  iec.TIME
	VALUE_12 iec.REAL
	TIME_13  iec.TIME
	VALUE_13 iec.REAL
	Y        iec.REAL
	RUN      iec.BOOL
	ET       iec.TIME

	edge       iec.BOOL
	state      iec.BYTE
	ta, tb, t0 iec.DWORD
	temp       iec.REAL
	va, vb     iec.REAL
}

// INIT resets the block and sets K and M to their initial values.
func (p *FT_PROFILE) INIT() { *p = FT_PROFILE{K: 1, M: 1} }

// Execute runs the block once.
func (p *FT_PROFILE) Execute(now time.Time) {
	tx := PLC_MS(now)
	span := func(t iec.TIME) iec.DWORD { return ms(td.MULTIME(t, p.M)) }
	if p.E && !p.edge {
		p.RUN = true
		p.ET = 0
		p.t0 = tx
		p.ta = tx
		p.tb = span(p.TIME_1)
		p.va = p.VALUE_0
		p.vb = p.VALUE_1
		p.temp = p.VALUE_0
		p.state = 1
	}
	p.edge = p.E
	if !p.RUN {
		return
	}
	// next moves on to the ramp to vb over the time t when the current
	// ramp has ended, and otherwise follows the current ramp.
	next := func(t iec.TIME, vb iec.REAL, state iec.BYTE) {
		if tx-p.ta >= p.tb {
			p.ta = p.ta + p.tb
			p.tb = span(t)
			p.va = p.vb
			p.temp = p.vb
			p.vb = vb
			p.state = state
		} else {
			p.temp = (p.vb-p.va)*iec.REAL(tx-p.ta)/iec.REAL(p.tb) + p.va
		}
	}
	switch p.state {
	case 1:
		next(p.TIME_2-p.TIME_1, p.VALUE_2, 2)
	case 2:
		next(p.TIME_3-p.TIME_2, p.VALUE_3, 3)
	case 3:
		next(p.TIME_10-p.TIME_3, p.VALUE_10, 4)
	case 4:
		next(p.TIME_11-p.TIME_10, p.VALUE_11, SEL[iec.BYTE](p.E, 6, 5))
	case 5:
		// Hold while E is true.
		if p.E {
			p.ta = tx
		} else {
			p.state = 6
		}
	case 6:
		next(p.TIME_12-p.TIME_11, p.VALUE_12, 7)
	case 7:
		next(p.TIME_13-p.TIME_12, p.VALUE_13, 8)
	case 8:
		if tx-p.ta >= p.tb {
			p.temp = p.VALUE_13
			p.RUN = false
		} else {
			p.temp = (p.vb-p.va)*iec.REAL(tx-p.ta)/iec.REAL(p.tb) + p.va
		}
	}
	p.Y = p.temp*p.K + p.O
	p.ET = DWORD_TO_TIME(tx - p.t0)
}

// INC_DEC decodes an incremental encoder with the channels CHA and CHB, with
// four counts per pulse. DIR is true when it turns up.
type INC_DEC struct {
	CHA, CHB iec.BOOL
	RST      iec.BOOL
	DIR      iec.BOOL
	CNT      iec.INT

	edgea, edgeb iec.BOOL
}

// INIT resets the block.
func (i *INC_DEC) INIT() { *i = INC_DEC{} }

// Execute runs the block once.
func (i *INC_DEC) Execute(now time.Time) {
	axb := i.CHA != i.CHB
	clka := i.CHA != i.edgea
	i.edgea = i.CHA
	clkb := i.CHB != i.edgeb
	i.edgeb = i.CHB
	clk := clka || clkb
	if axb && clka {
		i.DIR = true
	}
	if axb && clkb {
		i.DIR = false
	}
	if clk && bool(i.DIR) {
		i.CNT++
	}
	if clk && !bool(i.DIR) {
		i.CNT--
	}
	if i.RST {
		i.CNT = 0
	}
}

// INTERLOCK switches Q1 with I1 and Q2 with I2, but an output only while
// the other input has been off for TL.
type INTERLOCK struct {
	I1, I2 iec.BOOL
	TL     iec.TIME
	Q1, Q2 iec.BOOL

	t1, t2 timers.TOF
}

// INIT resets the block.
func (i *INTERLOCK) INIT() { *i = INTERLOCK{} }

// Execute runs the block once.
func (i *INTERLOCK) Execute(now time.Time) {
	i.t1.IN, i.t1.PT = i.I1, i.TL
	i.t1.Execute(now)
	i.t2.IN, i.t2.PT = i.I2, i.TL
	i.t2.Execute(now)
	i.Q1 = i.I1 && !i.t2.Q
	i.Q2 = i.I2 && !i.t1.Q
}

// INTERLOCK_4 reports which of 4 switches I0..I3 is pressed on OUT, as the
// bit of the switch, while E is true. TP is true for one scan when OUT
// changes. MODE 0 shows the inputs as they are, 1 the highest input, 2 the
// input pressed last and 3 the first input pressed, which locks the others
// out.
type INTERLOCK_4 struct {
	I0, I1, I2, I3 iec.BOOL
	E              iec.BOOL
	MODE           iec.INT
	OUT            iec.BYTE
	TP             iec.BOOL

	last, old iec.BYTE
	lmode     iec.INT
}

// INIT resets the block.
func (l *INTERLOCK_4) INIT() { *l = INTERLOCK_4{} }

// highest returns the highest bit of in, or in if it is bit 0 or none.
func highest(in iec.BYTE) iec.BYTE {
	switch {
	case in&8 != 0:
		return 8
	case in&4 != 0:
		return 4
	case in&2 != 0:
		return 2
	}
	return in
}

// Execute runs the block once.
func (l *INTERLOCK_4) Execute(now time.Time) {
	if !l.E {
		l.OUT, l.last, l.old, l.lmode = 0, 0, 0, 0
		l.TP = false
		return
	}
	if l.MODE != l.lmode {
		l.OUT, l.last, l.old = 0, 0, 0
		l.lmode = l.MODE
	}
	in := logic.BYTE_OF_BIT(l.I0, l.I1, l.I2, l.I3, false, false, false, false)
	if in != l.last {
		switch l.MODE {
		case 0:
			l.OUT = in
		case 1:
			l.OUT = highest(in)
		case 2:
			l.last = (in ^ l.last) & in
			l.OUT = highest(l.last)
		case 3:
			if l.OUT&in == 0 {
				l.OUT = highest(in)
			}
		}
		l.last = in
	}
	l.TP = l.OUT != l.old
	l.old = l.OUT
}

// MANUAL overrides a digital signal: ON forces it on and OFF forces it off;
// otherwise it is IN.
func MANUAL(in, on, off iec.BOOL) iec.BOOL {
	return !off && (in || on)
}

// MANUAL_1 overrides a digital signal: while MAN is false Q is IN, and while
// MAN is true Q is M_I, or set by a rising edge of SET or cleared by a
// rising edge of RST. STATUS is 100 automatic, 101 set, 102 reset and 103
// manual.
type MANUAL_1 struct {
	IN, MAN, M_I, SET, RST iec.BOOL
	Q                      iec.BOOL
	STATUS                 iec.BYTE

	sEdge, rEdge, edge iec.BOOL
}

// INIT resets the block.
func (m *MANUAL_1) INIT() { *m = MANUAL_1{} }

// Execute runs the block once.
func (m *MANUAL_1) Execute(now time.Time) {
	switch {
	case !bool(m.MAN):
		m.Q = m.IN
		m.STATUS = 100
		m.edge = false
	case bool(!m.sEdge && m.SET):
		m.Q = true
		m.edge = true
		m.STATUS = 101
	case bool(!m.rEdge && m.RST):
		m.Q = false
		m.edge = true
		m.STATUS = 102
	case !bool(m.edge):
		m.Q = m.M_I
		m.STATUS = 103
	}
	m.sEdge = m.SET
	m.rEdge = m.RST
}

// MANUAL_2 overrides a digital signal while ENA is true: Q is IN, forced on
// by ON, forced off by OFF, or MAN if both ON and OFF are true. STATUS is
// 100, 101, 102 and 103 for those, and 104 while ENA is false and Q is off.
type MANUAL_2 struct {
	IN, ENA, ON, OFF, MAN iec.BOOL
	Q                     iec.BOOL
	STATUS                iec.BYTE
}

// INIT resets the block.
func (m *MANUAL_2) INIT() { *m = MANUAL_2{} }

// Execute runs the block once.
func (m *MANUAL_2) Execute(now time.Time) {
	switch {
	case !bool(m.ENA):
		m.Q, m.STATUS = false, 104
	case !bool(m.ON) && !bool(m.OFF):
		m.Q, m.STATUS = m.IN, 100
	case bool(m.ON) && !bool(m.OFF):
		m.Q, m.STATUS = true, 101
	case !bool(m.ON) && bool(m.OFF):
		m.Q, m.STATUS = false, 102
	default:
		m.Q, m.STATUS = m.MAN, 103
	}
}

// MANUAL_4 overrides 4 digital signals: while MAN is false Q0..Q3 are
// I0..I3, and while MAN is true they are M0..M3, until a rising edge of STP
// steps through the outputs one at a time. STATUS is 100 automatic, 101
// manual and 110..113 for the step.
type MANUAL_4 struct {
	I0, I1, I2, I3 iec.BOOL
	MAN, STP       iec.BOOL
	M0, M1, M2, M3 iec.BOOL
	Q0, Q1, Q2, Q3 iec.BOOL
	STATUS         iec.BYTE

	edge iec.BOOL
	pos  iec.INT
	tog  iec.BOOL
}

// INIT resets the block.
func (m *MANUAL_4) INIT() { *m = MANUAL_4{} }

// Execute runs the block once.
func (m *MANUAL_4) Execute(now time.Time) {
	if m.MAN {
		if !m.tog {
			m.Q0, m.Q1, m.Q2, m.Q3 = m.M0, m.M1, m.M2, m.M3
			m.STATUS = 101
		}
		if m.STP && !m.edge {
			m.tog = true
			m.Q0, m.Q1, m.Q2, m.Q3 = m.pos == 0, m.pos == 1, m.pos == 2, m.pos == 3
			m.STATUS = 110 + iec.BYTE(m.pos)
			m.pos = math.INC(m.pos, 1, 3)
		}
	} else {
		m.Q0, m.Q1, m.Q2, m.Q3 = m.I0, m.I1, m.I2, m.I3
		m.STATUS = 100
		m.tog = false
		m.pos = 0
	}
	m.edge = m.STP
}

// PARSET selects one of 4 parameter sets, Xn1..Xn4 for the set n that A1,
// A0 address, on P1..P4. If TC is not 0, the outputs ramp to a new set over
// TC.
type PARSET struct {
	A0, A1             iec.BOOL
	X01, X02, X03, X04 iec.REAL
	X11, X12, X13, X14 iec.REAL
	X21, X22, X23, X24 iec.REAL
	X31, X32, X33, X34 iec.REAL
	TC                 iec.TIME
	P1, P2, P3, P4     iec.REAL

	x     [4][4]iec.REAL
	s     [4]iec.REAL
	last  iec.DWORD
	start iec.BOOL
	set   iec.BYTE
	init  iec.BOOL
}

// INIT resets the block.
func (p *PARSET) INIT() { *p = PARSET{} }

// Execute runs the block once.
func (p *PARSET) Execute(now time.Time) {
	tx := PLC_MS(now)
	out := [4]*iec.REAL{&p.P1, &p.P2, &p.P3, &p.P4}
	if !p.init {
		p.set = iec.BYTE(BOOL_TO_INT(!p.A0))
		p.init = true
		p.x = [4][4]iec.REAL{
			{p.X01, p.X02, p.X03, p.X04},
			{p.X11, p.X12, p.X13, p.X14},
			{p.X21, p.X22, p.X23, p.X24},
			{p.X31, p.X32, p.X33, p.X34},
		}
		p.P1, p.P2, p.P3, p.P4 = p.X01, p.X02, p.X03, p.X04
	}
	a0, a1 := iec.BYTE(BOOL_TO_INT(p.A0)), iec.BYTE(BOOL_TO_INT(p.A1))
	tc := ms(p.TC)
	switch {
	case a0 != p.set&1 || a1 != p.set>>1&1:
		p.set = a0 | a1<<1
		if p.TC > 0 {
			p.start = true
			p.last = tx
			for i, o := range out {
				p.s[i] = (p.x[p.set][i] - *o) / iec.REAL(tc)
			}
		}
	case bool(p.start) && tx-p.last < tc:
		for i, o := range out {
			*o = p.x[p.set][i] - p.s[i]*iec.REAL(tc-tx+p.last)
		}
	default:
		p.start = false
		for i, o := range out {
			*o = p.x[p.set][i]
		}
	}
}

// PARSET2 selects one of 4 parameter sets by the value of X: set 0 while
// |X| < L1, set 1 while |X| < L2, set 2 while |X| < L3 and set 3 above; see
// PARSET.
type PARSET2 struct {
	X                  iec.REAL
	X01, X02, X03, X04 iec.REAL
	X11, X12, X13, X14 iec.REAL
	X21, X22, X23, X24 iec.REAL
	X31, X32, X33, X34 iec.REAL
	L1, L2, L3         iec.REAL
	TC                 iec.TIME
	P1, P2, P3, P4     iec.REAL

	pset PARSET
	init iec.BOOL
}

// INIT resets the block.
func (p *PARSET2) INIT() { *p = PARSET2{} }

// Execute runs the block once.
func (p *PARSET2) Execute(now time.Time) {
	s := &p.pset
	if !p.init {
		p.init = true
		s.TC = p.TC
		s.X01, s.X02, s.X03, s.X04 = p.X01, p.X02, p.X03, p.X04
		s.X11, s.X12, s.X13, s.X14 = p.X11, p.X12, p.X13, p.X14
		s.X21, s.X22, s.X23, s.X24 = p.X21, p.X22, p.X23, p.X24
		s.X31, s.X32, s.X33, s.X34 = p.X31, p.X32, p.X33, p.X34
		s.Execute(now)
	}
	x := ABS(p.X)
	switch {
	case x < p.L1:
		s.A0, s.A1 = false, false
	case x < p.L2:
		s.A0, s.A1 = true, false
	case x < p.L3:
		s.A0, s.A1 = false, true
	default:
		s.A0, s.A1 = true, true
	}
	s.Execute(now)
	p.P1, p.P2, p.P3, p.P4 = s.P1, s.P2, s.P3, s.P4
}

// SIGNAL shows the bit pattern SIG on Q while IN is true, one bit each TS,
// or each 128 ms if TS is 0.
type SIGNAL struct {
	IN  iec.BOOL
	SIG iec.BYTE
	TS  iec.TIME
	Q   iec.BOOL
}

// INIT resets the block.
func (s *SIGNAL) INIT() { *s = SIGNAL{} }

// Execute runs the block once.
func (s *SIGNAL) Execute(now time.Time) {
	if !s.IN {
		s.Q = false
		return
	}
	tx := PLC_MS(now)
	var step iec.BYTE
	if s.TS > 0 {
		step = iec.BYTE(tx / ms(s.TS) & 7)
	} else {
		step = iec.BYTE(tx >> 7 & 7)
	}
	s.Q = iec.BYTE(1)<<step&s.SIG > 0
}

// SIGNAL_4 shows the pattern S1..S4 of the first input IN1..IN4 that is
// true on Q; see SIGNAL.
type SIGNAL_4 struct {
	IN1, IN2, IN3, IN4 iec.BOOL
	TS                 iec.TIME
	S1, S2, S3, S4     iec.BYTE // default 2#1111_1111, 2#1111_0000, 2#1010_1010, 2#1010_0000
	Q                  iec.BOOL

	sig SIGNAL
}

// INIT resets the block and sets S1..S4 to their initial values.
func (s *SIGNAL_4) INIT() {
	*s = SIGNAL_4{S1: 0b1111_1111, S2: 0b1111_0000, S3: 0b1010_1010, S4: 0b1010_0000}
}

// Execute runs the block once.
func (s *SIGNAL_4) Execute(now time.Time) {
	s.sig.IN = true
	s.sig.TS = s.TS
	switch {
	case bool(s.IN1):
		s.sig.SIG = s.S1
	case bool(s.IN2):
		s.sig.SIG = s.S2
	case bool(s.IN3):
		s.sig.SIG = s.S3
	case bool(s.IN4):
		s.sig.SIG = s.S4
	default:
		s.sig.IN = false
	}
	s.sig.Execute(now)
	s.Q = s.sig.Q
}

// SRAMP follows X with Y with a speed V limited to VU_MAX up and VD_MAX
// down and an acceleration limited to A_UP and A_DN, within LIMIT_LOW and
// LIMIT_HIGH.
type SRAMP struct {
	X          iec.REAL
	A_UP       iec.REAL
	A_DN       iec.REAL
	VU_MAX     iec.REAL
	VD_MAX     iec.REAL
	LIMIT_HIGH iec.REAL
	LIMIT_LOW  iec.REAL
	RST        iec.BOOL
	Y          iec.REAL
	V          iec.REAL

	cycleTime TC_S
	init      iec.BOOL
}

// INIT resets the block.
func (s *SRAMP) INIT() { *s = SRAMP{} }

// Execute runs the block once.
func (s *SRAMP) Execute(now time.Time) {
	s.cycleTime.Execute(now)
	tc := s.cycleTime.TC
	s.A_UP = max(0.0, s.A_UP)
	s.A_DN = min(0.0, s.A_DN)
	s.VU_MAX = max(0.0, s.VU_MAX)
	s.VD_MAX = min(0.0, s.VD_MAX)
	switch {
	case bool(s.RST || !s.init):
		s.init = true
		s.Y = 0.0
		s.V = 0.0
	case s.X == s.Y:
		s.V = 0.0
	case s.X > s.Y:
		s.V = min(s.V+s.A_UP*tc, s.VU_MAX)
		s.V = min(SQRT((s.Y-s.X)*2.0*s.A_DN), s.V)
		s.Y = LIMIT(s.LIMIT_LOW, s.Y+min(s.V*tc, s.X-s.Y), s.LIMIT_HIGH)
	case s.X < s.Y:
		s.V = max(s.V+s.A_DN*tc, s.VD_MAX)
		s.V = max(-SQRT((s.Y-s.X)*2.0*s.A_UP), s.V)
		s.Y = LIMIT(s.LIMIT_LOW, s.Y+max(s.V*tc, s.X-s.Y), s.LIMIT_HIGH)
	}
}

// TUNE sets Y with the keys SU and SD: a short press steps Y by SS, a press
// longer than T1 ramps it by S1 per second, and longer than T2 by S2 per
// second. SET sets Y to SET_VAL and RST to RST_VAL. Y stays within LIMIT_L
// and LIMIT_H.
type TUNE struct {
	SET     iec.BOOL
	SU, SD  iec.BOOL
	RST     iec.BOOL
	SS      iec.REAL // default 0.1
	LIMIT_L iec.REAL
	LIMIT_H iec.REAL // default 100.0
	RST_VAL iec.REAL
	SET_VAL iec.REAL // default 100.0
	T1      iec.TIME // default T#500ms
	T2      iec.TIME // default T#2s
	S1      iec.REAL // default 2.0
	S2      iec.REAL // default 10.0
	Y       iec.REAL

	start, start2   iec.DWORD
	state           iec.INT
	step, speed     iec.REAL
	yStart, yStart2 iec.REAL
}

// INIT resets the block and sets its inputs to their initial values.
func (t *TUNE) INIT() {
	*t = TUNE{SS: 0.1, LIMIT_H: 100, SET_VAL: 100, T1: iec.TIME(500 * time.Millisecond),
		T2: iec.TIME(2 * time.Second), S1: 2, S2: 10}
}

// Execute runs the block once.
func (t *TUNE) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case bool(t.RST):
		t.Y = t.RST_VAL
		t.state = 0
	case bool(t.SET):
		t.Y = t.SET_VAL
		t.state = 0
	case t.state > 0:
		in := SEL(t.state == 1, t.SD, t.SU)
		switch {
		case !bool(in) && tx-t.start <= ms(t.T1):
			t.Y = t.yStart + t.step
			t.state = 0
		case bool(in) && tx-t.start >= ms(t.T2):
			t.Y = t.yStart2 + iec.REAL(tx-t.start2)*t.S2/t.speed
		case bool(in) && tx-t.start >= ms(t.T1):
			t.Y = t.yStart + iec.REAL(tx-t.start-ms(t.T1))*t.S1/t.speed
			t.start2 = tx
			t.yStart2 = t.Y
		case !bool(in):
			t.state = 0
		}
	case bool(t.SU):
		t.state, t.start, t.step, t.speed, t.yStart = 1, tx, t.SS, 1000.0, t.Y
	case bool(t.SD):
		t.state, t.start, t.step, t.speed, t.yStart = 2, tx, -t.SS, -1000.0, t.Y
	}
	t.Y = LIMIT(t.LIMIT_L, t.Y, t.LIMIT_H)
}

// TUNE2 sets Y with the keys SU and SD, slow, and FU and FD, fast: a short
// press steps Y by SS or FS, and a press longer than TR ramps it by S1 or S2
// per second. SET sets Y to SET_VAL and RST to RST_VAL. Y stays within
// LIMIT_L and LIMIT_H.
type TUNE2 struct {
	SET     iec.BOOL
	SU, SD  iec.BOOL
	FU, FD  iec.BOOL
	RST     iec.BOOL
	SS      iec.REAL // default 0.1
	FS      iec.REAL // default 5.0
	LIMIT_L iec.REAL
	LIMIT_H iec.REAL // default 100.0
	RST_VAL iec.REAL
	SET_VAL iec.REAL // default 100.0
	TR      iec.TIME // default T#500ms
	S1      iec.REAL // default 2.0
	S2      iec.REAL // default 10.0
	Y       iec.REAL

	start       iec.DWORD
	state       iec.INT
	in          iec.BOOL
	step, speed iec.REAL
	yStart      iec.REAL
}

// INIT resets the block and sets its inputs to their initial values.
func (t *TUNE2) INIT() {
	*t = TUNE2{SS: 0.1, FS: 5, LIMIT_H: 100, SET_VAL: 100, TR: iec.TIME(500 * time.Millisecond), S1: 2, S2: 10}
}

// Execute runs the block once.
func (t *TUNE2) Execute(now time.Time) {
	tx := PLC_MS(now)
	begin := func(state iec.INT, step, speed iec.REAL) {
		t.state, t.start, t.step, t.speed, t.yStart = state, tx, step, speed, t.Y
	}
	switch {
	case bool(t.RST):
		t.Y = t.RST_VAL
		t.state = 0
	case bool(t.SET):
		t.Y = t.SET_VAL
		t.state = 0
	case t.state > 0:
		switch t.state {
		case 1:
			t.in = t.SU
		case 2:
			t.in = t.SD
		case 3:
			t.in = t.FU
		case 4:
			t.in = t.FD
		}
		switch {
		case !bool(t.in) && tx-t.start <= ms(t.TR):
			t.Y = t.yStart + t.step
			t.state = 0
		case bool(t.in) && tx-t.start >= ms(t.TR):
			t.Y = t.yStart + iec.REAL(tx-t.start-ms(t.TR))*t.speed
		case !bool(t.in):
			t.state = 0
		}
	case bool(t.SU):
		begin(1, t.SS, t.S1*1.0e-3)
	case bool(t.SD):
		begin(2, -t.SS, -t.S1*1.0e-3)
	case bool(t.FU):
		begin(3, t.FS, t.S2*1.0e-3)
	case bool(t.FD):
		begin(4, -t.FS, -t.S2*1.0e-3)
	}
	t.Y = LIMIT(t.LIMIT_L, t.Y, t.LIMIT_H)
}
