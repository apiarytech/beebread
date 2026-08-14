/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package engineering

import (
	"math"
	"time"

	. "beebread/basic"
)

// Astro converts between astronomical length units.
type Astro struct {
	Ym  float64 // Meters
	YAE float64 // Astronomical Units
	YPC float64 // Parsecs
	YLJ float64 // Lightyears
}

// Update executes the conversion logic.
func (a *Astro) Update(m, ae, pc, lj float64) {
	a.YAE = ae + m*6.6845871535e-12 + pc*206265.0 + lj*63240.0
	a.Ym = a.YAE * 149.597870e9
	a.YPC = a.YAE * 4.8481322570e-6
	a.YLJ = a.YAE * 1.5812776724e-5
}

// BftToMs converts wind speed from Beaufort to m/s.
func BftToMs(bft int) float64 {
	return math.Pow(float64(bft), 1.5) * 0.836
}

// CToF converts Celsius to Fahrenheit.
func CToF(celsius float64) float64 {
	return celsius*1.8 + 32.0
}

// CToK converts Celsius to Kelvin.
func CToK(celsius float64) float64 {
	return celsius - float64(Phys.T0)
}

// DegToDir converts degrees to a compass direction string.
// L is language, N is precision (1=4 dirs, 2=8 dirs, 3=16 dirs).
func DegToDir(deg, n, l int) string {
	if l == 0 {
		l = int(Language.Default)
	} else {
		l = int(math.Min(float64(l), float64(Language.Lmax)))
	}

	if n < 1 || n > 3 {
		return ""
	}

	// The original formula is a bit obscure. This is a clearer implementation.
	// It maps a 0-360 degree value to an index from 0-15.
	numDirs := 1 << (n + 1) // 4, 8, or 16
	sectorSize := 360.0 / float64(numDirs)
	index := int(math.Mod(math.Round(float64(deg)/sectorSize), float64(numDirs)))

	// The DIRS array is for 16 directions, so we need to scale the index.
	scale := 16 / numDirs
	return Language.Dirs[l-1][index*scale]
}

// DirToDeg converts a compass direction string to degrees.
func DirToDeg(dir string, l int) int {
	if l == 0 {
		l = int(Language.Default)
	} else {
		l = int(math.Min(float64(l), float64(Language.Lmax)))
	}

	for i := 0; i < 16; i++ {
		if Language.Dirs[l-1][i] == dir {
			return (i * 45) / 2
		}
	}
	return 0
}

// Energy converts between energy units.
type Energy struct {
	YJ  float64 // Joules
	YC  float64 // Calories
	YWh float64 // Watt-hours
}

// Update executes the conversion logic.
func (e *Energy) Update(j, c, wh float64) {
	e.YJ = j + wh*3600.0 + c*4.1868
	e.YC = e.YJ * 0.238845896627496
	e.YWh = e.YJ * 2.7777777778e-4
}

// FToC converts Fahrenheit to Celsius.
func FToC(fahrenheit float64) float64 {
	return (fahrenheit - 32.0) * 0.5555555555555
}

// FToOm converts frequency (Hz) to angular frequency (Omega).
func FToOm(f float64) float64 {
	return Math.Pi2 * f
}

// FToPt converts frequency (Hz) to period time.
func FToPt(f float64) time.Duration {
	if f > 0.0 {
		return time.Duration(1000.0/f) * time.Millisecond
	}
	return 0
}

// GeoToDeg converts degrees, minutes, and seconds to decimal degrees.
func GeoToDeg(d, m int, sec float64) float64 {
	return float64(d) + float64(m)/60.0 + sec/3600.0
}

// KToC converts Kelvin to Celsius.
func KToC(kelvin float64) float64 {
	return kelvin + float64(Phys.T0)
}

// KmhToMs converts kilometers per hour to meters per second.
func KmhToMs(kmh float64) float64 {
	return kmh * 0.2777777777777
}

// Length converts between various length units.
type Length struct {
	Ym    float64 // Meters
	Yp    float64 // Typographic points
	Yin   float64 // Inches
	Yft   float64 // Feet
	Yyd   float64 // Yards
	Ymile float64 // Miles
	Ysm   float64 // Nautical miles
	Yfm   float64 // Fathoms
}

// Update executes the conversion logic.
func (l *Length) Update(m, p, in, ft, yd, mile, sm, fm float64) {
	l.Ym = m + p*0.000376065 + in*0.0254 + ft*0.3048 + yd*0.9144 + mile*1609.344 + sm*1852.0 + fm*1.829
	l.Yp = l.Ym * 2659.11478
	l.Yin = l.Ym * 39.3700787
	l.Yft = l.Ym * 3.2808398
	l.Yyd = l.Ym * 1.0936132
	l.Ymile = l.Ym * 0.00062137
	l.Ysm = l.Ym * 0.00053995
	l.Yfm = l.Ym * 0.5467468
}

// MsToBft converts wind speed from m/s to Beaufort.
func MsToBft(ms float64) int {
	return int(math.Pow(ms*1.196172, 0.666667))
}

// MsToKmh converts meters per second to kilometers per hour.
func MsToKmh(ms float64) float64 {
	return ms * 3.6
}

// OmToF converts angular frequency (Omega) to frequency (Hz).
func OmToF(om float64) float64 {
	return om / Math.Pi2
}

// Pressure converts between various pressure units.
type Pressure struct {
	Ymws  float64 // Meter water column
	Ytorr float64 // Torr
	Yatt  float64 // Technical atmosphere
	Yatm  float64 // Physical atmosphere
	Ypa   float64 // Pascal
	Ybar  float64 // Bar
}

// Update executes the conversion logic.
func (p *Pressure) Update(mws, torr, att, atm, pa, bar float64) {
	p.Ybar = bar + pa*1.0e-5 + att*0.980665 + atm*1.01325 + torr*0.001333224 + mws*0.0980665
	p.Ymws = p.Ybar * 10.197162
	p.Ytorr = p.Ybar * 750.0615
	p.Yatt = p.Ybar * 1.019716
	p.Yatm = p.Ybar * 0.986923
	p.Ypa = p.Ybar * 100000.0
}

// PtToF converts period time to frequency (Hz).
func PtToF(pt time.Duration) float64 {
	if pt > 0 {
		return 1000.0 / float64(pt.Milliseconds())
	}
	return 0.0
}

// Speed converts between various speed units.
type Speed struct {
	Yms  float64 // Meters per second
	Ykmh float64 // Kilometers per hour
	Ykn  float64 // Knots
	Ymh  float64 // Miles per hour
}

// Update executes the conversion logic.
func (s *Speed) Update(ms, kmh, kn, mh float64) {
	s.Yms = ms + kmh*0.27777777777778 + kn*0.5144444 + mh*0.44704
	s.Ykmh = s.Yms * 3.6
	s.Ykn = s.Yms * 1.9438446
	s.Ymh = s.Yms * 2.2369362
}

// Temperature converts between various temperature units.
type Temperature struct {
	YK  float64 // Kelvin
	YC  float64 // Celsius
	YF  float64 // Fahrenheit
	YRe float64 // Réaumur
	YRa float64 // Rankine
}

// Update executes the conversion logic.
func (t *Temperature) Update(k, c, f, re, ra float64) {
	// The original ST code sums all inputs, which is unusual.
	// A better approach would be to handle one input at a time.
	// This is a direct translation of the original summing behavior.
	t.YK = k + (c + 273.15) + (f+459.67)*5.0/9.0 + (re*1.25 + 273.15) + ra*5.0/9.0
	t.YC = t.YK - 273.15
	t.YF = t.YK*1.8 - 459.67
	t.YRe = (t.YK - 273.15) * 0.8
	t.YRa = t.YK * 1.8
}
