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

import "testing"

// TestCoverage checks that every POU and type of OSCAT BASIC is ported.
func TestCoverage(t *testing.T) {
	pous, err := ReadPOUs("../../documents/oscat_basic_335.st")
	if err != nil {
		t.Fatal(err)
	}
	if len(pous) < 550 {
		t.Fatalf("read only %d POUs", len(pous))
	}
	decls, err := ReadDecls("../../basic")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pous {
		if prob := Problem(p, decls); prob != "" {
			t.Errorf("%s %s (%s, line %d): %s", p.Kind, p.Name, p.Subject, p.Line, prob)
		}
	}
}
