(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK CALENDAR_CALC
VAR_INPUT
	UTC : DT;
	LOCATION_NO : INT;
	DST_ENABLE : BOOL;
	LANGUAGE_NO : INT;
	LONGITUDE : REAL;
	LATITUDE : REAL;
END_VAR
VAR_OUTPUT
	CAL : CALENDAR;
END_VAR
VAR
	init : BOOL;
	last_day : DATE;
	loc : INT;
	lan : INT;
END_VAR

(*
version 1.3	18. jan. 2011
programmer 	hugo
tested by		oscat

calendar_calc calculates all calendar data based on UTC time and location data.

*)
(* @END_DECLARATION := '0' *)
(* startup initialisation *)
IF NOT init THEN
	init := TRUE;
	last_day := DWORD_TO_DATE(1);
END_IF;

(* set location and language *)
loc := SEL(LOCATION_NO = 0, LOCATION.DEFAULT, LOCATION_NO);
lan := SEL(LANGUAGE_NO = 0, LOCATION.LANGUAGE[loc], LANGUAGE_NO);

(* calculate calendar data *)
CAL.UTC := UTC;
CAL.DST_EN := DST_ENABLE;
CAL.LANGUAGE := lan;
CAL.LONGITUDE := LONGITUDE;
CAL.LATITUDE := LATITUDE;

(* calculate local time *)
CAL.OFFSET := UTC_TO_LTIME(UTC, loc, CAL.DST_ON);
CAL.LOCAL_DT := UTC + INT_TO_TIME(CAL.OFFSET * 60);

(* calculate date and time *)
sdt_to_date(dt_to_sdt(cal.local_dt), cal.local_date);
sdt_to_tod(dt_to_sdt(cal.local_dt), cal.local_tod);

(* calculate year, month, day and weekday *)
cal.YEAR := YEAR_OF_DATE(cal.local_date);
cal.MONTH := MONTH_OF_DATE(cal.local_date);
cal.DAY := DAY_OF_DATE(cal.local_date);
cal.WEEKDAY := DAY_OF_WEEK(cal.local_date);

(* calculate sun position only once a day *)
IF cal.local_date <> last_day THEN
	last_day := cal.local_date;
	sun_time(cal.local_date, longitude, latitude, cal.sun_rise, cal.sun_set, cal.sun_midday, cal.sun_heigth);
END_IF;

(* calculate current sun position *)
sun_pos(cal.local_dt, longitude, latitude, cal.sun_hor, cal.sun_ver);

(* calculate night and day *)
cal.NIGHT := cal.local_tod < cal.sun_rise OR cal.local_tod > cal.sun_set;

(* calculate holiday *)
cal.HOLIDAY := holiday(cal.local_date, loc, cal.holy_name);

(* calculate work week *)
cal.WORK_WEEK := WORK_WEEK(cal.local_date);


(* revision history
hm	21. nov. 2008	rev 1.0
	original version

hm	8. feb. 2009	rev 1.1
	added sun position calculation

hm	10. mar. 2009	rev 1.2
	added work_week, sun_midday, sun_heigth

hm	18. jan. 2011	rev 1.3
	location and language can be selected

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DATE_ADD : DATE
VAR_INPUT
	DA : DATE;
	T : TIME;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function adds a time to a date

*)
(* @END_DECLARATION := '0' *)
DATE_ADD := DWORD_TO_DATE(DATE_TO_DWORD(DA) + TIME_TO_DWORD(T));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAY_OF_DATE : INT
VAR_INPUT
	DI : DATE;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the day of a date

*)
(* @END_DECLARATION := '0' *)
DAY_OF_DATE := DATE_TO_DINT(DI) MOD 31 + 1;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAY_OF_MONTH : INT
VAR_INPUT
	DI : DATE;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the day of a month

*)
(* @END_DECLARATION := '0' *)
DAY_OF_MONTH := DATE_TO_DINT(DI) MOD 31 + 1;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAY_OF_WEEK : INT
VAR_INPUT
	DI : DATE;
END_VAR


(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function returns the day of the week of a date

*)
(* @END_DECLARATION := '0' *)
DAY_OF_WEEK := (DATE_TO_DINT(DI) / 31) MOD 7 + 1;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	renamed from weekday to day_of_week

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAY_OF_YEAR : INT
VAR_INPUT
	DI : DATE;
END_VAR
VAR
	temp : INT;
END_VAR

(*
version 1.1	18. oct. 2008
programmer 	hugo
tested by	oscat

this function returns the day of the year of a date

*)
(* @END_DECLARATION := '0' *)
temp := MONTH_OF_DATE(DI);
DAY_OF_YEAR := setup.MTH_OFS[temp] + DAY_OF_DATE(DI);
IF LEAP_YEAR(YEAR_OF_DATE(DI)) AND temp > 2 THEN
	DAY_OF_YEAR := DAY_OF_YEAR + 1;
END_IF;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	18. oct. 2008	rev 1.1
	renamed from yearday to day_of_year

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAY_TO_TIME : TIME
VAR_INPUT
	D : INT;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function converts days to time

*)
(* @END_DECLARATION := '0' *)
DAY_TO_TIME := DINT_TO_TIME(INT_TO_DINT(D) * 86400000);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAYS_DELTA : INT
VAR_INPUT
	D1 : DATE;
	D2 : DATE;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the number of days between two dates

*)
(* @END_DECLARATION := '0' *)
DAYS_DELTA := DINT_TO_INT(DATE_TO_DINT(D2) - DATE_TO_DINT(D1));


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAYS_IN_MONTH : INT
VAR_INPUT
	M : INT;
	Y : INT;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the number of days in a month

*)
(* @END_DECLARATION := '0' *)
IF M = 2 THEN
	DAYS_IN_MONTH := 28 + LEAP_DAY(Y);
ELSE
	DAYS_IN_MONTH := 30 + (M + (M/8)) MOD 2;
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DAYS_IN_YEAR : INT
VAR_INPUT
	Y : INT;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the number of days in a year

*)
(* @END_DECLARATION := '0' *)
DAYS_IN_YEAR := 365 + LEAP_DAY(Y);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DST : BOOL
VAR_INPUT
	DI : DATE;
END_VAR
VAR
	y : INT;
	m : INT;
	d : INT;
	w : INT;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns true if a date is in daylight saving time

*)
(* @END_DECLARATION := '0' *)
y := YEAR_OF_DATE(DI);
m := MONTH_OF_DATE(DI);
d := DAY_OF_DATE(DI);
w := DAY_OF_WEEK(DI);

IF m > 3 AND m < 10 THEN
	DST := TRUE;
ELSIF m = 3 AND d - w > 24 THEN
	DST := TRUE;
ELSIF m = 10 AND d - w < 25 THEN
	DST := TRUE;
END_IF;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DT2_TO_SDT : SDT
VAR_INPUT
	DT_IN : DT;
END_VAR
VAR
	temp : SDT;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function converts a DT into a structured data type SDT

*)
(* @END_DECLARATION := '0' *)
temp.YEAR := YEAR_OF_DATE(DT_TO_DATE(DT_IN));
temp.MONTH := MONTH_OF_DATE(DT_TO_DATE(DT_IN));
temp.DAY := DAY_OF_DATE(DT_TO_DATE(DT_IN));
temp.WEEKDAY := DAY_OF_WEEK(DT_TO_DATE(DT_IN));
temp.HOUR := HOUR_OF_DT(DT_IN);
temp.MINUTE := MINUTE_OF_DT(DT_IN);
temp.SECOND := SECOND_OF_DT(DT_IN);
DT2_TO_SDT := temp;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION DT_TO_SDT : SDT
VAR_INPUT
	DT_IN : DT;
END_VAR
VAR
	temp : SDT;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function converts a DT into a structured data type SDT

*)
(* @END_DECLARATION := '0' *)
temp.YEAR := YEAR_OF_DATE(DT_TO_DATE(DT_IN));
temp.MONTH := MONTH_OF_DATE(DT_TO_DATE(DT_IN));
temp.DAY := DAY_OF_DATE(DT_TO_DATE(DT_IN));
temp.WEEKDAY := DAY_OF_WEEK(DT_TO_DATE(DT_IN));
temp.HOUR := HOUR_OF_DT(DT_IN);
temp.MINUTE := MINUTE_OF_DT(DT_IN);
temp.SECOND := SECOND_OF_DT(DT_IN);
DT_TO_SDT := temp;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION EASTER : DATE
VAR_INPUT
	Y : INT;
END_VAR
VAR
	A, B, C, D, E, F, G, H, I, K, L, M, P : INT;
END_VAR

(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function calculates the date of easter for a given year

*)
(* @END_DECLARATION := '0' *)
A := Y MOD 19;
B := Y / 100;
C := Y MOD 100;
D := B / 4;
E := B MOD 4;
F := (B + 8) / 25;
G := (B - F + 1) / 3;
H := (19 * A + B - D - G + 15) MOD 30;
I := C / 4;
K := C MOD 4;
L := (32 + 2 * E + 2 * I - H - K) MOD 7;
M := (A + 11 * H + 22 * L) / 451;
P := (H + L - 7 * M + 114) MOD 31;
EASTER := SET_DATE(Y, (H + L - 7 * M + 114) / 31, P + 1);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION_BLOCK EVENTS
VAR_INPUT
	E : BOOL;
	DT_IN : DT;
END_VAR
VAR_IN_OUT
	EVENT : ARRAY[0..15] OF TIMER_EVENT;
END_VAR
VAR_OUTPUT
	Q : ARRAY[0..15] OF BOOL;
END_VAR
VAR
	i : INT;
	on : BOOL;
	td : TOD;
	di : DATE;
	dow : INT;
	loc : INT;
	on_time : TIME;
END_VAR

(*
version 1.1	26. jan. 2011
programmer 	hugo
tested by		oscat

EVENTS is a 16 channel timer switch.

*)
(* @END_DECLARATION := '0' *)
IF E THEN
	td := DT_TO_TOD(DT_IN);
	di := DT_TO_DATE(DT_IN);
	dow := DAY_OF_WEEK(di);
	FOR i := 0 TO 15 DO
		loc := SEL(event[i].land = 0, 1, event[i].land);
		on := FALSE;
		IF event[i].typ > 0 THEN
			on_time := event[i].start + event[i].duration;
			IF event[i].start < on_time THEN
				on := td >= event[i].start AND td < on_time;
			ELSE
				on := td >= event[i].start OR td < on_time;
			END_IF;
		END_IF;
		CASE event[i].typ OF
			1:	(* daily event *)
				Q[i] := on;
			2:	(* weekly event *)
				Q[i] := on AND BIT_OF_DWORD(BYTE_TO_DWORD(event[i].day),dow-1);
			3:	(* monthly event *)
				Q[i] := on AND DAY_OF_DATE(di) = event[i].day;
			4:	(* yearly event *)
				Q[i] := on AND DAY_OF_DATE(di) = event[i].day AND MONTH_OF_DATE(di) = event[i].lor;
			5:	(* holiday event *)
				Q[i] := on AND holiday(di,loc,s);
			6:	(* single event *)
				IF on AND event[i].last <> di THEN
					Q[i] := TRUE;
					event[i].last := di;
				ELSE
					Q[i] := FALSE;
				END_IF;
		END_CASE;
	END_FOR;
END_IF;


(* revision history
hm	22. jan. 2007	rev 1.0
	original version

hm	26. jan. 2011	rev 1.1
	changed type of last to be DT

*)
END_FUNCTION_BLOCK


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION HOLIDAY : BOOL
VAR_INPUT
	DI : DATE;
	L : INT;
END_VAR
VAR_IN_OUT
	NAME : STRING(30);
END_VAR
VAR
	i : INT;
	d : DATE;
	m : INT;
	y : INT;
	w : INT;
	hd : HOLIDAY_DATA;
	h_de : ARRAY[0..29] OF HOLIDAY_DATA := (name:='Neujahr',day:=1,month:=1,use:=1),(name:='Heilig Drei Knige',day:=6,month:=1,use:=1),(name:='Karfreitag',day:=-2,month:=0,use:=1),(name:='Ostersonntag',day:=0,month:=0,use:=1),(name:='Ostermontag',day:=1,month:=0,use:=1),(name:='Tag der Arbeit',day:=1,month:=5,use:=1),(name:='Christi Himmelfahrt',day:=39,month:=0,use:=1),(name:='Pfingstsonntag',day:=49,month:=0,use:=1),(name:='Pfingstmontag',day:=50,month:=0,use:=1),(name:='Fronleichnam',day:=60,month:=0,use:=1),(name:='Augsburger Friedensfest',day:=8,month:=8,use:=0),(name:='Maria Himmelfahrt',day:=15,month:=8,use:=1),(name:='Tag der Deutschen Einheit',day:=3,month:=10,use:=1),(name:='Reformationstag',day:=31,month:=10,use:=0),(name:='Allerheiligen',day:=1,month:=11,use:=1),(name:='Buss und Bettag',day:=-32,month:=11,use:=0),(name:='1. Weihnachtstag',day:=25,month:=12,use:=1),(name:='2. Weihnachtstag',day:=26,month:=12,use:=1);
	h_at : ARRAY[0..29] OF HOLIDAY_DATA := (name:='Neujahr',day:=1,month:=1,use:=1),(name:='Heilig Drei Knige',day:=6,month:=1,use:=1),(name:='Karfreitag',day:=-2,month:=0,use:=1),(name:='Ostersonntag',day:=0,month:=0,use:=1),(name:='Ostermontag',day:=1,month:=0,use:=1),(name:='Staatsfeiertag',day:=1,month:=5,use:=1),(name:='Christi Himmelfahrt',day:=39,month:=0,use:=1),(name:='Pfingstsonntag',day:=49,month:=0,use:=1),(name:='Pfingstmontag',day:=50,month:=0,use:=1),(name:='Fronleichnam',day:=60,month:=0,use:=1),(name:='Nationalfeiertag',day:=26,month:=10,use:=1),(name:='Maria Himmelfahrt',day:=15,month:=8,use:=1),(name:='Allerheiligen',day:=1,month:=11,use:=1),(name:='Maria Empfngnis',day:=8,month:=12,use:=1),(name:='1. Weihnachtstag',day:=25,month:=12,use:=1),(name:='2. Weihnachtstag',day:=26,month:=12,use:=1);
	h_fr : ARRAY[0..29] OF HOLIDAY_DATA := (name:='Nouvel an',day:=1,month:=1,use:=1),(name:='Vendredi Saint (alsace)',day:=-2,month:=0,use:=0),(name:='Dimanche de pques',day:=0,month:=0,use:=1),(name:='Lundi de pques',day:=1,month:=0,use:=1),(name:='Jeudi de Ascension',day:=39,month:=0,use:=1),(name:='dimanche de Pentecte ',day:=49,month:=0,use:=1),(name:='Fte du travail',day:=1,month:=5,use:=1),(name:='Victoire 1945 ',day:=8,month:=5,use:=1),(name:='Prise de La bastille',day:=14,month:=7,use:=1),(name:='15 Aot 1944',day:=15,month:=8,use:=1),(name:='Toussaint',day:=1,month:=11,use:=1),(name:='Armistice 1918',day:=11,month:=11,use:=1),(name:='Nol',day:=25,month:=12,use:=1),(name:='Saint tienne (alsace)',day:=26,month:=12,use:=0);
	h_bed : ARRAY[0..29] OF HOLIDAY_DATA := (name:='Neujahr',day:=1,month:=1,use:=1),(name:='Ostersonntag',day:=0,month:=0,use:=1),(name:='Ostermontag',day:=1,month:=0,use:=1),(name:='Tag der Arbeit',day:=1,month:=5,use:=1),(name:='Christi Himmelfahrt',day:=39,month:=0,use:=1),(name:='Pfingsten',day:=49,month:=0,use:=1),(name:='Pfingstmontag',day:=50,month:=0,use:=1),(name:='Nationalfeiertag',day:=21,month:=7,use:=1),(name:='Mari Himmelfahrt',day:=15,month:=8,use:=1),(name:='Allerheiligen',day:=1,month:=11,use:=1),(name:='Feiertag DG',day:=15,month:=11,use:=1),(name:='Heiligabend',day:=24,month:=12,use:=1),(name:='1. Weihnachtstag',day:=25,month:=12,use:=1),(name:='2. Weihnachtstag',day:=26,month:=12,use:=1),(name:='Silvester',day:=31,month:=12,use:=1);
	h_itd : ARRAY[0..29] OF HOLIDAY_DATA := (name:='Neujahr',day:=1,month:=1,use:=1),(name:='Heilig Drei Knige',day:=6,month:=1,use:=1),(name:='Ostersonntag',day:=0,month:=0,use:=1),(name:='Ostermontag',day:=1,month:=0,use:=1),(name:='Tag der Befeiung Italiens',day:=25,month:=4,use:=1),(name:='Tag der Arbeit',day:=1,month:=5,use:=1),(name:='Pfingsten',day:=49,month:=0,use:=1),(name:='Pfingstmontag',day:=50,month:=0,use:=1),(name:='Tag der Republik Italien',day:=2,month:=6,use:=1),(name:='Mari Himmelfahrt',day:=15,month:=8,use:=1),(name:='Allerheiligen',day:=1,month:=11,use:=1),(name:='Mari Empfngnis',day:=8,month:=12,use:=1),(name:='Heiligabend',day:=24,month:=12,use:=0),(name:='1. Weihnachtstag',day:=25,month:=12,use:=1),(name:='2. Stephanstag',day:=26,month:=12,use:=1);
END_VAR

(*
version 1.1	18. jan. 2011
programmer 	hugo
tested by		oscat

this function returns true if a date is a holiday

*)
(* @END_DECLARATION := '0' *)
y := YEAR_OF_DATE(DI);
d := EASTER(y);

FOR i := 0 TO 29 DO
	CASE L OF
	1: hd := h_de[i];
	2: hd := h_at[i];
	3: hd := h_fr[i];
	4: hd := h_bed[i];
	5: hd := h_itd[i];
	END_CASE;

	IF hd.use = 0 THEN
		CONTINUE;
	ELSIF hd.month > 0 THEN
		IF hd.use > 0 THEN
			IF hd.month = MONTH_OF_DATE(DI) AND hd.day = DAY_OF_DATE(DI) THEN
				HOLIDAY := TRUE;
				name := hd.name;
				RETURN;
			END_IF;
		ELSE
			m := DAY_OF_WEEK(SET_DATE(y,hd.month,hd.day));
			w := -hd.use;
			IF w > m THEN m := m + 7; END_IF;
			IF SET_DATE(y,hd.month,hd.day - m + w) = DI THEN
				HOLIDAY := TRUE;
				name := hd.name;
				RETURN;
			END_IF;
		END_IF;
	ELSIF d = DATE_SUB(DI,DINT_TO_TIME(SINT_TO_DINT(hd.day)*86400)) THEN
		HOLIDAY := TRUE;
		name := hd.name;
		RETURN;
	END_IF;
END_FOR;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

hm	18. jan. 2011	rev 1.1
	rewritten to be more flexible

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION HOUR : INT
VAR_INPUT
	T : TIME;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the hour of a time

*)
(* @END_DECLARATION := '0' *)
HOUR := TIME_TO_DINT(T) / 3600000;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION HOUR_OF_DT : INT
VAR_INPUT
	DI : DT;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function returns the hour of a date_and_time

*)
(* @END_DECLARATION := '0' *)
HOUR_OF_DT := DT_TO_DINT(DI) / 3600 MOD 24;


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION HOUR_TO_TIME : TIME
VAR_INPUT
	H : INT;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function converts hours to time

*)
(* @END_DECLARATION := '0' *)
HOUR_TO_TIME := DINT_TO_TIME(INT_TO_DINT(H) * 3600000);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION


(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '\/Time + Date' *)
(* @OBJECTFLAGS := '0, 8' *)
(* @SYMFILEFLAGS := '2048' *)
FUNCTION HOUR_TO_TOD : TOD
VAR_INPUT
	H : INT;
END_VAR


(*
version 1.0	22. jan. 2007
programmer 	hugo
tested by		tobias

this function converts hours to time_of_day

*)
(* @END_DECLARATION := '0' *)
HOUR_TO_TOD := DINT_TO_TOD(INT_TO_DINT(H) * 3600000);


(* revision history

hm	22. jan. 2007	rev 1.0
	original version

*)
END_FUNCTION