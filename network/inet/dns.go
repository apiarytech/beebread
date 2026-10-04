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

// Package inet is the port of the OSCAT NETWORK clients and servers of the
// internet protocols: DNS, HTTP, FTP, SMTP, SNTP, syslog, telnet and
// MySQL, and the services built on them.
//
// A client shares a connection IP_C, and its buffers S_BUF and R_BUF, with
// the other clients and an IP_CONTROL: it takes the connection in turn
// with IP_FIFO, sets it up, sends its request and reads the answer.
package inet

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/buffer"
	"github.com/apiarytech/beebread/basic/logic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// client is what the clients share: the connection and its buffers, and
// their place in its queue.
type client struct {
	IP_C  *network.IP_C
	S_BUF *network.NETWORK_BUFFER
	R_BUF *network.NETWORK_BUFFER

	queue
}

// queue is a block's place in the queue of a connection: ipState is 1 to
// ask for the connection, 3 when it has it, and 4 to give it back.
type queue struct {
	fifo    ip.IP_FIFO
	ipState iec.BYTE
	ipID    iec.BYTE
}

// runFIFO runs the IP_FIFO of the queue of the connection c.
func (q *queue) runFIFO(now time.Time, c *network.IP_C) {
	q.fifo.FIFO, q.fifo.STATE, q.fifo.ID = &c.FIFO, &q.ipState, &q.ipID
	q.fifo.Execute(now)
}

// ok reports whether the client has its connection and buffers.
func (c *client) ok() bool { return c.IP_C != nil && c.S_BUF != nil && c.R_BUF != nil }

// bind binds the connection and buffers of c to the client d.
func (c *client) bind(d *client) { d.IP_C, d.S_BUF, d.R_BUF = c.IP_C, c.S_BUF, c.R_BUF }

// dnsQuery writes into the send buffer a DNS query, with the transaction
// ID tid, of the type typ for the name d, and returns its size.
func dnsQuery(sb *network.NETWORK_BUFFER, tid iec.BYTE, d iec.STRING, typ iec.BYTE) iec.INT {
	l := LEN(d)
	for i := iec.INT(0); i <= 17+l; i++ {
		sb.BUFFER[i] = 0
	}
	sb.BUFFER[1] = tid  // transaction ID, low byte
	sb.BUFFER[2] = 0x01 // flags: recursion desired
	sb.BUFFER[5] = 0x01 // a question
	// The name, from its end: a length in place of each dot.
	i := l + 12
	count := iec.INT(0)
	for i > 12 {
		c := str.CODE(d, i-12)
		if c == 46 {
			sb.BUFFER[i] = iec.BYTE(count)
			count = 0
		} else {
			sb.BUFFER[i] = c
			count++
		}
		i--
	}
	sb.BUFFER[i] = iec.BYTE(count)
	i = l + 15
	sb.BUFFER[i] = typ // the type, low byte
	i += 2
	sb.BUFFER[i] = 0x01 // class IN, low byte
	return i + 1
}

// dnsConnect sets the connection up for a query to the DNS server ip.
func dnsConnect(c *network.IP_C, ip iec.DWORD) {
	c.C_PORT = 53
	c.C_IP = ip
	c.C_MODE = 1 // UDP client
	c.C_ENABLE = true
	c.TIME_RESET = true
	c.R_OBSERVE = true
}

// DNS_CLIENT finds the IPv4 address IP4 of DOMAIN with the DNS server
// IP4_DNS, on a rising edge of ACTIVATE: DONE when found, else ERROR, the
// DNS error, 255 for no address, or the error of the connection. A DOMAIN
// that is an address is its address.
type DNS_CLIENT struct {
	client
	ACTIVATE iec.BOOL
	DOMAIN   iec.STRING
	IP4_DNS  iec.DWORD
	IP4      iec.DWORD
	DONE     iec.BOOL
	ERROR    iec.DWORD

	i            iec.INT
	state        iec.INT
	domainCopy   iec.STRING
	activateLast iec.BOOL
	urlLength    iec.INT
	tid          iec.BYTE
	anc          iec.INT
	rrStop       iec.BOOL
	ttlLast      iec.DWORD
	ttlTmp       iec.DWORD
}

// INIT resets the block.
func (d *DNS_CLIENT) INIT() {
	*d = DNS_CLIENT{client: client{IP_C: d.IP_C, S_BUF: d.S_BUF, R_BUF: d.R_BUF}}
}

// Execute runs the block once.
func (d *DNS_CLIENT) Execute(now time.Time) {
	if !d.ok() {
		return
	}
	c, sb, rb := d.IP_C, d.S_BUF, d.R_BUF
	switch d.state {
	case 0: // wait for a rising edge
		if d.ACTIVATE && !d.activateLast {
			d.DONE = false
			d.ERROR = 0
			d.state = 5
		}
	case 5:
		// A domain that is an address.
		d.IP4 = encoding.IP4_DECODE(d.DOMAIN)
		if d.IP4 > 0 {
			d.DONE = true
			d.state = 0
		} else {
			d.ipState = 1 // ask for the connection
			d.ttlLast = 0
			d.domainCopy = d.DOMAIN
			d.state = 10
		}
	case 10:
		if d.ipState == 3 { // the connection is ours
			d.urlLength = LEN(d.domainCopy)
			d.tid++
			d.i = dnsQuery(sb, d.tid, d.domainCopy, 1) // type A
			dnsConnect(c, d.IP4_DNS)
			sb.SIZE = iec.UINT(d.i)
			rb.SIZE = 0
			d.state = 30
		}
	case 30:
		switch {
		case c.ERROR != 0:
			d.ERROR = c.ERROR
			d.state = 99
		case sb.SIZE == 0 && d.tid == rb.BUFFER[1] && rb.SIZE >= iec.UINT(34+d.urlLength):
			// an answer with the transaction ID, long enough
			d.ERROR = iec.DWORD(rb.BUFFER[3] & 0x0F) // the return code
			if d.ERROR == 0 {
				d.anc = iec.INT(rb.BUFFER[7]) // the answers
				c.R_OBSERVE = false
				d.state = 40
			} else {
				d.state = 99
			}
		}
	case 40:
		if d.anc > 0 {
			d.anc--
			for iec.UINT(d.i) < rb.SIZE {
				x := rb.BUFFER[d.i]
				if !d.rrStop {
					switch {
					case x > 63: // a pointer to a name
						d.i++
						d.rrStop = true
					case x == 0: // the end of a name
						d.rrStop = true
					default:
						d.i += iec.INT(x)
					}
					d.i++
					continue
				}
				d.rrStop = false
				if rb.BUFFER[d.i+1] == 1 && rb.BUFFER[d.i+3] == 1 { // type A, class IN
					d.i += 4
					d.ttlTmp = logic.DWORD_OF_BYTE(rb.BUFFER[d.i], rb.BUFFER[d.i+1], rb.BUFFER[d.i+2], rb.BUFFER[d.i+3])
					d.i += 6
					if d.ttlTmp >= d.ttlLast { // the longest time to live
						d.ttlLast = d.ttlTmp
						d.IP4 = logic.DWORD_OF_BYTE(rb.BUFFER[d.i], rb.BUFFER[d.i+1], rb.BUFFER[d.i+2], rb.BUFFER[d.i+3])
					}
					d.i += 4
				} else {
					d.i += 9
					d.i += iec.INT(rb.BUFFER[d.i]) + 1 // the next record
				}
				break
			}
		} else {
			if d.IP4 > 0 {
				d.DONE = true
			} else {
				d.ERROR = 255 // no address
			}
			d.state = 99
		}
	case 99:
		rb.SIZE = 0
		d.ipState = 4 // give the connection back
		d.state = 0
	}
	d.activateLast = d.ACTIVATE
	d.runFIFO(now, d.IP_C)
}

// DNS_REV_CLIENT finds the DOMAIN of the IPv4 address IP4 with the DNS
// server IP4_DNS, on a rising edge of ACTIVATE: DONE when found, else
// ERROR. OSCAT rotates the input IP4 as it reads it.
type DNS_REV_CLIENT struct {
	client
	ACTIVATE iec.BOOL
	IP4      iec.DWORD
	IP4_DNS  iec.DWORD
	DOMAIN   iec.STRING
	DONE     iec.BOOL
	ERROR    iec.DWORD

	state        iec.INT
	ip4Copy      iec.DWORD
	d            iec.STRING // STRING(27)
	activateLast iec.BOOL
	i            iec.INT
	dl           iec.INT
	p1           iec.INT
	tid          iec.BYTE
}

// INIT resets the block.
func (d *DNS_REV_CLIENT) INIT() {
	*d = DNS_REV_CLIENT{client: client{IP_C: d.IP_C, S_BUF: d.S_BUF, R_BUF: d.R_BUF}}
}

// Execute runs the block once.
func (d *DNS_REV_CLIENT) Execute(now time.Time) {
	if !d.ok() {
		return
	}
	c, sb, rb := d.IP_C, d.S_BUF, d.R_BUF
	switch d.state {
	case 0: // wait for a rising edge
		if d.ACTIVATE && !d.activateLast {
			d.ERROR = 0
			d.ip4Copy = d.IP4
			d.ipState = 1
			d.DOMAIN = ""
			d.DONE = false
			d.d = ""
			for i := 0; i < 4; i++ {
				d.d = CONCAT(d.d, network.BYTE_TO_STRING(iec.BYTE(d.IP4)), ".")
				d.IP4 = ROR(d.IP4, 8)
			}
			d.d = CONCAT(d.d, "in-addr.arpa")
			d.state = 10
		}
	case 10:
		if d.ipState == 3 {
			d.dl = LEN(d.d)
			d.tid++
			d.i = dnsQuery(sb, d.tid, d.d, 12) // type PTR
			dnsConnect(c, d.IP4_DNS)
			sb.SIZE = iec.UINT(d.i)
			rb.SIZE = 0
			d.state = 30
		}
	case 30:
		if c.ERROR != 0 {
			d.ERROR = c.ERROR
		} else if sb.SIZE == 0 && d.tid == rb.BUFFER[1] && rb.SIZE >= iec.UINT(d.i+12) {
			d.ERROR = iec.DWORD(rb.BUFFER[3] & 0x0F)
			if d.ERROR == 0 {
				// The name of the answer, the last in it.
				d.i = iec.INT(rb.SIZE) - 2
				for rb.BUFFER[d.i] > 0 {
					if d.i == 0 {
						break
					}
					d.i--
				}
				d.i += 2
				d.p1 = d.i + 1
				for iec.UINT(d.i) < rb.SIZE {
					d.dl = iec.INT(rb.BUFFER[d.i])
					if d.dl > 0 {
						rb.BUFFER[d.i] = 46 // a dot for each length
						d.i += d.dl + 1
					} else {
						d.DOMAIN = buffer.BUFFER_TO_STRING(rb.BUFFER[:], iec.UINT(len(rb.BUFFER)), iec.UINT(d.p1), iec.UINT(d.i-1))
						break
					}
				}
				d.state = 0
				d.DONE = true
			}
		}
		if d.ERROR != 0 {
			d.state = 0
		}
		if d.state == 0 {
			c.R_OBSERVE = false
			d.ipState = 4
		}
	}
	d.activateLast = d.ACTIVATE
	d.runFIFO(now, d.IP_C)
}

// DNS_DYN updates a dynamic DNS name, HOSTNAME at dyndns.org (MODE 0) or
// selfhost.de, with the account USERNAME and PASSWORD, to the address IP4,
// or the one the service sees if 0: on a rising edge of UPDATE, and every
// T_UPDATE, while ENABLE. DONE when done, else ERROR_C and ERROR_T: 1 DNS,
// 2 HTTP, 3 the service refused.
type DNS_DYN struct {
	client
	ENABLE   iec.BOOL
	UPDATE   iec.BOOL
	T_UPDATE iec.TIME // default T#1h
	MODE     iec.BYTE
	HOSTNAME iec.STRING // STRING(30)
	USERNAME iec.STRING // STRING(20)
	PASSWORD iec.STRING // STRING(20)
	IP4      iec.DWORD
	BUSY     iec.BOOL
	DONE     iec.BOOL
	ERROR_C  iec.DWORD
	ERROR_T  iec.BYTE

	dns         DNS_CLIENT
	http        HTTP_GET
	base64      encoding.BASE64_ENCODE_STR
	urlData     network.URL
	updateLast  iec.BOOL
	state       iec.INT
	base64Done  iec.BOOL
	s1          iec.STRING // STRING(144)
	s2          iec.STRING // STRING(192)
	w           timers.TON
	initialized bool
}

// INIT resets the block and sets T_UPDATE to its initial value.
func (d *DNS_DYN) INIT() {
	*d = DNS_DYN{client: client{IP_C: d.IP_C, S_BUF: d.S_BUF, R_BUF: d.R_BUF}, T_UPDATE: iec.TIME(time.Hour), initialized: true}
	d.http.INIT()
}

// Execute runs the block once.
func (d *DNS_DYN) Execute(now time.Time) {
	if !d.initialized {
		d.initialized = true
		d.http.INIT()
	}
	if !d.ok() {
		return
	}
	switch d.state {
	case 0:
		if d.ENABLE && (d.UPDATE && !d.updateLast || d.w.Q) {
			d.state = 20
			d.DONE = false
			d.BUSY = true
			d.ERROR_C, d.ERROR_T = 0, 0
			if d.MODE == 0 { // dyndns.org
				d.s2 = CONCAT("http://members.dyndns.org/nic/update?hostname=", d.HOSTNAME)
			} else { // selfhost.de
				d.s2 = "http://carol.selfhost.de/nic/update?hostname=1"
			}
			if d.IP4 > 0 {
				d.s1 = CONCAT("&myip=", encoding.IP4_TO_STRING(d.IP4))
				d.s2 = CONCAT(d.s2, d.s1)
			}
			d.urlData = encoding.STRING_TO_URL(d.s2, "", "")
			d.s1 = CONCAT(d.USERNAME, ":", d.PASSWORD)
		}
	case 20:
		if d.dns.DONE {
			d.state = 40
		} else if d.dns.ERROR > 0 {
			d.ERROR_C, d.ERROR_T = d.dns.ERROR, 1
			d.BUSY = false
			d.state = 0
		}
	case 40:
		if d.base64Done {
			d.s1 = CONCAT("Authorization: Basic ", d.s2, "\r\nUser-Agent: x\r\n")
			d.urlData.HEADER = network.STRING_N(d.s1, 160)
			d.state = 60
		}
	case 60:
		if d.http.DONE {
			d.state = 0
			d.BUSY = false
			d.DONE = true
			// good and nochg are the answers of a success.
			d.s2 = buffer.BUFFER_TO_STRING(d.R_BUF.BUFFER[:], iec.UINT(len(d.R_BUF.BUFFER)), d.http.BODY_START+3, d.http.BODY_STOP-7)
			if FIND(d.s2, "good") == 0 && FIND(d.s2, "nochg") == 0 {
				d.ERROR_C, d.ERROR_T = 1, 3
				d.DONE = false
			}
		} else if d.http.ERROR > 0 {
			d.ERROR_C, d.ERROR_T = d.http.ERROR, 2
			d.BUSY = false
			d.state = 0
		}
	}
	d.client.bind(&d.dns.client)
	d.dns.DOMAIN, d.dns.IP4_DNS, d.dns.ACTIVATE = d.urlData.DOMAIN, 0, d.state == 20
	d.dns.Execute(now)
	d.client.bind(&d.http.client)
	d.http.IP4, d.http.GET, d.http.MODE, d.http.UNLOCK_BUF, d.http.URL_DATA = d.dns.IP4, d.state == 60, 2, d.state != 60, &d.urlData
	d.http.Execute(now)
	d.base64.RUN, d.base64.STR1, d.base64.STR2 = d.state == 40, &d.s1, &d.s2
	d.base64.Execute(now)
	d.base64Done = d.base64.DONE
	if d.T_UPDATE > 0 {
		d.w.IN, d.w.PT = d.state == 0, d.T_UPDATE
		d.w.Execute(now)
	}
	d.updateLast = d.UPDATE
}
