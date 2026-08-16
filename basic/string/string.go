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

package str

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	. "beebread/basic"
	"beebread/basic/logic"
	beeMath "beebread/basic/math"
	. "beebread/basic/time_date"
)

// BIN_TO_BYTE converts a binary string into a byte.
func BIN_TO_BYTE(str string) byte {
	val, err := strconv.ParseUint(str, 2, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// BIN_TO_DWORD converts a binary string into a dword.
func BIN_TO_DWORD(str string) uint32 {
	val, err := strconv.ParseUint(str, 2, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// BYTE_TO_STRB converts a byte into a binary string.
func BYTE_TO_STRB(in byte) string {
	return fmt.Sprintf("%08b", in)
}

// BYTE_TO_STRH converts a byte into a hex string.
func BYTE_TO_STRH(in byte) string {
	return fmt.Sprintf("%02X", in)
}

// CAPITALIZE capitalizes the first letter of each word in a string.
func CAPITALIZE(str string) string {
	return strings.Title(strings.ToLower(str))
}

// CHARCODE returns the HTML character name for a given byte code.
// This is a simplified placeholder. A full implementation would require the charname data.
func CHARCODE(c byte) string {
	if c > 159 {
		// Placeholder for a complex lookup in beebread.Setup.Charnames
		return ""
	}
	return ""
}

// CHARNAME returns the byte code for a given HTML character name.
// This is a simplified placeholder. A full implementation would require the charname data.
func CHARNAME(str string) byte {
	// Placeholder for a complex lookup in beebread.Setup.Charnames
	return 0
}

// CHR_TO_STRING converts a byte into a string of length 1.
func CHR_TO_STRING(c byte) string {
	return string(c)
}

// CLEAN deletes all characters from a string except the ones specified in cx.
func CLEAN(in, cx string) string {
	var result strings.Builder
	result.Grow(len(in))
	for _, r := range in {
		if strings.ContainsRune(cx, r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// DEL_CHARS deletes all characters specified in cx from a string str.
func DEL_CHARS(str, cx string) string {
	// This is a more idiomatic way to implement the original DEL_CHARS
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(cx, r) {
			return -1
		}
		return r
	}, str)
}

// CODE returns the ASCII code of a character in a string at a given position (1-based).
func CODE(str string, pos int) byte {
	if pos > 0 && pos <= len(str) {
		return str[pos-1]
	}
	return 0
}

// COUNT_CHAR counts the number of characters c in a string str.
func COUNT_CHAR(str string, c byte) int {
	return strings.Count(str, string(c))
}

// COUNT_SUBSTRING counts the number of substrings sub in a string str.
func COUNT_SUBSTRING(str, sub string) int {
	return strings.Count(str, sub)
}

// DEC_TO_BYTE converts a decimal string into a byte.
func DEC_TO_BYTE(str string) byte {
	val, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// DEC_TO_DWORD converts a decimal string into a dword.
func DEC_TO_DWORD(str string) uint32 {
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// DEC_TO_INT converts a decimal string into an integer.
func DEC_TO_INT(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

// DT_TO_STRF converts a DT into a string with a given format.
func DT_TO_STRF(dtIn time.Time, format string) string {
	// This is a complex format string parser. A full implementation would be extensive.
	// Here's a simplified version using Go's standard time formatting as a base.
	// The OSCAT format codes don't map directly to Go's layout string.
	// Example: %Y-%m-%d %H:%M:%S
	replacer := strings.NewReplacer(
		"%d", fmt.Sprintf("%02d", dtIn.Day()),
		"%H", fmt.Sprintf("%02d", dtIn.Hour()),
		"%j", fmt.Sprintf("%03d", dtIn.YearDay()),
		"%m", fmt.Sprintf("%02d", dtIn.Month()),
		"%M", fmt.Sprintf("%02d", dtIn.Minute()),
		"%S", fmt.Sprintf("%02d", dtIn.Second()),
		"%w", strconv.Itoa(DAY_OF_WEEK(dtIn)-1), // OSCAT is 1-7, Go is 0-6 for some uses
		"%W", fmt.Sprintf("%02d", WORK_WEEK(dtIn)),
		"%y", dtIn.Format("06"),
		"%Y", dtIn.Format("2006"),
		"%%", "%",
	)
	return replacer.Replace(format)
}

// DWORD_TO_STRB converts a dword into a binary string.
func DWORD_TO_STRB(in uint32) string {
	return fmt.Sprintf("%032b", in)
}

// DWORD_TO_STRF converts a DWORD into a string with a given format.
func DWORD_TO_STRF(in uint32, format string) string {
	// This is a complex format string parser. A full implementation would be extensive.
	// Simplified version:
	replacer := strings.NewReplacer(
		"%b", DWORD_TO_STRB(in),
		"%d", strconv.FormatUint(uint64(in), 10),
		"%h", DWORD_TO_STRH(in),
		"%%", "%",
	)
	return replacer.Replace(format)
}

// DWORD_TO_STRH converts a dword into a hex string.
func DWORD_TO_STRH(in uint32) string {
	return fmt.Sprintf("%08X", in)
}

// EXEC executes a simple mathematical term.
// This is a placeholder for a very complex and unsafe function.
// A proper implementation would require a full expression parser.
func EXEC(str string) string {
	// Placeholder - a real implementation is a major task.
	return "ERROR"
}

// FILL creates a string of length L with character C.
func FILL(c byte, l int) string {
	if l <= 0 {
		return ""
	}
	if l > StringLength {
		l = StringLength
	}
	return strings.Repeat(string(c), l)
}

// FIND_CHAR finds the first character that is not a control character.
func FIND_CHAR(str string, pos int) int {
	if pos < 1 {
		pos = 1
	}
	for i := pos - 1; i < len(str); i++ {
		if !ISC_CTRL(str[i]) {
			return i + 1
		}
	}
	return 0
}

// FIND_CTRL finds the first control character in a string.
func FIND_CTRL(str string, pos int) int {
	if pos < 1 {
		pos = 1
	}
	for i := pos - 1; i < len(str); i++ {
		if ISC_CTRL(str[i]) {
			return i + 1
		}
	}
	return 0
}

// FIND_NONUM finds the first character that is not a number or a dot.
func FIND_NONUM(str string, pos int) int {
	if pos < 1 {
		pos = 1
	}
	for i := pos - 1; i < len(str); i++ {
		if !((str[i] >= '0' && str[i] <= '9') || str[i] == '.') {
			return i + 1
		}
	}
	return 0
}

// FIND_NUM finds the first character that is a number or a dot.
func FIND_NUM(str string, pos int) int {
	if pos < 1 {
		pos = 1
	}
	for i := pos - 1; i < len(str); i++ {
		if (str[i] >= '0' && str[i] <= '9') || str[i] == '.' {
			return i + 1
		}
	}
	return 0
}

// FINDB finds the last occurrence of str2 in str1.
func FINDB(str1, str2 string) int {
	pos := strings.LastIndex(str1, str2)
	if pos == -1 {
		return 0
	}
	return pos + 1 // 1-based index
}

// FINDB_NONUM finds the last character that is not a number or a dot.
func FINDB_NONUM(str string) int {
	for i := len(str) - 1; i >= 0; i-- {
		if !((str[i] >= '0' && str[i] <= '9') || str[i] == '.') {
			return i + 1
		}
	}
	return 0
}

// FINDB_NUM finds the last character that is a number or a dot.
func FINDB_NUM(str string) int {
	for i := len(str) - 1; i >= 0; i-- {
		if (str[i] >= '0' && str[i] <= '9') || str[i] == '.' {
			return i + 1
		}
	}
	return 0
}

// FINDP finds the first occurrence of src in str, starting from pos (1-based).
func FINDP(str, src string, pos int) int {
	if pos < 1 {
		pos = 1
	}
	if len(str) < pos-1 {
		return 0
	}
	foundPos := strings.Index(str[pos-1:], src)
	if foundPos == -1 {
		return 0
	}
	return foundPos + pos
}

// FIX adjusts a string to a fixed length L, padding or truncating as needed.
func FIX(str string, l int, c byte, m int) string {
	currentLen := len(str)
	if l <= currentLen {
		if m == 1 { // Right align (take from right)
			return str[currentLen-l:]
		}
		return str[:l] // Left align (take from left)
	}

	padding := strings.Repeat(string(c), l-currentLen)
	switch m {
	case 1: // Pad left
		return padding + str
	case 2: // Pad center
		padLen := l - currentLen
		leftPad := padLen / 2
		rightPad := padLen - leftPad
		return strings.Repeat(string(c), leftPad) + str + strings.Repeat(string(c), rightPad)
	default: // Pad right
		return str + padding
	}
}

// FLOAT_TO_REAL converts a string to a float64.
func FLOAT_TO_REAL(flt string) float64 {
	// A simplified version. The original is very complex and tries to parse manually.
	// Go's strconv is more robust.
	f, err := strconv.ParseFloat(strings.TrimSpace(flt), 64)
	if err != nil {
		return 0.0
	}
	return f
}

// FSTRING_TO_BYTE converts a formatted string (e.g., "16#FF", "2#1010") to a byte.
func FSTRING_TO_BYTE(in string) byte {
	if strings.HasPrefix(in, "16#") {
		return HEX_TO_BYTE(in[3:])
	} else if strings.HasPrefix(in, "8#") {
		return OCT_TO_BYTE(in[2:])
	} else if strings.HasPrefix(in, "2#") {
		return BIN_TO_BYTE(in[2:])
	}
	return DEC_TO_BYTE(in)
}

// FSTRING_TO_DWORD converts a formatted string to a dword.
func FSTRING_TO_DWORD(in string) uint32 {
	if strings.HasPrefix(in, "16#") {
		return HEX_TO_DWORD(in[3:])
	} else if strings.HasPrefix(in, "8#") {
		return OCT_TO_DWORD(in[2:])
	} else if strings.HasPrefix(in, "2#") {
		return BIN_TO_DWORD(in[2:])
	}
	return DEC_TO_DWORD(in)
}

// FstringToDt converts a formatted string into a DT (time.Time) value.
func FstringToDt(sdt, fmtStr string) time.Time {
	const (
		ignore = '*'
		fchar  = '#'
	)

	var (
		dy = 1970
		dm = 1
		dd = 1
		th = 0
		tm = 0
		ts = 0
	)

	sdtRunes := []rune(sdt)
	fmtRunes := []rune(fmtStr)

	sdtPos, fmtPos := 0, 0

	for fmtPos < len(fmtRunes) {
		if sdtPos >= len(sdtRunes) {
			break
		}

		switch fmtRunes[fmtPos] {
		case ignore:
			fmtPos++
			sdtPos++
		case fchar:
			fmtPos++ // Move past '#'
			if fmtPos >= len(fmtRunes) {
				break
			}
			formatCode := fmtRunes[fmtPos]
			fmtPos++

			end := sdtPos
			for end < len(sdtRunes) && (fmtPos >= len(fmtRunes) || sdtRunes[end] != fmtRunes[fmtPos]) {
				end++
			}

			val := string(sdtRunes[sdtPos:end])
			sdtPos = end

			switch formatCode {
			case 'Y':
				dy, _ = strconv.Atoi(val)
				if dy < 100 {
					dy += 2000
				}
			case 'M':
				dm, _ = strconv.Atoi(val)
			case 'N':
				dm = FSTRING_TO_MONTH(val, 0)
			case 'D':
				dd, _ = strconv.Atoi(val)
			case 'h':
				th, _ = strconv.Atoi(val)
			case 'm':
				tm, _ = strconv.Atoi(val)
			case 's':
				ts, _ = strconv.Atoi(val)
			}
		default:
			if fmtRunes[fmtPos] == sdtRunes[sdtPos] {
				fmtPos++
				sdtPos++
			} else {
				// Mismatch, stop parsing
				return time.Date(dy, time.Month(dm), dd, th, tm, ts, 0, time.UTC)
			}
		}
	}

	return time.Date(dy, time.Month(dm), dd, th, tm, ts, 0, time.UTC)
}

// FSTRING_TO_MONTH converts a month string (name or number) to an integer (1-12).
func FSTRING_TO_MONTH(mth string, lang int) int {
	// Placeholder for a complex lookup.
	if i, err := strconv.Atoi(mth); err == nil {
		return i
	}
	return 0
}

// FSTRING_TO_WEEK converts a comma-separated list of weekdays to a bitmask byte.
func FSTRING_TO_WEEK(week string, lang int) byte {
	// Placeholder for a complex lookup.
	return 0
}

// FSTRING_TO_WEEKDAY converts a weekday string to an integer (1-7).
func FSTRING_TO_WEEKDAY(wday string, lang int) int {
	// Placeholder for a complex lookup.
	if i, err := strconv.Atoi(wday); err == nil {
		return i
	}
	return 0
}

// HEX_TO_BYTE converts a hexadecimal string to a byte.
func HEX_TO_BYTE(hex string) byte {
	val, err := strconv.ParseUint(hex, 16, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// HEX_TO_DWORD converts a hexadecimal string to a dword.
func HEX_TO_DWORD(hex string) uint32 {
	val, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// IS_ALNUM checks if a string contains only alphanumeric characters.
func IS_ALNUM(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ISC_ALPHA(byte(r)) && !ISC_NUM(byte(r)) {
			return false
		}
	}
	return true
}

// IsAlpha checks if a string contains only alphabetic characters.
func IS_ALPHA(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ISC_ALPHA(byte(r)) {
			return false
		}
	}
	return true
}

// IsCc checks if a string contains only characters from the cmp string.
func IS_CC(str, cmp string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !strings.ContainsRune(cmp, r) {
			return false
		}
	}
	return true
}

func IS_CTRL(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ISC_CTRL(byte(r)) {
			return false
		}
	}
	return true
}

// IsHex checks if a string contains only hexadecimal characters.
func IS_HEX(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ISC_HEX(byte(r)) {
			return false
		}
	}
	return true
}

// IsLower checks if a string contains only lowercase characters.
func IS_LOWER(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ISC_LOWER(byte(r)) {
			return false
		}
	}
	return true
}

// IsNcc checks if a string contains no characters from the cmp string.
func IS_NCC(str, cmp string) bool {
	return !strings.ContainsAny(str, cmp)
}

// IsNum checks if a string contains only numeric characters.
func IS_NUM(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ISC_NUM(byte(r)) {
			return false
		}
	}
	return true
}

// IsUpper checks if a string contains only uppercase characters.
func IS_UPPER(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !ISC_UPPER(byte(r)) {
			return false
		}
	}
	return true
}

// IscAlpha checks if a character is a..z or A..Z.
func ISC_ALPHA(in byte) bool {
	return (in >= 'a' && in <= 'z') || (in >= 'A' && in <= 'Z')
}

// IscCtrl checks if a character is a control character.
func ISC_CTRL(in byte) bool {
	return in < 32 || in == 127
}

// IscHex checks if a character is 0..9, A..F, or a..f.
func ISC_HEX(in byte) bool {
	return (in >= '0' && in <= '9') || (in >= 'A' && in <= 'F') || (in >= 'a' && in <= 'f')
}

// IscLower checks if a character is lowercase.
func ISC_LOWER(in byte) bool {
	return in >= 'a' && in <= 'z'
}

// IscNum checks if a character is 0..9.
func ISC_NUM(in byte) bool {
	return in >= '0' && in <= '9'
}

// IscUpper checks if a character is uppercase.
func ISC_UPPER(in byte) bool {
	return in >= 'A' && in <= 'Z'
}

// Lowercase converts a string to lowercase.
func LOWERCASE(str string) string {
	return strings.ToLower(str)
}

// Message4R is a rotating message display.
type MESSAGE_4R struct {
	Mx string
	Mn int
	Tr bool

	// internal state
	timer logic.TON
	edge  bool
}

// Update executes the message rotation logic.
func (m *MESSAGE_4R) Update(m0, m1, m2, m3 string, mm int, enq, clk bool, t1 time.Duration) {
	m.Tr = false
	if enq {
		m.timer.Update(clk, t1)
		if (clk && !m.edge) || m.timer.Q {
			if mm > 0 {
				m.Mn = (m.Mn + 1) % (mm + 1)
			}
			m.Tr = true
			m.timer.IN = false // Reset timer
			switch m.Mn {
			case 0:
				m.Mx = m0
			case 1:
				m.Mx = m1
			case 2:
				m.Mx = m2
			case 3:
				m.Mx = m3
			}
		}
		m.edge = clk
	} else {
		m.Mx = ""
		m.Mn = 0
	}
}

// MESSAGE_8 selects one of 8 messages based on prioritized inputs.
func MESSAGE_8(in [8]bool, s [8]string) string {
	for i := 0; i < 8; i++ {
		if in[i] {
			return s[i]
		}
	}
	return ""
}

// MIRROR reverses an input string.
func MIRROR(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// MONTH_TO_STRING converts an integer (1-12) to a month name.
func MONTH_TO_STRING(mth, lang, lx int) string {
	// Placeholder for a complex lookup.
	if mth >= 1 && mth <= 12 {
		return time.Month(mth).String()
	}
	return ""
}

// OCT_TO_BYTE converts an octal string to a byte.
func OCT_TO_BYTE(oct string) byte {
	val, err := strconv.ParseUint(oct, 8, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// OCT_TO_DWORD converts an octal string to a dword.
func OCT_TO_DWORD(oct string) uint32 {
	val, err := strconv.ParseUint(oct, 8, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// REAL_TO_STRF converts a float to a string with N decimal places.
func REAL_TO_STRF(in float64, n int, d string) string {
	n = int(beeMath.LIMIT(0, float64(n), 7))

	// Scale and round
	multiplier := beeMath.EXP10(float64(n)) // ST: O := ABS(in) * EXP10(N);
	val := beeMath.D_TRUNC(math.Abs(in)*multiplier + 0.5)

	res := strconv.FormatInt(val, 10)

	// Pad with leading zeros if necessary
	for len(res) <= n {
		res = "0" + res
	}

	// Insert decimal separator
	if n > 0 {
		res = res[:len(res)-n] + d + res[len(res)-n:]
	}

	// Add sign for negative numbers
	if in < 0.0 {
		res = "-" + res
	}
	return res
}

// REPLACE_ALL replaces all occurrences of src in str with rep.
func REPLACE_ALL(str, src, rep string) string {
	return strings.ReplaceAll(str, src, rep)
}

// REPLACE_CHARS replaces characters in str based on a mapping from src to rep.
func REPLACE_CHARS(str, src, rep string) string {
	if len(src) == 0 || len(rep) == 0 {
		return str
	}
	// This is a simplified interpretation. The original ST code is complex.
	// A more robust Go version would use a map or a replacer.
	minLen := len(src)
	if len(rep) < minLen {
		minLen = len(rep)
	}
	for i := 0; i < minLen; i++ {
		str = strings.ReplaceAll(str, string(src[i]), string(rep[i]))
	}
	return str
}

// REPLACE_UML replaces German umlauts with their two-letter equivalents.
func REPLACE_UML(str string) string {
	r := strings.NewReplacer(
		"Ä", "Ae", "Ö", "Oe", "Ü", "Ue", "ß", "ss",
		"ä", "ae", "ö", "oe", "ü", "ue",
	)
	return r.Replace(str)
}

// Ticker creates a scrolling text effect.
type TICKER struct {
	Display string
	// internal state
	delay logic.TP
	step  int
}

// Update executes the ticker logic.
func (t *TICKER) Update(text string, n int, pt time.Duration) {
	if n <= 0 || n >= len(text) {
		t.Display = text
		return
	}

	// The TP timer will output Q=true for the duration of PT after a rising edge.
	// We can simulate a one-shot trigger for it.
	t.delay.Update(!t.delay.Q, pt)
	if !t.delay.Q { // When the timer is done (or on the first run)
		t.step++
		if t.step >= len(text) {
			t.step = 0
		}
	}

	// Create a circular view of the text
	circularText := text + text
	if t.step+n > len(circularText) {
		t.Display = circularText[t.step:]
	} else {
		t.Display = circularText[t.step : t.step+n]
	}
}

// TO_LOWER converts a character from uppercase to lowercase.
func TO_LOWER(in byte) byte {
	if in >= 'A' && in <= 'Z' {
		return in + ('a' - 'A')
	}
	// Placeholder for extended ASCII
	return in
}

// TO_UML converts a character to its two-letter Umlaut representation.
func TO_UML(in byte) string {
	switch in {
	case 196:
		return "Ae" // Ä
	case 214:
		return "Oe" // Ö
	case 220:
		return "Ue" // Ü
	case 223:
		return "ss" // ß
	case 228:
		return "ae" // ä
	case 246:
		return "oe" // ö
	case 252:
		return "ue" // ü
	default:
		return string(in)
	}
}

// TO_UPPER converts a character from lowercase to uppercase.
func TO_UPPER(in byte) byte {
	if in >= 'a' && in <= 'z' {
		return in - ('a' - 'A')
	}
	// Placeholder for extended ASCII
	return in
}

// TRIM removes all space characters from a string.
func TRIM(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

// TRIM1 replaces multiple spaces with a single space and trims leading/trailing spaces.
func TRIM1(str string) string {
	// Use Fields to split by whitespace and Join to put it back with single spaces.
	return strings.Join(strings.Fields(str), " ")
}

// TRIME removes leading and trailing space characters from a string.
func TRIME(str string) string {
	return strings.TrimSpace(str)
}

// UPPERCASE converts a string to uppercase.
func UPPERCASE(str string) string {
	return strings.ToUpper(str)
}

// WEEKDAY_TO_STRING converts an integer (1-7) to a weekday name.
func WEEKDAY_TO_STRING(wday, lang, lx int) string {
	// Placeholder for a complex lookup.
	if wday < 1 || wday > 7 {
		return ""
	}
	if lang <= 0 {
		lang = int(Language.Default)
	}
	return Language.Weekdays[lang-1][wday-1]
}

// INT_TO_STRF converts an integer to a string of a fixed length N.
func INT_TO_STRF(in, n int) string {
	return FIX(strconv.Itoa(in), n, '0', 1)
}
