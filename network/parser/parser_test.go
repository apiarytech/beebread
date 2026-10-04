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

package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/file"
	"github.com/apiarytech/royaljelly/iec"
)

func load(b *network.NETWORK_BUFFER, s string) {
	for i := range s {
		b.BUFFER[i] = iec.BYTE(s[i])
	}
	b.SIZE = iec.UINT(len(s))
}

type csvResult struct {
	value  iec.STRING
	result iec.BYTE
}

const csvData = "a;bb;c\r\n1;22;333\r\n"

var csvWant = []csvResult{{"a", 1}, {"bb", 1}, {"c", 2}, {"1", 1}, {"22", 1}, {"333", 2}, {"", 10}}

func TestCSVBuffer(t *testing.T) {
	var pt network.NETWORK_BUFFER
	load(&pt, csvData)
	var run iec.BYTE
	var offset iec.UDINT
	var value iec.STRING
	p := CSV_PARSER_BUF{RUN: &run, OFFSET: &offset, VALUE: &value, PT: &pt}
	p.INIT()
	p.SEP = ';'
	for _, want := range csvWant {
		run = 1
		for i := 0; i < 5 && run != 0; i++ {
			p.Execute(time.Unix(1000, 0))
		}
		if value != want.value || p.RESULT != want.result {
			t.Errorf("got %q %d, want %q %d", value, p.RESULT, want.value, want.result)
		}
	}
}

func TestCSVFile(t *testing.T) {
	name := filepath.Join(t.TempDir(), "data.csv")
	if err := os.WriteFile(name, []byte(csvData), 0o644); err != nil {
		t.Fatal(err)
	}
	var fsd network.FILE_SERVER_DATA
	var pt network.NETWORK_BUFFER
	fs := file.FILE_SERVER{FSD: &fsd, PT: &pt}
	fs.INIT()
	fn := iec.STRING(name)
	var run iec.BYTE
	var offset iec.UDINT
	var value iec.STRING
	p := CSV_PARSER_FILE{FILENAME: &fn, FSD: &fsd, RUN: &run, OFFSET: &offset, VALUE: &value, PT: &pt}
	p.INIT()
	p.SEP = ';'
	for _, want := range csvWant {
		run = 1
		for i := 0; i < 100 && run != 0; i++ {
			p.Execute(time.Unix(1000, 0))
			fs.Execute(time.Unix(1000, 0))
		}
		if value != want.value || p.RESULT != want.result {
			t.Errorf("got %q %d, want %q %d", value, p.RESULT, want.value, want.result)
		}
	}
	fsd.MODE = 5 // close the file
	for i := 0; i < 100 && fsd.MODE != 0; i++ {
		fs.Execute(time.Unix(1000, 0))
	}
}

const iniData = "; a comment\r\n[net]\r\nip=10.0.0.1\r\nport=80\r\n[log]\r\nlevel=3\r\n"

func TestINIBuffer(t *testing.T) {
	var pt network.NETWORK_BUFFER
	load(&pt, iniData)
	var run iec.BYTE
	var offset iec.UDINT
	var str, key, value iec.STRING
	p := INI_PARSER_BUF{STR: &str, RUN: &run, OFFSET: &offset, KEY: &key, VALUE: &value, PT: &pt}
	p.INIT()
	parse := func(r iec.BYTE, s iec.STRING) {
		t.Helper()
		run, str = r, s
		for i := 0; i < 10 && run != 0; i++ {
			p.Execute(time.Unix(1000, 0))
		}
	}
	parse(1, "net")
	if p.RESULT != 1 || key != "net" {
		t.Fatalf("section: %d %q", p.RESULT, key)
	}
	parse(2, "port")
	// The parse stops at ip=, a key that is not port, in OSCAT.
	if run == 0 {
		t.Fatalf("a key that is not the one asked for: RESULT %d, key %q", p.RESULT, key)
	}

	// Every section and key in turn.
	p.INIT()
	offset = 0
	var got []string
	for i := 0; i < 20; i++ {
		parse(3, "")
		if p.RESULT == 10 {
			break
		}
		got = append(got, string(p.RESULT)+":"+string(key)+"="+string(value))
	}
	want := []string{"\x01:net=", "\x02:ip=10.0.0.1", "\x02:port=80", "\x01:log=", "\x02:level=3"}
	if len(got) != len(want) {
		t.Fatalf("got %q", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%d: got %q, want %q", i, got[i], want[i])
		}
	}
}
