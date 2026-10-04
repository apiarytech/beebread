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
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/iec"
)

const loopback iec.DWORD = 0x7F000001

// conn is an IP_CONTROL and the connection it shares.
type conn struct {
	c    network.IP_C
	s, r network.NETWORK_BUFFER
	x    ip.IP_CONTROL
}

func newConn() *conn {
	k := &conn{}
	k.x = ip.IP_CONTROL{IP_C: &k.c, S_BUF: &k.s, R_BUF: &k.r}
	k.x.INIT()
	k.x.TIME_OUT = iec.TIME(2 * time.Second)
	return k
}

// scan runs the blocks every 2 ms until done or 5 s pass.
func scan(t *testing.T, done func() bool, blocks ...interface{ Execute(time.Time) }) {
	t.Helper()
	end := time.Now().Add(5 * time.Second)
	for time.Now().Before(end) {
		now := time.Now()
		for _, b := range blocks {
			b.Execute(now)
		}
		if done() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("timed out")
}

func TestHTTPGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.NotFound(w, r)
			return
		}
		user, pass, _ := r.BasicAuth()
		fmt.Fprintf(w, "path=%s q=%s user=%s pass=%s", r.URL.Path, r.URL.RawQuery, user, pass)
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	port, _ := strconv.Atoi(u.Port())

	k := newConn()
	urlData := encoding.STRING_TO_URL(iec.STRING(fmt.Sprintf("http://anna:secret@127.0.0.1:%d/a/b?x=1", port)), "", "")
	h := HTTP_GET{client: client{IP_C: &k.c, S_BUF: &k.s, R_BUF: &k.r}, URL_DATA: &urlData}
	h.INIT()
	h.IP4 = loopback
	h.GET = true
	scan(t, func() bool { return bool(h.DONE) || h.ERROR != 0 }, &k.x, &h)
	if h.ERROR != 0 {
		t.Fatalf("ERROR %08X", h.ERROR)
	}
	body := buffer.BUFFER_TO_STRING(k.r.BUFFER[:], k.r.SIZE, h.BODY_START, h.BODY_STOP)
	if h.HTTP_STATUS != "200 OK" || body != "path=/a/b q=x=1 user=anna pass=secret" {
		t.Errorf("status %q, body %q", h.HTTP_STATUS, body)
	}
	if h.HEADER_START != 0 || h.BODY_START != h.HEADER_STOP+1 || h.BODY_STOP != k.r.SIZE-1 {
		t.Errorf("header 0..%d, body %d..%d of %d", h.HEADER_STOP, h.BODY_START, h.BODY_STOP, k.r.SIZE)
	}

	// Release the connection, then a page that is not there.
	h.GET, h.UNLOCK_BUF = false, true
	scan(t, func() bool { return h.state == 0 }, &k.x, &h)
	// A request while the connection closes is lost, as in OSCAT: IP_CONTROL
	// clears S_BUF when it is closed.
	scan(t, func() bool { return k.c.C_STATE == 0 }, &k.x, &h)
	urlData = encoding.STRING_TO_URL(iec.STRING(fmt.Sprintf("http://127.0.0.1:%d/missing", port)), "", "")
	h.UNLOCK_BUF, h.GET = false, true
	scan(t, func() bool { return h.ERROR != 0 }, &k.x, &h)
	if h.ERROR != 0xFC {
		t.Errorf("a 404: ERROR %08X, status %q", h.ERROR, h.HTTP_STATUS)
	}
}

// dnsServer answers the DNS queries on 127.0.0.1:53 with the address a, or
// skips the test.
func dnsServer(t *testing.T, name string, a [4]byte) {
	pc, err := net.ListenPacket("udp4", "127.0.0.1:53")
	if err != nil {
		t.Skipf("no DNS server on port 53: %v", err)
	}
	t.Cleanup(func() { pc.Close() })
	go func() {
		buf := make([]byte, 512)
		for {
			n, from, err := pc.ReadFrom(buf)
			if err != nil {
				return
			}
			q := buf[:n]
			// The answer: the header, the question, and an A record.
			resp := append([]byte{}, q[:2]...)
			resp = append(resp, 0x81, 0x80, 0, 1, 0, 1, 0, 0, 0, 0)
			resp = append(resp, q[12:]...)
			resp = append(resp, 0xC0, 0x0C, 0, 1, 0, 1, 0, 0, 0x0E, 0x10, 0, 4, a[0], a[1], a[2], a[3])
			pc.WriteTo(resp, from)
		}
	}()
}

func TestDNSClient(t *testing.T) {
	dnsServer(t, "example.org", [4]byte{93, 184, 216, 34})
	k := newConn()
	d := DNS_CLIENT{client: client{IP_C: &k.c, S_BUF: &k.s, R_BUF: &k.r}}
	d.INIT()
	d.DOMAIN, d.IP4_DNS, d.ACTIVATE = "example.org", loopback, true
	scan(t, func() bool { return bool(d.DONE) || d.ERROR != 0 }, &k.x, &d)
	if d.ERROR != 0 || d.IP4 != 0x5DB8D822 {
		t.Errorf("IP4 %08X, ERROR %08X", d.IP4, d.ERROR)
	}
	// The query: example.org as labels.
	if got := string(byteString(k.s.BUFFER[12:24])); got != "\x07example\x03org" {
		t.Errorf("the query name: %q", got)
	}
}

func TestDNSAddress(t *testing.T) {
	var k conn
	d := DNS_CLIENT{client: client{IP_C: &k.c, S_BUF: &k.s, R_BUF: &k.r}}
	d.INIT()
	d.DOMAIN, d.ACTIVATE = "10.1.2.3", true
	d.Execute(time.Time{})
	d.Execute(time.Time{})
	if !d.DONE || d.IP4 != 0x0A010203 {
		t.Errorf("an address: DONE %v, IP4 %08X", d.DONE, d.IP4)
	}
}

func byteString(b []iec.BYTE) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = byte(b[i])
	}
	return out
}

func TestReadHTTP(t *testing.T) {
	var buf network.NW_BUF_LONG
	h := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\ncontent-length: 42\r\n\r\n"
	for i := range h {
		buf[i] = iec.BYTE(h[i])
	}
	r := READ_HTTP{PT: &buf, SIZE: iec.UINT(len(h))}
	r.Execute(time.Time{})
	if r.READ_HTTP != "1.1 200 OK" {
		t.Errorf("status: %q", r.READ_HTTP)
	}
	r.SRC = "CONTENT-LENGTH"
	r.Execute(time.Time{})
	if r.READ_HTTP != "42" {
		t.Errorf("Content-Length: %q", r.READ_HTTP)
	}
}
