(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CABS : REAL
VAR_INPUT
	X : COMPLEX;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the absolute value of a complex number

*)
(* @END_DECLARATION := '0' *)
CABS := HYPOT(X.re, X.im);


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using hypot function

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CACOS : COMPLEX
VAR_INPUT
	X : COMPLEX;
END_VAR
VAR
	t1 : REAL;
	t2 : REAL;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the arccos of a complex number

*)
(* @END_DECLARATION := '0' *)
t1 := 0.5 * CABS(CADD(X,CSET(1.0,0.0)));
t2 := 0.5 * CABS(CSUB(X,CSET(1.0,0.0)));
CACOS.re := 2.0 * ATAN(t2/t1);
CACOS.im := -LN(t1+t2) + LN(t1-t2);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CADD : COMPLEX
VAR_INPUT
	X1 : COMPLEX;
	X2 : COMPLEX;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function adds two complex numbers

*)
(* @END_DECLARATION := '0' *)
CADD.re := X1.re + X2.re;
CADD.im := X1.im + X2.im;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CASIN : COMPLEX
VAR_INPUT
	X : COMPLEX;
END_VAR
VAR
	t1 : REAL;
	t2 : REAL;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the arcsin of a complex number

*)
(* @END_DECLARATION := '0' *)
t1 := 0.5 * CABS(CADD(X,CSET(1.0,0.0)));
t2 := 0.5 * CABS(CSUB(X,CSET(1.0,0.0)));
CASIN.re := ATAN(X.re / (t1+t2));
CASIN.im := LN(t1+t2) + LN(SQRT(X.re*X.re + (t1+t2)*(t1+t2)));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CATAN : COMPLEX
VAR_INPUT
	X : COMPLEX;
END_VAR
VAR
	t1 : REAL;
	t2 : REAL;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the arctan of a complex number

*)
(* @END_DECLARATION := '0' *)
t1 := X.re*X.re + (1.0-X.im)*(1.0-X.im);
t2 := X.re*X.re + (1.0+X.im)*(1.0+X.im);
CATAN.re := 0.5 * ATAN(2.0*X.re/(1.0-X.re*X.re-X.im*X.im));
CATAN.im := 0.25 * LN(t2/t1);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CCOS : COMPLEX
VAR_INPUT
	X : COMPLEX;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the cos of a complex number

*)
(* @END_DECLARATION := '0' *)
CCOS.re := COS(X.re) * COSH(X.im);
CCOS.im := -SIN(X.re) * SINH(X.im);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CDIV : COMPLEX
VAR_INPUT
	X1 : COMPLEX;
	X2 : COMPLEX;
END_VAR
VAR
	t1 : REAL;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function divides two complex numbers

*)
(* @END_DECLARATION := '0' *)
t1 := X2.re*X2.re + X2.im*X2.im;
CDIV.re := (X1.re*X2.re + X1.im*X2.im)/t1;
CDIV.im := (X1.im*X2.re - X1.re*X2.im)/t1;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CEXP : COMPLEX
VAR_INPUT
	X : COMPLEX;
END_VAR
VAR
	t1 : REAL;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the exp of a complex number

*)
(* @END_DECLARATION := '0' *)
t1 := EXP(X.re);
CEXP.re := t1 * COS(X.im);
CEXP.im := t1 * SIN(X.im);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CLN : COMPLEX
VAR_INPUT
	X : COMPLEX;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the ln of a complex number

*)
(* @END_DECLARATION := '0' *)
CLN.re := 0.5 * LN(X.re*X.re + X.im*X.im);
CLN.im := ATAN2(X.im, X.re);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CMUL : COMPLEX
VAR_INPUT
	X1 : COMPLEX;
	X2 : COMPLEX;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function multiplies two complex numbers

*)
(* @END_DECLARATION := '0' *)
CMUL.re := X1.re*X2.re - X1.im*X2.im;
CMUL.im := X1.re*X2.im + X1.im*X2.re;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Complex' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CSET : COMPLEX
VAR_INPUT
	re : REAL;
	im : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function creates a complex number from two real values

*)
(* @END_DECLARATION := '0' *)
CSET.re := re;
CSET.im := im;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION