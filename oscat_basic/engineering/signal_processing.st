FUNCTION AIN : REAL
VAR_INPUT
	in : DWORD;
END_VAR
VAR_INPUT CONSTANT
	Bits : BYTE;
	sign : BYTE := 255;
	low : REAL;
	high : REAL := 10.0;
END_VAR
VAR CONSTANT
	ff : DWORD := 16#FFFFFFFF;
END_VAR
VAR
	temp1 : DWORD;
	temp2 : DWORD;
	sx : BOOL;
END_VAR

(*
version 1.5	16 mar 2008
programmer 	oscat
tested by		tobias

Ain converts signals from A/D converters or other digital sources to an internal real value
the lowest number of bits are extracted from the input word and the sign will be extracted if available separately.
the output signal is then conditioned to range from low to high values for a 0 to max value on the analog inputs:
for example a 15bit input (bits := 12) with sign at bit 15 (0..15) will deliver 0.0 (low value at 0) for an input value of 2#0000_0000_0000
an input value of 2#1111_1111_1111 will deliver 10.0 on the output (high value set to 10).

*) 
(* extract the sign bit *)
IF sign < 32 THEN
	temp1 := SHR(in,sign);
	sx := temp1.0;
ELSE
	sx := FALSE;
END_IF;
temp1 := SHR(ff, 32-bits);
temp2 := in AND temp1;
AIN := (high - low) * DWORD_TO_REAL(temp2) / DWORD_TO_REAL(temp1) + low;
IF sx THEN AIN := -AIN; END_IF;

(* revision history
hm 18.8.2006	rev 1.1
	fixed an error with low value negative and high value 0.

hm 19.1.2007	rev 1.2
	fixed an error with sign bit.

hm	13.9.2007	rev 1.3
	changed code to avoid warning under codesys 2.8.1

hm	2. dec 2007	rev 1.4
	changed code for better performance

hm	16. mar 2008	rev 1.5
	added type conversions to avoid warnings under codesys 30
*)
END_FUNCTION


FUNCTION_BLOCK AIN1
VAR_INPUT
	in : DWORD;
END_VAR
VAR_INPUT CONSTANT
	sign_bit : INT := 255;
	error_bit : INT := 255;
	error_code_en : BOOL;
	error_code : DWORD;
	overflow_bit : INT := 255;
	overflow_code_en : BOOL;
	overflow_code : DWORD;
	Bit_0 : INT;
	Bit_N : INT := 31;
	out_min : REAL;
	out_max : REAL := 10.0;
	code_min : DWORD;
	code_max : DWORD := 16#FFFFFFFF;
	error_output : REAL;
	overflow_output : REAL := 10.0;
END_VAR
VAR_OUTPUT
	out : REAL;
	sign : BOOL;
	error : BOOL;
	overflow : BOOL;
END_VAR
VAR
	tB: DWORD;
END_VAR

(*
version 1.3	10. mar. 2009
programmer 	oscat
tested by		tobias

Ain1 converts signals from A/D converters or other digital sources to an internal real value.

*) 
(* extract error bit *)
error := ((SHR(in,error_bit) AND 16#0000_0001) = 1) OR (error_code_en AND error_code = in);
IF error THEN
	out := error_output;
	RETURN;
END_IF;

(* strip off the data input *)
tb := SHR(SHL(in, 31 - bit_N), 31 - bit_N + Bit_0);

(* extract overflow bit *)
overflow := ((SHR(in,overflow_bit) AND 16#0000_0001) = 1) OR (overflow_code_en AND overflow_code = in) OR (tb < code_min OR tb > code_max);
IF overflow THEN
	out := overflow_output;
	RETURN;
END_IF;

(* extract sign bit *)
sign := (SHR(in,sign_bit) AND 16#0000_0001) = 1;

(* convert in to out *)
out := (DWORD_TO_REAL(tb - code_min) * (out_max - out_min) / DWORD_TO_REAL(code_max - code_min) + out_min);
IF sign THEN out := out * -1.0; END_IF;



(* revision history
hm	23. feb 2008	rev 1.0
	original version

hm	16. mar 2008	rev 1.1
	added type conversions to avoid warnngs under codesys 30

hm	22. apr. 2008	rev 1.2
	corrected error in formula when code_min was set
	corrected error when sign bit was used
	optimized code for better performance

hm	10. mar. 2009	rev 1.3
	real constants updated to new systax using dot

*)
END_FUNCTION_BLOCK


FUNCTION AOUT : DWORD
VAR_INPUT
	in : REAL;
END_VAR
VAR_INPUT CONSTANT
	Bits : BYTE;
	sign : BYTE := 255;
	low : REAL;
	high : REAL := 10.0;
END_VAR
VAR CONSTANT
	ff : DWORD := 2#1;
END_VAR
VAR
	in2 : REAL;
	sx : BOOL;
END_VAR

(*
version 1.4	23 feb 2008
programmer 	oscat
tested by		tobias

this module conditions an internal real value for D/A converter.
Aout converts and internal real value to a word for an D/A converter or other output devices.
The input value is converted to a n-bit wide output and a sign bit is set separately as specified.
the outout min value is set for the specified min input value and the max output is set for the max input value.
an input higher or lower then the max or min value will set the respective max or min value or the output.


*) 
(* if sign bit is defined *)
IF sign < 32 THEN
	sx := sign_R(in);
	in2 := ABS(in);
ELSE
	in2 := in;
END_IF;

(* begrenze in auf die zulässigen werte *)
in2 := LIMIT(low, in2, high);

(* Berechne den Ausgangswert *)
Aout := REAL_TO_DWORD((in2 - low) / (high - low) * DWORD_TO_REAL(SHL(ff,bits) -1));
IF sx THEN Aout := SHL(ff,sign) OR Aout; END_IF;

(*
revision history

hm 18.1.2007		rev 1.1
	renamed Modul to aout.
	changed Output to 32 Bit max.
	corrected error with sign bit.

hm	13.9.2007		rev 1.2
	changed code to avoid warning under codesys 2.8.1

hm	2. dec 2007		rev 1.3
	changed code for better performance

hm	23. feb 2008	rev 1.4
	changed code for better performance
*)
END_FUNCTION


FUNCTION AOUT1 : DWORD
VAR_INPUT
	in : REAL;
END_VAR
VAR_INPUT CONSTANT
	Bit_0 : INT;
	Bit_N : INT := 31;
	sign : INT := 255;
	low : REAL;
	high : REAL := 10.0;
END_VAR
VAR CONSTANT
	ff : DWORD := 2#1;
END_VAR
VAR
	sx: BOOL;
	in2 : REAL;
END_VAR

(*
version 1.0	23 feb 2008
programmer 	oscat
tested by		tobias

this module conditions an internal real value for a D/A converter.

*) 
(* if sign bit is defined *)
IF sign < 32 THEN
	sx := SIGN_R(in);
	in2 := ABS(in);
ELSE
	in2 := in;
END_IF;

(* begrenze in auf die zulässigen werte *)
in2 := LIMIT(low, in2, high);

(* Berechne den Ausgangswert *)
AOUT1 := SHL(REAL_TO_DWORD((in2 - low) / (high - low) * DWORD_TO_REAL(SHL(ff,bit_n - Bit_0 + 1) -1)),Bit_0);
IF sx THEN AOUT1 := SHL(ff,sign) OR AOUT1; END_IF;


(* revision history
hm	23. feb 2008	rev 1.0
	original version

*)
END_FUNCTION


FUNCTION BYTE_TO_RANGE : REAL
VAR_INPUT
	X : BYTE;
	low : REAL;
	high : REAL;
END_VAR


(*
version 1.0	9. jan 2008
programmer 	hugo
tested by		tobias

Byte_to_Range converts a Byte into a real between low and high.

*)
BYTE_TO_RANGE := (high - low) * X / 255.0 + low;


(* revision history
hm	9. jan 2008		rev 1.0
	original version

*)

END_FUNCTION


FUNCTION_BLOCK DELAY
VAR_INPUT
	IN : REAL;
	N : INT;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	OUT : REAL;
END_VAR
VAR
	buf : ARRAY[0..31] OF REAL;
	i : INT;
	init: BOOL;
	stop: INT;
END_VAR

(*
version 1.5	23. mar. 2009
programmer 	hugo
tested by		oscat

this function block delays input values by each programm cycle
after N+1 cycles the in value has shifted to the out.
N can be any alue from 0 .. 32
if n = 0 the input will be present on the output without a delay.
f N > 32 then the output will be delayed by 32 cycles.
any high on rst will load the buffer with in.

*)
stop := LIMIT(0,N,32) - 1;
IF rst OR NOT init THEN
	init := TRUE;
	FOR i := 0 TO stop DO buf[i] := in; END_FOR;
	out := in;
	i := 0;
ELSIF stop < 0 THEN
	out := in;
ELSE
	out := buf[i];
	buf[i] := in;
	i := INC1(i, N);
END_IF;



(* revision history
hm 1.10.2006		rev 1.1
	corrected error in buffer management

hm 19.1.2007		rev 1.2
	changed reset to load the value of in instead of 0

hm	27. oct. 2008	rev 1.3
	improved performance

hm	23. feb.2009	rev 1.4
	corrected an index problem

hm	23. mar. 2009	rev 1.5
	corrected non standard write to input N

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK DELAY_4
VAR_INPUT
	in : REAL;
END_VAR
VAR_OUTPUT
	out1 : REAL;
	out2 : REAL;
	out3 : REAL;
	out4 : REAL;
END_VAR
VAR
	temp: REAL;
END_VAR

(*
version 1.1	19 jan 2007
programmer 	hugo
tested by		tobias

this function block delays input values by each programm cycle
after 4 cycles the in value has shifted to the out 4 and will be discarded after the next cycle
the blocks can be cascaded.

*)
out4 := out3;
out3 := out2;
out2 := out1;
out1 := temp;
temp := in;

(* revision history

hm 19.1.2007	rev 1.1
	added variable temp to add 1 delay  for q1

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FADE
VAR_INPUT
	IN1 : REAL;
	IN2 : REAL;
	F : BOOL;
	TF : TIME;
	rst : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
END_VAR
VAR
	rmx : RMP_W;
END_VAR

(*
version 1.3	24. jan. 2009
programmer 	oscat
tested BY		oscat

FADE is used to crossfade between the inputs IN1 and IN2. The fade_over time is specified with TF.
When F = TRUE then Y = IN2 after the time TF, and when F = FALSE then Y = IN1.

*)

rmx(rst := rst AND NOT F, set := rst AND F, pt := TF, up := F);
Y := (in2 - In1) / 65535.0 * WORD_TO_REAL(rmx.out) + in1;


(* code for rev 1.1
IF rst THEN
	rmx(set := F, rst := NOT F);
ELSIF F AND (NOT rmx.high) THEN
	rmx(PT := TF, UP := TRUE, e := TRUE, rst := FALSE, set := FALSE);
ELSIF (NOT F) AND (NOT rmx.low) THEN
	rmx(PT := TF, UP := FALSE, e := TRUE, rst := FALSE, set := FALSE);
ELSE
	rmx(e := FALSE, rst := FALSE, set := FALSE);
END_IF;
Y := (WORD_TO_REAL(rmx.out) * in1 + WORD_TO_REAL(FF - rmx.out) * IN2) / FF;
*)


(* revision history
hm	26. dec 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	improved performance
	changed calls for rmp_w because rmp_w has chaged

hm	17. dec. 2008	rev 1.2
	function of input f was inverted

hm	24. jan. 2009	rev 1.3
	delted unused var FF
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FILTER_DW
VAR_INPUT
	X : DWORD;
	T : TIME;
END_VAR
VAR_OUTPUT
	Y : DWORD;
END_VAR
VAR
	last : DWORD;
	tx: DWORD;
	init: BOOL;
	Yi : REAL;
END_VAR

(*
version 1.1	3. nov. 2008
programmer 	hugo
tested by		oscat

FILTER_DW is an low pass filter with a programmable time T used for DWORD format.
 
*)
(* read system time *)
tx := T_PLC_MS();

(* startup initialisation *)
IF NOT init OR T = t#0s THEN
	init := TRUE;
	Yi := DWORD_TO_REAL(X);
ELSE
	Yi := Yi + (DWORD_TO_REAL(X) - DWORD_TO_REAL(Y)) * DWORD_TO_REAL(tx - last) / TIME_TO_REAL(T);
END_IF;
last := tx;

Y := REAL_TO_DWORD(Yi);



(*
hm 10. oct. 2008	rev 1.0
	original version

hm	3. nov. 2008	REV 1.1
	corrected an overflow problem
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FILTER_I
VAR_INPUT
	X : INT;
	T : TIME;
END_VAR
VAR_OUTPUT
	Y : INT;
END_VAR
VAR
	Yi : DINT;
	last : DWORD;
	tx: DWORD;
	init: BOOL;
END_VAR

(*
version 1.0	8. nov. 2008
programmer 	hugo
tested by	oscat

FILTER_I is a low pass filter with a programmable time T used for INT format.
 
*)
(* read system time *)
tx := T_PLC_MS();

(* startup initialisation *)
IF NOT init OR T = t#0s THEN
	init := TRUE;
	Yi := INT_TO_DINT(X) * 1000;
ELSE
	(* to increase accuracy of the filter we calculate internal Yi wich is Y * 1000 *)
	Yi := Yi + INT_TO_DINT(X - Y) * DWORD_TO_DINT(tx - last) * 1000 / TIME_TO_DINT(T);
END_IF;
last := tx;
Y := DINT_TO_INT(yi / 1000);



(*
hm 8. nov. 2008	rev 1.0
	original version

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FILTER_MAV_DW
VAR_INPUT
	X : DWORD;
	N : UINT;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : DWORD;
END_VAR
VAR
	init: BOOL;
	buffer : ARRAY[0..31] OF DWORD;
	i: INT;
END_VAR
VAR_TEMP
	tmp: INT;
END_VAR

(*
version 1.3	23. feb. 2009
programmer 	hugo
tested by		oscat

FILTER_MAV_DW is a moving average filter with programmable length N for DWORD Data.
 
*)
(* limit N to size of buffer *)
N := MIN(N, 32);

(* startup initialisation *)
IF NOT init OR rst OR N = 0 THEN
	init := TRUE;
	tmp := UINT_TO_INT(N)-1;
	FOR i := 0 TO tmp DO
		buffer[i] := X;
	END_FOR;
	Y := X;
ELSE
	tmp := UINT_TO_INT(N);
	i := INC1(i, tmp);
	Y := Y + (X - buffer[i]) / N;
	buffer[i] := X;
END_IF;




(*
hm 13. oct. 2008	rev 1.0
	original version

hm	27. oct. 2008	rev 1.1
	added typecast to avoid warnings

hm	24. nov. 2008	rev 1.2
	added typecasts to avoid warnings
	avoid divide by 0 if N = 0

hm	23. feb. 2009	rev 1.3
	limit N to max array size

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FILTER_MAV_W
VAR_INPUT
	X : WORD;
	N : UINT;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : WORD;
END_VAR
VAR
	init: BOOL;
	buffer : ARRAY[1..32] OF WORD;
	i: INT;
	sum : DWORD;
END_VAR
VAR_TEMP
	tmp : INT;
END_VAR

(*
version 1.4	26. MAR. 2011
programmer 	hugo
tested by		oscat

FILTER_MAV_W is a moving average filter with programmable length N for WORD Data.
 
*)
(* limit N to size of buffer *)
N := MIN(N, 32);

(* startup initialisation *)
IF NOT init OR rst OR N = 0 THEN
	init := TRUE;
	tmp := UINT_TO_INT(N) - 1;
	FOR i := 1 TO tmp DO
		buffer[i] := X;
	END_FOR;
	sum := Y * N;
	Y := X;
ELSE
	tmp := UINT_TO_INT(N);
	i := INC1(i, tmp);
	sum := sum + X - buffer[i];
	Y := DWORD_TO_WORD(sum / N);
	buffer[i] := X;
END_IF;



(*
hm 13. oct. 2008	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	added typecast to avoid warnings

hm	24. nov. 2008	rev 1.2
	added typecasts to avoid warnings
	avoid devide by 0 if N = 0

hm	23. feb. 2009	rev 1.3
	limit N to max array size

hm	26. mar. 2011	rev 1.4
	corrected error in calculation
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FILTER_W
VAR_INPUT
	X : WORD;
	T : TIME;
END_VAR
VAR_OUTPUT
	Y : WORD;
END_VAR
VAR
	last : DWORD;
	tx: DWORD;
	init: BOOL;
	tmp: DWORD;
END_VAR

(*
version 1.2	25. jan. 2011
programmer 	hugo
tested by	oscat

FILTER_W is an low pass filter with a programmable time T used for WORD format.
 
*)
(* read system time *)
tx := T_PLC_MS();

(* startup initialisation *)
IF NOT init OR T = T#0s THEN
	init := TRUE;
	last := tx;
	Y := X;
ELSIF Y = X THEN
	last := tx;
ELSE
	tmp :=  WORD_TO_DWORD(X - Y) * (tx - last) / TIME_TO_DWORD(T);
	IF tmp <> 0 THEN
		Y := DINT_TO_WORD(WORD_TO_DINT(Y) + DWORD_TO_DINT(tmp));
		last := tx;
	END_IF;
END_IF;

(*
hm 10. oct. 2008	rev 1.0
	original version

hm	3. nov. 2008	rev 1.1
	fixed overflow problem in formula

hm	25. jan. 2011	rev 1.2
	fixed error in formula
*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FILTER_WAV
VAR_INPUT
	X : REAL;
	W : ARRAY[0..15] OF REAL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Y : REAL;
END_VAR
VAR
	init: BOOL;
	buffer : ARRAY[0..15] OF REAL;
	i: INT;
	n: INT;
END_VAR

(*
version 1.2	10. mar. 2009
programmer 	hugo
tested by		oscat

FILTER_WAV is a moving average filter with programmable length N for DWORD Data.

*)
(* startup initialisation *)
IF NOT init OR rst THEN
	init := TRUE;
	FOR i := 0 TO 15 DO
		buffer[i] := X;
	END_FOR;
	i := 15;
	Y := X;
ELSE
	i := INC1(i, 16);
	buffer[i] := X;
END_IF;

(* calculate the weighted average *)
Y := 0.0;
FOR n := 0 TO 15 DO
	Y := buffer[i] * W[n] + Y;
	i := DEC1(i, 16);
END_FOR;


(*
hm 	14. oct. 2008	rev 1.0
	original version

hm	27. oct. 2008	rev 1.1
	changed _DEC and _INC to DEC1 and INC1

hm	10. mar. 2009	rev 1.2
	real constants updated to new systax using dot

*)

END_FUNCTION_BLOCK