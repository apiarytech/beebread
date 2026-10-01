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

package time_date

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/royaljelly/iec"
)

// HOLIDAY_DE, HOLIDAY_AT, HOLIDAY_FR, HOLIDAY_BED and HOLIDAY_ITD are the
// holiday lists OSCAT gives as samples for Germany, Austria, France, the
// German speaking part of Belgium and South Tyrol, for HOLIDAY and
// CALENDAR_CALC. An entry with USE 0 is not used.
var (
	HOLIDAY_DE = holidays(
		h("Neujahr", 1, 1, 1),
		h("Heilig Drei Könige", 6, 1, 1),
		h("Karfreitag", -2, 0, 1),
		h("Ostersonntag", 0, 0, 1),
		h("Ostermontag", 1, 0, 1),
		h("Tag der Arbeit", 1, 5, 1),
		h("Christi Himmelfahrt", 39, 0, 1),
		h("Pfingstsonntag", 49, 0, 1),
		h("Pfingstmontag", 50, 0, 1),
		h("Fronleichnam", 60, 0, 1),
		h("Augsburger Friedensfest", 8, 8, 0),
		h("Maria Himmelfahrt", 15, 8, 1),
		h("Tag der Deutschen Einheit", 3, 10, 1),
		h("Reformationstag", 31, 10, 0),
		h("Allerheiligen", 1, 11, 1),
		h("Buss und Bettag", 23, 11, 0),
		h("1. Weihnachtstag", 25, 12, 1),
		h("2. Weihnachtstag", 26, 12, 1),
	)
	HOLIDAY_AT = holidays(
		h("Neujahr", 1, 1, 1),
		h("Heilig Drei Könige", 6, 1, 1),
		h("Karfreitag", -2, 0, 1),
		h("Ostersonntag", 0, 0, 1),
		h("Ostermontag", 1, 0, 1),
		h("", 1, 5, 0),
		h("Christi Himmelfahrt", 39, 0, 1),
		h("Pfingstsonntag", 49, 0, 1),
		h("Pfingstmontag", 50, 0, 1),
		h("Fronleichnam", 60, 0, 1),
		h("", 8, 8, 0),
		h("Maria Himmelfahrt", 15, 8, 1),
		h("", 3, 10, 0),
		h("", 31, 10, 0),
		h("Allerheiligen", 1, 11, 1),
		h("Maria Empfängnis", 8, 12, 1),
		h("1. Weihnachtstag", 25, 12, 1),
		h("2. Weihnachtstag", 26, 12, 1),
	)
	HOLIDAY_FR = holidays(
		h("Nouvel an", 1, 1, 1),
		h("St Valentin", 14, 2, 0),
		h("Vendredi Saint (alsace)", -2, 0, 0),
		h("Dimanche de pâques", 0, 0, 1),
		h("Lundi de pâques", 1, 0, 1),
		h("Jeudi de Ascension", 39, 0, 1),
		h("dimanche de Pentecôte ", 49, 0, 1),
		h("jeudi de la Trinité", 60, 0, 0),
		h("Fête du travail", 1, 5, 1),
		h("Victoire 1945 ", 8, 5, 1),
		h("Prise de La bastille", 14, 7, 1),
		h("15 Août 1944", 15, 8, 1),
		h("Halloween", 31, 10, 0),
		h("Armistice 1918", 11, 11, 1),
		h("Noël", 25, 12, 1),
		h("Saint Étienne (alsace)", 26, 12, 0),
		h("Fête de la musique", 21, 6, 0),
	)
	HOLIDAY_BED = holidays(
		h("Neujahr", 1, 1, 1),
		h("Ostersonntag", 0, 0, 1),
		h("Ostermontag", 1, 0, 1),
		h("Tag der Arbeit", 1, 5, 1),
		h("Christi Himmelfahrt", 39, 0, 1),
		h("Pfingsten", 49, 0, 1),
		h("Pfingstmontag", 50, 0, 1),
		h("Nationalfeiertag", 21, 7, 1),
		h("Mariä Himmelfahrt", 15, 8, 1),
		h("Allerheiligen", 1, 11, 1),
		h("Feiertag DG", 15, 11, 1),
		h("Heiligabend", 24, 12, 1),
		h("1. Weihnachtstag", 25, 12, 1),
		h("2. Weihnachtstag", 26, 12, 1),
		h("Silvester", 31, 12, 1),
	)
	HOLIDAY_ITD = holidays(
		h("Neujahr", 1, 1, 1),
		h("Heilig Drei Könige", 6, 1, 1),
		h("Ostersonntag", 0, 0, 1),
		h("Ostermontag", 1, 0, 1),
		h("Tag der Befeiung Italiens", 25, 4, 1),
		h("Tag der Arbeit", 1, 5, 1),
		h("Pfingsten", 49, 0, 1),
		h("Pfingstmontag", 50, 0, 1),
		h("Tag der Republik Italien", 2, 6, 1),
		h("Mariä Himmelfahrt", 15, 8, 1),
		h("Allerheiligen", 1, 11, 1),
		h("Mariä Empfängnis", 8, 12, 1),
		h("Heiligabend", 24, 12, 0),
		h("1. Weihnachtstag", 25, 12, 1),
		h("2. Stephanstag", 26, 12, 1),
	)
)

func h(name iec.STRING, day, month, use iec.SINT) HOLIDAY_DATA {
	return HOLIDAY_DATA{NAME: name, DAY: day, MONTH: month, USE: use}
}

func holidays(list ...HOLIDAY_DATA) [30]HOLIDAY_DATA {
	var out [30]HOLIDAY_DATA
	copy(out[:], list)
	return out
}

// CALENDAR_CALC calculates the values of a CALENDAR from its world time
// UTC and location: local time, weekday, daylight saving time, holiday, sun
// rise and sun set, and, if SPE is true, the sun's position every 25
// seconds. H is the sun's height at sun rise and sun set.
type CALENDAR_CALC struct {
	SPE      iec.BOOL
	H        iec.REAL // default -0.83333333333
	XCAL     *CALENDAR
	HOLIDAYS *[30]HOLIDAY_DATA

	last     iec.DT
	lastDay  iec.DINT
	holy     HOLIDAY
	sun      SUN_TIME
	lastHour iec.INT
	pos      SUN_POS
	plast    iec.DT
}

// INIT resets the block and sets H to its initial value.
func (c *CALENDAR_CALC) INIT() {
	*c = CALENDAR_CALC{XCAL: c.XCAL, HOLIDAYS: c.HOLIDAYS, H: -0.83333333333}
}

// Execute runs the block once.
func (c *CALENDAR_CALC) Execute(now time.Time) {
	x := c.XCAL
	if x == nil || dtw(x.UTC) == dtw(c.last) {
		return
	}
	// Once a second.
	c.last = x.UTC
	x.LOCAL_DT = UTC_TO_LTIME(x.UTC, x.DST_EN, x.OFFSET)
	x.LOCAL_DATE = DT_TO_DATE(x.LOCAL_DT)
	x.LOCAL_TOD = DT_TO_TOD(x.LOCAL_DT)
	dtemp := DAY_OF_DATE(x.LOCAL_DATE)
	x.NIGHT = tw(x.LOCAL_TOD) < tw(x.SUN_RISE) || tw(x.LOCAL_TOD) > tw(x.SUN_SET)
	// Once an hour.
	if tmp := HOUR(x.LOCAL_TOD); tmp != c.lastHour {
		x.DST_ON = DST(x.UTC) && x.DST_EN
		c.lastHour = tmp
	}
	// Once a day.
	if dtemp != c.lastDay {
		c.lastDay = dtemp
		x.YEAR = YEAR_OF_DATE(x.LOCAL_DATE)
		x.MONTH = MONTH_OF_DATE(x.LOCAL_DATE)
		x.DAY = DAY_OF_MONTH(x.LOCAL_DATE)
		x.WEEKDAY = DAY_OF_WEEK(x.LOCAL_DATE)
		c.holy.DATE_IN = x.LOCAL_DATE
		c.holy.LANGU = x.LANGUAGE
		c.holy.HOLIDAYS = c.HOLIDAYS
		c.holy.Execute(now)
		x.HOLIDAY = c.holy.Y
		x.HOLY_NAME = c.holy.NAME
		c.sun.LATITUDE = x.LATITUDE
		c.sun.LONGITUDE = x.LONGITUDE
		c.sun.UTC = DT_TO_DATE(x.UTC)
		c.sun.H = c.H
		c.sun.Execute(now)
		local := func(t iec.TOD) iec.TOD {
			ms := iec.DINT(tw(t)) + iec.DINT(x.OFFSET)*60000 + SEL[iec.DINT](x.DST_ON, 0, 3600000)
			return DWORD_TO_TOD(iec.DWORD(ms))
		}
		x.SUN_RISE = local(c.sun.SUN_RISE)
		x.SUN_SET = local(c.sun.SUN_SET)
		x.SUN_MIDDAY = local(c.sun.MIDDAY)
		x.SUN_HEIGTH = c.sun.SUN_DECLINATION
		x.WORK_WEEK = WORK_WEEK(x.LOCAL_DATE)
	}
	// The sun's position every 25 seconds.
	if c.SPE && dtw(x.UTC)-dtw(c.plast) >= 25 {
		c.plast = c.last
		c.pos.LATITUDE = x.LATITUDE
		c.pos.LONGITUDE = x.LONGITUDE
		c.pos.UTC = x.UTC
		c.pos.Execute(now)
		x.SUN_HOR = c.pos.B
		x.SUN_VER = c.pos.HR
	}
}

// DCF77 decodes the time signal of a DCF77 receiver on REC. It needs two
// successive valid minutes before it clears ERROR, and then runs two clocks:
// RTC in world time and RTC1 in the time zone TIME_OFFSET hours away, with
// daylight saving time if DST_EN is true. TP is true for one scan when the
// clocks are set; SET sets them to SDT, with DSI as the daylight saving time
// flag. SYNC is false after SYNC_TIMEOUT without a valid signal, while the
// clocks keep running. MSEC is the milliseconds of the second.
type DCF77 struct {
	REC          iec.BOOL
	SET          iec.BOOL
	SDT          iec.DT
	DSI          iec.BOOL
	SYNC_TIMEOUT iec.TIME // default T#2m
	TIME_OFFSET  iec.INT  // default 1
	DST_EN       iec.BOOL // default TRUE
	TP           iec.BOOL
	DS           iec.BOOL
	WDAY         iec.INT
	ERROR        iec.BOOL // default TRUE
	RTC          iec.DT
	RTC1         iec.DT
	MSEC         iec.INT
	SYNC         iec.BOOL

	mez, utc    iec.DT
	state       iec.INT
	edge        iec.BOOL
	ty, last    iec.DWORD
	bits        [59]iec.BOOL
	oldTime     iec.DT
	lastSync    iec.DWORD
	init        iec.BOOL
	initialized bool
}

// INIT resets the block and sets SYNC_TIMEOUT, TIME_OFFSET and DST_EN to
// their initial values.
func (d *DCF77) INIT() {
	*d = DCF77{SYNC_TIMEOUT: iec.TIME(2 * time.Minute), TIME_OFFSET: 1, DST_EN: true, ERROR: true, initialized: true}
}

// Execute runs the block once.
func (d *DCF77) Execute(now time.Time) {
	if !d.initialized {
		d.initialized = true
		d.ERROR = true
	}
	d.TP = false
	t1 := PLC_MS(now)
	tx := t1 - d.last
	if d.REC != d.edge {
		d.edge = d.REC
		switch {
		case !bool(d.REC) && tx > 1700 && tx < 2000:
			// The start of a minute.
			d.state = 0
			d.TP = !d.ERROR
		case !bool(d.REC) && tx > 700 && tx < 1000:
			// The next second.
			if d.state < 58 {
				d.state++
			} else {
				d.state = 0
			}
		case bool(d.REC) && tx < 120:
			d.bits[d.state] = false
		case bool(d.REC) && tx > 120 && tx < 250:
			d.bits[d.state] = true
		default:
			// The signal is not valid.
			d.ERROR = true
			d.state = 0
		}
		d.last = d.last + tx
		if d.REC && d.state == 58 {
			d.decode()
		}
	}

	// A clock that runs by itself and is set by the DCF77 signal.
	tz := iec.DWORD(int32(max(d.TIME_OFFSET, -d.TIME_OFFSET))) * 3600
	if !d.init || d.SET {
		d.init = true
		d.utc = d.SDT
		d.TP = true
		d.DS = d.DSI
	}
	if d.TP {
		d.RTC = d.utc
		rtc1 := dtw(d.RTC)
		if d.TIME_OFFSET < 0 {
			rtc1 -= tz
		} else {
			rtc1 += tz
		}
		if d.DS && d.DST_EN {
			rtc1 += 3600
		}
		d.RTC1 = DWORD_TO_DT(rtc1)
		d.SYNC = true
		d.lastSync = d.last
		d.ty = d.last
	} else if dtw(d.RTC) > 0 && t1-d.ty >= 1000 {
		d.RTC = DWORD_TO_DT(dtw(d.RTC) + 1)
		d.RTC1 = DWORD_TO_DT(dtw(d.RTC1) + 1)
		d.ty += 1000
		d.SYNC = d.ty-d.lastSync < TIME_TO_DWORD(d.SYNC_TIMEOUT) && d.lastSync > 0
		d.WDAY = DAY_OF_WEEK(DT_TO_DATE(d.RTC1))
		d.DS = d.DST_EN && DST(d.utc)
	}
	d.MSEC = iec.INT(t1 - d.ty)
}

// decode decodes the bits of a minute.
func (d *DCF77) decode() {
	b := &d.bits
	num := func(from int, weights ...iec.INT) iec.INT {
		var n iec.INT
		for i, w := range weights {
			if b[from+i] {
				n += w
			}
		}
		return n
	}
	parity := func(from, to int) bool {
		p := false
		for i := from; i <= to; i++ {
			p = p != bool(b[i])
		}
		return p
	}
	d.ERROR = false
	if b[0] || b[17] == b[18] || !b[20] {
		d.ERROR = true
	}
	minute := num(21, 1, 2, 4, 8, 10, 20, 40)
	if minute > 59 || parity(21, 28) {
		d.ERROR = true
	}
	hour := num(29, 1, 2, 4, 8, 10, 20)
	if hour > 23 || parity(29, 35) {
		d.ERROR = true
	}
	day := num(36, 1, 2, 4, 8, 10, 20)
	if day > 31 {
		d.ERROR = true
	}
	d.WDAY = num(42, 1, 2, 4)
	if d.WDAY > 7 || d.WDAY < 1 {
		d.ERROR = true
	}
	month := num(45, 1, 2, 4, 8, 10)
	if month > 12 {
		d.ERROR = true
	}
	year := num(50, 1, 2, 4, 8, 10, 20, 40, 80)
	var cnt iec.DINT
	for i := 36; i <= 58; i++ {
		if b[i] {
			cnt++
		}
	}
	if !math.EVEN(cnt) {
		d.ERROR = true
	}
	// The time must be valid for two minutes to clear the error.
	if !d.ERROR {
		d.oldTime = d.mez
		if year >= 70 {
			year += 1900
		} else {
			year += 2000
		}
		d.mez = SET_DT(year, month, day, hour, minute, 0)
		d.DS = b[17]
		if d.DS {
			d.utc = DWORD_TO_DT(dtw(d.mez) - 7200)
		} else {
			d.utc = DWORD_TO_DT(dtw(d.mez) - 3600)
		}
		if dtw(d.mez) != dtw(d.oldTime)+60 {
			d.ERROR = true
		}
	}
}

// EVENTS reports whether DATE_IN is one of the events of ELIST, each
// USE days long from its DAY and MONTH, and gives its NAME. The outputs are
// off while ENA is false.
type EVENTS struct {
	DATE_IN iec.DATE
	ENA     iec.BOOL
	Y       iec.BOOL
	NAME    iec.STRING // STRING(30)
	ELIST   *[50]HOLIDAY_DATA

	lastActive iec.DATE
	yInt       iec.BOOL
	nameInt    iec.STRING
}

// INIT resets the block.
func (e *EVENTS) INIT() { *e = EVENTS{ELIST: e.ELIST} }

// Execute runs the block once.
func (e *EVENTS) Execute(now time.Time) {
	// Search the list once a day.
	if dw(e.lastActive) != dw(e.DATE_IN) && e.ELIST != nil {
		e.lastActive = e.DATE_IN
		e.yInt = false
		e.nameInt = ""
		dayIn := DAY_OF_DATE(e.DATE_IN)
		cyr := YEAR_OF_DATE(e.DATE_IN)
		for _, check := range e.ELIST {
			lday := DAY_OF_DATE(SET_DATE(cyr, iec.INT(check.MONTH), iec.INT(check.DAY)))
			if dayIn >= lday && dayIn <= lday+iec.DINT(check.USE)-1 {
				e.yInt = true
				e.nameInt = check.NAME
				break
			}
		}
	}
	if e.ENA {
		e.Y = e.yInt
		e.NAME = e.nameInt
	} else {
		e.Y = false
		e.NAME = ""
	}
}

// HOLIDAY reports whether DATE_IN is one of the HOLIDAYS and gives its NAME.
// A holiday has a fixed date, an offset from easter sunday, or is a weekday
// before a date; see HOLIDAY_DATA. FRIDAY, SATURDAY and SUNDAY make those
// weekdays holidays too, named in the language LANGU.
type HOLIDAY struct {
	DATE_IN  iec.DATE
	LANGU    iec.INT
	FRIDAY   iec.BOOL
	SATURDAY iec.BOOL
	SUNDAY   iec.BOOL
	HOLIDAYS *[30]HOLIDAY_DATA
	Y        iec.BOOL
	NAME     iec.STRING // STRING(30)

	lastActive iec.DATE
}

// INIT resets the block.
func (h *HOLIDAY) INIT() { *h = HOLIDAY{HOLIDAYS: h.HOLIDAYS} }

// Execute runs the block once.
func (h *HOLIDAY) Execute(now time.Time) {
	// Check once a day.
	if dw(h.lastActive) == dw(h.DATE_IN) && !time.Time(h.lastActive).IsZero() {
		return
	}
	h.lastActive = h.DATE_IN
	lx := LANGUAGE.DEFAULT
	if h.LANGU != 0 {
		lx = min(LANGUAGE.LMAX, h.LANGU)
	}
	jahr := YEAR_OF_DATE(h.DATE_IN)
	ostern := EASTER(jahr)
	wdx := DAY_OF_WEEK(h.DATE_IN)
	in := dw(h.DATE_IN)
	h.Y = false
	if h.HOLIDAYS != nil {
		for _, hd := range h.HOLIDAYS {
			xDate := SET_DATE(jahr, iec.INT(hd.MONTH), iec.INT(hd.DAY))
			found := false
			switch {
			case hd.USE == 1 && hd.MONTH > 0:
				found = dw(xDate) == in
			case hd.USE == 1 && hd.MONTH == 0:
				found = dw(DATE_ADD(ostern, iec.INT(hd.DAY), 0, 0, 0)) == in
			case hd.USE < 0:
				found = wdx == -iec.INT(hd.USE) && in < dw(xDate) && in >= dw(DATE_ADD(xDate, -7, 0, 0, 0))
			}
			if found {
				h.Y = true
				h.NAME = hd.NAME
				return
			}
		}
	}
	if wdx == 5 && h.FRIDAY || wdx == 6 && h.SATURDAY || wdx == 7 && h.SUNDAY {
		h.Y = true
		// OSCAT picks the language of the location numbered LANGU.
		lang := LOCATION.LANGUAGE[LIMIT(1, lx, 5)-1]
		h.NAME = LANGUAGE.WEEKDAYS[LIMIT(1, lang, 3)-1][wdx-1]
	} else {
		h.NAME = ""
	}
}

// REFRACTION calculates the atmospheric refraction in degrees for a sun
// elevation elev in degrees, 0 at the horizon.
func REFRACTION(elev iec.REAL) iec.REAL {
	elev = LIMIT(-1.9, elev, 80.0)
	return 0.0174532925199433 / TAN(0.0174532925199433*(elev+10.3/(elev+5.11)))
}

// RTC_2 is a real time clock of world time UDT, set to SDT and SMS
// milliseconds by SET, with the local time LOCAL_DT OFS minutes away and, if
// DEN is true, one hour more during daylight saving time, which DSO shows.
type RTC_2 struct {
	SET      iec.BOOL
	SDT      iec.DT
	SMS      iec.INT
	DEN      iec.BOOL
	OFS      iec.INT
	UDT      iec.DT
	LOCAL_DT iec.DT
	DSO      iec.BOOL
	XMS      iec.INT

	rt RTC_MS
}

// INIT resets the block.
func (r *RTC_2) INIT() { *r = RTC_2{} }

// Execute runs the block once.
func (r *RTC_2) Execute(now time.Time) {
	r.rt.SET = r.SET
	r.rt.SDT = r.SDT
	r.rt.SMS = r.SMS
	r.rt.Execute(now)
	r.UDT = r.rt.XDT
	r.XMS = r.rt.XMS
	r.DSO = DST(r.UDT) && r.DEN
	r.LOCAL_DT = DWORD_TO_DT(dtw(r.UDT) + intToDword(r.OFS+BOOL_TO_INT(r.DSO)*60)*60)
}

// RTC_MS is a real time clock with milliseconds, XDT and XMS, set to SDT and
// SMS by SET or when it first runs.
type RTC_MS struct {
	SET iec.BOOL
	SDT iec.DT
	SMS iec.INT
	XDT iec.DT
	XMS iec.INT

	init iec.BOOL
	last iec.DWORD
}

// INIT resets the block.
func (r *RTC_MS) INIT() { *r = RTC_MS{} }

// Execute runs the block once.
func (r *RTC_MS) Execute(now time.Time) {
	tx := PLC_MS(now)
	if r.SET || !r.init {
		r.init = true
		r.XDT = r.SDT
		r.XMS = r.SMS
	} else {
		r.XMS += iec.INT(tx - r.last)
		if r.XMS > 999 {
			r.XDT = DWORD_TO_DT(dtw(r.XDT) + 1)
			r.XMS -= 1000
		}
	}
	r.last = tx
}

// SUN_MIDDAY calculates the world time when the sun stands south of the
// longitude lon on the date utc.
func SUN_MIDDAY(lon iec.REAL, utc iec.DATE) iec.TOD {
	t := iec.REAL(DAY_OF_YEAR(utc))
	offset := -0.1752*SIN(0.033430*t+0.5474) - 0.1340*SIN(0.018234*t-0.1939)
	return HOUR_TO_TOD(12.0 - offset - lon*0.0666666666666)
}

// SUN_POS calculates the position of the sun at a location and world time:
// B is the angle from north, H the height above the horizon and HR the
// height corrected for refraction, all in degrees.
type SUN_POS struct {
	LATITUDE  iec.REAL
	LONGITUDE iec.REAL
	UTC       iec.DT
	B         iec.REAL
	H         iec.REAL
	HR        iec.REAL
}

// INIT resets the block.
func (s *SUN_POS) INIT() { *s = SUN_POS{} }

// Execute runs the block once.
func (s *SUN_POS) Execute(now time.Time) {
	pi, pi2 := MATH.PI, MATH.PI2
	// n is the days since 2000-01-01 12:00.
	n := iec.REAL(dtw(s.UTC)-946728000) * 0.000011574074074074
	g := math.MODR(6.240040768+0.01720197*n, pi2)
	d := math.MODR(4.89495042+0.017202792*n, pi2) + 0.033423055*SIN(g) + 0.000349066*SIN(2.0*g)
	e := 0.409087723 - 0.000000006981317008*n
	cosD := COS(d)
	sinD := SIN(d)
	a := ATAN(COS(e) * sinD / cosD)
	if cosD < 0.0 {
		a = a + pi
	}
	c := ASIN(SIN(e) * sinD)
	tau := math.RAD(math.MODR(6.697376+(n-0.25)*0.0657098245037645+iec.REAL(tw(DT_TO_TOD(s.UTC)))*0.0000002785383333, 24.0)*15.0+s.LONGITUDE) - a
	rlat := math.RAD(s.LATITUDE)
	sinLat := SIN(rlat)
	cosLat := COS(rlat)
	cosTau := COS(tau)
	t1 := cosTau*sinLat - TAN(c)*cosLat
	s.B = ATAN(SIN(tau) / t1)
	if t1 < 0.0 {
		s.B = s.B + pi2
	} else {
		s.B = s.B + pi
	}
	s.B = math.DEG(math.MODR(s.B, pi2))
	s.H = math.DEG(ASIN(COS(c)*cosTau*cosLat + SIN(c)*sinLat))
	if s.H > 180.0 {
		s.H = s.H - 360.0
	}
	s.HR = s.H + REFRACTION(s.H)
}

// SUN_TIME calculates, in world time, the sun rise, sun set and midday at a
// location on the date UTC, and the sun's height at midday in degrees. H is
// the sun's height at sun rise and sun set. The result is accurate to a few
// minutes.
type SUN_TIME struct {
	LATITUDE        iec.REAL
	LONGITUDE       iec.REAL
	UTC             iec.DATE
	H               iec.REAL // default -0.83333333333
	MIDDAY          iec.TOD
	SUN_RISE        iec.TOD
	SUN_SET         iec.TOD
	SUN_DECLINATION iec.REAL
}

// INIT resets the block and sets H to its initial value.
func (s *SUN_TIME) INIT() { *s = SUN_TIME{H: -0.83333333333} }

// Execute runs the block once.
func (s *SUN_TIME) Execute(now time.Time) {
	b := s.LATITUDE * 0.0174532925199433
	s.MIDDAY = SUN_MIDDAY(s.LONGITUDE, s.UTC)
	dk := 0.40954 * SIN(0.0172*(iec.REAL(DAY_OF_YEAR(s.UTC))-79.35))
	s.SUN_DECLINATION = math.DEG(dk)
	if s.SUN_DECLINATION > 180.0 {
		s.SUN_DECLINATION = s.SUN_DECLINATION - 360.0
	}
	s.SUN_DECLINATION = 90.0 - s.LATITUDE + s.SUN_DECLINATION
	delta := TIME_TO_DWORD(HOUR_TO_TIME(ACOS((SIN(math.RAD(s.H))-SIN(b)*SIN(dk))/(COS(b)*COS(dk))) * 3.819718632))
	s.SUN_RISE = DWORD_TO_TOD(tw(s.MIDDAY) - delta)
	s.SUN_SET = DWORD_TO_TOD(tw(s.MIDDAY) + delta)
}
