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
	"time"

	"github.com/apiarytech/beebread/basic/logic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// TN_RECEIVE reads the key the client sent in R_BUF, and empties it: a
// printable character into BY_INPUT_ASCII_CODE of the screen, with
// BO_INPUT_ASCII_ISNUM for a digit, and into BY_INPUT_EXTEN_CODE the
// cursor keys 65 to 68 (up, down, right, left), 72 home and 75 end, F1 to
// F4 as 80 to 83, and backspace 8, tab 9, return 13 and escape 27.
type TN_RECEIVE struct {
	R_BUF         *network.NETWORK_BUFFER
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_RECEIVE) INIT() {}

// Execute runs the block once.
func (b *TN_RECEIVE) Execute(now time.Time) {
	r, sc := b.R_BUF, b.XUS_TN_SCREEN
	if r == nil || sc == nil {
		return
	}
	sc.BY_INPUT_ASCII_CODE, sc.BY_INPUT_EXTEN_CODE = 0, 0
	if r.SIZE == 0 {
		return
	}
	k1, k2, k3 := r.BUFFER[0], r.BUFFER[1], r.BUFFER[2]
	switch {
	case k1 == 27 && r.SIZE >= 3: // an escape sequence
		switch {
		case k2 == '[' && (k3 >= 65 && k3 <= 68 || k3 == 72 || k3 == 75):
			sc.BY_INPUT_EXTEN_CODE = k3
		case k2 == 'O' && k3 >= 80 && k3 <= 83:
			sc.BY_INPUT_EXTEN_CODE = k3
		}
	case k1 >= 32 && k1 <= 126:
		sc.BY_INPUT_ASCII_CODE = k1
		sc.BO_INPUT_ASCII_ISNUM = str.ISC_NUM(k1)
	case k1 == 8 || k1 == 9 || k1 == 13 || k1 == 27:
		sc.BY_INPUT_EXTEN_CODE = k1
	}
	r.SIZE = 0
}

// TN_SEND_ROWS sends the rows of the screen that changed, when connected
// and S_BUF is empty, as many as fit in S_BUF_SIZE bytes, then the cursor
// position. At a new connection or BO_CLEAR_SCREEN it clears the screen of
// the client in BY_CLEAR_SCREEN_ATTR and sends every row.
type TN_SEND_ROWS struct {
	IP_C          *network.IP_C
	S_BUF         *network.NETWORK_BUFFER
	XUS_TN_SCREEN *network.US_TN_SCREEN
	S_BUF_SIZE    iec.INT

	rowCount  iec.INT
	lastColor iec.BYTE
}

// INIT resets the block.
func (b *TN_SEND_ROWS) INIT() {
	*b = TN_SEND_ROWS{IP_C: b.IP_C, S_BUF: b.S_BUF, XUS_TN_SCREEN: b.XUS_TN_SCREEN}
}

// Execute runs the block once.
func (b *TN_SEND_ROWS) Execute(now time.Time) {
	c, s, sc := b.IP_C, b.S_BUF, b.XUS_TN_SCREEN
	if c == nil || s == nil || sc == nil || c.C_STATE < 127 || s.SIZE > 0 {
		return // not connected, or still sending
	}
	w := iec.INT(s.SIZE) - 1
	start := w
	put := func(bytes ...iec.BYTE) {
		for _, x := range bytes {
			w++
			s.BUFFER[w] = x
		}
	}
	// bcd puts the two digits of n.
	bcd := func(n iec.INT) {
		d := logic.INT_TO_BCDC(n)
		put(d>>4|'0', d&0x0F|'0')
	}

	if c.C_STATE == 254 || sc.BO_CLEAR_SCREEN { // connected now
		sc.BO_CLEAR_SCREEN = false
		for i := range sc.BYA_LINE_UPDATE {
			sc.BYA_LINE_UPDATE[i] = true
		}
		a := sc.BY_CLEAR_SCREEN_ATTR
		put(0x1B, '[', '0', ';', '3', a>>4|'0', ';', '4', a&0x07|'0', 'm') // the colors
		put(0x1B, '[', '2', 'J')                                           // clear
		put(0x1B, '[', '?', '7', 'l')                                      // no wrap
	}

	lastRow := iec.INT(0) // OSCAT goes back to 0, not to start, when no row fits
	stop := b.S_BUF_SIZE - 20
	full := false
	for range rows {
		o := offset(0, b.rowCount)
		if sc.BYA_LINE_UPDATE[b.rowCount] {
			put(0x1B, '[')
			bcd(b.rowCount + 1)
			put(';', '1', 'H')
			for col := 0; col < columns; col++ {
				color := sc.BYA_COLOR[o]
				if color != b.lastColor || col == 0 {
					put(0x1B, '[')
					if (color^b.lastColor)&0b1000_1000 != 0 || col == 0 {
						put('0', ';')
						if color&0b0000_1000 == 0 {
							put('1', ';') // bright
						}
						if color&0b1000_0000 != 0 {
							put('5', ';') // blinking
						}
					}
					put('3', color>>4&0b0111|'0', ';', '4', color&0b0111|'0', 'm')
					b.lastColor = color
				}
				if w > stop {
					full = true
					break
				}
				put(sc.BYA_CHAR[o])
				o++
			}
			if full { // send the complete rows only
				w = lastRow
				break
			}
			lastRow = w
			sc.BYA_LINE_UPDATE[b.rowCount] = false
		}
		b.rowCount++
		if b.rowCount > rows-1 {
			b.rowCount = 0
		}
	}
	if !full {
		b.rowCount = 0
	}

	if w != start { // the cursor
		put(0x1B, '[')
		bcd(sc.IN_CURSOR_Y + 1)
		put(';')
		bcd(sc.IN_CURSOR_X + 1)
		put('H')
	}
	s.SIZE = iec.UINT(w + 1)
}

// TN_SC_VIEWPORT shows the messages of the log XUS_LOG_CONTROL that the
// viewport XUS_LOG_VIEWPORT selects, at IIN_X, IIN_Y, IIN_WIDTH wide, at
// most every IIT_TIME: the color of a message, the low 4 bits of its
// option, selects a byte of IDW_ATTR_1 (0 to 3), else of IDW_ATTR_2.
type TN_SC_VIEWPORT struct {
	XUS_LOG_VIEWPORT *network.US_LOG_VIEWPORT
	XUS_LOG_CONTROL  *network.LOG_CONTROL
	XUS_TN_SCREEN    *network.US_TN_SCREEN
	IIN_X            iec.INT
	IIN_Y            iec.INT
	IIN_WIDTH        iec.INT
	IDW_ATTR_1       iec.DWORD
	IDW_ATTR_2       iec.DWORD
	ITI_TIME         iec.TIME

	write TN_SC_WRITE
	ton   timers.TON
}

// INIT resets the block.
func (b *TN_SC_VIEWPORT) INIT() {
	*b = TN_SC_VIEWPORT{XUS_LOG_VIEWPORT: b.XUS_LOG_VIEWPORT, XUS_LOG_CONTROL: b.XUS_LOG_CONTROL, XUS_TN_SCREEN: b.XUS_TN_SCREEN}
}

// Execute runs the block once.
func (b *TN_SC_VIEWPORT) Execute(now time.Time) {
	lv, lc := b.XUS_LOG_VIEWPORT, b.XUS_LOG_CONTROL
	if lv == nil || lc == nil || b.XUS_TN_SCREEN == nil {
		return
	}
	if lv.UPDATE && b.ton.Q {
		lv.UPDATE = false
		for n := iec.INT(1); n <= lv.COUNT; n++ {
			i := lv.LINE_ARRAY[n-1]
			color := iec.INT(lc.MSG_OPTION[i] & 0b1111)
			var attr iec.BYTE
			switch {
			case color <= 3:
				attr = logic.BYTE_OF_DWORD(b.IDW_ATTR_1, iec.BYTE(color))
			case color <= 7:
				// OSCAT: the bytes 4 to 7 of a DWORD
				attr = logic.BYTE_OF_DWORD(b.IDW_ATTR_2, iec.BYTE(color))
			default:
				attr = logic.BYTE_OF_DWORD(b.IDW_ATTR_1, 0)
			}
			w := &b.write
			w.IIN_Y, w.IIN_X, w.IBY_ATTR, w.XUS_TN_SCREEN = b.IIN_Y+n-1, b.IIN_X, attr, b.XUS_TN_SCREEN
			w.IST_STRING = str.FIX(lc.MSG[i], b.IIN_WIDTH, ' ', 0)
			w.Execute(now)
		}
		b.ton.IN, b.ton.PT = false, b.ITI_TIME
		b.ton.Execute(now)
	}
	b.ton.IN, b.ton.PT = true, b.ITI_TIME
	b.ton.Execute(now)
}

// TN_FRAMEWORK is a telnet server on PORT (23 if 0): it runs the menu bar
// US_TN_MENU and the input elements US_TN_INPUT_CONTROL on the screen
// US_TN_SCREEN, sends the rows that changed and reads the keys, and resets
// the connection 5 s after an error.
type TN_FRAMEWORK struct {
	US_TN_INPUT_CONTROL *network.US_TN_INPUT_CONTROL
	US_TN_SCREEN        *network.US_TN_SCREEN
	US_TN_MENU          *network.US_TN_MENU
	S_BUF               *network.NETWORK_BUFFER
	R_BUF               *network.NETWORK_BUFFER
	IP_C                *network.IP_C
	PORT                iec.WORD

	ipc     ip.IP_CONTROL
	input   TN_INPUT_CONTROL
	menuBar TN_INPUT_MENU_BAR
	receive TN_RECEIVE
	send    TN_SEND_ROWS
	t       timers.TON
}

// INIT resets the block.
func (b *TN_FRAMEWORK) INIT() {
	*b = TN_FRAMEWORK{US_TN_INPUT_CONTROL: b.US_TN_INPUT_CONTROL, US_TN_SCREEN: b.US_TN_SCREEN,
		US_TN_MENU: b.US_TN_MENU, S_BUF: b.S_BUF, R_BUF: b.R_BUF, IP_C: b.IP_C, PORT: b.PORT}
	b.ipc.INIT()
}

// Execute runs the block once.
func (b *TN_FRAMEWORK) Execute(now time.Time) {
	c, sc := b.IP_C, b.US_TN_SCREEN
	if c == nil || sc == nil || b.US_TN_INPUT_CONTROL == nil || b.US_TN_MENU == nil || b.S_BUF == nil || b.R_BUF == nil {
		return
	}
	if !c.C_ENABLE { // listen
		c.C_PORT = b.PORT
		if b.PORT == 0 {
			c.C_PORT = 23
		}
		c.C_IP = 0
		c.C_MODE = 4 // TCP, passive, the port
		c.TIME_RESET = true
		c.C_ENABLE = true
		c.R_OBSERVE = false
	}
	b.t.IN, b.t.PT = c.ERROR > 0, iec.TIME(5*time.Second)
	b.t.Execute(now)
	if b.t.Q {
		c.TIME_RESET = true // reset the error
	}

	b.menuBar.XUS_TN_MENU, b.menuBar.XUS_TN_SCREEN = b.US_TN_MENU, sc
	b.menuBar.Execute(now)
	b.input.XUS_TN_SCREEN, b.input.XUS_TN_INPUT_CONTROL = sc, b.US_TN_INPUT_CONTROL
	b.input.Execute(now)
	s := &b.send
	s.IP_C, s.S_BUF, s.XUS_TN_SCREEN, s.S_BUF_SIZE = c, b.S_BUF, sc, iec.INT(len(b.S_BUF.BUFFER))
	s.Execute(now)
	p := &b.ipc
	p.IP_C, p.S_BUF, p.R_BUF, p.IP, p.PORT, p.TIME_OUT = c, b.S_BUF, b.R_BUF, 0, 0, iec.TIME(2*time.Second)
	p.Execute(now)
	b.receive.R_BUF, b.receive.XUS_TN_SCREEN = b.R_BUF, sc
	b.receive.Execute(now)
}
