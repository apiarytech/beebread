# ISA-88 procedural state machine

[isa88_state_machine.st](isa88_state_machine.st) is the ISA-88
(IEC 61512-1) procedural state model as an IEC 61131-3 Sequential Function
Chart in text form: function block `ISA88_STATE_MACHINE`, one step per
state, one transition per arrow of the state diagram.

| Class     | States |
|-----------|--------|
| Initial   | Idle |
| Quiescent | Idle, Paused, Held |
| Transient | Running, Pausing, Holding, Restarting, Stopping, Aborting |
| Final     | Complete, Stopped, Aborted |

| Command | Accepted in | Goes to |
|---------|-------------|---------|
| Start   | Idle | Running |
| Pause   | Running | Pausing |
| Resume  | Paused | Running |
| Hold    | Running, Pausing, Paused, Restarting | Holding |
| Restart | Held | Restarting |
| Stop    | Running, Pausing, Paused, Holding, Held, Restarting | Stopping |
| Abort   | all but Idle, Complete, Aborting, Aborted | Aborting |
| Reset   | Complete, Stopped, Aborted | Idle |
| SC      | Running, Pausing, Holding, Restarting, Stopping, Aborting | Complete, Paused, Held, Running, Stopped, Aborted |

SC (state complete) is set by the procedural logic when the current state's
logic is done. Commands are pulses. When a state receives more than one
command in the same scan, Abort wins, then Stop, then Hold, then the
others, then SC.

Outputs: `State` (1 Idle, 2 Running, 3 Complete, 4 Pausing, 5 Paused,
6 Holding, 7 Held, 8 Restarting, 9 Stopping, 10 Stopped, 11 Aborting,
12 Aborted; ISA-88 defines no numbers), and `IsQuiescent`, `IsTransient`,
`IsFinal`.

## Phase template

[isa88_phase.st](isa88_phase.st) declares `ISA88_PHASE`, the template for
a phase, the procedural element that runs equipment. A phase contains an
`ISA88_PHASE` and runs its own logic for each state:

```iecst
FUNCTION_BLOCK PH_HEAT
  VAR ph : ISA88_PHASE; END_VAR
  ph(Start := cmdStart, Hold := cmdHold, FailCode := fault);
  CASE ph.State OF
    ISA88_RUNNING:
      IF ph.Entry THEN heater := TRUE; END_IF;
      IF temp >= setpoint THEN ph.Done(); END_IF;
    ISA88_HOLDING:
      heater := FALSE;
      ph.Done();
  END_CASE
END_FUNCTION_BLOCK
```

| Member | |
|--------|---|
| `Done()` | The logic calls this when the current Running or transient state's work is finished (SC). The next call of `ph` acts on it and then forgets it, so it never carries into the next state. |
| `Entry` | TRUE on the first scan the logic runs a state. |
| `StateTime` | Time spent in the current state, for timeouts. |
| `FailCode` | Non-zero while the equipment has a failure. It holds the phase, and a Restart while the failure lasts holds again. |
| `Failure` | The first failure code, kept until the phase is back in Idle. |

The file also declares the state codes as constants, `ISA88_IDLE` to
`ISA88_ABORTED`, for `CASE` labels. Variables a phase keeps between scans,
such as a step number, are its restart point: Running picks up where it was
held.

## Modes

The ISA-88 procedural modes decide how the phase logic moves from one step
to the next. The logic asks the phase with `Advance(condition)`:

```iecst
ISA88_RUNNING:
  IF ph.JumpRequest THEN
    step := ph.JumpTarget;
  ELSIF ph.Advance(temp >= setpoint) THEN
    IF step >= 3 THEN ph.Done(); ELSE step := step + 1; END_IF;
  END_IF;
```

| `Mode` | `Advance(c)` is TRUE | Jump |
|--------|----------------------|------|
| `ISA88_AUTOMATIC` (default) | when `c` is TRUE | ignored |
| `ISA88_SEMI_AUTOMATIC` | when `c` is TRUE and the operator presses `Proceed` | ignored |
| `ISA88_MANUAL` | when the operator presses `Proceed`, whatever `c` is | `JumpRequest` for one scan, to `JumpTarget` |

- `Proceed` and `Jump` act on their rising edge. One `Proceed` allows one
  step, and a `Proceed` pressed before the condition is met is not kept.
- `Waiting` is TRUE while the logic waits on `Proceed`, for the HMI.
- `Mode` may change at any time. Any other value is ignored, and
  `ActiveMode` keeps the last valid mode.
- `Done()` is not tied to a mode, so Holding, Stopping and Aborting finish
  on their own in every mode. A step that should wait for the operator
  writes `IF ph.Advance(c) THEN ph.Done(); END_IF`.

Who may send commands, the operator or a program (the owner), is part of
the phase logic interface, not of the modes.

## Sequencer

[isa88_sequencer.st](isa88_sequencer.st) declares `ISA88_SEQUENCER`: up to
32 steps (0 to 31) run one after the other, forward or in reverse. Each
step is an ISA-88 phase whose Running has three timed parts:

```iecst
seq(Forward := pbFwd, Reverse := pbRev, Hold := pbHold, Stop := pbStop,
    Count := 3, StartTime := tStart, RunTime := tRun, StopTime := tStop);
conveyor1 := seq.Running[0];
```

| Part | Output | Lasts |
|------|--------|-------|
| Starting | `Starting[i]` | `StartTime[i]` from the step's start |
| Running  | `Running[i]`  | then `RunTime[i]` |
| Stopping | `Stopping[i]` | then `StopTime[i]`; the step is Complete and the next one starts |

| Input | |
|-------|---|
| `Forward`, `Reverse` | In Idle, start at step 0 going up, or at step `Count-1` going down. While the sequence runs, change the direction: the running step finishes and the next one is picked in the new direction. |
| `Pause`, `Resume`, `Hold`, `Restart`, `Stop`, `Abort` | Go to the running step, as ISA-88 says. A command given between two steps waits for the next one; Stop or Abort there ends the sequence at once. |
| `Reset` | From Complete, Stopped or Aborted back to Idle. |
| `Count` | Steps in use, 0 to 32. It may change while the sequence runs. 0 steps complete at once. |

`State` is the sequence's ISA-88 state: Idle, the running step's state,
then Complete, Stopped or Aborted. `Index` is the step in use (-1 for
none), `Reversing` the direction, `StepState[i]` each step's state. Held or
Paused, a step's time stands still and Restart or Resume carry on from where
it was. A Stop runs what is left of the step's `StopTime` (all of it if the
step had not reached Stopping). Pausing, Holding, Restarting and Aborting
finish at once.

The order of the steps is a Sequential Function Chart, `SEQ_CHART`:

| Chart step | Does | Goes to |
|------------|------|---------|
| `S_Idle` | | `S_Start` on Forward or Reverse (`S_Complete` with no steps) |
| `S_Start` | picks the step, asks it to start | `S_Run` once it has left Idle |
| `S_Run` | | `S_Next` when the step is Complete, `S_Stopped` or `S_Aborted` when it is |
| `S_Next` | resets the step, looks for the next one | `S_Start`, or `S_Complete` past the last one; `S_Stopped` or `S_Aborted` on a waiting Stop or Abort |
| `S_Complete`, `S_Stopped`, `S_Aborted` | | `S_Idle` on Reset |

Each scan `ISA88_SEQUENCER` calls every step once, as the chart asked, then
the chart on the steps' states. A chart action does not call the steps: in a
scan in which the chart moves, the actions of the steps it leaves and enters
both run, so a step would be called twice.

## Tests

```bash
beedance -test -engines all isa88_state_machine.st isa88_state_machine_test.st
beedance -test -engines all isa88_state_machine.st isa88_phase.st isa88_phase_test.st
beedance -test -engines all isa88_state_machine.st isa88_phase.st isa88_phase_mode_test.st
beedance -test -engines all isa88_state_machine.st isa88_phase.st isa88_sequencer.st isa88_sequencer_test.st
```

[isa88_state_machine_test.st](isa88_state_machine_test.st) walks the
complete, pause/hold/restart and stop/abort paths, the command ranking, and
commands a state ignores.
[isa88_phase_test.st](isa88_phase_test.st) runs a demo heating phase
through completion, hold and restart (Done not carrying over, Restarting
timed on `StateTime`, Running resuming where it was), and failures.
[isa88_phase_mode_test.st](isa88_phase_mode_test.st) steps a three-step
demo phase in each mode, with Proceed edges, jumps and mode changes while
running.
[isa88_sequencer_test.st](isa88_sequencer_test.st) runs sequences forward,
in reverse, turning back halfway, with no steps and with all 32, times a
step's three parts, and holds, restarts, stops and aborts them, also between
two steps. All run on the evaluator, the VM and transpiled Go.
