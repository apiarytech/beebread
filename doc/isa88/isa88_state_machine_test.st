(* Unit tests for ISA88_STATE_MACHINE:
   beedance -test isa88_state_machine.st isa88_state_machine_test.st *)

(* Idle -> Running -> Complete -> Idle *)
PROGRAM TEST_ISA88_Complete
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		sm   : ISA88_STATE_MACHINE;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: sm(); ok := ok AND sm.State = 1 AND sm.IsQuiescent;
		1: sm(Start := TRUE); ok := ok AND sm.State = 2 AND sm.IsTransient;
		2: sm(Start := FALSE); ok := ok AND sm.State = 2;
		3: sm(SC := TRUE); ok := ok AND sm.State = 3 AND sm.IsFinal;
		4: sm(SC := FALSE, Reset := TRUE); ok := ok AND sm.State = 1;
		5: sm(Reset := FALSE); ok := ok AND sm.State = 1;
		   done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM

(* Running -> Pausing -> Paused -> Running -> Holding -> Held -> Restarting -> Running *)
PROGRAM TEST_ISA88_PauseHold
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		sm   : ISA88_STATE_MACHINE;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: sm(Start := TRUE); ok := ok AND sm.State = 2;
		1: sm(Start := FALSE, Pause := TRUE); ok := ok AND sm.State = 4;
		2: sm(Pause := FALSE, SC := TRUE); ok := ok AND sm.State = 5 AND sm.IsQuiescent;
		3: sm(SC := FALSE, Resume := TRUE); ok := ok AND sm.State = 2;
		4: sm(Resume := FALSE, Hold := TRUE); ok := ok AND sm.State = 6;
		5: sm(Hold := FALSE, SC := TRUE); ok := ok AND sm.State = 7 AND sm.IsQuiescent;
		6: sm(SC := FALSE, Restart := TRUE); ok := ok AND sm.State = 8;
		7: sm(Restart := FALSE, Hold := TRUE); ok := ok AND sm.State = 6; (* Hold while Restarting *)
		8: sm(Hold := FALSE, SC := TRUE); ok := ok AND sm.State = 7;
		9: sm(SC := FALSE, Restart := TRUE); ok := ok AND sm.State = 8;
		10: sm(Restart := FALSE, SC := TRUE); ok := ok AND sm.State = 2;
		11: sm(SC := FALSE, Pause := TRUE); ok := ok AND sm.State = 4;
		12: sm(Pause := FALSE, Hold := TRUE); ok := ok AND sm.State = 6; (* Hold while Pausing *)
		    done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM

(* Stop and Abort paths, command priority, and commands a state ignores *)
PROGRAM TEST_ISA88_StopAbort
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		sm   : ISA88_STATE_MACHINE;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: sm(Stop := TRUE, Abort := TRUE, Hold := TRUE); ok := ok AND sm.State = 1; (* Idle ignores them *)
		1: sm(Stop := FALSE, Abort := FALSE, Hold := FALSE, Start := TRUE); ok := ok AND sm.State = 2;
		2: sm(Start := FALSE, Hold := TRUE, Stop := TRUE); ok := ok AND sm.State = 9; (* Stop beats Hold *)
		3: sm(Hold := FALSE, Stop := FALSE, SC := TRUE); ok := ok AND sm.State = 10 AND sm.IsFinal;
		4: sm(SC := FALSE, Abort := TRUE); ok := ok AND sm.State = 11; (* Abort from Stopped *)
		5: sm(Abort := FALSE, Reset := TRUE); ok := ok AND sm.State = 11; (* Aborting ignores Reset *)
		6: sm(Reset := FALSE, SC := TRUE); ok := ok AND sm.State = 12 AND sm.IsFinal;
		7: sm(SC := FALSE, Start := TRUE); ok := ok AND sm.State = 12; (* Aborted ignores Start *)
		8: sm(Start := FALSE, Reset := TRUE); ok := ok AND sm.State = 1;
		9: sm(Reset := FALSE, Start := TRUE); ok := ok AND sm.State = 2;
		10: sm(Start := FALSE, Stop := TRUE, Abort := TRUE); ok := ok AND sm.State = 11; (* Abort beats Stop *)
		11: sm(Stop := FALSE, Abort := FALSE, SC := TRUE); ok := ok AND sm.State = 12;
		12: sm(SC := FALSE, Reset := TRUE); ok := ok AND sm.State = 1;
		    done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM
