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

// Package tcpip emulates the function blocks of Beckhoff's TwinCAT TCP/IP
// library (Tc2_TcpIp) that OSCAT NETWORK uses, on Go's net package.
//
// As on TwinCAT, a block starts its operation on a rising edge of BEXECUTE
// and reports it done when BBUSY is false, with BERROR and NERRID; the
// outputs keep their values until the next rising edge. A socket is a
// handle, T_HSOCKET, into a table of the program's sockets. A connected
// socket reads what arrives in the background, so FB_SocketReceive,
// FB_SocketUdpReceiveFrom and FB_SocketAccept return at once with what has
// arrived, possibly nothing, as TwinCAT's do.
//
// The operations are real network I/O: their time outs are in wall clock
// time, unlike the time of the blocks' Execute(now). The TwinCAT TCP/IP
// connection server SSRVNETID is ignored.
package tcpip

import (
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// T_SOCKADDR is the address of a socket.
type T_SOCKADDR struct {
	SADDR iec.STRING // STRING(15): the IPv4 address
	NPORT iec.UDINT  // the port
}

// T_HSOCKET is the handle of a socket.
type T_HSOCKET struct {
	HANDLE     iec.UDINT
	REMOTEADDR T_SOCKADDR
	LOCALADDR  T_SOCKADDR
}

// TCPADS_NULL_HSOCKET is the handle of no socket.
var TCPADS_NULL_HSOCKET T_HSOCKET

// The error numbers the blocks report: the Windows socket errors, as
// TwinCAT does, and the closing of a connection by the remote.
const (
	// ERR_CLOSED_BY_REMOTE is the error of a receive on a TCP connection
	// the remote closed.
	ERR_CLOSED_BY_REMOTE iec.UDINT = 0x8004
	// ERR_INVALID_HANDLE is the error of an operation on no socket.
	ERR_INVALID_HANDLE iec.UDINT = 0x8007
	// ERR_CONN_RESET is WSAECONNRESET: the remote reset the connection,
	// or a UDP datagram was refused.
	ERR_CONN_RESET iec.UDINT = 0x80072746
	// ERR_TIMED_OUT is WSAETIMEDOUT.
	ERR_TIMED_OUT iec.UDINT = 0x8007274C
	// ERR_CONN_REFUSED is WSAECONNREFUSED.
	ERR_CONN_REFUSED iec.UDINT = 0x8007274D
	// ERR_ADDR_IN_USE is WSAEADDRINUSE.
	ERR_ADDR_IN_USE iec.UDINT = 0x80072740
	// ERR_HOST_UNREACH is WSAEHOSTUNREACH.
	ERR_HOST_UNREACH iec.UDINT = 0x80072751
	// ERR_FAILED is any other error.
	ERR_FAILED iec.UDINT = 0x80072745
)

// errID returns the error number of err.
func errID(err error) iec.UDINT {
	var ne net.Error
	switch {
	case err == nil:
		return 0
	case errors.Is(err, io.EOF):
		return ERR_CLOSED_BY_REMOTE
	case errors.Is(err, syscall.ECONNRESET):
		return ERR_CONN_RESET
	case errors.Is(err, syscall.ECONNREFUSED):
		return ERR_CONN_REFUSED
	case errors.Is(err, syscall.EADDRINUSE):
		return ERR_ADDR_IN_USE
	case errors.Is(err, syscall.EHOSTUNREACH), errors.Is(err, syscall.ENETUNREACH):
		return ERR_HOST_UNREACH
	case errors.Is(err, os.ErrDeadlineExceeded), errors.As(err, &ne) && ne.Timeout():
		return ERR_TIMED_OUT
	}
	return ERR_FAILED
}

// packet is a UDP datagram received.
type packet struct {
	data []byte
	from *net.UDPAddr
}

// socket is a socket of the table: a TCP connection, a UDP socket or a TCP
// listener, and what has arrived on it.
type socket struct {
	conn net.Conn
	udp  *net.UDPConn
	ln   net.Listener

	mu       sync.Mutex
	data     []byte     // the TCP data received
	packets  []packet   // the UDP datagrams received
	accepted []net.Conn // the TCP connections accepted
	err      error      // the error that ended the reading
}

var (
	tableMu sync.Mutex
	table   = map[iec.UDINT]*socket{}
	next    iec.UDINT
)

// add puts s in the table and starts reading on it.
func add(s *socket) iec.UDINT {
	tableMu.Lock()
	next++
	h := next
	table[h] = s
	tableMu.Unlock()
	go s.read()
	return h
}

func lookup(h iec.UDINT) *socket {
	tableMu.Lock()
	defer tableMu.Unlock()
	return table[h]
}

// remove takes the socket h out of the table and closes it.
func remove(h iec.UDINT) error {
	tableMu.Lock()
	s := table[h]
	delete(table, h)
	tableMu.Unlock()
	if s == nil {
		return errors.New("no such socket")
	}
	return s.close()
}

func (s *socket) close() error {
	switch {
	case s.conn != nil:
		return s.conn.Close()
	case s.udp != nil:
		return s.udp.Close()
	case s.ln != nil:
		return s.ln.Close()
	}
	return nil
}

// read reads what arrives on the socket until it fails.
func (s *socket) read() {
	buf := make([]byte, 65536)
	for {
		var err error
		switch {
		case s.conn != nil:
			var n int
			n, err = s.conn.Read(buf)
			s.mu.Lock()
			s.data = append(s.data, buf[:n]...)
			s.mu.Unlock()
		case s.udp != nil:
			var n int
			var from *net.UDPAddr
			n, from, err = s.udp.ReadFromUDP(buf)
			if err == nil {
				s.mu.Lock()
				s.packets = append(s.packets, packet{append([]byte(nil), buf[:n]...), from})
				s.mu.Unlock()
			} else if errors.Is(err, syscall.ECONNRESET) {
				// A datagram refused: Windows reports it on the next read.
				s.mu.Lock()
				s.err = err
				s.mu.Unlock()
				continue
			}
		case s.ln != nil:
			var c net.Conn
			c, err = s.ln.Accept()
			if err == nil {
				s.mu.Lock()
				s.accepted = append(s.accepted, c)
				s.mu.Unlock()
			}
		default:
			return
		}
		if err != nil {
			s.mu.Lock()
			s.err = err
			s.mu.Unlock()
			return
		}
	}
}

// addr returns the address of a, as a T_SOCKADDR.
func addr(a net.Addr) T_SOCKADDR {
	host, port, err := net.SplitHostPort(a.String())
	if err != nil {
		return T_SOCKADDR{}
	}
	p, _ := strconv.Atoi(port)
	if ip := net.ParseIP(host); ip != nil && ip.To4() != nil {
		host = ip.To4().String()
	}
	return T_SOCKADDR{iec.STRING(host), iec.UDINT(p)}
}

// handle returns the handle of the socket s in the table at h.
func handle(h iec.UDINT, s *socket) T_HSOCKET {
	hs := T_HSOCKET{HANDLE: h}
	switch {
	case s.conn != nil:
		hs.LOCALADDR, hs.REMOTEADDR = addr(s.conn.LocalAddr()), addr(s.conn.RemoteAddr())
	case s.udp != nil:
		hs.LOCALADDR = addr(s.udp.LocalAddr())
	case s.ln != nil:
		hs.LOCALADDR = addr(s.ln.Addr())
	}
	return hs
}

// hostPort returns the address host:port, ” meaning any address.
func hostPort(host iec.STRING, port iec.UDINT) string {
	return net.JoinHostPort(string(host), strconv.Itoa(int(port)))
}

// timeout returns the time out t, or a minute if it is 0.
func timeout(t iec.TIME) time.Duration {
	if t <= 0 {
		return time.Minute
	}
	return time.Duration(t)
}

// command is what the blocks share: the edge of BEXECUTE, the operation
// running and its outputs.
type command struct {
	BEXECUTE iec.BOOL
	TTIMEOUT iec.TIME
	BBUSY    iec.BOOL
	BERROR   iec.BOOL
	NERRID   iec.UDINT

	last iec.BOOL
	done chan func()
}

// edge reports a rising edge of BEXECUTE.
func (c *command) edge() bool {
	e := bool(c.BEXECUTE && !c.last)
	c.last = c.BEXECUTE
	return e
}

// start runs op in the background: the outputs are set from what it
// returns when the block runs after it ends.
func (c *command) start(op func() (func(), error)) {
	c.BBUSY, c.BERROR, c.NERRID = true, false, 0
	c.done = make(chan func(), 1)
	go func() {
		set, err := op()
		c.done <- func() {
			if set != nil {
				set()
			}
			c.BBUSY = false
			c.BERROR = err != nil
			c.NERRID = errID(err)
		}
	}()
}

// finish sets the outputs at once, of an operation that does not wait.
func (c *command) finish(err error) {
	c.done = nil
	c.BBUSY = false
	c.BERROR = err != nil
	c.NERRID = errID(err)
}

// poll sets the outputs of the operation if it has ended.
func (c *command) poll() {
	if c.done == nil {
		return
	}
	select {
	case set := <-c.done:
		c.done = nil
		set()
	default:
	}
}

// FB_SocketConnect connects to the TCP server SREMOTEHOST:NREMOTEPORT; the
// connection is HSOCKET.
type FB_SocketConnect struct {
	command
	SSRVNETID   iec.STRING
	SREMOTEHOST iec.STRING
	NREMOTEPORT iec.UDINT
	HSOCKET     T_HSOCKET
}

// INIT resets the block.
func (f *FB_SocketConnect) INIT() { *f = FB_SocketConnect{} }

// Execute runs the block once.
func (f *FB_SocketConnect) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	target, t := hostPort(f.SREMOTEHOST, f.NREMOTEPORT), timeout(f.TTIMEOUT)
	f.start(func() (func(), error) {
		c, err := net.DialTimeout("tcp4", target, t)
		if err != nil {
			return nil, err
		}
		s := &socket{conn: c}
		h := add(s)
		return func() { f.HSOCKET = handle(h, s) }, nil
	})
}

// FB_SocketClose closes the socket HSOCKET.
type FB_SocketClose struct {
	command
	SSRVNETID iec.STRING
	HSOCKET   T_HSOCKET
}

// INIT resets the block.
func (f *FB_SocketClose) INIT() { *f = FB_SocketClose{} }

// Execute runs the block once.
func (f *FB_SocketClose) Execute(now time.Time) {
	f.poll()
	if f.edge() {
		f.finish(remove(f.HSOCKET.HANDLE))
	}
}

// FB_SocketCloseAll closes all the sockets of the program.
type FB_SocketCloseAll struct {
	command
	SSRVNETID iec.STRING
}

// INIT resets the block.
func (f *FB_SocketCloseAll) INIT() { *f = FB_SocketCloseAll{} }

// Execute runs the block once.
func (f *FB_SocketCloseAll) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	tableMu.Lock()
	all := table
	table = map[iec.UDINT]*socket{}
	tableMu.Unlock()
	for _, s := range all {
		s.close()
	}
	f.finish(nil)
}

// FB_SocketListen opens a TCP listener HLISTENER on SLOCALHOST:NLOCALPORT,
// SLOCALHOST ” for every address.
type FB_SocketListen struct {
	command
	SSRVNETID  iec.STRING
	SLOCALHOST iec.STRING
	NLOCALPORT iec.UDINT
	HLISTENER  T_HSOCKET
}

// INIT resets the block.
func (f *FB_SocketListen) INIT() { *f = FB_SocketListen{} }

// Execute runs the block once.
func (f *FB_SocketListen) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	ln, err := net.Listen("tcp4", hostPort(f.SLOCALHOST, f.NLOCALPORT))
	if err == nil {
		s := &socket{ln: ln}
		f.HLISTENER = handle(add(s), s)
	}
	f.finish(err)
}

// FB_SocketAccept takes a connection HSOCKET a client made to the listener
// HLISTENER, if there is one: BACCEPTED tells.
type FB_SocketAccept struct {
	command
	SSRVNETID iec.STRING
	HLISTENER T_HSOCKET
	BACCEPTED iec.BOOL
	HSOCKET   T_HSOCKET
}

// INIT resets the block.
func (f *FB_SocketAccept) INIT() { *f = FB_SocketAccept{} }

// Execute runs the block once.
func (f *FB_SocketAccept) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	f.BACCEPTED = false
	s := lookup(f.HLISTENER.HANDLE)
	if s == nil || s.ln == nil {
		f.finish(errors.New("no such listener"))
		f.NERRID = ERR_INVALID_HANDLE
		return
	}
	s.mu.Lock()
	var c net.Conn
	if len(s.accepted) > 0 {
		c, s.accepted = s.accepted[0], s.accepted[1:]
	}
	err := s.err
	s.mu.Unlock()
	if c == nil {
		f.finish(err)
		return
	}
	cs := &socket{conn: c}
	f.HSOCKET = handle(add(cs), cs)
	f.BACCEPTED = true
	f.finish(nil)
}

// FB_SocketSend sends the CBLEN bytes of SRC on the TCP connection HSOCKET.
type FB_SocketSend struct {
	command
	SSRVNETID iec.STRING
	HSOCKET   T_HSOCKET
	CBLEN     iec.UDINT
	SRC       []iec.BYTE
}

// INIT resets the block.
func (f *FB_SocketSend) INIT() { *f = FB_SocketSend{} }

// Execute runs the block once.
func (f *FB_SocketSend) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	s := lookup(f.HSOCKET.HANDLE)
	if s == nil || s.conn == nil {
		f.finish(errors.New("no such connection"))
		f.NERRID = ERR_INVALID_HANDLE
		return
	}
	data := bytesOf(f.SRC, f.CBLEN)
	t := timeout(f.TTIMEOUT)
	f.start(func() (func(), error) {
		s.conn.SetWriteDeadline(time.Now().Add(t))
		_, err := s.conn.Write(data)
		return nil, err
	})
}

// bytesOf returns a copy of the first n bytes of b.
func bytesOf(b []iec.BYTE, n iec.UDINT) []byte {
	if int(n) > len(b) {
		n = iec.UDINT(len(b))
	}
	out := make([]byte, n)
	for i := range out {
		out[i] = byte(b[i])
	}
	return out
}

// FB_SocketReceive takes up to CBLEN bytes received on the TCP connection
// HSOCKET into DEST: NRECBYTES of them, 0 if none has arrived. A
// connection the remote closed is the error ERR_CLOSED_BY_REMOTE.
type FB_SocketReceive struct {
	command
	SSRVNETID iec.STRING
	HSOCKET   T_HSOCKET
	CBLEN     iec.UDINT
	DEST      []iec.BYTE
	NRECBYTES iec.UDINT
}

// INIT resets the block.
func (f *FB_SocketReceive) INIT() { *f = FB_SocketReceive{} }

// Execute runs the block once.
func (f *FB_SocketReceive) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	f.NRECBYTES = 0
	s := lookup(f.HSOCKET.HANDLE)
	if s == nil || s.conn == nil {
		f.finish(errors.New("no such connection"))
		f.NERRID = ERR_INVALID_HANDLE
		return
	}
	s.mu.Lock()
	n := min(int(f.CBLEN), len(f.DEST), len(s.data))
	for i := 0; i < n; i++ {
		f.DEST[i] = iec.BYTE(s.data[i])
	}
	s.data = s.data[n:]
	err := s.err
	if len(s.data) > 0 || n > 0 {
		err = nil
	}
	s.mu.Unlock()
	f.NRECBYTES = iec.UDINT(n)
	f.finish(err)
}

// FB_SocketUdpCreate opens a UDP socket HSOCKET on SLOCALHOST:NLOCALPORT,
// port 0 for any.
type FB_SocketUdpCreate struct {
	command
	SSRVNETID  iec.STRING
	SLOCALHOST iec.STRING
	NLOCALPORT iec.UDINT
	HSOCKET    T_HSOCKET
}

// INIT resets the block.
func (f *FB_SocketUdpCreate) INIT() { *f = FB_SocketUdpCreate{} }

// Execute runs the block once.
func (f *FB_SocketUdpCreate) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	a, err := net.ResolveUDPAddr("udp4", hostPort(f.SLOCALHOST, f.NLOCALPORT))
	if err == nil {
		var c *net.UDPConn
		if c, err = net.ListenUDP("udp4", a); err == nil {
			s := &socket{udp: c}
			f.HSOCKET = handle(add(s), s)
		}
	}
	f.finish(err)
}

// FB_SocketUdpSendTo sends the CBLEN bytes of SRC from the UDP socket
// HSOCKET to SREMOTEHOST:NREMOTEPORT.
type FB_SocketUdpSendTo struct {
	command
	SSRVNETID   iec.STRING
	HSOCKET     T_HSOCKET
	SREMOTEHOST iec.STRING
	NREMOTEPORT iec.UDINT
	CBLEN       iec.UDINT
	SRC         []iec.BYTE
}

// INIT resets the block.
func (f *FB_SocketUdpSendTo) INIT() { *f = FB_SocketUdpSendTo{} }

// Execute runs the block once.
func (f *FB_SocketUdpSendTo) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	s := lookup(f.HSOCKET.HANDLE)
	if s == nil || s.udp == nil {
		f.finish(errors.New("no such socket"))
		f.NERRID = ERR_INVALID_HANDLE
		return
	}
	a, err := net.ResolveUDPAddr("udp4", hostPort(f.SREMOTEHOST, f.NREMOTEPORT))
	if err == nil {
		_, err = s.udp.WriteToUDP(bytesOf(f.SRC, f.CBLEN), a)
	}
	f.finish(err)
}

// FB_SocketUdpReceiveFrom takes a datagram received on the UDP socket
// HSOCKET into DEST, up to CBLEN bytes: NRECBYTES of them, 0 if none has
// arrived, from SREMOTEHOST:NREMOTEPORT.
type FB_SocketUdpReceiveFrom struct {
	command
	SSRVNETID   iec.STRING
	HSOCKET     T_HSOCKET
	CBLEN       iec.UDINT
	DEST        []iec.BYTE
	SREMOTEHOST iec.STRING
	NREMOTEPORT iec.UDINT
	NRECBYTES   iec.UDINT
}

// INIT resets the block.
func (f *FB_SocketUdpReceiveFrom) INIT() { *f = FB_SocketUdpReceiveFrom{} }

// Execute runs the block once.
func (f *FB_SocketUdpReceiveFrom) Execute(now time.Time) {
	f.poll()
	if !f.edge() {
		return
	}
	f.NRECBYTES = 0
	s := lookup(f.HSOCKET.HANDLE)
	if s == nil || s.udp == nil {
		f.finish(errors.New("no such socket"))
		f.NERRID = ERR_INVALID_HANDLE
		return
	}
	s.mu.Lock()
	var p packet
	ok := len(s.packets) > 0
	if ok {
		p, s.packets = s.packets[0], s.packets[1:]
	}
	err := s.err
	s.err = nil
	s.mu.Unlock()
	if !ok {
		f.finish(err)
		return
	}
	n := min(int(f.CBLEN), len(f.DEST), len(p.data))
	for i := 0; i < n; i++ {
		f.DEST[i] = iec.BYTE(p.data[i])
	}
	f.NRECBYTES = iec.UDINT(n)
	a := addr(p.from)
	f.SREMOTEHOST, f.NREMOTEPORT = a.SADDR, a.NPORT
	f.finish(nil)
}
