(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE CALENDAR :
STRUCT
	UTC : DT;					(* world time UTC *)
	LOCAL_DT : DT;					(* local time *)
	LOCAL_DATE : DATE;				(* local date *)
	LOCAL_TOD : TOD;					(* local time of day *)
	YEAR : INT;					(* year of LDATE *)
	MONTH : INT;				(* month of LDATE *)
	DAY : INT;					(* day of LDATE *)
	WEEKDAY : INT;				(* weekday of LDATE *)
	OFFSET : INT;				(* Time Zone Offset for Local time in minutes *)
	DST_EN : BOOL;				(* daylight savings time enable *)
	DST_ON : BOOL;				(* true when daylight savings time os on *)
	NAME : STRING(5);			(* name of time zone *)
	LANGUAGE : INT;			(* location number pls see location setup *)
	LONGITUDE : REAL;			(* longitude of current location *)
	LATITUDE : REAL;			(* latitude of current location *)
	SUN_RISE : TOD;				(* sun_rise for current location *)
	SUN_SET : TOD;				(* sun_set for current location *)
	SUN_MIDDAY : TOD;			(* worldtime when sun stands at south position *)
	SUN_HEIGTH : REAL	;		(* suns heigth at midday, south position *)
	SUN_HOR : REAL;			(* sun angle horizontal 0 = north in degrees *)
	SUN_VER : REAL;			(* sun angle vertical above horizon in degrees *)
	NIGHT : BOOL;				(* true between sun_set and sun_rise *)
	HOLIDAY : BOOL;			(* true when holiday *)
	HOLY_NAME : STRING(30);	(* name of holiday *)
	WORK_WEEK : INT;			(* current work week *)
END_STRUCT
END_TYPE

(* revision history
hm	21. nov. 2008	rev 1.0
	original version

hm	8. feb. 2009	rev 1.1
	added sun position data

hm	10. mar. 2009	rev 1.2
	added work_week, sun_midday, sun_heigth

hm	23. jan 2010	rev 1.3
	sun_rise, sun_set and sun_midday are now calculated in local time

hm 3.2.2021    rev 1.4
   changed LDT >> LOCAL_DT, LDATE >> LOCAL_DATE, LTOD >> LOCAL_TOD
   reserved Words in IEC61131-3

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE COMPLEX :
STRUCT
	re : REAL;
	im : REAL;
END_STRUCT
END_TYPE


(* revision history
hm	18. oct. 2008
	original version

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE CONSTANTS_LANGUAGE :
STRUCT
	(* Language Setup *)
	DEFAULT : INT := 1; (* 1=english, 2=german 3=french *)
	LMAX : INT := 3;
	WEEKDAYS : ARRAY[1..3, 1..7] OF STRING(10) :=	'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday',
													'Montag', 'Dienstag', 'Mittwoch', 'Donnerstag', 'Freitag', 'Samstag', 'Sonntag',
													'Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche';
	WEEKDAYS2 : ARRAY[1..3, 1..7] OF STRING(2) :=	'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su',
													'Mo', 'Di', 'Mi', 'Do', 'Fr', 'Sa', 'So',
													'Lu', 'Ma', 'Me', 'Je', 'Ve', 'Sa', 'Di';
	MONTHS : ARRAY[1..3, 1..12] OF STRING(10) :=	'January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December',
													'Januar', 'Februar', 'Mrz', 'April', 'Mai', 'Juni', 'Juli', 'August', 'September', 'Oktober', 'November', 'Dezember',
													'Janvier', 'Fvrier', 'mars', 'Avril', 'Mai', 'Juin', 'Juillet', 'Aot', 'Septembre', 'Octobre', 'Novembre', 'Decembre';
	MONTHS3 : ARRAY[1..3, 1..12] OF STRING(3) :=	'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec',
													'Jan', 'Feb', 'Mrz', 'Apr', 'Mai', 'Jun', 'Jul', 'Aug', 'Sep', 'Okt', 'Nov', 'Dez',
													'Jan', 'Fev', 'Mar', 'Avr', 'Mai', 'Jun', 'Jul', 'Aou', 'Sep', 'Oct', 'Nov', 'Dec';
	DIRS : ARRAY[1..3,0..15] OF STRING(3) :=		'N', 'NNE', 'NE', 'ENE', 'E', 'ESE', 'SE', 'SSE', 'S', 'SSW', 'SW', 'WSW', 'W', 'WNW', 'NW', 'NNW',
													'N', 'NNO', 'NO', 'ONO', 'O', 'OSO', 'SO', 'SSO', 'S', 'SSW', 'SW', 'WSW', 'W', 'WNW', 'NW', 'NNW',
													'N', 'NNO', 'NO', 'ONO', 'O', 'OSO', 'SO', 'SSO', 'S', 'SSW', 'SW', 'WSW', 'W', 'WNW', 'NW', 'NNW';
END_STRUCT
END_TYPE

(* revision history
hm	19. oct. 2008	rev 1.0
	original version

hm	22. oct. 2008	rev 1.1
	added directions DIRS

hm	23. dec. 2008	rev 1.2
	added french language

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE CONSTANTS_LOCATION :
STRUCT
	(* location setup *)
	DEFAULT : INT := 1; (* 1=germany, 2=austria 3=france 4=belgium-german 5= italien-Sdtirol *)
	LMAX : INT := 5;

	(* language spoken in the location *)
	LANGUAGE : ARRAY[1..5] OF INT := 2, 2, 3, 2, 2;

END_STRUCT
END_TYPE


(* revision history
hm	20. oct. 2008	rev 1.0
	original version

hm	23. dec. 2008	rev 1.1
	added french holidays

hm	18. jan 2011	rev 1.2
	added spoken language
	deleted holidays from structure for more flexibility

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE CONSTANTS_MATH :
STRUCT
	PI : REAL := 	3.14159265358979323846264338327950288;		(* Kreiszahl PI *)
	PI2 : REAL := 	6.28318530717958647692528676655900576;	(* PI * 2 *)
	PI4 : REAL :=  12.56637061435917295385057353311801152;	(* PI * 4 *)
	PI05 : REAL :=  1.5707963267949;							(* PI / 2 *)
	PI025 : REAL := 0.785398163397448;							(* PI / 4 *)
	PI_INV : REAL := 0.318309886183791;						(* 1 / PI *)
	E : REAL := 	2.71828182845904523536028747135266249;		(* Euler constant e *)
	E_INV : REAL := 0.367879441171442;							(* 1 / e *)
	SQ2 : REAL :=	1.4142135623731;							(* Wurzel von 2 *)
	FACTS : ARRAY[0..12] OF DINT := 1,1,2,6,24,120,720,5040,40320,362880,3628800,39916800,479001600;
END_STRUCT
END_TYPE



(* revision history
hm	26.10.2008	rev 1.0
	original version

hm	26. mar. 2011	rev 1.1
	added array facts
*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE CONSTANTS_PHYS :
STRUCT
	C : REAL := 299792458.0;				(* Lichtgeschwindigkeit in m/s *)
	E : REAL := 1.60217653E-19;				(* elementarladung in Colomb = A * s *)
	G : REAL := 9.80665; 						(* Erdbeschleunigung in m / s *)
	T0 : REAL := -273.15;						(* absoluter Nullpunkt in C *)
	RU : REAL := 8.314472; 					(* Universelle Gaskonstante in J / (mol  K) *)
	PN : REAL := 101325.0;					(* NormalDruck in Pa *)
END_STRUCT
END_TYPE



(* revision history
hm	18.10.2008	rev 1.0
	original version

hm	13. mar. 2009	rev 1.1
	real constants updated to new systax using dot

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE CONSTANTS_SETUP :
STRUCT
	(* setup Parameters *)
	EXTENDED_ASCII : BOOL := TRUE;
	CHARNAMES : ARRAY[1..4] OF STRING(253) := ';"&quot;&&amp;<&lt;>&gt;&euro;&nbsp;&iexcl;&cent;&pound;&curren;&yen;&brvbar;&sect;&uml;&copy;&ordf;&laquo;&not;&shy;&reg;&macr;&deg;&plusmn;&sup2;&sup3;&acute;&micro;&para;&middot;&cedil;&sup1;&ordm;&raquo;&frac14;&Ucirc;',
		';&frac34;&iquest;&Agrave;&Aacute;&Acirc;&Atilde;&Auml;&Aring;&AElig;&Ccedil;&Egrave;&Eacute;&Ecirc;&Euml;&Igrave;&Iacute;&Icirc;&Iuml;&ETH;&Ntilde;&Ograve;&Oacute;&Ocirc;&Otilde;&Ouml;&times;&Oslash;&Ugrave;&Uacute;&frac12;',
		';&Uuml;&Yacute;&THORN;&szlig;&agrave;&aacute;&acirc;&atilde;&auml;&aring;&aelig;&ccedil;&egrave;&eacute;&ecirc;&euml;&igrave;&iacute;&icirc;&iuml;&eth;&ntilde;&ograve;&oacute;&ocirc;&otilde;&ouml;&divide;&oslash;&ugrave;',
		';&uacute;&ucirc;&uuml;&yacute;&thorn;&yuml;';
	MTH_OFS : ARRAY[1..12] OF INT := 0,31,59,90,120,151,181,212,243,273,304,334;
	DECADES : ARRAY[0..8] OF REAL := 1.0,10.0,100.0,1000.0,10000.0,100000.0,1000000.0,10000000.0,100000000.0;
END_STRUCT
END_TYPE

(* revision history
hm	24.10.2008	rev 1.0
	original version

hm	3. feb 2021	rev 1.1
	corrected error in decades
	thanks to Valentin Marte

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE ESR_DATA :
STRUCT
	TYP : BYTE;
	ADRESS : STRING(10);
	DS : DT;
	TS : TIME;
	DATA : ARRAY[0..7] OF BYTE;
END_STRUCT
END_TYPE


(* revision history
hm	18. oct. 2008
	original version

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE FRACTION :
STRUCT
	NUMERATOR : INT;
	DENOMINATOR : INT;
END_STRUCT
END_TYPE
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE HOLIDAY_DATA :
STRUCT
	NAME : STRING(30);
	DAY : SINT;
	MONTH : SINT;
	USE : SINT;
END_STRUCT
END_TYPE


(* revision history
hm	18. oct. 2008
	original version

*)

(*	if month = 0 then F_day is the offset in days from easter
	used :0 = NOT used, 1 = used, -1..-7 menas the weekday before the specified DATE example -3 means wednesday before 23.11
*)

(* samples for Holiday data

	HOLIDAY_DE : ARRAY[0..29] OF HOLIDAY_DATA := (name := 'Neujahr', day := 1, month := 1, use := 1),
		(name := 'Heilig Drei Knige', day := 6, month := 1, use := 1),
		(name := 'Karfreitag', day := -2, month := 0, use := 1),
		(name := 'Ostersonntag', day := 0, month := 0, use := 1),
		(name := 'Ostermontag', day := 1, month := 0, use := 1),
		(name := 'Tag der Arbeit', day := 1, month := 5, use := 1),
		(name := 'Christi Himmelfahrt', day := 39, month := 0, use := 1),
		(name := 'Pfingstsonntag', day := 49, month := 0, use := 1),
		(name := 'Pfingstmontag', day := 50, month := 0, use := 1),
		(name := 'Fronleichnam', day := 60, month := 0, use := 1),
		(name := 'Augsburger Friedensfest', day := 8, month := 8, use := 0),
		(name := 'Maria Himmelfahrt', day := 15, month := 8, use := 1),
		(name := 'Tag der Deutschen Einheit', day := 3, month := 10, use := 1),
		(name := 'Reformationstag', day := 31, month := 10, use := 0),
		(name := 'Allerheiligen', day := 1, month := 11, use := 1),
		(name := 'Buss und Bettag', day := 23, month := 11, use := 0),
		(name := '1. Weihnachtstag', day := 25, month := 12, use := 1),
		(name := '2. Weihnachtstag', day := 26, month := 12, use := 1);


	HOLIDAY_AT : ARRAY[0..29] OF HOLIDAY_DATA := (name := 'Neujahr', day := 1, month := 1, use := 1),
		(name := 'Heilig Drei Knige', day := 6, month := 1, use := 1),
		(name := 'Karfreitag', day := -2, month := 0, use := 1),
		(name := 'Ostersonntag', day := 0, month := 0, use := 1),
		(name := 'Ostermontag', day := 1, month := 0, use := 1),
		(name := '', day := 1, month := 5, use := 0),
		(name := 'Christi Himmelfahrt', day := 39, month := 0, use := 1),
		(name := 'Pfingstsonntag', day := 49, month := 0, use := 1),
		(name := 'Pfingstmontag', day := 50, month := 0, use := 1),
		(name := 'Fronleichnam', day := 60, month := 0, use := 1),
		(name := '', day := 8, month := 8, use := 0),
		(name := 'Maria Himmelfahrt', day := 15, month := 8, use := 1),
		(name := '', day := 3, month := 10, use := 0),
		(name := '', day := 31, month := 10, use := 0),
		(name := 'Allerheiligen', day := 1, month := 11, use := 1),
		(name := 'Maria Empfngnis', day := 8, month := 12, use := 1),
		(name := '1. Weihnachtstag', day := 25, month := 12, use := 1),
		(name := '2. Weihnachtstag', day := 26, month := 12, use := 1);

	HOLIDAY_FR : ARRAY[0..29] OF HOLIDAY_DATA := (name := 'Nouvel an', day := 1, month := 1, use := 1),
		(name := 'St Valentin', day := 14, month := 2, use := 0),
		(name := 'Vendredi Saint (alsace)', day := -2, month := 0, use := 0),
		(name := 'Dimanche de pques', day := 0, month := 0, use := 1),
		(name := 'Lundi de pques', day := 1, month := 0, use := 1),
		(name := 'Jeudi de Ascension', day := 39, month := 0, use := 1),
		(name := 'dimanche de Pentecte ', day := 49, month := 0, use := 1),
		(name := 'jeudi de la Trinit', day := 60, month := 0, use := 0),
		(name := 'Fte du travail', day := 1, month := 5, use := 1),
		(name := 'Victoire 1945 ', day := 8, month := 5, use := 1),
		(name := 'Prise de La bastille', day := 14, month := 7, use := 1),
		(name := '15 Aot 1944', day := 15, month := 8, use := 1),
		(name := 'Halloween', day := 31, month := 10, use := 0),
		(name := 'Armistice 1918', day := 11, month := 11, use := 1),
		(name := 'Nol', day := 25, month := 12, use := 1),
		(name := 'Saint tienne (alsace)', day := 26, month := 12, use := 0),
		(name := 'Fte de la musique', day := 21, month := 6, use := 0);


	HOLIDAY_BED : ARRAY[0..29] OF HOLIDAY_DATA :=	(name := 'Neujahr', day := 1, month := 1, use := 1),
		(name := 'Ostersonntag', day := 0, month := 0, use := 1),
		(name := 'Ostermontag', day := 1, month := 0, use := 1),
		(name := 'Tag der Arbeit', day := 1, month := 5, use := 1),
		(name := 'Christi Himmelfahrt', day := 39, month := 0, use := 1),
		(name := 'Pfingsten', day := 49, month := 0, use := 1),
		(name := 'Pfingstmontag', day := 50, month := 0, use := 1),
		(name := 'Nationalfeiertag', day := 21, month := 7, use := 1),
		(name := 'Mari Himmelfahrt', day := 15, month := 8, use := 1),
		(name := 'Allerheiligen', day := 1, month := 11, use := 1),
		(name := 'Feiertag DG', day := 15, month := 11, use := 1),
		(name := 'Heiligabend', day := 24, month := 12, use := 1),
		(name := '1. Weihnachtstag', day := 25, month := 12, use := 1),
		(name := '2. Weihnachtstag', day := 26, month := 12, use := 1),
		(name := 'Silvester', day := 31, month := 12, use := 1);


	HOLIDAY_ITD : ARRAY[0..29] OF HOLIDAY_DATA := (name := 'Neujahr', day := 1, month := 1, use := 1),
		(name := 'Heilig Drei Knige', day := 6, month := 1, use := 1),
		(name := 'Ostersonntag', day := 0, month := 0, use := 1),
		(name := 'Ostermontag', day := 1, month := 0, use := 1),
		(name := 'Tag der Befeiung Italiens', day := 25, month := 4, use := 1),
		(name := 'Tag der Arbeit', day := 1, month := 5, use := 1),
		(name := 'Pfingsten', day := 49, month := 0, use := 1),
		(name := 'Pfingstmontag', day := 50, month := 0, use := 1),
		(name := 'Tag der Republik Italien', day := 2, month := 6, use := 1),
		(name := 'Mari Himmelfahrt', day := 15, month := 8, use := 1),
		(name := 'Allerheiligen', day := 1, month := 11, use := 1),
		(name := 'Mari Empfngnis', day := 8, month := 12, use := 1),
		(name := 'Heiligabend', day := 24, month := 12, use := 0),
		(name := '1. Weihnachtstag', day := 25, month := 12, use := 1),
		(name := '2. Stephanstag', day := 26, month := 12, use := 1);

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE REAL2 :
STRUCT
	R1 : REAL;	(* small value *)
	RX : REAL;	(* big value *)
END_STRUCT
END_TYPE


(* revision history
hm	18. oct. 2008
	original version

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE SDT :
STRUCT
	YEAR : INT;
	MONTH : INT;
	DAY : INT;
	WEEKDAY : INT;
	HOUR : INT;
	MINUTE : INT;
	SECOND : INT;
	MS : INT;
END_STRUCT
END_TYPE


(* revision history
hm	18. oct. 2008
	original version

*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE TIMER_EVENT :
STRUCT
	TYP : BYTE;
	CHANNEL : BYTE;
	DAY : BYTE;
	START : TOD;
	DURATION : TIME;
	LAND : BYTE;
	LOR : BYTE;
	LAST : DT;
END_STRUCT
END_TYPE


(* revision history
hm	18. oct. 2008	rev 1.0
	original version

hm	26. jan. 2011	rev 1.1
	changed type of last to be DT
*)
(* @END_DECLARATION := '0' *)

(* @NESTEDCOMMENTS := 'Yes' *)
(* @PATH := '' *)
(* @OBJECTFLAGS := '0, 8' *)
TYPE VECTOR_3 :
STRUCT
	X : REAL;
	Y : REAL;
	Z : REAL;
END_STRUCT
END_TYPE


(* revision history
hm	18. oct. 2008
	original version

*)
(* @END_DECLARATION := '0' *)