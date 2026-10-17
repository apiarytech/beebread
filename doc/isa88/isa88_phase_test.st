(* Unit tests for ISA88_PHASE:
   beedance -test isa88_state_machine.st isa88_phase.st isa88_phase_test.st *)

(* A heating phase written with the template: Running heats one degree a
   scan up to 5, Restarting takes 30 ms, the other transient states finish
   at once. Temp is kept across Hold and Restart, as a restart point. *)
FUNCTION_BLOCK PH_DEMO_HEAT
	VAR_INPUT
		Start   : BOOL;
		Hold    : BOOL;
		Restart : BOOL;
		Stop    : BOOL;
		Reset   : BOOL;
		Fault   : DINT;
	END_VAR
	VAR_OUTPUT
		Heater  : BOOL;
		Temp    : INT;
		State   : INT;
		Entry   : BOOL;
		Failure : DINT;
	END_VAR
	VAR
		ph : ISA88_PHASE;
	END_VAR
	ph(Start := Start, Hold := Hold, Restart := Restart, Stop := Stop, Reset := Reset, FailCode := Fault);
	CASE ph.State OF
		ISA88_IDLE:
			IF ph.Entry THEN Temp := 0; END_IF;
		ISA88_RUNNING:
			Heater := TRUE;
			Temp := Temp + 1;
			IF Temp >= 5 THEN ph.Done(); END_IF;
		ISA88_COMPLETE:
			Heater := FALSE;
		ISA88_HOLDING, ISA88_STOPPING, ISA88_ABORTING:
			Heater := FALSE;
			ph.Done();
		ISA88_RESTARTING:
			IF ph.StateTime >= T#30ms THEN ph.Done(); END_IF;
	END_CASE
	State := ph.State;
	Entry := ph.Entry;
	Failure := ph.Failure;
END_FUNCTION_BLOCK

(* Idle -> Running until the logic is done -> Complete -> Idle *)
PROGRAM TEST_PHASE_Complete
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		h    : PH_DEMO_HEAT;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: h(); ok := h.State = ISA88_IDLE AND h.Entry;
		1: h(Start := TRUE); ok := h.State = ISA88_RUNNING AND h.Entry AND h.Heater AND h.Temp = 1;
		2: h(Start := FALSE); ok := h.State = ISA88_RUNNING AND NOT h.Entry AND h.Temp = 2;
		3: h(); ok := h.Temp = 3;
		4: h(); ok := h.Temp = 4;
		5: h(); ok := h.State = ISA88_RUNNING AND h.Temp = 5; (* logic sets Done *)
		6: h(); ok := h.State = ISA88_COMPLETE AND h.Entry AND NOT h.Heater;
		7: h(); ok := h.State = ISA88_COMPLETE AND NOT h.Entry;
		8: h(Reset := TRUE); ok := h.State = ISA88_IDLE AND h.Temp = 0;
		   done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM

(* Done does not carry into the next state; Restarting waits on StateTime;
   Running resumes where it was held *)
PROGRAM TEST_PHASE_HoldRestart
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		h     : PH_DEMO_HEAT;
		step  : INT;
		ok    : BOOL;
		next  : BOOL;
		waits : INT;
	END_VAR
	ok := TRUE;
	next := TRUE;
	CASE step OF
		0: h(Start := TRUE); ok := h.State = ISA88_RUNNING AND h.Temp = 1;
		1: h(Start := FALSE, Hold := TRUE); ok := h.State = ISA88_HOLDING AND NOT h.Heater;
		2: h(Hold := FALSE); ok := h.State = ISA88_HELD;
		3: h(); ok := h.State = ISA88_HELD; (* Holding's Done was cleared *)
		4: h(Restart := TRUE); ok := h.State = ISA88_RESTARTING;
		5: h(Restart := FALSE); ok := h.State = ISA88_RESTARTING;
		6: h(); ok := h.State = ISA88_RESTARTING;
		7: h();
		   IF h.State = ISA88_RESTARTING THEN
		       waits := waits + 1;
		       ok := waits < 10;
		       next := FALSE;
		   ELSE
		       ok := h.State = ISA88_RUNNING AND h.Temp = 2;
		       done := TRUE;
		   END_IF;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
		done := TRUE;
	END_IF;
	IF next THEN step := step + 1; END_IF;
END_PROGRAM

(* A failure holds the phase, Restart holds again while it lasts, Failure
   keeps the first code until Idle *)
PROGRAM TEST_PHASE_Failure
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		h    : PH_DEMO_HEAT;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: h(Start := TRUE); ok := h.State = ISA88_RUNNING AND h.Failure = 0;
		1: h(Start := FALSE, Fault := 42); ok := h.State = ISA88_HOLDING AND h.Failure = 42;
		2: h(Fault := 99); ok := h.State = ISA88_HELD AND h.Failure = 42;
		3: h(Restart := TRUE); ok := h.State = ISA88_RESTARTING;
		4: h(Restart := FALSE); ok := h.State = ISA88_HOLDING; (* failure still there *)
		5: h(Fault := 0); ok := h.State = ISA88_HELD AND h.Failure = 42;
		6: h(Restart := TRUE); ok := h.State = ISA88_RESTARTING;
		7: h(Restart := FALSE, Stop := TRUE); ok := h.State = ISA88_STOPPING;
		8: h(Stop := FALSE); ok := h.State = ISA88_STOPPED AND h.Failure = 42;
		9: h(Reset := TRUE); ok := h.State = ISA88_IDLE AND h.Failure = 0;
		   done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM
