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
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// scan is a test clock: run executes a block at the current time and
// advances it.
type scan struct{ now time.Time }

func newScan() *scan                            { return &scan{time.Unix(5000, 0)} }
func (s *scan) after(d time.Duration) time.Time { s.now = s.now.Add(d); return s.now }

func TestEdgeFlipFlops(t *testing.T) {
	var now time.Time

	var d2 FF_D2E
	d2.INIT()
	d2.D0, d2.D1, d2.CLK = true, false, true
	d2.Execute(now) // rising edge: take D
	if !d2.Q0 || d2.Q1 {
		t.Fatalf("FF_D2E after the edge: %v %v", d2.Q0, d2.Q1)
	}
	d2.D0, d2.D1 = false, true
	d2.Execute(now) // CLK still high: no edge, Q holds
	if !d2.Q0 || d2.Q1 {
		t.Fatal("FF_D2E took D without an edge")
	}
	d2.CLK = false
	d2.Execute(now)
	d2.CLK = true
	d2.Execute(now)
	if d2.Q0 || !d2.Q1 {
		t.Fatal("FF_D2E second edge")
	}
	d2.RST = true
	d2.Execute(now)
	if d2.Q0 || d2.Q1 {
		t.Fatal("FF_D2E reset")
	}

	var d4 FF_D4E
	d4.INIT()
	d4.D0, d4.D1, d4.D2, d4.D3, d4.CLK = true, false, true, true, true
	d4.Execute(now)
	if !d4.Q0 || d4.Q1 || !d4.Q2 || !d4.Q3 {
		t.Fatalf("FF_D4E after the edge: %v %v %v %v", d4.Q0, d4.Q1, d4.Q2, d4.Q3)
	}
	d4.D0 = false
	d4.Execute(now)
	if !d4.Q0 {
		t.Fatal("FF_D4E took D without an edge")
	}
	d4.RST = true
	d4.Execute(now)
	if d4.Q0 || d4.Q2 || d4.Q3 {
		t.Fatal("FF_D4E reset")
	}

	var dre FF_DRE
	dre.INIT()
	dre.D, dre.CLK = true, true
	dre.Execute(now)
	if !dre.Q {
		t.Fatal("FF_DRE edge")
	}
	dre.SET, dre.D, dre.CLK = true, false, false
	dre.Execute(now)
	if !dre.Q {
		t.Fatal("FF_DRE set")
	}
	dre.RST = true // reset wins over set
	dre.Execute(now)
	if dre.Q {
		t.Fatal("FF_DRE reset over set")
	}
	dre.SET, dre.RST, dre.D, dre.CLK = false, false, true, true
	dre.Execute(now)
	if !dre.Q {
		t.Fatal("FF_DRE edge after reset")
	}
}

func TestLTCH(t *testing.T) {
	var now time.Time
	var l LTCH
	l.INIT()
	l.D = true
	l.Execute(now) // L false: closed
	if l.Q {
		t.Fatal("LTCH followed D while closed")
	}
	l.L = true
	l.Execute(now)
	if !l.Q {
		t.Fatal("LTCH open did not follow D")
	}
	l.L, l.D = false, false
	l.Execute(now)
	if !l.Q {
		t.Fatal("LTCH did not hold when closed")
	}
	l.RST = true
	l.Execute(now)
	if l.Q {
		t.Fatal("LTCH reset")
	}
}

func TestSHR_4UDE(t *testing.T) {
	var now time.Time
	var s SHR_4UDE
	s.INIT()
	q := func() [4]iec.BOOL { return [4]iec.BOOL{s.Q0, s.Q1, s.Q2, s.Q3} }
	clock := func() { s.CLK = true; s.Execute(now); s.CLK = false; s.Execute(now) }

	s.D0 = true
	clock() // up: D0 into Q0
	if q() != [4]iec.BOOL{true, false, false, false} {
		t.Fatalf("after one shift up: %v", q())
	}
	s.D0 = false
	clock()
	if q() != [4]iec.BOOL{false, true, false, false} {
		t.Fatalf("after two shifts up: %v", q())
	}
	s.DN, s.D3 = true, true // down: D3 into Q3
	clock()
	if q() != [4]iec.BOOL{true, false, false, true} {
		t.Fatalf("after a shift down: %v", q())
	}
	s.SET = true
	s.Execute(now)
	if q() != [4]iec.BOOL{true, true, true, true} {
		t.Fatalf("after SET: %v", q())
	}
	s.SET, s.RST = false, true
	s.Execute(now)
	if q() != [4]iec.BOOL{} {
		t.Fatalf("after RST: %v", q())
	}
}

func TestFIFO_32(t *testing.T) {
	var now time.Time
	var f FIFO_32
	f.INIT()
	if !f.EMPTY || f.FULL {
		t.Fatal("FIFO_32 starts empty")
	}
	write := func(v iec.DWORD) { f.DIN, f.WD, f.RD = v, true, false; f.Execute(now); f.WD = false }
	read := func() iec.DWORD { f.RD, f.WD = true, false; f.Execute(now); f.RD = false; return f.DOUT }

	for v := iec.DWORD(1); v <= 3; v++ {
		write(v)
	}
	if got := []iec.DWORD{read(), read()}; got[0] != 1 || got[1] != 2 {
		t.Fatalf("FIFO_32 order: %v, want 1 2", got)
	}
	read()
	if !f.EMPTY {
		t.Fatal("FIFO_32 not empty after reading everything")
	}
	f.DOUT = 99
	read() // reading an empty FIFO changes nothing
	if f.DOUT != 99 {
		t.Fatal("FIFO_32 read while empty")
	}

	for v := range iec.DWORD(32) {
		write(v + 100)
	}
	if !f.FULL {
		t.Fatal("FIFO_32 not full after 32 writes")
	}
	write(999) // ignored while full
	if read() != 100 {
		t.Fatal("FIFO_32 oldest value after filling")
	}

	f.E = false // disabled: nothing happens
	f.WD, f.DIN = true, 7
	f.Execute(now)
	f.E, f.WD = true, false

	f.RST = true
	f.Execute(now)
	if !f.EMPTY || f.FULL || f.DOUT != 0 {
		t.Fatal("FIFO_32 reset")
	}

	var g FIFO_32 // without INIT it starts empty too, but disabled (E false)
	g.WD, g.DIN = true, 5
	g.Execute(now)
	if !g.EMPTY {
		t.Fatal("FIFO_32 without INIT wrote while E was false")
	}
}

func TestSTACK_16(t *testing.T) {
	var now time.Time
	var s STACK_16
	s.INIT()
	push := func(v iec.DWORD) { s.DIN, s.WD, s.RD = v, true, false; s.Execute(now); s.WD = false }
	pop := func() iec.DWORD { s.RD, s.WD = true, false; s.Execute(now); s.RD = false; return s.DOUT }

	push(1)
	push(2)
	push(3)
	if got := []iec.DWORD{pop(), pop(), pop()}; got[0] != 3 || got[1] != 2 || got[2] != 1 {
		t.Fatalf("STACK_16 order: %v, want 3 2 1", got)
	}
	if !s.EMPTY {
		t.Fatal("STACK_16 not empty")
	}
	for v := range iec.DWORD(16) {
		push(v)
	}
	if !s.FULL {
		t.Fatal("STACK_16 not full after 16 pushes")
	}
	push(99) // ignored while full
	if pop() != 15 {
		t.Fatal("STACK_16 newest value after filling")
	}
	s.RST = true
	s.Execute(now)
	if !s.EMPTY || s.FULL || s.DOUT != 0 {
		t.Fatal("STACK_16 reset")
	}

	var u STACK_16 // without INIT: empty and disabled
	u.WD, u.DIN = true, 1
	u.Execute(now)
	if !u.EMPTY {
		t.Fatal("STACK_16 without INIT pushed while E was false")
	}
}

func TestTP_1(t *testing.T) {
	s := newScan()
	var p TP_1
	p.INIT()
	p.PT = iec.TIME(100 * time.Millisecond)
	p.IN = true
	p.Execute(s.now)
	if !p.Q {
		t.Fatal("TP_1 rising edge")
	}
	p.IN = false
	p.Execute(s.after(50 * time.Millisecond))
	if !p.Q {
		t.Fatal("TP_1 ended before PT")
	}
	p.IN = true // a new edge restarts the pulse
	p.Execute(s.after(30 * time.Millisecond))
	p.Execute(s.after(90 * time.Millisecond)) // 170 ms after the first edge
	if !p.Q {
		t.Fatal("TP_1 did not restart on a new edge")
	}
	p.Execute(s.after(20 * time.Millisecond))
	if p.Q {
		t.Fatal("TP_1 still on after PT")
	}
	p.IN = false
	p.Execute(s.now)
	p.IN = true
	p.Execute(s.now)
	p.RST = true
	p.Execute(s.now)
	if p.Q {
		t.Fatal("TP_1 reset")
	}
}

func TestTMIN(t *testing.T) {
	s := newScan()
	var m TMIN
	m.INIT()
	m.PT = iec.TIME(100 * time.Millisecond)
	m.IN = true
	m.Execute(s.now)
	m.IN = false // a short pulse is stretched to PT
	m.Execute(s.after(10 * time.Millisecond))
	if !m.Q {
		t.Fatal("TMIN did not hold a short pulse")
	}
	m.Execute(s.after(100 * time.Millisecond))
	if m.Q {
		t.Fatal("TMIN held past PT")
	}
	m.IN = true // a long input is followed
	m.Execute(s.now)
	m.Execute(s.after(500 * time.Millisecond))
	if !m.Q {
		t.Fatal("TMIN dropped a long input")
	}
	m.IN = false
	m.Execute(s.after(time.Millisecond))
	if m.Q {
		t.Fatal("TMIN held after a long input ended")
	}
}

func TestCLK_N(t *testing.T) {
	s := newScan()
	var c CLK_N
	c.INIT()
	c.N = 3 // bit 3 of the millisecond timer changes every 8 ms
	pulses := 0
	for range 1000 {
		c.Execute(s.after(time.Millisecond))
		if c.Q {
			pulses++
		}
	}
	if pulses < 124 || pulses > 126 {
		t.Fatalf("CLK_N(3) gave %d pulses in 1000 ms, want 125", pulses)
	}
}

func TestCLICK_CNT(t *testing.T) {
	s := newScan()
	clicks := func(c *CLICK_CNT, n int) (fired int) {
		for range n {
			c.IN = true
			c.Execute(s.after(50 * time.Millisecond))
			c.IN = false
			c.Execute(s.after(50 * time.Millisecond))
		}
		for range 20 { // past TC
			c.Execute(s.after(100 * time.Millisecond))
			if c.Q {
				fired++
			}
		}
		return fired
	}
	var c CLICK_CNT
	c.INIT()
	c.N, c.TC = 2, iec.TIME(time.Second)
	if got := clicks(&c, 2); got != 1 {
		t.Fatalf("two clicks: Q set %d times, want once", got)
	}
	if got := clicks(&c, 3); got != 0 {
		t.Fatalf("three clicks for N 2: Q set %d times, want never", got)
	}
	if got := clicks(&c, 1); got != 0 {
		t.Fatalf("one click for N 2: Q set %d times, want never", got)
	}

	var d CLICK_CNT // without INIT: the count starts at -1 as well
	d.N, d.TC = 1, iec.TIME(time.Second)
	if got := clicks(&d, 1); got != 1 {
		t.Fatalf("CLICK_CNT without INIT: Q set %d times, want once", got)
	}
}
