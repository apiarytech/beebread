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

// Package netvar is the port of the OSCAT NETWORK network variables: two
// PLCs, a master and a slave, exchange the values of their NET_VAR blocks
// in a cycle, through NET_VAR_CONTROL and an IP_CONTROL. Each block has an
// ID, in the order the blocks first run, and writes its record, ID, type
// and value, into the data the master or slave sends, then reads the
// record of the same ID from the data it receives. The blocks must run in
// the same order on both sides; a record that does not match sets
// X.ERROR_ID to the ID of the block.
package netvar

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/basic/logic"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// The states of NET_VAR_DATA.
const (
	stateIdle  iec.BYTE = 0
	stateWrite iec.BYTE = 1 // the blocks write their records
	stateRead  iec.BYTE = 2 // the blocks read theirs
)

// NET_VAR_CONTROL exchanges the network variables of X with the remote
// PLC REMOTE_IP4:REMOTE_PORT, on a rising edge of ACTIVATE, over TCP or
// UDP: the MASTER sends its records every SCAN_TIME, the slave answers
// with its own. RUN while it runs; ERROR is the ID of a block whose record
// did not match, or the connection's error, after which it starts again.
type NET_VAR_CONTROL struct {
	ACTIVATE    iec.BOOL
	MASTER      iec.BOOL
	UDP         iec.BOOL
	REMOTE_IP4  iec.DWORD
	REMOTE_PORT iec.WORD // default 10000
	SCAN_TIME   iec.TIME // default T#1s
	WATCHDOG    iec.TIME // default T#2s
	RUN         iec.BOOL
	ERROR       iec.DWORD
	X           *network.NET_VAR_DATA

	ipc          ip.IP_CONTROL
	ipC          network.IP_C
	step         iec.INT
	activateLast iec.BOOL
	tscan        timers.TON
	reset        iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (n *NET_VAR_CONTROL) INIT() {
	*n = NET_VAR_CONTROL{X: n.X, REMOTE_PORT: 10000, SCAN_TIME: iec.TIME(time.Second), WATCHDOG: iec.TIME(2 * time.Second)}
}

// Execute runs the block once.
func (n *NET_VAR_CONTROL) Execute(now time.Time) {
	if n.X == nil {
		return
	}
	x, c := n.X, &n.ipC
	if x.ERROR_ID > 0 {
		// a record that does not match
		x.STATE = stateIdle
		n.ERROR = iec.DWORD(x.ERROR_ID)
		x.ERROR_ID = 0
		n.step = 999
	}
	switch n.step {
	case 0:
		if n.ACTIVATE && !n.activateLast || n.reset {
			// 0 TCP client, 1 UDP client, 2 TCP server, 3 UDP server
			if n.MASTER {
				c.C_MODE = 0
				n.step = 100
			} else {
				c.C_MODE = 2
				n.step = 200
			}
			if n.UDP {
				c.C_MODE |= 1
			}
			c.C_PORT = n.REMOTE_PORT
			c.C_IP = n.REMOTE_IP4
			c.TIME_RESET = true
			c.C_ENABLE = true
			c.R_OBSERVE = true
			x.S_BUF.SIZE, x.R_BUF.SIZE = 0, 0
			x.CYCLE = 0
			x.ERROR_ID = 0
			n.reset = false
		}

	// The master.
	case 100:
		if c.C_STATE > 127 { // connected: the blocks write
			x.STATE = stateWrite
			x.BUF_SIZE = iec.UINT(len(x.S_BUF.BUFFER))
			x.S_BUF.BUFFER[0] = iec.BYTE(x.CYCLE)
			x.INDEX = 1
			n.step = 110
		}
	case 110: // send
		x.STATE = stateIdle
		x.S_BUF.SIZE = iec.UINT(x.INDEX)
		n.step = 120
	case 120:
		if x.S_BUF.SIZE == 0 && x.R_BUF.SIZE >= 1 { // the answer: the blocks read
			x.STATE = stateRead
			x.BUF_SIZE = x.R_BUF.SIZE
			x.INDEX = 1
			x.ERROR_ID = 0
			x.CYCLE++
			n.ERROR = 0
			n.RUN = true
			n.step = 130
		}
	case 130:
		x.STATE = stateIdle
		x.R_BUF.SIZE = 0
		if !n.ACTIVATE {
			n.step = 999
		} else if n.tscan.Q {
			n.step = 100
		}

	// The slave.
	case 200:
		if x.R_BUF.SIZE >= 1 { // the blocks read
			x.STATE = stateRead
			x.BUF_SIZE = x.R_BUF.SIZE
			x.INDEX = 1
			n.step = 210
		}
	case 210: // the blocks write the answer, with the master's cycle
		x.STATE = stateWrite
		x.BUF_SIZE = iec.UINT(len(x.R_BUF.BUFFER))
		x.S_BUF.BUFFER[0] = x.R_BUF.BUFFER[0]
		x.INDEX = 1
		x.R_BUF.SIZE = 0
		n.step = 220
	case 220:
		x.STATE = stateIdle
		x.S_BUF.SIZE = iec.UINT(x.INDEX)
		n.step = 230
	case 230:
		if x.S_BUF.SIZE == 0 {
			n.ERROR = 0
			n.RUN = true
			if !n.ACTIVATE {
				n.step = 999
			} else {
				n.step = 200
			}
		}

	case 999:
		c.C_ENABLE = false
		n.RUN = false
		n.step = 0
	}

	n.ipc.IP_C, n.ipc.S_BUF, n.ipc.R_BUF = c, &x.S_BUF, &x.R_BUF
	n.ipc.IP, n.ipc.PORT, n.ipc.TIME_OUT = 0, 0, n.WATCHDOG
	n.ipc.Execute(now)
	if c.ERROR > 0 {
		// start again
		n.ERROR = c.ERROR
		n.RUN = false
		n.reset = true
		n.step = 0
	}
	n.activateLast = n.ACTIVATE
	n.tscan.IN, n.tscan.PT = n.step != 110, n.SCAN_TIME
	n.tscan.Execute(now)
}

// member is what the NET_VAR blocks share: their ID and the record they
// write and read.
type member struct {
	X  *network.NET_VAR_DATA
	ID iec.BYTE

	init iec.BOOL
}

// start gives the block its ID, the first time.
func (m *member) start() {
	if !m.init {
		m.init = true
		m.X.ID_MAX++
		m.ID = iec.BYTE(m.X.ID_MAX)
	}
}

// active reports whether the blocks write or read.
func (m *member) active() bool { return m.X.STATE > 0 && m.X.ERROR_ID == 0 }

// head reads the head of the record at index, checking its ID and type,
// and returns the index after it, or false for a record that does not
// match.
func (m *member) head(index iec.INT, typ iec.BYTE) (iec.INT, bool) {
	r := m.X.R_BUF.BUFFER[:]
	if r[index] != m.ID || r[index+1] != typ {
		m.X.ERROR_ID = m.ID
		return index, false
	}
	return index + 2, true
}

// dword8 writes or reads a record of 8 DWORDs of the type typ, little
// endian.
func (m *member) dword8(typ iec.BYTE, in [8]iec.DWORD, out *[8]iec.DWORD) {
	x := m.X
	index := x.INDEX
	if iec.UINT(index+34) > x.BUF_SIZE {
		x.ERROR_ID = m.ID
		return
	}
	switch x.STATE {
	case stateWrite:
		s := x.S_BUF.BUFFER[:]
		s[index], s[index+1] = m.ID, typ
		index += 2
		for _, d := range in {
			for i := 0; i < 4; i++ {
				s[index] = iec.BYTE(d)
				d = ROR(d, 8)
				index++
			}
		}
	case stateRead:
		var ok bool
		if index, ok = m.head(index, typ); !ok {
			break
		}
		r := x.R_BUF.BUFFER[:]
		for i2 := range out {
			var d iec.DWORD
			for i := 0; i < 4; i++ {
				d = ROR(d|iec.DWORD(r[index]), 8)
				index++
			}
			out[i2] = d
		}
	}
	x.INDEX = index
}

// NET_VAR_BOOL8 exchanges 8 BOOLs: IN1..IN8 to the remote's OUT1..OUT8.
type NET_VAR_BOOL8 struct {
	member
	IN1, IN2, IN3, IN4, IN5, IN6, IN7, IN8         iec.BOOL
	OUT1, OUT2, OUT3, OUT4, OUT5, OUT6, OUT7, OUT8 iec.BOOL
}

// INIT resets the block.
func (n *NET_VAR_BOOL8) INIT() { *n = NET_VAR_BOOL8{member: member{X: n.X}} }

// Execute runs the block once.
func (n *NET_VAR_BOOL8) Execute(now time.Time) {
	if n.X == nil {
		return
	}
	n.start()
	if !n.active() {
		return
	}
	x := n.X
	index := x.INDEX
	if iec.UINT(index+3) > x.BUF_SIZE {
		x.ERROR_ID = n.ID
		return
	}
	switch x.STATE {
	case stateWrite:
		s := x.S_BUF.BUFFER[:]
		s[index], s[index+1] = n.ID, 4 // BOOL8
		s[index+2] = logic.BYTE_OF_BIT(n.IN1, n.IN2, n.IN3, n.IN4, n.IN5, n.IN6, n.IN7, n.IN8)
		index += 3
	case stateRead:
		var ok bool
		if index, ok = n.head(index, 4); !ok {
			break
		}
		b := x.R_BUF.BUFFER[index]
		out := []*iec.BOOL{&n.OUT1, &n.OUT2, &n.OUT3, &n.OUT4, &n.OUT5, &n.OUT6, &n.OUT7, &n.OUT8}
		for i, o := range out {
			*o = iec.BOOL(BIT(b, i))
		}
		index++
	}
	x.INDEX = index
}

// NET_VAR_DWORD8 exchanges 8 DWORDs: IN1..IN8 to the remote's OUT1..OUT8.
type NET_VAR_DWORD8 struct {
	member
	IN1, IN2, IN3, IN4, IN5, IN6, IN7, IN8         iec.DWORD
	OUT1, OUT2, OUT3, OUT4, OUT5, OUT6, OUT7, OUT8 iec.DWORD
}

// INIT resets the block.
func (n *NET_VAR_DWORD8) INIT() { *n = NET_VAR_DWORD8{member: member{X: n.X}} }

// Execute runs the block once.
func (n *NET_VAR_DWORD8) Execute(now time.Time) {
	if n.X == nil {
		return
	}
	n.start()
	if !n.active() {
		return
	}
	out := [8]iec.DWORD{n.OUT1, n.OUT2, n.OUT3, n.OUT4, n.OUT5, n.OUT6, n.OUT7, n.OUT8}
	n.dword8(13, [8]iec.DWORD{n.IN1, n.IN2, n.IN3, n.IN4, n.IN5, n.IN6, n.IN7, n.IN8}, &out)
	n.OUT1, n.OUT2, n.OUT3, n.OUT4, n.OUT5, n.OUT6, n.OUT7, n.OUT8 = out[0], out[1], out[2], out[3], out[4], out[5], out[6], out[7]
}

// NET_VAR_REAL8 exchanges 8 REALs: IN1..IN8 to the remote's OUT1..OUT8.
type NET_VAR_REAL8 struct {
	member
	IN1, IN2, IN3, IN4, IN5, IN6, IN7, IN8         iec.REAL
	OUT1, OUT2, OUT3, OUT4, OUT5, OUT6, OUT7, OUT8 iec.REAL
}

// INIT resets the block.
func (n *NET_VAR_REAL8) INIT() { *n = NET_VAR_REAL8{member: member{X: n.X}} }

// Execute runs the block once.
func (n *NET_VAR_REAL8) Execute(now time.Time) {
	if n.X == nil {
		return
	}
	n.start()
	if !n.active() {
		return
	}
	in := [8]iec.REAL{n.IN1, n.IN2, n.IN3, n.IN4, n.IN5, n.IN6, n.IN7, n.IN8}
	outR := []*iec.REAL{&n.OUT1, &n.OUT2, &n.OUT3, &n.OUT4, &n.OUT5, &n.OUT6, &n.OUT7, &n.OUT8}
	var inD, out [8]iec.DWORD
	for i := range in {
		inD[i] = logic.REAL_TO_DW(in[i])
		out[i] = logic.REAL_TO_DW(*outR[i])
	}
	n.dword8(15, inD, &out)
	for i := range out {
		*outR[i] = logic.DW_TO_REAL(out[i])
	}
}

// NET_VAR_X8 exchanges 2 REALs, 2 DINTs, 2 UDINTs and 2 DWORDs.
type NET_VAR_X8 struct {
	member
	IN_REAL1, IN_REAL2     iec.REAL
	IN_DINT1, IN_DINT2     iec.DINT
	IN_UDINT1, IN_UDINT2   iec.UDINT
	IN_DWORD1, IN_DWORD2   iec.DWORD
	OUT_REAL1, OUT_REAL2   iec.REAL
	OUT_DINT1, OUT_DINT2   iec.DINT
	OUT_UDINT1, OUT_UDINT2 iec.UDINT
	OUT_DWORD1, OUT_DWORD2 iec.DWORD
}

// INIT resets the block.
func (n *NET_VAR_X8) INIT() { *n = NET_VAR_X8{member: member{X: n.X}} }

// Execute runs the block once.
func (n *NET_VAR_X8) Execute(now time.Time) {
	if n.X == nil {
		return
	}
	n.start()
	if !n.active() {
		return
	}
	in := [8]iec.DWORD{
		logic.REAL_TO_DW(n.IN_REAL1), logic.REAL_TO_DW(n.IN_REAL2),
		iec.DWORD(n.IN_DINT1), iec.DWORD(n.IN_DINT2),
		iec.DWORD(n.IN_UDINT1), iec.DWORD(n.IN_UDINT2),
		n.IN_DWORD1, n.IN_DWORD2,
	}
	out := [8]iec.DWORD{
		logic.REAL_TO_DW(n.OUT_REAL1), logic.REAL_TO_DW(n.OUT_REAL2),
		iec.DWORD(n.OUT_DINT1), iec.DWORD(n.OUT_DINT2),
		iec.DWORD(n.OUT_UDINT1), iec.DWORD(n.OUT_UDINT2),
		n.OUT_DWORD1, n.OUT_DWORD2,
	}
	n.dword8(21, in, &out)
	n.OUT_REAL1, n.OUT_REAL2 = logic.DW_TO_REAL(out[0]), logic.DW_TO_REAL(out[1])
	n.OUT_DINT1, n.OUT_DINT2 = iec.DINT(out[2]), iec.DINT(out[3])
	n.OUT_UDINT1, n.OUT_UDINT2 = iec.UDINT(out[4]), iec.UDINT(out[5])
	n.OUT_DWORD1, n.OUT_DWORD2 = out[6], out[7]
}

// NET_VAR_STRING exchanges a string: IN to the remote's OUT.
type NET_VAR_STRING struct {
	member
	IN  *iec.STRING // STRING(STRING_LENGTH)
	OUT *iec.STRING // STRING(STRING_LENGTH)
}

// INIT resets the block.
func (n *NET_VAR_STRING) INIT() { *n = NET_VAR_STRING{member: member{X: n.X}, IN: n.IN, OUT: n.OUT} }

// Execute runs the block once.
func (n *NET_VAR_STRING) Execute(now time.Time) {
	if n.X == nil || n.IN == nil || n.OUT == nil {
		return
	}
	n.start()
	if !n.active() {
		return
	}
	x := n.X
	index := x.INDEX
	switch x.STATE {
	case stateWrite:
		l := LEN(*n.IN)
		if iec.UINT(index+3+l) > x.BUF_SIZE {
			x.ERROR_ID = n.ID
			return
		}
		s := x.S_BUF.BUFFER[:]
		s[index], s[index+1], s[index+2] = n.ID, 6, iec.BYTE(l) // STRING
		index += 3
		buffer.STRING_TO_BUFFER_(*n.IN, index, s, x.BUF_SIZE)
		index += l
	case stateRead:
		var ok bool
		if index, ok = n.head(index, 6); !ok {
			break
		}
		l := iec.INT(x.R_BUF.BUFFER[index])
		start := index + 1
		index = start + l
		if l > 0 {
			*n.OUT = buffer.BUFFER_TO_STRING(x.R_BUF.BUFFER[:], x.BUF_SIZE, iec.UINT(start), iec.UINT(index-1))
		} else {
			*n.OUT = ""
		}
	}
	x.INDEX = index
}

// NET_VAR_BUFFER exchanges a buffer of 64 bytes: BUF_IN to the remote's
// BUF_OUT.
type NET_VAR_BUFFER struct {
	member
	BUF_IN  *[64]iec.BYTE // ARRAY[1..64]
	BUF_OUT *[64]iec.BYTE // ARRAY[1..64]
}

// INIT resets the block.
func (n *NET_VAR_BUFFER) INIT() {
	*n = NET_VAR_BUFFER{member: member{X: n.X}, BUF_IN: n.BUF_IN, BUF_OUT: n.BUF_OUT}
}

// Execute runs the block once.
func (n *NET_VAR_BUFFER) Execute(now time.Time) {
	if n.X == nil || n.BUF_IN == nil || n.BUF_OUT == nil {
		return
	}
	n.start()
	if !n.active() {
		return
	}
	x := n.X
	index := x.INDEX
	switch x.STATE {
	case stateWrite:
		size := iec.INT(len(n.BUF_IN))
		if iec.UINT(index+3+size) > x.BUF_SIZE {
			x.ERROR_ID = n.ID
			break
		}
		s := x.S_BUF.BUFFER[:]
		s[index], s[index+1], s[index+2] = n.ID, 20, iec.BYTE(size) // BUFFER
		index += 3
		copy(s[index:index+size], n.BUF_IN[:])
		index += size
	case stateRead:
		var ok bool
		if index, ok = n.head(index, 20); !ok {
			break
		}
		size := iec.INT(x.R_BUF.BUFFER[index])
		index++
		for i := iec.INT(0); i < size && int(i) < len(n.BUF_OUT); i++ {
			n.BUF_OUT[i] = x.R_BUF.BUFFER[index]
			index++
		}
	}
	x.INDEX = index
}
