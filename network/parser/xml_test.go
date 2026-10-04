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
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// readXML returns the items XML_READER finds in doc, as TYP:PATH:value.
func readXML(t *testing.T, doc string) []string {
	t.Helper()
	var buf network.NW_BUF_LONG
	for i := range doc {
		buf[i] = iec.BYTE(doc[i])
	}
	var ctrl network.XML_CONTROL
	x := XML_READER{CTRL: &ctrl, BUF: &buf}
	x.INIT()
	ctrl.START_POS, ctrl.STOP_POS = 0, iec.UINT(len(doc)-1)
	// Every type.
	ctrl.COMMAND = 0x8000 | 1<<1 | 1<<2 | 1<<3 | 1<<4 | 1<<5 | 1<<12 | 1<<13 | 1<<14
	var got []string
	for i := 0; i < 100; i++ {
		x.Execute(time.Unix(1000, 0))
		var v iec.STRING
		switch ctrl.TYP {
		case XML_ELEMENT, XML_PI, XML_CLOSE:
			v = ctrl.ELEMENT
		case XML_ATTRIBUTE:
			v = ctrl.ATTRIBUTE + "=" + ctrl.VALUE
		default:
			v = ctrl.VALUE
		}
		got = append(got, fmt.Sprintf("%d:%s:%s", ctrl.TYP, ctrl.PATH, v))
		if ctrl.TYP == XML_END {
			return got
		}
	}
	t.Fatal("no end")
	return nil
}

func TestXML(t *testing.T) {
	got := readXML(t, `<?xml version="1.0"?><root><a x="1" y='2'>text</a><b/><!-- note --><d>t2</d></root>`)
	want := []string{
		"5:/xml:xml",
		"4:/xml:version=1.0",
		"2::xml",
		"1:/root:root",
		"1:/root/a:a",
		"4:/root/a:x=1",
		"4:/root/a:y=2",
		"3:/root/a:text",
		"2:/root:a",
		"1:/root/b:b",
		"3:/root/b:",
		"2:/root:b",
		"13:/root: note ",
		"1:/root/d:d",
		"3:/root/d:t2",
		"2:/root:d",
		"2::root",
		"99::t2", // OSCAT leaves the last value at the end
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
