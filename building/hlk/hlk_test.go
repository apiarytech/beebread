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

package hlk

import (
	gomath "math"
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

func near(t *testing.T, name string, got, want, tol iec.REAL) {
	t.Helper()
	if gomath.Abs(float64(got-want)) > float64(tol) {
		t.Errorf("%s = %v, want %v ± %v", name, got, want, tol)
	}
}

// The physics against reference values: OSCAT's documentation and the
// usual tables.
func TestFunctions(t *testing.T) {
	near(t, "DEW_TEMP(60, 25)", DEW_TEMP(60, 25), 16.7, 0.1)
	near(t, "DEW_TEMP(0, 25)", DEW_TEMP(0, 25), -273.15, 0.01)
	near(t, "SDD(20)", SDD(20, false), 2339, 5)
	near(t, "SDD(-10, ice)", SDD(-10, true), 260, 2)
	near(t, "AIR_DENSITY(20, PN, 0)", AIR_DENSITY(20, 101325, 0), 1.204, 0.002)
	near(t, "DEW_CON(100, 20)", DEW_CON(100, 20), 17.3, 0.1)
	near(t, "DEW_RH(DEW_CON(50, 20), 20)", DEW_RH(DEW_CON(50, 20), 20), 50, 0.01)
	near(t, "AIR_ENTHALPY(20, 0)", AIR_ENTHALPY(20, 0), 20.1, 0.01)
	near(t, "WATER_DENSITY(4)", WATER_DENSITY(4, false), 999.97, 0.01)
	near(t, "WATER_DENSITY(20)", WATER_DENSITY(20, false), 998.2, 0.05)
	near(t, "WATER_CP(50)", WATER_CP(50), 4.181, 1e-4)
	near(t, "WATER_CP(55)", WATER_CP(55), 4.182, 1e-4)
	near(t, "WATER_ENTHALPY(50)", WATER_ENTHALPY(50), 209.4, 1e-3)
	near(t, "HEAT_INDEX(30, 70)", HEAT_INDEX(30, 70), 35.0, 0.5)
	near(t, "HEAT_INDEX(15, 70)", HEAT_INDEX(15, 70), 15, 0)
	near(t, "WCT(-10, 30)", WCT(-10, 30), -19.5, 0.1)
	near(t, "WCT(15, 30)", WCT(15, 30), 15, 0)
	near(t, "SDD_NH3(-33.4)", SDD_NH3(-33.4), 1.0, 0.05)
	near(t, "SDT_NH3(SDD_NH3(20))", SDT_NH3(SDD_NH3(20)), 20, 0.5)
	near(t, "SDT_NH3(0)", SDT_NH3(0), -110, 0)
	near(t, "TANK_VOL2(1, 2)", TANK_VOL2(1, 2), 4.18879, 1e-4)
	near(t, "TANK_VOL1(1, 2, 2)", TANK_VOL1(1, 2, 2), 2*3.14159, 1e-3)
}

func TestBoiler(t *testing.T) {
	var b BOILER
	b.INIT()
	b.T_UPPER = 40
	b.Execute(time.Time{})
	if !b.HEAT || b.STATUS != 102 {
		t.Fatalf("below T_UPPER_MIN: HEAT %v STATUS %d", b.HEAT, b.STATUS)
	}
	b.T_UPPER = 61
	b.Execute(time.Time{})
	if b.HEAT || b.STATUS != 100 {
		t.Fatalf("above T_UPPER_MAX: HEAT %v STATUS %d", b.HEAT, b.STATUS)
	}
	b.T_UPPER = 90
	b.Execute(time.Time{})
	if !b.ERROR || b.STATUS != 1 {
		t.Fatalf("over T_PROTECT_HIGH: ERROR %v STATUS %d", b.ERROR, b.STATUS)
	}
}

func TestBurner(t *testing.T) {
	now := time.Unix(1000, 0)
	var rt1, rt2, cycles iec.UDINT
	b := BURNER{RUNTIME1: &rt1, RUNTIME2: &rt2, CYCLES: &cycles}
	b.INIT()
	run := func(d time.Duration) {
		now = now.Add(d)
		b.Execute(now)
	}
	run(0) // power up
	b.IN = true
	run(time.Millisecond)
	if !b.PRE_HEAT {
		t.Fatal("IN starts the pre heating")
	}
	run(5 * time.Second) // OIL_TEMP is TRUE: the motor starts
	if !b.MOTOR {
		t.Fatal("the motor runs after the pre heating")
	}
	run(10 * time.Second) // PRE_VENT_TIME - PRE_IGNITE_TIME
	if !b.IGNITE {
		t.Fatal("the ignition starts")
	}
	run(5 * time.Second)
	if !b.COIL1 {
		t.Fatal("the oil valve opens")
	}
	run(5 * time.Second) // no flame within SAFETY_TIME
	if !b.FAIL || b.STATUS != 4 || b.MOTOR {
		t.Fatalf("no flame locks the burner out: FAIL %v STATUS %d", b.FAIL, b.STATUS)
	}
}

func TestHeatTemp(t *testing.T) {
	var h HEAT_TEMP
	h.INIT()
	h.T_INT = 20
	h.T_EXT = -15
	h.Execute(time.Time{})
	near(t, "TY at the design point", h.TY, 70, 0.01)
	h.T_EXT = 18
	h.Execute(time.Time{})
	if h.TY != 0 || h.HEAT {
		t.Errorf("no heat within H: TY %v", h.TY)
	}
}

func TestTempExt(t *testing.T) {
	var x TEMP_EXT
	x.INIT()
	x.T_EXT1 = 5
	x.T_EXT_CONFIG = 1
	x.DT_IN = td.SET_DT(2024, 1, 15, 12, 0, 0)
	// TEMP_EXT runs every CYCLE_TIME from the start of the PLC timer.
	x.Execute(time.Now().Add(time.Hour))
	if x.T_EXT != 5 || !x.HEAT || x.COOL {
		t.Errorf("January at 5 °C: T_EXT %v HEAT %v COOL %v", x.T_EXT, x.HEAT, x.COOL)
	}
	x.T_EXT1 = 30
	x.DT_IN = td.SET_DT(2024, 7, 15, 12, 0, 0)
	x.Execute(time.Now().Add(2 * time.Hour))
	if x.HEAT || !x.COOL {
		t.Errorf("July at 30 °C: HEAT %v COOL %v", x.HEAT, x.COOL)
	}
}

func TestTAvg24(t *testing.T) {
	var t24, tmax, tmin iec.REAL = -1000, 0, 0
	a := T_AVG24{T24: &t24, T24_MAX: &tmax, T24_MIN: &tmin}
	a.INIT()
	a.TS = 200 // 20.0 °C
	a.DTI = td.SET_DT(2024, 1, 1, 0, 10, 0)
	a.Execute(time.Unix(1000, 0))
	near(t, "T24 at the start", t24, 20, 0.01)
	if !a.TP {
		t.Error("TP at the start")
	}
	a.DTI = td.SET_DT(2024, 1, 1, 0, 30, 0)
	a.Execute(time.Unix(1001, 0))
	if !a.TP {
		t.Error("TP at a sample")
	}
	near(t, "T24_MAX", tmax, 20, 0.01)
}

func TestTankLevel(t *testing.T) {
	now := time.Unix(1000, 0)
	var x TANK_LEVEL
	x.INIT()
	x.MAX_VALVE_TIME = iec.TIME(time.Minute)
	x.LEVEL = true
	x.Execute(now)
	x.Execute(now.Add(time.Millisecond))
	if !x.VALVE || x.STATUS != 102 {
		t.Fatalf("a low level opens the valve: VALVE %v STATUS %d", x.VALVE, x.STATUS)
	}
	x.Execute(now.Add(2 * time.Minute))
	if x.VALVE || !x.ALARM || x.STATUS != 2 {
		t.Fatalf("the valve open too long: VALVE %v ALARM %v STATUS %d", x.VALVE, x.ALARM, x.STATUS)
	}
	x.ACLR = true
	x.Execute(now.Add(3 * time.Minute))
	if x.ALARM {
		t.Fatal("ACLR clears the alarm")
	}
}

func TestHeatMeter(t *testing.T) {
	var y iec.REAL
	h := HEAT_METER{Y: &y}
	h.INIT()
	h.PULSE_MODE = true
	h.TF, h.TR, h.LPH = 60, 40, 10
	h.E = true
	h.Execute(time.Unix(1000, 0))
	// The density is at TF, as the meter is not on the return line.
	want := WATER_DENSITY(60, false) * (WATER_ENTHALPY(60) - WATER_ENTHALPY(40)) * 10
	near(t, "Y after a pulse", y, want, want*1e-6)
	h.E = false
	h.Execute(time.Unix(1001, 0))
	h.E = true
	h.Execute(time.Unix(1002, 0))
	near(t, "Y after two pulses", y, 2*want, want*1e-6)
	h.RST = true
	h.Execute(time.Unix(1003, 0))
	if y != 0 {
		t.Errorf("RST clears Y: %v", y)
	}
}

func TestLegionella(t *testing.T) {
	var l LEGIONELLA
	l.INIT()
	l.MANUAL = true
	l.TEMP_BOILER = 50
	l.DT_IN = DWORD_TO_DT(0)
	l.Execute(time.Unix(1000, 0))
	if !l.RUN || !l.HEAT {
		t.Errorf("MANUAL starts a disinfection that heats: RUN %v HEAT %v STATUS %d", l.RUN, l.HEAT, l.STATUS)
	}
}
