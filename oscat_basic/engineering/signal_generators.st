(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK _RMP_B
VAR_INPUT
	DIR : BOOL;  (* true = up *)
	E : BOOL := TRUE;
	TR : TIME;
END_VAR
VAR_IN_OUT
	RMP : BYTE;
END_VAR
VAR
	tx, tl, tn : TIME;
	init : BOOL;
	last_dir : BOOL;
	start : BYTE;
END_VAR


(*
version 1.2		19. feb. 2011
programmer 	hugo
tested by		oscat

_RMP_B generates a ramp on an external var of type byte
tr is the total ramp time, E is an enable and dir the direction of the ramp

*)
(* @END_DECLARATION := '0' *)
(* read system timer *)
tx := DWORD_TO_TIME(T_PLC_MS());

IF E AND init AND (dir = last_dir) AND (RMP <> SEL(DIR, 0, 255)) AND TR = tn THEN
	RMP := FRMP_B(start, DIR, tx - tl, TR);
ELSE
	init := TRUE;
	tl := tx;
	tn := tr;
	start := RMP;
END_IF;
last_dir := dir;



(* revison history
hm	22. oct. 2008	rev 1.0
	original version

hm	20. nov. 2008	rev 1.1
	set default for E to be TRUE
	added init

hm	19. nov. 2011	rev 1.2
	new code
*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK _RMP_NEXT
VAR_INPUT
	E : BOOL := TRUE;
	IN : BYTE;
	TR : TIME;
	TF : TIME;
	TL : TIME;
END_VAR
VAR_OUTPUT
	DIR : BOOL;    (* upwards = TRUE *)
	UP : BOOL;
	DN : BOOL;
END_VAR
VAR_IN_OUT
	OUT : BYTE;
END_VAR
VAR
	rmx : _RMP_B;
	dirx : TREND_DW;
	t_lock :TP;
	xen: BOOL;
	xdir: BOOL;
END_VAR


(*
version 1.4		18. jul. 2009
programmer 	hugo
tested by		oscat


_RMP_NEXT  will generate a ramp output following the input IN.
on a rising ramp the output will stop as soon as it has surpassed the input in, and on a falling ramp it will stop when out is smaller than in.
a lockout time T_L will be added between up and down direction.

*)
(* @END_DECLARATION := '0' *)
dirx(X := in);

t_lock(in := FALSE, pt := TL);

IF dirx.TU AND (OUT < IN) THEN
	IF NOT xdir AND xen THEN t_lock(in := TRUE); END_IF;
	xen := TRUE;
	xdir := TRUE;
ELSIF dirx.TD AND (OUT > IN) THEN
	IF xdir AND xen THEN t_lock(in := TRUE); END_IF;
	xen := TRUE;
	xdir := FALSE;
ELSIF xen THEN
	IF (xdir AND (out >= in)) OR (NOT xdir AND (out <= in)) THEN
		xen := FALSE;
		IF tl > t#0s THEN t_lock(IN := TRUE); END_IF;
	END_IF;
END_IF;

IF NOT t_lock.Q AND xen THEN
		UP := XDIR;
		DIR := XDIR;
		DN := NOT XDIR;
ELSE
	UP := FALSE;
	DN := FALSE;
END_IF;

rmx(rmp := OUT, E := E AND (UP OR DN) , dir := DIR, tr := SEL(dir, TF, TR));


(* revison history
hm	23. nov. 2008	rev 1.0
	original version

hm	24. jan. 2009	rev 1.1
	deleted unused vars tmp1, tmp2

hm	20. feb. 2009	rev 1.2
	improved algorithm
	added TL

hm	9. mar. 2009	rev 1.3
	input E was ignored
	removed double assignments

hm	18. jul. 2009	rev 1.4
	improved performance

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK _RMP_W
VAR_INPUT
	DIR : BOOL;
	E : BOOL := TRUE;
	TR : TIME;
END_VAR
VAR_IN_OUT
	RMP : WORD;
END_VAR
VAR
	tx, tl : DWORD;
	step : DINT;
	init: BOOL;
	last_dir : BOOL;
END_VAR


(*
version 1.1	20. nov. 2008
programmer 	hugo
tested by	oscat


_RMP_B generates a ramp on an external var of type byte
tr is the total ramp time, E is an enable and dir the direction of the ramp

*)
(* @END_DECLARATION := '0' *)
tx := T_PLC_MS();
IF E AND init THEN
	(* we need to set tl = tx when direction changes *)
	IF dir XOR last_dir THEN
	 tl := tx;
	 last_dir := dir;
	END_IF;
	(* check_elapsed time and calculate only if necessary *)
	IF tr > t#0s THEN
		step := DWORD_TO_DINT(SHL(tx-tl, 16) / TIME_TO_DWORD(TR));
	ELSE
		step := 65535;
	END_IF;
	IF step > 0 THEN
		(* perform the step on the ramp *)
		tl := tx;
		(* calculate the step response *)
		IF NOT dir THEN step := - step; END_IF;
		rmp := DINT_TO_WORD(LIMIT(0, WORD_TO_DINT(rmp) + step, 65535));
	END_IF;
ELSE
	tl := tx;
	init := TRUE;
END_IF;


(* revison history
hm	22. oct. 2008	rev 1.0
	original version

hm	20. nov. 2008	rev 1.1
	set default for E to be TRUE
	added init

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK GEN_PULSE
VAR_INPUT
	ENQ : BOOL := TRUE;
	PTH : TIME;
	PTL : TIME;
END_VAR
VAR_OUTPUT
	Q : BOOL;
END_VAR
VAR
	tx: TIME;
	tn: TIME;
	init: BOOL;
END_VAR

(*
version 1.5	8. apr. 2011
programmer 	hugo
tested by		oscat

GEN_PULSE uses the internal sps timer to generate a continuous output waveform with programmable high and low time.
the accuracy of gen_pulse is depending on the system timer.
when time is 0 the high and low times are exactly one cycle.
ENQ = TRUE will start and ENQ = FALSE will stop the generator.

*)
(* @END_DECLARATION := '0' *)
IF enq THEN
	tx := DWORD_TO_TIME(T_PLC_MS());
	IF NOT init THEN init := TRUE; tn := tx; END_IF;
	IF tx - tn >= SEL(Q, PTL, PTH) THEN
		tn := tn + SEL(Q, PTL, PTH);
		Q := NOT Q;
	END_IF;
ELSE
	Q := FALSE;
	init := FALSE;
END_IF;



(* revision history
hm	29. jun. 2008	rev 1.0
	original version

hm	23. nov. 2008	rev 1.1
	set default for enq to be true

hm	18. jul. 2009	rev 1.2
	improved performance

hm	13. nov. 2009	rev 1.3
	corrected error

hm	16. feb. 2011	rev 1.4
	corrected an error when timer overflows

hm	8. apr. 2011	rev 1.5
	ptl and pth was exchanged
*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK GEN_PW2
VAR_INPUT
	ENQ : BOOL;
	TH1 : TIME;
	TL1 : TIME;
	TH2 : TIME;
	TL2 : TIME;
	TS : BOOL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	TH : TIME;
	TL : TIME;
END_VAR
VAR
	t_high : TIME;
	t_low : TIME;
	tx : TIME;
	start : TIME;
	init : BOOL;
	et : TIME;
END_VAR

(*
version 1.1	14. mar. 2009
programmer 	hugo
tested by		oscat

GEN_PW2 generates a time TH? followed by a time TL?.
the input ts selects between 2 sets of timings for the operation.


*)
(* @END_DECLARATION := '0' *)
(* read system timer *)
tx := DWORD_TO_TIME(T_PLC_MS());

(* startup initialization *)
IF NOT init THEN
	start := tx;
	init := TRUE;
	TH := T#0s;
	TL := T#0s;
END_IF;

(* timing selection *)
IF TS THEN
	t_high := TH2;
	t_low := TL2;
ELSE
	t_high := TH1;
	t_low := TL1;
END_IF;

(* normal operation *)
IF ENQ THEN
	et := tx - start;
	IF NOT Q THEN
		IF et >= t_low THEN
			Q := TRUE;
			start := tx;
			TL := T#0s;
		ELSE
			TL := et;
		END_IF;
	ELSE
		IF et >= t_high THEN
			Q := FALSE;
			start := tx;
			TH := T#0s;
		ELSE
			TH := et;
		END_IF;
	END_IF;
ELSE
	Q := FALSE;
	TH := T#0s;
	TL := T#0s;
	start := tx;
END_IF;



(* revision history
hm	26. sep. 2008	rev 1.0
	original version

hm	14. mar. 2009	rev 1.1
	removed double assignments

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK GEN_RDM
VAR_INPUT
	PT : TIME;
	AM : REAL := 1;
	OS : REAL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	Out : REAL;
END_VAR
VAR
	tx : TIME;
	last : TIME;
	init : BOOL;
END_VAR

(*
	version 1.1	16 sep 2007
	programmer 	oscat
	tested BY		oscat

this signal generator generates a random output. The signal is defined by period time (PT), 
amplitude (AM), offset (OS).
The Output waveform will have its max peak at AM/2 + OS and its minimum peak at -AM/2 + OS. 
The period time PT defines how often the output signal will jump to a new randow value.
The Output Q will be true for one cycle anytime the output OUT has changed

*)
(* @END_DECLARATION := '0' *)
(* read system time and prepare input data *)
tx := DWORD_TO_TIME(T_PLC_MS()) - last;

(* init section *)
IF NOT init THEN
	init := TRUE;
	last := tx;
	tx := t#0s;
END_IF;

(* add last if one cycle is finished *)
IF tx >= pt THEN
	last := last + pt;
	tx := tx - pt;

	(* generate output signal *)
	out := am * (RDM(0) - 0.5) + os;
	q := TRUE;
ELSE
	q := FALSE;
END_IF;

(* revision history

hm	7.2.2007		rev 1.0
	original version

hm	16.9.2007		rev 1.1
	changes time() to T_plc_ms() for compatibility reasons

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK GEN_RDT

VAR_INPUT
	Enable		 : BOOL := TRUE ;
	Min_Time_ms	: TIME := t#1s ; (* min Taktzeit *)
	Max_Time_ms	: TIME := t#1.2s ; (* Max Taktzeit *)
	TP_Q		 : TIME := t#100ms ; (* Zeit Ausgang auf TRUE *)
END_VAR
VAR_OUTPUT
	xQ		 : BOOL ;
END_VAR
VAR
	(* Taktgenerator Simulation *)
	tonRDMTimer	: TON ; (* Zeitbaustein Taktgenerator *)
	tof_xQ		: TOF ; (* Ausschaltverzögerung Taktgenerator *)
	tRDMTime	: TIME ; (* Sollzeit *)
	rRDMTime	: REAL ; (* Zufalswert Timer *)
END_VAR

(*
version 1.1	16 mar. 2008
programmer 	J. Schohaus
tested by		Hugo

RDMT generates a defined pulse with pulse width TP_Q at random times. the random time will be defined with an minimum and maximum time.

*)


(*
FUNCTION_BLOCK RDMTimer
###############################################################################
		 Ersteller / author :		 		 		 		 		 J.Schohaus
		 Datum / date:		 		 		 		 		 		 13.07.2007
###############################################################################
		 Änderungen / Datum / Ersteller :		 		 		 		 		 		 		 		 		 		 		 		 		 		 		 
		 moditication / date / author :		 		 
###############################################################################
		 Verwendete Bibliotheken		 ( * werden im Baustein nicht benötigt )
		 		 		 		 Oscat.lib
###############################################################################
Beschreibung:
*)
(* @END_DECLARATION := '0' *)
tonRDMTimer ( IN:= Enable , PT:= tRDMTime );
tof_xQ ( IN:= tonRDMTimer.Q , PT:= TP_Q );
XQ := tof_xq.Q;
IF tonRDMTimer.Q THEN
		 xQ := TRUE ;
		 rRDMTime := RDM ( last:= rRDMTime ) ;
		 tRDMTime := REAL_TO_TIME ( rRDMTime * DINT_TO_REAL(TIME_TO_DINT( Max_Time_ms - Min_Time_ms ) + TIME_TO_DINT(Min_Time_ms ))) ;
		 tonRDMTimer ( IN:= FALSE );
END_IF;

(* revision history
J. Schohaus	19. nov 2007	rev 1.0
	origial version

hm	16. mar. 2008			rev 1.1
	added type conversion to avoid warnings under codesys 3.0
*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/signal generators' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK GEN_RMP
VAR_INPUT
	PT : TIME := t#1s;
	AM : REAL := 1.0;
	OS : REAL;
	DL : REAL;
END_VAR
VAR_OUTPUT
	Q : BOOL;
	OUT : REAL;
END_VAR
VAR
	tx : TIME;
	last : TIME;
	init : BOOL;
	temp : REAL;
	ltemp: REAL;
END_VAR

(*
	version 1.4	10. mar. 2009
	programmer 	oscat
	tested BY		oscat

this signal generator generates a ramp wave output. The ramp wave signal is defined by period time (PT), 
amplitude (AM), offset (OS) and a specific delay for the output signal (DL).
The Output waveform will have its minimum peak at OS and its maximum peak at AM + OS. 
The delay input can delay a signal up to PT, this can be useful to synchronize different generators 
and generate interleaving signals.
in addition to the analog output Out there is a second boolean output Q with is true for one cycle when the ramp starts.
*)
(* @END_DECLARATION := '0' *)
(* read system time and prepare input data *)
tx := DWORD_TO_TIME(T_PLC_MS()) - last;
DL := MODR(dl,1.0);
IF dl < 0.0 THEN dl := 1.0 - dl; END_IF;

(* init section *)
IF NOT init THEN
	init := TRUE;
	last := tx;
	tx := t#0s;
END_IF;

(* add last if one cycle is finished *)
IF tx >= pt THEN
	last := last + pt;
	tx := tx - pt;
END_IF;

(* generate sine wave *)
ltemp := temp;
IF pt > t#0s THEN temp := FRACT(TIME_TO_REAL(tx + MULTIME(pt, dl)) / TIME_TO_REAL(pt)); END_IF;
out := am * temp + os;

(* boolean output Q *)
Q := temp < ltemp;

(* revision history
hm	3. mar 2007		rev 1.0
	original version

hm	17 sep 2007		rev 1.1
	replaced time() with t_plc_ms for compatibilitx reasons

hm	27. nov 2007	rev 1.2
	avoid divide by 0 when pt = 0

ks	26. oct. 2008	rev 1.3
	code optimization

hm	10. mar. 2009	rev 1.4
	changed real constants to use dot syntax
*)
END_FUNCTION_BLOCK