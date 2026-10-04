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

// Package ads emulates the ADS call FW_AdsRdWrt of Beckhoff TwinCAT for the
// file functions of its system service (port 10000), which OSCAT NETWORK's
// FILE_SERVER uses, on Go's os package: open, close, read, write, seek,
// tell and delete.
//
// As on TwinCAT, a rising edge of BEXECUTE starts the call; here it is done
// at once, and BBUSY is false after. A file's name is the bytes of the
// write buffer, up to a 0, in ISO 8859-1, a path of the operating system.
package ads

import (
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"os"
	"sync"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// The index groups of the file functions of the system service.
const (
	SYSTEMSERVICE_FOPEN   iec.UDINT = 120
	SYSTEMSERVICE_FCLOSE  iec.UDINT = 121
	SYSTEMSERVICE_FREAD   iec.UDINT = 122
	SYSTEMSERVICE_FWRITE  iec.UDINT = 123
	SYSTEMSERVICE_FSEEK   iec.UDINT = 124
	SYSTEMSERVICE_FTELL   iec.UDINT = 125
	SYSTEMSERVICE_FDELETE iec.UDINT = 131
)

// The modes of SYSTEMSERVICE_FOPEN, in the index offset.
const (
	FOPEN_MODEREAD   iec.UDINT = 0x01
	FOPEN_MODEWRITE  iec.UDINT = 0x02
	FOPEN_MODEAPPEND iec.UDINT = 0x04
	FOPEN_MODEPLUS   iec.UDINT = 0x08
	FOPEN_MODEBINARY iec.UDINT = 0x10
	FOPEN_MODETEXT   iec.UDINT = 0x20
)

// The ADS errors the calls report.
const (
	ADSERR_DEVICE_ERROR         iec.UDINT = 0x700
	ADSERR_DEVICE_SRVNOTSUPP    iec.UDINT = 0x701
	ADSERR_DEVICE_INVALIDGRP    iec.UDINT = 0x702
	ADSERR_DEVICE_INVALIDACCESS iec.UDINT = 0x704
	ADSERR_DEVICE_INVALIDSIZE   iec.UDINT = 0x705
	ADSERR_DEVICE_INVALIDPARM   iec.UDINT = 0x706
	ADSERR_DEVICE_NOTFOUND      iec.UDINT = 0x70C
)

var (
	filesMu sync.Mutex
	files   = map[iec.UDINT]*os.File{}
	next    iec.UDINT
)

// FW_AdsRdWrt writes CBWRITELEN bytes of WRITEBUFF to the index group
// NIDXGRP and offset NIDXOFFS of the ADS device NPORT, and reads up to
// CBREADLEN bytes into READBUFF: CBREAD of them.
type FW_AdsRdWrt struct {
	SNETID     iec.STRING
	NPORT      iec.UINT
	NIDXGRP    iec.UDINT
	NIDXOFFS   iec.UDINT
	CBWRITELEN iec.UDINT
	WRITEBUFF  []iec.BYTE
	CBREADLEN  iec.UDINT
	READBUFF   []iec.BYTE
	BEXECUTE   iec.BOOL
	TTIMEOUT   iec.TIME
	BBUSY      iec.BOOL
	BERROR     iec.BOOL
	NERRID     iec.UDINT
	CBREAD     iec.UDINT

	last iec.BOOL
}

// INIT resets the block.
func (f *FW_AdsRdWrt) INIT() { *f = FW_AdsRdWrt{} }

// Execute runs the block once.
func (f *FW_AdsRdWrt) Execute(now time.Time) {
	edge := f.BEXECUTE && !f.last
	f.last = f.BEXECUTE
	if !edge {
		return
	}
	f.CBREAD = 0
	f.NERRID = f.call()
	f.BERROR = f.NERRID != 0
	f.BBUSY = false
}

// write returns the bytes written.
func (f *FW_AdsRdWrt) write() []byte {
	n := min(int(f.CBWRITELEN), len(f.WRITEBUFF))
	out := make([]byte, n)
	for i := range out {
		out[i] = byte(f.WRITEBUFF[i])
	}
	return out
}

// name returns the file name written, up to a 0.
func (f *FW_AdsRdWrt) name() string {
	n := min(int(f.CBWRITELEN), len(f.WRITEBUFF))
	return string(STR(f.WRITEBUFF[:n]))
}

// read puts b in the read buffer.
func (f *FW_AdsRdWrt) read(b []byte) {
	n := min(len(b), int(f.CBREADLEN), len(f.READBUFF))
	for i := 0; i < n; i++ {
		f.READBUFF[i] = iec.BYTE(b[i])
	}
	f.CBREAD = iec.UDINT(n)
}

// readDword puts d in the read buffer, little endian.
func (f *FW_AdsRdWrt) readDword(d uint32) {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], d)
	f.read(b[:])
}

// adsErr returns the ADS error of err.
func adsErr(err error) iec.UDINT {
	switch {
	case err == nil:
		return 0
	case errors.Is(err, fs.ErrNotExist):
		return ADSERR_DEVICE_NOTFOUND
	case errors.Is(err, fs.ErrPermission):
		return ADSERR_DEVICE_INVALIDACCESS
	}
	return ADSERR_DEVICE_ERROR
}

func file(h iec.UDINT) *os.File {
	filesMu.Lock()
	defer filesMu.Unlock()
	return files[h]
}

// call makes the call and returns its error.
func (f *FW_AdsRdWrt) call() iec.UDINT {
	switch f.NIDXGRP {
	case SYSTEMSERVICE_FOPEN:
		mode := f.NIDXOFFS & 0xFFFF
		var flag int
		switch {
		case mode&FOPEN_MODEREAD != 0:
			flag = SEL(iec.BOOL(mode&FOPEN_MODEPLUS != 0), os.O_RDONLY, os.O_RDWR)
		case mode&FOPEN_MODEWRITE != 0:
			flag = SEL(iec.BOOL(mode&FOPEN_MODEPLUS != 0), os.O_WRONLY, os.O_RDWR) | os.O_CREATE | os.O_TRUNC
		case mode&FOPEN_MODEAPPEND != 0:
			flag = SEL(iec.BOOL(mode&FOPEN_MODEPLUS != 0), os.O_WRONLY, os.O_RDWR) | os.O_CREATE | os.O_APPEND
		default:
			return ADSERR_DEVICE_INVALIDPARM
		}
		fl, err := os.OpenFile(f.name(), flag, 0o644)
		if err != nil {
			return adsErr(err)
		}
		filesMu.Lock()
		next++
		h := next
		files[h] = fl
		filesMu.Unlock()
		f.readDword(uint32(h))
	case SYSTEMSERVICE_FCLOSE:
		filesMu.Lock()
		fl := files[f.NIDXOFFS]
		delete(files, f.NIDXOFFS)
		filesMu.Unlock()
		if fl == nil {
			return ADSERR_DEVICE_NOTFOUND
		}
		return adsErr(fl.Close())
	case SYSTEMSERVICE_FREAD:
		fl := file(f.NIDXOFFS)
		if fl == nil {
			return ADSERR_DEVICE_NOTFOUND
		}
		b := make([]byte, min(int(f.CBREADLEN), len(f.READBUFF)))
		n, err := io.ReadFull(fl, b)
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			err = nil
		}
		f.read(b[:n])
		return adsErr(err)
	case SYSTEMSERVICE_FWRITE:
		fl := file(f.NIDXOFFS)
		if fl == nil {
			return ADSERR_DEVICE_NOTFOUND
		}
		n, err := fl.Write(f.write())
		f.readDword(uint32(n))
		return adsErr(err)
	case SYSTEMSERVICE_FSEEK:
		fl := file(f.NIDXOFFS)
		if fl == nil {
			return ADSERR_DEVICE_NOTFOUND
		}
		w := f.write()
		if len(w) < 8 {
			return ADSERR_DEVICE_INVALIDSIZE
		}
		pos := int32(binary.LittleEndian.Uint32(w[0:4]))
		whence := int32(binary.LittleEndian.Uint32(w[4:8]))
		_, err := fl.Seek(int64(pos), int(whence))
		return adsErr(err)
	case SYSTEMSERVICE_FTELL:
		fl := file(f.NIDXOFFS)
		if fl == nil {
			return ADSERR_DEVICE_NOTFOUND
		}
		pos, err := fl.Seek(0, io.SeekCurrent)
		if err == nil {
			f.readDword(uint32(pos))
		}
		return adsErr(err)
	case SYSTEMSERVICE_FDELETE:
		return adsErr(os.Remove(f.name()))
	default:
		return ADSERR_DEVICE_SRVNOTSUPP
	}
	return 0
}
