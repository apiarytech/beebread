(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Functions' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_LIN : REAL
VAR_INPUT
	X : REAL;
	X1 : REAL;
	Y1 : REAL;
	X2 : REAL;
	Y2 : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

F_LIN is a linear interpolation function.
the function calculates a linear function through two points P1(X1,Y1) and P2(X2,Y2) and returns the Y value for a given X.

*)
(* @END_DECLARATION := '0' *)
IF X2 <> X1 THEN
	F_LIN := (Y2 - Y1) / (X2 - X1) * (X - X1) + Y1;
ELSE
	F_LIN := Y1;
END_IF;


(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	avoid divide by 0

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Functions' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_LIN2 : REAL
VAR_INPUT
	X : REAL;
	X1 : REAL;
	Y1 : REAL;
	X2 : REAL;
	Y2 : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

F_LIN2 is a linear interpolation function with limitation.
the function calculates a linear function through two points P1(X1,Y1) and P2(X2,Y2) and returns the Y value for a given X.
the output is limited to Y1 and Y2.

*)
(* @END_DECLARATION := '0' *)
F_LIN2 := LIMIT(Y1, F_LIN(X,X1,Y1,X2,Y2), Y2);


(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	code improvements

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Functions' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_POLY : REAL
VAR_INPUT
	X : REAL;
	A0 : REAL;
	A1 : REAL;
	A2 : REAL;
	A3 : REAL;
	A4 : REAL;
END_VAR


(*
version 1.1	2. dec 2007
programmer 	hugo
tested by		tobias

F_POLY calculates a polynom of 4th grade.
Y = A4*X^4 + A3*X^3 + A2*X^2 + A1*X + A0

*)
(* @END_DECLARATION := '0' *)
F_POLY := (((A4 * X + A3) * X + A2) * X + A1) * X + A0;


(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	2. dec 2007	rev 1.1
	changed code for better performance

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Functions' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_POWER : REAL
VAR_INPUT
	X : REAL;
	N : REAL;
	A : REAL;
	B : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

F_POWER calculates a power function.
Y = A*X^N + B

*)
(* @END_DECLARATION := '0' *)
F_POWER := A * EXPT(X,N) + B;


(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	code improvements

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Functions' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION F_QUAD : REAL
VAR_INPUT
	X : REAL;
	A : REAL;
	B : REAL;
	C : REAL;
END_VAR


(*
version 1.1	2. dec 2007
programmer 	hugo
tested by		tobias

F_QUAD calculates a quadratic function.
Y = A*X^2 + B*X + C

*)
(* @END_DECLARATION := '0' *)
F_QUAD := (A * X + B) * X + C;


(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	2. dec 2007	rev 1.1
	changed code for better performance

*)
END_FUNCTION