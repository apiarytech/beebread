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

// Package ip is the port of the OSCAT NETWORK connection management:
// IP_CONTROL keeps a TCP or UDP connection, which the blocks that share it
// take in turn through IP_FIFO, for Beckhoff TwinCAT's TCP/IP library,
// emulated by the package tcpip.
package ip

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/logic"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/tcpip"
	"github.com/apiarytech/royaljelly/iec"
)

// IP_CONTROL_RESET closes all the sockets once, at the start of the
// program, for the first IP_CONTROL that runs: READY is true after.
type IP_CONTROL_RESET struct {
	READY iec.BOOL

	closeAll tcpip.FB_SocketCloseAll
	step     iec.BYTE
}

// INIT resets the block.
func (r *IP_CONTROL_RESET) INIT() { *r = IP_CONTROL_RESET{} }

// Execute runs the block once.
func (r *IP_CONTROL_RESET) Execute(now time.Time) {
	switch network.TCP_SERVER_RESET {
	case 2:
		r.READY = true
		return
	case 1:
		if r.step == 0 {
			return
		}
	case 0:
		network.TCP_SERVER_RESET = 1 // the reset runs
		r.step = 1
	}
	switch r.step {
	case 1:
		r.closeAll.SSRVNETID, r.closeAll.BEXECUTE, r.closeAll.TTIMEOUT = "", true, iec.TIME(10*time.Second)
		r.closeAll.Execute(now)
		r.step = 2
	case 2:
		r.closeAll.BEXECUTE = false
		r.closeAll.Execute(now)
		if !r.closeAll.BBUSY {
			r.step = 0
			network.TCP_SERVER_RESET = 2 // ready
		}
	}
}

// The states of IP_CONTROL.
const (
	stateStop         iec.BYTE = 0
	stateUDPInit      iec.BYTE = 31
	stateUDPInitWait  iec.BYTE = 32
	stateTCInit       iec.BYTE = 51
	stateTCConnect    iec.BYTE = 52
	stateTSInit       iec.BYTE = 53
	stateTSListenWait iec.BYTE = 54
	stateTSAccept     iec.BYTE = 55
	stateTSAcceptWait iec.BYTE = 56
	stateClose        iec.BYTE = 190
	stateCloseWait    iec.BYTE = 191
	stateClose2       iec.BYTE = 192
	stateCloseWait2   iec.BYTE = 193
	stateWait         iec.BYTE = 200
	stateRec          iec.BYTE = 210
	stateRecWait      iec.BYTE = 211
	stateSnd          iec.BYTE = 220
	stateSndWait1     iec.BYTE = 221
	stateSndWait2     iec.BYTE = 222
)

// buffers are the send and receive buffers of an IP_CONTROL: the size and
// the data of each.
type buffers struct {
	sSize *iec.UINT
	sBuf  []iec.BYTE
	rSize *iec.UINT
	rBuf  []iec.BYTE
}

// control is IP_CONTROL, for buffers of either size.
type control struct {
	cTime, sTime, rTime iec.UDINT

	cEnable    iec.BOOL
	cIP        iec.DWORD
	cIPStr     iec.STRING
	cPort      iec.WORD
	cMode      iec.BYTE
	cStatus    iec.BYTE
	cReady     iec.BOOL
	cReadyOld  iec.BOOL
	sTotal     iec.INT
	sCurPos    iec.INT
	sCurSize   iec.INT
	sMaxSize   iec.INT
	sActive    iec.BOOL
	sStatus    iec.BYTE
	sIPStr     iec.STRING
	sIPExist   iec.BOOL
	rStatus    iec.BYTE
	rCount     iec.INT
	rOffset    iec.INT
	rMaxSize   iec.INT
	newConn    iec.BOOL
	tx         iec.DWORD
	errorTime  iec.DWORD
	state      iec.BYTE
	bytesRecv  iec.DINT
	udpMode    iec.BOOL
	socket     tcpip.T_HSOCKET
	server     tcpip.T_HSOCKET
	nErrID     iec.UDINT
	sendBusy   iec.BOOL
	sendError  iec.BOOL
	sendErrID  iec.UDINT
	ipcr       IP_CONTROL_RESET
	create     tcpip.FB_SocketUdpCreate
	connect    tcpip.FB_SocketConnect
	close      tcpip.FB_SocketClose
	sendTo     tcpip.FB_SocketUdpSendTo
	receiveFrm tcpip.FB_SocketUdpReceiveFrom
	send       tcpip.FB_SocketSend
	receive    tcpip.FB_SocketReceive
	listen     tcpip.FB_SocketListen
	accept     tcpip.FB_SocketAccept
}

// run runs IP_CONTROL once with the interface c, the buffers b and the
// inputs.
func (x *control) run(now time.Time, c *network.IP_C, b buffers, IP iec.DWORD, PORT iec.WORD, TIME_OUT iec.TIME) {
	srv := network.SSRVNETID
	// Close the old connections at the start of the program.
	x.ipcr.Execute(now)
	if !x.ipcr.READY {
		return
	}

	x.tx = PLC_MS(now)
	if c.C_PORT == 0 {
		c.C_PORT = PORT
	}
	if c.C_IP == 0 && c.C_MODE < 4 {
		c.C_IP = IP
	}

	// A new connection wanted?
	x.newConn = x.cIP != c.C_IP || x.cMode != c.C_MODE || x.cPort != c.C_PORT
	x.cEnable = c.C_ENABLE && !x.newConn

	if c.C_ENABLE && c.ERROR == 0 && x.state == stateStop {
		x.cIP, x.cMode, x.cPort = c.C_IP, c.C_MODE, c.C_PORT
		x.cIPStr = encoding.IP4_TO_STRING(x.cIP)
		x.udpMode = iec.BOOL(BIT(x.cMode, 0)) // modes 1, 3, 5 are UDP
		x.errorTime = TIME_TO_DWORD(max(iec.TIME(200*time.Millisecond), TIME_OUT))
		c.TIME_RESET = true
		x.rMaxSize = iec.INT(len(b.rBuf))
		x.sMaxSize = iec.INT(len(b.sBuf))
		switch {
		case bool(x.udpMode): // UDP client or server
			x.state = stateUDPInit
		case x.cMode == 0: // TCP client
			x.state = stateTCInit
		default: // TCP server, modes 2 and 4
			x.state = stateTSInit
		}
		// The receiver, unless mode 5 answers whoever sends.
		if x.cMode < 5 {
			x.sIPStr, x.sIPExist = x.cIPStr, true
		} else {
			x.sIPStr, x.sIPExist = "", false
		}
	}

	if c.TIME_RESET {
		c.TIME_RESET = false
		x.cTime, x.sTime, x.rTime = iec.UDINT(x.tx), iec.UDINT(x.tx), iec.UDINT(x.tx)
		x.cStatus, x.sStatus, x.rStatus = 0, 0, 0
	}

	switch x.state {
	case stateUDPInit:
		// A UDP client takes any local port, a server its port.
		x.create.BEXECUTE = false
		x.create.Execute(now)
		x.create.BEXECUTE, x.create.SSRVNETID, x.create.SLOCALHOST = true, srv, network.SLOCALHOST
		x.create.NLOCALPORT = iec.UDINT(SEL(x.cMode == 1, x.cPort, 0))
		x.create.TTIMEOUT = TIME_OUT
		x.create.Execute(now)
		x.create.BEXECUTE = false
		x.state = stateUDPInitWait

	case stateUDPInitWait:
		x.create.Execute(now)
		if !x.create.BBUSY {
			if !x.create.BERROR {
				x.socket = x.create.HSOCKET
				x.cReady = true
				x.state = stateWait
			} else {
				x.cStatus = 1 // FB_SocketUdpCreate failed
				x.nErrID = x.create.NERRID
				x.state = stateClose
			}
		}

	case stateTCInit: // TCP client
		x.connect.BEXECUTE = false
		x.connect.Execute(now)
		x.connect.BEXECUTE, x.connect.SSRVNETID = true, srv
		x.connect.SREMOTEHOST, x.connect.NREMOTEPORT, x.connect.TTIMEOUT = x.cIPStr, iec.UDINT(x.cPort), TIME_OUT
		x.connect.Execute(now)
		x.connect.BEXECUTE = false
		x.state = stateTCConnect

	case stateTCConnect:
		if !x.cEnable || x.cStatus == 255 {
			// disabled, or the connect timed out
			x.state = stateClose
		}
		x.connect.Execute(now)
		if !x.connect.BBUSY {
			if x.connect.BERROR {
				x.nErrID = x.connect.NERRID
				x.state = stateTCInit // try again
			} else {
				x.socket = x.connect.HSOCKET
				x.cReady = true
				x.state = stateWait
			}
		}

	case stateTSInit: // TCP server
		x.listen.BEXECUTE = false
		x.listen.Execute(now)
		x.listen.BEXECUTE, x.listen.SSRVNETID, x.listen.SLOCALHOST = true, srv, network.SLOCALHOST
		x.listen.NLOCALPORT, x.listen.TTIMEOUT = iec.UDINT(x.cPort), TIME_OUT
		x.listen.Execute(now)
		x.listen.BEXECUTE = false
		x.state = stateTSListenWait

	case stateTSListenWait:
		x.listen.Execute(now)
		if !x.listen.BBUSY {
			if x.listen.BERROR {
				x.cStatus = 3 // FB_SocketListen failed
				x.nErrID = x.listen.NERRID
				x.state = stateClose
			} else {
				x.server = x.listen.HLISTENER
				x.state = stateTSAccept
			}
		}

	case stateTSAccept:
		x.accept.BEXECUTE = false
		x.accept.Execute(now)
		x.accept.BEXECUTE, x.accept.SSRVNETID, x.accept.HLISTENER, x.accept.TTIMEOUT = true, srv, x.server, TIME_OUT
		x.accept.Execute(now)
		x.accept.BEXECUTE = false
		x.state = stateTSAcceptWait

	case stateTSAcceptWait:
		x.accept.Execute(now)
		if !x.accept.BBUSY {
			x.state = stateTSAccept
			switch {
			case bool(x.accept.BERROR):
				x.cStatus = 4 // FB_SocketAccept failed
				x.nErrID = x.accept.NERRID
				x.state = stateClose
			case bool(x.accept.BACCEPTED):
				x.socket = x.accept.HSOCKET
				if x.socket.REMOTEADDR.SADDR != x.cIPStr && x.cMode != 4 {
					x.state = stateClose // not the remote expected
				} else {
					x.state = stateWait
					x.cReady = true
				}
			}
		}

	case stateClose:
		x.cReady = false
		if x.socket.HANDLE > 0 {
			x.closeSocket(now, srv, x.socket)
			x.state = stateCloseWait
		} else {
			x.state = stateClose2
		}

	case stateCloseWait:
		x.close.Execute(now)
		if !x.close.BBUSY {
			x.socket = tcpip.TCPADS_NULL_HSOCKET
			x.nErrID = SEL(x.close.BERROR, 0, x.close.NERRID)
			x.state = stateClose2
		}

	case stateClose2:
		if x.server.HANDLE > 0 {
			x.closeSocket(now, srv, x.server)
			x.state = stateCloseWait2
		} else {
			x.state = stateStop
		}

	case stateCloseWait2:
		x.close.Execute(now)
		if !x.close.BBUSY {
			x.server = tcpip.TCPADS_NULL_HSOCKET
			x.nErrID = SEL(x.close.BERROR, 0, x.close.NERRID)
			x.state = stateStop
		}

	case stateWait: // wait for something to do
		switch {
		case bool(!x.cEnable || !x.cReady):
			x.state = stateClose
		case *b.sSize > 0 && c.MAILBOX[1] == 0 && bool(x.sIPExist):
			x.state = stateSnd
		case c.MAILBOX[2] == 0:
			x.state = stateRec
		}

	case stateRec: // receive
		if *b.rSize >= iec.UINT(x.rMaxSize) {
			*b.rSize = 0
			x.rStatus = 254 // the buffer overflowed and was reset
		}
		x.rOffset = SEL(x.udpMode, iec.INT(*b.rSize), 0)
		x.rCount = x.rMaxSize - x.rOffset
		x.bytesRecv = 0
		if x.udpMode {
			f := &x.receiveFrm
			f.BEXECUTE = false
			f.Execute(now)
			f.BEXECUTE, f.SSRVNETID, f.HSOCKET, f.CBLEN, f.DEST, f.TTIMEOUT = true, srv, x.socket, iec.UDINT(x.rCount), b.rBuf, TIME_OUT
			f.Execute(now)
			f.BEXECUTE = false
		} else {
			f := &x.receive
			f.BEXECUTE = false
			f.Execute(now)
			f.BEXECUTE, f.SSRVNETID, f.HSOCKET, f.CBLEN, f.DEST, f.TTIMEOUT = true, srv, x.socket, iec.UDINT(x.rCount), b.rBuf[x.rOffset:], TIME_OUT
			f.Execute(now)
			f.BEXECUTE = false
		}
		x.state = stateRecWait

	case stateRecWait:
		if x.udpMode {
			f := &x.receiveFrm
			f.Execute(now)
			if !f.BBUSY {
				switch {
				case bool(!f.BERROR):
					if f.NRECBYTES > 0 {
						if x.cMode != 5 && f.SREMOTEHOST != x.cIPStr {
							*b.rSize = 0
						} else {
							x.bytesRecv = iec.DINT(f.NRECBYTES)
							if !x.sIPExist {
								x.sIPStr, x.sIPExist = f.SREMOTEHOST, true
							}
						}
					}
				case f.NERRID == tcpip.ERR_CONN_RESET: // closed by the remote
					x.cStatus = 253
					x.cReady = false
				default:
					x.nErrID = f.NERRID
				}
				x.state = stateWait
			}
		} else {
			f := &x.receive
			f.Execute(now)
			if !f.BBUSY {
				switch {
				case bool(!f.BERROR):
					if f.NRECBYTES > 0 {
						x.bytesRecv = iec.DINT(f.NRECBYTES)
					}
				case f.NERRID == tcpip.ERR_CLOSED_BY_REMOTE:
					x.cStatus = 253
					x.cReady = false
				default:
					x.nErrID = f.NERRID
				}
				x.state = stateWait
			}
		}
		if x.bytesRecv > 0 {
			x.rTime = iec.UDINT(x.tx)
			*b.rSize = iec.UINT(x.rOffset) + iec.UINT(x.bytesRecv)
			c.MAILBOX[0]++ // the data received
			if c.MAILBOX[0] == 0 {
				c.MAILBOX[0] = 1
			}
		}

	case stateSnd: // send
		if *b.sSize > 0 && x.cReady && x.cEnable {
			if !x.sActive {
				x.sTotal = LIMIT(0, iec.INT(*b.sSize), x.rMaxSize)
				x.sCurPos, x.sCurSize = 0, 0
				x.sActive = true
			}
			x.sCurPos += x.sCurSize
			if x.sTotal > x.sCurPos {
				// more to send
				x.sCurSize = LIMIT(0, x.sTotal-x.sCurPos, x.sMaxSize)
				x.sTime, x.rTime = iec.UDINT(x.tx), iec.UDINT(x.tx)
				x.state = stateSndWait1
			} else {
				// all sent
				x.sActive = false
				*b.sSize = 0
				x.state = stateWait
			}
		}

	case stateSndWait1:
		src := b.sBuf[x.sCurPos:]
		if x.udpMode {
			f := &x.sendTo
			f.BEXECUTE = false
			f.Execute(now)
			f.BEXECUTE, f.SSRVNETID, f.HSOCKET, f.CBLEN, f.SRC = true, srv, x.socket, iec.UDINT(x.sCurSize), src
			f.SREMOTEHOST, f.NREMOTEPORT, f.TTIMEOUT = x.sIPStr, iec.UDINT(x.cPort), TIME_OUT
			f.Execute(now)
			f.BEXECUTE = false
		} else {
			f := &x.send
			f.BEXECUTE = false
			f.Execute(now)
			f.BEXECUTE, f.SSRVNETID, f.HSOCKET, f.CBLEN, f.SRC, f.TTIMEOUT = true, srv, x.socket, iec.UDINT(x.sCurSize), src, TIME_OUT
			f.Execute(now)
			f.BEXECUTE = false
		}
		x.state = stateSndWait2

	case stateSndWait2:
		if x.udpMode {
			x.sendTo.Execute(now)
			x.sendBusy, x.sendError, x.sendErrID = x.sendTo.BBUSY, x.sendTo.BERROR, x.sendTo.NERRID
		} else {
			x.send.Execute(now)
			x.sendBusy, x.sendError, x.sendErrID = x.send.BBUSY, x.send.BERROR, x.send.NERRID
		}
		if !x.sendBusy {
			if !x.sendError {
				x.state = stateSnd
			} else {
				x.nErrID = x.sendErrID
				x.state = stateWait
			}
		}
	}

	// The state of the connection, and its edges.
	switch {
	case bool(x.cReady && !x.cReadyOld): // connected
		c.C_STATE = 254
		c.TIME_RESET = true
	case bool(x.cReady):
		c.C_STATE = 255
	case bool(x.cReadyOld): // disconnected
		c.C_STATE = 1
		c.MAILBOX[0], c.MAILBOX[1], c.MAILBOX[2] = 0, 0, 0
		if !x.newConn {
			x.sActive = false
			*b.sSize = 0
		}
	default:
		c.C_STATE = 0
	}
	x.cReadyOld = x.cReady

	// The time outs.
	if !c.R_OBSERVE || x.sActive {
		x.rTime = iec.UDINT(x.tx)
	}
	if x.cStatus == 0 && x.tx-iec.DWORD(x.cTime) > x.errorTime && x.cEnable && !x.cReady && x.cMode < 2 {
		x.cStatus = 255 // connect time out
	}
	if x.sStatus == 0 && x.tx-iec.DWORD(x.sTime) > x.errorTime && x.sActive {
		x.sStatus = 255 // send time out
	}
	if x.rStatus == 0 && x.tx-iec.DWORD(x.rTime) > x.errorTime && x.cReady && c.R_OBSERVE {
		x.rStatus = 255 // receive time out
	}
	c.ERROR = logic.DWORD_OF_BYTE(x.cStatus, x.sStatus, x.rStatus, 0)
}

// closeSocket starts closing the socket h.
func (x *control) closeSocket(now time.Time, srv iec.STRING, h tcpip.T_HSOCKET) {
	x.close.BEXECUTE = false
	x.close.Execute(now)
	x.close.BEXECUTE, x.close.SSRVNETID, x.close.HSOCKET = true, srv, h
	x.close.Execute(now)
	x.close.BEXECUTE = false
}

// IP_CONTROL keeps the connection IP_C asks for: a TCP client or server, or
// UDP, with the remote IP and PORT unless IP_C has them, and sends S_BUF
// and receives into R_BUF. IP_C.ERROR has the errors of the connect, send
// and receive, in its bytes 3, 2 and 1: 255 a time out after TIME_OUT
// (200 ms at least).
type IP_CONTROL struct {
	IP_C     *network.IP_C
	S_BUF    *network.NETWORK_BUFFER
	R_BUF    *network.NETWORK_BUFFER
	IP       iec.DWORD
	PORT     iec.WORD
	TIME_OUT iec.TIME

	control
}

// INIT resets the block.
func (i *IP_CONTROL) INIT() {
	*i = IP_CONTROL{IP_C: i.IP_C, S_BUF: i.S_BUF, R_BUF: i.R_BUF}
}

// Execute runs the block once.
func (i *IP_CONTROL) Execute(now time.Time) {
	if i.IP_C == nil || i.S_BUF == nil || i.R_BUF == nil {
		return
	}
	i.run(now, i.IP_C, buffers{&i.S_BUF.SIZE, i.S_BUF.BUFFER[:], &i.R_BUF.SIZE, i.R_BUF.BUFFER[:]}, i.IP, i.PORT, i.TIME_OUT)
}

// IP_CONTROL2 is IP_CONTROL with short buffers.
type IP_CONTROL2 struct {
	IP_C     *network.IP_C
	S_BUF    *network.NETWORK_BUFFER_SHORT
	R_BUF    *network.NETWORK_BUFFER_SHORT
	IP       iec.DWORD
	PORT     iec.WORD
	TIME_OUT iec.TIME

	control
}

// INIT resets the block.
func (i *IP_CONTROL2) INIT() {
	*i = IP_CONTROL2{IP_C: i.IP_C, S_BUF: i.S_BUF, R_BUF: i.R_BUF}
}

// Execute runs the block once.
func (i *IP_CONTROL2) Execute(now time.Time) {
	if i.IP_C == nil || i.S_BUF == nil || i.R_BUF == nil {
		return
	}
	i.run(now, i.IP_C, buffers{&i.S_BUF.SIZE, i.S_BUF.BUFFER[:], &i.R_BUF.SIZE, i.R_BUF.BUFFER[:]}, i.IP, i.PORT, i.TIME_OUT)
}

// IP_FIFO gives the blocks that share an IP_C the connection in turn. A
// block with no ID gets one; it sets STATE to 1 to ask, gets 3 when it has
// the connection, and sets 4 to give it back, which makes it 5. STATE is
// 255 if there are more blocks than IDs.
type IP_FIFO struct {
	FIFO  *network.IP_FIFO_DATA
	ID    *iec.BYTE
	STATE *iec.BYTE
}

// INIT resets the block.
func (f *IP_FIFO) INIT() { *f = IP_FIFO{FIFO: f.FIFO, ID: f.ID, STATE: f.STATE} }

// Execute runs the block once.
func (f *IP_FIFO) Execute(now time.Time) {
	if f.FIFO == nil || f.ID == nil || f.STATE == nil {
		return
	}
	q, id, state := f.FIFO, f.ID, f.STATE
	if !q.INIT {
		q.INIT = true
		q.NW, q.NR = 1, 1
		q.EMPTY, q.FULL = true, false
		q.TOP = 128
		q.MAX_ID = 1 // the entries an ID may have
	}

	// A block with no ID gets the next one.
	if *id == 0 {
		if q.ID < iec.BYTE(q.TOP) {
			*id = q.ID + 1
			q.ID = *id
		} else if *state < 200 {
			// More blocks than IDs: every request fails.
			*state = 255
			return
		}
	}

	if *state == 1 && !q.FULL { // join the queue
		if q.Y[*id-1] < q.MAX_ID {
			q.Y[*id-1]++
			q.X[q.NW-1] = *id
			if q.NW == q.TOP {
				q.NW = 1
			} else {
				q.NW++
			}
			q.FULL = q.NW == q.NR
			q.EMPTY = false
			*state = 2 // wait for the turn
		}
	}
	if *state == 2 && !q.EMPTY && *id == q.X[q.NR-1] {
		*state = 3 // its turn
	}
	if *state == 4 && !q.EMPTY { // leave the queue
		if *id == q.X[q.NR-1] {
			q.Y[*id-1]--
			if q.NR == q.TOP {
				q.NR = 1
			} else {
				q.NR++
			}
			q.EMPTY = q.NR == q.NW
			q.FULL = false
		}
		*state = 5
	}
}
