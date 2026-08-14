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
	"strconv"
	"strings"
	"time"

	. "beebread/basic"
	. "beebread/basic/logic"
	. "beebread/basic/time_date"
)

// BinToByte converts a binary string into a byte.
func BinToByte(str string) byte {
	val, err := strconv.ParseUint(str, 2, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// BinToDword converts a binary string into a dword.
func BinToDword(str string) uint32 {
	val, err := strconv.ParseUint(str, 2, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// ByteToStrB converts a byte into a binary string.
func ByteToStrB(in byte) string {
	return fmt.Sprintf("%08b", in)
}

// ByteToStrH converts a byte into a hex string.
func ByteToStrH(in byte) string {
	return fmt.Sprintf("%02X", in)
}

// Capitalize capitalizes the first letter of each word in a string.
func Capitalize(str string) string {
	return strings.Title(strings.ToLower(str))
}

// CharCode returns the HTML character name for a given byte code.
// This is a simplified placeholder. A full implementation would require the charname data.
func CharCode(c byte) string {
	if c > 159 {
		// Placeholder for a complex lookup in beebread.Setup.Charnames
		return ""
	}
	return ""
}

// CharName returns the byte code for a given HTML character name.
// This is a simplified placeholder. A full implementation would require the charname data.
func CharName(str string) byte {
	// Placeholder for a complex lookup in beebread.Setup.Charnames
	return 0
}

// ChrToString converts a byte into a string of length 1.
func ChrToString(c byte) string {
	return string(c)
}

// Clean deletes all characters from a string except the ones specified in cx.
func Clean(in, cx string) string {
	var result strings.Builder
	result.Grow(len(in))
	for _, r := range in {
		if strings.ContainsRune(cx, r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// DelChars deletes all characters specified in cx from a string str.
func DelChars(str, cx string) string {
	// This is a more idiomatic way to implement the original DEL_CHARS
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(cx, r) {
			return -1
		}
		return r
	}, str)
}

// Code returns the ASCII code of a character in a string at a given position (1-based).
func Code(str string, pos int) byte {
	if pos > 0 && pos <= len(str) {
		return str[pos-1]
	}
	return 0
}

// CountChar counts the number of characters c in a string str.
func CountChar(str string, c byte) int {
	return strings.Count(str, string(c))
}

// CountSubstring counts the number of substrings sub in a string str.
func CountSubstring(str, sub string) int {
	return strings.Count(str, sub)
}

// DecToByte converts a decimal string into a byte.
func DecToByte(str string) byte {
	val, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// DecToDword converts a decimal string into a dword.
func DecToDword(str string) uint32 {
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// DecToInt converts a decimal string into an integer.
func DecToInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

// DtToStrF converts a DT into a string with a given format.
func DtToStrF(dtIn time.Time, format string) string {
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
		"%w", strconv.Itoa(DayOfWeek(dtIn)-1), // OSCAT is 1-7, Go is 0-6 for some uses
		"%W", fmt.Sprintf("%02d", WorkWeek(dtIn)),
		"%y", dtIn.Format("06"),
		"%Y", dtIn.Format("2006"),
		"%%", "%",
	)
	return replacer.Replace(format)
}

// DwordToStrB converts a dword into a binary string.
func DwordToStrB(in uint32) string {
	return fmt.Sprintf("%032b", in)
}

// DwordToStrF converts a DWORD into a string with a given format.
func DwordToStrF(in uint32, format string) string {
	// This is a complex format string parser. A full implementation would be extensive.
	// Simplified version:
	replacer := strings.NewReplacer(
		"%b", DwordToStrB(in),
		"%d", strconv.FormatUint(uint64(in), 10),
		"%h", DwordToStrH(in),
		"%%", "%",
	)
	return replacer.Replace(format)
}

// DwordToStrH converts a dword into a hex string.
func DwordToStrH(in uint32) string {
	return fmt.Sprintf("%08X", in)
}

// Exec executes a simple mathematical term.
// This is a placeholder for a very complex and unsafe function.
// A proper implementation would require a full expression parser.
func Exec(str string) string {
	// Placeholder - a real implementation is a major task.
	return "ERROR"
}

// Fill creates a string of length L with character C.
func Fill(c byte, l int) string {
	if l <= 0 {
		return ""
	}
	if l > StringLength {
		l = StringLength
	}
	return strings.Repeat(string(c), l)
}

// FindChar finds the first character that is not a control character.
func FindChar(str string, pos int) int {
	if pos < 1 {
		pos = 1
	}
	for i := pos - 1; i < len(str); i++ {
		if !IscCtrl(str[i]) {
			return i + 1
		}
	}
	return 0
}

// FindCtrl finds the first control character in a string.
func FindCtrl(str string, pos int) int {
	if pos < 1 {
		pos = 1
	}
	for i := pos - 1; i < len(str); i++ {
		if IscCtrl(str[i]) {
			return i + 1
		}
	}
	return 0
}

// FindNonum finds the first character that is not a number or a dot.
func FindNonum(str string, pos int) int {
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

// FindNum finds the first character that is a number or a dot.
func FindNum(str string, pos int) int {
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

// FindB finds the last occurrence of str2 in str1.
func FindB(str1, str2 string) int {
	pos := strings.LastIndex(str1, str2)
	if pos == -1 {
		return 0
	}
	return pos + 1 // 1-based index
}

// FindbNonum finds the last character that is not a number or a dot.
func FindbNonum(str string) int {
	for i := len(str) - 1; i >= 0; i-- {
		if !((str[i] >= '0' && str[i] <= '9') || str[i] == '.') {
			return i + 1
		}
	}
	return 0
}

// FindbNum finds the last character that is a number or a dot.
func FindbNum(str string) int {
	for i := len(str) - 1; i >= 0; i-- {
		if (str[i] >= '0' && str[i] <= '9') || str[i] == '.' {
			return i + 1
		}
	}
	return 0
}

// FindP finds the first occurrence of src in str, starting from pos (1-based).
func FindP(str, src string, pos int) int {
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

// Fix adjusts a string to a fixed length L, padding or truncating as needed.
func Fix(str string, l int, c byte, m int) string {
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

// FloatToReal converts a string to a float64.
func FloatToReal(flt string) float64 {
	// A simplified version. The original is very complex and tries to parse manually.
	// Go's strconv is more robust.
	f, err := strconv.ParseFloat(strings.TrimSpace(flt), 64)
	if err != nil {
		return 0.0
	}
	return f
}

// FstringToByte converts a formatted string (e.g., "16#FF", "2#1010") to a byte.
func FstringToByte(in string) byte {
	if strings.HasPrefix(in, "16#") {
		return HexToByte(in[3:])
	} else if strings.HasPrefix(in, "8#") {
		return OctToByte(in[2:])
	} else if strings.HasPrefix(in, "2#") {
		return BinToByte(in[2:])
	}
	return DecToByte(in)
}

// FstringToDword converts a formatted string to a dword.
func FstringToDword(in string) uint32 {
	if strings.HasPrefix(in, "16#") {
		return HexToDword(in[3:])
	} else if strings.HasPrefix(in, "8#") {
		return OctToDword(in[2:])
	} else if strings.HasPrefix(in, "2#") {
		return BinToDword(in[2:])
	}
	return DecToDword(in)
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
				dm = FstringToMonth(val, 0)
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

// FstringToMonth converts a month string (name or number) to an integer (1-12).
func FstringToMonth(mth string, lang int) int {
	// Placeholder for a complex lookup.
	if i, err := strconv.Atoi(mth); err == nil {
		return i
	}
	return 0
}

// FstringToWeek converts a comma-separated list of weekdays to a bitmask byte.
func FstringToWeek(week string, lang int) byte {
	// Placeholder for a complex lookup.
	return 0
}

// FstringToWeekday converts a weekday string to an integer (1-7).
func FstringToWeekday(wday string, lang int) int {
	// Placeholder for a complex lookup.
	if i, err := strconv.Atoi(wday); err == nil {
		return i
	}
	return 0
}

// HexToByte converts a hexadecimal string to a byte.
func HexToByte(hex string) byte {
	val, err := strconv.ParseUint(hex, 16, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// HexToDword converts a hexadecimal string to a dword.
func HexToDword(hex string) uint32 {
	val, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// IsAlnum checks if a string contains only alphanumeric characters.
func IsAlnum(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !IscAlpha(byte(r)) && !IscNum(byte(r)) {
			return false
		}
	}
	return true
}

// IsAlpha checks if a string contains only alphabetic characters.
func IsAlpha(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !IscAlpha(byte(r)) {
			return false
		}
	}
	return true
}

// IsCc checks if a string contains only characters from the cmp string.
func IsCc(str, cmp string) bool {
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

// IsCtrl checks if a string contains only control characters.
func IsCtrl(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !IscCtrl(byte(r)) {
			return false
		}
	}
	return true
}

// IsHex checks if a string contains only hexadecimal characters.
func IsHex(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !IscHex(byte(r)) {
			return false
		}
	}
	return true
}

// IsLower checks if a string contains only lowercase characters.
func IsLower(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !IscLower(byte(r)) {
			return false
		}
	}
	return true
}

// IsNcc checks if a string contains no characters from the cmp string.
func IsNcc(str, cmp string) bool {
	return !strings.ContainsAny(str, cmp)
}

// IsNum checks if a string contains only numeric characters.
func IsNum(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !IscNum(byte(r)) {
			return false
		}
	}
	return true
}

// IsUpper checks if a string contains only uppercase characters.
func IsUpper(str string) bool {
	if len(str) == 0 {
		return false
	}
	for _, r := range str {
		if !IscUpper(byte(r)) {
			return false
		}
	}
	return true
}

// IscAlpha checks if a character is a..z or A..Z.
func IscAlpha(in byte) bool {
	return (in >= 'a' && in <= 'z') || (in >= 'A' && in <= 'Z')
}

// IscCtrl checks if a character is a control character.
func IscCtrl(in byte) bool {
	return in < 32 || in == 127
}

// IscHex checks if a character is 0..9, A..F, or a..f.
func IscHex(in byte) bool {
	return (in >= '0' && in <= '9') || (in >= 'A' && in <= 'F') || (in >= 'a' && in <= 'f')
}

// IscLower checks if a character is lowercase.
func IscLower(in byte) bool {
	return in >= 'a' && in <= 'z'
}

// IscNum checks if a character is 0..9.
func IscNum(in byte) bool {
	return in >= '0' && in <= '9'
}

// IscUpper checks if a character is uppercase.
func IscUpper(in byte) bool {
	return in >= 'A' && in <= 'Z'
}

// Lowercase converts a string to lowercase.
func Lowercase(str string) string {
	return strings.ToLower(str)
}

// Message4R is a rotating message display.
type Message4R struct {
	Mx string
	Mn int
	Tr bool

	// internal state
	timer TON
	edge  bool
}

// Update executes the message rotation logic.
func (m *Message4R) Update(m0, m1, m2, m3 string, mm int, enq, clk bool, t1 time.Duration) {
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

// Message8 selects one of 8 messages based on prioritized inputs.
func Message8(in [8]bool, s [8]string) string {
	for i := 0; i < 8; i++ {
		if in[i] {
			return s[i]
		}
	}
	return ""
}

// Mirror reverses an input string.
func Mirror(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// MonthToString converts an integer (1-12) to a month name.
func MonthToString(mth, lang, lx int) string {
	// Placeholder for a complex lookup.
	if mth >= 1 && mth <= 12 {
		return time.Month(mth).String()
	}
	return ""
}

// OctToByte converts an octal string to a byte.
func OctToByte(oct string) byte {
	val, err := strconv.ParseUint(oct, 8, 8)
	if err != nil {
		return 0
	}
	return byte(val)
}

// OctToDword converts an octal string to a dword.
func OctToDword(oct string) uint32 {
	val, err := strconv.ParseUint(oct, 8, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

// RealToStrF converts a float to a string with N decimal places.
func RealToStrF(in float64, n int, d string) string {
	return strconv.FormatFloat(in, 'f', n, 64)
}

// ReplaceAll replaces all occurrences of src in str with rep.
func ReplaceAll(str, src, rep string) string {
	return strings.ReplaceAll(str, src, rep)
}

// ReplaceChars replaces characters in str based on a mapping from src to rep.
func ReplaceChars(str, src, rep string) string {
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

// ReplaceUml replaces German umlauts with their two-letter equivalents.
func ReplaceUml(str string) string {
	r := strings.NewReplacer(
		"Ä", "Ae", "Ö", "Oe", "Ü", "Ue", "ß", "ss",
		"ä", "ae", "ö", "oe", "ü", "ue",
	)
	return r.Replace(str)
}

// Ticker creates a scrolling text effect.
type Ticker struct {
	Display string
	// internal state
	delay TP
	step  int
}

// Update executes the ticker logic.
func (t *Ticker) Update(text string, n int, pt time.Duration) {
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

// ToLower converts a character from uppercase to lowercase.
func ToLower(in byte) byte {
	if in >= 'A' && in <= 'Z' {
		return in + ('a' - 'A')
	}
	// Placeholder for extended ASCII
	return in
}

// ToUml converts a character to its two-letter Umlaut representation.
func ToUml(in byte) string {
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

// ToUpper converts a character from lowercase to uppercase.
func ToUpper(in byte) byte {
	if in >= 'a' && in <= 'z' {
		return in - ('a' - 'A')
	}
	// Placeholder for extended ASCII
	return in
}

// Trim removes all space characters from a string.
func Trim(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

// Trim1 replaces multiple spaces with a single space and trims leading/trailing spaces.
func Trim1(str string) string {
	// Use Fields to split by whitespace and Join to put it back with single spaces.
	return strings.Join(strings.Fields(str), " ")
}

// Trime removes leading and trailing space characters from a string.
func Trime(str string) string {
	return strings.TrimSpace(str)
}

// Uppercase converts a string to uppercase.
func Uppercase(str string) string {
	return strings.ToUpper(str)
}

// WeekdayToString converts an integer (1-7) to a weekday name.
func WeekdayToString(wday, lang, lx int) string {
	// Placeholder for a complex lookup.
	if wday >= 1 && wday <= 7 {
		// Go's Sunday is 0, OSCAT's is 7. Adjusting for Go's standard.
		goWday := time.Weekday((wday % 7))
		return goWday.String()
	}
	return ""
}
