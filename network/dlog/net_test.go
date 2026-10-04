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

package dlog

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

const loopback iec.DWORD = 0x7F00_0001

// dnsServer answers the DNS queries on 127.0.0.1:53 with 127.0.0.1, or
// skips the test.
func dnsServer(t *testing.T) {
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
			resp := append([]byte{}, q[:2]...)
			resp = append(resp, 0x81, 0x80, 0, 1, 0, 1, 0, 0, 0, 0)
			resp = append(resp, q[12:]...)
			resp = append(resp, 0xC0, 0x0C, 0, 1, 0, 1, 0, 0, 0x0E, 0x10, 0, 4, 127, 0, 0, 1)
			pc.WriteTo(resp, from)
		}
	}()
}

// scan runs the blocks in real time until done, for 5 s at most.
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

func TestCronTab(t *testing.T) {
	c := DLOG_CRON_TAB{SECOND: "0,30", MINUTE: "*/15", HOUR: "8-17", DAY_OF_MONTH: "", DAY_OF_WEEK: "1-5", MONTH: "*"}
	c.ACTIVATE = true
	c.Execute(time.Time{}) // reads the table
	at := func(mo, d, h, mi, s int) iec.BOOL {
		c.DTI = td.SET_DT(2026, iec.INT(mo), iec.INT(d), iec.INT(h), iec.INT(mi), iec.INT(s))
		c.Execute(time.Time{})
		return c.Q
	}
	// 2026-10-05 is a Monday, 2026-10-04 a Sunday.
	for _, k := range []struct {
		mo, d, h, mi, s int
		q               iec.BOOL
	}{
		{10, 5, 8, 0, 0, true},
		{10, 5, 8, 15, 30, true},
		{10, 5, 8, 15, 31, false},
		{10, 5, 8, 16, 0, false},
		{10, 5, 18, 0, 0, false},
		{10, 4, 9, 0, 0, false},
		{10, 9, 17, 45, 30, true},
	} {
		if q := at(k.mo, k.d, k.h, k.mi, k.s); q != k.q {
			t.Errorf("%02d-%02d %02d:%02d:%02d: Q %v", k.mo, k.d, k.h, k.mi, k.s, q)
		}
	}
	// Both days given: either matches.
	c = DLOG_CRON_TAB{SECOND: "0", MINUTE: "0", HOUR: "0", DAY_OF_MONTH: "4", DAY_OF_WEEK: "1", MONTH: "*", ACTIVATE: true}
	c.Execute(time.Time{})
	if !at(10, 4, 0, 0, 0) || !at(10, 5, 0, 0, 0) || at(10, 6, 0, 0, 0) {
		t.Error("day of the month or of the week")
	}
}

func TestStoreRRD(t *testing.T) {
	dnsServer(t)
	query := make(chan string, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query <- r.URL.RawQuery
		fmt.Fprint(w, "0")
	}))
	defer srv.Close()
	u, _ := url.Parse(srv.URL)

	var x network.DLOG_DATA
	s := DLOG_STORE_RRD{}
	s.X = &x
	s.INIT()
	s.URL, s.SEP, s.DNS_IP4, s.TIMEOUT = iec.STRING("http://rrd.test:"+u.Port()+"/up?v="), ':', loopback, iec.TIME(2*time.Second)
	s.DTI = td.SET_DT(2026, 10, 3, 12, 0, 0)
	a := DLOG_DINT{VALUE: 5}
	a.X = &x
	b := DLOG_STRING{STR: "x"}
	b.X = &x
	s.ENABLE = true
	scan(t, func() bool { return s.step1 == 30 }, &s, &a, &b)
	s.TRIG_M = true
	scan(t, func() bool { return len(query) > 0 || s.ERROR_T != 0 }, &s, &a, &b)
	if s.ERROR_T != 0 {
		t.Fatalf("error %d %X", s.ERROR_T, s.ERROR_C)
	}
	if q := <-query; q != "v=5:x" {
		t.Errorf("query %q", q)
	}
	scan(t, func() bool { return s.step2 == 0 }, &s, &a, &b)
	if s.ERROR_T != 0 {
		t.Errorf("error %d %X", s.ERROR_T, s.ERROR_C)
	}
}

func TestFileToFTPRetry(t *testing.T) {
	// No DNS answer: the transfer fails, is tried once more, and dropped.
	var x network.DLOG_DATA
	f := DLOG_FILE_TO_FTP{}
	f.X = &x
	f.INIT()
	f.FTP_URL, f.DNS_IP4 = "ftp://u:p@ftp.test/dir", 0x7F00_00FE
	f.TIMEOUT, f.RETRY, f.RETRY_TIME = iec.TIME(200*time.Millisecond), 1, iec.TIME(50*time.Millisecond)
	if !f.DONE && f.step != 0 {
		t.Fatal("not idle")
	}
	x.NEW_FILE, x.NEW_FILE_RTRIG = "a.csv", true
	tries := 0
	busy := false
	scan(t, func() bool {
		if bool(f.BUSY) && !busy {
			tries++
		}
		busy = bool(f.BUSY)
		return tries == 2 && bool(f.DONE)
	}, &f)
	if f.ERROR_T == 0 {
		t.Error("no error reported")
	}
	if f.ucbd.BUF_COUNT != 0 {
		t.Errorf("queue %d", f.ucbd.BUF_COUNT)
	}
}
