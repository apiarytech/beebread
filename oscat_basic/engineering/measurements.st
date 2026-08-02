(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK ALARM_2
VAR_INPUT
	X : REAL;
	LO_1 : REAL;
	HI_1 : REAL;
	LO_2 : REAL;
	HI_2 : REAL;
	HYS : REAL;
END_VAR
VAR_OUTPUT
	Q1_LO : BOOL;
	Q1_HI : BOOL;
	Q2_LO : BOOL;
	Q2_HI : BOOL;
END_VAR
VAR
	tmp: REAL;
END_VAR

(*
	version 1.1	11. mar. 2009
	programmer 	hugo
	tested BY		tobias

ALARM_2 will check two pairs of limits and signal when the input is above or below a set limit.
with the input HYS a hysteresis can be defined for all limits. 

*)
(* @END_DECLARATION := '0' *)
tmp := X - Hys * 0.5;
IF tmp > LO_1 THEN Q1_LO := FALSE; END_IF;
IF tmp > LO_2 THEN Q2_LO := FALSE; END_IF;
IF tmp > HI_1 THEN Q1_HI := TRUE; END_IF;
IF tmp > HI_2 THEN Q2_HI := TRUE; END_IF;
tmp := tmp + hys;
IF tmp < LO_1 THEN Q1_LO := TRUE; END_IF;
IF tmp < LO_2 THEN Q2_LO := TRUE; END_IF;
IF tmp < HI_1 THEN Q1_HI := FALSE; END_IF;
IF tmp < HI_2 THEN Q2_HI := FALSE; END_IF;



(* revision history
hm	19. may. 2008		rev 1.0
	original version

hm	11. mar. 2009	rev 1.1
	improved code
*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK BAR_GRAPH
VAR_INPUT
	X : REAL;
	rst : BOOL;
END_VAR
VAR_INPUT CONSTANT
	trigger_Low : REAL;
	trigger_High : REAL;
	Alarm_low : BOOL;
	Alarm_high : BOOL;
	log_scale : BOOL;
END_VAR
VAR_OUTPUT
	LOW, Q1, Q2, Q3, Q4, Q5, Q6, HIGH : BOOL;
	Alarm : BOOL;
	Status : BYTE;
END_VAR
VAR
	init : BOOL;
	T1, T2, T3, T4, T5 : REAL;
END_VAR
VAR_TEMP
	temp : REAL;
END_VAR

(*
	version 1.2	6 jan 2008
	programmer 	hugo
	tested BY	hans

bar graph is a muti window comparator which displays an analog input signal on 8 digital outputs.
only one output is active a any given time depending on the value of the input signal.
the output can be of linear or logarithmic scale.

*)
(* @END_DECLARATION := '0' *)
IF NOT init THEN
	init := TRUE;
	IF log_scale THEN
		temp := EXP(LN(Trigger_high / Trigger_low) * 0.166666666666666666666);
		T1 := trigger_low * temp;
		T2 := T1 * temp;
		T3 := T2 * temp;
		T4 := T3 * temp;
		T5 := T4 * temp;
	ELSE
		temp := (trigger_high - trigger_low) * 0.142857142;
		T1 := trigger_low + temp;
		T2 := T1 + temp;
		T3 := T2 + temp;
		T4 := T3 + temp;
		T5 := T4 + temp;
	END_IF;
END_IF;

(* clear outputs before checking *)
Q1 := FALSE;
Q2 := FALSE;
Q3 := FALSE;
Q4 := FALSE;
Q5 := FALSE;
Q6 := FALSE;
status := 110;

(* low, high and alarm can only be cleared with rst depending on alarm_low and alarm_high *)
IF NOT alarm_low THEN low := FALSE; END_IF;
IF NOT alarm_high THEN high := FALSE; END_IF;
IF rst THEN
	alarm := FALSE;
	low := FALSE;
	high := FALSE;
END_IF;

(* check and set outputs *)
IF X < trigger_low THEN
	Low := TRUE;
	status := 111;
	IF alarm_low THEN
		alarm := TRUE;
		status := 1;
	END_IF;
ELSIF X < T1 THEN
	Q1 := TRUE;
ELSIF x < t2 THEN
	Q2 := TRUE;
ELSIF x < t3 THEN
	Q3 := TRUE;
ELSIF x < T4 THEN
	Q4 := TRUE;
ELSIF x < T5 THEN
	q5 := TRUE;
ELSIF x < trigger_high THEN
	q6 := TRUE;
ELSE
	high := TRUE;
	status := 112;
	IF alarm_high THEN
		alarm := TRUE;
		status := 2;
	END_IF;
END_IF;

(* revision history
hm	22. feb 2007	rev 1.0
	original version

hm	2. dec 2007		rev 1.1
	chaged code for better performance

hm	6. jan 2008		rev 1.2
	further performance improvement

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK CALIBRATE
VAR_INPUT
	X : REAL;
	CO : BOOL;
	CS : BOOL;
END_VAR
VAR_INPUT CONSTANT
	Y_Offset : REAL;
	Y_Scale : REAL;
END_VAR
VAR_OUTPUT
	Y : REAL;
END_VAR
VAR RETAIN
	offset : REAL;
	Scale : REAL := 1.0;
END_VAR

(*
	version 1.3	11. mar. 2009
	programmer 	hugo
	tested BY		hans

Calibrate allows for offset and scale calibration of an analog input.
offset is calibrated to a stored reference Y_offset while CO is true.
after the offset is calibrated, scale can be calibrated to a reference value Y_scale while CS is true.
*)
(* @END_DECLARATION := '0' *)
(* check for calibration *)
IF CO THEN
	OFFSET := Y_Offset - X;
ELSIF CS THEN
	SCALE := Y_scale / (X + OFFSET);
END_IF;
(* calculate output *)
Y := (X + OFFSET) * SCALE;


(* revision history
hm 22.2.2007		rev 1.2
	changed VAR RETAIN PERSISTENT to VAR RETAIN for better compatibility
	wago lon contollers do not support persisitent

hm	11. mar. 2009	rev 1.3
	changed real constants to use dot syntax

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK CYCLE_TIME
VAR_INPUT
	RST : BOOL;
END_VAR
VAR_OUTPUT
	ct_min : TIME;
	ct_max : TIME;
	ct_last : TIME;
	systime : TIME;
	sysdays : INT;
	cycles : DWORD;
END_VAR
VAR
	last_cycle : TIME;
	tx: TIME;
	init: BOOL;
END_VAR


(*
	version 1.2	16 sep 2007
	programmer 	hugo
	tested BY		hans

this function block measures the cycle time and displays the last, min and max cycle time of the current task.
the resolution is 1ms.
the cycles output is a dword counter which counts the cycles.
a rst pulse on the input will reset all data.


*)
(* @END_DECLARATION := '0' *)
tx := DWORD_TO_TIME(T_PLC_MS()) - last_cycle;
IF rst THEN
	ct_min := t#10h;
	ct_max := t#0ms;
	cycles := 0;
ELSIF last_cycle > t#0s THEN
	IF tx < ct_min THEN ct_min := tx;
	ELSIF tx > ct_max THEN ct_max := tx;
	END_IF;
	ct_last := tx;
ELSIF ct_min = t#0s THEN
	ct_min := t#0s - t#1ms;
END_IF;
IF init THEN
	systime := systime + tx;
		IF systime >= t#1d THEN
			systime := systime - t#1d;
			sysdays := sysdays + 1;
		END_IF;
END_IF;
init := TRUE;
last_cycle := last_cycle + tx;
cycles := cycles + 1;

(*	revision history
hm 12.12.2006		rev 1.1
	added cycles output, a dword cycle counter.
hm 10.3.2007			rev 1.2
	changed initialization of ct_min to t#10h for compatibility with siemens s7

hm	16.9.2007		rev 1.2
	changed Time() in T_PLC_MS() for compatibility resons

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK DT_SIMU
VAR_INPUT
	START : DT;
	SPEED : REAL := 1.0;
END_VAR
VAR_OUTPUT
	DTS : DT;
END_VAR
VAR
	tc : DWORD;
	init: BOOL;
	last: DWORD;
	tx: DWORD;
	td: DWORD;
END_VAR

(*
version 1.2	8. mar. 2009
programmer 	hugo
tested by		oscat

DT_SIMU simulates a real time clock and can be adjusted to different speeds
it can also be used in simulation to simulate a real time clock.
the peed of the clock can be increased or decreased to debug timers.
with the input start a start date-time can be specified.

*)
(* @END_DECLARATION := '0' *)
(* read system timer *)
tx := T_PLC_MS();
tc := REAL_TO_DWORD(DWORD_TO_REAL(tx - last) * speed);

IF NOT init THEN
	init := TRUE;
	DTS := Start;
	tc := 0;
	last := tx;
ELSIF SPEED = 0.0 THEN
	DTS := DWORD_TO_DT(DT_TO_DWORD(DTS) + 1);
ELSIF tc >= 1000 THEN
	td := (tc / 1000) * 1000;
	DTS := DTS + DWORD_TO_TIME(td);
	last := last + REAL_TO_DWORD(DWORD_TO_REAL(td) / speed);
END_IF;

(* revision history
hm	11. sep. 2008	rev 1.0
	original version

hm	16. nov	2008	rev 1.1
	added type conversions for compatibility reasons

hm	8.	mar. 2009	rev 1.2
	added increment by cycle mode

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK FLOW_METER
VAR_INPUT
	VX : REAL;
	E : BOOL;
	RST : BOOL;
END_VAR
VAR_INPUT CONSTANT
	PULSE_MODE : BOOL;
	UPDATE_TIME : TIME := t#1s;
END_VAR
VAR_OUTPUT
	F : REAL;
END_VAR
VAR_IN_OUT
	X : REAL;
	Y : UDINT;
END_VAR
VAR
	tx, tl : TIME;
	int1 : INTEGRATE;
	init: BOOL;
	e_last : BOOL;
	tmp: INT;
	x_last : REAL;
	y_last : UDINT;
END_VAR

(*
version 1.0	23. jan. 2011
programmer 	hugo
tested by	oscat

Flow meter measures flow according to gated time or pulses.

*)

(* @END_DECLARATION := '0' *)
IF NOT init THEN	(* init on power up *)
	init := TRUE;
	tl := tx;
	x_last := X;
	y_last := Y;
	int1.K := 2.7777777777777777E-4;
END_IF;

(* run integrator *)
int1(E := NOT (RST OR PULSE_MODE) AND E, X := VX, Y := X);	(* gated operation *)

IF RST THEN		(* reset *)
	X := 0.0;
	Y := 0;
	tl := tx;
	x_last := 0.0;
	y_last := 0;
ELSIF E AND PULSE_MODE THEN	(* check for pulse mode *)
	IF NOT e_last THEN X := X + VX; END_IF;
END_IF;
e_last := E;

(* reduce X to be less than 1 and increase Y respectively *)
IF X > 1.0 THEN
	tmp := FLOOR(X);
	Y := Y + INT_TO_UDINT(tmp);
	X := X - INT_TO_REAL(tmp);
END_IF;

(* calculate the current flow *)
tx := DWORD_TO_TIME(T_PLC_MS());
IF tx - tl >= UPDATE_TIME AND UPDATE_TIME > t#0s THEN
	F := (UDINT_TO_REAL(Y - y_last) + X - x_last) / TIME_TO_REAL(tx - tl) * 3.6E6;
	y_last := Y;
	x_last := X;
	tl := tx;
END_IF;



(* revision history
hm	23. jan. 2011	rev 1.0
	original version

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK M_D
VAR_INPUT
	start : BOOL;
	stop : BOOL;
	tmax : TIME := t#10d;
	rst: BOOL;
END_VAR
VAR_OUTPUT
	PT : TIME;
	ET : TIME;
	run : BOOL;
END_VAR
VAR
	edge : BOOL;
	T0 : TIME;
	tx: TIME;
	startup: BOOL;
END_VAR

(*
	version 1.2	27. feb. 2009
	programmer 	hugo
	tested BY	oscat

m_d measures the time between a rising edge on start to a rising edge on stop and returs the last measured value on the output PT.
a second output ET is starting from 0 at the rising edge of start and is counting up until the rising edge of stop occurs.
the asynchronous input rst can reset the outputs at any time.
tmax defines a timeout wich is the maximum measurable time between start and stop, if this time is exceeded the outputs will stay at 0.
*)
(* @END_DECLARATION := '0' *)
(* check for rst input *)
IF rst OR et >= tmax THEN
	pt := t#0ms;
	et := t#0ms;
	startup := FALSE;
	run := FALSE;
END_IF;

(* avoid timers to start when input is true at startup *)
IF NOT startup THEN
	edge := start;
	startup := TRUE;
END_IF;

(* read system timer *)
tx := DWORD_TO_TIME(T_PLC_MS());

(* detect rising edge on start *)
IF start AND NOT edge AND NOT stop THEN
	t0 := tx;
	run := TRUE;
	pt := t#0s;
ELSIF stop AND run THEN
	pt := et;
	run := FALSE;
END_IF;
edge := start;
IF run THEN et := tx - t0; END_IF;

(* revision history
hm		2.5.2007	rev 1.0
	original version

hm		16.9.2007	rev 1.1
	changes time() to T_plc_ms() for compatibility reasons

hm	27. feb 2009	rev 1.2
	deleted unnecessary init with 0
*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK M_T
VAR_INPUT
	IN : BOOL;
	TMAX : TIME := t#10d;
	RST : BOOL;
END_VAR
VAR_OUTPUT
	PT : TIME;
	ET: TIME;
END_VAR
VAR
	edge : BOOL;
	start : TIME;
	tx: TIME;
END_VAR

(*
version 1.4	10. mar. 2009
programmer 	hugo
tested BY		oscat

m_t measures the with of a high pulse and returs the last measured pulse width on output PT.
a second output ET is starting from 0 at the rising edge and counting up until the falling edge occurs and resetts et to 0.
the asynchrtonous input rst can reset the outputs at any time.
tmax defines a maximum measurable time, if this value is exceeded the outputs will be reset to 0.

*)
(* @END_DECLARATION := '0' *)
(* read system time *)
tx := DWORD_TO_TIME(T_PLC_MS());

IF RST OR ET >= TMAX THEN
	PT :=  t#0s;
	ET := t#0s;
ELSIF IN THEN
	IF NOT edge THEN start := tx; END_IF;
	ET := tx - start;
ELSE
	PT := ET;
END_IF;
edge := IN;



(* revision history
hm	4. aug. 2006	rev 1.0
	original version

hm	2. may. 2007	rev 1.1
	added init variable to avoid unreasonable measurement if falling edge occurs first.
	added et (elapsed time) output to count from 0 at rising edge until a falling edge resets et to 0.
	added reset input.

hm	16. sep. 2007	rev 1.2
	changes time() to T_plc_ms() for compatibility reasons

hm	17. dec. 2008	rev 1.3
	code optimized

hm	10. mar. 2009	rev 1.4
	removed nested comments
	removed double assignments

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK M_TX
VAR_INPUT
	in: BOOL;
	tmax : TIME := t#10d;
	rst : BOOL;
END_VAR
VAR_OUTPUT
	TH : TIME;
	TL : TIME;
	DC : REAL;
	F: REAL;
	ET : TIME;
END_VAR
VAR
	edge : BOOL;
	start: TIME;
	stop:TIME;
	tx: TIME;
	rise : BOOL;
	fall : BOOL;
	startup: BOOL;
END_VAR

(*
	version 1.4	11. mar. 2009
	programmer 	oscat
	tested BY		oscat

this measures the timing of a signal
TH is the high pulse width
TL is the low pulse width
DC is the duty cycle
F is the frequency
ET will start at 0 with a rising edge and count up till the next rising edge starts from 0 again. 
the acuracy depends on the cycle time of the sps and is limited to 1ms resolution.
an asynchronous reset can reset all outputs to 0 at any time.
the input tmax can configure a timeout for ET where the function block will outomatically reset itself.
tmx is predefined to 10 days.
*)
(* @END_DECLARATION := '0' *)
(* check FOR rst input *)
IF rst OR (et >= tmax) THEN
	rise := FALSE;
	fall := FALSE;
	startup := FALSE;
	th := t#0ms;
	tl := t#0ms;
	DC := 0.0;
	F := 0.0;
	ET := t#0s;
END_IF;

(* avoid timers to start when input is true at startup *)
IF NOT startup THEN
	edge := in;
	startup := TRUE;
END_IF;

(* read system timer *)
tx := DWORD_TO_TIME(T_PLC_MS());

(* edge trigger rising and falling edge *)
IF in XOR edge THEN
	edge := in;
	(* rising edge starts measurement *)
	IF in THEN
		start := Tx;
		rise := TRUE;
		IF fall THEN tl := start - stop; END_IF;
		IF th > t#0ms AND tl > t#0ms THEN
			dc := TIME_TO_REAL(th) / TIME_TO_REAL(th+tl);
			f := 1000.0 / TIME_TO_REAL(th + tl);
		END_IF;
	(* falling edge starts second cycle measurement *)
	ELSE
		stop := Tx;
		fall := TRUE;
		IF rise THEN th := stop - start; END_IF;
		IF th > t#0ms AND tl > t#0ms THEN
			dc := TIME_TO_REAL(th) /  TIME_TO_REAL(th+tl);
			f := 1000.0 / TIME_TO_REAL(th + tl);
		END_IF;
	END_IF;
END_IF;
IF rise THEN et := tx - start; END_IF;

(* revision history
hm	4. aug. 2006	rev 1.0
	original version

hm	2. mai. 2007	rev 1.1
	made sure that no undefined values would be measured at startup conditions
	added output et
	added rst input
	added tmax input

hm	16. sep. 2007	rev 1.2
	changes time() to T_plc_ms() for compatibility reasons

hm	27. feb. 2009	rev 1.3
	deleted unnecessary init with 0

hm	11. mar. 2009	rev 1.4
	changed real constants to use dot syntax

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK METER
VAR_INPUT
	M1, M2 : REAL;
	I1, I2 : BOOL;
	D : REAL := 1.0;
	RST : BOOL;
END_VAR
VAR_IN_OUT
	MX : REAL;
END_VAR
VAR
	MR : REAL2;
	MX1, MX2 : REAL; (* current consumption value on M1 and M2 *)
	tx: DWORD;
	last: DWORD;
	tc: REAL;
	init: BOOL;
END_VAR

(*
	version 1.4	8. mar. 2009
	programmer 	oscat
	tested BY	oscat

METER measures usage of power or similar values the output MX is the sum of the inputs over time.
MX := sum over time of (I1 * P1 + I2 * P2)/D.
a simple example is the consumption counter for a 2 stage oil burner where stage 1 is 10KW and stage 2 is 20KW.
I1 and I2 decide which value is accounted for, I1 = True only counts P1, I2 is True counts P2 
and I1 and I2 is true counts P1 and P2.
the meter can be used for all kind of consumption meters. P1 and P2 can change on the fly.

*)
(* @END_DECLARATION := '0' *)
(* measure the last cycle time and make sure we onle execute once every millisecond *)
tx := T_PLC_MS();
IF NOT init THEN
	init := TRUE;
	last := tx;
	mr.RX := mx;
	mr.R1 := 0.0;
ELSIF tx = last THEN
	RETURN;
ELSE
	tc := DWORD_TO_REAL(tx - last) * 0.001;
END_IF;
last := tx;

(* reset *)
IF rst THEN
	mr.R1 := 0.0;
	mr.RX := 0.0;
ELSE
	(* current consumption measurement *)
	IF I1 THEN MX1 := M1; ELSE MX1 := 0.0; END_IF;
	IF I2 THEN MX2 := M2; ELSE MX2 := 0.0; END_IF;
	(* add up the current values in a double real *)
	MR := R2_ADD(MR,(SEL(I1,0.0,mx1) + SEL(I2, 0.0, mx2)) / D * TC);
	(* set the current output value *)
	MX := mr.RX;
END_IF;




(* revision history

hm	16. sep.2007		rev 1.0
	original version

hm	7. feb 2008		rev 1.1
	use new version of ft_int to avoid a counting stop at high values
	deleted unnecessary limits

hm	24. mar. 2008	rev 1.2
	use data_type real2 to extend accuracy to 15 digits total
	do not use ft_int which adds unnecessary overhead

hm	8. feb. 2009	rev 1.3
	changed mx to be I/O
	make sure calculation works for cycle times < 1 ms

hm	8. mar. 2009	rev 1.4
	last was not updated
	code improvements

*)

END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Engineering\/measurements' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK METER_STAT
VAR_INPUT
	IN : REAL;
	DI : DATE;
	RST : BOOL;
END_VAR
VAR_IN_OUT
	Last_Day : REAL;
	Current_Day : REAL;
	Last_Week : REAL;
	Current_Week : REAL;
	Last_Month : REAL;
	Current_Month : REAL;
	Last_Year : REAL;
	Current_Year : REAL;
END_VAR
VAR RETAIN
	Year_Start : REAL;
	Month_Start : REAL;
	Week_Start : REAL;
	Day_Start : REAL;
	last_run: DATE;
END_VAR

(*
	version 1.3	18. jul. 2009
	programmer 	oscat
	tested BY		oscat

METER_STAT runs statistics of a metered value, it calculates the current week, day, month and year consumption and stores the corresponding last day, week, month and year value.

*)
(* @END_DECLARATION := '0' *)
IF rst THEN
	Last_Day := 0.0;
	Current_Day := 0.0;
	Day_start := IN;
	Last_week := 0.0;
	Current_week := 0.0;
	Week_start := in;
	Last_month := 0.0;
	Current_month := 0.0;
	month_start := in;
	last_year := 0.0;
	current_year := 0.0;
	year_start := in;
ELSE
	Current_Day := IN - Day_Start;
	Current_Week := In - Week_Start;
	Current_Month := IN - Month_Start;
	Current_Year := IN - Year_Start;
END_IF;

IF YEAR_OF_DATE(DI) > YEAR_OF_DATE(last_run) THEN
	(* a new year has started *)
	last_year := current_year;
	current_year := 0.0;
	year_start := in;
	last_month := current_month;
	current_month := 0.0;
	month_start := in;
	last_day := current_day;
	current_day := 0.0;
	day_start := in;
ELSIF MONTH_OF_DATE(DI) > MONTH_OF_DATE(last_run) THEN
	(* a new month has started, january is alrerady done by year change *)
	last_month := current_month;
	current_month := 0.0;
	month_start := in;
	last_day := current_day;
	current_day := 0.0;
	day_start := in;
ELSIF DAY_OF_YEAR(di) > DAY_OF_YEAR(last_run) THEN
	(* day has chaged, first day of year and first day of month has already been taken care of *)
	last_day := current_day;
	current_day := 0.0;
	day_start := in;
END_IF;
IF DAY_OF_WEEK(DI) < DAY_OF_WEEK(last_run) THEN
	(* a new week has started *)
	last_week := current_week;
	current_week := 0.0;
	week_start := in;
END_IF;
last_run := di;

(* revision history
hm	16. sep. 2007	rev 1.0
	original version

hm	7. oct. 2008	rev 1.1
	changed function year_to_year_of_date
	changed function month to month_of_date
	changed function weekday to day_of_week

hm	11. mar. 2009	rev 1.2
	changed real constants to use dot syntax

hm	18. jul. 2009	rev 1.3
	changes all outputs to be I/O
*)

END_FUNCTION_BLOCK