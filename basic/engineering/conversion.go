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

package engineering

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/royaljelly/iec"
)

// The unit conversion blocks add up all their inputs, so an unused input is
// left at 0.

// ASTRO converts astronomical lengths: meters, astronomical units, parsecs
// and light years.
type ASTRO struct {
	M, AE, PC, LJ     iec.REAL
	YM, YAE, YPC, YLJ iec.REAL
}

// INIT resets the block.
func (a *ASTRO) INIT() { *a = ASTRO{} }

// Execute runs the block once.
func (a *ASTRO) Execute(now time.Time) {
	a.YAE = a.AE + a.M*6.6845871535e-012 + a.PC*206265.0 + a.LJ*63240.0
	a.YM = a.YAE * 149.597870e9
	a.YPC = a.YAE * 4.8481322570e-006
	a.YLJ = a.YAE * 1.5812776724e-005
}

// BFT_TO_MS converts a wind speed in Beaufort to m/s.
func BFT_TO_MS(bft iec.INT) iec.REAL {
	return EXPT(iec.REAL(bft), 1.5) * 0.836
}

// C_TO_F converts °C to °F.
func C_TO_F(celsius iec.REAL) iec.REAL {
	return celsius*1.8 + 32.0
}

// C_TO_K converts °C to K.
func C_TO_K(celsius iec.REAL) iec.REAL {
	return celsius - PHYS.T0
}

// DEG_TO_DIR converts an angle in degrees to a compass direction in the
// language l, or the default language if l is 0, with n = 1 for N, E, S, W,
// 2 for 8 directions and 3 for 16.
func DEG_TO_DIR(deg, n, l iec.INT) iec.STRING {
	ly := LANGUAGE.DEFAULT
	if l != 0 {
		ly = min(l, LANGUAGE.LMAX)
	}
	shl := func(x, s iec.INT) iec.INT { return iec.INT(SHL(uint16(x), s)) }
	shr := func(x, s iec.INT) iec.INT { return iec.INT(SHR(uint16(x), s)) }
	i := (shl(deg, n-1) + 45) / 90 % shl(2, n) * shr(8, n)
	return LANGUAGE.DIRS[ly-1][LIMIT(0, i, 15)]
}

// DIR_TO_DEG converts a compass direction of up to 3 letters in the
// language l, or the default language if l is 0, to degrees.
func DIR_TO_DEG(dir iec.STRING, l iec.INT) iec.INT {
	ly := LANGUAGE.DEFAULT
	if l != 0 {
		ly = min(l, LANGUAGE.LMAX)
	}
	i := iec.INT(0)
	for ; i <= 15; i++ {
		if LANGUAGE.DIRS[ly-1][i] == dir {
			break
		}
	}
	return (i*45 + 1) >> 1
}

// ENERGY converts energies: Joule, calories and watt hours.
type ENERGY struct {
	J, C, WH    iec.REAL
	YJ, YC, YWH iec.REAL
}

// INIT resets the block.
func (e *ENERGY) INIT() { *e = ENERGY{} }

// Execute runs the block once.
func (e *ENERGY) Execute(now time.Time) {
	e.YJ = e.J + e.WH*3600.0 + e.C*4.1868
	e.YC = e.YJ * 0.238845896627496
	e.YWH = e.YJ * 2.7777777778e-004
}

// F_TO_C converts °F to °C.
func F_TO_C(fahrenheit iec.REAL) iec.REAL {
	return (fahrenheit - 32.0) * 0.5555555555555
}

// F_TO_OM converts a frequency to the angular frequency 2πF.
func F_TO_OM(f iec.REAL) iec.REAL {
	return MATH.PI2 * f
}

// F_TO_PT converts a frequency to its period.
func F_TO_PT(f iec.REAL) iec.TIME {
	return DWORD_TO_TIME(REAL_TO_DWORD(1.0 / f * 1000.0))
}

// GEO_TO_DEG converts degrees, minutes and seconds to decimal degrees.
func GEO_TO_DEG(d, m iec.INT, sec iec.REAL) iec.REAL {
	return iec.REAL(d) + iec.REAL(m)*0.016666666666667 + sec*0.00027777777777778
}

// K_TO_C converts K to °C.
func K_TO_C(kelvin iec.REAL) iec.REAL {
	return kelvin + PHYS.T0
}

// KMH_TO_MS converts km/h to m/s.
func KMH_TO_MS(kmh iec.REAL) iec.REAL {
	return kmh * 0.2777777777777
}

// LENGTH converts lengths: meters, typographic points, inches, feet, yards,
// miles, sea miles and fathoms.
type LENGTH struct {
	M, P, IN, FT, YD, MILE, SM, FM         iec.REAL
	YM, YP, YIN, YFT, YYD, YMILE, YSM, YFM iec.REAL
}

// INIT resets the block.
func (l *LENGTH) INIT() { *l = LENGTH{} }

// Execute runs the block once.
func (l *LENGTH) Execute(now time.Time) {
	l.YM = l.M + l.P*0.000376065 + l.IN*0.0254 + l.FT*0.3048 + l.YD*0.9144 +
		l.MILE*1609.344 + l.SM*1852.0 + l.FM*1.829
	l.YP = l.YM * 2659.11478068951
	l.YIN = l.YM * 39.37007874016
	l.YFT = l.YM * 3.28083989501
	l.YYD = l.YM * 1.09361329834
	l.YMILE = l.YM * 0.00062137119
	l.YSM = l.YM * 0.00053995680
	l.YFM = l.YM * 0.54674685621
}

// MS_TO_BFT converts a wind speed in m/s to Beaufort.
func MS_TO_BFT(ms iec.REAL) iec.INT {
	return REAL_TO_INT(EXPT(ms*1.196172, 0.666667))
}

// MS_TO_KMH converts m/s to km/h.
func MS_TO_KMH(ms iec.REAL) iec.REAL {
	return ms * 3.6
}

// OM_TO_F converts an angular frequency to a frequency.
func OM_TO_F(om iec.REAL) iec.REAL {
	return om / MATH.PI2
}

// PRESSURE converts pressures: meters of water, torr, technical and
// physical atmospheres, pascal and bar.
type PRESSURE struct {
	MWS, TORR, ATT, ATM, PA, BAR       iec.REAL
	YMWS, YTORR, YATT, YATM, YPA, YBAR iec.REAL
}

// INIT resets the block.
func (p *PRESSURE) INIT() { *p = PRESSURE{} }

// Execute runs the block once.
func (p *PRESSURE) Execute(now time.Time) {
	p.YBAR = p.BAR + p.PA*1.0e-5 + 0.980665*p.ATT + 1.01325*p.ATM + 0.001333224*p.TORR + 0.0980665*p.MWS
	p.YMWS = p.YBAR * 10.1971621297793
	p.YTORR = p.YBAR * 750.0615050434140
	p.YATT = p.YBAR * 1.0197162129779
	p.YATM = p.YBAR * 0.9869232667160
	p.YPA = p.YBAR * 100000.0
}

// PT_TO_F converts a period to its frequency.
func PT_TO_F(pt iec.TIME) iec.REAL {
	return 1000.0 / TIME_TO_REAL(pt)
}

// SPEED converts speeds: m/s, km/h, knots and miles per hour.
type SPEED struct {
	MS, KMH, KN, MH     iec.REAL
	YMS, YKMH, YKN, YMH iec.REAL
}

// INIT resets the block.
func (s *SPEED) INIT() { *s = SPEED{} }

// Execute runs the block once.
func (s *SPEED) Execute(now time.Time) {
	s.YMS = s.MS + s.KMH*0.27777777777778 + s.KN*0.5144444 + s.MH*0.44704
	s.YKMH = s.YMS * 3.6
	s.YKN = s.YMS * 1.94384466037535
	s.YMH = s.YMS * 2.23693629205440
}

// TEMPERATURE converts temperatures: Kelvin, Celsius, Fahrenheit, Réaumur
// and Rankine. The initial values of C, F and RE are absolute zero, which
// adds nothing.
type TEMPERATURE struct {
	K, C, F, RE, RA      iec.REAL
	YK, YC, YF, YRE, YRA iec.REAL
}

// INIT resets the block and sets C, F and RE to their initial values.
func (t *TEMPERATURE) INIT() { *t = TEMPERATURE{C: -273.15, F: -459.67, RE: -218.52} }

// Execute runs the block once.
func (t *TEMPERATURE) Execute(now time.Time) {
	t.YK = t.K + (t.C + 273.15) + (t.F+459.67)*0.555555555555 + (t.RE*1.25 + 273.15) + t.RA*0.555555555555
	t.YC = t.YK - 273.15
	t.YF = t.YK*1.8 - 459.67
	t.YRE = (t.YK - 273.15) * 0.8
	t.YRA = t.YK * 1.8
}
