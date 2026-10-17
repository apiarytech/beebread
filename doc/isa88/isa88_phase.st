(*
 * ISA-88 phase template: the lowest procedural element, the one that runs
 * equipment. ISA88_PHASE wraps ISA88_STATE_MACHINE and adds what a phase's
 * equipment logic needs from it. A phase is written by composition:
 *
 *   FUNCTION_BLOCK PH_HEAT
 *     VAR ph : ISA88_PHASE; END_VAR
 *     ph(Start := cmdStart, Hold := cmdHold, FailCode := fault);
 *     CASE ph.State OF
 *       ISA88_RUNNING:
 *         IF ph.Entry THEN heater := TRUE; END_IF;
 *         IF temp >= setpoint THEN ph.Done(); END_IF;
 *       ISA88_HOLDING:
 *         heater := FALSE;
 *         ph.Done();
 *     END_CASE
 *
 * Each scan the phase logic calls ph first, then runs the logic of ph.State.
 *   Done()    : called by the logic when the current Running or transient
 *               state's work is finished (SC). The next call of ph takes it
 *               and forgets it, so it never carries into the next state.
 *               A transient state whose logic never calls Done stays in it.
 *   Entry     : TRUE on the first scan the logic runs a state; use it to
 *               start that state's work (open a valve, latch a restart step).
 *   StateTime : time spent in the current state, for timeouts.
 *   FailCode  : non-zero while the equipment has a failure. It holds the
 *               phase (Running, Pausing, Paused and Restarting accept Hold),
 *               so a Restart while the failure is still there holds again.
 *               The first code is latched in Failure until the phase is
 *               back in Idle.
 *
 * Commands are pulses, as for ISA88_STATE_MACHINE.
 *
 * Modes (ISA-88 procedural modes) decide how the logic moves from one step
 * to the next. The logic asks with Advance(condition):
 *
 *   ISA88_RUNNING:
 *     IF ph.JumpRequest THEN step := ph.JumpTarget;
 *     ELSIF ph.Advance(temp >= setpoint) THEN step := step + 1;
 *     END_IF;
 *
 *   Mode        : ISA88_AUTOMATIC (the default), ISA88_SEMI_AUTOMATIC or
 *                 ISA88_MANUAL. It may change at any time; another value
 *                 is ignored and ActiveMode keeps the last valid one.
 *   Advance(c)  : TRUE when the logic may take its next step.
 *                   Automatic      : when c is TRUE.
 *                   Semi-automatic : when c is TRUE and the operator
 *                                    presses Proceed.
 *                   Manual         : when the operator presses Proceed,
 *                                    whatever c is.
 *                 One Proceed allows one step.
 *   Waiting     : the logic is waiting on Proceed (Semi-automatic with its
 *                 condition met, or Manual).
 *   Jump        : Manual only: JumpRequest is TRUE for one scan, and the
 *                 logic moves to JumpTarget (JumpStep when Jump was pressed).
 *
 * Proceed and Jump act on their rising edge. Done() is not tied to a mode:
 * holding, stopping and aborting finish on their own in every mode, and a
 * step that should wait for the operator writes IF ph.Advance(c) THEN
 * ph.Done(); END_IF.
 *)
VAR_GLOBAL CONSTANT
	ISA88_IDLE       : INT := 1;
	ISA88_RUNNING    : INT := 2;
	ISA88_COMPLETE   : INT := 3;
	ISA88_PAUSING    : INT := 4;
	ISA88_PAUSED     : INT := 5;
	ISA88_HOLDING    : INT := 6;
	ISA88_HELD       : INT := 7;
	ISA88_RESTARTING : INT := 8;
	ISA88_STOPPING   : INT := 9;
	ISA88_STOPPED    : INT := 10;
	ISA88_ABORTING   : INT := 11;
	ISA88_ABORTED    : INT := 12;

	ISA88_AUTOMATIC      : INT := 1;
	ISA88_SEMI_AUTOMATIC : INT := 2;
	ISA88_MANUAL         : INT := 3;
END_VAR

FUNCTION_BLOCK ISA88_PHASE
	VAR_INPUT
		Start    : BOOL;
		Pause    : BOOL;
		Resume   : BOOL;
		Hold     : BOOL;
		Restart  : BOOL;
		Stop     : BOOL;
		Abort    : BOOL;
		Reset    : BOOL;
		FailCode : DINT; (* non-zero: equipment failure *)
		Mode     : INT := 1; (* ISA88_AUTOMATIC *)
		Proceed  : BOOL;
		Jump     : BOOL;
		JumpStep : INT;
	END_VAR
	VAR_OUTPUT
		State       : INT;
		Entry       : BOOL;
		StateTime   : TIME;
		Failure     : DINT;
		IsQuiescent : BOOL;
		IsTransient : BOOL;
		IsFinal     : BOOL;
		ActiveMode  : INT := 1; (* ISA88_AUTOMATIC *)
		Waiting     : BOOL;
		JumpRequest : BOOL;
		JumpTarget  : INT;
	END_VAR
	VAR
		sm          : ISA88_STATE_MACHINE;
		complete    : BOOL; (* Done was called since the last scan *)
		lastState   : INT;
		inState     : TON;
		proceedEdge : R_TRIG;
		jumpEdge    : R_TRIG;
		proceeding  : BOOL; (* a Proceed this scan not yet used by Advance *)
	END_VAR

	(* The current state's work is finished (SC). *)
	METHOD Done : BOOL
		complete := TRUE;
		Done := TRUE;
	END_METHOD

	(* May the logic take its next step, its condition being Cond? *)
	METHOD Advance : BOOL
		VAR_INPUT
			Cond : BOOL;
		END_VAR
		IF ActiveMode = ISA88_AUTOMATIC THEN
			Advance := Cond;
		ELSIF ActiveMode = ISA88_SEMI_AUTOMATIC AND NOT Cond THEN
			Advance := FALSE;
		ELSE
			Advance := proceeding;
			proceeding := FALSE;
			Waiting := NOT Advance;
		END_IF;
	END_METHOD

	sm(Start := Start, Pause := Pause, Resume := Resume,
	   Hold := Hold OR FailCode <> 0, Restart := Restart,
	   Stop := Stop, Abort := Abort, Reset := Reset, SC := complete);
	complete := FALSE;

	State := sm.State;
	IsQuiescent := sm.IsQuiescent;
	IsTransient := sm.IsTransient;
	IsFinal := sm.IsFinal;

	Entry := State <> lastState;
	lastState := State;
	inState(IN := NOT Entry, PT := T#24D);
	StateTime := inState.ET;

	IF State = ISA88_IDLE THEN
		Failure := 0;
	ELSIF FailCode <> 0 AND Failure = 0 THEN
		Failure := FailCode;
	END_IF;

	IF Mode = ISA88_AUTOMATIC OR Mode = ISA88_SEMI_AUTOMATIC OR Mode = ISA88_MANUAL THEN
		ActiveMode := Mode;
	END_IF;
	proceedEdge(CLK := Proceed);
	proceeding := proceedEdge.Q AND ActiveMode <> ISA88_AUTOMATIC;
	jumpEdge(CLK := Jump);
	JumpRequest := jumpEdge.Q AND ActiveMode = ISA88_MANUAL;
	IF JumpRequest THEN
		JumpTarget := JumpStep;
	END_IF;
	Waiting := FALSE; (* until Advance says otherwise this scan *)
END_FUNCTION_BLOCK
