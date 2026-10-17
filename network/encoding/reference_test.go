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

package encoding

import (
	gobase64 "encoding/base64"
	"math/rand/v2"
	"strings"
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// latin1 returns the characters of s, which are ISO 8859-1, as bytes.
func latin1(s string) []byte {
	var b []byte
	for _, r := range s {
		b = append(b, byte(r))
	}
	return b
}

// randomText returns n random ISO 8859-1 characters.
func randomText(r *rand.Rand, n int) string {
	c := make([]rune, n)
	for i := range c {
		c[i] = rune(1 + r.IntN(255))
	}
	return string(c)
}

// TestBase64AgainstGo checks BASE64_ENCODE_STR and BASE64_DECODE_STR on
// random text of every length up to the 144 characters of STR1 against Go's
// encoding/base64.
func TestBase64AgainstGo(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for n := 0; n <= 144; n++ {
		s := randomText(r, n)
		in, enc, dec := iec.STRING(s), iec.STRING(""), iec.STRING("")
		e := BASE64_ENCODE_STR{STR1: &in, STR2: &enc}
		e.INIT()
		e.RUN = true
		for i := 0; i < 20 && !e.DONE; i++ {
			e.Execute(time.Time{})
		}
		if want := gobase64.StdEncoding.EncodeToString(latin1(s)); string(enc) != want {
			t.Fatalf("BASE64_ENCODE_STR of %d characters = %q, want %q", n, enc, want)
		}
		d := BASE64_DECODE_STR{STR1: &enc, STR2: &dec}
		d.INIT()
		d.RUN = true
		for i := 0; i < 20 && !d.DONE; i++ {
			d.Execute(time.Time{})
		}
		if string(dec) != s {
			t.Fatalf("BASE64_DECODE_STR of %d characters = %q, want %q", n, dec, s)
		}
	}
}

// TestURLRoundTrip checks that URL_DECODE undoes URL_ENCODE on random text
// without a %, and that URL_ENCODE leaves only characters a URL may hold.
func TestURLRoundTrip(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	for range 2000 {
		s := iec.STRING(strings.ReplaceAll(randomText(r, r.IntN(60)), "%", ""))
		enc := URL_ENCODE(s)
		for _, c := range latin1(string(enc)) {
			if c != '%' && !IS_URLCHR(iec.BYTE(c)) {
				t.Fatalf("URL_ENCODE(%q) = %q holds %q", s, enc, c)
			}
		}
		if dec := URL_DECODE(enc); dec != s {
			t.Fatalf("URL_DECODE(URL_ENCODE(%q)) = %q", s, dec)
		}
	}
}

// TestURLDecodePercent documents OSCAT's URL_DECODE, which looks for the
// next % from the start, so a % it has decoded is decoded again: text with a
// % followed by two hex digits does not come back from URL_ENCODE.
func TestURLDecodePercent(t *testing.T) {
	if got := URL_DECODE(URL_ENCODE("%41")); got != "A" {
		t.Errorf("URL_DECODE(URL_ENCODE('%%41')) = %q; OSCAT's gives 'A'", got)
	}
}

// TestHTMLRoundTrip checks that HTML_DECODE undoes HTML_ENCODE on random
// text, with and without the named entities.
func TestHTMLRoundTrip(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 6))
	for range 2000 {
		s := iec.STRING(randomText(r, r.IntN(40)))
		for _, m := range []iec.BOOL{false, true} {
			if dec := HTML_DECODE(HTML_ENCODE(s, m)); dec != s {
				t.Fatalf("HTML_DECODE(HTML_ENCODE(%q, %v)) = %q", s, m, dec)
			}
		}
	}
}
