(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION ACOSH : REAL
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the arccos hyperbolicus of X

*)
(* @END_DECLARATION := '0' *)
ACOSH := LN(X + SQRT(X*X - 1.0));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION ACOTH : REAL
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the arccot hyperbolicus of X

*)
(* @END_DECLARATION := '0' *)
ACOTH := 0.5 * LN((X+1.0)/(X-1.0));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION AGDF : REAL
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the inverse gudermann function of X

*)
(* @END_DECLARATION := '0' *)
AGDF := LN(TAN(math.PI025 + 0.5 * X));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION ASINH : REAL
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the arcsin hyperbolicus of X

*)
(* @END_DECLARATION := '0' *)
ASINH := LN(X + SQRT(X*X + 1.0));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION ATAN2 : REAL
VAR_INPUT
	Y : REAL;
	X : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function calculates the arctan of Y/X

*)
(* @END_DECLARATION := '0' *)
IF X > 0.0 THEN
	ATAN2 := ATAN(Y/X);
ELSIF Y >= 0.0 AND X < 0.0 THEN
	ATAN2 := ATAN(Y/X) + math.PI;
ELSIF Y < 0.0 AND X < 0.0 THEN
	ATAN2 := ATAN(Y/X) - math.PI;
ELSIF Y > 0.0 AND X = 0.0 THEN
	ATAN2 := math.PI05;
ELSIF Y < 0.0 AND X = 0.0 THEN
	ATAN2 := -math.PI05;
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION ATANH : REAL
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the arctan hyperbolicus of X

*)
(* @END_DECLARATION := '0' *)
ATANH := 0.5 * LN((1.0+X)/(1.0-X));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BETA : REAL
VAR_INPUT
	X : REAL;
	Y : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the beta function of X and Y

*)
(* @END_DECLARATION := '0' *)
BETA := GAMMA(X) * GAMMA(Y) / GAMMA(X+Y);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BINOM : DINT
VAR_INPUT
	N : INT;
	K : INT;
END_VAR


(*
version 1.1	26. mar. 2011
programmer 	hugo
tested by		oscat

this function calculates the binomial coefficient of N over K

*)
(* @END_DECLARATION := '0' *)
IF K >= 0 AND K <= N AND N < 13 THEN
	BINOM := math.FACTS[N] / (math.FACTS[K] * math.FACTS[N-K]);
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

hm	26. mar. 2011	rev 1.1
	using math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CAUCHY : REAL
VAR_INPUT
	X : REAL;
	M : REAL;
	G : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the cauchy distribution

*)
(* @END_DECLARATION := '0' *)
CAUCHY := math.PI_INV * G / ( (X-M)*(X-M) + G*G);


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CAUCHYCD : REAL
VAR_INPUT
	X : REAL;
	M : REAL;
	G : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the cumulative cauchy distribution

*)
(* @END_DECLARATION := '0' *)
CAUCHYCD := 0.5 + math.PI_INV * ATAN((X-M)/G);


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CEIL : INT
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function calculates the next higher integer of X

*)
(* @END_DECLARATION := '0' *)
CEIL := REAL_TO_INT(X);
IF X > 0.0 AND X <> INT_TO_REAL(CEIL) THEN
	CEIL := CEIL + 1;
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CEIL2 : REAL
VAR_INPUT
	X : REAL;
	Y : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function calculates the next higher integer of X/Y * Y

*)
(* @END_DECLARATION := '0' *)
IF Y <> 0.0 THEN
	CEIL2 := INT_TO_REAL(CEIL(X/Y)) * Y;
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CMP : INT
VAR_INPUT
	X1 : REAL;
	X2 : REAL;
	D : REAL;
END_VAR


(*
version 1.1	11. mar. 2009
programmer 	hugo
tested by		tobias

this function compares two real values with a given deviation D

*)
(* @END_DECLARATION := '0' *)
IF ABS(X1-X2) <= D THEN
	CMP := 0;
ELSIF X1 > X2 THEN
	CMP := 1;
ELSE
	CMP := -1;
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	code improvements

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION COSH : REAL
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the cosinus hyperbolicus of X

*)
(* @END_DECLARATION := '0' *)
COSH := 0.5 * (EXP(X) + EXP(-X));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION COTH : REAL
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the cotangens hyperbolicus of X

*)
(* @END_DECLARATION := '0' *)
COTH := (EXP(X) + EXP(-X)) / (EXP(X) - EXP(-X));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION D_TRUNC : DINT
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function truncates a real value to a dint

*)
(* @END_DECLARATION := '0' *)
D_TRUNC := REAL_TO_DINT(X);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION