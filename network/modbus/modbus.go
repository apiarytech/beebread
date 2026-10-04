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

// Package modbus is the port of the OSCAT NETWORK Modbus TCP and UDP client
// and servers, on an IP_CONTROL2 and its short buffers: the function codes
// 1 to 6, 15, 16, 22 and 23, on a data area of 256 words.
//
// An index outside the data area, which OSCAT does not always check, reads
// 0 and writes nothing.
package modbus

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/logic"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// DATA is the data area of the Modbus blocks.
type DATA = [256]iec.WORD // ARRAY[0..255] OF WORD

// VMAP is the map of the virtual addresses of a server.
type VMAP = [10]network.VMAP_DATA // ARRAY[1..10] OF VMAP_DATA

// area is a data area with guarded indexes.
type area struct{ d *DATA }

func (a area) get(i iec.INT) iec.WORD {
	if i < 0 || int(i) >= len(a.d) {
		return 0
	}
	return a.d[i]
}

func (a area) set(i iec.INT, v iec.WORD) {
	if i >= 0 && int(i) < len(a.d) {
		a.d[i] = v
	}
}

// word returns the word of the bytes hi and lo as an INT.
func word(hi, lo iec.BYTE) iec.INT { return iec.INT(logic.WORD_OF_BYTE(hi, lo)) }

// putWord writes w into b at i and i+1, high byte first.
func putWord(b []iec.BYTE, i iec.INT, w iec.WORD) {
	b[i] = iec.BYTE(SHR(w, 8))
	b[i+1] = iec.BYTE(w)
}

// MB_VMAP maps the virtual address V_ADR, of V_CNT points, of the function
// code FC to the address P_ADR (and bit P_BIT for the bit functions) of a
// data area of SIZE words, by VMAP: each entry maps the functions whose
// bits are set in FC, V_SIZE words from V_ADR to P_ADR. An empty map maps
// all 256 words. ERROR is 1 for a function the map does not know, 2 for an
// address outside it. A write sets the entry's TIME_OUT.
type MB_VMAP struct {
	VMAP  *VMAP
	FC    iec.INT
	V_ADR iec.INT
	V_CNT iec.INT
	SIZE  iec.INT
	P_ADR iec.INT
	P_BIT iec.INT
	ERROR iec.BYTE

	init iec.BOOL
}

// INIT resets the block.
func (m *MB_VMAP) INIT() { *m = MB_VMAP{VMAP: m.VMAP} }

// Execute runs the block once.
func (m *MB_VMAP) Execute(now time.Time) {
	if m.VMAP == nil {
		return
	}
	v := m.VMAP
	if !m.init {
		m.init = true
		if v[0].FC == 0 {
			v[0] = network.VMAP_DATA{FC: 0xFFFF_FFFF, V_ADR: 0, V_SIZE: 256, P_ADR: 0, TIME_OUT: v[0].TIME_OUT}
		}
	}
	m.ERROR = 2 // ILLEGAL DATA ADDRESS
	mask := SHL(iec.DWORD(1), m.FC)
	if mask&0x0040_0060 != 0 { // FC 5, 6, 22: a point
		m.V_CNT = 1
	} else if m.V_CNT == 0 {
		return
	}
	var i int
	for i = 0; i < len(v); i++ {
		if v[i].FC&mask == 0 {
			continue
		}
		vadr, vsize, padr := v[i].V_ADR, v[i].V_SIZE, v[i].P_ADR
		if mask&0x0000_8026 != 0 { // FC 1, 2, 5, 15: bits
			if iec.INT(SHR(iec.WORD(m.V_ADR+m.V_CNT-1), 4))+1 <= m.SIZE && iec.INT(SHR(iec.WORD(m.V_ADR+m.V_CNT-1), 4))+1 <= vsize {
				w := iec.WORD(m.V_ADR)
				m.P_ADR = iec.INT(SHR(w, 4)) + padr
				m.P_BIT = iec.INT(w & 0x000F)
				m.ERROR = 0
				break
			}
		} else if mask&0x00C1_0058 != 0 { // FC 3, 4, 6, 16, 22, 23: registers
			if m.V_ADR >= vadr && m.V_ADR+m.V_CNT <= vadr+vsize {
				m.P_ADR = m.V_ADR - vadr + padr
				if m.P_ADR+m.V_CNT <= m.SIZE {
					m.ERROR = 0
					break
				}
			}
		} else {
			m.ERROR = 1 // ILLEGAL FUNCTION
			return
		}
	}
	if m.ERROR == 0 && mask&0x00C1_8060 != 0 { // FC 5, 6, 15, 16, 22, 23: a write
		v[i].TIME_OUT = iec.TIME(time.Millisecond)
	}
}

// shortClient is what the Modbus blocks share: the connection, its short
// buffers and the data area, and their place in the connection's queue.
type shortClient struct {
	IP_C  *network.IP_C
	S_BUF *network.NETWORK_BUFFER_SHORT
	R_BUF *network.NETWORK_BUFFER_SHORT
	DATA  *DATA

	fifo    ip.IP_FIFO
	ipState iec.BYTE
	ipID    iec.BYTE
}

func (c *shortClient) ok() bool {
	return c.IP_C != nil && c.S_BUF != nil && c.R_BUF != nil && c.DATA != nil
}

func (c *shortClient) runFIFO(now time.Time) {
	c.fifo.FIFO, c.fifo.STATE, c.fifo.ID = &c.IP_C.FIFO, &c.ipState, &c.ipID
	c.fifo.Execute(now)
}

// MB_CLIENT is a Modbus client: while ENABLE, every DELAY, it sends the
// request of the function code FC to the server of the unit UNIT_ID that
// the IP_CONTROL2 of IP_C has, over TCP or UDP. It reads R_POINTS points
// from R_ADDR into DATA from R_DATA_ADR (and bit R_DATA_BITPOS), and writes
// W_POINTS points from DATA at W_DATA_ADR (and bit W_DATA_BITPOS) to
// W_ADDR. ERROR is the Modbus exception, or the connection's; BUSY while it
// has the connection.
type MB_CLIENT struct {
	shortClient
	DATA_SIZE     iec.INT
	ENABLE        iec.BOOL
	UDP           iec.BOOL
	FC            iec.INT
	UNIT_ID       iec.BYTE
	R_ADDR        iec.INT
	R_POINTS      iec.INT
	R_DATA_ADR    iec.INT
	R_DATA_BITPOS iec.INT
	W_ADDR        iec.INT
	W_POINTS      iec.INT
	W_DATA_ADR    iec.INT
	W_DATA_BITPOS iec.INT
	DELAY         iec.TIME
	ERROR         iec.DWORD
	BUSY          iec.BOOL

	state         iec.INT
	transactionID iec.INT
	response      iec.INT
	comp          iec.INT
	ton1          timers.TON
}

// INIT resets the block.
func (m *MB_CLIENT) INIT() {
	*m = MB_CLIENT{shortClient: shortClient{IP_C: m.IP_C, S_BUF: m.S_BUF, R_BUF: m.R_BUF, DATA: m.DATA}}
}

// Execute runs the block once.
func (m *MB_CLIENT) Execute(now time.Time) {
	if !m.ok() {
		return
	}
	c, sb, rb, data := m.IP_C, m.S_BUF, m.R_BUF, area{m.DATA}
	s := sb.BUFFER[:]
	switch m.state {
	case 0: // wait for DELAY
		m.ton1.IN, m.ton1.PT = m.ENABLE, m.DELAY
		m.ton1.Execute(now)
		if m.ton1.Q || m.DELAY == 0 {
			m.ton1.IN = false
			m.ton1.Execute(now)
			m.state = 20
			m.ipState = 1
		}
	case 20:
		if m.ipState != 3 {
			break
		}
		c.C_PORT = 0 // the port and address of the IP_CONTROL2
		c.C_IP = 0
		c.C_MODE = iec.BYTE(BOOL_TO_INT(m.UDP)) // 0 TCP, 1 UDP client
		c.TIME_RESET = true
		c.C_ENABLE = true
		c.R_OBSERVE = true
		m.ERROR = 0
		switch m.FC {
		case 1, 2: // Read Coil Status, Read Input Status
			if m.R_POINTS <= 2000 {
				putWord(s, 8, iec.WORD(m.R_ADDR))
				putWord(s, 10, iec.WORD(m.R_POINTS))
				sb.SIZE = 12
				m.response = iec.INT(SHR(iec.WORD(m.R_POINTS+7), 3))
				m.comp = 7
			} else {
				m.ERROR = 2 // ILLEGAL DATA ADDRESS
			}
		case 3, 4: // Read Holding Registers, Read Input Registers
			if m.R_POINTS <= 125 {
				putWord(s, 8, iec.WORD(m.R_ADDR))
				putWord(s, 10, iec.WORD(m.R_POINTS))
				sb.SIZE = 12
				m.response = iec.INT(SHL(iec.WORD(m.R_POINTS), 1))
				m.comp = 7
			} else {
				m.ERROR = 2
			}
		case 5: // Force Single Coil
			if m.W_DATA_ADR <= m.DATA_SIZE {
				putWord(s, 8, iec.WORD(m.W_ADDR))
				s[10] = SEL[iec.BYTE](data.get(m.W_DATA_ADR)&SHL(iec.WORD(1), m.W_DATA_BITPOS) > 0, 0x00, 0xFF)
				s[11] = 0
				sb.SIZE = 12
				m.response, m.comp = 3, 11
			} else {
				m.ERROR = 2
			}
		case 6: // Preset Single Register
			if m.W_DATA_ADR <= m.DATA_SIZE {
				putWord(s, 8, iec.WORD(m.W_ADDR))
				putWord(s, 10, data.get(m.W_DATA_ADR))
				sb.SIZE = 12
				m.response, m.comp = 3, 11
			} else {
				m.ERROR = 2
			}
		case 15: // Force Multiple Coils
			last := iec.INT(SHR(iec.WORD(iec.INT(SHL(iec.WORD(m.W_DATA_ADR), 4))+m.W_POINTS+m.W_DATA_BITPOS+15), 4))
			if last <= m.DATA_SIZE && m.W_POINTS <= 1968 {
				putWord(s, 8, iec.WORD(m.W_ADDR))
				putWord(s, 10, iec.WORD(m.W_POINTS))
				s[12] = iec.BYTE(SHR(iec.WORD(m.W_POINTS+7), 3))
				wMask := SHL(iec.WORD(1), m.W_DATA_BITPOS)
				bitPos := iec.INT(0)
				var w iec.WORD
				idx1, idx2 := iec.INT(12), m.W_DATA_ADR
				for i := iec.INT(1); i <= m.W_POINTS; i++ {
					w = logic.BIT_LOAD_W(w, data.get(idx2)&wMask > 0, bitPos)
					bitPos++
					if bitPos > 7 || i == m.W_POINTS {
						bitPos = 0
						idx1++
						s[idx1] = iec.BYTE(w)
						w = 0
					}
					wMask = ROL(wMask, 1)
					if wMask == 1 {
						idx2++
					}
				}
				sb.SIZE = iec.UINT(idx1 + 1)
				m.response, m.comp = 3, 11
			} else {
				m.ERROR = 2
			}
		case 16: // Write Multiple Registers
			if m.W_DATA_ADR+m.W_POINTS <= m.DATA_SIZE && m.W_POINTS <= 123 {
				putWord(s, 8, iec.WORD(m.W_ADDR))
				putWord(s, 10, iec.WORD(m.W_POINTS))
				s[12] = SHL(s[11], 1)
				idx1, idx2 := iec.INT(11), iec.INT(0)
				for i := m.W_DATA_ADR; i <= m.W_DATA_ADR+m.W_POINTS-1; i++ {
					idx1 += 2
					idx2 = idx1 + 1
					putWord(s, idx1, data.get(i))
				}
				sb.SIZE = iec.UINT(idx2 + 1)
				m.response, m.comp = 3, 11
			} else {
				m.ERROR = 2
			}
		case 22: // Mask Write Register
			if m.W_DATA_ADR+2 <= m.DATA_SIZE {
				putWord(s, 8, iec.WORD(m.W_ADDR))
				putWord(s, 10, data.get(m.W_DATA_ADR))   // the AND mask
				putWord(s, 12, data.get(m.W_DATA_ADR+1)) // the OR mask
				sb.SIZE = 14
				m.response, m.comp = 5, 13
			} else {
				m.ERROR = 2
			}
		case 23: // Read/Write Multiple Registers
			if m.W_DATA_ADR+m.W_POINTS <= m.DATA_SIZE && m.W_POINTS <= 121 && m.R_POINTS <= 125 {
				idx1, idx2 := iec.INT(15), iec.INT(0)
				for i := m.W_DATA_ADR; i <= m.W_DATA_ADR+m.W_POINTS-1; i++ {
					idx1 += 2
					idx2 = idx1 + 1
					putWord(s, idx1, data.get(i))
				}
				putWord(s, 12, iec.WORD(m.W_ADDR))
				putWord(s, 14, iec.WORD(m.W_POINTS))
				s[16] = SHL(s[15], 1)
				putWord(s, 8, iec.WORD(m.R_ADDR))
				putWord(s, 10, iec.WORD(m.R_POINTS))
				sb.SIZE = iec.UINT(idx2 + 1)
				m.response = iec.INT(SHL(iec.WORD(m.R_POINTS), 1))
				m.comp = 7
			} else {
				m.ERROR = 2
			}
		default:
			m.ERROR = 1 // ILLEGAL FUNCTION
		}
		if m.ERROR == 0 {
			// The header.
			m.transactionID++
			s[0] = m.ipID // the transaction ID: the block's ID and a count
			s[1] = iec.BYTE(m.transactionID)
			s[2], s[3] = 0, 0 // the protocol
			putWord(s, 4, iec.WORD(sb.SIZE-6))
			s[6] = m.UNIT_ID
			s[7] = iec.BYTE(m.FC)
			rb.SIZE = 0
			m.state = 30
		} else {
			m.ipState = 4
			m.state = 0
		}
	case 30: // the response
		if c.ERROR != 0 {
			m.ERROR = c.ERROR
			m.ipState = 4
			m.state = 0
			break
		}
		if !(sb.SIZE == 0 && rb.SIZE > 0) {
			break
		}
		r := rb.BUFFER[:]
		// The header expected, with the length of the response.
		putWord(s, 4, iec.WORD(m.response+3))
		var i iec.INT
		for i = 0; i <= m.comp; i++ {
			if s[i] != r[i] {
				break
			}
		}
		switch {
		case rb.SIZE >= 9 && r[7] > 128: // an exception
			m.ERROR = iec.DWORD(r[8])
		case rb.SIZE == iec.UINT(m.response+9) && i > m.comp:
			mask := SHL(iec.DWORD(1), m.FC)
			switch {
			case mask&0x0000_0006 > 0: // FC 1, 2
				last := iec.INT(SHR(iec.WORD(iec.INT(SHL(iec.WORD(m.R_DATA_ADR), 4))+m.R_POINTS+m.R_DATA_BITPOS+15), 4))
				if last <= m.DATA_SIZE {
					bMask := iec.BYTE(1)
					idx1, idx2, bitPos := iec.INT(9), m.R_DATA_ADR, m.R_DATA_BITPOS
					for j := iec.INT(1); j <= m.R_POINTS; j++ {
						data.set(idx2, logic.BIT_LOAD_W(data.get(idx2), r[idx1]&bMask > 0, bitPos))
						bitPos++
						if bitPos > 15 {
							idx2++
							bitPos = 0
						}
						bMask = ROL(bMask, 1)
						if bMask == 1 {
							idx1++
						}
					}
				} else {
					m.ERROR = 2
				}
			case mask&0x0080_0018 > 0: // FC 3, 4, 23
				if m.R_DATA_ADR+m.R_POINTS <= m.DATA_SIZE {
					idx1 := iec.INT(7)
					for j := m.R_DATA_ADR; j <= m.R_DATA_ADR+m.R_POINTS-1; j++ {
						idx1 += 2
						data.set(j, logic.WORD_OF_BYTE(r[idx1], r[idx1+1]))
					}
				} else {
					m.ERROR = 2
				}
			case mask&0xFF3E_7F81 > 0:
				m.ERROR = 1 // ILLEGAL FUNCTION
			}
		default:
			m.ERROR = 3 // ILLEGAL DATA VALUE
		}
		m.ipState = 4
		m.state = 0
		rb.SIZE = 0
	}
	m.runFIFO(now)
	m.BUSY = m.ipState == 3
}

// server is MB_SERVER and MB_SERVER_1, which differ in the watch of the
// receive, the time an error lasts and whether it stops the server.
type server struct {
	shortClient
	VMAP      *VMAP
	DATA_SIZE iec.INT
	ENABLE    iec.BOOL
	UDP       iec.BOOL
	ERROR     iec.DWORD

	vmapFB    MB_VMAP
	state     iec.INT
	lastCycle iec.DWORD
	t         timers.TON
	ipError   iec.BOOL
}

// run runs a scan of the server.
func (m *server) run(now time.Time, observe bool, errorTime time.Duration, stopOnError bool) {
	if !m.ok() || m.VMAP == nil {
		return
	}
	c, sb, rb := m.IP_C, m.S_BUF, m.R_BUF
	s, r := sb.BUFFER[:], rb.BUFFER[:]

	// The time since a write, of each entry of the map.
	tx := PLC_MS(now)
	for i := range m.VMAP {
		if m.VMAP[i].TIME_OUT > 0 {
			m.VMAP[i].TIME_OUT += DWORD_TO_TIME(tx - m.lastCycle)
		}
	}
	m.lastCycle = tx
	m.ipError = c.ERROR > 0

	switch m.state {
	case 0:
		if m.ENABLE {
			m.state = 10
			m.ipState = 1
		}
	case 10:
		if m.ipState == 3 {
			c.C_PORT = 0 // the port of the IP_CONTROL2
			c.C_IP = 0
			c.C_MODE = SEL[iec.BYTE](m.UDP, 4, 5) // a TCP or UDP server for any client
			c.TIME_RESET = true
			c.C_ENABLE = true
			c.R_OBSERVE = iec.BOOL(observe)
			m.state = 20
		}
	case 20:
		if m.ipError {
			m.ERROR = c.ERROR
		} else if sb.SIZE == 0 && rb.SIZE > 6 {
			// A request of the length it says; OSCAT answers one of another
			// length with the last exception, if any.
			if rb.SIZE == iec.UINT(word(r[4], r[5]))+6 {
				m.request(now)
			}
			if m.ERROR > 0 {
				s[7] |= 0x80 // an exception
				s[8] = iec.BYTE(m.ERROR)
				sb.SIZE = 9
			}
			s[4] = 0
			s[5] = iec.BYTE(sb.SIZE - 6) // the length, its low byte
		}
		m.t.IN, m.t.PT = m.ipError, iec.TIME(errorTime)
		m.t.Execute(now)
		if m.t.Q {
			c.TIME_RESET = true // reset the error
		}
		rb.SIZE = 0
		if !bool(m.ENABLE) || stopOnError && bool(m.t.Q) {
			m.ipState = 4
			m.state = 0
			if stopOnError {
				c.C_ENABLE = false
			}
		}
	}
	m.runFIFO(now)
}

// vmap maps an address with MB_VMAP: the address, the bit and the error.
func (m *server) vmap(now time.Time, fc, adr, cnt iec.INT) (iec.INT, iec.INT, iec.DWORD) {
	m.vmapFB.VMAP, m.vmapFB.FC, m.vmapFB.V_ADR, m.vmapFB.V_CNT, m.vmapFB.SIZE = m.VMAP, fc, adr, cnt, m.DATA_SIZE
	m.vmapFB.Execute(now)
	return m.vmapFB.P_ADR, m.vmapFB.P_BIT, iec.DWORD(m.vmapFB.ERROR)
}

// request answers the request in R_BUF into S_BUF.
func (m *server) request(now time.Time) {
	sb, rb, data := m.S_BUF, m.R_BUF, area{m.DATA}
	s, r := sb.BUFFER[:], rb.BUFFER[:]
	m.ERROR = 0
	copy(s[:14], r[:14]) // the header
	fc := iec.INT(r[7])
	adr1, points := word(r[8], r[9]), word(r[10], r[11])
	var bitPos iec.INT
	adr1, bitPos, m.ERROR = m.vmap(now, fc, adr1, points)
	if m.ERROR != 0 {
		return
	}
	switch fc {
	case 1, 2: // Read Coil Status, Read Input Status
		count := iec.INT(9)
		coils, mask := iec.BYTE(0), iec.BYTE(1)
		for i := iec.INT(1); i <= points; i++ {
			if mask == 0 {
				mask = 1
				s[count] = coils
				coils = 0
				count++
			}
			if data.get(adr1)&SHL(iec.WORD(1), bitPos) > 0 {
				coils |= mask
			}
			bitPos++
			if bitPos > 15 {
				adr1++
				bitPos = 0
			}
			mask = SHL(mask, 1)
		}
		s[count] = coils
		s[8] = iec.BYTE(count - 8)
		sb.SIZE = iec.UINT(count + 1)
	case 3, 4: // Read Holding Registers, Read Input Registers
		idx1 := iec.INT(7)
		for i := adr1; i <= adr1+points-1; i++ {
			idx1 += 2
			putWord(s, idx1, data.get(i))
		}
		s[8] = SHL(iec.BYTE(points), 1)
		sb.SIZE = iec.UINT(idx1 + 2)
	case 5: // Force Single Coil
		data.set(adr1, logic.BIT_LOAD_W(data.get(adr1), r[10] > 0, bitPos))
		sb.SIZE = 12
	case 6: // Preset Single Register
		data.set(adr1, logic.WORD_OF_BYTE(r[10], r[11]))
		sb.SIZE = 12
	case 15: // Force Multiple Coils
		mask, idx1 := iec.BYTE(1), iec.INT(13)
		for i := iec.INT(1); i <= points; i++ {
			data.set(adr1, logic.BIT_LOAD_W(data.get(adr1), r[idx1]&mask > 0, bitPos))
			bitPos++
			if bitPos > 15 {
				adr1++
				bitPos = 0
			}
			mask = ROL(mask, 1)
			if mask == 1 {
				idx1++
			}
		}
		sb.SIZE = 12
	case 16: // Preset Multiple Registers
		idx1 := iec.INT(11)
		for i := adr1; i <= adr1+points-1; i++ {
			idx1 += 2
			data.set(i, logic.WORD_OF_BYTE(r[idx1], r[idx1+1]))
		}
		sb.SIZE = 12
	case 22: // Mask Write Register
		and := logic.WORD_OF_BYTE(r[10], r[11])
		data.set(adr1, data.get(adr1)&and|logic.WORD_OF_BYTE(r[12], r[13])&^and)
		sb.SIZE = 14
	case 23: // Read/Write Registers: the write first
		adr2, points2 := word(r[12], r[13]), word(r[14], r[15])
		adr2, _, m.ERROR = m.vmap(now, 16, adr2, points2)
		if m.ERROR != 0 {
			return
		}
		idx1 := iec.INT(15)
		for i := adr2; i <= adr2+points2-1; i++ {
			idx1 += 2
			data.set(i, logic.WORD_OF_BYTE(r[idx1], r[idx1+1]))
		}
		idx1 = 7
		for i := adr1; i <= adr1+points-1; i++ {
			idx1 += 2
			putWord(s, idx1, data.get(i))
		}
		s[8] = SHL(iec.BYTE(points), 1)
		sb.SIZE = iec.UINT(idx1 + 2)
	default:
		m.ERROR = 1 // ILLEGAL FUNCTION
	}
}

// MB_SERVER is a Modbus server, over TCP or UDP, while ENABLE, on the port
// of the IP_CONTROL2 of IP_C, of the data area DATA, of DATA_SIZE words,
// through the map VMAP; see MB_VMAP. It watches the receive, and an error
// of 2 s stops it. ERROR is the last exception, or the connection's.
type MB_SERVER struct {
	server
}

// INIT resets the block.
func (m *MB_SERVER) INIT() {
	*m = MB_SERVER{server{shortClient: shortClient{IP_C: m.IP_C, S_BUF: m.S_BUF, R_BUF: m.R_BUF, DATA: m.DATA}, VMAP: m.VMAP}}
}

// Execute runs the block once.
func (m *MB_SERVER) Execute(now time.Time) { m.run(now, true, 2*time.Second, true) }

// MB_SERVER_1 is MB_SERVER that does not watch the receive, and resets an
// error after 5 s without stopping.
type MB_SERVER_1 struct {
	server
}

// INIT resets the block.
func (m *MB_SERVER_1) INIT() {
	*m = MB_SERVER_1{server{shortClient: shortClient{IP_C: m.IP_C, S_BUF: m.S_BUF, R_BUF: m.R_BUF, DATA: m.DATA}, VMAP: m.VMAP}}
}

// Execute runs the block once.
func (m *MB_SERVER_1) Execute(now time.Time) { m.run(now, false, 5*time.Second, false) }
