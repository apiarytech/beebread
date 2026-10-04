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

package logging

import (
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

func TestPrintSF(t *testing.T) {
	var args network.PRINTF_DATA
	args[0], args[1] = "42", "x"
	s := iec.STRING("a=~1, b=~2, c=~0, ~")
	p := PRINT_SF{PRINTF_DATA: &args, STR: &s}
	p.INIT()
	p.Execute(time.Time{})
	if s != "a=42, b=x, c=, " {
		t.Errorf("PRINT_SF = %q", s)
	}
}

func TestLogMsg(t *testing.T) {
	network.LOG_CL = network.NewLOG_CONTROL()
	network.LOG_CL.SIZE = 3
	var l LOG_MSG
	l.INIT()
	for i, m := range []string{"one ~1", "two", "three", "four"} {
		network.LOG_CL.NEW_MSG = iec.STRING(m)
		network.LOG_CL.PRINTF[0] = iec.STRING(rune('0' + i))
		l.Execute(time.Time{})
	}
	lc := network.LOG_CL
	// A ring of 3: four overwrites one.
	if lc.MSG[1] != "four" || lc.MSG[2] != "two" || lc.MSG[3] != "three" || !lc.RING_MODE || lc.IDX != 1 {
		t.Errorf("log: %q, ring %v, idx %d", lc.MSG[:4], lc.RING_MODE, lc.IDX)
	}
	if lc.NEW_MSG != "" || lc.UPDATE_COUNT != 4 {
		t.Errorf("NEW_MSG %q, UPDATE_COUNT %d", lc.NEW_MSG, lc.UPDATE_COUNT)
	}

	// A viewport of 2 lines. OSCAT names 30001 the newest messages and
	// 30000 the oldest, but its arithmetic shows the newest, three and
	// four, at 30000.
	lv := network.US_LOG_VIEWPORT{COUNT: 2, MOVE_TO_X: 30001}
	v := LOG_VIEWPORT{LC: &network.LOG_CL, LV: &lv}
	v.INIT()
	v.Execute(time.Time{})
	if !lv.UPDATE || lv.LINE_ARRAY[0] != 2 || lv.LINE_ARRAY[1] != 3 {
		t.Errorf("30001: %v %v", lv.UPDATE, lv.LINE_ARRAY[:2])
	}
	lv.MOVE_TO_X = 30000
	v.Execute(time.Time{})
	if lv.LINE_ARRAY[0] != 3 || lv.LINE_ARRAY[1] != 1 {
		t.Errorf("30000: %v", lv.LINE_ARRAY[:2])
	}
}
