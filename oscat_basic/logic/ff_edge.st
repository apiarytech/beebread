(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF edge triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK COUNT_BR
VAR_INPUT
	SET : BOOL;
	IN : BYTE;
	UP : BOOL;
	DN : BOOL;
	STEP : BYTE := 1;
	MX : BYTE := 255;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	CNT : BYTE;
END_VAR
VAR
	last_up: BOOL;
	last_dn: BOOL;
END_VAR

(*
	version 1.0	16 jan 2008
	programmer 	hugo
	tested BY	tobias

Count_BR is a byte counter with independen up and dn inputs. the counter counts from 0 to mx and continues at 0 after is has reached mx
a step input sets the counters stepping width.

*)
(* @END_DECLARATION := '0' *)
IF rst THEN
	cnt := 0;
ELSIF set THEN
	cnt := LIMIT(0,in,mx);
ELSIF up AND NOT last_up THEN
	cnt := INT_TO_BYTE(INC(cnt,step,mx));
ELSIF dn AND NOT last_dn THEN
	cnt := INT_TO_BYTE(inc(cnt,-step,mx));
END_IF;
last_up := up;
last_dn := dn;

(* revision history
hm	16. jan 2008	rev 1.0
	original version

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF edge triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK COUNT_DR
VAR_INPUT
	SET : BOOL;
	IN : DWORD;
	UP : BOOL;
	DN : BOOL;
	STEP : DWORD := 1;
	MX : DWORD := 16#FFFFFFFF;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	CNT : DWORD;
END_VAR
VAR
	last_up: BOOL;
	last_dn: BOOL;
END_VAR

(*
	version 1.1	20. jan. 2011
	programmer 	hugo
	tested BY		tobias

Count_DR is a DWORD counter with independen up and dn inputs. the counter counts from 0 to mx and continues at 0 after is has reached mx
a step input sets the counters stepping width.

*)
(* @END_DECLARATION := '0' *)
IF rst THEN
	cnt := 0;
ELSIF set THEN
	cnt := LIMIT(0,in,mx);
ELSIF up AND NOT last_up THEN
	IF STEP > MX - CNT THEN
		CNT := CNT - MX + STEP - 1;
	ELSE
		CNT := CNT + STEP;
	END_IF;
ELSIF dn AND NOT last_dn THEN
	IF STEP > CNT THEN
		CNT := CNT - STEP + MX + 1;
	ELSE
		CNT := CNT - STEP;
	END_IF;
END_IF;
last_up := up;
last_dn := dn;

(* revision history
hm	12. jun 2008	rev 1.0
	original version

hm	20. jan. 2011	rev 1.1
	corrected init value of MX to 16#FFFFFFFF

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF edge triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK FF_D2E
VAR_INPUT
	D0 : BOOL;
	D1 : BOOL;
	CLK : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q0 : BOOL;
	Q1 : BOOL;
END_VAR
VAR
	edge : BOOL;
END_VAR

(*
version 1.3	14. mar. 2009
programmer 	hugo
tested by		oscat

dual D-type flip flop with reset and rising clock trigger

*)

(* @END_DECLARATION := '0' *)
IF rst THEN
	Q0 := FALSE;
	Q1 := FALSE;
ELSIF clk AND NOT edge THEN
	Q0 := D0;
	Q1 := D1;
END_IF;
edge := CLK;

(* revision history
hm	25. dec 2007	rev 1.0
	original version

hm	27. dec 2007	rev 1.1
	changed code for better performance

hm	30. oct. 200	rev 1.2
	deleted unnecessary init with 0

hm	14. mar. 2009	rev 1.3
	removed double assignments

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF edge triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK FF_D4E
VAR_INPUT
	D0 : BOOL;
	D1 : BOOL;
	D2 : BOOL;
	D3 : BOOL;
	CLK : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q0 : BOOL;
	Q1 : BOOL;
	Q2 : BOOL;
	Q3 : BOOL;
END_VAR
VAR
	edge : BOOL;
END_VAR


(*
version 1.3	14. mar. 2009
programmer 	hugo
tested by		oscat

quad D-type flip flop with reset and rising clock trigger

*)

(* @END_DECLARATION := '0' *)
IF rst THEN
	Q0 := FALSE;
	Q1 := FALSE;
	Q2 := FALSE;
	Q3 := FALSE;
ELSIF clk AND NOT edge THEN
	Q0 := D0;
	Q1 := D1;
	Q2 := D2;
	Q3 := D3;
END_IF;
edge := CLK;

(* revision history
hm	4. aug 2006	rev 1.0
	original version

hm	27. dec 2007	rev 1.1
	changed code for better performance

hm	30. oct. 2008	rev 1.2
	deleted unnecessary init with 0

hm	14. mar. 2009	rev 1.3
	removed double assignments

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Logic\/FF edge triggered' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK FF_DRE
VAR_INPUT
	SET : BOOL;
	D : BOOL;
	CLK : BOOL;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
END_VAR
VAR
	edge : BOOL;
END_VAR


(*
version 1.2	30. oct. 2008
programmer 	hugo
tested by	oscat

D-type flip flop with set, reset and rising clock trigger

*)

(* @END_DECLARATION := '0' *)
IF rst OR set THEN
	Q := NOT rst;
ELSIF clk AND NOT edge THEN
	Q := D;
END_IF;
edge := CLK;

(* revision history
hm	4. aug 2006		rev 1.0
	original version

hm	27. dec 2007	rev 1.1
	changed code for better performance

hm	30. oct. 2008	rev 1.2
	optimized performance
*)
END_FUNCTION_BLOCK