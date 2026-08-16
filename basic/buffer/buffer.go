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
	str "beebread/basic/string"
	"bytes"
	"strings"
)

// BUFFER_CLEAR initializes a given slice of bytes with 0.
// It is the Go equivalent of the _BUFFER_CLEAR function from the OSCAT library.
// The function modifies the original slice in-place.
func BUFFER_CLEAR(buffer []byte, size uint) {
	clearLen := int(size)
	if clearLen > len(buffer) {
		clearLen = len(buffer)
	}

	// Efficiently zero out the specified portion of the slice.
	for i := 0; i < clearLen; i++ {
		buffer[i] = 0
	}
}

// BUFFER_INIT initializes a given slice of bytes with a specific value.
// It is the Go equivalent of the _BUFFER_INIT function from the OSCAT library.
// The function modifies the original slice in-place and returns true when finished.
func BUFFER_INIT(buffer []byte, size uint, init byte) bool {
	// The most idiomatic and performant way to fill a slice in Go is a simple
	// for loop or bytes.Repeat.
	initLen := int(size)
	if initLen > len(buffer) {
		initLen = len(buffer)
	}

	for i := 0; i < initLen; i++ {
		buffer[i] = init
	}
	return true
}

// BUFFER_INSERT inserts a string into a byte buffer at a given position.
// It is the Go equivalent of the _BUFFER_INSERT function from the OSCAT library.
// The function shifts existing data to make room and then copies the new string.
// It returns the position immediately after the inserted string.
func BUFFER_INSERT(buffer []byte, str string, pos int, size uint) int {
	bufferLen := int(size)
	if bufferLen > len(buffer) {
		bufferLen = len(buffer)
	}

	strLen := len(str)
	if pos < 0 || strLen == 0 || pos+strLen > bufferLen {
		return pos
	}

	// Shift existing data to make space for the new string.
	copy(buffer[pos+strLen:], buffer[pos:bufferLen-strLen])

	// Copy the new string into the created space.
	copy(buffer[pos:], str)

	// Return the position after the inserted string.
	return pos + strLen
}

// BUFFER_UPPERCASE converts all characters in a byte buffer to uppercase.
// It is the Go equivalent of the _BUFFER_UPPERCASE function from the OSCAT library.
func BUFFER_UPPERCASE(buffer []byte, size int) {
	// Determine the number of bytes to process, capped by the slice length.
	procSize := size
	if procSize > len(buffer) {
		procSize = len(buffer)
	}

	// bytes.ToUpper creates a new slice. To modify in-place, iterate.
	for i, b := range buffer[:procSize] {
		if b >= 'a' && b <= 'z' {
			buffer[i] = b - ('a' - 'A')
		}
	}
}

// BUFFER_SEARCH searches for a string within a byte buffer.
// It is the Go equivalent of the BUFFER_SEARCH function from the OSCAT library.
// The function returns the starting position of the found string, or -1 if not found.
func BUFFER_SEARCH(buffer []byte, size int, str string, pos int, ignoreCase bool) int {
	strLen := len(str)
	startPos := pos

	// Determine the effective buffer length, capped by the actual slice length.
	bufferLen := int(size)
	if bufferLen > len(buffer) {
		bufferLen = len(buffer)
	}

	// --- Boundary Checks ---
	if strLen == 0 || startPos < 0 || startPos >= bufferLen || startPos+strLen > bufferLen {
		return -1
	}

	// Use the standard library for searching, which is highly optimized.
	if ignoreCase {
		// For case-insensitive search, we search for the uppercase version of the string
		// in an uppercase version of the buffer.
		searchStr := strings.ToUpper(str)
		uppercaseBuffer := bytes.ToUpper(buffer[startPos:bufferLen])
		foundPos := bytes.Index(uppercaseBuffer, []byte(searchStr))
		if foundPos == -1 {
			return -1
		}
		return foundPos + startPos
	} else {
		foundPos := bytes.Index(buffer[startPos:bufferLen], []byte(str))
		if foundPos == -1 {
			return -1
		}
		return foundPos + startPos
	}
}

// BUFFER_TO_STRING retrieves a string from a byte buffer between a start and stop position.
// It is the Go equivalent of the BUFFER_TO_STRING function from the OSCAT library.
func BUFFER_TO_STRING(buffer []byte, size uint, start uint, stop uint) string {
	bufferLen := int(size)
	if bufferLen > len(buffer) {
		bufferLen = len(buffer)
	}

	// --- Boundary and Sanity Checks ---
	if bufferLen == 0 {
		return ""
	}

	// Ensure start and stop are within the buffer's bounds.
	// OSCAT is 1-based, Go is 0-based. The original ST code is complex,
	// but the intent is to use start and stop as inclusive indices.
	startPos := int(start)
	if startPos < 0 || startPos >= bufferLen {
		return ""
	}

	stopPos := int(stop)     // In Go, the end of a slice is exclusive.
	if stopPos > bufferLen { // So, if stop is 22 (last char), stopPos should be 22.
		stopPos = bufferLen
	}

	// If start is after stop, there's nothing to extract.
	if startPos > stopPos {
		return ""
	}

	// Extract the relevant portion of the buffer and convert it to a string.
	// Go slices make this operation simple and safe.
	return string(buffer[startPos:stopPos])
}

// BUFFER_TO_INT retrieves a string from a byte buffer and converts it to an integer.
// It is a logical implementation of what _BUFFER_INT from OSCAT would do.
func BUFFER_TO_INT(buffer []byte, size, start, stop uint) int {
	// Extract the string representation of the number from the buffer.
	s := BUFFER_TO_STRING(buffer, size, start, stop)
	// Use the existing string-to-integer conversion.
	return str.DEC_TO_INT(s)
}

// STRING_TO_BUFFER copies a string into a byte buffer starting at a specific position.
// It is the Go equivalent of the _STRING_TO_BUFFER function from the OSCAT library.
// The function returns the position in the buffer immediately after the inserted string.
func STRING_TO_BUFFER(buffer []byte, str string, pos int, size uint) int {
	bufferLen := int(size)
	if bufferLen > len(buffer) {
		bufferLen = len(buffer)
	}

	if pos < 0 {
		return pos
	}

	end := pos + len(str)
	if end > bufferLen {
		end = bufferLen
	}

	i := pos
	for ; i < end; i++ {
		buffer[i] = str[i-pos]
	}

	return i
}

// BUFFER_COMP compares two buffers to find the first occurrence of the second buffer within the first.
// It is the Go equivalent of the BUFFER_COMP function from the OSCAT library.
// The function returns the starting position of the found sequence, or -1 if not found.
func BUFFER_COMP(buffer1 []byte, size1 int, buffer2 []byte, size2 int, start int) int {
	len1 := int(size1)
	len2 := int(size2)
	startIndex := start

	// Ensure we don't go beyond the actual length of the slices.
	if len1 > len(buffer1) {
		len1 = len(buffer1)
	}
	if len2 > len(buffer2) {
		len2 = len(buffer2)
	}

	// --- Boundary Checks ---
	// If the second buffer is empty, or longer than the first, or start index is invalid.
	if len2 == 0 || len2 > len1 || startIndex < 0 || startIndex >= len1 {
		return -1
	}

	// Use the standard library's `bytes.Index`, which is highly optimized.
	foundPos := bytes.Index(buffer1[startIndex:len1], buffer2[:len2])
	if foundPos == -1 {
		return -1
	}
	return foundPos + startIndex
}
