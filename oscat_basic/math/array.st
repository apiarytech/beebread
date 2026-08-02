(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Array' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _ARRAY_ABS
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF REAL;
	SIZE : INT;
END_VAR
VAR
	i : INT;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function calculates the absolute value of every element in an array of REAL

*)
(* @END_DECLARATION := '0' *)
FOR i := 0 TO size - 1 DO
	pt^[i] := ABS(pt^[i]);
END_FOR;


(* revision history
hm	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Array' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _ARRAY_ADD
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF REAL;
	SIZE : INT;
	VAL : REAL;
END_VAR
VAR
	i : INT;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function adds a value to every element in an array of REAL

*)
(* @END_DECLARATION := '0' *)
FOR i := 0 TO size - 1 DO
	pt^[i] := pt^[i] + VAL;
END_FOR;


(* revision history
hm	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Array' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _ARRAY_INIT
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF REAL;
	SIZE : INT;
	VAL : REAL;
END_VAR
VAR
	i : INT;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function initializes every element in an array of REAL with a value

*)
(* @END_DECLARATION := '0' *)
FOR i := 0 TO size - 1 DO
	pt^[i] := VAL;
END_FOR;


(* revision history
hm	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Array' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _ARRAY_MEDIAN
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF REAL;
	SIZE : INT;
END_VAR
VAR
	i, j : INT;
	buf : REAL;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function calculates the median of an array of REAL

*)
(* @END_DECLARATION := '0' *)
(* sort the array first *)
FOR i := 0 TO size - 2 DO
	FOR j := i + 1 TO size - 1 DO
		IF pt^[i] > pt^[j] THEN
			buf := pt^[i];
			pt^[i] := pt^[j];
			pt^[j] := buf;
		END_IF;
	END_FOR;
END_FOR;

(* calculate the median *)
i := size / 2;
IF (size MOD 2) = 1 THEN
	_ARRAY_MEDIAN := pt^[i];
ELSE
	_ARRAY_MEDIAN := (pt^[i-1] + pt^[i]) * 0.5;
END_IF;


(* revision history
hm	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Array' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _ARRAY_MUL
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF REAL;
	SIZE : INT;
	VAL : REAL;
END_VAR
VAR
	i : INT;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function multiplies every element in an array of REAL with a value

*)
(* @END_DECLARATION := '0' *)
FOR i := 0 TO size - 1 DO
	pt^[i] := pt^[i] * VAL;
END_FOR;


(* revision history
hm	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Array' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _ARRAY_SHUFFLE
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF REAL;
	SIZE : INT;
END_VAR
VAR
	i, j : INT;
	buf : REAL;
	last : REAL;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function shuffles every element in an array of REAL

*)
(* @END_DECLARATION := '0' *)
FOR i := 0 TO size - 1 DO
	last := RDM(last);
	j := REAL_TO_INT(last * INT_TO_REAL(size - 1));
	buf := pt^[i];
	pt^[i] := pt^[j];
	pt^[j] := buf;
END_FOR;


(* revision history
hm	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Array' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _ARRAY_SORT
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF REAL;
	SIZE : INT;
	DESC : BOOL;
END_VAR
VAR
	i, j : INT;
	buf : REAL;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function sorts an array of REAL ascending or descending

*)
(* @END_DECLARATION := '0' *)
IF DESC THEN
	FOR i := 0 TO size - 2 DO
		FOR j := i + 1 TO size - 1 DO
			IF pt^[i] < pt^[j] THEN
				buf := pt^[i];
				pt^[i] := pt^[j];
				pt^[j] := buf;
			END_IF;
		END_FOR;
	END_FOR;
ELSE
	FOR i := 0 TO size - 2 DO
		FOR j := i + 1 TO size - 1 DO
			IF pt^[i] > pt^[j] THEN
				buf := pt^[i];
				pt^[i] := pt^[j];
				pt^[j] := buf;
			END_IF;
		END_FOR;
	END_FOR;
END_IF;


(* revision history
hm	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION