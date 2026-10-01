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
	gomath "math"
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

func near(a, b, eps iec.REAL) iec.BOOL { return gomath.Abs(float64(a-b)) <= float64(eps) }

func secs(s float64) iec.TIME { return iec.TIME(time.Duration(s * float64(time.Second))) }

// clock is a scan clock for the tests.
type clock struct{ now time.Time }

func newClock() *clock                { return &clock{time.Unix(10000, 0)} }
func (c *clock) step(d time.Duration) { c.now = c.now.Add(d) }
func (c *clock) stepS(s float64) time.Time {
	c.step(time.Duration(s * float64(time.Second)))
	return c.now
}

func TestDrivers(t *testing.T) {
	c := newClock()
	var d DRIVER_1
	d.TOGGLE_MODE = true
	d.TIMEOUT = secs(2)
	d.IN = true
	d.Execute(c.now)
	if !d.Q {
		t.Fatal("DRIVER_1 toggle on")
	}
	d.IN = false
	d.Execute(c.stepS(1))
	d.IN = true
	d.Execute(c.now)
	if d.Q {
		t.Fatal("DRIVER_1 toggle off")
	}
	d.IN = false
	d.Execute(c.now)
	d.IN = true
	d.Execute(c.now)
	d.Execute(c.stepS(2.5))
	d.Execute(c.now)
	if d.Q {
		t.Fatal("DRIVER_1 timeout")
	}

	var d4 DRIVER_4
	d4.IN2 = true
	d4.Execute(c.now)
	if !d4.Q2 || d4.Q0 {
		t.Fatalf("DRIVER_4 = %+v", d4)
	}
	d4.RST = true
	d4.Execute(c.now)
	if d4.Q2 {
		t.Fatal("DRIVER_4 RST")
	}

	var dc DRIVER_4C
	dc.INIT()
	want := []struct {
		sn     iec.INT
		q0, q1 iec.BOOL
		q3     iec.BOOL
	}{{1, true, false, false}, {2, true, true, false}, {3, true, true, false}, {4, true, true, true}, {0, false, false, false}}
	for i, w := range want {
		dc.IN = true
		dc.Execute(c.now)
		dc.IN = false
		dc.Execute(c.now)
		if dc.SN != w.sn || dc.Q0 != w.q0 || dc.Q1 != w.q1 || dc.Q3 != w.q3 {
			t.Fatalf("DRIVER_4C edge %d: SN %v Q %v %v %v %v", i+1, dc.SN, dc.Q0, dc.Q1, dc.Q2, dc.Q3)
		}
	}
}

func TestFlowControl(t *testing.T) {
	c := newClock()
	var f FLOW_CONTROL
	f.INIT()
	f.T_AUTO, f.T_DELAY = secs(10), secs(20)
	f.ENQ, f.REQ = true, true
	f.Execute(c.now)
	if !f.Q || f.STATUS != 102 {
		t.Fatalf("FLOW_CONTROL request = %v %v", f.Q, f.STATUS)
	}
	f.REQ = false
	f.Execute(c.stepS(5))
	if !f.Q {
		t.Fatal("FLOW_CONTROL stopped early")
	}
	f.Execute(c.stepS(6))
	if f.Q {
		t.Fatal("FLOW_CONTROL did not stop after T_AUTO")
	}
	f.IN = true
	f.Execute(c.now)
	if !f.Q || f.STATUS != 101 {
		t.Fatal("FLOW_CONTROL IN")
	}
}

func TestProfile(t *testing.T) {
	c := newClock()
	var p FT_PROFILE
	p.INIT()
	p.VALUE_0, p.TIME_1, p.VALUE_1 = 0, secs(10), 100
	p.TIME_2, p.VALUE_2 = secs(20), 100
	p.TIME_3, p.VALUE_3 = secs(30), 50
	p.TIME_10, p.VALUE_10 = secs(40), 50
	p.TIME_11, p.VALUE_11 = secs(50), 0
	p.TIME_12, p.TIME_13 = secs(60), secs(70)
	p.E = true
	p.Execute(c.now)
	p.E = false
	p.Execute(c.stepS(5))
	if !near(p.Y, 50, 0.1) || !p.RUN {
		t.Fatalf("FT_PROFILE at 5 s = %v", p.Y)
	}
	for i := 0; i < 20; i++ {
		p.Execute(c.stepS(1))
	}
	if !near(p.Y, 75, 0.1) {
		t.Fatalf("FT_PROFILE at 25 s = %v", p.Y)
	}
	for i := 0; i < 60; i++ {
		p.Execute(c.stepS(1))
	}
	if p.RUN || p.Y != 0 {
		t.Fatalf("FT_PROFILE end = %v %v", p.RUN, p.Y)
	}
}

func TestIncDec(t *testing.T) {
	var d INC_DEC
	seq := [][2]iec.BOOL{{true, false}, {true, true}, {false, true}, {false, false}}
	for _, s := range seq {
		d.CHA, d.CHB = s[0], s[1]
		d.Execute(time.Time{})
	}
	if d.CNT != 4 || !d.DIR {
		t.Fatalf("INC_DEC up = %v %v", d.CNT, d.DIR)
	}
	for i := len(seq) - 2; i >= 0; i-- {
		d.CHA, d.CHB = seq[i][0], seq[i][1]
		d.Execute(time.Time{})
	}
	if d.CNT != 1 || d.DIR {
		t.Fatalf("INC_DEC down = %v %v", d.CNT, d.DIR)
	}
}

func TestInterlocks(t *testing.T) {
	c := newClock()
	var il INTERLOCK
	il.TL = secs(1)
	il.I1 = true
	il.Execute(c.now)
	il.I1, il.I2 = false, true
	il.Execute(c.stepS(0.5))
	if il.Q2 {
		t.Fatal("INTERLOCK Q2 during dead time")
	}
	il.Execute(c.stepS(1))
	if !il.Q2 {
		t.Fatal("INTERLOCK Q2 after dead time")
	}

	var i4 INTERLOCK_4
	i4.E, i4.MODE = true, 3
	i4.I1 = true
	i4.Execute(c.now)
	if i4.OUT != 2 || !i4.TP {
		t.Fatalf("INTERLOCK_4 = %v", i4.OUT)
	}
	i4.I3 = true
	i4.Execute(c.now)
	if i4.OUT != 2 {
		t.Fatalf("INTERLOCK_4 mode 3 lock = %v", i4.OUT)
	}
	i4.MODE = 2
	i4.Execute(c.now)
	i4.I1, i4.I3 = true, false
	i4.Execute(c.now)
	i4.I0 = true
	i4.Execute(c.now)
	if i4.OUT != 1 {
		t.Fatalf("INTERLOCK_4 mode 2 last = %v", i4.OUT)
	}
}

func TestManual(t *testing.T) {
	if MANUAL(false, true, false) != true || MANUAL(true, true, true) != false {
		t.Error("MANUAL")
	}
	var m1 MANUAL_1
	m1.IN = true
	m1.Execute(time.Time{})
	m1.MAN, m1.RST = true, true
	m1.Execute(time.Time{})
	if m1.Q || m1.STATUS != 102 {
		t.Errorf("MANUAL_1 = %v %v", m1.Q, m1.STATUS)
	}
	var m2 MANUAL_2
	m2.ENA, m2.ON, m2.OFF, m2.MAN = true, true, true, true
	m2.Execute(time.Time{})
	if !m2.Q || m2.STATUS != 103 {
		t.Errorf("MANUAL_2 = %v %v", m2.Q, m2.STATUS)
	}
	var m4 MANUAL_4
	m4.MAN, m4.STP = true, true
	m4.Execute(time.Time{})
	m4.STP = false
	m4.Execute(time.Time{})
	m4.STP = true
	m4.Execute(time.Time{})
	if !m4.Q1 || m4.Q0 || m4.STATUS != 111 {
		t.Errorf("MANUAL_4 = %+v", m4)
	}
}

func TestParset(t *testing.T) {
	c := newClock()
	var p PARSET
	p.X01, p.X11, p.X21 = 1, 2, 3
	p.TC = secs(10)
	p.Execute(c.now)
	p.Execute(c.stepS(20))
	if p.P1 != 1 {
		t.Fatalf("PARSET set 0 = %v", p.P1)
	}
	p.A0 = true
	p.Execute(c.now)
	p.Execute(c.stepS(5))
	if !near(p.P1, 1.5, 0.01) {
		t.Fatalf("PARSET ramp = %v", p.P1)
	}
	p.Execute(c.stepS(6))
	if p.P1 != 2 {
		t.Fatalf("PARSET after ramp = %v", p.P1)
	}

	var p2 PARSET2
	p2.X01, p2.X11, p2.X21, p2.X31 = 1, 2, 3, 4
	p2.L1, p2.L2, p2.L3 = 10, 20, 30
	p2.X = -25
	p2.Execute(c.now)
	p2.Execute(c.stepS(1))
	if p2.P1 != 3 {
		t.Fatalf("PARSET2 = %v", p2.P1)
	}
}

func TestSignal(t *testing.T) {
	c := newClock()
	var s SIGNAL_4
	s.INIT()
	s.IN3 = true
	s.TS = secs(1)
	var q []bool
	for i := 0; i < 8; i++ {
		s.Execute(c.now)
		q = append(q, bool(s.Q))
		c.stepS(1)
	}
	ones := 0
	for _, b := range q {
		if b {
			ones++
		}
	}
	if ones != 4 {
		t.Fatalf("SIGNAL_4 pattern 2#1010_1010 = %v", q)
	}
}

func TestSramp(t *testing.T) {
	c := newClock()
	var s SRAMP
	s.X, s.A_UP, s.A_DN, s.VU_MAX, s.VD_MAX = 10, 1, -1, 2, -2
	s.LIMIT_HIGH, s.LIMIT_LOW = 100, -100
	s.Execute(c.now)
	for i := 0; i < 200; i++ {
		s.Execute(c.stepS(0.1))
		if s.V > 2.0001 {
			t.Fatalf("SRAMP speed %v over VU_MAX", s.V)
		}
	}
	if !near(s.Y, 10, 0.05) {
		t.Fatalf("SRAMP Y = %v", s.Y)
	}
}

func TestTune(t *testing.T) {
	c := newClock()
	var tn TUNE
	tn.INIT()
	tn.Y = 50
	tn.SU = true
	tn.Execute(c.now)
	tn.SU = false
	tn.Execute(c.stepS(0.1))
	if !near(tn.Y, 50.1, 1e-4) {
		t.Fatalf("TUNE step = %v", tn.Y)
	}
	tn.SD = true
	tn.Execute(c.now)
	tn.Execute(c.stepS(1.5))
	if !near(tn.Y, 48.1, 1e-3) {
		t.Fatalf("TUNE ramp = %v", tn.Y)
	}
	tn.SET = true
	tn.Execute(c.now)
	if tn.Y != 100 {
		t.Fatal("TUNE SET")
	}

	var t2 TUNE2
	t2.INIT()
	t2.FU = true
	t2.Execute(c.now)
	t2.FU = false
	t2.Execute(c.stepS(0.1))
	if t2.Y != 5 {
		t2.Execute(c.now)
		t.Fatalf("TUNE2 fast step = %v", t2.Y)
	}
	t2.FD = true
	t2.Execute(c.now)
	t2.Execute(c.stepS(1.5))
	if !near(t2.Y, -5, 1e-3) && t2.Y != 0 {
		t.Fatalf("TUNE2 ramp = %v", t2.Y)
	}
}

func TestControl(t *testing.T) {
	if BAND_B(3, 5) != 0 || BAND_B(252, 5) != 255 || BAND_B(100, 5) != 100 {
		t.Error("BAND_B")
	}
	var cs CONTROL_SET1
	cs.INIT()
	cs.KT, cs.TT, cs.PID = 10, 2, true
	cs.Execute(time.Time{})
	if cs.KP != 6 || cs.TN != 1 || cs.TV != 0.25 || cs.KI != 6 || cs.KD != 1.5 {
		t.Errorf("CONTROL_SET1 = %+v", cs)
	}
	var cs2 CONTROL_SET2
	cs2.INIT()
	cs2.KS, cs2.TU, cs2.TG, cs2.PI = 2, 1, 10, true
	cs2.Execute(time.Time{})
	if !near(cs2.KP, 4.5, 1e-5) || !near(cs2.TN, 3.33, 1e-5) {
		t.Errorf("CONTROL_SET2 = %+v", cs2)
	}
	if CTRL_IN(10, 9.9, 0.2) != 0 || CTRL_IN(10, 9, 0.2) != 1 {
		t.Error("CTRL_IN")
	}
	if DEAD_BAND(5, 2) != 3 || DEAD_BAND(-5, 2) != -3 || DEAD_BAND(1, 2) != 0 {
		t.Error("DEAD_BAND")
	}
	var co CTRL_OUT
	co.CI, co.OFFSET, co.LIM_L, co.LIM_H = 150, 10, 0, 100
	co.Execute(time.Time{})
	if co.Y != 100 || !co.LIM {
		t.Error("CTRL_OUT")
	}
	var dz DEAD_ZONE2
	dz.X, dz.L = 5, 2
	dz.Execute(time.Time{})
	dz.X = 1
	dz.Execute(time.Time{})
	if dz.Y != 2 {
		t.Errorf("DEAD_ZONE2 = %v", dz.Y)
	}
	h := HYST{ON: 10, OFF: 5}
	for _, in := range []iec.REAL{11, 7} {
		h.IN = in
		h.Execute(time.Time{})
	}
	if !h.Q || !h.WIN {
		t.Error("HYST")
	}
	h3 := HYST_3{HYST: 2, VAL1: 10, VAL2: 20}
	h3.IN = 8
	h3.Execute(time.Time{})
	h3.IN = 22
	h3.Execute(time.Time{})
	if h3.Q1 || !h3.Q2 {
		t.Error("HYST_3")
	}
}

func TestFilters(t *testing.T) {
	c := newClock()
	var pt1 FT_PT1
	pt1.INIT()
	pt1.T = secs(1)
	pt1.Execute(c.now)
	pt1.IN = 1
	for i := 0; i < 1000; i++ {
		pt1.Execute(c.stepS(0.001))
	}
	if !near(pt1.OUT, 0.632, 0.01) {
		t.Fatalf("FT_PT1 after T = %v", pt1.OUT)
	}

	var in INTEGRATE
	var y iec.REAL
	in.INIT()
	in.Y = &y
	in.X = 2
	in.Execute(c.now)
	for i := 0; i < 10; i++ {
		in.Execute(c.stepS(0.5))
	}
	if !near(y, 10, 1e-3) {
		t.Fatalf("INTEGRATE = %v", y)
	}

	var fi FT_INT
	fi.INIT()
	fi.IN, fi.OUT_MAX = 1, 3
	fi.Execute(c.now)
	for i := 0; i < 10; i++ {
		fi.Execute(c.stepS(1))
	}
	if fi.OUT != 3 || !fi.LIM {
		t.Fatalf("FT_INT limit = %v %v", fi.OUT, fi.LIM)
	}

	var d FT_DERIV
	d.INIT()
	d.Execute(c.now)
	d.IN = 5
	d.Execute(c.stepS(0.5))
	if !near(d.OUT, 10, 1e-3) {
		t.Fatalf("FT_DERIV = %v", d.OUT)
	}

	var pi FT_PIWL
	pi.INIT()
	pi.KP, pi.KI, pi.LIM_H = 1, 1, 5
	pi.IN = 1
	pi.Execute(c.now)
	for i := 0; i < 20; i++ {
		pi.Execute(c.stepS(1))
	}
	if pi.Y != 5 || !pi.LIM {
		t.Fatalf("FT_PIWL = %v %v", pi.Y, pi.LIM)
	}

	var ctrl CTRL_PI
	ctrl.INIT()
	ctrl.SET, ctrl.ACT = 10, 8
	ctrl.Execute(c.now)
	ctrl.Execute(c.stepS(1))
	if !near(ctrl.Y, 4, 0.01) || ctrl.DIFF != 2 {
		t.Fatalf("CTRL_PI = %v", ctrl.Y)
	}

	// Without INIT the inputs are 0, but the blocks inside start with
	// their initial values: the integrator's limits are not 0.
	var pid FT_PID
	pid.KP = 1
	pid.Execute(c.now)
	pid.IN = 1
	pid.Execute(c.stepS(1))
	if pid.Y != 1 || pid.LIM {
		t.Fatal("FT_PID without INIT did not run with its initial values")
	}

	var tn FT_TN8
	tn.T = secs(8)
	for i := 0; i < 20; i++ {
		tn.IN = iec.REAL(i)
		tn.Execute(c.now)
		c.stepS(1)
	}
	if tn.OUT != 11 {
		t.Fatalf("FT_TN8 = %v", tn.OUT)
	}
}

func TestConversions(t *testing.T) {
	tests := []struct {
		name      string
		got, want iec.REAL
		eps       iec.REAL
	}{
		{"C_TO_F", C_TO_F(100), 212, 1e-4},
		{"F_TO_C", F_TO_C(212), 100, 1e-3},
		{"C_TO_K", C_TO_K(0), 273.15, 1e-3},
		{"K_TO_C", K_TO_C(0), -273.15, 1e-3},
		{"KMH_TO_MS", KMH_TO_MS(36), 10, 1e-4},
		{"MS_TO_KMH", MS_TO_KMH(10), 36, 1e-4},
		{"BFT_TO_MS", BFT_TO_MS(4), 6.688, 1e-3},
		{"F_TO_OM", F_TO_OM(1), 2 * MATH.PI, 1e-5},
		{"OM_TO_F", OM_TO_F(MATH.PI2), 1, 1e-5},
		{"PT_TO_F", PT_TO_F(secs(0.5)), 2, 1e-5},
		{"GEO_TO_DEG", GEO_TO_DEG(10, 30, 36), 10.51, 1e-4},
		{"RES_PT", RES_PT(100, 100), 138.5, 0.01},
		{"RES_PT neg", RES_PT(-100, 100), 60.26, 0.02},
		{"TEMP_PT", TEMP_PT(138.5, 100), 100, 0.05},
		{"TEMP_PT neg", TEMP_PT(60.26, 100), -100, 0.05},
		{"RES_NI", RES_NI(100, 100), 161.8, 0.2},
		{"TEMP_NI", TEMP_NI(161.8, 100), 100, 0.5},
		{"RES_NTC", RES_NTC(25, 10000, 3950), 10000, 1},
		{"TEMP_NTC", TEMP_NTC(10000, 10000, 3950), 25, 0.01},
		{"RES_SI", RES_SI(25, 2000, 25), 2000, 1e-3},
		{"TEMP_SI", TEMP_SI(2000, 2000, 25), 25, 0.05},
		{"SENSOR_INT", SENSOR_INT(10, 0.1, 1e9, 0), 100, 0.01},
		{"MULTI_IN avg", MULTI_IN(1, 2, 300, 0, 0, 100, 0), 1.5, 0},
		{"MULTI_IN min", MULTI_IN(5, 2, 3, 0, 0, 100, 5), 2, 0},
		{"MULTI_IN max", MULTI_IN(5, 2, 3, 0, 0, 100, 6), 5, 0},
		{"MULTI_IN mid", MULTI_IN(5, 2, 3, 0, 0, 100, 7), 3, 0},
		{"MULTI_IN default", MULTI_IN(-5, 200, 300, 42, 0, 100, 0), 42, 0},
	}
	for _, tt := range tests {
		if !near(tt.got, tt.want, tt.eps) {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
	if F_TO_PT(4) != secs(0.25) || MS_TO_BFT(6.7) != 4 {
		t.Error("F_TO_PT, MS_TO_BFT")
	}
	for _, tt := range []struct {
		deg, n iec.INT
		want   iec.STRING
	}{{0, 1, "N"}, {90, 1, "E"}, {45, 2, "NE"}, {200, 3, "SSW"}, {359, 3, "N"}} {
		if got := DEG_TO_DIR(tt.deg, tt.n, 0); got != tt.want {
			t.Errorf("DEG_TO_DIR(%d, %d) = %q, want %q", tt.deg, tt.n, got, tt.want)
		}
	}
	if DIR_TO_DEG("SW", 1) != 225 || DIR_TO_DEG("ONO", 2) != 68 {
		t.Error("DIR_TO_DEG")
	}
	var tp TEMPERATURE
	tp.INIT()
	tp.C = 100
	tp.Execute(time.Time{})
	if !near(tp.YF, 212, 0.01) || !near(tp.YK, 373.15, 0.01) {
		t.Errorf("TEMPERATURE = %+v", tp)
	}
	var l LENGTH
	l.FT = 1
	l.Execute(time.Time{})
	if !near(l.YIN, 12, 1e-4) {
		t.Errorf("LENGTH = %v", l.YIN)
	}
	var p PRESSURE
	p.ATM = 1
	p.Execute(time.Time{})
	if !near(p.YPA, 101325, 1) {
		t.Errorf("PRESSURE = %v", p.YPA)
	}
	var s SPEED
	s.KN = 1
	s.Execute(time.Time{})
	if !near(s.YKMH, 1.852, 1e-3) {
		t.Errorf("SPEED = %v", s.YKMH)
	}
	var e ENERGY
	e.WH = 1
	e.Execute(time.Time{})
	if e.YJ != 3600 {
		t.Errorf("ENERGY = %v", e.YJ)
	}
	var a ASTRO
	a.LJ = 1
	a.Execute(time.Time{})
	if !near(a.YPC, 0.3066, 1e-3) {
		t.Errorf("ASTRO = %v", a.YPC)
	}
}

func TestMeasurements(t *testing.T) {
	c := newClock()
	var a ALARM_2
	a.LO_1, a.HI_1, a.LO_2, a.HI_2, a.HYS = 10, 90, 5, 95, 2
	a.X = 92
	a.Execute(c.now)
	if !a.Q1_HI || a.Q2_HI || a.Q1_LO {
		t.Errorf("ALARM_2 = %+v", a)
	}
	var b BAR_GRAPH
	b.TRIGGER_LOW, b.TRIGGER_HIGH, b.ALARM_HIGH = 0, 70, true
	b.X = 25
	b.Execute(c.now)
	if !b.Q3 || b.STATUS != 110 {
		t.Errorf("BAR_GRAPH = %+v", b)
	}
	b.X = 80
	b.Execute(c.now)
	b.X = 25
	b.Execute(c.now)
	if !b.HIGH || !b.ALARM {
		t.Error("BAR_GRAPH alarm does not latch")
	}
	var cal CALIBRATE
	cal.Y_OFFSET, cal.Y_SCALE = 0, 100
	cal.X, cal.CO = 2, true
	cal.Execute(c.now)
	cal.CO, cal.CS, cal.X = false, true, 52
	cal.Execute(c.now)
	cal.CS, cal.X = false, 27
	cal.Execute(c.now)
	if !near(cal.Y, 50, 1e-4) {
		t.Errorf("CALIBRATE = %v", cal.Y)
	}

	var m M_TX
	m.INIT()
	for i := 0; i < 3; i++ {
		m.IN = true
		m.Execute(c.now)
		c.stepS(0.3)
		m.IN = false
		m.Execute(c.now)
		c.stepS(0.7)
	}
	if !near(m.F, 1, 1e-4) || !near(m.DC, 0.3, 1e-4) || m.TH != secs(0.3) {
		t.Errorf("M_TX = F %v DC %v TH %v", m.F, m.DC, m.TH)
	}
	var mt M_T
	mt.INIT()
	mt.IN = true
	mt.Execute(c.now)
	mt.Execute(c.stepS(2))
	mt.IN = false
	mt.Execute(c.now)
	if mt.PT != secs(2) {
		t.Errorf("M_T = %v", mt.PT)
	}
	var md M_D
	md.INIT()
	md.START = true
	md.Execute(c.now)
	md.START = false
	md.Execute(c.now)
	md.START = true
	md.Execute(c.now)
	md.Execute(c.stepS(3))
	md.STOP = true
	md.Execute(c.now)
	if md.PT != secs(3) {
		t.Errorf("M_D = %v", md.PT)
	}

	var sec, cyc iec.UDINT
	var o ONTIME
	o.SECONDS, o.CYCLES = &sec, &cyc
	o.IN = true
	for i := 0; i < 25; i++ {
		o.Execute(c.stepS(0.1))
	}
	if sec != 2 || cyc != 1 {
		t.Errorf("ONTIME = %v s %v cycles", sec, cyc)
	}

	var mx iec.REAL
	var mtr METER
	mtr.INIT()
	mtr.MX = &mx
	mtr.M1, mtr.I1 = 3600, true
	mtr.Execute(c.now)
	for i := 0; i < 10; i++ {
		mtr.Execute(c.stepS(1))
	}
	if !near(mx, 36000, 1) {
		t.Errorf("METER = %v", mx)
	}

	var tcms TC_MS
	tcms.Execute(c.now)
	tcms.Execute(c.stepS(0.25))
	if tcms.TC != 250 {
		t.Errorf("TC_MS = %v", tcms.TC)
	}
	var ct CYCLE_TIME
	for i := 0; i < 5; i++ {
		ct.Execute(c.stepS(0.01 * float64(i+1)))
	}
	if ct.CYCLES != 5 || ct.CT_MAX != secs(0.05) || ct.CT_MIN != secs(0.02) {
		t.Errorf("CYCLE_TIME = %+v", ct)
	}
	var sim DT_SIMU
	sim.INIT()
	sim.START = DWORD_TO_DT(1000)
	sim.SPEED = 10
	sim.Execute(c.now)
	sim.Execute(c.stepS(1))
	if DT_TO_DWORD(sim.DTS) != 1010 {
		t.Errorf("DT_SIMU = %v", DT_TO_DWORD(sim.DTS))
	}
	var x iec.REAL
	var y iec.UDINT
	var fm FLOW_METER
	fm.INIT()
	fm.X, fm.Y = &x, &y
	fm.VX, fm.E = 3600, true
	fm.Execute(c.now)
	for i := 0; i < 5; i++ {
		fm.Execute(c.stepS(1))
	}
	if y+iec.UDINT(x) < 4 || !near(fm.F, 3600, 1) {
		t.Errorf("FLOW_METER = %v + %v, F %v", y, x, fm.F)
	}
}

func TestSignalProcessing(t *testing.T) {
	c := newClock()
	if got := AIN(0x0FFF, 12, 255, 0, 10); got != 10 {
		t.Errorf("AIN = %v", got)
	}
	if got := AIN(0x8800, 12, 15, 0, 4095); got != -2048 {
		t.Errorf("AIN with sign = %v", got)
	}
	if got := AOUT(5, 12, 255, 0, 10); got != 2048 {
		t.Errorf("AOUT = %v", got)
	}
	if got := AOUT(-10, 12, 15, 0, 10); got != 0x8FFF {
		t.Errorf("AOUT negative = %#x", got)
	}
	if got := AOUT1(10, 4, 11, 255, 0, 10); got != 0xFF0 {
		t.Errorf("AOUT1 = %#x", got)
	}
	var a AIN1
	a.INIT()
	a.IN = 0xFFFFFFFF
	a.Execute(c.now)
	if a.OUT != 10 {
		t.Errorf("AIN1 = %v", a.OUT)
	}
	checks := []struct {
		name      string
		got, want iec.REAL
	}{
		{"BYTE_TO_RANGE", BYTE_TO_RANGE(255, 0, 10), 10},
		{"WORD_TO_RANGE", WORD_TO_RANGE(65535, 0, 10), 10},
		{"MIX", MIX(0, 10, 0.25), 2.5},
		{"MUX_R4", MUX_R4(0, 1, 2, 3, false, true), 2},
		{"OFFSET", OFFSET(1, true, false, true, false, false, 10, 20, 30, 40, 0), 41},
		{"OFFSET2", OFFSET2(1, true, false, true, false, true, 10, 20, 30, 40, 5), 35},
		{"OVERRIDE", OVERRIDE(1, -5, 3, true, true, true), -5},
		{"SCALE", SCALE(5, 2, 1, 10, 0), 10},
		{"SCALE_D", SCALE_D(50, 0, 100, 0, 10), 5},
		{"SCALE_R", SCALE_R(0.5, 0, 1, 10, 20), 15},
		{"SCALE_B2", SCALE_B2(255, 0, 1, 0, 0, 1000, 0, 1000), 1000},
		{"SCALE_X2", SCALE_X2(true, false, 1, 0, 0, 1000, 5, 1000), 1005},
		{"STAIR", STAIR(7.4, 2), 8},
	}
	for _, tt := range checks {
		if !near(tt.got, tt.want, 1e-3) {
			t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
	if RANGE_TO_BYTE(5, 0, 10) != 127 || RANGE_TO_WORD(10, 0, 10) != 65535 {
		t.Error("RANGE_TO_BYTE, RANGE_TO_WORD")
	}

	var d DELAY
	d.N = 3
	for i := 1; i <= 5; i++ {
		d.IN = iec.REAL(i)
		d.Execute(c.now)
	}
	if d.OUT != 2 {
		t.Errorf("DELAY = %v", d.OUT)
	}
	var d4 DELAY_4
	for i := 1; i <= 5; i++ {
		d4.IN = iec.REAL(i)
		d4.Execute(c.now)
	}
	if d4.OUT1 != 4 || d4.OUT4 != 1 {
		t.Errorf("DELAY_4 = %+v", d4)
	}
	var mav FILTER_MAV_DW
	mav.N = 4
	for _, v := range []iec.DWORD{0, 4, 4, 4, 4} {
		mav.X = v
		mav.Execute(c.now)
	}
	if mav.Y != 4 {
		t.Errorf("FILTER_MAV_DW = %v", mav.Y)
	}
	var wav FILTER_WAV
	wav.W[0], wav.W[1] = 0.5, 0.5
	wav.X = 2
	wav.Execute(c.now)
	wav.X = 4
	wav.Execute(c.now)
	if wav.Y != 3 {
		t.Errorf("FILTER_WAV = %v", wav.Y)
	}
	var fw FILTER_W
	fw.T = secs(1)
	fw.Execute(c.now)
	fw.X = 1000
	fw.Execute(c.stepS(0.5))
	if fw.Y != 500 {
		t.Errorf("FILTER_W = %v", fw.Y)
	}
	var s3 SEL2_OF_3
	s3.IN1, s3.IN2, s3.IN3, s3.D = 10, 10.5, 20, 1
	s3.Execute(c.now)
	if s3.W != 3 || s3.Y != 10.25 {
		t.Errorf("SEL2_OF_3 = %+v", s3)
	}
	var s3b SEL2_OF_3B
	s3b.IN1, s3b.IN2, s3b.TD = true, true, secs(1)
	s3b.Execute(c.now)
	s3b.Execute(c.stepS(1))
	if !s3b.Q || !s3b.W {
		t.Errorf("SEL2_OF_3B = %+v", s3b)
	}
	var sh SH_2
	sh.INIT()
	sh.N, sh.DISC, sh.PT = 5, 2, secs(1)
	for _, v := range []iec.REAL{1, 100, 3, 4, 5} {
		sh.IN = v
		sh.Execute(c.stepS(1))
	}
	if sh.AVG != 4 || sh.LOW != 3 || sh.HIGH != 5 {
		t.Errorf("SH_2 = %+v", sh)
	}
	var st STAIR2
	st.D = 1
	for _, v := range []iec.REAL{2.2, 2.9, 3.3} {
		st.X = v
		st.Execute(c.now)
	}
	if st.Y != 3 {
		t.Errorf("STAIR2 = %v", st.Y)
	}
	var tr TREND_DW
	tr.X = 5
	tr.Execute(c.now)
	tr.X = 3
	tr.Execute(c.now)
	if !tr.TD || tr.D != 2 || tr.Q {
		t.Errorf("TREND_DW = %+v", tr)
	}
	var fd FADE
	fd.IN1, fd.IN2, fd.TF, fd.F = 0, 10, secs(1), true
	fd.Execute(c.now)
	fd.Execute(c.stepS(0.1))
	fd.Execute(c.stepS(0.5))
	if !near(fd.Y, 5, 0.01) {
		t.Errorf("FADE = %v", fd.Y)
	}
}

func TestGenerators(t *testing.T) {
	c := newClock()
	var gp GEN_PULSE
	gp.INIT()
	gp.PTH, gp.PTL = secs(1), secs(3)
	gp.Execute(c.now)
	highs := 0
	for i := 0; i < 40; i++ {
		gp.Execute(c.stepS(0.1))
		if gp.Q {
			highs++
		}
	}
	if highs != 10 {
		t.Errorf("GEN_PULSE high scans = %v", highs)
	}
	var sq GEN_SQR
	sq.INIT()
	sq.PT = secs(1)
	sq.Execute(c.now)
	highs = 0
	for i := 0; i < 100; i++ {
		sq.Execute(c.stepS(0.01))
		if sq.Q {
			highs++
		}
	}
	if highs < 49 || highs > 51 {
		t.Errorf("GEN_SQR high scans = %v", highs)
	}
	var sn GEN_SIN
	sn.INIT()
	sn.PT = secs(1)
	sn.Execute(c.now)
	sn.Execute(c.stepS(0.25))
	if !near(sn.OUT, 0.5, 1e-3) {
		t.Errorf("GEN_SIN = %v", sn.OUT)
	}
	var rm GEN_RMP
	rm.INIT()
	rm.Execute(c.now)
	rm.Execute(c.stepS(0.5))
	if !near(rm.OUT, 0.5, 1e-3) {
		t.Errorf("GEN_RMP = %v", rm.OUT)
	}
	var rd GEN_RDM
	rd.INIT()
	rd.PT = secs(1)
	rd.Execute(c.now)
	rd.Execute(c.stepS(1))
	if !rd.Q || rd.OUT < -0.5 || rd.OUT > 0.5 {
		t.Errorf("GEN_RDM = %v", rd.OUT)
	}
	var pw2 GEN_PW2
	pw2.ENQ, pw2.TH1, pw2.TL1 = true, secs(1), secs(1)
	pw2.Execute(c.now)
	pw2.Execute(c.stepS(1))
	if !pw2.Q {
		t.Error("GEN_PW2 did not go high")
	}

	var rb RMP_B
	rb.INIT()
	rb.PT = secs(2.55)
	rb.Execute(c.now)
	rb.Execute(c.stepS(1))
	if rb.OUT != 100 || !rb.BUSY {
		t.Errorf("RMP_B = %v", rb.OUT)
	}
	rb.Execute(c.stepS(2))
	if rb.OUT != 255 || !rb.HIGH {
		t.Errorf("RMP_B end = %v", rb.OUT)
	}
	var rw RMP_W
	rw.INIT()
	rw.PT = secs(1)
	// The ramp starts the scan after its direction is taken.
	rw.Execute(c.now)
	rw.Execute(c.stepS(0.1))
	rw.Execute(c.stepS(0.5))
	if rw.OUT != 32768 {
		t.Errorf("RMP_W = %v", rw.OUT)
	}
	var rs RMP_SOFT
	rs.INIT()
	rs.IN, rs.VAL, rs.PT_ON = true, 100, secs(2.55)
	rs.Execute(c.now)
	for i := 0; i < 20; i++ {
		rs.Execute(c.stepS(0.1))
	}
	if rs.OUT != 100 {
		t.Errorf("RMP_SOFT = %v", rs.OUT)
	}
	var out iec.BYTE
	var rn RMP_NEXT_
	rn.INIT()
	rn.OUT = &out
	rn.TR, rn.TF = secs(2.55), secs(2.55)
	rn.Execute(c.now)
	rn.IN = 50
	for i := 0; i < 20; i++ {
		rn.Execute(c.stepS(0.1))
	}
	if out < 50 || out > 52 {
		t.Errorf("RMP_NEXT_ = %v", out)
	}

	var pwm PWM_DC
	pwm.F, pwm.DC = 1, 0.25
	highs = 0
	for i := 0; i < 100; i++ {
		pwm.Execute(c.stepS(0.01))
		if pwm.Q {
			highs++
		}
	}
	if highs < 24 || highs > 26 {
		t.Errorf("PWM_DC high scans = %v", highs)
	}
	var rdt GEN_RDT
	rdt.INIT()
	pulses := 0
	for i := 0; i < 1000; i++ {
		rdt.Execute(c.stepS(0.01))
		if rdt.XQ {
			pulses++
		}
	}
	if pulses == 0 {
		t.Error("GEN_RDT made no pulses")
	}
}
