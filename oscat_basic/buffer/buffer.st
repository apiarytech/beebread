(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Buffer Management' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _BUFFER_CLEAR : BOOL
VAR_INPUT
	PT : POINTER TO BYTE;
	SIZE : UINT;
END_VAR
VAR
	ptw : POINTER TO DWORD;
	temp: DWORD;
	end, end32 : DWORD;
END_VAR


(*
version 1.2	31. oct. 2008
programmer 	hugo
tested by	oscat

this function will initialize a given array of byte with 0.
the function needs to be called:	_buffer_clear(adr("array"),sizeof("array"));
this function will manipulate a given array.
the function manipulates the original array, it rerturnes true when finished.
because this function works with pointers its very time efficient and it needs no extra memory.

*)
(* @END_DECLARATION := '0' *)
(* this routine uses 32 bit access to gain speed *)
(* first access bytes till pointer is aligned for 32 bit access *)
temp := pt;
end := temp + UINT_TO_DWORD(size);
end32 := end - 3;
WHILE (pt < end) AND ((temp AND 16#00000003) > 0) DO
	pt^ := 0;
	pt := pt + 1;
	temp := temp + 1;
END_WHILE;
(* pointer is aligned, now copy 32 bits at a time *)
ptw := pt;
WHILE ptw < end32 DO (* *)
	ptw^ := 0;
	ptw := ptw + 4;
END_WHILE;
(* copy the remaining bytes in byte mode *)
pt := ptw;
WHILE pt < end DO
	pt^ := 0;
	pt := pt + 1;
END_WHILE;

_BUFFER_CLEAR := TRUE;

(* revision History
hm 	5. mar. 2008	rev 1.0
	original version

hm	16. mar. 2008	rev 1.1
	added type conversion to avoid warnings under codesys 30
	chaged type of input size to uint
	deleted unused variable i

hm	31. oct. 2008	rev 1.2
	corrected an error while routine would write outside of arrays

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Buffer Management' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _BUFFER_INIT : BOOL
VAR_INPUT
	PT : POINTER TO BYTE;
	SIZE : UINT;
	INIT : BYTE;
END_VAR
VAR
	ptw : POINTER TO DWORD;
	temp : DWORD;
	end, end32 : DWORD;
END_VAR


(*
version 1.2	31. oct. 2008
programmer 	hugo
tested by	oscat

this function will initialize a given array of byte with init.
the function needs to be called:	_buffer_init(adr("array"),sizeof("array"),init);
this function will manipulate a given array.
the function manipulates the original array, it rerturnes true when finished.
because this function works with pointers its very time efficient and it needs no extra memory.

*)
(* @END_DECLARATION := '0' *)
(* this routine uses 32 bit access to gain speed *)
(* first access bytes till pointer is aligned for 32 bit access *)
temp := pt;
end := temp + UINT_TO_DWORD(size);
end32 := end - 3;
WHILE (pt < end) AND ((temp AND 16#00000003) > 0) DO
	pt^ := init;
	pt := pt + 1;
	temp := temp + 1;
END_WHILE;
(* pointer is aligned, now copy 32 bits at a time *)
ptw := pt;
temp := SHL(BYTE_TO_DWORD(init),24) OR SHL(BYTE_TO_DWORD(init),16) OR SHL(BYTE_TO_DWORD(init),8) OR BYTE_TO_DWORD(init);
WHILE ptw < end32 DO
	ptw^ := temp;
	ptw := ptw + 4;
END_WHILE;
(* copy the remaining bytes in byte mode *)
pt := ptw;
WHILE pt < end DO
	pt^ := init;
	pt := pt + 1;
END_WHILE;

_BUFFER_INIT := TRUE;

(* revision History
hm 	5. mar. 2008	rev 1.0
	original version

hm	16. mar. 2008	rev 1.1
	added type conversion to avoid warnings under codesys 3.0
	chaged type of input size to uint.

hm	31. oct. 2008	rev 1.2
	corrected an error while routine would write outside of arrays

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Buffer Management' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _BUFFER_INSERT : INT
VAR_INPUT
	STR : STRING(STRING_LENGTH);
	POS : INT;
	PT : POINTER TO ARRAY[0..32767] OF BYTE;
	SIZE : UINT;
END_VAR
VAR
	end : INT;
	lx: INT;
	i: INT;
END_VAR


(*
version 1.5	2. jan. 2012
programmer 	hugo
tested by	oscat

this function will insert a string at a given position in a buffer.
the function needs to be called:	_buffer_init(adr("array"),sizeof("array"),init);
this function will manipulate a given array.
the function manipulates the original array, it returnes the next position after the input string when finished.
because this function works with pointers its very time efficient and it needs no extra memory.

*)
(* @END_DECLARATION := '0' *)
lx := LEN(str);
end := pos + lx;
(* first move the upper part of the buffer to make space for the string *)
FOR i := UINT_TO_INT(size) - 1 TO end BY -1 DO
	pt^[i] := pt^[i-lx];
END_FOR;
(* copy the string *)
_BUFFER_INSERT := _STRING_TO_BUFFER(str, pos, PT, size);


(* revision History
hm 	9. mar. 2008	rev 1.0
	original version

hm	16. mar. 2008	rev 1.1
	changed type of input size to uint

hm	13. may. 2008	rev 1.2
	changed type of pointer to array[1..32767]
	changed size of string to STRING_LENGTH

hm	23. mar. 2009	rev 1.3
	avoid writing to input pos

hm	12. nov. 2009	rev 1.4
	code optimized

hm 2. jan 2012	rev 1.5
	function now returns the next position after the input string

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Buffer Management' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION _BUFFER_UPPERCASE : BOOL
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF BYTE;
	SIZE : INT;
END_VAR
VAR
	pos: INT;
END_VAR


(*
version 1.0	20. jan. 2011
programmer 	hugo
tested by	oscat

this function will convert an array of byte into uppercase
the function needs to be called:	_buffer_clear(adr("array"),sizeof("array"));
this function will manipulate a given array.
the function manipulates the original array, it rerturnes true when finished.
because this function works with pointers its very time efficient and it needs no extra memory.

*)
(* @END_DECLARATION := '0' *)
pos := 0;
WHILE pos < size DO
	PT^[pos] := TO_UPPER(pt^[pos]);
	pos := pos + 1;
END_WHILE;

_BUFFER_UPPERCASE := TRUE;

(* revision History
hm 	20. jan. 2011	rev 1.0
	original version

*)
END_FUNCTION