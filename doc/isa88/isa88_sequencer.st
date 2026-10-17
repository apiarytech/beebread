(*
 * ISA-88 sequencer: up to 32 steps (0 to 31) run one after the other,
 * forward (0, 1, 2 ...) or in reverse (Count-1 ... 1, 0). Each step is an
 * ISA-88 phase (ISA88_PHASE, so an ISA88_STATE_MACHINE) with three times:
 *
 *   Starting : StartTime after the step starts
 *   Running  : then RunTime
 *   Stopping : then StopTime; the step is then Complete and the next one
 *              starts.
 *
 *   seq(Forward := pbFwd, Hold := pbHold, Count := 3,
 *       StartTime := tStart, RunTime := tRun, StopTime := tStop);
 *   motor1 := seq.Running[0];
 *
 * Commands are pulses, as for ISA88_STATE_MACHINE.
 *   Forward, Reverse : in Idle, start the sequence in that direction.
 *                      While it runs, change the direction: the running
 *                      step finishes and the next one is picked in the new
 *                      direction (Reverse from step 0, or Forward from the
 *                      last step, then completes the sequence).
 *   Pause, Resume, Hold, Restart, Stop, Abort
 *                    : go to the running step and act as ISA-88 says. A
 *                      command given between two steps waits for the next
 *                      one; Stop or Abort there ends the sequence at once.
 *   Reset            : from Complete, Stopped or Aborted back to Idle.
 *   Count            : steps in use, 0 to 32; it may change while the
 *                      sequence runs. 0 steps complete at once.
 *
 * State is the ISA-88 state of the sequence: Idle, then the running step's
 * state, then Complete, Stopped or Aborted. A Stop runs what is left of the
 * step's StopTime (all of it if the step had not reached Stopping); Pausing,
 * Holding, Restarting and Aborting finish at once. Held or Paused, a step's
 * time stands still, and Restart or Resume carry on from where it was.
 *
 * The order of the steps is a Sequential Function Chart, SEQ_CHART. Each
 * scan ISA88_SEQUENCER calls every step once, as the chart asked, then the
 * chart, on the steps' states. (A chart action that called the steps would
 * run them twice in a scan in which the chart moves.) A step starts the
 * scan after the chart picks it.
 *)

(* One step: an ISA-88 phase whose Running is Starting, Running and
   Stopping, timed. *)
FUNCTION_BLOCK SEQ_STEP
	VAR_INPUT
		Start     : BOOL;
		Pause     : BOOL;
		Resume    : BOOL;
		Hold      : BOOL;
		Restart   : BOOL;
		Stop      : BOOL;
		Abort     : BOOL;
		Reset     : BOOL;
		StartTime : TIME;
		RunTime   : TIME;
		StopTime  : TIME;
	END_VAR
	VAR_OUTPUT
		State    : INT;
		Starting : BOOL;
		Running  : BOOL;
		Stopping : BOOL;
		Elapsed  : TIME; (* time run in Starting, Running and Stopping *)
	END_VAR
	VAR
		ph       : ISA88_PHASE;
		base     : TIME; (* Elapsed when Running was entered *)
		stopLeft : TIME; (* what a Stop has left to run *)
	END_VAR
	ph(Start := Start, Pause := Pause, Resume := Resume, Hold := Hold,
	   Restart := Restart, Stop := Stop, Abort := Abort, Reset := Reset);
	Starting := FALSE;
	Running := FALSE;
	Stopping := FALSE;
	CASE ph.State OF
		ISA88_IDLE:
			Elapsed := T#0s;
		ISA88_RUNNING:
			IF ph.Entry THEN base := Elapsed; END_IF;
			Elapsed := base + ph.StateTime;
			IF Elapsed < StartTime THEN
				Starting := TRUE;
			ELSIF Elapsed < StartTime + RunTime THEN
				Running := TRUE;
			ELSIF Elapsed < StartTime + RunTime + StopTime THEN
				Stopping := TRUE;
			ELSE
				ph.Done();
			END_IF;
		ISA88_STOPPING:
			IF ph.Entry THEN
				IF Elapsed >= StartTime + RunTime + StopTime THEN
					stopLeft := T#0s;
				ELSIF Elapsed > StartTime + RunTime THEN
					stopLeft := StartTime + RunTime + StopTime - Elapsed;
				ELSE
					stopLeft := StopTime;
				END_IF;
			END_IF;
			IF ph.StateTime >= stopLeft THEN
				ph.Done();
			ELSE
				Stopping := TRUE;
			END_IF;
		ISA88_PAUSING, ISA88_HOLDING, ISA88_RESTARTING, ISA88_ABORTING:
			ph.Done();
	END_CASE
	State := ph.State;
END_FUNCTION_BLOCK

(* The order of the steps. StepState is the state of step Index; the
   chart asks for StepStart and StepReset on step Index. StopReq and
   AbortReq are a Stop or Abort no step has taken yet. *)
FUNCTION_BLOCK SEQ_CHART
	VAR_INPUT
		Go        : BOOL; (* Forward or Reverse *)
		Reversing : BOOL;
		Count     : INT;
		StepState : INT;
		StopReq   : BOOL;
		AbortReq  : BOOL;
		Reset     : BOOL;
	END_VAR
	VAR_OUTPUT
		Index     : INT := -1; (* the step in use; -1: none *)
		StepStart : BOOL;
		StepReset : BOOL;
		Busy      : BOOL; (* from the first step's start to the last step's end *)
		State     : INT := 1; (* ISA88_IDLE; while Busy, see the step *)
	END_VAR
	VAR
		NextIndex : INT;
		HasNext   : BOOL;
	END_VAR

	ACTION A_Idle:
		State := ISA88_IDLE;
		Index := -1;
		IF Reversing THEN NextIndex := Count - 1; ELSE NextIndex := 0; END_IF;
	END_ACTION
	ACTION A_Enter:
		Index := NextIndex;
		State := ISA88_RUNNING;
	END_ACTION
	ACTION A_Look:
		IF Reversing THEN
			NextIndex := MIN(Index, Count) - 1;
		ELSE
			NextIndex := Index + 1;
		END_IF;
		HasNext := NextIndex >= 0 AND NextIndex < Count;
	END_ACTION
	ACTION A_Complete:
		State := ISA88_COMPLETE;
		Index := -1;
	END_ACTION
	ACTION A_Stopped:
		State := ISA88_STOPPED;
	END_ACTION
	ACTION A_Aborted:
		State := ISA88_ABORTED;
	END_ACTION

	INITIAL_STEP S_Idle: A_Idle(N); END_STEP
	TRANSITION FROM S_Idle TO S_Start := Go AND Count > 0; END_TRANSITION
	TRANSITION FROM S_Idle TO S_Complete := Go AND Count <= 0; END_TRANSITION

	(* Start step Index; it has started once it leaves Idle. *)
	STEP S_Start: A_Enter(P); StepStart(N); Busy(N); END_STEP
	TRANSITION FROM S_Start TO S_Run := StepState <> ISA88_IDLE; END_TRANSITION

	(* The step runs, through whatever commands it is given, to an end. *)
	STEP S_Run: Busy(N); END_STEP
	TRANSITION FROM S_Run TO S_Next := StepState = ISA88_COMPLETE; END_TRANSITION
	TRANSITION FROM S_Run TO S_Stopped := StepState = ISA88_STOPPED; END_TRANSITION
	TRANSITION FROM S_Run TO S_Aborted := StepState = ISA88_ABORTED; END_TRANSITION

	(* Reset the finished step, then start the next one, if there is one. *)
	STEP S_Next: StepReset(N); A_Look(N); Busy(N); END_STEP
	TRANSITION FROM S_Next TO S_Aborted := AbortReq; END_TRANSITION
	TRANSITION FROM S_Next TO S_Stopped := StopReq AND NOT AbortReq; END_TRANSITION
	TRANSITION FROM S_Next TO S_Start := StepState = ISA88_IDLE AND HasNext AND NOT StopReq AND NOT AbortReq; END_TRANSITION
	TRANSITION FROM S_Next TO S_Complete := StepState = ISA88_IDLE AND NOT HasNext AND NOT StopReq AND NOT AbortReq; END_TRANSITION

	STEP S_Complete: A_Complete(N); END_STEP
	TRANSITION FROM S_Complete TO S_Idle := Reset; END_TRANSITION

	STEP S_Stopped: A_Stopped(N); END_STEP
	TRANSITION FROM S_Stopped TO S_Idle := Reset; END_TRANSITION

	STEP S_Aborted: A_Aborted(N); END_STEP
	TRANSITION FROM S_Aborted TO S_Idle := Reset; END_TRANSITION
END_FUNCTION_BLOCK

FUNCTION_BLOCK ISA88_SEQUENCER
	VAR_INPUT
		Forward   : BOOL;
		Reverse   : BOOL;
		Pause     : BOOL;
		Resume    : BOOL;
		Hold      : BOOL;
		Restart   : BOOL;
		Stop      : BOOL;
		Abort     : BOOL;
		Reset     : BOOL;
		Count     : INT; (* steps in use, 0 to 32 *)
		StartTime : ARRAY[0..31] OF TIME;
		RunTime   : ARRAY[0..31] OF TIME;
		StopTime  : ARRAY[0..31] OF TIME;
	END_VAR
	VAR_OUTPUT
		State     : INT := 1; (* ISA88_IDLE *)
		Index     : INT := -1; (* the step in use; -1: none *)
		Reversing : BOOL;
		Busy      : BOOL;
		Starting  : ARRAY[0..31] OF BOOL;
		Running   : ARRAY[0..31] OF BOOL;
		Stopping  : ARRAY[0..31] OF BOOL;
		StepState : ARRAY[0..31] OF INT;
	END_VAR
	VAR
		chart    : SEQ_CHART;
		steps    : ARRAY[0..31] OF SEQ_STEP;
		n        : INT;
		i        : INT;
		resetAll : BOOL;
		was      : INT; (* the state of step Index before this scan's call *)
		taken    : BOOL; (* step Index can take the pending commands *)
		pauseReq : BOOL;
		holdReq  : BOOL;
		stopReq  : BOOL;
		abortReq : BOOL;
	END_VAR
	n := LIMIT(0, Count, 32);
	IF Forward THEN
		Reversing := FALSE;
	ELSIF Reverse THEN
		Reversing := TRUE;
	END_IF;

	(* Commands wait while no step can take them: between two steps the
	   finished one is Complete and the next one Idle. *)
	IF chart.Busy THEN
		pauseReq := (pauseReq OR Pause) AND NOT Resume;
		holdReq := (holdReq OR Hold) AND NOT Restart;
		stopReq := stopReq OR Stop;
		abortReq := abortReq OR Abort;
	ELSE
		pauseReq := FALSE;
		holdReq := FALSE;
		stopReq := FALSE;
		abortReq := FALSE;
	END_IF;
	resetAll := Reset AND NOT chart.Busy;

	(* The steps first, as the chart asked last scan, then the chart, on
	   the steps as they are now. *)
	Index := chart.Index;
	Busy := chart.Busy;
	IF Index >= 0 AND Index <= 31 THEN
		was := steps[Index].State;
	ELSE
		was := ISA88_IDLE;
	END_IF;
	taken := was <> ISA88_IDLE AND was <> ISA88_COMPLETE;
	FOR i := 0 TO 31 DO
		IF i = Index AND Busy THEN
			steps[i](Start := chart.StepStart, Pause := pauseReq, Resume := Resume,
			         Hold := holdReq, Restart := Restart, Stop := stopReq,
			         Abort := abortReq, Reset := chart.StepReset OR resetAll,
			         StartTime := StartTime[i], RunTime := RunTime[i], StopTime := StopTime[i]);
		ELSE
			steps[i](Start := FALSE, Pause := FALSE, Resume := FALSE, Hold := FALSE,
			         Restart := FALSE, Stop := FALSE, Abort := FALSE, Reset := resetAll,
			         StartTime := StartTime[i], RunTime := RunTime[i], StopTime := StopTime[i]);
		END_IF;
		Starting[i] := steps[i].Starting;
		Running[i] := steps[i].Running;
		Stopping[i] := steps[i].Stopping;
		StepState[i] := steps[i].State;
	END_FOR;
	IF taken THEN
		pauseReq := FALSE;
		holdReq := FALSE;
	END_IF;

	IF Index >= 0 AND Index <= 31 THEN
		was := steps[Index].State;
	ELSE
		was := ISA88_IDLE;
	END_IF;
	chart(Go := Forward OR Reverse, Reversing := Reversing, Count := n,
	      StepState := was, StopReq := stopReq, AbortReq := abortReq,
	      Reset := Reset);
	Index := chart.Index;
	Busy := chart.Busy;

	IF Busy AND Index >= 0 AND Index <= 31 THEN
		State := steps[Index].State;
		IF State = ISA88_IDLE OR State = ISA88_COMPLETE THEN
			State := ISA88_RUNNING; (* between two steps *)
		END_IF;
	ELSE
		State := chart.State;
	END_IF;
END_FUNCTION_BLOCK
