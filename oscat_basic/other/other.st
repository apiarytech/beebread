(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/other' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK ESR_COLLECT
VAR_INPUT
	E : BOOL;
	T : TIME;
	ADR : STRING(10);
END_VAR
VAR_OUTPUT
	Q : BOOL;
	CNT : INT;
END_VAR
VAR
	tx : TIME;
	last : TIME;
	init : BOOL;
	edge : BOOL;
	buf : ARRAY[0..15] OF ESR_DATA;
	i: INT;
END_VAR

(*
version 1.1	16 sep 2007
programmer 	hugo
tested by		oscat

ESR_COLLECT collects data from ESR_MON function blocks and stores them in a buffer.
the buffer can be read by a master station.

*)
(* @END_DECLARATION := '0' *)
(* read system time *)
tx := DWORD_TO_TIME(T_PLC_MS());

(* startup initialisation *)
IF NOT init THEN
	init := TRUE;
	last := tx;
END_IF;

(* check for trigger *)
IF E AND NOT edge THEN
	buf[cnt].ds := tx;
	buf[cnt].ts := tx - last;
	buf[cnt].adress := adr;
	cnt := cnt + 1;
	q := TRUE;
ELSE
	q := FALSE;
END_IF;
edge := E;

(* check for buffer full *)
IF cnt > 15 THEN cnt := 15; END_IF;

(* check for timeout *)
IF tx - last > T THEN
	last := tx;
	buf[cnt].ds := tx;
	buf[cnt].ts := T;
	buf[cnt].adress := 'TIMEOUT';
	cnt := cnt + 1;
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

hm	16. sep. 2007	rev 1.1
	replaced time() with T_PLC_MS for compatibility and performance reasons

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/other' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK ESR_MON_B8
VAR_INPUT
	IN : BYTE;
END_VAR
VAR_OUTPUT
	Q : BOOL;
END_VAR
VAR
	last : BYTE;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

ESR_MON_B8 monitors a byte for change and sets the output q for one cycle if a change has been detected

*)
(* @END_DECLARATION := '0' *)
Q := IN <> last;
last := IN;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/other' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION OSCAT_VERSION : STRING(10)


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the version of the oscat library

*)
(* @END_DECLARATION := '0' *)
OSCAT_VERSION := 'BASIC 3.35';


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION