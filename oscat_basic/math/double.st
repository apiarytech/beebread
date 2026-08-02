(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Double' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION R2_ABS : REAL2
VAR_INPUT
	X : REAL2;
END_VAR


(*
version 1.0	24. mar. 2008
programmer 	hugo
tested by	oscat

this function calculates the absolute value of a double real

*)
(* @END_DECLARATION := '0' *)
R2_ABS.RX := ABS(X.RX);
R2_ABS.R1 := ABS(X.R1);


(* revision history
hm	24. mar. 2008	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Double' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION R2_ADD : REAL2
VAR_INPUT
	X1 : REAL2;
	X2 : REAL;
END_VAR
VAR
	t1 : REAL;
	t2 : REAL;
END_VAR

(*
version 1.0	24. mar. 2008
programmer 	hugo
tested by	oscat

this function adds a real to a double real

*)
(* @END_DECLARATION := '0' *)
t1 := X1.RX + X2;
t2 := t1 - X1.RX;
R2_ADD.R1 := (X1.RX - (t1 - t2)) + (X2 - t2) + X1.R1;
R2_ADD.RX := t1 + R2_ADD.R1;
R2_ADD.R1 := (t1 - R2_ADD.RX) + R2_ADD.R1;


(* revision history
hm	24. mar. 2008	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Double' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION R2_ADD2 : REAL2
VAR_INPUT
	X1 : REAL2;
	X2 : REAL2;
END_VAR


(*
version 1.0	24. mar. 2008
programmer 	hugo
tested by	oscat

this function adds two double reals

*)
(* @END_DECLARATION := '0' *)
R2_ADD2 := R2_ADD(X1, X2.RX);
R2_ADD2 := R2_ADD(R2_ADD2, X2.R1);


(* revision history
hm	24. mar. 2008	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Double' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION R2_SET : REAL2
VAR_INPUT
	X : REAL;
END_VAR


(*
version 1.0	24. mar. 2008
programmer 	hugo
tested by	oscat

this function converts a real to a double real

*)
(* @END_DECLARATION := '0' *)
R2_SET.RX := X;
R2_SET.R1 := 0.0;


(* revision history
hm	24. mar. 2008	rev 1.0
	original version

*)
END_FUNCTION