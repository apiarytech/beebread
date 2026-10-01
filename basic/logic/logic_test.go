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
	gomath "math"
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

func TestGates(t *testing.T) {
	tests := []struct {
		name      string
		got, want uint64
	}{
		{"BCDC_TO_INT(16#42)", uint64(BCDC_TO_INT(0x42)), 42},
		{"INT_TO_BCDC(42)", uint64(INT_TO_BCDC(42)), 0x42},
		{"BIT_COUNT(3)", uint64(BIT_COUNT(3)), 2},
		{"BIT_COUNT(FFFFFFFF)", uint64(BIT_COUNT(0xFFFFFFFF)), 32},
		{"BIT_LOAD_B set", uint64(BIT_LOAD_B(0, true, 3)), 8},
		{"BIT_LOAD_B clear", uint64(BIT_LOAD_B(0xFF, false, 0)), 0xFE},
		{"BIT_LOAD_W set", uint64(BIT_LOAD_W(0, true, 15)), 0x8000},
		{"BIT_LOAD_DW clear", uint64(BIT_LOAD_DW(0xFFFFFFFF, false, 31)), 0x7FFFFFFF},
		{"BIT_LOAD_B2 set", uint64(BIT_LOAD_B2(0, true, 2, 3)), 0x1C},
		{"BIT_LOAD_B2 clear", uint64(BIT_LOAD_B2(0xFF, false, 2, 3)), 0xE3},
		{"BIT_LOAD_W2 set", uint64(BIT_LOAD_W2(0, true, 4, 8)), 0x0FF0},
		{"BIT_LOAD_DW2 clear", uint64(BIT_LOAD_DW2(0xFFFFFFFF, false, 28, 4)), 0x0FFFFFFF},
		{"BIT_TOGGLE_B", uint64(BIT_TOGGLE_B(1, 0)), 0},
		{"BIT_TOGGLE_W", uint64(BIT_TOGGLE_W(0, 15)), 0x8000},
		{"BIT_TOGGLE_DW", uint64(BIT_TOGGLE_DW(0, 31)), 0x80000000},
		{"BYTE_OF_BIT", uint64(BYTE_OF_BIT(true, false, false, false, false, false, false, true)), 0x81},
		{"BYTE_OF_DWORD", uint64(BYTE_OF_DWORD(0x12345678, 2)), 0x34},
		{"BYTE_TO_GRAY(3)", uint64(BYTE_TO_GRAY(3)), 2},
		{"GRAY_TO_BYTE(2)", uint64(GRAY_TO_BYTE(2)), 3},
		{"GRAY round trip", uint64(GRAY_TO_BYTE(BYTE_TO_GRAY(0xA7))), 0xA7},
		{"CHK_REAL(1)", uint64(CHK_REAL(1)), 0},
		{"CHK_REAL(+inf)", uint64(CHK_REAL(iec.REAL(gomath.Inf(1)))), 0x20},
		{"CHK_REAL(-inf)", uint64(CHK_REAL(iec.REAL(gomath.Inf(-1)))), 0x40},
		{"CHK_REAL(NaN)", uint64(CHK_REAL(iec.REAL(gomath.NaN()))), 0x80},
		{"DWORD_OF_BYTE", uint64(DWORD_OF_BYTE(0x12, 0x34, 0x56, 0x78)), 0x12345678},
		{"DWORD_OF_WORD", uint64(DWORD_OF_WORD(0x1234, 0x5678)), 0x12345678},
		{"WORD_OF_BYTE", uint64(WORD_OF_BYTE(0x12, 0x34)), 0x1234},
		{"WORD_OF_DWORD", uint64(WORD_OF_DWORD(0x12345678, 1)), 0x1234},
		{"REAL_TO_DW(1)", uint64(REAL_TO_DW(1)), 0x3F800000},
		{"REFLECT", uint64(REFLECT(0x1, 4)), 0x8},
		{"REFLECT keeps high bits", uint64(REFLECT(0x31, 4)), 0x38},
		{"REVERSE", uint64(REVERSE(0x01)), 0x80},
		{"REVERSE(16#C4)", uint64(REVERSE(0xC4)), 0x23},
		{"SHL1", uint64(SHL1(0, 4)), 0xF},
		{"SHR1", uint64(SHR1(0, 4)), 0xF0000000},
		{"SHL1 by 0", uint64(SHL1(5, 0)), 5},
		{"SWAP_BYTE", uint64(SWAP_BYTE(0x1234)), 0x3412},
		{"SWAP_BYTE2", uint64(SWAP_BYTE2(0x12345678)), 0x78563412},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %#x, want %#x", tt.name, tt.got, tt.want)
		}
	}
	if DW_TO_REAL(0x3F800000) != 1 {
		t.Error("DW_TO_REAL")
	}
	bools := []struct {
		name      string
		got, want iec.BOOL
	}{
		{"BIT_OF_DWORD", BIT_OF_DWORD(8, 3), true},
		{"PARITY(7)", PARITY(7), true},
		{"PARITY(3)", PARITY(3), false},
		{"CHECK_PARITY(3,false)", CHECK_PARITY(3, false), true},
		{"CHECK_PARITY(7,true)", CHECK_PARITY(7, true), true},
		{"CHECK_PARITY(7,false)", CHECK_PARITY(7, false), false},
		{"MUX_2", MUX_2(false, true, true), true},
		{"MUX_4", MUX_4(false, false, true, false, false, true), true},
	}
	for _, tt := range bools {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestDecoders(t *testing.T) {
	var d8 DEC_8
	d8.D, d8.A0, d8.A2 = true, true, true
	d8.Execute(time.Time{})
	if !d8.Q5 || d8.Q0 || d8.Q4 {
		t.Errorf("DEC_8 = %+v, want Q5", d8)
	}
	var d4 DEC_4
	d4.D, d4.A1 = true, true
	d4.Execute(time.Time{})
	if !d4.Q2 || d4.Q0 {
		t.Errorf("DEC_4 = %+v, want Q2", d4)
	}
	var d2 DEC_2
	d2.D = true
	d2.Execute(time.Time{})
	if !d2.Q0 || d2.Q1 {
		t.Errorf("DEC_2 = %+v, want Q0", d2)
	}
	var b BYTE_TO_BITS
	b.IN = 0x81
	b.Execute(time.Time{})
	if !b.B0 || !b.B7 || b.B1 {
		t.Errorf("BYTE_TO_BITS = %+v", b)
	}
}

func TestFlipFlops(t *testing.T) {
	var now time.Time

	var jk FF_JKE
	jk.J = true
	jk.CLK = true
	jk.Execute(now)
	if !jk.Q {
		t.Fatal("FF_JKE: J did not set Q")
	}
	jk.CLK = false
	jk.Execute(now)
	jk.K, jk.CLK = true, true
	jk.Execute(now)
	if jk.Q {
		t.Fatal("FF_JKE: J and K did not toggle Q")
	}
	jk.CLK = false
	jk.Execute(now)
	jk.CLK = true
	jk.Execute(now)
	if !jk.Q {
		t.Fatal("FF_JKE: second toggle")
	}
	jk.RST = true
	jk.Execute(now)
	if jk.Q {
		t.Fatal("FF_JKE: RST")
	}

	var tg TOGGLE
	for i := 0; i < 3; i++ {
		tg.CLK = true
		tg.Execute(now)
		tg.CLK = false
		tg.Execute(now)
	}
	if !tg.Q {
		t.Fatal("TOGGLE after 3 edges should be true")
	}

	var rs FF_RSE
	rs.CS, rs.CR = true, true
	rs.Execute(now)
	if rs.Q {
		t.Fatal("FF_RSE: CR has priority")
	}

	var cnt COUNT_BR
	cnt.INIT()
	cnt.MX = 3
	for i := 0; i < 5; i++ {
		cnt.UP = true
		cnt.Execute(now)
		cnt.UP = false
		cnt.Execute(now)
	}
	if cnt.CNT != 1 {
		t.Fatalf("COUNT_BR = %v, want 1", cnt.CNT)
	}
	cnt.DN = true
	cnt.Execute(now)
	cnt.DN = false
	cnt.Execute(now)
	cnt.DN = true
	cnt.Execute(now)
	if cnt.CNT != 3 {
		t.Fatalf("COUNT_BR down = %v, want 3", cnt.CNT)
	}

	var dr COUNT_DR
	dr.INIT()
	dr.MX = 10
	dr.STEP = 4
	for i := 0; i < 3; i++ {
		dr.UP = true
		dr.Execute(now)
		dr.UP = false
		dr.Execute(now)
	}
	if dr.CNT != 1 {
		t.Fatalf("COUNT_DR = %v, want 1", dr.CNT)
	}

	var sel SELECT_8
	sel.E = true
	sel.DN = true
	sel.Execute(now)
	if !sel.Q7 || sel.STATE != 7 {
		t.Fatalf("SELECT_8 down from 0 = %v", sel.STATE)
	}

	var sh SHR_4E
	sh.D0 = true
	sh.CLK = true
	sh.Execute(now)
	sh.CLK, sh.D0 = false, false
	sh.Execute(now)
	sh.CLK = true
	sh.Execute(now)
	if sh.Q0 || !sh.Q1 || sh.Q2 {
		t.Fatalf("SHR_4E = %+v", sh)
	}

	var ud SHR_8UDE
	ud.SET = true
	ud.Execute(now)
	ud.SET = false
	ud.DN = true
	ud.CLK = true
	ud.Execute(now)
	if ud.Q7 || !ud.Q6 {
		t.Fatalf("SHR_8UDE down = %+v", ud)
	}

	var pl SHR_8PLE
	pl.INIT()
	pl.DIN = true
	pl.CLK = true
	pl.Execute(now)
	for i := 0; i < 7; i++ {
		pl.CLK = false
		pl.DIN = false
		pl.Execute(now)
		pl.CLK = true
		pl.Execute(now)
	}
	if !pl.DOUT {
		t.Fatal("SHR_8PLE: the bit shifted in did not reach DOUT after 8 clocks")
	}

	var st STORE_8
	st.D2, st.D5 = true, true
	st.Execute(now)
	st.D2, st.D5 = false, false
	st.CLR = true
	st.Execute(now)
	if st.Q2 || !st.Q5 {
		t.Fatalf("STORE_8 CLR = %+v", st)
	}

	var l LTCH_4
	l.D1, l.L = true, true
	l.Execute(now)
	l.L, l.D1 = false, false
	l.Execute(now)
	if !l.Q1 {
		t.Fatal("LTCH_4 did not hold")
	}
}

func TestMemory(t *testing.T) {
	var f FIFO_16
	f.INIT()
	for i := 1; i <= 17; i++ {
		f.DIN = iec.DWORD(i)
		f.WD = true
		f.Execute(time.Time{})
	}
	if !f.FULL {
		t.Fatal("FIFO_16 not full after 16 writes")
	}
	f.WD, f.RD = false, true
	f.Execute(time.Time{})
	if f.DOUT != 1 {
		t.Fatalf("FIFO_16 first out = %v, want 1", f.DOUT)
	}

	var s STACK_32
	s.INIT()
	for i := 1; i <= 3; i++ {
		s.DIN = iec.DWORD(i)
		s.WD = true
		s.Execute(time.Time{})
	}
	s.WD, s.RD = false, true
	s.Execute(time.Time{})
	if s.DOUT != 3 || s.EMPTY {
		t.Fatalf("STACK_32 = %v, want 3", s.DOUT)
	}
	s.RST = true
	s.Execute(time.Time{})
	if !s.EMPTY {
		t.Fatal("STACK_32 RST")
	}
}

func TestGenerators(t *testing.T) {
	now := time.Unix(1000, 0)
	step := func(d time.Duration) { now = now.Add(d) }

	var gsq GEN_SQ
	gsq.PT = iec.TIME(100 * time.Millisecond)
	gsq.Execute(now)
	if !gsq.Q {
		t.Fatal("GEN_SQ starts high")
	}
	step(50 * time.Millisecond)
	gsq.Execute(now)
	if gsq.Q {
		t.Fatal("GEN_SQ after half a period")
	}

	var cp CLK_PRG
	cp.PT = iec.TIME(time.Second)
	cp.Execute(now)
	if !cp.Q {
		t.Fatal("CLK_PRG first pulse")
	}
	step(500 * time.Millisecond)
	cp.Execute(now)
	if cp.Q {
		t.Fatal("CLK_PRG pulse too early")
	}
	step(500 * time.Millisecond)
	cp.Execute(now)
	if !cp.Q {
		t.Fatal("CLK_PRG second pulse")
	}

	var pulse CLK_PULSE
	pulse.PT = iec.TIME(10 * time.Millisecond)
	pulse.N = 2
	pulses := 0
	for i := 0; i < 10; i++ {
		pulse.Execute(now)
		if pulse.Q {
			pulses++
		}
		step(10 * time.Millisecond)
	}
	if pulses != 2 || pulse.CNT != 2 {
		t.Fatalf("CLK_PULSE pulses = %v, want 2", pulses)
	}

	var c4 CYCLE_4
	c4.INIT()
	c4.T0, c4.T1, c4.T2, c4.T3 = iec.TIME(time.Second), iec.TIME(time.Second), iec.TIME(time.Second), iec.TIME(time.Second)
	c4.S0 = true
	states := []iec.INT{}
	for i := 0; i < 5; i++ {
		c4.Execute(now)
		states = append(states, c4.STATE)
		step(time.Second)
	}
	if states[0] != 0 || states[1] != 1 || states[4] != 0 {
		t.Fatalf("CYCLE_4 states = %v", states)
	}

	var tx TP_X
	tx.PT = iec.TIME(time.Second)
	tx.IN = true
	tx.Execute(now)
	step(600 * time.Millisecond)
	tx.IN = false
	tx.Execute(now)
	if !tx.Q || tx.ET != iec.TIME(600*time.Millisecond) {
		t.Fatalf("TP_X = %+v", tx)
	}
	step(600 * time.Millisecond)
	tx.Execute(now)
	if tx.Q {
		t.Fatal("TP_X did not end")
	}

	var tm TMAX
	tm.PT = iec.TIME(time.Second)
	tm.IN = true
	tm.Execute(now)
	step(2 * time.Second)
	tm.Execute(now)
	if tm.Q || !tm.Z {
		t.Fatalf("TMAX = %+v", tm)
	}

	var tof TOF_1
	tof.PT = iec.TIME(time.Second)
	tof.IN = true
	tof.Execute(now)
	tof.IN = false
	step(500 * time.Millisecond)
	tof.Execute(now)
	if !tof.Q {
		t.Fatal("TOF_1 released early")
	}
	step(600 * time.Millisecond)
	tof.Execute(now)
	if tof.Q {
		t.Fatal("TOF_1 did not release")
	}

	var tonof TONOF
	tonof.T_ON, tonof.T_OFF = iec.TIME(time.Second), iec.TIME(2*time.Second)
	tonof.IN = true
	tonof.Execute(now)
	if tonof.Q {
		t.Fatal("TONOF switched on without delay")
	}
	step(time.Second)
	tonof.Execute(now)
	if !tonof.Q {
		t.Fatal("TONOF did not switch on after T_ON")
	}
	tonof.IN = false
	tonof.Execute(now)
	step(time.Second)
	tonof.Execute(now)
	if !tonof.Q {
		t.Fatal("TONOF switched off before T_OFF")
	}
	step(time.Second)
	tonof.Execute(now)
	if tonof.Q {
		t.Fatal("TONOF did not switch off after T_OFF")
	}

	var d TP_1D
	d.PT1, d.PTD = iec.TIME(time.Second), iec.TIME(time.Second)
	d.IN = true
	d.Execute(now)
	if !d.Q || d.IN {
		t.Fatalf("TP_1D start = %+v", d)
	}
	step(time.Second)
	d.Execute(now)
	if d.Q || !d.W {
		t.Fatalf("TP_1D wait = %+v", d)
	}

	var seq SEQUENCE_4
	seq.INIT()
	seq.WAIT0, seq.DELAY0 = iec.TIME(time.Second), iec.TIME(time.Second)
	seq.WAIT1, seq.DELAY1 = iec.TIME(time.Second), iec.TIME(time.Second)
	seq.WAIT2, seq.DELAY2 = iec.TIME(time.Second), iec.TIME(time.Second)
	seq.WAIT3, seq.DELAY3 = iec.TIME(time.Second), iec.TIME(time.Second)
	seq.START = true
	seq.Execute(now)
	if !seq.Q0 || !seq.RUN || seq.STATUS != 111 {
		t.Fatalf("SEQUENCE_4 start = %+v", seq)
	}
	for i := 0; i < 3; i++ {
		step(time.Second)
		seq.Execute(now)
	}
	if !seq.Q3 || seq.Q2 || seq.STEP != 3 {
		t.Fatalf("SEQUENCE_4 step 3 = %+v", seq)
	}
	step(time.Second)
	seq.Execute(now)
	if seq.QX || seq.RUN || seq.STATUS != 110 || seq.STEP != -1 {
		t.Fatalf("SEQUENCE_4 end = %+v", seq)
	}

	var s8 SEQUENCE_8
	s8.INIT()
	s8.IN1 = false
	s8.WAIT0, s8.WAIT1 = iec.TIME(time.Second), iec.TIME(time.Second)
	s8.START = true
	s8.Execute(now)
	step(2 * time.Second)
	s8.Execute(now)
	step(2 * time.Second)
	s8.Execute(now)
	if s8.STATUS != 2 || s8.RUN {
		t.Fatalf("SEQUENCE_8 error = %+v", s8)
	}

	var s64 SEQUENCE_64
	s64.SMAX = 1
	s64.PROG[0], s64.PROG[1] = iec.TIME(time.Second), iec.TIME(time.Second)
	s64.START = true
	s64.Execute(now)
	step(time.Second)
	s64.Execute(now)
	if s64.STATE != 1 || !s64.TRIG {
		t.Fatalf("SEQUENCE_64 = %+v", s64)
	}
	step(time.Second)
	s64.Execute(now)
	if s64.STATE != -1 {
		t.Fatalf("SEQUENCE_64 end = %+v", s64)
	}

	var sch SCHEDULER
	sch.E0, sch.E1 = true, true
	sch.T0, sch.T1 = iec.TIME(time.Second), iec.TIME(time.Second)
	sch.Execute(now)
	if !sch.Q0 {
		t.Fatal("SCHEDULER Q0 first scan")
	}
	sch.Execute(now)
	if !sch.Q1 || sch.Q0 {
		t.Fatal("SCHEDULER Q1 second scan")
	}

	var s2 SCHEDULER_2
	s2.E0, s2.C0 = true, 3
	hits := 0
	for i := 0; i < 9; i++ {
		s2.Execute(now)
		if s2.Q0 {
			hits++
		}
	}
	if hits != 3 {
		t.Fatalf("SCHEDULER_2 hits = %v, want 3", hits)
	}
}

func TestClicksAndTriggers(t *testing.T) {
	now := time.Unix(1000, 0)
	scan := func(c *CLICK_DEC, in bool) {
		c.IN = iec.BOOL(in)
		c.Execute(now)
		now = now.Add(50 * time.Millisecond)
	}
	var c CLICK_DEC
	c.TC = iec.TIME(500 * time.Millisecond)
	scan(&c, true)
	scan(&c, false)
	scan(&c, true)
	scan(&c, false)
	scan(&c, true)
	for i := 0; i < 12; i++ {
		scan(&c, true)
	}
	if !c.Q2 {
		t.Fatalf("CLICK_DEC = %+v, want Q2", c)
	}

	var a A_TRIG
	a.RES = 1
	a.IN = 0.5
	a.Execute(now)
	if a.Q || a.D != 0.5 {
		t.Fatalf("A_TRIG small change = %+v", a)
	}
	a.IN = 2
	a.Execute(now)
	if !a.Q || a.D != 0 {
		t.Fatalf("A_TRIG = %+v", a)
	}

	var b B_TRIG
	b.CLK = true
	b.Execute(now)
	first := b.Q
	b.Execute(now)
	if !first || b.Q {
		t.Fatal("B_TRIG")
	}

	var d D_TRIG
	d.IN = 5
	d.Execute(now)
	if !d.Q || d.X != 5 {
		t.Fatal("D_TRIG")
	}

	var g GEN_BIT
	g.IN0 = 0x5
	g.STEPS = 3
	g.REP = 1
	var out []iec.BOOL
	for i := 0; i < 4; i++ {
		g.CLK = true
		g.Execute(now)
		out = append(out, g.Q0)
	}
	if !out[0] || out[1] || !out[2] || g.RUN {
		t.Fatalf("GEN_BIT = %v run %v", out, g.RUN)
	}

	var div CLK_DIV
	for i := 0; i < 3; i++ {
		div.CLK = true
		div.Execute(now)
	}
	if !div.Q0 || !div.Q1 {
		t.Fatal("CLK_DIV count 3")
	}
}

func TestOther(t *testing.T) {
	// CRC-16/XMODEM of "123456789" is 16#31C3.
	msg := []iec.BYTE("123456789")
	if got := CRC_GEN(msg, 9, 16, 0x1021, 0, false, false, 0); got != 0x31C3 {
		t.Errorf("CRC_GEN XMODEM = %#x, want 0x31c3", got)
	}
	// CRC-32 of "123456789" is 16#CBF43926.
	if got := CRC_GEN(msg, 9, 32, 0x04C11DB7, 0xFFFFFFFF, true, true, 0xFFFFFFFF); got != 0xCBF43926 {
		t.Errorf("CRC_GEN CRC-32 = %#x, want 0xcbf43926", got)
	}

	var p PIN_CODE
	p.PIN = "12"
	p.E = true
	p.CB = '1'
	p.Execute(time.Time{})
	p.CB = '2'
	p.Execute(time.Time{})
	if !p.TP {
		t.Fatal("PIN_CODE")
	}

	var m MATRIX
	m.INIT()
	// Row 0 is driven on the first scan; key in column 2 pressed.
	m.X2 = true
	m.Execute(time.Time{})
	if !m.TP || m.CODE != 0x82 {
		t.Fatalf("MATRIX press = %#x", m.CODE)
	}
	if m.Y1 || !m.Y2 {
		t.Fatal("MATRIX did not move to row 1")
	}
}
