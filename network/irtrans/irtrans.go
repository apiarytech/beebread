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

// Package irtrans is the port of the OSCAT NETWORK blocks of an IRTrans
// infrared gateway: IRTRANS_SERVER keeps the connection of an IP_CONTROL2,
// IRTRANS_DECODE reads the 'device,key' lines the gateway sends of the
// remote controls it sees, IRTRANS_RCV_1/4/8 decode the keys of a device,
// and IRTRANS_SND_1/4/8 send keys, 'snd device,key'.
package irtrans

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// IRTRANS_SERVER sets the IP_CONTROL2 of IP_C up as a server, UDP or (with
// UDP_TCP false) TCP, for any client: R_ENABLE when a client is
// connected, S_ENABLE when one can be sent to. It resets an ERROR after
// 5 s.
type IRTRANS_SERVER struct {
	IP_C     *network.IP_C
	S_BUF    *network.NETWORK_BUFFER_SHORT
	R_BUF    *network.NETWORK_BUFFER_SHORT
	UDP_TCP  iec.BOOL
	S_ENABLE iec.BOOL
	R_ENABLE iec.BOOL
	ERROR    iec.DWORD

	t timers.TON
}

// INIT resets the block.
func (s *IRTRANS_SERVER) INIT() { *s = IRTRANS_SERVER{IP_C: s.IP_C, S_BUF: s.S_BUF, R_BUF: s.R_BUF} }

// Execute runs the block once.
func (s *IRTRANS_SERVER) Execute(now time.Time) {
	if s.IP_C == nil || s.S_BUF == nil || s.R_BUF == nil {
		return
	}
	c := s.IP_C
	if !c.C_ENABLE {
		c.C_PORT = 0 // the port of the IP_CONTROL2
		c.C_IP = 0
		c.C_MODE = SEL[iec.BYTE](s.UDP_TCP, 5, 4) // TCP or UDP server for any client
		c.C_ENABLE = true
		c.R_OBSERVE = false
		c.TIME_RESET = true
		s.S_BUF.SIZE, s.R_BUF.SIZE = 0, 0
	}
	s.R_ENABLE = c.C_STATE > 127
	s.S_ENABLE = s.R_ENABLE && (c.MAILBOX[0] > 0 || c.C_MODE != 5)
	s.ERROR = c.ERROR
	s.t.IN, s.t.PT = c.ERROR > 0, iec.TIME(5*time.Second)
	s.t.Execute(now)
	if s.t.Q {
		c.TIME_RESET = true
	}
}

// IRTRANS_DECODE reads the line 'device,key' CR LF that the gateway sent
// into R_BUF, once a receive: CMD for a scan, with DEV and KEY (80
// characters each at most), or ERROR for a line that is not one.
type IRTRANS_DECODE struct {
	IP_C  *network.IP_C
	R_BUF *network.NETWORK_BUFFER_SHORT
	CMD   iec.BOOL
	DEV   iec.STRING
	KEY   iec.STRING
	ERROR iec.BOOL

	z iec.BYTE
}

// INIT resets the block.
func (d *IRTRANS_DECODE) INIT() { *d = IRTRANS_DECODE{IP_C: d.IP_C, R_BUF: d.R_BUF} }

// Execute runs the block once.
func (d *IRTRANS_DECODE) Execute(now time.Time) {
	if d.IP_C == nil || d.R_BUF == nil {
		return
	}
	b := d.R_BUF.BUFFER[:]
	size := d.R_BUF.SIZE
	if size > 0 && d.IP_C.MAILBOX[0] != d.z {
		d.ERROR = true
		i, stop := 0, int(size)-2
		var dev, key []iec.BYTE
		// The device, up to the comma.
		for ; i < stop && b[i] != 44; i++ {
			dev = append(dev, b[i])
		}
		// The key, after it.
		for i++; i < stop; i++ {
			key = append(key, b[i])
		}
		d.DEV = network.STRING_N(STR(dev), 80)
		d.KEY = network.STRING_N(STR(key), 80)
		// A line that ends with CR LF, OSCAT checks one of them, and has
		// a device and a key.
		if !(b[i] != 13 && b[i+1] != 10) && LEN(d.DEV) > 0 && LEN(d.KEY) > 0 {
			d.CMD = true
			d.ERROR = false
		}
	} else {
		d.ERROR = false
		d.CMD = false
	}
	d.z = d.IP_C.MAILBOX[0]
}

// receive reports the keys of keys that match KEY, for the device
// DEV_CODE, while CMD.
func receive(cmd iec.BOOL, devCode, dev, key iec.STRING, keys []iec.STRING, q []*iec.BOOL) {
	decode := cmd && devCode == dev
	for i := range keys {
		*q[i] = decode && key == keys[i]
	}
}

// IRTRANS_RCV_1 sets Q while CMD has the key KEY_CODE of the device
// DEV_CODE, as IRTRANS_DECODE reads DEV and KEY.
type IRTRANS_RCV_1 struct {
	DEV_CODE iec.STRING
	KEY_CODE iec.STRING
	CMD      iec.BOOL
	DEV      *iec.STRING
	KEY      *iec.STRING
	Q        iec.BOOL
}

// INIT resets the block.
func (r *IRTRANS_RCV_1) INIT() { *r = IRTRANS_RCV_1{DEV: r.DEV, KEY: r.KEY} }

// Execute runs the block once.
func (r *IRTRANS_RCV_1) Execute(now time.Time) {
	if r.DEV == nil || r.KEY == nil {
		return
	}
	receive(r.CMD, r.DEV_CODE, *r.DEV, *r.KEY, []iec.STRING{r.KEY_CODE}, []*iec.BOOL{&r.Q})
}

// IRTRANS_RCV_4 sets Q0..Q3 for the keys KEY_CODE_0..3 of the device
// DEV_CODE; see IRTRANS_RCV_1.
type IRTRANS_RCV_4 struct {
	DEV_CODE                                       iec.STRING
	KEY_CODE_0, KEY_CODE_1, KEY_CODE_2, KEY_CODE_3 iec.STRING
	CMD                                            iec.BOOL
	DEV                                            *iec.STRING
	KEY                                            *iec.STRING
	Q0, Q1, Q2, Q3                                 iec.BOOL
}

// INIT resets the block.
func (r *IRTRANS_RCV_4) INIT() { *r = IRTRANS_RCV_4{DEV: r.DEV, KEY: r.KEY} }

// Execute runs the block once.
func (r *IRTRANS_RCV_4) Execute(now time.Time) {
	if r.DEV == nil || r.KEY == nil {
		return
	}
	receive(r.CMD, r.DEV_CODE, *r.DEV, *r.KEY,
		[]iec.STRING{r.KEY_CODE_0, r.KEY_CODE_1, r.KEY_CODE_2, r.KEY_CODE_3},
		[]*iec.BOOL{&r.Q0, &r.Q1, &r.Q2, &r.Q3})
}

// IRTRANS_RCV_8 sets Q0..Q7 for the keys KEY_CODE_0..7 of the device
// DEV_CODE; see IRTRANS_RCV_1.
type IRTRANS_RCV_8 struct {
	DEV_CODE                                       iec.STRING
	KEY_CODE_0, KEY_CODE_1, KEY_CODE_2, KEY_CODE_3 iec.STRING
	KEY_CODE_4, KEY_CODE_5, KEY_CODE_6, KEY_CODE_7 iec.STRING
	CMD                                            iec.BOOL
	DEV                                            *iec.STRING
	KEY                                            *iec.STRING
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7                 iec.BOOL
}

// INIT resets the block.
func (r *IRTRANS_RCV_8) INIT() { *r = IRTRANS_RCV_8{DEV: r.DEV, KEY: r.KEY} }

// Execute runs the block once.
func (r *IRTRANS_RCV_8) Execute(now time.Time) {
	if r.DEV == nil || r.KEY == nil {
		return
	}
	receive(r.CMD, r.DEV_CODE, *r.DEV, *r.KEY,
		[]iec.STRING{r.KEY_CODE_0, r.KEY_CODE_1, r.KEY_CODE_2, r.KEY_CODE_3, r.KEY_CODE_4, r.KEY_CODE_5, r.KEY_CODE_6, r.KEY_CODE_7},
		[]*iec.BOOL{&r.Q0, &r.Q1, &r.Q2, &r.Q3, &r.Q4, &r.Q5, &r.Q6, &r.Q7})
}

// sender is IRTRANS_SND_1, 4 and 8: the key it sends and when.
type sender struct {
	IP_C     *network.IP_C
	S_BUF    *network.NETWORK_BUFFER_SHORT
	T_REPEAT iec.TIME
	KEY      iec.BYTE

	skey iec.STRING
	size iec.INT
	str  iec.STRING
	t    timers.TON
	k    iec.BYTE
	lk   iec.BYTE
	d    iec.BOOL
}

// run sends 'snd devCode,key' for the first input of in that is on, as
// KEY its number from 1, again every T_REPEAT while it stays on.
func (s *sender) run(now time.Time, devCode iec.STRING, in []iec.BOOL, keys []iec.STRING) {
	if s.IP_C == nil || s.S_BUF == nil {
		return
	}
	c := s.IP_C
	s.d = c.C_STATE > 127 && (c.MAILBOX[0] > 0 || c.C_MODE != 5)
	if !s.d {
		return
	}
	s.t.IN, s.t.PT = s.lk == s.k, s.T_REPEAT
	s.t.Execute(now)
	if s.S_BUF.SIZE != 0 {
		return
	}
	on := -1
	for i := range in {
		if in[i] {
			on = i
			break
		}
	}
	if on < 0 {
		// no key
		s.lk = 0
		s.KEY = 0
		return
	}
	s.skey, s.k = keys[on], iec.BYTE(on+1)
	if s.lk != s.k || s.t.Q { // another key, or the time to repeat
		s.str = CONCAT("snd ", devCode, ",", s.skey)
		s.size = LEN(s.str)
		buffer.STRING_TO_BUFFER_(s.str, 0, s.S_BUF.BUFFER[:], iec.UINT(s.size))
		s.S_BUF.SIZE = iec.UINT(s.size)
		s.KEY = s.k
		s.lk = s.k
		s.t.IN = false
		s.t.Execute(now)
	}
}

// IRTRANS_SND_1 sends the key KEY_CODE of the device DEV_CODE while IN,
// through the IP_CONTROL2 of IP_C, again every T_REPEAT; KEY is 1 when it
// sends.
type IRTRANS_SND_1 struct {
	sender
	DEV_CODE iec.STRING
	KEY_CODE iec.STRING
	IN       iec.BOOL
}

// INIT resets the block.
func (s *IRTRANS_SND_1) INIT() { *s = IRTRANS_SND_1{sender: sender{IP_C: s.IP_C, S_BUF: s.S_BUF}} }

// Execute runs the block once.
func (s *IRTRANS_SND_1) Execute(now time.Time) {
	s.run(now, s.DEV_CODE, []iec.BOOL{s.IN}, []iec.STRING{s.KEY_CODE})
}

// IRTRANS_SND_4 sends the key KEY_CODE_n of the device DEV_CODE of the
// first input IN_n that is on; KEY is n+1. See IRTRANS_SND_1.
type IRTRANS_SND_4 struct {
	sender
	DEV_CODE                                       iec.STRING
	KEY_CODE_0, KEY_CODE_1, KEY_CODE_2, KEY_CODE_3 iec.STRING
	IN_0, IN_1, IN_2, IN_3                         iec.BOOL
}

// INIT resets the block.
func (s *IRTRANS_SND_4) INIT() { *s = IRTRANS_SND_4{sender: sender{IP_C: s.IP_C, S_BUF: s.S_BUF}} }

// Execute runs the block once.
func (s *IRTRANS_SND_4) Execute(now time.Time) {
	s.run(now, s.DEV_CODE, []iec.BOOL{s.IN_0, s.IN_1, s.IN_2, s.IN_3},
		[]iec.STRING{s.KEY_CODE_0, s.KEY_CODE_1, s.KEY_CODE_2, s.KEY_CODE_3})
}

// IRTRANS_SND_8 is IRTRANS_SND_4 with 8 keys.
type IRTRANS_SND_8 struct {
	sender
	DEV_CODE                                       iec.STRING
	KEY_CODE_0, KEY_CODE_1, KEY_CODE_2, KEY_CODE_3 iec.STRING
	KEY_CODE_4, KEY_CODE_5, KEY_CODE_6, KEY_CODE_7 iec.STRING
	IN_0, IN_1, IN_2, IN_3, IN_4, IN_5, IN_6, IN_7 iec.BOOL
}

// INIT resets the block.
func (s *IRTRANS_SND_8) INIT() { *s = IRTRANS_SND_8{sender: sender{IP_C: s.IP_C, S_BUF: s.S_BUF}} }

// Execute runs the block once.
func (s *IRTRANS_SND_8) Execute(now time.Time) {
	s.run(now, s.DEV_CODE, []iec.BOOL{s.IN_0, s.IN_1, s.IN_2, s.IN_3, s.IN_4, s.IN_5, s.IN_6, s.IN_7},
		[]iec.STRING{s.KEY_CODE_0, s.KEY_CODE_1, s.KEY_CODE_2, s.KEY_CODE_3, s.KEY_CODE_4, s.KEY_CODE_5, s.KEY_CODE_6, s.KEY_CODE_7})
}
