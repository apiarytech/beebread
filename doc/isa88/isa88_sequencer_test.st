(* Unit tests for ISA88_SEQUENCER:
   beedance -test isa88_state_machine.st isa88_phase.st isa88_sequencer.st isa88_sequencer_test.st

   The test runner scans every 10 ms of simulated time. *)

(* Records the order the steps run in. *)
FUNCTION_BLOCK SEQ_TRACE
	VAR_INPUT
		Index : INT;
	END_VAR
	VAR_OUTPUT
		Order     : ARRAY[0..39] OF INT;
		Runs      : INT;
		Ascending : BOOL; (* ran 0, 1, 2 ... *)
	END_VAR
	VAR
		last : INT := -1;
		k    : INT;
	END_VAR
	IF Index <> last AND Index >= 0 AND Runs <= 39 THEN
		Order[Runs] := Index;
		Runs := Runs + 1;
	END_IF;
	last := Index;
	Ascending := TRUE;
	FOR k := 0 TO Runs - 1 DO
		IF Order[k] <> k THEN Ascending := FALSE; END_IF;
	END_FOR;
END_FUNCTION_BLOCK

(* Sets the same times on the first n steps. *)
FUNCTION_BLOCK SEQ_TIMES
	VAR_INPUT
		n     : INT;
		start : TIME;
		run   : TIME;
		stop  : TIME;
	END_VAR
	VAR_OUTPUT
		StartTime : ARRAY[0..31] OF TIME;
		RunTime   : ARRAY[0..31] OF TIME;
		StopTime  : ARRAY[0..31] OF TIME;
	END_VAR
	VAR
		k : INT;
	END_VAR
	FOR k := 0 TO n - 1 DO
		StartTime[k] := start;
		RunTime[k] := run;
		StopTime[k] := stop;
	END_FOR;
END_FUNCTION_BLOCK

(* Three steps forward: 0, 1, 2, each Starting, Running then Stopping,
   then Complete, then Reset to Idle *)
PROGRAM TEST_SEQ_Forward
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq      : ISA88_SEQUENCER;
		tm       : SEQ_TIMES;
		tr       : SEQ_TRACE;
		scan     : INT;
		sawStart : BOOL;
		sawRun   : BOOL;
		sawStop  : BOOL;
		overlap  : BOOL;
		on       : INT;
	END_VAR
	tm(n := 3, start := T#20ms, run := T#30ms, stop := T#20ms);
	seq(Forward := scan = 0, Reset := FALSE, Count := 3,
	    StartTime := tm.StartTime, RunTime := tm.RunTime, StopTime := tm.StopTime);
	tr(Index := seq.Index);
	IF seq.Starting[0] THEN sawStart := TRUE; END_IF;
	IF seq.Running[0] THEN sawRun := TRUE; END_IF;
	IF seq.Stopping[0] THEN sawStop := TRUE; END_IF;
	on := 0;
	IF seq.Starting[0] OR seq.Running[0] OR seq.Stopping[0] THEN on := on + 1; END_IF;
	IF seq.Starting[1] OR seq.Running[1] OR seq.Stopping[1] THEN on := on + 1; END_IF;
	IF seq.Starting[2] OR seq.Running[2] OR seq.Stopping[2] THEN on := on + 1; END_IF;
	IF on > 1 THEN overlap := TRUE; END_IF;
	IF scan = 0 AND NOT (seq.Busy AND seq.Index = 0 AND seq.State = ISA88_RUNNING) THEN
		failures := failures + 1; message := 'Forward did not start step 0'; done := TRUE;
	END_IF;
	IF seq.State = ISA88_COMPLETE AND NOT done THEN
		IF tr.Runs <> 3 OR tr.Order[0] <> 0 OR tr.Order[1] <> 1 OR tr.Order[2] <> 2 THEN
			failures := failures + 1; message := 'steps not run 0, 1, 2';
		ELSIF NOT (sawStart AND sawRun AND sawStop) THEN
			failures := failures + 1; message := 'step 0 missed a phase';
		ELSIF overlap THEN
			failures := failures + 1; message := 'two steps on at once';
		ELSIF seq.Busy OR seq.Index <> -1 THEN
			failures := failures + 1; message := 'Complete but busy';
		ELSE
			seq(Reset := TRUE);
			IF seq.State <> ISA88_IDLE OR seq.StepState[2] <> ISA88_IDLE THEN
				failures := failures + 1; message := 'Reset did not go to Idle';
			END_IF;
		END_IF;
		done := TRUE;
	END_IF;
	scan := scan + 1;
	IF scan > 200 AND NOT done THEN
		failures := failures + 1; message := 'never completed'; done := TRUE;
	END_IF;
END_PROGRAM

(* Step 0's phases follow its three times: 2 scans Starting, 3 Running,
   2 Stopping *)
PROGRAM TEST_SEQ_StepTimes
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq    : ISA88_SEQUENCER;
		scan   : INT;
		nStart : INT;
		nRun   : INT;
		nStop  : INT;
	END_VAR
	seq.StartTime[0] := T#20ms;
	seq.RunTime[0] := T#30ms;
	seq.StopTime[0] := T#20ms;
	seq(Forward := scan = 0, Count := 1);
	IF seq.Starting[0] THEN nStart := nStart + 1; END_IF;
	IF seq.Running[0] THEN nRun := nRun + 1; END_IF;
	IF seq.Stopping[0] THEN nStop := nStop + 1; END_IF;
	IF seq.State = ISA88_COMPLETE THEN
		IF nStart < 2 OR nStart > 3 OR nRun <> 3 OR nStop <> 2 THEN
			failures := failures + 1;
			message := CONCAT(CONCAT(CONCAT('scans Starting/Running/Stopping: ', INT_TO_STRING(nStart)),
			           CONCAT('/', INT_TO_STRING(nRun))), CONCAT('/', INT_TO_STRING(nStop)));
		END_IF;
		done := TRUE;
	END_IF;
	scan := scan + 1;
	IF scan > 100 AND NOT done THEN
		failures := failures + 1; message := 'never completed'; done := TRUE;
	END_IF;
END_PROGRAM

(* Reverse runs 2, 1, 0 *)
PROGRAM TEST_SEQ_Reverse
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq  : ISA88_SEQUENCER;
		tr   : SEQ_TRACE;
		scan : INT;
	END_VAR
	seq(Reverse := scan = 0, Count := 3);
	tr(Index := seq.Index);
	IF seq.State = ISA88_COMPLETE THEN
		IF NOT seq.Reversing OR tr.Runs <> 3 OR tr.Order[0] <> 2 OR tr.Order[1] <> 1 OR tr.Order[2] <> 0 THEN
			failures := failures + 1; message := 'steps not run 2, 1, 0';
		END_IF;
		done := TRUE;
	END_IF;
	scan := scan + 1;
	IF scan > 200 AND NOT done THEN
		failures := failures + 1; message := 'never completed'; done := TRUE;
	END_IF;
END_PROGRAM

(* Reverse while step 2 of 5 runs: 0, 1, 2, then back 1, 0 *)
PROGRAM TEST_SEQ_ChangeDirection
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq  : ISA88_SEQUENCER;
		tm   : SEQ_TIMES;
		tr   : SEQ_TRACE;
		scan : INT;
		sent : BOOL;
		rev  : BOOL;
	END_VAR
	tm(n := 5, run := T#30ms);
	rev := seq.Index = 2 AND NOT sent;
	IF rev THEN sent := TRUE; END_IF;
	seq(Forward := scan = 0, Reverse := rev, Count := 5, RunTime := tm.RunTime);
	tr(Index := seq.Index);
	IF seq.State = ISA88_COMPLETE THEN
		IF tr.Runs <> 5 OR tr.Order[0] <> 0 OR tr.Order[1] <> 1 OR tr.Order[2] <> 2
		   OR tr.Order[3] <> 1 OR tr.Order[4] <> 0 THEN
			failures := failures + 1; message := 'steps not run 0, 1, 2, 1, 0';
		END_IF;
		done := TRUE;
	END_IF;
	scan := scan + 1;
	IF scan > 300 AND NOT done THEN
		failures := failures + 1; message := 'never completed'; done := TRUE;
	END_IF;
END_PROGRAM

(* No steps: Complete at once; Reset back to Idle *)
PROGRAM TEST_SEQ_Empty
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq : ISA88_SEQUENCER;
		ok  : BOOL;
	END_VAR
	seq(Count := 0);
	ok := seq.State = ISA88_IDLE AND NOT seq.Busy;
	seq(Forward := TRUE);
	ok := ok AND seq.State = ISA88_COMPLETE AND seq.Index = -1 AND NOT seq.Busy;
	seq(Forward := FALSE, Reset := TRUE);
	ok := ok AND seq.State = ISA88_IDLE;
	IF NOT ok THEN
		failures := 1; message := 'empty sequence';
	END_IF;
	done := TRUE;
END_PROGRAM

(* All 32 steps, forward, with no times: each runs once, in order *)
PROGRAM TEST_SEQ_ThirtyTwo
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq  : ISA88_SEQUENCER;
		tr   : SEQ_TRACE;
		scan : INT;
	END_VAR
	seq(Forward := scan = 0, Count := 40); (* more than 32 is 32 *)
	tr(Index := seq.Index);
	IF seq.State = ISA88_COMPLETE THEN
		IF tr.Runs <> 32 OR NOT tr.Ascending THEN
			failures := failures + 1; message := 'steps not run 0 to 31';
		END_IF;
		done := TRUE;
	END_IF;
	scan := scan + 1;
	IF scan > 1000 AND NOT done THEN
		failures := failures + 1; message := 'never completed'; done := TRUE;
	END_IF;
END_PROGRAM

(* Hold while step 1 runs: Held, its time stands still; Restart carries on
   and the sequence completes *)
PROGRAM TEST_SEQ_HoldRestart
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq   : ISA88_SEQUENCER;
		tm    : SEQ_TIMES;
		tr    : SEQ_TRACE;
		scan  : INT;
		phase : INT;
		wait  : INT;
	END_VAR
	tm(n := 3, run := T#50ms);
	CASE phase OF
		0: seq(Forward := scan = 0, Count := 3, RunTime := tm.RunTime);
		   IF seq.Index = 1 AND seq.Running[1] THEN
		       seq(Hold := TRUE);
		       IF seq.State <> ISA88_HOLDING THEN
		           failures := failures + 1; message := 'Hold did not go to Holding'; done := TRUE;
		       END_IF;
		       phase := 1;
		   END_IF;
		1: seq(Hold := FALSE);
		   IF seq.State = ISA88_HELD THEN
		       wait := wait + 1;
		       IF seq.Running[1] OR seq.Index <> 1 THEN
		           failures := failures + 1; message := 'Held but running'; done := TRUE;
		       ELSIF wait = 10 THEN
		           seq(Restart := TRUE);
		           IF seq.State <> ISA88_RESTARTING THEN
		               failures := failures + 1; message := 'Restart did not go to Restarting'; done := TRUE;
		           END_IF;
		           phase := 2;
		       END_IF;
		   END_IF;
		2: seq(Restart := FALSE);
		   IF seq.State = ISA88_COMPLETE THEN
		       IF tr.Runs <> 3 THEN
		           failures := failures + 1; message := 'did not carry on to step 2';
		       END_IF;
		       done := TRUE;
		   END_IF;
	END_CASE
	tr(Index := seq.Index);
	scan := scan + 1;
	IF scan > 300 AND NOT done THEN
		failures := failures + 1; message := 'never completed'; done := TRUE;
	END_IF;
END_PROGRAM

(* Stop while step 0 runs: Stopping runs StopTime, then Stopped; step 1
   never starts; Reset makes every step Idle *)
PROGRAM TEST_SEQ_Stop
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq   : ISA88_SEQUENCER;
		scan  : INT;
		phase : INT;
		nStop : INT;
	END_VAR
	seq.RunTime[0] := T#100ms;
	seq.StopTime[0] := T#30ms;
	CASE phase OF
		0: seq(Forward := scan = 0, Count := 2);
		   IF scan = 3 THEN
		       seq(Stop := TRUE);
		       phase := 1;
		   END_IF;
		1: seq(Stop := FALSE);
		   IF seq.Stopping[0] THEN nStop := nStop + 1; END_IF;
		   IF seq.State = ISA88_STOPPED THEN
		       IF nStop <> 3 OR seq.StepState[1] <> ISA88_IDLE OR seq.Busy THEN
		           failures := failures + 1; message := CONCAT('wrong stop, scans Stopping: ', INT_TO_STRING(nStop));
		       END_IF;
		       seq(Reset := TRUE);
		       IF seq.State <> ISA88_IDLE OR seq.StepState[0] <> ISA88_IDLE THEN
		           failures := failures + 1; message := 'Reset did not go to Idle';
		       END_IF;
		       done := TRUE;
		   END_IF;
	END_CASE
	scan := scan + 1;
	IF scan > 100 AND NOT done THEN
		failures := failures + 1; message := 'never stopped'; done := TRUE;
	END_IF;
END_PROGRAM

(* Stop and Abort given between two steps are not lost *)
PROGRAM TEST_SEQ_BetweenSteps
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		a, b  : ISA88_SEQUENCER;
		scan  : INT;
		gapA  : BOOL;
		gapB  : BOOL;
		sentA : BOOL;
		sentB : BOOL;
	END_VAR
	(* the scan step 0 reports Complete, the chart is between steps *)
	gapA := a.StepState[0] = ISA88_COMPLETE AND NOT sentA;
	IF gapA THEN sentA := TRUE; END_IF;
	gapB := b.StepState[0] = ISA88_COMPLETE AND NOT sentB;
	IF gapB THEN sentB := TRUE; END_IF;
	a(Forward := scan = 0, Stop := gapA, Count := 3);
	b(Forward := scan = 0, Abort := gapB, Count := 3);
	IF scan = 40 THEN
		IF a.State <> ISA88_STOPPED OR a.StepState[1] = ISA88_RUNNING THEN
			failures := failures + 1; message := 'Stop between steps was lost';
		ELSIF b.State <> ISA88_ABORTED THEN
			failures := failures + 1; message := 'Abort between steps was lost';
		ELSIF NOT (sentA AND sentB) THEN
			failures := failures + 1; message := 'never between steps';
		END_IF;
		done := TRUE;
	END_IF;
	scan := scan + 1;
END_PROGRAM

(* Abort while a step runs: Aborted; Reset to Idle *)
PROGRAM TEST_SEQ_Abort
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		seq : ISA88_SEQUENCER;
		ok  : BOOL;
	END_VAR
	seq.RunTime[0] := T#1s;
	seq(Forward := TRUE, Count := 4);
	ok := seq.State = ISA88_RUNNING AND seq.Index = 0 AND NOT seq.Running[0]; (* starts next scan *)
	seq(Forward := FALSE);
	ok := ok AND seq.State = ISA88_RUNNING AND NOT seq.Starting[0] AND seq.Running[0];
	seq(Abort := TRUE);
	ok := ok AND seq.State = ISA88_ABORTING AND NOT seq.Running[0];
	seq(Abort := FALSE);
	ok := ok AND seq.State = ISA88_ABORTED;
	seq();
	ok := ok AND seq.State = ISA88_ABORTED AND NOT seq.Busy;
	seq(Reset := TRUE);
	ok := ok AND seq.State = ISA88_IDLE AND seq.StepState[0] = ISA88_IDLE;
	IF NOT ok THEN
		failures := 1; message := 'abort';
	END_IF;
	done := TRUE;
END_PROGRAM
