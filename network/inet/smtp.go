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
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/crypto"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/file"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// SMTP_CLIENT sends an email, on a rising edge of ACTIVATE, through the
// server SERVER, an URL 'smtp://user:password@host:port' ('%' for an '@' in
// the user; ssl:// and tls:// ask for SSL, which TwinCAT's IP_CONTROL does
// not give), with SMTP, or ESMTP with the authentication the server takes
// (PLAIN, LOGIN or CRAM-MD5) when there is a password:
//   - MAILFROM: 'address;name';
//   - MAILTO: 'to#cc#bcc', each a ';' list of addresses;
//   - SUBJECT and BODY;
//   - FILES: a ';' list of files to attach, '#DEL#' among them to delete
//     them after;
//   - DTI the time, DTI_OFFSET minutes from UTC.
//
// DONE when sent, BUSY while it works, else ERROR_C and ERROR_T: 1 DNS, 2
// the connection, 4 a file, 5 a step that took longer than TIMEOUT (10 s at
// least), with the step and the last reply.
type SMTP_CLIENT struct {
	ACTIVATE   iec.BOOL
	TIMEOUT    iec.TIME
	DTI        iec.DT
	DTI_OFFSET iec.INT
	DNS_IP4    iec.DWORD
	DONE       iec.BOOL
	BUSY       iec.BOOL
	ERROR_C    iec.DWORD
	ERROR_T    iec.BYTE
	SERVER     *iec.STRING
	MAILFROM   *iec.STRING
	MAILTO     *iec.STRING // STRING(STRING_LENGTH)
	SUBJECT    *iec.STRING
	BODY       *iec.STRING // STRING(STRING_LENGTH)
	FILES      *iec.STRING // STRING(STRING_LENGTH)

	ipC           network.IP_C
	sBuf          network.NETWORK_BUFFER
	rBuf          network.NETWORK_BUFFER
	fBuf          network.NETWORK_BUFFER
	urlData       network.URL
	dns           DNS_CLIENT
	ipc           ip.IP_CONTROL
	fs            file.FILE_SERVER
	fsd           network.FILE_SERVER_DATA
	b64Str        encoding.BASE64_ENCODE_STR
	b64Stream     encoding.BASE64_ENCODE_STREAM
	fpd           network.FILE_PATH_DATA
	md5Cram       crypto.MD5_CRAM_AUTH
	md5Run        iec.BOOL
	activateLast  iec.BOOL
	step          iec.INT
	sndStep       iec.INT
	str1, str2    iec.STRING // STRING(STRING_LENGTH)
	str3          iec.STRING
	sndText       iec.STRING // STRING(STRING_LENGTH)
	sndLfCnt      iec.INT
	rcvText       iec.STRING // STRING(STRING_LENGTH)
	rcvState      iec.INT
	authState     iec.BYTE
	tonWait       timers.TON
	lastStep      iec.INT
	nextStep      iec.INT
	timeout2      iec.TIME
	timeout3      iec.TIME
	lastBlock     iec.BOOL
	fileDelete    iec.BOOL
	esmtp         iec.BOOL
	cnt1, cnt2    iec.INT
	cnt3, cnt4    iec.INT
	idx1, idx2    iec.INT
	sndIdx        iec.INT
	sndEnable     iec.BOOL
	b64Done       iec.BOOL
	b64Start      iec.BOOL
	b64Buf1       [48]iec.BYTE
	b64Buf2       [64]iec.BYTE
	b64Size1      iec.INT
	b64Size2      iec.INT
	b64Str1       iec.STRING // STRING(144)
	b64Str2       iec.STRING // STRING(192)
	md5User       iec.STRING // STRING(64)
	md5Password   iec.STRING // STRING(64)
	md5B64TS      iec.STRING // STRING(64)
	md5AuthKey    iec.STRING // STRING(192)
	ipCRedDisable iec.BOOL   // ignore a close by the remote
	sslMode       iec.BYTE   // 0 off, 1 SSL, 2 TLS
}

// The authentications SMTP_CLIENT knows, in the order of its steps 1000,
// 1100, 1200 and 1300.
const (
	authNames    iec.STRING = " PLAIN; LOGIN; CRAM-MD5"
	authNamesTLS iec.STRING = " PLAIN; LOGIN; CRAM-MD5;STARTTLS"
	b64Max       iec.UINT   = 2880 // the bytes of a file a block
)

// INIT resets the block.
func (s *SMTP_CLIENT) INIT() {
	*s = SMTP_CLIENT{SERVER: s.SERVER, MAILFROM: s.MAILFROM, MAILTO: s.MAILTO, SUBJECT: s.SUBJECT, BODY: s.BODY, FILES: s.FILES}
	s.bind()
}

// bind binds the connection and the file server.
func (s *SMTP_CLIENT) bind() {
	s.ipc.IP_C, s.ipc.S_BUF, s.ipc.R_BUF = &s.ipC, &s.sBuf, &s.rBuf
	s.dns.IP_C, s.dns.S_BUF, s.dns.R_BUF = &s.ipC, &s.sBuf, &s.rBuf
	s.fs.FSD, s.fs.PT = &s.fsd, &s.fBuf
}

// Execute runs the block once.
func (s *SMTP_CLIENT) Execute(now time.Time) {
	if s.SERVER == nil || s.MAILFROM == nil || s.MAILTO == nil || s.SUBJECT == nil || s.BODY == nil || s.FILES == nil {
		return
	}
	s.bind()
	c, sb, fsd := &s.ipC, &s.sBuf, &s.fsd
	idle := s.sndStep == 0
	switch s.step {
	case 0:
		if s.ACTIVATE && !s.activateLast {
			s.TIMEOUT = max(iec.TIME(10*time.Second), s.TIMEOUT)
			s.timeout2 = s.TIMEOUT + iec.TIME(time.Second)
			s.timeout3 = s.TIMEOUT / 2
			s.DONE = false
			s.ipCRedDisable = false
			s.ERROR_C, s.ERROR_T = 0, 0
			fsd.ERROR = 0
			s.fileDelete = false
			s.urlData = encoding.STRING_TO_URL(*s.SERVER, "", "/")
			// % in the user is an @.
			if s.idx1 = FIND(s.urlData.USER, "%"); s.idx1 > 0 {
				s.urlData.USER = REPLACE(s.urlData.USER, "@", 1, s.idx1)
			}
			s.esmtp = LEN(s.urlData.PASSWORD) > 0
			s.step = 10
		}
	case 10:
		if s.dns.DONE {
			s.step = 20
		} else if s.dns.ERROR != 0 {
			s.ERROR_C, s.ERROR_T = s.dns.ERROR, 1
			s.step = 980
		}
	case 20: // the connection
		switch {
		case FIND(";SSL;ssl;Ssl", s.urlData.PROTOCOL) > 1:
			s.sslMode = 1
			c.C_PORT = SEL(s.urlData.PORT == 0, s.urlData.PORT, 465)
			c.MAILBOX[15] = 1 // SSL
		case FIND(";TLS;tls;Tls", s.urlData.PROTOCOL) > 1:
			s.sslMode = 2
			c.C_PORT = SEL(s.urlData.PORT == 0, s.urlData.PORT, 587)
			c.MAILBOX[15] = 0
		default:
			s.sslMode = 0
			c.C_PORT = SEL(s.urlData.PORT == 0, s.urlData.PORT, 25)
			c.MAILBOX[15] = 0
		}
		c.C_IP = s.dns.IP4
		c.C_MODE = 0 // TCP client
		c.TIME_RESET = true
		c.C_ENABLE = true
		c.R_OBSERVE = true
		c.MAILBOX[2] = 0
		s.rBuf.SIZE = 0
		s.sndEnable = true
		s.sndIdx = 0
		s.step = 30
	case 30:
		if s.rcvState == 220 { // service ready
			if s.esmtp {
				s.sndText = "EHLO OSCAT"
				s.str1 = SEL(s.sslMode == 2, authNames, authNamesTLS)
				s.authState = 0
				s.idx2 = 1900
				s.nextStep = 40
			} else {
				s.sndText = "HELO OSCAT"
				s.nextStep = 90
			}
			s.sndLfCnt = 1
		}
	case 40:
		if s.rcvState == 250 {
			// A server may not name an authentication with SSL, or send
			// the 250 lines in parts: after 1 s, assume LOGIN.
			if s.sslMode >= 1 && s.tonWait.ET > iec.TIME(time.Second) {
				s.rcvText = "250-AUTH LOGIN"
			}
			if FIND(s.rcvText, "AUTH") > 0 {
				c.MAILBOX[2] = 1 // no receive until the reply is read
				s.idx1 = 0
				s.step = 50
			}
		}
	case 50: // the authentications of the server
		s.str2 = encoding.ELEMENT_GET(59, s.idx1, &s.str1)
		if LEN(s.str2) > 0 {
			if FIND(s.rcvText, s.str2) > 0 {
				s.authState = logic.BIT_LOAD_B(s.authState, true, s.idx1)
				s.str3 = s.str2
				s.idx2 = 1000 + s.idx1*100
			}
			s.idx1++
		} else if s.authState != 0 {
			c.MAILBOX[2] = 0
			s.step = s.idx2 // 1000 PLAIN, 1100 LOGIN, 1200 CRAM-MD5, 1300 TLS
		} else {
			s.idx1 = 0
			s.rcvText = ""
			c.MAILBOX[2] = 0
			s.step = 40 // the next 250 line
		}

	case 90:
		if s.rcvState == 235 || s.rcvState == 250 {
			s.str2 = *s.MAILFROM
			s.str1 = encoding.ELEMENT_GET(59, 0, &s.str2)
			s.sndText = INSERT("MAIL FROM: <>", s.str1, 12)
			s.str3 = REPLACE(s.sndText, "From", 9, 1) // From: <address>, for the header
			s.nextStep = 95
			s.cnt1 = encoding.ELEMENT_COUNT(35, s.MAILTO)
			s.cnt2 = 0
		}
	case 95:
		if s.rcvState == 250 {
			if s.cnt2 < s.cnt1 {
				s.str1 = encoding.ELEMENT_GET(35, s.cnt2, s.MAILTO)
				s.cnt3 = encoding.ELEMENT_COUNT(59, &s.str1)
				s.cnt4 = 0
				s.cnt2++
				s.step = 100
			} else {
				s.step = 110
			}
		}
	case 100:
		if s.rcvState == 250 {
			if s.cnt4 < s.cnt3 {
				s.str2 = encoding.ELEMENT_GET(59, s.cnt4, &s.str1)
				s.sndText = INSERT("RCPT TO: <>", s.str2, 10)
				s.cnt4++
				s.nextStep = 100
			} else {
				s.step = 95
			}
		}
	case 110:
		if s.rcvState == 250 {
			s.sndText = "DATA"
			s.nextStep = 120
		}
	case 120:
		if s.rcvState == 354 { // start the mail
			s.sndEnable = false
			s.str2 = *s.MAILFROM
			s.str1 = encoding.ELEMENT_GET(59, 1, &s.str2) // the name shown
			if LEN(s.str1) > 0 {
				s.str2 = INSERT(`"" `, s.str1, 1)
				s.sndText = INSERT(s.str3, s.str2, 6)
			} else {
				s.sndText = s.str3
			}
			s.cnt2 = 0
			s.cnt1 = min(2, s.cnt1) // not the blind copies
			s.step = 200
		}
	case 200:
		if idle {
			if s.cnt2 < s.cnt1 {
				s.str1 = encoding.ELEMENT_GET(35, s.cnt2, s.MAILTO)
				s.cnt3 = encoding.ELEMENT_COUNT(59, &s.str1)
				s.cnt4 = 0
				s.cnt2++
				s.str3 = SEL[iec.STRING](s.cnt2 == 1, "Cc: <>", "To: <>")
				s.step = 205
			} else {
				s.step = 210
			}
		}
	case 205:
		if idle {
			if s.cnt4 < s.cnt3 {
				s.str2 = encoding.ELEMENT_GET(59, s.cnt4, &s.str1)
				s.sndText = INSERT(s.str3, s.str2, 5)
				s.cnt4++
			} else {
				s.step = 200
			}
		}
	case 210:
		if idle {
			s.sndText = CONCAT("Subject: ", *s.SUBJECT)
			s.step = 220
		}
	case 220:
		if idle {
			s.DTI = DWORD_TO_DT(DT_TO_DWORD(s.DTI) + iec.DWORD(iec.UDINT(-s.DTI_OFFSET*60)))
			if s.esmtp {
				// Date: Thu, 21 May 1998 05:33:29 +0000
				s.str2 = ";Mon;Tue;Wed;Thu;Fri;Sat;Sun"
				s.str1 = encoding.ELEMENT_GET(59, td.DAY_OF_WEEK(DT_TO_DATE(s.DTI)), &s.str2)
				s.sndText = CONCAT("Date: ", s.str1, str.DT_TO_STRF(s.DTI, 0, ", #H #E #A #N:#R:#T +0000", 1))
			} else {
				// Date: 21 May 98 05:33:29
				s.sndText = str.DT_TO_STRF(s.DTI, 0, "Date: #G #E #B #N:#R:#T", 1)
			}
			s.step = 230
		}
	case 230:
		if idle {
			s.sndText = "MIME-Version: 1.0"
			s.step = 250
		}
	case 250:
		if idle {
			s.sndText = `Content-Type: multipart/mixed;boundary="x"`
			s.sndLfCnt = 2
			s.step = 260
		}
	case 260:
		if idle {
			s.sndEnable = true
			s.sndText = "This is a multi-part message in MIME format."
			s.step = 270
		}
	case 270:
		if idle {
			s.sndEnable = false
			s.sndText = "--x"
			s.sndLfCnt = 1
			s.step = 280
		}
	case 280:
		if idle {
			s.sndText = `Content-Type: text/plain; format=flowed; charset="iso-8859-1"; reply-type=original`
			s.step = 290
		}
	case 290:
		if idle {
			s.sndText = "Content-Transfer-Encoding: 8bit"
			s.sndLfCnt = 2
			s.step = 300
		}
	case 300:
		if idle {
			s.sndEnable = true
			s.sndText = *s.BODY
			s.sndLfCnt = 2
			s.step = 310
		}
	case 310:
		if idle {
			s.cnt1 = encoding.ELEMENT_COUNT(59, s.FILES) - 1
			s.step = 400
		}

	// The attachments.
	case 400:
		if s.cnt1 >= 0 {
			s.str1 = encoding.ELEMENT_GET(59, s.cnt1, s.FILES)
			encoding.FILE_PATH_SPLIT(s.str1, &s.fpd)
			if s.fpd.FILENAME == "#DEL#" {
				s.fileDelete = true
			} else {
				s.step = 410
			}
			s.cnt1--
		} else {
			s.step = 800
		}
	case 410:
		if idle {
			s.sndEnable = false
			s.sndText = "--x"
			s.sndLfCnt = 1
			s.step = 420
		}
	case 420:
		if idle {
			s.sndText = "Content-Transfer-Encoding: BASE64"
			s.str2 = CONCAT(s.fpd.FILENAME, `"`)
			s.step = 430
		}
	case 430:
		if idle {
			s.sndText = CONCAT(`Content-Type: application/octet-stream; name="`, s.str2)
			s.step = 440
		}
	case 440:
		if idle {
			s.sndEnable = true
			s.sndText = CONCAT(`Content-Disposition: attachment; filename="`, s.str2)
			s.sndLfCnt = 2
			s.step = 450
		}
	case 450:
		if idle {
			fsd.FILENAME = s.str1
			fsd.OFFSET = 0
			s.step = 460
		}
	case 460:
		fsd.MODE = 1 // read
		s.fBuf.SIZE = b64Max
		s.idx1, s.idx2 = 0, 0
		s.step = 470
	case 470:
		if fsd.MODE == 0 && fsd.ERROR == 0 {
			s.lastBlock = fsd.FILE_SIZE == fsd.OFFSET
			s.b64Size1 = min(48, iec.INT(s.fBuf.SIZE)-s.idx1)
			var lf iec.BOOL
			if s.b64Size1 > 0 {
				// A line of Base64.
				for i1 := iec.INT(0); i1 < s.b64Size1; i1++ {
					s.b64Buf1[i1] = s.fBuf.BUFFER[s.idx1]
					s.idx1++
				}
				s.b64Stream.SIZE1, s.b64Stream.BUF1, s.b64Stream.BUF2 = s.b64Size1, &s.b64Buf1, &s.b64Buf2
				s.b64Stream.Execute(now)
				s.b64Size2 = s.b64Stream.SIZE2
				for i1 := iec.INT(0); i1 < s.b64Size2; i1++ {
					sb.BUFFER[s.idx2] = s.b64Buf2[i1]
					s.idx2++
				}
				lf = true
			} else {
				lf = s.lastBlock
				s.step = 480 // the block is done
			}
			if lf {
				sb.BUFFER[s.idx2], sb.BUFFER[s.idx2+1] = 0x0D, 0x0A
				s.idx2 += 2
			}
		}
	case 480:
		sb.SIZE = iec.UINT(s.idx2)
		s.step = 490
	case 490:
		if sb.SIZE == 0 { // all sent
			s.step = SEL[iec.INT](s.lastBlock, 460, 400)
		}

	case 800:
		if idle {
			s.sndEnable = false
			s.sndText = "--x--"
			s.sndLfCnt = 2
			s.step = 810
		}
	case 810:
		if idle {
			s.sndEnable = true
			s.sndText = "."
			s.sndLfCnt = 1
			s.step = 820
		}
	case 820:
		if s.rcvState == 250 {
			s.sndText = "QUIT"
			s.ipCRedDisable = true
			s.nextStep = 830
		}
	case 830:
		if s.rcvState == 221 || s.tonWait.ET > s.timeout3 {
			s.step = 900
		}
	case 900:
		if s.fileDelete {
			s.cnt1 = encoding.ELEMENT_COUNT(59, s.FILES) - 1
			s.step = 910
		} else {
			s.step = 950
		}
	case 910: // delete the files
		if fsd.MODE == 0 && fsd.ERROR == 0 {
			if s.cnt1 >= 0 {
				s.str1 = encoding.ELEMENT_GET(59, s.cnt1, s.FILES)
				if s.str1 != "#DEL#" {
					fsd.FILENAME = s.str1
					fsd.MODE = 4
				}
				s.cnt1--
			} else {
				s.step = 950
			}
		}
	case 950:
		s.DONE = true
		s.step = 980
	case 980:
		c.C_ENABLE = false
		sb.SIZE, s.rBuf.SIZE = 0, 0
		c.MAILBOX[2] = 0
		s.sndStep = 0
		s.nextStep = 0
		fsd.MODE = 5 // close the file
		s.step = 990
	case 990:
		if c.C_STATE == 0 && !fsd.FILE_OPEN {
			s.step = 0
		}

	// The authentications.
	case 1000, 1100, 1200: // AUTH PLAIN, LOGIN or CRAM-MD5
		s.sndText = CONCAT("AUTH", s.str3)
		s.sndLfCnt = 1
		s.step += 10
	case 1010:
		if s.rcvState == 334 {
			// Base64 of 0 user 0 password, 46 characters at most.
			s.cnt1, s.cnt2 = LEN(s.urlData.USER), LEN(s.urlData.PASSWORD)
			if s.cnt1+s.cnt2 <= 46 {
				s.b64Buf1[0] = 0
				for i1 := iec.INT(1); i1 <= s.cnt1; i1++ {
					s.b64Buf1[i1] = str.CODE(s.urlData.USER, i1)
				}
				s.b64Size1 = s.cnt1 + 1
				s.b64Buf1[s.b64Size1] = 0
				s.b64Size1++
				for i1 := iec.INT(1); i1 <= s.cnt2; i1++ {
					s.b64Buf1[s.b64Size1] = str.CODE(s.urlData.PASSWORD, i1)
					s.b64Size1++
				}
				s.b64Stream.SIZE1, s.b64Stream.BUF1, s.b64Stream.BUF2 = s.b64Size1, &s.b64Buf1, &s.b64Buf2
				s.b64Stream.Execute(now)
				s.b64Size2 = s.b64Stream.SIZE2
				copy(sb.BUFFER[:s.b64Size2], s.b64Buf2[:s.b64Size2])
				sb.BUFFER[s.b64Size2], sb.BUFFER[s.b64Size2+1] = 0x0D, 0x0A
				sb.SIZE = iec.UINT(s.b64Size2 + 2)
				s.rcvState = 0
			}
			s.step = 90
		}
	case 1110:
		if s.rcvState == 334 { // Username:
			s.b64Str1 = s.urlData.USER
			s.b64Start = true
			s.step = 1120
		}
	case 1120:
		if s.b64Start && s.b64Done {
			s.b64Start = false
			s.sndText = s.b64Str2
			s.nextStep = 1130
		}
	case 1130:
		if s.rcvState == 334 { // Password:
			s.b64Str1 = s.urlData.PASSWORD
			s.b64Start = true
			s.step = 1140
		}
	case 1140:
		if s.b64Start && s.b64Done {
			s.b64Start = false
			s.sndText = s.b64Str2
			s.nextStep = 90
		}
	case 1210:
		if s.rcvState == 334 { // the challenge
			s.md5B64TS = network.STRING_N(RIGHT(s.rcvText, LEN(s.rcvText)-4), 64)
			s.md5User = network.STRING_N(s.urlData.USER, 64)
			s.md5Password = network.STRING_N(s.urlData.PASSWORD, 64)
			s.md5Run = true
			s.step = 1220
		}
	case 1220:
		s.md5Cram.RUN, s.md5Cram.USERNAME, s.md5Cram.PASSWORD = &s.md5Run, &s.md5User, &s.md5Password
		s.md5Cram.B64_TS, s.md5Cram.AUTH_KEY = &s.md5B64TS, &s.md5AuthKey
		s.md5Cram.Execute(now)
		if !s.md5Run {
			s.sndText = s.md5AuthKey
			s.step = 90
		}
	case 1300: // STARTTLS
		s.sndText = "STARTTLS"
		s.sndLfCnt = 1
		s.step = 1310
	case 1310:
		if s.rcvState == 220 {
			c.MAILBOX[15] = 1 // SSL
			s.step = 1320
		}
	case 1320:
		if c.MAILBOX[15] == 2 { // SSL is on
			s.step = 30
		}
	case 1900: // no authentication: wait for the time out
	}

	// The reply received: the code of its last line.
	if s.step >= 30 && s.rBuf.SIZE > 0 {
		rb := &s.rBuf
		idx := iec.INT(rb.SIZE)
		for {
			idx--
			if idx == 0 || rb.BUFFER[idx] >= 32 {
				break
			}
		}
		// The lines of a reply of several, joined by ^.
		for i1 := iec.INT(0); i1 <= idx; i1++ {
			if rb.BUFFER[i1] < 32 {
				rb.BUFFER[i1] = 94
			}
		}
		s.rcvText = buffer.BUFFER_TO_STRING(rb.BUFFER[:], iec.UINT(idx+1), 0, iec.UINT(STRING_LENGTH-1))
		if LEN(s.rcvText) >= 3 {
			if code := LEFT(s.rcvText, 3); str.IS_NUM(code) {
				s.rcvState = STRING_TO_INT(code)
			}
			c.R_OBSERVE = false
		}
		rb.SIZE = 0
	}

	// Send the lines; a line is sent with those after it until one
	// enables the send.
	switch s.sndStep {
	case 0:
		if LEN(s.sndText) > 0 {
			buffer.STRING_TO_BUFFER_(s.sndText, s.sndIdx, sb.BUFFER[:], iec.UINT(len(sb.BUFFER)))
			s.sndIdx += LEN(s.sndText)
			for i1 := iec.INT(1); i1 <= s.sndLfCnt; i1++ {
				sb.BUFFER[s.sndIdx], sb.BUFFER[s.sndIdx+1] = 0x0D, 0x0A
				s.sndIdx += 2
			}
			s.sndText = ""
			s.rcvState = 9999 // no reply yet
			if s.sndEnable {
				sb.SIZE = iec.UINT(s.sndIdx)
				s.rBuf.SIZE = 0
				c.R_OBSERVE = s.step <= 120 || s.step >= 820
				s.sndIdx = 0
				s.sndStep = 10
			}
		}
	case 10:
		if sb.SIZE == 0 { // all sent
			if s.nextStep != 0 {
				s.step = s.nextStep
				s.nextStep = 0
			}
			s.sndStep = 0
		}
	}

	// The errors.
	if s.ERROR_T == 0 && s.step > 20 {
		if c.ERROR > 0 && c.C_ENABLE && !c.TIME_RESET && (!s.ipCRedDisable || c.ERROR != 0xFD00_0000) {
			s.ERROR_C, s.ERROR_T = c.ERROR, 2
			s.step = 980
		}
		if fsd.MODE == 0 && fsd.ERROR > 0 {
			s.ERROR_C, s.ERROR_T = iec.DWORD(fsd.ERROR), 4
			s.step = 980
		}
	}
	if s.tonWait.Q {
		s.ERROR_C = SHL(iec.DWORD(s.step), 16) | iec.DWORD(s.rcvState)
		s.ERROR_T = 5
		s.step = 980
	}

	s.dns.DOMAIN, s.dns.IP4_DNS, s.dns.ACTIVATE = s.urlData.DOMAIN, s.DNS_IP4, s.step == 10
	s.dns.Execute(now)
	s.ipc.IP, s.ipc.PORT, s.ipc.TIME_OUT = 0, 0, s.TIMEOUT
	s.ipc.Execute(now)
	s.fs.Execute(now)
	s.b64Str.RUN, s.b64Str.STR1, s.b64Str.STR2 = s.b64Start, &s.b64Str1, &s.b64Str2
	s.b64Str.Execute(now)
	s.b64Done = s.b64Str.DONE
	s.tonWait.IN, s.tonWait.PT = s.step == s.lastStep && s.step > 0, s.timeout2
	s.tonWait.Execute(now)
	s.lastStep = s.step
	s.activateLast = s.ACTIVATE
	s.BUSY = s.step != 0
}
