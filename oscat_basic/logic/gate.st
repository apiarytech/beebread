(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BCDC_TO_INT : INT
VAR_INPUT
	IN : BYTE;
END_VAR


(*
version 1.1 30. jun. 2008
programmer 	hugo
tested by	oscat

this function converts a two digit bcd number into an integer.

*)

(* @END_DECLARATION := '0' *)
BCDC_TO_INT := (in AND 16#0F) + (SHR(in,4) * 10);

(* revision history
hm	13.12.2007	rev 1.0
	original version

hm	30.6.2008	rev 1.1
	changed name BCD_TO_INT to BCDC_TO_INT to avoid collision with util.lib

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_COUNT : INT
VAR_INPUT
	IN : DWORD;
END_VAR


(*
version 1.1	10 sep 2007
programmer 	hugo
tested by		tobias

BIT_COUNT counts the amount True of bits in a dword.
for exabple: bit_count(3) returns 2 because two bits (bits 0 and 1) are true and all others are false. 

*)
(* @END_DECLARATION := '0' *)
WHILE in > 0 DO
	IF in.0 THEN Bit_Count := Bit_Count + 1; END_IF;
	in := SHR(in,1);
END_WHILE;



(* revision history
5.7.2007	rev 1.0		original version

10.9.2007	rev 1.1		hm
	changed algorithm for better performace
	the execution time has reduced by a factor of 5
	deleted unused variable temp
*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_LOAD_B : BYTE
VAR_INPUT
	IN : BYTE;
	VAL : BOOL;
	POS : INT;
END_VAR
VAR CONSTANT
	dat : BYTE := 1;
END_VAR

(*
version 1.1	16. mar 2008
programmer 	hugo
tested by	tobias

this function loads a bit into a byte at position pos.
 *)
(* @END_DECLARATION := '0' *)
IF VAL THEN
	BIT_LOAD_B := in OR SHL(dat,pos);
ELSE
	BIT_LOAD_B := in AND (NOT SHL(dat,pos));
END_IF;


(* revision history
hm	29. feb 2008	rev 1.0
	original version

hm	16. mar 2008	rev 1.1
	change input bit to val for compatibility reasons
*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_LOAD_B2 : BYTE
VAR_INPUT
	I : BYTE;
	D : BOOL;
	P : INT;
	N : INT;
END_VAR


(*
version 1.0 18. oct. 2008
programmer 	hugo
tested by	tobias

this function loads N bits of D at pos P in Byte I

*)
(* @END_DECLARATION := '0' *)
IF D THEN
	BIT_LOAD_B2 := ROL(SHR(BYTE#255, 8 - N) OR ROR(i, P), P);
ELSE
	BIT_LOAD_B2 := ROL(SHL(BYTE#255, N) AND ROR(I, P), P);
END_IF;



(* revision history
hm	18. oct. 2008	rev 1.0
	original version


*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_LOAD_DW : DWORD
VAR_INPUT
	IN : DWORD;
	VAL : BOOL;
	POS : INT;
END_VAR
VAR CONSTANT
	dat : DWORD := 1;
END_VAR

(*
version 1.1	16. mar 2008
programmer 	hugo
tested by	tobias

this function loads a bit into a DWord at position pos.
 *)
(* @END_DECLARATION := '0' *)
IF val THEN
	BIT_LOAD_DW := in OR SHL(dat,pos);
ELSE
	BIT_LOAD_DW := in AND (NOT SHL(dat,pos));
END_IF;


(* revision history
hm	29. feb 2008	rev 1.0
	original version

hm	16. mar 2008	rev 1.1
	change input bit to val for compatibility reasons
*)


END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_LOAD_DW2 : DWORD
VAR_INPUT
	I : DWORD;
	D : BOOL;
	P : INT;
	N : INT;
END_VAR


(*
version 1.0 18. oct. 2008
programmer 	hugo
tested by	tobias

this function loads N bits of D at pos P in DWORD I

*)
(* @END_DECLARATION := '0' *)
IF D THEN
	BIT_LOAD_DW2 := ROL(SHR(DWORD#4294967295, 32 - N) OR ROR(i, P), P);
ELSE
	BIT_LOAD_DW2 := ROL(SHL(DWORD#4294967295, N) AND ROR(I, P), P);
END_IF;



(* revision history
hm	18. oct. 2008	rev 1.0
	original version


*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_LOAD_W : WORD
VAR_INPUT
	IN : WORD;
	VAL : BOOL;
	POS : INT;
END_VAR
VAR CONSTANT
	dat : WORD := 1;
END_VAR

(*
version 1.1	16. mar 2008
programmer 	hugo
tested by	tobias

this function loads a bit into a Word at position pos.
 *)
(* @END_DECLARATION := '0' *)
IF val THEN
	BIT_LOAD_W := in OR SHL(dat,pos);
ELSE
	BIT_LOAD_W := in AND (NOT SHL(dat,pos));
END_IF;


(* revision history
hm	29. feb 2008	rev 1.0
	original version

hm	16. mar 2008	rev 1.1
	change input bit to val for compatibility reasons
*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_LOAD_W2 : WORD
VAR_INPUT
	I : WORD;
	D : BOOL;
	P : INT;
	N : INT;
END_VAR


(*
version 1.0 18. oct. 2008
programmer 	hugo
tested by	tobias

this function loads N bits of D at pos P in WORD I

*)
(* @END_DECLARATION := '0' *)
IF D THEN
	BIT_LOAD_W2 := ROL(SHR(WORD#65535, 16 - N) OR ROR(i, P), P);
ELSE
	BIT_LOAD_W2 := ROL(SHL(WORD#65535, N) AND ROR(I, P), P);
END_IF;



(* revision history
hm	18. oct. 2008	rev 1.0
	original version


*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_OF_DWORD : BOOL
VAR_INPUT
	in : DWORD;
	N : INT;
END_VAR


(*
version 1.2	6. jun. 2008
programmer 	hugo
tested by	tobias

this function extracts a single bit from the nth position from right (right is lowest bit)
the lowest Bit (Bit0 from in) is selected with N=0.
 *)
(* @END_DECLARATION := '0' *)
BIT_OF_DWORD := (SHR(in,N) AND 16#00000001) > 0;


(* old code used before rev 1.1
temp := SHR(in,n);
Bit_of_Dword := temp.0;
*)


(* revision history
hm	4. aug. 2006	rev 1.0
	original version

hm	29. feb 2008	rev 1.1
	improved performance

hm	6. jun. 2008	rev 1.2
	changes type of input N from byte to int
*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_TOGGLE_B : BYTE
VAR_INPUT
	IN : BYTE;
	POS : INT;
END_VAR


(*
version 1.0 18. oct. 2008
programmer 	hugo
tested by	tobias

this function toggles a bit of a BYTE at position pos.

*)
(* @END_DECLARATION := '0' *)
BIT_TOGGLE_B := SHL(BYTE#1, POS) XOR IN;




(* revision history
hm	18. oct. 2008	rev 1.0
	original version


*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_TOGGLE_DW : DWORD
VAR_INPUT
	IN : DWORD;
	POS : INT;
END_VAR


(*
version 1.0 18. oct. 2008
programmer 	hugo
tested by	tobias

this function toggles a bit of a WORD at position pos.

*)
(* @END_DECLARATION := '0' *)
BIT_TOGGLE_DW := SHL(DWORD#1, POS) XOR IN;



(* revision history
hm	18. oct. 2008	rev 1.0
	original version


*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BIT_TOGGLE_W : WORD
VAR_INPUT
	IN : WORD;
	POS : INT;
END_VAR


(*
version 1.0 18. oct. 2008
programmer 	hugo
tested by	tobias

this function toggles a bit of a WORD at position pos.

*)
(* @END_DECLARATION := '0' *)
BIT_TOGGLE_W := SHL(WORD#1, POS) XOR IN;



(* revision history
hm	18. oct. 2008	rev 1.0
	original version


*)

END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BYTE_OF_BIT : BYTE
VAR_INPUT
	B0:BOOL;
	B1:BOOL;
	B2:BOOL;
	B3:BOOL;
	B4:BOOL;
	B5:BOOL;
	B6:BOOL;
	B7:BOOL;
END_VAR


(*
version 1.1	 18. feb 2006
programmer 	hugo
tested by		tobias

this function creates a byte from 8 individual bits

*)
(* @END_DECLARATION := '0' *)
Byte_of_bit := SHL(SHL(SHL(SHL(SHL(SHL(SHL(BOOL_TO_BYTE(B7),1) OR BOOL_TO_BYTE(B6),1) OR BOOL_TO_BYTE(B5),1) OR BOOL_TO_BYTE(B4),1)
	OR BOOL_TO_BYTE(B3),1) OR BOOL_TO_BYTE(B2),1) OR BOOL_TO_BYTE(B1),1) OR BOOL_TO_BYTE(B0);


(* revision history

hm	4. aug 2006		rev 1.0
	original version

hm	18. feb. 2008	rev 1.1
	improved performance

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BYTE_OF_DWORD : BYTE
VAR_INPUT
	in : DWORD;
	N : BYTE;
END_VAR


(*
version 1.2	30. oct. 2008
programmer 	hugo
tested by	oscat

this function extracts a single byte from the nth position from right (right is lowest byte) 
the lower byte (starting with Bit0 from in) is selected with N=0.
*)

(* @END_DECLARATION := '0' *)
BYTE_OF_DWORD := DWORD_TO_BYTE(SHR(in,SHL(n,3)));

(* revision history
hm	17. jan 2007	rev 1.0
	original version

hm	2. jan 2008		rev 1.1
	improved performance

hm	30. oct. 2008	rev 1.2
	improved performance
*)


END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK BYTE_TO_BITS
VAR_INPUT
	IN: BYTE;
END_VAR
VAR_OUTPUT
	B0: BOOL;
	B1: BOOL;
	B2: BOOL;
	B3: BOOL;
	B4: BOOL;
	B5: BOOL;
	B6: BOOL;
	B7: BOOL;
END_VAR


(*
version 1.1	16. mar. 2008
programmer 	hugo
tested by		tobias

this Function Block extracts the 8 Bits from a byte


*)
(* @END_DECLARATION := '0' *)
B0 := IN.0;
B1 := IN.1;
B2 := IN.2;
B3 := IN.3;
B4 := IN.4;
B5 := IN.5;
B6 := IN.6;
B7 := IN.7;

(* revision history
hm		4. aug 2006	rev 1.0
	original version

hm		16. mar 2008	rev 1.1
	renamed from byte_to_bit to byte_to_bits for compatibility reasons

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/gate logic' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION BYTE_TO_GRAY : BYTE
VAR_INPUT
	IN : BYTE;
END_VAR


(*
version 1.0	9. nov. 2009
programmer 	hugo
tested by		oscat

this function converts a binary to gray code

*)

(* @END_DECLARATION := '0' *)
BYTE_TO_GRAY := IN XOR SHR(IN,1);

(* revision history
hm	9. nov. 2009	rev 1.0
	original version

*)
END_FUNCTION