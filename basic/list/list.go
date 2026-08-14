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

package math

import (
	"strings"

	. "beebread/basic"
)

// ListAdd appends an element to the end of a list string.
// The list is a string of elements, each prefixed by the separator character.
// It modifies the list pointer and returns false if the new element doesn't fit.
func ListAdd(list *string, sep byte, ins string) bool {
	sepStr := string(sep)
	// Ensure the element to be inserted starts with the separator.
	if !strings.HasPrefix(ins, sepStr) {
		ins = sepStr + ins
	}

	if len(*list)+len(ins) > ListLength {
		return false // Return false if element does not fit.
	}

	*list += ins
	return true
}

// ListClean removes empty elements from a list string.
// It modifies the list pointer and returns true.
func ListClean(list *string, sep byte) bool {
	sepStr := string(sep)
	// Split, filter empty strings, and join back.
	parts := strings.Split(*list, sepStr)
	var cleanedParts []string
	for _, p := range parts {
		if p != "" {
			cleanedParts = append(cleanedParts, p)
		}
	}

	// Rebuild the list string, ensuring it starts with a separator if it originally did.
	if strings.HasPrefix(*list, sepStr) {
		*list = sepStr + strings.Join(cleanedParts, sepStr)
	} else {
		*list = strings.Join(cleanedParts, sepStr)
	}

	return true
}

// ListGet retrieves the element at a specific position (1-based) from a list string.
func ListGet(list string, sep byte, pos int) string {
	sepStr := string(sep)
	parts := strings.Split(list, sepStr)

	// Adjust for leading separator
	var elements []string
	if len(parts) > 0 && parts[0] == "" {
		elements = parts[1:]
	} else {
		elements = parts
	}

	if pos > 0 && pos <= len(elements) {
		return elements[pos-1]
	}

	return ""
}

// ListInsert inserts an element at a specific position (1-based) in a list string.
// It modifies the list pointer and returns false if the new element causes an overflow.
func ListInsert(list *string, sep byte, pos int, ins string) bool {
	sepStr := string(sep)
	parts := strings.Split(*list, sepStr)

	// Adjust for leading separator
	var elements []string
	if len(parts) > 0 && parts[0] == "" {
		elements = parts[1:]
	} else {
		elements = parts
	}

	// Clamp the insertion position to be within the valid range.
	if pos < 1 {
		pos = 1
	}
	if pos > len(elements)+1 {
		pos = len(elements) + 1
	}

	// Insert the new element into the slice.
	elements = append(elements[:pos-1], append([]string{ins}, elements[pos-1:]...)...)

	// Rebuild the string and check length.
	newList := sepStr + strings.Join(elements, sepStr)
	if len(newList) > ListLength {
		return false
	}

	*list = newList
	return true
}

// ListLen returns the number of elements in a list string.
func ListLen(list string, sep byte) int {
	if list == "" {
		return 0
	}
	// The OSCAT list format assumes elements are separated by a separator,
	// and often the list starts with one. Counting the separators gives the number of elements.
	return strings.Count(list, string(sep))
}

// ListNext is a stateful block to iterate through elements of a list.
type ListNext struct {
	pos int
}

// Next retrieves the next element from the list string.
// It returns the element, and a boolean `nul` which is true when the end of the list is reached.
// If rst is true, the iterator is reset to the beginning of the list.
func (ln *ListNext) Next(list string, sep byte, rst bool) (lel string, nul bool) {
	if rst {
		ln.pos = 0
	}

	if ln.pos >= len(list) {
		return "", true
	}

	// Find the start of the next element (skip the current separator)
	if list[ln.pos] == sep {
		ln.pos++
	}

	if ln.pos >= len(list) {
		return "", true
	}

	// Find the end of the element (the next separator or end of string)
	end := strings.IndexByte(list[ln.pos:], sep)
	if end == -1 {
		// This is the last element
		lel = list[ln.pos:]
		ln.pos = len(list)
	} else {
		lel = list[ln.pos : ln.pos+end]
		ln.pos += end
	}

	return lel, false
}

// ListRetrieve retrieves an element at a specific position (1-based) and removes it from the list.
// It modifies the list pointer.
func ListRetrieve(list *string, sep byte, pos int) string {
	sepStr := string(sep)
	parts := strings.Split(*list, sepStr)

	var elements []string
	if len(parts) > 0 && parts[0] == "" {
		elements = parts[1:]
	} else {
		elements = parts
	}

	if pos < 1 || pos > len(elements) {
		return ""
	}

	retrieved := elements[pos-1]
	// Remove the element from the slice
	elements = append(elements[:pos-1], elements[pos:]...)

	*list = sepStr + strings.Join(elements, sepStr)
	return retrieved
}

// ListRetrieveLast retrieves the last element from a list and removes it.
// It modifies the list pointer.
func ListRetrieveLast(list *string, sep byte) string {
	sepStr := string(sep)
	lastSep := strings.LastIndex(*list, sepStr)
	if lastSep == -1 || lastSep == len(*list)-1 { // No separator or it's the last char
		return ""
	}

	retrieved := (*list)[lastSep+1:]
	*list = (*list)[:lastSep]
	return retrieved
}
