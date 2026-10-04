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

package electrical

import (
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

// clock is a simulated time for the tests.
type clock struct{ now time.Time }

func newClock() *clock                          { return &clock{time.Unix(1000, 0)} }
func (c *clock) step(d time.Duration) time.Time { c.now = c.now.Add(d); return c.now }

const msec = time.Millisecond

func TestStringToTime(t *testing.T) {
	for in, want := range map[iec.STRING]iec.TIME{
		"T#1h":         iec.TIME(time.Hour),
		"t#1h30m":      iec.TIME(90 * time.Minute),
		"TIME#2s500ms": iec.TIME(2500 * msec),
		"1.5s":         iec.TIME(1500 * msec),
		"T#1d_2h":      iec.TIME(26 * time.Hour),
		"T#-5s":        iec.TIME(-5 * time.Second),
		"x":            0,
		"T#5":          0,
	} {
		if got := STRING_TO_TIME(in); got != want {
			t.Errorf("STRING_TO_TIME(%q) = %v, want %v", in, time.Duration(got), time.Duration(want))
		}
	}
	for in, want := range map[iec.STRING]iec.DWORD{
		"TOD#06:30":         6*3600000 + 30*60000,
		"tod#23:59:59.5":    23*3600000 + 59*60000 + 59500,
		"TIME_OF_DAY#1:2:3": 3600000 + 2*60000 + 3000,
		"12:00:00":          12 * 3600000,
		"25:00":             0,
		"noon":              0,
	} {
		if got := TOD_TO_DWORD(STRING_TO_TOD(in)); got != want {
			t.Errorf("STRING_TO_TOD(%q) = %d ms, want %d", in, got, want)
		}
	}
}

func TestTimerEventDecode(t *testing.T) {
	ev := TIMER_EVENT_DECODE("<2;1;Mo,Mi;TOD#06:30;T#1h;16#0F;0>", 1)
	if ev.TYP != 2 || ev.CHANNEL != 1 || ev.LAND != 15 || ev.LOR != 0 {
		t.Errorf("numbers: %+v", ev)
	}
	if TOD_TO_DWORD(ev.START) != 6*3600000+30*60000 || ev.DURATION != iec.TIME(time.Hour) {
		t.Errorf("start %d, duration %v", TOD_TO_DWORD(ev.START), time.Duration(ev.DURATION))
	}
	if ev.DAY == 0 {
		t.Errorf("the weekdays of a TYP 2 event: %08b", ev.DAY)
	}
	if ev := TIMER_EVENT_DECODE("1;2;3;4;5;6;7", 1); ev != (TIMER_EVENT{}) {
		t.Errorf("no brackets: %+v", ev)
	}
}

func TestClickMode(t *testing.T) {
	c := newClock()
	var cm CLICK_MODE
	cm.INIT()
	press := func(in bool, d time.Duration) {
		cm.IN = iec.BOOL(in)
		cm.Execute(c.step(d))
	}
	// One short click.
	press(true, 0)
	press(false, 100*msec)
	press(false, 200*msec)
	press(false, 300*msec)
	if !cm.SINGLE_ {
		t.Error("a single click")
	}
	press(false, 10*msec)
	if cm.SINGLE_ {
		t.Error("SINGLE_ is a pulse")
	}
	// Two short clicks.
	press(true, 0)
	press(false, 50*msec)
	press(true, 50*msec)
	press(false, 50*msec)
	// royaljelly's TP restarts on the second press, which IEC's TP does
	// not, so the window ends T_LONG after it.
	press(false, 500*msec)
	if !cm.DOUBLE {
		t.Error("a double click")
	}
	// A long press.
	press(true, 0)
	press(true, 600*msec)
	if !cm.LONG || !cm.TP_LONG {
		t.Errorf("a long press: LONG %v TP_LONG %v", cm.LONG, cm.TP_LONG)
	}
	press(true, 10*msec)
	if !cm.LONG || cm.TP_LONG {
		t.Errorf("TP_LONG is a pulse: LONG %v TP_LONG %v", cm.LONG, cm.TP_LONG)
	}
}

func TestPulseLength(t *testing.T) {
	c := newClock()
	var p PULSE_LENGTH
	p.INIT()
	pulse := func(d time.Duration) {
		p.IN = true
		p.Execute(c.step(0))
		p.IN = false
		p.Execute(c.step(d))
	}
	pulse(50 * msec)
	if !p.SHORT {
		t.Error("a short pulse")
	}
	pulse(500 * msec)
	if !p.MIDDLE {
		t.Error("a middle pulse")
	}
	pulse(2 * time.Second)
	if !p.LONG {
		t.Error("a long pulse")
	}
}

func TestSwitchI(t *testing.T) {
	c := newClock()
	var s SWITCH_I
	s.INIT()
	s.Execute(c.step(0))
	s.IN = true
	s.Execute(c.step(msec))
	s.Execute(c.step(20 * msec))
	if !s.Q {
		t.Fatal("a debounced rising edge turns Q on")
	}
	s.IN = false
	s.Execute(c.step(msec))
	s.Execute(c.step(20 * msec))
	if !s.Q {
		t.Fatal("the falling edge of a push button leaves Q on")
	}
	s.RST = true
	s.Execute(c.step(msec))
	if s.Q {
		t.Fatal("RST")
	}
}

func TestDimm2(t *testing.T) {
	c := newClock()
	var out iec.BYTE = 100
	d := DIMM_2{OUT: &out}
	d.INIT()
	d.SET = true
	d.Execute(c.step(0))
	if !d.Q || out != 255 {
		t.Fatalf("SET: Q %v OUT %d", d.Q, out)
	}
	d.SET = false
	d.RST = true
	d.Execute(c.step(msec))
	if d.Q || out != 0 {
		t.Fatalf("RST: Q %v OUT %d", d.Q, out)
	}
	d.RST = false
	// A click of I1 turns the dimmer on at MIN_ON at least.
	d.I1 = true
	d.Execute(c.step(msec))
	d.I1 = false
	d.Execute(c.step(100 * msec))
	d.Execute(c.step(100 * msec)) // the end of the debounce time
	d.Execute(c.step(time.Second))
	if !d.Q || out != 50 {
		t.Fatalf("a click of I1: Q %v OUT %d", d.Q, out)
	}
	// OSCAT ramps in the direction dc2.LONG, so holding I1 dims down,
	// though its comments say up.
	d.I1 = true
	for range 20 {
		d.Execute(c.step(100 * msec))
	}
	if out >= 50 {
		t.Fatalf("holding I1: Q %v OUT %d", d.Q, out)
	}
}

func TestFLamp(t *testing.T) {
	c := newClock()
	var ontime, cycles iec.UDINT
	f := F_LAMP{ONTIME: &ontime, CYCLES: &cycles}
	f.INIT()
	f.SWITCH = true
	f.DIMM = 100
	f.Execute(c.step(0))
	f.Execute(c.step(2 * time.Second))
	if f.LAMP != 255 || f.STATUS != 111 || cycles != 1 {
		t.Errorf("a new lamp burns in at full power: LAMP %d STATUS %d CYCLES %d", f.LAMP, f.STATUS, cycles)
	}
	ontime = 100 * 3600
	f.Execute(c.step(time.Second))
	if f.LAMP != 100 || f.STATUS != 112 {
		t.Errorf("after T_NO_DIMM hours it dims: LAMP %d STATUS %d", f.LAMP, f.STATUS)
	}
}

func TestTimer1(t *testing.T) {
	var x TIMER_1
	x.INIT()
	x.START = DWORD_TO_TOD(8 * 3600000)
	x.DURATION = iec.TIME(time.Hour)
	x.DTI = td.SET_DT(2024, 1, 15, 8, 30, 0) // a Monday
	x.Execute(time.Time{})
	if !x.Q {
		t.Error("within the time on a Monday")
	}
	x.DAY = 0x3F // not Monday
	x.Execute(time.Time{})
	if x.Q {
		t.Error("Monday is not set")
	}
	x.DAY = 0x7F
	x.DTI = td.SET_DT(2024, 1, 15, 9, 30, 0)
	x.Execute(time.Time{})
	if x.Q {
		t.Error("after the time")
	}
}

func TestTimer2(t *testing.T) {
	c := newClock()
	var x TIMER_2
	x.INIT()
	x.MODE = 11 // every day
	x.START = DWORD_TO_TOD(8 * 3600000)
	x.DURATION = iec.TIME(time.Minute)
	x.DT_IN = td.SET_DT(2024, 1, 15, 8, 0, 1)
	x.Execute(c.step(0))
	if !x.Q {
		t.Fatal("at START")
	}
	x.Execute(c.step(2 * time.Minute))
	if x.Q {
		t.Fatal("after DURATION")
	}
	x.Execute(c.step(time.Second))
	if x.Q {
		t.Fatal("once a day")
	}
}

func TestTimerP4(t *testing.T) {
	var prog [TIMER_P4_ARRAY_MAX + 1]TIMER_EVENT
	prog[0] = TIMER_EVENT{TYP: 1, CHANNEL: 2, START: DWORD_TO_TOD(8 * 3600000), DURATION: iec.TIME(time.Hour)}
	x := TIMER_P4{PROG: &prog}
	x.INIT()
	x.ENQ = true
	x.DTIME = td.SET_DT(2024, 1, 15, 8, 30, 0)
	x.Execute(time.Time{})
	if !x.Q2 || x.Q0 || x.Q1 || x.Q3 || x.STATUS != 102 {
		t.Errorf("a daily event on channel 2: %v %v %v %v STATUS %d", x.Q0, x.Q1, x.Q2, x.Q3, x.STATUS)
	}
	x.DTIME = td.SET_DT(2024, 1, 15, 9, 30, 1)
	x.Execute(time.Time{})
	if x.Q2 {
		t.Error("after the event")
	}
	x.MAN, x.MI = true, 0x01
	x.Execute(time.Time{})
	if !x.Q0 || x.STATUS != 101 {
		t.Error("MAN and MI set the outputs by hand")
	}
}

func TestSwitchX(t *testing.T) {
	c := newClock()
	var s SWITCH_X
	s.INIT()
	s.Execute(c.step(0))
	s.IN1 = true
	s.Execute(c.step(10 * msec))
	s.IN1 = false
	s.Execute(c.step(10 * msec))
	s.Execute(c.step(100 * msec))
	if !s.Q1 {
		t.Error("IN1 alone sets Q1 after the debounce time")
	}
	s.IN1 = true
	s.Execute(c.step(10 * msec))
	s.IN3 = true
	s.Execute(c.step(10 * msec))
	if !s.Q31 {
		t.Error("IN3 while IN1 is held sets Q31")
	}
}
