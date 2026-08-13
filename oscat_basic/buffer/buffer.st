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
END_FUNCTION


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
END_FUNCTION


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

lx := LEN(str);
end := pos + lx;
(* first move the upper part of the buffer to make space for the string *)
FOR i := UINT_TO_INT(size) - 1 TO end BY -1 DO
	pt^[i] := pt^[i-lx];
END_FOR;
(* copy the string *)
_BUFFER_INSERT := _STRING_TO_BUFFER(str, pos, PT, size);

END_FUNCTION


FUNCTION _BUFFER_UPPERCASE : BOOL
VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF BYTE;
	SIZE : INT;
END_VAR
VAR
	pos: INT;
END_VAR

pos := 0;
WHILE pos < size DO
	PT^[pos] := TO_UPPER(pt^[pos]);
	pos := pos + 1;
END_WHILE;

_BUFFER_UPPERCASE := TRUE;

END_FUNCTION