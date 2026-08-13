FUNCTION BAND_B : BYTE
VAR_INPUT
	X : BYTE;
	B : BYTE;
END_VAR



(*
version 1.0	21. Nov. 2008
programmer 	hugo
tested by	oscat

BAND_B will limit X to B <= X <= 255-B. while X < B the resulkt will be 0 and while X > 255-B the output will be 255
otherwise the result = X.

*)
IF X < B THEN
	BAND_B := 0;
ELSIF X > BYTE#255-B THEN
	BAND_B := 255;
ELSE
	BAND_B := X;
END_IF;



(* revision history
hm	21. nov. 2008	rev 1.0
	original version

*)

END_FUNCTION


FUNCTION_BLOCK CONTROL_SET1
VAR_INPUT
	Kt : REAL;
	Tt : REAL;
	PI : BOOL;
	PID : BOOL;
END_VAR
VAR_INPUT CONSTANT
	P_K : REAL := 0.5;
	PI_K : REAL := 0.45;
	PI_TN : REAL := 0.83;
	PID_K : REAL := 0.6;
	PID_TN : REAL := 0.5;
	PID_TV : REAL := 0.125;
END_VAR
VAR_OUTPUT
	KP : REAL;
	TN : REAL;
	TV : REAL;
	KI : REAL;
	KD : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

takahashi calculates controller parameters for P, PI and PID controllers based on the ziegler nichols method.

*)
IF pi AND PID THEN
	KP := 0.0;
	TN := 0.0;
	TV := 0.0;
ELSIF PID THEN
	KP := PID_K * Kt;
	TN := PID_TN * Tt;
	TV := PID_TV * Tt;
ELSIF PI THEN
	KP := PI_K * Kt;
	TN := PI_TN * Tt;
ELSE
	KP := P_K * Kt;
END_IF;

(* KI and KD are calculated *)
IF tn > 0.0 THEN KI := KP / TN; ELSE KI := 0.0; END_IF;
KD := KP * TV;

(* revision history
hm	4. nov 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																											

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK CONTROL_SET2
VAR_INPUT
	KS : REAL;
	TU : REAL;
	TG : REAL;
	PI : BOOL;
	PID : BOOL;
END_VAR
VAR_INPUT CONSTANT
	P_K : REAL := 1.0;
	PI_K : REAL := 0.9;
	PI_TN : REAL := 3.33;
	PID_K : REAL := 1.2;
	PID_TN : REAL := 2.0;
	PID_TV : REAL := 0.5;
END_VAR
VAR_OUTPUT
	KP : REAL;
	TN : REAL;
	TV : REAL;
	KI : REAL;
	KD : REAL;
END_VAR
VAR
	TX : REAL;
END_VAR

(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

takahashi calculates controller parameters for P, PI and PID controllers based on the ziegler nichols method.

*)
IF TU > 0.0 AND KS > 0.0 THEN TX := TG / TU / KS; END_IF;
IF pi AND PID THEN
	KP := 0.0;
	TN := 0.0;
	TV := 0.0;
ELSIF PID THEN
	KP := PID_K * TX;
	TN := PID_TN * TU;
	TV := PID_TV * TU;
ELSIF PI THEN
	KP := PI_K * TX;
	TN := PI_TN * TU;
ELSE
	KP := P_K * TX;
END_IF;

(* KI and KD are calculated *)
IF TN > 0.0 THEN KI := KP / TN; ELSE KI := 0.0; END_IF;
KD := KP * TV;

(* revision history
hm	4. nov 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																																											

*)
END_FUNCTION_BLOCK


FUNCTION CTRL_IN : REAL
VAR_INPUT
	SET_POINT, ACTUAL, NOISE : REAL;
END_VAR


(*
version 1.0	2. jun 2008
programmer 	hugo
tested by		tobias


*)
(* calculate the process error DIFF *)
CTRL_IN := DEAD_ZONE(SET_POINT - ACTUAL, NOISE);


(* revision history
hm 	2. jun. 2008 	rev 1.0
	original version

*)
END_FUNCTION


FUNCTION_BLOCK CTRL_OUT
VAR_INPUT
	CI, OFFSET, MAN_IN, LIM_L, LIM_H : REAL;
	MANUAL : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	LIM : BOOL;
END_VAR


(*
version 1.1	5. nov 2008
programmer 	hugo
tested by	oscat


*)
Y := SEL(MANUAL, CI, MAN_IN) + OFFSET;

(* Limit the output *)
IF Y > LIM_L AND Y < LIM_H THEN
	LIM := FALSE;
ELSE
	Y := LIMIT(LIM_L, Y, LIM_H);
	LIM := TRUE;
END_IF;



(* revision history
hm 	2. jun. 2008 	rev 1.0
	original version

hm	5. nov. 2008	rev 1.1
	optimized code
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK CTRL_PI
VAR_INPUT
	ACT, SET, SUP, OFS, M_I : REAL;
	MAN : BOOL;
	RST : BOOL;
	KP : REAL := 1.0;
	KI : REAL := 1.0;
	LL : REAL := -1000.0;
	LH : REAL := 1000.0;
END_VAR
VAR_OUTPUT
	Y : REAL;
	DIFF : REAL;
	LIM : BOOL;
END_VAR
VAR
	pi : FT_PIWL;
	co : CTRL_OUT;
END_VAR

(*
version 2.0	30. jun 2008
programmer 	hugo
tested by	oscat

FT_PI is a PI controller with manual functionality.
The PID controller works according to the fomula Y = e *(KP+ KI * INTEG(e) ) + offset, while e = set_point - actual.
a rst will reset all internal data, while a switch to manual will cause the controller to follow the function Y = manual_in + offset.
limit_h and Limit_l set the possible output range of Y.
the output flags lim will signal that the output limits are active and overflow will signal that the integrator has reached its limits.

since rev 1.1 the "trapezregel is used for more accuracy.
rev 1.2 added selective integratin which means the integrative component is only active within a small range of the target value 
this avoids the integrator to go to limits while an input setpoint change happened and is only causing overshoots.
the int_band is by default 100 which means the int is active all the time and if set to for example to 0.1 the integrator is only active
while the input is between 0.9 and 1.1 of the set_point value.

default values for KP = 1, TN = 1, TV = 1, LIMIT_L = -1000, LIMIT_H = +1000.
*)
DIFF := CTRL_IN(SET, ACT, SUP);
pi(in := DIFF, kp := KP, ki := KI, lim_l := LL, lim_h := LH, rst := RST);
co(ci := pi.Y, OFFSET := OFS, man_in := M_I, lim_l := LL, lim_h := LH, MANUAL := MAN);
Y := co.Y;
LIM := co.LIM;


(* revision history
hm 	31.10.2007 		rev 1.0
	original version

hm	3.11.2007		rev 1.1
	added noise input to filter noise
	added output diff
	set limit output false when output is within limits
	overfolw was not set correctly

hm	5. jan 2008		rev 1.2
	improved performance

hm	20. jun. 2008	rev 2.0
	rewritten using new modular approach

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK CTRL_PID
VAR_INPUT
	ACT, SET, SUP, OFS, M_I : REAL;
	MAN : BOOL;
	RST : BOOL;
	KP : REAL := 1.0;
	TN : REAL := 1.0;
	TV : REAL := 1.0;
	LL : REAL := -1000.0;
	LH : REAL := 1000.0;
END_VAR
VAR_OUTPUT
	Y : REAL;
	DIFF : REAL;
	LIM : BOOL;
END_VAR
VAR
	pid : FT_PIDWL;
	co : CTRL_OUT;
END_VAR



(*
version 2.0	30. jun. 2008
programmer 	hugo
tested by	oscat

FT_PI is a PI controller with manual functionality.
The PID controller works according to the fomula Y = e *(KP+ KI * INTEG(e) ) + offset, while e = set_point - actual.
a rst will reset all internal data, while a switch to manual will cause the controller to follow the function Y = manual_in + offset.
limit_h and Limit_l set the possible output range of Y.
the output flags lim will signal that the output limits are active.

*)
DIFF := CTRL_IN(SET, ACT, SUP);
pid(in := DIFF, kp := KP, tn := TN, tv := TV, lim_l := LL, lim_h := LH, rst := RST);
co(ci := pid.Y, OFFSET := OFS, man_in := M_I, lim_l := LL, lim_h := LH, MANUAL := MAN);
Y := co.Y;
LIM := co.LIM;



(* revision history

hm 1.12.2006	 	rev 1.1
	changed algorithm to trapezregel for higher accuracy.

hm 3.1.2007		rev 1.2
	added integ_band to select when the integrator is active.

hm	3.3.2007		rev 1.3
	added default values to inputs KP, TN, TV, LIMIT_L und LIMIT_H.

hm 31.oct 2007	rev 1.4
	total rewrite of the module to avoid failures when one of the limits is 0

hm	3.11.2007	rev 1.5
	added noise input to filter noise
	added output diff
	set limit output false when output is within limits
	overfolw was not set correctly

hm	5. jan 2008	rev 1.6
	improved code for better performance
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK CTRL_PWM
VAR_INPUT
	CI, MAN_IN : REAL;
	MANUAL : BOOL;
	F : REAL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
END_VAR
VAR
	PW : PWM_DC;
END_VAR

(*
version 1.1	21. oct. 2008
programmer 	hugo
tested by	oscat


*)
PW(F := F, DC := SEL(MANUAL,CI,MAN_IN));
Q := PW.Q;


(* revision history
hm 3. jun. 2008 	rev 1.0
	original version

hm	21. oct. 2008	rev 1.1
	optimized code

*)
END_FUNCTION_BLOCK


FUNCTION DEAD_BAND : REAL
VAR_INPUT
	X : REAL;
	L : REAL;
END_VAR


(*
version 1.2	18. jan. 2011
programmer 	hugo
tested by		oscat

DEAD_BAND ist eine lineare Übertragungsfunktion mit Totzone. Die Funktion verschiebt den positiven Teil der Kurve um +L und den negativen Teil der Kurve um -L.
DEAD_BAND = X - L wenn X > L)
DEAD_BAND = X + L wenn X < -L
DEAD_BAND = 0 wenn Abs(X) <= L

*)
IF X > L THEN
	DEAD_BAND := X - L;
ELSIF X < -L THEN
	DEAD_BAND := X + L;
ELSE
	DEAD_BAND := 0.0;
END_IF;


(* revision history
hm	2. nov. 2007	rev 1.0
	original version

hm	14. jun. 2008	rev 1.1
	improved performance

hm	18. jan. 2011	rev 1.2
	assign 0 before return
*)
END_FUNCTION


FUNCTION_BLOCK DEAD_BAND_A
VAR_INPUT
	X : REAL;
	T : TIME;
	KL : REAL := 1.0;
	LM : REAL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	L : REAL;
END_VAR
VAR
	tp1 : FT_PT1;
	tp2 : FT_PT1;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		oscat

DEAD_BAND ist eine lineare Übertragungsfunktion mit Totzone. Die Funktion verschiebt den positiven Teil der Kurve um -L und den negativen Teil der Kurve um +L.
DEAD_BAND = X wenn Abs(X) > L wenn Abs(X) > L
DEAD_BAND = 0 wenn Abs(X) <= L
Die Breite der Totzone wird automatisch aus dem Signalrauschen errechnet.

*)
(* filter the input signal *)
tp1(in := X, T:= T);

(* filter the HF portion to generate a stable L *)
tp2(in := ABS(tp1.out - X), T := MULTIME(T, 4.0));

(* now we determine L which is half the width of the dead band *)
L := MIN(KL * tp2.out, LM);

IF X > L THEN
	Y := X - L;
ELSIF X < -L THEN
	Y := X + L;
ELSE
	Y := 0.0;
END_IF;


(* revision history
hm	14. jun. 2008	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)
END_FUNCTION_BLOCK


FUNCTION DEAD_ZONE : REAL
VAR_INPUT
	X : REAL;
	L : REAL;
END_VAR


(*
version 1.2	11. mar. 2009
programmer 	hugo
tested by		oscat

dead_zone2 is a linear transfer function which follows a linear function except for x is close to 0.
Y = X if abs(x) > L otherwise its 0.

*)
IF ABS(x) > L THEN
	dead_zone := X;
ELSE
	DEAD_ZONE := 0.0;
END_IF;

(* revision history
hm	12. feb. 2007	rev 1.0
	original version

hm	14. jun. 2008	rev 1.1
	improved performance

hm	11. mar. 2009	rev 1.2
	real constants updated to new systax using dot

*)
END_FUNCTION


FUNCTION_BLOCK DEAD_ZONE2
VAR_INPUT
	X : REAL;
	L : REAL;
END_VAR
VAR_OUTPUT
	Y: REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

dead_zone2 is a linear transfer function which follows a linear function except for x is close to 0.
Y = X if abs(x) > L.
for values of 0 +/- L a hysteresis function will hold the output at + or - L.

*)
IF ABS(x) > L THEN
	Y := X;
ELSIF Y > 0.0 THEN
	Y := L;
ELSE
	Y := -L;
END_IF;



(* revision history

hm	12. feb. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)

END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_DERIV
VAR_INPUT
	in : REAL;
	K : REAL := 1.0;
	run : BOOL := TRUE;
END_VAR
VAR_OUTPUT
	out : REAL;
END_VAR
VAR
	old: REAL;
	tx: DWORD;
	last: DWORD;
	init: BOOL;
	tc: REAL;
END_VAR


(*
version 1.5	11. mar. 2009
programmer 	hugo
tested by		oscat

FT_deriv calculates the derivate over the signal "in" with Faktor "K".
a run input enables or stops the calculation, if left unconnected its true and therfore the calculation is executed.
if K is not specified the default is 1.

*)



(* read system time *)
tx := T_PLC_US();
tc := DWORD_TO_REAL(tx - last);
last := tx;

(* init on first startup *)
IF NOT init THEN
	init := TRUE;
	old := in;
ELSIF run AND tc > 0.0 THEN
	out := (in - old) / tc * 1000000.0 * K;
	old := in;
ELSE
	out := 0.0;
END_IF;



(*
hm 3.1.2007			rev 1.1
	added init code for startup
	set the default for K to 1

hm	15. sep 2007	rev 1.2
	replaced Time() with T_PLC_US for compatibility and performance reasons
	increased accuracy and work in microseconds internally

hm 29 oct 2007	rev 1.3
	prohibit calculation when tx - last = 0 to avoid division by 0 and increase accuracy on fast systems

hm	6. nov. 2008	rev 1.4
	improved performance

hm	11. mar. 2009	rev 1.5
	inproved code
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_IMP
VAR_INPUT
	in : REAL;
	T : TIME;
	K : REAL := 1.0;
END_VAR
VAR_OUTPUT
	out : REAL;
END_VAR
VAR
	t1 : FT_PT1;
END_VAR

(*
version 1.1	3 Jan 2007
programmer 	hugo
tested by	oscat

FT_IMP is 	an impulse filter (high pass filter) with the time T and factor K.
 
*)
T1(in:= in, T:=T);
out := (in - t1.out) * K;

(*
hm 3.1.2007	rev 1.1
	added factor K

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_INT
VAR_INPUT
	IN : REAL;
	K : REAL := 1;
	RUN : BOOL := TRUE;
	RST : BOOL;
	OUT_MIN : REAL := -1E37;
	OUT_MAX : REAL := 1E37;
END_VAR
VAR_OUTPUT
	OUT : REAL;
	LIM : BOOL;
END_VAR
VAR
	INTeg : INTEGRATE;
END_VAR


(*
version 2.2	10. mar. 2009
programmer 	hugo
tested by		oscat

FT_int is an integrator with input IN and factor K.
integration is only done while run = true, if run is false integration is stopped and the out remains at the last value.
if run is left unconnected, run is considered true by the function block.
the output ET will display the total integration time from last reset.
the out_min and out_max values can be set to ,imit the output value range. 
the default for K is 1.

 
*)
(* @END_DECLARATION := '0' *)
IF rst THEN
	out := 0.0;
ELSE
	integ(X := IN, E := RUN, K := K, Y := out);
END_IF;

(* limit the outputs *)
IF out >= OUT_MAX THEN
	out := out_max;
	LIM := TRUE;
ELSIF out <= out_min THEN
	out := out_min;
	lim := TRUE;
ELSE
	lim := FALSE;
END_IF;


(*
hm 13.12.2006	rev 1.1
	changed to "trapezregel" which increases accuracy
	before out := out + in * time new: out := out + (in + in_last) / 2 * time

hm 15.1.2007		rev 1.2
	added default for k to be 1.

hm	15.9.2007		rev 1.3
	replaced time() with T_PLC_US for compatibility and performance reason
	increased internal accuracy to microseconds

hm 29. oct 2007		rev 1.4
	changed code so int will not be called is time difference is 0 to increase accuracy on systems
	with cycle times below 1ms

hm	2. dec 2007	 	rev 1.5
	changed code for better performance

hm	5. jan 2008		rev 1.6
	further improvements in performance

hm	8. feb 2008		rev 1.7
	deleted limits +/- 1000 now limit is the range of real
	added variables out1 and out2 to extend the resolution
	corrected an error with elapsed time
	deleted unusfull output ET

hm	13. mar 2008	rev 1.8
	changed preset value out_min to -1e37.
	stop integrator at the limits.

hm	16. mar 2008	rev 1.9
	added type conversion to avoid warning under codesys 3.0

hm	2. jun. 2008	rev 2.0
	rewritten with clear code
	limits can now be set without run

hm	3. nov. 2008	rev 2.1
	optimized code using INTEGRATE and LIMX

hm	10. mar. 2009	rev 2.2
	removed nested comments
*)

END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_INT2
VAR_INPUT
	IN : REAL;
	K : REAL := 1.0;
	RUN : BOOL := TRUE;
	RST : BOOL;
	OUT_MIN : REAL := -1.0E38;
	OUT_MAX : REAL := 1.0E38;
END_VAR
VAR_OUTPUT
	OUT : REAL;
	LIM : BOOL;
END_VAR
VAR
	integ: INTEGRATE;
	ix : REAL;
	val : REAL2;
END_VAR


(*
version 1.2	11. mar. 2009
programmer 	hugo
tested by		oscat

FT_int2 is an integrator with input IN and factor K.
integration is only done while run = true, if run is false integration is stopped and the out remains at the last value.
if run is left unconnected, run is considered true by the function block.
the out_min and out_max values can be set to ,imit the output value range. 
the default for K is 1.
FT_INT2 has double presision accuracy (14 digits).


*)
IF RST THEN
	val := R2_SET(0.0);
	out := 0.0;
ELSE
	integ(X := IN, E := RUN, K := K, Y := ix);
	val := R2_ADD(val, ix);
	ix := 0.0;
	OUT := val.RX;
END_IF;

(* check output for limits *)
IF out > OUT_MIN AND out < OUT_MAX THEN
	LIM := FALSE;
ELSE
	OUT := LIMIT(OUT_MIN, OUT, OUT_MAX);
	val := R2_SET(OUT);
	LIM := TRUE;
END_IF;



(*	revision history
hm	2. jun. 2008	rev 1.0
	original version

hm	5. nov. 2008	rev 1.1
	rewritten code using integrate

hm	11. mar. 2009	rev 1.2
	real constants updated to new systax using dot

*)

END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PD
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	TV : REAL := 1.0;
END_VAR
VAR_OUTPUT
	Y : REAL;
END_VAR
VAR
	diff : FT_DERIV;
END_VAR

(*
version 1.0	3. jun 2008
programmer 	hugo
tested by		tobias

FT_PD is a PD controller.
The PD controller works according to the fomula Y = KP *(IN + DERIV(e) ).

*)
(* run differentiator *)
diff(IN := IN, K := TV);

(* combine both values *)
Y := KP * (diff.out + IN);



(* revision history
hm 	3. jun. 2008 	rev 1.0
	original version


*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PDT1
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	TV : REAL := 1.0;
	T1 : REAL := 1.0;
END_VAR
VAR_OUTPUT
	Y : REAL;
END_VAR
VAR
	diff : FT_DERIV;
	TP : FT_PT1;
END_VAR

(*
version 1.0	3. jun 2008
programmer 	hugo
tested by		tobias

FT_PD is a PD controller.
The PD controller works according to the fomula Y = KP *(IN + DERIV(e) ).

*)
(* run differentiator *)
diff(IN := IN, K := TV);

(* Run PT1 filter *)
tp(in := diff.out, T := REAL_TO_TIME(T1));

(* combine both values *)
Y := KP * (tp.out + IN);



(* revision history
hm 	3. jun. 2008 	rev 1.0
	original version


*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PI
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	KI : REAL := 1.0;
	ILIM_L : REAL := -1E38;
	ILIM_H : REAL := 1E38;
	IEN : BOOL := TRUE;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	LIM : BOOL;
END_VAR
VAR
	integ : FT_INT;
END_VAR

(*
version 2.0	3. jun 2008
programmer 	hugo
tested by		tobias

FT_PI is a PI controller.
The PID controller works according to the fomula Y = IN *(KP+ KI * INTEG(e) ).
a rst will reset the integrator to 0
ilim_h and iLim_l set the possible output range of the internal integrator.
the output flags lim will signal that the output limits are active.

default values for KP = 1, KI = 1, ILIM_L = -1E37, iLIM_H = +1E38.
*)
(* run integrator *)
integ(IN := IN, K := KI, RUN := IEN, RST := RST, OUT_MIN := ILIM_L, OUT_MAX := ILIM_H);

(* check if integrator has reached its limits and set overflow *)
LIM := integ.LIM;
Y := KP * IN + integ.Out;



(* revision history
hm 	3. jun. 2008 	rev 2.0
	original version


*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PID
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	TN : REAL := 1.0;
	TV : REAL := 1.0;
	ILIM_L : REAL := -1.0E38;
	ILIM_H : REAL := 1.0E38;
	IEN : BOOL := TRUE;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	LIM : BOOL;
END_VAR
VAR
	integ : FT_INT;
	diff : FT_DERIV;
END_VAR

(*
version 2.1	11. mar. 2009
programmer 	hugo
tested by		tobias

FT_PI is a PI controller.
The PID controller works according to the fomula Y = IN *(KP+ KI * INTEG(e) ).
a rst will reset the integrator to 0
ilim_h and iLim_l set the possible output range of the internal integrator.
the output flags lim will signal that the output limits are active.

default values for KP = 1, KI = 1, ILIM_L = -1E37, iLIM_H = +1E38.
*)
(* run integrator only if TN > 0 *)
IF TN > 0.0 THEN
	integ(IN := IN, K := 1.0 / TN, RUN := IEN, RST := RST, OUT_MIN := ILIM_L, OUT_MAX := ILIM_H);
ELSE
	integ(RST := FALSE);
END_IF;

(* run differentiator *)
diff(IN := IN, K := TV);

(* combine both values *)
Y := KP * (integ.Out + diff.out + IN);

(* check if integrator has reached its limits and set overflow *)
LIM := integ.LIM;



(* revision history
hm 3. jun. 2008 	rev 2.0
	original version

hm	11. mar. 2009	rev 2.1
	real constants updated to new systax using dot

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PIDW
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	TN : REAL := 1.0;
	TV : REAL := 1.0;
	LIM_L : REAL := -1.0E38;
	LIM_H : REAL := 1.0E38;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	LIM : BOOL;
END_VAR
VAR
	integ : INTEGRATE;
	diff : FT_DERIV;
	YI: REAL;
END_VAR

(*
version 1.2	11. mar. 2009
programmer 	hugo
tested by		oscat

FT_PIDW is a PID controller with dynamic wind_up reset.
The PID controller works according to the fomula Y = KP *(IN + KI * INTEG(e) + DERIV(e) ).
a rst will reset the integrator to 0
lim_h and Lim_l set the possible output range of the internal integrator.
the output flags lim will signal that the output limits are active.

default values for KP = 1, KI = 1, ILIM_L = -1E37, iLIM_H = +1E38.
*)
(* run the integrator *)
IF tn = 0.0 OR rst THEN
	integ(E := FALSE, Y := YI);
	YI := 0.0;
ELSE
	integ(X := IN, K := 1.0 / TN, E := NOT LIM, Y := YI);
END_IF;

(* add up integrator and linear part *)
Y := KP * (IN + YI);

(* run differentiator *)
diff(IN := IN, K := TV);

(* set lim before differentiator is added to stop integrator if necessary *)
IF Y > LIM_L AND Y < LIM_H THEN
	LIM := FALSE;
ELSE
	LIM := TRUE;
END_IF;

(* add differential part and limit output Y *)

Y := LIMIT(LIM_L, Y + KP * diff.out, LIM_H);




(* revision history
hm 3. jun. 2008 	rev 1.0
	original version

hm	5. nov. 2008	rev 1.1
	changed code to use integrate

hm	11. mar. 2009	rev 1.2
	real constants updated to new systax using dot

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PIDWL
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	TN : REAL := 1.0;
	TV : REAL := 1.0;
	LIM_L : REAL := -1.0E38;
	LIM_H : REAL := 1.0E38;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	LIM : BOOL;
END_VAR
VAR
	piwl : FT_PIWL;
	diff : FT_DERIV;
END_VAR

(*
version 1.3	13. nov. 2009
programmer 	hugo
tested by		tobias

FT_PIDWL is a PID controller with dynamic wind_up reset.
The PID controller works according to the fomula Y = KP *(IN + KI * INTEG(e) + DERIV(e) ).
a rst will reset the integrator to 0
lim_h and Lim_l set the possible output range of the internal integrator.
the output flags lim will signal that the output limits are active.

default values for KP = 1, KI = 1, ILIM_L = -1E37, iLIM_H = +1E38.
*)
(* if rst then *)
IF rst THEN
	piwl(rst := TRUE);
	piwl.RST := FALSE;
ELSE
	(* run PIWL controller first *)
	(* we need to check if TN = 0 and do alternative calls *)
	IF TN = 0.0 THEN
		piwl(in := IN * KP, KP := 1.0, KI := 0.0, LIM_L := LIM_L, LIM_H := LIM_H);
	ELSE
		piwl(in := IN * KP, KP := 1.0, KI := 1.0 / TN, LIM_L := LIM_L, LIM_H := LIM_H);
	END_IF;

	(* run differentiator and add_to_output *)
	diff(IN := IN, K := KP * TV);
	Y := piwl.Y + diff.out;

	(* limit the output *)
	IF Y < LIM_L THEN
		LIM := TRUE;
		Y := LIM_L;
	ELSIF Y > LIM_H THEN
		LIM := TRUE;
		Y := LIM_H;
	ELSE
		LIM := FALSE;
	END_IF;
END_IF;



(* revision history
hm 13. jun. 2008 	rev 1.0
	original version

hm	25. jan 2008	rev 1.1
	multiply differential part with KP

hm	11. mar. 2009	rev 1.2
	real constants updated to new systax using dot

hm	13. nov. 2009	rev 1.3
	fixed code for negative KP

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PIW
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	KI : REAL := 1.0;
	LIM_L : REAL := -1E38;
	LIM_H : REAL := 1E38;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	LIM : BOOL;
END_VAR
VAR
	integ : FT_INT;
END_VAR

(*
version 1.0	3. jun 2008
programmer 	hugo
tested by		tobias

FT_PIW is a PI controller.
The PID controller works according to the fomula Y = IN *(KP+ KI * INTEG(e) ).
a rst will reset the integrator to 0
ilim_h and iLim_l set the possible output range of the internal integrator.
the output flag lim will signal that the output limits are active.
the controller is equipped with anti wind_up circuitry that stops the integrator when lim_h or lim_l is reached


default values for KP = 1, KI = 1, ILIM_L = -1E37, iLIM_H = +1E38.
*)
(* run integrator *)
integ(IN := IN, K := KI, RUN := NOT LIM, RST := RST);

(* set output value *)
Y := KP * IN + integ.Out;

(* check for limits and set integrator for anti wind up *)
IF Y < LIM_L THEN
	Y := LIM_L;
	LIM := TRUE;
ELSIF Y > LIM_H THEN
	Y := LIM_H;
	LIM := TRUE;
ELSE
	LIM := FALSE;
END_IF;



(* revision history
hm 	3. jun. 2008 	rev 1.0
	original version


*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PIWL
VAR_INPUT
	IN : REAL;
	KP : REAL := 1.0;
	KI : REAL := 1.0;
	LIM_L : REAL := -1.0E38;
	LIM_H : REAL := 1.0E38;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	LIM : BOOL;
END_VAR
VAR
	init: BOOL;
	tx: DWORD;
	tc : REAL;
	t_last: DWORD;
	in_last : REAL;
	i: REAL;
	p: REAL;
END_VAR

(*
version 1.3	11. mar. 2009
programmer 	hugo
tested by		oscat

FT_PIWL is a PI controller.
The PID controller works according to the fomula Y = IN *(KP+ KI * INTEG(e) ).
a rst will reset the integrator to 0
lim_h and Lim_l set the possible output range of the controller.
the output flag lim will signal that the output limits are active.
the integrator ist equipped with anti wind-up circuitry which limits trhe total output ranke to lim_l and lim_h

default values for KP = 1, KI = 1, ILIM_L = -1E37, iLIM_H = +1E38.
*)
(* initialize at power_up *)
IF NOT init OR RST THEN
	init := TRUE;
	in_last := in;
	t_last := T_PLC_US();
	i := 0.0;
	tc := 0.0;
ELSE
	(* read last cycle time in Microseconds *)
	tx := T_PLC_US();
	tc := DWORD_TO_REAL(tx - t_last);
	t_last := tx;

	(* calculate proportional part *)
	p := KP * IN;

	(* run integrator *)
	i := (IN + in_last) * 5.0E-7 * KI * tc + i;
	in_last := IN;

	(* calculate output Y *)
	Y := p + i;

	(* check output for limits *)
	IF Y >= LIM_H THEN
		Y := LIM_H;
		IF ki <> 0.0 THEN
			i := LIM_H - p;
		ELSE
			i := 0.0;
		END_IF;
		LIM := TRUE;
	ELSIF Y <= LIM_L THEN
		Y := LIM_L;
		IF ki <> 0.0 THEN
			i := LIM_L - p;
		ELSE
			i := 0.0;
		END_IF;
		LIM := TRUE;
	ELSE
		LIM := FALSE;
	END_IF;
END_IF;




(* revision history
hm 13. jun. 2008 	rev 1.0
	original version

hm	27. oct. 2008	rev 1.1
	integrator will not be adjusted when ki = 0

hm	25. jan 2009	rev 1.2
	module will also work with negative K

hm	11. mar. 2009	rev 1.3
	real constants updated to new systax using dot

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PT1
VAR_INPUT
	in : REAL;
	T : TIME;
	K : REAL := 1.0;
END_VAR
VAR_OUTPUT
	out : REAL;
END_VAR
VAR
	last : DWORD;
	tx: DWORD;
	init: BOOL;
END_VAR

(*
version 1.11	18. jan. 2011
programmer 	hugo
tested by		oscat

FT_PT1 is an low pass filter with a programmable time T and faktor K.
 
*)
(* read system time *)
tx := T_PLC_US();

(* startup initialisation *)
IF NOT init OR T = t#0s THEN
	init := TRUE;
	out := K * in;
ELSE
	out := out + (in * K - out) * DWORD_TO_REAL(Tx - last) / TIME_TO_REAL(T) * 1.0E-3;
	IF ABS(out) < 1.0E-20 THEN out := 0.0; END_IF;
END_IF;
last := tx;


(*
hm 1.1.2007	rev 1.1
	corrected error while startup value was not correct
	for very small time values real output would run out of range.

hm 3.1.2007	rev 1.2
	corrected an error for falling edge failures.
	added output faktor K.

hm 27. 2. 2007	rev 1.3
	output will be input during init for definitive startup condition.

hm	15.9.2007	rev 1.4
	changed time() to T_PLC_US() for compatibilitxy resons
	increased internal accuracy to Microseconds instead of Milliseconds

hm	23. oct 2007	rev 1.5
	added out := in to the init statements

hm	30. nov 2007	rev 1.6
	changed out to be K * in during initialization

hm	5. jan 2008	rev 1.7
	improved code for better performance

hm	16. mar. 2008	rev 1.8
	added type conversion to avoid warning under codesys 3.0

hm	14. jun. 2008	rev 1.9
	improved code

hm	11. mar. 2009	rev 1.10
	real constants updated to new systax using dot

hm	18. jan. 2011	rev 1.11
	avoid underrun of out
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_PT2
VAR_INPUT
	in : REAL;
	T : TIME;
	D : REAL;
	K : REAL := 1.0;
END_VAR
VAR_OUTPUT
	out : REAL;
END_VAR
VAR
	init: BOOL;
	int1 : INTEGRATE;
	int2 : INTEGRATE;
	tn: REAL;
	I1, I2 : REAL;
	tn2: REAL;
END_VAR

(*
version 1.5	11. mar. 2009
programmer 	hugo
tested by		oscat

FT_PT2 is a 2nd grade filter with programmable times T, D and faktor K.
 
*)
(* startup initialisation *)
IF NOT init OR T = T#0s THEN
	init := TRUE;
	out := K * in;
	I2 := out;
ELSE
	TN := TIME_TO_REAL(T) * 1.0E-3;
	tn2 := TN * TN;
	int1(X := in * K / tn2 - I1 * 0.5 * D / TN - I2 / TN2, Y := I1);
	int2(X := I1,Y := I2);
	out := I2;
END_IF;


(* revision history

15.1.2007 hm		rev 1.1
	changed formula to new more acurate formula

hm 15.9.2007		rev 1.2
	deleted unused code for init system time reading tx	

hm	30.11.2007	rev 1.3
	changed out to be K * in during initialization
	avoind divide by 0 if tn = 0

hm	3. nov. 2008	rev 1.4
	optimized code and fixed a problem with init

hm	11. mar. 2009	rev 1.5
	real constants updated to new systax using dot

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_TN16
VAR_INPUT
	in : REAL;
	T : TIME;
END_VAR
VAR_OUTPUT
	out : REAL;
	trig : BOOL;
END_VAR
VAR
	length : INT := 16;
	X : ARRAY[0..15] OF REAL;
	cnt : INT;
	last : TIME;
	tx: TIME;
	init: BOOL;
END_VAR

(*
version 1.1	16 sep 2007
programmer 	hugo
tested by		tobias

FT_TN16 is delay function, it will delay a signal by a specified time : T and will store 16 values of in before they are put thru to out.
for higher or lower resolution please use FT_TN8 or FT_TN64 instead.
 
*)
(* read system time *)
tx := DWORD_TO_TIME(T_PLC_MS());

trig := FALSE;
IF NOT init THEN
	x[cnt] := in;
	init := TRUE;
	last := tx;
ELSIF tx - last >= T / length THEN
	IF cnt = length - 1 THEN cnt := 0; ELSE cnt := cnt + 1; END_IF;
	Out := X[cnt];
	x[cnt] := in;
	last := tx;
	Trig := TRUE;
END_IF;

(* revision history
hm		1. jan 2007	rev 1.0
	original version

hm		16. sep 2007	rev 1.1
	changes time() to T_plc_ms() for compatibility reasons

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_TN64
VAR_INPUT
	in : REAL;
	T : TIME;
END_VAR
VAR_OUTPUT
	out : REAL;
	trig: BOOL;
END_VAR
VAR
	length : INT := 64;
	X : ARRAY[0..63] OF REAL;
	cnt : INT;
	last : TIME;
	tx: TIME;
	init: BOOL;
END_VAR

(*
version 1.1	15 sep 2007
programmer 	hugo
tested by		tobias

FT_TN7 is delay function, it will delay a signal by a specified time : T and will store 64 values of in before they are put thru to out.
if lower resolution is needed, pls use FT_TN8 or FT_TN16 instead.
 
*)
(* read system time *)
tx := DWORD_TO_TIME(T_PLC_MS());

trig := FALSE;
IF NOT init THEN
	x[cnt] := in;
	init := TRUE;
	last := tx;
ELSIF tx - last >= T / length THEN
	IF cnt = length - 1 THEN cnt := 0; ELSE cnt := cnt + 1; END_IF;
	Out := X[cnt];
	x[cnt] := in;
	last := tx;
	trig := TRUE;
END_IF;

(* revision history
hm		1. jan 2007		rev 1.0
	original version

hm		16. sep 2007	rev 1.1
	changes time() to T_plc_ms() for compatibility reasons

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_TN8
VAR_INPUT
	in : REAL;
	T : TIME;
END_VAR
VAR_OUTPUT
	out : REAL;
	trig: BOOL;
END_VAR
VAR
	length : INT := 8;
	X : ARRAY[0..7] OF REAL;
	cnt : INT;
	last : TIME;
	tx: TIME;
	init: BOOL;
END_VAR

(*
version 1.1	15 Sep 2007
programmer 	hugo
tested by		tobias

FT_TN8 is delay function, it will delay a signal by a specified time : T and will store 8 values of in before they are put thru to out.
if higher resolution is needed, pls use FT_TN16 or FT_TN64 instead.
 
*)
tx := DWORD_TO_TIME(T_PLC_MS());
trig := FALSE;
IF NOT init THEN
	x[cnt] := in;
	init := TRUE;
	last := tx;
ELSIF tx - last >= T / length THEN
	IF cnt = length - 1 THEN cnt := 0; ELSE cnt := cnt + 1; END_IF;
	Out := X[cnt];
	x[cnt] := in;
	last := tx;
	trig := TRUE;
END_IF;

(* revision history
hm		1. jan 2007		rev 1.0
	original version

hm		16. sep 2007	rev 1.1
	changes time() to T_plc_ms() for compatibility reasons

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK HYST
VAR_INPUT
	In : REAL;
	ON : REAL;
	OFF : REAL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	win : BOOL;
END_VAR


(*
version 1.0	2. jun. 2008
programmer 	hugo
tested BY	oscat

This Hystereses function has two modes:
1. if on > off then Q will be switched high when in > on and switched low when in < off.
2. if on < off then Q will be switched high when in < on and switched low when in > off.
the output win will be high when in is between low and high.

*)
IF ON >= OFF THEN
	IF IN < OFF THEN
		Q := FALSE;
		WIN := FALSE;
	ELSIF  IN > ON THEN
		Q := TRUE;
		WIN := FALSE;
	ELSE
		WIN := TRUE;
	END_IF;
ELSE
	IF IN > OFF THEN
		Q := FALSE;
		WIN := FALSE;
	ELSIF  IN < ON THEN
		Q := TRUE;
		WIN := FALSE;
	ELSE
		WIN := TRUE;
	END_IF;
END_IF;


(* revision history
hm		2.  jun 2008	rev 1.0
	original version

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK HYST_1
VAR_INPUT
	In : REAL;
	high : REAL;
	low : REAL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	win : BOOL;
END_VAR


(*
version 1.1	2. jun. 2008
programmer 	hugo
tested BY	oscat

this hysteresis function switches the output high if the input signal reaches obove high and will switch to low when the input falls back below low value.
a separate output mid is set if the input stays between low and high value.

*)
IF in < low THEN
	Q := FALSE;
	win := FALSE;
ELSIF in > high THEN
	Q := TRUE;
	win := FALSE;
ELSE
	win := TRUE;
END_IF;




(* revision history
hm		4.  aug. 2006	rev 1.0
	original version

hm		2. jun. 2008	rev 1.1
	improved performance
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK HYST_2
VAR_INPUT
	IN : REAL;
	VAL : REAL;
	HYS : REAL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	WIN: BOOL;
END_VAR
VAR
	tmp: REAL;
END_VAR


(*
version 1.2	2. jun. 2008
programmer 	oscat
tested BY	oscat

this hysteresis function switches the output high if the input signal reaches obove val + hys/2 and will switch to low when the input falls back below val - hys/2 value.
a separate output mid is set if the input stays between low and high value.

*)
tmp := val - hys * 0.5;
IF in < tmp THEN
	Q := FALSE;
	win := FALSE;
ELSIF in > tmp + hys THEN
	Q := TRUE;
	win := FALSE;
ELSE
	win := TRUE;
END_IF;



(* revision history
hm		4. aug 2006	rev 1.0
	original version

hm		5. jan 2008	rev 1.1
	improved code for better performance

hm		2. jun. 2008	rev 1.2
	improved performance
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK HYST_3
VAR_INPUT
	in : REAL;
	hyst : REAL;
	val1 : REAL;
	val2 : REAL;
END_VAR
VAR_OUTPUT
	Q1 : BOOL;
	Q2 : BOOL;
END_VAR
VAR
	X: REAL;
END_VAR

(*
version 1.2	5 jan 2008
programmer 	oscat
tested BY	oscat

this is a double hysteresis function. Out1 follows a hysteresis function defined by val1and hyst, while out 2 follows val2 and hyst.
if the input signal is between the two hysteresis switches (val1 and val2) then non of the outputs is active.

*)
X := hyst * 0.5;
IF in < val1 - X THEN
	q1 := TRUE;
ELSIF in > val1 + X THEN
	q1 := FALSE;
END_IF;
IF in < val2 - X THEN
	q2 := FALSE;
ELSIF in > val2 + X THEN
	q2 := TRUE;
END_IF;



(* revision history
hm	22. jan 2007	rev 1.0
	original version

hm	27. dec 2007	rev 1.1
	changed code for better performance

hm	5. jan 2008	rev 1.2
	further performance iprovements

*)

END_FUNCTION_BLOCK


FUNCTION_BLOCK INTEGRATE
VAR_INPUT
	E : BOOL := TRUE;
	X : REAL;
	K : REAL := 1.0;
END_VAR
VAR_IN_OUT
	Y : REAL;
END_VAR
VAR
	X_last : REAL;
	init: BOOL;
	last: DWORD;
	tx: DWORD;
END_VAR


(*
version 1.0	3. nov. 2008
programmer 	hugo
tested by	oscat

integrate is a plain integrator with I/O for out.

 
*)
(*read system time *)
tx := T_PLC_MS();

IF NOT init THEN
	init := TRUE;
	X_last := X;
ELSIF E THEN
	Y := (X + X_LAST) * 0.5E-3 * DWORD_TO_REAL(tx-last) * K + Y;
	X_last := X;
END_IF;
last := tx;



(*
hm 3. nov. 2008	rev 1.0
original version
	
*)

END_FUNCTION_BLOCK