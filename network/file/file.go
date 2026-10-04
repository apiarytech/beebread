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

// Package file is the port of the OSCAT NETWORK file access: FILE_SERVER
// reads, writes and deletes files through TwinCAT's file service, emulated
// by the package ads, FILE_SERVER_RUNTIME measures its operations, and
// FILE_BLOCK reads a file a byte at a time.
package file

import (
	"encoding/binary"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/ads"
	"github.com/apiarytech/royaljelly/iec"
)

// FILE_SERVER_RUNTIME measures the time of FILE_SERVER's operations into
// FSD.RUNTIME: MODE 0 with a new COMMAND starts the measure, MODE 1 ends
// it.
type FILE_SERVER_RUNTIME struct {
	MODE    iec.BYTE
	COMMAND iec.BYTE
	FSD     *network.FILE_SERVER_DATA

	lastCommand iec.BYTE
	init        iec.BOOL
	startTx     iec.DWORD
	stopTx      iec.DWORD
	runtime     iec.DWORD
}

// INIT resets the block.
func (r *FILE_SERVER_RUNTIME) INIT() { *r = FILE_SERVER_RUNTIME{FSD: r.FSD} }

// measure records runtime in the current, highest and lowest time.
func measure(runtime iec.UDINT, cur, mx, mn *iec.UDINT) {
	*cur = runtime
	if *mx < runtime {
		*mx = runtime
	}
	if *mn > runtime {
		*mn = runtime
	}
}

// Execute runs the block once.
func (r *FILE_SERVER_RUNTIME) Execute(now time.Time) {
	if r.FSD == nil {
		return
	}
	rt := &r.FSD.RUNTIME
	if !r.init {
		r.init = true
		const none = 99999999
		rt.TIME_FILE_OPEN_MIN, rt.TIME_FILE_CLOSE_MIN, rt.TIME_FILE_READ_MIN = none, none, none
		rt.TIME_FILE_WRITE_MIN, rt.TIME_FILE_TELL_MIN, rt.TIME_FILE_SEEK_MIN = none, none, none
		rt.TIME_FILE_DELETE_MIN = none
	}
	if r.MODE == 0 && r.lastCommand == 0 && r.COMMAND != 0 {
		// a new command
		r.startTx = PLC_MS(now)
	}
	if r.MODE == 1 {
		r.stopTx = PLC_MS(now)
		r.runtime = r.stopTx - r.startTx
		t := iec.UDINT(r.runtime)
		switch r.COMMAND {
		case 1:
			measure(t, &rt.TIME_FILE_OPEN_CUR, &rt.TIME_FILE_OPEN_MAX, &rt.TIME_FILE_OPEN_MIN)
		case 2:
			measure(t, &rt.TIME_FILE_CLOSE_CUR, &rt.TIME_FILE_CLOSE_MAX, &rt.TIME_FILE_CLOSE_MIN)
		case 3:
			measure(t, &rt.TIME_FILE_READ_CUR, &rt.TIME_FILE_READ_MAX, &rt.TIME_FILE_READ_MIN)
		case 4:
			measure(t, &rt.TIME_FILE_WRITE_CUR, &rt.TIME_FILE_WRITE_MAX, &rt.TIME_FILE_WRITE_MIN)
		case 5:
			measure(t, &rt.TIME_FILE_SEEK_CUR, &rt.TIME_FILE_SEEK_MAX, &rt.TIME_FILE_SEEK_MIN)
		case 6:
			measure(t, &rt.TIME_FILE_TELL_CUR, &rt.TIME_FILE_TELL_MAX, &rt.TIME_FILE_TELL_MIN)
		case 7:
			measure(t, &rt.TIME_FILE_DELETE_CUR, &rt.TIME_FILE_DELETE_MAX, &rt.TIME_FILE_DELETE_MIN)
		}
	}
	r.lastCommand = r.COMMAND
}

// The commands of FILE_SERVER to the file service.
const (
	cmdOpen   iec.BYTE = 1
	cmdClose  iec.BYTE = 2
	cmdRead   iec.BYTE = 3
	cmdWrite  iec.BYTE = 4
	cmdSeek   iec.BYTE = 5
	cmdTell   iec.BYTE = 6
	cmdRemove iec.BYTE = 7
)

// errorBase is what FILE_SERVER adds to the ADS error of a command.
var errorBase = map[iec.BYTE]iec.UDINT{
	cmdOpen: 0x10, cmdClose: 0x30, cmdRead: 0x40, cmdWrite: 0x50,
	cmdSeek: 0x60, cmdTell: 0x70, cmdRemove: 0x80,
}

// FILE_SERVER does what FSD.MODE asks with the file FSD.FILENAME, and sets
// FSD.MODE to 0 when done, with FSD.ERROR: 1 reads PT.SIZE bytes into PT
// from the offset FSD.OFFSET, 2 writes the PT.SIZE bytes of PT there, 3 does
// the same on a new file, 4 deletes the file and 5 closes it. An offset of
// 16#FFFF_FFFE is the end of the file. It keeps the file open, with its
// size in FSD.FILE_SIZE, and moves FSD.OFFSET past what it read or wrote.
// FSD.ERROR is 255 for an offset past the end, or the ADS error of the
// operation plus errorBase of its command, in a byte.
type FILE_SERVER struct {
	FSD *network.FILE_SERVER_DATA
	PT  *network.NETWORK_BUFFER

	bufSize       iec.UDINT
	handle        iec.UINT
	readMaxLength iec.UDINT
	filePosition  iec.UDINT
	usedFilename  iec.STRING
	lengthRead    iec.UDINT
	writeLength   iec.UDINT
	lengthWritten [4]iec.BYTE
	openHandle    iec.UINT
	seekPosition  iec.UDINT
	seekMode      iec.INT
	seekData      [8]iec.BYTE // ARRAY[0..1] OF DINT
	tellPosition  [4]iec.BYTE
	openMode      iec.DWORD
	lastMode      iec.BYTE
	command       iec.BYTE
	errorCode     iec.BYTE
	error         iec.BOOL
	fileChange    iec.BOOL
	step          iec.INT
	x             ads.FW_AdsRdWrt
	tmpDW         [4]iec.BYTE
	para          iec.BOOL
	fsr           FILE_SERVER_RUNTIME
	filename      []iec.BYTE
}

// INIT resets the block.
func (f *FILE_SERVER) INIT() { *f = FILE_SERVER{FSD: f.FSD, PT: f.PT} }

// dword reads a DWORD from b, little endian.
func dword(b [4]iec.BYTE) iec.UDINT {
	return iec.UDINT(b[0]) | iec.UDINT(b[1])<<8 | iec.UDINT(b[2])<<16 | iec.UDINT(b[3])<<24
}

// Execute runs the block once.
func (f *FILE_SERVER) Execute(now time.Time) {
	if f.FSD == nil || f.PT == nil {
		return
	}
	fsd, pt := f.FSD, f.PT
	switch f.step {
	case 0:
		if fsd.MODE == LIMIT(1, fsd.MODE, 5) {
			f.x.SNETID, f.x.TTIMEOUT, f.x.NPORT = "", iec.TIME(5*time.Second), 10000
			f.bufSize = iec.UDINT(len(pt.BUFFER))
			switch {
			case fsd.MODE == 4: // remove
				if f.handle > 0 {
					f.fileChange = true
					f.step = 32000 // close first
				} else {
					f.step = 90
				}
			case fsd.MODE == 5: // close
				f.step = 32000
			case f.handle == 0: // modes 1, 2, 3 open the file first
				f.step = 100
			case f.usedFilename == fsd.FILENAME && fsd.MODE == f.lastMode:
				f.step = 200 // the file is open
			default:
				f.step = 32000
				f.fileChange = true
			}
		}

	case 90: // remove
		f.command = cmdRemove
		f.step = 95
	case 95:
		if f.command == 0 {
			f.step = 30000
		}

	case 100: // open
		if fsd.MODE == 3 {
			f.openMode = 0x0001_001A // a new file: PATH_GENERIC, WRITE, PLUS, BINARY
		} else {
			f.openMode = 0x0001_0019 // PATH_GENERIC, READ, PLUS, BINARY
		}
		f.command = cmdOpen
		f.step = 105
	case 105:
		if f.command == 0 {
			if f.error {
				f.step = 30000
			} else {
				f.handle = f.openHandle
				f.usedFilename = fsd.FILENAME
				f.filePosition = 0xFFFF_FFFF // unknown
				f.lastMode = fsd.MODE
				f.step = 110
			}
		}
	case 110: // seek to the end
		f.command = cmdSeek
		f.seekPosition = 0
		f.seekMode = 2 // SEEK_END
		f.step = 115
	case 115:
		if f.command == 0 {
			if f.error {
				f.step = 30000
			} else {
				f.command = cmdTell
				f.step = 120
			}
		}
	case 120: // the size of the file
		if f.command == 0 {
			if f.error {
				f.step = 30000
			} else {
				fsd.FILE_SIZE = dword(f.tellPosition)
				fsd.FILE_OPEN = true
				f.step = 200
			}
		}

	case 200: // seek to the offset
		switch {
		case fsd.OFFSET == 0xFFFF_FFFE: // the end of the file
			f.command = cmdSeek
			f.seekPosition = fsd.FILE_SIZE
			f.seekMode = 0
			f.step = 210
		case fsd.OFFSET > fsd.FILE_SIZE:
			f.errorCode = 255
			f.step = 30000
		case fsd.OFFSET != f.filePosition:
			f.command = cmdSeek
			f.seekPosition = fsd.OFFSET
			f.seekMode = 0
			f.step = 210
		default:
			f.step = 300
		}
	case 210:
		if f.command == 0 {
			f.filePosition = f.seekPosition
			if f.error {
				f.step = 30000
			} else {
				f.step = 300
			}
		}

	case 300: // read or write
		f.step = SEL[iec.INT](fsd.MODE == 1, 500, 400)

	case 400: // read
		f.readMaxLength = min(iec.UDINT(pt.SIZE), f.bufSize)
		if f.filePosition+f.readMaxLength > fsd.FILE_SIZE {
			f.readMaxLength = fsd.FILE_SIZE - f.filePosition
		}
		if f.readMaxLength > 0 {
			f.command = cmdRead
			f.step = 410
		} else {
			pt.SIZE = 0
			f.step = 30000
		}
	case 410:
		if f.command == 0 {
			pt.SIZE = iec.UINT(f.lengthRead)
			if !f.error {
				f.filePosition += f.lengthRead
				fsd.OFFSET = f.filePosition
			}
			f.step = 30000
		}

	case 500: // write
		f.writeLength = min(iec.UDINT(pt.SIZE), f.bufSize)
		if f.writeLength > 0 {
			f.command = cmdWrite
			f.step = 510
		} else {
			f.step = 30000
		}
	case 510:
		if f.command == 0 {
			if !f.error {
				f.filePosition += dword(f.lengthWritten)
				fsd.OFFSET = f.filePosition
				if f.filePosition > fsd.FILE_SIZE {
					fsd.FILE_SIZE = f.filePosition
				}
			}
			f.step = 30000
		}

	case 30000: // done
		fsd.MODE = 0
		fsd.ERROR = f.errorCode
		f.step = 0

	case 32000: // close
		if f.handle > 0 {
			f.command = cmdClose
		}
		f.step = 32100
	case 32100:
		if f.command == 0 {
			fsd.FILE_OPEN = false
			f.filePosition = 0
			fsd.FILE_SIZE, fsd.OFFSET = 0, 0
			f.usedFilename = ""
			f.handle = 0
			if f.fileChange {
				// the file name changed: go on with it
				f.fileChange = false
				f.step = 0
			} else {
				f.step = 30000
			}
		}
	}

	f.x.BEXECUTE = false
	f.x.Execute(now)
	f.fsr.MODE, f.fsr.COMMAND, f.fsr.FSD = 0, f.command, fsd
	f.fsr.Execute(now)

	if f.command == 0 {
		return
	}
	if !f.para {
		// Start the call of the command.
		f.para = true
		x := &f.x
		x.NIDXOFFS = iec.UDINT(f.handle)
		x.CBWRITELEN, x.WRITEBUFF, x.CBREADLEN, x.READBUFF = 0, nil, 0, nil
		switch f.command {
		case cmdOpen:
			f.filename = network.BYTES(fsd.FILENAME, 81)
			x.NIDXGRP, x.NIDXOFFS = ads.SYSTEMSERVICE_FOPEN, iec.UDINT(f.openMode)
			x.CBWRITELEN, x.WRITEBUFF = iec.UDINT(len(f.filename)), f.filename
			x.CBREADLEN, x.READBUFF = 4, f.tmpDW[:]
		case cmdClose:
			x.NIDXGRP = ads.SYSTEMSERVICE_FCLOSE
		case cmdRead:
			x.NIDXGRP = ads.SYSTEMSERVICE_FREAD
			x.CBREADLEN, x.READBUFF = f.readMaxLength, pt.BUFFER[:]
		case cmdWrite:
			x.NIDXGRP = ads.SYSTEMSERVICE_FWRITE
			x.CBWRITELEN, x.WRITEBUFF = f.writeLength, pt.BUFFER[:]
			x.CBREADLEN, x.READBUFF = 4, f.lengthWritten[:]
		case cmdSeek:
			var b [8]byte
			binary.LittleEndian.PutUint32(b[0:4], uint32(int32(f.seekPosition)))
			binary.LittleEndian.PutUint32(b[4:8], uint32(int32(f.seekMode)))
			for i := range b {
				f.seekData[i] = iec.BYTE(b[i])
			}
			x.NIDXGRP = ads.SYSTEMSERVICE_FSEEK
			x.CBWRITELEN, x.WRITEBUFF = 8, f.seekData[:]
		case cmdTell:
			x.NIDXGRP = ads.SYSTEMSERVICE_FTELL
			x.CBREADLEN, x.READBUFF = 4, f.tellPosition[:]
		case cmdRemove:
			f.filename = network.BYTES(fsd.FILENAME, 81)
			x.NIDXGRP, x.NIDXOFFS = ads.SYSTEMSERVICE_FDELETE, 0x0001_0000 // PATH_GENERIC
			x.CBWRITELEN, x.WRITEBUFF = iec.UDINT(len(f.filename)), f.filename
		}
		x.BEXECUTE = true
		x.Execute(now)
		return
	}
	if f.x.BBUSY {
		return
	}
	// The call is done.
	f.para = false
	f.error = f.x.BERROR
	f.errorCode = SEL(f.error, 0, iec.BYTE(f.x.NERRID+errorBase[f.command]))
	switch f.command {
	case cmdOpen:
		if f.x.CBREAD >= 4 {
			f.openHandle = iec.UINT(dword(f.tmpDW))
		} else {
			f.openHandle = 0
		}
	case cmdRead:
		f.lengthRead = f.x.CBREAD
	}
	f.command = 0
	// OSCAT clears the command before it reports it, so
	// FILE_SERVER_RUNTIME measures nothing.
	f.fsr.MODE, f.fsr.COMMAND = 1, f.command
	f.fsr.Execute(now)
}

// FILE_BLOCK reads the byte at the position POS of the file FILENAME into
// DATA, with the FILE_SERVER of FSD and PT, which it reads a buffer at a
// time: MODE set to 1 starts it, and it is 0 when done, with ERROR. ERROR
// is 255 past the end of the file.
type FILE_BLOCK struct {
	MODE     *iec.BYTE
	FILENAME *iec.STRING
	FSD      *network.FILE_SERVER_DATA
	PT       *network.NETWORK_BUFFER
	POS      iec.UDINT
	ERROR    iec.BYTE
	DATA     iec.BYTE

	step      iec.INT
	dataStart iec.UDINT
	dataStop  iec.UDINT
}

// INIT resets the block.
func (b *FILE_BLOCK) INIT() {
	*b = FILE_BLOCK{MODE: b.MODE, FILENAME: b.FILENAME, FSD: b.FSD, PT: b.PT}
}

// Execute runs the block once.
func (b *FILE_BLOCK) Execute(now time.Time) {
	if b.MODE == nil || b.FILENAME == nil || b.FSD == nil || b.PT == nil {
		return
	}
	fsd, pt := b.FSD, b.PT
	switch b.step {
	case 0:
		if *b.MODE == 0 {
			return
		}
		b.ERROR = 0
		switch {
		case bool(fsd.FILE_OPEN && (fsd.FILE_SIZE == 0 || b.POS >= fsd.FILE_SIZE)):
			b.ERROR = 255
			*b.MODE = 0
		case bool(b.POS < b.dataStart || b.POS > b.dataStop || !fsd.FILE_OPEN):
			// read the buffer
			fsd.FILENAME = *b.FILENAME
			fsd.MODE = 1
			fsd.OFFSET = b.POS
			pt.SIZE = 65535
			b.dataStart, b.dataStop = 0, 0
			b.step = 10
		default:
			// the byte is in the buffer
			b.DATA = pt.BUFFER[b.POS-b.dataStart]
			*b.MODE = 0
		}
	case 10:
		if fsd.MODE == 0 {
			if fsd.ERROR > 0 {
				b.ERROR = fsd.ERROR
				*b.MODE = 0
			} else {
				b.dataStart = b.POS
				b.dataStop = b.POS + iec.UDINT(pt.SIZE) - 1
			}
			b.step = 0
		}
	}
}
