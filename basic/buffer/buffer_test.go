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

package buffer

import (
	"slices"
	"testing"

	"github.com/apiarytech/royaljelly/iec"
)

func bytesOf(s string) []iec.BYTE { return []iec.BYTE(s) }

func TestClearInit(t *testing.T) {
	b := bytesOf("abcde")
	BUFFER_CLEAR_(b, 3)
	if !slices.Equal(b, []iec.BYTE{0, 0, 0, 'd', 'e'}) {
		t.Errorf("BUFFER_CLEAR_ = %v", b)
	}
	BUFFER_INIT_(b, 10, 7)
	if !slices.Equal(b, []iec.BYTE{7, 7, 7, 7, 7}) {
		t.Errorf("BUFFER_INIT_ with a size past the slice = %v", b)
	}
}

func TestStringToBuffer(t *testing.T) {
	b := make([]iec.BYTE, 8)
	if next := STRING_TO_BUFFER_("abc", 2, b, 8); next != 5 {
		t.Errorf("STRING_TO_BUFFER_ = %v, want 5", next)
	}
	if string(b[2:5]) != "abc" || b[5] != 0 {
		t.Errorf("buffer = %q", b)
	}
	// The string is cut at the end of the buffer.
	if next := STRING_TO_BUFFER_("xyz", 6, b, 8); next != 8 || string(b[6:]) != "xy" {
		t.Errorf("STRING_TO_BUFFER_ at the end = %v, %q", next, b)
	}
	// Characters of ISO 8859-1 are one byte.
	c := make([]iec.BYTE, 2)
	STRING_TO_BUFFER_("ä", 0, c, 2)
	if c[0] != 228 {
		t.Errorf("STRING_TO_BUFFER_(ä) = %v", c)
	}
}

func TestInsert(t *testing.T) {
	b := append(bytesOf("abcd"), 0, 0, 0)
	if next := BUFFER_INSERT_("XY", 1, b, 7); next != 3 {
		t.Errorf("BUFFER_INSERT_ = %v, want 3", next)
	}
	if string(b[:6]) != "aXYbcd" {
		t.Errorf("buffer = %q", b)
	}
}

func TestUppercase(t *testing.T) {
	b := bytesOf("abc\xe4")
	BUFFER_UPPERCASE_(b, 4)
	if string(b) != "ABC\xc4" {
		t.Errorf("BUFFER_UPPERCASE_ = %q", b)
	}
}

func TestSearchCompare(t *testing.T) {
	b := bytesOf("hello world")
	if got := BUFFER_SEARCH(b, 11, "world", 0, false); got != 6 {
		t.Errorf("BUFFER_SEARCH = %v, want 6", got)
	}
	if got := BUFFER_SEARCH(b, 11, "WORLD", 0, true); got != 6 {
		t.Errorf("BUFFER_SEARCH ignoring case = %v, want 6", got)
	}
	if got := BUFFER_SEARCH(b, 11, "WORLD", 0, false); got != -1 {
		t.Errorf("BUFFER_SEARCH case = %v, want -1", got)
	}
	if got := BUFFER_SEARCH(b, 11, "o", 5, false); got != 7 {
		t.Errorf("BUFFER_SEARCH from 5 = %v, want 7", got)
	}
	if got := BUFFER_COMP(b, 11, bytesOf("lo w"), 4, 0); got != 3 {
		t.Errorf("BUFFER_COMP = %v, want 3", got)
	}
	if got := BUFFER_COMP(b, 11, bytesOf("xyz"), 3, 0); got != -1 {
		t.Errorf("BUFFER_COMP none = %v, want -1", got)
	}
	if got := BUFFER_COMP(b, 3, bytesOf("hello"), 5, 0); got != -1 {
		t.Errorf("BUFFER_COMP longer = %v, want -1", got)
	}
}

func TestBufferToString(t *testing.T) {
	b := bytesOf("hello\x00world")
	tests := []struct {
		start, stop iec.UINT
		want        iec.STRING
	}{
		{0, 4, "hello"},
		{6, 100, "world"},
		{2, 8, "llo"},
	}
	for _, tt := range tests {
		if got := BUFFER_TO_STRING(b, 11, tt.start, tt.stop); got != tt.want {
			t.Errorf("BUFFER_TO_STRING(%d, %d) = %q, want %q", tt.start, tt.stop, got, tt.want)
		}
	}
	if got := BUFFER_TO_STRING(b, 0, 0, 4); got != "" {
		t.Errorf("BUFFER_TO_STRING of size 0 = %q", got)
	}
	if got := BUFFER_TO_STRING([]iec.BYTE{0xe4}, 1, 0, 0); got != "ä" {
		t.Errorf("BUFFER_TO_STRING(16#E4) = %q", got)
	}
}
