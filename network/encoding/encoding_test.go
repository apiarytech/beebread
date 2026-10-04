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
	"strings"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

func TestBase64(t *testing.T) {
	for _, s := range []string{"", "f", "fo", "foo", "foob", "fooba", "foobar",
		"The quick brown fox jumps over the lazy dog and keeps running to the end of the line"} {
		in, enc, dec := iec.STRING(s), iec.STRING(""), iec.STRING("")
		e := BASE64_ENCODE_STR{STR1: &in, STR2: &enc}
		e.INIT()
		e.RUN = true
		for i := 0; i < 10 && !e.DONE; i++ {
			e.Execute(time.Time{})
		}
		if want := gobase64.StdEncoding.EncodeToString([]byte(s)); string(enc) != want {
			t.Errorf("BASE64_ENCODE_STR(%q) = %q, want %q", s, enc, want)
		}
		d := BASE64_DECODE_STR{STR1: &enc, STR2: &dec}
		d.INIT()
		d.RUN = true
		for i := 0; i < 10 && !d.DONE; i++ {
			d.Execute(time.Time{})
		}
		if string(dec) != s {
			t.Errorf("BASE64_DECODE_STR(%q) = %q, want %q", enc, dec, s)
		}
	}
}

func TestURL(t *testing.T) {
	if got := URL_ENCODE("a b&c/ü"); got != "a%20b%26c%2F%FC" {
		t.Errorf("URL_ENCODE = %q", got)
	}
	if got := URL_DECODE("a%20b%26c%2F%FC"); got != "a b&c/ü" {
		t.Errorf("URL_DECODE = %q", got)
	}
	u := STRING_TO_URL("http://hans:geheim@www.example.org:8080/demo/x.cgi?land=de#a1", "", "")
	want := network.URL{PROTOCOL: "http", USER: "hans", PASSWORD: "geheim", DOMAIN: "www.example.org",
		PORT: 8080, PATH: "/demo/x.cgi", QUERY: "land=de", ANCHOR: "a1"}
	if u != want {
		t.Errorf("STRING_TO_URL = %+v", u)
	}
	if got := URL_TO_STRING(u); got != "http://hans:geheim@www.example.org:8080/demo/x.cgi?land=de#a1" {
		t.Errorf("URL_TO_STRING = %q", got)
	}
	u = STRING_TO_URL("www.example.org", "http", "/")
	if u.PROTOCOL != "http" || u.PATH != "/" || u.DOMAIN != "www.example.org" {
		t.Errorf("defaults: %+v", u)
	}
	// OSCAT leaves out the last character of a user without a password.
	if u = STRING_TO_URL("ftp://anna@host", "", ""); u.USER != "ann" {
		t.Errorf("user without password: %q", u.USER)
	}
}

func TestHTML(t *testing.T) {
	if got := HTML_ENCODE(`<a href="x">&</a>`, false); got != "&lt;a href=&quot;x&quot;&gt;&amp;&lt;/a&gt;" {
		t.Errorf("HTML_ENCODE = %q", got)
	}
	if got := HTML_DECODE("&lt;b&gt; &#65;&#x42; &amp;"); got != "<b> AB &" {
		t.Errorf("HTML_DECODE = %q", got)
	}
	if got := HTML_ENCODE("ä", true); !strings.HasPrefix(string(got), "&") || !strings.HasSuffix(string(got), ";") {
		t.Errorf("HTML_ENCODE(ä, M) = %q", got)
	}
	if got := HTML_DECODE(HTML_ENCODE("ä", true)); got != "ä" {
		t.Errorf("round trip of ä: %q", got)
	}
}

func TestElements(t *testing.T) {
	list := iec.STRING("a,bb,,ccc")
	if n := ELEMENT_COUNT(',', &list); n != 4 {
		t.Errorf("ELEMENT_COUNT = %d", n)
	}
	for i, want := range []iec.STRING{"a", "bb", "", "ccc", ""} {
		if got := ELEMENT_GET(',', iec.INT(i), &list); got != want {
			t.Errorf("ELEMENT_GET(%d) = %q, want %q", i, got, want)
		}
	}
	empty := iec.STRING("")
	if n := ELEMENT_COUNT(',', &empty); n != 0 {
		t.Errorf("ELEMENT_COUNT('') = %d", n)
	}
}

func TestFilePath(t *testing.T) {
	var x network.FILE_PATH_DATA
	if !FILE_PATH_SPLIT(`C:\dir\sub\file.txt`, &x) || x.DRIVE != "C:" || x.DIRECTORY != `\dir\sub\` || x.FILENAME != "file.txt" {
		t.Errorf("FILE_PATH_SPLIT = %+v", x)
	}
	if FILE_PATH_SPLIT("", &x) {
		t.Error("an empty path")
	}
}

func TestIP4(t *testing.T) {
	if ip := IP4_DECODE("192.168.1.20"); ip != 0xC0A80114 {
		t.Errorf("IP4_DECODE = %08X", ip)
	}
	if s := IP4_TO_STRING(0xC0A80114); s != "192.168.1.20" {
		t.Errorf("IP4_TO_STRING = %q", s)
	}
	if !IS_IP4("10.0.0.1") || IS_IP4("example.org") {
		t.Error("IS_IP4")
	}
	if !IP4_CHECK(0xC0A80114, 0xC0A801FE, 0xFFFFFF00) || IP4_CHECK(0xC0A80114, 0xC0A802FE, 0xFFFFFF00) {
		t.Error("IP4_CHECK")
	}
}
