(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/Others' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION CRC_GEN : DWORD

VAR_INPUT
	PT : POINTER TO ARRAY[0..32000] OF BYTE;
	SIZE : INT;
	PL : INT;
	PN : DWORD;
	INIT : DWORD;
	REV_IN : BOOL;
	REV_OUT : BOOL;
	XOR_OUT : DWORD;
END_VAR
VAR
	pos : INT;
	shift : INT;
	dx: BYTE;
	bits: INT;
END_VAR

(*
version 2.1	3. feb. 2021
programmer 	hugo
tested by		tobias

CRC_GEN generates a CRC checksum from a block of data and returns  the checksum in a DWORD to be connected to the data for transmission.
the CRC Polynom is specified with the config variable PN and the length of the Polynom is specified by PL
A Polynom x4 + X + 1 is represented by 0011 with length 4, the highest order bit is not specified at all..
The input data is an array of byte of any size, the function is called by CRC_GEN(ADR(array),SIZEOF(array), ....).

*)

(* @END_DECLARATION := '0' *)
(* align polygon *)
shift := 32 - PL;
PN := SHL(PN, shift);

(* load first 4 bytes into register minimum message size is 4 bytes
	for smaller messages fill with 0#s at the beginning*)
FOR pos := 0 TO 3 DO
	IF REV_IN THEN CRC_GEN := SHL(CRC_GEN, 8) OR REVERSE(PT^[pos]); ELSE CRC_GEN := SHL(CRC_GEN, 8) OR PT^[pos]; END_IF;
END_FOR;
pos := 4;

(* xor with init value *)
CRC_GEN := CRC_GEN XOR SHL(init, shift);

(* calculate CRC for each byte *)
WHILE pos < size DO
	IF REV_IN THEN DX := REVERSE(PT^[pos]); ELSE DX := PT^[pos]; END_IF;
	pos := pos + 1;
	(* crc calculation for one byte *)
	FOR bits := 0 TO 7 DO
		IF CRC_GEN.31 THEN
			CRC_GEN := (SHL(CRC_GEN, 1) OR BOOL_TO_DWORD(DX.7)) XOR PN;
		ELSE
			CRC_GEN := SHL(CRC_GEN, 1) OR BOOL_TO_DWORD(DX.7);
		END_IF;
		dx := SHL(dx, 1);
	END_FOR;
END_WHILE;

(* all bytes are processed, need to finish the registers 32 bits *)
FOR bits := 0 TO 31 DO
	IF CRC_GEN.31 THEN
		CRC_GEN := (SHL(CRC_GEN, 1) ) XOR PN;
	ELSE
		CRC_GEN := SHL(CRC_GEN, 1);
	END_IF;
END_FOR;

(* final XOR *)
CRC_GEN := SHR(CRC_GEN, shift) XOR XOR_OUT;

(* reverse the crc_out put if necessary *)
IF REV_OUT THEN CRC_GEN := REFLECT(CRC_GEN, PL); END_IF;

(* suggested optimization by deninis klopke
die betreffs genannte Funktion kann folgendermaen optimiert werden:
In Zeilen 33 und 35 des Implementierungsteils kann jeweils der Teil OR BOOL_TO_DWORD(DX.7) entfallen, da nach 8maligem Linksshift von DX in Zeile 26 der Wert von DX = BYTE#0 ist. Eine weitere Veroderung mit diesem Wert hat keine Auswirkungen auf das Ergebnis mehr.
*)

(* typical crc polynoms

CRC-4-ITU 			x4 + x + 1 					(ITU G.704, p. 12) 	0x3 or 0xC (0x9)
CRC-5-ITU 			x5 + x4 + x2 + 1 			(ITU G.704, p. 9) 	0x15 or 0x15 (0x0B) Bluetooth
CRC-5-USB 			x5 + x2 + 1 				(use: USB token packets) 	0x05 or 0x14 (0x9)
CRC-6-ITU 			x6 + x + 1 					(ITU G.704, p. 3) 	0x03 or 0x30 (0x21)
CRC-7 				x7 + x3 + 1 				(use: telecom systems, MMC) 	0x09 or 0x48 (0x11)
CRC-8-ATM 			x8 + x2 + x + 1 				(use: ATM HEC) 	0x07 or 0xE0 (0xC1)
CRC-8-CCITT 		x8 + x7 + x3 + x2 + 1 		(use: 1-Wire bus) 	0x8D or 0xB1 (0x63)
CRC-8-Dallas/Maxim 	x8 + x5 + x4 + 1 			(use: 1-Wire bus) 	0x31 or 0x8C (0x19)
CRC-8 				x8 + x7 + x6 + x4 + x2 + 1 	0xD5 or 0xAB (0x57)
CRC-8-SAE J1850 	x8 + x4 + x3 + x2 + 1 		0x1D or 0xB8
CRC-10 				x10 + x9 + x5 + x4 + x + 1 	0x233 or 0x331 (0x263)
CRC-12 				x12 + x11 + x3 + x2 + x + 1 (use: telecom systems) 	0x80F or 0xF01 (0xE03)
CRC-15-CAN 			x15 + x14 + x10 + x8 + x7 + x4 + x3 + 1 	0x4599 or 0x4CD1 (0x19A3)
CRC-16-Fletcher 	Not a CRC; see Fletcher's checksum 	Used in Adler-32 A & B CRCs
CRC-16-CCITT 	x16 + x12 + x5 + 1 (XMODEM,X.25, V.41, Bluetooth, PPP, IrDA; known as "CRC-CCITT") 	0x1021 or 0x8408 (0x0811)
CRC-16-IBM 	x16 + x15 + x2 + 1 (USB, many others; also known as "CRC-16") 	0x8005 or 0xA001 (0x4003)
CRC-24-Radix-64 	x24 + x23 + x18 + x17 + x14 + x11 + x10 + x7 + x6 + x5 + x4 + x3 + x + 1 	0x864CFB or 0xDF3261 (0xBE64C3)
CRC-32-Adler 	Not a CRC; see Adler-32 	See Adler-32
CRC-32-MPEG2 	x32 + x26 + x23 + x22 + x16 + x12 + x11 + x10 + x8 + x7 + x5 + x4 + x2 + x + 1 	0x04C11DB7 or 0xEDB88320 (0xDB710641) Also used in IEEE 802.3
CRC-32-IEEE 802.3 	x32 + x26 + x23 + x22 + x16 + x12 + x11 + x10 + x8 + x7 + x5 + x4 + x2 + x + 1 (V.42) 	0x04C11DB7 or 0xEDB88320 (0xDB710641)
CRC-32C (Castagnoli) 	x32 + x28 + x27 + x26 + x25 + x23 + x22 + x20 + x19 + x18 + x14 + x13 + x11 + x10 + x9 + x8 + x6 + 1 	0x1EDC6F41 or 0x82F63B78 (0x05EC76F1)
CRC-64-ISO 	x64 + x4 + x3 + x + 1 (use: ISO 3309) 	0x000000000000001B or 0xD800000000000000 (0xB000000000000001)
CRC-64-ECMA-182 	x64 + x62 + x57 + x55 + x54 + x53 + x52 + x47 + x46 + x45 + x40 + x39 + x38 + x37 + x35 + x33 + x32 + x31 + x29 + x27 + x24 + x23 + x22 + x21 + x19 + x17 + x13 + x12 + x10 + x9 + x7 + x4 + x + 1
(as described in ECMA-182 p.63) 	0x42F0E1EBA9EA3693 or 0xC96C5795D7870F42 (0x92D8AF2BAF0E1E85)
*)



(* revision history

hm	9.6.2007		rev 1.0		
	original version 

hm	11.9.2007		rev 1.1
	deleted unused variable i

hm	9. oct 2007 	rev 1.2
	added init code for crc and xor_out
	added refelct in and reflect_out (rev_in und Rev_out)

hm	2. jan 2008	rev 1.3
	small changes for performance improvements

hm	16. mar. 2008	rev 1.4
	changed type of input size to uint

hm	10. mar. 2009	rev 1.5
	removed nested comments

hm	16. jan. 2011	rev 2.0
	new version

dk 3. feb 2021		rev 2.1
	improved code by dennis koepke

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/Others' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK MATRIX
VAR_INPUT
	x1, x2, x3, x4, X5 : BOOL;
END_VAR
VAR_INPUT CONSTANT
	Release : BOOL;
END_VAR
VAR_OUTPUT
	code : BYTE;
	TP : BOOL;
	y1: BOOL := TRUE;
	y2, y3, y4 : BOOL;
END_VAR
VAR
	line : BYTE;
	X : ARRAY[0..3] OF BYTE; (* scan line inputs *)
	L : ARRAY[0..3] OF BYTE; (* scan line status *)
	i: INT;
	temp: BYTE;
END_VAR

(*
version 1.3	26 oct. 2008
programmer 	hugo
tested by	oscat

MATRIX is a matrix keyboard encoder for 4 rows and up to 5 columns.
matrix can send a code while a key is pressed and it sends another code while a key is released when the setup variable 'release' is set to true.
the setup variable release is false the code is only sent when the key is pressed.
the output byte holds the 5 columns in the lowest bits 0..2  and the row number in bits 4..6, while bit 7 is true for a key pressed and false for a key released.

*)
(* @END_DECLARATION := '0' *)
TP := FALSE;
CODE := 0;
(* read scan lines *)
X[line].0 := X1;
X[line].1 := X2;
X[line].2 := X3;
X[line].3 := X4;
X[line].4 := X5;

(* compare for change *)
FOR i := 0 TO 3 DO
	IF X[i] <> L[i] THEN
		(* scan line information has changed code need to be found and generated *)
		temp := x[i] XOR L[i];
		IF temp.0 THEN
			CODE := 1;
			code.7 := X[i].0;
			L[i].0 := X[i].0;
		ELSIF temp.1 THEN
			code := 2;
			code.7 := X[i].1;
			L[i].1 := X[i].1;
		ELSIF temp.2 THEN
			code := 3;
			code.7 := X[i].2;
			L[i].2 := X[i].2;
		ELSIF temp.3 THEN
			CODE := 4;
			CODE.7 := X[i].3;
			L[i].3 := X[i].3;
		ELSIF temp.4 THEN
			CODE := 5;
			CODE.7 := X[i].4;
			L[i].4 := X[i].4;
		END_IF;
		TP := TRUE;
		CODE.4 := Line.0;
		CODE.5 := Line.1;
		CODE.6 := Line.2;

		(* check if release codes need to be killed *)
		IF NOT release AND CODE < 127 THEN
			CODE := 0;
			TP := FALSE;
		END_IF;
		EXIT;
	END_IF;
END_FOR;

(* increment scan line every cycle *)
line := (line + 1) AND 2#0000_0011;
temp := SHL(BYTE#1,line);
Y1 := temp.0;
Y2 := temp.1;
Y3 := temp.2;
Y4 := temp.3;




(* revision history
hm		10.6.2007	rev 1.0		
	original version 


hm		11.9.2007	rev 1.1		
	deleted unused variables k and old_code

hm		23.12.2007	rev 1.2
	added exit statement in for loop instead of i:=5

hm		26. oct. 2008	rev 1.3
	code optimized

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/Others' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK PIN_CODE
VAR_INPUT
	CB : BYTE;
	E : BOOL;
END_VAR
VAR_INPUT CONSTANT
	PIN : STRING(8);
END_VAR
VAR_OUTPUT
	TP : BOOL;
END_VAR
VAR
	POS: INT := 1;
END_VAR

(*
version 1.0	30. oct. 2008
programmer 	hugo
tested by	oscat

MATRIX_CODE scans the input of a key_pad (MATRIX) for a sequence of characters.

*)
(* @END_DECLARATION := '0' *)
tp := FALSE;
IF e THEN
	IF CB = CODE(pin, pos) THEN
		pos := pos + 1;
		IF pos > LEN(PIN) THEN
			(* proper code is detected *)
			TP := TRUE;
			pos := 1;
		END_IF;
	ELSE
		pos := 1;
	END_IF;
END_IF;


(* revision history
hm	30. oct. 2008	rev 1.0		
	original version 


*)

END_FUNCTION_BLOCK

