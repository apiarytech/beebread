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
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/basic/logic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// ntpEpoch is the seconds from 1900-01-01, the start of NTP time, to
// 1970-01-01.
const ntpEpoch iec.DWORD = 2208988800

// SNTP_CLIENT reads the time UDT, in UTC, and its milliseconds XMS from the
// time server IP4, on a rising edge of ACTIVATE, adding half the time of
// the round trip: DONE_P for a scan, else ERROR.
type SNTP_CLIENT struct {
	client
	IP4      iec.DWORD
	ACTIVATE iec.BOOL
	ERROR    iec.DWORD
	DONE_P   iec.BOOL
	UDT      iec.DT
	XMS      iec.INT

	last         iec.DWORD
	state        iec.INT
	tx           iec.DWORD
	activateLast iec.BOOL
	ip4Stored    iec.DWORD
}

// INIT resets the block.
func (s *SNTP_CLIENT) INIT() {
	*s = SNTP_CLIENT{client: client{IP_C: s.IP_C, S_BUF: s.S_BUF, R_BUF: s.R_BUF}}
}

// Execute runs the block once.
func (s *SNTP_CLIENT) Execute(now time.Time) {
	if !s.ok() {
		return
	}
	c, sb, rb := s.IP_C, s.S_BUF, s.R_BUF
	s.tx = PLC_MS(now)
	switch s.state {
	case 0: // wait for a rising edge
		s.DONE_P = false
		if s.ACTIVATE && !s.activateLast {
			s.ip4Stored = s.IP4
			s.state = 10
			s.ERROR = 0
			s.ipState = 1
		}
	case 10:
		if s.ipState == 3 {
			sb.BUFFER[0] = 0x1B // version 3, client
			for i := 1; i <= 47; i++ {
				sb.BUFFER[i] = 0
			}
			s.last = s.tx
			s.state = 30
			c.C_PORT = 123
			c.C_IP = s.ip4Stored
			c.C_MODE = 1 // UDP client
			c.C_ENABLE = true
			c.TIME_RESET = true
			c.R_OBSERVE = true
			sb.SIZE = 48
			rb.SIZE = 0
		}
	case 30:
		if c.ERROR != 0 {
			s.ERROR = c.ERROR
			sb.SIZE, rb.SIZE = 0, 0
			s.state = 0
		} else if sb.SIZE == 0 && rb.SIZE >= 48 {
			// The transmit time, from 1900.
			s.UDT = DWORD_TO_DT(logic.DWORD_OF_BYTE(rb.BUFFER[40], rb.BUFFER[41], rb.BUFFER[42], rb.BUFFER[43]) - ntpEpoch)
			// The milliseconds, from the high word of the fraction.
			s.XMS = iec.INT(SHR(iec.DWORD(iec.DINT(logic.WORD_OF_BYTE(rb.BUFFER[44], rb.BUFFER[45]))*1000), 16))
			// Half the round trip.
			s.XMS += iec.INT(SHR(s.tx-s.last, 1))
			s.UDT = DWORD_TO_DT(DT_TO_DWORD(s.UDT) + iec.DWORD(s.XMS/1000))
			s.XMS %= 1000
			s.state = 0
			s.DONE_P = true
		}
		if s.state == 0 {
			c.R_OBSERVE = false
			s.ipState = 4
		}
	}
	s.activateLast = s.ACTIVATE
	s.runFIFO(now, s.IP_C)
}

// SNTP_SERVER is a time server, while ENABLE, of the time UDT, in UTC, and
// its milliseconds XMS, with the stratum STRATUM, on the UDP port 123.
type SNTP_SERVER struct {
	client
	ENABLE  iec.BOOL
	STRATUM iec.BYTE // default 1
	UDT     iec.DT
	XMS     iec.INT

	timestampInt iec.DWORD
	timestampSek iec.DWORD
	state        iec.INT
	t            timers.TON
}

// INIT resets the block and sets STRATUM to its initial value.
func (s *SNTP_SERVER) INIT() {
	*s = SNTP_SERVER{client: client{IP_C: s.IP_C, S_BUF: s.S_BUF, R_BUF: s.R_BUF}, STRATUM: 1}
}

// Execute runs the block once.
func (s *SNTP_SERVER) Execute(now time.Time) {
	if !s.ok() {
		return
	}
	c, sb, rb := s.IP_C, s.S_BUF, s.R_BUF
	switch s.state {
	case 0:
		if s.ENABLE {
			s.state = 10
			s.ipState = 1
		}
	case 10:
		if s.ipState == 3 {
			c.C_PORT = 123
			c.C_IP = 0
			c.C_MODE = 5 // UDP server for any client
			c.C_ENABLE = true
			c.TIME_RESET = true
			c.R_OBSERVE = false
			sb.SIZE, rb.SIZE = 0, 0
			s.state = 20
		}
	case 20:
		if rb.SIZE > 0 {
			// Pad a short request with 0, and clear the answer.
			for i := 0; i <= 47; i++ {
				if iec.UINT(i) > rb.SIZE-1 {
					rb.BUFFER[i] = 0
				}
				sb.BUFFER[i] = 0
			}
			rb.SIZE = 0
			sb.BUFFER[0] = 0x1C // no warning, version 3, server
			sb.BUFFER[1] = s.STRATUM
			sb.BUFFER[2] = 10   // the poll interval, 2^10 s
			sb.BUFFER[3] = 0xFB // the precision, -5
			s.timestampInt = DT_TO_DWORD(s.UDT) + ntpEpoch
			// The high word of the fraction: XMS * 65536 / 1000.
			s.timestampSek = iec.DWORD(iec.DINT(SHL(iec.DWORD(s.XMS), 16)) / 1000)
			// The reference, receive and transmit times.
			for _, at := range []int{16, 32, 40} {
				sb.BUFFER[at] = logic.BYTE_OF_DWORD(s.timestampInt, 3)
				sb.BUFFER[at+1] = logic.BYTE_OF_DWORD(s.timestampInt, 2)
				sb.BUFFER[at+2] = logic.BYTE_OF_DWORD(s.timestampInt, 1)
				sb.BUFFER[at+3] = logic.BYTE_OF_DWORD(s.timestampInt, 0)
				sb.BUFFER[at+4] = logic.BYTE_OF_DWORD(s.timestampSek, 1)
				sb.BUFFER[at+5] = logic.BYTE_OF_DWORD(s.timestampSek, 0)
			}
			// The originate time is the client's transmit time.
			copy(sb.BUFFER[24:32], rb.BUFFER[40:48])
			sb.SIZE = 48
		} else if !s.ENABLE && sb.SIZE == 0 {
			s.ipState = 4
			s.state = 0
		}
	}
	// Reset an error after 5 s.
	s.t.IN, s.t.PT = c.ERROR > 0, iec.TIME(5*time.Second)
	s.t.Execute(now)
	if s.t.Q {
		c.TIME_RESET = true
	}
	s.runFIFO(now, s.IP_C)
}

// SYS_LOG sends the syslog message MESSAGE, with the facility FACILITY,
// the severity SEVERITY, the time LDT, HOSTNAME and TAG, to the syslog
// server SERVER_IP4, on a rising edge of ACTIVATE: DONE_P for a scan, else
// ERROR. The port is PORT, or 514 for UDP and 1468 for TCP. The bits of
// OPTION: 0 no priority, 1 no header, 2 a CR LF at the end, 3 TCP.
type SYS_LOG struct {
	client
	ACTIVATE   iec.BOOL
	LDT        iec.DT
	SERVER_IP4 iec.DWORD
	PORT       iec.WORD
	FACILITY   iec.BYTE
	SEVERITY   iec.BYTE
	TAG        iec.STRING // STRING(32)
	HOSTNAME   iec.STRING
	MESSAGE    iec.STRING // STRING(STRING_LENGTH)
	OPTION     iec.BYTE
	DONE_P     iec.BOOL
	ERROR      iec.DWORD

	activateLast iec.BOOL
	state        iec.INT
	s1           iec.STRING // STRING(20)
	i            iec.INT
}

// INIT resets the block.
func (s *SYS_LOG) INIT() { *s = SYS_LOG{client: client{IP_C: s.IP_C, S_BUF: s.S_BUF, R_BUF: s.R_BUF}} }

// put copies text into the send buffer at i, and moves i past it.
func (s *SYS_LOG) put(text iec.STRING) {
	buffer.STRING_TO_BUFFER_(text, s.i, s.S_BUF.BUFFER[:], 1024)
	s.i += LEN(text)
}

// Execute runs the block once.
func (s *SYS_LOG) Execute(now time.Time) {
	if !s.ok() {
		return
	}
	c, sb, rb := s.IP_C, s.S_BUF, s.R_BUF
	switch s.state {
	case 0:
		if s.ACTIVATE && !s.activateLast {
			s.state = 10
			s.ERROR = 0
			s.ipState = 1
		}
		s.DONE_P = false
	case 10:
		if s.ipState == 3 {
			s.i = 0
			if s.OPTION&0x01 == 0 {
				// The priority <facility * 8 + level>, in decimal.
				s.s1 = REPLACE("<P>", network.BYTE_TO_STRING(SHL(s.FACILITY, 3)|s.SEVERITY&0x07), 1, 2)
			} else {
				s.s1 = ""
			}
			if s.OPTION&0x02 == 0 {
				// The header: 'month day hh:mm:ss host tag '.
				s.s1 = network.STRING_N(CONCAT(s.s1, str.DT_TO_STRF(s.LDT, 0, "#E #W #N:#R:#T ", 1)), 20)
				s.put(s.s1)
				s.put(s.HOSTNAME)
				sb.BUFFER[s.i] = 32
				s.i++
				s.put(s.TAG)
				sb.BUFFER[s.i] = 32
				s.i++
			}
			s.put(s.MESSAGE)
			if s.OPTION&0x04 > 0 {
				sb.BUFFER[s.i], sb.BUFFER[s.i+1] = 13, 10
				s.i += 2
			}
			s.state = 30
			if s.OPTION&0x08 > 0 {
				c.C_MODE = 0 // TCP client
			} else {
				c.C_MODE = 1 // UDP client
			}
			switch {
			case s.PORT != 0:
				c.C_PORT = s.PORT
			case c.C_MODE == 0:
				c.C_PORT = 1468
			default:
				c.C_PORT = 514
			}
			c.C_IP = s.SERVER_IP4
			c.C_ENABLE = true
			c.TIME_RESET = true
			c.R_OBSERVE = false
			sb.SIZE = iec.UINT(s.i)
			rb.SIZE = 0
		}
	case 30:
		if c.ERROR != 0 {
			s.ERROR = c.ERROR
			s.state = 0
		} else if sb.SIZE == 0 { // all sent
			s.DONE_P = true
			s.state = 0
		}
		if s.state == 0 {
			s.ipState = 4
		}
	}
	s.activateLast = s.ACTIVATE
	s.runFIFO(now, s.IP_C)
}

// TELNET_PRINT is a telnet server on PORT that prints TEXT, while ENABLE,
// each scan SEND is true, at X_POS, Y_POS if not 0, in the colors
// FRONT_COLOR and BACK_COLOR, into S_BUF, which it sends when it is full or
// SEND goes false. READY while a client is connected, DONE for a scan when
// TEXT is in the buffer. The bits of OPTION: 0 clear the screen, 1 wrap
// lines, 2 colors, 3 a CR LF after each text, 7 send each text at once.
type TELNET_PRINT struct {
	IP_C        *network.IP_C
	S_BUF       *network.NETWORK_BUFFER
	TEXT        iec.STRING // STRING(STRING_LENGTH)
	ENABLE      iec.BOOL
	SEND        iec.BOOL
	OPTION      iec.BYTE // default 2#1000_1100
	BACK_COLOR  iec.BYTE
	FRONT_COLOR iec.BYTE
	X_POS       iec.BYTE
	Y_POS       iec.BYTE
	PORT        iec.WORD // default 23
	READY       iec.BOOL
	DONE        iec.BOOL

	lastFC iec.BYTE
	state  iec.INT
	x      iec.INT // default -1
	queue
	init          iec.BOOL
	b0ScreenClear iec.BOOL
	b1Autowrap    iec.BOOL
	b2Color       iec.BOOL
	b3CRLF        iec.BOOL
	b7NoFlush     iec.BOOL
	initialized   bool
}

// INIT resets the block and sets its inputs to their initial values.
func (t *TELNET_PRINT) INIT() {
	*t = TELNET_PRINT{IP_C: t.IP_C, S_BUF: t.S_BUF, OPTION: 0x8C, PORT: 23, x: -1, initialized: true}
}

// add writes b at the next place of the send buffer.
func (t *TELNET_PRINT) add(b ...iec.BYTE) {
	for _, c := range b {
		t.x++
		t.S_BUF.BUFFER[t.x] = c
	}
}

// Execute runs the block once.
func (t *TELNET_PRINT) Execute(now time.Time) {
	if !t.initialized {
		t.initialized = true
		t.x = -1
	}
	if t.IP_C == nil || t.S_BUF == nil {
		return
	}
	c, sb := t.IP_C, t.S_BUF
	t.READY = c.C_STATE > 127
	t.DONE = false
	switch t.state {
	case 0:
		if t.ENABLE {
			t.state = 10
			t.ipState = 1
		}
	case 10:
		if t.ipState == 3 {
			c.C_PORT = t.PORT
			c.C_IP = 0
			c.C_MODE = 4 // TCP server for any client
			c.TIME_RESET = true
			c.C_ENABLE = true
			c.R_OBSERVE = false
			t.state = 20
		}
	case 20:
		if c.ERROR > 0 {
			c.TIME_RESET = true
		}
		if sb.SIZE != 0 {
			break
		}
		if !t.ENABLE {
			t.ipState = 4
			c.C_ENABLE = false
			t.state = 0
		}
		if c.C_STATE == 1 { // disconnected
			t.init = false
			break
		}
		if !t.READY {
			break
		}
		if c.C_STATE == 254 && !t.init { // connected
			t.init = true
			t.b0ScreenClear = t.OPTION&0x01 > 0
			t.b1Autowrap = t.OPTION&0x02 > 0
			t.b2Color = t.OPTION&0x04 > 0
			t.b3CRLF = t.OPTION&0x08 > 0
			t.b7NoFlush = t.OPTION&0x80 > 0
			t.lastFC = ^t.FRONT_COLOR // force the colors
			t.x = -1
			// ESC[?7h or l: wrap lines or not.
			t.add(0x1B, 0x5B, 0x3F, 0x37, SEL[iec.BYTE](t.b1Autowrap, 0x6C, 0x68))
			if t.b0ScreenClear {
				if t.b2Color {
					// ESC[0;3f;4bm
					t.add(0x1B, 0x5B, 0x30, 0x3B, 0x33, t.FRONT_COLOR&0x07|0x30, 0x3B, 0x34, t.BACK_COLOR&0x07|0x30, 0x6D)
				}
				t.add(0x1B, 0x5B, 0x32, 0x4A) // ESC[2J: clear the screen
			} else {
				t.add(0x0A)
			}
		}
		// Send a full buffer, or the data when there is no more.
		if iec.UINT(t.x+LEN(t.TEXT)+21) > iec.UINT(len(sb.BUFFER)) || t.x >= 0 && !t.SEND {
			sb.SIZE = iec.UINT(t.x + 1)
			t.x = -1
		} else if t.SEND {
			if t.b2Color && t.FRONT_COLOR != t.lastFC {
				// ESC[0;[1;][5;]3fm
				t.add(0x1B, 0x5B, 0x30, 0x3B)
				if t.FRONT_COLOR&0x08 == 0 {
					t.add(0x31, 0x3B) // bold and bright
				}
				if t.FRONT_COLOR&0x10 > 0 {
					t.add(0x35, 0x3B) // blinking
				}
				t.add(0x33, t.FRONT_COLOR&0x07|0x30, 0x6D)
				t.lastFC = t.FRONT_COLOR
			}
			if t.X_POS > 0 && t.Y_POS > 0 {
				// ESC[yy;xxH: the cursor
				y, x := logic.INT_TO_BCDC(iec.INT(t.Y_POS)), logic.INT_TO_BCDC(iec.INT(t.X_POS))
				t.add(0x1B, 0x5B, SHR(y, 4)|0x30, y&0x0F|0x30, 0x3B, SHR(x, 4)|0x30, x&0x0F|0x30, 0x48)
			}
			if LEN(t.TEXT) > 0 {
				buffer.STRING_TO_BUFFER_(t.TEXT, t.x+1, sb.BUFFER[:], iec.UINT(len(sb.BUFFER)))
				t.x += LEN(t.TEXT)
			}
			if t.b3CRLF {
				t.add(0x0D, 0x0A)
			}
			t.DONE = true
			if t.b7NoFlush {
				sb.SIZE = iec.UINT(t.x + 1)
				t.x = -1
			}
		}
	}
	t.runFIFO(now, c)
}

// TELNET_LOG is a telnet server on PORT that shows the messages of the log
// LOG_CL, as they come, with TELNET_PRINT and its OPTION, while ENABLE.
// READY while a client is connected.
type TELNET_LOG struct {
	IP_C   *network.IP_C
	S_BUF  *network.NETWORK_BUFFER
	LOG_CL *network.LOG_CONTROL
	ENABLE iec.BOOL
	OPTION iec.BYTE // default 2#1000_1100
	PORT   iec.WORD // default 23
	READY  iec.BOOL

	print       TELNET_PRINT
	done        iec.BOOL
	init        iec.BOOL
	watchdog    timers.TON
	ci          iec.INT
	pi          iec.INT
	piLast      iec.INT
	send        iec.BOOL
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (t *TELNET_LOG) INIT() {
	*t = TELNET_LOG{IP_C: t.IP_C, S_BUF: t.S_BUF, LOG_CL: t.LOG_CL, OPTION: 0x8C, PORT: 23, initialized: true}
	t.print.INIT()
}

// Execute runs the block once.
func (t *TELNET_LOG) Execute(now time.Time) {
	if !t.initialized {
		t.initialized = true
		t.print.INIT()
	}
	if t.IP_C == nil || t.S_BUF == nil || t.LOG_CL == nil {
		return
	}
	c, lc := t.IP_C, t.LOG_CL
	if !t.init {
		t.init = true
		if t.OPTION > 0 {
			t.print.OPTION = t.OPTION
		}
		if t.PORT > 0 {
			t.print.PORT = t.PORT
		}
		t.watchdog.PT = iec.TIME(time.Millisecond)
	}
	t.watchdog.IN = false
	t.watchdog.Execute(now)
	t.ci = lc.IDX
	if c.C_STATE == 1 { // disconnected
		t.pi = 0
		t.init = false
	}
	for {
		if t.pi != t.ci && c.C_STATE > 127 {
			if t.pi == 0 && lc.RING_MODE {
				t.pi = t.ci
			}
			t.piLast = t.pi
			t.pi++
			if t.pi > lc.SIZE {
				t.pi = 1
			}
			t.send = true
		} else {
			t.send = false
		}
		p := &t.print
		p.IP_C, p.S_BUF, p.ENABLE, p.SEND = c, t.S_BUF, t.ENABLE, t.send
		p.TEXT = lc.MSG[t.pi]
		p.BACK_COLOR = iec.BYTE(SHL(lc.MSG_OPTION[t.pi], 8))
		p.FRONT_COLOR = iec.BYTE(lc.MSG_OPTION[t.pi])
		p.Execute(now)
		t.READY = p.READY
		t.done = p.DONE
		if t.send && !t.done {
			t.pi = t.piLast
		}
		t.watchdog.IN = true
		t.watchdog.Execute(now)
		if t.watchdog.Q || !t.done {
			break
		}
	}
}
