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

package irtrans

import (
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

func TestDecodeReceive(t *testing.T) {
	var c network.IP_C
	var r network.NETWORK_BUFFER_SHORT
	line := "tv,vol+\r\n"
	for i := range line {
		r.BUFFER[i] = iec.BYTE(line[i])
	}
	r.SIZE = iec.UINT(len(line))
	c.MAILBOX[0] = 1 // a receive
	d := IRTRANS_DECODE{IP_C: &c, R_BUF: &r}
	d.Execute(time.Time{})
	if !d.CMD || d.ERROR || d.DEV != "tv" || d.KEY != "vol+" {
		t.Fatalf("decode: CMD %v ERROR %v %q %q", d.CMD, d.ERROR, d.DEV, d.KEY)
	}
	rcv := IRTRANS_RCV_4{DEV: &d.DEV, KEY: &d.KEY, DEV_CODE: "tv", KEY_CODE_0: "vol-", KEY_CODE_1: "vol+"}
	rcv.CMD = d.CMD
	rcv.Execute(time.Time{})
	if rcv.Q0 || !rcv.Q1 || rcv.Q2 {
		t.Errorf("receive: %v %v %v", rcv.Q0, rcv.Q1, rcv.Q2)
	}
	// The same receive again is no command.
	d.Execute(time.Time{})
	if d.CMD {
		t.Error("CMD for an old receive")
	}
}

func TestSend(t *testing.T) {
	c := network.IP_C{C_STATE: 255, C_MODE: 4}
	var s network.NETWORK_BUFFER_SHORT
	snd := IRTRANS_SND_4{sender: sender{IP_C: &c, S_BUF: &s, T_REPEAT: iec.TIME(time.Second)}, DEV_CODE: "tv", KEY_CODE_2: "mute"}
	now := time.Unix(1000, 0)
	snd.IN_2 = true
	snd.Execute(now)
	got := make([]byte, s.SIZE)
	for i := range got {
		got[i] = byte(s.BUFFER[i])
	}
	if string(got) != "snd tv,mute" || snd.KEY != 3 {
		t.Fatalf("send %q, KEY %d", got, snd.KEY)
	}
	// Sent: no repeat before T_REPEAT.
	s.SIZE = 0
	snd.Execute(now.Add(500 * time.Millisecond))
	if s.SIZE != 0 {
		t.Error("a repeat before T_REPEAT")
	}
	snd.Execute(now.Add(1600 * time.Millisecond))
	if s.SIZE == 0 {
		t.Error("no repeat after T_REPEAT")
	}
}
