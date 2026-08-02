(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Vector' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION V3_ABS : REAL
VAR_INPUT
	V : VECTOR_3;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the absolute value of a vector

*)
(* @END_DECLARATION := '0' *)
V3_ABS := SQRT(V.X*V.X + V.Y*V.Y + V.Z*V.Z);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Vector' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION V3_ADD : VECTOR_3
VAR_INPUT
	V1 : VECTOR_3;
	V2 : VECTOR_3;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function adds two vectors

*)
(* @END_DECLARATION := '0' *)
V3_ADD.X := V1.X + V2.X;
V3_ADD.Y := V1.Y + V2.Y;
V3_ADD.Z := V1.Z + V2.Z;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Vector' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION V3_CROSS : VECTOR_3
VAR_INPUT
	V1 : VECTOR_3;
	V2 : VECTOR_3;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the cross product of two vectors

*)
(* @END_DECLARATION := '0' *)
V3_CROSS.X := V1.Y*V2.Z - V1.Z*V2.Y;
V3_CROSS.Y := V1.Z*V2.X - V1.X*V2.Z;
V3_CROSS.Z := V1.X*V2.Y - V1.Y*V2.X;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Vector' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION V3_DOT : REAL
VAR_INPUT
	V1 : VECTOR_3;
	V2 : VECTOR_3;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the dot product of two vectors

*)
(* @END_DECLARATION := '0' *)
V3_DOT := V1.X*V2.X + V1.Y*V2.Y + V1.Z*V2.Z;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Vector' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION V3_MUL : VECTOR_3
VAR_INPUT
	V : VECTOR_3;
	S : REAL;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function multiplies a vector with a scalar

*)
(* @END_DECLARATION := '0' *)
V3_MUL.X := V.X * S;
V3_MUL.Y := V.Y * S;
V3_MUL.Z := V.Z * S;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Vector' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION V3_SUB : VECTOR_3
VAR_INPUT
	V1 : VECTOR_3;
	V2 : VECTOR_3;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function subtracts two vectors

*)
(* @END_DECLARATION := '0' *)
V3_SUB.X := V1.X - V2.X;
V3_SUB.Y := V1.Y - V2.Y;
V3_SUB.Z := V1.Z - V2.Z;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION