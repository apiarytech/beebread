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

package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	gobase64 "encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// lengths are the lengths to hash: around the 55 and 64 bytes where the
// padding takes another block.
var lengths = []int{0, 1, 3, 54, 55, 56, 57, 63, 64, 65, 119, 120, 128, 200, 250}

func text(n int) string {
	return strings.Repeat("The quick brown fox jumps over the lazy dog. ", 6)[:n]
}

func TestMD5(t *testing.T) {
	for _, n := range lengths {
		s := iec.STRING(text(n))
		var sum [16]iec.BYTE
		m := MD5_STR{STR: &s, MD5: &sum}
		m.INIT()
		m.RUN = true
		for i := 0; i < 20 && !m.DONE; i++ {
			m.Execute(time.Time{})
		}
		want := md5.Sum([]byte(text(n)))
		if got := MD5_TO_STRH(sum); string(got) != hex.EncodeToString(want[:]) {
			t.Errorf("MD5 of %d bytes = %s, want %x", n, got, want)
		}
	}
}

func TestSHA1(t *testing.T) {
	for _, n := range lengths {
		s := iec.STRING(text(n))
		var sum [20]iec.BYTE
		h := SHA1_STR{STR: &s, SHA1: &sum}
		h.INIT()
		h.RUN = true
		for i := 0; i < 20 && !h.DONE; i++ {
			h.Execute(time.Time{})
		}
		want := sha1.Sum([]byte(text(n)))
		if got := SHA1_TO_STRH(sum); string(got) != hex.EncodeToString(want[:]) {
			t.Errorf("SHA1 of %d bytes = %s, want %x", n, got, want)
		}
	}
}

func TestRC4(t *testing.T) {
	key := iec.STRING("Secret")
	data := []byte(text(150))
	var buf [64]iec.BYTE
	var mode iec.INT = 1
	size := iec.UDINT(len(data))
	r := RC4_CRYPT_STREAM{MODE: &mode, KEY: &key, BUF: &buf, SIZE: &size}
	r.INIT()
	r.Execute(time.Time{}) // set up
	var got []byte
	for mode == 2 {
		pos := int(r.POS)
		for i := 0; i < int(size); i++ {
			buf[i] = iec.BYTE(data[pos+i])
		}
		n := int(size)
		r.Execute(time.Time{})
		for i := 0; i < n; i++ {
			got = append(got, byte(buf[i]))
		}
	}
	c, _ := rc4.NewCipher([]byte(key))
	want := make([]byte, len(data))
	c.XORKeyStream(want, data)
	if string(got) != string(want) {
		t.Errorf("RC4 = %x\nwant %x", got, want)
	}
}

func TestMD5CramAuth(t *testing.T) {
	challenge := "<1896.697170952@postoffice.reston.mci.net>"
	run := iec.BOOL(true)
	user, pass := iec.STRING("tim"), iec.STRING("tanstaaftanstaaf")
	ts := iec.STRING(gobase64.StdEncoding.EncodeToString([]byte(challenge)))
	var key iec.STRING
	m := MD5_CRAM_AUTH{RUN: &run, USERNAME: &user, PASSWORD: &pass, B64_TS: &ts, AUTH_KEY: &key}
	m.INIT()
	for i := 0; i < 50 && run; i++ {
		m.Execute(time.Time{})
	}
	mac := hmac.New(md5.New, []byte(pass))
	mac.Write([]byte(challenge))
	want := gobase64.StdEncoding.EncodeToString([]byte("tim " + hex.EncodeToString(mac.Sum(nil))))
	// RFC 2195's example answer.
	if want != "dGltIGI5MTNhNjAyYzdlZGE3YTQ5NWI0ZTZlNzMzNGQzODkw" {
		t.Fatalf("the reference: %s", want)
	}
	if run || string(key) != want {
		t.Errorf("AUTH_KEY = %q (RUN %v), want %q", key, run, want)
	}
}
