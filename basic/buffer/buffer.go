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
	"strings"
)

// BUFFER_CLEAR initializes a given slice of bytes with 0.
// It is the Go equivalent of the _BUFFER_CLEAR function from the OSCAT library.
// The function modifies the original slice in-place.
func BUFFER_CLEAR(buffer []byte, size uint) {
	// Determine the number of bytes to clear, ensuring we don't go
	// beyond the actual length of the slice to prevent a panic.
	clearLen := int(size)
	if clearLen > len(buffer) {
		clearLen = len(buffer)
	}
	// A simple loop is the most idiomatic and efficient way to zero a slice.
	for i := range buffer[:clearLen] {
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
	strLen := len(str)
	insertPos := pos

	// Determine the effective buffer length from SIZE, capped by the actual slice length.
	bufferLen := int(size)
	if bufferLen > len(buffer) {
		bufferLen = len(buffer)
	}

	// --- Boundary Checks ---
	// Ensure the insertion position is valid.
	if insertPos < 0 || insertPos > bufferLen-strLen {
		// If pos is invalid, do nothing and return the original position.
		return pos
	}

	// Ensure the buffer has enough space for the insertion.
	if insertPos+strLen > bufferLen {
		// Not enough space within the specified SIZE, do nothing.
		return pos
	}

	// --- Shift existing data ---
	// Make space for the new string by shifting the elements from the insertion
	// point to the right. Go's `copy` handles overlapping slices correctly.
	copy(buffer[insertPos+strLen:], buffer[insertPos:])

	// --- Copy the new string ---
	// Copy the string content into the newly created space.
	copy(buffer[insertPos:], []byte(str))

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
	startPos := int(start)
	if startPos >= bufferLen {
		return ""
	}

	stopPos := int(stop)
	if stopPos >= bufferLen {
		stopPos = bufferLen - 1
	}

	// If start is after stop, there's nothing to extract.
	if startPos > stopPos {
		return ""
	}

	// Extract the relevant portion of the buffer and convert it to a string.
	// Go slices make this operation simple and safe.
	return string(buffer[startPos : stopPos+1])
}

// STRING_TO_BUFFER copies a string into a byte buffer starting at a specific position.
// It is the Go equivalent of the _STRING_TO_BUFFER function from the OSCAT library.
// The function returns the position in the buffer immediately after the inserted string.
func STRING_TO_BUFFER(buffer []byte, str string, pos int, size uint) int {
	strLen := len(str)
	startPos := pos

	// Determine the effective buffer length from SIZE, capped by the actual slice length.
	bufferLen := int(size)
	if bufferLen > len(buffer) {
		bufferLen = len(buffer)
	}

	// --- Boundary Checks ---
	if strLen == 0 || startPos < 0 || startPos >= bufferLen {
		// If there's nothing to copy or the start position is invalid,
		// return the original position.
		return pos
	}

	// Determine how many bytes can be copied without overflowing the buffer.
	bytesToCopy := strLen
	if startPos+bytesToCopy > bufferLen {
		bytesToCopy = bufferLen - startPos
	}

	// Use Go's built-in `copy` function, which is highly optimized for this task.
	copy(buffer[startPos:], []byte(str[:bytesToCopy]))

	// Return the position after the copied string.
	return pos + bytesToCopy
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
