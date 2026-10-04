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

package jalousie

import (
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

func TestBlindActuator(t *testing.T) {
	start := time.Unix(1000, 0)
	var b BLIND_ACTUATOR
	b.INIT()
	b.T_UD = iec.TIME(2 * time.Second)
	b.T_ANGLE = iec.TIME(500 * time.Millisecond)
	b.UP = true
	for now := start; !now.After(start.Add(4 * time.Second)); now = now.Add(50 * time.Millisecond) {
		b.Execute(now)
	}
	if !b.QU || b.QD || b.STATUS != 121 {
		t.Fatalf("up: QU %v QD %v STATUS %d", b.QU, b.QD, b.STATUS)
	}
	if b.ANG != 255 || b.POS != 255 {
		t.Fatalf("up for longer than T_ANGLE + T_UD: ANG %d POS %d", b.ANG, b.POS)
	}
	b.DN = true
	b.Execute(start.Add(5 * time.Second))
	if b.QU || b.QD || b.STATUS != 1 {
		t.Fatalf("up and down: QU %v QD %v STATUS %d", b.QU, b.QD, b.STATUS)
	}
}

func TestBlindControl(t *testing.T) {
	start := time.Unix(1000, 0)
	var b BLIND_CONTROL
	b.INIT()
	b.T_UD = iec.TIME(2 * time.Second)
	b.T_ANGLE = iec.TIME(500 * time.Millisecond)
	// Automatic mode moves to PI and AI.
	b.UP, b.DN = true, true
	b.PI, b.AI = 128, 128
	for now := start; !now.After(start.Add(10 * time.Second)); now = now.Add(20 * time.Millisecond) {
		b.Execute(now)
	}
	if b.MU || b.MD {
		t.Fatalf("the position is reached: MU %v MD %v", b.MU, b.MD)
	}
	if d := int(b.POS) - 128; d < -5 || d > 5 {
		t.Errorf("POS %d, want 128 within SENS", b.POS)
	}
	if d := int(b.ANG) - 128; d < -10 || d > 10 {
		t.Errorf("ANG %d, want about 128", b.ANG)
	}
}

func TestBlindControlS(t *testing.T) {
	start := time.Unix(1000, 0)
	var b BLIND_CONTROL_S
	b.INIT()
	b.T_UP, b.T_DN = iec.TIME(2*time.Second), iec.TIME(2*time.Second)
	b.T_EXT = iec.TIME(time.Second)
	now := start
	run := func(d time.Duration) {
		for end := now.Add(d); now.Before(end); now = now.Add(20 * time.Millisecond) {
			b.Execute(now)
		}
	}
	// It calibrates when it powers up in automatic mode: in manual standby
	// STATUS is S_IN before the state machine sees 0.
	b.UP, b.DN, b.PI = true, true, 100
	run(40 * time.Millisecond)
	if b.STATUS != 128 || !b.MU {
		t.Fatalf("calibration at power up: STATUS %d MU %v", b.STATUS, b.MU)
	}
	// It runs up for T_UP + T_EXT, then to PI.
	run(3 * time.Second)
	if b.POS != 255 {
		t.Fatalf("after calibration: STATUS %d POS %d", b.STATUS, b.POS)
	}
	run(4 * time.Second)
	if d := int(b.POS) - 100; d < -2 || d > 2 {
		t.Errorf("POS %d, want 100", b.POS)
	}
	if b.MU || b.MD {
		t.Errorf("stopped: MU %v MD %v STATUS %d", b.MU, b.MD, b.STATUS)
	}
}

func TestBlindInput(t *testing.T) {
	start := time.Unix(1000, 0)
	var b BLIND_INPUT
	b.INIT()
	b.POS, b.ANG = 10, 20
	now := start
	run := func(d time.Duration) {
		for end := now.Add(d); now.Before(end); now = now.Add(10 * time.Millisecond) {
			b.Execute(now)
		}
	}
	run(50 * time.Millisecond)
	if b.STATUS != 130 || !b.QU || !b.QD || b.PO != 10 || b.AO != 20 {
		t.Fatalf("automatic operation: STATUS %d QU %v QD %v PO %d AO %d", b.STATUS, b.QU, b.QD, b.PO, b.AO)
	}
	// Holding S1 moves up.
	b.S1 = true
	run(time.Second)
	if b.STATUS != 132 || !b.QU || b.QD {
		t.Fatalf("holding S1: STATUS %d QU %v QD %v", b.STATUS, b.QU, b.QD)
	}
	b.S1 = false
	run(100 * time.Millisecond)
	if b.STATUS != 131 || b.QU || b.QD {
		t.Fatalf("released: STATUS %d QU %v QD %v", b.STATUS, b.QU, b.QD)
	}
	// IN forces PI and AI.
	b.IN, b.PI, b.AI = true, 200, 100
	run(20 * time.Millisecond)
	if b.STATUS != 136 || b.PO != 200 || b.AO != 100 {
		t.Fatalf("IN: STATUS %d PO %d AO %d", b.STATUS, b.PO, b.AO)
	}
}

func TestBlindNight(t *testing.T) {
	var b BLIND_NIGHT
	b.INIT()
	b.UP, b.DN = true, true
	b.PI, b.AI, b.S_IN = 1, 2, 3
	b.SUNSET = DWORD_TO_TOD(18 * 3600000)
	b.SUNRISE = DWORD_TO_TOD(7 * 3600000)
	b.NIGHT_POSITION, b.NIGHT_ANGLE = 0, 0
	b.DTIN = td.SET_DT(2024, 1, 15, 12, 0, 0)
	b.Execute(time.Time{})
	if b.STATUS != 3 || b.PO != 1 {
		t.Fatalf("day: STATUS %d PO %d", b.STATUS, b.PO)
	}
	b.DTIN = td.SET_DT(2024, 1, 15, 19, 0, 0)
	b.Execute(time.Time{})
	if b.STATUS != 141 || b.PO != 0 {
		t.Fatalf("night: STATUS %d PO %d", b.STATUS, b.PO)
	}
	b.DTIN = td.SET_DT(2024, 1, 16, 6, 0, 0)
	b.Execute(time.Time{})
	if b.STATUS != 141 {
		t.Fatal("the night lasts until sunrise")
	}
	b.DTIN = td.SET_DT(2024, 1, 16, 8, 0, 0)
	b.Execute(time.Time{})
	if b.STATUS != 3 {
		t.Fatalf("the next morning: STATUS %d", b.STATUS)
	}
}

func TestBlindScene(t *testing.T) {
	var b BLIND_SCENE
	b.INIT()
	b.UP, b.DN, b.ENABLE = true, true, true
	b.SCENE = 0x13 // scene 3
	b.PI, b.AI, b.SWRITE = 77, 88, true
	b.Execute(time.Time{})
	if b.STATUS != 176 {
		t.Fatalf("write: STATUS %d", b.STATUS)
	}
	b.SWRITE = false
	b.PI, b.AI = 0, 0
	b.Execute(time.Time{})
	if b.STATUS != 163 || b.PO != 77 || b.AO != 88 {
		t.Fatalf("scene 3: STATUS %d PO %d AO %d", b.STATUS, b.PO, b.AO)
	}
	b.UP = false
	b.Execute(time.Time{})
	if b.PO != 0 || b.QU {
		t.Fatalf("manual mode passes the inputs on: PO %d QU %v", b.PO, b.QU)
	}
}

func TestBlindSecurity(t *testing.T) {
	var b BLIND_SECURITY
	b.INIT()
	b.UP, b.DN, b.S_IN = true, true, 7
	for _, c := range []struct {
		fire, wind, alarm, door, rain bool
		qu, qd                        iec.BOOL
		status                        iec.BYTE
	}{
		{fire: true, wind: true, qu: true, status: 111},
		{wind: true, qu: true, status: 112},
		{alarm: true, qu: true, status: 113},
		{door: true, qu: true, status: 114},
		{rain: true, qd: true, status: 115},
		{qu: true, qd: true, status: 7},
	} {
		b.FIRE, b.WIND, b.ALARM, b.DOOR, b.RAIN = iec.BOOL(c.fire), iec.BOOL(c.wind), iec.BOOL(c.alarm), iec.BOOL(c.door), iec.BOOL(c.rain)
		b.Execute(time.Time{})
		if b.QU != c.qu || b.QD != c.qd || b.STATUS != c.status {
			t.Errorf("%+v: QU %v QD %v STATUS %d", c, b.QU, b.QD, b.STATUS)
		}
	}
}

func TestBlindSet(t *testing.T) {
	now := time.Unix(1000, 0)
	var b BLIND_SET
	b.INIT()
	b.RESTORE_POSITION = true
	b.UP, b.DN, b.S_IN = true, true, 9
	b.PI, b.AI = 10, 20
	b.Execute(now)
	b.Execute(now)
	b.IN, b.PX, b.AX = true, 200, 210
	b.Execute(now)
	if b.STATUS != 178 || b.PO != 200 || b.AO != 210 {
		t.Fatalf("forced: STATUS %d PO %d AO %d", b.STATUS, b.PO, b.AO)
	}
	b.IN = false
	b.Execute(now)
	b.Execute(now)
	if b.STATUS != 9 || b.PO != 10 {
		t.Fatalf("restored: STATUS %d PO %d", b.STATUS, b.PO)
	}
}

func TestBlindShade(t *testing.T) {
	now := time.Unix(1000, 0)
	cx := CALENDAR{
		UTC:      td.SET_DT(2024, 6, 21, 12, 0, 0),
		SUN_RISE: DWORD_TO_TOD(5 * 3600000), SUN_SET: DWORD_TO_TOD(21 * 3600000),
		SUN_HOR: 180, SUN_VER: 30,
	}
	b := BLIND_SHADE{CX: &cx}
	b.INIT()
	b.UP, b.DN, b.ENABLE, b.SUN = true, true, true, true
	b.SHADE_POS = 50
	b.Execute(now)
	if b.STATUS != 151 || b.PO != 50 {
		t.Fatalf("shading: STATUS %d PO %d", b.STATUS, b.PO)
	}
	if b.AO == 255 || b.AO == 0 {
		t.Errorf("the slat angle for a sun at 30°: AO %d", b.AO)
	}
	cx.SUN_HOR = 30
	b.Execute(now)
	if b.STATUS == 151 {
		t.Error("the sun on another side")
	}

	s := BLIND_SHADE_S{CX: &cx}
	s.INIT()
	s.UP, s.DN, s.ENABLE, s.SUN = true, true, true, true
	s.PI, s.SHADE_POS = 200, 50
	cx.SUN_HOR = 180
	s.Execute(now)
	if s.STATUS != 151 || s.PO != 50 {
		t.Fatalf("shading a shutter: STATUS %d PO %d", s.STATUS, s.PO)
	}
	s.ALERT = true
	s.Execute(now)
	if s.STATUS != 152 || !s.QU || s.QD {
		t.Fatalf("ALERT: STATUS %d", s.STATUS)
	}
}
