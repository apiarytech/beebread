(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION MULTI_IN : REAL
VAR_INPUT
	in_1 : REAL;
	in_2 : REAL;
	in_3 : REAL;
	default : REAL;
	in_min : REAL;
	in_max : REAL;
	mode : BYTE;
END_VAR
VAR
	count: INT;
	F1: BOOL;
	F2: BOOL;
	F3: BOOL;
END_VAR

(*

version 1.4	18. jul. 2009
programmer 	oscat
tested by		oscat

multi_in is a signal conditioning function which can be configured in 8 different ways.
it is used to read multiple sensors for the same value and protect the user from broken sensors or invalid data from sensors.

multi_in can be configured in different ways:
	mode = 0 means the the avg of the 3 inputs is used, this mode is the default mode.
	mode = 1 means in_1 is used.
	mode = 2 means in_2 is used.
	mode = 3 means in_3 is used.
	mode = 4 means default is used.
	mode = 5 means the lowest of the 3 external temperatures is used.
	mode = 6 means the higest externnal temperature is used.
	mode = 7 means the midlle input is used, if there are only two working, the lowest input is used.
	mode > 7 ,eans output is 0 at all times.
	in any config mode, an input higher then in_max or lower then in_min is ignored to prevent values from broken sensors or wires.
	if all inputs are higher then in_max or lower then in_min, a default value (default) is used.
*)

(* @END_DECLARATION := '0' *)
F1 :=  in_1 > in_min AND in_1 < in_max;
F2 :=  in_2 > in_min AND in_2 < in_max;
F3 :=  in_3 > in_min AND in_3 < in_max;

CASE mode OF
0:	count := 0;
	IF F1 THEN
		count := count + 1;
		MULTI_IN := in_1;
	ELSE
		MULTI_IN := 0.0;
	END_IF;
	IF F2 THEN
		count := count + 1;
		MULTI_IN := MULTI_IN + in_2;
	END_IF;
	IF F3 THEN
		count := count + 1;
		MULTI_IN := MULTI_IN + in_3;
	END_IF;
	MULTI_IN := SEL(count = 0, MULTI_IN / INT_TO_REAL(count), default);

1:	MULTI_IN := SEL(F1, default, IN_1);

2:	MULTI_IN := SEL(F2, default, IN_2);

3:	MULTI_IN := SEL(F3, default, IN_3);

4:	MULTI_IN := default;

5:	MULTI_IN := SEL(F1, in_max, IN_1);
	IF F2 AND in_2 < MULTI_IN THEN MULTI_IN := in_2; END_IF;
	IF F3 AND in_3 < MULTI_IN THEN MULTI_IN := in_3; END_IF;
	IF MULTI_IN = in_max THEN MULTI_IN := default; END_IF;

6:	MULTI_IN := SEL(F1, in_min, IN_1);
	IF F2 AND in_2 > MULTI_IN THEN MULTI_IN := in_2; END_IF;
	IF F3 AND in_3 > MULTI_IN THEN MULTI_IN := in_3; END_IF;
	IF MULTI_IN = in_min THEN MULTI_IN := default; END_IF;

7:	IF F1 AND F2 AND F3 THEN MULTI_IN := MID3(in_1, in_2, in_3);
	ELSIF F1 AND F2 THEN MULTI_IN := MIN(in_1, in_2);
	ELSIF F1 AND F3 THEN MULTI_IN := MIN(in_1, in_3);
	ELSIF F2 AND F3 THEN MULTI_IN := MIN(in_2, in_3);
	ELSIF F1 THEN MULTI_IN := in_1;
	ELSIF F2 THEN MULTI_IN := in_2;
	ELSIF F3 THEN MULTI_IN := in_3;
	ELSE MULTI_IN := default;
	END_IF;

END_CASE;

(*
hm 1.1.2007		rev 1.1
	changed midr to mid3 function

hm	14. 10. 2008	rev 1.2
	corrected an error for in_3 overrange detection
	improved performance

hm	11. mar. 2009	rev 1.3
	changed real constants to use dot syntax

hm	18. jul. 2009	rev 1.4
	improved performance
*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION RES_NI : REAL
VAR_INPUT
	T : REAL;
	R0 : REAL;
END_VAR
VAR CONSTANT
	A : REAL := 0.5485;
	B : REAL := 0.665E-3;
	C : REAL := 2.805E-9;
END_VAR
VAR
	T2: REAL;
END_VAR

(*
version 1.2	2 dec 2007
programmer 	hugo
tested by		tobias

this function returs the resistance for a nickel sensor for a temperature range from -60..+180 °C *)
(* @END_DECLARATION := '0' *)
T2 := T * T;
RES_NI := R0 + A * T + B * T2 + C * T2 * T2;

(* revision history

hm	4.8.2006		rev 1.0
	original version

hm	13.9.2007		rev 1.1
	changed coding for better performance

hm	2. dec. 2007	rev 1.2
	changed code for better performance
*)

(*

 Auszug aus DIN 43760 für Ni100

°C		R		°C		R		°C	R		°C		R		°C		R

-60 	69,5 	-10 	94,6 	40 	123,0 	90 		154,9 	140 	190,9
-55 	71,9 	-5 		97,3 	45 	126,0 	95 		158,3 	145 	194,8
-50 	74,3 	0 		100,0 	50 	129,1 	100 	161,8 	150 	198,7
-45 	76,7 	5 		102,8 	55 	132,2 	105 	165,3 	155 	202,6
-40 	79,1 	10 		105,6 	60 	135,3 	110 	168,8 	160 	206,6
-35 	81,6 	15 		108,4 	65 	138,5 	115 	172,4 	165 	210,7
-30 	84,2 	20 		111,2 	70 	141,7 	120 	176,0 	170 	214,8
-25 	86,7 	25 		114,1 	75 	145,0 	125 	179,6 	175 	219,0
-20 	89,3 	30 		117,1 	80 	148,3 	130 	183,3 	180 	223,2
-15 	91,9 	35 		120,0 	85 	151,6 	135 	187,1 	  	 

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION RES_NTC : REAL
VAR_INPUT
	T : REAL;
	RN : REAL;
	B : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function returs the resistance for a NTC sensor for a given temperature in °C.
RN is the resistance at 25 °C and B is a constant for the given sensor.

*)
(* @END_DECLARATION := '0' *)
RES_NTC := RN * EXP(B * (1.0 / (T+273.15) - 0.00335401643468053));



(* revision history

hm 30. dec. 2008	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	changed real constants to use dot syntax

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION RES_PT : REAL
VAR_INPUT
	T : REAL;
	R0 : REAL;
END_VAR
VAR CONSTANT
	A : REAL := 3.90802E-3 ;
	B : REAL := -5.802E-7;
	C : REAL := -4.27350E-12;
END_VAR
VAR
	T2: REAL;
END_VAR

(*
version 1.3	11. mar. 2009
programmer 	hugo
tested by		oscat


this function returs the resistance for a platinum sensor for a temperature range from -200..+850 °C *)
(* @END_DECLARATION := '0' *)
T2 := T * T;

IF T >= 0.0 THEN
	RES_PT := R0 * ( 1.0 + A * T + B * T2);
ELSE
	RES_PT := R0 * ( 1.0 + A * T + B * T2 + C * (T - 100.0) * T2 * T);
END_IF;



(* revision history
hm	4. aug. 2006	rev 1.1
	original version

hm	2. dec 2007	rev 1.2
	changed code for better performance

hm	11. mar. 2009	rev 1.3
	changed real constants to use dot syntax

*)


(*
 Auszug aus DIN 43760 für Pt100

°C		R		°C		R		°C	R			°C		R			°C		R			°C		R

-200 	18,49 	-100 	60,25 	0 	100,00 	100 	138,50 	200 	175,84 	300 	212,02
-195 	20,65 	-95 	62,28 	5 	101,95 	105 	140,39 	205 	177,68 	305 	213,80
-190 	22,80 	-90 	64,30 	10 	103,90 	110 	142,29 	210 	179,51 	310 	215,57
-185 	24,94 	-85 	66,31 	15 	105,85 	115 	144,17 	215 	181,34 	315 	217,35
-180 	27,08 	-80 	68,33 	20 	107,79 	120 	146,06 	220 	183,17 	320 	219,12
-175 	29,20 	-75 	70,33 	25 	109,73 	125 	147,94 	225 	184,99 	325 	220,88
-170 	31,32 	-70 	72,33 	30 	111,67 	130 	149,82 	230 	186,82 	330 	222,65
-165 	33,43 	-65 	74,33 	35 	113,61 	135 	151,70 	235 	188,63 	335 	224,41
-160 	35,33 	-60 	76,33 	40 	115,54 	140 	153,58 	240 	190,45 	340 	226,17
-155 	37,63 	-55 	78,32 	45 	117,47 	145 	155,45 	245 	192,26 	345 	227,92
-150 	39,71 	-50 	80,31 	50 	119,40 	150 	157,31 	250 	194,07 	350 	229,67
-145 	41,79 	-45 	82,29 	55 	121,32 	155 	159,18 	255 	195,88 	355 	231,42
-140 	43,87 	-40 	84,27 	60 	123,24 	160 	161,04 	260 	197,69 	360 	233,17
-135 	45,94 	-35 	86,25 	65 	125,16 	165 	162,90 	265 	199,49 	365 	234,91
-130 	48,00 	-30 	88,22 	70 	127,07 	170 	164,76 	270 	201,29 	370 	236,65
-125 	50,06 	-25 	90,19 	75 	128,98 	175 	166,61 	275 	203,08 	375 	238,39
-120 	52,11 	-20 	92,16 	80 	130,89 	180 	168,46 	280 	204,88 	380 	240,13
-115 	54,15 	-15 	94,12 	85 	132,80 	185 	170,31 	285 	206,67 	385 	241,86
-110 	56,19 	-10 	96,09 	90 	134,70 	190 	172,16 	290 	208,45 	390 	243,59
-105 	58,22 	-5 		98,04 	95 	136,60 	195 	174,00 	295 	210,24 	395 	245,31

*)
(*
revision history


*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION RES_SI : REAL
VAR_INPUT
	T : REAL;
	RS : REAL;
	TS : REAL;
END_VAR
VAR CONSTANT
	A : REAL := 7.64E-3;
	B : REAL := 1.66E-5;
END_VAR
VAR
	TX: REAL;
END_VAR

(*
version 1.2	11. mar. 2009
programmer 	hugo
tested by		tobias

this function returs the resistance for a silicon sensor for a temperature range from -50..+150 °C
for example KTY10 normaly Rs at 25 °C
*)
(* @END_DECLARATION := '0' *)
TX := T - TS;
RES_SI := RS * ( 1.0 + A * TX + B * TX * TX);

(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	2. dec 2007	rev 1.1
	changed code for better performance

hm	11. mar. 2009	rev 1.2
	changed real constants to use dot syntax

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION SENSOR_INT : REAL
VAR_INPUT
	Voltage : REAL;
	Current : REAL;
	RP : REAL;
	RS : REAL;
END_VAR
VAR
	RG: REAL;
END_VAR

(*
version 1.3	11. mar. 2009
programmer 	hugo
tested by		tobias

this function calculates the real resistance of a sensor RX given a parasitic resistor in parallel to the sensor and an additional serial resistor.
for example a parasitic parallel resistor for a PT100 can be ignored but a series resistor is the cable to the sensor.
for a humidity sensor, the parallel resistor is more dramatic and needs to be considered.
the function allows to consider a parallel and serial resistor at the same time.
with this function and known parasitic serial and parallel resistors the true resistance of the sensor can be calculated.

*)
(* @END_DECLARATION := '0' *)
IF current <> 0.0 THEN
	RG := voltage / current;
	SENSOR_INT := RP * ( RG - RS) / ( RP + RS - RG);
END_IF;

(* revisaion history

hm 20.1.2007		rev 1.1
	deleted input R0 which was not needed.

hm	2. dec 2007	rev 1.2
	corrected an error in formula and improoved speed

hm	11. mar. 2009	rev 1.3
	changed real constants to use dot syntax

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION TEMP_NI : REAL
VAR_INPUT
	Res : REAL;
	R0 : REAL;
END_VAR


(*
version 1.5	10. mar. 2009
programmer 	hugo
tested by		tobias

this function returns the temperature for a nickel sensor in a range from -60..+180 °C *)
(* @END_DECLARATION := '0' *)
TEMP_NI := (SQRT(0.30085225 - 2.66E-3 * (R0 - Res)) - 0.5485) * 751.8796992;


(* revision history

hm 24.1.2007		rev 1.1
	deleted unused variable A, B, C

hm 10.9.2007		rev 1.2
	changed accuracy to 0.02 degrees to improove performace

hm 17. 12 2007	rev 1.3
	improovements for better performance and higher accuracy

hm	6. jan 2008	rev 1.4
	further performance improvement

hm	10. mar. 2009	rev 1.5
	removed nested comments

*)


END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION TEMP_NTC : REAL
VAR_INPUT
	RES : REAL;
	RN : REAL;
	B : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function returs the temperature for a NTC sensor for a range from 0..85 °C.
RN is the resistance at 25 °C and B is a constant for the given sensor.

*)
(* @END_DECLARATION := '0' *)
IF res > 0.0 THEN
	TEMP_NTC := B * 298.15 / (B + LN(RES / RN) * 298.15) -273.15;
END_IF;


(* revision history

hm 30. dec. 2008	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	changed real constants to use dot syntax

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION TEMP_PT : REAL
VAR_INPUT
	Res : REAL;
	R0 : REAL;
END_VAR
VAR CONSTANT
	A : REAL := 3.9083E-3 ;
	B : REAL := -5.775E-7;
	accuracy : REAL := 0.01;
END_VAR
VAR
	step: REAL := 50.0;
	X: REAL;
	Y: REAL;
	t1: REAL;
	pt : POINTER TO DWORD;
END_VAR

(*
version 1.7	11. mar. 2009
programmer 	hugo
tested by		oscat

this function returs the temperature for a platinum sensor for a range from -200..+850 °C *)
(* @END_DECLARATION := '0' *)
X := A * R0;
Y := B * R0;
IF Res >= R0 THEN
	t1 := X * X - 4.0 * Y * (R0 - Res);
	IF t1 < 0.0 THEN
		TEMP_PT := 10000.0;
	ELSE
		TEMP_PT := (-X + SQRT(t1)) / (2.0 * Y);
	END_IF;
ELSE
	pt := ADR(step);
	(* since the formula cannot be solved this is a successive approximation *)
	TEMP_PT := -100.0;
	WHILE step > accuracy DO
		(* test if result greater or less *)
		IF RES_PT(TEMP_PT, R0) < res THEN TEMP_PT := TEMP_PT + step; ELSE TEMP_PT := TEMP_PT - step; END_IF;
		pt^ := pt^ - 8388608;  (* this is a super fast divide by 2 method for non floating point CPUs *)
		(* the alternative code: step := step * 0.5; *)
	END_WHILE;
END_IF;

(* revision history

rev 1.1 hm 24.1.2007	
	deleted unused variable C

rev 1.2 hm 10.9.2007
	reduced accuracy to 0.02 to shorten execution time

rev 1.3	hm	2. dec 2007
	changed code for better performance

rev	1.4	hm	23. dec 2007
	avoid a negative square root if input values are wrong

rev 1.5 hm	5. jan 2008
	replaced / 2 with * 0.5 for better performance

hm	31. oct. 2008	rev 1.6
	improved performance

hm	11. mar. 2009	rev 1.7
	changed real constants to use dot syntax

*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/sensor' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION TEMP_SI : REAL
VAR_INPUT
	Res : REAL;
	RS : REAL;
	TS : REAL;
END_VAR


(*
version 1.5	10. mar. 2009
programmer 	hugo
tested by		tobias

this function returs the temperature for a silicon sensor for a range from -60..+180 °C
for example KTY10 RS at TS normally 25 °C
*)
(* @END_DECLARATION := '0' *)
TEMP_SI := (-7.64E-3 + SQRT(Res / RS * 6.64E-5 - 0.803E-5)) * 30120.48193 + TS;



(* revision history

hm 24.1.2007		rev 1.1
	the function would only work for TS = 25, now it works for other values of ts.

hm 10.9.2007		rev 1.2
	changed the function to use successive approcimation to avoid an error generated by a sqrt from a negative result.
	changed accuracy to 0.02 degrees to shorten execution time

hm 17. 12 2007	rev 1.3
	improovements for better performance

hm	6. jan 2008	rev 1.4
	further performance improvement

hm	10. mar. 2009	rev 1.5
	removed nested comments

*)
END_FUNCTION