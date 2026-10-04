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

package inet

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// nativePassword is MySQL's native password scramble.
func nativePassword(pw string, nonce []byte) []byte {
	h1 := sha1.Sum([]byte(pw))
	h2 := sha1.Sum(h1[:])
	h3 := sha1.Sum(append(append([]byte{}, nonce...), h2[:]...))
	out := make([]byte, 20)
	for i := range out {
		out[i] = h1[i] ^ h3[i]
	}
	return out
}

func TestMySQLAuth(t *testing.T) {
	nonce := []byte("abcdefghijklmnopqrst")
	var msg, scramble [20]iec.BYTE
	for i := range msg {
		msg[i] = iec.BYTE(nonce[i])
	}
	run, pw := iec.BOOL(true), iec.STRING("geheim")
	a := MYSQL_AUTH{RUN: &run, PASSWORD: &pw, MESSAGE: &msg, SCRAMBLE: &scramble}
	a.INIT()
	for i := 0; i < 20 && run; i++ {
		a.Execute(time.Time{})
	}
	want := nativePassword("geheim", nonce)
	for i := range want {
		if byte(scramble[i]) != want[i] {
			t.Fatalf("scramble %x, want %x", scramble, want)
		}
	}
}

// mysqlServer accepts one login and checks its scramble.
func mysqlServer(t *testing.T, ln net.Listener, nonce []byte, ok chan<- string) {
	c, err := ln.Accept()
	if err != nil {
		return
	}
	defer c.Close()
	packet := func(no byte, b []byte) {
		h := []byte{byte(len(b)), byte(len(b) >> 8), byte(len(b) >> 16), no}
		c.Write(append(h, b...))
	}
	read := func() []byte {
		h := make([]byte, 4)
		if _, err := io.ReadFull(c, h); err != nil {
			return nil
		}
		b := make([]byte, int(h[0])|int(h[1])<<8|int(h[2])<<16)
		io.ReadFull(c, b)
		return b
	}
	// The handshake, protocol 10.
	var hs bytes.Buffer
	hs.WriteByte(10)
	hs.WriteString("5.0.96\x00")
	hs.Write([]byte{1, 0, 0, 0})
	hs.Write(nonce[:8])
	hs.WriteByte(0)
	hs.Write([]byte{0x2C, 0xA2, 8, 2, 0})
	hs.Write(make([]byte, 13))
	hs.Write(nonce[8:])
	hs.WriteByte(0)
	packet(0, hs.Bytes())
	auth := read()
	user, rest, _ := bytes.Cut(auth[32:], []byte{0})
	got := rest[1 : 1+rest[0]]
	if !bytes.Equal(got, nativePassword("geheim", nonce)) {
		ok <- fmt.Sprintf("bad scramble for %s", user)
		packet(2, append([]byte{0xFF, 0x15, 0x04}, "#28000Access denied for 'x'"...))
		return
	}
	packet(2, []byte{0, 0, 0, 2, 0, 0, 0})
	ok <- string(user)
	if q := read(); len(q) == 1 && q[0] == 1 {
		ok <- "QUIT"
	}
}

func TestMySQLControl(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	defer ln.Close()
	nonce := []byte("12345678ABCDEFGHIJKL")
	ok := make(chan string, 2)
	go mysqlServer(t, ln, nonce, ok)

	var com network.MYSQL_COM
	var info network.MYSQL_INFO
	m := MYSQL_CONTROL{COM: &com, INFO: &info}
	m.INIT()
	com.SQL_URL = iec.STRING(fmt.Sprintf("mysql://plc:geheim@127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port))
	com.SQL_CON = true
	scan(t, func() bool { return bool(info.SQL_CONNECTED) || com.ERROR_T != 0 }, &m)
	if !info.SQL_CONNECTED {
		t.Fatalf("not connected: ERROR %d/%d %q", com.ERROR_C, com.ERROR_T, info.SQL_ERROR)
	}
	if user := <-ok; user != "plc" {
		t.Fatalf("the server: %s", user)
	}
	if info.SERVER_PROTOCOL_VERSION != 10 || info.SERVER_LANGUAGE != 8 {
		t.Errorf("INFO %+v", info)
	}
	com.SQL_CON = false
	scan(t, func() bool { return m.step == 0 }, &m)
	select {
	case q := <-ok:
		if q != "QUIT" {
			t.Errorf("after the login: %s", q)
		}
	case <-time.After(2 * time.Second):
		t.Error("no QUIT")
	}
}
