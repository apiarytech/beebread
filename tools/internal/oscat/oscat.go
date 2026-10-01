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

// Package oscat reads the POUs of the OSCAT source, for the tools.
package oscat

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// A POU is a FUNCTION, FUNCTION_BLOCK or TYPE of the OSCAT source.
type POU struct {
	Subject string
	Kind    string // "FUNCTION", "FUNCTION_BLOCK" or "TYPE"
	Name    string
	Line    int
}

// GoName returns the name of a POU's port: a leading underscore, which Go
// does not export, moves to the end.
func (p POU) GoName() string {
	if strings.HasPrefix(p.Name, "_") {
		return strings.TrimPrefix(p.Name, "_") + "_"
	}
	return strings.ToUpper(p.Name)
}

var (
	rePath = regexp.MustCompile(`\(\* @PATH := '(.*)' \*\)`)
	// OSCAT misspells a few FUNCTION_BLOCKs as FUNCTIONBLOCK.
	rePOU = regexp.MustCompile(`^\s*(FUNCTION_BLOCK|FUNCTIONBLOCK|FUNCTION|TYPE)\s+([A-Za-z_0-9]+)`)
)

// ReadPOUs reads the POUs of an OSCAT source file, leaving out comments,
// which may nest.
func ReadPOUs(path string) ([]POU, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var pous []POU
	subject := ""
	depth := 0
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for n := 1; sc.Scan(); n++ {
		line := sc.Text()
		if m := rePath.FindStringSubmatch(line); m != nil {
			subject = strings.ReplaceAll(m[1], `\/`, "/")
			continue
		}
		if depth == 0 {
			if m := rePOU.FindStringSubmatch(line); m != nil {
				kind := m[1]
				if kind == "FUNCTIONBLOCK" {
					kind = "FUNCTION_BLOCK"
				}
				s := subject
				if s == "" {
					s = "/Types"
				}
				pous = append(pous, POU{Subject: s, Kind: kind, Name: m[2], Line: n})
			}
		}
		depth += strings.Count(line, "(*") - strings.Count(line, "*)")
		if depth < 0 {
			depth = 0
		}
	}
	return pous, sc.Err()
}

var (
	reSection = regexp.MustCompile(`^\s*(VAR_INPUT|VAR_OUTPUT|VAR_IN_OUT|VAR_TEMP|VAR|END_VAR)\b(\s+CONSTANT)?`)
	reVarDecl = regexp.MustCompile(`^\s*([A-Za-z_0-9, \t]+?)\s*:\s*([^:;]+?)\s*(?::=\s*([^;]*?))?\s*;`)
	reEnd     = regexp.MustCompile(`^\s*END_(FUNCTION|FUNCTION_BLOCK|FUNCTIONBLOCK)\b`)
)

// Inputs returns, for each FUNCTION of an OSCAT source file, the initial
// value of each of its VAR_INPUTs, VAR_INPUT CONSTANTs included, in their
// order and followed by its VAR_IN_OUTs: the ST literal, or "" if it has
// none.
func Inputs(path string) (map[string][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	out := map[string][]string{}
	inOuts := map[string][]string{}
	current := ""
	section := ""
	depth := 0
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		line := sc.Text()
		wasComment := depth > 0
		depth += strings.Count(line, "(*") - strings.Count(line, "*)")
		if depth < 0 {
			depth = 0
		}
		if wasComment {
			continue
		}
		// Leave out a comment at the end of the line.
		if i := strings.Index(line, "(*"); i >= 0 {
			line = line[:i]
		}
		if m := rePOU.FindStringSubmatch(line); m != nil {
			current = ""
			if m[1] == "FUNCTION" {
				current = m[2]
				out[current] = []string{}
			}
			continue
		}
		if reEnd.MatchString(line) {
			current = ""
			continue
		}
		if current == "" {
			continue
		}
		if m := reSection.FindStringSubmatch(line); m != nil {
			section = m[1]
			continue
		}
		m := reVarDecl.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		for range strings.Split(m[1], ",") {
			switch section {
			case "VAR_INPUT":
				out[current] = append(out[current], strings.TrimSpace(m[3]))
			case "VAR_IN_OUT":
				inOuts[current] = append(inOuts[current], "")
			}
		}
	}
	for name, io := range inOuts {
		out[name] = append(out[name], io...)
	}
	return out, sc.Err()
}

var (
	reArray   = regexp.MustCompile(`(?i)^\s*([A-Za-z_0-9, \t]+?)\s*:\s*ARRAY\s*\[([^\]]*)\]`)
	reTypeEnd = regexp.MustCompile(`^\s*END_(TYPE|FUNCTION_BLOCK|FUNCTIONBLOCK)\b`)
)

// ArrayBounds returns, for each FUNCTION_BLOCK and TYPE of an OSCAT source
// file, the lower bounds of its members that are arrays: for each member,
// the lower bound of each dimension.
func ArrayBounds(path string) (map[string]map[string][]int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	out := map[string]map[string][]int64{}
	current := ""
	depth := 0
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		line := sc.Text()
		wasComment := depth > 0
		depth += strings.Count(line, "(*") - strings.Count(line, "*)")
		if depth < 0 {
			depth = 0
		}
		if wasComment {
			continue
		}
		if i := strings.Index(line, "(*"); i >= 0 {
			line = line[:i]
		}
		if m := rePOU.FindStringSubmatch(line); m != nil {
			current = ""
			if m[1] != "FUNCTION" {
				current = m[2]
				out[current] = map[string][]int64{}
			}
			continue
		}
		if reTypeEnd.MatchString(line) || reEnd.MatchString(line) {
			current = ""
			continue
		}
		if current == "" {
			continue
		}
		a := reArray.FindStringSubmatch(line)
		if a == nil {
			continue
		}
		var lows []int64
		for _, r := range strings.Split(a[2], ",") {
			var low int64
			fmt.Sscan(strings.TrimSpace(strings.Split(r, "..")[0]), &low)
			lows = append(lows, low)
		}
		for _, name := range strings.Split(a[1], ",") {
			out[current][strings.ToUpper(strings.TrimSpace(name))] = lows
		}
	}
	return out, sc.Err()
}
