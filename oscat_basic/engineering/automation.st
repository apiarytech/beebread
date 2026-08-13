FUNCTION_BLOCK DRIVER_1
VAR_INPUT CONSTANT
	Toggle_Mode : BOOL;
	Timeout : TIME;
END_VAR
VAR_INPUT
	SET : BOOL;
	IN : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
END_VAR
VAR
	off : TON;
	edge: BOOL;
END_VAR

(*
version 1.0	2 jan 2008
programmer 	hugo
tested by		tobias

driver_1 is a multi purpose driver.
a rising edge on in sets the output high if toggle is flase. while toggle is true, a rising edge on in toggles the output Q.
if a timeout is specified the output q will be reset to false automatically after the timeout has elapsed.
a asynchronous reset and set will force the output high or low respectively.

*)
IF off.Q THEN Q := FALSE; END_IF;
IF rst THEN
	Q := FALSE;
ELSIF set THEN
	Q := TRUE;
ELSIF IN AND NOT edge THEN
	IF toggle_mode THEN q := NOT Q; ELSE q := TRUE; END_IF;
END_IF;
edge := in;
IF timeout > t#0s THEN off(in := Q, PT := Timeout); END_IF;


(* revision history
hm	2. jan 2008		rev 1.0
	original version

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK DRIVER_4
VAR_INPUT CONSTANT
	Toggle_Mode : BOOL;
	Timeout : TIME;
END_VAR
VAR_INPUT
	SET : BOOL;
	IN0, IN1, IN2, IN3 : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q0, Q1, Q2, Q3 : BOOL;
END_VAR
VAR
	d0, d1, d2, d3 : DRIVER_1;
END_VAR

(*
version 1.0	2 jan 2008
programmer 	hugo
tested by		tobias

driver_4 is a 4 channel multi purpose driver.
a rising edge on in? sets the output Q? high if toggle_mode is flase. while toggle_mode is true, a rising edge on in? toggles the output Q?.
if a timeout is specified the output Q? will be reset to false automatically after the timeout has elapsed.
a asynchronous reset and set will force the output high or low respectively.
an asynchronous set will force all outputs high simultaneously

*)
D0(Set:=set, in:=in0, rst:=rst, toggle_mode:=toggle_mode, timeout:=timeout);
D1(Set:=set, in:=in1, rst:=rst, toggle_mode:=toggle_mode, timeout:=timeout);
D2(Set:=set, in:=in2, rst:=rst, toggle_mode:=toggle_mode, timeout:=timeout);
D3(Set:=set, in:=in3, rst:=rst, toggle_mode:=toggle_mode, timeout:=timeout);
Q0 := D0.Q;
Q1 := D1.Q;
Q2 := D2.Q;
Q3 := D3.Q;


(* revision history
hm	2. jan 2008		rev 1.0
	original version

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK DRIVER_4C
VAR_INPUT
	IN : BOOL;
	RST : BOOL;
END_VAR
VAR_INPUT CONSTANT
	Timeout : TIME;
	SX : ARRAY[1..7] OF BYTE := 1,3,7,15;
END_VAR
VAR_OUTPUT
	SN : INT;
	Q0 : BOOL;
	Q1 : BOOL;
	Q2 : BOOL;
	Q3 : BOOL;
END_VAR
VAR
	off : TON;
	edge: BOOL;
END_VAR

(*
version 1.0	23. mar. 2009
programmer 	hugo
tested by	tobias

driver_4C is a multi purpose driver.
a rising edge on in switches from S0 state S1 and the next edge to state S2 and so on.
in state S0 all outputs Q are FALSE.
The stet of the Outputs in any state S? is configurable with setup variables.
The variables S1..S5 define the states, while the sequence is terminated when a state Variable S? = 0.
The lower bits 0..3 of the state vars S? are corresponding to the Outputs Q0..Q3

*)
IF RST OR off.Q THEN
	SN := 0;
ELSIF IN AND NOT edge THEN
	SN := SN + 1;
	IF SN > 7 OR SX[SN] = 0 THEN SN := 0; END_IF;
END_IF;
edge := in;
IF SN > 0 THEN
	Q0 := SX[SN].0;
	Q1 := SX[SN].1;
	Q2 := SX[SN].2;
	Q3 := SX[SN].3;
ELSE
	Q0 := FALSE;
	Q1 := FALSE;
	Q2 := FALSE;
	Q3 := FALSE;
END_IF;

(* maximaum timeout *)
IF timeout > t#0s THEN off(in := SN > 0, PT := Timeout); END_IF;


(* revision history
hm	23. mar. 2009	rev 1.0
	original version

*)
END_FUNCTION_BLOCK


 FUNCTION_BLOCK FLOW_CONTROL
VAR_INPUT
	IN : BOOL;
	REQ : BOOL;
	ENQ : BOOL;
	RST : BOOL;
END_VAR
VAR_INPUT CONSTANT
	T_AUTO : TIME := T#1h;
	T_DELAY : TIME := T#23h;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	STATUS : BYTE;
END_VAR
VAR
	timer : TP_1D;
END_VAR


(*
version 1.0	28. jun. 2008
programmer 	hugo
tested by	oscat

FLOW_CONTROL switches a valves depending on the input in.
flow control also limits the maximum ontime of the valve and controls pressure on the output side.

*)
STATUS := 100;
IF RST THEN
	Q := FALSE;
	timer(rst := TRUE);
	timer.RST := FALSE;
	STATUS := 103;
ELSIF ENQ THEN
	IF IN THEN
		status := 101;
	END_IF;
	IF REQ THEN
		(* timer will generate a timed pulse after TP goes high *)
		timer.PT1 := T_AUTO;
		timer.PTD := T_DELAY;
		timer.IN := TRUE;
		STATUS := 102;
	END_IF;
END_IF;


(* set output and run timer *)
timer();
timer.IN := FALSE;
Q := (IN AND ENQ) OR timer.Q;



(* revision history
hm 	28. jun. 2008 	rev 1.0
	original version

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK FT_Profile
VAR_INPUT
	K : REAL := 1.0;
	O : REAL;
	M : REAL := 1.0;
	E : BOOL;
END_VAR
VAR_INPUT CONSTANT
	value_0 : REAL;
	time_1 : TIME;
	value_1 : REAL;
	time_2 : TIME;
	value_2 : REAL;
	time_3 : TIME;
	value_3 : REAL;
	time_10 : TIME;
	value_10 : REAL;
	time_11 : TIME;
	value_11 : REAL;
	time_12 : TIME;
	value_12 : REAL;
	time_13 : TIME;
	value_13 : REAL;
END_VAR
VAR_OUTPUT
	Y : REAL;
	RUN : BOOL;
	ET : TIME;
END_VAR
VAR
	tx : TIME;
	edge : BOOL;
	state: BYTE;
	ta: TIME;
	tb: TIME;
	t0 : TIME;
	temp: REAL;
	va: REAL;
	vb: REAL;
END_VAR

(*
version 1.1	15 sep 2007
programmer 	tobias
tested by		hugo

FT_Profile generates an output signal which is defined by values over a time scale.
the different values are connected by ramps between the individual values.
a rising edge on E starts the output signal generation and E = True can delay time_3 / value_3 as long as it stays true.
an additional multiplier K can be used to multiply the output and an offset O can be added to the output dynamically.
*)
(* @END_DECLARATION := '0' *)
(* read system time *)
tx := DWORD_TO_TIME(T_PLC_MS());

(* determine start condition *)
IF E AND NOT edge THEN
	RUN := TRUE;
	ET := t#0s;
	t0 := tx;
	ta := tx;
	tb := multime(time_1, M);
	va := value_0;
	vb := value_1;
	temp := value_0;
	state := 1;
END_IF;
edge := E;

(* generate startup profile *)
IF run THEN
	CASE state OF
		1:	IF tx - ta >= tb THEN
				ta := ta + tb;
				tb := multime(time_2 - time_1, M);
				va := value_1;
				vb := value_2;
				temp := value_1;
				state := 2;
			ELSE
				temp := ((vb - va) * TIME_TO_REAL(tx - ta) / TIME_TO_REAL(tb) + va);
			END_IF;
		2:	IF tx - ta >= tb THEN
				ta := ta + tb;
				tb := multime(time_3 - time_2, M);
				va := value_2;
				vb := value_3;
				temp := value_2;
				state := 3;
			ELSE
				temp := ((vb - va) * TIME_TO_REAL(tx - ta) / TIME_TO_REAL(tb) + va);
			END_IF;
		3:	IF tx - ta >= tb THEN
				ta := ta + tb;
				tb := multime(time_10 - time_3, M);
				va := value_3;
				vb := value_10;
				temp := value_3;
				state := 4;
			ELSE
				temp := ((vb - va) * TIME_TO_REAL(tx - ta) / TIME_TO_REAL(tb) + va);
			END_IF;
		4 :	IF tx - ta >= tb THEN
				ta := ta + tb;
				tb := multime(time_11 - time_10, M);
				va := value_10;
				vb := value_11;
				temp := value_10;
				IF E THEN state := 5; ELSE state := 6; END_IF;
			ELSE
				temp := ((vb - va) * TIME_TO_REAL(tx - ta) / TIME_TO_REAL(tb) + va);
			END_IF;
		5:	(* extend the signal while E is true *)
			IF E THEN
				ta := tx;
			ELSE
				state := 6;
			END_IF;
		6:	IF tx - ta >= tb THEN
				ta := ta + tb;
				tb := multime(time_12 - time_11, M);
				va := value_11;
				vb := value_12;
				temp := value_11;
				state := 7;
			ELSE
				temp := ((vb - va) * TIME_TO_REAL(tx - ta) / TIME_TO_REAL(tb) + va);
			END_IF;
		7:	IF tx - ta >= tb THEN
				ta := ta + tb;
				tb := multime(time_13 - time_12, M);
				va := value_12;
				vb := value_13;
				temp := value_12;
				state := 8;
			ELSE
				temp := ((vb - va) * TIME_TO_REAL(tx - ta) / TIME_TO_REAL(tb) + va);
			END_IF;
		8:	IF tx - ta >= tb THEN
				temp := value_13;
				run := FALSE;
			ELSE
				temp := ((vb - va) * TIME_TO_REAL(tx - ta) / TIME_TO_REAL(tb) + va);
			END_IF;
	END_CASE;
	Y := temp * K + O;
	ET := tx - t0;
END_IF;

(* revision history
hm	27 feb 2007		rev 1.0
	original version

hm	15. sep2007		rev 1.1
	replaced Time() with T_PLC_MS for compatibility and performance reasons

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK INC_DEC
VAR_INPUT
	CHa : BOOL;
	CHb : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	dir : BOOL;
	cnt : INT;
END_VAR
VAR
	edgea : BOOL;
	clk: BOOL;
	clka: BOOL;
	clkb: BOOL;
	edgeb: BOOL;
	axb: BOOL;
END_VAR

(*
	version 1.0	4 aug 2006
	programmer 	oscat
	tested BY		oscat

incremental decoder with quadruple accuracy.
2 pulses for each channel are created for each directional pulse.
*)
(* @END_DECLARATION := '0' *)
axb := cha XOR chb;

(* create pulses for channel a *)
clka := cha XOR edgea;
edgea := cha;

clkb := chb XOR edgeb;
edgeb := chb;

(* create pulses for both channels *)
clk := clka OR clkb;

(* set the direction output *)
IF axb AND clka THEN dir := TRUE; END_IF;
IF axb AND clkb THEN dir := FALSE; END_IF;

(* increment or decrement the counter *)
IF clk AND dir THEN cnt := cnt + 1; END_IF;
IF clk AND NOT dir THEN cnt := cnt -1; END_IF;

(* reset the counter if rst active *)
IF rst THEN cnt := 0; END_IF;

END_FUNCTION_BLOCK


FUNCTION_BLOCK INTERLOCK
VAR_INPUT
	I1, I2 : BOOL;
	TL : TIME;
END_VAR
VAR_OUTPUT
	Q1, Q2 : BOOL;
END_VAR
VAR
	T1, T2 : TOF;
END_VAR

(*
version 1.0	28 sep 2007
programmer 	hugo
tested by		tobias

INTERLOCK has two inputs I1 and I2 which drive the corresponding outputs Q1 and Q2.
the inputs signals lock each other out and therfore I1 can only drive Q1 when I2 is Low and vice versa.
The input TL specifies a dead time between two outputs can become active.
an output can only become active when the other output was not active for the time TL.

*)
(* the input signal have a run delay to lockout the other input *)
T1(IN := I1, PT := TL);
T2(IN := I2, PT := TL);

Q1 := I1 AND NOT t2.Q;
Q2 := I2 AND NOT t1.Q;


(* revision history

hm	28 sep 2007		rev 1.0
	original version

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK INTERLOCK_4
VAR_INPUT CONSTANT
END_VAR
VAR_INPUT
	I0 : BOOL;
	I1 : BOOL;
	I2 : BOOL;
	I3 : BOOL;
	E : BOOL;
	MODE : INT;

END_VAR
VAR_OUTPUT
	OUT : BYTE;
	TP : BOOL;
END_VAR
VAR
	in: BYTE;
	last: BYTE;
	old: BYTE;
	lmode: INT;
END_VAR

(*
version 1.1	14. mar. 2009
programmer 	hugo
tested by		oscat

INTERLOCK_4 detects one of 4 switches and delivers the number of the switch pressed on the output out
a output tp is true for one cycle if the output has changed.
a setup variable MODE selects between 3 different modes:
MODE = 0, any input active will disable all other inputs
MODE = 1, the input with the highest number will be acepted
mode = 2, the input last pressed will disable all others 

*)
(* check if enable is active *)
IF E THEN
(* reset all vars when there is a mode change on thy fly *)
	IF mode <> lmode THEN
		out := 0;
		last := 0;
		old := 0;
		lmode := mode;
	END_IF;
	(* load inputs into in *)
	in.0 := I0;
	in.1 := I1;
	in.2 := I2;
	in.3 := I3;
	(* only execute when there is any change *)
	IF in <> last THEN
		(* only execute when inputs have chages state *)
		CASE mode OF
			0:	(* output directly display inputs as bits in byte out *)
				out := in;

			1:	(* the input with the highest number will be acepted *)
				IF in.3 THEN out := 8;
				ELSIF in.2 THEN out := 4;
				ELSIF in.1 THEN out := 2;
				ELSE out := in;
				END_IF;

			2:	(* input last pressed will be displayed only *)
				last := ((in XOR last) AND in);
				IF last.3 THEN out := 8;
				ELSIF last.2 THEN out := 4;
				ELSIF last.1 THEN out := 2;
				ELSE out := last;
				END_IF;

			3:	(* any input active will disable all other inputs *)
				IF (out AND in) = 0 THEN
					IF in.3 THEN out := 8;
					ELSIF in.2 THEN out := 4;
					ELSIF in.1 THEN out := 2;
					ELSE out := in;
					END_IF;
				END_IF;

		END_CASE;
		last := in;
	END_IF;
	tp := out <> old;
	old := out;
ELSE
	out := 0;
	last := 0;
	old := 0;
	lmode := 0;
	tp := FALSE;
END_IF;


(* revision history
hm	24. oct 2008	rev 1.0
	original version

hm	14. mar. 2009	rev 1.1
	replaced double assignments

*)
END_FUNCTION_BLOCK


FUNCTION MANUAL : BOOL
VAR_INPUT
	IN : BOOL;
	ON : BOOL;
	OFF : BOOL;
END_VAR


(*
version 1.0	21. nov. 2008
programmer 	hugo
tested by	oscat

MANUAL is a manual override for digital signals.
when on and off = FALSE, the output follows IN.
ON = TRUE forces the output high, and OFF = TRUE forces the output low.


*)
MANUAL := NOT OFF AND (IN OR ON);



(* revision history
hm	21. nov 2008	rev 1.0
	original version


*)
END_FUNCTION


FUNCTION_BLOCK MANUAL_1
VAR_INPUT
	IN : BOOL;
	MAN : BOOL;
	M_I : BOOL;
	SET : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	STATUS : BYTE;
END_VAR
VAR
	S_edge : BOOL;
	r_edge : BOOL;
	edge: BOOL;
END_VAR

(*
version 1.2	14. mar. 2009
programmer 	hugo
tested by		oscat

MANUAL_1 is a manual override for digital signals.
When MAN = FALSE, the output follows IN and when MAN = TRUE the Output follows M_I.


*)
IF NOT man THEN
	Q := IN;
	STATUS := 100;
	edge := FALSE;
ELSIF NOT s_edge AND set THEN
	Q := TRUE;
	edge := TRUE;
	status := 101;
ELSIF NOT r_edge AND rst THEN
	Q := FALSE;
	edge := TRUE;
	status := 102;
ELSIF NOT edge THEN
	Q := M_I;
	status := 103;
END_IF;

(* remember levels of manual signals *)
s_edge := SET;
r_edge := RST;



(* revision history
hm	17. jun 2008	rev 1.0
	original version

hm	17. oct 2008	rev 1.1
	deleted unnecessary variable m_edge

hm	14. mar. 2009	rev 1.2
	replaced double assignments

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK MANUAL_2
VAR_INPUT
	IN : BOOL;
	ENA : BOOL;
	ON : BOOL;
	OFF : BOOL;
	MAN : BOOL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	STATUS : BYTE;
END_VAR


(*
version 1.0	22. sep. 2008
programmer 	hugo
tested by	oscat

MANUAL_2 is a manual override for boolean signals.
it has static force high and low as well as a manual input.

*)
IF ena THEN
	IF NOT ON AND NOT OFF THEN
		Q := IN;
		STATUS := 100;
	ELSIF on AND NOT off THEN
		Q := TRUE;
		STATUS := 101;
	ELSIF NOT on AND off THEN
		q := FALSE;
		STATUS := 102;
	ELSE
		Q := MAN;
		STATUS := 103;
	END_IF;
ELSE
	Q := FALSE;
	STATUS := 104;
END_IF;



(* revision history
hm	22. sep. 2008		rev 1.0
	original version

*)
END_FUNCTION_BLOCK


FUNCTION_BLOCK MANUAL_4
VAR_INPUT
	I0, I1, I2, I3 : BOOL;
	MAN : BOOL;
	STP : BOOL;
	M0, M1, M2, M3 : BOOL;
END_VAR
VAR_OUTPUT
	Q0, Q1, Q2, Q3 : BOOL;
	STATUS : BYTE;
END_VAR
VAR
	edge: BOOL;
	pos : INT;
	tog: BOOL;
END_VAR


(*
version 1.0	17. jun 2008
programmer 	hugo
tested by	oscat

MANUAL_4 is a manual override for digital signals.
When MAN = FALSE, the output follows IN and when MAN = TRUE the Output follows M_I.


*)
IF man THEN
	IF NOT TOG THEN
		Q0 := M0;
		Q1 := M1;
		Q2 := M2;
		Q3 := M3;
		STATUS := 101;
	END_IF;
	IF STP AND NOT edge THEN
		tog := TRUE;
		CASE pos OF
			0:	Q0 := TRUE;
				Q1 := FALSE;
				Q2 := FALSE;
				Q3 := FALSE;
				STATUS := 110;
			1:	Q0 := FALSE;
				Q1 := TRUE;
				Q2 := FALSE;
				Q3 := FALSE;
				STATUS := 111;
			2:	Q0 := FALSE;
				Q1 := FALSE;
				Q2 := TRUE;
				Q3 := FALSE;
				STATUS := 112;
			3:	Q0 := FALSE;
				Q1 := FALSE;
				Q2 := FALSE;
				Q3 := TRUE;
				STATUS := 113;
		END_CASE;
		pos := INC(pos,1,3);
	END_IF;
ELSE
	Q0 := I0;
	Q1 := I1;
	Q2 := I2;
	Q3 := I3;
	STATUS := 100;
	tog := FALSE;
	pos := 0;
END_IF;

(* remember status of stp *)
edge := STP;



(* revision history
hm	17. jun 2008		rev 1.0
	original version

*)
END_FUNCTION_BLOCK