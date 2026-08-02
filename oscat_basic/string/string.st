(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIN_TO_BYTE : BYTE
VAR_INPUT
	STR : STRING(8);
END_VAR
VAR
	i : INT;
	j : INT;
	temp : BYTE;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a binary string into a byte

*)
(* @END_DECLARATION := '0' *)
j := LEN(str);
FOR i := 1 TO j DO
	temp := SHL(temp,1);
	IF CODE(str,i) = 49 THEN temp := temp OR 1; END_IF;
END_FOR;
BIN_TO_BYTE := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIN_TO_DWORD : DWORD
VAR_INPUT
	STR : STRING(32);
END_VAR
VAR
	i : INT;
	j : INT;
	temp : DWORD;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a binary string into a dword

*)
(* @END_DECLARATION := '0' *)
j := LEN(str);
FOR i := 1 TO j DO
	temp := SHL(temp,1);
	IF CODE(str,i) = 49 THEN temp := temp OR 1; END_IF;
END_FOR;
BIN_TO_DWORD := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BYTE_TO_STRB : STRING(8)
VAR_INPUT
	IN : BYTE;
END_VAR
VAR
	i : INT;
	temp : STRING(8);
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a byte into a binary string

*)
(* @END_DECLARATION := '0' *)
FOR i := 1 TO 8 DO
	IF in.0 THEN
		temp := CONCAT('1',temp);
	ELSE
		temp := CONCAT('0',temp);
	END_IF;
	in := SHR(in,1);
END_FOR;
BYTE_TO_STRB := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function chr to chr_to_string

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BYTE_TO_STRH : STRING(2)
VAR_INPUT
	IN : BYTE;
END_VAR
VAR
	temp : STRING(2);
	i : INT;
	j : BYTE;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a byte into a hex string

*)
(* @END_DECLARATION := '0' *)
FOR i := 1 TO 2 DO
	j := in AND 16#0F;
	IF j < 10 THEN
		temp := CONCAT(CHR_TO_STRING(j+48),temp);
	ELSE
		temp := CONCAT(CHR_TO_STRING(j+55),temp);
	END_IF;
	in := SHR(in,4);
END_FOR;
BYTE_TO_STRH := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function chr to chr_to_string

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CAPITALIZE : STRING(STRING_LENGTH)
VAR_INPUT
	STR : STRING(STRING_LENGTH);
END_VAR
VAR
	i, j : INT;
	pt : POINTER TO ARRAY[1..STRING_LENGTH] OF BYTE;
	last: BYTE;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function capitalizes a string

*)
(* @END_DECLARATION := '0' *)
CAPITALIZE := str;
pt := ADR(CAPITALIZE);
j := LEN(str);
FOR i := 1 TO j DO
	IF i = 1 OR last = 32 THEN
		pt^[i] := TO_UPPER(pt^[i]);
	ELSE
		pt^[i] := TO_LOWER(pt^[i]);
	END_IF;
	last := pt^[i];
END_FOR;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CHARCODE : STRING(253)
VAR_INPUT
	C : BYTE;
END_VAR
VAR
	i : INT;
	j : INT;
	pos : INT;
	str : STRING(253);
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function returns the html code for a given character

*)
(* @END_DECLARATION := '0' *)
IF C > 159 AND C < 256 THEN
	i := (C - 160) / 32 + 1;
	j := (C - 160) MOD 32;
	str := setup.CHARNAMES[i];
	pos := 1;
	WHILE j > 0 DO
		pos := FINDP(str,';',pos) + 1;
		j := j - 1;
	END_WHILE;
	CHARCODE := MID(str,FINDP(str,';',pos) - pos + 1, pos-1);
END_IF;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CHARNAME : BYTE
VAR_INPUT
	STR : STRING(10);
END_VAR
VAR
	i : INT;
	j : INT;
	pos : INT;
	str2 : STRING(253);
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the character for a given html code

*)
(* @END_DECLARATION := '0' *)
FOR i := 1 TO 4 DO
	str2 := setup.CHARNAMES[i];
	pos := FINDP(str2,str,1);
	IF pos > 0 THEN
		j := 0;
		WHILE pos > 1 DO
			pos := FINDB(str2,';',pos-1);
			j := j + 1;
		END_WHILE;
		CHARNAME := INT_TO_BYTE(160 + (i-1)*32 + j -1);
		RETURN;
	END_IF;
END_FOR;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CHR_TO_STRING : STRING(1)
VAR_INPUT
	C : BYTE;
END_VAR


(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a byte into a string of length 1

*)
(* @END_DECLARATION := '0' *)
CHR_TO_STRING := BYTE_TO_STRING(C);


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed name from chr to chr_to_string

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CLEAN : STRING(STRING_LENGTH)
VAR_INPUT
	STR : STRING(STRING_LENGTH);
END_VAR
VAR
	i, j, k : INT;
	pt : POINTER TO ARRAY[1..STRING_LENGTH] OF BYTE;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function cleans a string from all non printable characters

*)
(* @END_DECLARATION := '0' *)
CLEAN := str;
pt := ADR(CLEAN);
j := LEN(str);
k := 1;
FOR i := 1 TO j DO
	IF IS_NCC(pt^[i]) THEN
		pt^[k] := pt^[i];
		k := k + 1;
	END_IF;
END_FOR;
pt^[k] := 0;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CODE : BYTE
VAR_INPUT
	STR : STRING(STRING_LENGTH);
	POS : INT;
END_VAR


(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function returns the ascii code of a character in a string at a given position

*)
(* @END_DECLARATION := '0' *)
STRING_TO_BYTE(str,pos,CODE);


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed name from asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION COUNT_CHAR : INT
VAR_INPUT
	STR : STRING(STRING_LENGTH);
	C : BYTE;
END_VAR
VAR
	i, j : INT;
	pt : POINTER TO ARRAY[1..STRING_LENGTH] OF BYTE;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function counts the number of characters C in a string STR

*)
(* @END_DECLARATION := '0' *)
pt := ADR(str);
j := LEN(str);
FOR i := 1 TO j DO
	IF pt^[i] = C THEN COUNT_CHAR := COUNT_CHAR + 1; END_IF;
END_FOR;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION COUNT_SUBSTRING : INT
VAR_INPUT
	STR : STRING(STRING_LENGTH);
	SUB : STRING(STRING_LENGTH);
END_VAR
VAR
	pos : INT := 1;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function counts the number of substrings SUB in a string STR

*)
(* @END_DECLARATION := '0' *)
REPEAT
	pos := FIND(str,sub);
	IF pos > 0 THEN
		COUNT_SUBSTRING := COUNT_SUBSTRING + 1;
		str := DELETE(str,1,pos);
	END_IF;
UNTIL pos = 0 END_REPEAT;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DEC_TO_BYTE : BYTE
VAR_INPUT
	STR : STRING(3);
END_VAR
VAR
	i, j : INT;
	temp : INT;
	neg : BOOL;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a decimal string into a byte

*)
(* @END_DECLARATION := '0' *)
j := LEN(str);
FOR i := 1 TO j DO
	IF i = 1 AND CODE(str,1) = 45 THEN
		neg := TRUE;
	ELSE
		temp := temp * 10 + (CODE(str,i) - 48);
	END_IF;
END_FOR;
IF neg THEN temp := -temp; END_IF;
DEC_TO_BYTE := INT_TO_BYTE(temp);


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DEC_TO_DWORD : DWORD
VAR_INPUT
	STR : STRING(10);
END_VAR
VAR
	i, j : INT;
	temp : DWORD;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a decimal string into a dword

*)
(* @END_DECLARATION := '0' *)
j := LEN(str);
FOR i := 1 TO j DO
	temp := temp * 10 + (CODE(str,i) - 48);
END_FOR;
DEC_TO_DWORD := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DEC_TO_INT : INT
VAR_INPUT
	STR : STRING(6);
END_VAR
VAR
	i, j : INT;
	temp : INT;
	neg : BOOL;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a decimal string into an integer

*)
(* @END_DECLARATION := '0' *)
j := LEN(str);
FOR i := 1 TO j DO
	IF i = 1 AND CODE(str,1) = 45 THEN
		neg := TRUE;
	ELSE
		temp := temp * 10 + (CODE(str,i) - 48);
	END_IF;
END_FOR;
IF neg THEN temp := -temp; END_IF;
DEC_TO_INT := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DEL_CHARS : STRING(STRING_LENGTH)
VAR_INPUT
	STR : STRING(STRING_LENGTH);
	C : BYTE;
END_VAR
VAR
	i, j, k : INT;
	pt : POINTER TO ARRAY[1..STRING_LENGTH] OF BYTE;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function deletes all characters C from a string STR

*)
(* @END_DECLARATION := '0' *)
DEL_CHARS := str;
pt := ADR(DEL_CHARS);
j := LEN(str);
k := 1;
FOR i := 1 TO j DO
	IF pt^[i] <> C THEN
		pt^[k] := pt^[i];
		k := k + 1;
	END_IF;
END_FOR;
pt^[k] := 0;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DT_TO_STRF : STRING(STRING_LENGTH)
VAR_INPUT
	DT_IN : DT;
	FORMAT : STRING(STRING_LENGTH);
END_VAR
VAR
	sdt_temp : SDT;
	temp : STRING(STRING_LENGTH);
	pos : INT;
	c : STRING(1);
	i : INT;
	j : INT;
	X : INT;
END_VAR

(*
version 1.2	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a DT into a string with a given format

*)
(* @END_DECLARATION := '0' *)
sdt_temp := DT_TO_SDT(DT_IN);
j := LEN(format);
FOR i := 1 TO j DO
	c := MID(format,1,i);
	IF c = '%' THEN
		i := i + 1;
		c := MID(format,1,i);
		CASE c OF
		'd':	temp := CONCAT(temp,RIGHT(CONCAT('0',INT_TO_STRING(sdt_temp.DAY)),2));
		'H':	temp := CONCAT(temp,RIGHT(CONCAT('0',INT_TO_STRING(sdt_temp.HOUR)),2));
		'j':	temp := CONCAT(temp,RIGHT(CONCAT('00',INT_TO_STRING(DAY_OF_YEAR(DT_IN))),3));
		'm':	temp := CONCAT(temp,RIGHT(CONCAT('0',INT_TO_STRING(sdt_temp.MONTH)),2));
		'M':	temp := CONCAT(temp,RIGHT(CONCAT('0',INT_TO_STRING(sdt_temp.MINUTE)),2));
		'S':	temp := CONCAT(temp,RIGHT(CONCAT('0',INT_TO_STRING(sdt_temp.SECOND)),2));
		'w':	temp := CONCAT(temp,INT_TO_STRING(sdt_temp.WEEKDAY-1));
		'W':	temp := CONCAT(temp,RIGHT(CONCAT('0',INT_TO_STRING(WORK_WEEK(DT_IN))),2));
		'y':	temp := CONCAT(temp,RIGHT(INT_TO_STRING(sdt_temp.YEAR),2));
		'Y':	temp := CONCAT(temp,INT_TO_STRING(sdt_temp.YEAR));
		'%':	temp := CONCAT(temp,'%');
		ELSE
			temp := CONCAT(temp,c);
		END_CASE;
	ELSE
		temp := CONCAT(temp,c);
	END_IF;
END_FOR;
DT_TO_STRF := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	28. jun. 2008	rev 1.1
	corrected an error with day of year

hm	17. dec. 2008	rev 1.2
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DWORD_TO_STRB : STRING(32)
VAR_INPUT
	IN : DWORD;
END_VAR
VAR
	i : INT;
	temp : STRING(32);
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a dword into a binary string

*)
(* @END_DECLARATION := '0' *)
FOR i := 1 TO 32 DO
	IF in.0 THEN
		temp := CONCAT('1',temp);
	ELSE
		temp := CONCAT('0',temp);
	END_IF;
	in := SHR(in,1);
END_FOR;
DWORD_TO_STRB := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function chr to chr_to_string

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DWORD_TO_STRF : STRING(STRING_LENGTH)
VAR_INPUT
	IN : DWORD;
	FORMAT : STRING(STRING_LENGTH);
END_VAR
VAR
	temp : STRING(STRING_LENGTH);
	pos : INT;
	c : STRING(1);
	i : INT;
	j : INT;
	X : INT;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a DWORD into a string with a given format

*)
(* @END_DECLARATION := '0' *)
j := LEN(format);
FOR i := 1 TO j DO
	c := MID(format,1,i);
	IF c = '%' THEN
		i := i + 1;
		c := MID(format,1,i);
		CASE c OF
		'b':	temp := CONCAT(temp,DWORD_TO_STRB(IN));
		'd':	temp := CONCAT(temp,DWORD_TO_STRING(IN));
		'h':	temp := CONCAT(temp,DWORD_TO_STRH(IN));
		'%':	temp := CONCAT(temp,'%');
		ELSE
			temp := CONCAT(temp,c);
		END_CASE;
	ELSE
		temp := CONCAT(temp,c);
	END_IF;
END_FOR;
DWORD_TO_STRF := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function asc to code

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/String' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DWORD_TO_STRH : STRING(8)
VAR_INPUT
	IN : DWORD;
END_VAR
VAR
	temp : STRING(8);
	i : INT;
	j : BYTE;
END_VAR

(*
version 1.1	17. dec. 2008
programmer 	hugo
tested by	oscat

this function converts a dword into a hex string

*)
(* @END_DECLARATION := '0' *)
FOR i := 1 TO 8 DO
	j := DWORD_TO_BYTE(in AND 16#0F);
	IF j < 10 THEN
		temp := CONCAT(CHR_TO_STRING(j+48),temp);
	ELSE
		temp := CONCAT(CHR_TO_STRING(j+55),temp);
	END_IF;
	in := SHR(in,4);
END_FOR;
DWORD_TO_STRH := temp;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	17. dec. 2008	rev 1.1
	changed function chr to chr_to_string

*)
END_FUNCTION