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

// RMP_B_ ramps the byte RMP up to 255 if DIR is true, or down to 0, in the
// time TR for a full ramp, while E is true.
type RMP_B_ struct {
	DIR iec.BOOL
	E   iec.BOOL // default TRUE
	TR  iec.TIME
	RMP *iec.BYTE

	tl, tn  iec.DWORD
	init    iec.BOOL
	lastDir iec.BOOL
	start   iec.BYTE
}

// INIT resets the block and sets E to its initial value.
func (r *RMP_B_) INIT() { *r = RMP_B_{RMP: r.RMP, E: true} }

// Execute runs the block once.
func (r *RMP_B_) Execute(now time.Time) {
	if r.RMP == nil {
		return
	}
	tx := PLC_MS(now)
	if r.E && r.init && r.DIR == r.lastDir && *r.RMP != SEL[iec.BYTE](r.DIR, 0, 255) && ms(r.TR) == r.tn {
		*r.RMP = math.FRMP_B(r.start, r.DIR, DWORD_TO_TIME(tx-r.tl), r.TR)
	} else {
		r.init = true
		r.tl = tx
		r.tn = ms(r.TR)
		r.start = *r.RMP
	}
	r.lastDir = r.DIR
}

// RMP_NEXT_ ramps the byte OUT towards IN, up in the time TR or down in TF
// for a full ramp, with a lock of TL between the directions. DIR, UP and DN
// show the ramp's direction.
type RMP_NEXT_ struct {
	E   iec.BOOL // default TRUE
	IN  iec.BYTE
	TR  iec.TIME
	TF  iec.TIME
	TL  iec.TIME
	DIR iec.BOOL
	UP  iec.BOOL
	DN  iec.BOOL
	OUT *iec.BYTE

	rmx         RMP_B_
	dirx        TREND_DW
	tLock       timers.TP
	xen, xdir   iec.BOOL
	initialized bool
}

// INIT resets the block and sets E to its initial value.
func (r *RMP_NEXT_) INIT() {
	*r = RMP_NEXT_{OUT: r.OUT, E: true, initialized: true}
	r.rmx.INIT()
}

// Execute runs the block once.
func (r *RMP_NEXT_) Execute(now time.Time) {
	if !r.initialized {
		r.initialized = true
		r.rmx.INIT()
	}
	if r.OUT == nil {
		return
	}
	lock := func(in iec.BOOL) {
		r.tLock.IN = in
		r.tLock.Execute(now)
	}
	r.dirx.X = iec.DWORD(r.IN)
	r.dirx.Execute(now)
	r.tLock.PT = r.TL
	lock(false)
	out := *r.OUT
	switch {
	case bool(r.dirx.TU) && out < r.IN:
		if !r.xdir && r.xen {
			lock(true)
		}
		r.xen, r.xdir = true, true
	case bool(r.dirx.TD) && out > r.IN:
		if r.xdir && r.xen {
			lock(true)
		}
		r.xen, r.xdir = true, false
	case bool(r.xen):
		if r.xdir && out >= r.IN || !r.xdir && out <= r.IN {
			r.xen = false
			if r.TL > 0 {
				lock(true)
			}
		}
	}
	if !r.tLock.Q && r.xen {
		r.UP, r.DIR, r.DN = r.xdir, r.xdir, !r.xdir
	} else {
		r.UP, r.DN = false, false
	}
	r.rmx.RMP = r.OUT
	r.rmx.E = r.E && (r.UP || r.DN)
	r.rmx.DIR = r.DIR
	r.rmx.TR = SEL(r.DIR, r.TF, r.TR)
	r.rmx.Execute(now)
}

// RMP_W_ ramps the word RMP up to 65535 if DIR is true, or down to 0, in
// the time TR for a full ramp, while E is true.
type RMP_W_ struct {
	DIR iec.BOOL
	E   iec.BOOL // default TRUE
	TR  iec.TIME
	RMP *iec.WORD

	tl      iec.DWORD
	init    iec.BOOL
	lastDir iec.BOOL
}

// INIT resets the block and sets E to its initial value.
func (r *RMP_W_) INIT() { *r = RMP_W_{RMP: r.RMP, E: true} }

// Execute runs the block once.
func (r *RMP_W_) Execute(now time.Time) {
	if r.RMP == nil {
		return
	}
	tx := PLC_MS(now)
	if !r.E || !r.init {
		r.tl = tx
		r.init = true
		return
	}
	if r.DIR != r.lastDir {
		r.tl = tx
		r.lastDir = r.DIR
	}
	var step iec.DINT = 65535
	if r.TR > 0 {
		step = iec.DINT(((tx - r.tl) << 16) / ms(r.TR))
	}
	if step > 0 {
		r.tl = tx
		if !r.DIR {
			step = -step
		}
		*r.RMP = iec.WORD(LIMIT(0, iec.DINT(*r.RMP)+step, 65535))
	}
}

// GEN_PULSE generates a pulse wave, PTH high and PTL low, while ENQ is
// true.
type GEN_PULSE struct {
	ENQ iec.BOOL // default TRUE
	PTH iec.TIME
	PTL iec.TIME
	Q   iec.BOOL

	tn   iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets ENQ to its initial value.
func (g *GEN_PULSE) INIT() { *g = GEN_PULSE{ENQ: true} }

// Execute runs the block once.
func (g *GEN_PULSE) Execute(now time.Time) {
	if !g.ENQ {
		g.Q = false
		g.init = false
		return
	}
	tx := PLC_MS(now)
	if !g.init {
		g.init = true
		g.tn = tx
	}
	if t := ms(SEL(g.Q, g.PTL, g.PTH)); tx-g.tn >= t {
		g.tn += t
		g.Q = !g.Q
	}
}

// GEN_PW2 generates a pulse wave, TH1 high and TL1 low, or TH2 and TL2 if TS
// is true, while ENQ is true. TH and TL are the time in the high or low
// state.
type GEN_PW2 struct {
	ENQ      iec.BOOL
	TH1, TL1 iec.TIME
	TH2, TL2 iec.TIME
	TS       iec.BOOL
	Q        iec.BOOL
	TH, TL   iec.TIME

	start iec.DWORD
	init  iec.BOOL
}

// INIT resets the block.
func (g *GEN_PW2) INIT() { *g = GEN_PW2{} }

// Execute runs the block once.
func (g *GEN_PW2) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !g.init {
		g.start = tx
		g.init = true
		g.TH, g.TL = 0, 0
	}
	tHigh, tLow := g.TH1, g.TL1
	if g.TS {
		tHigh, tLow = g.TH2, g.TL2
	}
	if !g.ENQ {
		g.Q = false
		g.TH, g.TL = 0, 0
		g.start = tx
		return
	}
	et := tx - g.start
	if !g.Q {
		if et >= ms(tLow) {
			g.Q = true
			g.start = tx
			g.TL = 0
		} else {
			g.TL = DWORD_TO_TIME(et)
		}
	} else if et >= ms(tHigh) {
		g.Q = false
		g.start = tx
		g.TH = 0
	} else {
		g.TH = DWORD_TO_TIME(et)
	}
}

// GEN_RDM generates a random signal OUT from OS - AM/2 to OS + AM/2 that
// changes every PT. Q is true for one scan when it changes.
type GEN_RDM struct {
	PT  iec.TIME
	AM  iec.REAL // default 1
	OS  iec.REAL
	Q   iec.BOOL
	OUT iec.REAL

	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block and sets AM to its initial value.
func (g *GEN_RDM) INIT() { *g = GEN_RDM{AM: 1} }

// Execute runs the block once.
func (g *GEN_RDM) Execute(now time.Time) {
	tx := PLC_MS(now) - g.last
	if !g.init {
		g.init = true
		g.last = tx
		tx = 0
	}
	if tx >= ms(g.PT) {
		g.last += ms(g.PT)
		g.OUT = g.AM*(math.RDM(0)-0.5) + g.OS
		g.Q = true
	} else {
		g.Q = false
	}
}

// GEN_RDT generates a pulse of TP_Q at random times while ENABLE is true.
//
// OSCAT means the times to be between MIN_TIME_MS and MAX_TIME_MS but
// calculates them from 0 to MAX_TIME_MS; the port keeps that.
type GEN_RDT struct {
	ENABLE      iec.BOOL // default TRUE
	MIN_TIME_MS iec.TIME // default T#1s
	MAX_TIME_MS iec.TIME // default T#1.2s
	TP_Q        iec.TIME // default T#100ms
	XQ          iec.BOOL

	tonRDMTimer timers.TON
	tofXQ       timers.TOF
	tRDMTime    iec.TIME
	rRDMTime    iec.REAL
}

// INIT resets the block and sets its inputs to their initial values.
func (g *GEN_RDT) INIT() {
	*g = GEN_RDT{ENABLE: true, MIN_TIME_MS: iec.TIME(time.Second),
		MAX_TIME_MS: iec.TIME(1200 * time.Millisecond), TP_Q: iec.TIME(100 * time.Millisecond)}
}

// Execute runs the block once.
func (g *GEN_RDT) Execute(now time.Time) {
	g.tonRDMTimer.IN, g.tonRDMTimer.PT = g.ENABLE, g.tRDMTime
	g.tonRDMTimer.Execute(now)
	g.tofXQ.IN, g.tofXQ.PT = g.tonRDMTimer.Q, g.TP_Q
	g.tofXQ.Execute(now)
	g.XQ = g.tofXQ.Q
	if g.tonRDMTimer.Q {
		g.XQ = true
		g.rRDMTime = math.RDM(g.rRDMTime)
		span := iec.REAL(iec.DINT(ms(g.MAX_TIME_MS-g.MIN_TIME_MS)) + iec.DINT(ms(g.MIN_TIME_MS)))
		g.tRDMTime = REAL_TO_TIME(g.rRDMTime * span)
		g.tonRDMTimer.IN = false
		g.tonRDMTimer.Execute(now)
	}
}

// generator is the timing the wave generators GEN_RMP, GEN_SIN and GEN_SQR
// share: the time tx within the period PT.
type generator struct {
	last iec.DWORD
	init iec.BOOL
}

// period returns the time within the period pt.
func (g *generator) period(now time.Time, pt iec.TIME) iec.DWORD {
	tx := PLC_MS(now) - g.last
	if !g.init {
		g.init = true
		g.last = tx
		tx = 0
	}
	if tx >= ms(pt) {
		g.last += ms(pt)
		tx -= ms(pt)
	}
	return tx
}

// delay returns DL within 0..1, as OSCAT limits it.
func delay(dl iec.REAL) iec.REAL {
	dl = math.MODR(dl, 1.0)
	if dl < 0.0 {
		dl = 1.0 - dl
	}
	return dl
}

// GEN_RMP generates a sawtooth wave OUT from OS to OS + AM with the period
// PT, delayed by DL periods. Q is true for one scan when a ramp starts.
type GEN_RMP struct {
	PT  iec.TIME // default T#1s
	AM  iec.REAL // default 1.0
	OS  iec.REAL
	DL  iec.REAL
	Q   iec.BOOL
	OUT iec.REAL

	gen  generator
	temp iec.REAL
}

// INIT resets the block and sets PT and AM to their initial values.
func (g *GEN_RMP) INIT() { *g = GEN_RMP{PT: iec.TIME(time.Second), AM: 1} }

// Execute runs the block once.
func (g *GEN_RMP) Execute(now time.Time) {
	tx := g.gen.period(now, g.PT)
	g.DL = delay(g.DL)
	ltemp := g.temp
	if g.PT > 0 {
		g.temp = math.FRACT(iec.REAL(tx+ms(td.MULTIME(g.PT, g.DL))) / TIME_TO_REAL(g.PT))
	}
	g.OUT = g.AM*g.temp + g.OS
	g.Q = g.temp < ltemp
}

// GEN_SIN generates a sine wave OUT from OS - AM/2 to OS + AM/2 with the
// period PT, delayed by DL periods. Q is true while the wave is positive.
type GEN_SIN struct {
	PT  iec.TIME
	AM  iec.REAL // default 1.0
	OS  iec.REAL
	DL  iec.REAL
	Q   iec.BOOL
	OUT iec.REAL

	gen  generator
	temp iec.REAL
}

// INIT resets the block and sets AM to its initial value.
func (g *GEN_SIN) INIT() { *g = GEN_SIN{AM: 1} }

// Execute runs the block once.
func (g *GEN_SIN) Execute(now time.Time) {
	tx := g.gen.period(now, g.PT)
	g.DL = delay(g.DL)
	if g.PT > 0 {
		g.temp = SIN(MATH.PI2 * iec.REAL(tx+ms(td.MULTIME(g.PT, g.DL))) / TIME_TO_REAL(g.PT))
	}
	g.OUT = g.AM*0.5*g.temp + g.OS
	g.Q = !math.SIGN_R(g.temp)
}

// GEN_SQR generates a square wave OUT between OS - AM/2 and OS + AM/2 with
// the period PT and the duty cycle DC, delayed by DL periods. Q is true
// while OUT is high.
type GEN_SQR struct {
	PT  iec.TIME
	AM  iec.REAL // default 1.0
	OS  iec.REAL
	DC  iec.REAL // default 0.5
	DL  iec.REAL
	Q   iec.BOOL
	OUT iec.REAL

	gen generator
}

// INIT resets the block and sets AM and DC to their initial values.
func (g *GEN_SQR) INIT() { *g = GEN_SQR{AM: 1, DC: 0.5} }

// Execute runs the block once.
func (g *GEN_SQR) Execute(now time.Time) {
	high := func() { g.OUT, g.Q = g.AM*0.5+g.OS, true }
	low := func() { g.OUT, g.Q = -g.AM*0.5+g.OS, false }
	switch g.DC {
	case 0.0:
		low()
		return
	case 1.0:
		high()
		return
	}
	tx := g.gen.period(now, g.PT)
	g.DL = delay(g.DL)
	g.DC = delay(g.DC)
	at := func(f iec.REAL) iec.DWORD { return ms(td.MULTIME(g.PT, f)) }
	if at(g.DL+g.DC) >= ms(g.PT) {
		if tx >= at(g.DL+g.DC-1) {
			low()
		}
		if tx >= at(g.DL) {
			high()
		}
	} else {
		if tx >= at(g.DL) {
			high()
		}
		if tx >= at(g.DL+g.DC) {
			low()
		}
	}
}

// PWM_DC generates a pulse width modulated signal with the frequency F and
// the duty cycle DC.
type PWM_DC struct {
	F, DC iec.REAL
	Q     iec.BOOL

	clk   logic.CLK_PRG
	pulse logic.TP_X
}

// INIT resets the block.
func (p *PWM_DC) INIT() { *p = PWM_DC{} }

// Execute runs the block once.
func (p *PWM_DC) Execute(now time.Time) {
	if p.F <= 0.0 {
		return
	}
	tmp := 1000.0 / p.F
	p.clk.PT = REAL_TO_TIME(tmp)
	p.clk.Execute(now)
	p.pulse.IN, p.pulse.PT = p.clk.Q, REAL_TO_TIME(tmp*p.DC)
	p.pulse.Execute(now)
	p.Q = p.pulse.Q
}

// PWM_PW generates a pulse width modulated signal with the frequency F and
// the pulse width PW.
type PWM_PW struct {
	F  iec.REAL
	PW iec.TIME
	Q  iec.BOOL

	clk   logic.CLK_PRG
	pulse logic.TP_X
}

// INIT resets the block.
func (p *PWM_PW) INIT() { *p = PWM_PW{} }

// Execute runs the block once.
func (p *PWM_PW) Execute(now time.Time) {
	if p.F <= 0.0 {
		return
	}
	p.clk.PT = REAL_TO_TIME(1000.0 / p.F)
	p.clk.Execute(now)
	p.pulse.IN, p.pulse.PT = p.clk.Q, p.PW
	p.pulse.Execute(now)
	p.Q = p.pulse.Q
}

// RMP_B ramps OUT, a byte, up if UP is true or down, in PT for a full ramp,
// while E is true. SET sets it to 255 and RST to 0. BUSY is true while it
// ramps, HIGH at 255 and LOW at 0.
type RMP_B struct {
	SET  iec.BOOL
	PT   iec.TIME
	E    iec.BOOL // default TRUE
	UP   iec.BOOL // default TRUE
	RST  iec.BOOL
	OUT  iec.BYTE
	BUSY iec.BOOL
	HIGH iec.BOOL
	LOW  iec.BOOL

	rmp RMP_B_
}

// INIT resets the block and sets E and UP to their initial values.
func (r *RMP_B) INIT() {
	*r = RMP_B{E: true, UP: true}
	r.rmp.INIT()
}

// Execute runs the block once.
func (r *RMP_B) Execute(now time.Time) {
	r.rmp.DIR, r.rmp.E, r.rmp.TR, r.rmp.RMP = r.UP, r.E, r.PT, &r.OUT
	r.rmp.Execute(now)
	if r.RST {
		r.OUT = 0
	} else if r.SET {
		r.OUT = 255
	}
	r.LOW = r.OUT == 0
	r.HIGH = r.OUT == 255
	r.BUSY = !(r.LOW || r.HIGH) && r.E
}

// RMP_SOFT ramps OUT softly to VAL while IN is true, in PT_ON for a full
// ramp, and to 0 while IN is false, in PT_OFF.
type RMP_SOFT struct {
	IN     iec.BOOL
	VAL    iec.BYTE
	PT_ON  iec.TIME
	PT_OFF iec.TIME
	OUT    iec.BYTE

	rmp RMP_B_
}

// INIT resets the block.
func (r *RMP_SOFT) INIT() {
	*r = RMP_SOFT{}
	r.rmp.INIT()
}

// Execute runs the block once.
func (r *RMP_SOFT) Execute(now time.Time) {
	tmp := SEL(r.IN, 0, r.VAL)
	r.rmp.RMP = &r.OUT
	switch {
	case tmp > r.OUT:
		r.rmp.DIR, r.rmp.E, r.rmp.TR = true, true, r.PT_ON
		r.rmp.Execute(now)
		r.OUT = min(r.OUT, tmp)
	case tmp < r.OUT:
		r.rmp.DIR, r.rmp.E, r.rmp.TR = false, true, r.PT_OFF
		r.rmp.Execute(now)
		r.OUT = max(r.OUT, tmp)
	default:
		r.rmp.E = false
		r.rmp.Execute(now)
	}
}

// RMP_W ramps OUT, a word, up if UP is true or down, in PT for a full ramp,
// while E is true. SET sets it to 65535 and RST to 0. BUSY is true while it
// ramps, HIGH at 65535 and LOW at 0.
type RMP_W struct {
	SET  iec.BOOL
	PT   iec.TIME
	E    iec.BOOL // default TRUE
	UP   iec.BOOL // default TRUE
	RST  iec.BOOL
	OUT  iec.WORD
	BUSY iec.BOOL
	HIGH iec.BOOL
	LOW  iec.BOOL

	rmp RMP_W_
}

// INIT resets the block and sets E and UP to their initial values.
func (r *RMP_W) INIT() {
	*r = RMP_W{E: true, UP: true}
	r.rmp.INIT()
}

// Execute runs the block once.
func (r *RMP_W) Execute(now time.Time) {
	r.rmp.DIR, r.rmp.E, r.rmp.TR, r.rmp.RMP = r.UP, r.E, r.PT, &r.OUT
	r.rmp.Execute(now)
	if r.RST {
		r.OUT = 0
	} else if r.SET {
		r.OUT = 65535
	}
	r.LOW = r.OUT == 0
	r.HIGH = r.OUT == 65535
	r.BUSY = !(r.LOW || r.HIGH) && r.E
}
