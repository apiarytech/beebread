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
	"bufio"
	"net"
	"strings"
	"testing"
	"time"

	. "github.com/apiarytech/beebread/basic"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

func TestSNTP(t *testing.T) {
	// Port 123 must be free for the server.
	if pc, err := net.ListenPacket("udp4", "127.0.0.1:123"); err != nil {
		t.Skipf("port 123: %v", err)
	} else {
		pc.Close()
	}
	sk := newConn()
	srv := SNTP_SERVER{client: client{IP_C: &sk.c, S_BUF: &sk.s, R_BUF: &sk.r}}
	srv.INIT()
	srv.ENABLE = true
	srv.UDT = td.SET_DT(2024, 5, 6, 7, 8, 9)
	srv.XMS = 250

	ck := newConn()
	cl := SNTP_CLIENT{client: client{IP_C: &ck.c, S_BUF: &ck.s, R_BUF: &ck.r}}
	cl.INIT()
	cl.IP4 = loopback
	// The server first.
	scan(t, func() bool { return sk.c.C_STATE == 255 }, &sk.x, &srv)
	cl.ACTIVATE = true
	scan(t, func() bool { return bool(cl.DONE_P) || cl.ERROR != 0 }, &sk.x, &srv, &ck.x, &cl)
	if cl.ERROR != 0 {
		t.Fatalf("ERROR %08X", cl.ERROR)
	}
	if DT_TO_DWORD(cl.UDT) != DT_TO_DWORD(srv.UDT) || cl.XMS < 249 || cl.XMS > 400 {
		t.Errorf("UDT %d, XMS %d; want %d, 250 and half the round trip", DT_TO_DWORD(cl.UDT), cl.XMS, DT_TO_DWORD(srv.UDT))
	}
}

func TestSysLog(t *testing.T) {
	pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	defer pc.Close()
	got := make(chan string, 1)
	go func() {
		buf := make([]byte, 1500)
		if n, _, err := pc.ReadFrom(buf); err == nil {
			got <- string(buf[:n])
		}
	}()
	k := newConn()
	s := SYS_LOG{client: client{IP_C: &k.c, S_BUF: &k.s, R_BUF: &k.r}}
	s.INIT()
	s.SERVER_IP4, s.PORT = loopback, iec.WORD(pc.LocalAddr().(*net.UDPAddr).Port)
	s.FACILITY, s.SEVERITY = 1, 5 // user, notice: <13>
	s.LDT = td.SET_DT(2024, 1, 5, 9, 7, 3)
	s.HOSTNAME, s.TAG, s.MESSAGE = "plc1", "app:", "hello"
	s.ACTIVATE = true
	scan(t, func() bool { return len(got) > 0 || s.ERROR != 0 }, &k.x, &s)
	if s.ERROR != 0 {
		t.Fatalf("ERROR %08X", s.ERROR)
	}
	msg := <-got
	if msg != "<13>Jan  5 09:07:03 plc1 app: hello" {
		t.Errorf("the message: %q", msg)
	}

	// OSCAT builds the priority and the header in a STRING(20): a priority of
	// 3 digits cuts the space after the time.
	s.FACILITY, s.ACTIVATE = 16, false // local0: <133>
	scan(t, func() bool { return s.state == 0 }, &k.x, &s)
	s.ACTIVATE = true
	go func() {
		buf := make([]byte, 1500)
		if n, _, err := pc.ReadFrom(buf); err == nil {
			got <- string(buf[:n])
		}
	}()
	scan(t, func() bool { return len(got) > 0 }, &k.x, &s)
	if msg := <-got; msg != "<133>Jan  5 09:07:03plc1 app: hello" {
		t.Errorf("a priority of 3 digits: %q", msg)
	}
}

func TestTelnetPrint(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	k := newConn()
	p := TELNET_PRINT{IP_C: &k.c, S_BUF: &k.s}
	p.INIT()
	p.PORT = iec.WORD(port)
	p.OPTION = 0x88 // CR LF after each text, send at once
	p.ENABLE = true
	for i := 0; i < 20; i++ {
		now := time.Now()
		k.x.Execute(now)
		p.Execute(now)
	}
	conn, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	scan(t, func() bool { return bool(p.READY) }, &k.x, &p)
	p.TEXT, p.SEND = "hello, telnet", true
	scan(t, func() bool { return bool(p.DONE) }, &k.x, &p)
	p.SEND = false
	scan(t, func() bool { return k.s.SIZE == 0 }, &k.x, &p)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	r := bufio.NewReader(conn)
	var out strings.Builder
	for !strings.Contains(out.String(), "hello, telnet\r\n") {
		b, err := r.ReadByte()
		if err != nil {
			t.Fatalf("read %q: %v", out.String(), err)
		}
		out.WriteByte(b)
	}
	if !strings.HasPrefix(out.String(), "\x1b[?7l") {
		t.Errorf("the screen setup: %q", out.String())
	}
}
