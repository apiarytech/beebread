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

package main

import (
	"os"
	"strings"
	"testing"
)

// golden checks that doc/beedance_<lib>.st is what stclean -refs makes of
// it, with the POUs testdata/<lib>.exclude names and the commented out POUs
// of the cleaned libraries uses: that the POUs it comments out, and why, are
// those the exclude file gives, and every other POU is written with
// references. When beedance runs more, take the POUs it runs out of the
// exclude file and run stclean (see its doc).
func golden(t *testing.T, lib string, uses ...string) {
	t.Helper()
	path := "../../doc/beedance_" + lib + ".st"
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines, err := readSource(path)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := readExclude("testdata/" + lib + ".exclude")
	if err != nil {
		t.Fatal(err)
	}
	refs, excluded = true, ex
	defer func() { refs, excluded = false, map[string]string{} }()
	commented := map[string]string{}
	for _, u := range uses {
		used, err := readLines("../../doc/beedance_" + u + ".st")
		if err != nil {
			t.Fatal(err)
		}
		for _, l := range used {
			if m := reComment.FindStringSubmatch(l); m != nil {
				commented[strings.ToUpper(m[1])] = m[1]
			}
		}
	}
	got := clean(restoreCommented(lines), commented)
	wantLines := strings.Split(strings.TrimSuffix(strings.ReplaceAll(string(want), "\r\n", "\n"), "\n"), "\n")
	for i := 0; i < len(got) || i < len(wantLines); i++ {
		var g, w string
		if i < len(got) {
			g = got[i]
		}
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if g != w {
			t.Fatalf("%s is out of date at line %d: stclean -refs writes\n\t%q\nnot\n\t%q", path, i+1, g, w)
		}
	}
	if !strings.Contains(string(want), "\r\n") {
		t.Errorf("%s should end its lines with CRLF, as stclean writes it", path)
	}
}

// TestBasic checks doc/beedance_basic.st.
func TestBasic(t *testing.T) { golden(t, "basic") }

// TestBuilding checks doc/beedance_building.st, which uses OSCAT BASIC.
func TestBuilding(t *testing.T) { golden(t, "building", "basic") }

// TestNetwork checks doc/beedance_network.st, which uses OSCAT BASIC.
func TestNetwork(t *testing.T) { golden(t, "network", "basic") }

func TestRewrite(t *testing.T) {
	for in, want := range map[string]string{
		"x := TOD#8:0;":                  "x := TOD#08:00:00;",
		"d := D#2011-2-3;":               "d := D#2011-02-03;",
		"d := DT#1970-1-1-00:00;":        "d := DT#1970-01-01-00:00:00;",
		"t := tod#03:00:00;":             "t := tod#03:00:00;",
		"a := b := FALSE;":               "b := FALSE; a := b;",
		"x.single := TRUE; (* single *)": "x.single_ := TRUE; (* single *)",
		"s := 'SINGLE';":                 "s := 'SINGLE';",
		"SINGLE_SWITCH := TRUE;":         "SINGLE_SWITCH := TRUE;",
	} {
		out := rewrite([]string{in})
		if got := out[len(out)-1]; got != want {
			t.Errorf("rewrite(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestRefs checks -refs: a POU a previous run commented out is restored,
// pointers are written as references, and a POU -exclude names is commented
// out again, with the POUs that use it.
func TestRefs(t *testing.T) {
	refs, excluded = true, map[string]string{"STEP": "it computes with pointers as addresses."}
	defer func() { refs, excluded = false, map[string]string{} }()
	src := []string{
		"// Commented out for beedance: it uses POINTER TO or ADR, which IEC 61131-3 does not have.",
		"//FUNCTION GET : BYTE",
		"//VAR_INPUT pt : POINTER TO BYTE; END_VAR",
		"//GET := pt^;",
		"//END_FUNCTION",
		"",
		"FUNCTION STEP : BYTE",
		"VAR pt : POINTER TO BYTE; x : BYTE; END_VAR",
		"pt := ADR(x) + 1; (* the next POINTER TO byte *)",
		"END_FUNCTION",
		"",
		"FUNCTION USE : BYTE",
		"USE := STEP();",
		"END_FUNCTION",
	}
	got := strings.Join(clean(restoreCommented(src), map[string]string{}), "\n")
	for _, want := range []string{
		"FUNCTION GET : BYTE\nVAR_INPUT pt : REF_TO BYTE; END_VAR",
		"// Commented out for beedance: it computes with pointers as addresses.\n// FUNCTION STEP",
		"// pt := REF(x) + 1; (* the next POINTER TO byte *)",
		"// Commented out for beedance: it uses STEP, which is commented out.\n// FUNCTION USE",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "IEC 61131-3 does not have") {
		t.Errorf("the old note is kept:\n%s", got)
	}
}
