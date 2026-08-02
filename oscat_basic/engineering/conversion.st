(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK ASTRO
VAR_INPUT
	m : REAL;
	AE : REAL;
	PC : REAL;
	LJ : REAL;
END_VAR
VAR_OUTPUT
	Ym : REAL;
	YAE : REAL;
	YPC : REAL;
	YLJ : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		oscat

this function converts different length units
any unused input can simply be left open.
different inputs connected at the same time will be added up.

*)
(* @END_DECLARATION := '0' *)
YAE :=	AE
		+ m * 6.6845871535E-012
		+ PC * 206265.0
		+ LJ * 63240.0;
Ym := YAE * 149.597870E9;
YPC := YAE * 4.8481322570E-006;
YLJ := YAE * 1.5812776724E-005;


(*
Länge Meter m SI-Basiseinheit
Astronomische Einheit* AE 1 AE = 149,597 870 · E9 m
Parsec pc 1 pc = 206265 AE = 30,857 · E15 m
Lichtjahr Lj 1 Lj = 9,460 530 · E15 m = 63240 AE = 0,306 59 pc
Ångström Å 1 Å = E−l0 m
typograph. Punkt p 1 p = 0,376 065 mm · im Druckereigewerbe
inch** in 1 in = 2,54 · E−2 m = 25,4 mm***
foot ft 1 ft = 0,3048 m = 30,48 cm
yard yd 1 yd = 0,9144 m
mile mile 1 mile = 1609,344 m
Internat. Seemeile sm 1 sm = 1852 m
Fathom fm 1 fm = 1,829 m · in der Seeschifffahrt
*)

(* revision history

hm	27. mar. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	improved code

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BFT_TO_MS : REAL
VAR_INPUT
	BFT : INT;
END_VAR


(*
version 1.0	12. jun 2008
programmer 	hugo
tested by		oscdat

this function converts wind speed from beaufort to m/s
*)
(* @END_DECLARATION := '0' *)
BFT_TO_MS := EXPT(BFT, 1.5) * 0.836;


(* revision history
hm	12. 6. 2008		rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION C_TO_F : REAL
VAR_INPUT
	celsius : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function converts celsius  to fahrenheit 

*)
(* @END_DECLARATION := '0' *)
C_TO_F := celsius * 1.8 + 32.0;


(* revision history

hm	4. aug. 2006	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION C_TO_K : REAL
VAR_INPUT
	Celsius : REAL;
END_VAR


(*
version 1.1	19. aug 2009
programmer 	hugo
tested by		tobias

this function converts celsius to kelvin
*)
(* @END_DECLARATION := '0' *)
C_TO_K := Celsius - phys.T0;

(* revision history

hm	4. aug 2006	rev 1.0
	original version

hm	19. aug 2009	rev 1.1
	fixed calculation error

*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DEG_TO_DIR : STRING(3)
VAR_INPUT
	DEG : INT;
	N : INT;
	L : INT;
END_VAR
VAR
	ly: INT;
END_VAR


(*
version 1.1	22. oct. 2008
programmer 	hugo
tested by		oscat

this function converts degrees in compass direction.
the function supports output in english (L=0) and german (L=1).

*)
(* @END_DECLARATION := '0' *)
IF L = 0 THEN ly := LANGUAGE.DEFAULT; ELSE ly := MIN(L, LANGUAGE.LMAX); END_IF;
DEG_TO_DIR := LANGUAGE.DIRS[ly, ((SHL(DEG,N-1) + 45) / 90) MOD SHL(INT#2,N)*SHR(INT#8,N)];



(*
DIR := ((SHL(DEG,N-1) + 45) / 90) MOD SHL(INT#2,N);
explanation :
DIR is calculated BY the following formula:
DIR := ((DIR + 45) / 90) MOD 4 if N = 1 digit
North = 0, East = 1 ....
DIR := ((DIR + 22,5) / 45) MOD 8 if N = 2 digit
convert to integer calculation
DIR := ((DIR*2 + 45) / 90) MOD 8
N = 0, NE = 1 ....
ther above formula replaces 2^N with shift for performance
*)



(* revision histroy
hm	11. jun. 2008	rev 1.0
	original release

hm	22. oct. 2008	rev 1.1
	changed size of string variables to 30
*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DIR_TO_DEG : INT
VAR_INPUT
	DIR : STRING(3);
	L : INT;
END_VAR
VAR
	ly: INT;
	i : INT;
END_VAR

(*
version 1.1	22. oct. 2008
programmer 	hugo
tested by	oscat

this function converts compass directions to degrees
it will recognize up to 3 letter directions in english and german writing.

*)
(* @END_DECLARATION := '0' *)
IF L = 0 THEN ly := LANGUAGE.DEFAULT; ELSE ly := MIN(L, LANGUAGE.LMAX); END_IF;
FOR i := 0 TO 15 DO
	IF language.DIRS[ly, i] = DIR THEN EXIT; END_IF;
END_FOR;
DIR_TO_DEG := SHR(i * 45 + 1, 1);






(* revision histroy
hm	22. oct. 2008	rev 1.1
	original release

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK ENERGY
VAR_INPUT
	J : REAL;
	C : REAL;
	Wh : REAL;
END_VAR
VAR_OUTPUT
	YJ : REAL;
	YC : REAL;
	YWh : REAL;
END_VAR


(*
version 1.2	16. jan. 2010
programmer 	hugo
tested by		oscat

this function converts different energy units
any unused input can simply be left open.
different inputs connected at the same time will be added up.
*)
(* @END_DECLARATION := '0' *)
YJ := J + Wh * 3600.0 + C * 4.1868;
YC := YJ * 0.238845896627496;
YWh := YJ * 2.7777777778E-004;

(*
Arbeit, Energie, Joule* J 1 J = 1 N · m = 1 W · s = (1/3,6) E−6 kW · h = 1 kg · m2/s2
Wärmemenge Kilowattstunde kW · h 1 kW · h = 3,6 MJ = 860 kcal
Elektronvolt eV 1 eV = 160,218 92 E−21 J
Erg erg 1 erg = 1E−7 J
Kalorie calorie 1 calalorie = 4,1868 J = 1,163 E−3 W · h
Therm therm 1 therm = 105,50 · 106 J
*)

(* revision history

hm	27. mar. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	improved code

hm 16. jan 2010	rev 1.2
	avoid the string cal in comments for codesys import bug


*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_TO_C : REAL
VAR_INPUT
	fahrenheit : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested BY		tobias

this FUNCTION converts fahrenheit TO celsius

*)
(* @END_DECLARATION := '0' *)
F_TO_C := (fahrenheit - 32.0) * 0.5555555555555;

(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_TO_OM : REAL
VAR_INPUT
	F : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function converts frequency to Omega F
Omega = 2*PI*F
*)
(* @END_DECLARATION := '0' *)
F_TO_OM := math.PI2 * F;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	unsing math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_TO_PT : TIME
VAR_INPUT
	F : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function converts frequency to periode time  
*)
(* @END_DECLARATION := '0' *)
F_TO_PT := DWORD_TO_TIME(REAL_TO_DWORD(1.0 / F * 1000.0));


(* revision history

hm	4. aug. 2006	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION GEO_TO_DEG : REAL
VAR_INPUT
	D : INT;
	M : INT;
	SEC : REAL;
END_VAR


(*
version 1.0	22. jan. 2009
programmer 	hugo
tested by		oscat

this function converts degrees, minutes seconds to decimal degrees.

*)
(* @END_DECLARATION := '0' *)
GEO_TO_DEG := INT_TO_REAL(D) + INT_TO_REAL(M) * 0.016666666666667 + sec * 0.00027777777777778;


(* revision histroy
hm	22. jan. 2009	rev 1.0
	original release


*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION K_TO_C : REAL
VAR_INPUT
	Kelvin : REAL;
END_VAR


(*
version 1.1	19. aug 2009
programmer 	hugo
tested by		tobias

this function converts kelvin to celsius

*)
(* @END_DECLARATION := '0' *)
K_TO_C := Kelvin + phys.T0;


(* revision history

hm	4. aug 2006	rev 1.0
	original version

hm	19. aug 2009	rev 1.1
	fixed calculation error

*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION KMH_TO_MS : REAL
VAR_INPUT
	kmh : REAL;
END_VAR


(*
version 1.1	6 jan 2007
programmer 	hugo
tested by		tobias

this function converts velocities from Kilometers / hour to Meters / Second

*)
(* @END_DECLARATION := '0' *)
KMH_TO_MS := kmh * 0.2777777777777;

(* revision history
hm	4. feb 2007		rev 1.0
	original version

hm	6. jan 2008		rev 1.1
	improved performance

*)


END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/conversion' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK LENGTH
VAR_INPUT
	m : REAL;
	p : REAL;
	in : REAL;
	ft : REAL;
	yd : REAL;
	mile : REAL;
	sm : REAL;
	fm : REAL;
END_VAR
VAR_OUTPUT
	Ym : REAL;
	Yp : REAL;
	Yin : REAL;
	Yft : REAL;
	Yyd : REAL;
	Ymile : REAL;
	Ysm : REAL;
	Yfm : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		oscat

this function converts different length units
any unused input can simply be left open.
different inputs connected at the same time will be added up.

*)
(* @END_DECLARATION := '0' *)
Ym :=	m
		+ p * 0.000376065
		+ in * 0.0254
		+ ft * 0.3048
		+ yd * 0.9144
		+ mile * 1609.344
		+ sm * 1852.0
		+ fm * 1.829;
Yp := Ym * 2659.11478068951;
Yin := Ym * 39.37007874016;
Yft := Ym * 3.28083989501;
Yyd := Ym * 1.09361329834;
Ymile := Ym * 0.00062137119;
Ysm := Ym * 0.00053995680;
Yfm := Ym * 0.54674685621;

(*
Länge Meter m SI-Basiseinheit
Astronomische Einheit* AE 1 AE = 149,597 870 · E9 m
Parsec pc 1 pc = 206265 AE = 30,857 · E15 m
Lichtjahr Lj 1 Lj = 9,460 530 · E15 m = 63240 AE = 0,306 59 pc
Ångström Å 1 Å = E−l0 m
typograph. Punkt p 1 p = 0,376 065 mm · im Druckereigewerbe
inch** in 1 in = 2,54 · E−2 m = 25,4 mm***
foot ft 1 ft = 0,3048 m = 30,48 cm
yard yd 1 yd = 0,9144 m
mile mile 1 mile = 1609,344 m
Internat. Seemeile sm 1 sm = 1852 m
Fathom fm 1 fm = 1,829 m · in der Seeschifffahrt
*)

(* revision history

hm	27. mar. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	improved code

*)

END_FUNCTION_BLOCK