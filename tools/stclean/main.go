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

// Command stclean cleans an OSCAT library exported by CODESYS into the
// IEC 61131-3 source beedance reads, as beedance's
// reference/beedance_oscat_basic.st was cleaned:
//
//   - the CODESYS metadata, (* @... *), and the export's trailer (alarm
//     configuration, library list, resources) are left out;
//   - the version and description comment ahead of each body and the
//     revision history after it are left out, in commented out POUs too,
//     and so are the prose comment blocks of several lines in a body (tables
//     of units, explanations); comments on a line of code stay;
//   - a POU that uses POINTER TO or ADR, which IEC 61131-3 does not have, is
//     commented out, and so is a POU that uses one commented out, here or in
//     a library given with -uses.
//
// With -refs, pointers are written as beedance's references instead:
// POINTER TO as REF_TO and ADR( as REF(, which beedance also accepts. A
// cleaned source can be cleaned again with -refs: the POUs it commented out
// are restored, and only those -exclude names (with why), and the POUs that
// use them, are commented out, as beedance references refer to a whole
// variable of their type, with no arithmetic or byte access.
//
// The source is a CODESYS export, or a CODESYS 2.3 library (.lib), whose
// functions, function blocks, types and global variables are read from it;
// see lib.go.
//
// Usage:
//
//	go run ./tools/stclean -in building/oscat_building_100.st \
//		-uses ../beedance/reference/beedance_oscat_basic.st \
//		-out building/beedance_building_100.st
//	go run ./tools/stclean -in network/beckoff_network_135.st \
//		-uses ../beedance/reference/beedance_oscat_basic.st,building/beedance_building_100.st \
//		-out network/beedance_network_135.st
//	go run ./tools/stclean -refs -exclude tools/stclean/testdata/network.exclude \
//		-uses doc/beedance_basic.st -in doc/beedance_network.st -out doc/beedance_network.st
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	reMeta     = regexp.MustCompile(`^\s*\(\* @[A-Z_]+ := .*\*\)\s*$`)
	reMetaOpen = regexp.MustCompile(`^\s*\(\* @[A-Z_]+ :=`)
	reEndDecl  = regexp.MustCompile(`^\s*\(\* @END_DECLARATION`)
	// OSCAT misspells a few FUNCTION_BLOCKs as FUNCTIONBLOCK.
	rePOU     = regexp.MustCompile(`^\s*(FUNCTION_BLOCK|FUNCTIONBLOCK|FUNCTION|PROGRAM|TYPE)\s+([A-Za-z_0-9]+)`)
	reEndPOU  = regexp.MustCompile(`^\s*END_(FUNCTION_BLOCK|FUNCTIONBLOCK|FUNCTION|PROGRAM|TYPE)\b`)
	rePointer = regexp.MustCompile(`(?i)\bPOINTER\s+TO\b|\bADR\s*\(`)
	reWord    = regexp.MustCompile(`[A-Za-z_][A-Za-z_0-9]*`)
	// A declaration of variables, and a call or the type of a declaration.
	reDecl     = regexp.MustCompile(`(?m)^\s*([A-Za-z_]\w*(?:\s*,\s*[A-Za-z_]\w*)*)\s*:[^=]`)
	reUse      = regexp.MustCompile(`([A-Za-z_]\w*)\s*\(|:\s*([A-Za-z_]\w*)`)
	reComment  = regexp.MustCompile(`^// FUNCTION(?:_BLOCK)?\s+([A-Za-z_0-9]+)`)
	reTrailer  = regexp.MustCompile(`^_ALARMCONFIG\s*$`)
	reRevision = regexp.MustCompile(`(?i)revision\s+hist|\brev\.?\s*\d+\.\d+`)
)

func main() {
	in := flag.String("in", "", "the OSCAT source exported by CODESYS")
	out := flag.String("out", "", "the cleaned source to write")
	uses := flag.String("uses", "", "a cleaned library, comma separated, whose commented out POUs the source may use")
	flag.BoolVar(&refs, "refs", false, "write POINTER TO and ADR( as beedance's references, REF_TO and REF(, and restore the POUs a cleaned source commented out")
	exclude := flag.String("exclude", "", "with -refs, a file of POUs to comment out, one per line: the name, then why")
	flag.Parse()
	if *in == "" || *out == "" {
		flag.Usage()
		os.Exit(2)
	}
	lines, err := readSource(*in)
	if err != nil {
		fail(err)
	}
	if refs {
		lines = restoreCommented(lines)
	}
	if *exclude != "" {
		if excluded, err = readExclude(*exclude); err != nil {
			fail(err)
		}
	}
	commented := map[string]string{}
	for _, path := range strings.Split(*uses, ",") {
		if path == "" {
			continue
		}
		used, err := readLines(path)
		if err != nil {
			fail(err)
		}
		for _, l := range used {
			if m := reComment.FindStringSubmatch(l); m != nil {
				commented[strings.ToUpper(m[1])] = m[1]
			}
		}
	}
	text := clean(lines, commented)
	if err := os.WriteFile(*out, []byte(strings.Join(text, "\r\n")+"\r\n"), 0o644); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "stclean:", err)
	os.Exit(1)
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		lines = append(lines, strings.TrimRight(sc.Text(), "\r"))
	}
	return lines, sc.Err()
}

// depth returns the comment depth after line, starting at d.
func depth(d int, line string) int {
	d += strings.Count(line, "(*") - strings.Count(line, "*)")
	if d < 0 {
		return 0
	}
	return d
}

// lastComment returns where the comment block that ends the lines, past
// blank lines, starts, if one does and it starts a line.
func lastComment(lines []string) (int, bool) {
	end := len(lines)
	for end > 0 && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	if end == 0 || !strings.HasSuffix(strings.TrimSpace(lines[end-1]), "*)") {
		return 0, false
	}
	d := 0
	for start := end - 1; start >= 0; start-- {
		d += strings.Count(lines[start], "*)") - strings.Count(lines[start], "(*")
		if d > 0 {
			continue
		}
		if d < 0 || !strings.HasPrefix(strings.TrimSpace(lines[start]), "(*") {
			return 0, false
		}
		return start, true
	}
	return 0, false
}

// A pou is a POU's lines, its name and whether it is commented out, and why.
type pou struct {
	name   string
	lines  []string
	reason string
}

func clean(lines []string, commented map[string]string) []string {
	// Leave out the metadata and the trailer, and the comments around each
	// body.
	var out []string
	inMeta := false
	for _, line := range lines {
		if reTrailer.MatchString(line) {
			break
		}
		switch {
		case inMeta:
			inMeta = !strings.Contains(line, "*)")
			continue
		case reEndDecl.MatchString(line):
			// The version and description of the POU.
			for start, ok := lastComment(out); ok; start, ok = lastComment(out) {
				out = out[:start]
			}
			continue
		case reMeta.MatchString(line):
			continue
		case reMetaOpen.MatchString(line):
			inMeta = true
			continue
		case reEndPOU.MatchString(line):
			// The revision history, which other comments may follow.
			var kept [][]string
			for start, ok := lastComment(out); ok; start, ok = lastComment(out) {
				if !reRevision.MatchString(strings.Join(out[start:], "\n")) {
					kept = append([][]string{append([]string(nil), out[start:]...)}, kept...)
				}
				out = out[:start]
			}
			for _, block := range kept {
				out = append(out, block...)
			}
			for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
				out = out[:len(out)-1]
			}
			line = strings.Replace(strings.TrimSpace(line), "END_FUNCTIONBLOCK", "END_FUNCTION_BLOCK", 1)
		case rePOU.MatchString(line):
			line = strings.Replace(strings.TrimLeft(line, " \t"), "FUNCTIONBLOCK", "FUNCTION_BLOCK", 1)
		}
		// One blank line at a time.
		if strings.TrimSpace(line) == "" && (len(out) == 0 || strings.TrimSpace(out[len(out)-1]) == "") {
			continue
		}
		out = append(out, strings.TrimRight(line, " \t"))
	}

	// Split the POUs from what is between them, the globals.
	var parts []*pou
	var cur *pou
	d := 0
	for _, line := range out {
		if cur == nil && d == 0 {
			if m := rePOU.FindStringSubmatch(line); m != nil {
				cur = &pou{name: m[2]}
				parts = append(parts, cur)
			}
		}
		if cur != nil {
			cur.lines = append(cur.lines, line)
		} else {
			parts = append(parts, &pou{lines: []string{line}})
		}
		d = depth(d, line)
		if cur != nil && d == 0 && reEndPOU.MatchString(line) {
			cur = nil
		}
	}

	// Comment out the POUs that use pointers (or, with -refs, write them as
	// references and comment out those -exclude names), then those that use
	// them.
	for _, p := range parts {
		switch {
		case p.name == "":
		case refs:
			p.lines = mapCode(p.lines, func(s string) string {
				return reAdr.ReplaceAllString(rePointerTo.ReplaceAllString(s, "REF_TO"), "REF(")
			})
			if why, ok := excluded[strings.ToUpper(p.name)]; ok {
				p.reason = why
				commented[strings.ToUpper(p.name)] = p.name
			}
		case rePointer.MatchString(code(p.lines)):
			p.reason = "it uses POINTER TO or ADR, which IEC 61131-3 does not have."
			commented[strings.ToUpper(p.name)] = p.name
		}
	}
	for changed := true; changed; {
		changed = false
		for _, p := range parts {
			if p.name == "" || p.reason != "" {
				continue
			}
			body := code(p.lines[1:])
			locals := map[string]bool{}
			for _, m := range reDecl.FindAllStringSubmatch(body, -1) {
				for _, w := range reWord.FindAllString(m[1], -1) {
					locals[strings.ToUpper(w)] = true
				}
			}
			for _, m := range reUse.FindAllStringSubmatch(body, -1) {
				w := m[1] + m[2]
				if locals[strings.ToUpper(w)] {
					continue
				}
				if name, ok := commented[strings.ToUpper(w)]; ok {
					p.reason = fmt.Sprintf("it uses %s, which is commented out.", name)
					commented[strings.ToUpper(p.name)] = p.name
					changed = true
					break
				}
			}
		}
	}

	var text []string
	for _, p := range parts {
		if p.reason == "" {
			text = append(text, rewrite(p.lines)...)
			continue
		}
		text = append(text, "// Commented out for beedance: "+p.reason)
		for _, l := range p.lines {
			if l == "" {
				text = append(text, "//")
			} else {
				text = append(text, "// "+l)
			}
		}
	}
	text = dropCommentedProse(text)
	for len(text) > 0 && strings.TrimSpace(text[len(text)-1]) == "" {
		text = text[:len(text)-1]
	}
	for len(text) > 0 && strings.TrimSpace(text[0]) == "" {
		text = text[1:]
	}
	return text
}

var (
	// refs is -refs: pointers are written as references (see restoreCommented).
	refs bool
	// excluded are the POUs -exclude comments out, by upper case name, and why.
	excluded = map[string]string{}

	rePointerTo = regexp.MustCompile(`(?i)\bPOINTER\s+TO\b`)
	reAdr       = regexp.MustCompile(`(?i)\bADR\s*\(`)
	reNote      = regexp.MustCompile(`^// (Commented out|Changed) for beedance: `)
)

// readExclude reads an -exclude file: on each line a POU's name, then why it
// is commented out; a line starting with # is a comment.
func readExclude(path string) (map[string]string, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}
	ex := map[string]string{}
	for _, l := range lines {
		if name, why, ok := strings.Cut(strings.TrimSpace(l), " "); ok && !strings.HasPrefix(name, "#") {
			ex[strings.ToUpper(name)] = strings.TrimSpace(why)
		}
	}
	return ex, nil
}

// restoreCommented undoes what a previous run commented out, so that -refs
// can decide again with beedance's references: each POU after a "Commented
// out for beedance" note is uncommented, and the notes, which the run writes
// again, are left out.
func restoreCommented(lines []string) []string {
	var out []string
	for i := 0; i < len(lines); i++ {
		if reNote.MatchString(lines[i]) {
			if !strings.HasPrefix(lines[i], "// Commented out") {
				continue // a "Changed for beedance" note
			}
			// The POU's // lines, up to its END_ line.
			for i++; i < len(lines) && strings.HasPrefix(lines[i], "//"); i++ {
				line := strings.TrimPrefix(strings.TrimPrefix(lines[i], "//"), " ")
				out = append(out, line)
				if reEndPOU.MatchString(line) {
					break
				}
			}
			continue
		}
		out = append(out, lines[i])
	}
	return out
}

// reVersionNote matches a POU's version and author note.
var reVersionNote = regexp.MustCompile(`(?i)\bprogrammer\b|\btested\s+by\b|^\(\*\s*version\s+\d`)

// dropCommentedProse leaves out the comments that are prose rather than
// code. In the lines of POUs commented out with //, these are the revision
// history and the version and author note, as clean leaves them out of the
// other POUs. In the code, these are the comment blocks of several lines
// that stand on their own lines (tables of units, explanations of a
// formula); a comment on a line of code, or of one line, stays. It keeps the
// note saying why a POU is commented out; runs of empty // lines become one.
func dropCommentedProse(text []string) []string {
	var out []string
	for i := 0; i < len(text); i++ {
		if t := strings.TrimSpace(text[i]); strings.HasPrefix(t, "(*") && depth(0, t) > 0 {
			// A comment block starting a line of code: up to the line that
			// closes it, if nothing but the comment is on that line.
			end, d := -1, 0
			for j := i; j < len(text) && !strings.HasPrefix(text[j], "//"); j++ {
				if d = depth(d, text[j]); d == 0 {
					if strings.HasSuffix(strings.TrimSpace(text[j]), "*)") {
						end = j
					}
					break
				}
			}
			if end > i {
				i = end
				continue
			}
		}
		body, isComment := strings.CutPrefix(text[i], "//")
		body = strings.TrimSpace(body)
		if isComment && strings.HasPrefix(body, "(*") {
			// The (* ... *) block that starts here, within the // lines.
			end, d := i, 0
			for j := i; j < len(text); j++ {
				b, ok := strings.CutPrefix(text[j], "//")
				if !ok {
					break
				}
				d = depth(d, b)
				end = j
				if d == 0 {
					break
				}
			}
			var block []string
			for _, l := range text[i : end+1] {
				block = append(block, strings.TrimSpace(strings.TrimPrefix(l, "//")))
			}
			joined := strings.Join(block, "\n")
			if reRevision.MatchString(joined) || reVersionNote.MatchString(joined) {
				i = end
				continue
			}
		}
		if text[i] == "//" && len(out) > 0 && out[len(out)-1] == "//" {
			continue
		}
		out = append(out, text[i])
	}
	return out
}

// rewrite rewrites what CODESYS accepts and IEC 61131-3 does not, with a
// note ahead of the POU: a name that is a keyword gets an underscore, and a
// chained assignment, a := b := c, is two.
func rewrite(lines []string) []string {
	var notes []string
	renamed := map[string]bool{}
	chained, literals := false, false
	lines = mapCode(lines, func(s string) string {
		s = reKeyword.ReplaceAllStringFunc(s, func(w string) string {
			renamed[strings.ToUpper(w)] = true
			return w + "_"
		})
		s = reDateTime.ReplaceAllStringFunc(s, func(lit string) string {
			full := fullLiteral(lit)
			if full != lit {
				literals = true
			}
			return full
		})
		if reChained.MatchString(s) {
			chained = true
			s = reChained.ReplaceAllString(s, "$1$3 := $4; $2 := $3;")
		}
		return s
	})
	for _, w := range keywords {
		if renamed[w] {
			notes = append(notes, fmt.Sprintf("// Changed for beedance: %s is a keyword of IEC 61131-3, so it is %s_ here.", w, w))
		}
	}
	if literals {
		notes = append(notes, "// Changed for beedance: a date or time of day literal is written in full, as D#yyyy-mm-dd and TOD#hh:mm:ss.")
	}
	if chained {
		notes = append(notes, "// Changed for beedance: a chained assignment, a := b := c, is b := c; a := b.")
	}
	if moved := moveComments(&lines); moved {
		notes = append(notes, "// Changed for beedance: a comment before THEN, or between CASE ... OF and the first case, is moved out, where beedance's parser takes it.")
	}
	return append(notes, lines...)
}

// A comment before THEN, and comment lines after CASE ... OF.
var (
	reThenComment = regexp.MustCompile(`\s*(\(\*(?:[^*]|\*[^)])*\*\))\s*THEN\b`)
	reCaseOf      = regexp.MustCompile(`^\s*CASE\b.*\bOF\s*$`)
	reCommentLine = regexp.MustCompile(`^\s*\(\*.*\*\)\s*$`)
)

// moveComments moves a comment before THEN after it, and the comment lines
// between CASE ... OF and the first case before the CASE, and reports
// whether it moved any.
func moveComments(lines *[]string) bool {
	moved := false
	ls := *lines
	for i, l := range ls {
		if reThenComment.MatchString(l) {
			ls[i] = reThenComment.ReplaceAllString(l, " THEN $1")
			moved = true
		}
	}
	for i := 0; i < len(ls); i++ {
		if !reCaseOf.MatchString(ls[i]) {
			continue
		}
		j := i + 1
		for j < len(ls) && reCommentLine.MatchString(ls[j]) {
			j++
		}
		if j == i+1 {
			continue
		}
		// CASE, comments -> comments, CASE.
		caseLine := ls[i]
		copy(ls[i:j-1], ls[i+1:j])
		ls[j-1] = caseLine
		moved = true
		i = j - 1
	}
	return moved
}

// keywords are the IEC 61131-3 keywords OSCAT uses as names.
var (
	keywords   = []string{"SINGLE"}
	reKeyword  = regexp.MustCompile(`(?i)\b(` + strings.Join(keywords, "|") + `)\b`)
	reDateTime = regexp.MustCompile(`(?i)\b(?:DT|DATE_AND_TIME|D|DATE|TOD|TIME_OF_DAY)#[0-9:.-]+`)
	reChained  = regexp.MustCompile(`^(\s*)([A-Za-z_]\w*)\s*:=\s*([A-Za-z_]\w*)\s*:=\s*([^;]+);`)
)

// mapCode applies f to the code of the lines, outside comments and string
// literals.
func mapCode(lines []string, f func(string) string) []string {
	out := make([]string, len(lines))
	d := 0
	for n, line := range lines {
		var b, run strings.Builder
		flush := func() {
			b.WriteString(f(run.String()))
			run.Reset()
		}
		for i := 0; i < len(line); i++ {
			switch {
			case strings.HasPrefix(line[i:], "(*"):
				if d == 0 {
					flush()
				}
				d++
				b.WriteString("(*")
				i++
			case strings.HasPrefix(line[i:], "*)") && d > 0:
				d--
				b.WriteString("*)")
				i++
			case d > 0:
				b.WriteByte(line[i])
			case strings.HasPrefix(line[i:], "//"):
				flush()
				b.WriteString(line[i:])
				i = len(line)
			case line[i] == '\'':
				flush()
				j := i + 1
				for j < len(line) && line[j] != '\'' {
					if line[j] == '$' {
						j++
					}
					j++
				}
				if j >= len(line) {
					j = len(line) - 1
				}
				b.WriteString(line[i : j+1])
				i = j
			default:
				run.WriteByte(line[i])
			}
		}
		flush()
		out[n] = b.String()
	}
	return out
}

// code returns the lines without their comments and string literals.
func code(lines []string) string {
	var b strings.Builder
	d := 0
	for _, line := range lines {
		for i := 0; i < len(line); i++ {
			switch {
			case strings.HasPrefix(line[i:], "(*"):
				d++
				i++
			case strings.HasPrefix(line[i:], "*)") && d > 0:
				d--
				i++
			case d > 0:
			case strings.HasPrefix(line[i:], "//"):
				i = len(line)
			case line[i] == '\'':
				j := i + 1
				for j < len(line) && line[j] != '\'' {
					if line[j] == '$' {
						j++
					}
					j++
				}
				i = j
				b.WriteString("''")
			default:
				b.WriteByte(line[i])
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// fullLiteral writes a date, date and time or time of day literal in full,
// with two digit months, days, hours and minutes, and seconds.
func fullLiteral(lit string) string {
	hash := strings.Index(lit, "#")
	prefix, value := lit[:hash+1], lit[hash+1:]
	pad := func(parts []string) []string {
		for i, p := range parts {
			if len(p) == 1 {
				parts[i] = "0" + p
			}
		}
		return parts
	}
	tod := func(v string) string {
		parts := strings.Split(v, ":")
		if len(parts) == 2 {
			parts = append(parts, "00")
		}
		if len(parts) != 3 {
			return v
		}
		return strings.Join(pad(parts), ":")
	}
	parts := strings.Split(value, "-")
	switch strings.ToUpper(strings.TrimSuffix(prefix, "#")) {
	case "TOD", "TIME_OF_DAY":
		return prefix + tod(value)
	case "D", "DATE":
		if len(parts) != 3 {
			return lit
		}
		return prefix + strings.Join(pad(parts), "-")
	default:
		if len(parts) != 4 {
			return lit
		}
		return prefix + strings.Join(pad(parts[:3]), "-") + "-" + tod(parts[3])
	}
}

// readSource returns the lines of the export or library at path.
func readSource(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if isLib(data) {
		lines, programs := libExport(data)
		if len(programs) > 0 {
			fmt.Fprintf(os.Stderr, "stclean: left out the programs %s\n", strings.Join(programs, ", "))
		}
		return lines, nil
	}
	return readLines(path)
}
