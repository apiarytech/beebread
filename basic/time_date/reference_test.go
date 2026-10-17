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

package time_date

import (
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"

	"github.com/apiarytech/royaljelly/iec"
)

// computus returns the date of Easter Sunday in the Gregorian calendar, by
// the anonymous algorithm (Meeus, Jones and Butcher), which OSCAT's EASTER
// does not use.
func computus(y int) (month, day int) {
	a, b, c := y%19, y/100, y%100
	d, e := b/4, b%4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i, k := c/4, c%4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	return (h + l - 7*m + 114) / 31, (h+l-7*m+114)%31 + 1
}

// TestEasterAgainstComputus checks EASTER for every year OSCAT's dates
// hold.
func TestEasterAgainstComputus(t *testing.T) {
	for y := 1970; y <= 2099; y++ {
		m, d := computus(y)
		if got, want := EASTER(iec.INT(y)), date(y, m, d); DATE_TO_DWORD(got) != DATE_TO_DWORD(want) {
			t.Errorf("EASTER(%d) = %s, want %d-%02d-%02d", y, time.Time(got).Format("2006-01-02"), y, m, d)
		}
	}
}
