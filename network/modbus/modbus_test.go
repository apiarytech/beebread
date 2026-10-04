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

package modbus

import (
	"net"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/iec"
)

// end is a side of the connection: an IP_CONTROL2 and its data.
type end struct {
	c    network.IP_C
	s, r network.NETWORK_BUFFER_SHORT
	x    ip.IP_CONTROL2
	data DATA
}

func newEnd(port int) *end {
	e := &end{}
	e.x = ip.IP_CONTROL2{IP_C: &e.c, S_BUF: &e.s, R_BUF: &e.r}
	e.x.INIT()
	e.x.IP, e.x.PORT, e.x.TIME_OUT = 0x7F000001, iec.WORD(port), iec.TIME(2*time.Second)
	return e
}

func TestModbus(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	se, ce := newEnd(port), newEnd(port)
	var vm VMAP
	srv := MB_SERVER_1{server{shortClient: shortClient{IP_C: &se.c, S_BUF: &se.s, R_BUF: &se.r, DATA: &se.data}, VMAP: &vm}}
	srv.INIT()
	srv.DATA_SIZE, srv.ENABLE = 256, true
	cl := MB_CLIENT{shortClient: shortClient{IP_C: &ce.c, S_BUF: &ce.s, R_BUF: &ce.r, DATA: &ce.data}}
	cl.INIT()
	cl.DATA_SIZE = 256

	step := func() {
		now := time.Now()
		se.x.Execute(now)
		srv.Execute(now)
		ce.x.Execute(now)
		cl.Execute(now)
	}
	// The server listens first.
	for i := 0; i < 20; i++ {
		step()
	}
	// request runs one request of the client.
	request := func(fc iec.INT) {
		t.Helper()
		cl.FC, cl.ENABLE = fc, true
		end := time.Now().Add(3 * time.Second)
		started := false
		for time.Now().Before(end) {
			step()
			if cl.BUSY {
				started = true
			} else if started && cl.state == 0 {
				cl.ENABLE = false
				step()
				return
			}
			time.Sleep(time.Millisecond)
		}
		t.Fatalf("FC %d: no answer", fc)
	}

	// FC 3: read 3 registers.
	se.data[10], se.data[11], se.data[12] = 0x1234, 0xBEEF, 7
	cl.R_ADDR, cl.R_POINTS, cl.R_DATA_ADR = 10, 3, 100
	request(3)
	if cl.ERROR != 0 || ce.data[100] != 0x1234 || ce.data[101] != 0xBEEF || ce.data[102] != 7 {
		t.Fatalf("FC 3: ERROR %d, %04X", cl.ERROR, ce.data[100:103])
	}
	// FC 16: write 2 registers.
	ce.data[0], ce.data[1] = 0x1111, 0x2222
	cl.W_ADDR, cl.W_POINTS, cl.W_DATA_ADR = 20, 2, 0
	request(16)
	if cl.ERROR != 0 || se.data[20] != 0x1111 || se.data[21] != 0x2222 {
		t.Fatalf("FC 16: ERROR %d, %04X", cl.ERROR, se.data[20:22])
	}
	// FC 6: write a register.
	ce.data[5] = 0xABCD
	cl.W_ADDR, cl.W_DATA_ADR = 30, 5
	request(6)
	if cl.ERROR != 0 || se.data[30] != 0xABCD {
		t.Fatalf("FC 6: ERROR %d, %04X", cl.ERROR, se.data[30])
	}
	// FC 5: force coil 35, bit 3 of word 2.
	ce.data[6] = 1
	cl.W_ADDR, cl.W_DATA_ADR, cl.W_DATA_BITPOS = 35, 6, 0
	request(5)
	if cl.ERROR != 0 || se.data[2] != 1<<3 {
		t.Fatalf("FC 5: ERROR %d, %04X", cl.ERROR, se.data[2])
	}
	// FC 1: read 16 coils from 32, word 2.
	cl.R_ADDR, cl.R_POINTS, cl.R_DATA_ADR, cl.R_DATA_BITPOS = 32, 16, 110, 0
	request(1)
	if cl.ERROR != 0 || ce.data[110] != 1<<3 {
		t.Fatalf("FC 1: ERROR %d, %04X", cl.ERROR, ce.data[110])
	}
	// FC 22: mask write register 40.
	se.data[40] = 0x00FF
	ce.data[7], ce.data[8] = 0x0F0F, 0xF000 // AND, OR
	cl.W_ADDR, cl.W_DATA_ADR = 40, 7
	request(22)
	if want := iec.WORD(0x00FF&0x0F0F | 0xF000&^0x0F0F); cl.ERROR != 0 || se.data[40] != want {
		t.Fatalf("FC 22: ERROR %d, %04X want %04X", cl.ERROR, se.data[40], want)
	}
	// FC 23: write 50, read 10..11.
	ce.data[9] = 0x5555
	cl.W_ADDR, cl.W_POINTS, cl.W_DATA_ADR = 50, 1, 9
	cl.R_ADDR, cl.R_POINTS, cl.R_DATA_ADR = 10, 2, 120
	request(23)
	if cl.ERROR != 0 || se.data[50] != 0x5555 || ce.data[120] != 0x1234 || ce.data[121] != 0xBEEF {
		t.Fatalf("FC 23: ERROR %d, %04X %04X", cl.ERROR, se.data[50], ce.data[120:122])
	}
	// An address outside the map: the exception ILLEGAL DATA ADDRESS.
	cl.R_ADDR, cl.R_POINTS, cl.R_DATA_ADR = 250, 10, 0
	request(3)
	if cl.ERROR != 2 {
		t.Errorf("an address outside the map: ERROR %d, want 2", cl.ERROR)
	}
}

func TestVMAP(t *testing.T) {
	var vm VMAP
	vm[0] = network.VMAP_DATA{FC: 1 << 3, V_ADR: 1000, V_SIZE: 10, P_ADR: 50}
	m := MB_VMAP{VMAP: &vm}
	m.FC, m.V_ADR, m.V_CNT, m.SIZE = 3, 1004, 2, 256
	m.Execute(time.Time{})
	if m.ERROR != 0 || m.P_ADR != 54 {
		t.Errorf("FC 3 at 1004: P_ADR %d, ERROR %d", m.P_ADR, m.ERROR)
	}
	m.V_ADR = 1009
	m.Execute(time.Time{})
	if m.ERROR != 2 {
		t.Errorf("past the end of the map: ERROR %d", m.ERROR)
	}
	m.FC, m.V_ADR = 16, 1004
	m.Execute(time.Time{})
	if m.ERROR != 2 {
		t.Errorf("a function the map does not have: ERROR %d", m.ERROR)
	}
}
