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
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/parser"
	"github.com/apiarytech/royaljelly/iec"
)

// READ_HTTP reads the value of the header SRC from the HTTP header in the
// first SIZE bytes of PT, from POS: the rest of its line after 'SRC: ',
// the name in any case. SRC ” reads the status line after 'HTTP/'.
type READ_HTTP struct {
	SIZE      iec.UINT
	POS       iec.INT
	SRC       iec.STRING
	PT        *network.NW_BUF_LONG
	READ_HTTP iec.STRING // STRING(STRING_LENGTH)
}

// INIT resets the block.
func (r *READ_HTTP) INIT() { *r = READ_HTTP{PT: r.PT} }

// Execute runs the block once.
func (r *READ_HTTP) Execute(now time.Time) {
	if r.PT == nil {
		return
	}
	pt := r.PT
	stop := iec.INT(r.SIZE) - 1
	src := r.SRC
	if LEN(src) == 0 {
		src = "HTTP/"
	} else {
		src = CONCAT(src, ": ")
	}
	s1 := buffer.BUFFER_SEARCH(pt[:], iec.INT(r.SIZE), src, r.POS, true)
	if s1 < 0 {
		r.READ_HTTP = ""
		return
	}
	e1 := s1
	// The start of the line.
	for index := s1 - 1; index > 0; index-- {
		if pt[index] < 32 {
			break
		}
		s1--
	}
	// The end of the line.
	for index := e1 + 1; index <= stop; index++ {
		if pt[index] < 32 {
			break
		}
		e1++
	}
	s1 += LEN(src)
	r.READ_HTTP = buffer.BUFFER_TO_STRING(pt[:], r.SIZE, iec.UINT(s1), iec.UINT(e1))
}

// HTTP_GET gets the URL URL_DATA from the server IP4, on a rising edge of
// GET, with HTTP/1.1 (MODE 2, or 3 to close the connection after) or
// HTTP/1.0 (MODE 0, or 1 to keep it): DONE with the answer in R_BUF, its
// header from HEADER_START to HEADER_STOP and its body from BODY_START to
// BODY_STOP, and HTTP_STATUS, until UNLOCK_BUF gives the connection back.
// ERROR is 16#FC for a status other than 200, or the connection's.
type HTTP_GET struct {
	client
	URL_DATA     *network.URL
	IP4          iec.DWORD
	GET          iec.BOOL
	MODE         iec.BYTE // default 2
	UNLOCK_BUF   iec.BOOL
	HTTP_STATUS  iec.STRING
	HEADER_START iec.UINT
	HEADER_STOP  iec.UINT
	BODY_START   iec.UINT
	BODY_STOP    iec.UINT
	DONE         iec.BOOL
	ERROR        iec.DWORD

	readHTTP       READ_HTTP
	base64         encoding.BASE64_ENCODE_STR
	b64Str1        iec.STRING // STRING(144)
	b64Str2        iec.STRING // STRING(192)
	b64Done        iec.BOOL
	b64Start       iec.BOOL
	authentication iec.BOOL
	getLast        iec.BOOL
	ip4Stored      iec.DWORD
	state          iec.INT
	http11Host     iec.BOOL
	totalSize      iec.INT
	text           iec.STRING // STRING(STRING_LENGTH)
	delimiter2b    iec.STRING
	delimiter4b    iec.STRING
	delEnd         iec.STRING
	delPos         iec.INT
	idx            iec.INT
	rcvTimeout     iec.BOOL
	conRdce        iec.BOOL
}

// INIT resets the block and sets MODE to its initial value.
func (h *HTTP_GET) INIT() {
	*h = HTTP_GET{client: client{IP_C: h.IP_C, S_BUF: h.S_BUF, R_BUF: h.R_BUF}, URL_DATA: h.URL_DATA, MODE: 2}
}

// put copies s into the send buffer at idx, and moves idx past it.
func (h *HTTP_GET) put(s iec.STRING) {
	buffer.STRING_TO_BUFFER_(s, h.idx, h.S_BUF.BUFFER[:], iec.UINT(len(h.S_BUF.BUFFER)))
	h.idx += LEN(s)
}

// Execute runs the block once.
func (h *HTTP_GET) Execute(now time.Time) {
	if !h.ok() || h.URL_DATA == nil {
		return
	}
	c, sb, rb, u := h.IP_C, h.S_BUF, h.R_BUF, h.URL_DATA
	switch h.state {
	case 0:
		h.DONE = false
		if h.GET && !h.getLast {
			h.ip4Stored = h.IP4
			h.totalSize = 0
			h.HEADER_START, h.HEADER_STOP, h.BODY_START, h.BODY_STOP = 0, 0, 0, 0
			h.HTTP_STATUS = ""
			h.ERROR = 0
			h.ipState = 1
			h.authentication = false
			h.delimiter2b, h.delimiter4b = "\r\n", "\r\n\r\n"
			if u.USER != "" && u.PASSWORD != "" {
				// Basic authentication: Base64 of USER:PASSWORD.
				h.b64Str1 = CONCAT(u.USER, ":", u.PASSWORD)
				h.b64Start = true
				h.authentication = true
				h.state = 5
			} else {
				h.state = 10
			}
		}
	case 5:
		if h.b64Start && h.b64Done {
			h.b64Start = false
			h.state = 10
		}
	case 10:
		if h.ipState == 3 {
			// The request.
			h.idx = 0
			h.put(CONCAT("GET ", u.PATH))
			if LEN(u.QUERY) > 0 {
				sb.BUFFER[h.idx] = 63 // ?
				h.idx++
				h.put(u.QUERY)
			}
			h.text = SEL[iec.STRING](h.MODE >= 2, " HTTP/1.0", " HTTP/1.1")
			h.put(CONCAT(h.text, h.delimiter2b))
			switch h.MODE {
			case 1:
				h.text = CONCAT("Connection: Keep-Alive", h.delimiter2b)
			case 3:
				h.text = CONCAT("Connection: Close", h.delimiter2b)
			default:
				h.text = ""
			}
			h.put(CONCAT(h.text, u.HEADER))
			if h.authentication {
				h.put(CONCAT("Authorization: Basic ", h.b64Str2, h.delimiter2b))
			}
			h.put(CONCAT("Host: ", u.DOMAIN, h.delimiter4b))
			// The connection.
			if FIND(";HTTPS;https;Https;", CONCAT(u.PROTOCOL, ";")) > 1 {
				c.C_PORT = SEL(u.PORT == 0, u.PORT, 443)
				c.MAILBOX[15] = 1 // SSL
			} else {
				c.C_PORT = SEL(u.PORT == 0, u.PORT, 80)
			}
			c.C_IP = h.ip4Stored
			c.C_MODE = 0 // TCP client
			c.C_ENABLE = true
			c.TIME_RESET = true
			c.R_OBSERVE = true
			sb.SIZE = iec.UINT(h.idx)
			rb.SIZE = 0
			h.state = 30
		}
	case 30:
		h.rcvTimeout = c.ERROR&0x0000_FF00 == 0x0000_FF00
		h.conRdce = c.ERROR&0xFD00_0000 == 0xFD00_0000
		if sb.SIZE == 0 {
			if h.totalSize == 0 {
				if rb.SIZE > 15 {
					// The end of the header: CR LF CR LF, or LF LF.
					h.delEnd = "\r\n\r\n"
					h.delPos = buffer.BUFFER_SEARCH(rb.BUFFER[:], iec.INT(rb.SIZE), h.delEnd, 0, false)
					if h.delPos == 0 {
						h.delEnd = "\n\n"
						h.delPos = buffer.BUFFER_SEARCH(rb.BUFFER[:], iec.INT(rb.SIZE), h.delEnd, 0, false)
					}
					if h.delPos > 0 {
						h.HEADER_START = 0
						h.HEADER_STOP = iec.UINT(h.delPos + LEN(h.delEnd) - 1)
						h.totalSize = -1 // the header, but not the size yet
					}
					// The status.
					h.readHTTP.SIZE, h.readHTTP.POS, h.readHTTP.SRC, h.readHTTP.PT = h.HEADER_STOP, 0, "", &rb.BUFFER
					h.readHTTP.Execute(now)
					if i := LEN(h.readHTTP.READ_HTTP); i > 7 {
						h.text = MID(h.readHTTP.READ_HTTP, 1, 3)
						h.http11Host = h.text == "1"
						h.HTTP_STATUS = RIGHT(h.readHTTP.READ_HTTP, i-4)
					}
					if FIND(h.HTTP_STATUS, "200") > 0 {
						h.readHTTP.SIZE, h.readHTTP.POS, h.readHTTP.SRC = h.HEADER_STOP, 0, "CONTENT-LENGTH"
						h.readHTTP.Execute(now)
						h.text = h.readHTTP.READ_HTTP
						if LEN(h.text) >= 1 && str.IS_NUM(h.text) {
							h.totalSize = STRING_TO_INT(h.text) + iec.INT(h.HEADER_STOP) + 1
						}
					} else {
						h.ERROR = 0xFC // the status
					}
				}
			} else if rb.SIZE >= iec.UINT(h.totalSize) ||
				rb.SIZE > 0 && h.totalSize < 0 && (h.rcvTimeout || h.conRdce || c.C_STATE == 1) {
				// All the data: the size known, or a time out or a
				// close by the server.
				if rb.SIZE-1 > h.HEADER_STOP {
					h.BODY_START = h.HEADER_STOP + 1
					h.BODY_STOP = rb.SIZE - 1
				}
				h.DONE = true
				c.C_ENABLE = false
				h.state = 40
			}
		}
		if c.ERROR != 0 {
			// The time out or close that ends data of no size is no error.
			switch {
			case bool(h.rcvTimeout && h.DONE):
				c.ERROR &= 0xFFFF_00FF
			case bool(h.conRdce && h.DONE):
				c.ERROR &= 0x00FF_FFFF
			default:
				h.ERROR = c.ERROR
			}
		}
		if h.ERROR > 0 {
			c.C_ENABLE = false
			h.state = 40
		}
	case 40:
		if h.UNLOCK_BUF || !h.DONE {
			h.ipState = 4
			h.DONE = false
			c.MAILBOX[15] = 0
			h.state = 0
		}
	}
	h.getLast = h.GET
	h.runFIFO(now, h.IP_C)
	h.base64.RUN, h.base64.STR1, h.base64.STR2 = h.b64Start, &h.b64Str1, &h.b64Str2
	h.base64.Execute(now)
	h.b64Done = h.base64.DONE
}

// httpService is what the services on HTTP_GET share: a DNS_CLIENT and an
// HTTP_GET of the URL urlData.
type httpService struct {
	client
	urlData     network.URL
	dns         DNS_CLIENT
	http        HTTP_GET
	state       iec.INT
	initialized bool
}

// start sets the HTTP_GET up.
func (s *httpService) start() {
	if !s.initialized {
		s.initialized = true
		s.http.INIT()
	}
}

// run runs the DNS_CLIENT and HTTP_GET: the query while the state is
// dns, the get while it is get, and the release at done.
func (s *httpService) run(now time.Time, dns, get, done iec.INT) {
	s.client.bind(&s.dns.client)
	s.dns.DOMAIN, s.dns.IP4_DNS, s.dns.ACTIVATE = s.urlData.DOMAIN, 0, s.state == dns
	s.dns.Execute(now)
	s.client.bind(&s.http.client)
	s.http.IP4, s.http.GET, s.http.MODE, s.http.UNLOCK_BUF, s.http.URL_DATA = s.dns.IP4, s.state == get, 2, s.state == done, &s.urlData
	s.http.Execute(now)
}

// GET_WAN_IP finds the internet address WAN_IP4 of the network with
// checkip.dyndns.com, on a rising edge of ACTIVATE: DONE, and NEW_IP4 if
// it changed, else ERROR_C and ERROR_T: 1 DNS, 2 HTTP, 3 no address.
type GET_WAN_IP struct {
	httpService
	ACTIVATE iec.BOOL
	WAN_IP4  iec.DWORD
	DONE     iec.BOOL
	NEW_IP4  iec.BOOL
	ERROR_C  iec.DWORD
	ERROR_T  iec.BYTE

	activateLast iec.BOOL
	wanIP4Last   iec.DWORD
	stIP         iec.STRING // STRING(120)
}

// INIT resets the block.
func (g *GET_WAN_IP) INIT() {
	*g = GET_WAN_IP{httpService: httpService{client: client{IP_C: g.IP_C, S_BUF: g.S_BUF, R_BUF: g.R_BUF}}}
	g.start()
}

// Execute runs the block once.
func (g *GET_WAN_IP) Execute(now time.Time) {
	g.start()
	if !g.ok() {
		return
	}
	rb := g.R_BUF
	switch g.state {
	case 0:
		if g.ACTIVATE && !g.activateLast {
			g.state = 20
			g.DONE = false
			g.ERROR_C, g.ERROR_T = 0, 0
		}
	case 20:
		g.urlData = encoding.STRING_TO_URL("checkip.dyndns.com", "", "")
		g.state = 40
	case 40:
		if g.dns.DONE {
			g.state = 60
		} else if g.dns.ERROR > 0 {
			g.ERROR_C, g.ERROR_T = g.dns.ERROR, 1
			g.state = 100
		}
	case 60:
		if g.http.DONE {
			g.state = 80
		} else if g.http.ERROR > 0 {
			g.ERROR_C, g.ERROR_T = g.http.ERROR, 2
			g.state = 100
		}
	case 80:
		// The address after the ':' of 'Current IP Address: 1.2.3.4</body>'.
		g.stIP = ""
		p1 := iec.INT(0)
		for i1 := iec.INT(g.http.BODY_START); i1 <= iec.INT(g.http.BODY_STOP); i1++ {
			if rb.BUFFER[i1] == 58 {
				p1 = i1 + 2
				break
			}
		}
		p3 := p1 + 15
		if p1 > 0 && p3 < iec.INT(g.http.BODY_STOP) {
			for p2 := p1; p2 <= p3; p2++ {
				if rb.BUFFER[p2] == 60 {
					break
				}
				g.stIP = CONCAT(g.stIP, str.CHR_TO_STRING(rb.BUFFER[p2]))
			}
		}
		g.WAN_IP4 = encoding.IP4_DECODE(g.stIP)
		if g.WAN_IP4 != 0 {
			g.DONE = true
			g.NEW_IP4 = g.WAN_IP4 != g.wanIP4Last
			g.wanIP4Last = g.WAN_IP4
		} else {
			g.ERROR_C, g.ERROR_T = 1, 3
		}
		g.state = 100
	case 100:
		if !g.http.DONE {
			g.state = 0
			g.DONE = g.ERROR_T == 0
		}
	}
	g.run(now, 40, 60, 100)
	g.activateLast = g.ACTIVATE
}

// IP2GEO finds the location GEO of the IPv4 address IP, or of the network
// if 0, with ipinfodb.com, on a rising edge of ACTIVATE: DONE, else
// ERROR_C and ERROR_T: 1 DNS, 2 HTTP.
type IP2GEO struct {
	httpService
	GEO      *network.IP2GEO_DATA
	IP       iec.DWORD
	ACTIVATE iec.BOOL
	BUSY     iec.BOOL
	DONE     iec.BOOL
	ERROR_C  iec.DWORD
	ERROR_T  iec.BYTE

	ctrl      network.XML_CONTROL
	xml       parser.XML_READER
	lastState iec.BOOL
	valueInt  iec.INT
	valueReal iec.REAL
}

// INIT resets the block.
func (g *IP2GEO) INIT() {
	*g = IP2GEO{httpService: httpService{client: client{IP_C: g.IP_C, S_BUF: g.S_BUF, R_BUF: g.R_BUF}}, GEO: g.GEO}
	g.start()
}

// Execute runs the block once.
func (g *IP2GEO) Execute(now time.Time) {
	g.start()
	if !g.ok() || g.GEO == nil {
		return
	}
	switch g.state {
	case 0:
		if g.ACTIVATE && !g.lastState {
			g.state = 20
			g.DONE = false
			g.BUSY = true
			g.ERROR_C, g.ERROR_T = 0, 0
		}
	case 20:
		g.urlData = encoding.STRING_TO_URL("http://ipinfodb.com/ip_query.php?timezone=true&IP=", "", "")
		if g.IP > 0 {
			g.urlData.QUERY = CONCAT(g.urlData.QUERY, encoding.IP4_TO_STRING(g.IP))
		}
		g.state = 40
	case 40:
		if g.dns.DONE {
			g.state = 60
		} else if g.dns.ERROR > 0 {
			g.ERROR_C, g.ERROR_T = g.dns.ERROR, 1
			g.state = 100
		}
	case 60:
		if g.http.DONE {
			g.state = 80
			g.ctrl.START_POS, g.ctrl.STOP_POS = g.http.BODY_START, g.http.BODY_STOP
			g.ctrl.COMMAND = 0x8018 // the text and the attributes
			g.ctrl.WATCHDOG = iec.TIME(time.Millisecond)
		} else if g.http.ERROR > 0 {
			g.ERROR_C, g.ERROR_T = g.http.ERROR, 2
			g.state = 100
		}
	case 80:
		// The XML, an item a scan.
		g.xml.CTRL, g.xml.BUF = &g.ctrl, &g.R_BUF.BUFFER
		g.xml.Execute(now)
		if g.ctrl.TYP < 98 {
			g.valueInt, g.valueReal = 0, 0
			if LEN(g.ctrl.VALUE) <= 20 {
				if v := str.FLOAT_TO_REAL(g.ctrl.VALUE); logic.CHK_REAL(v) == 0 {
					g.valueReal = v
					g.valueInt = REAL_TO_INT(v)
				}
			}
			geo := g.GEO
			switch g.ctrl.COUNT {
			case 7:
				geo.IP4 = encoding.IP4_DECODE(g.ctrl.VALUE)
			case 10:
				geo.STATE = g.ctrl.VALUE == "OK"
			case 13:
				geo.COUNTRY_CODE = network.STRING_N(g.ctrl.VALUE, 2)
			case 16:
				geo.COUNTRY_NAME = network.STRING_N(g.ctrl.VALUE, 20)
			case 19:
				geo.REGION_CODE = network.STRING_N(g.ctrl.VALUE, 2)
			case 22:
				geo.REGION_NAME = network.STRING_N(g.ctrl.VALUE, 20)
			case 25:
				geo.CITY = network.STRING_N(g.ctrl.VALUE, 20)
			case 31:
				geo.GEO_LATITUDE = g.valueReal
			case 34:
				geo.GEO_LONGITUDE = g.valueReal
			case 37:
				geo.TIME_ZONE_NAME = network.STRING_N(g.ctrl.VALUE, 20)
			case 40:
				geo.GMT_OFFSET = g.valueInt
			case 43:
				geo.IS_DST = g.valueInt > 0
			}
		} else if g.ctrl.TYP == 99 {
			g.DONE = true
			g.state = 100
		}
	case 100:
		if !g.http.DONE {
			g.state = 0
			g.BUSY = false
			g.DONE = g.ERROR_T == 0
		}
	}
	g.run(now, 40, 60, 100)
	g.lastState = g.ACTIVATE
}

// SPIDER_ACCESS reads (MODE 1) or writes (MODE 2) the variable VAR_NAME
// of a Phoenix Contact web server through its CGI, whose address IP_C
// has: a name without a dot is a global variable. ERROR is 1 for a write
// the server did not take, or the error of HTTP_GET.
type SPIDER_ACCESS struct {
	client
	VALUE    *iec.STRING
	VAR_NAME *iec.STRING // STRING(40)
	MODE     iec.BYTE
	ERROR    iec.DWORD

	state       iec.INT
	stTmp       iec.STRING // STRING(STRING_LENGTH)
	urlData     network.URL
	http        HTTP_GET
	valueLen    iec.UINT
	bodyLen     iec.UINT
	modeSave    iec.BYTE
	initialized bool
}

// INIT resets the block.
func (s *SPIDER_ACCESS) INIT() {
	*s = SPIDER_ACCESS{client: client{IP_C: s.IP_C, S_BUF: s.S_BUF, R_BUF: s.R_BUF}, VALUE: s.VALUE, VAR_NAME: s.VAR_NAME, initialized: true}
	s.http.INIT()
}

// Execute runs the block once.
func (s *SPIDER_ACCESS) Execute(now time.Time) {
	if !s.initialized {
		s.initialized = true
		s.http.INIT()
	}
	if !s.ok() || s.VALUE == nil || s.VAR_NAME == nil {
		return
	}
	switch s.state {
	case 0:
		if s.MODE > 0 && s.MODE < 3 {
			s.modeSave = s.MODE
			if FIND(*s.VAR_NAME, ".") == 0 {
				s.stTmp = CONCAT("%40GV.", *s.VAR_NAME) // a global variable
			} else {
				s.stTmp = *s.VAR_NAME
			}
			if s.modeSave == 2 { // write
				s.valueLen = iec.UINT(LEN(*s.VALUE))
				s.stTmp = CONCAT(s.stTmp, "+", *s.VALUE)
				s.urlData = encoding.STRING_TO_URL(CONCAT("http://x/dir/cgi-bin/writeVal.exe?", s.stTmp), "", "")
			} else {
				s.urlData = encoding.STRING_TO_URL(CONCAT("http://x/dir/cgi-bin/readVal.exe?", s.stTmp), "", "")
			}
			s.state = 60
		}
	case 60:
		if s.http.DONE {
			if s.modeSave == 2 {
				s.bodyLen = s.http.BODY_STOP - s.http.BODY_START + SEL[iec.UINT](s.http.BODY_STOP > 0, 0, 1)
				s.ERROR = SEL[iec.DWORD](s.bodyLen != s.valueLen, 0, 1)
			} else if s.http.BODY_START > 0 {
				*s.VALUE = buffer.BUFFER_TO_STRING(s.R_BUF.BUFFER[:], s.R_BUF.SIZE, s.http.BODY_START, s.http.BODY_STOP)
			} else {
				*s.VALUE = ""
			}
			s.state = 100
		} else if s.http.ERROR > 0 {
			s.ERROR = s.http.ERROR
			s.state = 100
		}
	case 100:
		if !s.http.DONE {
			s.state = 0
		}
	}
	s.client.bind(&s.http.client)
	s.http.IP4, s.http.GET, s.http.MODE, s.http.UNLOCK_BUF, s.http.URL_DATA = 0, s.state == 60, 0, s.state == 100, &s.urlData
	s.http.Execute(now)
}
