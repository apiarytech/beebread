(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Geometry' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CIRCLE_A : REAL
VAR_INPUT
	R : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the area of a circle with radius R

*)
(* @END_DECLARATION := '0' *)
CIRCLE_A := R * R * math.PI;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Geometry' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CIRCLE_C : REAL
VAR_INPUT
	R : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the circumference of a circle with radius R

*)
(* @END_DECLARATION := '0' *)
CIRCLE_C := R * math.PI2;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Geometry' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CIRCLE_SEG : REAL
VAR_INPUT
	R : REAL;
	A : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the area of a circle segment with radius R and angle A

*)
(* @END_DECLARATION := '0' *)
CIRCLE_SEG := R * R * math.PI * A / 360.0;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Geometry' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CONE_V : REAL
VAR_INPUT
	R : REAL;
	H : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the volume of a cone with radius R and height H

*)
(* @END_DECLARATION := '0' *)
CONE_V := R * R * math.PI * H / 3.0;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using math constants

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Math\/Geometry' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION SPHERE_V : REAL
VAR_INPUT
	R : REAL;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function calculates the volume of a sphere with radius R

*)
(* @END_DECLARATION := '0' *)
SPHERE_V := R * R * R * math.PI4 / 3.0;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	using math constants

*)
END_FUNCTION