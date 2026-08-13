/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */
package buffer

import (
	"bytes"
	"testing"
)

func TestBufferClear(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5}
	BufferClear(buf, 5)
	expected := []byte{0, 0, 0, 0, 0}
	if !bytes.Equal(buf, expected) {
		t.Errorf("BufferClear failed: expected %v, got %v", expected, buf)
	}

	buf = []byte{1, 2, 3, 4, 5}
	BufferClear(buf, 3)
	expected = []byte{0, 0, 0, 4, 5}
	if !bytes.Equal(buf, expected) {
		t.Errorf("BufferClear partial failed: expected %v, got %v", expected, buf)
	}

	buf = []byte{1, 2, 3}
	BufferClear(buf, 10) // Size > len(buf)
	expected = []byte{0, 0, 0}
	if !bytes.Equal(buf, expected) {
		t.Errorf("BufferClear with oversized SIZE failed: expected %v, got %v", expected, buf)
	}
}

func TestBufferInit(t *testing.T) {
	buf := make([]byte, 5)
	BufferInit(buf, 5, 7)
	expected := []byte{7, 7, 7, 7, 7}
	if !bytes.Equal(buf, expected) {
		t.Errorf("BufferInit failed: expected %v, got %v", expected, buf)
	}

	buf = []byte{1, 2, 3, 4, 5}
	BufferInit(buf, 3, 8)
	expected = []byte{8, 8, 8, 4, 5}
	if !bytes.Equal(buf, expected) {
		t.Errorf("BufferInit partial failed: expected %v, got %v", expected, buf)
	}

	buf = make([]byte, 3)
	BufferInit(buf, 10, 9) // Size > len(buf)
	expected = []byte{9, 9, 9}
	if !bytes.Equal(buf, expected) {
		t.Errorf("BufferInit with oversized SIZE failed: expected %v, got %v", expected, buf)
	}
}

func TestBufferInsert(t *testing.T) {
	tests := []struct {
		name     string
		buffer   []byte
		str      string
		pos      int
		size     uint
		expected []byte
		wantPos  int
	}{
		{
			name: "Insert at beginning",
			// buffer with enough capacity to hold the inserted string
			buffer:   append([]byte("abc"), make([]byte, 3)...),
			str:      "XYZ",
			pos:      0,
			size:     6,
			expected: []byte("XYZabc"),
			wantPos:  3,
		},
		{
			name:     "Insert in middle",
			buffer:   append([]byte("abcdef"), make([]byte, 3)...),
			str:      "123",
			pos:      3,
			size:     9,
			expected: []byte("abc123def"),
			wantPos:  6,
		},
		{
			name:     "Insert at end",
			buffer:   append([]byte("abc"), make([]byte, 3)...),
			str:      "XYZ",
			pos:      3,
			size:     6,
			expected: []byte("abcXYZ"),
			wantPos:  6,
		},
		{
			name:     "Not enough space",
			buffer:   []byte("abcdef"),
			str:      "1234",
			pos:      3,
			size:     6,
			expected: []byte("abcdef"),
			wantPos:  3,
		},
		{
			name:     "Invalid position",
			buffer:   []byte("abc"),
			str:      "XYZ",
			pos:      -1,
			size:     6,
			expected: []byte("abc"),
			wantPos:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufCopy := make([]byte, len(tt.buffer))
			copy(bufCopy, tt.buffer)

			gotPos := BufferInsert(bufCopy, tt.str, tt.pos, tt.size) // bufCopy is the buffer to be modified
			if gotPos != tt.wantPos {
				t.Errorf("BufferInsert() returned position %v, want %v", gotPos, tt.wantPos)
			}

			// Only check buffer content if an insertion was expected
			if tt.wantPos != tt.pos {
				if !bytes.Equal(bufCopy[:len(tt.expected)], tt.expected) {
					t.Errorf("BufferInsert() buffer content = %q, want %q", bufCopy, tt.expected)
				}
			}
		})
	}
}

func TestBufferUppercase(t *testing.T) {
	buf := []byte("aBcDeFg-123")
	BufferUppercase(buf, len(buf))
	expected := []byte("ABCDEFG-123")
	if !bytes.Equal(buf, expected) {
		t.Errorf("BufferUppercase failed: expected %q, got %q", expected, buf)
	}
}

func TestBufferSearch(t *testing.T) {
	buffer := []byte("Hello World, this is a test. hello world again.")

	tests := []struct {
		name       string
		str        string
		pos        int
		ignoreCase bool
		want       int
	}{
		{"Found case-sensitive", "World", 0, false, 6},
		{"Not found case-sensitive", "WoRlD", 0, false, -1}, // Changed search string to ensure it's not found case-sensitively
		{"Found case-insensitive", "world", 0, true, 6},     // Case-insensitive search should find it
		{"Found second occurrence", "world", 10, true, 35},  // Search from pos 10
		{"Not found", "galaxy", 0, false, -1},
		{"Empty string", "", 0, false, -1},
		{"Out of bounds", "test", 100, false, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BufferSearch(buffer, len(buffer), tt.str, tt.pos, tt.ignoreCase)
			if got != tt.want {
				t.Errorf("BufferSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBufferToString(t *testing.T) {
	buffer := []byte("This is a test string.")

	tests := []struct {
		name  string
		start uint
		stop  uint
		want  string
	}{
		{"Extract middle", 5, 11, "is a te"},
		{"Extract beginning", 0, 3, "This"},
		{"Extract end", 15, 20, "string"},
		{"Single character", 5, 5, "i"},
		{"Invalid start", 100, 101, ""},
		{"Invalid stop", 0, 100, "This is a test string."},
		{"Start after stop", 10, 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BufferToString(buffer, uint(len(buffer)), tt.start, tt.stop)
			if got != tt.want {
				t.Errorf("BufferToString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringToBuffer(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		pos      int
		size     uint
		expected []byte
		wantPos  int
	}{
		{
			name:     "Simple copy",
			str:      "World",
			pos:      6,
			size:     20,
			expected: []byte("Hello World."),
			wantPos:  11,
		},
		{
			name:     "Copy with truncation",
			str:      "a long string",
			pos:      0,
			size:     7,
			expected: []byte("a long "),
			wantPos:  7,
		},
		{
			name:     "Invalid position",
			str:      "test",
			pos:      -1,
			size:     10,
			expected: []byte("Hello....."),
			wantPos:  -1,
		},
		{
			name:     "Empty string",
			str:      "",
			pos:      0,
			size:     10,
			expected: []byte("Hello....."),
			wantPos:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initial buffer state for each test
			buffer := []byte("Hello World....") // Ensure enough space and correct initial content
			gotPos := StringToBuffer(buffer, tt.str, tt.pos, tt.size)

			if gotPos != tt.wantPos {
				t.Errorf("StringToBuffer() returned pos %v, want %v", gotPos, tt.wantPos)
			}

			// Only check buffer if a change was expected
			if tt.pos >= 0 && len(tt.str) > 0 {
				// Determine the part of the buffer that should have changed
				end := tt.pos + len(tt.str)
				if end > int(tt.size) {
					end = int(tt.size)
				}
				if end > len(buffer) {
					end = len(buffer)
				}
				if !bytes.Equal(buffer[:end], tt.expected[:end]) {
					t.Errorf("StringToBuffer() buffer content mismatch, got %q, want prefix %q", buffer, tt.expected)
				}
			}
		})
	}
}

func TestBufferComp(t *testing.T) {
	buffer1 := []byte("This is a test buffer for testing.")

	tests := []struct {
		name    string
		buffer2 []byte
		start   int
		want    int
	}{
		{
			name:    "Found in middle",
			buffer2: []byte("test"),
			start:   0,
			want:    10,
		},
		{
			name:    "Found at beginning",
			buffer2: []byte("This"),
			start:   0,
			want:    0,
		},
		{
			name:    "Found second occurrence",
			buffer2: []byte("test"),
			start:   11, // Start searching after the first "test"
			want:    26, // It should find "test" inside "testing"
		},
		{
			name:    "Not found",
			buffer2: []byte("galaxy"),
			start:   0,
			want:    -1,
		},
		{
			name:    "Empty search buffer",
			buffer2: []byte{},
			start:   0,
			want:    -1,
		},
		{
			name:    "Search buffer too long",
			buffer2: []byte("This is a very long search buffer that cannot possibly be in the original."),
			start:   0,
			want:    -1,
		},
		{
			name:    "Invalid start index",
			buffer2: []byte("test"),
			start:   -5,
			want:    -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BufferComp(buffer1, len(buffer1), tt.buffer2, len(tt.buffer2), tt.start)
			if got != tt.want {
				t.Errorf("BufferComp() = %v, want %v", got, tt.want)
			}
		})
	}
}
