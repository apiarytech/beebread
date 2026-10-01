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

package list

import (
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

func TestList(t *testing.T) {
	l := iec.STRING(",a,bb")
	if !LIST_ADD(',', "ccc", &l) || l != ",a,bb,ccc" {
		t.Fatalf("LIST_ADD = %q", l)
	}
	if n := LIST_LEN(',', &l); n != 3 {
		t.Fatalf("LIST_LEN = %v", n)
	}
	for pos, want := range map[iec.INT]iec.STRING{1: "a", 2: "bb", 3: "ccc", 4: ""} {
		if got := LIST_GET(',', pos, &l); got != want {
			t.Errorf("LIST_GET(%d) = %q, want %q", pos, got, want)
		}
	}
	if !LIST_INSERT(',', 2, "x", &l) || l != ",a,x,bb,ccc" {
		t.Fatalf("LIST_INSERT = %q", l)
	}
	short := iec.STRING(",a")
	if !LIST_INSERT(',', 3, "z", &short) || short != ",a,,z" {
		t.Fatalf("LIST_INSERT past the end = %q", short)
	}
	if got := LIST_RETRIEVE(',', 2, &l); got != "x" || l != ",a,bb,ccc" {
		t.Fatalf("LIST_RETRIEVE = %q, list %q", got, l)
	}
	if got := LIST_RETRIEVE_LAST(',', &l); got != "ccc" || l != ",a,bb" {
		t.Fatalf("LIST_RETRIEVE_LAST = %q, list %q", got, l)
	}
	dirty := iec.STRING(",a,,b,")
	if !LIST_CLEAN(',', &dirty) || dirty != ",a,b" {
		t.Fatalf("LIST_CLEAN = %q", dirty)
	}
	long := iec.STRING(",")
	for i := 0; i < 300; i++ {
		long += "x"
	}
	if LIST_ADD(',', "y", &long) {
		t.Fatal("LIST_ADD past LIST_LENGTH")
	}
}

func TestListNext(t *testing.T) {
	l := iec.STRING(",a,bb,c")
	var n LIST_NEXT
	n.INIT()
	n.LIST = &l
	n.SEP = ','
	var got []iec.STRING
	for i := 0; i < 4; i++ {
		n.Execute(time.Time{})
		got = append(got, n.LEL)
	}
	if got[0] != "a" || got[1] != "bb" || got[2] != "c" || got[3] != "" || !n.NUL {
		t.Fatalf("LIST_NEXT = %q, NUL %v", got, n.NUL)
	}
	n.RST = true
	n.Execute(time.Time{})
	if n.LEL != "a" || n.NUL {
		t.Fatalf("LIST_NEXT after RST = %q", n.LEL)
	}
}
