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

package ip

import (
	"io"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

const loopback iec.DWORD = 0x7F000001

// scan runs step every 2 ms until done or 3 s pass.
func scan(t *testing.T, step func(time.Time), done func() bool) {
	t.Helper()
	end := time.Now().Add(3 * time.Second)
	for time.Now().Before(end) {
		step(time.Now())
		if done() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("timed out")
}

func put(b *network.NETWORK_BUFFER, s string) {
	for i := range s {
		b.BUFFER[i] = iec.BYTE(s[i])
	}
	b.SIZE = iec.UINT(len(s))
}

func get(b *network.NETWORK_BUFFER) string {
	out := make([]byte, b.SIZE)
	for i := range out {
		out[i] = byte(b.BUFFER[i])
	}
	return string(out)
}

func TestTCPClient(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	defer ln.Close()
	go func() {
		c, err := ln.Accept()
		if err == nil {
			io.Copy(c, c) // echo
			c.Close()
		}
	}()
	port := iec.WORD(ln.Addr().(*net.TCPAddr).Port)

	var c network.IP_C
	var s, r network.NETWORK_BUFFER
	x := IP_CONTROL{IP_C: &c, S_BUF: &s, R_BUF: &r}
	x.INIT()
	x.IP, x.PORT, x.TIME_OUT = loopback, port, iec.TIME(time.Second)
	c.C_ENABLE = true
	scan(t, x.Execute, func() bool { return c.C_STATE == 255 })

	put(&s, "hello, plc")
	scan(t, x.Execute, func() bool { return r.SIZE == 10 })
	if got := get(&r); got != "hello, plc" || s.SIZE != 0 || c.MAILBOX[0] == 0 {
		t.Errorf("received %q, S_BUF.SIZE %d, MAILBOX[1] %d", got, s.SIZE, c.MAILBOX[0])
	}

	c.C_ENABLE = false
	scan(t, x.Execute, func() bool { return c.C_STATE == 0 })
}

func TestUDPClient(t *testing.T) {
	pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	defer pc.Close()
	go func() {
		buf := make([]byte, 1500)
		n, from, err := pc.ReadFrom(buf)
		if err == nil {
			pc.WriteTo(append([]byte("re: "), buf[:n]...), from)
		}
	}()
	port := iec.WORD(pc.LocalAddr().(*net.UDPAddr).Port)

	var c network.IP_C
	var s, r network.NETWORK_BUFFER
	x := IP_CONTROL{IP_C: &c, S_BUF: &s, R_BUF: &r}
	x.INIT()
	x.IP, x.PORT, x.TIME_OUT = loopback, port, iec.TIME(time.Second)
	c.C_MODE = 1 // UDP client
	c.C_ENABLE = true
	scan(t, x.Execute, func() bool { return c.C_STATE == 255 })
	put(&s, "ping")
	scan(t, x.Execute, func() bool { return r.SIZE > 0 })
	if got := get(&r); got != "re: ping" {
		t.Errorf("received %q", got)
	}
}

func TestTCPServer(t *testing.T) {
	// A free port.
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	var c network.IP_C
	var s, r network.NETWORK_BUFFER
	x := IP_CONTROL{IP_C: &c, S_BUF: &s, R_BUF: &r}
	x.INIT()
	x.IP, x.PORT, x.TIME_OUT = loopback, iec.WORD(port), iec.TIME(time.Second)
	c.C_MODE = 2 // TCP server for the remote IP
	c.C_ENABLE = true
	for i := 0; i < 10; i++ {
		x.Execute(time.Now())
	}
	conn, err := net.Dial("tcp4", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	scan(t, x.Execute, func() bool { return c.C_STATE == 255 })
	conn.Write([]byte("from client"))
	scan(t, x.Execute, func() bool { return r.SIZE == 11 })
	if got := get(&r); got != "from client" {
		t.Errorf("received %q", got)
	}
}

func TestConnectTimeout(t *testing.T) {
	// A port nothing listens on.
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	var c network.IP_C
	var s, r network.NETWORK_BUFFER
	x := IP_CONTROL{IP_C: &c, S_BUF: &s, R_BUF: &r}
	x.INIT()
	x.IP, x.PORT, x.TIME_OUT = loopback, iec.WORD(port), iec.TIME(200*time.Millisecond)
	c.C_ENABLE = true
	scan(t, x.Execute, func() bool { return c.ERROR != 0 })
	if c.ERROR>>24 != 255 {
		t.Errorf("ERROR %08X, want a connect time out", c.ERROR)
	}
}

func TestFIFO(t *testing.T) {
	var q network.IP_FIFO_DATA
	var id1, id2, st1, st2 iec.BYTE
	f1 := IP_FIFO{FIFO: &q, ID: &id1, STATE: &st1}
	f2 := IP_FIFO{FIFO: &q, ID: &id2, STATE: &st2}
	st1, st2 = 1, 1
	f1.Execute(time.Time{})
	f2.Execute(time.Time{})
	if id1 != 1 || id2 != 2 || st1 != 3 || st2 != 2 {
		t.Fatalf("ids %d %d, states %d %d", id1, id2, st1, st2)
	}
	st1 = 4
	f1.Execute(time.Time{})
	f2.Execute(time.Time{})
	if st1 != 5 || st2 != 3 {
		t.Fatalf("after the first leaves: states %d %d", st1, st2)
	}
}
