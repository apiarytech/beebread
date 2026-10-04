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

// Package hlk is the port of the OSCAT BUILDING HLK (heating, ventilation
// and air conditioning) functions and blocks: the physics of air and water,
// boilers and burners, heat meters, heating curves, legionella protection,
// tanks and the outside temperature.
package hlk

import (
	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/engineering"
	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/royaljelly/iec"
)

// AIR_DENSITY returns the density of air in kg/m³ at the temperature T in
// °C, the pressure P in Pa and the relative humidity RH in %.
func AIR_DENSITY(T, P, RH iec.REAL) iec.REAL {
	const (
		RL iec.REAL = 287.05
		RX iec.REAL = 3.773319e-3
	)
	return PHYS.PN * (1.0 - RH*SDD(T, true)*RX/P) / (RL * (T - PHYS.T0))
}

// AIR_ENTHALPY returns the enthalpy of air in kJ/kg at the temperature T in
// °C and the relative humidity RH in %.
func AIR_ENTHALPY(T, RH iec.REAL) iec.REAL {
	const (
		CPL iec.REAL = 1.00482
		CW  iec.REAL = 1.86
		LH  iec.REAL = 2500.78
	)
	return CPL*T + (CW*T+LH)*DEW_CON(RH, T)
}

// DEW_CON returns the concentration of water vapor in air in g/m³ at the
// relative humidity RH in % and the temperature T in °C, 0 for RH = 0 or T
// below -50 °C.
func DEW_CON(RH, T iec.REAL) iec.REAL {
	if RH > 0.0 && T > -50.0 {
		return 2.166824303e-2 * RH * SDD(T, false) / (T - PHYS.T0)
	}
	return 0.0
}

// DEW_RH returns the relative humidity in % of air with the concentration
// of water vapor VC in g/m³ at the temperature T in °C.
func DEW_RH(VC, T iec.REAL) iec.REAL {
	return LIMIT(0.0, VC/DEW_CON(1.0, T), 100.0)
}

// DEW_TEMP returns the dew point in °C of air at the relative humidity RH
// in % and the temperature T in °C, and absolute zero for RH = 0.
func DEW_TEMP(RH, T iec.REAL) iec.REAL {
	const (
		a iec.REAL = 7.5
		b iec.REAL = 237.3
	)
	if RH > 0.0 {
		v := LOG(RH * 0.01 * math.EXP10((a*T)/(b+T)))
		return b * v / (a - v)
	}
	return PHYS.T0
}

// HEAT_INDEX returns the heat index in °C, the temperature felt at the
// temperature T in °C and the relative humidity RH in %, T itself below
// 20 °C or 20%.
func HEAT_INDEX(T, RH iec.REAL) iec.REAL {
	if RH < 20.0 || T < 20.0 {
		return T
	}
	rh2 := RH * RH
	T = engineering.C_TO_F(T)
	t2 := T * T
	hi := -42.379 + 2.04901523*T - 6.83783e-3*t2 +
		RH*(10.1433127-0.22475541*T+1.22874e-3*t2) +
		rh2*(8.5282e-4*T-5.481717e-2-1.99e-6*t2)
	return engineering.F_TO_C(hi)
}

// SDD returns the saturation vapor pressure of water in Pa at the
// temperature T in °C, over ice if ICE is true.
func SDD(T iec.REAL, ICE iec.BOOL) iec.REAL {
	if ICE {
		return 611.153 * EXP(22.4433*T/(272.186+T))
	}
	return 611.213 * EXP(17.5043*T/(241.2+T))
}

// SDD_NH3 returns the saturation vapor pressure of ammonia in bar at the
// temperature T in °C, by the Antoine equation with NIST's parameters.
func SDD_NH3(T iec.REAL) iec.REAL {
	if T < -33.65 {
		return EXP(7.3396511649 - (1166.7498002 / (T + 192.37)))
	}
	return EXP(11.210964456 - (2564.9140075 / (T + 262.741)))
}

// SDT_NH3 returns the saturation temperature of ammonia in °C at the vapor
// pressure P in bar, -110 °C below 0.001 bar.
func SDT_NH3(P iec.REAL) iec.REAL {
	switch {
	case P < 1.0e-3:
		return -110.0
	case P < 1.0:
		return 506.713/(3.18757-LOG(P)) - 192.37
	}
	return 1113.928/(4.86886-LOG(P)) - 262.71
}

// TANK_VOL1 returns the volume of a horizontal cylindrical tank of the
// radius TR and the length TL filled to the height H.
func TANK_VOL1(TR, TL, H iec.REAL) iec.REAL {
	return math.CIRCLE_SEG(TR, H) * TL
}

// TANK_VOL2 returns the volume of a spherical tank of the radius TR filled
// to the height H.
func TANK_VOL2(TR, H iec.REAL) iec.REAL {
	return MATH.PI * H * H * (TR - H/3.0)
}

// table returns the points of a table for LINEAR_INT.
func table(pts ...iec.REAL) (xy [20][2]iec.REAL) {
	for i := 0; i+1 < len(pts); i += 2 {
		xy[i/2] = [2]iec.REAL{pts[i], pts[i+1]}
	}
	return xy
}

// The specific heat capacity and enthalpy of water by temperature.
var (
	waterCP = table(0.0, 4.228, 5.0, 4.20, 10.0, 4.188, 15.0, 4.184, 50.0, 4.181,
		60.0, 4.183, 70.0, 4.187, 80.0, 4.194, 90.0, 4.204, 100.0, 4.22)
	waterEnthalpy = table(0.0, 0.06, 10.0, 42.1, 20.0, 83.9, 30.0, 125.8, 40.0, 167.58,
		50.0, 209.4, 60.0, 251.2, 70.0, 293.1, 80.0, 335.0, 90.0, 377.0, 100.0, 419.1)
)

// WATER_CP returns the specific heat capacity of water in kJ/(kg·K) at the
// temperature T in °C, 0..100 °C.
func WATER_CP(T iec.REAL) iec.REAL {
	return math.LINEAR_INT(T, waterCP, 10)
}

// WATER_DENSITY returns the density of water in kg/m³ at the temperature T
// in °C, free of air or, if SAT is true, saturated with it.
func WATER_DENSITY(T iec.REAL, SAT iec.BOOL) iec.REAL {
	const (
		a0 iec.REAL = 999.83952
		a1 iec.REAL = 16.952577
		a2 iec.REAL = -7.9905127e-3
		a3 iec.REAL = -4.6241757e-5
		a4 iec.REAL = 1.0584601e-7
		a5 iec.REAL = -2.8103006e-10
		b  iec.REAL = 0.0168872
	)
	t2 := T * T
	t4 := t2 * t2
	d := (a0 + a1*T + a2*t2 + a3*t2*T + a4*t4 + a5*t4*T) / (1.0 + b*T)
	if SAT {
		d = d - 0.004612 + 0.000106*T
	}
	return d
}

// WATER_ENTHALPY returns the specific enthalpy of water in kJ/kg at the
// temperature T in °C, 0..100 °C.
func WATER_ENTHALPY(T iec.REAL) iec.REAL {
	return math.LINEAR_INT(T, waterEnthalpy, 11)
}

// WCT returns the wind chill temperature in °C at the temperature T in °C
// and the wind speed V in km/h, T itself above 10 °C or below 5 km/h.
func WCT(T, V iec.REAL) iec.REAL {
	if V < 5.0 || T > 10.0 {
		return T
	}
	return 13.12 + 0.6215*T + (0.3965*T-11.37)*EXP(LN(V)*0.16)
}
