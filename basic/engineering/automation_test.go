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
	"testing"
	"time"
)

func TestMANUAL(t *testing.T) {
	testCases := []struct {
		name   string
		in     bool
		on     bool
		off    bool
		expect bool
	}{
		{"IN true, no override", true, false, false, true},
		{"IN false, no override", false, false, false, false},
		{"ON true forces true", false, true, false, true},
		{"OFF true forces false", true, false, true, false},
		{"ON and OFF, OFF wins", true, true, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MANUAL(tc.in, tc.on, tc.off)
			if got != tc.expect {
				t.Errorf("MANUAL(in: %v, on: %v, off: %v) = %v; want %v", tc.in, tc.on, tc.off, got, tc.expect)
			}
		})
	}
}

func TestDRIVER_1(t *testing.T) {
	t.Run("Normal Mode", func(t *testing.T) {
		d := DRIVER_1{}
		// Initial state
		if got := d.Update(false, false, false); got != false {
			t.Errorf("Initial state should be false, got %v", got)
		}
		// Rising edge on IN
		if got := d.Update(false, true, false); got != true {
			t.Errorf("Rising edge on IN should set Q to true, got %v", got)
		}
		// IN stays high
		if got := d.Update(false, true, false); got != true {
			t.Errorf("Q should remain true while IN is high, got %v", got)
		}
		// RST forces false
		if got := d.Update(false, true, true); got != false {
			t.Errorf("RST should force Q to false, got %v", got)
		}
	})

	t.Run("Toggle Mode", func(t *testing.T) {
		d := DRIVER_1{ToggleMode: true}
		// Initial state
		if got := d.Update(false, false, false); got != false {
			t.Errorf("Initial state should be false, got %v", got)
		}
		// First rising edge
		if got := d.Update(false, true, false); got != true {
			t.Errorf("First rising edge should toggle Q to true, got %v", got)
		}
		// Falling edge
		if got := d.Update(false, false, false); got != true {
			t.Errorf("Q should remain true after falling edge, got %v", got)
		}
		// Second rising edge
		if got := d.Update(false, true, false); got != false {
			t.Errorf("Second rising edge should toggle Q to false, got %v", got)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		d := DRIVER_1{Timeout: 50 * time.Millisecond}
		// Rising edge on IN
		if got := d.Update(false, true, false); got != true {
			t.Errorf("Rising edge on IN should set Q to true, got %v", got)
		}
		// Wait for less than timeout
		time.Sleep(25 * time.Millisecond)
		if got := d.Update(false, false, false); got != true {
			t.Errorf("Q should be true before timeout, got %v", got)
		}
		// Wait for timeout to expire
		time.Sleep(70 * time.Millisecond)
		if got := d.Update(false, false, false); got != false {
			t.Errorf("Q should be false after timeout, got %v", got)
		}
	})

	t.Run("Set Overrides Timeout", func(t *testing.T) {
		d := DRIVER_1{Timeout: 50 * time.Millisecond}
		// Rising edge on IN
		if got := d.Update(false, true, false); got != true {
			t.Errorf("Rising edge on IN should set Q to true, got %v", got)
		}
		// Wait for timeout to expire
		time.Sleep(55 * time.Millisecond)
		// SET should force Q to true even though timer has expired
		if got := d.Update(true, false, false); got != true {
			t.Errorf("SET should force Q to true, overriding timeout; got %v", got)
		}
	})

	t.Run("Timeout Retrigger", func(t *testing.T) {
		d := DRIVER_1{Timeout: 100 * time.Millisecond}
		// 1. First rising edge on IN starts the timer
		d.Update(false, true, false)
		d.Update(false, false, false) // IN goes low
		if got := d.Update(false, false, false); got != true {
			t.Fatalf("Q should be true after first pulse, got %v", got)
		}

		// 2. Wait for 60ms, Q should still be true
		time.Sleep(60 * time.Millisecond)
		if got := d.Update(false, false, false); got != true {
			t.Errorf("Q should be true before timeout is halfway, got %v", got)
		}

		// 3. Second rising edge should re-trigger the timeout
		d.Update(false, true, false)
		d.Update(false, false, false) // IN goes low

		// 4. Wait for another 60ms. Total time is 120ms, but since the timer was
		// re-triggered, Q should still be true.
		time.Sleep(10 * time.Millisecond)
		if got := d.Update(false, false, false); got != true {
			t.Errorf("Q should be true after re-trigger, got %v", got)
		}
	})

	t.Run("Timeout in Toggle Mode", func(t *testing.T) {
		d := DRIVER_1{ToggleMode: true, Timeout: 50 * time.Millisecond}
		// Toggle ON
		d.Update(false, true, false)
		time.Sleep(60 * time.Millisecond)
		if got := d.Update(false, false, false); got != false {
			t.Errorf("Q should be false after timeout in toggle mode, got %v", got)
		}
	})
}

func TestINC_DEC(t *testing.T) {
	t.Run("Reset", func(t *testing.T) {
		id := INC_DEC{Cnt: 10}
		id.Update(false, false, true)
		if id.Cnt != 0 {
			t.Errorf("RST should reset Cnt to 0, got %d", id.Cnt)
		}
	})

	t.Run("Quadrature Decoding", func(t *testing.T) {
		id := INC_DEC{}
		// One full quadrature cycle forward
		// Step 1: CHa rises
		id.Update(true, false, false)
		if id.Cnt != 1 {
			t.Errorf("Step 1: Expected Cnt 1, got %d", id.Cnt)
		}
		// Step 2: CHb rises
		id.Update(true, true, false)
		if id.Cnt != 2 {
			t.Errorf("Step 2: Expected Cnt 2, got %d", id.Cnt)
		}
		// Step 3: CHa falls
		id.Update(false, true, false)
		if id.Cnt != 3 {
			t.Errorf("Step 3: Expected Cnt 3, got %d", id.Cnt)
		}
		// Step 4: CHb falls
		id.Update(false, false, false)
		if id.Cnt != 4 {
			t.Errorf("Step 4: Expected Cnt 4, got %d", id.Cnt)
		}

		// One full quadrature cycle backward
		// Step 5: CHb rises
		id.Update(false, true, false)
		if id.Cnt != 3 {
			t.Errorf("Step 5: Expected Cnt 3, got %d", id.Cnt)
		}
		// Step 6: CHa rises
		id.Update(true, true, false)
		if id.Cnt != 2 {
			t.Errorf("Step 6: Expected Cnt 2, got %d", id.Cnt)
		}
		// Step 7: CHb falls
		id.Update(true, false, false)
		if id.Cnt != 1 {
			t.Errorf("Step 7: Expected Cnt 1, got %d", id.Cnt)
		}
		// Step 8: CHa falls
		id.Update(false, false, false)
		if id.Cnt != 0 {
			t.Errorf("Step 8: Expected Cnt 0, got %d", id.Cnt)
		}
	})
}

func TestINTERLOCK(t *testing.T) {
	il := INTERLOCK{}
	tl := 50 * time.Millisecond

	// Test basic lockout
	q1, q2 := il.Update(true, false, tl)
	if !q1 || q2 {
		t.Errorf("Basic I1: expected q1=true, q2=false; got q1=%v, q2=%v", q1, q2)
	}

	q1, q2 = il.Update(false, true, tl)
	if q1 || !q2 {
		t.Errorf("Basic I2: expected q1=false, q2=true; got q1=%v, q2=%v", q1, q2)
	}

	// Test dead time
	il.Update(true, false, tl) // Q1 is on, T1.Q is true
	time.Sleep(20 * time.Millisecond)
	q1, q2 = il.Update(false, true, tl) // Try to switch to Q2 within dead time
	if q1 || q2 {
		t.Errorf("Dead time fail: expected q1=false, q2=false; got q1=%v, q2=%v", q1, q2)
	}

	time.Sleep(tl) // Wait for dead time to expire
	q1, q2 = il.Update(false, true, tl)
	if q1 || !q2 {
		t.Errorf("After dead time: expected q1=false, q2=true; got q1=%v, q2=%v", q1, q2)
	}
}

func TestMANUAL_1(t *testing.T) {
	m := MANUAL_1{}

	// Auto mode
	m.Update(true, false, false, false, false)
	if !m.Q || m.Status != 100 {
		t.Errorf("Auto mode (true): expected Q=true, Status=100; got Q=%v, Status=%d", m.Q, m.Status)
	}
	m.Update(false, false, false, false, false)
	if m.Q || m.Status != 100 {
		t.Errorf("Auto mode (false): expected Q=false, Status=100; got Q=%v, Status=%d", m.Q, m.Status)
	}

	// Manual mode
	m.Update(false, true, true, false, false)
	if !m.Q || m.Status != 103 {
		t.Errorf("Manual mode (M_I=true): expected Q=true, Status=103; got Q=%v, Status=%d", m.Q, m.Status)
	}

	// Manual Set
	m.Update(false, true, false, true, false) // Rising edge on SET
	if !m.Q || m.Status != 101 {
		t.Errorf("Manual Set: expected Q=true, Status=101; got Q=%v, Status=%d", m.Q, m.Status)
	}
	m.Update(false, true, false, false, false) // SET goes low, Q should hold
	if !m.Q || m.Status != 103 {
		t.Errorf("Manual Set hold: expected Q=true, Status=103; got Q=%v, Status=%d", m.Q, m.Status)
	}

	// Manual Reset
	m.Update(false, true, true, false, true) // Rising edge on RST
	if m.Q || m.Status != 102 {
		t.Errorf("Manual Reset: expected Q=false, Status=102; got Q=%v, Status=%d", m.Q, m.Status)
	}
}

func TestFT_PROFILE(t *testing.T) {
	p := FT_PROFILE{
		Time1:  100 * time.Millisecond,
		Value1: 10.0,
		Time2:  200 * time.Millisecond,
		Value2: 10.0,
		Time3:  300 * time.Millisecond,
		Value3: 0.0,
		// The rest of the times/values are zero by default
	}

	// Initial state
	if p.Run || p.Y != 0.0 {
		t.Fatalf("Initial state should be Run=false, Y=0.0")
	}

	// Trigger profile
	p.Update(1.0, 0.0, 1.0, true)
	if !p.Run {
		t.Fatalf("Profile should be running after trigger")
	}

	// Check ramp up
	time.Sleep(50 * time.Millisecond)
	p.Update(1.0, 0.0, 1.0, true)
	if p.Y <= 4.0 || p.Y >= 6.0 {
		t.Errorf("Ramp up @ 50ms: expected ~5.0, got %f", p.Y)
	}

	// Check plateau
	time.Sleep(100 * time.Millisecond) // Total elapsed ~150ms
	p.Update(1.0, 0.0, 1.0, true)
	if p.Y <= 9.9 || p.Y > 10.0 {
		t.Errorf("Plateau @ 150ms: expected 10.0, got %f", p.Y)
	}

	// Check ramp down
	time.Sleep(100 * time.Millisecond) // Total elapsed ~250ms
	p.Update(1.0, 0.0, 1.0, true)
	if p.Y <= 4.0 || p.Y >= 6.0 {
		t.Errorf("Ramp down @ 250ms: expected ~5.0, got %f", p.Y)
	}

	// Check end of profile
	time.Sleep(100 * time.Millisecond) // Total elapsed > 300ms
	p.Update(1.0, 0.0, 1.0, true)
	if p.Run || p.Y != 0.0 {
		t.Errorf("End of profile: expected Run=false, Y=0.0, got Run=%v, Y=%f", p.Run, p.Y)
	}
}
