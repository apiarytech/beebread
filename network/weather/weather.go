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

// Package weather is the port of the OSCAT NETWORK weather services:
// WORLD_WEATHER reads the forecast of worldweatheronline.com and
// YAHOO_WEATHER the one of Yahoo, with their descriptions in German and
// OSCAT's icons, and MOON_PHASE computes the phase of the moon. The
// services have closed their old interfaces since the library was written.
package weather

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/logic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/beebread/network/inet"
	"github.com/apiarytech/beebread/network/parser"
	"github.com/apiarytech/royaljelly/iec"
)

// service is what the weather services share: the connection, a
// DNS_CLIENT and an HTTP_GET of the URL urlData.
type service struct {
	IP_C  *network.IP_C
	S_BUF *network.NETWORK_BUFFER
	R_BUF *network.NETWORK_BUFFER

	urlData     network.URL
	dns         inet.DNS_CLIENT
	http        inet.HTTP_GET
	state       iec.INT
	initialized bool
}

func (s *service) ok() bool { return s.IP_C != nil && s.S_BUF != nil && s.R_BUF != nil }

// start sets the HTTP_GET up, the first time.
func (s *service) start() {
	if !s.initialized {
		s.initialized = true
		s.http.INIT()
	}
}

// run runs the DNS_CLIENT and HTTP_GET: the query while the state is 40,
// the get while it is 60, and the release at 100.
func (s *service) run(now time.Time) {
	s.dns.IP_C, s.dns.S_BUF, s.dns.R_BUF = s.IP_C, s.S_BUF, s.R_BUF
	s.dns.DOMAIN, s.dns.IP4_DNS, s.dns.ACTIVATE = s.urlData.DOMAIN, 0, s.state == 40
	s.dns.Execute(now)
	s.http.IP_C, s.http.S_BUF, s.http.R_BUF = s.IP_C, s.S_BUF, s.R_BUF
	s.http.IP4, s.http.GET, s.http.MODE, s.http.UNLOCK_BUF, s.http.URL_DATA = s.dns.IP4, s.state == 60, 2, s.state == 100, &s.urlData
	s.http.Execute(now)
}

// number reads value as a REAL and an INT, 0 if it is not a number of 20
// characters at most.
func number(value iec.STRING) (iec.REAL, iec.INT) {
	if LEN(value) > 20 {
		return 0, 0
	}
	v := str.FLOAT_TO_REAL(value)
	if logic.CHK_REAL(v) != 0 {
		return 0, 0
	}
	return v, REAL_TO_INT(v)
}

// WORLD_WEATHER reads the weather WW at LATITUDE, LONGITUDE from the CSV
// of worldweatheronline.com, with the key KEY, on a rising edge of
// ACTIVATE: DONE, BUSY while it works, else ERROR_C and ERROR_T: 1 DNS, 2
// HTTP.
type WORLD_WEATHER struct {
	service
	WW        *network.WORLD_WEATHER_DATA
	ACTIVATE  iec.BOOL
	LATITUDE  iec.REAL
	LONGITUDE iec.REAL
	KEY       iec.STRING // STRING(30)
	BUSY      iec.BOOL
	DONE      iec.BOOL
	ERROR_C   iec.DWORD
	ERROR_T   iec.BYTE

	cpb       parser.CSV_PARSER_BUF
	lastState iec.BOOL
	urlStr    iec.STRING // STRING(STRING_LENGTH)
	offset    iec.UDINT
	sep       iec.BYTE
	cnt       iec.INT
	day       iec.INT
	idx       iec.INT
	runCSV    iec.BYTE
	result    iec.BYTE
	value     iec.STRING // STRING(STRING_LENGTH)
}

// INIT resets the block.
func (w *WORLD_WEATHER) INIT() {
	*w = WORLD_WEATHER{service: service{IP_C: w.IP_C, S_BUF: w.S_BUF, R_BUF: w.R_BUF}, WW: w.WW}
	w.start()
}

// icon reads the icon of the URL of an icon, '.../wsymbol_0001_...'.
func icon(value iec.STRING) iec.INT {
	pos := FIND(value, "/wsymbol_")
	return str.DEC_TO_INT(MID(value, 4, pos+9))
}

// Execute runs the block once.
func (w *WORLD_WEATHER) Execute(now time.Time) {
	w.start()
	if !w.ok() || w.WW == nil {
		return
	}
	ww := w.WW
	switch w.state {
	case 0:
		if w.ACTIVATE && !w.lastState {
			w.state = 20
			w.DONE = false
			w.BUSY = true
			w.ERROR_C, w.ERROR_T = 0, 0
		}
	case 20:
		w.urlStr = CONCAT("http://api.worldweatheronline.com/free/v1/weather.ashx?q=", str.REAL_TO_STRF(w.LATITUDE, 6, "."),
			",", str.REAL_TO_STRF(w.LONGITUDE, 6, "."), "&format=csv&num_of_days=5&key=", w.KEY)
		w.urlData = encoding.STRING_TO_URL(w.urlStr, "", "")
		w.state = 40
	case 40:
		if w.dns.DONE {
			w.state = 60
		} else if w.dns.ERROR > 0 {
			w.ERROR_C, w.ERROR_T = w.dns.ERROR, 1
			w.state = 100
		}
	case 60:
		if w.http.DONE {
			w.sep, w.idx, w.day, w.cnt = 0, 0, 0, 0
			w.offset = iec.UDINT(w.http.BODY_START)
			w.runCSV = 1
			w.state = 80
		} else if w.http.ERROR > 0 {
			w.ERROR_C, w.ERROR_T = w.http.ERROR, 2
			w.state = 100
		}
	case 80:
		if w.runCSV != 0 {
			break
		}
		w.runCSV = 1 // the next element
		switch w.result {
		case 1, 2:
			vReal, vInt := number(w.value)
			w.cnt++
			// The comment lines '#...' first.
			if w.cnt == 1 && LEFT(w.value, 1) != "#" {
				w.cnt = 0
			}
			if w.cnt < 8 || w.cnt > 87 {
				break
			}
			if w.cnt > 22 {
				// The days.
				d := &ww.DAY[w.day]
				switch w.idx {
				case 0:
					d.DATE_OF_DAY = network.STRING_N(w.value, 10)
				case 1:
					d.TEMP_MAX_C = vInt
				case 2:
					d.TEMP_MAX_F = vInt
				case 3:
					d.TEMP_MIN_C = vInt
				case 4:
					d.TEMP_MIN_F = vInt
				case 5:
					d.WIND_SPEED_MILES = vInt
				case 6:
					d.WIND_SPEED_KMPH = vInt
				case 7:
					d.WIND_DIR_DEGREE = vInt
				case 8:
					d.WIND_DIR16POINT = network.STRING_N(w.value, 3)
				case 9:
					d.WEATHER_CODE = vInt
				case 10:
					d.WEATHER_ICON = icon(w.value)
				case 11:
					d.WEATHER_DESC = network.STRING_N(w.value, 60)
				case 12:
					d.PRECIPMM = vReal
				}
				w.idx++
				if w.idx > 12 {
					w.idx = 0
					w.day++
				}
				break
			}
			// The current weather.
			cur := &ww.CUR
			switch w.cnt {
			case 8:
				w.sep = 44 // the elements are separated by commas from here
			case 9:
				cur.OBSERVATION_TIME = network.STRING_N(w.value, 8)
			case 10:
				cur.TEMP_C = vInt
			case 11:
				cur.WEATHER_CODE = vInt
			case 12:
				cur.WEATHER_ICON = icon(w.value)
			case 13:
				cur.WEATHER_DESC = network.STRING_N(w.value, 60)
			case 14:
				cur.WIND_SPEED_MILES = vInt
			case 15:
				cur.WIND_SPEED_KMPH = vInt
			case 16:
				cur.WIND_DIR_DEGREE = vInt
			case 17:
				cur.WIND_DIR16POINT = network.STRING_N(w.value, 3)
			case 18:
				cur.PRECIPMM = vReal
			case 19:
				cur.HUMIDITY = vInt
			case 20:
				cur.VISIBILITY = vInt
			case 21:
				cur.PRESSURE = vInt
			case 22:
				cur.CLOUDOVER = vInt
			}
		case 10: // the end of the data
			w.DONE = true
			w.runCSV = 0
			w.state = 100
		}
	case 100:
		if !w.http.DONE {
			w.state = 0
			w.BUSY = false
			w.DONE = w.ERROR_T == 0
		}
	}
	w.run(now)
	w.cpb.SEP, w.cpb.RUN, w.cpb.OFFSET, w.cpb.VALUE, w.cpb.PT = w.sep, &w.runCSV, &w.offset, &w.value, w.R_BUF
	w.cpb.Execute(now)
	w.result = w.cpb.RESULT
	w.lastState = w.ACTIVATE
}

// worldIcons are OSCAT's icons of the icons of worldweatheronline.com.
var worldIcons = [42]iec.INT{0, 14, 12, 9, 9, 0, 7, 7, 0, 2, 2, 3, 3, 1, 0, 0, 17, 4, 4, 3, 3, 1, 0, 0, 17, 2, 2, 3, 3, 1, 0, 0, 17, 2, 4, 3, 3, 1, 0, 0, 17, 0}

// WORLD_WEATHER_ICON_OSCAT turns the icons of WW into OSCAT's icons, on a
// rising edge of ACTIVATE.
type WORLD_WEATHER_ICON_OSCAT struct {
	WW       *network.WORLD_WEATHER_DATA
	ACTIVATE iec.BOOL

	actLast iec.BOOL
}

// INIT resets the block.
func (w *WORLD_WEATHER_ICON_OSCAT) INIT() { *w = WORLD_WEATHER_ICON_OSCAT{WW: w.WW} }

// Execute runs the block once.
func (w *WORLD_WEATHER_ICON_OSCAT) Execute(now time.Time) {
	if w.WW == nil {
		return
	}
	if w.ACTIVATE && !w.actLast {
		w.WW.CUR.WEATHER_ICON = worldIcons[LIMIT(0, w.WW.CUR.WEATHER_ICON, 41)]
		for i := range w.WW.DAY {
			w.WW.DAY[i].WEATHER_ICON = worldIcons[LIMIT(0, w.WW.DAY[i].WEATHER_ICON, 41)]
		}
	}
	w.actLast = w.ACTIVATE
}

// notAvailable is the description of a code the tables do not have.
const notAvailable iec.STRING = "nicht verfügbar"

// describe returns the description of code in table.
func describe(table map[iec.INT]iec.STRING, code iec.INT) iec.STRING {
	if d, ok := table[code]; ok {
		return d
	}
	return notAvailable
}

// worldDesc are the German descriptions of the weather codes of
// worldweatheronline.com.
var worldDesc = map[iec.INT]iec.STRING{
	113: "heiter/sonnig",
	116: "stellenweise wolkig",
	119: "wolkig",
	122: "bewölkt",
	143: "Nebel",
	176: "stellenweise nahegelegener Regen",
	179: "stellenweise nahegelegener Schneefall",
	182: "stellenweise nahegelegener Eisregen",
	185: "stellenweise gefrierender Niesel",
	200: "nahegelegene Sturm",
	227: "Schneetreiben",
	230: "Schneesturm",
	248: "Nebel",
	260: "gefrierender Nebel",
	263: "stellenweise leichter Nieselregen",
	266: "leichter Niesel",
	281: "gefrierender Nieselregen",
	284: "stark gefrierender Nieselregen",
	293: "stellenweise leichter Regen",
	296: "leichter Regen",
	299: "zeitweise etwas Regen",
	302: "mäßiger Regen",
	305: "zeitweise etwas starker Regen",
	308: "starker Regen",
	311: "leichter gefrierender Regen",
	314: "mäßiger oder starker Eisregen",
	317: "leichter Eisregen",
	320: "mäßiger oder starker Eisregen",
	323: "stellenweise leichter Schneefall",
	326: "leichter Scheefall",
	329: "stellenweise moderater Schneefall",
	332: "mäßiger Schnee",
	335: "stellenweise starker Scheefall",
	338: "starker Schneefall",
	350: "Hagel",
	353: "leichte Regenfälle",
	356: "mäßiger oder starker Regenfall",
	359: "sintflutartige Regenfälle",
	362: "leichter Eisregen",
	365: "mäßiger oder starker Eisregen",
	368: "leichte Schneeschauer",
	371: "mäßiger oder starker Schneefall",
	374: "leichter Hagel",
	377: "mäßiger oder starker Hagel",
	386: "stellenweise leichter Regen in Gewittergebiet",
	389: "mäßiger oder starker Regen in Gewittergebiet",
	392: "stellenweise leichter Schneefall in Gewittergebiet",
	395: "mäßiger oder starker Schnee in Gewittergebiet",
}

// WORLD_WEATHER_DESC_GE sets the descriptions of WW in German, on a rising
// edge of ACTIVATE.
type WORLD_WEATHER_DESC_GE struct {
	WW       *network.WORLD_WEATHER_DATA
	ACTIVATE iec.BOOL

	actLast iec.BOOL
}

// INIT resets the block.
func (w *WORLD_WEATHER_DESC_GE) INIT() { *w = WORLD_WEATHER_DESC_GE{WW: w.WW} }

// Execute runs the block once.
func (w *WORLD_WEATHER_DESC_GE) Execute(now time.Time) {
	if w.WW == nil {
		return
	}
	if w.ACTIVATE && !w.actLast {
		for i := range w.WW.DAY {
			w.WW.DAY[i].WEATHER_DESC = describe(worldDesc, w.WW.DAY[i].WEATHER_CODE)
		}
		w.WW.CUR.WEATHER_DESC = describe(worldDesc, w.WW.CUR.WEATHER_CODE)
	}
	w.actLast = w.ACTIVATE
}

// YAHOO_WEATHER reads the weather YW at the place LOCATION, a WOEID of
// Yahoo, in Fahrenheit (UNITS) or Celsius, on a rising edge of ACTIVATE,
// in two queries: DONE, BUSY while it works, else ERROR_C and ERROR_T: 1
// DNS, 2 HTTP.
type YAHOO_WEATHER struct {
	service
	YW       *network.YAHOO_WEATHER_DATA
	ACTIVATE iec.BOOL
	UNITS    iec.BOOL
	LOCATION iec.STRING // STRING(20)
	BUSY     iec.BOOL
	DONE     iec.BOOL
	ERROR_C  iec.DWORD
	ERROR_T  iec.BYTE

	ctrl         network.XML_CONTROL
	xml          parser.XML_READER
	activateLast iec.BOOL
	cycle        iec.INT
}

// INIT resets the block.
func (y *YAHOO_WEATHER) INIT() {
	*y = YAHOO_WEATHER{service: service{IP_C: y.IP_C, S_BUF: y.S_BUF, R_BUF: y.R_BUF}, YW: y.YW}
	y.start()
}

// Execute runs the block once.
func (y *YAHOO_WEATHER) Execute(now time.Time) {
	y.start()
	if !y.ok() || y.YW == nil {
		return
	}
	yw := y.YW
	switch y.state {
	case 0:
		if y.ACTIVATE && !y.activateLast {
			y.state = 40
			y.DONE = false
			y.BUSY = true
			y.ERROR_C, y.ERROR_T = 0, 0
			y.cycle = 0
			y.urlData.DOMAIN = "query.yahooapis.com"
		}
	case 40:
		if y.dns.DONE {
			y.state = 50
		} else if y.dns.ERROR > 0 {
			y.ERROR_C, y.ERROR_T = y.dns.ERROR, 1
			y.state = 100
		}
	case 50: // the query: the information, then the days
		if y.cycle == 0 {
			y.urlData.QUERY = "q=select%20units,wind,atmosphere,astronomy,location%20from%20weather.forecast%20where%20woeid="
		} else {
			y.urlData.QUERY = "q=select%20item%20from%20weather.forecast%20where%20woeid="
		}
		y.urlData.QUERY = CONCAT(y.urlData.QUERY, y.LOCATION, SEL[iec.STRING](y.UNITS, "%20and%20u=%27c%27", "%20and%20u=%27f%27"))
		y.urlData.PATH = "/v1/public/yql"
		y.state = 60
	case 60:
		if y.http.DONE {
			y.state = 80
			y.ctrl.START_POS, y.ctrl.STOP_POS = y.http.BODY_START, y.http.BODY_STOP
			y.ctrl.COMMAND = 0x8018 // the text and the attributes
			y.ctrl.WATCHDOG = iec.TIME(time.Millisecond)
		} else if y.http.ERROR > 0 {
			y.ERROR_C, y.ERROR_T = y.http.ERROR, 2
			y.state = 100
		}
	case 80:
		y.xml.CTRL, y.xml.BUF = &y.ctrl, &y.R_BUF.BUFFER
		y.xml.Execute(now)
		if y.ctrl.TYP < 98 {
			vReal, vInt := number(y.ctrl.VALUE)
			v := y.ctrl.VALUE
			count := iec.INT(y.ctrl.COUNT)
			if y.cycle == 0 {
				switch count {
				case 14:
					yw.LOCATION_CITY = network.STRING_N(v, 40)
				case 15:
					yw.LOCATION_COUNTRY = network.STRING_N(v, 20)
				case 16:
					yw.LOCATION_REGION = network.STRING_N(v, 20)
				case 20:
					yw.WIND_CHILL = vInt
				case 21:
					yw.WIND_DIRECTION = vInt
				case 22:
					yw.WIND_SPEED = vReal
				case 26:
					yw.ATMOSPHERE_HUMIDITY = vInt
				case 27:
					yw.ATMOSPHERE_PRESSURE = vReal
				case 28:
					yw.ATMOSPHERE_RISING = vInt
				case 29:
					yw.ATMOSPHERE_VISIBILITY = vReal
				case 33:
					yw.ASTRONOMY_SUNRISE = network.STRING_N(v, 10)
				case 34:
					yw.ASTRONOMY_SUNSET = network.STRING_N(v, 10)
				case 38:
					yw.UNIT_DISTANCE = network.STRING_N(v, 2)
				case 39:
					yw.UNIT_PRESSURE = network.STRING_N(v, 2)
				case 40:
					yw.UNIT_SPEED = network.STRING_N(v, 4)
				case 41:
					yw.UNIT_TEMPERATURE = network.STRING_N(v, 1)
				}
			} else {
				switch {
				case count == 18:
					yw.GEO_LATITUDE = vReal
				case count == 22:
					yw.GEO_LONGITUDE = vReal
				case count == 32:
					yw.CUR_CONDITIONS_CODE, yw.CUR_CONDITIONS_ICON = vInt, vInt
				case count == 34:
					yw.CUR_CONDITIONS_TEMP = vInt
				case count == 35:
					yw.CUR_CONDITIONS_TEXT = network.STRING_N(v, 40)
				case count == 39:
					yw.FORECAST_TODAY_CODE, yw.FORECAST_TODAY_ICON = vInt, vInt
				case count == 40:
					yw.FORECAST_TODAY_DATE_LONG = network.STRING_N(v, 20)
				case count == 42:
					yw.FORECAST_TODAY_HIGH_TEMP = vInt
				case count == 43:
					yw.FORECAST_TODAY_LOW_TEMP = vInt
				case count == 44:
					yw.FORECAST_TODAY_TEXT = network.STRING_N(v, 40)
				case count >= 48 && count <= 125: // the days 1 to 9
					cnt := count - 48
					d := &yw.FORECAST_DAY[cnt/9] // FORECAST_DAY[cnt / 9 + 1]
					switch cnt % 9 {
					case 0:
						d.CODE, d.ICON = vInt, vInt
					case 1:
						d.DATE_LONG = network.STRING_N(v, 20)
					case 3:
						d.HIGH_TEMP = vInt
					case 4:
						d.LOW_TEMP = vInt
					case 5:
						d.TEXT = network.STRING_N(v, 40)
					}
				}
			}
		} else if y.ctrl.TYP == 99 { // the end
			y.state = 100
		}
	case 100:
		if !y.http.DONE {
			if y.cycle == 0 && y.ERROR_T == 0 {
				// the second query
				y.state = 50
				y.cycle++
			} else {
				y.state = 0
				y.BUSY = false
				y.DONE = y.ERROR_T == 0
			}
		}
	}
	y.run(now)
	y.activateLast = y.ACTIVATE
}

// yahooIcons are OSCAT's icons of the weather codes of Yahoo.
var yahooIcons = [49]iec.INT{4, 2, 2, 4, 17, 1, 3, 1, 2, 2, 2, 2, 2, 3, 3, 3, 3, 5, 3, 6, 7, 14, 8, 9, 9, 10, 9, 11, 12, 11, 12, 13, 14, 15, 16, 4, 14, 4, 17, 2, 2, 3, 3, 3, 12, 2, 3, 4, 0}

// YAHOO_WEATHER_ICON_OSCAT sets the icons of YW to OSCAT's icons of their
// codes, on a rising edge of ACTIVATE.
type YAHOO_WEATHER_ICON_OSCAT struct {
	YW       *network.YAHOO_WEATHER_DATA
	ACTIVATE iec.BOOL

	actLast iec.BOOL
}

// INIT resets the block.
func (y *YAHOO_WEATHER_ICON_OSCAT) INIT() { *y = YAHOO_WEATHER_ICON_OSCAT{YW: y.YW} }

// yahooIcon returns OSCAT's icon of the code c.
func yahooIcon(c iec.INT) iec.INT { return yahooIcons[max(0, min(c, 48))] }

// Execute runs the block once.
func (y *YAHOO_WEATHER_ICON_OSCAT) Execute(now time.Time) {
	if y.YW == nil {
		return
	}
	yw := y.YW
	if y.ACTIVATE && !y.actLast {
		yw.CUR_CONDITIONS_ICON = yahooIcon(yw.CUR_CONDITIONS_CODE)
		yw.FORECAST_TODAY_ICON = yahooIcon(yw.FORECAST_TODAY_CODE)
		for i := range yw.FORECAST_DAY {
			yw.FORECAST_DAY[i].ICON = yahooIcon(yw.FORECAST_DAY[i].CODE)
		}
	}
	y.actLast = y.ACTIVATE
}

// yahooDesc are the German descriptions of the weather codes of Yahoo.
var yahooDesc = map[iec.INT]iec.STRING{
	0:  "Wirbelsturm",
	1:  "Tropensturm",
	2:  "Wirbelsturm",
	3:  "starkes Gewitter ",
	4:  "Gewitter",
	5:  "Schneeregen",
	6:  "Schneeregen",
	7:  "Schneeregen",
	8:  "gefrierender Nieselregen",
	9:  "Nieselregen",
	10: "gefrierender Regen",
	11: "Regenfälle",
	12: "Regenfälle",
	13: "Schneegestöber",
	14: "leichte Schneefälle",
	15: "Schneetreiben",
	16: "Schnee",
	17: "Hagel",
	18: "Schneeregen",
	19: "Staub",
	20: "nebelig",
	21: "Dunst",
	22: "rauchig",
	23: "stürmisch",
	24: "windig",
	25: "kalt",
	26: "bewölkt",
	27: "meist bewölkt",
	28: "meist bewölkt",
	29: "teilweise bewölkt",
	30: "teilweise bewölkt",
	31: "klar",
	32: "sonnig",
	33: "heiter",
	34: "heiter",
	35: "Regen und Hagel",
	36: "heiß",
	37: "örtliche Gewitter",
	38: "vereinzelte Gewitter",
	39: "vereinzelte Gewitter",
	40: "vereinzelte Regenfälle",
	41: "schwere Schneefälle",
	42: "vereinzelte Schneefälle",
	43: "schwere Schneefälle",
	44: "teilweise bewölkt",
	45: "gewittrige Regenfälle",
	46: "Schneefälle",
	47: "örtliche gewittrige Regenfälle",
}

// YAHOO_WEATHER_DESC_GE sets the descriptions of YW in German, on a rising
// edge of ACTIVATE. OSCAT has the case 7 twice, and the first one leaves
// out the code of day 5: day 5 has the description of day 4.
type YAHOO_WEATHER_DESC_GE struct {
	YW       *network.YAHOO_WEATHER_DATA
	ACTIVATE iec.BOOL

	actLast iec.BOOL
}

// INIT resets the block.
func (y *YAHOO_WEATHER_DESC_GE) INIT() { *y = YAHOO_WEATHER_DESC_GE{YW: y.YW} }

// Execute runs the block once.
func (y *YAHOO_WEATHER_DESC_GE) Execute(now time.Time) {
	if y.YW == nil {
		return
	}
	yw := y.YW
	if y.ACTIVATE && !y.actLast {
		var x iec.INT
		var s iec.STRING
		for i := 1; i <= 12; i++ {
			switch i {
			case 1:
				x = yw.CUR_CONDITIONS_CODE
			case 2:
				yw.CUR_CONDITIONS_TEXT = s
				x = yw.FORECAST_TODAY_CODE
			case 3:
				yw.FORECAST_TODAY_TEXT = s
				x = yw.FORECAST_DAY[0].CODE
			case 7:
				// OSCAT's first case 7: the code of day 5 is not read.
				yw.FORECAST_DAY[3].TEXT = s
			default: // 4..6, 8..12
				yw.FORECAST_DAY[i-4].TEXT = s
				if i < 12 {
					x = yw.FORECAST_DAY[i-3].CODE
				}
			}
			s = network.STRING_N(describe(yahooDesc, x), 40)
		}
	}
	y.actLast = y.ACTIVATE
}

// MOON_PHASE computes the phase of the moon PHASE at the time XDT, from 0,
// the new moon, to SCALE, every UPDATE.
type MOON_PHASE struct {
	XDT    iec.DT
	SCALE  iec.BYTE // default 12
	UPDATE iec.TIME // default T#1h
	PHASE  iec.BYTE

	x      iec.UDINT
	lastDT iec.DT
}

// INIT resets the block and sets its inputs to their initial values.
func (m *MOON_PHASE) INIT() { *m = MOON_PHASE{SCALE: 12, UPDATE: iec.TIME(time.Hour)} }

// Execute runs the block once.
func (m *MOON_PHASE) Execute(now time.Time) {
	// DT - DT is a TIME, in milliseconds modulo 2^32.
	if (DT_TO_DWORD(m.XDT)-DT_TO_DWORD(m.lastDT))*1000 > TIME_TO_DWORD(m.UPDATE) {
		m.x = (iec.UDINT(DT_TO_DWORD(m.XDT)) - 603240) % 2551392
		m.PHASE = iec.BYTE(iec.UDINT(m.SCALE) * m.x / 2551392)
		m.lastDT = m.XDT
	}
}
