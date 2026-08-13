FUNCTION LIST_ADD : BOOL
VAR_INPUT
	SEP : BYTE;
	INS : STRING(LIST_LENGTH);
END_VAR
VAR_IN_OUT
	LIST : STRING(LIST_LENGTH);
END_VAR
VAR
	sx : STRING(1);
END_VAR


(*
version 2.0		21. mar. 2011
programmer 	hugo
tested by		kurt

LIST_ADD hngt ein weiteres element ans ende einer liste.
die liste liegt als string von elementen vor die mit den zeichen SEP am anfang makiert sind.

*)
sx := CHR_TO_STRING(SEP);			(* convert separation character to string *)
INS := CONCAT(sx, INS);					(* start element with separation character *)

IF LEN(LIST) + LEN(INS) > LIST_LENGTH THEN
	LIST_ADD := FALSE;					(* return false if element does not fit *)
ELSE
	LIST := CONCAT(LIST, INS);
	LIST_ADD := TRUE;
END_IF;


(* revision histroy

hm	21. mar. 2011	rev 2.0
	original release
*)	

END_FUNCTION


FUNCTION LIST_CLEAN : BOOL
VAR_INPUT
	SEP : BYTE;
END_VAR
VAR_IN_OUT
	LIST : STRING(LIST_LENGTH);
END_VAR
VAR
	pt : POINTER TO ARRAY[1..LIST_LENGTH] OF BYTE;
	read : INT := 1;
	write : INT := 1;
	last: BYTE;
	c: BYTE;
END_VAR


(*
version 2.0	21. mar. 2011
programmer 	hugo
tested by		oscat

LIST_CLEAN bereinigt eine liste von leeren Elementen.
die liste liegt als string von elementen vor die mit den zeichen SEP am anfang markiert sind.

*)
pt := ADR(LIST);

FOR read := 1 TO LIST_LENGTH DO
	c := pt^[read];								(* read character from list *)
	IF c = 0 THEN EXIT;						(* exit the loop if character definbes end of string *)
	ELSIF c <> SEP OR sep <> last THEN	(* copy element from read to write position unless a double sep character is present *)
		pt^[write] := c;
		write := write + 1;
	END_IF;
	last := c;									(* remember last character *)
END_FOR;

IF last = SEP THEN write := write - 1; END_IF;		(* if last character is sep then delete empty element at end *)
IF write <= STRING_LENGTH THEN pt^[write] := 0;	 END_IF;		(* terminate string with 0 *)

LIST_CLEAN := TRUE;						(* retrun TRUE *)


(* revision histroy
hm	28. jun. 2008	rev 1.0
	original release

hm	19. jan. 2011	rev 1.1
	changed string(255) to string(LIST_LENGTH)

hm	21. mar. 2011	rev 2.0
	all elements start with SEP

*)	

END_FUNCTION


FUNCTION LIST_GET : STRING(LIST_LENGTH)
VAR_INPUT
	SEP : BYTE;
	POS : INT;
END_VAR
VAR_IN_OUT
	LIST : STRING(LIST_LENGTH);
END_VAR
VAR
	i : INT := 1;
	o : INT := 1;
	pt : POINTER TO ARRAY[1..LIST_LENGTH] OF BYTE;
	po : POINTER TO ARRAY[1..LIST_LENGTH] OF BYTE;
	cnt: INT := 0;
	c: BYTE;
END_VAR


(*
version 2.0	21. mar. 2011
programmer 	hugo
tested by		oscat

LIST_GET liefert das element an der stelle pos einer liste.
die liste liegt als string von elementen vor die mit den zeichen SEP am anfang makiert sind.

*)
pt := ADR(LIST);							(* load pointers *)
po := ADR(LIST_GET);

REPEAT
	c := pt^[i];
	IF cnt = pos THEN						(* copy element to output string *)
		IF c = SEP  THEN EXIT; END_IF;		(* next element stop copy and finish *)
		po^[o] := c;
		o := o + 1;
	ELSIF c = SEP THEN
		cnt := cnt + 1;
	END_IF;
	i := i + 1;
UNTIL
	(i = LIST_LENGTH) OR (c = 0)
END_REPEAT;
po^[o] := 0;						(* close the output string *)



(* revision histroy
hm	20. jun. 2008	rev 1.0
	original release

hm	19. jan. 2011	rev 1.1
	changed string(255) to strring(LIST_LENGTH)

hm	21. mar. 2011	rev 2.0
	all list elements start with SEP
*)	

END_FUNCTION


FUNCTION LIST_INSERT : BOOL
VAR_INPUT
	SEP : BYTE;
	POS : INT;
	INS : STRING(LIST_LENGTH);
END_VAR
VAR_IN_OUT
	LIST : STRING(LIST_LENGTH);
END_VAR
VAR
	pt : POINTER TO ARRAY[1..LIST_LENGTH] OF BYTE;
	read : INT := 1;
	cnt : INT := 1;
	sx : STRING(1);
END_VAR


(*
version 2.0	21. mar. 2011
programmer 	hugo
tested by		oscat

LIST_INSERT setzt ein element an der stelle pos in eine liste ein.
die liste liegt als string von elementen vor die mit den zeichen SEP am anfang markiert sind.

*)
(* load pointers *)
pt := ADR(LIST);

IF LEN(ins) + 1 + LEN(LIST) > LIST_LENGTH THEN
	LIST_INSERT := FALSE;
ELSE
	LIST_INSERT := TRUE;
	(* search for position *)
	WHILE read < LIST_LENGTH DO
		IF cnt >= POS THEN
			sx := CHR_TO_STRING(SEP);
			INS := CONCAT(sx, INS);
			LIST := INSERT(LIST, INS, read - 1);
			LIST_INSERT := TRUE;
			RETURN;
		END_IF;
		IF pt^[read]  = 0 THEN				(* list is too short add empty element *)
			pt^[read] := SEP;
			pt^[read + 1] := 0;
		END_IF;
		read := read + 1;
		IF pt^[read] = SEP OR pt^[read] = 0 THEN cnt := cnt + 1; END_IF;
	END_WHILE;
END_IF;



(* revision histroy
hm	28. jun. 2008	rev 1.0
	original release

hm	17. dec. 2008	rev 1.1
	changes name of function chr to chr_to_string

hm	19. jan. 2011	rev 1.2
	changed string(255) to string(LIST_LENGTH)

hm	21. mar. 2011	rev 2.0
	all list elements start with SEP
	
*)	

END_FUNCTION