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

package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// do runs FILE_SERVER until it is done with MODE, and returns FSD.ERROR.
func do(t *testing.T, f *FILE_SERVER, mode iec.BYTE) iec.BYTE {
	t.Helper()
	f.FSD.MODE = mode
	for i := 0; i < 100; i++ {
		f.Execute(time.Unix(1000, int64(i)*1e6))
		if f.FSD.MODE == 0 {
			return f.FSD.ERROR
		}
	}
	t.Fatalf("mode %d not done", mode)
	return 0
}

func setBuf(b *network.NETWORK_BUFFER, s string) {
	for i := range s {
		b.BUFFER[i] = iec.BYTE(s[i])
	}
	b.SIZE = iec.UINT(len(s))
}

func bufString(b *network.NETWORK_BUFFER) string {
	out := make([]byte, b.SIZE)
	for i := range out {
		out[i] = byte(b.BUFFER[i])
	}
	return string(out)
}

func TestFileServer(t *testing.T) {
	name := filepath.Join(t.TempDir(), "log.csv")
	var fsd network.FILE_SERVER_DATA
	var pt network.NETWORK_BUFFER
	f := FILE_SERVER{FSD: &fsd, PT: &pt}
	f.INIT()
	fsd.FILENAME = iec.STRING(name)

	// A new file.
	setBuf(&pt, "hello ")
	if e := do(t, &f, 3); e != 0 {
		t.Fatalf("write a new file: error %02X", e)
	}
	if fsd.FILE_SIZE != 6 || fsd.OFFSET != 6 || !fsd.FILE_OPEN {
		t.Fatalf("after writing: %+v", fsd)
	}
	// Append at the end.
	setBuf(&pt, "world")
	fsd.OFFSET = 0xFFFF_FFFE
	if e := do(t, &f, 3); e != 0 || fsd.FILE_SIZE != 11 {
		t.Fatalf("append: error %02X, size %d", e, fsd.FILE_SIZE)
	}
	if e := do(t, &f, 5); e != 0 || fsd.FILE_OPEN {
		t.Fatalf("close: error %02X", e)
	}
	if got, _ := os.ReadFile(name); string(got) != "hello world" {
		t.Fatalf("the file: %q", got)
	}

	// Read from an offset.
	fsd.OFFSET = 6
	pt.SIZE = 100
	if e := do(t, &f, 1); e != 0 || bufString(&pt) != "world" || fsd.OFFSET != 11 {
		t.Fatalf("read: error %02X, %q, offset %d", e, bufString(&pt), fsd.OFFSET)
	}
	// Past the end.
	fsd.OFFSET = 20
	if e := do(t, &f, 1); e != 255 {
		t.Fatalf("read past the end: error %d", e)
	}

	// Delete it.
	if e := do(t, &f, 4); e != 0 {
		t.Fatalf("delete: error %02X", e)
	}
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Fatal("the file is still there")
	}
	// A file that is not there.
	if e := do(t, &f, 1); e != 0x1C {
		t.Fatalf("open a missing file: error %02X, want 1C (ADS 70C + 10)", e)
	}
	// OSCAT reports command 0 to FILE_SERVER_RUNTIME, which measures nothing.
	if fsd.RUNTIME.TIME_FILE_OPEN_MIN != 99999999 || fsd.RUNTIME.TIME_FILE_OPEN_CUR != 0 {
		t.Errorf("FILE_SERVER_RUNTIME: %+v", fsd.RUNTIME)
	}
}

func TestFileBlock(t *testing.T) {
	name := filepath.Join(t.TempDir(), "data.txt")
	if err := os.WriteFile(name, []byte("ABCDEF"), 0o644); err != nil {
		t.Fatal(err)
	}
	var fsd network.FILE_SERVER_DATA
	var pt network.NETWORK_BUFFER
	fs := FILE_SERVER{FSD: &fsd, PT: &pt}
	fs.INIT()
	var mode iec.BYTE
	fn := iec.STRING(name)
	b := FILE_BLOCK{MODE: &mode, FILENAME: &fn, FSD: &fsd, PT: &pt}
	b.INIT()
	read := func(pos iec.UDINT) (iec.BYTE, iec.BYTE) {
		b.POS, mode = pos, 1
		for i := 0; i < 100 && mode != 0; i++ {
			b.Execute(time.Unix(1000, 0))
			fs.Execute(time.Unix(1000, 0))
		}
		return b.DATA, b.ERROR
	}
	for i, want := range "ABCDEF" {
		if got, e := read(iec.UDINT(i)); e != 0 || got != iec.BYTE(want) {
			t.Errorf("byte %d: %c, error %d", i, got, e)
		}
	}
	if _, e := read(6); e != 255 {
		t.Errorf("past the end: error %d", e)
	}
	// FILE_SERVER keeps the file open: close it.
	fsd.MODE = 5
	for i := 0; i < 100 && fsd.MODE != 0; i++ {
		fs.Execute(time.Unix(1000, 0))
	}
}
