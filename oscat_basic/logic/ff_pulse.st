(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF pulse triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK LTCH
VAR_INPUT
	D : BOOL;
	L : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
END_VAR


(*
version 1.2	30. oct. 2008
programmer 	hugo
tested by	oscat

Transparent Latch with asynchronous reset


*)

(* @END_DECLARATION := '0' *)
(* as long as L is true, the latch is transparent and and change of D is transferred to Q *)
(* of course only when there is no asynchronous reset *)
IF rst THEN				(* if asynchronous reset then Q=0 *)
	Q := FALSE;
ELSIF L THEN			(* latch is transparent *)
	Q := D;
END_IF;


(* revision history
hm	4. aug 2006		rev 1.0
	original version

hm	27. dec 2007	rev 1.1
	changed code for better performance

hm	30. oct. 2008	rev 1.2
	deleted unnecessary init with 0

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF pulse triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK LTCH_4
VAR_INPUT
	D0 : BOOL;
	D1 : BOOL;
	D2 : BOOL;
	D3 : BOOL;
	L : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q0 : BOOL;
	Q1 : BOOL;
	Q2 : BOOL;
	Q3 : BOOL;
END_VAR


(*
version 1.3	14. mar. 2009
programmer 	hugo
tested by		oscat

Quad Transparent Latch with asynchronous reset


*)

(* @END_DECLARATION := '0' *)
(* as long as L is true, the latch is transparent and and change of D is transferred to Q *)
(* of course only when there is no asynchronous reset *)
IF rst THEN			(* if asynchronous reset then Q=0 *)
	Q0 := FALSE;
	Q1 := FALSE;
	Q2 := FALSE;
	Q3 := FALSE;
ELSIF L THEN			(* latch is transparent *)
	Q0 := D0;
	Q1 := D1;
	Q2 := D2;
	Q3 := D3;
END_IF;

(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	27. dec 2007	rev 1.1
	changed code for better performance

hm	30. oct. 2008	rev 1.2
	deleted unnecessary init with 0

hm	14. mar. 2009	rev 1.3
	removed double assignments

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF pulse triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK STORE_8
VAR_INPUT
	Set : BOOL;
	D0,D1, D2, D3, D4, D5, D6, D7 : BOOL;
	Clr : BOOL;
	Rst : BOOL;
END_VAR
VAR_OUTPUT
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7 : BOOL;
END_VAR
VAR
	edge: BOOL;
END_VAR

(*
version 1.2	14. mar. 2009
programmer 	hugo
tested by		oscat

Store_4 stores up to 4 boolean inputs until a reset is clearing the outputs.
the respective output is set with a true at the respective input and stays true until a reset is generated.
a set input is used to set all outputs true at the same time, independent of their inputs.
a clear input is used to reset only the lowest priority of the outputs one at a time.
the lowest priority is q3.

*)
(* @END_DECLARATION := '0' *)
IF rst OR set THEN
	q0 := NOT rst;
	q1 := q0;
	q2 := q0;
	q3 := q0;
	q4 := q0;
	q5 := q0;
	q6 := q0;
	q7 := q0;
ELSE
	IF D0 THEN Q0 := TRUE; END_IF;
	IF D1 THEN Q1 := TRUE; END_IF;
	IF D2 THEN Q2 := TRUE; END_IF;
	IF D3 THEN Q3 := TRUE; END_IF;
	IF D4 THEN Q4 := TRUE; END_IF;
	IF D5 THEN Q5 := TRUE; END_IF;
	IF D6 THEN Q6 := TRUE; END_IF;
	IF D7 THEN Q7 := TRUE; END_IF;

	IF clr AND NOT edge THEN
		IF q0 THEN q0 := FALSE;
		ELSIF q1 THEN q1 := FALSE;
		ELSIF q2 THEN q2 := FALSE;
		ELSIF q3 THEN q3 := FALSE;
		ELSIF q4 THEN q4 := FALSE;
		ELSIF q5 THEN q5 := FALSE;
		ELSIF q6 THEN q6 := FALSE;
		ELSE q7 := FALSE;
		END_IF;
	END_IF;
	edge := clr;
END_IF;


(* revision history
hm	25.12.2007	rev 1.0
	original version

hm	30. oct. 2008	rev 1.1
	optimized performance

hm	14. mar. 2009	rev 1.2
	removed double assignments

*)


END_FUNCTION_BLOCK