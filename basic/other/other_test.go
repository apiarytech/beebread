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

package other

import (
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

func TestESR(t *testing.T) {
	now := time.Unix(100, 0)
	var b8out, x8out [4]ESR_DATA

	var b8 ESR_MON_B8
	b8.ESR_OUT = &b8out
	b8.A1 = "door"
	b8.S1 = true
	b8.Execute(now)
	if !b8.ESR_FLAG || b8out[0].TYP != 11 || b8out[0].ADRESS != "door" || b8out[1].TYP != 0 {
		t.Fatalf("ESR_MON_B8 = %+v", b8out[0])
	}
	b8.Execute(now)
	if b8.ESR_FLAG || b8out[0].TYP != 0 {
		t.Fatal("ESR_MON_B8 reported no change")
	}
	// Only 4 changes fit in a run.
	b8.S0, b8.S1, b8.S2, b8.S3, b8.S4 = true, false, true, true, true
	b8.Execute(now)
	b8.Execute(now)
	if b8out[0].TYP != 11 || b8out[1].TYP != 0 {
		t.Fatalf("ESR_MON_B8 did not report S4 in the next run: %+v", b8out)
	}

	var x8 ESR_MON_X8
	x8.INIT()
	x8.ESR_OUT = &x8out
	x8.MODE = 1
	x8.S0, x8.S1 = 150, 5
	x8.Execute(now)
	if x8out[0].TYP != 1 || x8out[0].DATA[0] != 5 || x8out[1].TYP != 0 {
		t.Fatalf("ESR_MON_X8 = %+v", x8out)
	}

	var r4out [4]ESR_DATA
	var r4 ESR_MON_R4
	r4.ESR_OUT = &r4out
	r4.R2, r4.S2 = 1.0, 0.5
	r4.Execute(now)
	if r4out[0].TYP != 20 || r4out[0].DATA[3] != 0x3F || r4out[0].DATA[2] != 0x80 {
		t.Fatalf("ESR_MON_R4 = %+v", r4out[0])
	}

	var pos iec.INT
	var c ESR_COLLECT
	c.POS = &pos
	c.Execute(now)
	if pos != -1 {
		t.Fatalf("ESR_COLLECT first run POS = %v", pos)
	}
	c.ESR_3 = x8out
	c.ESR_5 = r4out
	c.Execute(now)
	if pos != 1 || c.ESR_OUT[0].TYP != 1 || c.ESR_OUT[1].TYP != 20 {
		t.Fatalf("ESR_COLLECT POS %v OUT %+v", pos, c.ESR_OUT[:2])
	}
}

func TestVersion(t *testing.T) {
	if OSCAT_VERSION(false) != 335 {
		t.Error("OSCAT_VERSION")
	}
	if got := OSCAT_VERSION(true); got != 1721088000 {
		t.Errorf("OSCAT_VERSION(TRUE) = %v, want 1721088000", got)
	}
	if e := STATUS_TO_ESR(250, "x", iec.DT{}, 0); e.TYP != 3 || e.DATA[0] != 250 {
		t.Errorf("STATUS_TO_ESR = %+v", e)
	}
}
