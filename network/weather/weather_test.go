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

package weather

import (
	"testing"
	"time"

	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

func TestWorldTables(t *testing.T) {
	var ww network.WORLD_WEATHER_DATA
	ww.CUR.WEATHER_CODE, ww.CUR.WEATHER_ICON = 113, 1
	ww.DAY[2].WEATHER_CODE, ww.DAY[2].WEATHER_ICON = 999, 50
	d := WORLD_WEATHER_DESC_GE{WW: &ww, ACTIVATE: true}
	d.Execute(time.Time{})
	if ww.CUR.WEATHER_DESC != "heiter/sonnig" || ww.DAY[2].WEATHER_DESC != "nicht verfügbar" {
		t.Errorf("descriptions %q %q", ww.CUR.WEATHER_DESC, ww.DAY[2].WEATHER_DESC)
	}
	i := WORLD_WEATHER_ICON_OSCAT{WW: &ww, ACTIVATE: true}
	i.Execute(time.Time{})
	if ww.CUR.WEATHER_ICON != 14 || ww.DAY[2].WEATHER_ICON != 0 {
		t.Errorf("icons %d %d", ww.CUR.WEATHER_ICON, ww.DAY[2].WEATHER_ICON)
	}
}

func TestYahooTables(t *testing.T) {
	var yw network.YAHOO_WEATHER_DATA
	yw.CUR_CONDITIONS_CODE = 0
	for i := range yw.FORECAST_DAY {
		yw.FORECAST_DAY[i].CODE = iec.INT(i + 1)
	}
	d := YAHOO_WEATHER_DESC_GE{YW: &yw, ACTIVATE: true}
	d.Execute(time.Time{})
	if yw.CUR_CONDITIONS_TEXT != describe(yahooDesc, 0) || yw.FORECAST_DAY[0].TEXT != describe(yahooDesc, 1) {
		t.Errorf("texts %q %q", yw.CUR_CONDITIONS_TEXT, yw.FORECAST_DAY[0].TEXT)
	}
	// OSCAT's double case 7: day 5 has the text of day 4.
	if yw.FORECAST_DAY[4].TEXT != yw.FORECAST_DAY[3].TEXT || yw.FORECAST_DAY[8].TEXT != describe(yahooDesc, 9) {
		t.Errorf("days 4, 5, 9: %q %q %q", yw.FORECAST_DAY[3].TEXT, yw.FORECAST_DAY[4].TEXT, yw.FORECAST_DAY[8].TEXT)
	}
	yw.CUR_CONDITIONS_CODE = 200 // past the table
	i := YAHOO_WEATHER_ICON_OSCAT{YW: &yw, ACTIVATE: true}
	i.Execute(time.Time{})
	if yw.CUR_CONDITIONS_ICON != 0 || yw.FORECAST_DAY[0].ICON != 2 {
		t.Errorf("icons %d %d", yw.CUR_CONDITIONS_ICON, yw.FORECAST_DAY[0].ICON)
	}
}

func TestMoonPhase(t *testing.T) {
	var m MOON_PHASE
	m.INIT()
	// A new moon, 2024-01-11 11:57 UTC, and a full moon, 2024-01-25 17:54.
	m.XDT = td.SET_DT(2024, 1, 11, 12, 0, 0)
	m.Execute(time.Time{})
	if m.PHASE != 0 && m.PHASE != 11 {
		t.Errorf("new moon: %d", m.PHASE)
	}
	m.XDT = td.SET_DT(2024, 1, 25, 18, 0, 0)
	m.Execute(time.Time{})
	if m.PHASE != 6 && m.PHASE != 5 {
		t.Errorf("full moon: %d", m.PHASE)
	}
}
