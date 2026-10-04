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
	"github.com/apiarytech/beebread/network/file"
	"github.com/apiarytech/beebread/network/ip"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// control is a control connection of a client of its own: its IP_CONTROL,
// and the lines it sends and the reply codes it reads.
type control struct {
	IP_C  network.IP_C
	S_BUF network.NETWORK_BUFFER
	R_BUF network.NETWORK_BUFFER
	ipc   ip.IP_CONTROL

	sndStep  iec.INT
	sndText  iec.STRING // STRING(STRING_LENGTH)
	rcvText  iec.STRING // STRING(STRING_LENGTH)
	rcvState iec.INT
}

// receive reads the reply code of the line received, if any.
func (k *control) receive() {
	if k.R_BUF.SIZE == 0 {
		return
	}
	// The end of the text, before the end of the line.
	idx1 := iec.INT(k.R_BUF.SIZE)
	for {
		idx1--
		if idx1 == 0 || k.R_BUF.BUFFER[idx1] >= 32 {
			break
		}
	}
	idx1++
	k.rcvText = buffer.BUFFER_TO_STRING(k.R_BUF.BUFFER[:], iec.UINT(idx1), 0, iec.UINT(STRING_LENGTH-1))
	if LEN(k.rcvText) >= 3 {
		if code := LEFT(k.rcvText, 3); str.IS_NUM(code) {
			k.rcvState = STRING_TO_INT(code)
		}
		k.IP_C.R_OBSERVE = false
	}
	k.R_BUF.SIZE = 0
}

// send sends the line sndText, and then moves step to next.
func (k *control) send(step *iec.INT, next iec.INT) {
	switch k.sndStep {
	case 0:
		if LEN(k.sndText) > 0 {
			buffer.STRING_TO_BUFFER_(k.sndText, 0, k.S_BUF.BUFFER[:], iec.UINT(len(k.S_BUF.BUFFER)))
			i := LEN(k.sndText)
			k.S_BUF.BUFFER[i], k.S_BUF.BUFFER[i+1] = 0x0D, 0x0A
			k.S_BUF.SIZE = iec.UINT(i + 2)
			k.R_BUF.SIZE = 0
			k.IP_C.R_OBSERVE = true
			k.rcvState = 9999 // no reply yet
			k.sndText = ""
			k.sndStep = 10
		}
	case 10:
		if k.S_BUF.SIZE == 0 { // all sent
			*step = next
			k.sndStep = 0
		}
	}
}

// FTP_CLIENT sends the file FILENAME to the FTP URL FTP_URL, or with
// FTP_DOWNLOAD gets it from there, on a rising edge of ACTIVATE, in the
// passive mode, or the active mode with FTP_ACTIV and the PLC's address
// PLC_IP4, deleting the file after if FILE_DELETE. A URL without a file
// name takes FILENAME's. DNS_IP4 is the DNS server. DONE when done, BUSY
// while it works, else ERROR_C and ERROR_T: 1 DNS, 2 the control
// connection, 3 the data connection, 4 the file, 5 a step that took longer
// than TIMEOUT (10 s at least), with the step and the last reply.
type FTP_CLIENT struct {
	ACTIVATE     iec.BOOL
	FILENAME     iec.STRING
	FTP_URL      iec.STRING // STRING(STRING_LENGTH)
	FTP_DOWNLOAD iec.BOOL
	FTP_ACTIV    iec.BOOL
	FILE_DELETE  iec.BOOL
	TIMEOUT      iec.TIME
	DNS_IP4      iec.DWORD
	PLC_IP4      iec.DWORD
	DONE         iec.BOOL
	BUSY         iec.BOOL
	ERROR_C      iec.DWORD
	ERROR_T      iec.BYTE

	c1           control // the control connection
	ipC2         network.IP_C
	sBuf2        network.NETWORK_BUFFER
	rBuf2        network.NETWORK_BUFFER
	ipc2         ip.IP_CONTROL
	urlData      network.URL
	dns          DNS_CLIENT
	fs           file.FILE_SERVER
	fsd          network.FILE_SERVER_DATA
	activateLast iec.BOOL
	step         iec.INT
	tonWait      timers.TON
	lastStep     iec.INT
	nextStep     iec.INT
	str1, str2   iec.STRING // STRING(STRING_LENGTH)
	ftpPath      iec.STRING
	byar         [6]iec.BYTE
	dwTmp        iec.DWORD
	wTmp         iec.WORD
	fpd          network.FILE_PATH_DATA
	timeout2     iec.TIME
	timeout3     iec.TIME
	ftpFileSize  iec.UDINT
	c1RedDisable iec.BOOL // ignore a close by the remote of the control connection
	c2RedDisable iec.BOOL // ignore a close by the remote of the data connection
}

// INIT resets the block.
func (f *FTP_CLIENT) INIT() {
	*f = FTP_CLIENT{}
	f.bind()
}

// bind binds the connections and the file server.
func (f *FTP_CLIENT) bind() {
	f.c1.ipc.IP_C, f.c1.ipc.S_BUF, f.c1.ipc.R_BUF = &f.c1.IP_C, &f.c1.S_BUF, &f.c1.R_BUF
	f.ipc2.IP_C, f.ipc2.S_BUF, f.ipc2.R_BUF = &f.ipC2, &f.sBuf2, &f.rBuf2
	f.dns.IP_C, f.dns.S_BUF, f.dns.R_BUF = &f.c1.IP_C, &f.c1.S_BUF, &f.c1.R_BUF
	f.fs.FSD = &f.fsd
}

// Execute runs the block once.
func (f *FTP_CLIENT) Execute(now time.Time) {
	f.bind()
	c1, c2, fsd := &f.c1.IP_C, &f.ipC2, &f.fsd
	rcv := f.c1.rcvState
	switch f.step {
	case 0:
		if f.ACTIVATE && !f.activateLast {
			f.ftpFileSize = 0
			f.TIMEOUT = max(iec.TIME(10*time.Second), f.TIMEOUT)
			f.timeout2 = f.TIMEOUT + iec.TIME(time.Second)
			f.timeout3 = f.TIMEOUT / 2
			f.DONE = false
			f.c1RedDisable, f.c2RedDisable = false, false
			f.ERROR_C, f.ERROR_T = 0, 0
			fsd.ERROR = 0
			f.urlData = encoding.STRING_TO_URL(f.FTP_URL, "", "/")
			if LEN(f.urlData.USER) == 0 {
				f.urlData.USER, f.urlData.PASSWORD = "Anonymous", "User@"
			}
			// The file name on the server: the URL's, or the local file's.
			encoding.FILE_PATH_SPLIT(f.FILENAME, &f.fpd)
			f.str1 = f.fpd.FILENAME
			encoding.FILE_PATH_SPLIT(f.urlData.PATH, &f.fpd)
			if LEN(f.fpd.FILENAME) == 0 {
				f.ftpPath = CONCAT(f.urlData.PATH, f.str1)
			} else {
				f.ftpPath = f.urlData.PATH
			}
			f.step = 10
		}
	case 10:
		if f.dns.DONE {
			f.step = 20
		} else if f.dns.ERROR != 0 {
			f.ERROR_C, f.ERROR_T = f.dns.ERROR, 1
			f.step = 980
		}
	case 20: // the control connection
		c1.C_PORT = 21
		c1.C_IP = f.dns.IP4
		c1.C_MODE = 0 // TCP client
		c1.TIME_RESET = true
		c1.C_ENABLE = true
		c1.R_OBSERVE = true
		f.c1.R_BUF.SIZE = 0
		c2.C_ENABLE = false
		f.step = 30
	case 30:
		if rcv == 220 { // ready for a new user
			f.c1.sndText = CONCAT("USER ", f.urlData.USER)
			f.nextStep = 40
		}
	case 40:
		switch rcv {
		case 331: // the password
			f.c1.sndText = CONCAT("PASS ", f.urlData.PASSWORD)
			f.nextStep = 50
		case 230: // logged in
			f.nextStep = 50
		case 220:
			f.nextStep = 20
		}
	case 50:
		if rcv == 230 { // logged in
			f.c1.sndText = "TYPE I" // binary
			f.nextStep = 54
		}
	case 54:
		if rcv == 200 {
			if f.FTP_DOWNLOAD {
				f.c1.sndText = CONCAT("SIZE ", f.ftpPath)
				f.nextStep = 56
			} else {
				f.step = 60
			}
		}
	case 56:
		if rcv == 213 { // the size of the file
			f.str1 = encoding.ELEMENT_GET(32, 1, &f.c1.rcvText)
			f.ftpFileSize = network.STRING_TO_UDINT(f.str1)
			f.step = 60
		}
	case 60:
		if f.FTP_ACTIV {
			// PORT h1,h2,h3,h4,p1,p2: the server connects to the PLC.
			if f.urlData.PORT == 0 {
				f.urlData.PORT = 20
			}
			f.str2 = ""
			f.dwTmp = f.PLC_IP4
			for idx1 := 0; idx1 <= 5; idx1++ {
				if idx1 == 4 {
					f.dwTmp = ROR(iec.DWORD(f.urlData.PORT), 8)
				} else {
					f.dwTmp = ROL(f.dwTmp, 8)
				}
				f.str2 = CONCAT(f.str2, network.BYTE_TO_STRING(iec.BYTE(f.dwTmp)))
				if idx1 < 5 {
					f.str2 = CONCAT(f.str2, ",")
				}
			}
			f.dwTmp = f.dns.IP4
			f.wTmp = f.urlData.PORT
			f.c1.sndText = CONCAT("PORT ", f.str2)
			f.nextStep = 80
		} else {
			f.c1.sndText = "PASV"
			f.nextStep = 70
		}
		f.step = 999 // wait for nextStep
	case 70:
		if rcv == 227 { // Entered Passive Mode (h1,h2,h3,h4,p1,p2)
			idx1, idx2 := FIND(f.c1.rcvText, "("), FIND(f.c1.rcvText, ")")
			if idx1 > 0 && idx2 > idx1+1 {
				f.str1 = MID(f.c1.rcvText, idx2-idx1-1, idx1+1)
				for i := iec.INT(0); i <= 5; i++ {
					f.str2 = encoding.ELEMENT_GET(44, i, &f.str1)
					f.byar[i] = network.STRING_TO_BYTE(f.str2)
				}
				f.dwTmp = logic.DWORD_OF_BYTE(f.byar[0], f.byar[1], f.byar[2], f.byar[3])
				f.wTmp = logic.WORD_OF_BYTE(f.byar[4], f.byar[5])
			}
			f.step = 100
		}
	case 80:
		if rcv == 200 { // PORT command successful
			f.step = 100
		}
	case 100: // the data connection
		c2.C_MODE = SEL[iec.BYTE](f.FTP_ACTIV, 0, 2)
		c2.C_IP = f.dwTmp
		c2.C_PORT = f.wTmp
		c2.TIME_RESET = true
		c2.C_ENABLE = true
		c2.R_OBSERVE = false
		f.rBuf2.SIZE = 0
		f.step = 120
	case 120:
		if c2.C_STATE > 127 { // connected
			if f.FTP_DOWNLOAD {
				c2.MAILBOX[2] = 1 // no receive yet
				f.step = 400
			} else {
				f.c1.sndText = CONCAT("STOR ", f.ftpPath)
				f.nextStep = 140
				f.step = 999
			}
		}

	// Send the file.
	case 140:
		if rcv == 125 || rcv == 150 { // the data connection is open
			f.step = 160
		}
	case 160:
		fsd.FILENAME = f.FILENAME
		fsd.MODE = 1 // read
		fsd.OFFSET = 0
		f.sBuf2.SIZE = 65535
		c2.MAILBOX[1] = 1 // no send yet
		f.step = 200
	case 200:
		if fsd.MODE == 0 && fsd.ERROR == 0 {
			c2.MAILBOX[1] = 0 // send
			f.step = 210
		}
	case 210:
		if f.sBuf2.SIZE == 0 { // all sent
			if fsd.FILE_SIZE-fsd.OFFSET > 0 {
				fsd.MODE = 1
				f.sBuf2.SIZE = 65535
				c2.MAILBOX[1] = 1
				f.step = 200
			} else {
				c2.C_ENABLE = false
				f.step = 300
			}
		}
	case 300:
		if rcv == 226 { // the transfer is done
			f.c1.sndText = "QUIT"
			f.c1RedDisable = true
			f.nextStep = 320
		}
	case 320:
		if rcv == 221 || f.tonWait.ET > f.timeout3 {
			c1.C_ENABLE = false
			if f.FILE_DELETE {
				fsd.MODE = 4 // delete the file
				f.step = 340
			} else {
				f.step = 900
			}
		}
	case 340:
		if fsd.MODE == 0 && fsd.ERROR == 0 {
			f.step = 900
		}

	// Get the file.
	case 400:
		f.c1.sndText = CONCAT("RETR ", f.ftpPath)
		f.nextStep = 410
		f.step = 999
	case 410:
		if rcv == 125 || rcv == 150 {
			c2.R_OBSERVE = true
			f.step = 420
		}
	case 420:
		fsd.FILENAME = f.FILENAME
		fsd.MODE = 3 // a new file
		fsd.OFFSET = 0
		f.step = 460
	case 440:
		if f.rBuf2.SIZE > 0 {
			c2.MAILBOX[2] = 1 // no receive until the data are written
			fsd.MODE = 3
			f.step = 460
		} else if fsd.FILE_SIZE == f.ftpFileSize {
			c2.C_ENABLE = false
			f.c2RedDisable = true
			f.step = 700
		}
	case 460:
		if fsd.MODE == 0 && fsd.ERROR == 0 { // the data are written
			f.rBuf2.SIZE = 0
			c2.MAILBOX[2] = 0
			f.step = 440
		}

	case 700:
		if rcv == 226 || f.tonWait.ET > f.timeout3 {
			if f.FTP_DOWNLOAD && f.FILE_DELETE {
				f.c1.sndText = CONCAT("DELE ", f.ftpPath)
				f.nextStep = 720
			} else {
				f.step = 740
			}
		}
	case 720:
		if rcv == 250 { // deleted
			f.step = 740
		}
	case 740:
		f.c1.sndText = "QUIT"
		f.c1RedDisable = true
		f.nextStep = 760
		f.step = 999
	case 760:
		if rcv == 221 || f.tonWait.ET > f.timeout3 {
			c1.C_ENABLE = false
			if !f.FTP_DOWNLOAD && f.FILE_DELETE {
				fsd.MODE = 4
				f.step = 780
			} else {
				f.step = 900
			}
		}
	case 780:
		if fsd.MODE == 0 && fsd.ERROR == 0 {
			f.step = 900
		}

	case 900:
		f.DONE = true
		f.step = 980
	case 980:
		c1.C_ENABLE, c2.C_ENABLE = false, false
		c2.MAILBOX[1], c2.MAILBOX[2] = 0, 0
		f.c1.S_BUF.SIZE, f.sBuf2.SIZE, f.c1.R_BUF.SIZE, f.rBuf2.SIZE = 0, 0, 0, 0
		f.ftpFileSize = 0
		f.c1.sndStep = 0
		f.nextStep = 0
		fsd.MODE = 5 // close the file
		f.step = 990
	case 990:
		if c1.C_STATE == 0 && c2.C_STATE == 0 && !fsd.FILE_OPEN {
			f.step = 0
		}
	case 999: // wait for nextStep
	}

	if f.step >= 30 {
		f.c1.receive()
	}
	f.c1.send(&f.step, f.nextStep)

	// The errors.
	if f.ERROR_T == 0 && f.step > 20 {
		if c1.ERROR > 0 && c1.C_ENABLE && !c1.TIME_RESET && (!f.c1RedDisable || c1.ERROR != 0xFD00_0000) {
			f.ERROR_C, f.ERROR_T = c1.ERROR, 2
			f.step = 980
		}
		if c2.ERROR > 0 && c2.C_ENABLE && !c2.TIME_RESET && (!f.c2RedDisable || c2.ERROR != 0xFD00_0000) {
			f.ERROR_C, f.ERROR_T = c2.ERROR, 3
			f.step = 980
		}
		if fsd.MODE == 0 && fsd.ERROR > 0 {
			f.ERROR_C, f.ERROR_T = iec.DWORD(fsd.ERROR), 4
			f.step = 980
		}
	}
	if f.tonWait.Q {
		// A step that hangs: the step and the last reply.
		f.ERROR_C = SHL(iec.DWORD(f.step), 16) | iec.DWORD(f.c1.rcvState)
		f.ERROR_T = 5
		f.step = 980
	}

	f.dns.DOMAIN, f.dns.IP4_DNS, f.dns.ACTIVATE = f.urlData.DOMAIN, f.DNS_IP4, f.step == 10
	f.dns.Execute(now)
	f.c1.ipc.IP, f.c1.ipc.PORT, f.c1.ipc.TIME_OUT = 0, 0, f.TIMEOUT
	f.c1.ipc.Execute(now)
	f.ipc2.IP, f.ipc2.PORT, f.ipc2.TIME_OUT = 0, 0, f.TIMEOUT
	f.ipc2.Execute(now)
	f.tonWait.IN, f.tonWait.PT = f.step == f.lastStep && f.step > 0, f.timeout2
	f.tonWait.Execute(now)
	f.lastStep = f.step
	f.activateLast = f.ACTIVATE
	f.BUSY = f.step != 0
	if f.FTP_DOWNLOAD {
		f.fs.PT = &f.rBuf2
	} else {
		f.fs.PT = &f.sBuf2
	}
	f.fs.Execute(now)
}
