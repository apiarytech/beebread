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

package dlog

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

func TestRing(t *testing.T) {
	var d network.UNI_CIRCULAR_BUFFER_DATA
	var u UNI_CIRCULAR_BUFFER
	u.DATA = &d
	for _, s := range []iec.STRING{"one", "two"} {
		d.D_HEAD, d.D_STRING, d.D_MODE = 0x0101, s, ucbAdd
		u.Execute(time.Time{})
	}
	if d.BUF_COUNT != 2 {
		t.Fatalf("count %d", d.BUF_COUNT)
	}
	d.D_MODE = ucbReadRemove
	u.Execute(time.Time{})
	if d.D_STRING != "one" || d.D_HEAD != 0x0101 || d.BUF_COUNT != 1 {
		t.Errorf("read %q %x %d", d.D_STRING, d.D_HEAD, d.BUF_COUNT)
	}
}

// logger is a store with an INT and a STRING column.
type logger interface{ Execute(time.Time) }

// logTo runs the store s, on the data x, with two values: a log at the
// rising edge of TRIG_M, then the end of the log.
func logTo(t *testing.T, x *network.DLOG_DATA, s logger, enable, trig *iec.BOOL) {
	t.Helper()
	a := DLOG_DINT{VALUE: 5, COLUMN: "A"}
	a.X = x
	b := DLOG_STRING{STR: "x", COLUMN: "B"}
	b.X = x
	now := time.Time{}
	scan := func(n int) {
		for range n {
			s.Execute(now)
			a.Execute(now)
			b.Execute(now)
		}
	}
	*enable = true
	scan(5)
	*trig = true
	scan(5)
	*trig = false
	a.VALUE, b.STR = 6, "y"
	scan(1)
	*trig = true
	scan(5)
	*enable = false
	scan(100) // each file command takes a few scans
}

func TestStoreCSV(t *testing.T) {
	name := filepath.Join(t.TempDir(), "log.csv")
	var x network.DLOG_DATA
	var save network.DLOG_SAVE
	s := DLOG_STORE_FILE_CSV{}
	s.X, s.SAVE_DATA = &x, &save
	s.INIT()
	s.SEP = ';'
	s.FILENAME, s.DTI = iec.STRING(name), td.SET_DT(2026, 10, 3, 12, 0, 0)
	logTo(t, &x, &s, &s.ENABLE, &s.TRIG_M)
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err, s.ERROR_C, s.ERROR_T)
	}
	if got := string(b); got != "A;B\r\n5;x\r\n6;y\r\n" {
		t.Errorf("file %q", got)
	}
	if save.TRIG_CNT != 2 || x.ID_MAX != 2 || s.ERROR_T != 0 {
		t.Errorf("count %d, columns %d, error %d", save.TRIG_CNT, x.ID_MAX, s.ERROR_T)
	}
}

func TestStoreXML(t *testing.T) {
	name := filepath.Join(t.TempDir(), "log.xml")
	var x network.DLOG_DATA
	var save network.DLOG_SAVE
	s := DLOG_STORE_FILE_XML{}
	s.X, s.SAVE_DATA = &x, &save
	s.INIT()
	s.FILENAME, s.DTI = iec.STRING(name), td.SET_DT(2026, 10, 3, 12, 0, 0)
	logTo(t, &x, &s, &s.ENABLE, &s.TRIG_M)
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	want := `<?xml version="1.0" encoding="UTF-8"?><table>` +
		"<row><entry>A</entry><entry>B</entry></row>" +
		"<row><entry>5</entry><entry>x</entry></row>" +
		"<row><entry>6</entry><entry>y</entry></row></table>"
	if string(b) != want {
		t.Errorf("file %q", b)
	}
}

func TestStoreHTML(t *testing.T) {
	name := filepath.Join(t.TempDir(), "log.html")
	var x network.DLOG_DATA
	var save network.DLOG_SAVE
	s := DLOG_STORE_FILE_HTML{}
	s.X, s.SAVE_DATA = &x, &save
	s.INIT()
	s.HTML_CAPTION, s.HTML_TABLE, s.HTML_TR_HEAD, s.HTML_TR_EVEN, s.HTML_TR_ODD = "c", "t", "h", "e", "o"
	s.FILENAME, s.DTI = iec.STRING(name), td.SET_DT(2026, 10, 3, 12, 0, 0)
	logTo(t, &x, &s, &s.ENABLE, &s.TRIG_M)
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	want := "<html><body><table t><caption>c</caption>" +
		"<TR h><TD>A</TD><TD>B</TD></TR>" +
		"<TR e><TD>5</TD><TD>x</TD></TR>" +
		"<TR o><TD>6</TD><TD>y</TD></TR></table></body></html>"
	if string(b) != want {
		t.Errorf("file %q", b)
	}
}
