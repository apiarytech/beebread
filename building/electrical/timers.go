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

package electrical

import (
	"strconv"
	"strings"
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/engineering"
	str "github.com/apiarytech/beebread/basic/string"
	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

// day is a day in milliseconds.
const day = 86400000

// TIMER_1 is a daily timer: while E is true, Q is on from START for
// DURATION on the days DAY sets, bit 6 Monday to bit 0 Sunday, at the date
// and time DTI. STOP is the time of day Q goes off.
type TIMER_1 struct {
	E        iec.BOOL // default TRUE
	DTI      iec.DT
	START    iec.TOD
	DURATION iec.TIME
	DAY      iec.BYTE // default 2#0111_1111
	Q        iec.BOOL
	STOP     iec.TIME
}

// INIT resets the block and sets its inputs to their initial values.
func (t *TIMER_1) INIT() { *t = TIMER_1{E: true, DAY: 0x7F} }

// Execute runs the block once.
func (t *TIMER_1) Execute(now time.Time) {
	if !t.E {
		t.Q = false
		return
	}
	t.STOP = TOD_TO_TIME(t.START) + t.DURATION
	if t.STOP > iec.TIME(24*time.Hour) {
		t.STOP -= iec.TIME(24 * time.Hour)
	}
	weekday := SHR(iec.BYTE(128), td.DAY_OF_WEEK(DT_TO_DATE(t.DTI)))
	t.Q = td.TIMECHECK(DT_TO_TOD(t.DTI), t.START, DWORD_TO_TOD(ms(t.STOP))) && weekday&t.DAY > 0
}

// TIMER_2 switches Q on at START for DURATION, once a day on the days MODE
// selects, at the date and time DT_IN: 1..7 a day of the week (1 Monday),
// 11 every day, 12..16 every 2nd to 6th day, 20 Monday to Friday, 21 the
// weekend, 22 workdays not HOLIDAY, 23 weekends and holidays, 24 holidays,
// 25 the first and 26 the last day of the month, 27 December 31 and 28
// January 1.
type TIMER_2 struct {
	DT_IN    iec.DT
	START    iec.TOD
	DURATION iec.TIME
	MODE     iec.BYTE
	HOLIDAY  iec.BOOL
	Q        iec.BOOL

	lastCheck  iec.DWORD
	activation iec.DWORD
	init       iec.BOOL
	// runDate is RETAIN in OSCAT: the day Q last went on.
	runDate iec.DATE
}

// INIT resets the block.
func (t *TIMER_2) INIT() { *t = TIMER_2{} }

// Execute runs the block once.
func (t *TIMER_2) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !t.init {
		t.init = true
		t.lastCheck = tx - 100
	}
	// The block runs every 100 ms at most.
	if tx-t.lastCheck < 100 {
		return
	}

	dat := DT_TO_DATE(t.DT_IN)
	daytime := DT_TO_TOD(t.DT_IN)
	wday := td.DAY_OF_WEEK(dat)

	// Whether Q goes on today.
	var enabled bool
	switch m := t.MODE; {
	case m >= 1 && m <= 7:
		enabled = wday == iec.INT(m)
	case m == 11:
		enabled = true
	case m >= 12 && m <= 16:
		enabled = DATE_TO_DWORD(dat)/86400%iec.DWORD(m-10) == 0
	case m == 20:
		enabled = wday <= 5
	case m == 21:
		enabled = wday > 5
	case m == 22:
		enabled = wday <= 5 && !bool(t.HOLIDAY)
	case m == 23:
		enabled = wday > 5 || bool(t.HOLIDAY)
	case m == 24:
		enabled = bool(t.HOLIDAY)
	case m == 25:
		enabled = td.DAY_OF_MONTH(dat) == 1
	case m == 26:
		enabled = td.DAY_OF_MONTH(DWORD_TO_DATE(DATE_TO_DWORD(dat)+86400)) == 1
	case m == 27:
		enabled = td.DAY_OF_MONTH(dat) == 31 && td.MONTH_OF_DATE(dat) == 12
	case m == 28:
		enabled = td.DAY_OF_YEAR(dat) == 1
	}

	switch {
	case enabled && !bool(t.Q) && TOD_TO_DWORD(daytime) >= TOD_TO_DWORD(t.START) && DATE_TO_DWORD(t.runDate) != DATE_TO_DWORD(dat):
		t.Q = true
		t.activation = tx
		t.runDate = dat
	case bool(t.Q) && tx-t.activation >= ms(t.DURATION):
		t.Q = false
	}
	t.lastCheck = tx
}

// TIMER_EVENT_DECODE decodes a timer event written as
// '<typ;channel;day;start;duration;land;lor>': the numbers in decimal or
// as FSTRING_TO_BYTE reads them, the day also as a weekday or (for typ 2)
// a list of weekdays in the language LANG, start as a TOD and duration as a
// TIME, in IEC 61131-3 literals with or without their prefix.
func TIMER_EVENT_DECODE(EVENT iec.STRING, LANG iec.INT) TIMER_EVENT {
	var ev TIMER_EVENT
	// The start and end characters.
	if LEFT(EVENT, 1) != "<" && RIGHT(EVENT, 1) != ">" {
		return ev
	}
	stop := LEN(EVENT)
	pt := CHARS(EVENT)
	start := iec.INT(2)
	step := 0

	// The fields end with ; or >.
	for pos := iec.INT(2); pos <= stop; pos++ {
		if c := pt[pos-1]; c != 59 && c != 62 {
			continue
		}
		tmp := MID(EVENT, pos-start, start)
		switch step {
		case 0:
			ev.TYP = str.FSTRING_TO_BYTE(tmp)
		case 1:
			ev.CHANNEL = str.FSTRING_TO_BYTE(tmp)
		case 2:
			switch {
			case bool(str.IS_CC(tmp, "0123456789abcdefABCDEF#")):
				ev.DAY = str.FSTRING_TO_BYTE(tmp)
			case ev.TYP == 2:
				ev.DAY = str.FSTRING_TO_WEEK(tmp, LANG)
			default:
				ev.DAY = iec.BYTE(str.FSTRING_TO_WEEKDAY(tmp, LANG))
			}
		case 3:
			ev.START = STRING_TO_TOD(tmp)
		case 4:
			ev.DURATION = STRING_TO_TIME(tmp)
		case 5:
			ev.LAND = str.FSTRING_TO_BYTE(tmp)
		case 6:
			ev.LOR = str.FSTRING_TO_BYTE(tmp)
		}
		start = pos + 1
		step++
	}
	return ev
}

// STRING_TO_TOD reads a time of day, as TOD#hh:mm:ss.fff or without the
// prefix (TIME_OF_DAY# too), seconds and milliseconds optional. It is
// TOD#00:00 if s is not a time of day.
func STRING_TO_TOD(s iec.STRING) iec.TOD {
	v := strings.TrimSpace(string(s))
	for _, prefix := range []string{"TIME_OF_DAY#", "TOD#"} {
		if len(v) >= len(prefix) && strings.EqualFold(v[:len(prefix)], prefix) {
			v = v[len(prefix):]
			break
		}
	}
	parts := strings.Split(strings.ReplaceAll(v, "_", ""), ":")
	if len(parts) < 2 || len(parts) > 3 {
		return DWORD_TO_TOD(0)
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	sec := 0.0
	var err3 error
	if len(parts) == 3 {
		sec, err3 = strconv.ParseFloat(parts[2], 64)
	}
	if err1 != nil || err2 != nil || err3 != nil || h < 0 || h > 23 || m < 0 || m > 59 || sec < 0 || sec >= 60 {
		return DWORD_TO_TOD(0)
	}
	return DWORD_TO_TOD(iec.DWORD(h*3600000 + m*60000 + int(sec*1000+0.5)))
}

// STRING_TO_TIME reads a duration, as T#1d2h3m4s5ms or without the prefix
// (TIME# too), each unit optional, the last with a fraction. It is T#0s if
// s is not a duration.
func STRING_TO_TIME(s iec.STRING) iec.TIME {
	v := strings.TrimSpace(string(s))
	for _, prefix := range []string{"TIME#", "T#"} {
		if len(v) >= len(prefix) && strings.EqualFold(v[:len(prefix)], prefix) {
			v = v[len(prefix):]
			break
		}
	}
	v = strings.ToLower(strings.ReplaceAll(v, "_", ""))
	neg := strings.HasPrefix(v, "-")
	v = strings.TrimPrefix(v, "-")
	if v == "" {
		return 0
	}
	units := []struct {
		name string
		d    time.Duration
	}{{"ms", time.Millisecond}, {"us", time.Microsecond}, {"ns", time.Nanosecond},
		{"d", 24 * time.Hour}, {"h", time.Hour}, {"m", time.Minute}, {"s", time.Second}}
	var total time.Duration
	for v != "" {
		i := 0
		for i < len(v) && (v[i] >= '0' && v[i] <= '9' || v[i] == '.') {
			i++
		}
		n, err := strconv.ParseFloat(v[:i], 64)
		if i == 0 || err != nil {
			return 0
		}
		v = v[i:]
		found := false
		for _, u := range units {
			if strings.HasPrefix(v, u.name) {
				total += time.Duration(n * float64(u.d))
				v = v[len(u.name):]
				found = true
				break
			}
		}
		if !found {
			return 0
		}
	}
	if neg {
		total = -total
	}
	return iec.TIME(total)
}

// TIMER_EXT is a timer for a light or a blind: Q goes on and off at the
// times of day T_DAY_START and T_DAY_STOP, T_RISE_START before and
// T_RISE_STOP after SUN_RISE, and T_SET_START before and T_SET_STOP after
// SUN_SET, each if not 0, at the date and time DT_IN. Saturdays, Sundays and
// holidays turn it off unless enabled. A click of SWITCH toggles it, and a
// MANUAL_2 gives the manual control by ENA, ON, OFF and MAN. It runs every
// 200 ms at most.
type TIMER_EXT struct {
	ENA             iec.BOOL // default TRUE
	ON              iec.BOOL
	OFF             iec.BOOL
	MAN             iec.BOOL
	SWITCH          iec.BOOL
	DT_IN           iec.DT
	SUN_RISE        iec.TOD
	SUN_SET         iec.TOD
	HOLIDAY         iec.BOOL
	T_DEBOUNCE      iec.TIME // default T#100ms
	T_RISE_START    iec.TIME
	T_RISE_STOP     iec.TIME
	T_SET_START     iec.TIME
	T_SET_STOP      iec.TIME
	T_DAY_START     iec.TOD
	T_DAY_STOP      iec.TOD
	ENABLE_SATURDAY iec.BOOL
	ENABLE_SUNDAY   iec.BOOL
	ENABLE_HOLIDAY  iec.BOOL
	Q               iec.BOOL
	STATUS          iec.BYTE

	mx   engineering.MANUAL_2
	deb  DEBOUNCE
	tl   iec.DWORD
	qx   iec.BOOL
	init iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (t *TIMER_EXT) INIT() {
	*t = TIMER_EXT{ENA: true, T_DEBOUNCE: iec.TIME(100 * time.Millisecond)}
}

// Execute runs the block once.
func (t *TIMER_EXT) Execute(now time.Time) {
	// Sunrise and sunset in whole seconds.
	t.SUN_RISE = DWORD_TO_TOD(TOD_TO_DWORD(t.SUN_RISE) / 1000 * 1000)
	t.SUN_SET = DWORD_TO_TOD(TOD_TO_DWORD(t.SUN_SET) / 1000 * 1000)

	// The cycle time tc; the block runs every 200 ms at most.
	tx := PLC_MS(now)
	if !t.init {
		t.init = true
		t.tl = tx
	}
	tc := tx - t.tl
	if tc < 200 {
		return
	}
	t.tl = tx

	t.deb.IN, t.deb.TD, t.deb.PM = t.SWITCH, t.T_DEBOUNCE, true
	t.deb.Execute(now)

	// TOD - TOD is a TIME, in milliseconds modulo 2^32.
	tdx := TOD_TO_DWORD(DT_TO_TOD(t.DT_IN))
	wdx := td.DAY_OF_WEEK(DT_TO_DATE(t.DT_IN))
	rise, set := TOD_TO_DWORD(t.SUN_RISE), TOD_TO_DWORD(t.SUN_SET)
	dayStart, dayStop := TOD_TO_DWORD(t.T_DAY_START), TOD_TO_DWORD(t.T_DAY_STOP)

	// The automatic control.
	switch {
	case bool(t.deb.Q):
		t.qx = !t.qx
		t.STATUS = 110
	case bool(t.HOLIDAY) && !bool(t.ENABLE_HOLIDAY):
		t.qx = false
	case wdx == 6 && !bool(t.ENABLE_SATURDAY):
		t.qx = false
	case wdx == 7 && !bool(t.ENABLE_SUNDAY):
		t.qx = false
	case dayStart > 0 && tdx-dayStart <= tc:
		// on at a time of day
		t.qx = true
		t.STATUS = 111
	case dayStop > 0 && tdx-dayStop <= tc:
		// off at a time of day
		t.qx = false
		t.STATUS = 112
	case t.T_RISE_START > 0 && tdx-rise+ms(t.T_RISE_START) <= tc:
		// on before sunrise
		t.qx = true
		t.STATUS = 113
	case t.T_RISE_STOP > 0 && tdx-rise-ms(t.T_RISE_STOP) <= tc:
		// off after sunrise
		t.qx = false
		t.STATUS = 114
	case t.T_SET_START > 0 && tdx-set+ms(t.T_SET_START) <= tc:
		// on before sunset
		t.qx = true
		t.STATUS = 115
	case t.T_SET_STOP > 0 && tdx-set-ms(t.T_SET_STOP) <= tc:
		// off after sunset
		t.qx = false
		t.STATUS = 116
	}

	// The manual control.
	t.mx.IN, t.mx.ENA, t.mx.ON, t.mx.OFF, t.mx.MAN = t.qx, t.ENA, t.ON, t.OFF, t.MAN
	t.mx.Execute(now)
	t.Q = t.mx.Q
	if t.mx.STATUS > 100 {
		t.STATUS = t.mx.STATUS
	}
}

// TIMER_P4 is a programmable timer of four channels, Q0..Q3, for the
// channels OFS..OFS+3 of the timer events in PROG, at the date and time
// DTIME. An event's TYP selects when it is on:
//
//	1 daily, 2 the weekdays DAY sets (bit 6 Monday), 3 every DAY days,
//	10 weekly on the weekday DAY, 20 monthly on the day DAY, 21 the last day
//	of the month, 30 yearly on the day of the year DAY, 31 the last day of
//	the year, 40 leap days, 41 holidays (HOLY), 42 holidays and weekends,
//	43 Monday to Friday: from START for DURATION;
//	50 and 51 START after or before the reference time TREF_0 or TREF_1
//	(by DAY), for DURATION; 52 and 53 set and clear the output START after
//	it, 54 and 55 START before it.
//
// An output is on if its event is and the inputs L0..L3 match the event's
// mask LAND, or they match LOR; ENQ enables the outputs and MAN with MI
// sets them by hand. RST makes the set and clear events ready to run again.
type TIMER_P4 struct {
	DTIME          iec.DT
	TREF_0         iec.TOD
	TREF_1         iec.TOD
	HOLY           iec.BOOL
	L0, L1, L2, L3 iec.BOOL
	OFS            iec.BYTE
	ENQ            iec.BOOL
	MAN            iec.BOOL
	MI             iec.BYTE
	RST            iec.BOOL
	PROG           *[TIMER_P4_ARRAY_MAX + 1]TIMER_EVENT
	Q0, Q1, Q2, Q3 iec.BOOL
	STATUS         iec.BYTE

	lastExecute iec.DT
	// currentDay is never set in OSCAT: TYP 3 is on every day.
	currentDay iec.DINT
	ma, mo     [TIMER_P4_CHANNEL_MAX + 1]iec.BYTE
	qn, qs     [TIMER_P4_CHANNEL_MAX + 1]iec.BOOL
}

// The constants of TIMER_P4: the last index of PROG and of the channels.
const (
	TIMER_P4_ARRAY_MAX   = 63
	TIMER_P4_CHANNEL_MAX = 3
)

// INIT resets the block.
func (t *TIMER_P4) INIT() { *t = TIMER_P4{PROG: t.PROG} }

// Execute runs the block once.
func (t *TIMER_P4) Execute(now time.Time) {
	if t.PROG == nil {
		return
	}
	prog := t.PROG
	dtime := DT_TO_DWORD(t.DTIME)

	if t.RST {
		// The events of the channels were last active on 1970-01-01.
		t.lastExecute = DWORD_TO_DT(0)
		for pos := range prog {
			if prog[pos].CHANNEL >= t.OFS && iec.INT(prog[pos].CHANNEL) < iec.INT(t.OFS)+4 {
				prog[pos].LAST = t.lastExecute
			}
		}
		t.qs = [TIMER_P4_CHANNEL_MAX + 1]iec.BOOL{}
	} else if dtime != DT_TO_DWORD(t.lastExecute) {
		// The block runs once a second.
		t.lastExecute = t.DTIME
		date := DT_TO_DATE(t.DTIME)
		dayStart := DATE_TO_DWORD(date)
		wday := td.DAY_OF_WEEK(date)
		t.qn = [TIMER_P4_CHANNEL_MAX + 1]iec.BOOL{}

		// DT + TIME, in seconds.
		at := func(offset iec.TIME, sign iec.DWORD) iec.DWORD { return dayStart + sign*(ms(offset)/1000) }
		ref := func(ev TIMER_EVENT) iec.TIME {
			switch ev.DAY {
			case 0:
				return TOD_TO_TIME(t.TREF_0)
			case 1:
				return TOD_TO_TIME(t.TREF_1)
			}
			return 0
		}
		for pos := range prog {
			ev := prog[pos]
			channel := iec.INT(ev.CHANNEL) - iec.INT(t.OFS)
			if ev.TYP == 0 || channel < 0 || channel > TIMER_P4_CHANNEL_MAX {
				continue
			}
			t.ma[channel] = ev.LAND
			t.mo[channel] = ev.LOR
			// on sets qn for the event's DURATION from start.
			on := func(start iec.DWORD) {
				t.qn[channel] = dtime >= start && dtime <= start+ms(ev.DURATION)/1000
			}
			// set sets or clears qs once a day from start.
			set := func(start iec.DWORD, q iec.BOOL) {
				if dtime >= start && dayStart > DT_TO_DWORD(ev.LAST) {
					t.qs[channel] = q
					prog[pos].LAST = DWORD_TO_DT(dayStart)
				}
			}
			start := at(TOD_TO_TIME(ev.START), 1)
			switch ev.TYP {
			case 1: // daily
				on(start)
			case 2: // the weekdays DAY sets
				if SHR(iec.BYTE(128), wday)&ev.DAY > 0 {
					on(start)
				}
			case 3: // every DAY days
				if ev.DAY == 0 || t.currentDay%iec.DINT(ev.DAY) == 0 {
					on(start)
				}
			case 10: // weekly
				if wday == iec.INT(ev.DAY) {
					on(start)
				}
			case 20: // monthly
				if td.DAY_OF_MONTH(date) == iec.INT(ev.DAY) {
					on(start)
				}
			case 21: // the last day of the month
				if DATE_TO_DWORD(date) == DATE_TO_DWORD(td.MONTH_END(date)) {
					on(start)
				}
			case 30: // yearly
				if td.DAY_OF_YEAR(date) == iec.INT(ev.DAY) {
					on(start)
				}
			case 31: // the last day of the year
				if DATE_TO_DWORD(date) == DATE_TO_DWORD(td.YEAR_END(td.YEAR_OF_DATE(date))) {
					on(start)
				}
			case 40: // leap days
				if td.LEAP_DAY(date) {
					on(start)
				}
			case 41: // holidays
				if t.HOLY {
					on(start)
				}
			case 42: // holidays and weekends
				if bool(t.HOLY) || wday == 6 || wday == 7 {
					on(start)
				}
			case 43: // Monday to Friday
				if wday < 6 {
					on(start)
				}
			case 50: // after the reference time
				on(at(ref(ev), 1) + ms(TOD_TO_TIME(ev.START))/1000)
			case 51: // before the reference time
				on(at(ref(ev), 1) - ms(TOD_TO_TIME(ev.START))/1000)
			case 52: // set after the reference time
				set(at(ref(ev), 1)+ms(TOD_TO_TIME(ev.START))/1000, true)
			case 53: // clear after the reference time
				set(at(ref(ev), 1)+ms(TOD_TO_TIME(ev.START))/1000, false)
			case 54: // set before the reference time
				set(at(ref(ev), 1)-ms(TOD_TO_TIME(ev.START))/1000, true)
			case 55: // clear before the reference time
				set(at(ref(ev), 1)-ms(TOD_TO_TIME(ev.START))/1000, false)
			}
		}
	}

	// The logic input mask.
	mask := iec.BYTE(0xF0)
	for i, l := range [4]iec.BOOL{t.L0, t.L1, t.L2, t.L3} {
		if l {
			mask |= 1 << i
		}
	}

	// The outputs; OSCAT compares LOR AND mask with the channel number.
	out := func(c int) iec.BOOL {
		return t.ENQ && ((t.qn[c] || t.qs[c]) && t.ma[c]&mask == t.ma[c] ||
			t.mo[c]&mask > iec.BYTE(c) || t.MAN && iec.BOOL(BIT(t.MI, c)))
	}
	t.Q0, t.Q1, t.Q2, t.Q3 = out(0), out(1), out(2), out(3)

	switch {
	case bool(!t.ENQ):
		t.STATUS = 100
	case bool(t.MAN):
		t.STATUS = 101
	default:
		t.STATUS = 102
	}
}
