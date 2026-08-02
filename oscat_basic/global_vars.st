(* @NESTEDCOMMENTS := 'Yes' *)
(* @GLOBAL_VARIABLE_LIST := 'Constants' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '4096' *)
VAR_GLOBAL CONSTANT
	STRING_LENGTH : INT := 250;
	LIST_LENGTH : INT := 250;
END_VAR




(* revisionhistory
hm 3. mar. 2007
	extended accuracy of pi2
	renamed T0 to Tk for better compatibility with other libraries
	renamed C t0 c0 for better compatibility

hm	2. may. 2007
	added RU universal gas constant

hm	1. oct. 2007
	added setup parameter STEP7

hm 31. oct. 2007
	added PN

hm	6. mar. 2008
	added setup Parameter EXTENDED_ASCII

hm	28. Mar. 2008
	added setup parameter STRING_LENGTH

hm	7. apr. 2008
	Deleted unused setup parameter Step7

hm	19. apr. 2008
	added setup parameter NETWORK_BUFFER_SIZE

hm	27. apr. 2008
	added PI05 and PI025

hm	18. may. 2008
	added charname lists
	renamed E to E1

hm	30. jun. 2008
	reduced size of string constants to 253 to avoid problems with target aixia DCUF

hm	21. sep. 2008
	added section language setup
	added section location setup

hm	20. oct. 2008
	changed to structured constants

hm	16. nov. 2008
	moved structured constants to global vars for compatibility reasons

hm	13. nov. 2009
	added network_buffer_short_size
	renamed network_buffer_long_size
	changed string length to 250

hm	20. jan. 2011
	removed network buffers
	added list_length
*)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
VAR_GLOBAL
	MATH : CONSTANTS_MATH;
	PHYS : CONSTANTS_PHYS;
	LANGUAGE : CONSTANTS_LANGUAGE;
	SETUP : CONSTANTS_SETUP;
	LOCATION : CONSTANTS_LOCATION;
END_VAR

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
VAR_GLOBAL RETAIN

END_VAR