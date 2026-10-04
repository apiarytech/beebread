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

package actuators

import (
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

func TestActuatorCoil(t *testing.T) {
	now := time.Unix(1000, 0)
	var a ACTUATOR_COIL
	a.INIT()
	a.SELF_ACT_CYCLE = iec.TIME(time.Minute)
	a.Execute(now)
	a.IN = true
	a.Execute(now.Add(time.Second))
	if !a.OUT || a.STATUS != 101 {
		t.Fatalf("IN: OUT %v STATUS %d", a.OUT, a.STATUS)
	}
	a.IN = false
	a.Execute(now.Add(2 * time.Second))
	if a.OUT || a.STATUS != 100 {
		t.Fatalf("idle: OUT %v STATUS %d", a.OUT, a.STATUS)
	}
	// A minute off activates it for SELF_ACT_TIME.
	a.Execute(now.Add(61 * time.Second))
	if !a.OUT || a.STATUS != 102 {
		t.Fatalf("self activation: OUT %v STATUS %d", a.OUT, a.STATUS)
	}
	a.Execute(now.Add(62*time.Second + time.Millisecond))
	if a.OUT || a.STATUS != 100 {
		t.Fatalf("after self activation: OUT %v STATUS %d", a.OUT, a.STATUS)
	}
}

func TestActuatorPump(t *testing.T) {
	start := time.Unix(1000, 0)
	now := start
	var runtime, cycles iec.UDINT
	p := ACTUATOR_PUMP{RUNTIME: &runtime, CYCLES: &cycles}
	p.INIT()
	// run scans every 100 ms until d after the start; ONTIME counts a
	// second a scan at most.
	run := func(d time.Duration) {
		for !now.After(start.Add(d)) {
			p.Execute(now)
			now = now.Add(100 * time.Millisecond)
		}
	}
	run(0)
	p.IN = true
	run(5 * time.Second)
	if p.PUMP {
		t.Fatal("the pump waits for MIN_OFFTIME")
	}
	run(10 * time.Second)
	if !p.PUMP || cycles != 1 {
		t.Fatalf("after MIN_OFFTIME: PUMP %v CYCLES %d", p.PUMP, cycles)
	}
	p.IN = false
	run(15 * time.Second)
	if !p.PUMP {
		t.Fatal("the pump runs for MIN_ONTIME")
	}
	run(20 * time.Second)
	if p.PUMP || runtime != 10 {
		t.Fatalf("after MIN_ONTIME: PUMP %v RUNTIME %d", p.PUMP, runtime)
	}
	p.RST = true
	run(21 * time.Second)
	if runtime != 0 || cycles != 0 || p.RST {
		t.Fatalf("RST: RUNTIME %d CYCLES %d RST %v", runtime, cycles, p.RST)
	}
}

func TestActuatorUD(t *testing.T) {
	now := time.Unix(1000, 0)
	var a ACTUATOR_UD
	a.INIT()
	a.TON = iec.TIME(time.Second)
	a.TOFF = iec.TIME(time.Second)
	a.Execute(now)
	a.UD, a.ON = true, true
	a.Execute(now.Add(2 * time.Second))
	if !a.YUP || a.YDN || a.STATUS != 111 {
		t.Fatalf("up: YUP %v YDN %v STATUS %d", a.YUP, a.YDN, a.STATUS)
	}
	// A change of direction turns both off first, and waits for TOFF.
	a.UD = false
	a.Execute(now.Add(4 * time.Second))
	if a.YUP || a.YDN {
		t.Fatal("both off on a change of direction")
	}
	a.Execute(now.Add(4500 * time.Millisecond))
	if a.YDN {
		t.Fatal("down within TOFF")
	}
	a.Execute(now.Add(6 * time.Second))
	if !a.YDN || a.STATUS != 112 {
		t.Fatalf("down: YDN %v STATUS %d", a.YDN, a.STATUS)
	}
	a.OFF = true
	a.Execute(now.Add(8 * time.Second))
	if a.YUP || a.YDN || a.STATUS != 101 {
		t.Fatal("OFF")
	}
}

func TestActuator2P(t *testing.T) {
	now := time.Unix(1000, 0)
	var arx iec.BOOL
	a := ACTUATOR_2P{ARX: &arx}
	a.INIT()
	a.CYCLE_TIME = iec.TIME(10 * time.Second)
	a.SENS = 10
	a.IN = 5
	a.Execute(now)
	if a.OUT || a.ARO {
		t.Fatal("IN below SENS is off")
	}
	a.IN = 250
	a.Execute(now.Add(time.Second))
	if !a.OUT {
		t.Fatal("IN above 255 - SENS is on")
	}
	// TEST self activates, and ARX tells the other actuators.
	a.TEST = true
	a.Execute(now.Add(2 * time.Second))
	if !a.ARO || !arx {
		t.Fatalf("TEST: ARO %v ARX %v", a.ARO, arx)
	}
}

func TestActuatorA(t *testing.T) {
	now := time.Unix(1000, 0)
	var a ACTUATOR_A
	a.INIT()
	a.RUNTIME = iec.TIME(time.Second)
	a.OUT_MIN, a.OUT_MAX = 1000, 2000
	a.I1 = 255
	// Without SELF_ACT_TIME the cycle starts in state 0: OUT_MIN, then
	// OUT_MAX, then normal operation.
	a.Execute(now)
	if a.Y != 1000 {
		t.Fatalf("state 0: Y %d", a.Y)
	}
	a.Execute(now.Add(1500 * time.Millisecond))
	if a.Y != 2000 {
		t.Fatalf("state 1: Y %d", a.Y)
	}
	a.Execute(now.Add(3 * time.Second))
	a.Execute(now.Add(4 * time.Second))
	a.I1 = 51
	a.Execute(now.Add(5 * time.Second))
	if a.Y != 1200 {
		t.Fatalf("normal operation, I1 = 51: Y %d", a.Y)
	}
	a.RV = true
	a.Execute(now.Add(6 * time.Second))
	if a.Y != 1800 {
		t.Fatalf("reversed: Y %d", a.Y)
	}
}
