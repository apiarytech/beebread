(*
 * ISA-88 (IEC 61512-1) procedural state machine as a Sequential Function Chart.
 *
 * One step per procedural state; one transition per arrow of the ISA-88
 * state diagram. Commands are inputs, pulsed by the caller for one scan.
 * SC (state complete) is raised by the procedural logic when the logic of
 * the current Running or transient state has finished.
 *
 *   Initial    : Idle
 *   Quiescent  : Idle, Paused, Held
 *   Transient  : Running, Pausing, Holding, Restarting, Stopping, Aborting
 *   Final      : Complete, Stopped, Aborted
 *
 * Which commands a state accepts (the nested boxes of the diagram):
 *   Pause  : Running
 *   Resume : Paused
 *   Hold   : Running, Pausing, Paused, Restarting
 *   Restart: Held
 *   Stop   : Running, Pausing, Paused, Holding, Held, Restarting
 *   Abort  : every state except Idle, Complete, Aborting and Aborted
 *   Reset  : Complete, Stopped, Aborted
 *
 * A step fires at most one transition per scan, so where a state accepts
 * more than one command they are ranked: Abort, Stop, Hold, Pause or
 * Resume or Restart, then SC.
 *
 * State codes (ISA-88 names no numbers; these are beebread's):
 *   1 Idle      2 Running   3 Complete
 *   4 Pausing   5 Paused    6 Holding   7 Held      8 Restarting
 *   9 Stopping 10 Stopped  11 Aborting 12 Aborted
 *)
FUNCTION_BLOCK ISA88_STATE_MACHINE
	VAR_INPUT
		Start   : BOOL;
		Pause   : BOOL;
		Resume  : BOOL;
		Hold    : BOOL;
		Restart : BOOL;
		Stop    : BOOL;
		Abort   : BOOL;
		Reset   : BOOL;
		SC      : BOOL; (* state complete *)
	END_VAR
	VAR_OUTPUT
		State       : INT;
		IsQuiescent : BOOL;
		IsTransient : BOOL;
		IsFinal     : BOOL;
	END_VAR

	ACTION A_Idle:
		State := 1; IsQuiescent := TRUE; IsTransient := FALSE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Running:
		State := 2; IsQuiescent := FALSE; IsTransient := TRUE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Complete:
		State := 3; IsQuiescent := FALSE; IsTransient := FALSE; IsFinal := TRUE;
	END_ACTION
	ACTION A_Pausing:
		State := 4; IsQuiescent := FALSE; IsTransient := TRUE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Paused:
		State := 5; IsQuiescent := TRUE; IsTransient := FALSE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Holding:
		State := 6; IsQuiescent := FALSE; IsTransient := TRUE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Held:
		State := 7; IsQuiescent := TRUE; IsTransient := FALSE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Restarting:
		State := 8; IsQuiescent := FALSE; IsTransient := TRUE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Stopping:
		State := 9; IsQuiescent := FALSE; IsTransient := TRUE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Stopped:
		State := 10; IsQuiescent := FALSE; IsTransient := FALSE; IsFinal := TRUE;
	END_ACTION
	ACTION A_Aborting:
		State := 11; IsQuiescent := FALSE; IsTransient := TRUE; IsFinal := FALSE;
	END_ACTION
	ACTION A_Aborted:
		State := 12; IsQuiescent := FALSE; IsTransient := FALSE; IsFinal := TRUE;
	END_ACTION

	(* ---- Initial state ---- *)
	INITIAL_STEP Idle: A_Idle(N); END_STEP
	TRANSITION FROM Idle TO Running := Start; END_TRANSITION

	(* ---- Running ---- *)
	STEP Running: A_Running(N); END_STEP
	TRANSITION FROM Running TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Running TO Stopping := Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Running TO Holding := Hold AND NOT Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Running TO Pausing := Pause AND NOT Hold AND NOT Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Running TO Complete := SC AND NOT Pause AND NOT Hold AND NOT Stop AND NOT Abort; END_TRANSITION

	(* ---- Pause ---- *)
	STEP Pausing: A_Pausing(N); END_STEP
	TRANSITION FROM Pausing TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Pausing TO Stopping := Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Pausing TO Holding := Hold AND NOT Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Pausing TO Paused := SC AND NOT Hold AND NOT Stop AND NOT Abort; END_TRANSITION

	STEP Paused: A_Paused(N); END_STEP
	TRANSITION FROM Paused TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Paused TO Stopping := Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Paused TO Holding := Hold AND NOT Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Paused TO Running := Resume AND NOT Hold AND NOT Stop AND NOT Abort; END_TRANSITION

	(* ---- Hold ---- *)
	STEP Holding: A_Holding(N); END_STEP
	TRANSITION FROM Holding TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Holding TO Stopping := Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Holding TO Held := SC AND NOT Stop AND NOT Abort; END_TRANSITION

	STEP Held: A_Held(N); END_STEP
	TRANSITION FROM Held TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Held TO Stopping := Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Held TO Restarting := Restart AND NOT Stop AND NOT Abort; END_TRANSITION

	STEP Restarting: A_Restarting(N); END_STEP
	TRANSITION FROM Restarting TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Restarting TO Stopping := Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Restarting TO Holding := Hold AND NOT Stop AND NOT Abort; END_TRANSITION
	TRANSITION FROM Restarting TO Running := SC AND NOT Hold AND NOT Stop AND NOT Abort; END_TRANSITION

	(* ---- Stop ---- *)
	STEP Stopping: A_Stopping(N); END_STEP
	TRANSITION FROM Stopping TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Stopping TO Stopped := SC AND NOT Abort; END_TRANSITION

	STEP Stopped: A_Stopped(N); END_STEP
	TRANSITION FROM Stopped TO Aborting := Abort; END_TRANSITION
	TRANSITION FROM Stopped TO Idle := Reset AND NOT Abort; END_TRANSITION

	(* ---- Abort ---- *)
	STEP Aborting: A_Aborting(N); END_STEP
	TRANSITION FROM Aborting TO Aborted := SC; END_TRANSITION

	STEP Aborted: A_Aborted(N); END_STEP
	TRANSITION FROM Aborted TO Idle := Reset; END_TRANSITION

	(* ---- Complete ---- *)
	STEP Complete: A_Complete(N); END_STEP
	TRANSITION FROM Complete TO Idle := Reset; END_TRANSITION
END_FUNCTION_BLOCK
