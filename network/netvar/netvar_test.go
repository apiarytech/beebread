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

package netvar

import (
	"net"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// plc is a side: its control and its variables.
type plc struct {
	x   network.NET_VAR_DATA
	ctl NET_VAR_CONTROL
	b   NET_VAR_BOOL8
	r   NET_VAR_REAL8
	x8  NET_VAR_X8
	s   NET_VAR_STRING
	in  iec.STRING
	out iec.STRING
}

func newPLC(master bool, port int) *plc {
	p := &plc{}
	p.ctl = NET_VAR_CONTROL{X: &p.x}
	p.ctl.INIT()
	p.ctl.MASTER, p.ctl.REMOTE_IP4, p.ctl.REMOTE_PORT = iec.BOOL(master), 0x7F000001, iec.WORD(port)
	p.ctl.SCAN_TIME = iec.TIME(20 * time.Millisecond)
	p.b = NET_VAR_BOOL8{member: member{X: &p.x}}
	p.r = NET_VAR_REAL8{member: member{X: &p.x}}
	p.x8 = NET_VAR_X8{member: member{X: &p.x}}
	p.s = NET_VAR_STRING{member: member{X: &p.x}, IN: &p.in, OUT: &p.out}
	return p
}

// scan runs the side once: the control, then the variables.
func (p *plc) scan(now time.Time) {
	p.ctl.Execute(now)
	p.b.Execute(now)
	p.r.Execute(now)
	p.x8.Execute(now)
	p.s.Execute(now)
}

func TestNetVar(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	m, s := newPLC(true, port), newPLC(false, port)
	m.b.IN1, m.b.IN8 = true, true
	m.r.IN3 = 3.25
	m.x8.IN_DINT2, m.x8.IN_REAL1 = -42, 1.5
	m.in = "from the master"
	s.b.IN2 = true
	s.r.IN1 = -7.5
	s.x8.IN_UDINT1, s.x8.IN_DWORD2 = 4000000000, 0xDEADBEEF
	s.in = "from the slave"

	s.ctl.ACTIVATE, m.ctl.ACTIVATE = true, true
	end := time.Now().Add(5 * time.Second)
	for time.Now().Before(end) && m.x.CYCLE < 3 {
		now := time.Now()
		s.scan(now)
		m.scan(now)
		time.Sleep(time.Millisecond)
	}
	if m.x.CYCLE < 3 || !m.ctl.RUN || !s.ctl.RUN {
		t.Fatalf("cycles %d, RUN %v %v, ERROR %08X %08X", m.x.CYCLE, m.ctl.RUN, s.ctl.RUN, m.ctl.ERROR, s.ctl.ERROR)
	}
	if !s.b.OUT1 || s.b.OUT2 || !s.b.OUT8 || !m.b.OUT2 || m.b.OUT1 {
		t.Errorf("BOOL8: slave %v %v %v, master %v %v", s.b.OUT1, s.b.OUT2, s.b.OUT8, m.b.OUT2, m.b.OUT1)
	}
	if s.r.OUT3 != 3.25 || m.r.OUT1 != -7.5 {
		t.Errorf("REAL8: %v %v", s.r.OUT3, m.r.OUT1)
	}
	if s.x8.OUT_DINT2 != -42 || s.x8.OUT_REAL1 != 1.5 || m.x8.OUT_UDINT1 != 4000000000 || m.x8.OUT_DWORD2 != 0xDEADBEEF {
		t.Errorf("X8: %d %v %d %08X", s.x8.OUT_DINT2, s.x8.OUT_REAL1, m.x8.OUT_UDINT1, m.x8.OUT_DWORD2)
	}
	if s.out != "from the master" || m.out != "from the slave" {
		t.Errorf("STRING: %q %q", s.out, m.out)
	}
}
