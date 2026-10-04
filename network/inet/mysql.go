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
	"github.com/apiarytech/beebread/network/crypto"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// MYSQL_AUTH computes the answer SCRAMBLE of MySQL's native password
// authentication to the message MESSAGE of the server, with PASSWORD (55
// characters at most): SHA1(PASSWORD) XOR SHA1(MESSAGE + SHA1(SHA1(PASSWORD))).
// RUN starts it and is cleared at the end.
type MYSQL_AUTH struct {
	RUN      *iec.BOOL
	PASSWORD *iec.STRING // STRING(64)
	MESSAGE  *[20]iec.BYTE
	SCRAMBLE *[20]iec.BYTE

	sha       crypto.SHA1_STREAM
	buf       [64]iec.BYTE
	shaMode   iec.INT
	step      iec.INT
	size      iec.UDINT
	pos       iec.UDINT
	scramble1 [20]iec.BYTE
	scramble2 [20]iec.BYTE
}

// INIT resets the block.
func (m *MYSQL_AUTH) INIT() {
	*m = MYSQL_AUTH{RUN: m.RUN, PASSWORD: m.PASSWORD, MESSAGE: m.MESSAGE, SCRAMBLE: m.SCRAMBLE}
}

// Execute runs the block once.
func (m *MYSQL_AUTH) Execute(now time.Time) {
	if m.RUN == nil || m.PASSWORD == nil || m.MESSAGE == nil || m.SCRAMBLE == nil || !*m.RUN {
		return
	}
	switch m.step {
	case 0:
		switch m.shaMode {
		case 0: // SHA1(PASSWORD)
			m.shaMode = 1
			pw := network.STRING_N(*m.PASSWORD, 64)
			m.size = iec.UDINT(LEN(pw))
			for i := 0; i < 64; i++ {
				if i < int(m.size) {
					m.buf[i] = str.CODE(pw, iec.INT(i+1))
				} else {
					m.buf[i] = 0
				}
			}
		case 3: // SHA1(SHA1(PASSWORD))
			m.scramble1 = *m.SCRAMBLE
			m.shaMode = 1
			m.size = 20
			copy(m.buf[:20], m.scramble1[:])
			m.step = 10
		}
	case 10:
		if m.shaMode == 3 { // SHA1(MESSAGE + SHA1(SHA1(PASSWORD)))
			m.shaMode = 1
			m.scramble2 = *m.SCRAMBLE
			m.size = 40
			copy(m.buf[:20], m.MESSAGE[:])
			copy(m.buf[20:40], m.scramble2[:])
			m.step = 20
		}
	case 20:
		if m.shaMode == 3 {
			m.shaMode = 0
			for i := range m.SCRAMBLE {
				m.SCRAMBLE[i] ^= m.scramble1[i]
			}
			m.step = 0
			*m.RUN = false
		}
	}
	m.sha.SIZE, m.sha.MODE, m.sha.BUF, m.sha.SHA1 = &m.size, &m.shaMode, &m.buf, m.SCRAMBLE
	m.sha.Execute(now)
	m.pos = m.sha.POS
}

// MYSQL_CONTROL keeps a connection to the MySQL server COM.SQL_URL, a URL
// 'mysql://user:password@host:port', port 3306 by default, while
// COM.SQL_CON: it logs in, and INFO.SQL_CONNECTED is true. It sends the
// packets the users put in COM.S_BUF, from byte 4, numbering them, and
// COM.SQL_RCV_STATE is 1 when an answer is in COM.R_BUF, 2 for an error,
// with INFO.SQL_ERROR. COM.ERROR_C and COM.ERROR_T: 1 DNS, 2 the connection
// or a protocol before 10, 3 an error of the server, 5 a step that took
// longer than COM.TIMEOUT (10 s at least).
type MYSQL_CONTROL struct {
	COM  *network.MYSQL_COM
	INFO *network.MYSQL_INFO

	ipC           network.IP_C
	urlData       network.URL
	dns           DNS_CLIENT
	ipc           ip.IP_CONTROL
	step          iec.INT
	sndStep       iec.INT
	tonWait       timers.TON
	lastStep      iec.INT
	dwTmp         iec.DWORD
	timeout1      iec.TIME
	timeout2      iec.TIME
	ipCRedDisable iec.BOOL
	authPassword  iec.STRING // STRING(64)
	authRun       iec.BOOL
	auth          MYSQL_AUTH
	scramble1     [20]iec.BYTE
	scramble2     [20]iec.BYTE
	charsetNumber iec.BYTE  // default 16#08
	maxPacketSize iec.DWORD // default 1024
	clientFlags   iec.DWORD // default 16#0003_8601
	initialized   bool
}

// INIT resets the block.
func (m *MYSQL_CONTROL) INIT() {
	*m = MYSQL_CONTROL{COM: m.COM, INFO: m.INFO}
	m.defaults()
}

// defaults sets the initial values of the internal state.
func (m *MYSQL_CONTROL) defaults() {
	m.initialized = true
	m.charsetNumber, m.maxPacketSize, m.clientFlags = 0x08, 1024, 0x0003_8601
}

// Execute runs the block once.
func (m *MYSQL_CONTROL) Execute(now time.Time) {
	if !m.initialized {
		m.defaults()
	}
	if m.COM == nil || m.INFO == nil {
		return
	}
	com, info, c := m.COM, m.INFO, &m.ipC
	sb, rb := &com.S_BUF, &com.R_BUF
	switch m.step {
	case 0:
		if com.SQL_CON {
			m.timeout1 = max(iec.TIME(10*time.Second), com.TIMEOUT)
			m.timeout2 = m.timeout1 + iec.TIME(time.Second)
			info.SQL_ERROR = ""
			com.ERROR_C, com.ERROR_T = 0, 0
			m.urlData = encoding.STRING_TO_URL(com.SQL_URL, "", "/")
			m.step = 10
		}
	case 10:
		if m.dns.DONE {
			m.step = 20
		} else if m.dns.ERROR != 0 {
			com.ERROR_C, com.ERROR_T = m.dns.ERROR, 1
			m.step = 980
		}
	case 20: // the connection
		if m.urlData.PORT == 0 {
			m.urlData.PORT = 3306
		}
		c.C_PORT = m.urlData.PORT
		c.C_IP = m.dns.IP4
		c.C_MODE = 0 // TCP client
		c.TIME_RESET = true
		c.C_ENABLE = true
		c.R_OBSERVE = true
		rb.SIZE = 0
		com.SQL_RCV_STATE = 0
		m.step = 30
	case 30:
		if com.SQL_RCV_STATE == 1 { // the handshake
			info.SERVER_PROTOCOL_VERSION = rb.BUFFER[4]
			// The end of the server version.
			idx1 := iec.INT(5)
			for ; idx1 <= iec.INT(rb.SIZE)-1; idx1++ {
				if rb.BUFFER[idx1] == 0 {
					break
				}
			}
			idx1 += 5 // the 0 and the thread ID
			for i := 0; i < 8; i++ {
				m.scramble1[i] = rb.BUFFER[idx1]
				idx1++
			}
			idx1++
			info.SERVER_CAPABILITIES = logic.WORD_OF_BYTE(rb.BUFFER[idx1+1], rb.BUFFER[idx1])
			idx1 += 2
			info.SERVER_LANGUAGE = rb.BUFFER[idx1]
			idx1++
			info.SERVER_STATUS = logic.WORD_OF_BYTE(rb.BUFFER[idx1+1], rb.BUFFER[idx1])
			idx1 += 15
			for i := 8; i < 20; i++ {
				m.scramble1[i] = rb.BUFFER[idx1]
				idx1++
			}
			switch {
			case info.SERVER_PROTOCOL_VERSION < 10:
				com.ERROR_T, com.ERROR_C = 2, 1
				m.step = 980
			case m.urlData.PASSWORD == "":
				m.step = 50
			default:
				m.authPassword = network.STRING_N(m.urlData.PASSWORD, 64)
				m.authRun = true
				m.step = 40
			}
		}
	case 40:
		m.auth.RUN, m.auth.PASSWORD, m.auth.MESSAGE, m.auth.SCRAMBLE = &m.authRun, &m.authPassword, &m.scramble1, &m.scramble2
		m.auth.Execute(now)
		if !m.authRun {
			m.step = 50
		}
	case 50: // the authentication packet, from byte 4
		for i := 0; i < 4; i++ {
			sb.BUFFER[4+i] = logic.BYTE_OF_DWORD(m.clientFlags, iec.BYTE(i))
			sb.BUFFER[8+i] = logic.BYTE_OF_DWORD(m.maxPacketSize, iec.BYTE(i))
		}
		sb.BUFFER[12] = m.charsetNumber
		idx1 := iec.INT(13)
		for ; idx1 <= 35; idx1++ { // 23 bytes of 0
			sb.BUFFER[idx1] = 0
		}
		for i := iec.INT(1); i <= LEN(m.urlData.USER); i++ {
			sb.BUFFER[idx1] = str.CODE(m.urlData.USER, i)
			idx1++
		}
		sb.BUFFER[idx1] = 0
		idx1++
		if LEN(m.urlData.PASSWORD) == 0 {
			sb.BUFFER[idx1] = 0
		} else {
			sb.BUFFER[idx1] = 20
			for i := 0; i < 20; i++ {
				idx1++
				sb.BUFFER[idx1] = m.scramble2[i]
			}
		}
		sb.SIZE = iec.UINT(idx1 + 1)
		m.step = 60
	case 60:
		if com.SQL_RCV_STATE == 1 {
			info.SQL_CONNECTED = true
			m.step = 300
		}
	case 300:
		m.step = SEL[iec.INT](com.SQL_CON, 700, 310)
	case 310:
		m.step = 300
	case 700: // COM_QUIT
		sb.BUFFER[4] = 1
		sb.SIZE = 5
		com.SQL_PACKET_NO = 255
		m.ipCRedDisable = true
		m.step = 710
	case 710:
		if c.C_STATE == 0 { // the server closed the connection
			m.step = 980
		}
	case 980:
		c.C_ENABLE = false
		m.ipCRedDisable = false
		info.SQL_CONNECTED = false
		com.SQL_CON = false
		sb.SIZE, rb.SIZE = 0, 0
		m.sndStep = 0
		m.step = 990
	case 990:
		if c.C_STATE == 0 {
			m.step = 0
		}
	}

	// A packet received: its length, and an error.
	if m.step >= 30 && rb.SIZE >= 3 && com.SQL_RCV_STATE == 0 {
		if iec.UINT(logic.DWORD_OF_BYTE(0, rb.BUFFER[2], rb.BUFFER[1], rb.BUFFER[0]))+4 == rb.SIZE {
			com.SQL_PACKET_NO = rb.BUFFER[3]
			if rb.BUFFER[4] == 0xFF {
				idx2, idx3 := iec.INT(13), iec.INT(rb.SIZE)-1
				for i := idx2; i <= idx3; i++ {
					if rb.BUFFER[i] == 39 { // ' to `
						rb.BUFFER[i] = 96
					}
				}
				info.SQL_ERROR = buffer.BUFFER_TO_STRING(rb.BUFFER[:], rb.SIZE, iec.UINT(idx2), iec.UINT(idx3))
				com.ERROR_T = 3
				com.ERROR_C = logic.DWORD_OF_BYTE(0, 0, rb.BUFFER[6], rb.BUFFER[5])
				com.SQL_RCV_STATE = 2
				rb.SIZE = 0
				m.step = 980
			} else {
				com.SQL_RCV_STATE = 1
				c.R_OBSERVE = false
			}
		}
	}

	// Send a packet: its length and number first.
	switch m.sndStep {
	case 0:
		if m.step >= 20 && sb.SIZE > 0 {
			m.dwTmp = iec.DWORD(sb.SIZE - 4)
			sb.BUFFER[0] = logic.BYTE_OF_DWORD(m.dwTmp, 0)
			sb.BUFFER[1] = logic.BYTE_OF_DWORD(m.dwTmp, 1)
			sb.BUFFER[2] = 0
			com.SQL_PACKET_NO++
			sb.BUFFER[3] = com.SQL_PACKET_NO
			c.R_OBSERVE = true
			rb.SIZE = 0
			com.SQL_RCV_STATE = 0
			m.sndStep = 10
		}
	case 10:
		if sb.SIZE == 0 {
			m.sndStep = 0
		}
	}

	if com.ERROR_T == 0 && m.step > 20 && c.ERROR > 0 && c.C_ENABLE && !c.TIME_RESET &&
		(!m.ipCRedDisable || c.ERROR != 0xFD00_0000) {
		com.ERROR_C, com.ERROR_T = c.ERROR, 2
		m.step = 980
	}
	if m.tonWait.Q {
		com.ERROR_C, com.ERROR_T = iec.DWORD(m.step), 5
		m.step = 980
	}

	m.dns.IP_C, m.dns.S_BUF, m.dns.R_BUF = c, sb, rb
	m.dns.DOMAIN, m.dns.IP4_DNS, m.dns.ACTIVATE = m.urlData.DOMAIN, com.DNS_IP4, m.step == 10
	m.dns.Execute(now)
	m.ipc.IP_C, m.ipc.S_BUF, m.ipc.R_BUF = c, sb, rb
	m.ipc.IP, m.ipc.PORT, m.ipc.TIME_OUT = 0, 0, m.timeout1
	m.ipc.Execute(now)
	m.tonWait.IN, m.tonWait.PT = m.step == m.lastStep && m.step > 0, m.timeout2
	m.tonWait.Execute(now)
	m.lastStep = m.step
}
