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
	"time"

	. "github.com/apiarytech/beebread/basic"
	str "github.com/apiarytech/beebread/basic/string"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/royaljelly/iec"
)

// DLOG_CRON_TAB sets Q, at each new second of DTI, when DTI matches the
// cron table SECOND, MINUTE, HOUR, DAY_OF_MONTH, DAY_OF_WEEK (1 Monday)
// and MONTH: a field is '*' or empty for all, or a list by ',' of values,
// ranges a-b and steps */n. It reads the table at the rising edge of
// ACTIVATE, and holds Q while ACTIVATE is FALSE. Like cron, a day matches
// the day of the month or the day of the week when both are given.
type DLOG_CRON_TAB struct {
	ACTIVATE     iec.BOOL
	DTI          iec.DT
	SECOND       iec.STRING // STRING(20)
	MINUTE       iec.STRING // STRING(20)
	HOUR         iec.STRING // STRING(20)
	DAY_OF_MONTH iec.STRING // STRING(20)
	DAY_OF_WEEK  iec.STRING // STRING(20)
	MONTH        iec.STRING // STRING(20)
	Q            iec.BOOL

	cd          network.DLOG_CRON_DATA
	activateOld iec.BOOL
	dtiLast     iec.DT
	pos1        iec.INT
	pos2        iec.INT
}

// INIT resets the block.
func (d *DLOG_CRON_TAB) INIT() { *d = DLOG_CRON_TAB{} }

// Execute runs the block once.
func (d *DLOG_CRON_TAB) Execute(now time.Time) {
	if d.ACTIVATE {
		if !d.activateOld {
			d.parse()
		} else if d.DTI != d.dtiLast {
			d.match()
		}
	}
	d.activateOld = d.ACTIVATE
}

// parse reads the table into the elements of the fields.
func (d *DLOG_CRON_TAB) parse() {
	ce := &d.cd.CE
	// OSCAT allows hours up to 59.
	for i, r := range [6][2]iec.INT{{0, 59}, {0, 59}, {0, 59}, {1, 31}, {1, 7}, {1, 12}} {
		ce[i].VALUE_MIN, ce[i].VALUE_MAX = r[0], r[1]
		for j := r[0]; j <= r[1]; j++ {
			ce[i].ELEMENTS[j] = false
		}
		ce[i].ALL_SELECTED = false
	}

	table := CONCAT(d.SECOND, "#", d.MINUTE, "#", d.HOUR, "#", d.DAY_OF_MONTH, "#", d.DAY_OF_WEEK, "#", d.MONTH)
	for i := range iec.INT(6) {
		e := &ce[i]
		field := encoding.ELEMENT_GET('#', i, &table)
		count := iec.INT(1)
		if LEN(field) > 0 {
			count = encoding.ELEMENT_COUNT(',', &field)
		}
		for j := range count {
			item := encoding.ELEMENT_GET(',', j, &field)
			if FIND(item, "*/") > 0 && LEN(item) <= 10 {
				// every step: DEC_TO_INT skips the */
				if step := str.DEC_TO_INT(item); step > 0 {
					for k := e.VALUE_MIN; k <= e.VALUE_MAX; k += step {
						e.ELEMENTS[k] = true
					}
				}
				continue
			}
			if FIND(item, "*") > 0 || LEN(item) == 0 {
				d.pos1, d.pos2 = e.VALUE_MIN, e.VALUE_MAX
				e.ALL_SELECTED = true
			} else if n := LEN(item); n >= 1 && n <= 10 {
				// a value or a range; a longer item keeps the last range
				d.pos1 = str.DEC_TO_INT(encoding.ELEMENT_GET('-', 0, &item))
				d.pos2 = d.pos1
				if encoding.ELEMENT_COUNT('-', &item) == 2 {
					d.pos2 = str.DEC_TO_INT(encoding.ELEMENT_GET('-', 1, &item))
				}
			}
			if d.pos1 >= e.VALUE_MIN && d.pos2 <= e.VALUE_MAX {
				for k := d.pos1; k <= d.pos2; k++ {
					e.ELEMENTS[k] = true
				}
			}
		}
	}
}

// match sets Q for the time DTI.
func (d *DLOG_CRON_TAB) match() {
	ce := &d.cd.CE
	date := DT_TO_DATE(d.DTI)
	ce[0].VALUE = td.SECOND_OF_DT(d.DTI)
	ce[1].VALUE = td.MINUTE_OF_DT(d.DTI)
	ce[2].VALUE = td.HOUR_OF_DT(d.DTI)
	ce[3].VALUE = td.DAY_OF_MONTH(date)
	ce[4].VALUE = td.DAY_OF_WEEK(date)
	ce[5].VALUE = td.MONTH_OF_DATE(date)
	d.dtiLast = d.DTI

	set := func(i int) iec.BOOL { return ce[i].ELEMENTS[ce[i].VALUE] }
	d.Q = set(0) && set(1) && set(2) && set(5)
	day, weekday := set(3), set(4)
	switch {
	case bool(!ce[3].ALL_SELECTED && ce[4].ALL_SELECTED):
		d.Q = d.Q && day
	case bool(ce[3].ALL_SELECTED && !ce[4].ALL_SELECTED):
		d.Q = d.Q && weekday
	case bool(!ce[3].ALL_SELECTED && !ce[4].ALL_SELECTED):
		d.Q = d.Q && (weekday || day)
	}
}
