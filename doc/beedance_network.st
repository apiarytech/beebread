// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK BASE64_DECODE_STR
// VAR_INPUT
// 	RUN : BOOL;
// END_VAR
// VAR_OUTPUT
// 	DONE : BOOL;
// END_VAR
// VAR_IN_OUT
// 	STR1 : STRING(192);
// 	STR2 : STRING(144);
// END_VAR
// VAR
// 	BASE64_DECODE_STREAM : BASE64_DECODE_STREAM;
// 	run_last : BOOL;
// 	buf1 : ARRAY [0..63] OF BYTE;
// 	buf2 : ARRAY [0..47] OF BYTE;
// 	mode : INT;
// 	size1 : INT;
// 	size2 : INT;
// 	end : INT;
// 	pos : INT;
// END_VAR
//
// CASE mode OF
// 0:	IF RUN AND NOT run_last THEN
// 		DONE := FALSE;
// 		mode := 1;
// 		end := LEN(STR1);
// 		pos := 0;
// 		STR2 := '';
// 	END_IF;
// 1:
// 	size1 := MIN(48,end - pos);
// 	IF size1 > 0 THEN
// 		_STRING_TO_BUFFER(STR:=MID(STR1, size1, pos + 1), POS:=0, PT:=REF(buf1), SIZE:=SIZEOF(buf1));
// 		BASE64_DECODE_STREAM(BUF1:=buf1 ,BUF2:=buf2 ,SIZE1:=size1 ,SIZE2=>size2);
// 		STR2 := CONCAT(STR2,BUFFER_TO_STRING(PT:=REF(buf2), SIZE:=SIZEOF(buf2), START:=0, STOP:=size2 - 1));
// 	ELSE
// 		mode := 0;
// 		DONE := TRUE;
// 	END_IF;
// 	pos := pos + size1;
// END_CASE;
//
// run_last := RUN;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK BASE64_DECODE_STREAM
VAR_IN_OUT
	BUF1 : ARRAY [0..63] OF BYTE;
	BUF2 : ARRAY [0..47] OF BYTE;
END_VAR
VAR_INPUT
	SIZE1 : INT;
END_VAR
VAR_OUTPUT
	SIZE2 : INT;
END_VAR
VAR_TEMP
	a: INT;
	b: WORD;
	o: INT;
	off: INT;
	i1: INT;
	i2: INT;
	c: INT;
END_VAR

IF SIZE1 <= 0 THEN SIZE2 := 0; RETURN; END_IF;
c := SIZE1 - 1;
i2 := 0;
a := 0;
b := 0;
FOR i1 := 0 TO c DO
	o := BUF1[i1];
	IF o = 61 THEN (* = *)	EXIT;
	ELSIF o > 96 THEN (* a-z *) off := 71;
	ELSIF o > 64 THEN (* A-Z *) off := 65;
	ELSIF o > 47 THEN (* 0-9 *) off := -4;
	ELSIF o = 43 THEN (* + *) off := -19;
	ELSIF o = 47 THEN (* / *) off := -16;
	END_IF;
	b := SHL(b,6) OR INT_TO_BYTE(o - off);
	a := a + 6;
	IF a >= 8 THEN
		a := a - 8;
		BUF2[i2] := WORD_TO_BYTE(SHR(b,a));
		i2 := i2 + 1;
		b := b AND NOT SHL(WORD#16#FFFF,a);
	END_IF;
END_FOR;
SIZE2 := i2;
END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK BASE64_ENCODE_STR
// VAR_INPUT
// 	RUN : BOOL;
// END_VAR
// VAR_OUTPUT
// 	DONE : BOOL;
// END_VAR
// VAR_IN_OUT
// 	STR1 : STRING(144);
// 	STR2 : STRING(192);
// END_VAR
// VAR
// 	BASE64_ENCODE_STREAM : BASE64_ENCODE_STREAM;
// 	run_last : BOOL;
// 	buf1 : ARRAY [0..47] OF BYTE;
// 	buf2 : ARRAY [0..63] OF BYTE;
// 	mode : INT;
// 	size1 : INT;
// 	size2 : INT;
// 	end : INT;
// 	pos : INT;
// END_VAR
//
// CASE mode OF
// 0:	IF RUN AND NOT run_last THEN
// 		DONE := FALSE;
// 		mode := 1;
// 		end := LEN(STR1);
// 		pos := 0;
// 		STR2 := '';
// 	END_IF;
// 1:
// 	size1 := MIN(48,end - pos);
// 	IF size1 > 0 THEN
// 		_STRING_TO_BUFFER(STR:=MID(STR1, size1, pos + 1), POS:=0, PT:=REF(buf1), SIZE:=SIZEOF(buf1));
// 		BASE64_ENCODE_STREAM(BUF1:=buf1 ,BUF2:=buf2 ,SIZE1:=size1 ,SIZE2=>size2);
// 		STR2 := CONCAT(STR2,BUFFER_TO_STRING(PT:=REF(buf2), SIZE:=SIZEOF(buf2), START:=0, STOP:=size2 - 1));
// 	ELSE
// 		mode := 0;
// 		DONE := TRUE;
// 	END_IF;
// 	pos := pos + size1;
// END_CASE;
//
// run_last := RUN;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK BASE64_ENCODE_STREAM
// VAR_IN_OUT
// 	BUF1 : ARRAY [0..47] OF BYTE;
// 	BUF2 : ARRAY [0..63] OF BYTE;
// END_VAR
// VAR_INPUT
// 	SIZE1 : INT;
// END_VAR
// VAR_OUTPUT
// 	SIZE2 : INT;
// END_VAR
// VAR
// 	a: INT;
// 	b: WORD;
// 	i: INT;
// 	i2: INT;
// 	i3: INT;
// 	c: INT;
// 	BASE64: STRING(64) := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
// END_VAR
//
// IF SIZE1 <= 0 THEN
// 	SIZE2 := 0;
// 	RETURN;
// END_IF;
// c := SIZE1 - 1;
// i2 := 0;
// a := 0;
// b := 0;
// FOR i := 0 TO c DO
// 	b := SHL(b,8) OR BUF1[i];
// 	a := a + 8;
// 	WHILE a >= 6 DO
// 		a := a - 6;
// 		buf2[i2] := CODE(BASE64,WORD_TO_INT(SHR(b,a)) + 1);
// 		i2 := i2 + 1;
// 		b := b AND NOT SHL(WORD#16#FFFF,a);
// 	END_WHILE;
// END_FOR;
// IF a > 0 THEN
// 	buf2[i2] := CODE(BASE64,WORD_TO_INT(SHL(b,6 - a)) + 1);
// 	i2 := i2 + 1;
// 	i3 := i2 + (NOT INT_TO_WORD(i2) AND 2#11);
// 	FOR i := i2 TO i3 DO
// 		buf2[i] := 61;
// 	END_FOR;
// 	SIZE2 := i3 + 1;
// ELSE
// 	SIZE2 := i2;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION_BLOCK CSV_PARSER_BUF
// VAR_INPUT
// 	SEP : BYTE;
// END_VAR
// VAR_IN_OUT
// 	RUN : BYTE;
// 	OFFSET : UDINT;
// 	VALUE : STRING(STRING_LENGTH);
// 	PT : NETWORK_BUFFER;
// END_VAR
// VAR_OUTPUT
// 	RESULT : BYTE;
// END_VAR
// VAR
// 	step : INT;
// 	size : INT;
// 	c : BYTE;
// 	v_add : BYTE;
// 	i : INT;
// 	state : BYTE;
// 	eof : BOOL;
// 	pv : REF_TO BYTE;
// 	pve : DWORD;
// 	watchdog : TON;
// 	x : BOOL;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	step_last : INT;*)
// END_VAR
//
// IF RUN > 0 THEN
// 	IF state = 0 THEN (* restart evaluation *)
//
// 		pv := REF(VALUE);
// 		pv^ := 0; (* VALUE = '' *)
// 		pve := pv + INT_TO_DWORD(STRING_LENGTH);
// 		i := UDINT_TO_INT(OFFSET);
// 		size := UINT_TO_INT(PT.SIZE) - 1;
// 		RESULT := 0;
// 		state := 5; (* command in progress (busy) *)
// 		step := 0; (* restart evaluation *)
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'CPF: Startindex: ~1';
// 			LOG_CL.PRINTF[1] := INT_TO_STRING(i);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	END_IF;
//
// 	watchdog(IN:=FALSE,PT:=T#1ms); (* restart timer *)
// 	watchdog.IN := TRUE;
//
// 	WHILE (state = 5) DO (* loop as long as there is no result *)
//
// 		watchdog(); (* update timer *)
//
// 		IF watchdog.Q THEN (* watchdog triggered *)
// 			state := 5;
// 			EXIT;
// 		ELSIF i > size THEN (* no further character available *)
// 			state := 10; (* nothing found - end of data reached *)
// 			EXIT;
// 		ELSE
//
// 			c := PT.BUFFER[i]; (* read new character from buffer *)
// 			eof := i + 1 = PT.SIZE; (* End of File *)
//
// 			v_add := 0;
//
// 			IF c = SEP THEN (* new element detected *)
// 				state := 1;
// 				x := TRUE;
// 			ELSIF c >= 32 THEN (* element detected *)
// 				v_add := c; (* store character *)
// 				IF eof THEN
// 					state := 2; (* element and end of line detected *)
// 				ELSE
// 					x := TRUE;
// 				END_IF;
// 			ELSE
// 				IF x THEN
// 					state := 2; (* element and end of line detected *)
// 					x := FALSE;
// 				END_IF;
// 			END_IF;
//
// 			IF v_add > 0  AND pv < pve THEN (* add character to value string *)
// 				pv^ := v_add; (* enter character *)
// 				pv := pv + 1;
// 				pv^ := 0; (* enter string termination *)
// 			END_IF;
//
// 			i := i + 1; (* increment data pointer *)
//
// 		END_IF;
// 	END_WHILE;
//
// 	CASE state OF
// 	01:	(* element found *)
// 		RESULT := 1; (* element found *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
// 	02:	(* element found *)
// 		RESULT := 2; (* element and end of line detected *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
//
// 	05:	RESULT := 5; (* in progress *)
//
// 	10:	(* index out of range - end of buffer *)
// 		RESULT := 10; (* nothing found - end of data reached *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
// 	END_CASE;
//
// 	(* ---------------- Debug-Message -----------------------*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'CPF: intern Result ~1 , state ~2, step ~3 , offset ~4, run ~5';
// 		LOG_CL.PRINTF[1] := BYTE_TO_STRING(RESULT);
// 		LOG_CL.PRINTF[2] := BYTE_TO_STRING(state);
// 		LOG_CL.PRINTF[3] := INT_TO_STRING(step);
// 		LOG_CL.PRINTF[4] := UDINT_TO_STRING(OFFSET);
// 		LOG_CL.PRINTF[5] := BYTE_TO_STRING(RUN);
// 		LOG_MSG();
// 	END_IF;*)
// 	(* -----------------------------------------------------------------*)
//
// END_IF;
//
// (* ---------------- Debug-Message -----------------------*)
// (*IF _debug_enable THEN
// 	IF step_last <> step THEN
// 		LOG_CL.NEW_MSG := 'IPB: Step: ~1';
// 		LOG_CL.PRINTF[1]  := INT_TO_STRING(step);
// 		LOG_MSG();
// 	END_IF;
// 	step_last := step;
// END_IF;*)
// (* -----------------------------------------------------------------*)
// END_FUNCTION_BLOCK

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION_BLOCK CSV_PARSER_FILE
// VAR_INPUT
// 	SEP : BYTE;
// END_VAR
// VAR_IN_OUT
// 	FILENAME : STRING;
// 	FSD : FILE_SERVER_DATA;
// 	RUN : BYTE;
// 	OFFSET : UDINT;
// 	VALUE : STRING(STRING_LENGTH);
// 	PT	: NETWORK_BUFFER;
// END_VAR
// VAR_OUTPUT
// 	RESULT : BYTE;
// END_VAR
// VAR
// 	i : UDINT;
// 	state : BYTE;
// 	e : BYTE;
// 	mode : BYTE;
// 	step : INT;
// 	c : BYTE;
// 	v_add : BYTE;
// 	eof : BOOL;
// 	pv : REF_TO BYTE;
// 	pve : DWORD;
// 	watchdog : TON;
// 	x : BOOL;
// 	FB : FILE_BLOCK;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := FALSE;
// 	step_last : INT;*)
// END_VAR
//
// IF RUN > 0 THEN
//
// 	IF state = 0 THEN (* restart evaluation *)
// 		pv := REF(VALUE);
// 		pv^ := 0; (* VALUE = '' *)
// 		pve := pv + INT_TO_DWORD(STRING_LENGTH);
// 		i := OFFSET;
// 		RESULT := 0;
// 		state := 5; (* command in progress (busy) *)
// 		step := 0; (* restart evaluation *)
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'CPF: Startindex: ~1';
// 			LOG_CL.PRINTF[1] := UDINT_TO_STRING(i);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	END_IF;
//
// 	watchdog(IN:=FALSE,PT:=T#1ms); (* restart timer *)
// 	watchdog.IN := TRUE;
//
// 	WHILE (state = 5) DO (* loop as long as there is no result *)
//
// 		mode := 1;
// 		FB(MODE:=mode,FSD:=fsd,PT:=PT,FILENAME:=FILENAME,POS:=i );
// 		e:=FB.ERROR;
// 		c:=FB.DATA;
//
// 		eof := i + 1 = FSD.FILE_SIZE; (* End of File *)
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'CPF: mode: ~1, error: ~2 , >~7<>~3< ~4 step ~5 index ~6';
// 			LOG_CL.PRINTF[1] := BYTE_TO_STRING(mode);
// 			LOG_CL.PRINTF[2] := BYTE_TO_STRING(e);
// 			LOG_CL.PRINTF[3] := BYTE_TO_STRING(c);
// 			LOG_CL.PRINTF[4] := BOOL_TO_STRING(eof);
// 			LOG_CL.PRINTF[5] := INT_TO_STRING(step);
// 			LOG_CL.PRINTF[6] := UDINT_TO_STRING(i);
// 			LOG_CL.PRINTF[7] := CHR_TO_STRING(c);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
// 		watchdog(); (* update timer *)
//
// 		IF mode > 0 OR watchdog.Q THEN (* block still active, or watchdog triggered *)
// 			state := 5;
// 			EXIT;
// 		ELSIF e > 0 THEN (* File_error or no further character available *)
// 			state := 10; (* nothing found - end of data reached *)
// 			EXIT;
// 		ELSE
//
// 			v_add := 0;
//
// 			IF c = SEP THEN (* new element detected *)
// 				state := 1;
// 				x := TRUE;
// 			ELSIF c >= 32 THEN (* element detected *)
// 				v_add := c; (* store character *)
// 				IF eof THEN
// 					state := 2; (* element and end of line detected *)
// 				ELSE
// 					x := TRUE;
// 				END_IF;
// 			ELSE
// 				IF x THEN
// 					state := 2; (* element and end of line detected *)
// 					x := FALSE;
// 				END_IF;
// 			END_IF;
//
// 			IF v_add > 0  AND pv < pve THEN (* add character to value string *)
// 				pv^ := v_add; (* enter character *)
// 				pv := pv + 1;
// 				pv^ := 0; (* enter string termination *)
// 			END_IF;
//
// 			i := i + 1; (* increment data pointer *)
//
// 		END_IF;
// 	END_WHILE;
//
// 	CASE state OF
// 	01:	(* element found *)
// 		RESULT := 1; (* element found *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
// 	02:	(* element found *)
// 		RESULT := 2; (* element and end of line detected *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
//
// 	05:	RESULT := 5; (* in progress *)
//
// 	10:	(* index out of range - end of buffer *)
// 		RESULT := 10; (* nothing found - end of data reached *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
// 	END_CASE;
//
// 	(* ---------------- Debug-Message -----------------------*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'CPF: intern Result ~1 , state ~2, step ~3 , offset ~4, run ~5';
// 		LOG_CL.PRINTF[1] := BYTE_TO_STRING(RESULT);
// 		LOG_CL.PRINTF[2] := BYTE_TO_STRING(state);
// 		LOG_CL.PRINTF[3] := INT_TO_STRING(step);
// 		LOG_CL.PRINTF[4] := UDINT_TO_STRING(OFFSET);
// 		LOG_CL.PRINTF[5] := BYTE_TO_STRING(RUN);
// 		LOG_MSG();
// 	END_IF;*)
// 	(* -----------------------------------------------------------------*)
//
// END_IF;
//
// (* ---------------- Debug-Message -----------------------*)
// (*IF _debug_enable THEN
// 	IF step_last <> step THEN
// 		LOG_CL.NEW_MSG := 'CPF: Step: ~1';
// 		LOG_CL.PRINTF[1]  := INT_TO_STRING(step);
// 		LOG_MSG();
// 	END_IF;
// 	step_last := step;
// END_IF;*)
// (* -----------------------------------------------------------------*)
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_BOOL
// VAR_IN_OUT
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	STATE : BOOL;
// 	OFF : STRING;
// 	ON : STRING;
// 	COLUMN : STRING(40);
// 	DELTA : BOOL;
// END_VAR
// VAR_OUTPUT
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	id : WORD;
// 	last_delta : BOOL;
// 	state_last: BOOL;
// END_VAR
//
// CASE x.ADD_COM OF
//
// 01:	(* ADD INFO *)
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	id := WORD#16#0601; (* source type BOOL, target type STRING *)
// 02:	(* ADD HEADER *)
// 	X.UCB.D_STRING := COLUMN;
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 03:	(* ADD DATA *)
// 	IF X.STORE_TYPE = BYTE#05 THEN (* MySQL mode *)
// 		X.UCB.D_STRING := SEL(state_last,'0','1');
// 	ELSE
// 		X.UCB.D_STRING := SEL(state_last,OFF,ON);
// 	END_IF;
// 	X.UCB.D_HEAD := id;
//
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 	last_delta := STATE;
// 04:	(* ADD DATA REQ *)
// 	IF DELTA AND (STATE XOR last_delta) THEN
// 		X.ADD_DATA_REQ := TRUE;
// 		last_delta := STATE;
// 	END_IF;
// END_CASE;
// state_last := STATE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses ELEMENT_GET, which is commented out.
// FUNCTION_BLOCK DLOG_CRON_TAB
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	DTI : DT;
// 	SECOND : STRING(20);
// 	MINUTE : STRING(20);
// 	HOUR : STRING(20);
// 	DAY_OF_MONTH : STRING(20);
// 	DAY_OF_WEEK : STRING(20);
// 	MONTH : STRING(20);
// END_VAR
// VAR_OUTPUT
// 	Q : BOOL;
// END_VAR
// VAR
// 	CD :  DLOG_CRON_DATA;
// 	activate_old : BOOL;
// 	datei : DATE;
// 	dti_last : DT;
// 	cnt1 : INT;
// 	cnt2 : INT;
// 	idx1 : INT;
// 	idx2 : INT;
// 	idx3 : INT;
// 	idx4 : INT;
// 	pos1 : INT;
// 	pos2 : INT;
// 	char_cnt : INT;
// 	trig_day : BOOL;
// 	trig_weekday : BOOL;
// 	str1 : STRING(STRING_LENGTH);
// 	str2 : STRING(STRING_LENGTH);
// 	str3 : STRING(STRING_LENGTH);
// 	num_str1 : STRING(10);
// END_VAR
//
// IF ACTIVATE THEN
// 	IF NOT activate_old THEN
// 		CD.CE[0].VALUE_MIN := 0;	(* seconds *)
// 		CD.CE[0].VALUE_MAX := 59;
// 		CD.CE[1].VALUE_MIN := 0;	(* minutes *)
// 		CD.CE[1].VALUE_MAX := 59;
// 		CD.CE[2].VALUE_MIN := 0;	(* hours *)
// 		CD.CE[2].VALUE_MAX := 59;
// 		CD.CE[3].VALUE_MIN := 1;	(* days *)
// 		CD.CE[3].VALUE_MAX := 31;
// 		CD.CE[4].VALUE_MIN := 1;	(* weekday *)
// 		CD.CE[4].VALUE_MAX := 7;
// 		CD.CE[5].VALUE_MIN := 1;	(* month *)
// 		CD.CE[5].VALUE_MAX := 12;
//
// 		FOR idx1 := 0 TO 5 DO
// 			idx3 := CD.CE[idx1].VALUE_MIN;
// 			idx4 := CD.CE[idx1].VALUE_MAx;
// 			FOR idx2 := idx3 TO idx4 DO
// 				CD.CE[idx1].ELEMENTS[idx2] := FALSE;
// 			END_FOR;
// 			CD.CE[idx1].ALL_SELECTED := FALSE;
// 		END_FOR;
//
// 		str1 := CONCAT(SECOND,'#');
// 		str1 := CONCAT(str1,MINUTE);
// 		str1 := CONCAT(str1,'#');
// 		str1 := CONCAT(str1,HOUR);
// 		str1 := CONCAT(str1,'#');
// 		str1 := CONCAT(str1,DAY_OF_MONTH);
// 		str1 := CONCAT(str1,'#');
// 		str1 := CONCAT(str1,DAY_OF_WEEK);
// 		str1 := CONCAT(str1,'#');
// 		str1 := CONCAT(str1,MONTH);
//
// 		FOR idx1 := 0 TO 5 DO
// 			str2 := ELEMENT_GET(BYTE#35, idx1, str1);
// 			IF LEN(str2) > 0 THEN
// 				cnt1 := ELEMENT_COUNT(BYTE#44, str2);	(* separator ',' *)
// 			ELSE
// 				cnt1 := 1;
// 			END_IF;
//
// 			FOR idx2 := 0 TO cnt1 -1 DO
// 				str3 := ELEMENT_GET(BYTE#44, idx2, str2);
// 				IF FIND(str3, '*/') > 0 AND LEN(str3) <= 10 THEN
// 					pos1 := DEC_TO_INT(str3);
// 					IF pos1 > 0 THEN
// 						idx3 := CD.CE[idx1].VALUE_MIN;
// 						idx4 := CD.CE[idx1].VALUE_MAX;
// 						WHILE idx3 <= idx4 DO
// 							CD.CE[idx1].ELEMENTS[idx3] := TRUE;
// 							idx3 := idx3 + pos1;
// 						END_WHILE;
// 					END_IF;
//
// 				ELSE
// 					IF FIND(str3, '*') > 0 OR LEN(str3) = 0 THEN
// 						pos1 := CD.CE[idx1].VALUE_MIN;
// 						pos2 := CD.CE[idx1].VALUE_MAX;
// 						CD.CE[idx1].ALL_SELECTED := TRUE;
// 					ELSE
// 						cnt2 := ELEMENT_COUNT(BYTE#45, str3);	(* separator '-' *)
// 						char_cnt := LEN(str3);
// 						IF char_cnt >= 1 AND char_cnt <= 10 THEN
// 							num_str1 := ELEMENT_GET(BYTE#45, 0, str3);
// 							pos1 := DEC_TO_INT(num_str1);
// 							pos2 := pos1;
// 							IF cnt2 = 2 THEN
// 								num_str1 := ELEMENT_GET(BYTE#45, 1, str3);
// 								pos2 := DEC_TO_INT(num_str1);
// 							END_IF;
// 						END_IF;
// 					END_IF;
//
// 					IF pos1 >= CD.CE[idx1].VALUE_MIN AND pos2 <= CD.CE[idx1].VALUE_MAX THEN
// 						FOR idx3 := pos1 TO pos2 DO
// 							CD.CE[idx1].ELEMENTS[idx3] := TRUE;
// 						END_FOR;
// 					END_IF;
// 				END_IF;
// 			END_FOR;
// 		END_FOR;
// 	ELSIF DTI <> dti_last THEN
// 		CD.CE[0].VALUE := SECOND_OF_DT(DTI); (* seconds *)
// 		CD.CE[1].VALUE := MINUTE_OF_DT(DTI); (* minutes *)
// 		CD.CE[2].VALUE := HOUR_OF_DT(DTI); (* hours *)
// 		datei := DT_TO_DATE(DTI);
// 		CD.CE[3].VALUE := DAY_OF_MONTH(datei); (* Tag *)
// 		CD.CE[4].VALUE := DAY_OF_WEEK(datei); (* weekday *)
// 		CD.CE[5].VALUE := MONTH_OF_DATE(datei); (*  month *)
// 		dti_last := DTI;
//
// 		idx2 := CD.CE[0].VALUE;	(* seconds *)
// 		Q := CD.CE[0].ELEMENTS[idx2] = TRUE;
// 		idx2 := CD.CE[1].VALUE;	(* minutes *)
// 		Q := Q AND CD.CE[1].ELEMENTS[idx2] = TRUE;
// 		idx2 := CD.CE[2].VALUE;	(* hours *)
// 		Q := Q AND CD.CE[2].ELEMENTS[idx2] = TRUE;
// 		idx2 := CD.CE[5].VALUE;	(* month *)
// 		Q := Q AND CD.CE[5].ELEMENTS[idx2] = TRUE;
//
// 		(* evaluate day and weekday *)
// 		idx2 := CD.CE[3].VALUE;	(* Tag *)
// 		trig_day:= CD.CE[3].ELEMENTS[idx2] = TRUE;
// 		idx2 := CD.CE[4].VALUE;	(* weekday *)
// 		trig_weekday := CD.CE[4].ELEMENTS[idx2] = TRUE;
//
// 		IF CD.CE[3].ALL_SELECTED = FALSE AND CD.CE[4].ALL_SELECTED = TRUE THEN
// 			Q := Q AND trig_day; (* evaluate day only *)
// 		ELSIF CD.CE[3].ALL_SELECTED = TRUE AND CD.CE[4].ALL_SELECTED = FALSE THEN
// 			Q := Q AND trig_weekday; (* evaluate weekday only *)
// 		ELSIF CD.CE[3].ALL_SELECTED = FALSE AND CD.CE[4].ALL_SELECTED = FALSE THEN
// 			Q := Q AND (trig_weekday OR trig_day); (* evaluate day and weekday *)
// 		END_IF;
// 	END_IF;
//
// END_IF;
// activate_old := ACTIVATE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_DINT
// VAR_IN_OUT
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	VALUE : DINT;
// 	COLUMN : STRING(40);
// 	DELTA : DINT;
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	id : WORD;
// 	delta_last : DINT;
// 	value_last : DINT;
// END_VAR
//
// CASE x.ADD_COM OF
//
// 01:	(* ADD INFO *)
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	id := WORD#16#0801; (* source type DINT, target type STRING *)
// 02:	(* ADD HEADER *)
// 	X.UCB.D_STRING := COLUMN;
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 03:	(* ADD DATA *)
// 	X.UCB.D_STRING := DINT_TO_STRING(value_last);
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 	delta_last := value_last;
// 04:	(* ADD DATA REQ *)
// 	IF DELTA <> DINT#0 THEN
// 		IF VALUE <= (delta_last - DELTA) OR VALUE >= (delta_last + DELTA) THEN
// 			X.ADD_DATA_REQ := TRUE;
// 			delta_last := VALUE;
// 		END_IF;
// 	END_IF;
// END_CASE;
// value_last := VALUE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_DT
// VAR_IN_OUT
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	FMT : STRING := '#A-#D-#H #N:#R:#T';
// 	COLUMN : STRING(40);
// 	DELTA : UDINT;
// END_VAR
//
// VAR_OUTPUT
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	id : WORD;
// 	delta_last : DT;
// END_VAR
//
// CASE x.ADD_COM OF
//
// 01:	(* ADD INFO *)
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	id := WORD#16#0D01; (* source type DT, target type STRING *)
// 02:	(* ADD HEADER *)
// 	X.UCB.D_STRING := COLUMN;
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 03:	(* ADD DATA *)
// 	X.UCB.D_STRING := DT_TO_STRF(DTI:=X.DTI,MS:=0,FMT:=FMT,LANG:=0);
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 	delta_last := X.DTI;
// 04:	(* ADD DATA REQ *)
// 	IF DELTA <> UDINT#0 THEN
// 		IF DT_TO_UDINT(X.DTI) >= (DT_TO_UDINT(delta_last) + DELTA) THEN
// 			X.ADD_DATA_REQ := TRUE;
// 			delta_last := X.DTI;
// 		END_IF;
// 	END_IF;
// END_CASE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_FILE_TO_FTP
// VAR_IN_OUT
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	FTP_URL : STRING(STRING_LENGTH);
// 	FTP_ACTIV : BOOL;
// 	FILE_DELETE : BOOL;
// 	TIMEOUT : TIME := T#30s;
// 	RETRY : INT;
// 	RETRY_TIME : TIME := T#30s;
// 	DNS_IP4 : DWORD;
// 	PLC_IP4 : DWORD;
// END_VAR
// VAR_OUTPUT
// 	DONE : BOOL;
// 	BUSY : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	UCBD : UNI_CIRCULAR_BUFFER_DATA;
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	FTP_CLIENT : FTP_CLIENT;
// 	WT_1 : TON;
// 	ftrig_old : BOOL;
// 	filename : STRING;
// 	step : INT;
// 	ftp_busy : BOOL;
// 	ftp_done : BOOL;
// 	cnt : INT;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;*)
// END_VAR
//
// (* store new file name *)
// IF X.NEW_FILE_RTRIG AND ftrig_old = FALSE THEN
// 	UCBD.D_HEAD := 1;
// 	UCBD.D_MODE := 1;
// 	UCBD.D_STRING := X.NEW_FILE;
// 	UCB(DATA:=UCBD); (* enter data *)
// END_IF;
// ftrig_old := X.NEW_FILE_RTRIG;
//
// CASE step OF
//
// 00:	IF UCBD.BUF_COUNT > 0 AND ftp_busy = FALSE THEN
// 		UCBD.D_MODE := 12; (* read and delete data *)
// 		UCB(DATA:=UCBD);
// 		filename := UCBD.D_STRING;
// 		cnt := RETRY;
// 		step := 10;
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'F_T_FTP: Auftrag aktiv /file: ~2';
// 			LOG_CL.PRINTF[2] := files;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
// 	END_IF;
//
// 10:	IF ftp_busy = FALSE THEN
// 		IF ftp_done THEN
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_FTP: Auftrag OK /file: ~2';
// 				LOG_CL.PRINTF[2] := filename;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 0;
// 		ELSE
// 			step := 20; (* pause after error *)
// 		END_IF;
// 	END_IF;
//
// 20:	IF WT_1.Q THEN (* wait time *)
// 		cnt := cnt - 1;
// 		IF RETRY = 0 THEN (* repeat FTP transfer indefinitely *)
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_FTP: Auftrag unbegrenzt wiederholen /file: ~2';
// 				LOG_CL.PRINTF[2] := filename;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 10; (* repeat FTP transfer *)
// 		ELSIF cnt >= 0 THEN
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_FTP: Auftrag wiederholen /cnt: ~1 /file: ~2';
// 				LOG_CL.PRINTF[1] := INT_TO_STRING(cnt);
// 				LOG_CL.PRINTF[2] := filename;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 10; (* repeat FTP transfer *)
// 		ELSE
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_FTP: Auftrag verworfen /file: ~2';
// 				LOG_CL.PRINTF[2] := filename;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 0; (* discard FTP job *)
// 		END_IF;
// 	END_IF;
//
// END_CASE;
//
// BUSY := step = 10;
// DONE := step = 0;
// WT_1(IN:= step = 20, PT:= RETRY_TIME);
//
// FTP_CLIENT(
// 	ACTIVATE:= BUSY,
// 	FILENAME:= filename,
// 	FTP_URL:= FTP_URL,
// 	FTP_DOWNLOAD:= FALSE,
// 	FTP_ACTIV:= FTP_ACTIV,
// 	FILE_DELETE:= FILE_DELETE,
// 	TIMEOUT:= TIMEOUT,
// 	DNS_IP4:= DNS_IP4,
// 	PLC_IP4:= PLC_IP4,
// 	DONE=> ftp_done,
// 	BUSY=> ftp_busy
// 	);
//
// (* error evaluation *)
// IF FTP_CLIENT.ERROR_T <> BYTE#0 THEN
// 	ERROR_C := FTP_CLIENT.ERROR_C;
// 	ERROR_T := FTP_CLIENT.ERROR_T;
// ELSIF UCBD.BUF_DATA_LOST > UDINT#0 THEN
// 	ERROR_C := DWORD#2; (* error: buffer overflow - reduce amount of data *)
// 	ERROR_T := BYTE#6; 	(* Type: Buffer *)
// ELSIF UCBD.BUF_USED_MAX > USINT#90 THEN
// 	ERROR_C := DWORD#1; (* warning: buffer utilization > 90% *)
// 	ERROR_T := BYTE#6;	(* Type: Buffer *)
// ELSE
// 	ERROR_C := DWORD#0;
// 	ERROR_T := BYTE#0;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_FILE_TO_SMTP
// VAR_IN_OUT
// 	SERVER : STRING;
// 	MAILFROM : STRING;
// 	MAILTO : STRING(STRING_LENGTH);
// 	SUBJECT : STRING;
// 	BODY : STRING(STRING_LENGTH);
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	FILE_DELETE : BOOL;
// 	TIMEOUT : TIME := T#30s;
// 	DTI : DT;
// 	DTI_OFFSET : INT;
// 	RETRY : INT;
// 	RETRY_TIME : TIME := T#30s;
// 	DNS_IP4 : DWORD;
// END_VAR
// VAR_OUTPUT
// 	DONE : BOOL;
// 	BUSY : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	UCBD : UNI_CIRCULAR_BUFFER_DATA;
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	SMTP_CLIENT : SMTP_CLIENT;
// 	WT_1 : TON;
// 	ftrig_old : BOOL;
// 	files : STRING(STRING_LENGTH);
// 	step : INT;
// 	smtp_busy : BOOL;
// 	smtp_done : BOOL;
// 	cnt : INT;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;*)
// END_VAR
//
// (* store new file name *)
// IF X.NEW_FILE_RTRIG AND ftrig_old = FALSE THEN
// 	UCBD.D_HEAD := 1;
// 	UCBD.D_MODE := 1;
// 	UCBD.D_STRING := X.NEW_FILE;
// 	UCB(DATA:=UCBD); (* enter data *)
// END_IF;
// ftrig_old := X.NEW_FILE_RTRIG;
//
// CASE step OF
//
// 00:	IF UCBD.BUF_COUNT > 0 AND smtp_busy = FALSE THEN
// 		UCBD.D_MODE := 12; (* read and delete data *)
// 		UCB(DATA:=UCBD);
// 		files := UCBD.D_STRING;
// 		IF FILE_DELETE THEN
// 			files := CONCAT(files,';#DEL#');
// 		END_IF;
// 		cnt := RETRY;
// 		step := 10;
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'F_T_SMTP: Auftrag aktiv /file: ~2';
// 			LOG_CL.PRINTF[2] := files;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
// 	END_IF;
//
// 10:	IF smtp_busy = FALSE THEN
// 		IF smtp_done THEN
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_SMTP: Auftrag OK /file: ~2';
// 				LOG_CL.PRINTF[2] := files;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 0;
// 		ELSE
// 			step := 20; (* pause after error *)
// 		END_IF;
// 	END_IF;
//
// 20:	IF WT_1.Q THEN (* wait time *)
// 		cnt := cnt - 1;
// 		IF RETRY = 0 THEN (* repeat transfer indefinitely *)
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_SMTP: Auftrag unbegrenzt wiederholen /file: ~2';
// 				LOG_CL.PRINTF[2] := files;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 10; (* repeat transfer *)
// 		ELSIF cnt >= 0 THEN
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_SMTP: Auftrag wiederholen /cnt: ~1 /file: ~2';
// 				LOG_CL.PRINTF[1] := INT_TO_STRING(cnt);
// 				LOG_CL.PRINTF[2] := files;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 10; (* repeat transfer *)
// 		ELSE
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'F_T_SMTP: Auftrag verworfen /file: ~2';
// 				LOG_CL.PRINTF[2] := files;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 			step := 0; (* discard job *)
// 		END_IF;
// 	END_IF;
//
// END_CASE;
//
// BUSY := step = 10;
// DONE := step = 0;
// WT_1(IN:= step = 20, PT:= RETRY_TIME);
//
// SMTP_CLIENT(
// 	ACTIVATE:= BUSY,
// 	TIMEOUT:= TIMEOUT,
// 	DTI:= DTI,
// 	DTI_OFFSET:= DTI_OFFSET,
// 	DNS_IP4:= DNS_IP4,
// 	SERVER:= SERVER,
// 	MAILFROM:= MAILFROM,
// 	MAILTO:= MAILTO,
// 	SUBJECT:= SUBJECT,
// 	BODY:= BODY,
// 	FILES:= FILES,
// 	DONE=> smtp_done,
// 	BUSY=> smtp_busy
// 	);
//
// (* error evaluation *)
// IF SMTP_CLIENT.ERROR_T <> BYTE#0 THEN
// 	ERROR_C := SMTP_CLIENT.ERROR_C;
// 	ERROR_T := SMTP_CLIENT.ERROR_T;
// ELSIF UCBD.BUF_DATA_LOST > UDINT#0 THEN
// 	ERROR_C := DWORD#2; (* error: buffer overflow - reduce amount of data *)
// 	ERROR_T := BYTE#6; 	(* Type: Buffer *)
// ELSIF UCBD.BUF_USED_MAX > USINT#90 THEN
// 	ERROR_C := DWORD#1; (* warning: buffer utilization > 90% *)
// 	ERROR_T := BYTE#6;	(* Type: Buffer *)
// ELSE
// 	ERROR_C := DWORD#0;
// 	ERROR_T := BYTE#0;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_REAL
// VAR_IN_OUT
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	VALUE : REAL;
// 	N : INT;
// 	D : STRING(1) := ',';
// 	COLUMN : STRING(40);
// 	DELTA : REAL;
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	id : WORD;
// 	delta_last : REAL;
// 	value_last : REAL;
// END_VAR
//
// CASE x.ADD_COM OF
//
// 01:	(* ADD INFO *)
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	id := WORD#16#0201; (* source type REAL, target type STRING *)
// 02:	(* ADD HEADER *)
// 	X.UCB.D_STRING := COLUMN;
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 03:	(* ADD DATA *)
// 	X.UCB.D_STRING := REAL_TO_STRF(IN:=value_last,N:=N,D:=D);
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 	delta_last := value_last;
// 04:	(* ADD DATA REQ *)
// 	IF DELTA <> 0.0 THEN
// 		IF VALUE <= (delta_last - DELTA) OR VALUE >= (delta_last + DELTA) THEN
// 			X.ADD_DATA_REQ := TRUE;
// 			delta_last := VALUE;
// 		END_IF;
// 	END_IF;
// END_CASE;
// value_last := VALUE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_REAL_ARRAY
// VAR_IN_OUT
// 	R_ARRAY : DLOG_REAL_ARRAY_DATA;
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	N : INT;
// 	D : STRING(1);
// 	R_COUNT : INT;
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	id : WORD;
// 	idx : INT;
// END_VAR
//
// FOR idx := 1 TO R_COUNT DO
//
// 	CASE x.ADD_COM OF
//
// 	01:	(* ADD INFO *)
// 		X.ID_MAX := X.ID_MAX + USINT#1;
// 		id := WORD#16#0201; (* source type REAL, target type STRING *)
// 		(* assign name automatically *)
// 		IF LEN(R_ARRAY[idx].COLUMN) = 0 THEN
// 			R_ARRAY[idx].COLUMN := CONCAT('Spalte_',USINT_TO_STRING(X.ID_MAX));
// 		END_IF;
// 	02:	(* ADD HEADER *)
// 		X.UCB.D_STRING := R_ARRAY[idx].COLUMN;
// 		X.UCB.D_HEAD := id;
// 		X.UCB.D_MODE := 1;
// 		UCB(DATA:=X.UCB); (* enter data *)
// 	03:	(* ADD DATA *)
// 		X.UCB.D_STRING := REAL_TO_STRF(IN:=R_ARRAY[idx]._value_last,N:=N,D:=D);
// 		X.UCB.D_HEAD := id;
// 		X.UCB.D_MODE := 1;
// 		UCB(DATA:=X.UCB); (* enter data *)
// 		R_ARRAY[idx]._delta_last := R_ARRAY[idx]._value_last;
// 	04:	(* ADD DATA REQ *)
// 		IF R_ARRAY[idx].DELTA <> 0.0 THEN
// 			IF R_ARRAY[idx].VALUE <= (R_ARRAY[idx]._delta_last - R_ARRAY[idx].DELTA) OR R_ARRAY[idx].VALUE >= (R_ARRAY[idx]._delta_last + R_ARRAY[idx].DELTA) THEN
// 				X.ADD_DATA_REQ := TRUE;
// 				R_ARRAY[idx]._delta_last := R_ARRAY[idx].VALUE;
// 			END_IF;
// 		END_IF;
// 	END_CASE;
// 	R_ARRAY[idx]._value_last := R_ARRAY[idx].VALUE;
//
// END_FOR;
// END_FUNCTION_BLOCK

// Commented out for beedance: it names a variable fn, which beedance reads as its function keyword.
// FUNCTION_BLOCK DLOG_STORE_FILE_CSV
// VAR_IN_OUT
// 	X : DLOG_DATA;
// 	SAVE_DATA : DLOG_SAVE;
// END_VAR
// VAR_INPUT
// 	ENABLE : BOOL;
// 	TRIG_M : BOOL;
// 	TRIG_T : TIME;
// 	FILENAME : STRING;
// 	DTI : DT;
// 	SEP : BYTE;
// 	AUTO_CLOSE : TIME := T#15s;
// END_VAR
// VAR_OUTPUT
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	FS : FILE_SERVER;
// 	FSD : FILE_SERVER_DATA;
// 	PT : NETWORK_BUFFER;
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	trig_auto : BOOL;
// 	trig_stored : BOOL;
// 	trig_tmp : BOOL;
// 	trig_m_last : BOOL;
// 	dti_last : DT := DT#2070-02-06-06:28:15;
// 	fn_last : STRING;
// 	fn : STRING;
// 	step_1: INT;
// 	step_2: INT;
// 	idx : INT;
// 	log_stop : BOOL;
// 	n : USINT;
// 	wd_ton : TON;
// 	wd_time : TIME;
// 	aw_ton : TON;
// 	aw_enable : BOOL;
// 	aw_aktiv : BOOL;
// 	total_bytes : UDINT;
//
// 	(*
// 	LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	*)
// END_VAR
//
// X.DTI := DTI; (* enter date / time *)
// X.CLOCK_TRIG := DTI <> dti_last;
// dti_last := DTI;
// trig_auto := FALSE;
// IF X.CLOCK_TRIG THEN (* pulse every second *)
// 	fn := DT_TO_STRF(DTI, 0, FILENAME, 0); (* automatic file name generation *)
// 	IF TRIG_T >= T#1000ms THEN (* automatic trigger *)
// 		trig_auto := DT_TO_UDINT(DTI) MOD (TIME_TO_UDINT(TRIG_T) / UDINT#1000) = UDINT#0;
// 	END_IF;
// END_IF;
//
// CASE step_1 OF
//
// 00:	IF ENABLE THEN
// 		aw_enable := AUTO_CLOSE >= T#5s;
// 		wd_time := SEL(X.LOAD_TIME_MAX > T#0s, T#5ms,X.LOAD_TIME_MAX);
// 		fn_last := fn; (* store current file name *)
// 		IF X.ID_MAX = USINT#0 THEN (* scan for satellite blocks only once *)
// 			X.ADD_COM := 01; (* ADD INFO *)
// 			X.STORE_TYPE := BYTE#2; (* CSV mode *)
// 		END_IF;
// 		step_1 := 10;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_CSV1 - Start: File [~1], [~2]';
// 			LOG_CL.PRINTF[1] := fn_last;
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	END_IF;
//
// 10:	X.UCB.D_STRING := fn_last;
// 	IF fn_last = SAVE_DATA.FN_REM THEN
// 		X.UCB.D_HEAD := WORD#16#F201; (* Filemode open + write *)
// 		X.ADD_COM := 00;
// 	ELSE
// 		X.UCB.D_HEAD := WORD#16#F301; (* Filemode create + write *)
// 		X.ADD_COM := 02; (* ADD HEADER *)
// 	END_IF;
//
// 	X.UCB.D_MODE := 1; (* store file name as string *)
// 	UCB(DATA:=X.UCB); (* enter data *)
// 	step_1 := 30;
//
// 30:	(* process LOG trigger *)
// 	trig_tmp := (TRIG_M AND NOT trig_m_last) OR trig_auto OR X.ADD_DATA_REQ;
//
// 	IF ENABLE = FALSE OR fn <> fn_last THEN (* ENABLE = 0 or new file name -> log stop *)
// 		trig_stored := trig_tmp; (* store current trigger pulse *)
// 		X.UCB.D_HEAD := WORD#16#EE00;
// 		X.UCB.D_MODE := 1;
// 		UCB(DATA:=X.UCB); (* write element *)
// 		X.ADD_COM := 00;
// 		step_1 := 00;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_CSV1 - Stop: [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	ELSIF ENABLE THEN
// 		IF trig_tmp OR trig_stored THEN (* manual, automatic or external trigger *)
// 			trig_stored := FALSE;
// 			X.ADD_COM := 03; (* ADD DATA *)
// 			SAVE_DATA.TRIG_CNT := SAVE_DATA.TRIG_CNT + UDINT#1;
// 			SAVE_DATA.TRIG_CNT_TOTAL  := SAVE_DATA.TRIG_CNT_TOTAL + UDINT#1;
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'FILE_CSV1 - Trig: [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		ELSE
// 			X.ADD_COM := 04;	(* ADD DATA REQ *)
// 		END_IF;
// 	END_IF;
//
// END_CASE;
// X.ADD_DATA_REQ := FALSE; (* reset external request *)
//
// (* ##################################### *)
//
// CASE step_2 OF
//
// 00:	aw_ton(IN:= idx > 0 AND aw_enable, PT:= AUTO_CLOSE); (* AutoWrite Timer *)
// 	IF aw_ton.Q THEN (* time-controlled data writing *)
// 		aw_aktiv := TRUE;
// 		step_2 := 10; (* write data *)
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_CSV2 - AutoWrite: [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	ELSE
// 		wd_ton(IN:= FALSE, PT:= wd_time);
// 		wd_ton(IN:= TRUE);
// 		WHILE X.UCB.BUF_COUNT > 0 AND step_2 = 0 AND wd_ton.Q = FALSE DO
//
// 			X.UCB.D_MODE := 10;
// 			UCB(DATA:=X.UCB); (* read element, but do not remove it yet *)
// 			IF X.UCB.D_HEAD > WORD#16#F000 THEN (* F301 or F201 *)
// 				X.UCB.D_MODE := 11;
// 				UCB(DATA:=X.UCB); (* remove element *)
//
// 				FSD.MODE := WORD_TO_BYTE(SHR(SHL(X.UCB.D_HEAD,4),12)); (* 0X00 -> 000X = 3 create + write / 2 open + write *)
// 				FSD.OFFSET := UDINT#0; (* start at beginning of file *)
// 				FSD.FILENAME := X.UCB.D_STRING; (* file name *)
// 				PT.SIZE := UINT#0;
// 				SAVE_DATA.FN_REM := FSD.FILENAME; (* save current file name as retained *)
//
// 				IF FSD.MODE = BYTE#3 THEN	(* create new file -> then reset counter *)
// 					SAVE_DATA.TRIG_CNT := UDINT#0;
// 				END_IF;
//
// 				X.NEW_FILE := FSD.FILENAME;
// 				X.NEW_FILE_SIZE := UDINT#0;
// 				X.NEW_FILE_RTRIG := FALSE;
// 				step_2 := 20;
// 				(* ---------------- Debug-Message -----------------------*)
// 				(*IF _debug_enable THEN
// 					LOG_CL.NEW_MSG := 'FILE_CSV2 - Create: Mode [~1] , [~2]';
// 					LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 					LOG_CL.PRINTF[1] := BYTE_TO_STRING(FSD.MODE);
// 					LOG_MSG();
// 				END_IF;*)
// 				(* -----------------------------------------------------------------*)
//
// 			ELSIF X.UCB.D_HEAD = WORD#16#EE00 THEN (* LOG_STOP Command *)
// 				X.UCB.D_MODE := 11;
// 				UCB(DATA:=X.UCB); (* remove element *)
// 				log_stop := TRUE;
// 				step_2 := 10; (* write data *)
//
// 			ELSIF idx + X.UCB.D_SIZE + 2 < SIZEOF(PT.BUFFER) THEN (* buffer space available *)
// 				X.UCB.D_MODE := 11;
// 				UCB(DATA:=X.UCB); (* remove element *)
// 				(* ---------------- Debug-Message -----------------------*)
// 				(*IF _debug_enable THEN
// 					LOG_CL.NEW_MSG := 'FILE_CSV2 - BUF: TXT [~1],Pos [~2], [~3]';
// 					LOG_CL.PRINTF[1] := X.UCB.D_STRING;
// 					LOG_CL.PRINTF[2] := INT_TO_STRING(idx);
// 					LOG_CL.PRINTF[3] := DT_TO_STRING(DTI);
// 					LOG_MSG();
// 				END_IF;*)
// 				(* -----------------------------------------------------------------*)
// 				_STRING_TO_BUFFER(STR:=X.UCB.D_STRING, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 				idx := idx + LEN(X.UCB.D_STRING);
// 				n := n + USINT#1;
// 				IF n = X.ID_MAX THEN (* last element of a line *)
// 					n := USINT#0;
// 					PT.BUFFER[idx] := BYTE#16#0D;
// 					idx := idx + 1;
// 					PT.BUFFER[idx] := BYTE#16#0A;
// 				ELSE
// 					PT.BUFFER[idx] := SEP;
// 				END_IF;
// 				idx := idx + 1;
// 				X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
// 			ELSE (* no space available *)
// 				step_2 := 10; (* write data *)
// 			END_IF;
// 			wd_ton();
// 		END_WHILE;
// 	END_IF;
//
// 10:	IF idx > 0 THEN (* write data *)
// 		FSD.MODE := BYTE#2; (* File_Open_Exist + File_Write *)
// 		FSD.OFFSET := UDINT#4294967294; (* 16#FFFF_FFFE -> special code: automatically append data to end of file *)
// 		PT.SIZE := INT_TO_UINT(idx); (* enter data length *)
// 		idx := 0;
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_CSV2 - WriteData: Mode [~3] Size [~1] , [~2]';
// 			LOG_CL.PRINTF[1] := UINT_TO_STRING(PT.SIZE);
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_CL.PRINTF[3] := BYTE_TO_STRING(FSD.MODE);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
// 	END_IF;
// 	step_2 := 20;
//
// 20:	IF FSD.MODE = BYTE#0 THEN
// 		total_bytes := FSD.FILE_SIZE;
// 		X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
// 		IF log_stop OR aw_aktiv THEN
// 			aw_aktiv := FALSE;
// 			FSD.MODE := BYTE#5; (* close file *)
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'FILE_CSV2 - Close: [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		END_IF;
// 		step_2 := 30;
// 	END_IF;
//
// 30:	IF FSD.MODE = BYTE#0 THEN (* command finished? *)
// 		IF log_stop THEN
// 			X.NEW_FILE_RTRIG := TRUE;
// 			log_stop := FALSE;
// 			SAVE_DATA.FN_REM := ''; (* clear retained file name *)
// 		END_IF;
// 		step_2 := 00;
// 	END_IF;
//
// END_CASE;
//
// FS(FSD:=FSD,PT:=PT); (* FILE_SERVER *)
// (* error evaluation *)
// IF FSD.ERROR <> BYTE#0 THEN
// 	ERROR_C := BYTE_TO_DWORD(FSD.ERROR);
// 	ERROR_T := BYTE#1;
// ELSIF X.UCB.BUF_DATA_LOST > UDINT#0 THEN
// 	ERROR_C := DWORD#2; (* error: buffer overflow - reduce amount of data *)
// 	ERROR_T := BYTE#6; 	(* Type: Buffer *)
// ELSIF X.UCB.BUF_USED_MAX > USINT#90 THEN
// 	ERROR_C := DWORD#1; (* warning: buffer utilization > 90% *)
// 	ERROR_T := BYTE#6;	(* Type: Buffer *)
// ELSE
// 	ERROR_C := DWORD#0;
// 	ERROR_T := BYTE#0;
// END_IF;
//
// trig_m_last := TRIG_M;
// END_FUNCTION_BLOCK

// Commented out for beedance: it names a variable fn, which beedance reads as its function keyword.
// FUNCTION_BLOCK DLOG_STORE_FILE_HTML
// VAR_IN_OUT
// 	X : DLOG_DATA;
// 	SAVE_DATA : DLOG_SAVE;
// END_VAR
// VAR_INPUT
// 	ENABLE : BOOL;
// 	TRIG_M : BOOL;
// 	TRIG_T : TIME;
// 	FILENAME : STRING;
// 	DTI : DT;
// 	AUTO_CLOSE : TIME := T#15s;
// 	HTML_CAPTION : STRING;
// 	HTML_TABLE : STRING;
// 	HTML_TR_HEAD : STRING;
// 	HTML_TR_EVEN : STRING;
// 	HTML_TR_ODD : STRING;
// END_VAR
// VAR_OUTPUT
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	FS : FILE_SERVER;
// 	FSD : FILE_SERVER_DATA;
// 	PT : NETWORK_BUFFER;
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	trig_auto : BOOL;
// 	trig_stored : BOOL;
// 	trig_tmp : BOOL;
// 	trig_m_last : BOOL;
// 	dti_last : DT := DT#2070-02-06-06:28:15;
// 	fn_last : STRING;
// 	str : STRING(STRING_LENGTH);
// 	fn : STRING;
// 	step_1 : INT;
// 	step_2 : INT;
// 	idx : INT;
// 	log_state : USINT;
// 	log_stop : BOOL;
// 	n : INT;
// 	wd_ton : TON;
// 	wd_time : TIME;
// 	aw_ton : TON;
// 	aw_enable : BOOL;
// 	aw_aktiv : BOOL;
// 	total_bytes : UDINT;
// 	cur_ode : BYTE;
//
// 	(*
// 	LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	*)
//
// END_VAR
//
// X.DTI := DTI; (* enter date / time *)
// X.CLOCK_TRIG := DTI <> dti_last;
// dti_last := DTI;
// trig_auto := FALSE;
// IF X.CLOCK_TRIG THEN (* pulse every second *)
// 	fn := DT_TO_STRF(DTI, 0, FILENAME, 0); (* automatic file name generation *)
// 	IF TRIG_T >= T#1000ms THEN (* automatic trigger *)
// 		trig_auto := DT_TO_UDINT(DTI) MOD (TIME_TO_UDINT(TRIG_T) / UDINT#1000) = UDINT#0;
// 	END_IF;
// END_IF;
//
// CASE step_1 OF
//
// 00:	IF ENABLE THEN
// 		aw_enable := AUTO_CLOSE >= T#5s;
// 		wd_time := SEL(X.LOAD_TIME_MAX > T#0s, T#5ms,X.LOAD_TIME_MAX);
// 		fn_last := fn; (* store current file name *)
// 		IF X.ID_MAX = USINT#0 THEN (* scan for satellite blocks only once *)
// 			X.ADD_COM := 01; (* ADD INFO *)
// 			X.STORE_TYPE := BYTE#4; (* HTML mode *)
// 		END_IF;
// 		step_1 := 10;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_HTML1 - Start: File [~1], [~2]';
// 			LOG_CL.PRINTF[1] := fn_last;
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	END_IF;
//
// 10:	X.UCB.D_STRING := fn_last;
// 	IF fn_last = SAVE_DATA.FN_REM THEN
// 		X.UCB.D_HEAD := WORD#16#F201; (* Filemode open + write *)
// 		X.ADD_COM := 00;
// 	ELSE
// 		X.UCB.D_HEAD := WORD#16#F301; (* Filemode create + write *)
// 		X.ADD_COM := 02; (* ADD HEADER *)
// 	END_IF;
//
// 	X.UCB.D_MODE := 1; (* store file name as string *)
// 	UCB(DATA:=X.UCB); (* enter data *)
// 	step_1 := 30;
//
// 30:	(* process LOG trigger *)
// 	trig_tmp := (TRIG_M AND NOT trig_m_last) OR trig_auto OR X.ADD_DATA_REQ;
//
// 	IF ENABLE = FALSE OR fn <> fn_last THEN (* ENABLE = 0 or new file name -> log stop *)
// 		trig_stored := trig_tmp; (* store current trigger pulse *)
// 		X.UCB.D_HEAD := WORD#16#EE00;
// 		X.UCB.D_MODE := 1;
// 		UCB(DATA:=X.UCB); (* write element *)
// 		X.ADD_COM := 00;
// 		step_1 := 00;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_HTML1 - Stop: [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	ELSIF ENABLE THEN
// 		IF trig_tmp OR trig_stored THEN (* manual, automatic or external trigger *)
// 			trig_stored := FALSE;
// 			X.ADD_COM := 03; (* ADD DATA *)
// 			SAVE_DATA.TRIG_CNT := SAVE_DATA.TRIG_CNT + UDINT#1;
// 			SAVE_DATA.TRIG_CNT_TOTAL  := SAVE_DATA.TRIG_CNT_TOTAL + UDINT#1;
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'FILE_HTML1 - Trig: [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		ELSE
// 			X.ADD_COM := 04;	(* ADD DATA REQ *)
// 		END_IF;
// 	END_IF;
//
// END_CASE;
// X.ADD_DATA_REQ := FALSE; (* reset external request *)
//
// (* ##################################### *)
//
// CASE step_2 OF
//
// 00:	aw_ton(IN:= idx > 0 AND aw_enable , PT:= AUTO_CLOSE); (* AutoWrite Timer *)
// 	IF aw_ton.Q THEN (* time-controlled data writing *)
// 		aw_aktiv := TRUE;
// 		step_2 := 10; (* write data *)
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_HTML2 - AutoWrite: [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	ELSE
// 		wd_ton(IN:= FALSE, PT:= wd_time);
// 		wd_ton(IN:= TRUE);
// 		WHILE X.UCB.BUF_COUNT > 0 AND step_2 = 0 AND wd_ton.Q = FALSE DO
//
// 			X.UCB.D_MODE := 10;
// 			UCB(DATA:=X.UCB); (* read element, but do not remove it yet *)
//
// 			IF X.UCB.D_HEAD > WORD#16#F000 THEN (* F301 or F201 *)
//
// 				str := CONCAT('<html><body><table ',HTML_TABLE);
// 		        		str := CONCAT(str,'><caption>');
// 		        		str := CONCAT(str,HTML_CAPTION);
// 		        		str := CONCAT(str,'</caption>');
//
// 				IF idx + LEN(str) < SIZEOF(PT.BUFFER) THEN (* space for element available? *)
//
// 					FSD.MODE := WORD_TO_BYTE(SHR(SHL(X.UCB.D_HEAD,4),12)); (* 0X00 -> 000X = 3 create + write / 2 open + write *)
// 					FSD.OFFSET := UDINT#0; (* start at beginning of file *)
// 					FSD.FILENAME := X.UCB.D_STRING; (* file name *)
// 					PT.SIZE := UINT#0;
// 					SAVE_DATA.FN_REM := FSD.FILENAME; (* save current file name as retained *)
//
// 					X.NEW_FILE := FSD.FILENAME;
// 					X.NEW_FILE_SIZE := UDINT#0;
// 					X.NEW_FILE_RTRIG := FALSE;
// 					step_2 := 20;
//
// 					(* write header data *)
// 					IF FSD.MODE = BYTE#3 THEN
// 						SAVE_DATA.TRIG_CNT := UDINT#0; (* create new file -> then reset counter *)
// 						SAVE_DATA.COLOR := FALSE;
// 						SAVE_DATA.HEAD := FALSE;
//
// 						_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 						idx := idx + LEN(str);
// 						X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
// 					END_IF;
//
// 					X.UCB.D_MODE := 11;
// 					UCB(DATA:=X.UCB); (* remove element *)
//
// 					(* ---------------- Debug-Message -----------------------*)
// 					(*IF _debug_enable THEN
// 						LOG_CL.NEW_MSG := 'FILE_HTML2 - Create: Mode [~1] , [~2]';
// 						LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 						LOG_CL.PRINTF[1] := BYTE_TO_STRING(FSD.MODE);
// 						LOG_MSG();
// 					END_IF;*)
// 					(* -----------------------------------------------------------------*)
// 				END_IF;
//
// 			ELSIF X.UCB.D_HEAD = WORD#16#EE00 THEN (* LOG_STOP Command *)
//
// 				str := '</table></body></html>';
//
// 				IF idx + LEN(str) < SIZEOF(PT.BUFFER) THEN (* space for element available? *)
//
// 					_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 					idx := idx + LEN(str);
// 					log_stop := TRUE;
// 					step_2 := 10; (* write data *)
//
// 					X.UCB.D_MODE := 11;
// 					UCB(DATA:=X.UCB); (* remove element *)
// 				END_IF;
//
// 			ELSIF idx + LEN(X.UCB.D_STRING) + 120 < SIZEOF(PT.BUFFER) THEN (* space for element available? *)
//
// 				(* ---------------- Debug-Message -----------------------*)
// 				(*IF _debug_enable THEN
// 					LOG_CL.NEW_MSG := 'FILE_CSV2 - BUF: TXT [~1],Pos [~2], [~3]';
// 					LOG_CL.PRINTF[1] := X.UCB.D_STRING;
// 					LOG_CL.PRINTF[2] := INT_TO_STRING(idx);
// 					LOG_CL.PRINTF[3] := DT_TO_STRING(DTI);
// 					LOG_MSG();
// 				END_IF;*)
// 				(* -----------------------------------------------------------------*)
//
// 				IF log_state = USINT#0 THEN
// 					str := '<TR ';
// 					IF NOT SAVE_DATA.HEAD THEN
// 						str := CONCAT(str,HTML_TR_HEAD);
// 						SAVE_DATA.HEAD := TRUE;
// 					ELSIF SAVE_DATA.COLOR THEN
// 						str := CONCAT(str,HTML_TR_EVEN);
// 					ELSE
// 						str := CONCAT(str,HTML_TR_ODD);
// 					END_IF;
// 					str := CONCAT(str,'><TD>');
// 				ELSE
// 					str := '<TD>';
// 				END_IF;
//
// 				_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 				idx := idx + LEN(str);
// 				log_state := log_state + USINT#1;
//
// 				_STRING_TO_BUFFER(STR:=X.UCB.D_STRING, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 				idx := idx + LEN(X.UCB.D_STRING);
// 				n := n + USINT#1;
// 				IF n >= X.ID_MAX THEN
// 					n := USINT#0;
// 					str := '</TD></TR>';
// 					SAVE_DATA.COLOR := NOT SAVE_DATA.COLOR;
// 					log_state := USINT#0;
// 				ELSE
// 					str := '</TD>';
// 				END_IF;
// 				_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 				idx := idx + LEN(str);
// 				X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
//
// 				X.UCB.D_MODE := 11;
// 				UCB(DATA:=X.UCB); (* remove element *)
//
// 			ELSE (* no space available *)
// 				step_2 := 10; (* write data *)
//
// 			END_IF;
// 			wd_ton();
// 		END_WHILE;
// 	END_IF;
//
// 10:	IF idx > 0 THEN (* write data *)
// 		FSD.MODE := BYTE#2; (* File_Open_Exist + File_Write *)
// 		FSD.OFFSET := UDINT#4294967294; (* 16#FFFF_FFFE -> special code: automatically append data to end of file *)
// 		PT.SIZE := INT_TO_UINT(idx); (* enter data length *)
// 		idx := 0;
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_HTML2 - WriteData: Mode [~3] Size [~1] , [~2]';
// 			LOG_CL.PRINTF[1] := UINT_TO_STRING(PT.SIZE);
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_CL.PRINTF[3] := BYTE_TO_STRING(FSD.MODE);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
// 	END_IF;
// 	step_2 := 20;
//
// 20:	IF FSD.MODE = BYTE#0 THEN
// 		total_bytes := FSD.FILE_SIZE;
// 		X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
// 		IF log_stop OR aw_aktiv THEN
// 			aw_aktiv := FALSE;
// 			FSD.MODE := BYTE#5; (* close file *)
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'FILE_HTML2 - Close: [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		END_IF;
// 		step_2 := 30;
// 	END_IF;
//
// 30:	IF FSD.MODE = BYTE#0 THEN (* file close finished *)
// 		IF log_stop THEN
// 			X.NEW_FILE_RTRIG := TRUE;
// 			log_stop := FALSE;
// 			SAVE_DATA.FN_REM := ''; (* clear retained file name *)
// 		END_IF;
// 		step_2 := 00;
// 	END_IF;
//
// END_CASE;
//
// FS(FSD:=FSD,PT:=PT); (* FILE_SERVER *)
//
// (* error evaluation *)
// IF FSD.ERROR <> BYTE#0 THEN
// 	ERROR_C := BYTE_TO_DWORD(FSD.ERROR);
// 	ERROR_T := BYTE#1;
// ELSIF X.UCB.BUF_DATA_LOST > UDINT#0 THEN
// 	ERROR_C := DWORD#2; (* error: buffer overflow - reduce amount of data *)
// 	ERROR_T := BYTE#6; 	(* Type: Buffer *)
// ELSIF X.UCB.BUF_USED_MAX > USINT#90 THEN
// 	ERROR_C := DWORD#1; (* warning: buffer utilization > 90% *)
// 	ERROR_T := BYTE#6;	(* Type: Buffer *)
// ELSE
// 	ERROR_C := DWORD#0;
// 	ERROR_T := BYTE#0;
// END_IF;
//
// trig_m_last := TRIG_M;
// END_FUNCTION_BLOCK

// Commented out for beedance: it names a variable fn, which beedance reads as its function keyword.
// FUNCTION_BLOCK DLOG_STORE_FILE_XML
// VAR_IN_OUT
// 	X : DLOG_DATA;
// 	SAVE_DATA : DLOG_SAVE;
// END_VAR
// VAR_INPUT
// 	ENABLE : BOOL;
// 	TRIG_M : BOOL;
// 	TRIG_T : TIME;
// 	FILENAME : STRING;
// 	DTI : DT;
// 	AUTO_CLOSE : TIME := T#15s;
// END_VAR
// VAR_OUTPUT
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	FS : FILE_SERVER;
// 	FSD : FILE_SERVER_DATA;
// 	PT : NETWORK_BUFFER;
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	trig_auto : BOOL;
// 	trig_stored : BOOL;
// 	trig_tmp : BOOL;
// 	trig_m_last : BOOL;
// 	dti_last : DT := DT#2070-02-06-06:28:15;
// 	fn_last : STRING;
// 	fn : STRING;
// 	str : STRING(STRING_LENGTH);
// 	step_1 : INT;
// 	step_2 : INT;
// 	idx : INT;
// 	log_state : USINT;
// 	log_stop : BOOL;
// 	n : INT;
// 	wd_ton : TON;
// 	wd_time : TIME;
// 	aw_ton : TON;
// 	aw_enable : BOOL;
// 	aw_aktiv : BOOL;
// 	total_bytes : UDINT;
//
// 	(*
// 	LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	*)
// END_VAR
//
// X.DTI := DTI; (* enter date / time *)
// X.CLOCK_TRIG := DTI <> dti_last;
// dti_last := DTI;
// trig_auto := FALSE;
// IF X.CLOCK_TRIG THEN (* pulse every second *)
// 	fn := DT_TO_STRF(DTI, 0, FILENAME, 0); (* automatic file name generation *)
// 	IF TRIG_T >= T#1000ms THEN (* automatic trigger *)
// 		trig_auto := DT_TO_UDINT(DTI) MOD (TIME_TO_UDINT(TRIG_T) / UDINT#1000) = UDINT#0;
// 	END_IF;
// END_IF;
//
// CASE step_1 OF
//
// 00:	IF ENABLE THEN
// 		aw_enable := AUTO_CLOSE >= T#5s;
// 		wd_time := SEL(X.LOAD_TIME_MAX > T#0s, T#5ms,X.LOAD_TIME_MAX);
// 		fn_last := fn; (* store current file name *)
// 		IF X.ID_MAX = USINT#0 THEN (* scan for satellite blocks only once *)
// 			X.ADD_COM := 01; (* ADD INFO *)
// 			X.STORE_TYPE := BYTE#3; (* XML mode *)
// 		END_IF;
// 		step_1 := 10;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_XML1 - Start: File [~1], [~2]';
// 			LOG_CL.PRINTF[1] := fn_last;
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	END_IF;
//
// 10:	X.UCB.D_STRING := fn_last;
// 	IF fn_last = SAVE_DATA.FN_REM THEN
// 		X.UCB.D_HEAD := WORD#16#F201; (* Filemode open + write *)
// 		X.ADD_COM := 00;
// 	ELSE
// 		X.UCB.D_HEAD := WORD#16#F301; (* Filemode create + write *)
// 		X.ADD_COM := 02; (* ADD HEADER *)
// 	END_IF;
//
// 	X.UCB.D_MODE := 1; (* store file name as string *)
// 	UCB(DATA:=X.UCB); (* enter data *)
// 	step_1 := 30;
//
// 30:	(* process LOG trigger *)
// 	trig_tmp := (TRIG_M AND NOT trig_m_last) OR trig_auto OR X.ADD_DATA_REQ;
//
// 	IF ENABLE = FALSE OR fn <> fn_last THEN (* ENABLE = 0 or new file name -> log stop *)
// 		trig_stored := trig_tmp; (* store current trigger pulse *)
// 		X.UCB.D_HEAD := WORD#16#EE00;
// 		X.UCB.D_MODE := 1;
// 		UCB(DATA:=X.UCB); (* write element *)
// 		X.ADD_COM := 00;
// 		step_1 := 00;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_XML1 - Stop: [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	ELSIF ENABLE THEN
// 		IF trig_tmp OR trig_stored THEN (* manual, automatic or external trigger *)
// 			trig_stored := FALSE;
// 			X.ADD_COM := 03; (* ADD DATA *)
// 			SAVE_DATA.TRIG_CNT := SAVE_DATA.TRIG_CNT + UDINT#1;
// 			SAVE_DATA.TRIG_CNT_TOTAL  := SAVE_DATA.TRIG_CNT_TOTAL + UDINT#1;
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'FILE_XML1 - Trig: [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		ELSE
// 			X.ADD_COM := 04;	(* ADD DATA REQ *)
// 		END_IF;
// 	END_IF;
//
// END_CASE;
// X.ADD_DATA_REQ := FALSE; (* reset external request *)
//
// (* ##################################### *)
//
// CASE step_2 OF
//
// 00:	aw_ton(IN:= idx > 0 AND aw_enable , PT:= AUTO_CLOSE); (* AutoWrite Timer *)
// 	IF aw_ton.Q THEN (* time-controlled data writing *)
// 		aw_aktiv := TRUE;
// 		step_2 := 10; (* write data *)
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_XML2 - AutoWrite: [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	ELSE
// 		wd_ton(IN:= FALSE, PT:= wd_time);
// 		wd_ton(IN:= TRUE);
// 		WHILE X.UCB.BUF_COUNT > 0 AND step_2 = 0 AND wd_ton.Q = FALSE DO
//
// 			X.UCB.D_MODE := 10;
// 			UCB(DATA:=X.UCB); (* read element, but do not remove it yet *)
//
// 			IF X.UCB.D_HEAD > WORD#16#F000 THEN (* F301 or F201 *)
//
// 				str := '<?xml version="1.0" encoding="UTF-8"?><table>';
//
// 				IF idx + LEN(str) < SIZEOF(PT.BUFFER) THEN (* space for element available? *)
//
// 					FSD.MODE := WORD_TO_BYTE(SHR(SHL(X.UCB.D_HEAD,4),12)); (* 0X00 -> 000X = 3 create + write / 2 open + write *)
// 					FSD.OFFSET := UDINT#0; (* start at beginning of file *)
// 					FSD.FILENAME := X.UCB.D_STRING; (* file name *)
// 					PT.SIZE := UINT#0;
// 					SAVE_DATA.FN_REM := FSD.FILENAME; (* save current file name as retained *)
//
// 					X.NEW_FILE := FSD.FILENAME;
// 					X.NEW_FILE_SIZE := UDINT#0;
// 					X.NEW_FILE_RTRIG := FALSE;
// 					step_2 := 20;
//
// 					(* write header data *)
// 					IF FSD.MODE = BYTE#3 THEN
// 						SAVE_DATA.TRIG_CNT := UDINT#0; (* create new file -> then reset counter *)
//
// 						_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 						idx := idx + LEN(str);
// 						X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
// 					END_IF;
//
// 					X.UCB.D_MODE := 11;
// 					UCB(DATA:=X.UCB); (* remove element *)
//
// 					(* ---------------- Debug-Message -----------------------*)
// 					(*IF _debug_enable THEN
// 						LOG_CL.NEW_MSG := 'FILE_XML2 - Create: Mode [~1] , [~2]';
// 						LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 						LOG_CL.PRINTF[1] := BYTE_TO_STRING(FSD.MODE);
// 						LOG_MSG();
// 					END_IF;*)
// 					(* -----------------------------------------------------------------*)
// 				END_IF;
//
// 			ELSIF X.UCB.D_HEAD = WORD#16#EE00 THEN (* LOG_STOP Command *)
//
// 				str := '</table>';
//
// 				IF idx + LEN(str) < SIZEOF(PT.BUFFER) THEN (* space for element available? *)
//
// 					_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 					idx := idx + LEN(str);
// 					log_stop := TRUE;
// 					step_2 := 10; (* write data *)
//
// 					X.UCB.D_MODE := 11;
// 					UCB(DATA:=X.UCB); (* remove element *)
//
// 				END_IF;
//
// 			ELSIF idx + LEN(X.UCB.D_STRING) + 120 < SIZEOF(PT.BUFFER) THEN (* space for element available? *)
//
// 				(* ---------------- Debug-Message -----------------------*)
// 				(*IF _debug_enable THEN
// 					LOG_CL.NEW_MSG := 'FILE_XML2 - BUF: TXT [~1],Pos [~2], [~3]';
// 					LOG_CL.PRINTF[1] := X.UCB.D_STRING;
// 					LOG_CL.PRINTF[2] := INT_TO_STRING(idx);
// 					LOG_CL.PRINTF[3] := DT_TO_STRING(DTI);
// 					LOG_MSG();
// 				END_IF;*)
// 				(* -----------------------------------------------------------------*)
//
// 				IF log_state = USINT#0 THEN
// 					str := '<row><entry>';
// 				ELSE
// 					str := '<entry>';
// 				END_IF;
//
// 				_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 				idx := idx + LEN(str);
// 				log_state := log_state + USINT#1;
//
// 				_STRING_TO_BUFFER(STR:=X.UCB.D_STRING, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 				idx := idx + LEN(X.UCB.D_STRING);
//
// 				n := n + USINT#1;
// 				IF n >= X.ID_MAX THEN
// 					n := USINT#0;
// 					str := '</entry></row>';
// 					log_state := USINT#0;
// 				ELSE
// 					str := '</entry>';
// 				END_IF;
// 				_STRING_TO_BUFFER(STR:=str, POS:=idx, PT:=REF(PT.BUFFER), SIZE:=SIZEOF(PT.BUFFER));
// 				idx := idx + LEN(str);
// 				X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
//
// 				X.UCB.D_MODE := 11;
// 				UCB(DATA:=X.UCB); (* remove element *)
//
// 			ELSE (* no space available *)
// 				step_2 := 10; (* write data *)
// 			END_IF;
// 			wd_ton();
// 		END_WHILE;
// 	END_IF;
//
// 10:	IF idx > 0 THEN (* write data *)
// 		FSD.MODE := BYTE#2; (* File_Open_Exist + File_Write *)
// 		FSD.OFFSET := UDINT#4294967294; (* 16#FFFF_FFFE -> special code: automatically append data to end of file *)
// 		PT.SIZE := INT_TO_UINT(idx); (* enter data length *)
// 		idx := 0;
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_XML2 - WriteData: Mode [~3] Size [~1] , [~2]';
// 			LOG_CL.PRINTF[1] := UINT_TO_STRING(PT.SIZE);
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_CL.PRINTF[3] := BYTE_TO_STRING(FSD.MODE);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
// 	END_IF;
// 	step_2 := 20;
//
// 20:	IF FSD.MODE = BYTE#0 THEN
// 		total_bytes := FSD.FILE_SIZE;
// 		X.NEW_FILE_SIZE := total_bytes + INT_TO_UDINT(idx);
// 		IF log_stop OR aw_aktiv THEN
// 			aw_aktiv := FALSE;
// 			FSD.MODE := BYTE#5; (* close file *)
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'FILE_XML2 - Close: [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		END_IF;
// 		step_2 := 30;
// 	END_IF;
//
// 30:	IF FSD.MODE = BYTE#0 THEN (* file close finished *)
// 		IF log_stop THEN
// 			X.NEW_FILE_RTRIG := TRUE;
// 			log_stop := FALSE;
// 			SAVE_DATA.FN_REM := ''; (* clear retained file name *)
// 		END_IF;
// 		step_2 := 00;
// 	END_IF;
//
// END_CASE;
//
// FS(FSD:=FSD,PT:=PT); (* FILE_SERVER *)
//
// (* error evaluation *)
// IF FSD.ERROR <> BYTE#0 THEN
// 	ERROR_C := BYTE_TO_DWORD(FSD.ERROR);
// 	ERROR_T := BYTE#1;
// ELSIF X.UCB.BUF_DATA_LOST > UDINT#0 THEN
// 	ERROR_C := DWORD#2; (* error: buffer overflow - reduce amount of data *)
// 	ERROR_T := BYTE#6; 	(* Type: Buffer *)
// ELSIF X.UCB.BUF_USED_MAX > USINT#90 THEN
// 	ERROR_C := DWORD#1; (* warning: buffer utilization > 90% *)
// 	ERROR_T := BYTE#6;	(* Type: Buffer *)
// ELSE
// 	ERROR_C := DWORD#0;
// 	ERROR_T := BYTE#0;
// END_IF;
//
// trig_m_last := TRIG_M;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_STORE_MYSQL
// VAR_IN_OUT
// 	SQL_INFO : MYSQL_INFO;
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	ENABLE : BOOL;
// 	TRIG_M : BOOL;
// 	TRIG_T : TIME;
// 	DTI : DT;
// 	URL : STRING(STRING_LENGTH);
// 	DB_NAME : STRING(64);
// 	TB_NAME : STRING(64);
// 	AUTO_DELETE : UDINT;
// 	TIMEOUT : TIME;
// 	DNS_IP4 : DWORD;
// END_VAR
// VAR_OUTPUT
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	MYSQL_CONTROL: MYSQL_CONTROL;
// 	trig_auto : BOOL;
// 	trig_tmp : BOOL;
// 	trig_m_last : BOOL;
// 	dti_last : DT := DT#2070-02-06-06:28:15;
// 	enable_last : BOOL;
// 	step : INT;
// 	idx : INT;
// 	n : USINT;
// 	snd_text : STRING(STRING_LENGTH);
// 	snd_size : INT;
// 	sql_command : BYTE;
// 	next_step : INT;
// 	i1 : INT;
// 	i2 : INT;
// 	y: MYSQL_COM;
// 	del_dt : UDINT;
// 	del_start : BOOL;
// 	run : BOOL;
//
// 	(*
// 	LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	*)
// END_VAR
//
// X.DTI := DTI; (* enter date / time *)
// X.CLOCK_TRIG := DTI <> dti_last;
// dti_last := DTI;
// trig_auto := FALSE;
// IF X.CLOCK_TRIG THEN (* pulse every second *)
// 	del_start := DT_TO_UDINT(DTI) - del_dt > AUTO_DELETE AND AUTO_DELETE > UDINT#0;
// 	IF TRIG_T >= T#1000ms THEN (* automatic trigger *)
// 		trig_auto := DT_TO_UDINT(DTI) MOD (TIME_TO_UDINT(TRIG_T) / UDINT#1000) = UDINT#0;
// 	END_IF;
// END_IF;
// X.ADD_COM := 00;
//
// CASE step OF
//
// 00:	IF ENABLE AND enable_last = FALSE THEN
// 		y.SQL_URL := URL;
// 		y.TIMEOUT := TIMEOUT;
// 		y.DNS_IP4 := DNS_IP4;
// 		y.SQL_CON := TRUE;
//
// 		X.UCB.D_MODE := 99; (* buffer reset - delete everything *)
// 		UCB(DATA:=X.UCB);
//
// 		IF X.ID_MAX = USINT#0 THEN (* scan for satellite blocks only once *)
// 			X.ADD_COM := 01; (* ADD INFO *)
// 			X.STORE_TYPE := BYTE#5; (* MySQL mode *)
// 		END_IF;
//
// 		step := 05;
// 	END_IF;
//
// 05:	X.ADD_COM := 02; (* ADD HEADER *)
// 	step := 10;
//
// 10:	IF SQL_INFO.SQL_CONNECTED THEN
// 		sql_command := BYTE#3; (* CREATE DATABASE IF NOT EXISTS `db_name` *)
// 		snd_text :=  CONCAT('CREATE DATABASE IF NOT EXISTS `' ,DB_NAME);
// 		snd_text :=  CONCAT(snd_text, '`');
// 		next_step := 20;
// 	END_IF;
//
// 20:	sql_command := BYTE#2; (* 0x02 COM_INIT_DB mysql_select_db *)
// 	snd_text := DB_NAME;
// 	next_step := 30;
//
// 30:	(* CREATE TABLE IF NOT EXISTS tb_name (Timestamp TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP, n1,n2) *)
// 	sql_command := BYTE#0;
// 	idx := 5;
// 	n := USINT#0;
// 	WHILE sql_command = BYTE#0 DO
// 		IF n = USINT#0 THEN
// 			snd_text := CONCAT('CREATE TABLE IF NOT EXISTS `', TB_NAME);
// 			snd_text := CONCAT(snd_text,'` (`ID` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY ,');
// 			snd_text := CONCAT(snd_text,' store_timestamp TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP');
// 		ELSIF n <= X.ID_MAX THEN
// 			X.UCB.D_MODE := 12; (* read and delete element *)
// 			UCB(DATA:=X.UCB);
// 			snd_text := CONCAT(', `',X.UCB.D_STRING);
// 			CASE WORD_TO_INT(SHR(X.UCB.D_HEAD,8)) OF
// 			02:	snd_text := CONCAT(snd_text, '` FLOAT'); (* REAL *)
// 			06:	snd_text := CONCAT(snd_text, '` BOOL'); (* BOOL *)
// 			08:	snd_text := CONCAT(snd_text, '` INT'); (* DINT *)
// 			13:	snd_text := CONCAT(snd_text, '` DATETIME'); (* DT *)
// 			ELSE
// 				snd_text := CONCAT(snd_text, '` VARCHAR(80)'); (* STRING(80) *)
// 			END_CASE;
// 		ELSE
// 			snd_text := ')';
// 			sql_command := BYTE#3;
// 			snd_size := idx + LEN(snd_text);
// 			next_step := 40;
// 		END_IF;
// 		_STRING_TO_BUFFER(str:=snd_text,pos:=idx,pt:=REF(y.S_BUF.BUFFER),size:=SIZEOF(y.S_BUF.BUFFER)); (* copy string into buffer *)
// 		idx := idx + LEN(snd_text);
// 		n := n + USINT#1;
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'dlog_mysql: create >~1<';
// 			LOG_CL.PRINTF[1] := snd_text;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
// 	END_WHILE;
//
// 40:	IF del_start THEN
// 		del_start := FALSE;
// 		sql_command := BYTE#3; (* DELETE FROM `tb_name` WHERE TIMESTAMPDIFF(SECOND, Timestamp, NOW()) > xxx *)
// 		snd_text := CONCAT('DELETE FROM `', TB_NAME);
// 		snd_text := CONCAT(snd_text, '` WHERE TIMESTAMPDIFF(SECOND, store_timestamp, NOW()) > ');
// 		snd_text := CONCAT(snd_text, UDINT_TO_STRING(AUTO_DELETE));
// 		next_step := 40;
// 		del_dt := DT_TO_UDINT(DTI);
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'dlog_mysql: # delete time:~1 >~2<';
// 			LOG_CL.PRINTF[1] := UDINT_TO_STRING(AUTO_DELETE);
// 			LOG_CL.PRINTF[2] := snd_text;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
//          (* INSERT INTO `tb_name` VALUES (CURRENT_TIMESTAMP,  "n1","n2") *)
// 	ELSIF X.UCB.BUF_COUNT >= USINT_TO_INT(X.ID_MAX) AND step > 0 THEN (* new data row(s) available *)
// 		idx := 5;
// 		n := USINT#0;
// 		sql_command := BYTE#0;
// 		WHILE sql_command = BYTE#0 DO
// 			IF n = USINT#0 THEN
// 				snd_text := CONCAT('INSERT INTO `', TB_NAME);
// 				snd_text := CONCAT(snd_text,'` VALUES (NULL, CURRENT_TIMESTAMP');
// 			ELSIF n <= X.ID_MAX THEN
// 				X.UCB.D_MODE := 12; (* read and delete element *)
// 				UCB(DATA:=X.UCB);
// 				snd_text := CONCAT(', "',X.UCB.D_STRING);
// 				snd_text := CONCAT(snd_text, '"');
// 			ELSE
// 				snd_text := ')';
// 				sql_command := BYTE#3;
// 				snd_size := idx + LEN(snd_text);
// 				next_step := 40;
// 			END_IF;
// 			_STRING_TO_BUFFER(str:=snd_text,pos:=idx,pt:=REF(y.S_BUF.BUFFER),size:=SIZEOF(y.S_BUF.BUFFER)); (* copy string into buffer *)
// 			idx := idx + LEN(snd_text);
// 			n := n + USINT#1;
//
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'dlog_mysql: insert >~1<';
// 				LOG_CL.PRINTF[1] := snd_text;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 		END_WHILE;
// 	END_IF;
//
// 	IF ENABLE THEN 	(* process LOG trigger *)
// 		trig_tmp := (TRIG_M AND NOT trig_m_last) OR trig_auto OR X.ADD_DATA_REQ;
//
// 		IF trig_tmp THEN (* manual, automatic or external trigger *)
// 			X.ADD_COM := 03; (* ADD DATA *)
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'dlog_mysql: Trigger: DATE [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		ELSE
// 			X.ADD_COM := 04;	(* ADD DATA REQ *)
// 		END_IF;
// 	ELSE
// 		y.SQL_CON := FALSE; (* close MySQL connection *)
// 		step := 0;
// 	END_IF;
// 	X.ADD_DATA_REQ := FALSE;
//
// (* ~~~~~~~~~~~~~~~~~~~~~~ *)
//
// 199:IF y.SQL_RCV_STATE = BYTE#1 THEN (* OK *)
// 		step := next_step;
// 	ELSIF y.SQL_RCV_STATE = BYTE#2 THEN  (* ERROR *)
// 		step := 0;
// 	END_IF;
//
// END_CASE;
//
// (* ########################## *)
// (* execute SQL command *)
// IF sql_command > BYTE#0 THEN
//
// 	IF snd_size > 0 THEN
// 		y.S_BUF.SIZE := INT_TO_UINT(snd_size);
// 		snd_size := 0;
// 	ELSE
// 		_STRING_TO_BUFFER(str:=snd_text,pos:=5,pt:=REF(y.S_BUF.BUFFER),size:=SIZEOF(y.S_BUF.BUFFER)); (* copy string into buffer *)
// 		y.S_BUF.SIZE := INT_TO_UINT(5 + LEN(snd_text));
// 		snd_text := '';
// 	END_IF;
//
// 	IF y.S_BUF.SIZE > SIZEOF(y.S_BUF.BUFFER) THEN (* buffer overflow detected -> abort *)
// 		y.S_BUF.SIZE := 0;
// 		sql_command := BYTE#0;
// 		y.SQL_CON := FALSE; (* close MySQL connection *)
// 		step := 999; (* dummy step *)
// 	ELSE
//
// 		i2 := UINT_TO_INT(y.S_BUF.SIZE) - 1; (* replace all " characters with ' *)
// 		FOR i1 := 5 TO i2 DO
// 			IF y.S_BUF.BUFFER[i1] = BYTE#34 THEN
// 				y.S_BUF.BUFFER[i1] := BYTE#39;
// 			END_IF;
// 		END_FOR;
//
// 		y.S_BUF.BUFFER[4] := sql_command;
// 		y.SQL_PACKET_NO := BYTE#255;
// 		sql_command := BYTE#0;
// 		step := 199;
// 	END_IF;
// END_IF;
//
// (* ########################## *)
//
// MYSQL_CONTROL(COM:= y, INFO:=SQL_INFO);
//
// (* error evaluation *)
// IF step = 999 THEN	(* buffer overflow *)
// 	ERROR_C := DWORD#3; (* network buffer overflow *)
// 	ERROR_T := BYTE#6; 	(* Type: Buffer *)
// ELSIF y.ERROR_T > BYTE#0 THEN	(* SQL Error *)
// 	ERROR_C := y.ERROR_C;
// 	ERROR_T := y.ERROR_T;
// 	step := 0;
// ELSIF X.UCB.BUF_DATA_LOST > UDINT#0 THEN
// 	ERROR_C := DWORD#2; (* error: buffer overflow - reduce amount of data *)
// 	ERROR_T := BYTE#6; 	(* Type: Buffer *)
// ELSIF X.UCB.BUF_USED_MAX > USINT#90 THEN
// 	ERROR_C := DWORD#1; (* warning: buffer utilization > 90% *)
// 	ERROR_T := BYTE#6;	(* Type: Buffer *)
// ELSE
// 	ERROR_C := DWORD#0;
// 	ERROR_T := BYTE#0;
// END_IF;
//
// (* ########################## *)
//
// trig_m_last := TRIG_M;
// enable_last := ENABLE;
// run := step > 0;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_STORE_RRD
// VAR_IN_OUT
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	ENABLE : BOOL;
// 	TRIG_M : BOOL;
// 	TRIG_T : TIME;
// 	URL : STRING(STRING_LENGTH);
// 	DTI : DT;
// 	SEP : BYTE;
// 	DNS_IP4 : DWORD;
// 	TIMEOUT : TIME;
// END_VAR
// VAR_OUTPUT
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	HTTP_GET : HTTP_GET;
// 	IPC : IP_CONTROL;
//
// 	trig_auto : BOOL;
// 	trig_stored : BOOL;
// 	trig_tmp : BOOL;
// 	trig_m_last : BOOL;
// 	dti_last : DT := DT#2070-02-06-06:28:15;
// 	step_1 : INT;
// 	step_2 : INT;
// 	sep_char : STRING(1);
// 	url_query : STRING(STRING_LENGTH);
// 	n : INT;
//
// 	(*
// 	LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	*)
// END_VAR
//
// X.DTI := DTI; (* enter date / time *)
// X.CLOCK_TRIG := DTI <> dti_last;
// dti_last := DTI;
// trig_auto := FALSE;
// IF X.CLOCK_TRIG THEN (* pulse every second *)
// 	IF TRIG_T >= T#1000ms THEN (* automatic trigger *)
// 		trig_auto := DT_TO_UDINT(DTI) MOD (TIME_TO_UDINT(TRIG_T) / UDINT#1000) = UDINT#0;
// 	END_IF;
// END_IF;
//
// CASE step_1 OF
//
// 00:	IF ENABLE THEN
// 		IF X.ID_MAX = USINT#0 THEN (* scan for satellite blocks only once *)
// 			X.ADD_COM := 01; (* ADD INFO *)
// 			X.STORE_TYPE := BYTE#1; (* RRD mode *)
// 		END_IF;
//
// 		ERROR_T := BYTE#0;
// 		ERROR_C := DWORD#0;
//
// 		X.UCB.D_MODE := 99;
// 		UCB(DATA:=X.UCB); (* clear buffer *)
//
// 		URL_DATA:=STRING_TO_URL(STR:=URL,DEFAULT_PROTOCOL:='',DEFAULT_PATH:=''); (* split URL *)
// 		sep_char := CHR_TO_STRING(SEP);
// 		url_query := URL_DATA.QUERY;
// 		step_1 := 30;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_RRD1 - Start: DATE [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	END_IF;
//
// 30:	(* process LOG trigger *)
// 	trig_tmp := (TRIG_M AND NOT trig_m_last) OR trig_auto OR X.ADD_DATA_REQ;
//
// 	IF ENABLE = FALSE THEN (* ENABLE = 0 -> Log Stop *)
// 		X.ADD_COM := 00;
// 		step_1 := 00;
//
// 		(* ---------------- Debug-Message -----------------------*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'FILE_RRD1 - Stop: DATE [~2]';
// 			LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 			LOG_MSG();
// 		END_IF;*)
// 		(* -----------------------------------------------------------------*)
//
// 	ELSIF ENABLE THEN
// 		IF trig_tmp OR trig_stored THEN (* manual, automatic or external trigger *)
// 			trig_stored := FALSE;
// 			X.ADD_COM := 03; (* ADD DATA *)
// 			(* ---------------- Debug-Message -----------------------*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'FILE_RRD1 - Trig: DATE [~2]';
// 				LOG_CL.PRINTF[2] := DT_TO_STRING(DTI);
// 				LOG_MSG();
// 			END_IF;*)
// 			(* -----------------------------------------------------------------*)
// 		ELSE
// 			X.ADD_COM := 04;	(* ADD DATA REQ *)
// 		END_IF;
// 	END_IF;
//
// END_CASE;
// X.ADD_DATA_REQ := FALSE; (* reset external request *)
//
// (* ##################################### *)
//
// CASE step_2 OF
//
// 00:	IF X.UCB.BUF_DATA_LOST > UDINT#0 THEN
// 		ERROR_C := DWORD#2; (* error: buffer overflow - reduce amount of data *)
// 		ERROR_T := BYTE#6; 	(* Type: Buffer *)
// 		step_2 := 100;
// 	ELSIF X.UCB.BUF_USED_MAX > USINT#90 THEN
// 		ERROR_C := DWORD#1; (* warning: buffer utilization > 90% *)
// 		ERROR_T := BYTE#6;	(* Type: Buffer *)
// 	END_IF;
//
// 	IF X.UCB.BUF_COUNT >= USINT_TO_INT(X.ID_MAX) THEN (* enough entries available *)
// 		URL_DATA.QUERY := url_query; (* restore original query *)
// 		n := 1;
// 		WHILE n <= USINT_TO_INT(X.ID_MAX) DO
// 			X.UCB.D_MODE := 12; (* read and delete element *)
// 			UCB(DATA:=X.UCB);
// 			IF n < USINT_TO_INT(X.ID_MAX) THEN (* add separator *)
// 				X.UCB.D_STRING := CONCAT(X.UCB.D_STRING,sep_char);
// 			END_IF;
// 			IF LEN(X.UCB.D_STRING) + LEN(URL_DATA.QUERY) <= STRING_LENGTH THEN
// 				URL_DATA.QUERY := CONCAT(URL_DATA.QUERY,X.UCB.D_STRING); (* append LOG data to query *)
// 				step_2 := 40;
// 			ELSE
// 				ERROR_T := BYTE#4; (* LOG data cannot be packed into the URL *)
// 				ERROR_C := DWORD#1;
// 				step_2 := 100;
// 			END_IF;
// 			n := n + 1;
// 		END_WHILE;;
// 	END_IF;
//
// 40:	IF DNS_CLIENT.DONE THEN
// 		step_2 := 60;
// 	ELSIF (DNS_CLIENT.ERROR > DWORD#0) THEN
// 		ERROR_C := DNS_CLIENT.ERROR; (* error handling *)
// 		ERROR_T := BYTE#1;
// 		step_2 := 100;
// 	END_IF;
//
// 60:	IF HTTP_GET.DONE THEN
// 		n := UINT_TO_INT(HTTP_GET.BODY_START);
// 		IF R_BUF.BUFFER[n] <> BYTE#48 THEN (* ASCII '0' *)
// 			ERROR_C := DWORD#1;
// 			ERROR_T := BYTE#3;
// 		END_IF;
// 		step_2 := 100;
//
// 	ELSIF HTTP_GET.ERROR > DWORD#0 THEN (* error handling *)
// 		ERROR_C := HTTP_GET.ERROR;
// 		ERROR_T := BYTE#2;
// 		step_2 := 100;
// 	END_IF;
//
// 100:	IF (NOT HTTP_GET.DONE) THEN (* UNLOCK HTTP DATA *)
// 		step_2 := SEL(ERROR_T = BYTE#0, 120,0);
// 	END_IF;
//
// 120:	IF NOT ENABLE THEN
// 		step_2 := 0;
// 	END_IF;
//
// END_CASE;
//
// HTTP_GET(IP_C:=IP_C,S_BUF:=S_BUF,R_BUF:=R_BUF,IP4:=DNS_CLIENT.IP4,GET:=step_2=60,MODE:=2,UNLOCK_BUF:=step_2=100,URL_DATA:=URL_DATA); (* HTTP_GET *)
// DNS_CLIENT(IP_C:=IP_C, S_BUF:=S_BUF, R_BUF:=R_BUF, DOMAIN:=URL_DATA.DOMAIN, IP4_DNS:= DNS_IP4, ACTIVATE:=step_2=40); (* DNS-Client *)
// IPC(IP:= 0, PORT:= 0 , TIME_OUT:= TIMEOUT, IP_C:= IP_C, S_BUF:= S_BUF, R_BUF:= R_BUF); (* IP_CONTROL *)
//
// trig_m_last := TRIG_M;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses UNI_CIRCULAR_BUFFER, which is commented out.
// FUNCTION_BLOCK DLOG_STRING
// VAR_IN_OUT
// 	X : DLOG_DATA;
// END_VAR
// VAR_INPUT
// 	STR : STRING;
// 	COLUMN : STRING(40);
// END_VAR
// VAR_OUTPUT
// END_VAR
// VAR
// 	UCB : UNI_CIRCULAR_BUFFER;
// 	id : WORD;
// 	str_last : STRING;
// END_VAR
//
// CASE x.ADD_COM OF
//
// 01:	(* ADD INFO *)
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	id := WORD#16#0101; (* source type STRING, target type STRING *)
// 02:	(* ADD HEADER *)
// 	X.UCB.D_STRING := COLUMN;
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// 03:	(* ADD DATA *)
// 	X.UCB.D_STRING := str_last;
// 	X.UCB.D_HEAD := id;
// 	X.UCB.D_MODE := 1;
// 	UCB(DATA:=X.UCB); (* enter data *)
// END_CASE;
// str_last := STR;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK DNS_CLIENT
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// END_VAR
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	DOMAIN : STRING;
// 	IP4_DNS : DWORD;
// END_VAR
// VAR_OUTPUT
// 	IP4 : DWORD;
// 	DONE : BOOL;
// 	ERROR : DWORD;
// END_VAR
// VAR
// 	IP_FIFO : IP_FIFO;
// 	ip_state : BYTE;
// 	ip_id : BYTE;
// 	i : INT;
// 	state : INT;
// 	domain_copy : STRING;
// 	activate_last : BOOL;
// 	url_length : INT;
// 	char : BYTE;
// 	char_count : INT;
// 	tid : BYTE;
// 	anc : INT;
// 	x : BYTE;
// 	rr_stop : BOOL;
// 	ttl_last : DWORD;
// 	ttl_tmp : DWORD;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;*)
//
// END_VAR
//
// (* Domain Name System (query) *)
//
// CASE state OF
//
// 00:	IF ACTIVATE AND NOT activate_last THEN (* wait for rising edge on GET *)
// 		DONE := FALSE;
// 		ERROR := DWORD#0;
// 		state := 5;
// 	END_IF;
//
// 05:	IP4 := IP4_DECODE(str:=DOMAIN); (* check whether domain is already an IP4 address *)
// 	IF IP4 > DWORD#00 THEN
// 		DONE := TRUE;
// 		state := 0;
//         		(*####################  DEBUG-MESSAGE  ###################################*)
//         		(*IF _debug_enable THEN
//             		LOG_CL.NEW_MSG := 'DNS: S05 IP4 ~1';
//             		LOG_CL.PRINTF[1]  := DOMAIN;
//             		LOG_MSG();
//         		END_IF;*)
//         		(*########################################################################*)
// 	ELSE
// 		ip_state := BYTE#1;
// 		ttl_last := DWORD#0;
// 		domain_copy := DOMAIN; (* save DOMAIN string *)
// 		state := 10;
//         		(*####################  DEBUG-MESSAGE  ###################################*)
//         		(*IF _debug_enable THEN
//             		LOG_CL.NEW_MSG := 'DNS: S05 START';
//             		LOG_MSG();
//         		END_IF;*)
//         		(*########################################################################*)
// 	END_IF;
//
// 10:	IF ip_state = BYTE#3 THEN (* register with IP_Control *)
// 		url_length := LEN(domain_copy);
// 		char_count := 17 + url_length; (* initialize send buffer *)
// 		FOR i := 0 TO char_count DO
// 			S_BUF.BUFFER[i] := BYTE#00;
// 		END_FOR;
// 		tid := UINT_TO_BYTE(BYTE_TO_UINT(tid) + UINT#1);
// 		S_BUF.BUFFER[01] := tid; (* Transaction ID lowbyte *)
// 		S_BUF.BUFFER[02] := BYTE#16#01; (* Flags: highbyte *)
// 		S_BUF.BUFFER[05] := BYTE#16#01; (* Questions: lowbyte *)
// 		i := url_length + 12;	(* process the server name from back to front *)
// 		char_count := 0;
// 		WHILE i > 12 DO
// 			char := CODE(STR:=domain_copy,POS:=i - 12);
// 			IF char = BYTE#46 THEN (* check for main/sub domain name separator "." *)
// 				S_BUF.BUFFER[i] := INT_TO_BYTE(char_count); (* instead of the separator "." the length of the detected main/sub domain name is entered *)
// 				char_count := 0;
// 			ELSE
// 				S_BUF.BUFFER[i] := char; (* enter char into buffer *)
// 				char_count := char_count + 1; (* count number of chars *)
// 			END_IF;
// 			i := i - 1;
// 		END_WHILE;
// 		S_BUF.BUFFER[i] := INT_TO_BYTE(char_count); (* enter the length of the leftmost (first) domain name *)
// 		i := url_length + 15; (* position after the server name *)
// 		S_BUF.BUFFER[i] := BYTE#01; (* Type: A Hostaddress lowbyte *)
// 		i := i + 2;
// 		S_BUF.BUFFER[i] := BYTE#01; (* Class: IN lowbyte *)
// 		i := i + 1;
// 		(* set up IP data traffic *)
// 		IP_C.C_PORT := WORD#53; (* enter port number *)
// 		IP_C.C_IP := IP4_DNS; (* enter IP *)
// 		IP_C.C_MODE := BYTE#1; (* Mode: UDP+ACTIVE+PORT+IP *)
// 		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
// 		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
// 		IP_C.R_OBSERVE  := TRUE; (* monitor data reception *)
// 		S_BUF.SIZE := INT_TO_UINT(i); (* enter send length *)
// 		R_BUF.SIZE := UINT#0; (* reset receive length *)
// 		state := 30;
//
//         		(*####################  DEBUG-MESSAGE  ###################################*)
//         		(*IF _debug_enable THEN
//             		LOG_CL.NEW_MSG := 'DNS: S10 SEND';
//             		LOG_MSG();
//         		END_IF;*)
//         		(*########################################################################*)
// 	END_IF;
//
// 30:	IF IP_C.ERROR <> DWORD#00 THEN
// 		ERROR := IP_C.ERROR;
// 		state := 99;
//         		(*####################  DEBUG-MESSAGE  ###################################*)
//         		(*IF _debug_enable THEN
//             		LOG_CL.NEW_MSG := 'DNS: S30 ERROR ~1';
//             		LOG_CL.PRINTF[1]  := DWORD_TO_STRING(ERROR);
//             		LOG_MSG();
//         		END_IF;*)
//         		(*########################################################################*)
//
// 	ELSIF (S_BUF.SIZE = UINT#0) AND (tid = R_BUF.BUFFER[1]) AND (R_BUF.SIZE >= INT_TO_UINT(34 + url_length)) THEN	(* check whether the response telegram has the correct transaction ID and minimum length *)
// 		ERROR := BYTE_TO_DWORD(R_BUF.BUFFER[3] AND BYTE#2#0000_1111); (* check return code of the flags. 4 bits *)
// 		IF ERROR = DWORD#00 THEN
// 			anc := R_BUF.BUFFER[7]; (* number of answer RRs *)
// 			IP_C.R_OBSERVE := FALSE; (* monitor data reception *)
// 			state := 40;
// 		ELSE
// 			state := 99;
// 		END_IF;
// 	END_IF;
//
// 40:	IF anc > 0 THEN
// 		anc := anc - 1;
// 		WHILE i < R_BUF.SIZE DO
// 			x := R_BUF.BUFFER[i];
// 			IF rr_stop = FALSE THEN
// 				IF x > BYTE#63 THEN (* pointer to name *)
// 					i := i + 1;
// 					rr_stop := TRUE;
// 				ELSIF x= BYTE#0 THEN (* end of name *)
// 					rr_stop := TRUE;
// 				ELSE
// 					i := i + x; (* names *)
// 				END_IF;
// 				i := i + 1;
// 			ELSE
// 				rr_stop := FALSE;
// 				IF R_BUF.BUFFER[i+1] = BYTE#1 AND R_BUF.BUFFER[i+3] = BYTE#1 THEN (* Type A / Class IN *)
// 					i := i + 4;
// 					ttl_tmp := DWORD_OF_BYTE(R_BUF.BUFFER[i],R_BUF.BUFFER[i+1],R_BUF.BUFFER[i+2],R_BUF.BUFFER[i+3]); (* Time to live *)
// 					i := i + 6;
// 					IF ttl_tmp >= ttl_last THEN (* TTL is longer *)
// 						ttl_last := ttl_tmp;
// 						IP4 := DWORD_OF_BYTE(R_BUF.BUFFER[i],R_BUF.BUFFER[i+1],R_BUF.BUFFER[i+2],R_BUF.BUFFER[i+3]); (* IP address *)
// 					END_IF;
// 					i := i + 4;
// 				ELSE
// 					i := i + 9;
// 					i := i + R_BUF.BUFFER[i] + 1; (* offset to next block *)
// 				END_IF;
// 				EXIT;
// 			END_IF;
// 		END_WHILE;
// 	ELSE
// 		IF IP4 > DWORD#0 THEN (* valid RR with IP found *)
// 			DONE := TRUE;
//             		(*####################  DEBUG-MESSAGE  ###################################*)
//             		(*IF _debug_enable THEN
//                 			LOG_CL.NEW_MSG := 'DNS: S40 OK';
//                 			LOG_MSG();
//             		END_IF;*)
//             		(*########################################################################*)
// 		ELSE
// 			ERROR := DWORD#255; (* data faulty - no IP address present *)
// 		END_IF;
// 		state := 99;
// 	END_IF;
//
// 99:	R_BUF.SIZE := UINT#0; (* reset receive length *)
// 	ip_state := BYTE#4; (* unregister *)
// 	state := 00;
//
// END_CASE;
//
// activate_last := ACTIVATE; (* save last status for edge detection *)
// IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id); (* IP_FIFO *)
//
// (*----------------------------------------------------------------------------------------------------------------*)
// (* Flags: 0x0100 (Standard query)                                                                                            *)
// (* 0... .... .... .... = Response: Message is a query                                                                   *)
// (* .000 0... ....  = Opcode: Standard query (0)                                                                          *)
// (* .... ..0. .... .... = Truncated: Message is not truncated                                                          *)
// (* .... ...1 .... .... = Recursion desired: Do query recursively                                                    *)
// (* .... .... .0.. .... = Z: reserved (0)                                                                                                  *)
// (* .... .... ...0 .... = Non-authenticated data OK: Non-authenticated data is unacceptable *)
// (*-----------------------------------------------------------------------------------------------------------------*)
//
// (*      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5  *)
// (*    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+ *)
// (*    |                      ID                       | idx: 0/1   *)
// (*    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+ *)
// (*    |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   | idx: 2/3   *)
// (*    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+ *)
// (*    |                    QDCOUNT                    | idx: 4/5   *)
// (*    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+ *)
// (*    |                    ANCOUNT                    | idx: 6/7   *)
// (*    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+ *)
// (*    |                    NSCOUNT                    | idx: 8/9   *)
// (*    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+ *)
// (*    |                    ARCOUNT                    | idx: 10/11 *)
// (*    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+ *)
//
// (*
// Errorcodes  Description	References
// 000	No error. 		 	The request completed successfully.
// 001	Format error. 		The name server was unable to interpret the query.
// 002	Server failure.  	The name server was unable to process this query due to a problem with the name server.
// 003	Name Error. 	 	Meaningful only for responses from an authoritative name server, this code signifies that the domain name referenced in the query does not exist.
// 004	Not Implemented. 	The name server does not support the requested kind of query.
// 005	Refused. 			The name server refuses to perform the specified operation for policy reasons.
//   						 For example, a name server may not wish to provide the information to the particular requester, or a name server may not wish to perform a particular operation (e.g., zone transfer) for particular data.
// 006	YXDomain.			Name Exists when it should not.
// 007	YXRRSet. RR			Set Exists when it should not.
// 008	NXRRSet. RR			Set that should exist does not.
// 009	NotAuth.			Server Not Authoritative for zone.
// 010	NotZone,			Name not contained in zone.
// 255	Timeout				No Response
//
// example of a DNS request for google.at
// a request with 3 answer blocks
//
// Domain Name System (response)
// [Request In: 3]
// [Time: 0.027533000 seconds]
// Transaction ID: 0x5e1a
// Flags: 0x8180 (Standard query response, No error)
// 1... .... .... .... = Response: Message is a response
// .000 0... .... .... = Opcode: Standard query (0)
// .... .0.. .... .... = Authoritative: Server is not an authority for domain
// .... ..0. .... .... = Truncated: Message is not truncated
// .... ...1 .... .... = Recursion desired: Do query recursively
// .... .... 1... .... = Recursion available: Server can do recursive queries
// .... .... .0.. .... = Z: reserved (0)
// .... .... ..0. .... = Answer authenticated: Answer/authority portion was not authenticated
// .... .... .... 0000 = Reply code: No error (0)
// Questions: 1
// Answer RRs: 3
// Authority RRs: 0
// Additional RRs: 0
//
// Queries
// google.at: type A, class IN
// Name: google.at
// Type: A (Host address)
// Class: IN (0x0001)
//
// Answers
// google.at: type A, class IN, addr 72.14.221.104
// Name: google.at
// Type: A (Host address)
// Class: IN (0x0001)
// Time to live: 19 minutes, 9 seconds
// Data length: 4
// Addr: 72.14.221.104
//
// google.at: type A, class IN, addr 66.249.93.104
// Name: google.at
// Type: A (Host address)
// Class: IN (0x0001)
// Time to live: 19 minutes, 9 seconds
// Data length: 4
// Addr: 66.249.93.104
//
// google.at: type A, class IN, addr 216.239.59.104
// Name: google.at
// Type: A (Host address)
// Class: IN (0x0001)
// Time to live: 19 minutes, 9 seconds
// Data length: 4
// Addr: 216.239.59.104
// *)
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses BUFFER_TO_STRING, which is commented out.
// FUNCTION_BLOCK DNS_DYN
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// END_VAR
// VAR_INPUT
// 	ENABLE : BOOL;
// 	UPDATE : BOOL;
// 	T_UPDATE : TIME := T#1h;
// 	MODE : BYTE;
// 	HOSTNAME : STRING(30);
// 	USERNAME : STRING(20);
// 	PASSWORD : STRING(20);
// 	IP4 : DWORD;
// END_VAR
// VAR_OUTPUT
// 	BUSY : BOOL;
// 	DONE : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	DNS_CLIENT : DNS_CLIENT;
// 	HTTP_GET : HTTP_GET;
// 	BASE64_ENCODE_STR : BASE64_ENCODE_STR;
// 	URL_DATA : url;
// 	update_last : BOOL;
// 	state : INT;
// 	base64_done : BOOL;
// 	s1 : STRING(144);
// 	s2 : STRING(192);
// 	w : TON;
// END_VAR
//
// CASE state OF
//
// 00:	IF ENABLE AND ((UPDATE AND NOT update_last) OR w.Q) THEN
// 		state := 20;
// 		DONE := FALSE;
// 		BUSY := TRUE;
// 		ERROR_C := DWORD#0;
// 		ERROR_T := BYTE#0;
//
// 		IF MODE = 0 THEN (* dyndns.org *)
// 			s2 := CONCAT('http://members.dyndns.org/nic/update?hostname=',HOSTNAME);
// 		ELSE (* selfhost.de *)
// 			s2 := 'http://carol.selfhost.de/nic/update?hostname=1';
// 		END_IF;
// 		IF IP4 > 0 THEN (* use configured IP *)
// 			s1 := CONCAT('&myip=',IP4_TO_STRING(IP4));
// 			s2 := CONCAT(s2,s1);
// 		END_IF;
// 		URL_DATA:=STRING_TO_URL(STR:=s2, DEFAULT_PROTOCOL:='', DEFAULT_PATH:=''); (* URL for DNS AND HTTP-GET *)
// 		s1 := CONCAT(USERNAME,':');
// 		s1 := CONCAT(s1,PASSWORD); (* BASE64 Encode Text *)
// 	END_IF;
//
// 20:	IF DNS_CLIENT.DONE THEN (* DNS_CLIENT *)
// 			state := 40;
// 	ELSIF (DNS_CLIENT.ERROR > DWORD#00) THEN (* error handling *)
// 		ERROR_C := DNS_CLIENT.ERROR;
// 		ERROR_T := BYTE#01;
// 		BUSY := FALSE;
// 		state := 0;
// 	END_IF;
//
// 40: IF base64_done THEN (* BASE64 encode finished? *)
// 		s1 := CONCAT('Authorization: Basic ',s2);
// 		s1 := CONCAT(s1,'$0D$0AUser-Agent: x$0D$0A');
// 		URL_DATA.HEADER := s1;
// 		state := 60;
// 	END_IF;
//
// 60:	IF HTTP_GET.DONE THEN
// 		state := 0;
// 		BUSY := FALSE;
// 		DONE := TRUE;
//
// 		(* !! Codes !!
// 		Successful Updates: 		good , nchg
// 		Update Syntax Errors:		badsys , badagent
// 		Account-Related Errors:		badauth , !donator
// 		Hostname-Related Errors:	notfqdn , nohost , !yours , numhost , abuse
// 		Server Error Conditions:	dnserr , 911
// 		*)
// 		s2 := BUFFER_TO_STRING(PT:=REF(R_BUF.BUFFER),SIZE:=SIZEOF(R_BUF.BUFFER),START:=HTTP_GET.BODY_START + 3,STOP:=HTTP_GET.BODY_STOP - 7);
// 		IF FIND(s2,'good') = 0 THEN
// 			IF FIND(s2,'nochg') = 0 THEN
// 				ERROR_C := DWORD#1;
// 				ERROR_T := BYTE#03;
// 				DONE := FALSE;
// 			END_IF;
// 		END_IF;
// 	ELSIF (HTTP_GET.ERROR > DWORD#00) THEN (* error handling *)
// 		ERROR_C := HTTP_GET.ERROR;
// 		ERROR_T := BYTE#02;
// 		BUSY := FALSE;
// 		state := 0;
// 	END_IF;
//
// END_CASE;
//
// DNS_CLIENT(IP_C:=IP_C,	S_BUF:=S_BUF, R_BUF:=R_BUF,DOMAIN:=URL_DATA.DOMAIN, IP4_DNS:=0, ACTIVATE:=state=20); (* DNS_CLIENT *)
// HTTP_GET( IP_C:=IP_C, S_BUF:=S_BUF, R_BUF:=R_BUF, IP4:=DNS_CLIENT.IP4, GET:=state=60, MODE:=2, UNLOCK_BUF:=state<>60,URL_DATA:=URL_DATA); (* HTTP_GET *)
// BASE64_ENCODE_STR( RUN:=state=40, STR1:=s1, STR2:=s2, DONE=>base64_done); (* BASE64 Encode *)
//
// IF T_UPDATE > T#0ms THEN w(IN:=state=0, PT:=T_UPDATE); END_IF;
// update_last := UPDATE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK DNS_REV_CLIENT
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// END_VAR
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	IP4 : DWORD;
// 	IP4_DNS : DWORD;
// END_VAR
// VAR_OUTPUT
// 	DOMAIN : STRING;
// 	DONE : BOOL;
// 	ERROR : DWORD;
// END_VAR
// VAR
// 	IP_FIFO : IP_FIFO;
// 	ip_state : BYTE;
// 	ip_id : BYTE;
// 	state : INT;
// 	ip4_copy : DWORD;
// 	d : STRING(27);
// 	activate_last : BOOL;
// 	i : INT;
// 	dl : INT;
// 	p1 : INT;
// 	char : BYTE;
// 	char_count : INT;
// 	tid : BYTE;
// END_VAR
//
// (* Domain Name System (query) *)
//
// CASE state OF
// 00: (* wait for rising edge on GET *)
// 	IF ACTIVATE AND NOT activate_last THEN
// 		ERROR := DWORD#0;
// 		ip4_copy := ip4; (* save IP4 *)
// 		ip_state := BYTE#1;
// 		DOMAIN := '';
// 		DONE := FALSE;
//
// 		d := '';
// 		FOR i := 1 TO 4 DO
// 			d := CONCAT(d,BYTE_TO_STRING(DWORD_TO_BYTE(IP4)));
// 			d := CONCAT(d,'.');
// 			IP4 := ROR(IP4,8);
// 		END_FOR;
// 		d := CONCAT(d,'in-addr.arpa');
// 		state := 10;
// 	END_IF;
//
// 10: (* register with IP_Control *)
// 	IF ip_state = BYTE#3 THEN
// 		dl := LEN(d);
//
// 		(* initialize send buffer *)
// 		char_count := 17 + dl;
// 		FOR i := 0 TO char_count DO
// 			S_BUF.BUFFER[i] := BYTE#00;
// 		END_FOR;
//
// 		tid := UINT_TO_BYTE(BYTE_TO_UINT(tid) + UINT#1);
// 		S_BUF.BUFFER[01] := tid; (* Transaction ID lowbyte *)
// 		S_BUF.BUFFER[02] := BYTE#16#01; (* Flags: highbyte *)
// 		S_BUF.BUFFER[05] := BYTE#16#01; (* Questions: lowbyte *)
//
// 		(* process the server name from back to front *)
// 		i := dl + 12;
// 		char_count := 0;
//
// 		WHILE i > 12 DO
// 			char := CODE(STR:=d,POS:=i - 12);
//
// 			(* check for main/sub domain name separator "." *)
// 			IF char = BYTE#46 THEN
// 				(* instead of the separator "." the length of the detected main/sub domain name is entered *)
// 				S_BUF.BUFFER[i] := INT_TO_BYTE(char_count);
// 				char_count := 0;
// 			ELSE
// 				S_BUF.BUFFER[i] := char; (* enter char into buffer *)
// 				char_count := char_count + 1; (* count number of chars *)
// 			END_IF;
// 			i := i - 1;
// 		END_WHILE;
//
// 		S_BUF.BUFFER[i] := INT_TO_BYTE(char_count); (* enter the length of the leftmost (first) domain name *)
// 		i := dl + 15; (* position after the server name *)
// 		S_BUF.BUFFER[i] := BYTE#12; (* Type: PTR (Domain name pointer) lowbyte *)
// 		i := i + 2;
// 		S_BUF.BUFFER[i] := BYTE#01; (* Class: IN lowbyte *)
// 		i := i + 1;
//
// 		(* set up IP data traffic *)
// 		IP_C.C_PORT := WORD#53; (* enter port number *)
// 		IP_C.C_IP := IP4_DNS; (* enter IP *)
// 		IP_C.C_MODE := BYTE#1; (* Mode: UDP+ACTIVE+PORT+IP *)
// 		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
// 		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
// 		IP_C.R_OBSERVE  := TRUE; (* monitor data reception *)
// 		S_BUF.SIZE := INT_TO_UINT(i); (* enter send length *)
// 		R_BUF.SIZE := UINT#0; (* reset receive length *)
// 		state := 30;
// 	END_IF;
//
// 30:
// 	IF IP_C.ERROR <> DWORD#00 THEN
// 		ERROR := IP_C.ERROR;
// 	(* check whether the response telegram has the correct transaction ID and minimum length *)
// 	ELSIF (S_BUF.SIZE = UINT#0) AND (tid = R_BUF.BUFFER[1]) AND (R_BUF.SIZE >= INT_TO_UINT(i+12)) THEN
//
// 		(* check return code of the flags. 4 bits *)
// 		ERROR := BYTE_TO_DWORD(R_BUF.BUFFER[3] AND BYTE#2#0000_1111);
//
// 		IF ERROR = DWORD#00 THEN
//
// 			i := R_BUF.SIZE - 2;
// 			WHILE R_BUF.BUFFER[i] > 0 DO
// 				IF i = 0 THEN EXIT; END_IF;
// 				i := i - 1;
// 			END_WHILE;
//
// 			i := i + 2;
// 			p1 := i + 1;
// 			WHILE i < R_BUF.SIZE DO
// 				dl := R_BUF.BUFFER[i];
// 				IF dl > 0 THEN
// 					R_BUF.BUFFER[i] := BYTE#46;
// 					i := i + dl + 1;
// 				ELSE
// 					DOMAIN := BUFFER_TO_STRING(PT:=REF(R_BUF.BUFFER),SIZE:=SIZEOF(R_BUF.BUFFER),START:=p1,STOP:=i-1);
// 					EXIT;
// 				END_IF;
// 			END_WHILE;
// 			state := 0;
// 			DONE := TRUE;
// 		END_IF;
// 	END_IF;
//
// 	IF ERROR <> DWORD#0 THEN
// 		state := 0;
// 	END_IF;
//
// 	IF state = 0 THEN
// 		IP_C.R_OBSERVE := FALSE; (* monitor data reception *)
// 		ip_state := BYTE#4; (* unregister *)
// 	END_IF;
//
// END_CASE;
//
// activate_last := ACTIVATE; (* save last status for edge detection *)
//
// (* IP_FIFO *)
// IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id);
//
// (*----------------------------------------------------------------------------------------------------------------*)
// (* Flags: 0x0100 (Standard query)                                                                                            *)
// (* 0... .... .... .... = Response: Message is a query                                                                   *)
// (* .000 0... ....  = Opcode: Standard query (0)                                                                          *)
// (* .... ..0. .... .... = Truncated: Message is not truncated                                                          *)
// (* .... ...1 .... .... = Recursion desired: Do query recursively                                                    *)
// (* .... .... .0.. .... = Z: reserved (0)                                                                                                  *)
// (* .... .... ...0 .... = Non-authenticated data OK: Non-authenticated data is unacceptable *)
// (*-----------------------------------------------------------------------------------------------------------------*)
//
// (*
// Errorcodes  Description	References
// 000	No error. 		 	The request completed successfully.
// 001	Format error. 		The name server was unable to interpret the query.
// 002	Server failure.  	The name server was unable to process this query due to a problem with the name server.
// 003	Name Error. 	 	Meaningful only for responses from an authoritative name server, this code signifies that the domain name referenced in the query does not exist.
// 004	Not Implemented. 	The name server does not support the requested kind of query.
// 005	Refused. 			The name server refuses to perform the specified operation for policy reasons.
//   						 For example, a name server may not wish to provide the information to the particular requester, or a name server may not wish to perform a particular operation (e.g., zone transfer) for particular data.
// 006	YXDomain.			Name Exists when it should not.
// 007	YXRRSet. RR			Set Exists when it should not.
// 008	NXRRSet. RR			Set that should exist does not.
// 009	NotAuth.			Server Not Authoritative for zone.
// 010	NotZone,			Name not contained in zone.
// 255	Timeout				No Response
// *)
// END_FUNCTION_BLOCK

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION ELEMENT_COUNT : INT
// VAR_INPUT
// 	SEP : BYTE;
// END_VAR
// VAR_IN_OUT
// 	ELEMENT : STRING(ELEMENT_LENGTH);
// END_VAR
// VAR
// 	pt : REF_TO ARRAY[1..ELEMENT_LENGTH] OF BYTE;
// 	l: INT;
// 	pos: INT;
// END_VAR
//
// IF ELEMENT = '' THEN
// 	ELEMENT_COUNT := 0;
// ELSE
// 	PT := REF(ELEMENT);
// 	ELEMENT_COUNT := 1;
// 	l := LEN(ELEMENT);
// 	FOR pos := 1 TO l DO
// 		IF pt^[pos] = SEP THEN ELEMENT_COUNT := ELEMENT_COUNT + 1; END_IF;
// 	END_FOR;
// END_IF;
// END_FUNCTION

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION ELEMENT_GET : STRING(ELEMENT_LENGTH)
// VAR_INPUT
// 	SEP : BYTE;
// 	POS : INT;
// END_VAR
// VAR_IN_OUT
// 	ELEMENT : STRING(ELEMENT_LENGTH);
// END_VAR
// VAR
// 	i : INT := 1;
// 	o : INT := 1;
// 	pt : REF_TO ARRAY[1..ELEMENT_LENGTH] OF BYTE;
// 	po : REF_TO ARRAY[1..ELEMENT_LENGTH] OF BYTE;
// 	cnt: INT;
// 	c: BYTE;
// END_VAR
//
// (* load pointers *)
// pt := REF(ELEMENT);
// po := REF(ELEMENT_GET);
//
// (* search for the n-th occurence of a separation character and store position in i *)
// WHILE (cnt < pos) AND (i < ELEMENT_LENGTH) AND (pt^[i] > 0) DO
// 	IF pt^[i] = sep THEN cnt := cnt + 1; END_IF;
// 	i := i + 1;
// END_WHILE;
//
// (* copy the list element to output *)
// IF i < ELEMENT_LENGTH THEN c := pt^[i]; END_IF;
// WHILE (c <> sep) AND (c > 0) AND (o < ELEMENT_LENGTH) AND (i < ELEMENT_LENGTH) DO
// 	po^[o] :=pt^[i];
// 	o := o + 1;
// 	i := i + 1;
// 	c := pt^[i];
// END_WHILE;
//
// (* terminate the output string *)
// po^[o] := 0;
// END_FUNCTION

FUNCTION_BLOCK FILE_BLOCK
VAR_IN_OUT
	MODE : BYTE;
	FILENAME : STRING;
	FSD : FILE_SERVER_DATA;
	PT : NETWORK_BUFFER;
END_VAR
VAR_INPUT
	POS : UDINT;
END_VAR
VAR_OUTPUT
	ERROR : BYTE;
	DATA : BYTE;
END_VAR
VAR
	step : INT;
	i : INT;
	data_start : UDINT;
	data_stop : UDINT;

END_VAR

CASE step OF
0:	IF MODE > 0 THEN
		ERROR := 0;
		IF FSD.FILE_OPEN AND (FSD.FILE_SIZE = 0 OR POS >= FSD.FILE_SIZE) THEN
			ERROR := 255;
			MODE := 0;
		ELSIF POS < data_start OR POS > data_stop OR NOT FSD.FILE_OPEN THEN (* reading data necessary *)
			FSD.FILENAME := FILENAME;
			FSD.MODE := 1; (* open + read *)
			FSD.OFFSET := POS;
			PT.SIZE := 65535; (* read maximum number *);
			data_start := 0;
			data_stop := 0;
			step := 10;
		ELSE (* data is present in the buffer *)
			i := UDINT_TO_INT(POS - data_start); (* calculate index *)
			DATA := PT.BUFFER[i];
			MODE := 0;
			(* ---------------------- Debug-Message ----------------------------*)
			(* -----------------------------------------------------------------*)
		END_IF;
	END_IF;

10:	IF FSD.MODE = 0 THEN (* command finished *)
		IF FSD.ERROR > 0 THEN
			ERROR := FSD.ERROR; (* output error *)
			MODE := 0;
		ELSE
			data_start := POS;
			data_stop := POS + UINT_TO_UDINT(PT.SIZE) - 1;
		END_IF;
		step := 0;
	END_IF;
END_CASE;
END_FUNCTION_BLOCK

FUNCTION FILE_PATH_SPLIT : BOOL
VAR_INPUT
	FILENAME : STRING(STRING_LENGTH);
END_VAR
VAR_IN_OUT
	X : FILE_PATH_DATA;
END_VAR
VAR
	b : INT;
	c : INT;
	d : STRING(1);
	p1 : INT;
	p2 : INT;
END_VAR

X.DRIVE := '';
X.DIRECTORY := '';
X.FILENAME := '';
c := LEN(FILENAME);
IF c = 0 THEN
	FILE_PATH_SPLIT := FALSE;
	RETURN;
END_IF;

FOR b := 1 TO c DO
	d := MID(FILENAME, 1, b);
	IF d = ':' THEN
		p1 := b;
	ELSIF d = '/' OR d = '\' THEN
		p2 := b;
	END_IF;
END_FOR;

IF p1 = 2 THEN
	X.DRIVE := LEFT(FILENAME,p1);
END_IF;

IF p2 > 0 AND p2 > p1 THEN
	X.DIRECTORY := MID(FILENAME,p2 - p1,p1+1); (* len , pos *)
END_IF;

X.FILENAME := RIGHT(FILENAME,c - MAX (p1,p2));
FILE_PATH_SPLIT := TRUE;
END_FUNCTION

FUNCTION_BLOCK FILE_SERVER
VAR_IN_OUT
	FSD : FILE_SERVER_DATA;
	PT : NETWORK_BUFFER;
END_VAR
VAR
	buf_size : UDINT;
	handle : UINT; (* Filehandle *)
	read_max_length : UDINT;(* maximum data length *)
	file_position : UDINT;(* current write/read pointer *)
	used_filename : STRING;(* current file name *)
	length_read : UDINT;
	write_length : UDINT;
	length_written : UDINT;
	open_handle : UINT;
	seek_position : UDINT;
	seek_mode : INT;
	seek_data : ARRAY[0..1] OF DINT;
	tell_position : UDINT;
	open_mode : DWORD;
	last_mode: BYTE;
	command : BYTE;
	error_code : BYTE;
	error : BOOL;
	file_change : BOOL;
	step : INT;
	X : FW_AdsRdWrt;
	tmpDW : UDINT;
	para : BOOL;
	FSR : FILE_SERVER_RUNTIME;

END_VAR

CASE step OF
0:
	IF FSD.MODE = LIMIT(1,FSD.MODE,5) THEN

		(* ---------------------- Debug-Message ----------------------------*)
		(* -----------------------------------------------------------------------------*)

		(* initialization *)
		X.sNetId := '';
		X.tTimeout := T#5s;
		X.nPort := 10000; (* AMSPORT_R3_SYSSERV *)

		buf_size := SIZEOF(PT.BUFFER);

		IF FSD.MODE = 4 THEN (* File_Remove *)
			IF handle > 0 THEN
				File_change := TRUE;
				step := 32000; (* execute File_Close *)
			ELSE
				step := 90; (* execute File_Remove *)
			END_IF;

		ELSIF FSD.MODE = 5 THEN (* File_Close *)
			step := 32000; (* execute File_Close *)

		ELSIF handle = 0 THEN (* MODE 1,2,3 *)
			step := 100; (* file must be opened first *)

		ELSE (* file is open *)
			IF used_filename = FSD.FILENAME AND FSD.MODE = last_mode THEN (* no mode change? *)
				step := 200; (* continue without open procedure *)
			ELSE
				step := 32000;
				File_change := TRUE;
			END_IF;
		END_IF;
	END_IF;

(* ----------------------------------- *)
090: (* File_Remove *)
	command := 7; (* File Remove *)
	step := 95;

095: (* File_Remove finished *)
	IF command = 0 THEN
		step := 30000; (* exit *)
	END_IF;
(* ----------------------------------- *)

100: (* File_Open *)
	IF FSD.MODE = 3 THEN (* Open_New + File Write *)
		open_mode := 16#0001_001A; (* PATH_GENERIC + WRITEPLUS + BINARY *)
	ELSE (* Mode 1 Open_Exist + File _Read / Mode 2 Open_Exist + File_Write *)
		open_mode := 16#0001_0019; (* PATH_GENERIC + READPLUS + BINARY *)
	END_IF;
	command := 1; (* FileOpen *)
	step := 105;

105: (* file open finished *)
	IF command = 0 THEN
		IF error THEN
			step := 30000; (* exit *)
		ELSE
			(* enter file data *)
			handle := open_handle;
			used_filename := FSD.FILENAME;
			file_position := UDINT#4294967295; (* 16#FFFF_FFFF -> maximum = unknown *)
			last_mode := FSD.MODE;
			step := 110;
		END_IF;
	END_IF;

110: (* File Seek *)
	command := 5; (* File-Seek *)
	seek_position := 0;
	seek_mode := 2; (* SEEK_END - Seek from the end of file *)
	step := 115;

115: (* file seek finished *)
	IF command = 0 THEN
		IF error THEN
			step := 30000; (* exit *)
		ELSE
			command := 6; (* start file tell *)
			step := 120;
		END_IF;
	END_IF;

120: (* file tell finished *)
	IF command = 0 THEN
		IF error THEN
			step := 30000; (* exit *)
		ELSE
			FSD.FILE_SIZE := tell_position; (* current file size *)
			FSD.FILE_OPEN := TRUE;
			step := 200;
		END_IF;
	END_IF;

200: (* File Seek *)
	IF FSD.OFFSET = UDINT#4294967294 THEN (* 16#FFFF_FFFE -> set pointer to end of file *)
		command := 5; (* File-Seek *)
		seek_position := FSD.FILE_SIZE;
		seek_mode := 0; (* Seek from beginning of file *)
		step := 210;
	ELSIF FSD.OFFSET > FSD.FILE_SIZE THEN (* check for valid offset setting *)
		error_code := 255;
		step := 30000; (* exit *)
	ELSIF FSD.OFFSET <> file_position THEN
		command := 5; (* File-Seek *)
		seek_position := FSD.OFFSET;
		seek_mode := 0; (* Seek from beginning of file *)
		step := 210;
	ELSE
		step := 300;
	END_IF;

210: (* File Seek *)
	IF command = 0 THEN
		file_position := seek_position;
		IF error THEN
			step := 30000; (* exit *)
		ELSE
			step := 300;
		END_IF;
	END_IF;

300: (* Read or Write ? *)
	IF FSD.MODE = 1 THEN
		step := 400; (* Read *)
	ELSE
		step := 500; (* Write *)
	END_IF;

400: (* File Read *)
	IF UINT_TO_UDINT(PT.SIZE) > buf_size THEN (* limit read length to buffer size *)
		read_max_length := buf_size;
	ELSE
		read_max_length := UINT_TO_UDINT(PT.SIZE);
	END_IF;

	IF read_max_length > buf_size THEN
		read_max_length := buf_size;  (* limit read length to buffer size *)
	END_IF;

	(* adjust read length to file length *)
	IF file_position + read_max_length > FSD.FILE_SIZE THEN
		read_max_length := FSD.FILE_SIZE - file_position;
	END_IF;

	IF read_max_length > 0 THEN
		command := 3; (* File_Read *)
		step := 410;
	ELSE
		PT.SIZE := UINT#0;
		step := 30000; (* exit *)
	END_IF;

410: (* File Read *)
	IF command = 0 THEN
		PT.SIZE := UDINT_TO_UINT(length_read);
		IF NOT Error THEN
			file_position := file_position + length_read; (* correct current File_Position *)
			FSD.OFFSET := file_position;
		END_IF;
		step := 30000; (* exit *)
	END_IF;

500: (* File Write *)
	IF UINT_TO_UDINT(PT.SIZE) > buf_size THEN (* limit write length to buffer size *)
		write_length := buf_size;
	ELSE
		write_length := UINT_TO_UDINT(PT.SIZE);
	END_IF;

	IF write_length > 0 THEN
		command := 4; (* File_Write *)
		step := 510;
	ELSE
		step := 30000; (* exit *)
	END_IF;

510:
	IF command = 0 THEN
		IF NOT Error THEN
			file_position := file_position + length_written; (* correct current File_Position *)
			FSD.OFFSET := file_position;
			IF file_position > FSD.FILE_SIZE THEN (* file has grown *)
				FSD.FILE_SIZE := file_position;
			END_IF;
		END_IF;
		step := 30000; (* exit *)
	END_IF;

 (* --------------------------------------------- *)

30000: (* command done *)
	FSD.MODE := 0;
	FSD.ERROR := error_code;
	step := 0;

	(* ---------------- Debug-Message -----------------------*)
	(* --------------------------------------------- *)

32000:
	IF handle > 0 THEN
		command := 2; (* execute File_Close *)
	END_IF;
	step := 32100;

32100:
	IF command = 0 THEN (* file close finished *)
		(* clear file data *)
		FSD.FILE_OPEN := FALSE;
		file_position := 0;
		FSD.FILE_SIZE := 0;
		FSD.OFFSET := 0;
		used_filename := '';
		handle := 0;

		(* File_Name changed on the fly? *)
		IF file_change THEN
			file_change := FALSE;
			step := 0;
		ELSE
			step := 30000;
		END_IF;
	END_IF;

END_CASE;

X(bExecute := FALSE); (* cyclic FW_AdsRdWrt call *)

FSR(MODE:=BYTE#0,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command start *)

(* ------------ File Open ------------- *)
IF command = 1 THEN
	IF NOT para THEN
		para := TRUE;
		X(	nIdxGrp := 120, (* SYSTEMSERVICE_FOPEN *)
			nidxOffs := open_mode,
			cbWriteLen := SIZEOF(FSD.FILENAME),
			pWriteBuff := REF(FSD.FILENAME),
			cbReadLen := SIZEOF(tmpDW),
			pReadBuff := REF(tmpDW),
			bExecute := TRUE
			);
	ELSIF X.bBusy = FALSE THEN
		para := FALSE;
		error := X.bError;
		error_code := SEL(error,0,UDINT_TO_BYTE(X.nErrId + 16#10));

		IF X.cbRead >= SIZEOF(tmpDW) THEN
			open_handle := UDINT_TO_UINT(tmpDW);
		ELSE
			open_handle := 0;
		END_IF;
		command := 0;
		FSR(MODE:=BYTE#1,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command stop *)
		(* ---------------- Debug-Message -----------------------*)
		(* -----------------------------------------------------------------*)
	END_IF;
END_IF;

(* ----------- File Close ------------- *)
IF command = 2 THEN
	IF NOT para THEN
		para := TRUE;
		X(	nIdxGrp := 121, (* SYSTEMSERVICE_FCLOSE *)
			nidxOffs := handle,
			cbWriteLen := 0,
			pWriteBuff := 0,
			cbReadLen := 0,
			pReadBuff := 0,
			bExecute := TRUE
			);
	ELSIF X.bBusy = FALSE THEN
		para := FALSE;
		error := X.bError;
		error_code := SEL(error,0,UDINT_TO_BYTE(X.nErrId + 16#30));
		command := 0;
		FSR(MODE:=BYTE#1,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command stop *)
		(* ---------------- Debug-Message -----------------------*)
		(* -----------------------------------------------------------------*)
	END_IF;
END_IF;

(* ------------ File Read ------------- *)
IF command = 3 THEN
	IF NOT para THEN
		para := TRUE;
		X(	nIdxGrp := 122, (* SYSTEMSERVICE_FREAD *)
			nidxOffs := handle,
			cbWriteLen := 0,
			pWriteBuff := 0,
			cbReadLen := read_max_length,
			pReadBuff := REF(PT.BUFFER),
			bExecute := TRUE
			);
	ELSIF X.bBusy = FALSE THEN
		length_read := X.cbRead;
		para := FALSE;
		error := X.bError;
		error_code := SEL(error,0,UDINT_TO_BYTE(X.nErrId + 16#40));
		command := 0;
		FSR(MODE:=BYTE#1,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command stop *)
		(* ---------------- Debug-Message -----------------------*)
		(* -----------------------------------------------------------------*)
	END_IF;
END_IF;

(* ------------ File Write ------------- *)
IF command = 4 THEN
	IF NOT para THEN
		para := TRUE;
		X(	nIdxGrp := 123, (* SYSTEMSERVICE_FWRITE *)
			nidxOffs := handle,
			cbWriteLen := write_length,
			pWriteBuff := REF(PT.BUFFER),
			cbReadLen := SIZEOF(length_written),
			pReadBuff := REF(length_written),
			bExecute := TRUE
			);
	ELSIF X.bBusy = FALSE THEN
		para := FALSE;
		error := X.bError;
		error_code := SEL(error,0,UDINT_TO_BYTE(X.nErrId + 16#50));
		command := 0;
		FSR(MODE:=BYTE#1,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command stop *)
		(* ---------------- Debug-Message -----------------------*)
		(* -----------------------------------------------------------------*)
	END_IF;
END_IF;

(* ------------ File Seek ------------- *)
IF command = 5 THEN
	IF NOT para THEN
		para := TRUE;
		seek_data[0] := UDINT_TO_DINT(seek_position);
		seek_data[1] := INT_TO_DINT(seek_mode);

		X(	nIdxGrp := 124, (* SYSTEMSERVICE_FSEEK *)
			nidxOffs := handle,
			cbWriteLen := SIZEOF(seek_data),
			pWriteBuff := REF(seek_data),
			cbReadLen := 0,
			pReadBuff := 0,
			bExecute := TRUE
			);
	ELSIF X.bBusy = FALSE THEN
		para := FALSE;
		error := X.bError;
		error_code := SEL(error,0,UDINT_TO_BYTE(X.nErrId + 16#60));
		command := 0;
		FSR(MODE:=BYTE#1,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command stop *)
		(* ---------------- Debug-Message -----------------------*)
		(* -----------------------------------------------------------------*)
	END_IF;
END_IF;

(* ------------ File Tell ------------- *)
IF command = 6 THEN
	IF NOT para THEN
		para := TRUE;
		X(	nIdxGrp := 125, (* SYSTEMSERVICE_FTELL *)
			nidxOffs := handle,
			cbWriteLen := 0,
			pWriteBuff := 0,
			cbReadLen := SIZEOF(tell_position),
			pReadBuff := REF(tell_position),
			bExecute := TRUE
			);
	ELSIF X.bBusy = FALSE THEN
		para := FALSE;
		error := X.bError;
		error_code := SEL(error,0,UDINT_TO_BYTE(X.nErrId + 16#70));
		command := 0;
		FSR(MODE:=BYTE#1,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command stop *)
		(* ---------------- Debug-Message -----------------------*)
		(* -----------------------------------------------------------------*)
	END_IF;
END_IF;

(* ------------ File Remove ------------- *)
IF command = 7 THEN
	IF NOT para THEN
		para := TRUE;
		X(	nIdxGrp := 131, (* SYSTEMSERVICE_FDELETE *)
			nidxOffs := 16#0001_0000, (* PATH_GENERIC *)
			cbWriteLen := SIZEOF(FSD.FILENAME),
			pWriteBuff := REF(FSD.FILENAME),
			cbReadLen := 0,
			pReadBuff := 0,
			bExecute := TRUE
			);
	ELSIF X.bBusy = FALSE THEN
		para := FALSE;
		error := X.bError;
		error_code := SEL(error,0,UDINT_TO_BYTE(X.nErrId + 16#80));
		command := 0;
		FSR(MODE:=BYTE#1,COMMAND:=command,FSD:=FSD); (* monitor file server runtime / command stop *)
		(* ---------------- Debug-Message -----------------------*)
		(* -----------------------------------------------------------------*)
	END_IF;
END_IF;

(* ---------------- Debug-Message -----------------------*)
(* -----------------------------------------------------------------*)
END_FUNCTION_BLOCK

FUNCTION_BLOCK FILE_SERVER_RUNTIME
VAR_INPUT
	MODE : BYTE;
	COMMAND : BYTE;
END_VAR
VAR_IN_OUT
	FSD : FILE_SERVER_DATA;
END_VAR
VAR
	last_command : BYTE;
	init : BOOL;
	start_tx : DWORD;
	stop_tx : DWORD;
	runtime : DWORD;
END_VAR

IF init = FALSE THEN
	init := TRUE;
	FSD.RUNTIME.TIME_FILE_OPEN_MIN := UDINT#99999999;
	FSD.RUNTIME.TIME_FILE_CLOSE_MIN := UDINT#99999999;
	FSD.RUNTIME.TIME_FILE_READ_MIN := UDINT#99999999;
	FSD.RUNTIME.TIME_FILE_WRITE_MIN := UDINT#99999999;
	FSD.RUNTIME.TIME_FILE_TELL_MIN := UDINT#99999999;
	FSD.RUNTIME.TIME_FILE_SEEK_MIN := UDINT#99999999;
	FSD.RUNTIME.TIME_FILE_DELETE_MIN := UDINT#99999999;
END_IF;

(* wait for File_Server command *)
IF MODE = BYTE#0 THEN

	(* new command detected *)
	IF last_command = BYTE#0 AND COMMAND <> BYTE#0 THEN

		(* read system time *)
		start_tx := T_PLC_MS();

	END_IF;
END_IF;

IF MODE = BYTE#1 THEN

	(* read system time *)
	stop_tx:=  T_PLC_MS();

	runtime:= (stop_tx - start_tx);

	CASE BYTE_TO_INT(COMMAND) OF

	1:	FSD.RUNTIME.TIME_FILE_OPEN_CUR := runtime; 			(* current runtime *)
		IF FSD.RUNTIME.TIME_FILE_OPEN_MAX < runtime THEN
			FSD.RUNTIME.TIME_FILE_OPEN_MAX := runtime;		(* new max runtime *)
		END_IF;
		IF FSD.RUNTIME.TIME_FILE_OPEN_MIN > runtime THEN
			FSD.RUNTIME.TIME_FILE_OPEN_MIN := runtime;		(* new min runtime *)
		END_IF;

	2:	FSD.RUNTIME.TIME_FILE_CLOSE_CUR := runtime; 			(* current runtime *)
		IF FSD.RUNTIME.TIME_FILE_CLOSE_MAX < runtime THEN
			FSD.RUNTIME.TIME_FILE_CLOSE_MAX := runtime;		(* new max runtime *)
		END_IF;
		IF FSD.RUNTIME.TIME_FILE_CLOSE_MIN > runtime THEN
			FSD.RUNTIME.TIME_FILE_CLOSE_MIN := runtime;		(* new min runtime *)
		END_IF;

	3:	FSD.RUNTIME.TIME_FILE_READ_CUR := runtime; 			(* current runtime *)
		IF FSD.RUNTIME.TIME_FILE_READ_MAX < runtime THEN
			FSD.RUNTIME.TIME_FILE_READ_MAX := runtime;		(* new max runtime *)
		END_IF;
		IF FSD.RUNTIME.TIME_FILE_READ_MIN > runtime THEN
			FSD.RUNTIME.TIME_FILE_READ_MIN := runtime;		(* new min runtime *)
		END_IF;

	4:	FSD.RUNTIME.TIME_FILE_WRITE_CUR := runtime; 			(* current runtime *)
		IF FSD.RUNTIME.TIME_FILE_WRITE_MAX < runtime THEN
			FSD.RUNTIME.TIME_FILE_WRITE_MAX := runtime;		(* new max runtime *)
		END_IF;
		IF FSD.RUNTIME.TIME_FILE_WRITE_MIN > runtime THEN
			FSD.RUNTIME.TIME_FILE_WRITE_MIN := runtime;		(* new min runtime *)
		END_IF;

	5:	FSD.RUNTIME.TIME_FILE_SEEK_CUR := runtime; 			(* current runtime *)
		IF FSD.RUNTIME.TIME_FILE_SEEK_MAX < runtime THEN
			FSD.RUNTIME.TIME_FILE_SEEK_MAX := runtime;		(* new max runtime *)
		END_IF;
		IF FSD.RUNTIME.TIME_FILE_SEEK_MIN > runtime THEN
			FSD.RUNTIME.TIME_FILE_SEEK_MIN := runtime;		(* new min runtime *)
		END_IF;

	6:	FSD.RUNTIME.TIME_FILE_TELL_CUR := runtime; 			(* current runtime *)
		IF FSD.RUNTIME.TIME_FILE_TELL_MAX < runtime THEN
			FSD.RUNTIME.TIME_FILE_TELL_MAX := runtime;		(* new max runtime *)
		END_IF;
		IF FSD.RUNTIME.TIME_FILE_TELL_MIN > runtime THEN
			FSD.RUNTIME.TIME_FILE_TELL_MIN := runtime;		(* new min runtime *)
		END_IF;

	7:	FSD.RUNTIME.TIME_FILE_DELETE_CUR := runtime; 			(* current runtime *)
		IF FSD.RUNTIME.TIME_FILE_DELETE_MAX < runtime THEN
			FSD.RUNTIME.TIME_FILE_DELETE_MAX := runtime;		(* new max runtime *)
		END_IF;
		IF FSD.RUNTIME.TIME_FILE_DELETE_MIN > runtime THEN
			FSD.RUNTIME.TIME_FILE_DELETE_MIN := runtime;		(* new min runtime *)
		END_IF;

	END_CASE;

END_IF;

last_command := COMMAND;
END_FUNCTION_BLOCK

// Commented out for beedance: it has a comment between UNTIL and END_REPEAT, which beedance does not parse yet.
// FUNCTION_BLOCK FTP_CLIENT
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	FILENAME : STRING;
// 	FTP_URL : STRING(STRING_LENGTH);
// 	FTP_DOWNLOAD : BOOL;
// 	FTP_ACTIV : BOOL;
// 	FILE_DELETE : BOOL;
// 	TIMEOUT : TIME;
// 	DNS_IP4 : DWORD;
// 	PLC_IP4 : DWORD;
// END_VAR
// VAR_OUTPUT
// 	DONE : BOOL;
// 	BUSY : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	IP_C1 : IP_C; (* IP_Control management structure *)
// 	IP_C2 : IP_C; (* IP_Control management structure *)
// 	S_BUF1 : NETWORK_BUFFER;
// 	S_BUF2 : NETWORK_BUFFER;
// 	R_BUF1 : NETWORK_BUFFER;
// 	R_BUF2 : NETWORK_BUFFER;
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	IPC1 : IP_CONTROL;
// 	IPC2 : IP_CONTROL;
// 	FS : FILE_SERVER;
// 	FSD : FILE_SERVER_DATA;
//
// 	activate_last : BOOL;
// 	step : INT;
// 	snd_step : INT;
// 	snd_text :	STRING(STRING_LENGTH);
// 	rcv_text :	STRING(STRING_LENGTH);
// 	rcv_state : INT;
// 	TON_WAIT : TON;
// 	last_step : INT;
// 	next_step : INT;
// 	idx1 : INT;
// 	idx2 : INT;
// 	str1 : STRING(STRING_LENGTH);
// 	str2 : STRING(STRING_LENGTH);
// 	ftp_path : STRING;
// 	byar : ARRAY [0..5] OF BYTE;
// 	dw_tmp : DWORD;
// 	w_tmp : WORD;
// 	fpd : FILE_PATH_DATA;
// 	timeout2 : TIME;
// 	timeout3 : TIME;
// 	ftp_file_size : UDINT;
// 	ip_c1_red_disable : BOOL; (* IP_CONTROL 1 - suppress remote disconnect error *)
// 	ip_c2_red_disable : BOOL; (* IP_CONTROL 1 - suppress remote disconnect error *)
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := FALSE;*)
//
// END_VAR
// CASE step OF
//
// 00:	IF ACTIVATE AND NOT activate_last THEN
// 		ftp_file_size := UDINT#0;
// 		TIMEOUT := MAX(T#10s,TIMEOUT);
// 		timeout2 := TIMEOUT + T#1s;
// 		timeout3 := TIMEOUT / 2;
// 		DONE := FALSE;
// 		ip_c1_red_disable := FALSE;
// 		ip_c2_red_disable := FALSE;
// 		ERROR_C := DWORD#0;
// 		ERROR_T := BYTE#0;
// 		FSD.ERROR := BYTE#0;
// 		URL_DATA := STRING_TO_URL(STR:=FTP_URL,DEFAULT_PROTOCOL:='',DEFAULT_PATH:='/');
// 		IF LEN(URL_DATA.USER) = 0 THEN
// 			URL_DATA.USER := 'Anonymous';
// 			URL_DATA.PASSWORD := 'User@';
// 		END_IF;
//
// 		FILE_PATH_SPLIT(FILENAME:=FILENAME,X:=fpd); (* split file path *)
// 		str1 := fpd.FILENAME; (* buffer temporarily *)
// 		FILE_PATH_SPLIT(FILENAME:=URL_DATA.PATH,X:=fpd); (* split FTP path *)
// 		IF LEN(fpd.FILENAME) = 0 THEN (* if the FTP path contains no file name, use the name of the local file *)
// 			ftp_path := CONCAT(URL_DATA.PATH,str1);
// 		ELSE
// 			ftp_path := URL_DATA.PATH;
// 		END_IF;
// 		step := 10;
// 	END_IF;
//
// 10:	IF DNS_CLIENT.DONE THEN
// 		step := 20;
// 	ELSIF DNS_CLIENT.ERROR <> DWORD#0 THEN
// 		ERROR_C := DNS_CLIENT.ERROR;
// 		ERROR_T := BYTE#1;
// 		step := 980;
// 	END_IF;
//
// 20:	(* set up FTP control channel *)
// 	IP_C1.C_PORT := WORD#21; (* enter port number *)
//   	IP_C1.C_IP := DNS_CLIENT.IP4; (* enter IP *)
//   	IP_C1.C_MODE := BYTE#0; (* Mode: TCP+ACTIV+PORT *)
//   	IP_C1.TIME_RESET := TRUE; (* reset timeout monitoring *)
//   	IP_C1.C_ENABLE := TRUE; (* enable connection setup *)
//   	IP_C1.R_OBSERVE := TRUE; (* receive monitoring *)
// 	R_BUF1.SIZE := UINT#0;
//
//   	IP_C2.C_ENABLE   := FALSE; (* block connection setup *)
// 	step := 30;
//
// 30:	IF rcv_state = 220 THEN (* 220 = Service ready for new user *)
// 		snd_text := CONCAT('USER ',URL_DATA.USER);
// 		next_step := 40;
// 	END_IF;
//
// 40:	IF rcv_state = 331 THEN (* 331 = User name okay, need password *)
// 		snd_text := CONCAT('PASS ',URL_DATA.PASSWORD);
// 		next_step := 50;
// 	ELSIF rcv_state = 230 THEN (* 230 = User logged in, proceed *)
// 		next_step := 50;
// 	ELSIF rcv_state = 220 THEN (* 220 = Service ready for new user *)
// 		next_step := 20; (* WORK connection should be closed/opened *)
// 	END_IF;
//
// 50:	IF rcv_state = 230 THEN (* 230 = User logged in, proceed *)
// 		snd_text := 'TYPE I';  (* activate binary mode *)
// 		next_step := 54;
// 	END_IF;
//
// 54:	IF rcv_state = 200 THEN (* 200 = Command okay *)
// 		IF FTP_DOWNLOAD THEN (* receive file from FTP *)
// 			snd_text := CONCAT('SIZE ', ftp_path); (* SIZE /PATH/FILENAME *)
// 			next_step := 56;
// 		ELSE
// 			step := 60;
// 		END_IF;
// 	END_IF;
//
// 56:	IF rcv_state = 213 THEN (* 213 Filesize *)
// 		str1 := ELEMENT_GET(SEP:=BYTE#32,POS:=1,ELEMENT:=rcv_text);
// 		ftp_file_size := STRING_TO_UDINT(str1); (* determine file size in bytes *)
// 		step := 60;
// 	END_IF;
//
// 60:	IF FTP_ACTIV THEN (* FTP server active *)
// 		IF URL_DATA.PORT = WORD#0 THEN
// 			URL_DATA.PORT := WORD#20; (* default data channel *)
// 		END_IF;
// 		str2 := '';
// 		dw_tmp := PLC_IP4;
// 		FOR idx1 := 0 TO 5 DO
// 			IF idx1 = 4 THEN
// 				dw_tmp := ROR(WORD_TO_DWORD(URL_DATA.PORT),8);
// 			ELSE
// 				dw_tmp := ROL(dw_tmp,8);
// 			END_IF;
// 			str2 := CONCAT(str2,BYTE_TO_STRING(DWORD_TO_BYTE(dw_tmp)));
// 			IF idx1 < 5 THEN
// 				str2 := CONCAT(str2,',');
// 			END_IF;
// 		END_FOR;
// 		dw_tmp := DNS_CLIENT.IP4; (* enter IP *)
// 		w_tmp := URL_DATA.PORT; (* port number of the data channel *)
// 		snd_text := CONCAT('PORT ',str2); (* switch FTP server to active mode *)
// 		next_step := 80;
// 	ELSE
// 		snd_text := 'PASV';  (* switch FTP server to passive mode *)
// 		next_step := 70;
// 	END_IF;
// 	step := 999; (* wait for next_step *)
//
// 70:	IF rcv_state = 227 THEN (* 227 = Entered Passive Mode (213,202,225,43,212,186) *)
// 		idx1 := FIND(rcv_text , '(' );
// 		idx2 := FIND(rcv_text , ')' );
// 		IF idx1 > 0 AND idx2 > idx1+1 THEN
// 			str1 := MID(rcv_text, idx2 - idx1 - 1, idx1+1);
// 			FOR idx1 := 0 TO 5 DO
// 				str2 := ELEMENT_GET(SEP:=BYTE#44,POS:=idx1,ELEMENT:=str1);
// 				byar[idx1] := STRING_TO_BYTE(str2);
// 			END_FOR;
// 			dw_tmp := DWORD_OF_BYTE(byar[0],byar[1],byar[2],byar[3]); (* IP address of the data channel *)
// 			w_tmp := WORD_OF_BYTE(byar[4],byar[5]); (* port number of the data channel *)
// 		END_IF;
// 		step := 100;
// 	END_IF;
//
// 80:	IF rcv_state = 200 THEN (* 200 = PORT command successful *)
// 		step := 100;
// 	END_IF;
//
// 100:IP_C2.C_MODE := SEL(FTP_ACTIV,BYTE#0,BYTE#2); (* active = mode 2: TCP+PASSIVE+PORT / passive = mode 0: TCP+ACTIVE+PORT *)
// 	IP_C2.C_IP := dw_tmp; (* IP address of the data channel *)
// 	IP_C2.C_PORT := w_tmp; (* port number of the data channel *)
//   	IP_C2.TIME_RESET := TRUE; (* reset timeout monitoring *)
//   	IP_C2.C_ENABLE := TRUE; (* enable connection setup *)
// 	IP_C2.R_OBSERVE := FALSE; (* receive monitoring *)
// 	R_BUF2.SIZE := UINT#0;
// 	step := 120;
//
// 120:IF IP_C2.C_STATE > BYTE#127 THEN (* FTP data channel established *)
// 		IF FTP_DOWNLOAD THEN (* receive file from FTP *)
// 			IP_C2.MAILBOX[3] := BYTE#1; (* block IP receive *)
// 			step := 400;
// 		ELSE
// 			snd_text := 'STOR ';
// 			next_step := 140;
// 			snd_text := CONCAT(snd_text, ftp_path); (* STORE /PATH/FILENAME *)
// 			step := 999; (* wait for next_step *)
// 		END_IF;
// 	END_IF;
//
// (* ============== send file to FTP ================ *)
//
// 140:IF rcv_state = 125 OR rcv_state = 150 THEN (* 125 = Data connection already open / 150 = File status okay; about to open data connection *)
// 		step := 160;
// 	END_IF;
//
// 160:FSD.FILENAME := FILENAME; (* file name *)
// 	FSD.MODE := BYTE#1; (* open existing file for reading *)
// 	FSD.OFFSET := UDINT#0; (* start at beginning of file *)
// 	S_BUF2.SIZE := UINT#65535; (* read maximum length *)
// 	IP_C2.MAILBOX[2] := BYTE#1; (* block IP send *)
// 	step := 200;
//
// 200:IF FSD.MODE = BYTE#0 AND FSD.ERROR = BYTE#0 THEN
// 		IP_C2.MAILBOX[2] := BYTE#0; (* allow IP send *)
// 		step := 210;
// 	END_IF;
//
// 210:IF S_BUF2.SIZE = UINT#0 THEN (* wait until all data is sent *)
// 		IF FSD.FILE_SIZE - FSD.OFFSET > UDINT#0 THEN
// 			FSD.MODE := BYTE#1; (* open existing file for reading *)
// 			S_BUF2.SIZE := UINT#65535; (* read maximum length *)
// 			IP_C2.MAILBOX[2] := BYTE#1; (* block IP send *)
// 			step := 200;
// 		ELSE
// 			IP_C2.C_ENABLE := FALSE; (* close FTP data channel *)
// 			step := 300;
// 		END_IF;
// 	END_IF;
//
// 300:IF rcv_state = 226 THEN (* 226 Closing data connection. Requested file action successful  *)
// 		snd_text := 'QUIT';
// 		ip_c1_red_disable := TRUE;
// 		next_step := 320;
// 	END_IF;
//
// 320:IF rcv_state = 221 OR TON_WAIT.ET > timeout3  THEN (* 221 Service closing control connection *)
// 		IP_C1.C_ENABLE := FALSE; (* close FTP control channel *)
// 		IF FILE_DELETE THEN (* delete file after transfer *)
// 			FSD.MODE := BYTE#4; (* File Remove *)
// 			step := 340;
// 		ELSE
// 			step := 900;
// 		END_IF;
// 	END_IF;
//
// 340:IF FSD.MODE = BYTE#0 AND FSD.ERROR = BYTE#0 THEN
// 		step := 900;
// 	END_IF;
//
// (*  ============== receive file from FTP  ============== *)
//
// 400:snd_text := 'RETR ';
// 	next_step := 410;
// 	snd_text := CONCAT(snd_text, ftp_path); (* RETR /PATH/FILENAME *)
// 	step := 999; (* wait for next_step *)
//
// 410:IF rcv_state = 125 OR rcv_state = 150 THEN (* 125 = Data connection already open / 150 = File status okay; about to open data connection *)
// 		IP_C2.R_OBSERVE := TRUE; (* receive monitoring *)
// 		step := 420;
// 	END_IF;
//
// 420:FSD.FILENAME := FILENAME; (* file name *)
// 	FSD.MODE := BYTE#3; (* create new file *)
// 	FSD.OFFSET := UDINT#0; (* start at beginning of file *)
// 	step := 460;
//
// 440:IF R_BUF2.SIZE > UINT#0 THEN (* wait until data has been received *)
// 		IP_C2.MAILBOX[3] := BYTE#1; (* block IP receive until data has been written to file *)
// 		FSD.MODE := BYTE#3; (* write data to file *)
// 		step := 460;
// 	ELSIF FSD.FILE_SIZE = ftp_file_size THEN (* connection closed?, or 226 = Closing data connection. Requested file action successful *)
// 		IP_C2.C_ENABLE := FALSE; (* close FTP data channel *)
// 		ip_c2_red_disable := TRUE;
// 		step := 700;
// 	END_IF;
//
// 460:IF FSD.MODE = BYTE#0 AND FSD.ERROR = BYTE#0 THEN (* wait until data has been written *)
// 		R_BUF2.SIZE := UINT#0; (* delete data *)
// 		IP_C2.MAILBOX[3] := BYTE#0; (* allow IP receive again *)
// 		step := 440;
// 	END_IF;
//
// (* ================================================== *)
//
// 700:IF rcv_state = 226 OR  TON_WAIT.ET > timeout3 THEN (* 226 Closing control connection. Requested file action successful  *)
// 		IF FTP_DOWNLOAD AND FILE_DELETE THEN (* FTP download and delete file? *)
// 			snd_text := CONCAT('DELE ', ftp_path); (* DELE /PATH/FILENAME *)
// 			next_step := 720;
// 		ELSE
// 			step := 740;
// 		END_IF;
// 	END_IF;
//
// 720:IF rcv_state = 250 THEN (* 250 Requested file action okay, completed *)
// 		step := 740;
// 	END_IF;
//
// 740:snd_text := 'QUIT';
// 	ip_c1_red_disable := TRUE;
// 	next_step := 760;
// 	step := 999; (* wait for next_step *)
//
// 760:IF rcv_state = 221 OR TON_WAIT.ET > timeout3 THEN (* 221 Service closing control connection  *)
// 		IP_C1.C_ENABLE := FALSE; (* close FTP control channel *)
// 		IF NOT FTP_DOWNLOAD AND FILE_DELETE THEN (* FTP upload and delete file after transfer *)
// 			FSD.MODE := BYTE#4; (* File Remove *)
// 			step := 780;
// 		ELSE
// 			step := 900;
// 		END_IF;
// 	END_IF;
//
// 780:IF FSD.MODE = BYTE#0 AND FSD.ERROR = BYTE#0 THEN
// 		step := 900;
// 	END_IF;
//
// (* ================================================== *)
//
// 900:DONE := TRUE;
// 	step := 980;
//
// 980:IP_C1.C_ENABLE := FALSE; (* close FTP control channel *)
// 	IP_C2.C_ENABLE := FALSE; (* close FTP data channel *)
// 	IP_C2.MAILBOX[2] := BYTE#0; (* allow IP send again *)
// 	IP_C2.MAILBOX[3] := BYTE#0; (* allow IP receive again *)
// 	S_BUF1.SIZE := UINT#0;
// 	S_BUF2.SIZE := UINT#0;
// 	R_BUF1.SIZE := UINT#0;
// 	R_BUF2.SIZE := UINT#0;
// 	ftp_file_size := UDINT#0;
// 	snd_step := 0;
// 	next_step := 0;
// 	FSD.MODE := BYTE#5; (* close file *)
// 	step := 990;
//
// 990:IF IP_C1.C_STATE = BYTE#0 AND IP_C2.C_STATE = BYTE#0 AND FSD.FILE_OPEN = FALSE THEN
// 		step := 0;
// 	END_IF;
//
// 999: ; (* dummy step, wait until a new step becomes active through "next_step" *)
//
// END_CASE;
//
// (* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// IF step >= 30 THEN
// 	IF R_BUF1.SIZE > UINT#0 THEN (* data received? *)
// 		idx1 := UINT_TO_INT(R_BUF1.SIZE);
// 		REPEAT (* detect line termination and determine real end of TEXT *)
// 			idx1 := idx1 - 1;
// 		UNTIL idx1 = 0 OR R_BUF1.BUFFER[idx1] >= BYTE#32 (* TRUE = EXIT *)
// 		END_REPEAT;
// 		idx1 := idx1 + 1; (* convert index to SIZE *)
// 		rcv_text := BUFFER_TO_STRING(PT:=REF(R_BUF1.BUFFER),SIZE:=idx1,START:=0,STOP:=STRING_LENGTH-1);
// 		IF LEN(rcv_text) >= 3 THEN (* check and evaluate status code *)
// 			str1 := LEFT(rcv_text,3);
// 			IF IS_NUM(STR:=str1) THEN (* contains only digits? *)
// 				rcv_state := STRING_TO_INT(str1);
// 			END_IF;
// 			IP_C1.R_OBSERVE := FALSE;
// 		END_IF;
// 		R_BUF1.SIZE := UINT#0; (* delete data *)
//
// 	        (*####################  DEBUG-MESSAGE  ###################################*)
// 	        (*IF _debug_enable THEN
// 	          LOG_CL.NEW_MSG := 'RCV: >~1<';
// 	          LOG_CL.PRINTF[1] := rcv_text;
// 	          LOG_MSG();
// 	        END_IF;*)
// 	        (*########################################################################*)
//
// 	END_IF;
// END_IF;
//
// (* automatic data sending *)
//
// CASE snd_step OF
//
// 00:	IF LEN(snd_text) > 0 THEN
// 	        (*####################  DEBUG-MESSAGE  ###################################*)
// 	        (*IF _debug_enable THEN
// 	          LOG_CL.NEW_MSG := 'SND: >~1<';
// 	          LOG_CL.PRINTF[1] := snd_text;
// 	          LOG_MSG();
// 	        END_IF;*)
// 	        (*########################################################################*)
// 		_STRING_TO_BUFFER(str:=snd_text,pos:=0,pt:=REF(S_BUF1.BUFFER),size:=SIZEOF(S_BUF1.BUFFER)); (* copy string into buffer *)
// 		idx1 := LEN(snd_text);
// 		S_BUF1.BUFFER[idx1] := BYTE#16#0D;
// 		idx1 := idx1 + 1;
// 		S_BUF1.BUFFER[idx1] := BYTE#16#0A;
// 		idx1 := idx1 + 1;
// 		S_BUF1.SIZE := INT_TO_UINT(idx1); (* enter send length *)
// 		R_BUF1.SIZE := UINT#0; (* reset receive buffer *)
// 		IP_C1.R_OBSERVE := TRUE;
// 		rcv_state := 9999; (* clear last FTP status *)
// 		snd_text := '';
// 		snd_step := 10;
// 	END_IF;
//
// 10:	IF S_BUF1.SIZE = UINT#0 THEN (* wait until all data is sent *)
// 		step := next_step;
// 		snd_step := 0;
// 	END_IF;
//
// END_CASE;
//
// (* global error handling for main routine *)
// IF ERROR_T = BYTE#0 AND step > 20 THEN (* no error stored yet and main routine active *)
//
// 	IF IP_C1.ERROR > DWORD#0 AND IP_C1.C_ENABLE AND NOT IP_C1.TIME_RESET THEN  (* IP Control 1 *)
// 		IF ip_c1_red_disable = FALSE OR IP_C1.ERROR <> DWORD#16#FD000000 THEN
// 			ERROR_C := IP_C1.ERROR;
// 			ERROR_T := BYTE#2;
// 			step := 980;
// 		END_IF;
// 	END_IF;
//
// 	IF IP_C2.ERROR > DWORD#0 AND IP_C2.C_ENABLE AND NOT IP_C2.TIME_RESET THEN  (* IP Control 2 *)
// 		IF ip_c2_red_disable = FALSE OR IP_C2.ERROR <> DWORD#16#FD000000 THEN
// 			ERROR_C := IP_C2.ERROR;
// 			ERROR_T := BYTE#3;
// 			step := 980;
// 		END_IF;
// 	END_IF;
//
// 	IF FSD.MODE = BYTE#0 THEN
// 		IF FSD.ERROR > BYTE#0 THEN
// 			ERROR_C := BYTE_TO_DWORD(FSD.ERROR);
// 			ERROR_T := BYTE#4; (* file access *)
// 			step := 980;
// 		END_IF;
// 	END_IF;
// END_IF;
//
// IF TON_WAIT.Q THEN (* sequence fault - step chain *)
// 	ERROR_C := SHL(INT_TO_DWORD(step),16) OR INT_TO_DWORD(rcv_state); (* enter current step and last FTP status *)
// 	ERROR_T := BYTE#5;
// 	step := 980;
// END_IF;
//
// (* ######################################################## *)
//
// DNS_CLIENT(IP_C:=IP_C1, S_BUF:=S_BUF1, R_BUF:=R_BUF1, DOMAIN:=URL_DATA.DOMAIN, IP4_DNS:= DNS_IP4, ACTIVATE:=step=10); (* DNS-Client *)
// IPC1(IP:= DWORD#0, PORT:= WORD#0 , TIME_OUT:= TIMEOUT, IP_C:= IP_C1, S_BUF:= S_BUF1, R_BUF:= R_BUF1); (* FTP control channel *)
// IPC2(IP:= DWORD#0, PORT:= WORD#0 , TIME_OUT:= TIMEOUT, IP_C:= IP_C2, S_BUF:= S_BUF2, R_BUF:= R_BUF2); (* FTP data channel *)
//
// TON_WAIT(IN:=step = last_step AND step > 0,PT:=timeout2);
// last_step := step;
// activate_last := ACTIVATE;
// BUSY := step <> 0;
//
// IF FTP_DOWNLOAD THEN (* FILE_SERVER *)
// 	FS(FSD:=FSD,PT:=R_BUF2); (* receive file from FTP *)
// ELSE
// 	FS(FSD:=FSD,PT:=S_BUF2); (* send file to FTP *)
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CHR_TO_STRING, which is commented out.
// FUNCTION_BLOCK GET_WAN_IP
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// END_VAR
// VAR_INPUT
// 	ACTIVATE : BOOL;
// END_VAR
// VAR_OUTPUT
// 	WAN_IP4 : DWORD;
// 	DONE : BOOL;
// 	NEW_IP4 : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	HTTP_GET : HTTP_GET;
// 	activate_last : BOOL;
// 	wan_ip4_last : DWORD;
// 	state : INT;
// 	st_ip : STRING(120);
// 	i1 :	INT;
// 	p1 : INT;
// 	p2 : INT;
// 	p3 : INT;
// END_VAR
//
// CASE state OF
//
// 00:
// 	IF ACTIVATE AND NOT activate_last THEN
// 		state := 20;
// 		DONE := FALSE;
// 		ERROR_C := DWORD#0;
// 		ERROR_T := BYTE#0;
// 	END_IF;
//
// 20: (* URL for DNS AND HTTP-GET *)
// 	URL_DATA := STRING_TO_URL(STR:='checkip.dyndns.com',DEFAULT_PROTOCOL:='',DEFAULT_PATH:='');
// 	state := 40;
//
// 40:
// 	IF DNS_CLIENT.DONE THEN
// 		state := 60;
// 	ELSIF (DNS_CLIENT.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := DNS_CLIENT.ERROR;
// 		ERROR_T := BYTE#01;
// 		state := 100;
// 	END_IF;
//
// 60:
// 	IF HTTP_GET.DONE THEN
// 		state := 80;
// 	ELSIF (HTTP_GET.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := HTTP_GET.ERROR;
// 		ERROR_T := BYTE#02;
// 		state := 100;
// 	END_IF;
//
// 80:
// 	(* the following HTML page is received as the response; the IP address must be copied out of it *)
// 	(* <html><head><title>Current IP Check</title></head><body>Current IP Address: 85.127.217.33</body></html> *)
// 	st_ip := '';
// 	p1 := 0;
// 	FOR i1 := UINT_TO_INT(HTTP_GET.BODY_START) TO UINT_TO_INT(HTTP_GET.BODY_STOP) DO
// 		IF R_BUF.BUFFER[i1] = BYTE#58 THEN (* search ':' *)
// 			p1 := i1 + 2;
// 			EXIT;
// 		END_IF;
// 	END_FOR;
//
// 	p3 := p1 + 15;
// 	IF p1 > 0 AND p3 < UINT_TO_INT(HTTP_GET.BODY_STOP) THEN
// 		FOR p2 := p1 TO p3 DO
// 			IF R_BUF.BUFFER[p2] <> BYTE#60 THEN (* search '<' *)
// 				st_ip := CONCAT(st_ip,CHR_TO_STRING(C:=R_BUF.BUFFER[p2]));
// 			ELSE
// 				EXIT;
// 			END_IF;
// 		END_FOR;
// 	END_IF;
//
// 	WAN_IP4 := IP4_DECODE(STR:=st_ip);
//
// 	IF WAN_IP4 <> DWORD#0 THEN
// 		DONE := TRUE;
// 		NEW_IP4 := WAN_IP4 <> wan_ip4_last;
// 		wan_ip4_last := WAN_IP4;
// 	ELSE
// 		ERROR_C := DWORD#01; (* no valid IP found *)
// 		ERROR_T := BYTE#03;
// 	END_IF;
// 	state := 100;
//
// 100:
// 	(* UNLOCK HTTP DATA *)
// 	IF (NOT HTTP_GET.DONE) THEN
// 		state := 0;
// 		DONE := ERROR_T = BYTE#0;
// 	END_IF;
//
// END_CASE;
//
// (* ---------------- DNS --------------- *)
// DNS_CLIENT(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			DOMAIN:=URL_DATA.DOMAIN,
// 			IP4_DNS:=DWORD#00,
// 			ACTIVATE:=state=40 (* DNS Start *)
// 			);
// (* parameters not used, or not used directly *)
// (* DWORD :=DNS_CLIENT.IP4;
//    BOOL  :=DNS_CLIENT.DONE;
//    BYTE  :=DNS_CLIENT.ERROR; *)
// (* -------------------------------------*)
//
// (* ------------- HTTP_GET --------------- *)
//
// HTTP_GET(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			IP4:=DNS_CLIENT.IP4,
// 			GET:=state=60,
// 			MODE:=BYTE#2, (* HTTP 1.1 with persistent connection *)
// 			UNLOCK_BUF:=state=100,
// 			URL_DATA:=URL_DATA
// 			);
//
// (* parameters not used, or not used directly *)
// (*STRING := HTTP_GET.HTTP_STATUS;
//   UINT :=HTTP_GET.HEADER_START;
//   UINT :=HTTP_GET.HEADER_STOP;
//   UINT :=HTTP_GET.BODY_START;
//   UINT :=HTTP_GET.BODY_STOP;
//   BOOL :=HTTP_GET.DONE;
//   BYTE :=HTTP_GET.ERROR;
// *)
// (* -------------------------------------*)
//
// activate_last := ACTIVATE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CHR_TO_STRING, which is commented out.
// FUNCTION HTML_DECODE : STRING(STRING_LENGTH)
// VAR_INPUT
// 	IN : STRING(STRING_LENGTH);
// END_VAR
// VAR
// 	pos : INT;
// 	end : INT;
// 	tmp : STRING(10);
// 	code : STRING(1);
// END_VAR
//
// HTML_DECODE := IN;
//
// (* docode characters encoded by character position *)
// pos := FIND(HTML_DECODE, '&');
// WHILE pos > 0 DO
// 	tmp := MID(HTML_DECODE, 2, pos + 1);
// 	IF tmp = '#x' OR tmp = '#X' THEN
// 		(* hexadecimal representation found *)
// 		tmp := MID(HTML_DECODE, 10, pos + 3);
// 		end := FIND(tmp, ';');
// 		CODE := CHR_TO_STRING(DWORD_TO_BYTE(HEX_TO_DWORD(LEFT(tmp, end - 1))));
// 		HTML_DECODE := REPLACE(HTML_DECODE, CODE, end + 3, pos);
// 	ELSIF LEFT(tmp,1) = '#' THEN
// 		tmp := MID(HTML_DECODE, 10, pos + 2);
// 		end := FIND(tmp, ';');
// 		CODE := CHR_TO_STRING(INT_TO_BYTE(STRING_TO_INT(LEFT(tmp, end - 1))));
// 		HTML_DECODE := REPLACE(HTML_DECODE, CODE, end + 2, pos);
// 	ELSE
// 		(* replace all named characters *)
// 		tmp := MID(HTML_DECODE, 10, pos + 1);
// 		end := FIND(tmp, ';');
// 		CODE := CHR_TO_STRING(CHARCODE(LEFT(tmp, end - 1)));
// 		HTML_DECODE := REPLACE(HTML_DECODE, CODE, end + 1, pos);
// 	END_IF;
// 	(* search for next occurence of encoded characters *)
// 	pos := FINDP(HTML_DECODE, '&', pos + 1);
// END_WHILE;
// END_FUNCTION

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION HTML_ENCODE : STRING(STRING_LENGTH)
// VAR_INPUT
// 	IN : STRING(STRING_LENGTH);
// 	M : BOOL;
// END_VAR
// VAR
// 	pt_in : REF_TO ARRAY[1..255] OF BYTE;
// 	pt_out : REF_TO ARRAY[1..256] OF BYTE;
// 	pos_in : INT;
// 	pos_out : INT;
// 	stop : INT;
// 	B: BYTE;
// 	tmp : STRING(10);
// 	i: INT;
// END_VAR
//
// pt_in  := REF(in);
// pt_out := REF(HTML_ENCODE);
//
// (* scan the input string and convert unallowed characters *)
// pos_in := 1;
// stop := LEN(in);
// FOR pos_out := 1 TO 250 DO
// 	IF pos_in > stop THEN EXIT; END_IF;
// 	B := pt_in^[pos_in];
// 	CASE BYTE_TO_INT(B) OF
// 		34: (* " detected *)
// 			IF pos_out > STRING_LENGTH - 6 THEN EXIT; END_IF;
// 			pt_out^[pos_out] := 38;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 113;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 117;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 111;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 116;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 59;
// 		38: (* & detected *)
// 			IF pos_out > STRING_LENGTH - 5 THEN EXIT; END_IF;
// 			pt_out^[pos_out] := 38;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 97;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 109;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 112;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 59;
// 		60: (* < detected *)
// 			IF pos_out > STRING_LENGTH - 4 THEN EXIT; END_IF;
// 			pt_out^[pos_out] := 38;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 108;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 116;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 59;
// 		62: (* > detected *)
// 			IF pos_out > STRING_LENGTH - 4 THEN EXIT; END_IF;
// 			pt_out^[pos_out] := 38;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 103;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 116;
// 			pos_out := pos_out + 1;
// 			pt_out^[pos_out] := 59;
// 	ELSE
// 		IF M AND B > 127 THEN
// 			tmp := CHARNAME(B);
// 			(* exit if output string has not enough space *)
// 			IF pos_out + LEN(tmp) + 2 <= 250 THEN
// 				pt_out^[pos_out] := 38;
// 				pos_out := pos_out + 1;
// 				FOR i := 1 TO LEN(tmp) DO
// 					pt_out^[pos_out] := CODE(tmp, i);
// 					pos_out := pos_out + 1;
// 				END_FOR;
// 				pt_out^[pos_out] := 59;
// 			END_IF;
// 		ELSE
// 			pt_out^[pos_out] := B;
// 		END_IF;
// 	END_CASE;
// 	(* increment in by 1 *)
// 	pos_in := pos_in + 1;
// END_FOR;
//
// (* terminate the output string *)
// pt_out^[pos_out] := 0;
// END_FUNCTION

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK HTTP_GET
// VAR_IN_OUT
// 	URL_DATA : url;
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// END_VAR
// VAR_INPUT
// 	IP4 : DWORD;
// 	GET : BOOL;
// 	MODE : BYTE := BYTE#2;
// 	UNLOCK_BUF : BOOL;
// END_VAR
// VAR_OUTPUT
// 	HTTP_STATUS : STRING;
// 	HEADER_START : UINT; (* HTTP Header (Message Header) *)
// 	HEADER_STOP : UINT; (* HTTP Header (Message Header) *)
// 	BODY_START : UINT; (* Message Body *)
// 	BODY_STOP : UINT; (* Message Body *)
// 	DONE : BOOL;
// 	ERROR : DWORD;
// END_VAR
// VAR
// 	READ_HTTP :	READ_HTTP;
// 	BASE64_ENCODE_STR :	BASE64_ENCODE_STR;
// 	b64_str1 :	STRING(144);
// 	b64_str2 :	STRING(192);
// 	b64_done :	BOOL;
// 	b64_start :	BOOL;
// 	authentication :	BOOL;
// 	IP_FIFO : IP_FIFO;
// 	ip_state : BYTE;
// 	ip_id : BYTE;
// 	get_last : BOOL;
// 	ip4_stored : DWORD;
// 	state : INT;
// 	i : INT;
// 	http_1_1_host : BOOL;
// 	total_size : INT;
// 	text : STRING(STRING_LENGTH);
// 	delimiter2b : STRING(2);
// 	delimiter4b : STRING(4);
// 	del_end : STRING(4);
// 	del_pos : INT;
// 	idx : INT;
// 	rcv_timeout : BOOL;
// 	con_rdce : BOOL;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	_debug_state_last : INT := -1;*)
//
// END_VAR
//
// CASE state OF
// 00:	DONE := FALSE;
// 	IF GET AND NOT GET_LAST THEN (* wait for rising edge on SET *)
// 		state := 10;
// 		ip4_stored := IP4; (* store current IP4 *)
// 		total_size := 0;
// 		HEADER_START := UINT#0;
// 		HEADER_STOP := UINT#0;
// 		BODY_START := UINT#0;
// 		BODY_STOP := UINT#0;
// 		HTTP_STATUS := '';
// 		ERROR := DWORD#0;
// 		ip_state := BYTE#1; (* register *)
// 		authentication := FALSE;
// 		delimiter2b := '$0D$0A';
// 		delimiter4b := '$0D$0A$0D$0A';
//
// 		(* LOGIN parameters present *)
// 		(* LOGIN parameters present *)
// 		IF URL_DATA.USER <> '' AND URL_DATA.PASSWORD <> '' THEN
// 			b64_str1 := CONCAT(URL_DATA.USER,':');
// 			b64_str1 := CONCAT(b64_str1,URL_DATA.PASSWORD); (* USER:PASSWORD *)
// 			b64_start := TRUE; (* BASE64 encoding of username *)
// 			authentication := TRUE;
// 			state := 5;
// 		ELSE
// 			state := 10;
// 		END_IF;
// 	        	(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 	       		LOG_CL.NEW_MSG := 'HTTP: S00 START /IP_ID ~4';
// 	            	LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
// 	            	LOG_MSG();
// 	        	END_IF;*)
// 	        	(*########################################################################*)
// 	END_IF;
//
// 5:	IF b64_start AND b64_done THEN (* wait for base64 conversion *)
// 		b64_start := FALSE;
// 		state := 10;
// 	END_IF;
//
// 10:	IF ip_state = BYTE#3 THEN (* register with IP_Control and wait for release *)
// 		idx := 0;
// 		text := CONCAT('GET ',URL_DATA.PATH);
// 		_STRING_TO_BUFFER(STR:=text,POS:=idx,PT:=REF(S_BUF.BUFFER),SIZE:=SIZEOF(S_BUF.BUFFER)); (* copy GET path into buffer *)
// 		idx := idx + LEN(text);
//
// 		IF LEN(URL_DATA.QUERY) > 0 THEN
// 			S_BUF.BUFFER[idx] := 63; (* ASCII '?' *)
// 			idx := idx + 1;
// 			_STRING_TO_BUFFER(STR:=URL_DATA.QUERY,POS:=idx,PT:=REF(S_BUF.BUFFER),SIZE:=SIZEOF(S_BUF.BUFFER)); (* copy QUERY into buffer *)
// 			idx := idx + LEN(URL_DATA.QUERY);
// 		END_IF;
//
// 		IF MODE >= BYTE#2 THEN (* pass requested HTTP version *)
// 			text := ' HTTP/1.1';
// 		ELSE
// 			text := ' HTTP/1.0';
// 		END_IF;
// 		text := CONCAT(text,delimiter2b);
// 		_STRING_TO_BUFFER(STR:=text,POS:=idx,PT:=REF(S_BUF.BUFFER),SIZE:=SIZEOF(S_BUF.BUFFER)); (* copy HTTP version into buffer *)
// 		idx := idx + LEN(text);
//
// 		IF MODE = BYTE#1 THEN (* HTTP/1.0 with persistent connection requested *)
// 			text := CONCAT('Connection: Keep-Alive',delimiter2b); (* workaround for pcworx -> $0D$0A not possible as a string constant *)
// 		ELSIF MODE = BYTE#3 THEN (* HTTP/1.1 without persistent connection requested *)
// 			text := CONCAT('Connection: Close',delimiter2b); (* workaround for pcworx -> $0D$0A not possible as a string constant *)
// 		ELSE
// 			text := '';
// 		END_IF;
//
// 		(* take additional parameters *)
// 		text := CONCAT(text,URL_DATA.HEADER);
// 		_STRING_TO_BUFFER(STR:=text,POS:=idx,PT:=REF(S_BUF.BUFFER),SIZE:=SIZEOF(S_BUF.BUFFER)); (* copy HTTP version into buffer *)
// 		idx := idx + LEN(text);
//
// 		IF authentication THEN
// 			text := CONCAT('Authorization: Basic ',b64_str2);
// 			text := CONCAT(text,delimiter2b); (* workaround for pcworx -> $0D$0A not possible as a string constant *)
// 			_STRING_TO_BUFFER(STR:=text,POS:=idx,PT:=REF(S_BUF.BUFFER),SIZE:=SIZEOF(S_BUF.BUFFER)); (* copy HTTP host into buffer *)
// 			idx := idx + LEN(text);
// 		END_IF;
//
// 		text := CONCAT('Host: ',URL_DATA.DOMAIN);
// 		text := CONCAT(text,delimiter4b); (* workaround for pcworx -> $0D$0A not possible as a string constant *)
// 		_STRING_TO_BUFFER(STR:=text,POS:=idx,PT:=REF(S_BUF.BUFFER),SIZE:=SIZEOF(S_BUF.BUFFER)); (* copy HTTP host into buffer *)
// 		idx := idx + LEN(text);
//
// 		(* set up IP data traffic *)
// 		IF FIND(';HTTPS;https;Https;', CONCAT(URL_DATA.PROTOCOL,';')) > 1 THEN
// 			IP_C.C_PORT := SEL(URL_DATA.PORT = WORD#0,URL_DATA.PORT,WORD#443); (* enter port number *)
// 			IP_C.MAILBOX[16] := BYTE#1; (* enable SSL mode *)
// 		ELSE
// 			IP_C.C_PORT := SEL(URL_DATA.PORT = WORD#0,URL_DATA.PORT,WORD#80);  (* enter port number *)
// 		END_IF;
//
// 		IP_C.C_IP := ip4_stored; (* enter IP *)
// 		IP_C.C_MODE := BYTE#0; (* Mode: TCP+ACTIVE+PORT+IP *)
// 		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
// 		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
// 		IP_C.R_OBSERVE := TRUE; (* monitor data reception *)
// 		S_BUF.SIZE := INT_TO_UINT(idx); (* enter send length *)
// 		R_BUF.SIZE := UINT#0; (* reset receive length *)
//
// 		state := 30;
// 	        	(*####################  DEBUG-MESSAGE  ###################################*)
// 	        	(*IF _debug_enable THEN
// 	        		LOG_CL.NEW_MSG := 'HTTP: S20 SEND /IP_ID ~4';
// 	            	LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
// 	            	LOG_MSG();
// 	        	END_IF;*)
// 	        	(*########################################################################*)
// 	END_IF;
//
// 30:	rcv_timeout := (IP_C.ERROR AND DWORD#16#0000_FF00) = DWORD#16#0000_FF00;
// 	con_rdce :=     (IP_C.ERROR AND DWORD#16#FD00_0000) = DWORD#16#FD00_0000;
//
// 	IF S_BUF.SIZE = UINT#0 THEN
// 		IF total_size = 0 THEN
// 			IF R_BUF.SIZE > UINT#15 THEN; (* check length of received data *)
//
// 				(* search end of HTTP HEADER -> recognizable by $0D$A$0D$A on DEC RT-11,CP/M,MS-DOS,Windows *)
// 				del_end := '$0D$0A$0D$0A';
// 				del_pos := BUFFER_SEARCH(PT:=REF(R_BUF.BUFFER),SIZE:=R_BUF.SIZE,STR:=del_end,POS:=0,IGN:=FALSE);
//
// 				(* no end found? -> search other end encoding *)
// 				IF del_pos = 0 THEN
// 					del_end := '$0A$0A';
// 					(* search end of HTTP HEADER -> recognizable by $0A$0A on Multics, Unix derivatives, Mac OS X *)
// 					del_pos := BUFFER_SEARCH(PT:=REF(R_BUF.BUFFER),SIZE:=R_BUF.SIZE,STR:=del_end,POS:=0,IGN:=FALSE);
// 				END_IF;
//
// 				IF del_pos > 0 THEN
// 					HEADER_START := UINT#0;
// 					HEADER_STOP := INT_TO_UINT(del_pos + LEN(del_end) - 1);
// 					TOTAL_SIZE := -1; (* HTTP header received but total length not known yet *)
// 				END_IF;
// 		                	(*####################  DEBUG-MESSAGE  ###################################*)
// 		                	(*IF _debug_enable THEN
// 		                   	LOG_CL.NEW_MSG := 'HTTP: S20 HEADER_START ~1 / HEADER_STOP ~2/IP_ID ~4';
// 		                   	LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
// 		                   	LOG_CL.PRINTF[1]  := UINT_TO_STRING(HEADER_START);
// 		                   	LOG_CL.PRINTF[2]  := UINT_TO_STRING(HEADER_STOP);
// 		                   	LOG_MSG();
// 		                	END_IF;*)
// 		                	(*########################################################################*)
//
// 				(* evaluate HTTP status code: (can be 1.0 or 1.1 *)
// 				READ_HTTP(size:=HEADER_STOP,pos:=0,SRC:='',pt:=R_BUF.BUFFER);
//
// 				i := LEN(READ_HTTP.READ_HTTP);
// 				IF i > 7 THEN
// 					(* evaluate HTTP version of HOST 1.0 or 1.1 *)
// 					text:=MID(READ_HTTP.READ_HTTP,1,3);
// 					http_1_1_host := text = '1';
//
// 					(* isolate code and message *)
// 					HTTP_STATUS := RIGHT(READ_HTTP.READ_HTTP,i-4);
// 				END_IF;
//
// 				(* positive HTTP request? *)
// 				IF FIND(HTTP_STATUS,'200') > 0 THEN
//
// 					(* 'Content-Length: read and evaluate value *)
// 					READ_HTTP(size:=HEADER_STOP,pos:=0,SRC:='CONTENT-LENGTH',pt:=R_BUF.BUFFER);
// 					text:=READ_HTTP.READ_HTTP;
//
// 					IF LEN(text) >= 1 THEN
// 						IF IS_NUM(str:=text) THEN
// 							(* calculate expected total number of bytes *)
// 							total_size := STRING_TO_INT(text) + UINT_TO_INT(HEADER_STOP) + 1;
//                            					(*####################  DEBUG-MESSAGE  ###################################*)
//                            					(*IF _debug_enable THEN
//                               					LOG_CL.NEW_MSG := 'HTTP: S20 HEADER ~1, BODY ~2 TOTAL ~3 /IP_ID ~4';
//                              					LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
//                                					LOG_CL.PRINTF[1]  := UINT_TO_STRING(HEADER_STOP);
//                                					LOG_CL.PRINTF[2]  := INT_TO_STRING(TOTAL_SIZE-UINT_TO_INT(HEADER_STOP));
//                                					LOG_CL.PRINTF[3]  := INT_TO_STRING(TOTAL_SIZE);
//                               					LOG_MSG();
//                           					END_IF;*)
//                           					(*########################################################################*)
// 						END_IF;
// 					END_IF;
// 				ELSE
// 					ERROR := DWORD#16#000000FC; (* error from HTTP status codes *)
// 				END_IF;
// 			END_IF;
// 		(* 1. with a known Content-Length the connection is closed after the expected data is received *)
// 		(* 2. with a persistent connection and unknown data amount, wait for the receive timeout *)
// 		(* 3. if the host closes the connection after data was received *)
// 		ELSIF (R_BUF.SIZE >= INT_TO_UINT(total_size)) OR
// 			(R_BUF.SIZE > UINT#0 AND total_size < 0 AND ((rcv_timeout OR con_rdce) OR (IP_C.C_STATE = BYTE#1))) THEN
//
// 			IF R_BUF.SIZE - 1 > HEADER_STOP THEN (* BODY data present *)
// 				BODY_START := HEADER_STOP + UINT#1;
// 				BODY_STOP := R_BUF.SIZE - 1;
// 		                	(*####################  DEBUG-MESSAGE  ###################################*)
// 		                	(*IF _debug_enable THEN
// 		                   	LOG_CL.NEW_MSG := 'HTTP: S20 BODY ~1, START ~2 END ~3 /IP_ID ~4';
// 		                   	LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
// 		                   	LOG_CL.PRINTF[1]  := UINT_TO_STRING(HEADER_STOP);
// 		                   	LOG_CL.PRINTF[2]  := UINT_TO_STRING(BODY_START);
// 		                   	LOG_CL.PRINTF[3]  := UINT_TO_STRING(BODY_STOP);
// 		                   	LOG_MSG();
// 		               	END_IF;*)
// 		               	(*########################################################################*)
// 			END_IF;
//
// 			DONE := TRUE;
// 			IP_C.C_ENABLE := FALSE; (* close connection *)
// 			state := 40;
//             		(*####################  DEBUG-MESSAGE  ###################################*)
//             		(*IF _debug_enable THEN
//                 			LOG_CL.NEW_MSG := 'HTTP: S30 OK /IP_ID ~4';
//                 			LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
//                 			LOG_MSG();
//             		END_IF;*)
//             		(*########################################################################*)
// 		END_IF;
// 	END_IF;
//
// 	IF IP_C.ERROR <> DWORD#0 THEN
// 		(* with unknown content size it is not known whether or how many telegrams will still follow *)
// 		(* so after a receive timeout or remote disconnect error the error is filtered out again (logical error) *)
// 		IF rcv_timeout AND DONE THEN
// 			IP_C.ERROR := IP_C.ERROR AND DWORD#16#FFFF_00FF; (* reset error again *)
// 		ELSIF con_rdce AND DONE THEN
// 			IP_C.ERROR := IP_C.ERROR AND DWORD#16#00FF_FFFF; (* reset error again *)
// 		ELSE
// 			ERROR := IP_C.ERROR;
// 		END_IF;
// 	END_IF;
//
// 	IF ERROR > DWORD#0 THEN
// 		IP_C.C_ENABLE := FALSE; (* block connection setup *)
// 		state := 40;
// 	        	(*####################  DEBUG-MESSAGE  ###################################*)
// 	        	(*IF _debug_enable THEN
// 	            	LOG_CL.NEW_MSG := 'HTTP: S30 ERROR ~1 /IP_ID ~4';
// 	            	LOG_CL.PRINTF[1]  := DWORD_TO_STRING(ERROR);
// 	            	LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
// 	            	LOG_MSG();
// 	        	END_IF;*)
// 	        	(*########################################################################*)
// 	END_IF;
//
// 40:	IF (UNLOCK_BUF OR NOT DONE) THEN (* wait for release and connection closed -> remove IP_ID from stack *)
// 		IP_STATE := BYTE#4; (* unregister *)
// 		DONE := FALSE;
// 		IP_C.MAILBOX[16] := BYTE#0; (* end SSL mode *)
// 		state := 0;
// 	END_IF;
// END_CASE;
//
// get_last := GET; (* save last status for edge detection *)
//
// IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id); (* IP_FIFO *)
//
// BASE64_ENCODE_STR(RUN:=b64_start,STR1:=b64_str1,STR2:=b64_str2); (* BASE64 encoding *)
// b64_done:=BASE64_ENCODE_STR.DONE;
//
// (*####################  DEBUG-MESSAGE  ###################################*)
// (*IF _debug_enable THEN
// 	IF state <> _debug_state_last THEN
// 		LOG_CL.NEW_MSG := 'http_get: /last ~1 /New ~2';
// 		LOG_CL.PRINTF[1] := INT_TO_STRING(_debug_state_last);
// 		LOG_CL.PRINTF[2] := INT_TO_STRING(state);
// 		LOG_MSG();
// 		_debug_state_last := state;
// 	END_IF;
// END_IF;*)
// (*########################################################################*)
// END_FUNCTION_BLOCK

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION_BLOCK INI_PARSER_BUF
// VAR_IN_OUT
// 	STR : STRING(STRING_LENGTH);
// 	RUN : BYTE;
// 	OFFSET : UDINT;
// 	KEY : STRING(STRING_LENGTH);
// 	VALUE : STRING(STRING_LENGTH);
// 	PT : NETWORK_BUFFER;
// END_VAR
// VAR_OUTPUT
// 	RESULT : BYTE;
// END_VAR
// VAR
// 	step : INT;
// 	size : INT;
// 	c : BYTE;
// 	k_add : BYTE;
// 	v_add : BYTE;
// 	i : INT;
// 	state : BYTE;
// 	eof : BOOL;
// 	pv : REF_TO BYTE;
// 	pk : REF_TO BYTE;
// 	pve : DWORD;
// 	pke : DWORD;
// 	watchdog : TON;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	step_last : INT;*)
// END_VAR
//
// IF RUN > 0 THEN
// 	IF state = 0 THEN (* restart evaluation *)
//
// 		pv := REF(VALUE);
// 		pv^ := 0; (* VALUE = '' *)
// 		pve := pv + INT_TO_DWORD(STRING_LENGTH);
// 		pk := REF(KEY);
// 		pk^ := 0; (* KEY = '' *)
// 		pke := pk + INT_TO_DWORD(STRING_LENGTH);
//
// 		i := UDINT_TO_INT(OFFSET);
// 		size := UINT_TO_INT(PT.SIZE) - 1;
// 		RESULT := 0;
// 		state := 5; (* command in progress (busy) *)
// 		step := 0; (* restart evaluation *)
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'IPB: Startindex: ~1';
// 			LOG_CL.PRINTF[1]  := INT_TO_STRING(i);
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
// 	END_IF;
//
// 	watchdog(IN:=FALSE,PT:=T#1ms); (* restart timer *)
// 	watchdog.IN := TRUE;
//
// 	WHILE (state = 5) DO (* loop as long as there is no result *)
//
// 		watchdog(); (* update timer *)
//
// 		IF watchdog.Q THEN (* watchdog triggered *)
// 			state := 5;
// 			EXIT;
// 		ELSIF i > size THEN (* no further character available *)
// 			state := 10; (* nothing found - end of data reached *)
// 			EXIT;
// 		ELSE
//
// 			c := PT.BUFFER[i]; (* read new character from buffer *)
// 			eof := i + 1 = PT.SIZE; (* End of File *)
//
// 			k_add := 0;
// 			v_add := 0;
//
// 			CASE step OF
//
// 			0: (* no mode detected yet *)
// 				IF c = 91 THEN (* compare with "[" *)
// 					step := 200; (* activate mode for section *)
// 				ELSIF c = 59 OR c = 35 THEN (* compare with ";" and '#' *)
// 					step := 100; (* activate mode for comment *)
// 				ELSIF c > 32 THEN (* key detected *)
// 					k_add := c; (* store character *)
// 					step := 300; (* start key evaluation *)
// 				END_IF;
// 			100: (* comment detected -> search end of line *)
// 				IF c < 32 THEN (* check for control characters *)
// 					state  := 3; (* exit -> comment was found *)
// 				END_IF;
// 			200: (* section start detected -> search section end *)
// 				IF c = 93 THEN (* compare with "]" *)
// 					IF pk > REF(KEY) > 0 THEN (* there must be at least one character inside [] *)
// 						state  := 1; (* exit -> section was found *)
// 					ELSE
// 						pk := REF(KEY);
// 						pk^ := 0; (* KEY = '' *)
// 						step := 0; (* faulty section line present -> continue searching *)
// 					END_IF;
// 				ELSIF c < 32 OR eof THEN (* control character - faulty section line -> continue searching *)
// 					pk := REF(KEY);
// 					pk^ := 0; (* KEY = '' *)
// 					step := 0;
// 				ELSE
// 					k_add := c; (* store character *)
// 				END_IF;
// 			300: (* begin: key detected -> evaluate name *)
// 				IF c < 32 THEN (* control character - faulty key line -> continue searching *)
// 					pk := REF(KEY);
// 					pk^ := 0; (* KEY = '' *)
// 					step := 0;
// 				ELSIF c = 61 THEN (* check for '=' *)
// 					step := 330;
// 				ELSE
// 					k_add := c; (* store character *)
// 				END_IF;
// 			330: (* search KeyValue start *)
// 				IF c < 32 THEN (* control character  *)
// 					state := 2; (* key without value found *)
// 				ELSIF eof THEN
// 					state := 2; (* key without value found *)
// 					v_add := c; (* store character *)
// 				ELSE
// 					v_add := c; (* store character *)
// 					step := 350;
// 				END_IF;
// 			350: (* search KeyValue end *)
// 				IF c < 32 THEN (* control character *)
// 					state := 2; (* key found *)
// 				ELSIF eof THEN
// 					v_add := c; (* store character *)
// 					state := 2; (* key found *)
// 				ELSE
// 					v_add := c; (* store character *)
// 				END_IF;
// 			END_CASE;
//
// 			IF v_add > 0  AND pv < pve THEN (* add character to value string *)
// 				pv^ := v_add; (* enter character *)
// 				pv := pv + 1;
// 				pv^ := 0; (* enter string termination *)
// 			END_IF;
//
// 			IF k_add > 0 AND pk < pke THEN (* add character to key string *)
// 				pk^ := k_add; (* enter character *)
// 				pk := pk + 1;
// 				pk^ := 0; (* enter string termination *)
// 			END_IF;
//
// 			i := i + 1; (* increment data pointer *)
//
// 		END_IF;
// 	END_WHILE;
//
// 	CASE state OF
// 	01:	(* section found *)
// 		IF RUN = 1 OR RUN = 3 THEN (* section preselected or report everything *)
// 			IF KEY = STR OR LEN(STR) = 0 THEN (* compare with searched section or no section specified *)
// 				RESULT := 1; (* section found *)
// 				OFFSET := i; (* take offset *)
// 				RUN := 0;
// 				state := 0;
// 			END_IF;
// 		ELSE
// 			pk := REF(KEY);
// 			pk^ := 0; (* KEY = '' *)
// 			RESULT := 11; (* key not found, section end - do not change offset *)
// 			RUN := 0;
// 			state := 0;
// 		END_IF;
// 	02:	(* key found *)
// 		IF RUN = 2 OR RUN = 3 THEN (* key preselected or report everything *)
// 			IF KEY = STR OR LEN(STR) = 0 THEN (* compare with searched key *)
// 				RESULT := 2; (* key found *)
// 				OFFSET := i; (* take offset *)
// 				RUN := 0;
// 				state := 0;
// 			END_IF;
// 		END_IF;
// 	03:	(* comment found *)
// 		OFFSET := i; (* take offset *)
// 		state := 0;
//
// 	05:	RESULT := 5; (* in progress *)
//
// 	10:	(* index out of range - end of buffer *)
// 		RESULT := 10; (* nothing found - end of data reached *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
// 	END_CASE;
//
// 	(*####################  DEBUG-MESSAGE  ###################################*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'IPF: intern Result ~1 , state ~2, step ~3 , offset ~4, run ~5';
// 		LOG_CL.PRINTF[1] := BYTE_TO_STRING(RESULT);
// 		LOG_CL.PRINTF[2] := BYTE_TO_STRING(state);
// 		LOG_CL.PRINTF[3] := INT_TO_STRING(step);
// 		LOG_CL.PRINTF[4] := UDINT_TO_STRING(OFFSET);
// 		LOG_CL.PRINTF[5] := BYTE_TO_STRING(RUN);
// 		LOG_MSG();
// 	END_IF;*)
// 	(*########################################################################*)
//
// END_IF;
//
// (*####################  DEBUG-MESSAGE  ###################################*)
// (*IF _debug_enable THEN
// 	IF step_last <> step THEN
// 		LOG_CL.NEW_MSG := 'IPB: Step: ~1';
// 		LOG_CL.PRINTF[1]  := INT_TO_STRING(step);
// 		LOG_MSG();
// 	END_IF;
// 	step_last := step;
// END_IF;*)
// (*########################################################################*)
// END_FUNCTION_BLOCK

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION_BLOCK INI_PARSER_FILE
// VAR_IN_OUT
// 	FILENAME : STRING;
// 	FSD : FILE_SERVER_DATA	;
// 	STR : STRING(STRING_LENGTH);
// 	RUN : BYTE;
// 	OFFSET : UDINT;
// 	KEY : STRING(STRING_LENGTH);
// 	VALUE : STRING(STRING_LENGTH);
// 	PT	: NETWORK_BUFFER;
// END_VAR
// VAR_OUTPUT
// 	RESULT : BYTE;
// END_VAR
// VAR
// 	FB : FILE_BLOCK;
// 	i : UDINT;
// 	state : BYTE;
// 	e : BYTE;
// 	mode : BYTE;
// 	step : INT;
// 	c : BYTE;
// 	k_add : BYTE;
// 	v_add : BYTE;
// 	eof : BOOL;
// 	pv : REF_TO BYTE;
// 	pk : REF_TO BYTE;
// 	pve : DWORD;
// 	pke : DWORD;
// 	watchdog : TON;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	step_last : INT;*)
// END_VAR
//
// IF RUN > 0 THEN
//
// 	IF state = 0 THEN (* restart evaluation *)
// 		pv := REF(VALUE);
// 		pv^ := 0; (* VALUE = '' *)
// 		pve := pv + INT_TO_DWORD(STRING_LENGTH);
// 		pk := REF(KEY);
// 		pk^ := 0; (* KEY = '' *)
// 		pke := pk + INT_TO_DWORD(STRING_LENGTH);
//
// 		i := OFFSET;
// 		RESULT := 0;
// 		state := 5; (* command in progress (busy) *)
// 		step := 0; (* restart evaluation *)
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'IPF: Startindex: ~1';
// 			LOG_CL.PRINTF[1] := UDINT_TO_STRING(i);
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
// 	END_IF;
//
// 	watchdog(IN:=FALSE,PT:=T#1ms); (* restart timer *)
// 	watchdog.IN := TRUE;
//
// 	WHILE (state = 5) DO (* loop as long as there is no result *)
//
// 		mode := 1;
// 		FB(MODE:=mode,FSD:=fsd,PT:=PT,FILENAME:=FILENAME,POS:=i );
// 		e:=FB.ERROR;
// 		c:=FB.DATA;
//
// 		eof := i + 1 = FSD.FILE_SIZE; (* End of File *)
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'IPF: mode: ~1, error: ~2 , >~7<>~3< ~4 step ~5 index ~6';
// 			LOG_CL.PRINTF[1] := BYTE_TO_STRING(mode);
// 			LOG_CL.PRINTF[2] := BYTE_TO_STRING(e);
// 			LOG_CL.PRINTF[3] := BYTE_TO_STRING(c);
// 			LOG_CL.PRINTF[4] := BOOL_TO_STRING(eof);
// 			LOG_CL.PRINTF[5] := INT_TO_STRING(step);
// 			LOG_CL.PRINTF[6] := UDINT_TO_STRING(i);
// 			LOG_CL.PRINTF[7] := CHR_TO_STRING(c);
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
// 		watchdog(); (* update timer *)
//
// 		IF mode > 0 OR watchdog.Q THEN (* block still active, or watchdog triggered *)
// 			state := 5;
// 			EXIT;
// 		ELSIF e > 0 THEN (* File_error or no further character available *)
// 			state := 10; (* nothing found - end of data reached *)
// 			EXIT;
// 		ELSE
//
// 			k_add := 0;
// 			v_add := 0;
//
// 			CASE step OF
//
// 			0: (* no mode detected yet *)
// 				IF c = 91 THEN (* compare with "[" *)
// 					step := 200; (* activate mode for section *)
// 				ELSIF c = 59 OR c = 35 THEN (* compare with ";" and '#' *)
// 					step := 100; (* activate mode for comment *)
// 				ELSIF c > 32 THEN (* key detected *)
// 					k_add := c; (* store character *)
// 					step := 300; (* start key evaluation *)
// 				END_IF;
// 			100: (* comment detected -> search end of line *)
// 				IF c < 32 THEN (* check for control characters *)
// 					state  := 3; (* exit -> comment was found *)
// 				END_IF;
// 			200: (* section start detected -> search section end *)
// 				IF c = 93 THEN (* compare with "]" *)
// 					IF pk > REF(KEY) > 0 THEN (* there must be at least one character inside [] *)
// 						state  := 1; (* exit -> section was found *)
// 					ELSE
// 						pk := REF(KEY);
// 						pk^ := 0; (* KEY = '' *)
// 						step := 0; (* faulty section line present -> continue searching *)
// 					END_IF;
// 				ELSIF c < 32 OR eof THEN (* control character - faulty section line -> continue searching *)
// 					pk := REF(KEY);
// 					pk^ := 0; (* KEY = '' *)
// 					step := 0;
// 				ELSE
// 					k_add := c; (* store character *)
// 				END_IF;
// 			300: (* begin: key detected -> evaluate name *)
// 				IF c < 32 THEN (* control character - faulty key line -> continue searching *)
// 					pk := REF(KEY);
// 					pk^ := 0; (* KEY = '' *)
// 					step := 0;
// 				ELSIF c = 61 THEN (* check for '=' *)
// 					step := 330;
// 				ELSE
// 					k_add := c; (* store character *)
// 				END_IF;
// 			330: (* search KeyValue start *)
// 				IF c < 32 THEN (* control character  *)
// 					state := 2; (* key without value found *)
// 				ELSIF eof THEN
// 					state := 2; (* key without value found *)
// 					v_add := c; (* store character *)
// 				ELSE
// 					v_add := c; (* store character *)
// 					step := 350;
// 				END_IF;
// 			350: (* search KeyValue end *)
// 				IF c < 32 THEN (* control character *)
// 					state := 2; (* key found *)
// 				ELSIF eof THEN
// 					v_add := c; (* store character *)
// 					state := 2; (* key found *)
// 				ELSE
// 					v_add := c; (* store character *)
// 				END_IF;
// 			END_CASE;
//
// 			IF v_add > 0  AND pv < pve THEN (* add character to value string *)
// 				pv^ := v_add; (* enter character *)
// 				pv := pv + 1;
// 				pv^ := 0; (* enter string termination *)
// 			END_IF;
//
// 			IF k_add > 0 AND pk < pke THEN (* add character to key string *)
// 				pk^ := k_add; (* enter character *)
// 				pk := pk + 1;
// 				pk^ := 0; (* enter string termination *)
// 			END_IF;
//
// 			i := i + 1; (* increment data pointer *)
//
// 		END_IF;
// 	END_WHILE;
//
// 	CASE state OF
// 	01:	(* section found *)
// 		IF RUN = 1 OR RUN = 3 THEN (* section preselected or report everything *)
// 			IF KEY = STR OR LEN(STR) = 0 THEN (* compare with searched section or no section specified *)
// 				RESULT := 1; (* section found *)
// 				OFFSET := i; (* take offset *)
// 				RUN := 0;
// 				state := 0;
// 			END_IF;
// 		ELSE
// 			pk := REF(KEY);
// 			pk^ := 0; (* KEY = '' *)
// 			RESULT := 11; (* key not found, section end - do not change offset *)
// 			RUN := 0;
// 			state := 0;
// 		END_IF;
// 	02:	(* key found *)
// 		IF RUN = 2 OR RUN = 3 THEN (* key preselected or report everything *)
// 			IF KEY = STR OR LEN(STR) = 0 THEN (* compare with searched key *)
// 				RESULT := 2; (* key found *)
// 				OFFSET := i; (* take offset *)
// 				RUN := 0;
// 				state := 0;
// 			END_IF;
// 		END_IF;
// 	03:	(* comment found *)
// 		OFFSET := i; (* take offset *)
// 		state := 0;
//
// 	05:	RESULT := 5; (* in progress *)
//
// 	10:	(* index out of range - end of buffer *)
// 		RESULT := 10; (* nothing found - end of data reached *)
// 		OFFSET := i; (* take offset *)
// 		RUN := 0;
// 		state := 0;
// 	END_CASE;
//
// 	(*####################  DEBUG-MESSAGE  ###################################*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'IPF: intern Result ~1 , state ~2, step ~3 , offset ~4, run ~5';
// 		LOG_CL.PRINTF[1] := BYTE_TO_STRING(RESULT);
// 		LOG_CL.PRINTF[2] := BYTE_TO_STRING(state);
// 		LOG_CL.PRINTF[3] := INT_TO_STRING(step);
// 		LOG_CL.PRINTF[4] := UDINT_TO_STRING(OFFSET);
// 		LOG_CL.PRINTF[5] := BYTE_TO_STRING(RUN);
// 		LOG_MSG();
// 	END_IF;*)
// 	(*########################################################################*)
//
// END_IF;
//
// (*####################  DEBUG-MESSAGE  ###################################*)
// (*IF _debug_enable THEN
// 	IF step_last <> step THEN
// 		LOG_CL.NEW_MSG := 'IPB: Step: ~1';
// 		LOG_CL.PRINTF[1]  := INT_TO_STRING(step);
// 		LOG_MSG();
// 	END_IF;
// 	step_last := step;
// END_IF;*)
// (*########################################################################*)
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses FLOAT_TO_REAL, which is commented out.
// FUNCTION_BLOCK IP2GEO
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// 	GEO : IP2GEO_DATA;
// END_VAR
// VAR_INPUT
// 	IP : DWORD;
// 	ACTIVATE : BOOL;
// END_VAR
// VAR_OUTPUT
// 	BUSY : BOOL;
// 	DONE : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	CTRL : XML_CONTROL;
// 	XML_READER : XML_READER;
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	HTTP_GET : HTTP_GET;
// 	last_state : BOOL;
// 	value_int :	INT;
// 	value_real : REAL;
// 	v_real : REAL;
// 	state : INT;
// END_VAR
// CASE state OF
//
// 00:
// 	IF ACTIVATE AND NOT last_state THEN
// 		state := 20;
// 		DONE := FALSE;
// 		BUSY := TRUE;
// 		ERROR_C := DWORD#0;
// 		ERROR_T := BYTE#0;
// 	END_IF;
//
// 20: (* URL for DNS AND HTTP-GET *)
// 	(* Example: http://ipinfodb.com/ip_query.php?ip=10.10.10.10 *)
// 	(* Example: http://ipinfodb.com/ip_query.php?timezone=true&IP=*)
// 	URL_DATA:=STRING_TO_URL(STR:='http://ipinfodb.com/ip_query.php?timezone=true&IP=',DEFAULT_PROTOCOL:='',DEFAULT_PATH:='');
//
// 	IF IP > DWORD#0 THEN
// 		URL_DATA.QUERY := CONCAT(URL_DATA.QUERY,IP4_TO_STRING(IP4:=IP));
// 	END_IF;
//
// 	state := 40;
//
// 40:
// 	IF DNS_CLIENT.DONE THEN
// 		state := 60;
// 	ELSIF (DNS_CLIENT.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := DNS_CLIENT.ERROR;
// 		ERROR_T := BYTE#01;
// 		state := 100;
// 	END_IF;
//
// 60:
// 	IF HTTP_GET.DONE THEN
// 		state := 80;
// 		CTRL.START_POS := HTTP_GET.BODY_START;
// 		CTRL.STOP_POS := HTTP_GET.BODY_STOP;
// 		CTRL.COMMAND := WORD#2#10000000_00011000; (* ONLY TEXT AND ATTRIBUTE *)
// 		CTRL.WATCHDOG := T#1ms;
//
// 	ELSIF (HTTP_GET.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := HTTP_GET.ERROR;
// 		ERROR_T := BYTE#02;
// 		state := 100;
// 	END_IF;
//
// 80:
// 	(* read XML data serially *)
// 	XML_READER(CTRL:=CTRL,BUF:=R_BUF.BUFFER);
// 	IF CTRL.TYP < 98 THEN (* evaluate only if not a watchdog pass *)
// 		(* perform automatic conversion to real and int *)
// 		value_int := 0;
// 		value_real := 0.0;
// 		IF LEN(CTRL.VALUE) <= 20 THEN
// 			v_real :=FLOAT_TO_REAL(FLT:=CTRL.VALUE);
// 			IF CHK_REAL(v_real) = 0 THEN (* !!! checks for a valid REAL value, otherwise TwinCAT crashes !!!!  *)
// 				value_real := v_real;
// 				value_int := REAL_TO_INT(value_real);
// 			END_IF;
// 		END_IF;
//
// 		CASE UINT_TO_INT(CTRL.COUNT) OF
//
// 		07: GEO.IP4 := IP4_DECODE(STR:=CTRL.VALUE); (* IP address *)
// 		10: GEO.STATE := 'OK' = CTRL.VALUE; (* STATUS *)
// 		13: GEO.COUNTRY_CODE := CTRL.VALUE; (* Country-Code *)
// 		16: GEO.COUNTRY_NAME := CTRL.VALUE; (* Country-Name *)
// 		19: GEO.REGION_CODE := CTRL.VALUE; (* Region-Code *)
// 		22: GEO.REGION_NAME := CTRL.VALUE; (* Region-Name *)
// 		25: GEO.CITY := CTRL.VALUE; (* City *)
// 		31: GEO.GEO_LATITUDE := value_real; (* Latitude *)
// 		34: GEO.GEO_LONGITUDE := value_real; (* Longitude *)
// 		37: GEO.TIME_ZONE_NAME := CTRL.VALUE; (* TimeZone Name *)
// 		40: GEO.GMT_OFFSET := value_int; (* GMT OFFSET in seconds *)
// 		43: GEO.IS_DST := value_int > 0; (* daylight saving time active *)
// 		END_CASE;
//
// 	ELSIF CTRL.TYP = 99 THEN (* EXIT - last element read *)
// 		DONE := TRUE;
// 		state := 100;
// 	END_IF;
//
// 100:
// 	(* UNLOCK HTTP DATA *)
// 	IF (NOT HTTP_GET.DONE) THEN
// 		state := 0;
// 		BUSY := FALSE;
// 		DONE := ERROR_T = BYTE#0;
// 	END_IF;
//
// END_CASE;
//
// (* ---------------- DNS --------------- *)
// DNS_CLIENT(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			DOMAIN:=URL_DATA.DOMAIN,
// 			IP4_DNS:=DWORD#00,
// 			ACTIVATE:=state=40 (* DNS Start *)
// 			);
//
// (* parameters not used, or not used directly *)
// (* 	DWORD :=DNS_CLIENT.IP4;
// 	BOOL  :=DNS_CLIENT.DONE;
// 	BYTE  :=DNS_CLIENT.ERROR; *)
// (* -------------------------------------*)
//
// (* ------------- HTTP_GET --------------- *)
// HTTP_GET(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			IP4:=DNS_CLIENT.IP4,
// 			GET:=state=60,
// 			MODE:=BYTE#2, (* HTTP 1.1 with persistent connection *)
// 			UNLOCK_BUF:=state=100,
// 			URL_DATA:=URL_DATA
// 			);
//
// (* parameters not used, or not used directly *)
// (*	STRING := HTTP_GET.HTTP_STATUS;
// 	UINT :=HTTP_GET.HEADER_START;
// 	UINT :=HTTP_GET.HEADER_STOP;
// 	UINT :=HTTP_GET.BODY_START;
// 	UINT :=HTTP_GET.BODY_STOP;
// 	BOOL :=HTTP_GET.DONE;
// 	BYTE :=HTTP_GET.ERROR;
// *)
// (* -------------------------------------*)
//
// last_state := ACTIVATE;
// END_FUNCTION_BLOCK

FUNCTION IP4_CHECK : BOOL
VAR_INPUT
	NIP : DWORD;
	LIP : DWORD;
	SM : DWORD;
END_VAR

IP4_CHECK := (NIP AND sm) = (LIP AND SM);
END_FUNCTION

FUNCTION IP4_DECODE : DWORD
VAR_INPUT
	str : STRING(15);
END_VAR
VAR
	pos: INT;
END_VAR

pos := FIND(str, '.');
WHILE pos > 0 DO
	IP4_DECODE := SHL(IP4_DECODE,8) OR STRING_TO_DWORD(LEFT(str,pos - 1));
	str := DELETE(str, pos, 1);
	pos := FIND(str, '.');
END_WHILE;
IP4_DECODE := SHL(IP4_DECODE,8) OR STRING_TO_DWORD(str);
END_FUNCTION

FUNCTION IP4_TO_STRING : STRING(15)
VAR_INPUT
	IP4 : DWORD;
END_VAR
VAR
	tmp : STRING(3);
END_VAR

(* mask bits AND INSERT AT given position *)
tmp := BYTE_TO_STRING(DWORD_TO_BYTE(IP4));
IP4_TO_STRING := CONCAT('...', tmp);
tmp := BYTE_TO_STRING(DWORD_TO_BYTE(SHR(IP4,8)));
IP4_TO_STRING := INSERT(IP4_TO_STRING, tmp, 2);
tmp := BYTE_TO_STRING(DWORD_TO_BYTE(SHR(IP4,16)));
IP4_TO_STRING := INSERT(IP4_TO_STRING, tmp, 1);
tmp := BYTE_TO_STRING(DWORD_TO_BYTE(SHR(IP4,24)));
IP4_TO_STRING := CONCAT(tmp, IP4_TO_STRING);
END_FUNCTION

FUNCTION_BLOCK IP_CONTROL
VAR_IN_OUT
	IP_C : IP_C; (* IP_Control management structure *)
	S_BUF : NETWORK_BUFFER; (* external send buffer *)
	R_BUF : NETWORK_BUFFER; (* external receive buffer *)
END_VAR
VAR_INPUT
	IP : DWORD; (* IP address *)
	PORT : WORD; (* port number *)
	TIME_OUT : TIME; (* monitoring time connect / send / receive *)
END_VAR
VAR
	c_time : UDINT; (* CONNECT time (time without a connection) *)
	s_time : UDINT; (* SEND time    (time since last data sent) *)
	r_time : UDINT; (* RECEIVE time (time since last data received) *)

	c_enable : BOOL; (* Connect: enable *)
	c_ip : DWORD; (* Connect: current IP address *)
	c_ip_str : STRING(15);
	c_port : WORD; (* Connect: current port number *)
	c_mode : BYTE; (* Connect: current mode UDP/TCP...... *)
	c_status : BYTE; (* Connect: Status (Errorcode) *)
	c_ready : BOOL; (* Connect: connection established *)
	c_ready_old : BOOL; (* Connect: last state of c_ready for edge evaluation *)

	s_total : INT; (* Send: total number of bytes to send *)
	s_cur_pos : INT; (* Send: current send position *)
	s_cur_size : INT; (* Send: current send size *)
	s_max_size : INT; (* Send: maximum telegram length for sending *)
	s_active : BOOL; (* Send: data sending is active *)
	s_status : BYTE; (* Send: Status (Errorcodes) *)
	s_ip_str : STRING(15);
	s_ip_exist : BOOL;

	r_status : BYTE; (* Receive: Status (Errorcodes) *)
	r_count : INT; (* Receive: number of received bytes *)
	r_offset : INT; (* Receive: current position in external receive buffer *)
	r_max_size : INT; (* Receive: size of the external receive buffer*)

	new_connection : BOOL; (* new connection parameters available *)
	tx : DWORD; (* last time value *)
	error_time : DWORD; (* current timeout value *)

	state: BYTE; (* step flag for connect *)
	bytes_received : DINT;
	udp_mode : BOOL;

	socket : T_HSOCKET;
	server : T_HSOCKET;
	fbCreate	: FB_SocketUdpCreate;
	fbConnect : FB_SocketConnect;
	fbClose : FB_SocketClose;
	fbSendTo : FB_SocketUdpSendTo;
	fbReceiveFrom : FB_SocketUdpReceiveFrom;
	fbSend : FB_SocketSend;
	fbReceive : FB_SocketReceive;
	fbListen : FB_SocketListen;
	fbAccept : FB_SocketAccept;
	nErrId : UDINT;
	Send_Busy : BOOL;
	Send_Error : BOOL;
	Send_ErrId : UDINT;

	IPCR : IP_CONTROL_RESET;

END_VAR

VAR CONSTANT
	STOP : BYTE := 0;
	UDP_INIT : BYTE:= 31;
	UDP_INIT_WAIT : BYTE:= 32;
	TC_INIT : BYTE := 51;
	TC_CONNECT : BYTE:= 52;
	TS_INIT : BYTE:= 53;
	TS_LISTEN_WAIT : BYTE:= 54;
	TS_ACCEPT : BYTE:= 55;
	TS_ACCEPT_WAIT : BYTE:= 56;
	C_CLOSE : BYTE := 190;
	C_CLOSE_WAIT : BYTE:= 191;
	C_CLOSE2 : BYTE := 192;
	C_CLOSE_WAIT2 : BYTE:= 193;
	C_WAIT : BYTE := 200;
	C_REC : BYTE := 210;
	C_REC_WAIT : BYTE := 211;
	C_SND : BYTE := 220;
	C_SND_WAIT1 : BYTE := 221;
	C_SND_WAIT2 : BYTE := 222;
END_VAR

(* on PLC reset or program download all old connections are closed *)
IPCR();
IF NOT IPCR.READY THEN RETURN; END_IF;

(*####################  DEBUG-MESSAGE  ###################################*)
(*########################################################################*)

(* read system time *)
tx := T_PLC_MS();

IF IP_C.C_PORT = 0 THEN
	IP_C.C_PORT := PORT; (* replace the given PORT with the configured PORT *)
END_IF;
IF IP_C.C_IP = 0 AND IP_C.C_MODE < 4 THEN
	IP_C.C_IP := IP; (* replace the given IP with the configured IP *)
END_IF;

new_connection := c_ip <> IP_C.C_IP OR c_mode <> IP_C.C_MODE OR c_port <> IP_C.C_PORT; (* new connection requested? *)
c_enable := IP_C.C_ENABLE AND NOT new_connection;

IF IP_C.C_ENABLE AND IP_C.ERROR = 0 AND state = STOP THEN (* enabled, no error, no socket open *)
	c_ip   := IP_C.C_IP;
	c_mode := IP_C.C_MODE;
	c_port := IP_C.C_PORT;
	c_ip_str := IP4_TO_STRING(c_ip);
	udp_mode := C_MODE.0; (* Mode 1,3,5 = udp *)

	error_time := TIME_TO_DWORD(MAX(T#200ms,TIME_OUT));
	IP_C.TIME_RESET := TRUE; (* reset all timers *)
	r_max_size := SIZEOF(R_BUF.BUFFER);
	s_max_size := SIZEOF(S_BUF.BUFFER);

	IF udp_mode THEN (* Mode: 1,3,5 = UDP CLIENT/SERVER *)
		state := UDP_INIT;
	ELSIF c_mode = 0 THEN (* Mode: 0 = TCP CLIENT *)
		state := TC_INIT;
	ELSE
		state := TS_INIT; (* Mode: 2,4 = TCP SERVER *)
	END_IF;

	IF c_mode < 5 THEN (* set receiver IP *)
		s_ip_str := c_ip_str;
		s_ip_exist := TRUE;
	ELSE
		s_ip_str := '';
		s_ip_exist := FALSE;
	END_IF;

    	(*####################  DEBUG-MESSAGE  ###################################*)
	(*########################################################################*)

END_IF;

(*####################  DEBUG-MESSAGE  ###################################*)
(*########################################################################*)

IF IP_C.TIME_RESET THEN (* reset all timers *)
	IP_C.TIME_RESET := FALSE; (* clear reset *)
	c_time := tx; (* reset connect timer *)
	s_time := tx; (* reset send timer *)
	r_time := tx; (* reset receive timer *)

	c_status := BYTE#0; (* reset last connect error *)
	s_status := BYTE#0; (* reset last send error *)
	r_status := BYTE#0; (* reset last receive error *)
END_IF;

(* ---------------------------------------------- *)

CASE state OF

UDP_INIT: (* UDP_CLIENT + SERVER *)

	fbCreate(bExecute := FALSE);
	fbCreate(	bExecute := TRUE , sSrvNetId := sSrvNetId, sLocalHost := sLocalHost, nLocalPort := SEL(c_mode = 1,c_port,0) , tTimeOut := TIME_OUT);
	fbCreate.bExecute := FALSE; (* presetting *)
	state := UDP_INIT_WAIT;

UDP_INIT_WAIT:
	fbCreate();
	IF NOT fbCreate.bBusy THEN
		IF NOT fbCreate.bError THEN
			socket := fbCreate.hSocket;
			c_ready := TRUE; (* Connected *)
			state := C_WAIT;
		ELSE
			c_status := 1; (* FB_SocketUdpCreate failed *)
			nErrId := fbCreate.nErrId;
			state := C_CLOSE;
		END_IF;
	END_IF;

(* ---------------------------------------------- *)
TC_INIT: (* TCP CLIENT (ACTIVE) *)
	fbConnect(bExecute := FALSE);
	fbConnect(bExecute := TRUE, sSrvNetId := sSrvNetId, sRemoteHost := c_ip_str, nRemotePort := c_port, tTimeout := TIME_OUT);
	fbConnect.bExecute := FALSE; (* presetting *)
	state := TC_CONNECT;

(* ---------------------------------------------- *)
TC_CONNECT:
	IF NOT c_enable OR c_status = 255 THEN (* not enabled or connect timeout *)
		state := C_CLOSE;
	END_IF;

	fbConnect();
	IF NOT fbConnect.bBusy THEN
		IF fbConnect.bError THEN
			nErrID := fbConnect.nErrId;
			state :=TC_INIT; (* try again *)
		ELSE
			socket := fbConnect.hSocket;
			c_ready := TRUE; (* Connected *)
			state := C_WAIT;
		END_IF;
	END_IF;
(* ---------------------------------------------- *)
TS_INIT: (* TCP SERVER (PASSIVE) *)
	fbListen(bExecute := FALSE);
	fbListen(bExecute := TRUE , sSrvNetId:= sSrvNetId , sLocalHost:= sLocalHost , nLocalPort:= c_port , tTimeout:= TIME_OUT);
	fbListen.bExecute := FALSE;
	state := TS_LISTEN_WAIT;
(* ---------------------------------------------- *)
TS_LISTEN_WAIT:
	fbListen();
	IF NOT fbListen.bBusy THEN
		IF fbListen.bError THEN
			c_status := 3; (* FB_SocketListen failed *)
			nErrID := fbListen.nErrId;
			state := C_CLOSE;
		ELSE
			server := fbListen.hListener;
			state := TS_ACCEPT;
		END_IF;
	END_IF;
(* ---------------------------------------------- *)
TS_ACCEPT:
	fbAccept(bExecute := FALSE);
	fbAccept(bExecute := TRUE, sSrvNetId := sSrvNetId, hListener := server, tTimeout := TIME_OUT);
	fbAccept.bExecute := FALSE;
	state := TS_ACCEPT_WAIT;
(* ---------------------------------------------- *)
TS_ACCEPT_WAIT:
	fbAccept();
	IF NOT fbAccept.bBusy THEN

		state := TS_ACCEPT; (* presetting  *)

		IF fbAccept.bError THEN
			c_status := 4; (* FB_SocketAccept failed *)
			nErrID := fbAccept.nErrId;
			state := C_CLOSE;
		ELSE (* no error *)
			IF fbAccept.bAccepted THEN
				socket := fbAccept.hSocket;
				IF socket.remoteAddr.sAddr  <> c_ip_str AND c_mode <> 4 THEN
					state := C_CLOSE; (* reject remote - close connection *)
				ELSE
					state := C_WAIT;
					c_ready := TRUE; (* Connected *)
				END_IF;
			END_IF; (* fbAccept.bAccepted *)
		END_IF; (* NOT fbAccept.bError *)
	END_IF; (*  NOT fbAccept.bBusy *)
(* ---------------------------------------------- *)
C_CLOSE:
	c_ready := FALSE;
	IF socket.handle > 0 THEN
		fbClose(bExecute := FALSE );
		fbClose(bExecute := TRUE, sSrvNetId := sSrvNetId,	hSocket := socket);
		fbClose.bExecute := FALSE;
		state := C_CLOSE_WAIT;
	ELSE
		state:= C_CLOSE2;
	END_IF;

C_CLOSE_WAIT:
	fbClose();
	IF NOT fbClose.bBusy THEN
		socket := TCPADS_NULL_HSOCKET; (* reset socket *)
		nErrId:= SEL(fbClose.bError,0,fbClose.nErrId);
		state:= C_CLOSE2;
	END_IF;

C_CLOSE2:
	IF server.handle > 0 THEN
		fbClose(bExecute := FALSE );
		fbClose(bExecute := TRUE, sSrvNetId := sSrvNetId,	hSocket := server);
		fbClose.bExecute := FALSE;
		state := C_CLOSE_WAIT2;
	ELSE
		state:= STOP;
	END_IF;

C_CLOSE_WAIT2:
	fbClose();
	IF NOT fbClose.bBusy THEN
		server := TCPADS_NULL_HSOCKET; (* reset socket *)
		nErrId:= SEL(fbClose.bError,0,fbClose.nErrId);
		state:= STOP;
	END_IF;

(* ---------------------------------------------- *)

C_WAIT: (* wait for request to close *)
	IF NOT c_enable OR NOT c_ready THEN
		state:= C_CLOSE;
	ELSIF S_BUF.SIZE > 0 AND IP_C.MAILBOX[2] = 0 AND s_ip_exist THEN
		state:= C_SND;
	ELSIF IP_C.MAILBOX[3] = 0 THEN
		state:= C_REC;
	END_IF;

(* ---------------------------------------------- *)

C_REC: (* receive data *)
	IF R_BUF.SIZE >= INT_TO_UINT(r_max_size) THEN
		R_BUF.SIZE := UINT#0;
		r_status := 254; (* receive buffer reset due to overflow *)
	END_IF;

	r_offset := SEL(udp_mode,UINT_TO_INT(R_BUF.SIZE),0);
	r_count := r_max_size - r_offset;
	bytes_received := 0;

	IF udp_mode THEN
		fbReceiveFrom(bExecute := FALSE );
		fbReceiveFrom(bExecute:= TRUE, sSrvNetId := sSrvNetId, hSocket:= socket, cbLen := r_count, pDest:= REF(R_BUF.BUFFER), tTimeout:= TIME_OUT);
		fbReceiveFrom.bExecute := FALSE; (* presetting *)
	ELSE
		fbReceive(bExecute:= FALSE);
		fbReceive(bExecute:= TRUE, sSrvNetId := sSrvNetId , hSocket:= socket , cbLen:= r_count, pDest:= REF(R_BUF.BUFFER[r_offset]) , tTimeout:= TIME_OUT);
		fbReceive.bExecute:= FALSE;
	END_IF;
	state := C_REC_WAIT;

C_REC_WAIT:
	IF udp_mode THEN
		fbReceiveFrom(); (* UDP Mode *)
		IF NOT fbReceiveFrom.bBusy THEN
			IF NOT fbReceiveFrom.bError THEN
				IF fbReceiveFrom.nRecBytes > 0 THEN
					IF c_mode <> BYTE#5 AND fbReceiveFrom.sRemoteHost  <> c_ip_str THEN
						R_BUF.SIZE := UINT#0;
					ELSE
						bytes_received := fbReceiveFrom.nRecBytes;
						IF NOT s_ip_exist THEN
							s_ip_str := fbReceiveFrom.sRemoteHost;
							s_ip_exist := TRUE;
						END_IF;
					END_IF;
				END_IF;
				state := C_WAIT;
			ELSIF fbReceiveFrom.nErrId = 16#80072746 THEN (* connection closed by remote *)
				c_status := 253;
				c_ready := FALSE;
				state := C_WAIT;
			ELSE
				nErrId := fbReceiveFrom.nErrId;
				state := C_WAIT;
			END_IF;
		END_IF;
	ELSE
		fbReceive(); (* TCP Mode *)
		IF NOT fbReceive.bBusy THEN
			IF NOT fbReceive.bError THEN
				IF fbReceive.nRecBytes > 0 THEN
					bytes_received := fbReceive.nRecBytes;
				END_IF;
			ELSIF fbReceive.nErrId = 16#8004 THEN (* connection closed by remote *)
				c_status := 253;
				c_ready := FALSE;
			ELSE
				nErrId := fbReceive.nErrId;
			END_IF;
			state := C_WAIT;
		END_IF;
	END_IF;

	IF bytes_received > 0 THEN
		r_time := tx; (* reset receive timer *)
		R_BUF.SIZE := INT_TO_UINT(r_offset + DINT_TO_UINT(bytes_received)); (* enter current buffer size *)
		IP_C.MAILBOX[1] := IP_C.MAILBOX[1] + 1; (* Receive Info *)
		IF IP_C.MAILBOX[1] = 0 THEN IP_C.MAILBOX[1] := 1; END_IF;
    		(*####################  DEBUG-MESSAGE  ###################################*)
		(*########################################################################*)
	END_IF;

(* ---------------------------------------------- *)
C_SND:
	IF S_BUF.SIZE > 0 THEN
		IF c_ready AND c_enable THEN
			IF NOT s_active THEN
				s_total := LIMIT(0,UINT_TO_INT(S_BUF.SIZE),r_max_size); (* limit and pass total number of bytes *)
				s_cur_pos  := 0;
				s_cur_size := 0;
				s_active := TRUE; (* activate send mode *)
			END_IF;

			s_cur_pos := s_cur_pos + s_cur_size; (* determine new send data position *)
			IF s_total > s_cur_pos THEN (* more data to send? *)
				s_cur_size := s_total - s_cur_pos;
				s_cur_size := LIMIT(0, s_cur_size, s_max_size); (* limit telegram length *)
				s_time := tx; (* reset send timer *)
				r_time := tx; (* reset receive timer *)
				state := C_SND_WAIT1;
			ELSE (* no more send data -> stop sending *)
				s_active := FALSE;
				S_BUF.SIZE := 0;
				state := C_WAIT;
			END_IF;
		END_IF;
	END_IF;

C_SND_WAIT1:
	IF udp_mode THEN
		fbSendTo(bExecute := FALSE );
		fbSendTo(bExecute := TRUE, sSrvNetId := sSrvNetId, hSocket := socket , cbLen := s_cur_size, pSrc := REF(S_BUF.BUFFER[s_cur_pos]), sRemoteHost := s_ip_str, nRemotePort := c_port, tTimeout:= TIME_OUT);
		fbSendTo.bExecute := FALSE;
	ELSE
		fbSend(bExecute := FALSE);
		fbSend(bExecute := TRUE, sSrvNetId := sSrvNetId, hSocket := socket , cbLen:= s_cur_size, pSrc:=REF(S_BUF.BUFFER[s_cur_pos]), tTimeout:= TIME_OUT);
		fbSend.bExecute := FALSE;
	END_IF;
	state := C_SND_WAIT2;

C_SND_WAIT2:
	IF udp_mode THEN
		fbSendTo();
		Send_Busy := fbSendTo.bBusy;
		Send_Error := fbSendTo.bError;
		Send_ErrId := fbSendTo.nErrId;
	ELSE
		fbSend();
		Send_Busy := fbSend.bBusy;
		Send_Error := fbSend.bError;
		Send_ErrId := fbSend.nErrId;
	END_IF;

	IF NOT Send_Busy THEN
		IF NOT Send_Error THEN
    			(*####################  DEBUG-MESSAGE  ###################################*)
			(*########################################################################*)
			state := C_SND;
		ELSE
			nErrId := Send_ErrId;
			state := C_WAIT;
    			(*####################  DEBUG-MESSAGE  ###################################*)
			(*########################################################################*)
		END_IF;
	END_IF;

END_CASE;

(* ---------------------------------------------- *)

IF c_ready THEN
	IF NOT c_ready_old THEN (* connection established - rising edge detected *)
		IP_C.C_STATE := 254;
		IP_C.TIME_RESET := TRUE; (* reset all timers *)
    		(*####################  DEBUG-MESSAGE  ###################################*)
		(*########################################################################*)
	ELSE
		IP_C.C_STATE := 255; (* connection is established *)
	END_IF;
ELSE
	IF c_ready_old THEN (* connection closed - falling edge detected *)
		IP_C.C_STATE := 1;
		IP_C.MAILBOX[1] := 0; (* reset receive info *)
		IP_C.MAILBOX[2] := 0; (* reset send pause *)
		IP_C.MAILBOX[3] := 0; (* reset receive pause *)
		IF NOT new_connection THEN
			s_active := FALSE; (* stop send activities *)
			S_BUF.SIZE := 0; (* clear send length *)
			(*IP_C.TIME_RESET := TRUE; *) (* reset all timers *)
		END_IF;
    		(*####################  DEBUG-MESSAGE  ###################################*)
		(*########################################################################*)

	ELSE
		IP_C.C_STATE := 0; (* connection is closed *)
	END_IF;
END_IF;

c_ready_old := c_ready;

(* ---------------------------------------------- *)

IF IP_C.R_OBSERVE = FALSE OR s_active THEN
    r_time := tx;
END_IF;

IF c_status = 0 AND tx - c_time > error_time AND c_enable AND NOT c_ready AND c_mode < 2 THEN (* Connect Timeout *)
	c_status := 255; (* output timeout error *)
END_IF;
IF s_status = 0 AND tx - s_time > error_time AND s_active THEN (* Send Timeout *)
	s_status := 255; (* output timeout error *)
END_IF;
IF r_status = 0 AND tx - r_time > error_time AND c_ready AND IP_C.R_OBSERVE THEN (* Receive Timeout *)
	r_status := 255; (* output timeout error *)
END_IF;

IP_C.ERROR := DWORD_OF_BYTE(c_status,s_status,r_status,0);

(*####################  DEBUG-MESSAGE  ###################################*)
(*########################################################################*)
END_FUNCTION_BLOCK

FUNCTION_BLOCK IP_CONTROL2
VAR_IN_OUT
	IP_C : IP_C; (* IP_Control management structure *)
	S_BUF : NETWORK_BUFFER_SHORT; (* external send buffer *)
	R_BUF : NETWORK_BUFFER_SHORT; (* external receive buffer *)
END_VAR
VAR_INPUT
	IP : DWORD; (* IP address *)
	PORT : WORD; (* port number *)
	TIME_OUT : TIME; (* monitoring time connect / send / receive *)
END_VAR
VAR
	c_time : UDINT; (* CONNECT time (time without a connection) *)
	s_time : UDINT; (* SEND time    (time since last data sent) *)
	r_time : UDINT; (* RECEIVE time (time since last data received) *)

	c_enable : BOOL; (* Connect: enable *)
	c_ip : DWORD; (* Connect: current IP address *)
	c_ip_str : STRING(15);
	c_port : WORD; (* Connect: current port number *)
	c_mode : BYTE; (* Connect: current mode UDP/TCP...... *)
	c_status : BYTE; (* Connect: Status (Errorcode) *)
	c_ready : BOOL; (* Connect: connection established *)
	c_ready_old : BOOL; (* Connect: last state of c_ready for edge evaluation *)

	s_total : INT; (* Send: total number of bytes to send *)
	s_cur_pos : INT; (* Send: current send position *)
	s_cur_size : INT; (* Send: current send size *)
	s_max_size : INT; (* Send: maximum telegram length for sending *)
	s_active : BOOL; (* Send: data sending is active *)
	s_status : BYTE; (* Send: Status (Errorcodes) *)
	s_ip_str : STRING(15);
	s_ip_exist : BOOL;

	r_status : BYTE; (* Receive: Status (Errorcodes) *)
	r_count : INT; (* Receive: number of received bytes *)
	r_offset : INT; (* Receive: current position in external receive buffer *)
	r_max_size : INT; (* Receive: size of the external receive buffer*)

	new_connection : BOOL; (* new connection parameters available *)
	tx : DWORD; (* last time value *)
	error_time : DWORD; (* current timeout value *)

	state: BYTE; (* step flag for connect *)
	bytes_received : DINT;
	udp_mode : BOOL;

	socket : T_HSOCKET;
	server : T_HSOCKET;
	fbCreate	: FB_SocketUdpCreate;
	fbConnect : FB_SocketConnect;
	fbClose : FB_SocketClose;
	fbSendTo : FB_SocketUdpSendTo;
	fbReceiveFrom : FB_SocketUdpReceiveFrom;
	fbSend : FB_SocketSend;
	fbReceive : FB_SocketReceive;
	fbListen : FB_SocketListen;
	fbAccept : FB_SocketAccept;
	nErrId : UDINT;
	Send_Busy : BOOL;
	Send_Error : BOOL;
	Send_ErrId : UDINT;

	IPCR : IP_CONTROL_RESET;

END_VAR

VAR CONSTANT
	STOP : BYTE := 0;
	UDP_INIT : BYTE:= 31;
	UDP_INIT_WAIT : BYTE:= 32;
	TC_INIT : BYTE := 51;
	TC_CONNECT : BYTE:= 52;
	TS_INIT : BYTE:= 53;
	TS_LISTEN_WAIT : BYTE:= 54;
	TS_ACCEPT : BYTE:= 55;
	TS_ACCEPT_WAIT : BYTE:= 56;
	C_CLOSE : BYTE := 190;
	C_CLOSE_WAIT : BYTE:= 191;
	C_CLOSE2 : BYTE := 192;
	C_CLOSE_WAIT2 : BYTE:= 193;
	C_WAIT : BYTE := 200;
	C_REC : BYTE := 210;
	C_REC_WAIT : BYTE := 211;
	C_SND : BYTE := 220;
	C_SND_WAIT1 : BYTE := 221;
	C_SND_WAIT2 : BYTE := 222;
END_VAR

(* on PLC reset or program download all old connections are closed *)
IPCR();
IF NOT IPCR.READY THEN RETURN; END_IF;

(* ---------------------- Debug-Message ----------------------------*)
(* -----------------------------------------------------------------------------*)

(* read system time *)
tx := T_PLC_MS();

IF IP_C.C_PORT = 0 THEN
	IP_C.C_PORT := PORT; (* replace the given PORT with the configured PORT *)
END_IF;
IF IP_C.C_IP = 0 AND IP_C.C_MODE < 4 THEN
	IP_C.C_IP := IP; (* replace the given IP with the configured IP *)
END_IF;

new_connection := c_ip <> IP_C.C_IP OR c_mode <> IP_C.C_MODE OR c_port <> IP_C.C_PORT; (* new connection requested? *)
c_enable := IP_C.C_ENABLE AND NOT new_connection;

IF IP_C.C_ENABLE AND IP_C.ERROR = 0 AND state = STOP THEN (* enabled, no error, no socket open *)
	c_ip   := IP_C.C_IP;
	c_mode := IP_C.C_MODE;
	c_port := IP_C.C_PORT;
	c_ip_str := IP4_TO_STRING(c_ip);
	udp_mode := C_MODE.0; (* Mode 1,3,5 = udp *)

	error_time := TIME_TO_DWORD(MAX(T#200ms,TIME_OUT));
	IP_C.TIME_RESET := TRUE; (* reset all timers *)
	r_max_size := SIZEOF(R_BUF.BUFFER);
	s_max_size := SIZEOF(S_BUF.BUFFER);

	IF udp_mode THEN (* Mode: 1,3,5 = UDP CLIENT/SERVER *)
		state := UDP_INIT;
	ELSIF c_mode = 0 THEN (* Mode: 0 = TCP CLIENT *)
		state := TC_INIT;
	ELSE
		state := TS_INIT; (* Mode: 2,4 = TCP SERVER *)
	END_IF;

	IF c_mode < 5 THEN (* set receiver IP *)
		s_ip_str := c_ip_str;
		s_ip_exist := TRUE;
	ELSE
		s_ip_str := '';
		s_ip_exist := FALSE;
	END_IF;

	(* ---------------------- Debug-Message ----------------------------*)
	(* -----------------------------------------------------------------------------*)

END_IF;

(* ---------------------- Debug-Message ----------------------------*)
(* -----------------------------------------------------------------------------*)

IF IP_C.TIME_RESET THEN (* reset all timers *)
	IP_C.TIME_RESET := FALSE; (* clear reset *)
	c_time := tx; (* reset connect timer *)
	s_time := tx; (* reset send timer *)
	r_time := tx; (* reset receive timer *)

	c_status := BYTE#0; (* reset last connect error *)
	s_status := BYTE#0; (* reset last send error *)
	r_status := BYTE#0; (* reset last receive error *)
END_IF;

(* ---------------------------------------------- *)

CASE state OF

UDP_INIT: (* UDP_CLIENT + SERVER *)

	fbCreate(bExecute := FALSE);
	fbCreate(	bExecute := TRUE , sSrvNetId := sSrvNetId, sLocalHost := sLocalHost, nLocalPort := SEL(c_mode = 1,c_port,0) , tTimeOut := TIME_OUT);
	fbCreate.bExecute := FALSE; (* presetting *)
	state := UDP_INIT_WAIT;

UDP_INIT_WAIT:
	fbCreate();
	IF NOT fbCreate.bBusy THEN
		IF NOT fbCreate.bError THEN
			socket := fbCreate.hSocket;
			c_ready := TRUE; (* Connected *)
			state := C_WAIT;
		ELSE
			c_status := 1; (* FB_SocketUdpCreate failed *)
			nErrId := fbCreate.nErrId;
			state := C_CLOSE;
		END_IF;
	END_IF;

(* ---------------------------------------------- *)
TC_INIT: (* TCP CLIENT (ACTIVE) *)
	fbConnect(bExecute := FALSE);
	fbConnect(bExecute := TRUE, sSrvNetId := sSrvNetId, sRemoteHost := c_ip_str, nRemotePort := c_port, tTimeout := TIME_OUT);
	fbConnect.bExecute := FALSE; (* presetting *)
	state := TC_CONNECT;

(* ---------------------------------------------- *)
TC_CONNECT:
	IF NOT c_enable OR c_status = 255 THEN (* not enabled or connect timeout *)
		state := C_CLOSE;
	END_IF;

	fbConnect();
	IF NOT fbConnect.bBusy THEN
		IF fbConnect.bError THEN
			nErrID := fbConnect.nErrId;
			state :=TC_INIT; (* try again *)
		ELSE
			socket := fbConnect.hSocket;
			c_ready := TRUE; (* Connected *)
			state := C_WAIT;
		END_IF;
	END_IF;
(* ---------------------------------------------- *)
TS_INIT: (* TCP SERVER (PASSIVE) *)
	fbListen(bExecute := FALSE);
	fbListen(bExecute := TRUE , sSrvNetId:= sSrvNetId , sLocalHost:= sLocalHost , nLocalPort:= c_port , tTimeout:= TIME_OUT);
	fbListen.bExecute := FALSE;
	state := TS_LISTEN_WAIT;
(* ---------------------------------------------- *)
TS_LISTEN_WAIT:
	fbListen();
	IF NOT fbListen.bBusy THEN
		IF fbListen.bError THEN
			c_status := 3; (* FB_SocketListen failed *)
			nErrID := fbListen.nErrId;
			state := C_CLOSE;
		ELSE
			server := fbListen.hListener;
			state := TS_ACCEPT;
		END_IF;
	END_IF;
(* ---------------------------------------------- *)
TS_ACCEPT:
	fbAccept(bExecute := FALSE);
	fbAccept(bExecute := TRUE, sSrvNetId := sSrvNetId, hListener := server, tTimeout := TIME_OUT);
	fbAccept.bExecute := FALSE;
	state := TS_ACCEPT_WAIT;
(* ---------------------------------------------- *)
TS_ACCEPT_WAIT:
	fbAccept();
	IF NOT fbAccept.bBusy THEN

		state := TS_ACCEPT; (* presetting  *)

		IF fbAccept.bError THEN
			c_status := 4; (* FB_SocketAccept failed *)
			nErrID := fbAccept.nErrId;
			state := C_CLOSE;
		ELSE (* no error *)
			IF fbAccept.bAccepted THEN
				socket := fbAccept.hSocket;
				IF socket.remoteAddr.sAddr  <> c_ip_str AND c_mode <> 4 THEN
					state := C_CLOSE; (* reject remote - close connection *)
				ELSE
					state := C_WAIT;
					c_ready := TRUE; (* Connected *)
				END_IF;
			END_IF; (* fbAccept.bAccepted *)
		END_IF; (* NOT fbAccept.bError *)
	END_IF; (*  NOT fbAccept.bBusy *)
(* ---------------------------------------------- *)
C_CLOSE:
	c_ready := FALSE;
	IF socket.handle > 0 THEN
		fbClose(bExecute := FALSE );
		fbClose(bExecute := TRUE, sSrvNetId := sSrvNetId,	hSocket := socket);
		fbClose.bExecute := FALSE;
		state := C_CLOSE_WAIT;
	ELSE
		state:= C_CLOSE2;
	END_IF;

C_CLOSE_WAIT:
	fbClose();
	IF NOT fbClose.bBusy THEN
		socket := TCPADS_NULL_HSOCKET; (* reset socket *)
		nErrId:= SEL(fbClose.bError,0,fbClose.nErrId);
		state:= C_CLOSE2;
	END_IF;

C_CLOSE2:
	IF server.handle > 0 THEN
		fbClose(bExecute := FALSE );
		fbClose(bExecute := TRUE, sSrvNetId := sSrvNetId,	hSocket := server);
		fbClose.bExecute := FALSE;
		state := C_CLOSE_WAIT2;
	ELSE
		state:= STOP;
	END_IF;

C_CLOSE_WAIT2:
	fbClose();
	IF NOT fbClose.bBusy THEN
		server := TCPADS_NULL_HSOCKET; (* reset socket *)
		nErrId:= SEL(fbClose.bError,0,fbClose.nErrId);
		state:= STOP;
	END_IF;

(* ---------------------------------------------- *)

C_WAIT: (* wait for request to close *)
	IF NOT c_enable OR NOT c_ready THEN
		state:= C_CLOSE;
	ELSIF S_BUF.SIZE > 0 AND IP_C.MAILBOX[2] = 0 AND s_ip_exist THEN
		state:= C_SND;
	ELSIF IP_C.MAILBOX[3] = 0 THEN
		state:= C_REC;
	END_IF;

(* ---------------------------------------------- *)

C_REC: (* receive data *)
	IF R_BUF.SIZE >= INT_TO_UINT(r_max_size) THEN
		R_BUF.SIZE := UINT#0;
		r_status := 254; (* receive buffer reset due to overflow *)
	END_IF;

	r_offset := SEL(udp_mode,UINT_TO_INT(R_BUF.SIZE),0);
	r_count := r_max_size - r_offset;
	bytes_received := 0;

	IF udp_mode THEN
		fbReceiveFrom(bExecute := FALSE );
		fbReceiveFrom(bExecute:= TRUE, sSrvNetId := sSrvNetId, hSocket:= socket, cbLen := r_count, pDest:= REF(R_BUF.BUFFER), tTimeout:= TIME_OUT);
		fbReceiveFrom.bExecute := FALSE; (* presetting *)
	ELSE
		fbReceive(bExecute:= FALSE);
		fbReceive(bExecute:= TRUE, sSrvNetId := sSrvNetId , hSocket:= socket , cbLen:= r_count, pDest:= REF(R_BUF.BUFFER[r_offset]) , tTimeout:= TIME_OUT);
		fbReceive.bExecute:= FALSE;
	END_IF;
	state := C_REC_WAIT;

C_REC_WAIT:
	IF udp_mode THEN
		fbReceiveFrom(); (* UDP Mode *)
		IF NOT fbReceiveFrom.bBusy THEN
			IF NOT fbReceiveFrom.bError THEN
				IF fbReceiveFrom.nRecBytes > 0 THEN
					IF c_mode <> BYTE#5 AND fbReceiveFrom.sRemoteHost  <> c_ip_str THEN
						R_BUF.SIZE := UINT#0;
					ELSE
						bytes_received := fbReceiveFrom.nRecBytes;
						IF NOT s_ip_exist THEN
							s_ip_str := fbReceiveFrom.sRemoteHost;
							s_ip_exist := TRUE;
						END_IF;
					END_IF;
				END_IF;
				state := C_WAIT;
			ELSIF fbReceiveFrom.nErrId = 16#80072746 THEN (* connection closed by remote *)
				c_status := 253;
				c_ready := FALSE;
				state := C_WAIT;
			ELSE
				nErrId := fbReceiveFrom.nErrId;
				state := C_WAIT;
			END_IF;
		END_IF;
	ELSE
		fbReceive(); (* TCP Mode *)
		IF NOT fbReceive.bBusy THEN
			IF NOT fbReceive.bError THEN
				IF fbReceive.nRecBytes > 0 THEN
					bytes_received := fbReceive.nRecBytes;
				END_IF;
			ELSIF fbReceive.nErrId = 16#8004 THEN (* connection closed by remote *)
				c_status := 253;
				c_ready := FALSE;
			ELSE
				nErrId := fbReceive.nErrId;
			END_IF;
			state := C_WAIT;
		END_IF;
	END_IF;

	IF bytes_received > 0 THEN
		r_time := tx; (* reset receive timer *)
		R_BUF.SIZE := INT_TO_UINT(r_offset + DINT_TO_UINT(bytes_received)); (* enter current buffer size *)
		IP_C.MAILBOX[1] := IP_C.MAILBOX[1] + 1; (* Receive Info *)
		IF IP_C.MAILBOX[1] = 0 THEN IP_C.MAILBOX[1] := 1; END_IF;
		(* ---------------------- Debug-Message ----------------------------*)
		(* -----------------------------------------------------------------------------*)
	END_IF;

(* ---------------------------------------------- *)
C_SND:
	IF S_BUF.SIZE > 0 THEN
		IF c_ready AND c_enable THEN
			IF NOT s_active THEN
				s_total := LIMIT(0,UINT_TO_INT(S_BUF.SIZE),r_max_size); (* limit and pass total number of bytes *)
				s_cur_pos  := 0;
				s_cur_size := 0;
				s_active := TRUE; (* activate send mode *)
			END_IF;

			s_cur_pos := s_cur_pos + s_cur_size; (* determine new send data position *)
			IF s_total > s_cur_pos THEN (* more data to send? *)
				s_cur_size := s_total - s_cur_pos;
				s_cur_size := LIMIT(0, s_cur_size, s_max_size); (* limit telegram length *)
				s_time := tx; (* reset send timer *)
				r_time := tx; (* reset receive timer *)
				state := C_SND_WAIT1;
			ELSE (* no more send data -> stop sending *)
				s_active := FALSE;
				S_BUF.SIZE := 0;
				state := C_WAIT;
			END_IF;
		END_IF;
	END_IF;

C_SND_WAIT1:
	IF udp_mode THEN
		fbSendTo(bExecute := FALSE );
		fbSendTo(bExecute := TRUE, sSrvNetId := sSrvNetId, hSocket := socket , cbLen := s_cur_size, pSrc := REF(S_BUF.BUFFER[s_cur_pos]), sRemoteHost := s_ip_str, nRemotePort := c_port, tTimeout:= TIME_OUT);
		fbSendTo.bExecute := FALSE;
	ELSE
		fbSend(bExecute := FALSE);
		fbSend(bExecute := TRUE, sSrvNetId := sSrvNetId, hSocket := socket , cbLen:= s_cur_size, pSrc:=REF(S_BUF.BUFFER[s_cur_pos]), tTimeout:= TIME_OUT);
		fbSend.bExecute := FALSE;
	END_IF;
	state := C_SND_WAIT2;

C_SND_WAIT2:
	IF udp_mode THEN
		fbSendTo();
		Send_Busy := fbSendTo.bBusy;
		Send_Error := fbSendTo.bError;
		Send_ErrId := fbSendTo.nErrId;
	ELSE
		fbSend();
		Send_Busy := fbSend.bBusy;
		Send_Error := fbSend.bError;
		Send_ErrId := fbSend.nErrId;
	END_IF;

	IF NOT Send_Busy THEN
		IF NOT Send_Error THEN
			(* ---------------------- Debug-Message ----------------------------*)
			(* -----------------------------------------------------------------------------*)
			state := C_SND;
		ELSE
			nErrId := Send_ErrId;
			state := C_WAIT;
			(* ---------------------- Debug-Message ----------------------------*)
			(* -----------------------------------------------------------------------------*)
		END_IF;
	END_IF;

END_CASE;

(* ---------------------------------------------- *)

IF c_ready THEN
	IF NOT c_ready_old THEN (* connection established - rising edge detected *)
		IP_C.C_STATE := 254;
		IP_C.TIME_RESET := TRUE; (* reset all timers *)
		(* ---------------------- Debug-Message ----------------------------*)
		(* -----------------------------------------------------------------------------*)
	ELSE
		IP_C.C_STATE := 255; (* connection is established *)
	END_IF;
ELSE
	IF c_ready_old THEN (* connection closed - falling edge detected *)
		IP_C.C_STATE := 1;
		IP_C.MAILBOX[1] := 0; (* reset receive info *)
		IP_C.MAILBOX[2] := 0; (* reset send pause *)
		IP_C.MAILBOX[3] := 0; (* reset receive pause *)
		IF NOT new_connection THEN
			s_active := FALSE; (* stop send activities *)
			S_BUF.SIZE := 0; (* clear send length *)
			(* IP_C.TIME_RESET := TRUE; *) (* reset all timers *)
		END_IF;
		(* ---------------------- Debug-Message ----------------------------*)
		(* -----------------------------------------------------------------------------*)

	ELSE
		IP_C.C_STATE := 0; (* connection is closed *)
	END_IF;
END_IF;

c_ready_old := c_ready;

(* ---------------------------------------------- *)

IF IP_C.R_OBSERVE = FALSE OR s_active THEN
    r_time := tx;
END_IF;

IF c_status = 0 AND tx - c_time > error_time AND c_enable AND NOT c_ready AND c_mode < 2 THEN (* Connect Timeout *)
	c_status := 255; (* output timeout error *)
END_IF;
IF s_status = 0 AND tx - s_time > error_time AND s_active THEN (* Send Timeout *)
	s_status := 255; (* output timeout error *)
END_IF;
IF r_status = 0 AND tx - r_time > error_time AND c_ready AND IP_C.R_OBSERVE THEN (* Receive Timeout *)
	r_status := 255; (* output timeout error *)
END_IF;

IP_C.ERROR := DWORD_OF_BYTE(c_status,s_status,r_status,0);

(* ---------------------- Debug-Message ----------------------------*)
(* -----------------------------------------------------------------------------*)
END_FUNCTION_BLOCK

FUNCTION_BLOCK IP_CONTROL_RESET
VAR_INPUT
END_VAR
VAR_OUTPUT
	READY : BOOL;
END_VAR
VAR
	fbSocketCloseAll : FB_SocketCloseAll;
	step : BYTE;
END_VAR
(* On PLC reset or program download close all old connections *)

CASE TCP_SERVER_RESET OF

2:	READY := TRUE;
	RETURN;

1:	IF step = 0 THEN
		RETURN;
	END_IF;

0:	TCP_SERVER_RESET := 1; (* reset is active *)
	step := 1;

END_CASE;

CASE step OF

0:	RETURN;

1:	fbSocketCloseAll(sSrvNetId:= '', bExecute:= TRUE, tTimeout:= T#10s );
	step := 2;

2:	fbSocketCloseAll(bExecute:= FALSE);
	IF fbSocketCloseAll.bBusy = FALSE THEN
		step := 0;
		TCP_SERVER_RESET := 2; (* reset done, ready for use *)
	END_IF;

END_CASE;
END_FUNCTION_BLOCK

FUNCTION_BLOCK IP_FIFO
VAR_IN_OUT
	FIFO : IP_FIFO_DATA;
	ID : BYTE;
	STATE : BYTE;
END_VAR
VAR
	tmp : INT;

END_VAR

IF NOT FIFO.INIT THEN
	FIFO.INIT := TRUE;
	FIFO.NW := 1;
	FIFO.NR := 1;
	FIFO.EMPTY := TRUE;
	FIFO.FULL := FALSE;
	FIFO.TOP := 128; (* array size 1 to x *)
	FIFO.MAX_ID := BYTE#1; (* maximum number of registrations of an ID *)
END_IF;

(* if no ID is present, assign the next free ID *)
IF ID = BYTE#00 THEN
	IF FIFO.ID < INT_TO_BYTE(FIFO.TOP) THEN
		ID := USINT_TO_BYTE(BYTE_TO_USINT(FIFO.ID) + USINT#1); (* ID +1 *)
		FIFO.ID := ID; (* enter ID *)
	    	(*####################  DEBUG-MESSAGE  ###################################*)
	    	(*########################################################################*)
	ELSIF STATE < BYTE#200 THEN
		(*####################  DEBUG-MESSAGE  ###################################*)
		(*########################################################################*)

		(* should never happen! -> too many participants are registered *)
		(* no number is assigned and all requests are answered negatively *)
		STATE := BYTE#255;
		RETURN;
	END_IF;
END_IF;

IF STATE = BYTE#1 AND NOT FIFO.FULL THEN (* register *)
	(* enter only if the maximum number was not exceeded *)
	tmp := BYTE_TO_INT(ID);
	IF FIFO.Y[tmp] < FIFO.MAX_ID THEN
		FIFO.Y[tmp] := USINT_TO_BYTE(BYTE_TO_USINT(FIFO.Y[tmp]) + USINT#1); (* IP_ID counter +1 *)
		tmp := FIFO.NW; (* unfortunately no other way works for me *)
		FIFO.X[tmp] := ID;
		IF FIFO.NW = FIFO.TOP THEN FIFO.NW := 1; ELSE FIFO.NW := FIFO.NW + 1; END_IF;
		FIFO.FULL := FIFO.NW = FIFO.NR;
		FIFO.EMPTY := FALSE;
		STATE := BYTE#2; (* advance to "wait for release" *)
		(*####################  DEBUG-MESSAGE  ###################################*)
		(*########################################################################*)

	END_IF;
END_IF;

IF STATE = BYTE#2 AND NOT FIFO.EMPTY THEN (* check current stack element *)
	tmp := FIFO.NR; (* unfortunately no other way works for me *)
	IF ID = FIFO.X[tmp] THEN
	STATE := BYTE#3; (* advance to "access release" *)
	END_IF;
END_IF;

IF STATE = BYTE#4 AND NOT FIFO.EMPTY THEN (* remove from stack *)
	tmp := FIFO.NR; (* unfortunately no other way works for me *)
	IF ID = FIFO.X[tmp] THEN
		tmp := BYTE_TO_INT(ID); (* unfortunately no other way works for me *)
		FIFO.Y[tmp] := USINT_TO_BYTE(BYTE_TO_USINT(FIFO.Y[tmp]) - USINT#1); (* ID counter -1 *)
		IF FIFO.NR = FIFO.TOP THEN FIFO.NR := 1; ELSE FIFO.NR := FIFO.NR + 1; END_IF;
		FIFO.EMPTY := FIFO.NR = FIFO.NW;
		FIFO.FULL := FALSE;
		(*####################  DEBUG-MESSAGE  ###################################*)
		(*########################################################################*)
	END_IF;
	STATE := BYTE#5; (* "Remove done" *)
END_IF;
END_FUNCTION_BLOCK

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION_BLOCK IRTRANS_DECODE
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	R_BUF : NETWORK_BUFFER_SHORT;
// END_VAR
// VAR_OUTPUT
// 	CMD : BOOL;
// 	DEV : STRING;
// 	KEY : STRING;
// 	ERROR : BOOL;
// END_VAR
// VAR
// 	i : INT;
// 	stop : INT;
// 	pt : REF_TO BYTE;
// 	buffer : REF_TO ARRAY[0..32000] OF BYTE;
// 	size : UINT;
// 	z : BYTE;
// END_VAR
//
// size := R_BUF.SIZE;
//
// (* do nothing if the buffer is empty *)
// IF size > 0 AND IP_C.MAILBOX[1] <> z THEN
// 	i := 0;
// 	error := TRUE;
// 	stop := size - 2;
// 	(* read characters before first comma and store in device *)
// 	pt := REF(dev);
// 	buffer := REF(R_BUF.BUFFER);
// 	WHILE i < stop DO
// 		IF buffer^[i] = 44 THEN EXIT; END_IF;
// 		pt^ := buffer^[i];
// 		pt := pt + 1;
// 		i := i + 1;
// 	END_WHILE;
// 	(* terminate the device string *)
// 	pt^ := 0;
// 	(* read after first comma and store in key *)
// 	i := i + 1;
// 	pt := REF(key);
// 	WHILE i < stop DO
// 		pt^ := buffer^[i];
// 		pt := pt + 1;
// 		i := i + 1;
// 	END_WHILE;
// 	(* terminate the key string *)
// 	pt^:= 0;
// 	(* check if last two characters are $R$N 13, 10 *)
// 	IF buffer^[i] <> 13 AND buffer^[i+1] <> 10 THEN	RETURN;	END_IF;
// 	(* check if device and key are present *)
// 	IF LEN(dev) = 0 OR LEN(key) = 0 THEN RETURN; END_IF;
// 	(* a valid transmission has been received *)
// 	cmd := TRUE;
// 	error := FALSE;
// ELSE
// 	error := FALSE;
// 	cmd := FALSE;
// END_IF;
//
// z := IP_C.MAILBOX[1];
// END_FUNCTION_BLOCK

FUNCTION_BLOCK IRTRANS_RCV_1
VAR_INPUT CONSTANT
	DEV_CODE : STRING;
	KEY_CODE : STRING;
END_VAR
VAR_INPUT
	CMD : BOOL;
END_VAR
VAR_IN_OUT
	DEV : STRING;
	KEY : STRING;
END_VAR
VAR_OUTPUT
	Q : BOOL;
END_VAR

Q := cmd AND dev_code = dev AND key_code = key;
END_FUNCTION_BLOCK

FUNCTION_BLOCK IRTRANS_RCV_4
VAR_INPUT CONSTANT
	DEV_CODE : STRING;
	KEY_CODE_0 : STRING;
	KEY_CODE_1 : STRING;
	KEY_CODE_2 : STRING;
	KEY_CODE_3 : STRING;
END_VAR
VAR_INPUT
	CMD : BOOL;
END_VAR
VAR_IN_OUT
	DEV : STRING;
	KEY : STRING;
END_VAR
VAR_OUTPUT
	Q0, Q1, Q2, Q3 : BOOL;
END_VAR
VAR_TEMP
	decode: BOOL;
END_VAR

decode := cmd AND dev_code = dev;

Q0 := key = key_code_0 AND decode;
Q1 := key = key_code_1 AND decode;
Q2 := key = key_code_2 AND decode;
Q3 := key = key_code_3 AND decode;
END_FUNCTION_BLOCK

FUNCTION_BLOCK IRTRANS_RCV_8
VAR_INPUT CONSTANT
	dev_code : STRING;
	key_code_0 : STRING;
	key_code_1 : STRING;
	key_code_2 : STRING;
	key_code_3 : STRING;
	key_code_4 : STRING;
	key_code_5 : STRING;
	key_code_6 : STRING;
	key_code_7 : STRING;
END_VAR
VAR_INPUT
	CMD : BOOL;
END_VAR
VAR_IN_OUT
	DEV : STRING;
	KEY : STRING;
END_VAR
VAR_OUTPUT
	Q0, Q1, Q2, Q3, Q4, Q5, Q6, Q7 : BOOL;
END_VAR
VAR_TEMP
	decode: BOOL;
END_VAR

decode := cmd AND dev_code = dev;

Q0 := key = key_code_0 AND decode;
Q1 := key = key_code_1 AND decode;
Q2 := key = key_code_2 AND decode;
Q3 := key = key_code_3 AND decode;
Q4 := key = key_code_4 AND decode;
Q5 := key = key_code_5 AND decode;
Q6 := key = key_code_6 AND decode;
Q7 := key = key_code_7 AND decode;
END_FUNCTION_BLOCK

FUNCTION_BLOCK IRTRANS_SERVER
VAR_IN_OUT
	IP_C : IP_C; (* IP_Control management structure *)
	S_BUF : NETWORK_BUFFER_SHORT;
	R_BUF : NETWORK_BUFFER_SHORT;
END_VAR
VAR_INPUT
	UDP_TCP : BOOL;
END_VAR
VAR_OUTPUT
	S_ENABLE : BOOL;
	R_ENABLE : BOOL;
	ERROR : DWORD;
END_VAR
VAR
	t : TON;
END_VAR

IF NOT IP_C.C_ENABLE THEN
	(* set up IP data traffic *)
	IP_C.C_PORT := WORD#0; (* enter port number *)
	IP_C.C_IP := DWORD#00; (* enter IP *)
	IP_C.C_MODE := SEL(UDP_TCP,5,4); (* 4 = Mode: TCP+PASSIV+PORT / 5 = Mode: UDP+PASSIV+PORT*)
	IP_C.C_ENABLE := TRUE; (* enable connection setup *)
	IP_C.R_OBSERVE := FALSE; (* monitor data reception *)
	IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
	S_BUF.SIZE := UINT#0; (* reset send data *)
	R_BUF.SIZE := UINT#0; (* reset receive data *)
END_IF;

R_ENABLE := IP_C.C_STATE > 127;
S_ENABLE := R_ENABLE AND (IP_C.MAILBOX[1] > 0 OR IP_C.C_MODE <> 5);
ERROR := IP_C.ERROR;

t(IN:= IP_C.ERROR > 0, PT:=T#5s);
IF t.Q THEN
	IP_C.TIME_RESET := TRUE; (* Error Reset *)
END_IF;
END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK IRTRANS_SND_1
// VAR_INPUT CONSTANT
// 	DEV_CODE : STRING;
// 	KEY_CODE : STRING;
// END_VAR
// VAR_INPUT
// 	IN : BOOL;
// 	T_REPEAT : TIME;
// END_VAR
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER_SHORT;
// END_VAR
// VAR_OUTPUT
// 	KEY : BYTE;
// END_VAR
// VAR
// 	skey : STRING;
// 	size : INT;
// 	str : STRING;
// 	t : TON;
// 	k : BYTE;
// 	lk : BYTE;
// 	d : BOOL;
// END_VAR
//
// d := IP_C.C_STATE > 127 AND (IP_C.MAILBOX[1] > 0 OR IP_C.C_MODE <> 5);
// IF NOT d THEN RETURN; END_IF;
//
// t(IN:= lk = k , PT:=T_REPEAT);
//
// IF S_BUF.SIZE = 0 THEN
// 	IF IN THEN
// 		skey := KEY_CODE; k := 1;
// 	ELSE
// 		lk := 0;
// 		KEY := 0;
// 		RETURN; (* exit if no key is active *)
// 	END_IF;
//
// 	IF lk <> k OR T.Q THEN (* other key or repeat wait time elapsed *)
// 		str := CONCAT('snd ',dev_code);
// 		str := CONCAT(str,',');
// 		str := CONCAT(str,skey);
// 		size := LEN(str);
// 		_STRING_TO_BUFFER(STR:=str, POS:=0, PT:=REF(S_BUF.BUFFER), SIZE:=size);
// 		S_BUF.SIZE := size; (* enter send length *)
// 		KEY := k;
// 		lk := k; (* store last key *)
// 		t(IN:=FALSE); (* Reset Time *)
// 	END_IF;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK IRTRANS_SND_4
// VAR_INPUT CONSTANT
// 	DEV_CODE : STRING;
// 	KEY_CODE_0 : STRING;
// 	KEY_CODE_1 : STRING;
// 	KEY_CODE_2 : STRING;
// 	KEY_CODE_3 : STRING;
// END_VAR
// VAR_INPUT
// 	IN_0 : BOOL;
// 	IN_1 : BOOL;
// 	IN_2 : BOOL;
// 	IN_3 : BOOL;
// 	T_REPEAT : TIME;
// END_VAR
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER_SHORT;
// END_VAR
// VAR_OUTPUT
// 	KEY : BYTE;
// END_VAR
// VAR
// 	skey : STRING;
// 	size : INT;
// 	str : STRING;
// 	t : TON;
// 	k : BYTE;
// 	lk : BYTE;
// 	d : BOOL;
// END_VAR
//
// d := IP_C.C_STATE > 127 AND (IP_C.MAILBOX[1] > 0 OR IP_C.C_MODE <> 5);
// IF NOT d THEN RETURN; END_IF;
//
// t(IN:= lk = k , PT:=T_REPEAT);
//
// IF S_BUF.SIZE = 0 THEN
// 	IF IN_0 THEN
// 		skey := KEY_CODE_0; k := 1;
// 	ELSIF IN_1 THEN
// 		skey := KEY_CODE_1; k := 2;
// 	ELSIF IN_2 THEN
// 		skey := KEY_CODE_2; k := 3;
// 	ELSIF IN_3 THEN
// 		skey := KEY_CODE_3; k := 4;
// 	ELSE
// 		lk := 0;
// 		KEY := 0;
// 		RETURN; (* exit if no key is active *)
// 	END_IF;
//
// 	IF lk <> k OR T.Q THEN (* other key or repeat wait time elapsed *)
// 		str := CONCAT('snd ',dev_code);
// 		str := CONCAT(str,',');
// 		str := CONCAT(str,skey);
// 		size := LEN(str);
// 		_STRING_TO_BUFFER(STR:=str, POS:=0, PT:=REF(S_BUF.BUFFER), SIZE:=size);
// 		S_BUF.SIZE := size; (* enter send length *)
// 		KEY := k;
// 		lk := k; (* store last key *)
// 		t(IN:=FALSE); (* Reset Time *)
// 	END_IF;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK IRTRANS_SND_8
// VAR_INPUT CONSTANT
// 	DEV_CODE : STRING;
// 	KEY_CODE_0 : STRING;
// 	KEY_CODE_1 : STRING;
// 	KEY_CODE_2 : STRING;
// 	KEY_CODE_3 : STRING;
// 	KEY_CODE_4 : STRING;
// 	KEY_CODE_5 : STRING;
// 	KEY_CODE_6 : STRING;
// 	KEY_CODE_7 : STRING;
// END_VAR
// VAR_INPUT
// 	IN_0 : BOOL;
// 	IN_1 : BOOL;
// 	IN_2 : BOOL;
// 	IN_3 : BOOL;
// 	IN_4 : BOOL;
// 	IN_5 : BOOL;
// 	IN_6 : BOOL;
// 	IN_7 : BOOL;
// 	T_REPEAT : TIME;
// END_VAR
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER_SHORT;
// END_VAR
// VAR_OUTPUT
// 	KEY : BYTE;
// END_VAR
// VAR
// 	skey : STRING;
// 	size : INT;
// 	str : STRING;
// 	t : TON;
// 	k : BYTE;
// 	lk : BYTE;
// 	d : BOOL;
// END_VAR
//
// d := IP_C.C_STATE > 127 AND (IP_C.MAILBOX[1] > 0 OR IP_C.C_MODE <> 5);
// IF NOT d THEN RETURN; END_IF;
//
// t(IN:= lk = k , PT:=T_REPEAT);
//
// IF S_BUF.SIZE = 0 THEN
// 	IF IN_0 THEN
// 		skey := KEY_CODE_0; k := 1;
// 	ELSIF IN_1 THEN
// 		skey := KEY_CODE_1; k := 2;
// 	ELSIF IN_2 THEN
// 		skey := KEY_CODE_2; k := 3;
// 	ELSIF IN_3 THEN
// 		skey := KEY_CODE_3; k := 4;
// 	ELSIF IN_4 THEN
// 		skey := KEY_CODE_4; k := 5;
// 	ELSIF IN_5 THEN
// 		skey := KEY_CODE_5; k := 6;
// 	ELSIF IN_6 THEN
// 		skey := KEY_CODE_6; k := 7;
// 	ELSIF IN_7 THEN
// 		skey := KEY_CODE_7; k := 8;
// 	ELSE
// 		lk := 0;
// 		KEY := 0;
// 		RETURN; (* exit if no key is active *)
// 	END_IF;
//
// 	IF lk <> k OR T.Q THEN (* other key or repeat wait time elapsed *)
// 		str := CONCAT('snd ',dev_code);
// 		str := CONCAT(str,',');
// 		str := CONCAT(str,skey);
// 		size := LEN(str);
// 		_STRING_TO_BUFFER(STR:=str, POS:=0, PT:=REF(S_BUF.BUFFER), SIZE:=size);
// 		S_BUF.SIZE := size; (* enter send length *)
// 		KEY := k;
// 		lk := k; (* store last key *)
// 		t(IN:=FALSE); (* Reset Time *)
// 	END_IF;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses COUNT_CHAR, which is commented out.
// FUNCTION IS_IP4 : BOOL
// VAR_INPUT
// 	str : STRING;
// END_VAR
//
// IS_IP4 := COUNT_CHAR(str, 46) = 3;
// END_FUNCTION

FUNCTION IS_URLCHR : BOOL
VAR_INPUT
	IN : BYTE;
END_VAR

IS_URLCHR := ((IN > 47) AND (IN < 58)) OR ((IN > 64) AND (IN < 91)) OR ((IN > 96) AND (IN < 123)) OR (in = 126) OR (in = 95) OR (in =45) OR (in =46);
END_FUNCTION

// Commented out for beedance: it uses PRINT_SF, which is commented out.
// FUNCTION_BLOCK LOG_MSG
// (* VAR_EXTERNAL
// 	LOG_CL :	LOG_CONTROL;
// END_VAR *)
// VAR
// 	FB_PRINT_SF :PRINT_SF;
// 	idX : INT;
// END_VAR
//
// IF DWORD_TO_BYTE(SHL(LOG_CL.NEW_MSG_OPTION,16)) > LOG_CL.LEVEL THEN RETURN; END_IF; (* logging enabled?*)
//
// FB_PRINT_SF(PRINTF_DATA := LOG_CL.PRINTF,STR := LOG_CL.NEW_MSG);
//
// IF LOG_CL.RESET THEN (* delete log entries *)
// 	LOG_CL.RESET := FALSE;
// 	LOG_CL.RING_MODE := FALSE;
// 	LOG_CL.IDX := 0;
// 	LOG_CL.UPDATE_COUNT := LOG_CL.UPDATE_COUNT + UINT#1; (* increment UpdateCount *)
// END_IF;
//
// IF LOG_CL.SIZE > 0 THEN (* size value is entered *)
// 	IF LEN(LOG_CL.NEW_MSG) > 0 THEN (* enter message into buffer *)
// 		IF (LOG_CL.IDX >= LOG_CL.SIZE) THEN
// 			LOG_CL.RING_MODE := TRUE; (* on overflow = ring buffer mode *)
// 			LOG_CL.IDX := 0; (* position at start of array again *)
// 		END_IF;
//
// 		LOG_CL.UPDATE_COUNT := LOG_CL.UPDATE_COUNT + UINT#1; (* increment UpdateCount *)
// 		LOG_CL.IDX := LOG_CL.IDX + INT#1; (* next free position *)
// 		idx := LOG_CL.IDX; (* load current log buffer position *)
// 		LOG_CL.MSG[idx] := LOG_CL.NEW_MSG; (* store new log string in buffer  *)
// 		LOG_CL.MSG_OPTION[idx] := LOG_CL.NEW_MSG_OPTION; (* store new log string option in buffer  *)
// 		LOG_CL.NEW_MSG := ''; (* delete last log string *)
// 	END_IF;
// END_IF;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK LOG_VIEWPORT
VAR_IN_OUT
	LC : LOG_CONTROL;
	LV : us_LOG_VIEWPORT;
END_VAR
VAR
	pos : INT;
	count : INT;
	idx : INT;
	base : INT;
	update_pos : BOOL;
END_VAR

(* process scroll button *)
IF LV.MOVE_TO_X <> 0 THEN
	CASE LV.MOVE_TO_X OF
	30000: pos := 1 ; (* show oldest message *)
	30001: pos := LC.SIZE; (* show newest message *)
	30002: pos := pos + LV.COUNT; (* page down -> towards newer messages *)
	30003: pos := pos - LV.COUNT; (* page up -> towards older messages *)
	ELSE
		pos := pos + LV.MOVE_TO_X;
	END_CASE;
	update_pos := TRUE;
	LV.MOVE_TO_X := 0;
END_IF;

IF (LV.UPDATE_COUNT <> LC.UPDATE_COUNT) OR update_pos THEN (* monitor update *)
	LV.UPDATE := TRUE;
	LV.UPDATE_COUNT := LC.UPDATE_COUNT;
	update_pos := FALSE;

	(* check and take pos *)
	IF LC.RING_MODE THEN
		pos := LIMIT(1,pos,LC.SIZE - LV.COUNT + 1);
	ELSE
		IF LC.IDX > LV.COUNT THEN
			pos := LIMIT(1,pos,LC.IDX - LV.COUNT + 1);
		ELSE
			pos := 1;
		END_IF;
	END_IF;

	(* calculate idxBasis *)
	base := LC.IDX - pos - LV.COUNT + 1;
	IF LC.RING_MODE THEN
		(* display in ring buffer mode *)
		FOR count:= 1 TO LV.COUNT DO
			idx := base + count;
			IF idx < 1 THEN
				idx := idx + LC.SIZE;
			END_IF;
			LV.LINE_ARRAY[count] := idx;
		END_FOR;
	ELSE
		(* display in normal mode *)
		FOR count:= 1 TO LV.COUNT DO
			idx := base + count;
			IF idx < 1 THEN
				(* output empty line *)
				LV.LINE_ARRAY[count] := 0;
			ELSE
				LV.LINE_ARRAY[count] := idx;
			END_IF;
		END_FOR;
	END_IF;
END_IF;
END_FUNCTION_BLOCK

FUNCTION_BLOCK MB_CLIENT
VAR_IN_OUT
	IP_C : IP_C;
	S_BUF : NETWORK_BUFFER_SHORT;
	R_BUF : NETWORK_BUFFER_SHORT;
	DATA : ARRAY [0..255] OF WORD;
END_VAR
VAR_INPUT
	DATA_SIZE : INT;
	ENABLE : BOOL;
	UDP : BOOL;
	FC : INT;
	UNIT_ID : BYTE;
	R_ADDR : INT;
	R_POINTS : INT;
	R_DATA_ADR : INT;
	R_DATA_BITPOS : INT;
	W_ADDR : INT;
	W_POINTS : INT;
	W_DATA_ADR : INT;
	W_DATA_BITPOS : INT;
	DELAY : TIME;
END_VAR
VAR_OUTPUT
	ERROR : DWORD;
	BUSY : BOOL;
END_VAR
VAR
	ip_state : BYTE;
	ip_id : BYTE;
	ip_fifo : ip_fifo;
	mask : DWORD;
	w_tmp : WORD;
	idx1 : INT;
	idx2 : INT;
	stop : INT;
	i : INT;
	bit_pos : INT;
	w_mask : WORD;
	b_mask : BYTE;
	count : INT;
	state : INT;
	transaction_id : INT;
	response : INT;
	comp : INT;
	ton_1 : TON;
END_VAR

CASE state OF

00: (* wait for delay time *)
	TON_1(IN:=ENABLE,PT:=DELAY);

	IF TON_1.Q OR DELAY = T#0s THEN
		TON_1(IN:=FALSE);
		state := 20;
		IP_STATE := BYTE#1; (* register *)
	END_IF;

20: (* register with IP_Control *)
	IF IP_STATE = BYTE#3 THEN
		(* set up IP data traffic *)
		IP_C.C_PORT := WORD#0; (* use port from IP_CONTROL *)
		IP_C.C_IP := DWORD#0; (* use IP from IP_CONTROL *)
		IP_C.C_MODE := BOOL_TO_BYTE(UDP); (* Mode: 0=TCP+ACTIVE+IP , 1=UDP+ACTIVE+IP *)
		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
		IP_C.R_OBSERVE := TRUE; (* receive monitoring *)
		ERROR := DWORD#00; (* reset last error *)

		(******************************************************************)
		CASE FC OF
		1..2: (* 0x01 Read Coil Status , 0x02 Read Input Status *)
			IF R_POINTS <= 2000 THEN
				w_tmp := INT_TO_WORD(R_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				w_tmp := INT_TO_WORD(R_POINTS);
				S_BUF.BUFFER[10] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[11] := WORD_TO_BYTE(w_tmp);
				S_BUF.SIZE := UINT#12;
				Response :=	WORD_TO_INT(SHR(INT_TO_WORD(R_POINTS + 7),3));
				comp := 7;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)
		3..4: (* 0x03 Read Holding Registers , 0x04 Read Input Registers *)
			IF R_POINTS <= 125 THEN
				w_tmp := INT_TO_WORD(R_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				w_tmp := INT_TO_WORD(R_POINTS);
				S_BUF.BUFFER[10] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[11] := WORD_TO_BYTE(w_tmp);
				S_BUF.SIZE := UINT#12;
				Response :=	WORD_TO_INT(SHL(w_tmp,1));
				comp := 7;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)
		5:	(* 0x05 Force Single Coil *)
			IF W_DATA_ADR <= DATA_SIZE THEN
				w_tmp := INT_TO_WORD(W_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				S_BUF.BUFFER[10] := SEL((DATA[W_DATA_ADR] AND SHL(WORD#1,W_DATA_BITPOS))>WORD#00,BYTE#16#00,BYTE#16#FF);
				S_BUF.BUFFER[11] := BYTE#00;
				S_BUF.SIZE := UINT#12;
				Response := 3;
				comp := 11;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)
		6:	(* 0x06 Preset Single Register *)
			IF W_DATA_ADR <= DATA_SIZE THEN
				w_tmp := INT_TO_WORD(W_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				w_tmp := DATA[W_DATA_ADR];
				S_BUF.BUFFER[10] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[11] := WORD_TO_BYTE(w_tmp);
				S_BUF.SIZE := UINT#12;
				Response := 3;
				comp := 11;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)
		15:	(* 0x0F Force Multiple Coils *)
			IF (WORD_TO_INT(SHR(INT_TO_WORD(WORD_TO_INT(SHL(INT_TO_WORD(W_DATA_ADR),4)) + W_POINTS + W_DATA_BITPOS + 15),4)) <= DATA_SIZE) AND W_POINTS <= 1968 THEN
				w_tmp := INT_TO_WORD(W_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				w_tmp := INT_TO_WORD(W_POINTS);
				S_BUF.BUFFER[10] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[11] := WORD_TO_BYTE(w_tmp);
				S_BUF.BUFFER[12] := WORD_TO_BYTE(SHR(INT_TO_WORD(W_POINTS + 7),3));

				w_mask := SHL(WORD#1,W_DATA_BITPOS);
				bit_pos := 0;
				w_tmp := WORD#00;
				idx1 := 12;
				idx2 := W_DATA_ADR;
				FOR i := 1 TO W_POINTS DO
					w_tmp :=BIT_LOAD_W(w_tmp,(DATA[idx2] AND w_mask) > WORD#0,bit_pos);
					bit_pos := bit_pos + 1;
					IF bit_pos > 7 OR i=W_POINTS THEN
						bit_pos := 0;
						idx1 := idx1 + 1;
						S_BUF.BUFFER[idx1] := WORD_TO_BYTE(w_tmp);
						w_tmp := WORD#0;
					END_IF;
					w_mask := ROL(w_mask,1);
					IF w_mask = WORD#1 THEN
						idx2 := idx2 +1;
					END_IF;
				END_FOR;
				S_BUF.SIZE := INT_TO_UINT(idx1 + 1);
				Response := 3;
				comp := 11;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)
		16:	(* 0x10 Write Multiple Registers *)
			IF W_DATA_ADR + W_POINTS <= DATA_SIZE AND W_POINTS <= 123 THEN
				w_tmp := INT_TO_WORD(W_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				w_tmp := INT_TO_WORD(W_POINTS);
				S_BUF.BUFFER[10] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[11] := WORD_TO_BYTE(w_tmp);
				S_BUF.BUFFER[12] := SHL(S_BUF.BUFFER[11],1);

				idx1 := 11;
				count := W_DATA_ADR + W_POINTS - 1;
				FOR i := W_DATA_ADR TO count DO
					idx1 := idx1 + 2;
					idx2 := idx1 + 1; (* second index necessary for pcworx! *)
					w_tmp := DATA[i];
					S_BUF.BUFFER[idx1] := WORD_TO_BYTE(SHR(w_tmp,8));
					S_BUF.BUFFER[idx2] := WORD_TO_BYTE(w_tmp);
				END_FOR;

				S_BUF.SIZE := INT_TO_UINT(idx2 + 1);
				Response := 3;
				comp := 11;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)
		22:	(* 0x16 Mask Write Register *)
			IF W_DATA_ADR + 2 <= DATA_SIZE THEN
				w_tmp := INT_TO_WORD(W_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				idx1 := W_DATA_ADR;
				w_tmp := DATA[idx1];
				S_BUF.BUFFER[10] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[11] := WORD_TO_BYTE(w_tmp);
				idx1 := idx1 +1;
				w_tmp := DATA[idx1];
				S_BUF.BUFFER[12] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[13] := WORD_TO_BYTE(w_tmp);

				S_BUF.SIZE := UINT#14;
				Response := 5;
				comp := 13;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)
		23:	(* 0x17 Read/Write Multible Registers *)
			IF W_DATA_ADR + W_POINTS <= DATA_SIZE AND W_POINTS <= 121 AND R_POINTS <= 125 THEN
				idx1 := 15;
				count := W_DATA_ADR + W_POINTS - 1;
				FOR i := W_DATA_ADR TO count DO
					idx1 := idx1 + 2;
					idx2 := idx1 + 1; (* second index necessary for pcworx! *)
					w_tmp := DATA[i];
					S_BUF.BUFFER[idx1] := WORD_TO_BYTE(SHR(w_tmp,8));
					S_BUF.BUFFER[idx2] := WORD_TO_BYTE(w_tmp);
				END_FOR;

				(* WRITE *)
				w_tmp := INT_TO_WORD(W_ADDR);
				S_BUF.BUFFER[12] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[13] := WORD_TO_BYTE(w_tmp);
				w_tmp := INT_TO_WORD(W_POINTS);
				S_BUF.BUFFER[14] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[15] := WORD_TO_BYTE(w_tmp);
				S_BUF.BUFFER[16] := SHL(S_BUF.BUFFER[15],1);

				(* READ *)
				w_tmp := INT_TO_WORD(R_ADDR);
				S_BUF.BUFFER[08] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[09] := WORD_TO_BYTE(w_tmp);
				w_tmp := INT_TO_WORD(R_POINTS);
				S_BUF.BUFFER[10] := WORD_TO_BYTE(SHR(w_tmp,8));
				S_BUF.BUFFER[11] := WORD_TO_BYTE(w_tmp);

				S_BUF.SIZE := INT_TO_UINT(idx2 + 1);
				Response := WORD_TO_INT(SHL(w_tmp,1));
				comp := 7;
			ELSE
				ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
			END_IF;
		(******************************************************************)

		ELSE
			ERROR := DWORD#00000001; (* ILLEGAL FUNCTION *)
		END_CASE;

		(* send data *)
		IF ERROR = DWORD#00 THEN
			Transaction_ID := Transaction_ID + 1;
			w_tmp := UINT_TO_WORD(S_BUF.SIZE - UINT#6);
			(* ---- H E A D E R --------------------------------------------------------------- *)
			S_BUF.BUFFER[0] := IP_ID; (* TRANSACTION-ID HI-Byte *)
			S_BUF.BUFFER[1] := INT_TO_BYTE(Transaction_ID); (* TRANSACTION-ID LO-Byte *)
			S_BUF.BUFFER[2] := BYTE#00; (* PROTOCOL ID HI byte *)
			S_BUF.BUFFER[3] := BYTE#00; (* PROTOCOL ID LO byte *)
			S_BUF.BUFFER[4] := WORD_TO_BYTE(SHR(w_tmp,8)); (* length HI byte *)
			S_BUF.BUFFER[5] := WORD_TO_BYTE(w_tmp); (* length LO byte *)
			S_BUF.BUFFER[6] := UNIT_ID; (* UNIT *)
			S_BUF.BUFFER[7] := INT_TO_BYTE(FC); (* FUNCTION-CODE *)
			(*----- END OF HEADER --------------------------------------------------------------*)
			R_BUF.SIZE := UINT#0; (* clear receive buffer *)
			state := 30; (* wait for response *)
		ELSE
			IP_STATE := BYTE#4; (* unregister *)
			state := 00;
		END_IF;
	END_IF;

30: (* response *)
	IF IP_C.ERROR <> DWORD#00 THEN
		ERROR := IP_C.ERROR;
		IP_STATE := BYTE#4; (* unregister *)
		state := 00;
	ELSIF S_BUF.SIZE = UINT#0 AND R_BUF.SIZE > UINT#0 THEN
		(* enter expected length into send buffer afterwards to allow comparison *)
		w_tmp := INT_TO_WORD(Response + 3);
		S_BUF.BUFFER[4] := WORD_TO_BYTE(SHR(w_tmp,8));
		S_BUF.BUFFER[5] := WORD_TO_BYTE(w_tmp);

		FOR i := 0 TO comp DO
			IF S_BUF.BUFFER[i] <> R_BUF.BUFFER[i] THEN EXIT; END_IF;
		END_FOR;

		IF R_BUF.SIZE >= UINT#9 AND R_BUF.BUFFER[7] > BYTE#128 THEN  (* Bit 7 = Error Message *)
			ERROR := BYTE_TO_DWORD(R_BUF.BUFFER[8]);
			IP_STATE := BYTE#4; (* unregister *)
			state := 00;
		ELSIF R_BUF.SIZE = INT_TO_UINT(Response + 9) AND i > comp THEN
			mask := SHL(DWORD#1,FC);
			(* READ_COILS, READ_INPUT_DISCRETES *)
			IF (mask AND DWORD#2#00000000_00000000_00000000_00000110) > DWORD#0 THEN (* FC 1,2 *)
				IF WORD_TO_INT(SHR(INT_TO_WORD(WORD_TO_INT(SHL(INT_TO_WORD(R_DATA_ADR),4)) + R_POINTS + R_DATA_BITPOS + 15),4)) <= DATA_SIZE THEN
					b_mask := BYTE#1;
					idx1 := 09;
					idx2 := R_DATA_ADR;
					bit_pos := R_DATA_BITPOS;
					FOR i := 1 TO R_POINTS DO
						DATA[idx2]:=BIT_LOAD_W(DATA[idx2],(R_BUF.BUFFER[idx1] AND b_mask) > BYTE#0,bit_pos);
						bit_pos := bit_pos + 1;
						IF bit_pos > 15 THEN
							idx2 := idx2 + 1;
							bit_pos := 0;
						END_IF;
						b_mask := ROL(b_mask,1);
						IF b_mask = BYTE#1 THEN
							idx1 := idx1 + 1;
						END_IF;
					END_FOR;
				ELSE
					ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
				END_IF;

			(* READ_MULTIPLE_REGISTERS, READ_INPUT_REGISTERS, READ_WRITE_REGISTERS *)
			ELSIF (mask AND DWORD#2#00000000_10000000_00000000_00011000) > DWORD#0 THEN (* FC 3,4,23 *)
				IF R_DATA_ADR + R_POINTS <= DATA_SIZE THEN
					stop := R_DATA_ADR + R_POINTS -1 ;
					idx1 := 07;
					FOR i := R_DATA_ADR TO stop DO
						idx1 := idx1 + 2;
						idx2 := idx1 + 1;
						DATA[i] := WORD_OF_BYTE(R_BUF.BUFFER[idx1],R_BUF.BUFFER[idx2]);
					END_FOR;
				ELSE
					ERROR := DWORD#2; (* ILLEGAL DATA ADDRESS *)
				END_IF;
			ELSIF (mask AND DWORD#2#11111111_00111110_01111111_10000001) > DWORD#0 THEN (* ILLEGAL FUNCTION CODE ? *)
				ERROR := DWORD#1; (* ILLEGAL FUNCTION *)
			END_IF;
		ELSE
			ERROR := DWORD#3; (* ILLEGAL DATA VALUE *)
		END_IF;

		IP_STATE := BYTE#4; (* unregister *)
		state := 00;
		R_BUF.SIZE := UINT#0;

	END_IF;

END_CASE;

(* IP_FIFO *)
IP_FIFO(FIFO:=IP_C.FIFO,STATE:=IP_STATE,ID:=IP_ID);
BUSY := IP_STATE = BYTE#3;
END_FUNCTION_BLOCK

FUNCTION_BLOCK MB_SERVER
VAR_IN_OUT
	IP_C : IP_C;
	S_BUF : NETWORK_BUFFER_SHORT;
	R_BUF : NETWORK_BUFFER_SHORT;
	VMAP : ARRAY [1..10] OF VMAP_DATA;
	DATA : ARRAY [0..255] OF WORD;
END_VAR
VAR_INPUT
	DATA_SIZE : INT;
	ENABLE : BOOL;
	UDP : BOOL;
END_VAR
VAR_OUTPUT
	ERROR : DWORD;
END_VAR
VAR
	MB_VMAP : MB_VMAP;
	adr1 : INT;
	adr2 : INT;
	points : INT;
	points2 : INT;
	coils : BYTE;
	mask : BYTE;
	idx1 : INT;
	idx2 : INT;
	bit_pos : INT;
	count : INT;
	fc : INT;
	i : INT;
	state : INT;
	ip_state : BYTE;
	ip_id : BYTE;
	IP_FIFO : IP_FIFO;
	w_tmp : WORD;
	last_cycle : TIME;
	tx : TIME;
	t : TON;
	ip_error : BOOL;
END_VAR

(* read system time *)
tx := DWORD_TO_TIME(T_PLC_MS());

FOR i := 1 TO 10 DO
	IF VMAP[i].TIME_OUT > T#0s THEN
		VMAP[i].TIME_OUT := VMAP[i].TIME_OUT + tx - last_cycle;
	END_IF;
END_FOR;
last_cycle := tx;
ip_error := IP_C.ERROR > DWORD#0;

CASE state OF
00: (* wait for release *)
	IF ENABLE THEN
		state := 10;
		ip_state := BYTE#1; (* register *)
	END_IF;

10: (* register with IP_Control and wait for release *)
	IF ip_state = BYTE#3 THEN

		(* set up IP data traffic *)
		IP_C.C_PORT := WORD#0; (* use port from IP_CONTROL *)
		IP_C.C_IP := DWORD#0; (* use IP from IP_CONTROL *)
		IP_C.C_MODE := SEL(UDP,BYTE#4,BYTE#5); (* Mode: 4=TCP+PASSIV+PORT , 5=UDP+PASSIV+PORT *)
		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
		IP_C.R_OBSERVE := TRUE; (* no receive monitoring *)
		state := 20;
	END_IF;

20:	IF NOT ip_error THEN
		IF S_BUF.SIZE = UINT#0 AND R_BUF.SIZE > UINT#6 THEN	(* data sending not active and data was received *)
			IF  R_BUF.SIZE = WORD_TO_UINT(WORD_OF_BYTE(R_BUF.BUFFER[4],R_BUF.BUFFER[5])) + UINT#6 THEN (* valid data length *)
				ERROR := DWORD#00; (* clear error *)

				(* copy header *)
				FOR i := 0 TO 13 DO
					S_BUF.BUFFER[i] := R_BUF.BUFFER[i];
				END_FOR;
				fc := BYTE_TO_INT(R_BUF.BUFFER[7]);
				adr1 := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[8],R_BUF.BUFFER[9]));
				points := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[10],R_BUF.BUFFER[11]));

				(* Virtual Address *)
				MB_VMAP(VMAP:=VMAP,FC:=fc,V_ADR:=adr1,V_CNT:=points,SIZE:=DATA_SIZE);
				adr1:=MB_VMAP.P_ADR;
				bit_pos:=MB_VMAP.P_BIT;
				ERROR:=BYTE_TO_DWORD(MB_VMAP.ERROR);

				IF ERROR = DWORD#0 THEN
					(*---------------------------------------------*)
					CASE fc OF
					1..2 : (* Read Coil Status / Read Input Status *)
						count := 9; (* first data byte *)
						coils := BYTE#0;
						mask := BYTE#1;
						FOR i := 1 TO points DO
							IF mask = BYTE#0 THEN
								mask := BYTE#1;
								S_BUF.BUFFER[count] := coils;
								coils := BYTE#0;
								count := count + 1;
							END_IF;
							IF (DATA[adr1] AND SHL(WORD#1,bit_pos)) > WORD#0 THEN
								coils := coils OR mask;
							END_IF;

							bit_Pos := bit_Pos + 1;
							IF bit_pos > 15 THEN
								ADR1 := ADR1 + 1;
								bit_pos := 0;
							END_IF;

							mask := SHL(mask,1);
						END_FOR;
						S_BUF.BUFFER[Count] := coils;

						S_BUF.BUFFER[8] := INT_TO_BYTE(count - 8); (* Bytecount *)
						S_BUF.SIZE := INT_TO_UINT(count + 1);

					(*--------------------------------------------------*)
					3..4 : (* Read Holding Registers / Read Input Registers *)
						idx1 := 07;
						count := adr1 + points - 1;
						FOR i := adr1 TO count DO
							idx1 := idx1 + 2;
							idx2 := idx1 + 1; (* second index necessary for pcworx! *)
							w_tmp := DATA[i];
							S_BUF.BUFFER[idx1] := WORD_TO_BYTE(SHR(w_tmp,8));
							S_BUF.BUFFER[idx2] := WORD_TO_BYTE(w_tmp);
						END_FOR;
						S_BUF.BUFFER[8] := SHL(INT_TO_BYTE(points),1); (* number of registers in bytes *)
						S_BUF.SIZE := INT_TO_UINT(idx2 + 1);
					(*--------------------------------------------------*)
					5 : (* Force Single Coil *)
						DATA[adr1]:=BIT_LOAD_W(DATA[adr1],R_BUF.BUFFER[10] > BYTE#0,bit_pos);
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					6 : (* Preset Single Register *)
						DATA[adr1] := WORD_OF_BYTE(R_BUF.BUFFER[10],R_BUF.BUFFER[11]);
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					15 : (* Force Multible Coils *)
						mask := BYTE#1;
						idx1 := 13;
						FOR i := 1 TO points DO
							DATA[adr1]:=BIT_LOAD_W(DATA[adr1],(R_BUF.BUFFER[idx1] AND mask) > BYTE#0,bit_pos);
							bit_pos := bit_pos + 1;
							IF bit_pos > 15 THEN
								adr1:= adr1 + 1;
								bit_pos := 0;
							END_IF;
							mask := ROL(mask,1);
							IF mask = BYTE#1 THEN
								idx1 := idx1 +1;
							END_IF;
						END_FOR;
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					16 : (* Preset Multiple Registers	*)
						idx1 := 11;
						count := adr1 + points - 1;
						FOR i := adr1 TO Count DO
							idx1 := idx1 + 2;
	 						idx2 := idx1 + 1;
							Data[i] := WORD_OF_BYTE(R_BUF.BUFFER[idx1],R_BUF.BUFFER[idx2]);
						END_FOR;
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					22 : (* Mask Write Register *)
						w_tmp := WORD_OF_BYTE(R_BUF.BUFFER[10],R_BUF.BUFFER[11]); (* and mask *)
						DATA[adr1] := (DATA[adr1] AND w_tmp) OR (WORD_OF_BYTE(R_BUF.BUFFER[12],R_BUF.BUFFER[13]) AND (NOT w_tmp));
						S_BUF.SIZE := UINT#14;
					(*--------------------------------------------------*)
					23 : (* Read / Write Registers	*)
						(*--- Write-Register -----------------------------------*)
						adr2 := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[12],R_BUF.BUFFER[13]));
						points2 := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[14],R_BUF.BUFFER[15]));

						(* Virtual Address Write Register *)
						MB_VMAP(VMAP:=VMAP,FC:=16,V_ADR:=adr2,V_CNT:=points2,SIZE:=DATA_SIZE);
						adr2:=MB_VMAP.P_ADR;
						ERROR:=BYTE_TO_DWORD(MB_VMAP.ERROR);

						IF ERROR = DWORD#0 THEN
							idx1 := 15;
							count := adr2 + points2 - 1;
							FOR i := adr2 TO count DO
								idx1 := idx1 + 2;
	 							idx2 := idx1 + 1; (* second index necessary for pcworx! *)
								Data[i] := WORD_OF_BYTE(R_BUF.BUFFER[idx1],R_BUF.BUFFER[idx2]);
							END_FOR;

							(*--- Read-Register -----------------------------------*)
							idx1 := 07;
							count := adr1 + points - 1;
							FOR i := adr1 TO count DO
								idx1 := idx1 + 2;
								idx2 := idx1 + 1; (* second index necessary for pcworx! *)
								w_tmp := DATA[i];
								S_BUF.BUFFER[idx1] := WORD_TO_BYTE(SHR(w_tmp,8));
								S_BUF.BUFFER[idx2] := WORD_TO_BYTE(w_tmp);
							END_FOR;
							S_BUF.BUFFER[8] := SHL(INT_TO_BYTE(points),1); (* number of registers in bytes *)
							S_BUF.SIZE := INT_TO_UINT(idx2 + 1);

						END_IF;

					(*--------------------------------------------------*)
					ELSE
						ERROR := DWORD#01; (* ILLEGAL FUNCTION *)
					END_CASE;
				END_IF;
			END_IF;

			IF ERROR > DWORD#0 THEN
				S_BUF.BUFFER[7] := S_BUF.BUFFER[7] OR BYTE#2#1000_0000; (* set error bit *)
				S_BUF.BUFFER[8] := DWORD_TO_BYTE(ERROR);
				S_BUF.SIZE := UINT#9;
			END_IF;

			S_BUF.BUFFER[4] := BYTE#0; (* high Byte *)
			S_BUF.BUFFER[5] := UINT_TO_BYTE(S_BUF.SIZE - UINT#6); (* enter length field - low byte only *)

		END_IF;
	ELSE
		ERROR := IP_C.ERROR;
	END_IF;

	t(IN:= ip_error, PT:=T#2s);
	IF t.Q THEN
		IP_C.TIME_RESET := TRUE; (* Error Reset *)
	END_IF;

	R_BUF.SIZE := UINT#0; (* destroy receive data *)

	IF ENABLE = FALSE OR t.Q THEN
	    ip_state := BYTE#4; (* unregister *)
		state := 00;
		IP_C.C_ENABLE := FALSE;
	END_IF;
END_CASE;

(* IP_FIFO *)
IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id);
END_FUNCTION_BLOCK

FUNCTION_BLOCK MB_SERVER_1
VAR_IN_OUT
	IP_C : IP_C;
	S_BUF : NETWORK_BUFFER_SHORT;
	R_BUF : NETWORK_BUFFER_SHORT;
	VMAP : ARRAY [1..10] OF VMAP_DATA;
	DATA : ARRAY [0..255] OF WORD;
END_VAR
VAR_INPUT
	DATA_SIZE : INT;
	ENABLE : BOOL;
	UDP : BOOL;
END_VAR
VAR_OUTPUT
	ERROR : DWORD;
END_VAR
VAR
	MB_VMAP : MB_VMAP;
	adr1 : INT;
	adr2 : INT;
	points : INT;
	points2 : INT;
	coils : BYTE;
	mask : BYTE;
	idx1 : INT;
	idx2 : INT;
	bit_pos : INT;
	count : INT;
	fc : INT;
	i : INT;
	state : INT;
	ip_state : BYTE;
	ip_id : BYTE;
	IP_FIFO : IP_FIFO;
	w_tmp : WORD;
	last_cycle : TIME;
	tx : TIME;
	t : TON;
	ip_error : BOOL;
END_VAR

(* read system time *)
tx := DWORD_TO_TIME(T_PLC_MS());

FOR i := 1 TO 10 DO
	IF VMAP[i].TIME_OUT > T#0s THEN
		VMAP[i].TIME_OUT := VMAP[i].TIME_OUT + tx - last_cycle;
	END_IF;
END_FOR;
last_cycle := tx;
ip_error := IP_C.ERROR > DWORD#0;

CASE state OF
00: (* wait for release *)
	IF ENABLE THEN
		state := 10;
		ip_state := BYTE#1; (* register *)
	END_IF;

10: (* register with IP_Control and wait for release *)
	IF ip_state = BYTE#3 THEN

		(* set up IP data traffic *)
		IP_C.C_PORT := WORD#0; (* use port from IP_CONTROL *)
		IP_C.C_IP := DWORD#0; (* use IP from IP_CONTROL *)
		IP_C.C_MODE := SEL(UDP,BYTE#4,BYTE#5); (* Mode: 4=TCP+PASSIV+PORT , 5=UDP+PASSIV+PORT *)
		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
		IP_C.R_OBSERVE := FALSE; (* no receive monitoring *)
		state := 20;
	END_IF;

20:	IF NOT ip_error THEN
		IF S_BUF.SIZE = UINT#0 AND R_BUF.SIZE > UINT#6 THEN	(* data sending not active and data was received *)
			IF  R_BUF.SIZE = WORD_TO_UINT(WORD_OF_BYTE(R_BUF.BUFFER[4],R_BUF.BUFFER[5])) + UINT#6 THEN (* valid data length *)
				ERROR := DWORD#00; (* clear error *)

				(* copy header *)
				FOR i := 0 TO 13 DO
					S_BUF.BUFFER[i] := R_BUF.BUFFER[i];
				END_FOR;
				fc := BYTE_TO_INT(R_BUF.BUFFER[7]);
				adr1 := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[8],R_BUF.BUFFER[9]));
				points := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[10],R_BUF.BUFFER[11]));

				(* Virtual Address *)
				MB_VMAP(VMAP:=VMAP,FC:=fc,V_ADR:=adr1,V_CNT:=points,SIZE:=DATA_SIZE);
				adr1:=MB_VMAP.P_ADR;
				bit_pos:=MB_VMAP.P_BIT;
				ERROR:=BYTE_TO_DWORD(MB_VMAP.ERROR);

				IF ERROR = DWORD#0 THEN
					(*---------------------------------------------*)
					CASE fc OF
					1..2 : (* Read Coil Status / Read Input Status *)
						count := 9; (* first data byte *)
						coils := BYTE#0;
						mask := BYTE#1;
						FOR i := 1 TO points DO
							IF mask = BYTE#0 THEN
								mask := BYTE#1;
								S_BUF.BUFFER[count] := coils;
								coils := BYTE#0;
								count := count + 1;
							END_IF;
							IF (DATA[adr1] AND SHL(WORD#1,bit_pos)) > WORD#0 THEN
								coils := coils OR mask;
							END_IF;

							bit_Pos := bit_Pos + 1;
							IF bit_pos > 15 THEN
								ADR1 := ADR1 + 1;
								bit_pos := 0;
							END_IF;

							mask := SHL(mask,1);
						END_FOR;
						S_BUF.BUFFER[Count] := coils;

						S_BUF.BUFFER[8] := INT_TO_BYTE(count - 8); (* Bytecount *)
						S_BUF.SIZE := INT_TO_UINT(count + 1);

					(*--------------------------------------------------*)
					3..4 : (* Read Holding Registers / Read Input Registers *)
						idx1 := 07;
						count := adr1 + points - 1;
						FOR i := adr1 TO count DO
							idx1 := idx1 + 2;
							idx2 := idx1 + 1; (* second index necessary for pcworx! *)
							w_tmp := DATA[i];
							S_BUF.BUFFER[idx1] := WORD_TO_BYTE(SHR(w_tmp,8));
							S_BUF.BUFFER[idx2] := WORD_TO_BYTE(w_tmp);
						END_FOR;
						S_BUF.BUFFER[8] := SHL(INT_TO_BYTE(points),1); (* number of registers in bytes *)
						S_BUF.SIZE := INT_TO_UINT(idx2 + 1);
					(*--------------------------------------------------*)
					5 : (* Force Single Coil *)
						DATA[adr1]:=BIT_LOAD_W(DATA[adr1],R_BUF.BUFFER[10] > BYTE#0,bit_pos);
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					6 : (* Preset Single Register *)
						DATA[adr1] := WORD_OF_BYTE(R_BUF.BUFFER[10],R_BUF.BUFFER[11]);
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					15 : (* Force Multible Coils *)
						mask := BYTE#1;
						idx1 := 13;
						FOR i := 1 TO points DO
							DATA[adr1]:=BIT_LOAD_W(DATA[adr1],(R_BUF.BUFFER[idx1] AND mask) > BYTE#0,bit_pos);
							bit_pos := bit_pos + 1;
							IF bit_pos > 15 THEN
								adr1:= adr1 + 1;
								bit_pos := 0;
							END_IF;
							mask := ROL(mask,1);
							IF mask = BYTE#1 THEN
								idx1 := idx1 +1;
							END_IF;
						END_FOR;
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					16 : (* Preset Multiple Registers	*)
						idx1 := 11;
						count := adr1 + points - 1;
						FOR i := adr1 TO Count DO
							idx1 := idx1 + 2;
	 						idx2 := idx1 + 1;
							Data[i] := WORD_OF_BYTE(R_BUF.BUFFER[idx1],R_BUF.BUFFER[idx2]);
						END_FOR;
						S_BUF.SIZE := UINT#12;
					(*--------------------------------------------------*)
					22 : (* Mask Write Register *)
						w_tmp := WORD_OF_BYTE(R_BUF.BUFFER[10],R_BUF.BUFFER[11]); (* and mask *)
						DATA[adr1] := (DATA[adr1] AND w_tmp) OR (WORD_OF_BYTE(R_BUF.BUFFER[12],R_BUF.BUFFER[13]) AND (NOT w_tmp));
						S_BUF.SIZE := UINT#14;
					(*--------------------------------------------------*)
					23 : (* Read / Write Registers	*)
						(*--- Write-Register -----------------------------------*)
						adr2 := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[12],R_BUF.BUFFER[13]));
						points2 := WORD_TO_INT(WORD_OF_BYTE(R_BUF.BUFFER[14],R_BUF.BUFFER[15]));

						(* Virtual Address Write Register *)
						MB_VMAP(VMAP:=VMAP,FC:=16,V_ADR:=adr2,V_CNT:=points2,SIZE:=DATA_SIZE);
						adr2:=MB_VMAP.P_ADR;
						ERROR:=BYTE_TO_DWORD(MB_VMAP.ERROR);

						IF ERROR = DWORD#0 THEN
							idx1 := 15;
							count := adr2 + points2 - 1;
							FOR i := adr2 TO count DO
								idx1 := idx1 + 2;
	 							idx2 := idx1 + 1; (* second index necessary for pcworx! *)
								Data[i] := WORD_OF_BYTE(R_BUF.BUFFER[idx1],R_BUF.BUFFER[idx2]);
							END_FOR;

							(*--- Read-Register -----------------------------------*)
							idx1 := 07;
							count := adr1 + points - 1;
							FOR i := adr1 TO count DO
								idx1 := idx1 + 2;
								idx2 := idx1 + 1; (* second index necessary for pcworx! *)
								w_tmp := DATA[i];
								S_BUF.BUFFER[idx1] := WORD_TO_BYTE(SHR(w_tmp,8));
								S_BUF.BUFFER[idx2] := WORD_TO_BYTE(w_tmp);
							END_FOR;
							S_BUF.BUFFER[8] := SHL(INT_TO_BYTE(points),1); (* number of registers in bytes *)
							S_BUF.SIZE := INT_TO_UINT(idx2 + 1);

						END_IF;

					(*--------------------------------------------------*)
					ELSE
						ERROR := DWORD#01; (* ILLEGAL FUNCTION *)
					END_CASE;
				END_IF;
			END_IF;

			IF ERROR > DWORD#0 THEN
				S_BUF.BUFFER[7] := S_BUF.BUFFER[7] OR BYTE#2#1000_0000; (* set error bit *)
				S_BUF.BUFFER[8] := DWORD_TO_BYTE(ERROR);
				S_BUF.SIZE := UINT#9;
			END_IF;

			S_BUF.BUFFER[4] := BYTE#0; (* high Byte *)
			S_BUF.BUFFER[5] := UINT_TO_BYTE(S_BUF.SIZE - UINT#6); (* enter length field - low byte only *)

		END_IF;
	ELSE
		ERROR := IP_C.ERROR;
	END_IF;

	t(IN:= ip_error, PT:=T#5s);
	IF t.Q THEN
		IP_C.TIME_RESET := TRUE; (* Error Reset *)
	END_IF;

	R_BUF.SIZE := UINT#0; (* destroy receive data *)

	IF ENABLE = FALSE THEN
	    ip_state := BYTE#4; (* unregister *)
		state := 00;
	END_IF;
END_CASE;

(* IP_FIFO *)
IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id);
END_FUNCTION_BLOCK

FUNCTION_BLOCK MB_VMAP
VAR_IN_OUT
	VMAP :  ARRAY [1..10] OF VMAP_DATA;
END_VAR
VAR_INPUT
	FC : INT;
	V_ADR : INT;
	V_CNT : INT;
	SIZE : INT;
END_VAR
VAR_OUTPUT
	P_ADR : INT;
	P_BIT : INT;
	ERROR : BYTE;
END_VAR
VAR
	i : INT;
	vopt : DWORD;
	vadr : INT;
	padr : INT;
	mask : DWORD;
	vsize : INT;
	w_tmp : WORD;
	i_tmp : INT;
	init : BOOL;
END_VAR

IF init = FALSE THEN
	init := TRUE;
	IF VMAP[1].FC = DWORD#0 THEN
		VMAP[1].FC := DWORD#16#FFFF_FFFF;
		VMAP[1].V_ADR := 0;
		VMAP[1].V_SIZE := 256;
		VMAP[1].P_ADR := 0;
	END_IF;
END_IF;

ERROR := BYTE#02; (* ILLEGAL DATA ADDRESS *)
mask := SHL(DWORD#1,FC);

IF (mask AND DWORD#2#00000000_01000000_00000000_01100000) <> DWORD#0 THEN (* FC 5,6,22 *)
	V_CNT := 1;
ELSIF V_CNT = 0 THEN
	RETURN;
END_IF;

FOR i := 1 TO 10 DO
	vopt := VMAP[i].FC;
	IF (vopt AND mask) <> DWORD#0 THEN
		vadr := VMAP[i].V_ADR;
		vsize := VMAP[i].V_SIZE;
		padr := VMAP[i].P_ADR;

		IF (mask AND DWORD#2#00000000_00000000_10000000_00100110) <> DWORD#0 THEN (* FC 1,2,5,15 for bit access *)
			i_tmp := WORD_TO_INT(SHR(INT_TO_WORD(V_ADR + V_CNT -1),4))+1;
			IF i_tmp <= SIZE AND i_tmp <= vsize THEN
				w_tmp := INT_TO_WORD(V_ADR);
				P_ADR := WORD_TO_INT(SHR(w_tmp,4)) + padr;
				P_BIT := WORD_TO_INT(w_tmp AND WORD#16#000F);
				ERROR := BYTE#0;
				EXIT;
			END_IF;
		ELSIF (mask AND DWORD#2#00000000_11000001_00000000_01011000) <> DWORD#0 THEN (* FC 3,4,6,16,22,23 for register access *)
			IF V_ADR >= vadr AND V_ADR + V_CNT <= vadr + vsize THEN (* valid virtual range *)
				P_ADR := V_ADR - vadr + padr; (* calculate real address *)
				IF P_ADR + V_CNT <= SIZE THEN (* out of range? *)
					ERROR := BYTE#0;
					EXIT;
				END_IF;
			END_IF;
		ELSE
			ERROR := BYTE#01; (* ILLEGAL FUNCTION *)
			RETURN;
		END_IF;
	END_IF;
END_FOR;

IF ERROR = BYTE#0 THEN
	IF (mask AND DWORD#2#00000000_11000001_10000000_01100000) <> DWORD#0 THEN (* FC 5,6,15,16,22,23 write access *)
		VMAP[i].TIME_OUT := T#1ms;
	END_IF;
END_IF;
END_FUNCTION_BLOCK

FUNCTION MD5_AUX : DWORD
VAR_INPUT
	N : INT;
	A : DWORD;
	B : DWORD;
	C : DWORD;
	D : DWORD;
	X : DWORD;
	U : INT;
	T : DWORD;
END_VAR
VAR
	w : DWORD;
END_VAR

CASE N OF
	1: w := (B AND C) OR (NOT(B) AND D);
	2: w := (B AND D) OR (C AND NOT D);
	3: w := (B XOR C XOR D);
	4: w := (C XOR (B OR  NOT(D)));
END_CASE

MD5_AUX := B + ROL(A + w + X + T,U);
END_FUNCTION

// Commented out for beedance: it uses MD5_TO_STRH, which is commented out.
// FUNCTION_BLOCK MD5_CRAM_AUTH
// VAR_IN_OUT
// 	RUN : BOOL;
// 	USERNAME : STRING(64);
// 	PASSWORD : STRING(64);
// 	B64_TS : STRING(64);  (* Base64 coded ServerTimeStamp *)
// 	AUTH_KEY : STRING(192);
// END_VAR
// VAR
// 	BASE64_DECODE_STR : BASE64_DECODE_STR;
// 	BASE64_ENCODE_STR : BASE64_ENCODE_STR;
// 	MD5_STREAM : MD5_STREAM;
// 	step : INT;
// 	md5_mode : INT;
// 	buf : ARRAY [0..63] OF BYTE;
// 	md5 : ARRAY [0..15] OF BYTE;
// 	md5_first : ARRAY [0..15] OF BYTE;
// 	md5_str1 : STRING(64);
// 	md5_str2 : STRING(64);
// 	size : UDINT;
// 	n1 : INT;
// 	n2 : INT;
// 	pos : UDINT;
// 	b64d_run : BOOL;
// 	b64d_done : BOOL;
// 	b64e_run : BOOL;
// 	b64e_done : BOOL;
// 	str_192 : STRING(192);
// 	str_144 : STRING(144);
// 	xpad : BYTE;
// 	pt : REF_TO BYTE;
// END_VAR
//
// CASE step OF
//
// 0:	IF RUN THEN
// 		step := 5;
// 		str_192 := B64_TS;
// 		b64d_run := TRUE;
// 	END_IF;
//
// 5: IF B64D_DONE THEN (* IPAD cycle *)
// 		b64d_run := FALSE;
// 		xpad := 16#36;
// 		md5_str1 := PASSWORD;
// 		md5_str2 := str_144;
// 		size := 64 + LEN(md5_str2);
// 		md5_mode := 1;
// 		step := 10;
// 	END_IF;
//
// 10: IF md5_mode = 3 THEN (* OPAD Zyklus *)
// 		md5_first := md5; (* save md5 array of IPAD *)
// 		xpad := 16#5C;
// 		size := 64 + 16;
// 		md5_mode := 1;
// 		step := 20;
// 	END_IF;
//
// 20: IF md5_mode = 3 THEN
// 		md5_mode := 0;
// 		str_144 := CONCAT(USERNAME , ' ');
// 		str_144 := CONCAT(str_144 , MD5_TO_STRH(md5));
// 		b64e_run := TRUE;
// 		step := 30;
// 	END_IF;
//
// 30: IF B64E_DONE THEN
// 		AUTH_KEY := str_192;
// 		b64e_run := FALSE;
// 		RUN := FALSE;
// 		step := 0; (* all finished !! *)
// 	END_IF;
//
// END_CASE;
//
// (* ----------- md5 -------------- *)
// CASE md5_mode OF
//
// 1:	n2 := LEN(md5_str1) - 1;
// 	pt := REF(md5_str1);
// 	FOR n1 := 0 TO n2 DO (* copy string into buf *)
// 		buf[n1] := pt^ XOR xpad;
// 		pt := pt + 1;
// 	END_FOR;
// 	FOR n1 := n2 + 1 TO 63 DO (* fill with ipad *)
// 		buf[n1] := xpad;
// 	END_FOR;
//
// 2:	(* Data copy and Calc *)
// 	IF size > 0 AND pos >= 64 THEN
// 		IF step = 20 THEN
// 			FOR n1 := 0 TO 15 DO
// 				buf[n1] := md5_first[n1];
// 			END_FOR;
// 		ELSE
// 			n2 := LEN(md5_str2) - 1;
// 			pt := REF(md5_str2);
// 			FOR n1 := 0 TO n2 DO (* copy string into buf *)
// 				buf[n1] := pt^;
// 				pt := pt + 1;
// 			END_FOR;
// 		END_IF;
// 	END_IF;
//
// END_CASE;
//
// IF md5_mode > 0 THEN
// 	MD5_STREAM(SIZE:=size, MODE:=md5_mode, BUF:=buf, MD5:=md5, POS=>pos);
// END_IF;
//
// BASE64_DECODE_STR(RUN:=b64d_run, STR1:=str_192, STR2:=str_144, DONE=>b64d_done);
// BASE64_ENCODE_STR(RUN:=b64e_run, STR1:=str_144, STR2:=str_192, DONE=>b64e_done);
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK MD5_STR
// VAR_INPUT
// 	RUN: BOOL;
// END_VAR
// VAR_OUTPUT
// 	DONE: BOOL;
// END_VAR
// VAR_IN_OUT
// 	STR: STRING(STRING_LENGTH);
// 	MD5 : ARRAY [0..15] OF BYTE;
// END_VAR
// VAR
// 	run_last: BOOL;
// 	MD5_STREAM : MD5_STREAM;
// 	buf : ARRAY [0..63] OF BYTE;
// 	mode : INT;
// 	size : UDINT;
// 	pos : UDINT;
// END_VAR
//
// CASE mode OF
// 0:	(* wait for start *)
// 	IF RUN AND NOT run_last THEN
// 		DONE := FALSE;
// 		mode := 1; (* initialization *)
// 		size := LEN(STR);
// 	END_IF;
//
// 2:	(* Data copy and Calc *)
// 	IF size > 0 THEN
// 		_STRING_TO_BUFFER(STR:=MID(STR, UDINT_TO_INT(SIZE), UDINT_TO_INT(POS) + 1), POS:=0, PT:=REF(buf), SIZE:=SIZEOF(buf));
// 	END_IF;
//
// 3:	DONE := TRUE;
// 	mode := 0; (* Stop *)
//
// END_CASE;
//
// IF mode > 0 THEN
// 	MD5_STREAM(SIZE:=size, MODE:=mode, BUF:=buf, MD5:=MD5, POS=>pos);
// END_IF;
//
// run_last := RUN;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK MD5_STREAM
VAR_IN_OUT
	MODE : INT;
	BUF : ARRAY [0..63] OF BYTE;
	MD5 : ARRAY [0..15] OF BYTE;
	SIZE : UDINT;
END_VAR

VAR_OUTPUT
	POS : UDINT;
END_VAR

VAR
	hash : ARRAY [0..3] OF DWORD;
	a : DWORD;
	b : DWORD;
	c : DWORD;
	d : DWORD;
	end : UDINT;
	block : UDINT;
	x : ARRAY [0..15] OF DWORD;
	n : INT;
	n1 : INT;
	z : INT;
	pad_1 : BOOL;
END_VAR

CASE MODE OF

1:	(* initialization of the hash values *)
	hash[0] := DWORD#16#67452301;
	hash[1] := DWORD#16#EFCDAB89;
	hash[2] := DWORD#16#98BADCFE;
	hash[3] := DWORD#16#10325476;

	block := SHR(SIZE,6) + BOOL_TO_UDINT((UDINT_TO_BYTE(SIZE) AND 16#3F) > 55);
	pad_1 := FALSE;
	POS := 0;
	End := SIZE;
	SIZE := MIN(64,End);
	MODE := 2;

2:	(* append padding bytes *)
	FOR n := UDINT_TO_INT(SIZE) TO 63 DO
		BUF[n] := 0;
	END_FOR;

	IF SIZE < 64 AND NOT pad_1 THEN
		BUF[SIZE] := 2#1000_0000;
		pad_1 := TRUE;
	END_IF;

	POS := POS + SIZE;
	SIZE := MIN(64,End - POS);

	n := 0;	(* copy byte array into DWORD array *)
	FOR n1 := 0 TO 15 DO
		x[n1] := DWORD_OF_BYTE(BUF[n+3],BUF[n+2],BUF[n+1],BUF[n]); (* Little Endian Format *)
		n := n + 4;
	END_FOR;

	IF block = 0 THEN (* enter the data length at the end of the last block *)
		x[14] := SHL(UDINT_TO_DWORD(End),3); (* enter number of data bits in little-endian format *)
		MODE := 3;
	ELSE
		block := block - 1;
	END_IF;

	a := hash[0];
	b := hash[1];
	c := hash[2];
	d := hash[3];

	a := MD5_AUX(1, a, b, c, d, X[00], 07, 16#D76AA478); (* 01 *)
	d := MD5_AUX(1, d, a, b, c, X[01], 12, 16#E8C7B756); (* 02 *)
	c := MD5_AUX(1, c, d, a, b, X[02], 17, 16#242070DB); (* 03 *)
	b := MD5_AUX(1, b, c, d, a, X[03], 22, 16#C1BDCEEE); (* 04 *)
	a := MD5_AUX(1, a, b, c, d, X[04], 07, 16#F57C0FAF); (* 05 *)
	d := MD5_AUX(1, d, a, b, c, X[05], 12, 16#4787C62A); (* 06 *)
	c := MD5_AUX(1, c, d, a, b, X[06], 17, 16#A8304613); (* 07 *)
	b := MD5_AUX(1, b, c, d, a, X[07], 22, 16#FD469501); (* 08 *)
	a := MD5_AUX(1, a, b, c, d, X[08], 07, 16#698098D8); (* 09 *)
	d := MD5_AUX(1, d, a, b, c, X[09], 12, 16#8B44F7AF); (* 10 *)
	c := MD5_AUX(1, c, d, a, b, X[10], 17, 16#FFFF5BB1); (* 11 *)
	b := MD5_AUX(1, b, c, d, a, X[11], 22, 16#895CD7BE); (* 12 *)
	a := MD5_AUX(1, a, b, c, d, X[12], 07, 16#6B901122); (* 13 *)
	d := MD5_AUX(1, d, a, b, c, X[13], 12, 16#FD987193); (* 14 *)
	c := MD5_AUX(1, c, d, a, b, X[14], 17, 16#A679438E); (* 15 *)
	b := MD5_AUX(1, b, c, d, a, X[15], 22, 16#49B40821); (* 16 *)

	a := MD5_AUX(2, a, b, c, d, X[01], 05, 16#F61E2562); (* 17 *)
	d := MD5_AUX(2, d, a, b, c, X[06], 09, 16#C040B340); (* 18 *)
	c := MD5_AUX(2, c, d, a, b, X[11], 14, 16#265E5A51); (* 19 *)
	b := MD5_AUX(2, b, c, d, a, X[00], 20, 16#E9B6C7AA); (* 20 *)
	a := MD5_AUX(2, a, b, c, d, X[05], 05, 16#D62F105D); (* 21 *)
	d := MD5_AUX(2, d, a, b, c, X[10], 09, 16#02441453); (* 22 *)
	c := MD5_AUX(2, c, d, a, b, X[15], 14, 16#D8A1E681); (* 23 *)
	b := MD5_AUX(2, b, c, d, a, X[04], 20, 16#E7D3FBC8); (* 24 *)
	a := MD5_AUX(2, a, b, c, d, X[09], 05, 16#21E1CDE6); (* 25 *)
	d := MD5_AUX(2, d, a, b, c, X[14], 09, 16#C33707D6); (* 26 *)
	c := MD5_AUX(2, c, d, a, b, X[03], 14, 16#F4D50D87); (* 27 *)
	b := MD5_AUX(2, b, c, d, a, X[08], 20, 16#455A14ED); (* 28 *)
	a := MD5_AUX(2, a, b, c, d, X[13], 05, 16#A9E3E905); (* 29 *)
	d := MD5_AUX(2, d, a, b, c, X[02], 09, 16#FCEFA3F8); (* 30 *)
	c := MD5_AUX(2, c, d, a, b, X[07], 14, 16#676F02D9); (* 31 *)
	b := MD5_AUX(2, b, c, d, a, X[12], 20, 16#8D2A4C8A); (* 32 *)

	a := MD5_AUX(3, a, b, c, d, X[05], 04, 16#FFFA3942); (* 33 *)
	d := MD5_AUX(3, d, a, b, c, X[08], 11, 16#8771F681); (* 34 *)
	c := MD5_AUX(3, c, d, a, b, X[11], 16, 16#6D9D6122); (* 35 *)
	b := MD5_AUX(3, b, c, d, a, X[14], 23, 16#FDE5380C); (* 36 *)
	a := MD5_AUX(3, a, b, c, d, X[01], 04, 16#A4BEEA44); (* 37 *)
	d := MD5_AUX(3, d, a, b, c, X[04], 11, 16#4BDECFA9); (* 38 *)
	c := MD5_AUX(3, c, d, a, b, X[07], 16, 16#F6BB4B60); (* 39 *)
	b := MD5_AUX(3, b, c, d, a, X[10], 23, 16#BEBFBC70); (* 40 *)
	a := MD5_AUX(3, a, b, c, d, X[13], 04, 16#289B7EC6); (* 41 *)
	d := MD5_AUX(3, d, a, b, c, X[00], 11, 16#EAA127FA); (* 42 *)
	c := MD5_AUX(3, c, d, a, b, X[03], 16, 16#D4EF3085); (* 43 *)
	b := MD5_AUX(3, b, c, d, a, X[06], 23, 16#04881D05); (* 44 *)
	a := MD5_AUX(3, a, b, c, d, X[09], 04, 16#D9D4D039); (* 45 *)
	d := MD5_AUX(3, d, a, b, c, X[12], 11, 16#E6DB99E5); (* 46 *)
	c := MD5_AUX(3, c, d, a, b, X[15], 16, 16#1FA27CF8); (* 47 *)
	b := MD5_AUX(3, b, c, d, a, X[02], 23, 16#C4AC5665); (* 48 *)

	a := MD5_AUX(4, a, b, c, d, X[00], 06, 16#F4292244); (* 49 *)
	d := MD5_AUX(4, d, a, b, c, X[07], 10, 16#432AFF97); (* 50 *)
	c := MD5_AUX(4, c, d, a, b, X[14], 15, 16#AB9423A7); (* 51 *)
	b := MD5_AUX(4, b, c, d, a, X[05], 21, 16#FC93A039); (* 52 *)
	a := MD5_AUX(4, a, b, c, d, X[12], 06, 16#655B59C3); (* 53 *)
	d := MD5_AUX(4, d, a, b, c, X[03], 10, 16#8F0CCC92); (* 54 *)
	c := MD5_AUX(4, c, d, a, b, X[10], 15, 16#FFEFF47D); (* 55 *)
	b := MD5_AUX(4, b, c, d, a, X[01], 21, 16#85845DD1); (* 56 *)
	a := MD5_AUX(4, a, b, c, d, X[08], 06, 16#6FA87E4F); (* 57 *)
	d := MD5_AUX(4, d, a, b, c, X[15], 10, 16#FE2CE6E0); (* 58 *)
	c := MD5_AUX(4, c, d, a, b, X[06], 15, 16#A3014314); (* 59 *)
	b := MD5_AUX(4, b, c, d, a, X[13], 21, 16#4E0811A1); (* 60 *)
	a := MD5_AUX(4, a, b, c, d, X[04], 06, 16#F7537E82); (* 61 *)
	d := MD5_AUX(4, d, a, b, c, X[11], 10, 16#BD3AF235); (* 62 *)
	c := MD5_AUX(4, c, d, a, b, X[02], 15, 16#2AD7D2BB); (* 63 *)
	b := MD5_AUX(4, b, c, d, a, X[09], 21, 16#EB86D391); (* 64 *)

	hash[0] := hash[0] + a;
	hash[1] := hash[1] + b;
	hash[2] := hash[2] + c;
	hash[3] := hash[3] + d;

END_CASE;

IF MODE = 3 THEN
	z := 0;
	FOR n := 0 TO 3 DO (* convert from little-endian to big-endian format and DWORD to BYTE array *)
		FOR n1 := 0 TO 3 DO
			MD5[z] := DWORD_TO_BYTE(hash[n]);
			hash[n] := ROR(hash[n],8);
			z := z + 1;
		END_FOR;
	END_FOR;
END_IF;
END_FUNCTION_BLOCK

// Commented out for beedance: it computes with pointers as addresses (pt := pt + 1); beedance references (REF_TO) have no arithmetic.
// FUNCTION MD5_TO_STRH : STRING(32)
// VAR_INPUT
// 	MD5 : ARRAY [0..15] OF BYTE;
// END_VAR
// VAR
// 	i : INT;
// 	tmp : BYTE;
// 	pt : REF_TO BYTE;
// END_VAR
// pt := REF(MD5_TO_STRH); (* read output adress to pointer *)
//
// FOR i := 0 TO 15 DO; (* write the 8 hex characters backwards *)
// 	tmp := SHR(MD5[i],4);
// 	PT^ := tmp + SEL(tmp <= 9,87,48);
// 	pt := pt + 1; (* decrement the pointer *)
// 	tmp := MD5[i] AND 2#1111;
// 	PT^ := tmp + SEL(tmp <= 9,87,48);
// 	pt := pt + 1; (* decrement the pointer *)
// END_FOR;
// pt^:= 0; (* write the closing byte for the string *)
// END_FUNCTION

FUNCTION_BLOCK MOON_PHASE
VAR_INPUT
	XDT : DT;
	SCALE : BYTE := 12;
	UPDATE : TIME := T#1h;
END_VAR
VAR_OUTPUT
	PHASE : BYTE;
END_VAR
VAR
	x : UDINT;
	last_dt : DT;
END_VAR

IF XDT - last_dt > UPDATE THEN
	x := (DT_TO_UDINT(XDT) - 603240) MOD 2551392;
	PHASE := UDINT_TO_BYTE(SCALE * x / 2551392) ;
	last_dt := XDT;
END_IF;
END_FUNCTION_BLOCK

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK MYSQL_AUTH
// VAR_IN_OUT
// 	RUN : BOOL;
// 	PASSWORD : STRING(64);
// 	MESSAGE : ARRAY [0..19] OF BYTE;
// 	SCRAMBLE : ARRAY [0..19] OF BYTE;
// END_VAR
// VAR
// 	SHA1_STREAM : SHA1_STREAM;
// 	buf : ARRAY [0..63] OF BYTE;
// 	sha_mode : INT;
// 	step : INT;
// 	size : UDINT;
// 	pos : UDINT;
// 	i : INT;
// 	i2: INT;
// 	scramble1 : ARRAY [0..19] OF BYTE;
// 	scramble2 : ARRAY [0..19] OF BYTE;
// END_VAR
//
// IF NOT RUN THEN RETURN; END_IF;
// CASE step OF
//
// 0:	CASE sha_mode OF
// 	0:	sha_mode := 1; (* initialization *)
// 		size := INT_TO_UDINT(LEN(PASSWORD));
// 		i2 := UDINT_TO_INT(size) - 1;
// 		FOR i := 0 TO 63 DO
// 			IF i <= i2 THEN
// 				buf[i] := CODE(STR:=PASSWORD,POS:=i+1);
// 			ELSE
// 				buf[i] := BYTE#0;
// 			END_IF;
// 		END_FOR;
//
// 	3:	scramble1 := SCRAMBLE;
// 		sha_mode := 1; (* initialization *)
// 		size := UDINT#20;
// 		FOR i := 0 TO 19 DO
// 			buf[i] := scramble1[i];
// 		END_FOR;
// 		step := 10;
// 	END_CASE;
//
// 10:	IF sha_mode = 3 THEN
// 		sha_mode := 1; (* initialization *)
// 		scramble2 := SCRAMBLE;
// 		size := UDINT#40;
// 		FOR i := 0 TO 19 DO
// 			buf[i] := message[i];
// 			buf[i+20] := scramble2[i];
// 		END_FOR;
// 		step :=20;
// 	END_IF;
//
// 20:	IF sha_mode = 3 THEN
// 		sha_mode := 0; (* Stop *)
// 		FOR i:= 0 TO 19 DO
// 			SCRAMBLE[i] := SCRAMBLE[i] XOR scramble1[i];
// 		END_FOR;
// 		step :=0;
// 		RUN := FALSE;
// 	END_IF;
// END_CASE;
//
// SHA1_STREAM(SIZE:=size, MODE:=sha_mode, BUF:=buf, SHA1:=SCRAMBLE, POS=>pos);
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK MYSQL_CONTROL
// VAR_IN_OUT
// 	COM : MYSQL_COM;
// 	INFO : MYSQL_INFO;
// END_VAR
// VAR_OUTPUT
// END_VAR
// VAR
// 	IP_C : IP_C;
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	IPC : IP_CONTROL;
//
// 	step : INT;
// 	snd_step : INT;
// 	TON_WAIT : TON;
// 	last_step : INT;
// 	idx1 : INT;
// 	idx2 : INT;
// 	idx3 : INT;
// 	dw_tmp : DWORD;
// 	timeout1 : TIME;
// 	timeout2 : TIME;
// 	ip_c_red_disable : BOOL;
// 	auth_password : STRING(64);
// 	auth_run : BOOL;
// 	MYSQL_AUTH : MYSQL_AUTH;
// 	scramble1_buffer : ARRAY [0..19] OF BYTE;
// 	scramble2_buffer : ARRAY [0..19] OF BYTE;
// 	charset_number : BYTE := 16#08;
// 	max_packet_size : DWORD := 1024;
// 	client_flags : DWORD := 16#0003_8601;
//
// 	(*
// 	LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;
// 	*)
//
// END_VAR
//
// CASE step OF
//
// 00:	IF COM.SQL_CON THEN
// 		timeout1 := MAX(T#10s,COM.TIMEOUT);
// 		timeout2 := timeout1 + T#1s;
// 		INFO.SQL_ERROR := '';
// 		COM.ERROR_C := DWORD#0;
// 		COM.ERROR_T := BYTE#0;
// 		URL_DATA := STRING_TO_URL(STR:=COM.SQL_URL,DEFAULT_PROTOCOL:='',DEFAULT_PATH:='/');
// 		step := 10;
// 	END_IF;
//
// 10:	IF DNS_CLIENT.DONE THEN
// 		step := 20;
// 	ELSIF DNS_CLIENT.ERROR <> DWORD#0 THEN
// 		COM.ERROR_C := DNS_CLIENT.ERROR;
// 		COM.ERROR_T := BYTE#1;
// 		step := 980;
// 	END_IF;
//
// 20:	(* set up MySQL control channel *)
// 	IF URL_DATA.PORT = WORD#0 THEN
// 		URL_DATA.PORT := WORD#3306; (* default port *)
// 	END_IF;
// 	IP_C.C_PORT := URL_DATA.PORT; (* enter port number *)
//   	IP_C.C_IP := DNS_CLIENT.IP4; (* enter IP *)
//   	IP_C.C_MODE := BYTE#0; (* Mode: TCP+ACTIV+PORT *)
//   	IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
//   	IP_C.C_ENABLE := TRUE; (* enable connection setup *)
//   	IP_C.R_OBSERVE := TRUE; (* receive monitoring *)
// 	COM.R_BUF.SIZE := UINT#0;
// 	COM.SQL_RCV_STATE := BYTE#0;
// 	step := 30;
//
// 30:	IF COM.SQL_RCV_STATE = BYTE#1 THEN
// 		INFO.SERVER_PROTOCOL_VERSION := COM.R_BUF.BUFFER[4];
//
// 		(* search end marker of string server_version *)
// 		idx2 := UINT_TO_INT(COM.R_BUF.SIZE) - 1;
// 		FOR idx1 := 5 TO idx2 DO
// 			IF COM.R_BUF.BUFFER[idx1] = BYTE#0 THEN EXIT; END_IF;
// 		END_FOR;
// 		idx1 := idx1 + 5; (* skip termination and thread id *)
//
// 		FOR idx2:=0 TO 7 DO	(* Byte 0-7 *)
// 			scramble1_buffer[idx2] := COM.R_BUF.BUFFER[idx1];
// 			idx1 := idx1 + 1;
// 		END_FOR;
// 		idx1 := idx1 + 1;
// 		idx2 := idx1 + 1;
// 		INFO.SERVER_CAPABILITIES := WORD_OF_BYTE(COM.R_BUF.BUFFER[idx2], COM.R_BUF.BUFFER[idx1]);
// 		idx1 := idx1 + 2;
// 		INFO.SERVER_LANGUAGE := COM.R_BUF.BUFFER[idx1];
// 		idx1 := idx1 + 1;
// 		idx2 := idx1 + 1;
// 		INFO.SERVER_STATUS := WORD_OF_BYTE(COM.R_BUF.BUFFER[idx2], COM.R_BUF.BUFFER[idx1]);
// 		idx1 := idx1 + 15;
//
// 		FOR idx2:=8 TO 19 DO	(* Byte 8-19 *)
// 			scramble1_buffer[idx2] := COM.R_BUF.BUFFER[idx1];
// 			idx1 := idx1 + 1;
// 		END_FOR;
//
// 		IF INFO.SERVER_PROTOCOL_VERSION < BYTE#10 THEN  (* Error: protocol version too low *)
// 			COM.ERROR_T := BYTE#2;
// 			COM.ERROR_C := DWORD#1;
// 			step := 980;
// 		ELSIF URL_DATA.PASSWORD = '' THEN
// 			step := 50;
// 		ELSE
// 			auth_password := URL_DATA.PASSWORD;
// 			auth_run := TRUE;
// 			step := 40;
// 		END_IF;
// 	END_IF;
//
// 40:	MYSQL_AUTH(RUN:= auth_run,PASSWORD:= auth_password ,MESSAGE:= scramble1_buffer,SCRAMBLE:= scramble2_buffer);
// 	IF NOT auth_run THEN
// 		step := 50;
// 	END_IF;
//
// 50: (* Client Authentication Packet *)
// 	(* bytes 0-3 are entered later *)
// 	COM.S_BUF.BUFFER[04] := BYTE_OF_DWORD(client_flags, BYTE#0); (* client_flags *)
// 	COM.S_BUF.BUFFER[05] := BYTE_OF_DWORD(client_flags, BYTE#1);
// 	COM.S_BUF.BUFFER[06] := BYTE_OF_DWORD(client_flags, BYTE#2);
// 	COM.S_BUF.BUFFER[07] := BYTE_OF_DWORD(client_flags, BYTE#3);
// 	COM.S_BUF.BUFFER[08] := BYTE_OF_DWORD(max_packet_size, BYTE#0); (* max_packet_size *)
// 	COM.S_BUF.BUFFER[09] := BYTE_OF_DWORD(max_packet_size, BYTE#1);
// 	COM.S_BUF.BUFFER[10] := BYTE_OF_DWORD(max_packet_size, BYTE#2);
// 	COM.S_BUF.BUFFER[11] := BYTE_OF_DWORD(max_packet_size, BYTE#3);
// 	COM.S_BUF.BUFFER[12] := charset_number; (* charset_number *)
// 	FOR idx1 := 13 TO 35 DO (* 23 bytes filler always 0x00.. *)
// 		COM.S_BUF.BUFFER[idx1] := BYTE#0;
// 	END_FOR;
// 	idx3 := LEN(URL_DATA.USER);
// 	FOR idx2 :=1 TO idx3 DO
// 		COM.S_BUF.BUFFER[idx1] := CODE(URL_DATA.USER, idx2);
// 		idx1 := idx1 + 1;
// 	END_FOR;
// 	COM.S_BUF.BUFFER[idx1] := BYTE#0; (* enter string termination *)
// 	idx1 := idx1 + 1;
// 	IF LEN(URL_DATA.PASSWORD) = 0 THEN
// 		COM.S_BUF.BUFFER[idx1] := BYTE#0; (* scramble_buff length *)
// 	ELSE
// 		COM.S_BUF.BUFFER[idx1] := BYTE#20; (* scramble_buff length *)
// 		FOR idx2 :=0 TO 19 DO
// 			idx1 := idx1 + 1;
// 			COM.S_BUF.BUFFER[idx1] := scramble2_buffer[idx2];
// 		END_FOR;
// 	END_IF;
// 	COM.S_BUF.SIZE  := INT_TO_UINT(idx1 + 1); (* enter send length *)
//
// 	step := 60;
//
// 60:	IF COM.SQL_RCV_STATE = BYTE#1 THEN
// 		INFO.SQL_CONNECTED := TRUE;
// 		step := 300;
// 	END_IF;
//
// (* ================================================== *)
//
// 300:	step := SEL(COM.SQL_CON,700,310); (* step selection *)
//
// 310:	step := 300; (* dummy step *)
//
// (* ================================================== *)
//
// 700: (* COM_QUIT *)
// 	COM.S_BUF.BUFFER[04] := BYTE#1; (* 0x01 COM_QUIT mysql_close *)
// 	COM.S_BUF.SIZE := UINT#5; (* enter send length *)
// 	COM.SQL_PACKET_NO := BYTE#255;
// 	ip_c_red_disable := TRUE;
// 	step := 710;
//
// 710:	IF IP_C.C_STATE = BYTE#0 THEN (* connection closed by SQL server? *)
// 		step := 980;
// 	END_IF;
//
// (* ================================================== *)
//
// 980:IP_C.C_ENABLE := FALSE; (* close MYSQL control channel *)
// 	ip_c_red_disable := FALSE;
// 	INFO.SQL_CONNECTED := FALSE;
// 	COM.SQL_CON := FALSE;
// 	COM.S_BUF.SIZE := UINT#0;
// 	COM.R_BUF.SIZE := UINT#0;
// 	snd_step := 0;
// 	step := 990;
// 	(*####################  DEBUG-MESSAGE  ###################################*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'QUIT: ERROR CODE >~1< TYPE >~2<';
// 		LOG_CL.PRINTF[1] := DWORD_TO_STRH(COM.ERROR_C);
// 		LOG_CL.PRINTF[2] := BYTE_TO_STRING(COM.ERROR_T);
// 		LOG_MSG();
// 	END_IF;*)
// 	(*########################################################################*)
//
// 990:IF IP_C.C_STATE = BYTE#0 THEN
// 		step := 0;
// 	END_IF;
//
// END_CASE;
//
// (* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// IF step >= 30 THEN
// 	IF COM.R_BUF.SIZE >= UINT#3 AND COM.SQL_RCV_STATE = BYTE#0 THEN (* data received? *)
// 		IF DWORD_TO_UINT(DWORD_OF_BYTE(BYTE#0, COM.R_BUF.BUFFER[2], COM.R_BUF.BUFFER[1], COM.R_BUF.BUFFER[0])) + UINT#4 = COM.R_BUF.SIZE THEN
// 			COM.SQL_PACKET_NO := COM.R_BUF.BUFFER[3];
// 			IF COM.R_BUF.BUFFER[4] = BYTE#16#FF THEN
// 				idx2 := 13; (* error string - first character *)
// 				idx3 := UINT_TO_INT(COM.R_BUF.SIZE) - 1; (* error string - last character *)
// 				FOR idx1 := idx2 TO idx3 DO (* replace all ' with ´ *)
// 					IF COM.R_BUF.BUFFER[idx1] = BYTE#39 THEN
// 						COM.R_BUF.BUFFER[idx1] := BYTE#96;
// 					END_IF;
// 				END_FOR;
// 				INFO.SQL_ERROR := BUFFER_TO_STRING(PT:=REF(COM.R_BUF.BUFFER),SIZE:=COM.R_BUF.SIZE,START:=INT_TO_UINT(idx2),STOP:=INT_TO_UINT(idx3));
// 				COM.ERROR_T := BYTE#3;
// 				COM.ERROR_C := DWORD_OF_BYTE(BYTE#0,BYTE#0,COM.R_BUF.BUFFER[6], COM.R_BUF.BUFFER[5]);
// 				COM.SQL_RCV_STATE := BYTE#2;
// 				COM.R_BUF.SIZE := UINT#0; (* reset receive buffer after error *)
// 				step := 980;
// 			ELSE
// 				COM.SQL_RCV_STATE := BYTE#1;
// 				IP_C.R_OBSERVE := FALSE;
//
// 				(*####################  DEBUG-MESSAGE  ###################################*)
// 				(*IF _debug_enable THEN
// 					LOG_CL.NEW_MSG := 'RCV: size >~1<';
// 					LOG_CL.PRINTF[1] := UINT_TO_STRING(COM.R_BUF.SIZE);
// 					LOG_MSG();
// 				END_IF;*)
// 				(*########################################################################*)
// 			END_IF;
// 		END_IF;
// 	END_IF;
// END_IF;
//
// (* automatic data sending *)
//
// CASE snd_step OF
//
// 00:	IF step >= 20 AND COM.S_BUF.SIZE > UINT#0 THEN
// 		dw_tmp := UINT_TO_DWORD(COM.S_BUF.SIZE - UINT#4); (* packet_length *)
// 		COM.S_BUF.BUFFER[0] := BYTE_OF_DWORD(dw_tmp, BYTE#0);
// 		COM.S_BUF.BUFFER[1] := BYTE_OF_DWORD(dw_tmp, BYTE#1);
// 		COM.S_BUF.BUFFER[2] := 0;
// 		COM.SQL_PACKET_NO := COM.SQL_PACKET_NO + 1;
// 		COM.S_BUF.BUFFER[3] := COM.SQL_PACKET_NO;
//
// 		IP_C.R_OBSERVE := TRUE;
// 		COM.R_BUF.SIZE := UINT#0; (* reset receive buffer *)
// 		COM.SQL_RCV_STATE := BYTE#0;
//
// 	        (*####################  DEBUG-MESSAGE  ###################################*)
// 	        (*IF _debug_enable THEN
// 	          LOG_CL.NEW_MSG := 'SND: size >~1<';
// 	          LOG_CL.PRINTF[1] := UINT_TO_STRING(COM.S_BUF.SIZE);
// 	          LOG_MSG();
// 	        END_IF;*)
// 	        (*########################################################################*)
//
// 		snd_step := 10;
// 	END_IF;
//
// 10:	IF COM.S_BUF.SIZE = UINT#0 THEN (* wait until all data is sent *)
// 		snd_step := 0;
// 	END_IF;
//
// END_CASE;
//
// (* global error handling for main routine *)
// IF COM.ERROR_T = BYTE#0 AND step > 20 THEN (* no error stored yet and main routine active *)
// 	IF IP_C.ERROR > DWORD#0 AND IP_C.C_ENABLE AND NOT IP_C.TIME_RESET THEN  (* IP Control *)
// 		IF ip_c_red_disable = FALSE OR IP_C.ERROR <> DWORD#16#FD000000 THEN
// 			COM.ERROR_C := IP_C.ERROR;
// 			COM.ERROR_T := BYTE#2;
// 			step := 980;
// 		END_IF;
// 	END_IF;
// END_IF;
//
// IF TON_WAIT.Q THEN (* sequence fault - step chain *)
// 	COM.ERROR_C := INT_TO_DWORD(step); (* enter current step *)
// 	COM.ERROR_T := BYTE#5;
// 	step := 980;
// END_IF;
//
// (* ######################################################## *)
//
// DNS_CLIENT(IP_C:=IP_C, S_BUF:=COM.S_BUF, R_BUF:=COM.R_BUF, DOMAIN:=URL_DATA.DOMAIN, IP4_DNS:= COM.DNS_IP4, ACTIVATE:=step=10); (* DNS-Client *)
// IPC(IP:= DWORD#0, PORT:= WORD#0 , TIME_OUT:= timeout1, IP_C:= IP_C, S_BUF:= COM.S_BUF, R_BUF:= COM.R_BUF); (* MYSQL control channel *)
//
// TON_WAIT(IN:=step = last_step AND step > 0,PT:=timeout2);
// last_step := step;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK NET_VAR_BOOL8
VAR_IN_OUT
	X : NET_VAR_DATA;
END_VAR
VAR_INPUT
	IN1 : BOOL;
	IN2 : BOOL;
	IN3 : BOOL;
	IN4 : BOOL;
	IN5 : BOOL;
	IN6 : BOOL;
	IN7 : BOOL;
	IN8 : BOOL;
END_VAR
VAR_OUTPUT
	OUT1 : BOOL;
	OUT2 : BOOL;
	OUT3 : BOOL;
	OUT4 : BOOL;
	OUT5 : BOOL;
	OUT6 : BOOL;
	OUT7 : BOOL;
	OUT8 : BOOL;
	ID : BYTE;
END_VAR
VAR
	init : BOOL;
	index : INT;
	index2 : INT;
	BYTE_TO_BITS : BYTE_TO_BITS;
END_VAR

IF NOT init THEN
	init := TRUE;
	X.ID_MAX := X.ID_MAX + USINT#1;
	ID := USINT_TO_BYTE(X.ID_MAX);
END_IF;

IF X.STATE > BYTE#0 AND X.ERROR_ID = BYTE#0 THEN
	index := X.index; (* PCWORX Workaround *)
	IF INT_TO_UINT(index + 3) > X.BUF_SIZE THEN
		X.ERROR_ID := ID;
	ELSE
		IF X.STATE = BYTE#1 THEN (* write data *)
			X.S_BUF.BUFFER[index] := ID; (* ID *)
			index := index + 1;
			X.S_BUF.BUFFER[index] := BYTE#4; (* TYPE BOOL8 *)
			index := index + 1;
			X.S_BUF.BUFFER[index] := BYTE_OF_BIT(IN1,IN2,IN3,IN4,IN5,IN6,IN7,IN8);
			index := index + 1;

		ELSIF X.STATE = BYTE#2 THEN (* read data *)
			index2 := index + 1; (* PCWORX Workaround *)
			IF X.R_BUF.BUFFER[index] <> ID OR X.R_BUF.BUFFER[index2] <> BYTE#4 THEN
				X.ERROR_ID := ID;
			ELSE
				index := index + 2;
				BYTE_TO_BITS(IN:=X.R_BUF.BUFFER[index]);
				OUT1 := BYTE_TO_BITS.B0;
				OUT2 := BYTE_TO_BITS.B1;
				OUT3 := BYTE_TO_BITS.B2;
				OUT4 := BYTE_TO_BITS.B3;
				OUT5 := BYTE_TO_BITS.B4;
				OUT6 := BYTE_TO_BITS.B5;
				OUT7 := BYTE_TO_BITS.B6;
				OUT8 := BYTE_TO_BITS.B7;
				index := index + 1;
			END_IF;
		END_IF;
		X.index := index;
	END_IF;
END_IF;
END_FUNCTION_BLOCK

FUNCTION_BLOCK NET_VAR_BUFFER
VAR_IN_OUT
	BUF_IN : ARRAY[1..64] OF BYTE;
	BUF_OUT : ARRAY[1..64] OF BYTE;
	X : NET_VAR_DATA;
END_VAR
VAR_OUTPUT
	ID : BYTE;
END_VAR
VAR
	size : INT;
	init : BOOL;
	index : INT;
	index2 : INT;
END_VAR

IF NOT init THEN
	init := TRUE;
	X.ID_MAX := X.ID_MAX + USINT#1;
	ID := USINT_TO_BYTE(X.ID_MAX);
END_IF;

IF X.STATE > BYTE#0 AND X.ERROR_ID = BYTE#0 THEN
	index := X.index; (* PCWORX Workaround *)
	IF X.STATE = BYTE#1 THEN (* write data *)
		size := SIZEOF(BUF_IN);
		IF INT_TO_UINT(index + 3 + size) > X.BUF_SIZE THEN
			X.ERROR_ID := ID;
		ELSE
			X.S_BUF.BUFFER[index] := ID; (* ID *)
			index := index + 1;
			X.S_BUF.BUFFER[index] := BYTE#20; (* TYPE BUFFER *)
			index := index + 1;
			X.S_BUF.BUFFER[index] := INT_TO_BYTE(size); (* number of bytes *)
				index := index + 1;
			FOR index2 := 1 TO size DO
				X.S_BUF.BUFFER[index] := BUF_IN[index2];
				index := index + 1;
			END_FOR;
		END_IF;
	ELSIF X.STATE = BYTE#2 THEN (* read data *)
		index2 := index + 1; (* PCWORX Workaround *)
		IF X.R_BUF.BUFFER[index] <> ID OR X.R_BUF.BUFFER[index2] <> BYTE#20 THEN
			X.ERROR_ID := ID;
		ELSE
			index := index + 2;
			size := BYTE_TO_INT(X.R_BUF.BUFFER[index]);
			index := index + 1;
			FOR index2 := 1 TO size DO
				BUF_OUT[index2] := X.R_BUF.BUFFER[index];
				index := index + 1;
			END_FOR;
		END_IF;
	END_IF;
	X.index := index;
END_IF;
END_FUNCTION_BLOCK

// Commented out for beedance: beedance cannot yet tell SIZEOF(X.S_BUF.BUFFER), the size of a member of the VAR_IN_OUT X.
// FUNCTION_BLOCK NET_VAR_CONTROL
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	MASTER : BOOL;
// 	UDP : BOOL;
// 	REMOTE_IP4 : DWORD;
// 	REMOTE_PORT : WORD := 10000;
// 	SCAN_TIME : TIME := T#1s;
// 	WATCHDOG : TIME := T#2s;
// END_VAR
// VAR_OUTPUT
// 	RUN : BOOL;
// 	ERROR : DWORD;
// END_VAR
// VAR_IN_OUT
// 	X : NET_VAR_DATA;
// END_VAR
// VAR
// 	IPC : IP_CONTROL;
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	step : INT;
// 	activate_last : BOOL;
// 	tscan : TON;
// 	reset : BOOL;
// END_VAR
//
// IF X.ERROR_ID > BYTE#0 THEN (* severe errors in variable assignment present  !! *)
// 	X.STATE := BYTE#0;
// 	ERROR := BYTE_TO_DWORD(X.ERROR_ID);
// 	X.ERROR_ID := BYTE#0;
// 	step := 999;
// END_IF;
//
// CASE step OF
//
// 00:	IF ACTIVATE AND NOT activate_last OR reset THEN
//
//         (*  TCP + MASTER = C_MODE 0 (TCP + ACTIVE + PORT + IP) *)
//         (*  UDP + MASTER = C_MODE 1 (UDP + ACTIVE + PORT + IP) *)
//         (*  TCP + SLAVE  = C_MODE 2 (TCP + PASSIVE + PORT + IP) *)
//         (*  UDP + SLAVE  = C_MODE 3 (UDP + PASSIVE + PORT + IP) *)
// 	IF MASTER THEN
// 		IP_C.C_MODE := BYTE#0;
// 		step := 100;
// 	ELSE
// 		IP_C.C_MODE := BYTE#2;
// 		step := 200;
// 	END_IF;
//         IF UDP THEN IP_C.C_MODE := IP_C.C_MODE OR BYTE#1; END_IF;
// 	  	IP_C.C_PORT := REMOTE_PORT; (* enter port number *)
// 	  	IP_C.C_IP := REMOTE_IP4; (* enter IP *)
// 	  	IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
// 	  	IP_C.C_ENABLE := TRUE; (* enable connection setup *)
// 	  	IP_C.R_OBSERVE := TRUE; (* receive monitoring *)
// 		X.S_BUF.SIZE := UINT#0;
// 		X.R_BUF.SIZE := UINT#0;
// 		X.CYCLE := UDINT#0;
// 		X.ERROR_ID := BYTE#0;
// 		reset := FALSE;
// 	END_IF;
//
// (* ++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// 100:IF IP_C.C_STATE > BYTE#127 THEN (* connection present *)
// 		X.STATE := BYTE#1; (* collect data *)
// 		X.BUF_SIZE := SIZEOF(X.S_BUF.BUFFER);
// 		X.S_BUF.BUFFER[0] := UDINT_TO_BYTE(X.CYCLE);
// 		X.INDEX := 1;
// 		step := 110;
// 	END_IF;
//
// 110:X.STATE := BYTE#0;
// 	X.S_BUF.SIZE := INT_TO_UINT(X.INDEX); (* enter data length *)
// 	step := 120;
//
// 120:IF X.S_BUF.SIZE = UINT#0 AND X.R_BUF.SIZE >= UINT#1 THEN (* check length of received data *)
// 		X.STATE := BYTE#2; (* output data *)
// 		X.BUF_SIZE := X.R_BUF.SIZE;
// 		X.INDEX := 1;
// 		X.ERROR_ID := BYTE#0;
// 		X.CYCLE := X.CYCLE + UDINT#1;
// 		ERROR := DWORD#0; (* clear last error code after a successful cycle *)
// 		RUN := TRUE;
// 		step := 130;
// 	END_IF;
//
// 130:X.STATE := BYTE#0;
// 	X.R_BUF.SIZE := UINT#0;
// 	IF NOT ACTIVATE THEN
// 		step := 999; (* exit *)
// 	ELSIF tscan.Q THEN
// 		step := 100;
// 	END_IF;
//
// (* ++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// 200:IF X.R_BUF.SIZE >= UINT#1 THEN
// 		X.STATE := BYTE#2; (* output data *)
// 		X.BUF_SIZE := X.R_BUF.SIZE;
// 		X.INDEX := 1;
// 		step := 210;
// 	END_IF;
//
// 210:X.STATE := BYTE#1; (* collect data *)
// 	X.BUF_SIZE := SIZEOF(X.R_BUF.BUFFER);;
// 	X.S_BUF.BUFFER[0] := X.R_BUF.BUFFER[0]; (* take cycle *)
// 	X.INDEX := 1;
// 	X.R_BUF.SIZE := UINT#0;
// 	step := 220;
//
// 220:X.STATE := BYTE#0;
// 	X.S_BUF.SIZE := INT_TO_UINT(X.INDEX); (* enter data length *)
// 	step := 230;
//
// 230:IF X.S_BUF.SIZE = UINT#0 THEN
// 		ERROR := DWORD#0; (* clear last error code after a successful cycle *)
// 		RUN := TRUE;
// 		IF NOT ACTIVATE THEN
// 			step := 999; (* exit *)
// 		ELSE
// 			step := 200;
// 		END_IF;
// 	END_IF;
//
// 999:IP_C.C_ENABLE := FALSE; (* exit *)
// 	RUN := FALSE;
// 	step := 0;
//
// END_CASE;
//
// (* ++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// IPC(IP:= DWORD#0, PORT:= WORD#0 , TIME_OUT:= WATCHDOG, IP_C:= IP_C, S_BUF:= X.S_BUF, R_BUF:= X.R_BUF);
//
// IF IP_C.ERROR > DWORD#0 THEN
// 	ERROR := IP_C.ERROR;
// 	RUN := FALSE;
// 	reset := TRUE;
// 	step := 0;
// END_IF;
//
// ACTIVATE_LAST := ACTIVATE; (* save last status for edge detection *)
// tscan(IN:= step <> 110 ,PT:=SCAN_TIME);
// END_FUNCTION_BLOCK

FUNCTION_BLOCK NET_VAR_DWORD8
VAR_IN_OUT
	X : NET_VAR_DATA;
END_VAR
VAR_INPUT
	IN1 : DWORD;
	IN2 : DWORD;
	IN3 : DWORD;
	IN4 : DWORD;
	IN5 : DWORD;
	IN6 : DWORD;
	IN7 : DWORD;
	IN8 : DWORD;
END_VAR
VAR_OUTPUT
	OUT1 : DWORD;
	OUT2 : DWORD;
	OUT3 : DWORD;
	OUT4 : DWORD;
	OUT5 : DWORD;
	OUT6 : DWORD;
	OUT7 : DWORD;
	OUT8 : DWORD;
	ID : BYTE;
END_VAR
VAR
	init : BOOL;
	i : INT;
	i2 : INT;
	index : INT;
	index2 : INT;
	ra : ARRAY[1..8] OF DWORD;
END_VAR

IF NOT init THEN
	init := TRUE;
	X.ID_MAX := X.ID_MAX + USINT#1;
	ID := USINT_TO_BYTE(X.ID_MAX);
END_IF;

IF X.STATE > BYTE#0 AND X.ERROR_ID = BYTE#0 THEN
	index := X.index; (* PCWORX Workaround *)
	IF INT_TO_UINT(index + 34) > X.BUF_SIZE THEN  (* ID (1) + TYPE (1) + 8 x DWORD (4) *)
		X.ERROR_ID := ID;
	ELSE
		IF X.STATE = BYTE#1 THEN (* write data *)
			ra[1] := IN1;
			ra[2] := IN2;
			ra[3] := IN3;
			ra[4] := IN4;
			ra[5] := IN5;
			ra[6] := IN6;
			ra[7] := IN7;
			ra[8] := IN8;
			X.S_BUF.BUFFER[index] := ID; (* ID *)
			index := index + 1;
			X.S_BUF.BUFFER[index] := BYTE#13; (* TYPE DWORD8 *)
			index := index + 1;
			FOR i2 := 1 TO 8 DO;
				FOR i := 1 TO 4 DO;
					X.S_BUF.BUFFER[index] := DWORD_TO_BYTE(ra[i2]);
					ra[i2] := ROR(ra[i2],8);
					index := index + 1;
				END_FOR;
			END_FOR;

		ELSIF X.STATE = BYTE#2 THEN (* read data *)
			index2 := index + 1; (* PCWORX Workaround *)
			IF X.R_BUF.BUFFER[index] <> ID OR X.R_BUF.BUFFER[index2] <> BYTE#13 THEN (* TYPE DWORD8 *)
				X.ERROR_ID := ID;
			ELSE
				index := index + 2;
				FOR i2 := 1 TO 8 DO;
					ra[i2] := DWORD#0;
					FOR i := 1 TO 4 DO;
						ra[i2] := ra[i2] OR BYTE_TO_DWORD(X.R_BUF.BUFFER[index]);
						ra[i2] := ROR(ra[i2],8);
						index := index + 1;
					END_FOR;
				END_FOR;
				OUT1 := ra[1];
				OUT2 := ra[2];
				OUT3 := ra[3];
				OUT4 := ra[4];
				OUT5 := ra[5];
				OUT6 := ra[6];
				OUT7 := ra[7];
				OUT8 := ra[8];
			END_IF;
		END_IF;
		X.index := index;
	END_IF;
END_IF;
END_FUNCTION_BLOCK

// Commented out for beedance: it uses REAL_TO_DW, which is commented out.
// FUNCTION_BLOCK NET_VAR_REAL8
// VAR_IN_OUT
// 	X : NET_VAR_DATA;
// END_VAR
// VAR_INPUT
// 	IN1 : REAL;
// 	IN2 : REAL;
// 	IN3 : REAL;
// 	IN4 : REAL;
// 	IN5 : REAL;
// 	IN6 : REAL;
// 	IN7 : REAL;
// 	IN8 : REAL;
// END_VAR
// VAR_OUTPUT
// 	OUT1 : REAL;
// 	OUT2 : REAL;
// 	OUT3 : REAL;
// 	OUT4 : REAL;
// 	OUT5 : REAL;
// 	OUT6 : REAL;
// 	OUT7 : REAL;
// 	OUT8 : REAL;
// 	ID : BYTE;
// END_VAR
// VAR
// 	init : BOOL;
// 	i : INT;
// 	i2 : INT;
// 	index : INT;
// 	index2 : INT;
// 	r_dw : DWORD;
// 	ra : ARRAY[1..8] OF REAL;
// END_VAR
//
// IF NOT init THEN
// 	init := TRUE;
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	ID := USINT_TO_BYTE(X.ID_MAX);
// END_IF;
//
// IF X.STATE > BYTE#0 AND X.ERROR_ID = BYTE#0 THEN
// 	index := X.index; (* PCWORX Workaround *)
// 	IF INT_TO_UINT(index + 34) > X.BUF_SIZE THEN (* ID (1) + TYPE (1) + 8 x REAL (4) *)
// 		X.ERROR_ID := ID;
// 	ELSE
// 		IF X.STATE = BYTE#1 THEN (* write data *)
// 			ra[1] := IN1;
// 			ra[2] := IN2;
// 			ra[3] := IN3;
// 			ra[4] := IN4;
// 			ra[5] := IN5;
// 			ra[6] := IN6;
// 			ra[7] := IN7;
// 			ra[8] := IN8;
// 			X.S_BUF.BUFFER[index] := ID; (* ID *)
// 			index := index + 1;
// 			X.S_BUF.BUFFER[index] := BYTE#15; (* TYPE REAL8 *)
// 			index := index + 1;
// 			FOR i2 := 1 TO 8 DO;
// 				r_dw := REAL_TO_DW(X:=ra[i2]);
// 				FOR i := 1 TO 4 DO;
// 					X.S_BUF.BUFFER[index] := DWORD_TO_BYTE(r_dw);
// 					r_dw := ROR(r_dw,8);
// 					index := index + 1;
// 				END_FOR;
// 			END_FOR;
//
// 		ELSIF X.STATE = BYTE#2 THEN (* read data *)
// 			index2 := index + 1; (* PCWORX Workaround *)
// 			IF X.R_BUF.BUFFER[index] <> ID OR X.R_BUF.BUFFER[index2] <> BYTE#15 THEN (* TYPE REAL8 *)
// 				X.ERROR_ID := ID;
// 			ELSE
// 				index := index + 2;
// 				FOR i2 := 1 TO 8 DO;
// 					r_dw := DWORD#0;
// 					FOR i := 1 TO 4 DO;
// 						r_dw := r_dw OR BYTE_TO_DWORD(X.R_BUF.BUFFER[index]);
// 						r_dw := ROR(r_dw,8);
// 						index := index + 1;
// 					END_FOR;
// 					ra[i2] := DW_TO_REAL(X:=r_dw);
// 				END_FOR;
// 				OUT1 := ra[1];
// 				OUT2 := ra[2];
// 				OUT3 := ra[3];
// 				OUT4 := ra[4];
// 				OUT5 := ra[5];
// 				OUT6 := ra[6];
// 				OUT7 := ra[7];
// 				OUT8 := ra[8];
// 			END_IF;
// 		END_IF;
// 		X.index := index;
// 	END_IF;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK NET_VAR_STRING
// VAR_IN_OUT
// 	IN : STRING(STRING_LENGTH);
// 	OUT : STRING(STRING_LENGTH);
// 	X : NET_VAR_DATA;
// END_VAR
// VAR_OUTPUT
// 	ID : BYTE;
// END_VAR
// VAR
// 	init : BOOL;
// 	len_str : INT;
// 	index : INT;
// 	index2 : INT;
// END_VAR
//
// IF NOT init THEN
// 	init := TRUE;
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	ID := USINT_TO_BYTE(X.ID_MAX);
// END_IF;
//
// IF X.STATE > BYTE#0 AND X.ERROR_ID = BYTE#0 THEN
// 	index := X.index; (* PCWORX Workaround *)
// 	IF X.STATE = BYTE#1 THEN (* write data *)
// 		len_str := LEN(IN);
// 		IF INT_TO_UINT(index + 3 + len_str) > X.BUF_SIZE THEN
// 			X.ERROR_ID := ID;
// 		ELSE
// 			X.S_BUF.BUFFER[index] := ID; (* ID *)
// 			index := index + 1;
// 			X.S_BUF.BUFFER[index] := BYTE#6; (* TYPE REAL *)
// 			index := index + 1;
// 			X.S_BUF.BUFFER[index] := INT_TO_BYTE(len_str); (* number of bytes *)
// 			index := index + 1;
// 			_STRING_TO_BUFFER(STR:=IN,POS:=index,PT:=REF(X.S_BUF.BUFFER),SIZE:=X.BUF_SIZE);
// 			index := index + len_str;
// 		END_IF;
// 	ELSIF X.STATE = BYTE#2 THEN (* read data *)
// 		index2 := index + 1; (* PCWORX Workaround *)
// 		IF X.R_BUF.BUFFER[index] <> ID OR X.R_BUF.BUFFER[index2] <> BYTE#6 THEN
// 			X.ERROR_ID := ID;
// 		ELSE
// 			index := index + 2;
// 			len_str := BYTE_TO_INT(X.R_BUF.BUFFER[index]);
// 			index2 := index + 1;
// 			index := index2 + len_str;
// 			IF len_STR > 0 THEN
// 				OUT := BUFFER_TO_STRING(PT:=REF(X.R_BUF.BUFFER),SIZE:=X.BUF_SIZE,START:=INT_TO_UINT(index2),STOP:=INT_TO_UINT(index - 1));
// 			ELSE
// 				OUT := '';
// 			END_IF;
// 		END_IF;
// 	END_IF;
// 	X.index := index;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses REAL_TO_DW, which is commented out.
// FUNCTION_BLOCK NET_VAR_X8
// VAR_IN_OUT
// 	X : NET_VAR_DATA;
// END_VAR
// VAR_INPUT
// 	IN_REAL1 : REAL;
// 	IN_REAL2 : REAL;
// 	IN_DINT1 : DINT;
// 	IN_DINT2 : DINT;
// 	IN_UDINT1 : UDINT;
// 	IN_UDINT2 : UDINT;
// 	IN_DWORD1 : DWORD;
// 	IN_DWORD2 : DWORD;
// END_VAR
// VAR_OUTPUT
// 	OUT_REAL1 : REAL;
// 	OUT_REAL2: REAL;
// 	OUT_DINT1 : DINT;
// 	OUT_DINT2 : DINT;
// 	OUT_UDINT1 : UDINT;
// 	OUT_UDINT2 : UDINT;
// 	OUT_DWORD1 : DWORD;
// 	OUT_DWORD2 : DWORD;
// 	ID : BYTE;
// END_VAR
// VAR
// 	init : BOOL;
// 	i : INT;
// 	i2 : INT;
// 	index : INT;
// 	index2 : INT;
// 	ra : ARRAY[1..8] OF DWORD;
// END_VAR
//
// IF NOT init THEN
// 	init := TRUE;
// 	X.ID_MAX := X.ID_MAX + USINT#1;
// 	ID := USINT_TO_BYTE(X.ID_MAX);
// END_IF;
//
// IF X.STATE > BYTE#0 AND X.ERROR_ID = BYTE#0 THEN
// 	index := X.index; (* PCWORX Workaround *)
// 	IF INT_TO_UINT(index + 34) > X.BUF_SIZE THEN  (* ID (1) + TYPE (1) + 8 x DWORD (4) *)
// 		X.ERROR_ID := ID;
// 	ELSE
// 		IF X.STATE = BYTE#1 THEN (* write data *)
// 			ra[1] := REAL_TO_DW(X:=IN_REAL1);
// 			ra[2] := REAL_TO_DW(X:=IN_REAL2);
// 			ra[3] := DINT_TO_DWORD(IN_DINT1);
// 			ra[4] := DINT_TO_DWORD(IN_DINT2);
// 			ra[5] := UDINT_TO_DWORD(IN_UDINT1);
// 			ra[6] := UDINT_TO_DWORD(IN_UDINT2);
// 			ra[7] := IN_DWORD1;
// 			ra[8] := IN_DWORD2;
// 			X.S_BUF.BUFFER[index] := ID; (* ID *)
// 			index := index + 1;
// 			X.S_BUF.BUFFER[index] := BYTE#21; (* TYPE X8 *)
// 			index := index + 1;
// 			FOR i2 := 1 TO 8 DO;
// 				FOR i := 1 TO 4 DO;
// 					X.S_BUF.BUFFER[index] := DWORD_TO_BYTE(ra[i2]);
// 					ra[i2] := ROR(ra[i2],8);
// 					index := index + 1;
// 				END_FOR;
// 			END_FOR;
//
// 		ELSIF X.STATE = BYTE#2 THEN (* read data *)
// 			index2 := index + 1; (* PCWORX Workaround *)
// 			IF X.R_BUF.BUFFER[index] <> ID OR X.R_BUF.BUFFER[index2] <> BYTE#21 THEN (* TYPE X8 *)
// 				X.ERROR_ID := ID;
// 			ELSE
// 				index := index + 2;
// 				FOR i2 := 1 TO 8 DO;
// 					ra[i2] := DWORD#0;
// 					FOR i := 1 TO 4 DO;
// 						ra[i2] := ra[i2] OR BYTE_TO_DWORD(X.R_BUF.BUFFER[index]);
// 						ra[i2] := ROR(ra[i2],8);
// 						index := index + 1;
// 					END_FOR;
// 				END_FOR;
// 				OUT_REAL1 := DW_TO_REAL(X:=ra[1]);
// 				OUT_REAL2 := DW_TO_REAL(X:=ra[2]);
// 				OUT_DINT1 := DWORD_TO_DINT(ra[3]);
// 				OUT_DINT2 := DWORD_TO_DINT(ra[4]);
// 				OUT_UDINT1 := DWORD_TO_UDINT(ra[5]);
// 				OUT_UDINT2 := DWORD_TO_UDINT(ra[6]);
// 				OUT_DWORD1 := ra[7];
// 				OUT_DWORD2 := ra[8];
// 			END_IF;
// 		END_IF;
// 		X.index := index;
// 	END_IF;
// END_IF;
// END_FUNCTION_BLOCK

FUNCTION NETWORK_VERSION : DWORD
VAR_INPUT
	IN : BOOL;
END_VAR

IF in THEN
	NETWORK_VERSION := DATE_TO_DWORD(D#2016-07-01);
ELSE
	NETWORK_VERSION := 135;
END_IF;
END_FUNCTION

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK PRINT_SF
// VAR_IN_OUT
// 	PRINTF_DATA : PRINTF_DATA;
// 	STR : STRING(LOG_SIZE);
// END_VAR
// VAR
// 	pos : INT;
// 	c : INT;
// 	src : STRING(1) := '~';
// 	run : BOOL;
// END_VAR
//
// IF LEN(STR) > 0 THEN
// 	run := TRUE;
//
// 	WHILE (run) DO
// 		pos := FIND(STR,src);
// 		PRINTF_DATA[10] := '';
//
// 		IF pos = 0 THEN
// 			run:= FALSE;
// 		ELSE
// 			c := CODE(STR:=STR,POS:=pos + 1);
// 			IF ISC_NUM(INT_TO_BYTE(c)) THEN
// 				c := c - 48; (* necessary because only one var is allowed inside an index expression *)
// 				PRINTF_DATA[10] := PRINTF_DATA[c]; (* command 1-9 gives parameter index 1-9 *)
// 				(* replace parameter with text *)
// 			END_IF;
// 			IF LEN(PRINTF_DATA[10]) + LEN(STR) > LOG_SIZE THEN
// 				PRINTF_DATA[10] := '..';
// 			END_IF;
// 			(* pcworx: output and input string must be different *)
// 			PRINTF_DATA[11] := STR;
// 			STR := REPLACE(PRINTF_DATA[11],PRINTF_DATA[10],2,pos);
// 		END_IF;
// 	END_WHILE;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK RC4_CRYPT_STREAM
// VAR_IN_OUT
// 	MODE : INT;
// 	KEY : STRING(40);
// 	BUF : ARRAY [0..63] OF BYTE;
// 	SIZE : UDINT;
// END_VAR
//
// VAR_OUTPUT
// 	POS : UDINT;
// END_VAR
//
// VAR
// 	sbox : ARRAY [0..255] OF BYTE;
// 	skey : ARRAY [0..255] OF BYTE;
// 	a : INT;
// 	b : INT;
// 	c : INT;
// 	d : USINT;
// 	e : USINT;
// 	end : UDINT;
// 	v : BYTE;
// END_VAR
//
// CASE MODE OF
//
// 1:	(* initialization *)
// 	b := LEN(KEY) - 1;
// 	IF b < 0 OR SIZE < 1 THEN
// 		MODE := 3;
// 		RETURN;
// 	END_IF;
//
// 	FOR a := 0 TO b DO
// 		skey[a] := CODE(KEY, a + 1);
// 		sbox[a] := INT_TO_BYTE(a);
// 	END_FOR;
//
// 	c := 0;
// 	FOR a := b + 1 TO 255 DO
// 		skey[a] := skey[c];
// 		sbox[a] := INT_TO_BYTE(a);
// 		c := c + 1;
// 		IF c > b THEN c := 0; END_IF;
// 	END_FOR;
//
// 	d := 0;
// 	FOR a := 0 TO 255 DO
// 		d := d + sbox[a] + skey[a];
// 		v := sbox[a];
// 		sbox[a] := sbox[d];
// 		sbox[d] := v;
// 	END_FOR;
//
// 	POS := 0;
// 	end := SIZE;
// 	SIZE := MIN(64,end);
// 	d := 0;
// 	e := 0;
// 	MODE := 2;
//
// 2:	(* encrypt/decrypt data *)
//
// 	b := UDINT_TO_INT(SIZE) - 1;
// 	FOR a := 0 TO b DO
// 		d := d + 1;
// 		e := e + sbox[d];
// 		v := sbox[d];
// 		sbox[d] := sbox[e];
// 		sbox[e] := v;
// 		v := sbox[d] + sbox[e];
// 		BUF[a] :=BUF[a] XOR sbox[v];
// 	END_FOR;
//
// 	POS := POS + SIZE;
// 	SIZE := MIN(64,end - POS);
//
// 	IF SIZE = 0 THEN
// 		MODE := 3;
// 	END_IF;
// END_CASE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses BUFFER_SEARCH, which is commented out.
// FUNCTION_BLOCK READ_HTTP
// VAR_INPUT
// 	SIZE : UINT;
// 	POS : INT;
// 	SRC : STRING;
// END_VAR
// VAR_IN_OUT
// 	PT : NW_BUF_LONG;
// END_VAR
// VAR_OUTPUT
// 	READ_HTTP : STRING(STRING_LENGTH);
// END_VAR
// VAR
// 	index : INT;
// 	s1 : INT;
// 	e1 : INT;
// 	stop : INT;
// END_VAR
//
// stop := UINT_TO_INT(SIZE) - 1;
//
// IF LEN(SRC) = 0 THEN
// 	(* if there is no search term, evaluate the HTTP STATUS *)
// 	SRC := 'HTTP/'; (* new search term *)
// ELSE
// 	(* append the searched parameter text *)
// 	SRC := CONCAT(SRC,': ');
// END_IF;
//
// s1 := BUFFER_SEARCH(pt:=REF(PT),SIZE:=SIZE,STR:=SRC,POS:=POS,IGN:=TRUE);
//
// IF s1 >= 0 THEN (* Text found *)
// 	e1 := s1;
//
// 	(* search start of line *)
// 	index := s1 - 1;
// 	WHILE (index > 0) DO
// 		IF pt[index] < BYTE#32 THEN EXIT; END_IF;
// 		s1 := s1 - 1;
// 		index := index - 1;
// 	END_WHILE;
//
// 	(* search end of line *)
// 	index := e1 + 1;
// 	WHILE (index <= stop) DO
// 		IF pt[index] < BYTE#32 THEN EXIT; END_IF;
// 		e1 := e1 + 1;
// 		index := index + 1;
// 	END_WHILE;
//
// 	s1 := s1 + LEN(SRC); (* do not copy 'search term: ' *)
//
// 	(* copy string *)
// 	READ_HTTP := BUFFER_TO_STRING(PT:=REF(PT),SIZE:=SIZE,START:=s1,STOP:=e1);
// ELSE
// 	READ_HTTP := '';
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK SHA1_STR
// VAR_INPUT
// 	RUN: BOOL;
// END_VAR
// VAR_OUTPUT
// 	DONE: BOOL;
// END_VAR
// VAR_IN_OUT
// 	STR: STRING(STRING_LENGTH);
// 	SHA1 : ARRAY [0..19] OF BYTE;
// END_VAR
//
// VAR
// 	run_last: BOOL;
// 	SHA1_STREAM : SHA1_STREAM;
// 	buf : ARRAY [0..63] OF BYTE;
// 	mode : INT;
// 	size : UDINT;
// 	pos : UDINT;
// END_VAR
//
// CASE mode OF
// 0:	(* wait for start *)
// 	IF RUN AND NOT run_last THEN
// 		DONE := FALSE;
// 		mode := 1; (* initialization *)
// 		size := LEN(STR);
// 	END_IF;
//
// 2:	(* Data copy and Calc *)
// 	IF size > 0 THEN
// 		_STRING_TO_BUFFER(STR:=MID(STR, UDINT_TO_INT(SIZE), UDINT_TO_INT(POS) + 1), POS:=0, PT:=REF(buf), SIZE:=SIZEOF(buf));
// 	END_IF;
//
// 3:	DONE := TRUE;
// 	mode := 0; (* Stop *)
//
// END_CASE;
//
// IF mode > 0 THEN
// 	SHA1_STREAM(SIZE:=size, MODE:=mode, BUF:=buf, SHA1:=SHA1, POS=>pos);
// END_IF;
//
// run_last := RUN;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK SHA1_STREAM
VAR_IN_OUT
	MODE : INT;
	BUF : ARRAY [0..63] OF BYTE;
	SHA1 : ARRAY [0..19] OF BYTE;
	SIZE : UDINT;
END_VAR

VAR_OUTPUT
	POS : UDINT;
END_VAR

VAR
	hash : ARRAY [0..4] OF DWORD;
	end : UDINT;
	block : UDINT;
	pad_1 : BOOL;
	w : ARRAY [0..79] OF DWORD;
	n : INT;
	n1 : INT;
	a : DWORD;
	b : DWORD;
	c : DWORD;
	d : DWORD;
	e : DWORD;
	f : DWORD;
	k : DWORD;
	x : DWORD;
	z : INT;
END_VAR

CASE MODE OF

1:	(* initialization of the hash values *)
	hash[0] := DWORD#16#67452301;
	hash[1] := DWORD#16#EFCDAB89;
	hash[2] := DWORD#16#98BADCFE;
	hash[3] := DWORD#16#10325476;
	hash[4] := DWORD#16#C3D2E1F0;

	(* initialize buffer *)
	FOR n := 16 TO 79 DO
   		w[n] := DWORD#0;
	END_FOR;

	block := SHR(SIZE,6) + BOOL_TO_UDINT((UDINT_TO_BYTE(SIZE) AND 16#3F) > 55);	(* calculate the required blocks *)

	pad_1 := FALSE;
	POS := 0;
	End := SIZE;
	SIZE := MIN(64,End);
	MODE := 2;

2:	FOR n := UDINT_TO_INT(SIZE) TO 63 DO (* append padding bytes *)
		BUF[n] := 0;
	END_FOR;

	IF SIZE < 64 AND NOT pad_1 THEN
		BUF[SIZE] := 2#1000_0000;
		pad_1 := TRUE;
	END_IF;

	POS := POS + SIZE;
	SIZE := MIN(64,End - POS);

	n := 0;	(* copy byte array into DWORD array *)
	FOR n1 := 0 TO 15 DO
		w[n1] := DWORD_OF_BYTE(BUF[n],BUF[n+1],BUF[n+2],BUF[n+3]);
		n := n + 4;
	END_FOR;

	IF block = 0 THEN (* enter the data length at the end of the last block *)
		w[15] := SHL(UDINT_TO_DWORD(End),3); (* enter the number of data bits *)
		MODE := 3;
	ELSE
		block := block - 1;
	END_IF;

	FOR n := 16 TO 79 DO (* extend the 16 32-bit words to 80 32-bit words *)
		w[n] := ROL(w[n-3] XOR w[n-8] XOR w[n-14] XOR w[n-16],1);
	END_FOR;

	a := hash[0];
	b := hash[1];
	c := hash[2];
	d := hash[3];
	e := hash[4];

	FOR n := 0 TO 79 DO
		IF n <= 19 THEN
			f := (b AND c) OR (NOT(b) AND d);
			k := DWORD#16#5A827999;
		ELSIF n <= 39 THEN
			f := b XOR c XOR d;
			k := DWORD#16#6ED9EBA1;
		ELSIF n <= 59 THEN
			f := (b AND c) OR (b AND d) OR (c AND d);
			k := DWORD#16#8F1BBCDC;
		ELSE
			f := b XOR c XOR d;
			k := DWORD#16#CA62C1D6;
		END_IF;

		x := ROL(a,5) + f + e + k + w[n];
		e := d;
		d := c;
		c := ROL(b,30);
		b := a;
		a := x;
	END_FOR;

	(* add the current hash values to the last value *)
	hash[0] := hash[0] + a;
	hash[1] := hash[1] + b;
	hash[2] := hash[2] + c;
	hash[3] := hash[3] + d;
	hash[4] := hash[4] + e;

END_CASE;

IF MODE = 3 THEN
	z := 0;
	FOR n := 0 TO 4 DO (* convert DWORD to BYTE array *)
		FOR n1 := 0 TO 3 DO
			hash[n] := ROL(hash[n],8);
			SHA1[z] := DWORD_TO_BYTE(hash[n]);
			z := z + 1;
		END_FOR;
	END_FOR;
END_IF;
END_FUNCTION_BLOCK

// Commented out for beedance: it computes with pointers as addresses (pt := pt + 1); beedance references (REF_TO) have no arithmetic.
// FUNCTION SHA1_TO_STRH : STRING(40)
// VAR_INPUT
// 	SHA1 : ARRAY [0..19] OF BYTE;
// END_VAR
// VAR
// 	i : INT;
// 	tmp : BYTE;
// 	pt : REF_TO BYTE;
// END_VAR
// pt := REF(SHA1_TO_STRH); (* read output adress to pointer *)
//
// FOR i := 0 TO 19 DO; (* write the 10 hex characters backwards *)
// 	tmp := SHR(SHA1[i],4);
// 	PT^ := tmp + SEL(tmp <= 9,87,48);
// 	pt := pt + 1; (* decrement the pointer *)
// 	tmp := SHA1[i] AND 2#1111;
// 	PT^ := tmp + SEL(tmp <= 9,87,48);
// 	pt := pt + 1; (* decrement the pointer *)
// END_FOR;
// pt^:= 0; (* write the closing byte for the string *)
// END_FUNCTION

// Commented out for beedance: it has a comment between UNTIL and END_REPEAT, which beedance does not parse yet.
// FUNCTION_BLOCK SMTP_CLIENT
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	TIMEOUT : TIME;
// 	DTI : DT;
// 	DTI_OFFSET : INT;
// 	DNS_IP4 : DWORD;
// END_VAR
// VAR_OUTPUT
// 	DONE : BOOL;
// 	BUSY : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR_IN_OUT
// 	SERVER : STRING;
// 	MAILFROM : STRING;
// 	MAILTO : STRING(STRING_LENGTH);
// 	SUBJECT : STRING;
// 	BODY : STRING(STRING_LENGTH);
// 	FILES : STRING(STRING_LENGTH);
// END_VAR
// VAR
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// 	F_BUF : NETWORK_BUFFER;
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	IPC : IP_CONTROL;
// 	FS : FILE_SERVER;
// 	FSD : FILE_SERVER_DATA;
// 	BASE64_ENCODE_STR : BASE64_ENCODE_STR;
// 	BASE64_ENCODE_STREAM : BASE64_ENCODE_STREAM;
// 	fpd : FILE_PATH_DATA;
// 	MD5_CRAM_AUTH : MD5_CRAM_AUTH;
// 	MD5_RUN : BOOL;
// 	activate_last : BOOL;
// 	step : INT;
// 	snd_step : INT;
// 	str1 : STRING(STRING_LENGTH);
// 	str2 : STRING(STRING_LENGTH);
// 	str3 : STRING;
// 	snd_text :	STRING(STRING_LENGTH);
// 	snd_lf_cnt : INT;
// 	rcv_text :	STRING(STRING_LENGTH);
// 	rcv_s : STRING(3);
// 	rcv_state : INT;
// 	auth_names : STRING := ' PLAIN; LOGIN; CRAM-MD5';
// 	auth_names_tls : STRING := ' PLAIN; LOGIN; CRAM-MD5;STARTTLS';
// 	auth_state : BYTE;
// 	TON_WAIT : TON;
// 	last_step : INT;
// 	next_step : INT;
// 	timeout2 : TIME;
// 	timeout3 : TIME;
// 	last_block : BOOL;
// 	file_delete : BOOL;
// 	esmtp : BOOL;
// 	cnt1 : INT;
// 	cnt2 : INT;
// 	cnt3 : INT;
// 	cnt4 : INT;
// 	idx1 : INT;
// 	idx2 : INT;
// 	idx3 : INT;
// 	snd_idx : INT;
// 	snd_enable : BOOL;
// 	i1 : INT;
// 	b64_done : BOOL;
// 	b64_start : BOOL;
// 	b64_max : UINT := 2880; (* short = 960 / long = 2880 *)
// 	b64_buf1 : ARRAY [0..47] OF BYTE;
// 	b64_buf2 : ARRAY [0..63] OF BYTE;
// 	b64_size1 : INT;
// 	b64_size2 : INT;
// 	b64_lf : BOOL;
// 	b64_str1	: STRING(144);
// 	b64_str2	: STRING(192);
// 	md5_user : STRING(64);
// 	md5_password : STRING(64);
// 	md5_b64_ts : STRING(64);
// 	md5_auth_key : STRING(192);
// 	ip_c_red_disable : BOOL; (* IP_CONTROL - suppress remote disconnect error *)
// 	ssl_mode :	BYTE;(*SSL/TLS Mode -   0 = off / 1 = SSL / 2 = TLS*)
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;*)
// END_VAR
//
// CASE step OF
//
// 00:
// 	IF ACTIVATE AND NOT activate_last THEN
// 		TIMEOUT := MAX(T#10s,TIMEOUT);
// 		timeout2 := TIMEOUT + T#1s;
// 		timeout3 := TIMEOUT / 2;
// 		DONE := FALSE;
// 		ip_c_red_disable := FALSE;
// 		ERROR_C := DWORD#0;
// 		ERROR_T := BYTE#0;
// 		FSD.ERROR := BYTE#0;
// 		file_delete := FALSE;
// 		URL_DATA := STRING_TO_URL(STR:=SERVER,DEFAULT_PROTOCOL:='',DEFAULT_PATH:='/');
// 		idx1 := FIND(URL_DATA.USER,'%'); (* change % character to @ *)
// 		IF idx1 > 0 THEN
// 			URL_DATA.USER := REPLACE(URL_DATA.USER,'@',1,idx1);
// 		END_IF;
// 		esmtp := LEN(URL_DATA.PASSWORD) > 0; (* SMTP or ExtendedSMTP *)
// 		step := 10;
// 	END_IF;
//
// 10:
// 	IF DNS_CLIENT.DONE THEN
// 		step := 20;
// 	ELSIF DNS_CLIENT.ERROR <> DWORD#0 THEN
// 		ERROR_C := DNS_CLIENT.ERROR;
// 		ERROR_T := BYTE#1;
// 		step := 980;
// 	END_IF;
//
// 20:	(* set up SMTP communication *)
// 	IF FIND(';SSL;ssl;Ssl', URL_DATA.PROTOCOL) > 1 THEN
// 		ssl_mode := BYTE#1;
// 		IP_C.C_PORT := SEL(URL_DATA.PORT = WORD#0,URL_DATA.PORT,WORD#465); (* enter port number *)
// 		IP_C.MAILBOX[16] := BYTE#1; (* enable SSL mode *)
// 	ELSIF FIND(';TLS;tls;Tls', URL_DATA.PROTOCOL) > 1 THEN
// 		ssl_mode := BYTE#2;
// 		IP_C.C_PORT := SEL(URL_DATA.PORT = WORD#0,URL_DATA.PORT,WORD#587); (* enter port number *)
// 		IP_C.MAILBOX[16] := BYTE#0; (* no SSL mode yet*)
// 	ELSE
// 		ssl_mode := BYTE#0;
// 		IP_C.C_PORT := SEL(URL_DATA.PORT = WORD#0,URL_DATA.PORT,WORD#25); (* enter port number *)
// 		IP_C.MAILBOX[16] := BYTE#0; (* no SSL mode*)
// 	END_IF;
//
//   	IP_C.C_IP := DNS_CLIENT.IP4; (* enter IP *)
//   	IP_C.C_MODE := BYTE#0; (* Mode: TCP+ACTIV+PORT *)
//   	IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
//   	IP_C.C_ENABLE := TRUE; (* enable connection setup *)
//   	IP_C.R_OBSERVE := TRUE; (* receive monitoring *)
// 	IP_C.MAILBOX[3] := BYTE#0; (* re-enable IP receive *)
// 	R_BUF.SIZE := UINT#0;
// 	snd_enable := TRUE; (* allow immediate sending *)
// 	snd_idx := 0;
// 	step := 30;
//
// 30:
// 	IF rcv_state = 220 THEN (* 220 = Service ready *)
// 		IF esmtp THEN
// 			snd_text := 'EHLO OSCAT'; (* ExtendedSMTP *)
//
// 			str1:= SEL(ssl_mode = BYTE#2,auth_names,auth_names_tls);
// 			auth_state := BYTE#0;
// 			idx2 := 1900;
//
// 			next_step := 40;
// 		ELSE
// 			snd_text := 'HELO OSCAT'; (* SMTP *)
// 			next_step := 90;
// 		END_IF;
// 		snd_lf_cnt := 1;
// 	END_IF;
//
// 40:
// 	IF rcv_state = 250 THEN
// 		(* 250 = Requested mail action okay, complete
// 			1. some providers do not announce an AUTH method with SSL.
// 			2. since some SMTP servers send the 250 xxx messages in more than one telegram,
// 		       a special SSL handling must be done via timeout monitoring instead *)
// 		IF ssl_mode >= BYTE#1 AND TON_WAIT.ET > T#1s THEN
// 			rcv_text := '250-AUTH LOGIN'; (* fake AUTH *)
//
//         	(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'SMTP: PROBLEM /rcv_text ~1';
// 				LOG_CL.PRINTF[1] := rcv_text;
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
//
// 		END_IF;
//
// 		IF FIND(rcv_text,'AUTH') > 0 THEN (* BUGFIX  !*)
// 			IP_C.MAILBOX[3] := BYTE#1; (* block IP receive until data has been checked *)
// 			idx1 := 0;
// 			step := 50;
// 		END_IF;
//
// 	END_IF;
//
// 50:
// 	str2 := ELEMENT_GET(59, idx1, str1); (* separator >;< *)
//
// 	(*####################  DEBUG-MESSAGE  ###################################*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'ST s1>~1< s2>~2< i1>~2< ~3';
// 		LOG_CL.PRINTF[1] := str1;
// 		LOG_CL.PRINTF[2] := str2;
// 		LOG_CL.PRINTF[3] := INT_TO_STRING(idx1);
// 		LOG_CL.PRINTF[4] := rcv_text;
// 		LOG_MSG();
// 	END_IF;*)
// 	(*########################################################################*)
//
// 	IF LEN(str2) > 0 THEN
// 		IF FIND(rcv_text,str2) > 0 THEN
// 			auth_state := BIT_LOAD_B(auth_state, TRUE, idx1);
// 			str3 := str2; (* save current AUTH type name *)
// 			idx2 := 1000 + (idx1 * 100); (* calculate step *)
// 		END_IF;
// 		idx1 := idx1 + 1;
// 	ELSE
// 		IF auth_state <> BYTE#0 THEN (* AUTH type found *)
// 			IP_C.MAILBOX[3] := BYTE#0; (* re-enable IP receive *)
// 			(* mit den jeweiligen Authentifizierung Schritten fortfahren:
// 			   1000 = PLAIN / 1100 = LOGIN / 1200 = CRAM-MD5 / 1300 TLS *)
// 			step := idx2;
// 		ELSE
// 			idx1 := 0;
// 			rcv_text := '';
// 			IP_C.MAILBOX[3] := BYTE#0; (* re-enable IP receive *)
// 			step := 40; (* wait for the next 250 auth message *)
// 		END_IF;
// 	END_IF;
//
// 	(*####################  DEBUG-MESSAGE  ###################################*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'SMTP: /auth ~1 /idx1 ~2 /idx2 ~3 /str3 ~4 /LEN ~5 /step ~6 /str1 ~7';
// 		LOG_CL.PRINTF[1] := BYTE_TO_STRING(auth_state);
// 		LOG_CL.PRINTF[2] := INT_TO_STRING(idx1);
// 		LOG_CL.PRINTF[3] := INT_TO_STRING(idx2);
// 		LOG_CL.PRINTF[4] := str3;
// 		LOG_CL.PRINTF[5] := INT_TO_STRING(LEN(str2));
// 		LOG_CL.PRINTF[6] := INT_TO_STRING(step);
// 		LOG_CL.PRINTF[7] := str1;
// 		LOG_MSG();
// 	END_IF;*)
// 	(*########################################################################*)
//
//  (* ++++++++++++++++++++++++++++++++++++++++++++++++++++++*)
//
// 90:
// 	IF rcv_state = 235 OR rcv_state = 250 THEN (* 235 = Authentication successful / 250 Requested mail action okay *)
// 		str2 := MAILFROM;
// 		str1 := ELEMENT_GET(59, 0, str2);(* read MAILFROM, separator >;< *)
// 		snd_text := INSERT('MAIL FROM: <>', str1, 12);
// 		str3 := REPLACE(snd_text, 'From', 9, 1); (* replace 'MAIL FROM' with "From" and buffer it *)
// 		next_step := 95;
// 		cnt1 := ELEMENT_COUNT(35, MAILTO); (* separator >#<, determine the number of blocks *)
// 		cnt2 := 0;
// 	END_IF;
//
// 95:
// 	IF rcv_state = 250 THEN (* 250 = Requested mail action okay *)
// 		IF cnt2 < cnt1 THEN
// 			str1 := ELEMENT_GET(35, cnt2, MAILTO);(* separator >#< *)
// 			cnt3 := ELEMENT_COUNT(59, str1); (* separator >;<, determine the number of sub-blocks *)
// 			cnt4 := 0;
// 			cnt2 := cnt2 + 1;
// 			step := 100; (* evaluate block content *)
// 		ELSE
// 			step := 110; (* all blocks processed *)
// 		END_IF;
// 	END_IF;
//
// 100:
// 	IF rcv_state = 250 THEN (* 250 = Requested mail action okay *)
// 		IF cnt4 < cnt3 THEN
// 			str2 := ELEMENT_GET(59, cnt4, str1); (* separator >;< *)
// 			snd_text := INSERT('RCPT TO: <>', str2, 10);
// 			cnt4 := cnt4 + 1;
// 			next_step := 100;
// 		ELSE
// 			step := 95;
// 		END_IF;
// 	END_IF;
//
// 110:
// 	IF rcv_state = 250 THEN (* 250 = Requested mail action okay *)
// 		snd_text := 'DATA';
// 		next_step := 120;
// 	END_IF;
//
// 120:
// 	IF rcv_state = 354 THEN (* 354 = Start mail input; end with. *)
// 		snd_enable := FALSE;
// 		str2 := MAILFROM;
// 		str1 := ELEMENT_GET(59, 1, str2);(* read display name, separator >;< *)
// 		IF LEN(str1) > 0 THEN
// 			str2 := INSERT('"" ', str1, 1);
// 			snd_text := INSERT(str3, str2, 6);
// 		ELSE
// 			snd_text := str3;
// 		END_IF;
// 		cnt2 := 0;
// 		cnt1 := MIN(2,cnt1); (* hide blind copy block *)
// 		step := 200;
// 	END_IF;
//
// 200:
// 	IF snd_step = 0 THEN
// 		IF cnt2 < cnt1 THEN
// 			str1 := ELEMENT_GET(35, cnt2, MAILTO);(* separator >#< *)
// 			cnt3 := ELEMENT_COUNT(59, str1); (* separator >;<, determine the number of sub-blocks *)
// 			cnt4 := 0;
// 			cnt2 := cnt2 + 1;
// 			str3 := SEL(cnt2 = 1,'Cc: <>','To: <>');
// 			step := 205; (* evaluate block content *)
// 		ELSE
// 			step := 210; (* all blocks processed *)
// 		END_IF;
// 	END_IF;
//
// 205:
// 	IF snd_step = 0 THEN
// 		IF cnt4 < cnt3 THEN
// 			str2 := ELEMENT_GET(59, cnt4, str1); (* separator >;< *)
// 			snd_text := INSERT(str3, str2, 5);
// 			cnt4 := cnt4 + 1;
// 		ELSE
// 			step := 200;
// 		END_IF;
// 	END_IF;
//
// 210:
// 	IF snd_step = 0 THEN
// 		snd_text := CONCAT('Subject: ',SUBJECT);
// 		step := 220;
// 	END_IF;
//
// 220:
// 	IF snd_step = 0 THEN
// 		DTI := UDINT_TO_DT(DT_TO_UDINT(DTI) + INT_TO_UDINT(- DTI_OFFSET * 60)); (* convert minute offset to seconds and add *)
// 		IF esmtp THEN (* ESMTP Format >Date: Thu, 21 May 1998 05:33:29 -0700< *)
// 			str2 := ';Mon;Tue;Wed;Thu;Fri;Sat;Sun';
// 			str1 := ELEMENT_GET(59, DAY_OF_WEEK(DT_TO_DATE(DTI)), str2);
// 			snd_text := CONCAT('Date: ',str1);
// 			str1 := DT_TO_STRF(DTI, 0, ', #H #E #A #N:#R:#T +0000', 1);
// 			snd_text := CONCAT(snd_text,str1);
// 		ELSE (* SMTP Format >Date: 21 May 98 05:33:29< *)
// 			snd_text := DT_TO_STRF(DTI, 0, 'Date: #G #E #B #N:#R:#T', 1);
// 		END_IF;
// 		step := 230;
// 	END_IF;
//
// 230:
// 	IF snd_step = 0 THEN
// 		snd_text := 'MIME-Version: 1.0';
// 		step := 250;
// 	END_IF;
//
// 250:
// 	IF snd_step = 0 THEN
// 		snd_text := 'Content-Type: multipart/mixed;boundary="x"';
// 		snd_lf_cnt := 2;
// 		step := 260;
// 	END_IF;
//
// 260:
// 	IF snd_step = 0 THEN
// 		snd_enable := TRUE;
// 		snd_text := 'This is a multi-part message in MIME format.';
// 		step := 270;
// 	END_IF;
//
// 270:
// 	IF snd_step = 0 THEN
// 		snd_enable := FALSE;
// 		snd_text := '--x';
// 		snd_lf_cnt := 1;
// 		step := 280;
// 	END_IF;
//
// 280:
// 	IF snd_step = 0 THEN
// 		snd_text := 'Content-Type: text/plain; format=flowed; charset="iso-8859-1";';
// 		snd_text := CONCAT(snd_text,' reply-type=original'); (* PCWORX only supports constants with max 80 characters *)
// 		step := 290;
// 	END_IF;
//
// 290:
// 	IF snd_step = 0 THEN
// 		snd_text := 'Content-Transfer-Encoding: 8bit';
// 		snd_lf_cnt := 2;
// 		step := 300;
// 	END_IF;
//
// 300:
// 	IF snd_step = 0 THEN
// 		snd_enable := TRUE;
// 		snd_text := BODY;
// 		snd_lf_cnt := 2;
// 		step := 310;
// 	END_IF;
//
// 310:
// 	IF snd_step = 0 THEN
// 		cnt1 := ELEMENT_COUNT(59, FILES) - 1; (* number of files *)
// 		step := 400;
// 	END_IF;
//
// (* +++++++++++++ File-Attachment ++++++++++++ *)
// 400:
// 	IF cnt1 >= 0 THEN
// 		str1 := ELEMENT_GET(59, cnt1, FILES);
// 		FILE_PATH_SPLIT(FILENAME:=str1,X:=fpd); (* split file path *)
// 		IF fpd.FILENAME = '#DEL#' THEN
// 			file_delete := TRUE;
// 		ELSE
// 			step := 410;
// 		END_IF;
// 		cnt1 := cnt1 - 1;
// 	ELSE
// 		step := 800;
// 	END_IF;
//
// 410:
// 	IF snd_step = 0 THEN
// 		snd_enable := FALSE;
// 		snd_text := '--x';
// 		snd_lf_cnt := 1;
// 		step := 420;
// 	END_IF;
//
// 420:
// 	IF snd_step = 0 THEN
// 		snd_text := 'Content-Transfer-Encoding: BASE64';
// 		str2 := CONCAT(fpd.FILENAME,'"');
// 		step := 430;
// 	END_IF;
//
// 430:
// 	IF snd_step = 0 THEN
// 		snd_text := CONCAT('Content-Type: application/octet-stream; name="',str2);
// 		step := 440;
// 	END_IF;
//
// 440:
// 	IF snd_step = 0 THEN
// 		snd_enable := TRUE;
// 		snd_text := CONCAT('Content-Disposition: attachment; filename="',str2);
// 		snd_lf_cnt := 2;
// 		step := 450;
// 	END_IF;
//
// 450:
// 	IF snd_step = 0 THEN
// 		FSD.FILENAME := str1; (* file name *)
// 		FSD.OFFSET := UDINT#0; (* start at beginning of file *)
// 		step := 460;
// 	END_IF;
//
// 460:
// 	FSD.MODE := BYTE#1; (* open existing file for reading *)
// 	F_BUF.SIZE := b64_max; (* read maximum length *)
// 	idx1 := 0;
// 	idx2 := 0;
// 	step := 470;
//
// 470:
// 	IF FSD.MODE = BYTE#0 AND FSD.ERROR = BYTE#0 THEN
// 		last_block := FSD.FILE_SIZE = FSD.OFFSET;
// 		b64_size1 := MIN(48,UINT_TO_INT(F_BUF.SIZE) - idx1);
// 		IF b64_size1 > 0 THEN
// 			idx3 := b64_size1 - 1;
// 			FOR i1 := 0 TO idx3 DO (* copy data into BASE64 buffer *)
// 				b64_buf1[i1] := F_BUF.BUFFER[idx1];
// 				idx1 := idx1 + 1;
// 			END_FOR;
// 			BASE64_ENCODE_STREAM( SIZE1:= b64_size1, BUF1:= b64_buf1, BUF2:= b64_buf2, SIZE2=> b64_size2); (* BASE64 encoding *)
// 			idx3 := b64_size2 - 1;
// 			FOR i1 := 0 TO idx3 DO (* copy data out of BASE64 buffer *)
// 				S_BUF.BUFFER[idx2] := b64_buf2[i1];
// 				idx2 := idx2 + 1;
// 			END_FOR;
// 			b64_lf := TRUE;
// 			(*####################  DEBUG-MESSAGE  ###################################*)
// 			(*IF _debug_enable THEN
// 				LOG_CL.NEW_MSG := 'B64: FILE ~1 /SIZE ~2 /OFF ~3 /BSIZE ~4 /idx1 ~5 /idx2 ~6';
// 				LOG_CL.PRINTF[1] := FSD.FILENAME;
// 				LOG_CL.PRINTF[2] := UDINT_TO_STRING(FSD.FILE_SIZE);
// 				LOG_CL.PRINTF[3] := UDINT_TO_STRING(FSD.OFFSET);
// 				LOG_CL.PRINTF[4] := INT_TO_STRING(b64_size1);
// 				LOG_CL.PRINTF[5] := INT_TO_STRING(idx1);
// 				LOG_CL.PRINTF[6] := INT_TO_STRING(idx2);
// 				LOG_MSG();
// 			END_IF;*)
// 			(*########################################################################*)
// 		ELSE
// 			b64_lf := last_block;
// 			step := 480; (* block finished  *)
// 		END_IF;
// 		IF b64_lf THEN
// 			S_BUF.BUFFER[idx2] := BYTE#16#0D;
// 			idx2 := idx2 + 1;
// 			S_BUF.BUFFER[idx2] := BYTE#16#0A;
// 			idx2 := idx2 + 1;
// 		END_IF;
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'B64: exit ~1 /SIZE ~2 /OFF ~3 /eof ~4 /lf_cnt ~5';
// 			LOG_CL.PRINTF[1] := FSD.FILENAME;
// 			LOG_CL.PRINTF[2] := UDINT_TO_STRING(FSD.FILE_SIZE);
// 			LOG_CL.PRINTF[3] := UDINT_TO_STRING(FSD.OFFSET);
// 			LOG_CL.PRINTF[4] := BOOL_TO_STRING(last_block);
// 			LOG_CL.PRINTF[5] := INT_TO_STRING(snd_lf_cnt);
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
// 	END_IF;
//
// 480:
// 	S_BUF.SIZE := INT_TO_UINT(idx2); (* send data *)
// 	(*####################  DEBUG-MESSAGE  ###################################*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'B64: tcp_send ~1';
// 		LOG_CL.PRINTF[1] := INT_TO_STRING(idx2);
// 		LOG_MSG();
// 	END_IF;*)
// 	(*########################################################################*)
// 	step := 490;
//
// 490:
// 	IF S_BUF.SIZE = UINT#0 THEN (* wait until all data is sent *)
// 		step := SEL(last_block ,460,400);
// 	END_IF;
//
// (* ++++++++++++++++++++++++++++++++++++++++ *)
//
// 800:
// 	IF snd_step = 0 THEN
// 		snd_enable := FALSE;
// 		snd_text := '--x--';
// 		snd_lf_cnt := 2;
// 		step := 810;
// 	END_IF;
//
// 810:
// 	IF snd_step = 0 THEN
// 		snd_enable := TRUE;
// 		snd_text := '.';
// 		snd_lf_cnt := 1;
// 		step := 820;
// 	END_IF;
//
// 820:
// 	IF rcv_state = 250 THEN (* 250 = Requested mail action okay *)
// 		snd_text := 'QUIT';
// 		ip_c_red_disable := TRUE;
// 		next_step := 830;
// 	END_IF;
//
// 830:
// 	IF rcv_state = 221  OR TON_WAIT.ET > timeout3 THEN (* 221 = Service closing transmission channel *)
// 		step := 900;
// 	END_IF;
//
// (* ================================================== *)
//
// 900:
// 	IF FILE_DELETE THEN (* delete files after transfer *)
// 		cnt1 := ELEMENT_COUNT(59, FILES) - 1; (* separator >;<, number of files *)
// 		step := 910;
// 	ELSE
// 		step := 950;
// 	END_IF;
//
// 910:
// 	IF FSD.MODE = BYTE#0 AND FSD.ERROR = BYTE#0 THEN
// 		IF cnt1 >= 0 THEN
// 			str1 := ELEMENT_GET(59, cnt1, FILES); (* separator >;< *)
// 			IF str1 <> '#DEL#' THEN
// 				FSD.FILENAME := str1;
// 				FSD.MODE := BYTE#4; (* File Remove *)
// 			END_IF;
// 			cnt1 := cnt1 - 1;
// 		ELSE
// 			step := 950;
// 		END_IF;
// 	END_IF;
//
// 950:
// 	DONE := TRUE;
// 	step := 980;
//
// 980:
// 	IP_C.C_ENABLE := FALSE; (* close SMTP control channel *)
// 	S_BUF.SIZE := UINT#0;
// 	R_BUF.SIZE := UINT#0;
// 	IP_C.MAILBOX[3] := BYTE#0; (* re-enable IP receive *)
// 	snd_step := 0;
// 	next_step := 0;
// 	FSD.MODE := BYTE#5; (* close file *)
// 	step := 990;
//
// 990:
// 	IF IP_C.C_STATE = BYTE#0 AND FSD.FILE_OPEN = FALSE THEN
// 		step := 0;
// 	END_IF;
//
// (*  +++++++++++++++++ authentication +++++++++++++++++ *)
//
// (* PLAIN ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// 1000:
// 	snd_text := CONCAT('AUTH',str3);
// 	snd_lf_cnt := 1;
//  	step := 1010;
//
// 1010:
// 	IF rcv_state = 334 THEN (* 334 = authentication  *)
// 		cnt1 := LEN(URL_DATA.USER);
// 		cnt2 := LEN(URL_DATA.PASSWORD);
// 		IF cnt1 + cnt2 <= 46 THEN (* USER + PASSWORD together may have at most 46 characters *)
// 			b64_buf1[0] := BYTE#0;
// 			FOR i1 := 1 TO cnt1 DO
// 				b64_buf1[i1] := CODE(URL_DATA.USER, i1);
// 			END_FOR;
// 			b64_size1 := cnt1 + 1;
// 			b64_buf1[b64_size1] := BYTE#0;
// 			b64_size1 := b64_size1 + 1;
// 			FOR i1 := 1 TO cnt2 DO
// 				b64_buf1[b64_size1] := CODE(URL_DATA.PASSWORD, i1);
// 				b64_size1 := b64_size1 + 1;
// 			END_FOR;
// 			BASE64_ENCODE_STREAM( SIZE1:= b64_size1, BUF1:= b64_buf1, BUF2:= b64_buf2, SIZE2=> b64_size2); (* BASE64 encoding *)
// 			idx1 := b64_size2 - 1;
// 			FOR i1 := 0 TO idx1 DO
// 				S_BUF.BUFFER[i1] := b64_buf2[i1];
// 			END_FOR;
// 			S_BUF.BUFFER[b64_size2] := BYTE#16#0D;
// 			b64_size2 := b64_size2 + 1;
// 			S_BUF.BUFFER[b64_size2] := BYTE#16#0A;
// 			b64_size2 := b64_size2 + 1;
// 			S_BUF.SIZE := INT_TO_UINT(b64_size2);
// 			rcv_state := 0;
// 		END_IF;
// 		step := 90;
// 	END_IF;
//
// (* LOGIN ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// 1100:
// 	snd_text := CONCAT('AUTH',str3);
// 	snd_lf_cnt := 1;
//  	step := 1110;
//
// 1110:
// 	IF rcv_state = 334 THEN (* 334 = authentication  *)
// 		b64_str1 := URL_DATA.USER;
// 		b64_start := TRUE; (* BASE64 encoding of username *)
// 		step := 1120;
// 	END_IF;
//
// 1120:
// 	IF b64_start AND b64_done THEN
// 		b64_start := FALSE;
// 		snd_text := b64_str2;
// 		next_step := 1130;
// 	END_IF;
//
// 1130:
// 	IF rcv_state = 334 THEN (* 334 = Password:  *)
// 		b64_str1 := URL_DATA.PASSWORD;
// 		b64_start := TRUE; (* BASE64 encoding of password *)
// 		step := 1140;
// 	END_IF;
//
// 1140:
// 	IF b64_start AND b64_done THEN
// 		b64_start := FALSE;
// 		snd_text := b64_str2;
// 		next_step := 90;
// 	END_IF;
//
// (* CRAM-MD5  ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// (Challenge-Response Authentication Mechanism, Message Digest 5) authentication based on the MD5-HMAC algorithm. *)
//
// 1200:
// 	snd_text := CONCAT('AUTH',str3);
// 	snd_lf_cnt := 1;
//  	step := 1210;
//
// 1210:
// 	IF rcv_state = 334 THEN (* 334 = authentication  *)
// 		md5_b64_ts := RIGHT(rcv_text,LEN(rcv_text) - 4);
// 		md5_user := URL_DATA.USER;
// 		md5_password := URL_DATA.PASSWORD;
// 		MD5_RUN := TRUE;
// 		step := 1220;
// 	END_IF;
//
// 1220:
// 	MD5_CRAM_AUTH(RUN:=MD5_RUN,USERNAME:=md5_user,PASSWORD:=md5_password,B64_TS:=md5_b64_ts,AUTH_KEY:=md5_auth_key);
// 	IF NOT MD5_RUN THEN
// 		snd_text := md5_auth_key;
// 		step := 90;
// 	END_IF;
//
// (* STARTTLS  ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ *)
//
// 1300:
// 	snd_text := 'STARTTLS';
// 	snd_lf_cnt := 1;
//  	step := 1310;
//
// 1310:
// 	IF rcv_state = 220 THEN (* 220 = Service ready - Ready to start TLS *)
// 		IP_C.MAILBOX[16] := BYTE#1; (* activate SSL/TLS *)
// 		step := 1320;
// 	END_IF;
//
// 1320:
// 	IF IP_C.MAILBOX[16] = BYTE#2 THEN (* wait until SSL/TLS has been activated *)
// 		step := 30; (* restart the whole sequence from the beginning *)
// 	END_IF;
//
// (* no authentication found +++++++++++++++++++++++++++++++ *)
// 1900: ; (* wait for runtime error *)
//
// END_CASE;
// (* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ *)
// IF step >= 30 THEN
// 	IF R_BUF.SIZE > UINT#0 THEN (* data received? *)
// 		idx1 := UINT_TO_INT(R_BUF.SIZE);
// 		REPEAT (* detect line termination and determine real end of TEXT *)
// 			idx1 := idx1 - 1;
// 		UNTIL idx1 = 0 OR R_BUF.BUFFER[idx1] >= BYTE#32 (* TRUE = EXIT *)
// 		END_REPEAT;
// 		FOR i1 := 0 TO idx1 DO (* for multiline replace all CRLF with >^< *)
// 			IF R_BUF.BUFFER[i1] < BYTE#32 THEN
// 				R_BUF.BUFFER[i1] := BYTE#94; (* >^< *)
// 			END_IF;
// 		END_FOR;
// 		rcv_text := BUFFER_TO_STRING(PT:=REF(R_BUF.BUFFER),SIZE:=idx1+1,START:=0,STOP:=STRING_LENGTH-1);
// 		IF LEN(rcv_text) >= 3 THEN (* check and evaluate status code *)
// 			rcv_s := LEFT(rcv_text,3);
// 			IF IS_NUM(STR:=rcv_s) THEN (* contains only digits? *)
// 				rcv_state := STRING_TO_INT(rcv_s);
// 			END_IF;
// 			IP_C.R_OBSERVE := FALSE;
// 		END_IF;
// 		R_BUF.SIZE := UINT#0; (* delete data *)
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'RCV: >~1<';
// 			LOG_CL.PRINTF[1] := rcv_text;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
// 	END_IF;
// END_IF;
//
// (* automatic data sending *)
//
// CASE snd_step OF
//
// 00:
// 	IF LEN(snd_text) > 0 THEN
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'SND: >~1< /lf ~2';
// 			LOG_CL.PRINTF[1] := snd_text;
// 			LOG_CL.PRINTF[2] := INT_TO_STRING(snd_lf_cnt);
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
// 		_STRING_TO_BUFFER(str:=snd_text,pos:=snd_idx,pt:=REF(S_BUF.BUFFER),size:=SIZEOF(S_BUF.BUFFER)); (* copy string into buffer *)
// 		snd_idx := snd_idx + LEN(snd_text);
// 		FOR i1 := 1 TO snd_lf_cnt DO
// 			S_BUF.BUFFER[snd_idx] := BYTE#16#0D;
// 			snd_idx := snd_idx + 1;
// 			S_BUF.BUFFER[snd_idx] := BYTE#16#0A;
// 			snd_idx := snd_idx + 1;
// 		END_FOR;
// 		snd_text := '';
// 		rcv_state := 9999; (* clear last SMTP status *)
//
// 		IF snd_enable THEN
// 			S_BUF.SIZE := INT_TO_UINT(snd_idx); (* enter send length *)
// 			R_BUF.SIZE := UINT#0; (* reset receive buffer *)
// 			IP_C.R_OBSERVE := step <= 120 OR step >= 820;
// 			snd_idx := 0;
// 			snd_step := 10;
// 		END_IF;
// 	END_IF;
//
// 10:
// 	IF S_BUF.SIZE = UINT#0 THEN (* wait until all data is sent *)
// 		IF next_step <> 0 THEN
// 			step := next_step;
// 			next_step := 0;
// 		END_IF;
// 		snd_step := 0;
// 	END_IF;
//
// END_CASE;
//
// (* global error handling for main routine *)
// IF ERROR_T = BYTE#0 AND step > 20 THEN (* no error stored yet and main routine active *)
//
// 	IF IP_C.ERROR > DWORD#0 AND IP_C.C_ENABLE AND NOT IP_C.TIME_RESET THEN  (* IP Control 1 *)
// 		IF ip_c_red_disable = FALSE OR IP_C.ERROR <> DWORD#16#FD000000 THEN
// 			ERROR_C := IP_C.ERROR;
// 			ERROR_T := BYTE#2;
// 			step := 980;
// 		END_IF;
// 	END_IF;
//
// 	IF FSD.MODE = BYTE#0 THEN
// 		IF FSD.ERROR > BYTE#0 THEN
// 			ERROR_C := BYTE_TO_DWORD(FSD.ERROR);
// 			ERROR_T := BYTE#4; (* file access *)
// 			step := 980;
// 		END_IF;
// 	END_IF;
// END_IF;
//
// IF TON_WAIT.Q THEN (* sequence fault - step chain *)
// 	ERROR_C := SHL(INT_TO_DWORD(step),16) OR INT_TO_DWORD(rcv_state); (* enter current step and last SMTP status *)
// 	ERROR_T := BYTE#5;
// 	step := 980;
// END_IF;
//
// (* ######################################################## *)
//
// DNS_CLIENT(IP_C:=IP_C, S_BUF:=S_BUF, R_BUF:=R_BUF, DOMAIN:=URL_DATA.DOMAIN, IP4_DNS:= DNS_IP4, ACTIVATE:=step=10); (* DNS-Client *)
// IPC(IP:= DWORD#0, PORT:= WORD#0 , TIME_OUT:= TIMEOUT, IP_C:= IP_C, S_BUF:= S_BUF, R_BUF:= R_BUF); (* SMTP control channel *)
// FS(FSD:=FSD,PT:=F_BUF); (* FILE_SERVER *)
// BASE64_ENCODE_STR(RUN:= b64_start, STR1:= b64_str1, STR2:= b64_str2, DONE=> b64_done); (* BASE64 encoding *)
//
// TON_WAIT(IN:=step = last_step AND step > 0,PT:=timeout2);
// last_step := step;
// activate_last := ACTIVATE;
// BUSY := step <> 0;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK SNTP_CLIENT
VAR_IN_OUT
	IP_C : IP_C; (* IP_Control management structure *)
	S_BUF : NETWORK_BUFFER;
	R_BUF : NETWORK_BUFFER;
END_VAR
VAR_INPUT
	IP4 : DWORD;
	ACTIVATE : BOOL;
END_VAR
VAR_OUTPUT
	ERROR : DWORD;
	DONE_P : BOOL;
	UDT : DT;
	XMS : INT;
END_VAR
VAR
	IP_FIFO : IP_FIFO;
	ip_state : BYTE;
	ip_id : BYTE;
	i : INT;
	last : DWORD;
	state : INT;
	tx : DWORD;
	activate_last : BOOL;
	ip4_stored : DWORD;

END_VAR

(* Domain Name System (query) *)
tx := T_PLC_MS(); (* read system time *)

CASE state OF
00: (* wait for rising edge on SET *)
	DONE_P := FALSE;
	IF ACTIVATE AND NOT activate_last THEN
		ip4_stored := IP4; (* store current IP4 *)
		state := 10;
		ERROR := DWORD#0;
		ip_state := BYTE#1; (* register *)
		(*####################  DEBUG-MESSAGE  ###################################*)
		(*########################################################################*)
	END_IF;

10: (* register with IP_Control and wait for release *)
	IF ip_state = BYTE#3 THEN

		S_BUF.BUFFER[00] := BYTE#16#1B;
		FOR i := 1 TO 47 DO
			S_BUF.BUFFER[i] := BYTE#16#00;
		END_FOR;

		last := tx; (* save current time value *)
		state := 30;

		(* set up IP data traffic *)
		IP_C.C_PORT := WORD#123; (* enter port number *)
		IP_C.C_IP := ip4_stored; (* enter IP *)
		IP_C.C_MODE := BYTE#1; (* Mode: UDP+ACTIVE+PORT+IP *)
		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
		IP_C.R_OBSERVE := TRUE; (* monitor data reception *)
		S_BUF.SIZE := UINT#48; (* enter send length *)
		R_BUF.SIZE := UINT#0; (* reset receive length *)
    		(*####################  DEBUG-MESSAGE  ###################################*)
    		(*########################################################################*)
	END_IF;

30:
	IF IP_C.ERROR <> DWORD#00 THEN
		ERROR := IP_C.ERROR;
		S_BUF.SIZE := UINT#0;
		R_BUF.SIZE := UINT#0;
		state := 0;
    		(*####################  DEBUG-MESSAGE  ###################################*)
    		(*########################################################################*)

	ELSIF S_BUF.SIZE = UINT#0 AND R_BUF.SIZE >= UINT#48 THEN (* check length of received data *)

		(* convert SNTP timestamp base 1900-01-01 -> Unix 1970-01-01 = +2208988800 seconds*)
		UDT := DWORD_TO_DT(DWORD_OF_BYTE(R_BUF.BUFFER[40],R_BUF.BUFFER[41],R_BUF.BUFFER[42],R_BUF.BUFFER[43]) - DWORD#2208988800);

		(* original fractional seconds are represented as a 32 bit fraction *)
		(* calculate ms only with reduced resolution 1/65535 -> use high word *)
		XMS := DWORD_TO_INT(SHR(DINT_TO_DWORD(WORD_TO_DINT(WORD_OF_BYTE(R_BUF.BUFFER[44],R_BUF.BUFFER[45])) * DINT#1000),16));
		(* divide the time difference (runtime error) between send and receive by 2 and add *)
		XMS := XMS + DWORD_TO_INT(SHR(tx - last,1));

		UDT := DWORD_TO_DT(DT_TO_DWORD(UDT) + (XMS / 1000)); (* add the fractional seconds of the time difference to UDT *)
		XMS := XMS MOD 1000;  (* take over ms *)
		state := 0;
		DONE_P := TRUE;
    		(*####################  DEBUG-MESSAGE  ###################################*)
    		(*########################################################################*)

	END_IF;

	IF state = 0 THEN
		IP_C.R_OBSERVE := FALSE; (* monitor data reception *)
		ip_state := BYTE#4;
	END_IF;

END_CASE;

activate_last := ACTIVATE; (* save last status for edge detection *)

(* IP_FIFO *)
IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id);

(***************************************************************************************)
END_FUNCTION_BLOCK

FUNCTION_BLOCK SNTP_SERVER
VAR_IN_OUT
	IP_C : IP_C; (* IP_Control management structure *)
	S_BUF : NETWORK_BUFFER;
	R_BUF : NETWORK_BUFFER;
END_VAR
VAR_INPUT
	ENABLE : BOOL;
	STRATUM : BYTE := BYTE#1;
	UDT : DT;
	XMS : INT;
END_VAR
VAR
	IP_FIFO : IP_FIFO;
	ip_state : BYTE;
	ip_id : BYTE;
	i2 : INT;
	i : INT;
	timestamp_int : DWORD; (* Timestamp integer part *)
	timestamp_sek : DWORD; (* Timestamp fraction part (HIGH_WORD) *)
	ts_b : BYTE;
	state : INT;
	tx : UDINT;
	t : TON;

END_VAR

(* Domain Name System (query) *)
tx := T_PLC_MS(); (* read system time *)

CASE state OF

00: (* wait for release *)
	IF ENABLE THEN
		state := 10;
		ip_state := BYTE#1; (* register *)
    		(*####################  DEBUG-MESSAGE  ###################################*)
    		(*########################################################################*)
	END_IF;

10: (* register with IP_Control and wait for release *)
	IF ip_state = BYTE#3 THEN
		(* set up IP data traffic *)
		IP_C.C_PORT := WORD#123; (* enter port number *)
		IP_C.C_IP := DWORD#00; (* enter IP *)
		IP_C.C_MODE := BYTE#5; (* Mode: UDP+PASSIV+PORT *)
		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
		IP_C.R_OBSERVE := FALSE; (* monitor data reception *)
		S_BUF.SIZE := UINT#0; (* enter send length *)
		R_BUF.SIZE := UINT#0; (* reset receive length *)
		state := 20;
    		(*####################  DEBUG-MESSAGE  ###################################*)
    		(*########################################################################*)
	END_IF;
20:
	IF R_BUF.SIZE > UINT#0 THEN
		state := 20;
    		(*####################  DEBUG-MESSAGE  ###################################*)
    		(*########################################################################*)
		(* pad received telegram with 00 if too short! and preset send buffer with 00 *)
		FOR i := 0 TO 47 DO
			IF INT_TO_UINT(i) > R_BUF.SIZE - 1 THEN
				R_BUF.BUFFER[i] := BYTE#00;
			END_IF;
			S_BUF.BUFFER[i] := BYTE#00;
		END_FOR;

		R_BUF.SIZE := UINT#0;

		(* prepare send data *)
		S_BUF.BUFFER[0] := BYTE#2#11100; (* Flags: NO WARNING, NTP VERSION 3 , SERVER *)
		S_BUF.BUFFER[1] := STRATUM; (* Peer Clock Stratum: *)
		S_BUF.BUFFER[2] := BYTE#10; (* peer polling interval 2 to the power of 10 = 1024 sec. *)
		S_BUF.BUFFER[3] := BYTE#16#FB; (* Peer Clock Precision: -5 = -31,25ms *)

		(* convert SNTP timestamp base 1900-01-01 = Unix 1970-01-01 + 2208988800 seconds *)
		timestamp_int := DT_TO_DWORD(UDT) + DWORD#2208988800;

		(* fractional seconds are represented as a 32 bit fraction *)
		(* calculate only with reduced resolution 1/65535 -> use high word *)
		(* calculation of the fractional seconds: xms * 65535 / 1000 *)
		timestamp_sek := DINT_TO_DWORD(DWORD_TO_DINT(SHL(INT_TO_DWORD(XMS),16))/DINT#1000);

		ts_b := BYTE_OF_DWORD(timestamp_int,3); (* Timestamp Byte 3 *)
		S_BUF.BUFFER[16] := ts_b; (* Reference Timestamp *)
		S_BUF.BUFFER[32] := ts_b; (* Receive   Timestamp *)
		S_BUF.BUFFER[40] := ts_b; (* Transmit  Timestamp *)

		ts_b := BYTE_OF_DWORD(TIMESTAMP_INT,2); (* Timestamp Byte 2 *)
		S_BUF.BUFFER[17] := ts_b; (* Reference Timestamp *)
		S_BUF.BUFFER[33] := ts_b; (* Receive   Timestamp *)
		S_BUF.BUFFER[41] := ts_b; (* Transmit  Timestamp *)

		ts_b := BYTE_OF_DWORD(TIMESTAMP_INT,1); (* Timestamp Byte 1 *)
		S_BUF.BUFFER[18] := ts_b; (* Reference Timestamp *)
		S_BUF.BUFFER[34] := ts_b; (* Receive   Timestamp *)
		S_BUF.BUFFER[42] := ts_b; (* Transmit  Timestamp *)

		ts_b := BYTE_OF_DWORD(TIMESTAMP_INT,0); (* Timestamp Byte 0 *)
		S_BUF.BUFFER[19] := ts_b; (* Reference Timestamp *)
		S_BUF.BUFFER[35] := ts_b; (* Receive   Timestamp *)
		S_BUF.BUFFER[43] := ts_b; (* Transmit  Timestamp *)

		ts_b := BYTE_OF_DWORD(TIMESTAMP_SEK,1); (* Timestamp Byte 1 *)
		S_BUF.BUFFER[20] := ts_b; (* Reference Timestamp *)
		S_BUF.BUFFER[36] := ts_b; (* Receive   Timestamp *)
		S_BUF.BUFFER[44] := ts_b; (* Transmit  Timestamp *)

		ts_b := BYTE_OF_DWORD(timestamp_sek,0);	(* Timestamp Byte 0 *)
		S_BUF.BUFFER[21] := ts_b; (* Reference Timestamp *)
		S_BUF.BUFFER[37] := ts_b; (* Receive   Timestamp *)
		S_BUF.BUFFER[45] := ts_b; (* Transmit  Timestamp *)

		(* copy bytes 1-8 of transmit time to originate time *)
		i2 := 40;
		FOR i := 24 TO 31 DO
			S_BUF.BUFFER[i] := R_BUF.BUFFER[i2];
			i2 := i2 +1;
		END_FOR;

		S_BUF.SIZE := UINT#48; (* send data *)

	ELSIF NOT ENABLE AND S_BUF.SIZE = UINT#0 THEN
		ip_state := BYTE#4; (* unregister *)
		state := 0;
	END_IF;

END_CASE;

t(IN:= IP_C.ERROR > DWORD#0, PT:=T#5s);
IF t.Q THEN
  IP_C.TIME_RESET := TRUE; (* Error Reset *)
END_IF;

(* IP_FIFO *)
IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id);
END_FUNCTION_BLOCK

// Commented out for beedance: it uses BUFFER_TO_STRING, which is commented out.
// FUNCTION_BLOCK SPIDER_ACCESS
// VAR_IN_OUT
// 	IP_C : IP_C;(*IP_Control management structure*)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// 	VALUE : STRING;
// 	VAR_NAME : STRING(40);
// END_VAR
// VAR_INPUT
// 	MODE : BYTE;
// END_VAR
// VAR_OUTPUT
// 	ERROR : DWORD;
// END_VAR
// VAR
// 	state : INT;
// 	st_tmp : STRING(STRING_LENGTH);
// 	URL_DATA : url;
// 	HTTP_GET : HTTP_GET;
// 	value_len :	UINT;
// 	body_len :	UINT;
// 	mode_save : BYTE;
// END_VAR
//
// CASE state OF
//
// 00:
// 	IF MODE > BYTE#0 AND MODE < BYTE#3 THEN
// 		MODE_SAVE := MODE; (* store mode *)
// 		(* if there is no '.' in the name, it is a global variable *)
// 		IF FIND(VAR_NAME,'.') = 0 THEN
// 			st_tmp := CONCAT('%40GV.',VAR_NAME);
// 		ELSE
// 			st_tmp := VAR_NAME;
// 		END_IF;
//
// 		IF MODE_SAVE = BYTE#2 THEN (* 2 = Write *)
// 			VALUE_LEN := INT_TO_UINT(LEN(VALUE));
// 			st_tmp := CONCAT(st_tmp,'+');
// 			st_tmp := CONCAT(st_tmp,VALUE);
// 			URL_DATA:=STRING_TO_URL(STR:= CONCAT('http://x/dir/cgi-bin/writeVal.exe?',st_tmp),DEFAULT_PROTOCOL:='',DEFAULT_PATH:='');
// 		ELSE
// 			URL_DATA:=STRING_TO_URL(STR:= CONCAT('http://x/dir/cgi-bin/readVal.exe?',st_tmp),DEFAULT_PROTOCOL:='',DEFAULT_PATH:='');
// 		END_IF;
//
// 		state := 60;
// 	END_IF;
// 60:
// 	IF HTTP_GET.DONE THEN
// 		IF MODE_SAVE = BYTE#2 THEN (* Write *)
// 			BODY_LEN := HTTP_GET.BODY_STOP - HTTP_GET.BODY_START + BOOL_TO_UINT(HTTP_GET.BODY_STOP > UINT#0);
// 			ERROR := BOOL_TO_DWORD(BODY_LEN <> VALUE_LEN);
// 		ELSE (* Read *)
// 			IF HTTP_GET.BODY_START > 0 THEN
// 				VALUE := BUFFER_TO_STRING(PT:=REF(R_BUF.BUFFER),SIZE:=R_BUF.SIZE,START:=UINT_TO_INT(HTTP_GET.BODY_START),STOP:=UINT_TO_INT(HTTP_GET.BODY_STOP));
// 			ELSE
// 				VALUE := '';
// 			END_IF;
// 		END_IF;
// 		STATE   := 100;
//
// 	ELSIF (HTTP_GET.ERROR > DWORD#00) THEN
// 		ERROR := HTTP_GET.ERROR;
// 		STATE   := 100;
// 	END_IF;
//
// 100:
// 	(* UNLOCK HTTP DATA *)
// 	IF HTTP_GET.DONE = FALSE THEN
// 		STATE := 0;
// 	END_IF;
//
// END_CASE;
//
// (* HTTP_GET *)
// HTTP_GET(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			IP4:=DWORD#00,
// 			GET:=state=60,
// 			MODE:=BYTE#0,
// 			UNLOCK_BUF:=state=100,
// 			URL_DATA:=URL_DATA
// 			);
// END_FUNCTION_BLOCK

FUNCTION STRING_TO_URL : URL
VAR_INPUT
	STR : STRING(STRING_LENGTH);
	DEFAULT_PROTOCOL : STRING(10);
	DEFAULT_PATH : STRING(80);
END_VAR
VAR
	pos: INT;
	pos2: INT;
	x : STRING(STRING_LENGTH);
END_VAR

(* check if string is empty and return empty url *)
IF STR = '' THEN RETURN; END_IF;

(* extract the protocol *)
(* first we check if the optional protocol is specified ending with a :// *)
pos := FIND(STR, '://');
IF pos > 0 THEN
	(* store protocol and restore the rest of url in temp *)
	STRING_TO_URL.PROTOCOL := LEFT(STR, pos - 1);
	x := RIGHT(STR, LEN(STR) - pos - 2);
ELSE
	x := STR;
	STRING_TO_URL.PROTOCOL := DEFAULT_PROTOCOL;
END_IF;

(* extract user and password if present ending with a @ *)
pos := FINDB(STR1:=x,STR2:='@'); (* Find last @ *)

IF pos > 0 THEN
	(* extract user *)
	pos2 := FIND(x,':')+1;
	IF pos2 > 1 AND pos2 <= pos THEN
		STRING_TO_URL.USER := LEFT(x,pos2 - 2);
		STRING_TO_URL.PASSWORD := MID(x, pos - pos2, pos2);
	ELSE
		STRING_TO_URL.USER := LEFT(x, pos - 2);
		STRING_TO_URL.PASSWORD := '';
	END_IF;
	x := RIGHT(x, LEN(x) - pos);
ELSE
	STRING_TO_URL.USER := '';
	STRING_TO_URL.PASSWORD := '';
END_IF;

(* extract anchor *)
pos := FIND(x, '#');
IF pos > 0 THEN
	STRING_TO_URL.ANCHOR := RIGHT(x, LEN(x) - pos);
	x := LEFT(x, pos - 1);
ELSE
	STRING_TO_URL.ANCHOR := '';
END_IF;

(* extract query  *)
pos := FIND(x, '?');
IF pos > 0 THEN
	STRING_TO_URL.QUERY := RIGHT(x, LEN(x) - pos);
	x := LEFT(x, pos - 1);
ELSE
	STRING_TO_URL.QUERY := '';
END_IF;

(* extract path  *)
pos := FIND(x, '/');
IF pos > 0 THEN
	STRING_TO_URL.PATH := RIGHT(x, LEN(x) - pos + 1);
	x := LEFT(x, pos - 1);
ELSE
	STRING_TO_URL.PATH := DEFAULT_PATH;
END_IF;

(* extract port and leave domain *)
pos := FIND(x, ':');
IF pos > 0 THEN
	STRING_TO_URL.PORT := STRING_TO_WORD(RIGHT(x, LEN(x) - pos));
	STRING_TO_URL.DOMAIN := LEFT(x, pos - 1);
ELSE
	STRING_TO_URL.PORT := 0;
	STRING_TO_URL.DOMAIN := x;
END_IF;
END_FUNCTION

// Commented out for beedance: it uses DT_TO_STRF, which is commented out.
// FUNCTION_BLOCK SYS_LOG
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// END_VAR
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	LDT : DT;
// 	SERVER_IP4 : DWORD;
// 	PORT : WORD;
// 	FACILITY : BYTE;
// 	SEVERITY : BYTE;
// 	TAG : STRING(32);
// 	HOSTNAME : STRING;
// 	MESSAGE : STRING(STRING_LENGTH);
// 	OPTION : BYTE;
// END_VAR
// VAR_OUTPUT
// 	DONE_P : BOOL;
// 	ERROR : DWORD;
// END_VAR
// VAR
// 	ip_state : BYTE;
// 	ip_id : BYTE;
// 	IP_FIFO : IP_FIFO;
// 	activate_last : BOOL;
// 	state : INT;
// 	s1 : STRING(20);
// 	i : INT;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;*)
//
// END_VAR
//
// CASE state OF
// 00: (* wait for rising edge on SET *)
// 	IF ACTIVATE AND NOT ACTIVATE_LAST THEN
// 		state := 10;
// 		ERROR := DWORD#0;
// 		IP_STATE := BYTE#1; (* register *)
//     		(*####################  DEBUG-MESSAGE  ###################################*)
//     		(*IF _debug_enable THEN
//       			LOG_CL.NEW_MSG := 'SYS_LOG: S00 START /IP_ID ~4';
//       			LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
//       			LOG_MSG();
//     		END_IF;*)
//     		(*########################################################################*)
// 	END_IF;
// 	DONE_P := FALSE;
//
// 10: (* register with IP_Control and wait for release *)
// 	IF IP_STATE = BYTE#3 THEN
// 		(* prepare send data *)
// 		i := 0;
// 		IF (OPTION AND BYTE#0000_0001) = BYTE#00 THEN (* Option: no priority Code *)
// 			(* Priority:  '<xx>' as decimal value *)
// 			(* BYTE: xxxx x...  = Facility *)
// 			(*       .... .xxx  = Level *)
// 			s1 := REPLACE('<P>',BYTE_TO_STRING(SHL(FACILITY,3) OR (SEVERITY AND BYTE#0000_0111)),1,2);
// 		ELSE
// 			s1 := '';
// 		END_IF;
//
// 		IF (OPTION AND BYTE#0000_0010) = BYTE#00 THEN (* Option: No Use Syslog RFC Header *)
// 			(* 'month day hh:mm:ss ' *)
// 			s1 := CONCAT(s1,DT_TO_STRF(DTI:=LDT,MS:=0,FMT:='#E #W #N:#R:#T ',LANG:=1));
// 			_STRING_TO_BUFFER(STR:=s1,POS:=i,PT:=REF(S_BUF.BUFFER),SIZE:=UINT#1024);
// 			i := i + LEN(s1);
// 			(* ' enter sourcename *)
// 			_STRING_TO_BUFFER(STR:=HOSTNAME,POS:=i,PT:=REF(S_BUF.BUFFER),SIZE:=UINT#1024);
// 			i := i + LEN(HOSTNAME);
// 			S_BUF.BUFFER[i] := BYTE#32; (* insert space *)
// 			i := i + 1;
// 			(* ' enter tag *)
// 			_STRING_TO_BUFFER(STR:=TAG,POS:=i,PT:=REF(S_BUF.BUFFER),SIZE:=UINT#1024);
// 			i := i + LEN(TAG);
// 			S_BUF.BUFFER[i] := BYTE#32; (* insert space *)
// 			i := i + 1;
// 		END_IF;
//
// 		(* ' enter message *)
// 		_STRING_TO_BUFFER(STR:=MESSAGE,POS:=i,PT:=REF(S_BUF.BUFFER),SIZE:=UINT#1024);
// 		i := i + LEN(MESSAGE);
//
// 		IF (OPTION AND BYTE#0000_0100) > BYTE#00 THEN (* Option: Add a CR & LF to end of message *)
// 			S_BUF.BUFFER[i] := BYTE#13; (* insert carriage return *)
// 			i := i + 1;
// 			S_BUF.BUFFER[i] := BYTE#10; (* insert line feed *)
// 			i := i + 1;
// 		END_IF;
//
// 		state := 30;
//
// 		(* set up IP data traffic *)
// 		IF (OPTION AND BYTE#2#0000_1000) > BYTE#00 THEN (* Option: TCP Mode *)
// 			IP_C.C_MODE := BYTE#0; (* Mode: TCP/ACTIVE *)
// 		ELSE
// 			IP_C.C_MODE := BYTE#1; (* Mode: UDP/ACTIVE *)
// 		END_IF;
//
// 		IF PORT = WORD#0 THEN (* no port number specified? *)
// 			IF IP_C.C_MODE = BYTE#0 THEN (* determine mode  TCP or UDP? *)
// 				IP_C.C_PORT := WORD#1468; (* default port number for syslog over TCP *)
// 			ELSE
// 				IP_C.C_PORT := WORD#514; (* default port number for syslog over UDP *)
// 			END_IF;
// 		ELSE
// 			IP_C.C_PORT := PORT; (* enter the given port number *)
// 		END_IF;
//
// 		IP_C.C_IP := SERVER_IP4; (* enter IP *)
// 		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
// 		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
// 		IP_C.R_OBSERVE := FALSE; (* monitor data reception *)
// 		S_BUF.SIZE := INT_TO_UINT(i); (* enter send length *)
// 		R_BUF.SIZE := UINT#0; (* reset receive length *)
//     		(*####################  DEBUG-MESSAGE  ###################################*)
//     		(*IF _debug_enable THEN
//       			LOG_CL.NEW_MSG := 'SYS_LOG: S20 SEND /IP_ID ~4';
//       			LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
//       			LOG_MSG();
//     		END_IF;*)
//     		(*########################################################################*)
// 	END_IF;
//
// 30:
// 	IF IP_C.ERROR <> DWORD#00 THEN
// 		ERROR := IP_C.ERROR;
// 		state := 0;
//     		(*####################  DEBUG-MESSAGE  ###################################*)
//     		(*IF _debug_enable THEN
//       			LOG_CL.NEW_MSG := 'SYS_LOG: S30 ERROR ~1 IP_ID ~4';
//       			LOG_CL.PRINTF[1]  := DWORD_TO_STRING(ERROR);
//       			LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
//       			LOG_MSG();
//     		END_IF;*)
//     		(*########################################################################*)
// 	ELSIF S_BUF.SIZE = UINT#0 THEN (* wait until all data is sent *)
// 		DONE_P := TRUE;
// 		state := 0;
//     		(*####################  DEBUG-MESSAGE  ###################################*)
//     		(*IF _debug_enable THEN
//       			LOG_CL.NEW_MSG := 'SYS_LOG: S30 OK /IP_ID ~4';
//       			LOG_CL.PRINTF[4]  := BYTE_TO_STRING(IP_ID);
//       			LOG_MSG();
//     		END_IF;*)
//     		(*########################################################################*)
// 	END_IF;
//
// 	IF state = 0 THEN (* release connection *)
// 		IP_STATE := BYTE#4;
// 	END_IF;
//
// END_CASE;
//
// ACTIVATE_LAST := ACTIVATE; (* save last status for edge detection *)
//
// (* IP_FIFO *)
// IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=IP_ID);
// END_FUNCTION_BLOCK

FUNCTION_BLOCK TELNET_LOG
VAR_IN_OUT
	IP_C : IP_C;
	S_BUF : NETWORK_BUFFER;
	LOG_CL : LOG_CONTROL;
END_VAR
VAR_INPUT
	ENABLE : BOOL;
	OPTION : BYTE := BYTE#2#1000_1100; (* no flush buffer (7) , CR&LF (3), Colors (2), AUTOWRAP (1), Screen-Init (0) *)
	PORT : WORD := 23; (* port number *)
END_VAR
VAR_OUTPUT
	READY : BOOL;
END_VAR
VAR
	TELNET_PRINT : TELNET_PRINT;
	done : BOOL;
	init : BOOL;
	watchdog : TON;
	ci : INT;
	pi : INT;
	pi_last : INT;
	send : BOOL;
END_VAR

IF NOT init THEN
	init := TRUE;
	IF OPTION > 0 THEN
		TELNET_PRINT.OPTION := OPTION;
	END_IF;
	IF PORT > 0 THEN (* a specified port overrides the default port *)
		TELNET_PRINT.PORT := PORT;
	END_IF;
  	WATCHDOG.PT := T#1ms;
END_IF;

WATCHDOG(IN:=FALSE); (* stop timer *)

ci := LOG_CL.IDX;

(* evaluate CONNECT_OFF edge *)
IF IP_C.C_STATE = BYTE#1 THEN
	pi := 0;
	init := FALSE;
END_IF;

REPEAT
	IF pi <> ci AND IP_C.C_STATE > BYTE#127 THEN
		IF pi = 0 AND LOG_CL.RING_MODE THEN pi := ci; END_IF;
		pi_last := pi;
		pi := pi + 1;
		IF pi > LOG_CL.SIZE THEN pi := 1; END_IF;
		SEND := TRUE;
	ELSE
		SEND := FALSE;
	END_IF;

	TELNET_PRINT(	IP_C := IP_C,
					S_BUF := S_BUF,
					ENABLE := ENABLE,
					SEND := SEND,
					TEXT := LOG_CL.MSG[pi],
					BACK_COLOR := DWORD_TO_BYTE(SHL(LOG_CL.MSG_OPTION[pi],8)),
					FRONT_COLOR := DWORD_TO_BYTE(LOG_CL.MSG_OPTION[pi] ));
	READY := TELNET_PRINT.READY;
	done := TELNET_PRINT.DONE;

	IF SEND AND NOT done THEN
		pi := pi_last;
	END_IF;

	WATCHDOG(IN:=TRUE);

UNTIL (WATCHDOG.Q OR NOT done )
END_REPEAT;
END_FUNCTION_BLOCK

// Commented out for beedance: it uses _STRING_TO_BUFFER, which is commented out.
// FUNCTION_BLOCK TELNET_PRINT
// VAR_IN_OUT
// 	IP_C : IP_C;
// 	S_BUF : NETWORK_BUFFER;
// END_VAR
// VAR_INPUT
// 	TEXT : STRING(STRING_LENGTH);
// 	ENABLE : BOOL;
// 	SEND : BOOL;
// 	OPTION : BYTE := BYTE#2#1000_1100; (* no flush buffer (7) , CR&LF (3), Colors (2), AUTOWRAP (1), Screen-Init (0) *)
// 	BACK_COLOR : BYTE;
// 	FRONT_COLOR : BYTE;
// 	X_POS : BYTE;
// 	Y_POS : BYTE;
// 	PORT : WORD := 23; (* port number *)
// END_VAR
// VAR_OUTPUT
// 	READY : BOOL;
// 	DONE : BOOL;
// END_VAR
// VAR
// 	last_fc : BYTE;
// 	state : INT;
// 	x : INT := -1;
// 	ip_state : BYTE;
// 	ip_id : BYTE;
// 	IP_FIFO : IP_FIFO;
// 	init : BOOL;
// 	b0_screen_clear : BOOL;
// 	b1_autowrap : BOOL;
// 	b2_color : BOOL;
// 	b3_cr_lr : BOOL;
// 	b7_no_flush : BOOL;
// END_VAR
//
// READY := IP_C.C_STATE > BYTE#127;
// DONE := FALSE;
//
// CASE state OF
//
// 00: (* wait for release *)
// 	IF ENABLE THEN
// 		state := 10;
// 		ip_state := BYTE#1; (* register *)
// 	END_IF;
//
// 10: (* register with IP_Control and wait for release *)
// 	IF ip_state = BYTE#3 THEN
//
// 		(* set up IP data traffic *)
// 		IP_C.C_PORT := PORT; (* enter port number *)
// 		IP_C.C_IP := DWORD#00; (* enter IP *)
// 		IP_C.C_MODE := BYTE#4; (* Mode: TCP+PASSIV+PORT *)
// 		IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
// 		IP_C.C_ENABLE := TRUE; (* enable connection setup *)
// 		IP_C.R_OBSERVE := FALSE; (* no receive monitoring *)
// 		state := 20;
// 	END_IF;
//
// 20:	IF IP_C.ERROR > DWORD#0 THEN
// 		IP_C.TIME_RESET := TRUE; (* reset error *)
// 	END_IF;
// 	IF S_BUF.SIZE = UINT#0 THEN (* data sending not active *)
// 		IF NOT ENABLE THEN
// 			ip_state := BYTE#4; (* close connection - release resource *)
// 			IP_C.C_ENABLE := FALSE; (* connection shutdown *)
// 			state := 0;
// 		END_IF;
// 		IF IP_C.C_STATE = BYTE#1 THEN (* rising edge connect off *)
// 			init := FALSE;
// 		ELSIF READY THEN (* connection present *)
// 			IF IP_C.C_STATE = BYTE#254 AND NOT INIT THEN (* rising edge connect on *)
// 				(* var init is necessary so the block can be called several times in the same cycle *)
// 				init := TRUE;
// 				b0_screen_clear := (OPTION AND BYTE#2#0000_0001) > BYTE#0;
// 				b1_autowrap := (OPTION AND BYTE#2#0000_0010) > BYTE#0;
// 				b2_color := (OPTION AND BYTE#2#0000_0100) > BYTE#0;
// 				b3_cr_lr := (OPTION AND BYTE#2#0000_1000) > BYTE#0;
// 				b7_no_flush := (OPTION AND BYTE#2#1000_0000) > BYTE#0;
//
// 				last_fc := NOT FRONT_COLOR; (* force color update!*)
// 				x := 0;
//
// 				(* AUTOWRAP OFF/ON *)
// 				S_BUF.BUFFER[x] := BYTE#16#1B; (* ESC *)
// 				x := x +1;
// 				S_BUF.BUFFER[x] := BYTE#16#5B; (* [ *)
// 				x := x +1;
// 				S_BUF.BUFFER[x] := BYTE#16#3F; (* ? *)
// 				x := x +1;
// 				S_BUF.BUFFER[x] := BYTE#16#37; (* 7 *)
// 				x := x +1;
// 				S_BUF.BUFFER[x] := SEL(b1_autowrap, BYTE#16#6C (* l *),BYTE#16#68 (* h *));
//
// 				IF b0_screen_clear THEN
// 					IF b2_color THEN
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#1B; (* ESC *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#5B; (* [ *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#30; (* 0 *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#3B; (* ; *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#33; (* 3 *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := (FRONT_COLOR AND BYTE#2#00000111) OR BYTE#16#30; (* output ones digit of foreground color *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#3B; (* ; *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#34; (* 4 *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := (BACK_COLOR AND BYTE#2#00000111) OR BYTE#16#30; (* output ones digit of background color *)
// 						x := x +1;
// 						S_BUF.BUFFER[x] := BYTE#16#6D; (* m *)
// 					END_IF;
// 					(* Screen Clear *)
// 					x := x +1;
// 					S_BUF.BUFFER[x] := BYTE#16#1B; (* ESC *)
// 					x := x +1;
// 					S_BUF.BUFFER[x] := BYTE#16#5B; (* [ *)
// 					x := x +1;
// 					S_BUF.BUFFER[x] := BYTE#16#32; (* 2 *)
// 					x := x +1;
// 					S_BUF.BUFFER[x] := BYTE#16#4A; (* J *)
// 				ELSE
// 					x := x +1;
// 					S_BUF.BUFFER[x] := BYTE#16#0A; (* new line *)
// 				END_IF;
// 			END_IF;
//
// 			(* buffer full, or data present and no further send request *)
// 			IF (INT_TO_UINT(x + LEN(TEXT) + 21) > SIZEOF(S_BUF.BUFFER)) OR (x >= 0 AND NOT SEND) THEN
// 				S_BUF.SIZE := INT_TO_UINT(x + 1); (* enter length -> start sending data *)
// 				x := -1;
// 			ELSIF SEND THEN
//
// 				(* determine current color  *)
// 				IF b2_color AND FRONT_COLOR <> last_FC THEN
// 					(* color setting: output ANSI sequence ESC[X;Y;Zm *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#1B; (* ESC *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#5B; (* [ *)
// 					(* color setting: determine dark or bright and blinking *)
// 					IF ((FRONT_COLOR AND BYTE#2#0001_1000) XOR (last_FC AND BYTE#2#0001_1000)) <> BYTE#00 OR TRUE THEN
// 						x := x + 1;
// 						S_BUF.BUFFER[x] := BYTE#16#30; (* 0  *)
// 						x := x + 1;
// 						S_BUF.BUFFER[x] := BYTE#16#3B; (* ; *)
// 						(* no dark colors wanted -> then bold & bright *)
// 						IF (FRONT_COLOR AND BYTE#2#00001000) = BYTE#0 THEN
// 							(* switch to bold & bright colors *)
// 							x := x + 1;
// 							S_BUF.BUFFER[x] := BYTE#16#31; (* 1 *)
// 							x := x + 1;
// 							S_BUF.BUFFER[x] := BYTE#16#3B; (* ; *)
// 						END_IF;
// 						(* blinking colors wanted *)
// 						IF (FRONT_COLOR AND BYTE#2#00010000) > BYTE#0 THEN
// 							(* switch to blinking color *)
// 							x := x + 1;
// 							S_BUF.BUFFER[x] := BYTE#16#35; (* 5 *)
// 							x := x + 1;
// 							S_BUF.BUFFER[x] := BYTE#16#3B; (* ; *)
// 						END_IF;
// 					END_IF;
// 					(* output tens digit of foreground color *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#33; (* 3 *)
// 					(* output ones digit of foreground color *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := (FRONT_COLOR AND BYTE#2#00000111) OR BYTE#16#30;
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#6D; (* m *)
// 					(* store new color as last used *)
// 					last_FC := FRONT_COLOR;
// 				END_IF;
//
// 				(* position cursor *)
// 				IF X_POS > BYTE#0 AND Y_POS > BYTE#0 THEN
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#1B; (* ESC *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#5B; (* [ *)
// 					(* output tens digit of Y_POS *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := SHR(INT_TO_BCDC(BYTE_TO_INT(Y_POS)),4) OR BYTE#16#30;
// 					(* output ones digit of Y_POS *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := (INT_TO_BCDC(BYTE_TO_INT(Y_POS)) AND BYTE#16#0F) OR BYTE#16#30;
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#3B; (* ; *)
// 					(* output tens digit of X_POS *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := SHR(INT_TO_BCDC(BYTE_TO_INT(X_POS)),4) OR BYTE#16#30;
// 					(* output ones digit of X_POS *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := (INT_TO_BCDC(BYTE_TO_INT(X_POS)) AND BYTE#16#0F) OR BYTE#16#30;
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#48; (* H *)
// 				END_IF;
//
// 				(* output text *)
// 				IF LEN(TEXT) > 0 THEN
// 					_STRING_TO_BUFFER(str:=TEXT,pos:=x + 1,pt :=REF(S_BUF.BUFFER), size:=SIZEOF(S_BUF.BUFFER));
// 					x := x + LEN(TEXT);
// 				END_IF;
//
// 				IF b3_cr_lr THEN (* Option: output CR,LF *)
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#0D;
// 					x := x + 1;
// 					S_BUF.BUFFER[x] := BYTE#16#0A;
// 				END_IF;
//
// 				DONE := TRUE; (* finished entering into buffer *)
//
// 				IF b7_no_flush THEN (* Option: no buffer flush *)
// 					(* send buffer immediately - store current write index -> trigger send *)
// 					S_BUF.SIZE := INT_TO_UINT(x + 1);
// 					x := -1;
// 				END_IF;
// 			END_IF;
// 		END_IF;
// 	END_IF;
// END_CASE;
//
// (* IP_FIFO *)
// IP_FIFO(FIFO:=IP_C.FIFO,STATE:=ip_state,ID:=ip_id);
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses TN_INPUT_MENU_BAR, which is commented out.
// FUNCTION_BLOCK TN_FRAMEWORK
// VAR
// 	FB_IP_CONTROL : IP_CONTROL;
// 	FB_TN_INPUT_CONTROL : TN_INPUT_CONTROL;
// 	FB_TN_INPUT_MENU_BAR : TN_INPUT_MENU_BAR;
// 	FB_TN_RECEIVE :	TN_RECEIVE;
// 	FB_TN_SEND_ROWS : TN_SEND_ROWS;
// 	t : TON;
// END_VAR
// VAR_IN_OUT
// 	us_TN_INPUT_CONTROL : us_TN_INPUT_CONTROL;
// 	us_TN_SCREEN : us_TN_SCREEN;
// 	us_TN_MENU : us_TN_MENU;
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// 	IP_C : IP_C;
// END_VAR
// VAR_INPUT
// 	PORT : WORD;
// END_VAR
//
// (* initialization / release of the connection *)
// IF (NOT IP_C.C_ENABLE) THEN
//
// 	(* set up IP data traffic *)
// 	IP_C.C_PORT := SEL(PORT=BYTE#0,PORT,WORD#23); (* enter port number *)
// 	IP_C.C_IP := DWORD#00; (* enter IP *)
// 	IP_C.C_MODE := BYTE#4; (* Mode: TCP+PASSIV+PORT *)
// 	IP_C.TIME_RESET := TRUE; (* reset timeout monitoring *)
// 	IP_C.C_ENABLE := TRUE; (* enable connection setup *)
// 	IP_C.R_OBSERVE := FALSE; (* no receive monitoring *)
//
// END_IF;
//
// t(IN:= IP_C.ERROR > DWORD#0, PT:=T#5s);
// IF t.Q THEN
// 	IP_C.TIME_RESET := TRUE; (* Error Reset *)
// END_IF;
//
// (* Telnet-Menu-Bar *)
// FB_TN_INPUT_MENU_BAR(
// 	Xus_TN_MENU := us_TN_MENU,
// 	Xus_TN_SCREEN := us_TN_SCREEN);
//
// (* Telnet-Input-Control *)
// FB_TN_INPUT_CONTROL(
// 	Xus_TN_SCREEN := us_TN_SCREEN,
// 	Xus_TN_INPUT_CONTROL := us_TN_INPUT_CONTROL);
//
// (* Telnet-Data-Send Line-Updates *)
// FB_TN_SEND_ROWS(
// 	IP_C :=IP_C,
// 	S_BUF :=S_BUF,
// 	Xus_TN_SCREEN :=us_TN_SCREEN,
// 	S_BUF_SIZE :=SIZEOF(S_BUF.BUFFER));
//
// (* Telnet-Server *)
// FB_IP_CONTROL(
// 	IP_C :=IP_C,
// 	S_BUF :=S_BUF,
// 	R_BUF :=R_BUF,
// 	IP :=DWORD#00,
// 	PORT :=WORD#00,
// 	TIME_OUT :=T#2s);
//
// (* Telnet-Data-Receive *)
// FB_TN_RECEIVE(
// 	R_BUF := R_BUF,
// 	Xus_TN_SCREEN := us_TN_SCREEN);
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses TN_INPUT_SELECT_POPUP, which is commented out.
// FUNCTION_BLOCK TN_INPUT_CONTROL
// VAR_IN_OUT
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// 	Xus_TN_INPUT_CONTROL : us_TN_INPUT_CONTROL;
// END_VAR
// VAR
// 	FB_TN_INPUT_SELECT_POPUP : TN_INPUT_SELECT_POPUP;
// 	FB_TN_INPUT_SELECT_TEXT : TN_INPUT_SELECT_TEXT;
// 	FB_TN_SC_WRITE_C : TN_SC_WRITE_C;
// 	FB_TN_INPUT_EDIT_LINE : TN_INPUT_EDIT_LINE;
// END_VAR
// VAR
// 	in_index : INT;
// 	in_Scroll_Offset : INT;
// 	in_Focus_new : INT;
// 	in_Focus_alt : INT;
// 	bo_Tool_Tip_Update : BOOL;
// END_VAR
//
// (* processing enabled? *)
// IF (Xus_TN_INPUT_CONTROL.bo_Enable = FALSE) OR (Xus_TN_SCREEN.bo_Menue_Bar_Dialog = TRUE) THEN
// 	RETURN;
// END_IF;
//
// (* perform focus reset *)
// IF (Xus_TN_INPUT_CONTROL.bo_Reset_Fokus = TRUE) THEN
// 	Xus_TN_INPUT_CONTROL.bo_Reset_Fokus := FALSE;
//
// 	FOR in_index := INT#01 TO Xus_TN_INPUT_CONTROL.in_Count DO
// 		Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].bo_Focus := FALSE;
// 	END_FOR;
//
// 	IF (Xus_TN_INPUT_CONTROL.in_Count >= INT#01) THEN
// 		Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[1].bo_Focus := TRUE;
// 		Xus_TN_INPUT_CONTROL.in_Focus_at := INT#01;
// 	END_IF;
//
// 	Xus_TN_INPUT_CONTROL.bo_Update_all := TRUE;
// END_IF;
//
// IF (Xus_TN_SCREEN.bo_Modal_Dialog = FALSE) THEN
// 	(* evaluate TAB key for switching control element *)
// 	IF (Xus_TN_SCREEN.by_Input_Exten_Code = BYTE#09) THEN
//
// 		in_Scroll_Offset := INT#01;
//
// 		(* cursor up key *)
// 	ELSIF (Xus_TN_SCREEN.by_Input_Exten_Code = BYTE#65) THEN
//
// 		in_Scroll_Offset := -1;
//
// 		(* cursor down key *)
// 	ELSIF (Xus_TN_SCREEN.by_Input_Exten_Code = BYTE#66) THEN
//
// 		in_Scroll_Offset := INT#01;
//
// 	END_IF;
// END_IF;
//
// IF (in_Scroll_Offset <> INT#00) THEN
// 	(* check limits                                            *)
//
// 	Xus_TN_INPUT_CONTROL.in_Focus_at:=LIMIT(INT#01,Xus_TN_INPUT_CONTROL.in_Focus_at,Xus_TN_INPUT_CONTROL.in_Count);
//
// 	in_Focus_new := Xus_TN_INPUT_CONTROL.in_Focus_at + in_Scroll_Offset;
//
// 	IF (in_Focus_new > Xus_TN_INPUT_CONTROL.in_Count) THEN
// 		in_Focus_new := INT#1;
// 	END_IF;
//
// 	IF (in_Focus_new < INT#01) THEN
// 		in_Focus_new := Xus_TN_INPUT_CONTROL.in_Count;
// 	END_IF;
//
// 	in_Focus_alt := Xus_TN_INPUT_CONTROL.in_Focus_at;
//
// 	(* remove focus from the current element *)
// 	Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_Focus_alt].bo_Focus := FALSE;
// 	Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_Focus_alt].bo_Update_Input := TRUE;
//
// 	(* give focus to the new element *)
// 	Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_Focus_new].bo_Focus  := TRUE;
// 	Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_Focus_new].bo_Update_Input := TRUE;
//
// 	Xus_TN_INPUT_CONTROL.in_Focus_at := in_Focus_new;
//
// 	bo_Tool_Tip_Update := TRUE;
// 	in_Scroll_Offset := INT#00;
//
// END_IF;
//
// (* update all elements? *)
// IF (Xus_TN_INPUT_CONTROL.bo_Update_all = TRUE) THEN
// 	FOR in_index := 1 TO Xus_TN_INPUT_CONTROL.in_Count DO
// 		Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].bo_Update_All := TRUE;
// 	END_FOR;
// 	Xus_TN_INPUT_CONTROL.bo_Update_all := FALSE;
// 	bo_Tool_Tip_Update := TRUE;
// END_IF;
//
// FOR in_index := 1 TO Xus_TN_INPUT_CONTROL.in_Count DO
//
// 	IF (Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].in_Type > INT#00) THEN
//
// 		IF (Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].bo_Focus = TRUE) THEN
//
// 			Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].by_Input_ASCII_Code  := Xus_TN_SCREEN.by_Input_ASCII_Code;
// 			Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].by_Input_Exten_Code  := Xus_TN_SCREEN.by_Input_Exten_Code;
// 			Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].bo_Input_ASCII_IsNum := Xus_TN_SCREEN.bo_Input_ASCII_IsNum;
//
// 			(* output tooltip text *)
// 			IF (bo_Tool_Tip_Update = TRUE) THEN
//
// 				FB_TN_SC_WRITE_C(	Iin_Y := Xus_TN_INPUT_CONTROL.in_ToolTip_Y,
// 									Iin_X := Xus_TN_INPUT_CONTROL.in_ToolTip_X,
// 									Iby_ATTR := Xus_TN_INPUT_CONTROL.by_ToolTip_Attr,
// 									Iin_OPTION := INT#02, (* Center-String *)
// 									Iin_LENGTH := Xus_TN_INPUT_CONTROL.in_ToolTip_Size,
// 									Ist_STRING := Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].st_Input_ToolTip,
// 									Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 				(* reset flag *)
// 				bo_Tool_Tip_Update := FALSE;
// 			END_IF;
// 		END_IF;
//
// 		CASE Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].in_Type OF
// 		1:
// 			(* process edit line *)
// 			FB_TN_INPUT_EDIT_LINE(	Xus_TN_INPUT_CONTROL_DATA := Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index],
// 									Xus_TN_SCREEN := Xus_TN_SCREEN);
// 		2:
// 			(* process select text *)
// 			FB_TN_INPUT_SELECT_TEXT(	Xus_TN_INPUT_CONTROL_DATA := Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index],
// 										Xus_TN_SCREEN := Xus_TN_SCREEN);
// 		3:
// 			(* process select popup *)
// 			FB_TN_INPUT_SELECT_POPUP(	Xus_TN_INPUT_CONTROL_DATA := Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index],
// 										Xus_TN_SCREEN := Xus_TN_SCREEN);
// 		END_CASE;
//
// 		IF (Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].bo_Focus = TRUE) THEN
//
// 			(* transfer cursor position of element to screen *)
// 			Xus_TN_SCREEN.in_Cursor_X := Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].in_Cursor_X;
// 			Xus_TN_SCREEN.in_Cursor_Y := Xus_TN_INPUT_CONTROL.usa_TN_INPUT_CONTROL_DATA[in_index].in_Cursor_Y;
//
// 		END_IF;
// 	END_IF;
// END_FOR;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CHR_TO_STRING, which is commented out.
// FUNCTION_BLOCK TN_INPUT_EDIT_LINE
// VAR
// 	in_count : INT;
// 	in_Count_End : INT;
// 	in_CursorOffset : INT;
// 	by_Attr : BYTE;
// 	st_String : STRING;
// 	bo_Pos_Found : BOOL;
// 	st_CHAR : STRING(1);
// 	in_old_charpos : INT;
// END_VAR
// VAR
// 	FB_TN_SC_WRITE : TN_SC_WRITE;
// END_VAR
// VAR_IN_OUT
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// 	Xus_TN_INPUT_CONTROL_DATA : us_TN_INPUT_CONTROL_DATA;
// END_VAR
//
// (* reset cursor control *)
// in_CursorOffset := INT#0;
//
// IF Xus_TN_INPUT_CONTROL_DATA.bo_Focus = TRUE THEN
//
// 	(* text input detected? *)
// 	IF Xus_TN_INPUT_CONTROL_DATA.by_Input_ASCII_Code > BYTE#0 THEN
//
// 		IF (Xus_TN_INPUT_CONTROL_DATA.bo_Input_Only_Num = TRUE AND Xus_TN_INPUT_CONTROL_DATA.bo_Input_ASCII_IsNum = TRUE)
// 			OR (Xus_TN_INPUT_CONTROL_DATA.bo_Input_Only_Num = FALSE) THEN
//
// 			(* convert ASCII value into a string *)
// 			st_CHAR := CHR_TO_STRING(C:=Xus_TN_INPUT_CONTROL_DATA.by_Input_ASCII_Code);
//
// 			(* accept input *)
// 			Xus_TN_INPUT_CONTROL_DATA.st_Input_Data := REPLACE(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data,
// 			st_CHAR,INT#1,Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos);
//
// 			in_CursorOffset := INT#1;
// 			Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// 		END_IF;
// 	END_IF;
//
// 	CASE BYTE_TO_INT(Xus_TN_INPUT_CONTROL_DATA.by_Input_Exten_Code) OF
//
// 	67: (* cursor right pressed *)
// 		in_CursorOffset := INT#1;
//
// 	68: (* cursor left pressed *)
// 		in_CursorOffset := -1;
//
// 	08: (* backspace pressed *)
// 		Xus_TN_INPUT_CONTROL_DATA.st_Input_Data := REPLACE(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data, ' ', 1,Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos);
// 		Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// 		in_CursorOffset := -1;
//
// 	13: (* enter key pressed *)
// 		IF LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask) > INT#00 THEN
// 			(* accept input string *)
// 			Xus_TN_INPUT_CONTROL_DATA.st_Input_String := Xus_TN_INPUT_CONTROL_DATA.st_Input_Data;
//
// 			(* input string -> remove leading and trailing spaces *)
// 			Xus_TN_INPUT_CONTROL_DATA.st_Input_String := TRIME(str:=Xus_TN_INPUT_CONTROL_DATA.st_Input_String);
// 			Xus_TN_INPUT_CONTROL_DATA.bo_Input_Entered := TRUE;
//
// 			(* prepare new input *)
// 			Xus_TN_INPUT_CONTROL_DATA.st_Input_Data := Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask;
//
// 			Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos := INT#1;
// 			Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
//
// 			(* search first position in input mask *)
// 			st_CHAR := MID(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask,1,Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos);
// 			IF st_CHAR <> ' ' THEN
// 				(* search next input position *)
// 				in_CursorOffset := INT#1;
// 			END_IF;
// 		END_IF;
// 	END_CASE;
// END_IF;
//
// (* move cursor to the right                                   *)
// IF (in_CursorOffset > INT#0) OR (Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos = INT#0) THEN
// 	bo_Pos_Found := FALSE;
//
// 	in_Count := Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos + INT#1;
// 	in_Count_End := LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask);
//
// 	WHILE (in_count <= in_Count_End) AND (bo_Pos_Found = FALSE) DO
//
// 		st_CHAR := MID(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask,1,in_count);
// 		IF st_CHAR = ' ' THEN
// 			Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos := in_count;
// 			Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// 			bo_Pos_Found := TRUE;
// 		END_IF;
//
// 		in_count := in_count + INT#1;
//
// 	END_WHILE;
// END_IF;
//
// (* move cursor to the left *)
// IF in_CursorOffset < INT#0 THEN
// 	bo_Pos_Found := FALSE;
//
// 	in_Count     := Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos - INT#01;
// 	in_Count_End := INT#1;
//
// 	WHILE (in_count >= in_Count_End) AND (bo_Pos_Found = FALSE) DO
//
// 		st_CHAR := MID(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask,1,in_count);
// 		IF st_CHAR = ' ' THEN
// 			Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos := in_count;
// 			Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// 			bo_Pos_Found := TRUE;
// 		END_IF;
//
// 		in_count := in_count - INT#1;
// 	END_WHILE;
// END_IF;
//
// (* maximum cursor position? *)
// IF Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos > LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data) THEN
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos := LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data);
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// END_IF;
//
// (* minimum cursor position? *)
// IF Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos < INT#1 THEN
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos := INT#1;
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// END_IF;
//
// (* update of the display *)
// IF (Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input = TRUE) OR (Xus_TN_INPUT_CONTROL_DATA.bo_Update_All = TRUE) THEN
//
// 	(* output title text *)
// 	IF (Xus_TN_INPUT_CONTROL_DATA.bo_Update_All = TRUE) THEN
//
// 		IF (Xus_TN_INPUT_CONTROL_DATA.in_Title_X_Offset = INT#00) AND (Xus_TN_INPUT_CONTROL_DATA.in_Title_Y_Offset = INT#00) THEN
// 			(* if offset X,Y = 0 then position the title string left of the element *)
// 			FB_TN_SC_WRITE.Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y;
// 			FB_TN_SC_WRITE.Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X - LEN(Xus_TN_INPUT_CONTROL_DATA.st_Title_String);
// 		ELSE
// 			(* output title text via row, color offset *)
// 			FB_TN_SC_WRITE.Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y + Xus_TN_INPUT_CONTROL_DATA.in_Title_Y_Offset;
// 			FB_TN_SC_WRITE.Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X + Xus_TN_INPUT_CONTROL_DATA.in_Title_X_Offset;
// 		END_IF;
//
// 		(* output title text via row, color offset *)
// 		FB_TN_SC_WRITE(	Iby_ATTR := Xus_TN_INPUT_CONTROL_DATA.by_Title_Attr,
// 						Ist_STRING := Xus_TN_INPUT_CONTROL_DATA.st_Title_String,
// 						Xus_TN_SCREEN := Xus_TN_SCREEN);
// 	END_IF;
//
// 	(* output input text *)
// 	IF (Xus_TN_INPUT_CONTROL_DATA.bo_Input_Hidden = TRUE) THEN
//
// 		st_String := '';
//
// 		IF LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data) = LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask) THEN
// 			(* input mode: HIDDEN *)
// 			FOR in_count := 1 TO LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data) DO
//
// 				(* read character from input mask *)
// 				st_CHAR := MID(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask,1,in_count);
//
// 				(* output input or mask character *)
// 				IF st_CHAR = ' ' THEN
//
// 					(* read character from input data *)
// 					st_CHAR := MID(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data,1,in_count);
//
// 					IF st_CHAR = ' ' THEN
// 						(* output spaces *)
// 						st_String := CONCAT(st_String, ' ');
//
// 					ELSE
// 						(* output hidden character *)
// 						st_String := CONCAT(st_String, '*');
//
// 					END_IF;
//
// 				ELSE
// 					(* take character from input mask *)
// 					st_CHAR := MID(Xus_TN_INPUT_CONTROL_DATA.st_Input_Data,1,in_count);
// 					st_String := CONCAT(st_String,st_CHAR);
//
// 				END_IF;
// 			END_FOR;
// 		END_IF;
// 	ELSE
//
// 		(* input mode: NORMAL *)
// 		st_String := Xus_TN_INPUT_CONTROL_DATA.st_Input_Data;
//
// 	END_IF;
//
// 	(* output string *)
//
// 	IF Xus_TN_INPUT_CONTROL_DATA.bo_Focus = TRUE THEN
// 		by_Attr := Xus_TN_INPUT_CONTROL_DATA.by_Attr_mF;
// 	ELSE
// 		by_Attr := Xus_TN_INPUT_CONTROL_DATA.by_Attr_oF;
// 	END_IF;
//
// 	FB_TN_SC_WRITE(	Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y,
// 					Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X,
// 					Iby_ATTR := by_Attr,
// 					Ist_STRING := st_String,
// 					Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 	(* set cursor to current position *)
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y;
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_X := Xus_TN_INPUT_CONTROL_DATA.in_X + Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Pos - INT#1;
//
// 	in_old_charpos := INT#0;
//
// 	(* reset update flags *)
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_All := FALSE;
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := FALSE;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses ELEMENT_GET, which is commented out.
// FUNCTION_BLOCK TN_INPUT_MENU_BAR
// VAR
// 	FB_TN_INPUT_MENU_POPUP : TN_INPUT_MENU_POPUP;
// 	us_TN_MENU_POPUP : us_TN_MENU_POPUP;
// 	FB_TN_SC_WRITE_C : TN_SC_WRITE_C;
// END_VAR
// VAR
// 	in_ROW : INT;
// 	in_COL : INT;
// 	in_TITLE_Length : INT;
// 	st_MENU_BAR : STRING(STRING_LENGTH);
// 	st_MENU_TITLE : STRING;
// 	in_Count : INT;
// 	in_Index : INT;
// 	by_New_KeyCode : BYTE;
// 	in_Scroll_Offset : INT;
// 	by_ATTR : BYTE;
// 	bo_Reset_Position :	BOOL;
// END_VAR
// VAR_IN_OUT
// 	Xus_TN_MENU : us_TN_MENU;
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// END_VAR
//
// (* make sure the rectangle is always positioned inside the screen *)
//
// IF (Xus_TN_MENU.in_State = INT#00) AND (Xus_TN_MENU.bo_Create = TRUE) THEN
// 	Xus_TN_MENU.in_State := INT#01;
// 	Xus_TN_MENU.bo_Create := FALSE;
// 	Xus_TN_MENU.in_Menu_Selected := INT#00;
// 	Xus_TN_MENU.bo_Update := TRUE;
// 	bo_Reset_Position := TRUE;
//
// 	(* menu bar - main menu *)
// 	st_MENU_BAR := ELEMENT_GET(SEP:=BYTE#37,POS:=0,ELEMENT:=Xus_TN_MENU.st_MENU_TEXT);
//
// 	(* menu bar - determine number of elements *)
// 	Xus_TN_MENU.in_Menu_E_Count := ELEMENT_COUNT(SEP:=BYTE#35,ELEMENT:=st_MENU_BAR);
//
// END_IF;
//
// (* store keycode *)
// by_New_KeyCode := Xus_TN_SCREEN.by_Input_Exten_Code;
//
// (* popup menu key control (only allow cursor up + down *)
// IF (by_New_KeyCode = BYTE#65) OR (by_New_KeyCode = BYTE#66) THEN
// 	us_TN_MENU_POPUP.by_Input_Exten_Code := by_New_KeyCode;
// ELSE
// 	us_TN_MENU_POPUP.by_Input_Exten_Code := BYTE#00;
// END_IF;
//
// (* Escape *)
// IF (by_New_KeyCode = BYTE#27) THEN
// 	(* menu was inactive -> escape -> activate menu *)
// 	IF (Xus_TN_MENU.in_State = INT#01) AND (Xus_TN_SCREEN.bo_Modal_Dialog = FALSE) THEN
// 		Xus_TN_MENU.in_State := INT#02;
// 		Xus_TN_MENU.bo_Update := TRUE;
// 		Xus_TN_MENU.in_Menu_Selected := INT#00;
//
// 		Xus_TN_SCREEN.bo_Menue_Bar_Dialog := TRUE;
//
// 		(* menu was active -> escape -> deactivate menu *)
// 	ELSIF (Xus_TN_MENU.in_State = INT#02) THEN
// 		us_TN_MENU_POPUP.bo_Destroy := TRUE;
// 		Xus_TN_MENU.bo_Update := TRUE;
// 		Xus_TN_MENU.in_State := INT#01;
// 		bo_Reset_Position := TRUE;
//
// 		Xus_TN_SCREEN.bo_Menue_Bar_Dialog := FALSE;
// 	END_IF;
// END_IF;
//
// (* menu active *)
// IF (Xus_TN_MENU.in_State = INT#02) THEN
//
// 	(* Return *)
// 	IF (by_New_KeyCode = BYTE#13)  THEN
// 		us_TN_MENU_POPUP.bo_Destroy := TRUE;
// 		Xus_TN_MENU.bo_Update := TRUE;
// 		Xus_TN_MENU.in_State := INT#01;
// 		bo_Reset_Position := TRUE;
// 		Xus_TN_MENU.in_Menu_Selected := (Xus_TN_MENU.in_Cur_Menu_Item * INT#10) + Xus_TN_MENU.in_Cur_Sub_Item;
// 		Xus_TN_SCREEN.by_Input_Exten_Code := BYTE#00;
//
// 		Xus_TN_SCREEN.bo_Menue_Bar_Dialog := FALSE;
//
// 		(* destroy command -> close popup menu *)
// 	ELSIF (Xus_TN_MENU.bo_Destroy = TRUE) THEN
// 		us_TN_MENU_POPUP.bo_Destroy := TRUE;
// 		Xus_TN_MENU.bo_Update := TRUE;
// 		Xus_TN_MENU.in_State := INT#01;
// 		bo_Reset_Position := TRUE;
// 		Xus_TN_MENU.in_Menu_Selected := INT#00;
// 		Xus_TN_SCREEN.by_Input_Exten_Code := BYTE#00;
//
// 		Xus_TN_SCREEN.bo_Menue_Bar_Dialog := FALSE;
//
// 	END_IF;
//
// 	IF (by_New_KeyCode = BYTE#67) THEN (* cursor right *)
// 		in_Scroll_Offset := INT#1;
// 		Xus_TN_MENU.bo_Update := TRUE;
// 		Xus_TN_SCREEN.by_Input_Exten_Code := BYTE#00;
//
// 	ELSIF (by_New_KeyCode = BYTE#68) THEN (* cursor left *)
// 		in_Scroll_Offset := -1;
// 		Xus_TN_MENU.bo_Update := TRUE;
// 		Xus_TN_SCREEN.by_Input_Exten_Code := BYTE#00;
// 	END_IF;
// END_IF;
//
// IF bo_Reset_Position = TRUE THEN
// 	bo_Reset_Position := FALSE;
// 	Xus_TN_MENU.in_X_SM_old := -1;
// 	Xus_TN_MENU.in_Y_SM_old := -1;
// 	Xus_TN_MENU.in_X_SM_new := -1;
// 	Xus_TN_MENU.in_Y_SM_new := -1;
// END_IF;
//
// (* update POPUP menu *)
// IF (Xus_TN_MENU.in_State > INT#00) AND (Xus_TN_MENU.bo_Update = TRUE) THEN
//
// 	Xus_TN_MENU.bo_Update := FALSE;
//
// 	in_ROW := Xus_TN_MENU.in_Y;
// 	in_COL := Xus_TN_MENU.in_X;
//
// 	(* Scrolloffset *)
// 	IF in_Scroll_Offset <> INT#00 THEN
// 		Xus_TN_MENU.in_Cur_Menu_Item := INC2(Xus_TN_MENU.in_Cur_Menu_Item,in_Scroll_Offset,1,Xus_TN_MENU.in_Menu_E_Count);
// 		in_Scroll_Offset := INT#00;
// 	ELSE
// 		(* item = 0 and no scroll offset *)
// 		IF (Xus_TN_MENU.in_Cur_Menu_Item = INT#00) THEN
// 			Xus_TN_MENU.in_Cur_Menu_Item := INT#01;
// 		END_IF;
// 	END_IF;
//
// 	(* draw all menu texts *)
// 	FOR in_Count := INT#01 TO Xus_TN_MENU.in_Menu_E_Count DO;
//
// 		(* determine character color *)
// 		IF (in_Count = Xus_TN_MENU.in_Cur_Menu_Item) AND (Xus_TN_MENU.in_State > INT#01 ) THEN
// 			by_ATTR := Xus_TN_MENU.by_Attr_mF;
//
// 			Xus_TN_MENU.in_Y_SM_new := in_ROW + INT#01;
// 			Xus_TN_MENU.in_X_SM_new := in_COL;
//
// 		ELSE
// 			by_ATTR := Xus_TN_MENU.by_Attr_oF;
// 		END_IF;
//
// 		st_MENU_TITLE := ELEMENT_GET(SEP:=BYTE#35,POS:=in_Count - 1,ELEMENT:=st_MENU_BAR);
//
// 		in_TITLE_Length := LEN(st_MENU_TITLE) + INT#02;
//
// 		FB_TN_SC_WRITE_C(	Iin_Y := in_ROW,
// 							Iin_X := in_COL,
// 							Iby_ATTR := by_ATTR,
// 							Iin_LENGTH := in_TITLE_Length,
// 							Iin_OPTION := INT#02,
// 							Ist_STRING := st_MENU_TITLE,
// 							Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 		in_COL := in_COL + in_TITLE_Length;
//
// 	END_FOR;
// END_IF;
//
// IF Xus_TN_MENU.in_X_SM_new <> Xus_TN_MENU.in_X_SM_old THEN
//
// 	(* submenu is active -> exit                                     *)
// 	IF us_TN_MENU_POPUP.bo_Activ = TRUE THEN
// 		us_TN_MENU_POPUP.bo_Destroy := TRUE;
//
// 		(* process submenu *)
// 		FB_TN_INPUT_MENU_POPUP(	Xus_TN_MENU_POPUP := us_TN_MENU_POPUP,
// 								Xus_TN_SCREEN := Xus_TN_SCREEN);
// 	END_IF;
//
// 	IF us_TN_MENU_POPUP.bo_Activ = FALSE THEN
// 		(* configure and run the popup menu for the current menu         *)
// 		in_Index := Xus_TN_MENU.in_Cur_Menu_Item;
//
// 		us_TN_MENU_POPUP.in_Y := Xus_TN_MENU.in_Y_SM_new;
// 		us_TN_MENU_POPUP.in_X := Xus_TN_MENU.in_X_SM_new;
// 		us_TN_MENU_POPUP.by_Attr_mF := Xus_TN_MENU.by_Attr_mF;
// 		us_TN_MENU_POPUP.by_Attr_oF := Xus_TN_MENU.by_Attr_oF;
// 		us_TN_MENU_POPUP.st_Menu_Text := ELEMENT_GET(SEP:=BYTE#37,POS:=in_Index,ELEMENT:=Xus_TN_MENU.st_MENU_TEXT);
// 		us_TN_MENU_POPUP.bo_Create := TRUE;
//
// 		Xus_TN_MENU.in_X_SM_old := Xus_TN_MENU.in_X_SM_new;
// 		Xus_TN_MENU.in_Y_SM_old := Xus_TN_MENU.in_Y_SM_new;
//
// 	END_IF;
// END_IF;
//
// (* call popup menu cyclically *)
// FB_TN_INPUT_MENU_POPUP(	Xus_TN_MENU_POPUP := us_TN_MENU_POPUP,
// 						Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// (* copy current item from popup menu *)
// Xus_TN_MENU.in_Cur_Sub_Item := us_TN_MENU_POPUP.in_Cur_Item;
//
// (* destroy command *)
// IF (Xus_TN_MENU.in_State = INT#01) THEN
// 	IF (Xus_TN_MENU.bo_Destroy = TRUE) THEN
// 		Xus_TN_MENU.in_State := INT#00;
// 		Xus_TN_MENU.bo_Destroy := FALSE;
// 		Xus_TN_MENU.in_Cur_Menu_Item := INT#00;
// 		Xus_TN_MENU.in_Cur_Sub_Item := INT#00;
// 	END_IF;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses ELEMENT_COUNT, which is commented out.
// FUNCTION_BLOCK TN_INPUT_MENU_POPUP
// VAR
// 	FB_TN_SC_AREA_SAVE : TN_SC_AREA_SAVE;
// 	FB_TN_SC_AREA_RESTORE : TN_SC_AREA_RESTORE;
// 	FB_TN_SC_LINE : TN_SC_LINE;
// 	FB_TN_SC_ADD_SHADOW : TN_SC_ADD_SHADOW;
// 	FB_TN_SC_BOX : TN_SC_BOX;
// 	FB_TN_SC_WRITE_C : TN_SC_WRITE_C;
// END_VAR
// VAR
// 	in_ROW : INT;
// 	in_ROWS : INT;
// 	in_COL : INT;
// 	in_COLS : INT;
// 	in_index : INT;
// 	in_Count : INT;
// 	in_Scroll_Offset : INT;
// 	by_ATTR : BYTE;
// END_VAR
// VAR_IN_OUT
// 	Xus_TN_MENU_POPUP : us_TN_MENU_POPUP;
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// END_VAR
//
// (* make sure the rectangle is always positioned inside the screen *)
//
// IF (Xus_TN_MENU_POPUP.bo_Activ = FALSE) AND (Xus_TN_MENU_POPUP.bo_Create = TRUE) THEN
// 	Xus_TN_MENU_POPUP.bo_Activ := TRUE;
// 	Xus_TN_MENU_POPUP.bo_Create := FALSE;
// 	Xus_TN_MENU_POPUP.bo_Destroy := FALSE;
// 	Xus_TN_MENU_POPUP.bo_Update := TRUE;
//
//   (* submenu texts *)
//
// 	Xus_TN_MENU_POPUP.in_Menu_E_Count := ELEMENT_COUNT(SEP:=BYTE#35,ELEMENT:=Xus_TN_MENU_POPUP.st_Menu_Text);
//
// 	(* determine number of lines (submenu texts) *)
// 	Xus_TN_MENU_POPUP.in_Rows := Xus_TN_MENU_POPUP.in_Menu_E_Count;
//
// 	(* find longest submenu text *)
// 	Xus_TN_MENU_POPUP.in_Cols := INT#00;
//
// 	FOR in_Count := INT#01 TO Xus_TN_MENU_POPUP.in_Rows DO;
// 		Xus_TN_MENU_POPUP.in_Cols:=MAX(LEN(ELEMENT_GET(SEP:=BYTE#35,POS:=in_Count - 1,ELEMENT:=Xus_TN_MENU_POPUP.st_Menu_Text)),Xus_TN_MENU_POPUP.in_Cols);
// 	END_FOR;
//
// 	(* check and correct horizontal position *)
// 	IF (Xus_TN_MENU_POPUP.in_X + Xus_TN_MENU_POPUP.in_Cols) > INT#75 THEN
// 		Xus_TN_MENU_POPUP.in_X := INT#79 - Xus_TN_MENU_POPUP.in_Cols - INT#04;
// 	END_IF;
//
//   (* check and correct vertical position *)
// 	IF (Xus_TN_MENU_POPUP.in_Y + Xus_TN_MENU_POPUP.in_ROWS) > INT#21 THEN
// 		Xus_TN_MENU_POPUP.in_Y := INT#23 - Xus_TN_MENU_POPUP.in_ROWS - INT#2;
// 	END_IF;
//
// 	Xus_TN_SCREEN.bo_Modal_Dialog := TRUE;
//
// 	(* save background *)
// 	FB_TN_SC_AREA_SAVE(	Iin_X1 := Xus_TN_MENU_POPUP.in_X,
// 									Iin_Y1 := Xus_TN_MENU_POPUP.in_Y,
// 									Iin_X2 := Xus_TN_MENU_POPUP.in_X + Xus_TN_MENU_POPUP.in_Cols + INT#04,
// 									Iin_Y2 := Xus_TN_MENU_POPUP.in_Y + Xus_TN_MENU_POPUP.in_Rows + INT#02,
// 									Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 	(* draw frame for submenu *)
// 	FB_TN_SC_BOX(	Iin_X1 := Xus_TN_MENU_POPUP.in_X,
// 					Iin_Y1 := Xus_TN_MENU_POPUP.in_Y,
// 					Iin_X2 := Xus_TN_MENU_POPUP.in_X + Xus_TN_MENU_POPUP.in_Cols + INT#03,
// 					Iin_Y2 := Xus_TN_MENU_POPUP.in_Y + Xus_TN_MENU_POPUP.in_Rows + INT#01,
// 					Iby_FILL := BYTE#32, (* fill *)
// 					Iin_BORDER := INT#01, (* Single-Line Border *)
// 					Iby_ATTR := Xus_TN_MENU_POPUP.by_Attr_oF,
// 					Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 	(* draw shadow for submenu *)
// 	FB_TN_SC_ADD_SHADOW(Iin_X1 := Xus_TN_MENU_POPUP.in_X + INT#1,
// 						Iin_Y1 := Xus_TN_MENU_POPUP.in_Y + INT#1,
// 						Iin_X2 := Xus_TN_MENU_POPUP.in_X + Xus_TN_MENU_POPUP.in_Cols + INT#04,
// 						Iin_Y2 := Xus_TN_MENU_POPUP.in_Y + Xus_TN_MENU_POPUP.in_Rows + INT#02,
// 						Iin_OPTION := INT#00,
// 						Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 	(* load parameters *)
// 	in_COL := Xus_TN_MENU_POPUP.in_X;
// 	in_ROW := Xus_TN_MENU_POPUP.in_Y;
// 	in_ROWS := Xus_TN_MENU_POPUP.in_Rows;
// 	in_COLS := Xus_TN_MENU_POPUP.in_Cols;
//
// 	(* draw separator lines *)
// 	in_COLS := in_COLS + INT#02; (* extend text by 2 characters *)
//
// 	FOR in_Count := INT#01 TO Xus_TN_MENU_POPUP.in_Menu_E_Count DO;
// 		in_ROW := in_ROW + INT#01;
//
// 		IF ELEMENT_GET(SEP:=BYTE#35,POS:=in_Count - 1,ELEMENT:=Xus_TN_MENU_POPUP.st_Menu_Text) = '-' THEN
//
// 			FB_TN_SC_LINE(	Iin_Y1 := in_ROW,
// 							Iin_X1 := in_COL,
// 							Iin_Y2 := in_ROW,
// 							Iin_X2 := in_COL + in_COLS + INT#01,
// 							Iby_ATTR := Xus_TN_MENU_POPUP.by_Attr_oF,
// 							Iby_BORDER := BYTE#01,
// 							Xus_TN_SCREEN := Xus_TN_SCREEN);
// 		END_IF;
// 	END_FOR;
// END_IF;
//
// IF (Xus_TN_MENU_POPUP.bo_Activ = TRUE) THEN
//
// 	CASE BYTE_TO_INT(Xus_TN_MENU_POPUP.by_Input_Exten_Code) OF
//
// 	65: (* cursor up *)
// 		in_Scroll_Offset := -1;
// 		Xus_TN_MENU_POPUP.bo_Update := TRUE;
//
// 	66: (* cursor down *)
// 		in_Scroll_Offset := 1;
// 		Xus_TN_MENU_POPUP.bo_Update := TRUE;
//
// 	13:  (* Enter *)
// 		Xus_TN_MENU_POPUP.bo_Destroy := TRUE;
//
// 	27:  (* Escape*)
// 		Xus_TN_MENU_POPUP.bo_Destroy := TRUE;
// 		Xus_TN_MENU_POPUP.in_Cur_Item := INT#00;
//
// 	67..68: (* cursor left or right *)
// 		Xus_TN_MENU_POPUP.bo_Destroy := TRUE;
// 		Xus_TN_MENU_POPUP.in_Cur_Item := INT#00;
//
// 	END_CASE;
// END_IF;
//
// (* update POPUP menu *)
// IF (Xus_TN_MENU_POPUP.bo_Activ = TRUE) AND (Xus_TN_MENU_POPUP.bo_Update = TRUE) THEN
//
// 	Xus_TN_MENU_POPUP.bo_Update := FALSE;
//
// 	(* check + find next selectable submenu *)
// 	in_index := Xus_TN_MENU_POPUP.in_Cur_Item;
//
// 	IF (in_index < INT#01) OR (in_index > Xus_TN_MENU_POPUP.in_Menu_E_Count) THEN
// 		in_Scroll_Offset := INT#01;
// 		in_index := Xus_TN_MENU_POPUP.in_Menu_E_Count;
// 	END_IF;
//
// 	(* current index is on a separator -> change *)
// 	IF ELEMENT_GET(SEP:=BYTE#35,POS:=in_index-1,ELEMENT:=Xus_TN_MENU_POPUP.st_Menu_Text) = '-' THEN
// 		in_Scroll_Offset := INT#01;
// 	END_IF;
//
// 	(* Scroll-Offset *)
// 	IF in_Scroll_Offset <> INT#00 THEN
//
// 		FOR in_Count := INT#01 TO Xus_TN_MENU_POPUP.in_Menu_E_Count DO;
//
// 			in_index:=INC2(in_index,in_Scroll_Offset,1,Xus_TN_MENU_POPUP.in_Menu_E_Count);
//
// 			IF ELEMENT_GET(SEP:=BYTE#35,POS:=in_index - 1,ELEMENT:=Xus_TN_MENU_POPUP.st_Menu_Text) <> '-' THEN
// 				EXIT;
// 			END_IF;
//
// 		END_FOR;
// 	END_IF;
//
// 	Xus_TN_MENU_POPUP.in_Cur_Item := in_index;
// 	in_Scroll_Offset := INT#00;
//
// 	(* load parameters *)
// 	in_COL  := Xus_TN_MENU_POPUP.in_X;
// 	in_ROW  := Xus_TN_MENU_POPUP.in_Y;
// 	in_ROWS := Xus_TN_MENU_POPUP.in_Rows;
// 	in_COLS := Xus_TN_MENU_POPUP.in_Cols;
//
// 	(* output SUB-MENU title *)
// 	in_COLS := in_COLS + INT#02; (* extend text by 2 characters *)
//
// 	FOR in_Count := INT#01 TO Xus_TN_MENU_POPUP.in_Menu_E_Count DO;
//
// 		in_ROW := in_ROW + INT#01;
//
// 		(* determine character color *)
// 		IF in_Count = Xus_TN_MENU_POPUP.in_Cur_Item THEN
// 			by_ATTR := Xus_TN_MENU_POPUP.by_Attr_mF;
// 		ELSE
// 			by_ATTR := Xus_TN_MENU_POPUP.by_Attr_oF;
// 		END_IF;
//
// 		FB_TN_SC_WRITE_C.Ist_STRING := ELEMENT_GET(SEP:=BYTE#35,POS:=in_Count - 1,ELEMENT:=Xus_TN_MENU_POPUP.st_Menu_Text);
//
// 		IF FB_TN_SC_WRITE_C.Ist_STRING <> '-' THEN
//
// 			FB_TN_SC_WRITE_C(	Iin_Y := in_ROW,
// 								Iin_X := in_COL + INT#01,
// 								Iby_ATTR := by_ATTR,
// 								Iin_LENGTH := in_COLS,
// 								Iin_OPTION := INT#02, (* Center *)
// 								Xus_TN_SCREEN := Xus_TN_SCREEN);
// 		END_IF;
// 	END_FOR;
// END_IF;
//
// (* remove POPUP menu (destroy) *)
// IF (Xus_TN_MENU_POPUP.bo_Activ = TRUE) AND (Xus_TN_MENU_POPUP.bo_Destroy = TRUE) THEN
//
// 	Xus_TN_MENU_POPUP.bo_Destroy := FALSE;
// 	Xus_TN_MENU_POPUP.bo_Activ := FALSE;
//
// 	(* restore old background from memory *)
// 	FB_TN_SC_AREA_RESTORE(Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 	Xus_TN_SCREEN.bo_Modal_Dialog := FALSE;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses TN_INPUT_MENU_POPUP, which is commented out.
// FUNCTION_BLOCK TN_INPUT_SELECT_POPUP
// VAR
// 	in_index : INT;
// 	in_Element_Count : INT;
// 	by_Cur_Key : BYTE;
// 	by_Attr : BYTE;
// 	st_String : STRING;
// END_VAR
// VAR
// 	us_TN_MENU_POPUP : us_TN_MENU_POPUP;
// 	FB_TN_INPUT_MENU_POPUP : TN_INPUT_MENU_POPUP;
// 	FB_TN_SC_WRITE : TN_SC_WRITE;
// END_VAR
// VAR_IN_OUT
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// 	Xus_TN_INPUT_CONTROL_DATA : us_TN_INPUT_CONTROL_DATA;
// END_VAR
//
// (* evaluate user input -> selection text n *)
// IF (Xus_TN_INPUT_CONTROL_DATA.bo_Focus = TRUE) THEN
//
// 	by_Cur_Key := Xus_TN_INPUT_CONTROL_DATA.by_Input_Exten_Code;
//
// 	IF (us_TN_MENU_POPUP.bo_Activ = FALSE) THEN
// 		us_TN_MENU_POPUP.by_Input_Exten_Code := BYTE#00;
//
// 		(* evaluate return key *)
// 		IF (by_Cur_Key = BYTE#13) THEN
// 			(* pass popup menu parameters *)
// 			us_TN_MENU_POPUP.in_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y;
// 			us_TN_MENU_POPUP.in_X := Xus_TN_INPUT_CONTROL_DATA.in_X;
// 			us_TN_MENU_POPUP.by_Attr_mF := Xus_TN_INPUT_CONTROL_DATA.by_Attr_mF;
// 			us_TN_MENU_POPUP.by_Attr_oF := Xus_TN_INPUT_CONTROL_DATA.by_Attr_oF;
// 			us_TN_MENU_POPUP.st_Menu_Text := Xus_TN_INPUT_CONTROL_DATA.st_Input_Data;
// 			us_TN_MENU_POPUP.bo_Create := TRUE;
// 			us_TN_MENU_POPUP.in_Cur_Item := INT#00;
// 		END_IF;
//
// 		IF (us_TN_MENU_POPUP.in_Cur_Item > INT#00) THEN
// 			Xus_TN_INPUT_CONTROL_DATA.in_selected := us_TN_MENU_POPUP.in_Cur_Item;
// 			Xus_TN_INPUT_CONTROL_DATA.bo_Input_Entered := TRUE;
// 			Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// 			us_TN_MENU_POPUP.in_Cur_Item := INT#00;
// 		END_IF;
//
// 	ELSE
// 		us_TN_MENU_POPUP.by_Input_Exten_Code := by_Cur_Key;
// 	END_IF;
//
// 	(* call popup menu cyclically *)
// 	FB_TN_INPUT_MENU_POPUP(Xus_TN_MENU_POPUP := us_TN_MENU_POPUP, Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// END_IF;
//
// (* update of the display *)
// IF (Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input = TRUE) OR (Xus_TN_INPUT_CONTROL_DATA.bo_Update_All = TRUE) THEN
//
// 	(* output title text *)
// 	IF Xus_TN_INPUT_CONTROL_DATA.bo_Update_All = TRUE THEN
//
// 		IF (Xus_TN_INPUT_CONTROL_DATA.in_Title_X_Offset = INT#00) AND (Xus_TN_INPUT_CONTROL_DATA.in_Title_Y_Offset = INT#00) THEN
// 			(* if offset X,Y = 0 then position the title string left of the element *)
// 			FB_TN_SC_WRITE.Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y;
// 			FB_TN_SC_WRITE.Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X - LEN(Xus_TN_INPUT_CONTROL_DATA.st_Title_String);
//
// 		ELSE
// 			(* output title text via row, color offset *)
// 			FB_TN_SC_WRITE.Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y + Xus_TN_INPUT_CONTROL_DATA.in_Title_Y_Offset;
// 			FB_TN_SC_WRITE.Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X + Xus_TN_INPUT_CONTROL_DATA.in_Title_X_Offset;
// 		END_IF;
//
// 		(* output title text via row, color offset *)
// 		FB_TN_SC_WRITE(	Iby_ATTR := Xus_TN_INPUT_CONTROL_DATA.by_Title_Attr,
// 						Ist_STRING := Xus_TN_INPUT_CONTROL_DATA.st_Title_String,
// 						Xus_TN_SCREEN := Xus_TN_SCREEN);
// 	END_IF;
//
// 	in_Element_Count := ELEMENT_COUNT(SEP:=BYTE#35,ELEMENT:=Xus_TN_INPUT_CONTROL_DATA.st_Input_Data);
//
// 	IF in_Element_Count > INT#00 THEN
//
// 		Xus_TN_INPUT_CONTROL_DATA.in_selected:=LIMIT(1,Xus_TN_INPUT_CONTROL_DATA.in_selected,in_Element_Count);
//
// 		in_index := Xus_TN_INPUT_CONTROL_DATA.in_selected;
//
// 		st_String := ELEMENT_GET(SEP:=BYTE#35,POS:=in_index - 1,ELEMENT:=Xus_TN_INPUT_CONTROL_DATA.st_Input_Data);
//
// 		st_String:= FIX(str:=st_String,
// 						L:=LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask),
// 						C:=BYTE#32,
// 						M:=Xus_TN_INPUT_CONTROL_DATA.in_Input_Option);
//
// 		Xus_TN_INPUT_CONTROL_DATA.st_Input_String := st_String;
//
// 		IF Xus_TN_INPUT_CONTROL_DATA.bo_Focus = TRUE THEN
// 			by_Attr := Xus_TN_INPUT_CONTROL_DATA.by_Attr_mF;
// 		ELSE
// 			by_Attr := Xus_TN_INPUT_CONTROL_DATA.by_Attr_oF;
// 		END_IF;
//
// 		FB_TN_SC_WRITE(	Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X,
// 								Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y,
// 								Iby_ATTR := by_Attr,
// 								Ist_STRING := st_String,
// 								Xus_TN_SCREEN := Xus_TN_SCREEN );
// 	END_IF;
//
// 	(* set cursor to current position *)
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_X := Xus_TN_INPUT_CONTROL_DATA.in_X;
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y;
//
// 	(* reset update flags *)
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_All   := FALSE;
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := FALSE;
//
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses ELEMENT_COUNT, which is commented out.
// FUNCTION_BLOCK TN_INPUT_SELECT_TEXT
// VAR
// 	in_count : INT;
// 	by_Attr : BYTE;
// 	st_String : STRING;
// END_VAR
// VAR
// 	FB_TN_SC_WRITE : TN_SC_WRITE;
// END_VAR
// VAR_IN_OUT
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// 	Xus_TN_INPUT_CONTROL_DATA : us_TN_INPUT_CONTROL_DATA;
// END_VAR
//
// (* evaluate user input -> selection text n *)
// IF Xus_TN_INPUT_CONTROL_DATA.bo_Focus = TRUE THEN
//
// (* evaluate space key *)
// 	IF Xus_TN_INPUT_CONTROL_DATA.by_Input_ASCII_Code = BYTE#32 THEN
// 		Xus_TN_INPUT_CONTROL_DATA.in_selected := Xus_TN_INPUT_CONTROL_DATA.in_selected + INT#01;
// 		Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := TRUE;
// 	END_IF;
//
// 	(* evaluate return key *)
// 	IF Xus_TN_INPUT_CONTROL_DATA.by_Input_Exten_Code = BYTE#13 THEN
// 		Xus_TN_INPUT_CONTROL_DATA.bo_Input_Entered := TRUE;
// 	END_IF;
//
// END_IF;
//
// (* update of the display *)
// IF (Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input = TRUE) OR (Xus_TN_INPUT_CONTROL_DATA.bo_Update_All = TRUE) THEN
//
// 	(* output title text *)
// 	IF Xus_TN_INPUT_CONTROL_DATA.bo_Update_All = TRUE THEN
//
// 		IF (Xus_TN_INPUT_CONTROL_DATA.in_Title_X_Offset = INT#00) AND (Xus_TN_INPUT_CONTROL_DATA.in_Title_Y_Offset = INT#00) THEN
// 			(* if offset X,Y = 0 then position the title string left of the element *)
// 			FB_TN_SC_WRITE.Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y;
// 			FB_TN_SC_WRITE.Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X - LEN(Xus_TN_INPUT_CONTROL_DATA.st_Title_String);
//
// 		ELSE
// 			(* output title text via row, color offset *)
// 			FB_TN_SC_WRITE.Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y + Xus_TN_INPUT_CONTROL_DATA.in_Title_Y_Offset;
// 			FB_TN_SC_WRITE.Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X + Xus_TN_INPUT_CONTROL_DATA.in_Title_X_Offset;
// 		END_IF;
//
// 		(* output title text via row, color offset *)
// 		FB_TN_SC_WRITE(	Iby_ATTR := Xus_TN_INPUT_CONTROL_DATA.by_Title_Attr,
// 						Ist_STRING := Xus_TN_INPUT_CONTROL_DATA.st_Title_String,
// 						Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 	END_IF;
//
// 	(* process input line *)
//
// 	(* determine number of elements *)
// 	in_count := ELEMENT_COUNT(SEP:=BYTE#35,ELEMENT:=Xus_TN_INPUT_CONTROL_DATA.st_Input_Data);
//
// 	IF in_Count > INT#00 THEN
//
// 		IF Xus_TN_INPUT_CONTROL_DATA.in_selected < INT#1 THEN
// 			Xus_TN_INPUT_CONTROL_DATA.in_selected := INT#01;
// 		END_IF;
//
// 		IF Xus_TN_INPUT_CONTROL_DATA.in_selected > in_count THEN
// 			Xus_TN_INPUT_CONTROL_DATA.in_selected := INT#01;
// 		END_IF;
//
// 		st_String := ELEMENT_GET(SEP:=BYTE#35,POS:=Xus_TN_INPUT_CONTROL_DATA.in_selected - 1,ELEMENT:=Xus_TN_INPUT_CONTROL_DATA.st_Input_Data);
//
// 		st_String := FIX(	str:=st_String,
// 							L:=LEN(Xus_TN_INPUT_CONTROL_DATA.st_Input_Mask),
// 							C:=BYTE#32,
// 							M:=Xus_TN_INPUT_CONTROL_DATA.in_Input_Option);
//
// 		Xus_TN_INPUT_CONTROL_DATA.st_Input_String := st_String;
//
// 		IF Xus_TN_INPUT_CONTROL_DATA.bo_Focus = TRUE THEN
// 			by_Attr := Xus_TN_INPUT_CONTROL_DATA.by_Attr_mF;
// 		ELSE
// 			by_Attr := Xus_TN_INPUT_CONTROL_DATA.by_Attr_oF;
// 		END_IF;
//
// 		FB_TN_SC_WRITE(	Iin_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y,
// 						Iin_X := Xus_TN_INPUT_CONTROL_DATA.in_X,
// 						Iby_ATTR := by_Attr,
// 						Ist_STRING := st_String,
// 						Xus_TN_SCREEN := Xus_TN_SCREEN);
// 	END_IF;
//
// 	(* set cursor to current position *)
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_Y := Xus_TN_INPUT_CONTROL_DATA.in_Y;
// 	Xus_TN_INPUT_CONTROL_DATA.in_Cursor_X := Xus_TN_INPUT_CONTROL_DATA.in_X;
//
// 	(* reset update flags *)
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_All := FALSE;
// 	Xus_TN_INPUT_CONTROL_DATA.bo_Update_Input := FALSE;
//
// END_IF;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_RECEIVE
VAR_IN_OUT
	R_BUF :	NETWORK_BUFFER;
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR
VAR
	REQ_Size : UINT;
	REQ1 : BYTE;
	REQ2 : BYTE;
	REQ3 : BYTE;
END_VAR

(* evaluate data reception *)
Xus_TN_SCREEN.by_Input_ASCII_Code := BYTE#00;
Xus_TN_SCREEN.by_Input_Exten_Code := BYTE#00;

IF R_BUF.SIZE > UINT#0 THEN
	(* evaluate received data *)
	REQ_Size := R_BUF.SIZE;
	REQ1 := R_BUF.BUFFER[0];
	REQ2 := R_BUF.BUFFER[1];
	REQ3 := R_BUF.BUFFER[2];

	IF REQ1 = BYTE#27 AND REQ_Size >= UINT#3 THEN (* Escape *)
		IF REQ2 = BYTE#91 THEN (* Cursor code *)
			IF (REQ3 >= BYTE#65 AND REQ3 <= BYTE#68) OR (REQ3 = BYTE#72) OR (REQ3 = BYTE#75) THEN (* 65..68 = cursor up,down,right,left / 72 = HOME / 75 = END *)
				Xus_TN_SCREEN.by_Input_Exten_Code := REQ3;
			END_IF;
		ELSIF REQ2 = BYTE#79 THEN (* special keys F1/F2/F3/F4 *)
			IF REQ3 >= BYTE#80 AND REQ3 <= BYTE#83 THEN
				Xus_TN_SCREEN.by_Input_Exten_Code := REQ3;
			END_IF;
		END_IF;
	ELSIF REQ1 >= BYTE#32 AND REQ1 <= BYTE#126 THEN
		Xus_TN_SCREEN.by_Input_ASCII_Code := REQ1;
		Xus_TN_SCREEN.bo_Input_ASCII_IsNum := ISC_NUM(REQ1); (* check for numeric input *)
	ELSIF REQ1 = BYTE#8 OR REQ1 = BYTE#9 OR REQ1 = BYTE#13 OR REQ1 = BYTE#27 THEN (* check for control characters 8=Backspace / 9=Tab / 13=Return / 27=Escape *)
		Xus_TN_SCREEN.by_Input_Exten_Code := REQ1;
	END_IF;

	(* set receive buffer data length to zero (destroy telegram *)
	R_BUF.SIZE := UINT#00;
END_IF;
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_ADD_SHADOW
VAR
	in_Offset : INT;
	in_X : INT;
	in_Y : INT;
END_VAR
VAR_INPUT
	Iin_Y1 : INT;
	Iin_X1 : INT;
	Iin_Y2 : INT;
	Iin_X2 : INT;
	Iin_OPTION : INT;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY2_ERROR(Iin_X1,Iin_Y1,Iin_X2,Iin_Y2) THEN RETURN; END_IF;

IF Iin_Y1 < 0 OR Iin_X1 < 0 OR Iin_Y2 < 0 OR Iin_X2 < 0 THEN RETURN; END_IF; (* check min values *)
IF Iin_Y1 > 23 OR Iin_X1 > 79 OR Iin_Y2 > 23 OR Iin_X2 > 79 THEN RETURN; END_IF; (* check max values *)

in_Offset := (Iin_Y1 * 80) + Iin_X2;

IF Iin_OPTION = INT#00 THEN
	(* create shadow by adjusting color *)
	FOR in_Y := Iin_Y1 TO Iin_Y2 DO
		Xus_TN_SCREEN.bya_COLOR[in_Offset] := TN_SC_SHADOW_ATTR(Xus_TN_SCREEN.bya_COLOR[in_Offset]);
		in_Offset := in_Offset + 80;
		Xus_TN_SCREEN.bya_Line_Update[in_Y] := TRUE; (* mark line for update *)
	END_FOR;
ELSE
	(* create shadow with black color and spaces *)
	FOR in_Y := Iin_Y1 TO Iin_Y2 DO
		Xus_TN_SCREEN.bya_COLOR[in_Offset] := BYTE#00;
		Xus_TN_SCREEN.bya_CHAR[in_Offset] := BYTE#32;
		in_Offset := in_Offset + 80;
		Xus_TN_SCREEN.bya_Line_Update[in_Y] := TRUE; (* mark line for update *)
	END_FOR;
END_IF;

in_Offset := (Iin_Y2 * 80) + Iin_X1;

IF Iin_OPTION = INT#00 THEN
	(* create shadow by adjusting color *)
	FOR in_X  := Iin_X1 TO Iin_X2 DO
		Xus_TN_SCREEN.bya_COLOR[in_Offset] := TN_SC_SHADOW_ATTR(Xus_TN_SCREEN.bya_COLOR[in_Offset]);
		in_Offset := in_Offset + 1;
	END_FOR;
ELSE
	(* create shadow with black color and spaces *)
	FOR in_X  := Iin_X1 TO Iin_X2 DO
		Xus_TN_SCREEN.bya_COLOR[in_Offset] := BYTE#00;
		Xus_TN_SCREEN.bya_CHAR[in_Offset] := BYTE#32;
		in_Offset := in_Offset + 1;
	END_FOR;
END_IF;
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_AREA_RESTORE
VAR
	in_Offset_1 : INT;
	in_Offset_2 : INT;
	in_Y : INT;
	in_X : INT;
	in_Y1 : INT;
	in_X1 : INT;
	in_Y2 : INT;
	in_X2 : INT;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF Xus_TN_SCREEN.bya_BACKUP[0] > BYTE#00 THEN

	(* read parameters *)
	in_X1 := BYTE_TO_INT(Xus_TN_SCREEN.bya_BACKUP[1]);
	in_Y1 := BYTE_TO_INT(Xus_TN_SCREEN.bya_BACKUP[2]);
	in_X2 := BYTE_TO_INT(Xus_TN_SCREEN.bya_BACKUP[3]);
	in_Y2 := BYTE_TO_INT(Xus_TN_SCREEN.bya_BACKUP[4]);

	IF TN_SC_XY2_ERROR(in_X1,in_Y1,in_X2,in_Y2) THEN RETURN; END_IF;
	in_Offset_2 := 5;

	(* restore CHAR and COLOR from BUFFER *)
	FOR in_Y := in_Y1 TO in_Y2 DO

		Xus_TN_SCREEN.bya_Line_Update[in_Y] := TRUE; (* release line for update *)
		in_Offset_1 := (in_Y * 80) + in_X1;

		FOR in_X := in_X1 TO in_X2 DO

			Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := Xus_TN_SCREEN.bya_BACKUP[in_Offset_2];
			in_Offset_2 := in_Offset_2 + INT#01;

			Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Xus_TN_SCREEN.bya_BACKUP[in_Offset_2];
			in_Offset_2 := in_Offset_2 + INT#01;
			in_Offset_1 := in_Offset_1 + INT#01;

		END_FOR;
	END_FOR;

	(* mark data as invalid *)
	Xus_TN_SCREEN.bya_BACKUP[0] := BYTE#00;
END_IF;
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_AREA_SAVE
VAR
	in_Offset_1 : INT;
	in_Offset_2 : INT;
	in_Y : INT;
	in_X : INT;
END_VAR
VAR_INPUT
	Iin_Y1 : INT;
	Iin_X1 : INT;
	Iin_Y2 : INT;
	Iin_X2 : INT;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY2_ERROR(Iin_X1,Iin_Y1,Iin_X2,Iin_Y2) THEN RETURN; END_IF;

(* store parameters *)
Xus_TN_SCREEN.bya_BACKUP[0] := BYTE#1;
Xus_TN_SCREEN.bya_BACKUP[1] := INT_TO_BYTE(Iin_X1);
Xus_TN_SCREEN.bya_BACKUP[2] := INT_TO_BYTE(Iin_Y1);
Xus_TN_SCREEN.bya_BACKUP[3] := INT_TO_BYTE(Iin_X2);
Xus_TN_SCREEN.bya_BACKUP[4] := INT_TO_BYTE(Iin_Y2);

in_Offset_2 := 5;
(* save CHAR and COLOR in BUFFER *)
FOR in_Y := Iin_Y1 TO Iin_Y2 DO

	in_Offset_1 := (in_Y * 80) + Iin_X1;
	FOR in_X := Iin_X1 TO Iin_X2 DO
		Xus_TN_SCREEN.bya_BACKUP[in_Offset_2] := Xus_TN_SCREEN.bya_CHAR[in_Offset_1];
		in_Offset_2 := in_Offset_2 + INT#01;
		Xus_TN_SCREEN.bya_BACKUP[in_Offset_2] := Xus_TN_SCREEN.bya_COLOR[in_Offset_1];
		in_Offset_2 := in_Offset_2 + INT#01;
		in_Offset_1 := in_Offset_1 + INT#01;
	END_FOR;
END_FOR;
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_BOX
VAR
	FB_TN_SC_FILL : TN_SC_FILL;
	by_BorderChar_L_T : BYTE;
	by_BorderChar_T : BYTE;
	by_BorderChar_B : BYTE;
	by_BorderChar_R_T : BYTE;
	by_BorderChar_L : BYTE;
	by_BorderChar_R : BYTE;
	by_BorderChar_L_B : BYTE;
	by_BorderChar_R_B : BYTE;
	in_Y : INT;
	in_X : INT;
	in_Border_Offset : INT;
	in_Offset_1 : INT;
	in_Offset_2 : INT;
END_VAR
VAR_INPUT
	Iin_Y1 : INT;
	Iin_X1 : INT;
	Iin_Y2 : INT;
	Iin_X2 : INT;
	Iby_FILL : BYTE;
	Iby_ATTR : BYTE;
	Iin_BORDER : INT;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY2_ERROR(Iin_X1,Iin_Y1,Iin_X2,Iin_Y2) THEN RETURN; END_IF;

(* determine border type *)
CASE Iin_Border OF
0:
	in_Border_Offset := 0;
1:
	by_BorderChar_L_T := BYTE#218; (* top left *)
	by_BorderChar_T := BYTE#196; (* top *)
	by_BorderChar_R_T := BYTE#191; (* top right *)
	by_BorderChar_L := BYTE#179; (* left *)
	by_BorderChar_R := BYTE#179; (* right *)
	by_BorderChar_B := BYTE#196; (* bottom *)
	by_BorderChar_L_B := BYTE#217; (* bottom left *)
	by_BorderChar_R_B := BYTE#192; (* bottom right *)
	in_Border_Offset := 1;
 2:
	by_BorderChar_L_T := BYTE#201;
	by_BorderChar_T := BYTE#205;
	by_BorderChar_R_T := BYTE#187;
	by_BorderChar_L := BYTE#186;
	by_BorderChar_R := BYTE#186;
	by_BorderChar_B := BYTE#205;
	by_BorderChar_L_B := BYTE#188;
	by_BorderChar_R_B := BYTE#200;
	in_Border_Offset := 1;
ELSE
	by_BorderChar_L_T := BYTE#32;
	by_BorderChar_T := BYTE#32;
	by_BorderChar_R_T := BYTE#32;
	by_BorderChar_L := BYTE#32;
	by_BorderChar_R := BYTE#32;
	by_BorderChar_B := BYTE#32;
	by_BorderChar_L_B := BYTE#32;
	by_BorderChar_R_B := BYTE#32;
	in_Border_Offset := 1;

END_CASE;

IF Iby_FILL > BYTE#0 THEN (* FILL character present? -> draw box *)

	(* draw fill function one character row smaller *)
	FB_TN_SC_FILL(	Iin_Y1 := Iin_Y1 + in_Border_Offset,
					Iin_X1 := Iin_X1 + in_Border_Offset,
					Iin_Y2 := Iin_Y2 - in_Border_Offset,
					Iin_X2 := Iin_X2 - in_Border_Offset,
					Iby_CHAR := Iby_FILL,
					Iby_Attr := Iby_ATTR,
					Xus_TN_SCREEN := Xus_TN_SCREEN);
END_IF;

IF Iin_BORDER = 0 THEN RETURN; END_IF;

(* output left and right border *)

in_Offset_1 := (Iin_Y1 * 80) + Iin_X1;
in_Offset_2 := (Iin_Y1 * 80) + Iin_X2;

FOR in_Y := Iin_Y1 TO Iin_Y2 DO

	(* release line for update *)
	Xus_TN_SCREEN.bya_Line_Update[in_Y] := TRUE;

	Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_BorderChar_L;
	Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;

	Xus_TN_SCREEN.bya_CHAR[in_Offset_2] := by_BorderChar_R;
	Xus_TN_SCREEN.bya_COLOR[in_Offset_2] := Iby_Attr;

	in_Offset_1 := in_Offset_1 + 80;
	in_Offset_2 := in_Offset_2 + 80;

END_FOR;

(* output bottom and top border *)
in_Offset_1 := (Iin_Y1 * 80) + Iin_X1;
in_Offset_2 := (Iin_Y2 * 80) + Iin_X1;

FOR in_X := Iin_X1 TO Iin_X2 DO

	Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_BorderChar_T;
	Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;

	Xus_TN_SCREEN.bya_CHAR[in_Offset_2] := by_BorderChar_B;
	Xus_TN_SCREEN.bya_COLOR[in_Offset_2] := Iby_Attr;

	in_Offset_1 := in_Offset_1 + INT#1;
	in_Offset_2 := in_Offset_2 + INT#1;

END_FOR;

(* output top left corner *)
in_Offset_1 := (Iin_Y1 * 80) + Iin_X1;
Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_BorderChar_L_T;
Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;

(* output top right corner *)
in_Offset_1 := (Iin_Y1 * 80) + Iin_X2;
Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_BorderChar_R_T;
Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;

(* output bottom left corner *)
in_Offset_1 := (Iin_Y2 * 80) + Iin_X1;
Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_BorderChar_R_B;
Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;

(* output bottom right corner *)
in_Offset_1 := (Iin_Y2 * 80) + Iin_X2;
Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_BorderChar_L_B;
Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_FILL
VAR
	in_X : INT;
	in_Y : INT;
	in_Offset : INT;
END_VAR
VAR_INPUT
	Iin_Y1 : INT;
	Iin_X1 : INT;
	Iin_Y2 : INT;
	Iin_X2 : INT;
	Iby_CHAR : BYTE;
	Iby_Attr : BYTE;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY2_ERROR(Iin_X1,Iin_Y1,Iin_X2,Iin_Y2) THEN RETURN; END_IF;

FOR in_Y := Iin_Y1 TO Iin_Y2 DO
	Xus_TN_SCREEN.bya_Line_Update[in_Y] := TRUE; (* release line for update *)
	in_Offset := (in_Y * 80) + Iin_X1;
	FOR in_X := Iin_X1 TO Iin_X2 DO
		Xus_TN_SCREEN.bya_CHAR[in_Offset]  := Iby_CHAR; (* enter character *)
		Xus_TN_SCREEN.bya_COLOR[in_Offset] := Iby_Attr; (* enter color *)
		in_Offset := in_Offset + 1;
	END_FOR;
END_FOR;
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_LINE
VAR
	by_Connect_Left_Char : BYTE;
	by_Connect_Right_Char : BYTE;
	by_Connect_Top_Char : BYTE;
	by_Connect_Bottom_Char : BYTE;
	by_Connect_Cruz_Char : BYTE;
	by_Vertikal_Line_Char : BYTE;
	by_Horzontal_Line_Char : BYTE;
	by_old_Char : BYTE;
	in_Offset_1 : INT;
	in_ROW_Index : INT;
	in_COL_Index : INT;
	by_Char : BYTE;
END_VAR
VAR_INPUT
	Iin_X1 : INT;
	Iin_Y1 : INT;
	Iin_X2 : INT;
	Iin_Y2 : INT;
	Iby_ATTR : BYTE;
	Iby_BORDER : BYTE;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY2_ERROR(Iin_X1,Iin_Y1,Iin_X2,Iin_Y2) THEN RETURN; END_IF;

(* determine border type *)
IF Iby_Border = BYTE#01 THEN (* Single Line *)
	by_Connect_Top_Char := BYTE#16#C2;
	by_Connect_Bottom_Char := BYTE#16#C1;
	by_Connect_Left_Char := BYTE#16#C3;
	by_Connect_Right_Char := BYTE#16#B4;
	by_Horzontal_Line_Char := BYTE#16#C4;
	by_Vertikal_Line_Char := BYTE#16#B3;
	by_Connect_Cruz_Char := BYTE#16#C5;
ELSIF Iby_Border = BYTE#02 THEN	(* double line *)
	by_Connect_Top_Char := BYTE#16#CB;
	by_Connect_Bottom_Char := BYTE#16#CA;
	by_Connect_Left_Char := BYTE#16#CC;
	by_Connect_Right_Char := BYTE#16#B9;
	by_Horzontal_Line_Char := BYTE#16#CD;
	by_Vertikal_Line_Char := BYTE#16#BA;
	by_Connect_Cruz_Char := BYTE#16#CE;
ELSE (* full line *)
	by_Connect_Top_Char := Iby_Border;
	by_Connect_Bottom_Char := Iby_Border;
	by_Connect_Left_Char := Iby_Border;
	by_Connect_Right_Char := Iby_Border;
	by_Horzontal_Line_Char := Iby_Border;
	by_Vertikal_Line_Char := Iby_Border;
	by_Connect_Cruz_Char := Iby_Border;
END_IF;

(* draw vertical line *)
IF Iin_X1 = Iin_X2 THEN
	in_Offset_1 := (Iin_Y1 * 80) + Iin_X1;
	FOR in_ROW_Index := Iin_Y1 TO Iin_Y2 DO
		by_old_Char := Xus_TN_SCREEN.bya_CHAR[in_Offset_1]; (* read existing character *)
		by_Char := by_Vertikal_Line_Char; (* preset default characters *)
		IF by_old_Char = by_Horzontal_Line_Char THEN (* intersection detected *)
			IF in_ROW_Index = Iin_Y1 THEN (* handle start of line *)
				by_Char := by_Connect_Top_Char;
			ELSIF in_ROW_Index = Iin_Y2 THEN (* handle end of line *)
				by_Char := by_Connect_Bottom_Char;
			ELSE (* handle crossing point *)
				by_Char := by_Connect_Cruz_Char;
			END_IF;
		END_IF;
		Xus_TN_SCREEN.bya_Line_Update[in_ROW_Index] := TRUE; (* release line for update *)
		Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_Char;
		Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;
		in_Offset_1 := in_Offset_1 + 80;
	END_FOR;

ELSIF Iin_Y1 = Iin_Y2 THEN (* draw horizontal line *)
	Xus_TN_SCREEN.bya_Line_Update[Iin_Y1] := TRUE; (* release line for update *)
	in_Offset_1 := (Iin_Y1 * 80) + Iin_X1;
	FOR in_COL_Index := Iin_X1 TO Iin_X2 DO
		by_old_Char := Xus_TN_SCREEN.bya_CHAR[in_Offset_1]; (* read existing character *)
		by_Char := by_Horzontal_Line_Char; (* preset default characters *)
		IF by_old_Char = by_Vertikal_Line_Char THEN (* intersection detected *)
			IF in_COL_Index = Iin_X1 THEN (* handle start of line *)
				by_Char := by_Connect_Left_Char;
			ELSIF in_COL_Index = Iin_X2 THEN (* handle end of line *)
				by_Char := by_Connect_Right_Char;
			ELSE (* handle crossing point *)
				by_Char := by_Connect_Cruz_Char;
			END_IF;
		END_IF;
		Xus_TN_SCREEN.bya_CHAR[in_Offset_1] := by_Char;
		Xus_TN_SCREEN.bya_COLOR[in_Offset_1] := Iby_Attr;
		in_Offset_1 := in_Offset_1 + 1;
	END_FOR;
END_IF;
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_READ_ATTR
VAR
	in_Offset : INT;
END_VAR
VAR_INPUT
	Iin_Y : INT;
	Iin_X : INT;
END_VAR
VAR_OUTPUT
	Oby_ATTR : BYTE;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY_ERROR(Iin_X,Iin_Y) THEN RETURN; END_IF;
in_Offset := (Iin_Y * 80) + Iin_X;
Oby_ATTR := Xus_TN_SCREEN.bya_COLOR[in_Offset]; (* read attributes *)
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_READ_CHAR
VAR
	in_Offset : INT;
END_VAR
VAR_INPUT
	Iin_Y : INT;
	Iin_X : INT;
END_VAR
VAR_OUTPUT
	Oby_CHAR : BYTE;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY_ERROR(Iin_X,Iin_Y) THEN RETURN; END_IF;
in_Offset := (Iin_Y * 80) + Iin_X;
Oby_CHAR := Xus_TN_SCREEN.bya_CHAR[in_Offset]; (* read character *)
END_FUNCTION_BLOCK

FUNCTION TN_SC_SHADOW_ATTR : BYTE
VAR_INPUT
	Iby_ATTR : BYTE;
END_VAR

(* change background color to black      *)
(* and switch on dark foreground colors *)
TN_SC_SHADOW_ATTR := ((Iby_ATTR AND BYTE#2#11110000) OR BYTE#2#00001000);
END_FUNCTION

// Commented out for beedance: it uses FIX, which is commented out.
// FUNCTION_BLOCK TN_SC_VIEWPORT
// VAR
// 	FB_TN_SC_WRITE : TN_SC_WRITE;
// 	FB_TON : TON;
// END_VAR
// VAR
// 	in_index : INT;
// 	in_Color : INT;
// 	by_Attr : BYTE;
// 	in_count : INT;
// END_VAR
// VAR_IN_OUT
// 	Xus_LOG_VIEWPORT : us_LOG_VIEWPORT;
// 	Xus_LOG_CONTROL : LOG_CONTROL;
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// END_VAR
// VAR_INPUT
// 	Iin_X : INT;
// 	Iin_Y : INT;
// 	Iin_Width : INT;
// 	Idw_ATTR_1 : DWORD;
// 	Idw_ATTR_2 : DWORD;
// 	Iti_TIME : TIME;
// END_VAR
//
// (* output LOG messages *)
//
// IF Xus_LOG_VIEWPORT.UPDATE AND FB_TON.Q THEN (* update requested and delay time elapsed *)
// 	Xus_LOG_VIEWPORT.UPDATE := FALSE;
// 	FOR in_count := 1 TO Xus_LOG_VIEWPORT.COUNT DO (* process all viewport lines *)
// 		in_index := Xus_LOG_VIEWPORT.LINE_ARRAY[in_count];
// 		(* read color from the option parameter >0,Level,Back,Front< *)
// 		in_Color := DWORD_TO_INT(Xus_LOG_CONTROL.MSG_OPTION[in_index] AND DWORD#2#1111);
//
// 		CASE in_Color OF
// 		0..3: by_Attr := BYTE_OF_DWORD(Idw_ATTR_1,INT_TO_BYTE(in_Color));
// 		4..7: by_Attr := BYTE_OF_DWORD(Idw_ATTR_2,INT_TO_BYTE(in_Color));
// 		ELSE
// 			by_Attr := BYTE_OF_DWORD(Idw_ATTR_1,BYTE#0);
// 		END_CASE;
//
// 		(* output message (text + color) *)
// 		FB_TN_SC_WRITE(	Iin_Y := Iin_Y + in_count - INT#01,
// 						Iin_X := Iin_X,
// 						Iby_ATTR := by_Attr,
// 						Ist_STRING := FIX(str := Xus_LOG_CONTROL.MSG[in_index], L := Iin_Width, C := BYTE#32, M := INT#00), (* pad or truncate string on the right *)
// 						Xus_TN_SCREEN := Xus_TN_SCREEN);
//
// 	END_FOR;
//
// 	(* stop/reset timer *)
// 	FB_TON.IN := FALSE;
// 	FB_TON.PT := Iti_TIME;
// 	FB_TON();
//
// END_IF;
//
// (* update time delay *)
// FB_TON.IN := TRUE;
// FB_TON.PT := Iti_TIME;
// FB_TON();
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CODE, which is commented out.
// FUNCTION_BLOCK TN_SC_WRITE
// VAR
// 	bo_Color : BOOL;
// 	in_Offset : INT;
// 	in_Char_Index : INT;
// 	in_String_Length : INT;
// END_VAR
// VAR_INPUT
// 	Iin_Y : INT;
// 	Iin_X : INT;
// 	Iby_ATTR : BYTE;
// 	Ist_STRING : STRING;
// END_VAR
// VAR_IN_OUT
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// END_VAR
//
// in_String_Length := LEN(Ist_STRING);
// IF TN_SC_XY_ERROR(Iin_X + in_String_Length - 1,Iin_Y) THEN RETURN; END_IF;
// Xus_TN_SCREEN.bya_Line_Update[Iin_Y] := TRUE; (* release line for update *)
// in_Offset := (Iin_Y * 80) + Iin_X; (* calculate offset of the first character *)
// bo_Color := Iby_ATTR <> BYTE#00;
//
// (* copy characters and attributes of the string *)
// FOR in_Char_Index := 1 TO in_String_Length DO
// 	IF bo_Color THEN Xus_TN_SCREEN.bya_COLOR[in_Offset] := Iby_ATTR; END_IF;
// 	Xus_TN_SCREEN.bya_CHAR[in_Offset] := CODE(STR:=Ist_STRING,POS:=in_Char_Index);
// 	in_Offset := in_Offset + INT#1;
// END_FOR;
//
// Xus_TN_SCREEN.in_EOS_Offset := in_Offset; (* save end of string offset *)
// END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_WRITE_ATTR
VAR
	in_Offset : INT;
END_VAR
VAR_INPUT
	Iin_Y : INT;
	Iin_X : INT;
	Iby_ATTR : BYTE;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY_ERROR(Iin_X,Iin_Y) THEN RETURN; END_IF;
in_Offset := (Iin_Y * 80) + Iin_X;
Xus_TN_SCREEN.bya_COLOR[in_Offset] := Iby_ATTR; (* write attributes *)
END_FUNCTION_BLOCK

// Commented out for beedance: it uses TN_SC_WRITE, which is commented out.
// FUNCTION_BLOCK TN_SC_WRITE_C
// VAR
// 	FB_TN_SC_WRITE : TN_SC_WRITE;
// END_VAR
// VAR_INPUT
// 	Iin_Y : INT;
// 	Iin_X : INT;
// 	Iby_ATTR : BYTE;
// 	Ist_STRING : STRING;
// 	Iin_LENGTH : INT;
// 	Iin_OPTION : INT;
// END_VAR
// VAR_IN_OUT
// 	Xus_TN_SCREEN : us_TN_SCREEN;
// END_VAR
//
// (* output text *)
// FB_TN_SC_WRITE( Iin_Y := Iin_Y,
// 				Iin_X := Iin_X,
// 				Iby_ATTR := Iby_ATTR,
// 				Ist_STRING := FIX(str:=Ist_STRING,L:=Iin_LENGTH,C:=BYTE#32,M:=Iin_OPTION),
// 				Xus_TN_SCREEN := Xus_TN_SCREEN);
// END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_WRITE_CHAR
VAR
	in_Offset : INT;
END_VAR
VAR_INPUT
	Iin_Y : INT;
	Iin_X : INT;
	Iby_CHAR : BYTE;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

IF TN_SC_XY_ERROR(Iin_X,Iin_Y) THEN RETURN; END_IF;
in_Offset := (Iin_Y * 80) + Iin_X;
Xus_TN_SCREEN.bya_CHAR[in_Offset] := Iby_CHAR; (* write character *)
END_FUNCTION_BLOCK

FUNCTION_BLOCK TN_SC_WRITE_EOS
VAR
	TN_SC_WRITE : TN_SC_WRITE;
END_VAR
VAR_INPUT
	Iby_ATTR : BYTE;
	Ist_STRING : STRING;
END_VAR
VAR_IN_OUT
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR

TN_SC_WRITE(Iin_Y := Xus_TN_SCREEN.in_EOS_Offset / 80, (* calculate current line *)
			Iin_X := Xus_TN_SCREEN.in_EOS_Offset MOD 80, (* calculate current column *)
			Iby_ATTR := Iby_ATTR,
			Ist_STRING := Ist_STRING,
			Xus_TN_SCREEN := Xus_TN_SCREEN);
END_FUNCTION_BLOCK

FUNCTION TN_SC_XY2_ERROR : BOOL
VAR_INPUT
	X1 : INT;
	Y1 : INT;
	X2 : INT;
	Y2 : INT;
END_VAR

(* check min and max values *)
TN_SC_XY2_ERROR := Y1 < 0 OR X1 < 0 OR Y2 < 0 OR X2 < 0 OR Y1 > 23 OR X1 > 79 OR Y2 > 23 OR X2 > 79;
END_FUNCTION

FUNCTION TN_SC_XY_ERROR : BOOL
VAR_INPUT
	X : INT;
	Y : INT;
END_VAR

(* check min and max values *)
TN_SC_XY_ERROR := Y < 0 OR X < 0 OR Y > 23 OR X > 79;
END_FUNCTION

FUNCTION_BLOCK TN_SEND_ROWS
VAR_IN_OUT
	IP_C : IP_C;
	S_BUF : NETWORK_BUFFER;
	Xus_TN_SCREEN : us_TN_SCREEN;
END_VAR
VAR_INPUT
	S_BUF_SIZE : INT;
END_VAR
VAR
	in_Line_Count : INT;(*counter to run through at most all 24 lines*)
	in_Col_Count : INT;
	in_Last_Row_Index : INT;
	bo_Break_Loop : BOOL;
	in_Row_Count : INT;
	in_Offset : INT;
	in_Write_Index : INT;
	in_Write_Index_Backup : INT;
	by_Cur_Color : BYTE;
	by_last_Color : BYTE;
	in_Write_Stop : INT;
END_VAR

(* connection present and no data sending active *)
IF (IP_C.C_STATE < BYTE#127 OR S_BUF.SIZE > UINT#0) THEN
	RETURN; (* no -> exit block *)
END_IF;

(* load current write index *)
in_Write_Index := UINT_TO_INT(S_BUF.SIZE) - 1;

(* buffer current write index *)
in_Write_Index_Backup := in_Write_Index;

(* rising edge connect on *)
IF (IP_C.C_STATE = BYTE#254) OR (Xus_TN_SCREEN.bo_Clear_Screen) THEN

	(* reset clear-screen message *)
	Xus_TN_SCREEN.bo_Clear_Screen := FALSE;

	(* mark all lines for update *)
	FOR in_Line_Count := INT#00 TO INT#23 DO
		Xus_TN_SCREEN.bya_Line_Update[in_Line_Count] := TRUE;
	END_FOR;

	(* ESC [0;30,40m  Reset Attributes + Black,Black *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#1B; (* ESC *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#5B; (* [ *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#30; (* 0 *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#33; (* 3 *)

	(* output ones digit of foreground color *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := (SHR(Xus_TN_SCREEN.by_Clear_Screen_Attr,4) OR BYTE#16#30);

	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#34; (* 4 *)

	(* output ones digit of background color *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := ((Xus_TN_SCREEN.by_Clear_Screen_Attr AND BYTE#16#07) OR BYTE#16#30);

	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#6D; (* m *)

	(* ESC[2J Clear screen and home cursor *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#1B; (* ESC *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#5B; (* [ *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#32; (* 2 *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#4A; (* J *)

	(* ESC [?7l Off *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#1B; (* ESC *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#5B; (* [ *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#3F; (* ? *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#37; (* 7 *)
	in_Write_Index := in_Write_Index + INT#1;
	S_BUF.BUFFER[in_Write_Index] := BYTE#16#6C; (* l *)

END_IF;

in_Last_Row_Index := INT#0;
bo_Break_Loop := FALSE;
in_Write_Stop := S_BUF_SIZE - 20;

(* run through all lines *)
FOR in_Line_Count := 0 TO 23 DO

	in_Offset := in_Row_Count * INT#80;

	(* check whether an update of a line is necessary, or after reconnect *)
	IF (Xus_TN_SCREEN.bya_Line_Update[in_Row_Count] = TRUE) THEN

		(* position write cursor *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#1B; (* ESC *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#5B; (* [ *)

		(* output tens digit of line number *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] :=
		SHR(INT_TO_BCDC(in_Row_Count + INT#1),INT#4) OR BYTE#16#30;

		(* output ones digit of line number *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] :=
		(INT_TO_BCDC(in_Row_Count + INT#1) AND BYTE#16#0F) OR BYTE#16#30;

		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#31; (* 1 *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#48; (* H *)

		(* run through all columns of a line *)
		FOR in_Col_Count := 0 TO 79 DO

			(* determine current color *)
			by_Cur_Color := Xus_TN_SCREEN.bya_COLOR[in_Offset];
			IF (by_Cur_Color <> by_last_Color) OR (in_Col_Count = INT#00) THEN

				(* color setting: output ANSI sequence ESC[X;Y;Zm *)
				in_Write_Index := in_Write_Index + INT#1;
				S_BUF.BUFFER[in_Write_Index] := BYTE#16#1B; (* ESC *)

				in_Write_Index := in_Write_Index + INT#1;
				S_BUF.BUFFER[in_Write_Index] := BYTE#16#5B; (* [ *)

				(* color setting: determine dark or bright and blinking *)

					IF (((by_Cur_Color AND BYTE#2#10001000) XOR (by_last_Color AND BYTE#2#10001000)) <> BYTE#00) OR (in_Col_Count = 0) THEN

						in_Write_Index := in_Write_Index + INT#1;
						S_BUF.BUFFER[in_Write_Index] := BYTE#16#30; (* 0 *)
						in_Write_Index := in_Write_Index + INT#1;
						S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)

						(* no dark colors wanted -> then bold & bright *)
						IF (by_Cur_Color AND BYTE#2#00001000) = BYTE#2#00000000 THEN

							(* switch to bold & bright colors *)
							in_Write_Index := in_Write_Index + INT#1;
							S_BUF.BUFFER[in_Write_Index] := BYTE#16#31; (* 1 *)
							in_Write_Index := in_Write_Index + INT#1;
							S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)

						END_IF;

						(* blinking colors wanted *)
						IF (by_Cur_Color AND BYTE#2#10000000) = BYTE#2#10000000 THEN

							(* switch to bold & bright colors *)
							in_Write_Index := in_Write_Index + INT#1;
							S_BUF.BUFFER[in_Write_Index] := BYTE#16#35; (* 5 *)
							in_Write_Index := in_Write_Index + INT#1;
							S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)

						END_IF;

					END_IF;

					(* output tens digit of foreground color *)
					in_Write_Index := in_Write_Index + INT#1;
					S_BUF.BUFFER[in_Write_Index] := BYTE#16#33; (* 3 *)

					(* output ones digit of foreground color *)
					in_Write_Index := in_Write_Index + INT#1;
					S_BUF.BUFFER[in_Write_Index] := ((SHR(by_Cur_Color,4) AND BYTE#2#00000111) OR BYTE#16#30);

					in_Write_Index := in_Write_Index + INT#1;
					S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)

					(* output tens digit of background color *)
					in_Write_Index := in_Write_Index + INT#1;
					S_BUF.BUFFER[in_Write_Index] := BYTE#16#34; (* 4 *)

					(* output ones digit of background color *)
					in_Write_Index := in_Write_Index + INT#1;
					S_BUF.BUFFER[in_Write_Index] := ((by_Cur_Color AND BYTE#2#00000111) OR BYTE#16#30);

					in_Write_Index := in_Write_Index + INT#1;
					S_BUF.BUFFER[in_Write_Index] := BYTE#16#6D; (* m *)

					(* store new color as last used *)
					by_last_Color := by_Cur_Color;

				END_IF;

				(* end of buffer reached early -> end *)
				IF (in_Write_Index > in_Write_Stop) THEN

					bo_Break_Loop := TRUE;
					EXIT;

				END_IF;

				(* transfer one TEXT character from the screen buffer to the send buffer *)
				in_Write_Index := in_Write_Index + INT#1;
				S_BUF.BUFFER[in_Write_Index] := Xus_TN_SCREEN.bya_CHAR[in_Offset];

				(* select next column *)
				in_Offset := in_Offset + INT#01;
			END_FOR;

			IF (bo_Break_Loop = TRUE) THEN;

				(* end of buffer reached early -> transfer only up to the last complete line *)
				in_Write_Index := in_Last_Row_Index;

				EXIT;
			END_IF;

			(* buffer is not full yet -> take new line *)
			in_Last_Row_Index := in_Write_Index;
			Xus_TN_SCREEN.bya_Line_Update[in_Row_Count] := FALSE;

		END_IF;

		(* address next line *)
		in_Row_Count := in_Row_Count + INT#01;

		(* after the last line start again from the beginning *)
		IF (in_Row_Count > INT#23) THEN
			in_Row_Count := INT#00;
		END_IF;
		(*-----------------------------------------------*)

	END_FOR;

	(* if all lines were done at once in one update cycle *)
	(* Row_Count is set to the first line *)
	(* meant to prevent ugly graphics flicker *)
	(* otherwise continue at the last line *)
	IF (bo_Break_Loop = FALSE) THEN
		in_Row_Count := INT#00;
	END_IF;

	(* position edit cursor *)
	IF (in_Write_Index <> in_Write_Index_Backup) THEN

		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#1B; (* ESC *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#5B; (* [ *)

		(* output tens digit of line number *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := SHR(INT_TO_BCDC(Xus_TN_SCREEN.in_Cursor_Y + 1),INT#4) OR BYTE#16#30;

		(* output ones digit of line number *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := (INT_TO_BCDC(Xus_TN_SCREEN.in_Cursor_Y + 1) AND BYTE#16#0F) OR BYTE#16#30;

		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#3B; (* ; *)

		(* output tens digit of column number *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := SHR(INT_TO_BCDC(Xus_TN_SCREEN.in_Cursor_X + 1),INT#4) OR BYTE#16#30;

		(* output ones digit of column number *)
		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := (INT_TO_BCDC(Xus_TN_SCREEN.in_Cursor_X + 1) AND BYTE#16#0F) OR BYTE#16#30;

		in_Write_Index := in_Write_Index + INT#1;
		S_BUF.BUFFER[in_Write_Index] := BYTE#16#48; (* H *)

END_IF;

(* store current write index -> trigger send *)
S_BUF.SIZE := INT_TO_UINT(in_Write_Index + 1);
END_FUNCTION_BLOCK

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION_BLOCK UNI_CIRCULAR_BUFFER
// VAR_IN_OUT
// 	DATA : UNI_CIRCULAR_BUFFER_DATA;
// END_VAR
// VAR_TEMP
// 	tmp_dw : DWORD;
// 	pTmp : REF_TO BYTE;
// 	n : INT;
// 	n1 : INT;
// 	n2 : INT;
// 	x: BYTE;
// 	buf_totalsize : UINT;
// END_VAR
// VAR CONSTANT
// 	HeaderSize : INT := 4;
// END_VAR
// VAR
// 	cnt : UINT;
// 	mode_last : INT;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := FALSE;*)
// END_VAR
//
// CASE DATA.D_MODE OF
//
// 01:	x := WORD_TO_BYTE(DATA.D_HEAD); (* evaluate only the right byte *)
// 	CASE x OF
// 	01:	cnt := LEN(DATA.D_STRING); (* process string *)
// 		pTmp := REF(DATA.D_STRING);
// 	02:	cnt := 4; (* process real *)
// 		pTmp := REF(DATA.D_REAL);
// 	03:	cnt := 4; (* process DWORD *)
// 		pTmp := REF(DATA.D_DWORD);
// 	ELSE
// 		cnt := 0;
// 	END_CASE;
//
// 	DATA.BUF_DATA_CNT := DATA.BUF_DATA_CNT + UDINT#1;
//
// 	(* ------ enter data into buffer ------ *)
// 	buf_totalsize := SIZEOF(DATA.BUF);
//
// 	IF (cnt + HeaderSize) > (buf_totalsize - DATA.BUF_SIZE) THEN (* enough memory free? *)
// 		DATA.BUF_DATA_LOST := DATA.BUF_DATA_LOST + UDINT#1;
// 		RETURN;
// 	END_IF;
//
// 	DATA.D_SIZE := cnt;
//
// 	mode_last := DATA.D_MODE;
// 	tmp_dw := SHL(WORD_TO_DWORD(DATA.D_HEAD),16) OR UINT_TO_DWORD(cnt);
//
// 	(* enter header into buffer *)
// 	n1 := 0;
// 	REPEAT
// 		n2 := MIN(HeaderSize - n1, buf_totalsize - DATA._Last);
// 		FOR n := 1 TO n2 DO
// 			tmp_dw := ROL(tmp_dw,8);
// 			data.BUF[DATA._Last] := DWORD_TO_BYTE(tmp_dw);
// 			DATA._Last := DATA._Last + 1;
// 		END_FOR;
// 		IF DATA._Last >= buf_totalsize THEN
// 			DATA._Last := 0;
// 		END_IF
// 		n1 := n1 + n2;
// 	UNTIL n1 >= HeaderSize
// 	END_REPEAT;
//
// 	(* enter data into buffer *)
// 	n1 := 0;
// 	REPEAT
// 		n2 := MIN(cnt - n1, buf_totalsize - DATA._Last);
// 		FOR n := 1 TO n2 DO
// 			data.BUF[DATA._Last] := pTmp^;
// 			pTmp := pTmp + 1;
// 			DATA._Last := DATA._Last + 1;
// 		END_FOR;
// 		IF DATA._Last >= buf_totalsize THEN
// 			DATA._Last := 0;
// 		END_IF
// 		n1 := n1 + n2;
// 	UNTIL n1 >= cnt
// 	END_REPEAT;
//
// 	DATA.BUF_SIZE := DATA.BUF_SIZE + cnt + HeaderSize;
// 	DATA.BUF_COUNT := DATA.BUF_COUNT + 1;
// 	DATA.BUF_USED := UINT_TO_USINT(100 * DATA.BUF_SIZE / buf_totalsize); (* used memory in % *)
// 	DATA.D_MODE := 0;
//
// 	(*####################  DEBUG-MESSAGE  ###################################*)
// 	(*IF _debug_enable THEN
// 		LOG_CL.NEW_MSG := 'RP: ADD /HEAD ~1 / SIZE ~2 / CNT ~3 / >~4<';
// 		LOG_CL.PRINTF[1] := DWORD_TO_STRH(DATA.D_HEAD);
// 		LOG_CL.PRINTF[2] := UINT_TO_STRING(cnt);
// 		LOG_CL.PRINTF[3] := UINT_TO_STRING(DATA.BUF_COUNT);
// 		CASE x OF
// 		01:	LOG_CL.PRINTF[4] := DATA.D_STRING;
// 		02:	LOG_CL.PRINTF[4] := REAL_TO_STRING(DATA.D_REAL);
// 		03:	LOG_CL.PRINTF[4] := DWORD_TO_STRH(DATA.D_DWORD);
// 		ELSE
// 			LOG_CL.PRINTF[4] := '';
// 		END_CASE;
// 		LOG_MSG();
// 	END_IF;*)
// 	(*########################################################################*)
//
// 10..12:	(* ------ read data from buffer and optionally delete ------ *)
// 		(* 10 = read only / 11 = delete after read / 12 = read and delete *)
// 	IF DATA.D_MODE <> 11 THEN
// 		IF DATA.BUF_COUNT = 0 THEN RETURN; END_IF;
// 		buf_totalsize := SIZEOF(DATA.BUF);
// 		DATA._GetEnd := DATA._First;(* remember index of the oldest entry *)
// 		mode_last := DATA.D_MODE;
//
// 		(* read header from buffer *)
// 		n1 := 0;
// 		tmp_dw := 0;
// 		REPEAT
// 			n2 := MIN(HeaderSize - n1, buf_totalsize - DATA._GetEnd);
// 			FOR n := 1 TO n2 DO
// 				tmp_dw := ROL(tmp_dw,8) OR BYTE_TO_DWORD(data.BUF[DATA._GetEnd]);
// 				DATA._GetEnd := DATA._GetEnd + 1;
// 			END_FOR;
// 			IF DATA._GetEnd >= buf_totalsize THEN
// 				DATA._GetEnd := 0;
// 			END_IF
// 			n1 := n1 + n2;
// 		UNTIL n1 >= HeaderSize
// 		END_REPEAT;
//
// 		DATA._GetStart := DATA._GetEnd;
// 		cnt := DWORD_TO_UINT(tmp_dw AND 16#0000_FFFF);
// 		DATA.D_SIZE := cnt;
// 		DATA.D_HEAD := DWORD_TO_WORD(SHR(tmp_dw,16));
//
// 		x := WORD_TO_BYTE(DATA.D_HEAD); (* evaluate only the right byte *)
// 		CASE x OF
// 		01:	pTmp := REF(DATA.D_STRING); (* process string *)
// 		02:	pTmp := REF(DATA.D_REAL); (* process real *)
// 		03:	pTmp := REF(DATA.D_DWORD); (* process DWORD *)
// 		END_CASE;
//
// 		(* read data from buffer *)
// 		n1 := 0;
// 		REPEAT
// 			n2 := MIN(cnt - n1, buf_totalsize - DATA._GetEnd);
// 			FOR n := 1 TO n2 DO
// 				pTmp^ := data.BUF[DATA._GetEnd];
// 				pTmp := pTmp + 1;
// 				DATA._GetEnd := DATA._GetEnd + 1;
// 			END_FOR;
// 			IF DATA._GetEnd >= buf_totalsize THEN
// 				DATA._GetEnd := 0;
// 			END_IF
// 			n1 := n1 + n2;
// 		UNTIL n1 >= cnt
// 		END_REPEAT;
// 		IF x = 1 THEN
// 			pTmp^ := 0; (* terminate string *)
// 		END_IF;
// 		DATA.D_MODE := 0;
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'RP: READ /HEAD ~1 / SIZE ~2 / CNT ~3 / >~4<';
// 			LOG_CL.PRINTF[1] := DWORD_TO_STRH(DATA.D_HEAD);
// 			LOG_CL.PRINTF[2] := UINT_TO_STRING(cnt);
// 			LOG_CL.PRINTF[3] := UINT_TO_STRING(DATA.BUF_COUNT);
// 			CASE x OF
// 			01:	LOG_CL.PRINTF[4] := DATA.D_STRING;
// 			02:	LOG_CL.PRINTF[4] := REAL_TO_STRING(DATA.D_REAL);
// 			03:	LOG_CL.PRINTF[4] := DWORD_TO_STRH(DATA.D_DWORD);
// 			ELSE
// 				LOG_CL.PRINTF[4] := '';
// 			END_CASE;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
// 	END_IF;
//
// 	IF (mode_last = 10 AND DATA.D_MODE = 11) OR (mode_last = 12 AND DATA.D_MODE = 0) THEN (* delete only possible directly after read *)
// 		mode_last := DATA.D_MODE;
// 		DATA.BUF_SIZE := DATA.BUF_SIZE - cnt - HeaderSize;
// 		DATA.BUF_COUNT := DATA.BUF_COUNT - 1;
// 		DATA._First := DATA._GetEnd;
// 		DATA.BUF_USED := UINT_TO_USINT(100 * DATA.BUF_SIZE / SIZEOF(DATA.BUF)); (* used memory in % *)
// 		DATA.D_MODE := 0;
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'RP: DEL /HEAD ~1 / SIZE ~2 / CNT ~3 / >~4<';
// 			LOG_CL.PRINTF[1] := DWORD_TO_STRH(DATA.D_HEAD);
// 			LOG_CL.PRINTF[2] := UINT_TO_STRING(cnt);
// 			LOG_CL.PRINTF[3] := UINT_TO_STRING(DATA.BUF_COUNT);
// 			CASE x OF
// 			01:	LOG_CL.PRINTF[4] := DATA.D_STRING;
// 			02:	LOG_CL.PRINTF[4] := REAL_TO_STRING(DATA.D_REAL);
// 			03:	LOG_CL.PRINTF[4] := DWORD_TO_STRH(DATA.D_DWORD);
// 			ELSE
// 				LOG_CL.PRINTF[4] := '';
// 			END_CASE;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
// 	END_IF;
//
// 99:	(* buffer reset - delete everything *)
// 	DATA._First := 0;
// 	DATA._Last := 0;
// 	DATA._GetStart := 0;
// 	DATA._GetEnd := 0;
// 	DATA.BUF_COUNT := 0;
// 	DATA.BUF_SIZE := 0;
// 	DATA.BUF_USED := 0;
// 	DATA.D_STRING := '';
// 	DATA.D_DWORD := 0;
// 	DATA.D_REAL := 0.0;
// 	DATA.D_HEAD := 0;
// 	DATA.D_MODE := 0;
// 	DATA.D_SIZE := 0;
// 	mode_last := 0;
// END_CASE;
//
// (* store highest utilization value *)
// IF DATA.BUF_USED > DATA.BUF_USED_MAX THEN
// 	DATA.BUF_USED_MAX := DATA.BUF_USED;
// END_IF;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses CHR_TO_STRING, which is commented out.
// FUNCTION URL_DECODE : STRING(STRING_LENGTH)
// VAR_INPUT
// 	IN : STRING(STRING_LENGTH);
// END_VAR
// VAR
// 	pos : INT;
// 	seq : STRING(3);
// 	new : STRING(1);
// END_VAR
//
// URL_DECODE := IN;
// pos := FIND(URL_DECODE,'%');
// WHILE pos > 0 DO
// 	seq := MID(URL_DECODE,2,pos + 1);
// 	new := CHR_TO_STRING(DWORD_TO_BYTE(HEX_TO_DWORD(seq)));
// 	URL_DECODE := REPLACE(URL_DECODE, new, 3, pos);
// 	pos := FIND(URL_DECODE,'%');
// END_WHILE;
// END_FUNCTION

// Commented out for beedance: it reads or writes the bytes of a STRING through a pointer; a beedance reference (REF_TO) refers to a whole variable of its type.
// FUNCTION URL_ENCODE : STRING(STRING_LENGTH)
// VAR_INPUT
// 	IN : STRING(STRING_LENGTH);
// END_VAR
// VAR
// 	pt_in : REF_TO ARRAY[1..255] OF BYTE;
// 	pt_out : REF_TO ARRAY[1..256] OF BYTE;
// 	pos_in : INT;
// 	pos_out : INT;
// 	stop : INT;
// 	tb : BYTE;
// END_VAR
//
// pt_in  := REF(in);
// pt_out := REF(URL_ENCODE);
//
// (* scan the input string and convert unallowed characters *)
// pos_in := 1;
// stop := LEN(in);
// FOR pos_out := 1 TO 250 DO
// 	IF pos_in > stop THEN
// 		EXIT;
// 	ELSIF IS_URLCHR(pt_in^[pos_in]) THEN
// 		pt_out^[pos_out] := pt_in^[pos_in];
// 		(* increment in by 1 *)
// 		pos_in := pos_in + 1;
// 	ELSE
// 		(* cannot fit 3 digits, exit *)
// 		IF pos_out > 248 THEN EXIT; END_IF;
// 		(* invalid character that needs to be converted *)
// 		pt_out^[pos_out] := 37;    (* start with % sign *)
// 		pos_out := pos_out + 1;
// 		(* need to add the two hex digits *)
// 		tb := SHR(pt_in^[pos_in],4);
// 		IF tb > 9 THEN
// 			pt_out^[pos_out] := tb + 55;
// 		ELSE
// 			pt_out^[pos_out] := tb + 48;
// 		END_IF;
// 		pos_out := pos_out + 1;
// 		tb := pt_in^[pos_in] AND 16#0F;
// 		IF tb > 9 THEN
// 			pt_out^[pos_out] := tb + 55;
// 		ELSE
// 			pt_out^[pos_out] := tb + 48;
// 		END_IF;
// 		(* increment in and out because out has to move 3 digits *)
// 		pos_in := pos_in + 1;
// 	END_IF;
// END_FOR;
//
// (* terminate the output string *)
// pt_out^[pos_out] := 0;
// END_FUNCTION

FUNCTION URL_TO_STRING : STRING(STRING_LENGTH)
VAR_INPUT
	IN : URL;
END_VAR

(* attach protocol *)
IF IN.PROTOCOL <> '' THEN URL_TO_STRING := CONCAT(IN.PROTOCOL,'://'); END_IF;

(* attach user and password *)
IF IN.USER <> '' THEN
	URL_TO_STRING := CONCAT(URL_TO_STRING, IN.USER);
	IF IN.PASSWORD <> '' THEN
		URL_TO_STRING := CONCAT(URL_TO_STRING, ':');
		URL_TO_STRING := CONCAT(URL_TO_STRING, IN.PASSWORD);
	END_IF;
	URL_TO_STRING := CONCAT(URL_TO_STRING, '@');
END_IF;

(* attach domain *)
URL_TO_STRING := CONCAT(URL_TO_STRING, IN.DOMAIN);

(* attach Port *)
IF IN.PORT > 0 THEN
	URL_TO_STRING := CONCAT(URL_TO_STRING, ':');
	URL_TO_STRING := CONCAT(URL_TO_STRING, WORD_TO_STRING(IN.PORT));
END_IF;

(* attach path *)
URL_TO_STRING := CONCAT(URL_TO_STRING, IN.PATH);

(* attach query *)
IF IN.QUERY <> '' THEN
	URL_TO_STRING := CONCAT(URL_TO_STRING, '?');
	URL_TO_STRING := CONCAT(URL_TO_STRING, IN.QUERY);
END_IF;

(* attach anchor *)
IF IN.ANCHOR <> '' THEN
	URL_TO_STRING := CONCAT(URL_TO_STRING, '#');
	URL_TO_STRING := CONCAT(URL_TO_STRING, IN.ANCHOR);
END_IF;
END_FUNCTION

// Commented out for beedance: it uses CSV_PARSER_BUF, which is commented out.
// FUNCTION_BLOCK WORLD_WEATHER
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// 	WW : WORLD_WEATHER_DATA;
// END_VAR
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	LATITUDE : REAL;
// 	LONGITUDE : REAL;
// 	KEY : STRING(30);
// END_VAR
// VAR_OUTPUT
// 	BUSY : BOOL;
// 	DONE : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	HTTP_GET : HTTP_GET;
// 	CPB : CSV_PARSER_BUF;
// 	last_state : BOOL;
// 	value_int : INT;
// 	value_real : REAL;
// 	v_real : REAL;
// 	state : INT;
// 	url_str : STRING(STRING_LENGTH);
// 	offset : UDINT;
// 	size : INT;
// 	sep : BYTE;
// 	cnt : INT;
// 	day : INT;
// 	idx : INT;
// 	run : BYTE;
// 	result : BYTE;
// 	value : STRING(STRING_LENGTH);
// 	pos : INT;
// END_VAR
//
// CASE state OF
//
// 00:
// 	IF ACTIVATE AND NOT last_state THEN
// 		state := 20;
// 		DONE := FALSE;
// 		BUSY := TRUE;
// 		ERROR_C := DWORD#0;
// 		ERROR_T := BYTE#0;
// 	END_IF;
//
// 20: (* URL for DNS AND HTTP-GET *)
// 	(* Example: http://api.worldweatheronline.com/free/v1/weather.ashx?q=48.21,16.37&format=csv&num_of_days=5&key=xxxxxxxxxxxxxxxxxxxxxx *)
// 	url_str := CONCAT('http://api.worldweatheronline.com/free/v1/weather.ashx?q=',REAL_TO_STRF(LATITUDE,6,D:='.'));
// 	url_str := CONCAT(url_str,',');
// 	url_str := CONCAT(url_str,REAL_TO_STRF(LONGITUDE,6,D:='.'));
// 	url_str := CONCAT(url_str,'&format=csv&num_of_days=5&key=');
// 	url_str := CONCAT(url_str,key);
//
// 	URL_DATA:=STRING_TO_URL(STR:=url_str, DEFAULT_PROTOCOL:='', DEFAULT_PATH:='' );
// 	state := 40;
//
// 40:
// 	IF DNS_CLIENT.DONE THEN
// 		state := 60;
// 	ELSIF (DNS_CLIENT.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := DNS_CLIENT.ERROR;
// 	ERROR_T := BYTE#01;
// 		state   := 100;
// 	END_IF;
//
// 60:
// 	IF HTTP_GET.DONE THEN
// 		sep := 0;
// 		idx := 0;
// 		day := 0;
// 		cnt := 0;
// 		OFFSET := HTTP_GET.BODY_START; (* set start offset *)
// 		size := HTTP_GET.BODY_STOP; (* End Position *)
// 		run := 1; (* search element *)
// 		state := 80;
//
// 	ELSIF (HTTP_GET.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := HTTP_GET.ERROR;
// 		ERROR_T := BYTE#02;
// 		state := 100;
// 	END_IF;
//
// 80:	IF run = 0 THEN (* no command active anymore? *)
// 		run := 1; (*  search next element  *)
// 		IF result = 1 OR result = 2 THEN
// 			value_int := 0;
// 			value_real := 0.0;
// 			(* perform automatic conversion to real and int *)
// 			IF LEN(value) <= 20 THEN
// 				v_real:=FLOAT_TO_REAL(FLT:=value);
// 				IF CHK_REAL(v_real) = 0 THEN (* !!! checks for a valid REAL value, otherwise TwinCAT crashes !!!!  *)
// 					value_real := v_real;
// 					value_int := REAL_TO_INT(value_real);
// 				END_IF;
// 			END_IF;
// 			cnt := cnt + 1;
//
// 			IF cnt = 1 AND LEFT(value,1) <> '#' THEN
// 				cnt := 0;
// 			END_IF;
//
// 			IF cnt >=8 AND cnt <= 87 THEN
// 				IF cnt > 22 THEN
// 					CASE idx OF
// 					00:	WW.DAY[day].DATE_OF_DAY := value;
// 					01:	WW.DAY[day].TEMP_MAX_C := value_int;
// 					02:	WW.DAY[day].TEMP_MAX_F := value_int;
// 					03:	WW.DAY[day].TEMP_MIN_C := value_int;
// 					04:	WW.DAY[day].TEMP_MIN_F := value_int;
// 					05:	WW.DAY[day].WIND_SPEED_MILES := value_int;
// 					06:	WW.DAY[day].WIND_SPEED_KMPH := value_int;
// 					07:	WW.DAY[day].WIND_DIR_DEGREE := value_int;
// 					08:	WW.DAY[day].WIND_DIR16POINT := value;
// 					09:	WW.DAY[day].WEATHER_CODE := value_int;
// 					10:	pos := FIND(value, '/wsymbol_');
// 						WW.DAY[day].WEATHER_ICON := DEC_TO_INT(MID(value,4,pos+9));
// 					11:	WW.DAY[day].WEATHER_DESC := value;
// 					12:	WW.DAY[day].PRECIPMM := value_real;
// 					END_CASE;
// 					idx := idx + 1;
// 					IF idx > 12 THEN
// 						idx := 0;
// 						day := day + 1;
// 					END_IF;
// 				ELSE
// 					CASE cnt OF
// 					08:	sep :=44; (* from now on work with separators, *)
// 					09:	WW.CUR.OBSERVATION_TIME := value;
// 					10:	WW.CUR.TEMP_C := value_int;
// 					11:	WW.CUR.WEATHER_CODE := value_int;
// 					12:	pos := FIND(value, '/wsymbol_');
// 						WW.CUR.WEATHER_ICON := DEC_TO_INT(MID(value,4,pos+9));
// 					13:	WW.CUR.WEATHER_DESC := value;
// 					14:	WW.CUR.WIND_SPEED_MILES := value_int;
// 					15:	WW.CUR.WIND_SPEED_KMPH := value_int;
// 					16:	WW.CUR.WIND_DIR_DEGREE := value_int;
// 					17:	WW.CUR.WIND_DIR16POINT := value;
// 					18:	WW.CUR.PRECIPMM := value_real;
// 					19:	WW.CUR.HUMIDITY := value_int;
// 					20:	WW.CUR.VISIBILITY := value_int;
// 					21:	WW.CUR.PRESSURE := value_int;
// 					22:	WW.CUR.CLOUDOVER := value_int;
// 					END_CASE;
// 				END_IF;
// 			END_IF;
// 		ELSIF result = 10 THEN (* nothing more found *)
// 			DONE  := TRUE;
// 			run := 0; (* end evaluation *)
// 			state := 100;
// 		END_IF;
// 	END_IF;
//
// 100:
// 	(* UNLOCK HTTP DATA *)
// 	IF (NOT HTTP_GET.DONE) THEN
// 		state := 0;
// 		BUSY := FALSE;
// 		DONE := ERROR_T = BYTE#0;
// 	END_IF;
//
// END_CASE;
//
// (* ---------------- DNS --------------- *)
// DNS_CLIENT(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			DOMAIN:=URL_DATA.DOMAIN,
// 			IP4_DNS:=DWORD#00,
// 			ACTIVATE:=state=40 (* DNS Start *)
// 			);
// (* parameters not used, or not used directly *)
// (*	DWORD :=DNS_CLIENT.IP4;
// 	BOOL  :=DNS_CLIENT.DONE;
// 	BYTE  :=DNS_CLIENT.ERROR; *)
// (* -------------------------------------*)
//
// (* ------------- HTTP_GET --------------- *)
// HTTP_GET(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			IP4:=DNS_CLIENT.IP4,
// 			GET:=state=60,
// 			MODE:=BYTE#2, (* HTTP 1.1 with persistent connection *)
// 			UNLOCK_BUF:=state=100,
// 			URL_DATA:=URL_DATA
// 			);
//
// (* parameters not used, or not used directly *)
// (*	STRING := HTTP_GET.HTTP_STATUS;
// 	UINT :=HTTP_GET.HEADER_START;
// 	UINT :=HTTP_GET.HEADER_STOP;
// 	UINT :=HTTP_GET.BODY_START;
// 	UINT :=HTTP_GET.BODY_STOP;
// 	BOOL :=HTTP_GET.DONE;
// 	BYTE :=HTTP_GET.ERROR;
// *)
// (* -------------------------------------*)
// CPB(SEP:=sep,RUN:=run,OFFSET:=OFFSET,VALUE:=value,PT:=R_BUF,RESULT=>result); (* CSV PARAMETER BUF *)
//
// last_state := ACTIVATE;
// END_FUNCTION_BLOCK

FUNCTION_BLOCK WORLD_WEATHER_DESC_GE
VAR_IN_OUT
	WW : WORLD_WEATHER_DATA;
END_VAR
VAR_INPUT
	ACTIVATE : BOOL;
END_VAR
VAR
	act_last : BOOL;
	x : INT;
	y : INT;
	str : STRING(60);
	i : INT;
END_VAR

IF ACTIVATE AND act_last = FALSE THEN
	i := 0;
	FOR y := 1 TO 12 DO
		CASE y OF
		1,3,5,7,9:		x :=  WW.DAY[i].WEATHER_CODE;
		2,4,6,8,10:	WW.DAY[i].WEATHER_DESC := str;
					i := i + 1;
		11:			x :=  WW.CUR.WEATHER_CODE;
		12:			WW.CUR.WEATHER_DESC := str;
		END_CASE;

		CASE x OF
			113	:	str := 'heiter/sonnig';
			116	:	str := 'stellenweise wolkig';
			119	:	str := 'wolkig';
			122	:	str := 'bewölkt';
			143	:	str := 'Nebel';
			176	:	str := 'stellenweise nahegelegener Regen';
			179	:	str := 'stellenweise nahegelegener Schneefall';
			182	:	str := 'stellenweise nahegelegener Eisregen';
			185	:	str := 'stellenweise gefrierender Niesel';
			200	:	str := 'nahegelegene Sturm';
			227	:	str := 'Schneetreiben';
			230	:	str := 'Schneesturm';
			248	:	str := 'Nebel';
			260	:	str := 'gefrierender Nebel';
			263	:	str := 'stellenweise leichter Nieselregen';
			266	:	str := 'leichter Niesel';
			281	:	str := 'gefrierender Nieselregen';
			284	:	str := 'stark gefrierender Nieselregen';
			293	:	str := 'stellenweise leichter Regen';
			296	:	str := 'leichter Regen';
			299	:	str := 'zeitweise etwas Regen';
			302	:	str := 'mäßiger Regen';
			305	:	str := 'zeitweise etwas starker Regen';
			308	:	str := 'starker Regen';
			311	:	str := 'leichter gefrierender Regen';
			314	:	str := 'mäßiger oder starker Eisregen';
			317	:	str := 'leichter Eisregen';
			320	:	str := 'mäßiger oder starker Eisregen';
			323	:	str := 'stellenweise leichter Schneefall';
			326	:	str := 'leichter Scheefall';
			329	:	str := 'stellenweise moderater Schneefall';
			332	:	str := 'mäßiger Schnee';
			335	:	str := 'stellenweise starker Scheefall';
			338	:	str := 'starker Schneefall';
			350	:	str := 'Hagel';
			353	:	str := 'leichte Regenfälle';
			356	:	str := 'mäßiger oder starker Regenfall';
			359	:	str := 'sintflutartige Regenfälle';
			362	:	str := 'leichter Eisregen';
			365	:	str := 'mäßiger oder starker Eisregen';
			368	:	str := 'leichte Schneeschauer';
			371	:	str := 'mäßiger oder starker Schneefall';
			374	:	str := 'leichter Hagel';
			377	:	str := 'mäßiger oder starker Hagel';
			386	:	str := 'stellenweise leichter Regen in Gewittergebiet';
			389	:	str := 'mäßiger oder starker Regen in Gewittergebiet';
			392	:	str := 'stellenweise leichter Schneefall in Gewittergebiet';
			395	:	str := 'mäßiger oder starker Schnee in Gewittergebiet';
		ELSE
			str := 	'nicht verfügbar';
		END_CASE;
	END_FOR;
END_IF;
act_last := ACTIVATE;
END_FUNCTION_BLOCK

FUNCTION_BLOCK WORLD_WEATHER_ICON_OSCAT
VAR_IN_OUT
	WW : WORLD_WEATHER_DATA;
END_VAR
VAR_INPUT
	ACTIVATE : BOOL;
END_VAR
VAR
	dat : ARRAY[0..41] OF INT := 0,14,12,9,9,0,7,7,0,2,2,3,3,1,0,0,17,4,4,3,3,1,0,0,17,2,2,3,3,1,0,0,17,2,4,3,3,1,0,0,17,0;
	act_last : BOOL;
	x : INT;
	i : INT;
END_VAR

IF ACTIVATE AND act_last = FALSE THEN
	x := LIMIT(0,WW.CUR.WEATHER_ICON,41);
	WW.CUR.WEATHER_ICON := dat[x];
	FOR i := 0 TO 4 DO
		x := LIMIT(0,WW.DAY[i].WEATHER_ICON,41);
		WW.DAY[i].WEATHER_ICON := dat[x];
	END_FOR;
END_IF;
act_last := ACTIVATE;
END_FUNCTION_BLOCK

// Commented out for beedance: it uses BUFFER_TO_STRING, which is commented out.
// FUNCTION_BLOCK XML_READER
// VAR_IN_OUT
// 	CTRL : XML_CONTROL;
// 	BUF : NW_BUF_LONG;
// END_VAR
// VAR
// 	index :	INT;
// 	index2 : INT;
// 	stop : INT;
// 	mode : INT;
// 	command : WORD;
// 	c :	BYTE;
// 	s1 : INT;
// 	e1 : INT;
// 	path_overflow : BOOL;
// 	empty_tag : BOOL;
// 	s_code : STRING(10);
// 	e_code : STRING(10);
// 	watchdog : TON;
//
// 	(*LOG_MSG : LOG_MSG;
// 	_debug_enable : BOOL := TRUE;*)
//
// END_VAR
//
// (*
//
// TYP
// 00 = UNKNOWN
// 01 = ELEMENT (STANDARD)
// 02 = CLOSE ELEMENT
// 03 = TEXT
// 04 = ATTRIBUTE
// 05 = ELEMENT (PROCESSING INSTRUCTION)
//
// 12 = !CDATA
// 13 = !COMMENT
// 14 = !DTD ELEMENT (Document Type Declarations )
//
// 98 = WATCHDOG
// 99 = END
//
// no block data is output for type 02!
//
// *)
//
// (* ---------- Init -------------- *)
// IF BIT_OF_DWORD(WORD_TO_DWORD(CTRL.COMMAND),15) THEN
// 	COMMAND := CTRL.COMMAND;
// 	CTRL.COMMAND := WORD#0;
//
// 	index := UINT_TO_INT(CTRL.START_POS);
// 	stop := UINT_TO_INT(CTRL.STOP_POS);
// 	mode := 100;
// 	path_overflow := FALSE;
//
// 	CTRL.TYP := 0;
// 	CTRL.COUNT := UINT#00;
// 	CTRL.LEVEL := UINT#00;
// 	CTRL.ATTRIBUTE := '';
// 	CTRL.ELEMENT := '';
// 	CTRL.PATH := '';
// 	CTRL.VALUE := '';
//
// 	WATCHDOG.PT := CTRL.WATCHDOG;
// END_IF;
//
// (*----------------------------------------------- *)
//
// IF index < 0 THEN RETURN;END_IF;
//
// CTRL.BLOCK1_START := UINT#0;
// CTRL.BLOCK1_STOP := UINT#0;
// CTRL.BLOCK2_START := UINT#0;
// CTRL.BLOCK2_STOP := UINT#0;
//
// WATCHDOG(IN:=FALSE); (* stop timer *)
//
// WHILE (TRUE) DO
//
// 	(* -------------- WATCHDOG ------------------ *)
// 	WATCHDOG(IN:=CTRL.WATCHDOG > T#0s);
// 	IF WATCHDOG.Q THEN
// 		CTRL.TYP := 98;
//     		(*####################  DEBUG-MESSAGE  ###################################*)
//     		(*IF _debug_enable THEN
//       			LOG_CL.NEW_MSG := 'WATCHDOG:';
//       			LOG_MSG();
//     		END_IF;*)
//     		(*########################################################################*)
// 		RETURN;
// 	END_IF;
//
// 	(* -------------- PARSING END ------------------ *)
// 	IF index > stop THEN
// 		CTRL.COUNT := CTRL.COUNT + UINT#1;
// 		CTRL.TYP := 99;
//     		(*####################  DEBUG-MESSAGE  ###################################*)
//     		(*IF _debug_enable THEN
//       			LOG_CL.NEW_MSG := 'EX: #T:~1 #C:~2 #P:~6 #L:~7';
//       			LOG_CL.PRINTF[1] := INT_TO_STRING(CTRL.TYP);
//       			LOG_CL.PRINTF[2] := UINT_TO_STRING(CTRL.COUNT);
//       			LOG_CL.PRINTF[6] := CTRL.PATH;
//       			LOG_CL.PRINTF[7] := UINT_TO_STRING(CTRL.LEVEL);
//       			LOG_MSG();
//     		END_IF;*)
//     		(*########################################################################*)
// 		mode  := 0; (* reset step *)
// 		index := 0; (* stop processing *)
// 		RETURN;
// 	END_IF;
//
// 	CASE mode OF
//
// 	100: (* ---------- ELEMENT SECTION ------------------ *)
// 		(* search begin of element *)
// 		s1 := 0;
// 		e1 := 0;
// 		WHILE (index <= stop) DO
// 			IF BUF[INDEX]=BYTE#60 THEN (* '<' TagBegin *)
// 				index := index + 1;
// 				C := BUF[index];
// 				IF C = BYTE#47 THEN (* '/' TagTerm *)
// 					mode := 300; (* End Element Section *)
// 				ELSIF C = BYTE#33 THEN (* '!' for Special Section *)
// 					index := index + 1;
// 					mode := 500;
// 				ELSE
// 					s1 := index;
// 				END_IF;
// 				EXIT;
// 			END_IF;
// 			index := index + 1;
// 		END_WHILE;
//
// 		(* analyse element *)
// 		CTRL.TYP := 0; (* clear TYPE *)
// 		IF s1 > 0 THEN
// 			WHILE (index <= stop AND e1 = 0) DO
// 				c := BUF[index];
// 				(* whitespace character = element with attributes *)
// 				IF c <= BYTE#32 THEN
// 					(* search next printable character *)
// 					index2 := index;
// 					WHILE (index2 <= stop) AND (BUF[index2] <= BYTE#32) DO
// 						index2 := index2 + 1;
// 					END_WHILE;
// 					IF BUF[index2] = BYTE#47 THEN (* check for empty element '/' *)
// 						index := index2;
// 						e1 := index - 2;
// 						CTRL.TYP := 1; (* structure element *)
// 						mode := 300;   (* Close Element *)
// 						EXIT;
// 					END_IF;
//
// 					e1 := index - 1;
// 					CTRL.TYP := 1; (* element with attributes *)
// 					mode := 400;
//
// 				ELSIF c = BYTE#62 THEN (* '>' TagEnd *)
//
// 					(* check for empty element <xx/> *)
// 					index2 := index - 1;
// 					IF BUF[index2] = BYTE#47 THEN (* TagTerm '/' *)
// 						e1 := index2 - 1;
// 						CTRL.TYP := 1;
// 						EMPTY_TAG := TRUE; (* EMPTY TAG Mode *)
// 						mode := 200; (* TEXT Element *)
// 						EXIT;
// 					END_IF;
//
// 					e1 := index - 1;
// 					index2 := index + 1;
//
// 					(* examine following element *)
// 					WHILE (index2 <= stop) DO
// 						c := BUF[index2];
// 						IF c = BYTE#60 THEN (* '<' TagBegin *)
//
// 							index2 := index2 + 1; (* unfortunately no other way on PCWORX!*)
// 							c := BUF[index2];
// 							index2 := index2 - 1;
// 							IF c = BYTE#47 THEN (* TagTerm '/' *)
// 								empty_tag := TRUE;
// 								mode := 200; (* TEXT Element *)
// 							ELSE
// 								mode := 100; (* Element *)
// 							END_IF;
// 							CTRL.TYP := 1; (* structure element *)
// 							EXIT;
// 						ELSIF c > BYTE#32  THEN (* Printable-Character *)
// 							CTRL.TYP := 1; (* element with text *)
// 							mode := 200;
// 							EXIT;
// 						END_IF;
// 						index2 := index2 + 1;
// 					END_WHILE;
// 				END_IF;
// 				index := index + 1;
// 			END_WHILE;
// 		END_IF;
//
// 		(* new element found*)
// 		IF CTRL.TYP > 0 THEN
//
// 			(* special check for processing instruction (PI) *)
// 			IF BUF[s1] = BYTE#63 THEN (* '?' for Processing Instruction *)
// 				s1 := s1 + 1; (* skip leading '?' *)
// 				CTRL.TYP := 5;
// 			END_IF;
// 			CTRL.BLOCK1_START := INT_TO_UINT(s1);
// 			CTRL.BLOCK1_STOP  := INT_TO_UINT(e1);
// 			CTRL.ELEMENT := BUFFER_TO_STRING(PT:=REF(BUF),SIZE:=INT_TO_UINT(stop + 1),START:=CTRL.BLOCK1_START,STOP:=CTRL.BLOCK1_STOP);
// 			CTRL.COUNT := CTRL.COUNT + UINT#1;
// 			CTRL.LEVEL := CTRL.LEVEL + UINT#1;
//
// 			IF path_overflow = FALSE THEN
// 				(* monitor string length *)
// 				IF LEN(CTRL.PATH) + LEN(CTRL.ELEMENT) + 1 > 250 THEN
// 					path_overflow := TRUE;
// 					CTRL.PATH := 'OVERFLOW';
// 				ELSE
// 					CTRL.PATH := CONCAT(CTRL.PATH,'/');
// 					CTRL.PATH  := CONCAT(CTRL.PATH,CTRL.ELEMENT);
// 				END_IF;
// 			END_IF;
//
// 			IF BIT_OF_DWORD(WORD_TO_DWORD(COMMAND),CTRL.TYP) THEN
//         				(*####################  DEBUG-MESSAGE  ###################################*)
//         				(*IF _debug_enable THEN
//           				LOG_CL.NEW_MSG := 'EL: #T:~1 #C:~2 #E:~3 #P:~6 #L:~7';
//           				LOG_CL.PRINTF[1] := INT_TO_STRING(CTRL.TYP);
//           				LOG_CL.PRINTF[2] := UINT_TO_STRING(CTRL.COUNT);
//           				LOG_CL.PRINTF[3] := CTRL.ELEMENT;
//           				LOG_CL.PRINTF[6] := CTRL.PATH;
//           				LOG_CL.PRINTF[7] := UINT_TO_STRING(CTRL.LEVEL);
//           				LOG_MSG();
//         				END_IF;*)
//         				(*########################################################################*)
// 				RETURN;
// 			END_IF;
// 		END_IF;
//
// 	200: (* ---------- TEXT SECTION ------------------ *)
// 		CTRL.VALUE := '';
// 		IF empty_tag = FALSE THEN
// 			s1 := index;
// 			(* search '<' end tag*)
// 			WHILE (index <= stop) AND (BUF[index] <> BYTE#60) DO
// 				index := index + 1;
// 			END_WHILE;
//
// 			E1 := index - 1;
// 			CTRL.BLOCK1_START := INT_TO_UINT(s1);
// 			CTRL.BLOCK1_STOP  := INT_TO_UINT(e1);
// 			CTRL.VALUE := BUFFER_TO_STRING(PT:=REF(BUF),SIZE:=INT_TO_UINT(stop + 1),START:=CTRL.BLOCK1_START,STOP:=CTRL.BLOCK1_STOP);
// 		END_IF;
//
// 		empty_tag := FALSE;
// 		mode := 300; (* -> NEXT-JOB: END ELEMENT *)
// 		CTRL.COUNT := CTRL.COUNT + UINT#1;
// 		CTRL.TYP := 3;
//
// 		IF BIT_OF_DWORD(WORD_TO_DWORD(command),CTRL.TYP) THEN
//       			(*####################  DEBUG-MESSAGE  ###################################*)
//       			(*IF _debug_enable THEN
//         				LOG_CL.NEW_MSG := 'TE: #T:~1 #C:~2 #E:~3 #V:~4 #P:~6 #L:~7';
//         				LOG_CL.PRINTF[1] := INT_TO_STRING(CTRL.TYP);
//         				LOG_CL.PRINTF[2] := UINT_TO_STRING(CTRL.COUNT);
//         				LOG_CL.PRINTF[3] := CTRL.ELEMENT;
//         				LOG_CL.PRINTF[4] := CTRL.VALUE;
//         				LOG_CL.PRINTF[6] := CTRL.PATH;
//         				LOG_CL.PRINTF[7] := UINT_TO_STRING(CTRL.LEVEL);
//         				LOG_MSG();
//       			END_IF;*)
//       			(*########################################################################*)
// 			RETURN;
// 		END_IF;
//
// 	300: (* ---------- END ELEMENT SECTION ------------------ *)
// 		WHILE (index <= stop) DO
// 			c := BUF[index];
// 			IF c = BYTE#62 THEN (* '>' TagEnd *)
// 				index := index + 1;
// 				s1 := FINDB(str1:=CTRL.PATH,str2:='/');
// 				IF path_overflow = FALSE THEN
// 					IF s1 > 1 THEN
// 						CTRL.ELEMENT := RIGHT(CTRL.PATH,LEN(CTRL.PATH) - s1);
// 						CTRL.PATH := LEFT(CTRL.PATH,s1 - 1);
// 					ELSE
// 						CTRL.ELEMENT := RIGHT(CTRL.PATH,LEN(CTRL.PATH) - 1);
// 						CTRL.PATH :='';
// 					END_IF;
// 				END_IF;
//
// 				CTRL.LEVEL := CTRL.LEVEL - UINT#1;
// 				CTRL.COUNT := CTRL.COUNT + UINT#1;
// 				CTRL.TYP := 2;
// 				mode := 100;
//
// 				IF BIT_OF_DWORD(WORD_TO_DWORD(COMMAND),CTRL.TYP) THEN
//          	 				(*####################  DEBUG-MESSAGE  ###################################*)
//           				(*IF _debug_enable THEN
//             					LOG_CL.NEW_MSG := 'CE: #T:~1 #C:~2 #E:~3 #P:~6 #L:~7';
//             					LOG_CL.PRINTF[1] := INT_TO_STRING(CTRL.TYP);
//             					LOG_CL.PRINTF[2] := UINT_TO_STRING(CTRL.COUNT);
//             					LOG_CL.PRINTF[3] := CTRL.ELEMENT;
//             					LOG_CL.PRINTF[6] := CTRL.PATH;
//             					LOG_CL.PRINTF[7] := UINT_TO_STRING(CTRL.LEVEL);
//             					LOG_MSG();
//           				END_IF;*)
//           				(*########################################################################*)
// 					RETURN;
// 				ELSE
// 					EXIT;
// 				END_IF;
// 			END_IF;
// 			index := index + 1;
// 		END_WHILE;
//
// 	400: (* ---------- ATTRIBUTE SECTION ------------------ *)
// 		(* search first normal charakter *)
// 		WHILE (index <= stop AND BUF[index] <= BYTE#32) DO
// 			index := index + 1;
// 		END_WHILE;
//
// 		(* begin Attribe name *)
// 		e1 := 0;
// 		s1 := index;
// 		(* search end of Attrib Name '=' *)
// 		WHILE (index <= stop) DO
// 			IF BUF[index] = BYTE#61 THEN (* '=' Equals *)
// 				e1 := index - 1;
// 				EXIT;
// 			END_IF;
// 			index := index + 1;
// 		END_WHILE;
//
// 		IF e1 > 0 THEN
// 			(* Attribute Name *)
// 			CTRL.BLOCK1_START := INT_TO_UINT(s1);
// 			CTRL.BLOCK1_STOP  := INT_TO_UINT(e1);
// 			CTRL.ATTRIBUTE := BUFFER_TO_STRING(PT:=REF(BUF),SIZE:=INT_TO_UINT(STOP + 1),START:=CTRL.BLOCK1_START,STOP:=CTRL.BLOCK1_STOP);
// 			(* search end of Attrib Value '=' *)
// 			index := index + 2;
// 			e1 := 0;
// 			s1 := INDEX;
// 			WHILE (index <= stop AND e1 = 0) DO
// 				c := BUF[index];
// 				IF c = BYTE#34 OR c = BYTE#39 THEN (* " OR ' Equals *)
// 					e1 := index - 1;
// 				END_IF;
// 				index := index + 1;
// 			END_WHILE;
//
// 			IF e1 > 0 THEN
// 				(* Value *)
// 				CTRL.BLOCK2_START := INT_TO_UINT(s1);
// 				CTRL.BLOCK2_STOP  := INT_TO_UINT(e1);
//
// 				(* for empty "" start is greater than stop and so it becomes an empty string *)
// 				CTRL.VALUE := BUFFER_TO_STRING(PT:=REF(BUF),SIZE:=INT_TO_UINT(stop + 1),START:=CTRL.BLOCK2_START,STOP:=CTRL.BLOCK2_STOP);
// 				(* search next non white space charakter *)
// 				WHILE (index <= stop AND BUF[index] <= BYTE#32) DO
// 					index := index + 1;
// 				END_WHILE;
//
// 					c := BUF[index];
// 					(* TagEnd '>' *)
// 					IF c = BYTE#62 THEN
// 						index := index + 1;
// 						index2 := index;
//
// 					(* search next non white space charakter *)
// 					WHILE (index <= stop AND BUF[index] <= BYTE#32) DO
// 						index := index + 1;
// 					END_WHILE;
//
// 					IF BUF[index] = BYTE#60 THEN (* TagBegin '<' *)
// 						mode := 100; (* -> NEXT-JOB: ELEMENT-SECTION *)
// 					ELSE
// 						index := index2;
// 						mode := 200; (* -> NEXT-JOB: TEXT-SECTION *)
// 					END_IF;
//
// 				ELSIF c = BYTE#47 OR c = BYTE#63 THEN (* TagTerm '/' OR '?' end Processing Instruction *)
// 					mode := 300; (* -> NEXT-JOB: END ELEMENT *)
// 				ELSE
// 					mode := 400; (* -> NEXT-JOB: ATTRIBUTE-SECTION *)
// 				END_IF;
//
// 				CTRL.COUNT := CTRL.COUNT + UINT#1;
// 				CTRL.TYP := 4;
// 				IF BIT_OF_DWORD(WORD_TO_DWORD(command),CTRL.TYP) THEN
//           				(*####################  DEBUG-MESSAGE  ###################################*)
//           				(*IF _debug_enable THEN
//             					LOG_CL.NEW_MSG := 'AT: #T:~1 #C:~2 #E:~3 #A:~4 #V:~5 #P:~6 #L:~7';
//             					LOG_CL.PRINTF[1] := INT_TO_STRING(CTRL.TYP);
//             					LOG_CL.PRINTF[2] := UINT_TO_STRING(CTRL.COUNT);
//             					LOG_CL.PRINTF[3] := CTRL.ELEMENT;
//             					LOG_CL.PRINTF[4] := CTRL.ATTRIBUTE;
//             					LOG_CL.PRINTF[5] := CTRL.VALUE;
//             					LOG_CL.PRINTF[6] := CTRL.PATH;
//             					LOG_CL.PRINTF[7] := UINT_TO_STRING(CTRL.LEVEL);
//             					LOG_MSG();
//           				END_IF;*)
//           				(*########################################################################*)
// 					RETURN;
// 				END_IF;
// 			END_IF;
// 		END_IF;
//
// 	500: (* ---------- SPECIAL DATA SECTION (Document Type Declarations) ---------- *)
//
// 		CTRL.TYP := 0;
// 		s1 := 0;
// 		index2 := 11;
//
// 		WHILE (s1 = 0) DO
// 			index2 := index2 + 1;
// 			CTRL.TYP := index2;
//
// 			CASE index2 OF
// 			12: (* CDATA *)
// 				s_code := '[CDATA[';
// 				e_code := ']]>';
// 			13: (* COMMENT *)
// 				s_code := '--';
// 				e_code := '-->';
// 			14: (* ANY DTD *)
// 				e_code := '';
// 				s1 := index;
// 				EXIT;
// 			END_CASE;
//
// 			e1 := LEN(s_code);
// 			IF BUFFER_SEARCH(PT:=REF(BUF),SIZE:=INT_TO_UINT(index + e1),STR:=S_CODE,POS:=index,IGN:=FALSE) = INDEX THEN
// 				s1 := index + e1;
// 			END_IF;
// 		END_WHILE;
//
// 		CTRL.BLOCK1_START := INT_TO_UINT(s1);
// 		IF LEN(e_code) > 0 THEN
// 			(* search end marker *)
// 			e1 := BUFFER_SEARCH(PT:=REF(BUF),SIZE:=INT_TO_UINT(stop + 1),STR:=e_code,POS:=index,IGN:=FALSE);
// 			CTRL.BLOCK1_STOP := INT_TO_UINT(e1 - 1);
// 			index := e1 + LEN(e_code); (* next element *)
// 		ELSE
// 			(* search '>' end tag*)
// 			index2 := 1; (* <> pair counter *)
// 			WHILE (index <= stop AND index2 > 0) DO
// 				c := BUF[index];
// 				IF c = BYTE#60 THEN (* '<' TagBegin *)
// 					index2 := index2 + 1;
// 				ELSIF c = BYTE#62 THEN (* '>' TagEnd   *)
// 					index2 := index2 - 1;
// 				ELSIF c = BYTE#34 OR c = BYTE#39 THEN (* " Double Quote OR ' Single Quote *)
// 					index := index + 1;
// 					WHILE (index <= stop) DO
// 						c := BUF[index];
// 						IF c = BYTE#34 OR c = BYTE#39 THEN (* " Double Quote OR ' Single Quote *)
// 							EXIT;
// 						END_IF;
// 						index := index + 1;
// 					END_WHILE;
// 				END_IF;
// 				index := index + 1;
// 			END_WHILE;
// 			CTRL.BLOCK1_STOP := INT_TO_UINT(index - 2);
// 		END_IF;
//
// 		CTRL.VALUE := BUFFER_TO_STRING(PT:=REF(BUF),SIZE:=INT_TO_UINT(stop + 1),START:=CTRL.BLOCK1_START,STOP:=CTRL.BLOCK1_STOP);
// 		CTRL.COUNT := CTRL.COUNT + UINT#1;
// 		mode := 100; (* -> NEXT-JOB: ELEMENT *)
//
// 		IF BIT_OF_DWORD(WORD_TO_DWORD(command),CTRL.TYP) THEN
//       			(*####################  DEBUG-MESSAGE  ###################################*)
//       			(*IF _debug_enable THEN
//         				LOG_CL.NEW_MSG := 'DTD: #T:~1 #C:~2 #S:~3 #E:~4 #V:~5 #P:~6 #L:~7';
//         				LOG_CL.PRINTF[1] := INT_TO_STRING(CTRL.TYP);
//         				LOG_CL.PRINTF[2] := UINT_TO_STRING(CTRL.COUNT);
//         				LOG_CL.PRINTF[3] := UINT_TO_STRING(CTRL.BLOCK1_START);
//         				LOG_CL.PRINTF[4] := UINT_TO_STRING(CTRL.BLOCK1_STOP);
//         				LOG_CL.PRINTF[5] := CTRL.VALUE;
//         				LOG_CL.PRINTF[6] := CTRL.PATH;
//         				LOG_CL.PRINTF[7] := UINT_TO_STRING(CTRL.LEVEL);
//         				LOG_MSG();
//       			END_IF;*)
//       			(*########################################################################*)
// 			RETURN;
// 		END_IF;
//
// 	END_CASE;
// END_WHILE;
// END_FUNCTION_BLOCK

// Commented out for beedance: it uses FLOAT_TO_REAL, which is commented out.
// FUNCTION_BLOCK YAHOO_WEATHER
// VAR_IN_OUT
// 	IP_C : IP_C; (* IP_Control management structure *)
// 	S_BUF : NETWORK_BUFFER;
// 	R_BUF : NETWORK_BUFFER;
// 	YW : YAHOO_WEATHER_DATA;
// END_VAR
// VAR_INPUT
// 	ACTIVATE : BOOL;
// 	UNITS : BOOL;
// 	LOCATION : STRING(20);
// END_VAR
// VAR_OUTPUT
// 	BUSY : BOOL;
// 	DONE : BOOL;
// 	ERROR_C : DWORD;
// 	ERROR_T : BYTE;
// END_VAR
// VAR
// 	CTRL : XML_CONTROL;
// 	XML_READER : XML_READER;
// 	URL_DATA : url;
// 	DNS_CLIENT : DNS_CLIENT;
// 	HTTP_GET : HTTP_GET;
// 	activate_last : BOOL;
// 	value_int : INT;
// 	value_real : REAL;
// 	v_real : REAL;
// 	state : INT;
// 	cycle : INT;
// 	day_idx :	INT;
// 	idx : INT;
// 	cnt :	INT;
//
// 	(*
// 	_debug_state_last : INT := -1;
// 	_debug_enable : BOOL := TRUE;
// 	LOG_MSG : LOG_MSG;
// 	*)
// END_VAR
//
// CASE state OF
//
// 00:	IF ACTIVATE AND NOT activate_last THEN
// 		state := 40;
// 		DONE := FALSE;
// 		BUSY := TRUE;
// 		ERROR_C := DWORD#0;
// 		ERROR_T := BYTE#0;
// 		cycle := 0;
//   		URL_DATA.DOMAIN := 'query.yahooapis.com'; (* domain name for DNS_CLIENT *)
// 	END_IF;
//
// 40:	IF DNS_CLIENT.DONE THEN
// 		state := 50;
// 	ELSIF (DNS_CLIENT.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := DNS_CLIENT.ERROR;
// 	ERROR_T := BYTE#01;
// 		state   := 100;
// 	END_IF;
//
// 50:   	(* Teil 1: http://query.yahooapis.com/v1/public/yql?q=select%20units,wind,atmosphere,astronomy,location%20from%20weather.forecast%20where%20woeid=551801%20AND%20u=%27c%27
//              part 2: http://query.yahooapis.com/v1/public/yql?q=select%20item%20from%20weather.forecast%20where%20woeid=551801%20AND%20u=%27c%27 *)
//
// 	IF cycle = 0 THEN
//   		URL_DATA.QUERY := 'q=select%20units,wind,atmosphere,astronomy,location%20from%20weather.forecast%20where%20woeid='; (* part 1 - infos *)
//   	ELSE
//   		URL_DATA.QUERY := 'q=select%20item%20from%20weather.forecast%20where%20woeid='; (* part 2 - daily data *)
//   	END_IF;
//
//   	URL_DATA.QUERY := CONCAT(URL_DATA.QUERY,LOCATION); (* location *)
//
//   	IF UNITS THEN
//     		URL_DATA.QUERY := CONCAT(URL_DATA.QUERY,'%20and%20u=%27f%27'); (* fahrenheit *)
//   	ELSE
//     		URL_DATA.QUERY := CONCAT(URL_DATA.QUERY,'%20and%20u=%27c%27'); (* celsius *)
//   	END_IF;
//
//   	URL_DATA.PATH := '/v1/public/yql';
//
//   	state := 60;
//
// 60:	IF HTTP_GET.DONE THEN
// 		state := 80;
// 		CTRL.START_POS := HTTP_GET.BODY_START;
// 		CTRL.STOP_POS := HTTP_GET.BODY_STOP;
// 		CTRL.COMMAND := WORD#2#10000000_00011000; (* ONLY TEXT AND ATTRIBUTE *)
// 		CTRL.WATCHDOG  := T#1ms;
//
// 	ELSIF (HTTP_GET.ERROR > DWORD#00) THEN
// 		(* error handling *)
// 		ERROR_C := HTTP_GET.ERROR;
// 		ERROR_T := BYTE#02;
// 		state := 100;
// 	END_IF;
//
// 80: XML_READER(CTRL:=CTRL,BUF:=R_BUF.BUFFER); (* read XML data serially *)
// 	IF CTRL.TYP < 98 THEN (* evaluate only if not a watchdog pass *)
// 		value_int := 0;
// 		value_real := 0.0;
// 		(* perform automatic conversion to real and int *)
// 		IF LEN(CTRL.VALUE) <= 20 THEN
// 			v_real:=FLOAT_TO_REAL(FLT:=CTRL.VALUE);
// 			IF CHK_REAL(v_real) = 0 THEN (* !!! checks for a valid REAL value, otherwise TwinCAT 2.x can crash !!!!  *)
// 				value_real := v_real;
// 				value_int := REAL_TO_INT(value_real);
// 			END_IF;
// 		END_IF;
//
// 		(*####################  DEBUG-MESSAGE  ###################################*)
// 		(*IF _debug_enable THEN
// 			LOG_CL.NEW_MSG := 'XML: ~2 # ~3 # ~4 # ~5';
// 			LOG_CL.PRINTF[2] := UINT_TO_STRING(CTRL.LEVEL);
// 			LOG_CL.PRINTF[3] := UINT_TO_STRING(CTRL.COUNT);
// 			LOG_CL.PRINTF[4] := CTRL.ELEMENT;
// 			LOG_CL.PRINTF[5] := CTRL.VALUE;
// 			LOG_MSG();
// 		END_IF;*)
// 		(*########################################################################*)
//
// 		IF cycle = 0 THEN
// 	      		(* ------------ part 1 --------------*)
// 	      		CASE UINT_TO_INT(CTRL.COUNT) OF
// 	        		014: YW.location_city := CTRL.VALUE;
// 	        		015: YW.location_country := CTRL.VALUE;
// 	        		016: YW.location_region := CTRL.VALUE;
// 	        		020: YW.wind_chill := value_int;
// 	        		021: YW.wind_direction := value_int;
// 	        		022: YW.wind_speed := value_real;
// 	        		026: YW.atmosphere_humidity := value_int;
// 	        		027: YW.atmosphere_pressure := value_real;
// 	        		028: YW.atmosphere_rising := value_int;
// 	        		029: YW.atmosphere_visibility := value_real;
// 	        		033: YW.astronomy_sunrise := CTRL.VALUE;
// 	        		034: YW.astronomy_sunset := CTRL.VALUE;
// 	        		038: YW.unit_distance := CTRL.VALUE;
// 	        		039: YW.unit_pressure := CTRL.VALUE;
// 	        		040: YW.unit_speed := CTRL.VALUE;
// 	        		041: YW.unit_temperature := CTRL.VALUE;
// 	      		END_CASE;
// 	    	ELSE
// 	      		(* ------------ part 2 --------------*)
// 	      		CASE UINT_TO_INT(CTRL.COUNT) OF
// 	        		018: YW.GEO_LATITUDE := value_real;
// 	        		022: YW.GEO_LONGITUDE := value_real;
//
// 			032: YW.CUR_CONDITIONS_CODE := value_int;
// 			         YW.CUR_CONDITIONS_ICON := value_int;
// 	        		034: YW.CUR_CONDITIONS_TEMP := value_int;
// 	        		035: YW.CUR_CONDITIONS_TEXT := CTRL.VALUE;
//
// 			039: YW.FORECAST_TODAY_CODE := value_int;
// 			         YW.FORECAST_TODAY_ICON := value_int;
// 			040: YW.FORECAST_TODAY_DATE_LONG := CTRL.VALUE;
// 	        		042: YW.FORECAST_TODAY_HIGH_TEMP := value_int;
// 	        		043: YW.FORECAST_TODAY_LOW_TEMP := value_int;
// 	        		044: YW.FORECAST_TODAY_TEXT := CTRL.VALUE;
//
// 			048..125: (* forcast day 1-10 *)
// 				cnt := UINT_TO_INT(CTRL.COUNT) - 48;
// 				day_idx := cnt / 9 + 1;
// 				idx := cnt MOD 9;
// 	      	  		CASE idx OF
// 					0:	YW.FORECAST_DAY[day_idx].CODE := value_int;
// 						YW.FORECAST_DAY[day_idx].ICON := value_int;
// 					1:	YW.FORECAST_DAY[day_idx].DATE_LONG := CTRL.VALUE;
// 					3:	YW.FORECAST_DAY[day_idx].HIGH_TEMP := value_int;
// 					4:	YW.FORECAST_DAY[day_idx].LOW_TEMP := value_int;
// 					5:	YW.FORECAST_DAY[day_idx].TEXT := CTRL.VALUE;
// 	      	  		END_CASE;
// 	      		END_CASE;
// 	    	END_IF;
// 	ELSIF CTRL.TYP = 99 THEN (* EXIT - last element read *)
// 		state := 100;
// 	END_IF;
//
// 100:	(* UNLOCK HTTP DATA *)
// 	IF (NOT HTTP_GET.DONE) THEN
// 		IF cycle = 0 AND ERROR_T = BYTE#0 THEN
// 			(* run part 2 *)
//     			state := 50;
// 			cycle := cycle + 1;
// 		ELSE
// 	    		state := 0;
// 	    		BUSY  := FALSE;
// 	    		DONE := ERROR_T = BYTE#0;
// 		END_IF;
// 	END_IF;
//
// END_CASE;
//
// (* ---------------- DNS --------------- *)
// DNS_CLIENT(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			DOMAIN:=URL_DATA.DOMAIN,
// 			IP4_DNS:=DWORD#00,
// 			ACTIVATE:=state=40 (* DNS Start *)
// 			);
// (* parameters not used, or not used directly *)
// (*	DWORD :=DNS_CLIENT.IP4;
// 	BOOL  :=DNS_CLIENT.DONE;
// 	BYTE  :=DNS_CLIENT.ERROR; *)
// (* -------------------------------------*)
//
// (* ------------- HTTP_GET --------------- *)
// HTTP_GET(	IP_C:=IP_C,
// 			S_BUF:=S_BUF,
// 			R_BUF:=R_BUF,
// 			IP4:=DNS_CLIENT.IP4,
// 			GET:=state=60,
// 			MODE:=BYTE#2, (* HTTP 1.1 with persistent connection *)
// 			UNLOCK_BUF:=state=100,
// 			URL_DATA:=URL_DATA
// 			);
//
// (* parameters not used, or not used directly *)
// (*	STRING := HTTP_GET.HTTP_STATUS;
// 	UINT :=HTTP_GET.HEADER_START;
// 	UINT :=HTTP_GET.HEADER_STOP;
// 	UINT :=HTTP_GET.BODY_START;
// 	UINT :=HTTP_GET.BODY_STOP;
// 	BOOL :=HTTP_GET.DONE;
// 	BYTE :=HTTP_GET.ERROR;
// *)
// (* -------------------------------------*)
//
// activate_last := ACTIVATE;
//
// (*####################  DEBUG-MESSAGE  ###################################*)
// (*IF _debug_enable THEN
// 	IF state <> _debug_state_last THEN
// 		LOG_CL.NEW_MSG := 'yw-step: /last ~1 /New ~2';
// 		LOG_CL.PRINTF[1] := INT_TO_STRING(_debug_state_last);
// 		LOG_CL.PRINTF[2] := INT_TO_STRING(state);
// 		LOG_MSG();
// 		_debug_state_last := state;
// 	END_IF;
// END_IF;*)
// (*########################################################################*)
// END_FUNCTION_BLOCK

FUNCTION_BLOCK YAHOO_WEATHER_DESC_GE
VAR_IN_OUT
	YW : YAHOO_WEATHER_DATA;
END_VAR
VAR_INPUT
	ACTIVATE : BOOL;
END_VAR
VAR
	act_last : BOOL;
	x : INT;
	y : INT;
	str : STRING(40);
END_VAR

IF ACTIVATE AND act_last = FALSE THEN
	FOR y := 1 TO 12 DO
		CASE y OF
		01:	x := YW.CUR_CONDITIONS_CODE;

		02:	YW.CUR_CONDITIONS_TEXT := str;
			x := YW.FORECAST_TODAY_CODE;
		03:	YW.FORECAST_TODAY_TEXT := str;
			x := YW.FORECAST_DAY[1].CODE;
		04:	YW.FORECAST_DAY[1].TEXT := str;
			x := YW.FORECAST_DAY[2].CODE;
		05:	YW.FORECAST_DAY[2].TEXT := str;
			x := YW.FORECAST_DAY[3].CODE;
		06:	YW.FORECAST_DAY[3].TEXT := str;
			x := YW.FORECAST_DAY[4].CODE;
		07:	YW.FORECAST_DAY[4].TEXT := str;
			x := YW.FORECAST_DAY[5].CODE;
		08:	YW.FORECAST_DAY[5].TEXT := str;
			x := YW.FORECAST_DAY[6].CODE;
		09:	YW.FORECAST_DAY[6].TEXT := str;
			x := YW.FORECAST_DAY[7].CODE;
		10:	YW.FORECAST_DAY[7].TEXT := str;
			x := YW.FORECAST_DAY[8].CODE;
		11:	YW.FORECAST_DAY[8].TEXT := str;
			x := YW.FORECAST_DAY[9].CODE;
		12:	YW.FORECAST_DAY[9].TEXT := str;
		END_CASE;

		CASE x OF
		00:	str := 'Wirbelsturm';
		01:	str := 'Tropensturm';
		02:	str := 'Wirbelsturm';
		03:	str := 'starkes Gewitter ';
		04:	str := 'Gewitter';
		05:	str := 'Schneeregen';
		06:	str := 'Schneeregen';
		07:	str := 'Schneeregen';
		08:	str := 'gefrierender Nieselregen';
		09:	str := 'Nieselregen';
		10:	str := 'gefrierender Regen';
		11:	str := 'Regenfälle';
		12:	str := 'Regenfälle';
		13:	str := 'Schneegestöber';
		14:	str := 'leichte Schneefälle';
		15:	str := 'Schneetreiben';
		16:	str := 'Schnee';
		17:	str := 'Hagel';
		18:	str := 'Schneeregen';
		19:	str := 'Staub';
		20:	str := 'nebelig';
		21:	str := 'Dunst';
		22:	str := 'rauchig';
		23:	str := 'stürmisch';
		24:	str := 'windig';
		25:	str := 'kalt';
		26:	str := 'bewölkt';
		27:	str := 'meist bewölkt';
		28:	str := 'meist bewölkt';
		29:	str := 'teilweise bewölkt';
		30:	str := 'teilweise bewölkt';
		31:	str := 'klar';
		32:	str := 'sonnig';
		33:	str := 'heiter';
		34:	str := 'heiter';
		35:	str := 'Regen und Hagel';
		36:	str := 'heiß';
		37:	str := 'örtliche Gewitter';
		38:	str := 'vereinzelte Gewitter';
		39:	str := 'vereinzelte Gewitter';
		40:	str := 'vereinzelte Regenfälle';
		41:	str := 'schwere Schneefälle';
		42:	str := 'vereinzelte Schneefälle';
		43:	str := 'schwere Schneefälle';
		44:	str := 'teilweise bewölkt';
		45:	str := 'gewittrige Regenfälle';
		46:	str := 'Schneefälle';
		47:	str := 'örtliche gewittrige Regenfälle';
		ELSE
			str := 'nicht verfügbar';
		END_CASE;
	END_FOR;
END_IF;
act_last := ACTIVATE;
END_FUNCTION_BLOCK

FUNCTION_BLOCK YAHOO_WEATHER_ICON_OSCAT
VAR_IN_OUT
	YW : YAHOO_WEATHER_DATA;
END_VAR
VAR_INPUT
	ACTIVATE : BOOL;
END_VAR
VAR
	dat : ARRAY[0..48] OF INT := 4,2,2,4,17,1,3,1,2,2,2,2,2,3,3,3,3,5,3,6,7,14,8,9,9,10,9,11,12,11,12,13,14,15,16,4,14,4,17,2,2,3,3,3,12,2,3,4,0;
	act_last : BOOL;
	x : INT;
	idx : INT;
END_VAR

IF ACTIVATE AND act_last = FALSE THEN
	x := MIN(YW.CUR_CONDITIONS_CODE,48);
	YW.CUR_CONDITIONS_ICON := dat[x];
	x := MIN(YW.FORECAST_TODAY_CODE,48);
	YW.FORECAST_TODAY_ICON := dat[x];
	FOR idx := 1 TO 9 DO
		x := MIN(YW.FORECAST_DAY[idx].CODE,48);
		YW.FORECAST_DAY[idx].ICON := dat[x];
	END_FOR;
END_IF;
act_last := ACTIVATE;
END_FUNCTION_BLOCK

TYPE DLOG_CRON_DATA :
STRUCT
	CE : ARRAY [0..5] OF DLOG_CRON_ELEMENT;
END_STRUCT
END_TYPE

TYPE DLOG_CRON_ELEMENT :
STRUCT
	ELEMENTS : ARRAY [0..59] OF BOOL;
	VALUE : INT;
	VALUE_MIN : INT;
	VALUE_MAX : INT;
	ALL_SELECTED : BOOL;
END_STRUCT
END_TYPE

TYPE DLOG_DATA :
STRUCT
	STORE_TYPE : BYTE;
	ADD_COM : INT;
	ADD_DATA_REQ : BOOL;
	CLOCK_TRIG : BOOL;
	ID_MAX : USINT;
	DTI : DT;
	UCB : UNI_CIRCULAR_BUFFER_DATA;
	LOAD_TIME_MAX : TIME;
	NEW_FILE : STRING;
	NEW_FILE_SIZE : UDINT;
	NEW_FILE_RTRIG : BOOL;
END_STRUCT
END_TYPE

TYPE
	DLOG_REAL_ARRAY_DATA : ARRAY [1..64] OF DLOG_REAL_ARRAY_ELEMENT;
END_TYPE

TYPE DLOG_REAL_ARRAY_ELEMENT :
STRUCT
	VALUE : REAL;
	DELTA : REAL;
	COLUMN : STRING(10);
	_DELTA_LAST : 	REAL;
	_VALUE_LAST : 	REAL;
END_STRUCT
END_TYPE

TYPE DLOG_RETAIN :
STRUCT
	FN_REM : STRING;
	COLOR : BOOL;
	HEAD: BOOL;
	TRIG_CNT : UDINT;
	TRIG_CNT_TOTAL : UDINT;
END_STRUCT
END_TYPE

TYPE DLOG_SAVE :
STRUCT
	FN_REM : STRING;
	COLOR : BOOL;
	HEAD: BOOL;
	TRIG_CNT : UDINT;
	TRIG_CNT_TOTAL : UDINT;
END_STRUCT
END_TYPE

TYPE FILE_PATH_DATA :
STRUCT
	DRIVE : STRING(3);
	DIRECTORY : STRING(STRING_LENGTH);
	FILENAME : STRING;
END_STRUCT
END_TYPE

TYPE FILE_SERVER_DATA :
STRUCT
	FILE_OPEN : BOOL;
	FILENAME : STRING;
	MODE : BYTE;
	OFFSET : UDINT;
	FILE_SIZE : UDINT;
	ERROR : BYTE;
	RUNTIME : FILE_SERVER_RUNTIME_DATA;
END_STRUCT
END_TYPE

TYPE FILE_SERVER_RUNTIME_DATA :
STRUCT
	TIME_FILE_OPEN_CUR : UDINT;
	TIME_FILE_OPEN_MAX : UDINT;
	TIME_FILE_OPEN_MIN : UDINT;
	TIME_FILE_CLOSE_CUR : UDINT;
	TIME_FILE_CLOSE_MAX : UDINT;
	TIME_FILE_CLOSE_MIN : UDINT;
	TIME_FILE_READ_CUR: UDINT;
	TIME_FILE_READ_MAX : UDINT;
	TIME_FILE_READ_MIN : UDINT;
	TIME_FILE_WRITE_CUR : UDINT;
	TIME_FILE_WRITE_MAX : UDINT;
	TIME_FILE_WRITE_MIN : UDINT;
	TIME_FILE_SEEK_CUR : UDINT;
	TIME_FILE_SEEK_MAX : UDINT;
	TIME_FILE_SEEK_MIN : UDINT;
	TIME_FILE_TELL_CUR : UDINT;
	TIME_FILE_TELL_MAX : UDINT;
	TIME_FILE_TELL_MIN : UDINT;
	TIME_FILE_DELETE_CUR : UDINT;
	TIME_FILE_DELETE_MAX : UDINT;
	TIME_FILE_DELETE_MIN : UDINT;
END_STRUCT
END_TYPE

TYPE IP2GEO_DATA :
STRUCT
	STATE : BOOL;
	IP4 : DWORD;
	COUNTRY_CODE : STRING(2);
	COUNTRY_NAME : STRING(20);
	REGION_CODE : STRING(2);
	REGION_NAME : STRING(20);
	CITY : STRING(20);
	GEO_LATITUDE : REAL;
	GEO_LONGITUDE : REAL;
	TIME_ZONE_NAME : STRING(20);
	GMT_OFFSET : INT;
	IS_DST : BOOL;
END_STRUCT
END_TYPE

TYPE IP_C :
STRUCT
	C_MODE : BYTE;						(*W Mode: TCP/UDP ACTIVE/PASSIVE                               *)
	C_PORT : WORD;						(*W port number                                               *)
	C_IP : DWORD;						(*W IP address packed wwxxyyzz = www.xxx.yyy.zzz            *)
	C_STATE : BYTE;						(*R connection status ON/OFF + edge ON/OFF                 *)
	C_ENABLE : BOOL;					(*W enable for connect                                     *)
	R_OBSERVE : BOOL;					(*W monitor receive                                      *)
	TIME_RESET : BOOL;					(*W reset all timers                                   *)
	ERROR : DWORD;					(*R vvwwxxyy (vv = CON_ERROR, ww = SEN_ERROR, xx = REC_ERROR *)
	FIFO : IP_FIFO_DATA;					(*I IP FIFO structure                                         *)
	MAILBOX : ARRAY [1..16] OF BYTE;	(*I mailbox: data area for data exchange between blocks         *)
END_STRUCT
END_TYPE

TYPE IP_FIFO_DATA :
STRUCT
	X : ARRAY [1..128] OF BYTE; 	(* IP_ID FIFO memory          *)
	Y : ARRAY [1..128] OF BYTE; 	(* IP_ID counter                 *)
	ID : BYTE;					(* highest assigned id         *)
	MAX_ID : BYTE;				(* maximum registrations per ID   *)
	INIT : BOOL;					(* initialization done *)
	EMPTY : BOOL;				(* FIFO empty                    *)
	FULL : BOOL;					(* FIFO full                    *)
	TOP : INT;					(* maximum number in FIFO      *)
	NW : INT;					(* write pointer                *)
	NR : INT;						(* read pointer                   *)
END_STRUCT
END_TYPE

TYPE LOG_CONTROL :
STRUCT
	NEW_MSG : STRING(LOG_SIZE); (* new message for the data buffer *)
	NEW_MSG_OPTION : DWORD;
	LEVEL : BYTE;
	SIZE : INT := LOG_MAX; (* number of strings in LOG_DATA *)
	RESET : BOOL; (* clear message buffer *)
	PRINTF : PRINTF_DATA; (* Array PRINTF Parameter *)
	MSG : ARRAY [0..LOG_MAX] OF STRING(LOG_SIZE); (* Array Message *)
	MSG_OPTION : ARRAY [0..LOG_MAX] OF DWORD; (* Array Options *)
	UPDATE_COUNT : UINT; (* update counter *)
	IDX : INT; (* current message index *)
	RING_MODE : BOOL; (* buffer overflow detected *)
END_STRUCT
END_TYPE

TYPE MYSQL_COM :
STRUCT
	S_BUF : NETWORK_BUFFER;
	R_BUF : NETWORK_BUFFER;
	SQL_CON : BOOL;
	SQL_URL : STRING(STRING_LENGTH);
	SQL_PACKET_NO : BYTE;
	TIMEOUT : TIME;
	DNS_IP4 : DWORD;
	SQL_RCV_STATE : BYTE;
	ERROR_C : DWORD;
	ERROR_T : BYTE;
END_STRUCT
END_TYPE

TYPE MYSQL_INFO :
STRUCT
	SERVER_PROTOCOL_VERSION : BYTE;
	SERVER_CAPABILITIES : WORD;
	SERVER_STATUS : WORD;
	SERVER_LANGUAGE : BYTE;
	SQL_CONNECTED : BOOL;
	SQL_ERROR : STRING(STRING_LENGTH);
	DATA_INSERT_AKTIV : BOOL;
	DATA_INSERT_OK_CNT : UDINT;
	DATA_INSERT_NOK_CNT : UDINT;
END_STRUCT
END_TYPE

TYPE NET_VAR_DATA :
STRUCT
	CYCLE : UDINT;
	STATE : BYTE;
	INDEX : INT;
	ID_MAX : USINT;
	ERROR_ID : BYTE;
	BUF_SIZE : UINT;
	S_BUF : NETWORK_BUFFER;
	R_BUF : NETWORK_BUFFER;
END_STRUCT
END_TYPE

TYPE NETWORK_BUFFER :
STRUCT
	SIZE : UINT;
	BUFFER : NW_BUF_LONG;
END_STRUCT
END_TYPE

TYPE NETWORK_BUFFER_SHORT :
STRUCT
	SIZE : UINT;
	BUFFER : NW_BUF_SHORT;
END_STRUCT
END_TYPE

TYPE
	NW_BUF_LONG  : ARRAY[0..NETWORK_BUFFER_LONG_SIZE] OF BYTE;
END_TYPE

TYPE
	NW_BUF_SHORT : ARRAY[0..NETWORK_BUFFER_SHORT_SIZE] OF BYTE;
END_TYPE

TYPE PRINTF_DATA : ARRAY [1..11]  OF STRING(LOG_SIZE);
END_TYPE

TYPE
	UNI_CIRCULAR_BUF : ARRAY [0..8191] OF BYTE;
END_TYPE

TYPE UNI_CIRCULAR_BUFFER_DATA :
STRUCT
	D_MODE : INT;
	D_SIZE : INT;
	D_HEAD : WORD;
	D_STRING : STRING(STRING_LENGTH);
	D_REAL : REAL;
	D_DWORD : DWORD;
	BUF_SIZE : UINT;
	BUF_COUNT : UINT;
	BUF_USED : USINT;
	BUF_USED_MAX : USINT;
	BUF_DATA_CNT : UDINT;
	BUF_DATA_LOST : UDINT;
	BUF : UNI_CIRCULAR_BUF;
	_GetStart : UINT;
	_GetEnd : UINT;
	_Last : UINT;
	_First : UINT;
END_STRUCT
END_TYPE

TYPE URL :
STRUCT
	PROTOCOL : STRING(10);
	USER : STRING(32);
	PASSWORD : STRING(32);
	DOMAIN : STRING(80);
	PORT : WORD;
	PATH : STRING(80);
	QUERY : STRING(STRING_LENGTH);
	ANCHOR : STRING(40);
	HEADER : STRING(160);
END_STRUCT
END_TYPE

TYPE us_LOG_VIEWPORT :
STRUCT
	LINE_ARRAY : ARRAY [1..40] OF INT;
	COUNT : INT; (* number of visible messages *)
	UPDATE_COUNT : UINT; (* update counter *)
	MOVE_TO_X : INT; (* control of the message display *)
	UPDATE : BOOL; (* data was changed -> redraw *)
END_STRUCT
END_TYPE

TYPE us_TN_INPUT_CONTROL :
STRUCT
	bo_Enable: BOOL; (* enable/disable processing *)
	bo_Update_all: BOOL; (* redraw all elements *)
	bo_Reset_Fokus: BOOL; (* set focus on 1st element *)
	in_Focus_at: INT; (* control element with active focus *)
	in_Count: INT; (* number of used control elements *)
	in_ToolTip_X: INT; (* Tooltip Text X-Offset *)
	in_ToolTip_Y: INT; (* Tooltip Text Y-Offset *)
	by_ToolTip_Attr: BYTE; (* Tooltip Text Attribute *)
	in_ToolTip_Size: INT; (* tooltip text length *)
	usa_TN_INPUT_CONTROL_DATA: ARRAY [1..20] OF us_TN_INPUT_CONTROL_DATA;
END_STRUCT
END_TYPE

TYPE us_TN_INPUT_CONTROL_DATA :
STRUCT
 	by_Input_Exten_Code: BYTE;
	by_Input_ASCII_Code: BYTE;
	bo_Input_ASCII_IsNum: BOOL;
	in_Title_X_Offset: INT; (* Title Text X-Offset *)
	in_Title_Y_Offset: INT; (* Title Text Y-Offset *)
	by_Title_Attr: BYTE; (* Title Text Attribute *)
	st_Title_String: STRING; (* Title Text String *)
	in_Cursor_X: INT; (* current cursor X position *)
	in_Cursor_Y: INT; (* current cursor Y position *)
	in_Type: INT; (* type of the element *)
	in_X: INT; (* Element X-Position *)
	in_Y: INT; (* Element Y-Position *)
	in_Cursor_Pos: INT; (* current cursor position *)
	by_Attr_mF: BYTE; (* text attributes with focus *)
	by_Attr_oF: BYTE; (* text attributes without focus *)
	in_selected: INT; (* text item was selected *)
	st_Input_Mask: STRING; (* text input mask *)
	st_Input_Data: STRING(STRING_LENGTH); (* text current input data *)
	st_Input_String: STRING; (* text copy after return key *)
	st_Input_ToolTip: STRING; (* text tooltip (input information) *)
	in_Input_Option: INT; (* text output options *)
	bo_Input_Entered: BOOL; (* RETURN key was pressed *)
	bo_Input_Hidden: BOOL; (* text hidden input with '*' *)
	bo_Input_Only_Num: BOOL; (* text allow numeric input only *)
	bo_Focus: BOOL; (* element has the focus *)
	bo_Update_Input: BOOL; (* redraw element because of user input *)
	bo_Update_All: BOOL; (* redraw element completely *)

END_STRUCT
END_TYPE

TYPE us_TN_MENU :
STRUCT
	st_Menu_Text: STRING(STRING_LENGTH); (* X menu texts *)
	in_Menu_E_Count: INT; (* menu element count *)
	in_Y: INT; (* menu line number  Y position *)
	in_X: INT; (* menu column number X position *)
	by_Attr_mF: BYTE; (* text attributes with focus *)
	by_Attr_oF: BYTE; (* text attributes without focus *)
	in_X_SM_new: INT; (* new submenu X position *)
	in_Y_SM_new: INT; (* new submenu Y position *)
	in_X_SM_old: INT; (* old submenu X position *)
	in_Y_SM_old: INT; (* old submenu Y position *)
	in_Cur_Menu_Item: INT; (* current main menu *)
	in_Cur_Sub_Item: INT; (* current submenu *)
	in_State: INT; (*   Menu status   *)
	in_Menu_Selected: INT; (* X menu item selected with enter *)
	bo_Create: BOOL; (* X create menu *)
	bo_Destroy: BOOL; (* X remove menu *)
	bo_Update: BOOL; (* X update menu *)
END_STRUCT
END_TYPE

TYPE us_TN_MENU_POPUP :
STRUCT
	st_Menu_Text: STRING(STRING_LENGTH); (* X menu texts *)
	in_Menu_E_Count: INT; (* menu element count *)
	in_X: INT; (* X menu X position *)
	in_Y: INT; (* X menu Y position *)
	in_Cols: INT; (* menu column count *)
	in_Rows: INT; (* menu line count *)
	in_Cur_Item: INT; (* X current submenu *)
	by_Attr_mF: BYTE; (* X text attributes with focus	*)
	by_Attr_oF: BYTE; (* X text attributes without focus *)
	by_Input_Exten_Code: BYTE; (* X key code: special keys *)
	bo_Create: BOOL; (* X create popup menu *)
	bo_Destroy: BOOL; (* X remove popup menu *)
	bo_Update: BOOL; (* X update popup menu *)
	bo_Activ: BOOL; (*   popup menu is active *)
END_STRUCT
END_TYPE

TYPE us_TN_SCREEN :
STRUCT
	bya_CHAR: ARRAY [0..1919] OF BYTE;
	bya_COLOR: ARRAY [0..1919] OF BYTE;
	bya_BACKUP: ARRAY [0..1919] OF BYTE;
	bya_Line_Update: ARRAY [0..23] OF BOOL;
	by_Input_Exten_Code: BYTE;
	by_Input_ASCII_Code: BYTE;
	bo_Input_ASCII_IsNum: BOOL;
	in_Page_Number: INT; (* active page no. *)
	in_Cursor_X: INT; (* Cursor X-Position *)
	in_Cursor_Y: INT; (* Cursor Y-Position *)
	in_EOS_Offset: INT; (* End of String Offset *)
	by_Clear_Screen_Attr: BYTE; (* Clear-Screen Attr *)
	bo_Clear_Screen: BOOL; (* trigger clear screen *)
	bo_Modal_Dialog: BOOL; (* modal dialog is active *)
	bo_Menue_Bar_Dialog: BOOL; (* menu dialog is active *)
END_STRUCT
END_TYPE

TYPE VMAP_DATA :
STRUCT
	FC : DWORD;
	V_ADR : INT;
	V_SIZE : INT;
	P_ADR : INT;
	TIME_OUT : TIME;
END_STRUCT
END_TYPE

TYPE WORLD_WEATHER_CUR :
STRUCT
	OBSERVATION_TIME : STRING(8);
	TEMP_C : INT;
	WEATHER_CODE : INT;
	WEATHER_DESC : STRING(60);
	WEATHER_ICON : INT;
	WIND_SPEED_MILES : INT;
	WIND_SPEED_KMPH : INT;
	WIND_DIR_DEGREE : INT;
	WIND_DIR16POINT : STRING(3);
	PRECIPMM : REAL;
	HUMIDITY : INT;
	VISIBILITY : INT;
	PRESSURE : INT;
	CLOUDOVER : INT;
END_STRUCT
END_TYPE

TYPE WORLD_WEATHER_DATA :
STRUCT
	CUR : WORLD_WEATHER_CUR;
	DAY : ARRAY [0..4] OF WORLD_WEATHER_DAY;
END_STRUCT
END_TYPE

TYPE WORLD_WEATHER_DAY :
STRUCT
	DATE_OF_DAY : STRING(10);
	TEMP_MAX_C : INT;
	TEMP_MAX_F : INT;
	TEMP_MIN_C : INT;
	TEMP_MIN_F : INT;
	WIND_SPEED_MILES : INT;
	WIND_SPEED_KMPH : INT;
	WIND_DIR_DEGREE : INT;
	WIND_DIR16POINT : STRING(3);
	WEATHER_CODE : INT;
	WEATHER_DESC : STRING(60);
	WEATHER_ICON : INT;
	PRECIPMM : REAL;
END_STRUCT
END_TYPE

TYPE XML_CONTROL :
STRUCT
	COMMAND : WORD;
	WATCHDOG : TIME;
	START_POS : UINT;
	STOP_POS : UINT;
	COUNT : UINT;
	TYP : INT;
	LEVEL : UINT;
	PATH : STRING(STRING_LENGTH);
	ELEMENT : STRING(STRING_LENGTH);
	ATTRIBUTE : STRING(STRING_LENGTH);
	VALUE : STRING(STRING_LENGTH);
	BLOCK1_START : UINT;
	BLOCK1_STOP : UINT;
	BLOCK2_START : UINT;
	BLOCK2_STOP : UINT;
END_STRUCT
END_TYPE

TYPE YAHOO_WEATHER_DATA :
STRUCT
	LOCATION_CITY : STRING(40);   					(* THE LOCATION OF THIS FORECAST: CITY: CITY NAME  *)
	LOCATION_REGION : STRING(20);   					(* THE LOCATION OF THIS FORECAST: REGION: STATE, TERRITORY, OR REGION, IF GIVEN *)
	LOCATION_COUNTRY : STRING(20);   				(* THE LOCATION OF THIS FORECAST: COUNTRY: TWO-CHARACTER COUNTRY CODE *)
	UNIT_TEMPERATURE : STRING(1);    				(* TEMPERATURE: DEGREE UNITS, F FOR FAHRENHEIT OR C FOR CELSIUS *)
	UNIT_DISTANCE : STRING(2);    					(* DISTANCE: UNITS FOR DISTANCE, MI FOR MILES OR KM FOR KILOMETERS *)
	UNIT_PRESSURE : STRING(2);    					(* PRESSURE: UNITS OF BAROMETRIC PRESSURE, IN FOR POUNDS PER SQUARE INCH OR MB FOR MILLIBARS *)
	UNIT_SPEED : STRING(4);    						(* SPEED: UNITS OF SPEED, MPH FOR MILES PER HOUR OR KPH FOR KILOMETERS PER HOUR *)
	WIND_CHILL : INT;              							(* FORECAST INFORMATION ABOUT WIND CHILL IN DEGREES *)
	WIND_DIRECTION : INT;              						(* FORECAST INFORMATION ABOUT WIND DIRECTION, IN DEGREES *)
	WIND_SPEED : REAL;             						(* FORECAST INFORMATION ABOUT WIND SPEED, IN THE UNITS (MPH OR KPH)*)
	ATMOSPHERE_HUMIDITY : INT;              				(* FORECAST INFORMATION ABOUT CURRENT ATMOSPHERIC HUMIDITY: HUMIDITY, IN PERCENT  *)
	ATMOSPHERE_PRESSURE : REAL;             			(* FORECAST INFORMATION ABOUT CURRENT ATMOSPHERIC PRESSURE: BAROMETRIC PRESSURE, IN THE UNITS (IN OR MB) *)
	ATMOSPHERE_VISIBILITY : REAL;             			(* FORECAST INFORMATION ABOUT CURRENT ATMOSPHERIC VISIBILITY, IN THE UNITS (MI OR KM) *)
	ATMOSPHERE_RISING : INT;              					(* FORECAST INFORMATION ABOUT RISING: STATE OF THE BAROMETRIC PRESSURE: STEADY (0), RISING (1), OR FALLING (2). (INTEGER: 0, 1, 2) *)
	ASTRONOMY_SUNRISE : STRING(10);   				(* SUNRISE: TODAY'S SUNRISE TIME. THE TIME IS A STRING IN A LOCAL TIME FORMAT OF "H:MM AM/PM" *)
	ASTRONOMY_SUNSET : STRING(10);   				(* SUNSET: TODAY'S SUNSET TIME. THE TIME IS A STRING IN A LOCAL TIME FORMAT OF "H:MM AM/PM", *)
	GEO_LATITUDE : REAL;             						(* THE LATITUDE OF THE LOCATION *)
	GEO_LONGITUDE : REAL;             					(* THE LONGITUDE OF THE LOCATION *)
	CUR_CONDITIONS_TEMP : INT;              				(* THE CURRENT WEATHER CONDITIONS: TEMP: THE CURRENT TEMPERATURE, IN THE UNITS (F OR C) *)
	CUR_CONDITIONS_TEXT : STRING(40);   			(* THE CURRENT WEATHER CONDITIONS: TEXT: A TEXTUAL DESCRIPTION OF CONDITIONS *)
	CUR_CONDITIONS_CODE : INT;              				(* THE CURRENT WEATHER CONDITIONS: CODE: THE CONDITION CODE FOR THIS FORECAST *)
	CUR_CONDITIONS_ICON : INT;			  			(* THE CURRENT WEATHER CONDITIONS: ICON: THE CONDITION ICON FOR THIS FORECAST *)
	FORECAST_TODAY_LOW_TEMP : INT;              		(* THE FORECAST TODAY WEATHER CONDITIONS: THE FORECASTED LOW TEMPERATURE FOR THIS DAY IN THE UNITS (F OR C) *)
	FORECAST_TODAY_HIGH_TEMP : INT;              		(* THE FORECAST TODAY WEATHER CONDITIONS: THE FORECASTED HIGH TEMPERATURE FOR THIS DAY IN THE UNITS (F OR C) *)
	FORECAST_TODAY_TEXT : STRING(40);   			(* THE FORECAST TODAY WEATHER CONDITIONS: TEXT: A TEXTUAL DESCRIPTION OF CONDITIONS *)
	FORECAST_TODAY_CODE : INT;              				(* THE CURRENT WEATHER CONDITIONS: CODE: THE CONDITION CODE FOR THIS FORECAST *)
	FORECAST_TODAY_ICON : INT;			  		(* THE CURRENT WEATHER CONDITIONS: ICON: THE CONDITION ICON FOR THIS FORECAST *)
	FORECAST_TODAY_DATE_LONG : STRING(20);	  	(* THE CURRENT WEATHER CONDITIONS: DATE: DATE OF THE DAY *)
	FORECAST_DAY : ARRAY [1..9] OF YAHOO_WEATHER_FORECAST_DAY;	 (* THE CURRENT WEATHER DAYS ARRAY *)
END_STRUCT
END_TYPE

TYPE YAHOO_WEATHER_FORECAST_DAY :
STRUCT
	LOW_TEMP : INT;					(* THE FORECAST WEATHER CONDITIONS: THE FORECASTED LOW TEMPERATURE FOR THIS DAY IN THE UNITS (F OR C) *)
	HIGH_TEMP : INT;              			(* THE FORECAST WEATHER CONDITIONS: THE FORECASTED HIGH TEMPERATURE FOR THIS DAY IN THE UNITS (F OR C) *)
	TEXT : STRING(40);   				(* THE FORECAST WEATHER CONDITIONS: TEXT: A TEXTUAL DESCRIPTION OF CONDITIONS *)
	DATE_LONG : STRING(20);   		(* THE FORECAST WEATHER CONDITIONS: DATE: DATE OF THE DAY *)
	CODE : INT;              				(* THE FORECAST WEATHER CONDITIONS: CODE: THE CONDITION CODE FOR THIS FORECAST *)
	ICON : INT;			  			(* THE FORECAST WEATHER CONDITIONS: ICON: THE CONDITION ICON FOR THIS FORECAST *)
END_STRUCT
END_TYPE

VAR_GLOBAL
END_VAR

VAR_GLOBAL CONSTANT
	NETWORK_BUFFER_LONG_SIZE : UINT := 4095;
	NETWORK_BUFFER_SHORT_SIZE : UINT := 1407;
	LOG_MAX : INT := 40;
	LOG_SIZE : INT := 80;
	ELEMENT_LENGTH : INT := 250;
END_VAR

VAR_GLOBAL
	(* ------- only for Beckhoff - TwinCat  --------- *)
	TCP_SERVER_RESET : BYTE; (* only for internaly use - do't modify *)
	sSrvNetId : STRING(23) := ''; (* String containing the network address of the TwinCAT TCP/IP Connection Server (e.g. '1.1.1.2.7.1'). For the local computer (default) an empty string may be specified *)
	sLocalHost : STRING(15) := ''; (* The local IP address (Ipv4) of the UDP client/server socket as a string (e.g. '172.33.5.1'). An empty string may be specified for the default network adapter *)
	(* --------------------------------------------- *)

	(* ------- only for CoDeSys --------------------        *)
	(* Configuration for SysLibSocket.lib                   *)
    (* Bit 0 = 1 -> PLCWinNT                                *)
    (* Bit 1 = 1 -> PLC WAGO 750-841                        *)
    (* Bit 2 = 1 -> PLC WAGO 750-881                        *)
	SYSLIBSOCKETS_OPTION : BYTE := BYTE#2#0000_0000;
	(* ---------------------------------------------------- *)

	LOG_CL : LOG_CONTROL;
END_VAR
