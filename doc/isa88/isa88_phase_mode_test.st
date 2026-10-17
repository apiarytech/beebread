(* Unit tests for ISA88_PHASE modes:
   beedance -test isa88_state_machine.st isa88_phase.st isa88_phase_mode_test.st *)

(* A phase whose Running has steps 1 to 3. Each step's condition is Ready;
   after step 3 the phase is done. *)
FUNCTION_BLOCK PH_DEMO_STEPS
	VAR_INPUT
		Start    : BOOL;
		Reset    : BOOL;
		Mode     : INT := 1;
		Proceed  : BOOL;
		Jump     : BOOL;
		JumpStep : INT;
		Ready    : BOOL;
	END_VAR
	VAR_OUTPUT
		State      : INT;
		Step       : INT := 1;
		Waiting    : BOOL;
		ActiveMode : INT;
	END_VAR
	VAR
		ph : ISA88_PHASE;
	END_VAR
	ph(Start := Start, Reset := Reset, Mode := Mode, Proceed := Proceed, Jump := Jump, JumpStep := JumpStep);
	CASE ph.State OF
		ISA88_IDLE:
			IF ph.Entry THEN Step := 1; END_IF;
		ISA88_RUNNING:
			IF ph.JumpRequest THEN
				Step := ph.JumpTarget;
			ELSIF ph.Advance(Ready) THEN
				IF Step >= 3 THEN ph.Done(); ELSE Step := Step + 1; END_IF;
			END_IF;
	END_CASE
	State := ph.State;
	Waiting := ph.Waiting;
	ActiveMode := ph.ActiveMode;
END_FUNCTION_BLOCK

(* Automatic: a step is taken as soon as its condition is met; Proceed and
   Jump do nothing *)
PROGRAM TEST_MODE_Automatic
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		p    : PH_DEMO_STEPS;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: p(); ok := p.State = ISA88_IDLE AND p.ActiveMode = ISA88_AUTOMATIC;
		1: p(Start := TRUE); ok := p.State = ISA88_RUNNING AND p.Step = 1 AND NOT p.Waiting;
		2: p(Start := FALSE, Proceed := TRUE, Jump := TRUE, JumpStep := 3); ok := p.Step = 1;
		3: p(Proceed := FALSE, Jump := FALSE, Ready := TRUE); ok := p.Step = 2 AND NOT p.Waiting;
		4: p(); ok := p.Step = 3;
		5: p(); ok := p.State = ISA88_RUNNING; (* step 3 met: Done *)
		6: p(); ok := p.State = ISA88_COMPLETE;
		   done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM

(* Semi-automatic: a step needs its condition and a Proceed; one Proceed is
   one step; a Proceed before the condition is met is not kept *)
PROGRAM TEST_MODE_SemiAutomatic
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		p    : PH_DEMO_STEPS;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: p(Mode := ISA88_SEMI_AUTOMATIC, Start := TRUE); ok := p.State = ISA88_RUNNING AND p.ActiveMode = ISA88_SEMI_AUTOMATIC AND NOT p.Waiting;
		1: p(Start := FALSE, Ready := TRUE); ok := p.Step = 1 AND p.Waiting;
		2: p(Proceed := TRUE); ok := p.Step = 2 AND NOT p.Waiting;
		3: p(); ok := p.Step = 2 AND p.Waiting; (* Proceed still held: no edge *)
		4: p(Proceed := FALSE, Ready := FALSE); ok := p.Step = 2 AND NOT p.Waiting;
		5: p(Proceed := TRUE); ok := p.Step = 2; (* condition not met *)
		6: p(Proceed := FALSE, Ready := TRUE); ok := p.Step = 2 AND p.Waiting; (* that Proceed was not kept *)
		7: p(Proceed := TRUE); ok := p.Step = 3;
		8: p(Proceed := FALSE); ok := p.Step = 3 AND p.Waiting;
		9: p(Proceed := TRUE); ok := p.State = ISA88_RUNNING; (* Done *)
		10: p(Proceed := FALSE); ok := p.State = ISA88_COMPLETE;
		    done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM

(* Manual: the operator steps whatever the condition and may jump; the mode
   may change while running; an unknown mode is ignored *)
PROGRAM TEST_MODE_Manual
	VAR_OUTPUT
		failures : INT;
		message  : STRING;
		done     : BOOL;
	END_VAR
	VAR
		p    : PH_DEMO_STEPS;
		step : INT;
		ok   : BOOL;
	END_VAR
	ok := TRUE;
	CASE step OF
		0: p(Mode := ISA88_MANUAL, Start := TRUE); ok := p.State = ISA88_RUNNING AND p.Step = 1 AND p.Waiting;
		1: p(Start := FALSE, Proceed := TRUE); ok := p.Step = 2; (* Ready is FALSE *)
		2: p(Proceed := FALSE, Jump := TRUE, JumpStep := 1); ok := p.Step = 1;
		3: p(Jump := FALSE, Mode := 7); ok := p.ActiveMode = ISA88_MANUAL AND p.Step = 1;
		4: p(Proceed := TRUE); ok := p.Step = 2;
		5: p(Proceed := FALSE, Mode := ISA88_AUTOMATIC, Ready := TRUE); ok := p.ActiveMode = ISA88_AUTOMATIC AND p.Step = 3 AND NOT p.Waiting;
		6: p(Jump := TRUE, JumpStep := 1); ok := p.State = ISA88_RUNNING; (* no jump in Automatic: Done *)
		7: p(Jump := FALSE); ok := p.State = ISA88_COMPLETE AND p.Step = 3;
		   done := TRUE;
	END_CASE
	IF NOT ok THEN
		failures := failures + 1;
		message := CONCAT('wrong state after step ', INT_TO_STRING(step));
	END_IF;
	step := step + 1;
END_PROGRAM
