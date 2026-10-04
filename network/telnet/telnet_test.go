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

package telnet

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// row returns the characters of the row y of the screen.
func row(sc *network.US_TN_SCREEN, y int) string {
	var b strings.Builder
	for x := 0; x < columns; x++ {
		b.WriteByte(byte(sc.BYA_CHAR[y*columns+x]))
	}
	return b.String()
}

func bytesOf(b []iec.BYTE) []byte {
	out := make([]byte, len(b))
	for i, x := range b {
		out[i] = byte(x)
	}
	return out
}

func TestWrite(t *testing.T) {
	var sc network.US_TN_SCREEN
	w := TN_SC_WRITE{IIN_Y: 2, IIN_X: 78, IBY_ATTR: 0x70, IST_STRING: "ab", XUS_TN_SCREEN: &sc}
	w.Execute(time.Time{})
	if row(&sc, 2)[78:] != "ab" || sc.BYA_COLOR[2*80+79] != 0x70 || !sc.BYA_LINE_UPDATE[2] || sc.IN_EOS_OFFSET != 3*80 {
		t.Fatalf("write: %q, eos %d", row(&sc, 2)[78:], sc.IN_EOS_OFFSET)
	}
	w.IIN_Y, w.IIN_X = 3, 79 // past the end: nothing
	w.Execute(time.Time{})
	if sc.BYA_LINE_UPDATE[3] {
		t.Error("wrote past the end")
	}
	e := TN_SC_WRITE_EOS{IST_STRING: "cd", XUS_TN_SCREEN: &sc}
	e.Execute(time.Time{})
	if row(&sc, 3)[:2] != "cd" || sc.BYA_COLOR[3*80] != 0 {
		t.Errorf("eos: %q", row(&sc, 3)[:2])
	}
	c := TN_SC_WRITE_C{IIN_Y: 4, IIN_X: 0, IST_STRING: "ab", IIN_LENGTH: 6, IIN_OPTION: 2, XUS_TN_SCREEN: &sc}
	c.Execute(time.Time{})
	if row(&sc, 4)[:6] != "  ab  " {
		t.Errorf("centered: %q", row(&sc, 4)[:6])
	}
	ra := TN_SC_READ_ATTR{IIN_Y: 2, IIN_X: 78, XUS_TN_SCREEN: &sc}
	ra.Execute(time.Time{})
	rc := TN_SC_READ_CHAR{IIN_Y: 2, IIN_X: 78, XUS_TN_SCREEN: &sc}
	rc.Execute(time.Time{})
	if ra.OBY_ATTR != 0x70 || rc.OBY_CHAR != 'a' {
		t.Errorf("read %x %c", ra.OBY_ATTR, rc.OBY_CHAR)
	}
}

func TestBoxLineShadow(t *testing.T) {
	var sc network.US_TN_SCREEN
	b := TN_SC_BOX{IIN_X1: 0, IIN_Y1: 0, IIN_X2: 4, IIN_Y2: 3, IBY_FILL: '.', IBY_ATTR: 0x17, IIN_BORDER: 1, XUS_TN_SCREEN: &sc}
	b.Execute(time.Time{})
	want := []string{"\xda\xc4\xc4\xc4\xbf", "\xb3...\xb3", "\xb3...\xb3", "\xc0\xc4\xc4\xc4\xd9"}
	for y, w := range want {
		if row(&sc, y)[:5] != w {
			t.Errorf("box row %d: % x", y, row(&sc, y)[:5])
		}
	}
	// A line across the box joins its sides.
	l := TN_SC_LINE{IIN_X1: 0, IIN_Y1: 2, IIN_X2: 4, IIN_Y2: 2, IBY_BORDER: 1, IBY_ATTR: 0x17, XUS_TN_SCREEN: &sc}
	l.Execute(time.Time{})
	if row(&sc, 2)[:5] != "\xc3\xc4\xc4\xc4\xb4" {
		t.Errorf("line: % x", row(&sc, 2)[:5])
	}
	s := TN_SC_ADD_SHADOW{IIN_X1: 1, IIN_Y1: 1, IIN_X2: 5, IIN_Y2: 4, XUS_TN_SCREEN: &sc}
	s.Execute(time.Time{})
	if sc.BYA_COLOR[1*80+5] != 0x08 || sc.BYA_COLOR[4*80+1] != 0x08 || sc.BYA_COLOR[1*80+4] != 0x17 {
		t.Errorf("shadow %x %x %x", sc.BYA_COLOR[1*80+5], sc.BYA_COLOR[4*80+1], sc.BYA_COLOR[1*80+4])
	}
}

func TestAreaSaveRestore(t *testing.T) {
	var sc network.US_TN_SCREEN
	w := TN_SC_WRITE{IIN_Y: 1, IIN_X: 1, IBY_ATTR: 0x20, IST_STRING: "abc", XUS_TN_SCREEN: &sc}
	w.Execute(time.Time{})
	before := row(&sc, 1)[:6]
	s := TN_SC_AREA_SAVE{IIN_X1: 0, IIN_Y1: 0, IIN_X2: 5, IIN_Y2: 2, XUS_TN_SCREEN: &sc}
	s.Execute(time.Time{})
	f := TN_SC_FILL{IIN_X1: 0, IIN_Y1: 0, IIN_X2: 79, IIN_Y2: 23, IBY_CHAR: '#', XUS_TN_SCREEN: &sc}
	f.Execute(time.Time{})
	r := TN_SC_AREA_RESTORE{XUS_TN_SCREEN: &sc}
	r.Execute(time.Time{})
	if row(&sc, 1)[:6] != before || sc.BYA_COLOR[81] != 0x20 || row(&sc, 3)[:1] != "#" || sc.BYA_BACKUP[0] != 0 {
		t.Errorf("restored %q", row(&sc, 1)[:6])
	}
}

func TestReceive(t *testing.T) {
	var sc network.US_TN_SCREEN
	var r network.NETWORK_BUFFER
	rcv := TN_RECEIVE{R_BUF: &r, XUS_TN_SCREEN: &sc}
	for _, k := range []struct {
		in           string
		ascii, exten iec.BYTE
	}{
		{"\x1b[A", 0, keyUp}, {"\x1bOP", 0, 80}, {"\x1b[Z", 0, 0}, {"7", '7', 0}, {"\r", 0, keyReturn}, {"\x1b", 0, keyEscape},
	} {
		r.SIZE = iec.UINT(len(k.in))
		for i := range k.in {
			r.BUFFER[i] = iec.BYTE(k.in[i])
		}
		rcv.Execute(time.Time{})
		if sc.BY_INPUT_ASCII_CODE != k.ascii || sc.BY_INPUT_EXTEN_CODE != k.exten || r.SIZE != 0 {
			t.Errorf("%q: %d %d", k.in, sc.BY_INPUT_ASCII_CODE, sc.BY_INPUT_EXTEN_CODE)
		}
	}
	if !sc.BO_INPUT_ASCII_ISNUM {
		t.Error("7 is a digit")
	}
}

func TestSendRows(t *testing.T) {
	var sc network.US_TN_SCREEN
	var c network.IP_C
	var s network.NETWORK_BUFFER
	for i := range sc.BYA_CHAR {
		sc.BYA_CHAR[i], sc.BYA_COLOR[i] = ' ', 0x70
	}
	sc.BYA_CHAR[0] = 'X'
	sc.BY_CLEAR_SCREEN_ATTR = 0x70
	sc.IN_CURSOR_X, sc.IN_CURSOR_Y = 4, 9
	snd := TN_SEND_ROWS{IP_C: &c, S_BUF: &s, XUS_TN_SCREEN: &sc, S_BUF_SIZE: 300}
	snd.Execute(time.Time{})
	if s.SIZE != 0 {
		t.Fatal("sent while not connected")
	}
	c.C_STATE = 254
	snd.Execute(time.Time{})
	out := string(bytesOf(s.BUFFER[:s.SIZE]))
	// The clear screen, then two rows; the third does not fit in 300 - 20.
	head := "\x1b[0;37;40m\x1b[2J\x1b[?7l"
	row1 := "\x1b[01;1H\x1b[0;1;37;40mX" + strings.Repeat(" ", 79)
	row2 := "\x1b[02;1H\x1b[0;1;37;40m" + strings.Repeat(" ", 80)
	if out != head+row1+row2+"\x1b[10;05H" {
		t.Fatalf("sent %q", out)
	}
	if sc.BYA_LINE_UPDATE[0] || sc.BYA_LINE_UPDATE[1] || !sc.BYA_LINE_UPDATE[2] {
		t.Error("row updates")
	}
	// Then the rest, from row 3, without the clear screen.
	s.SIZE, c.C_STATE = 0, 255
	snd.Execute(time.Time{})
	out = string(bytesOf(s.BUFFER[:s.SIZE]))
	if !strings.HasPrefix(out, "\x1b[03;1H\x1b[0;1;37;40m ") {
		t.Errorf("next %q", out[:20])
	}
}

func TestEditLine(t *testing.T) {
	var sc network.US_TN_SCREEN
	d := network.US_TN_INPUT_CONTROL_DATA{IN_X: 10, IN_Y: 5, ST_INPUT_MASK: "  -  ", ST_INPUT_DATA: "  -  ", BO_FOCUS: true,
		ST_TITLE_STRING: "Code:", BY_ATTR_MF: 0x70, BO_UPDATE_ALL: true}
	e := TN_INPUT_EDIT_LINE{XUS_TN_SCREEN: &sc, XUS_TN_INPUT_CONTROL_DATA: &d}
	key := func(ascii, exten iec.BYTE) {
		d.BY_INPUT_ASCII_CODE, d.BY_INPUT_EXTEN_CODE = ascii, exten
		e.Execute(time.Time{})
	}
	key(0, 0)
	if row(&sc, 5)[5:15] != "Code:  -  " || d.IN_CURSOR_POS != 1 {
		t.Fatalf("shown %q, cursor %d", row(&sc, 5)[5:15], d.IN_CURSOR_POS)
	}
	for _, k := range "123" {
		key(iec.BYTE(k), 0)
	}
	if d.ST_INPUT_DATA != "12-3 " || d.IN_CURSOR_POS != 5 || d.IN_CURSOR_X != 14 {
		t.Fatalf("data %q, cursor %d", d.ST_INPUT_DATA, d.IN_CURSOR_POS)
	}
	key(0, keyBackspace)
	key(0, keyBackspace)
	if d.ST_INPUT_DATA != "12-  " || d.IN_CURSOR_POS != 2 {
		t.Fatalf("backspace: %q, cursor %d", d.ST_INPUT_DATA, d.IN_CURSOR_POS)
	}
	key('9', 0)
	key(0, keyReturn)
	if d.ST_INPUT_STRING != "19-" || !d.BO_INPUT_ENTERED || d.ST_INPUT_DATA != "  -  " || d.IN_CURSOR_POS != 1 {
		t.Errorf("entered %q, data %q, cursor %d", d.ST_INPUT_STRING, d.ST_INPUT_DATA, d.IN_CURSOR_POS)
	}
	d.BO_INPUT_HIDDEN = true
	key('7', 0)
	// Backspace clears the position under the cursor, then moves left; the
	// hidden line shows '*' for the 7 and the mask's '-'.
	if row(&sc, 5)[10:15] != "* -  " {
		t.Errorf("hidden %q", row(&sc, 5)[10:15])
	}
}

func TestInputControl(t *testing.T) {
	var sc network.US_TN_SCREEN
	var ic network.US_TN_INPUT_CONTROL
	ic.BO_ENABLE, ic.BO_RESET_FOKUS, ic.IN_COUNT = true, true, 2
	ic.IN_TOOLTIP_Y, ic.IN_TOOLTIP_SIZE = 23, 10
	ic.USA_TN_INPUT_CONTROL_DATA[0] = network.US_TN_INPUT_CONTROL_DATA{IN_TYPE: 2, IN_X: 0, IN_Y: 1,
		ST_INPUT_DATA: "red#green#blue", ST_INPUT_MASK: "     ", ST_INPUT_TOOLTIP: "color"}
	ic.USA_TN_INPUT_CONTROL_DATA[1] = network.US_TN_INPUT_CONTROL_DATA{IN_TYPE: 3, IN_X: 0, IN_Y: 2,
		ST_INPUT_DATA: "one#-#two", ST_INPUT_MASK: "   ", ST_INPUT_TOOLTIP: "number", BY_ATTR_MF: 0x70, BY_ATTR_OF: 0x07}
	b := TN_INPUT_CONTROL{XUS_TN_SCREEN: &sc, XUS_TN_INPUT_CONTROL: &ic}
	key := func(ascii, exten iec.BYTE) {
		sc.BY_INPUT_ASCII_CODE, sc.BY_INPUT_EXTEN_CODE = ascii, exten
		b.Execute(time.Time{})
	}
	key(0, 0)
	text := &ic.USA_TN_INPUT_CONTROL_DATA[0]
	if row(&sc, 1)[:5] != "red  " || row(&sc, 23)[:10] != "  color   " || !text.BO_FOCUS {
		t.Fatalf("shown %q, tooltip %q", row(&sc, 1)[:5], row(&sc, 23)[:10])
	}
	key(' ', 0)
	key(' ', 0)
	if text.ST_INPUT_STRING != "blue " || text.IN_SELECTED != 3 {
		t.Errorf("selected %q", text.ST_INPUT_STRING)
	}
	// Tab to the popup, open it, down (over the line), return.
	key(0, keyTab)
	pop := &ic.USA_TN_INPUT_CONTROL_DATA[1]
	if !pop.BO_FOCUS || text.BO_FOCUS || row(&sc, 23)[:10] != "  number  " {
		t.Fatal("focus")
	}
	key(0, keyReturn)
	if !sc.BO_MODAL_DIALOG {
		t.Fatal("no popup")
	}
	key(0, keyDown)
	key(0, keyReturn)
	key(0, 0)
	if sc.BO_MODAL_DIALOG || pop.IN_SELECTED != 3 || !pop.BO_INPUT_ENTERED || pop.ST_INPUT_STRING != "two" {
		t.Errorf("popup: modal %v, selected %d %q", sc.BO_MODAL_DIALOG, pop.IN_SELECTED, pop.ST_INPUT_STRING)
	}
	if row(&sc, 3)[:3] == "\xda\xc4\xc4" {
		t.Error("the popup was not removed")
	}
}

func TestMenuBar(t *testing.T) {
	var sc network.US_TN_SCREEN
	m := network.US_TN_MENU{ST_MENU_TEXT: "File#Edit%Open#Quit%Cut#Copy", BO_CREATE: true, BY_ATTR_MF: 0x70, BY_ATTR_OF: 0x07}
	b := TN_INPUT_MENU_BAR{XUS_TN_MENU: &m, XUS_TN_SCREEN: &sc}
	key := func(k iec.BYTE) {
		sc.BY_INPUT_EXTEN_CODE = k
		b.Execute(time.Time{})
	}
	key(0)
	if row(&sc, 0)[:12] != " File  Edit " || m.IN_STATE != 1 {
		t.Fatalf("bar %q", row(&sc, 0)[:12])
	}
	key(keyEscape)
	if m.IN_STATE != 2 || !sc.BO_MENUE_BAR_DIALOG || sc.BYA_COLOR[1] != 0x70 {
		t.Fatal("not open")
	}
	key(0)
	if !strings.Contains(row(&sc, 2), "Open") {
		t.Fatalf("popup %q", row(&sc, 2)[:10])
	}
	key(keyRight)
	key(0)
	key(keyDown)
	key(keyReturn)
	if m.IN_MENU_SELECTED != 22 || m.IN_STATE != 1 || sc.BO_MENUE_BAR_DIALOG {
		t.Errorf("selected %d, state %d", m.IN_MENU_SELECTED, m.IN_STATE)
	}
	key(0)
	if strings.Contains(row(&sc, 2), "Copy") {
		t.Error("the popup was not removed")
	}
}

func TestFramework(t *testing.T) {
	var ic network.US_TN_INPUT_CONTROL
	var sc network.US_TN_SCREEN
	var m network.US_TN_MENU
	var s, r network.NETWORK_BUFFER
	var c network.IP_C
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	f := TN_FRAMEWORK{US_TN_INPUT_CONTROL: &ic, US_TN_SCREEN: &sc, US_TN_MENU: &m, S_BUF: &s, R_BUF: &r, IP_C: &c, PORT: iec.WORD(port)}
	f.INIT()
	w := TN_SC_WRITE{IIN_Y: 0, IIN_X: 0, IBY_ATTR: 0x70, IST_STRING: "hello", XUS_TN_SCREEN: &sc}
	w.Execute(time.Time{})

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
			}
			f.Execute(time.Now())
			time.Sleep(time.Millisecond)
		}
	}()
	defer func() { close(stop); <-done }()

	var conn net.Conn
	for range 200 {
		if conn, err = net.Dial("tcp4", fmt.Sprintf("127.0.0.1:%d", port)); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	var got []byte
	buf := make([]byte, 4096)
	for !bytes.Contains(got, []byte("hello")) {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatalf("read: %v, got %q", err, got)
		}
		got = append(got, buf[:n]...)
	}
	if !bytes.HasPrefix(got, []byte("\x1b[0;30;40m\x1b[2J")) {
		t.Errorf("start %q", got[:16])
	}
}
