/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - EPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package electrical

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/engineering"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// DIMM_2 is a dimmer with two push buttons: a click of I1 turns it on and a
// click of I2 off, holding I1 or I2 longer than T_DIMM_START dims OUT up or
// down, over T_DIMM from 0 to 255. SET sets OUT to VAL, RST turns the dimmer
// off (and OUT to 0 if RST_OUT). A double click toggles D1 or D2, or sets OUT
// to DBL1_POS or DBL2_POS if DBL1_SET or DBL2_SET. Q goes off after T_ON_MAX
// if that is not 0.
type DIMM_2 struct {
	SET          iec.BOOL
	VAL          iec.BYTE // default 255
	I1           iec.BOOL
	I2           iec.BOOL
	RST          iec.BOOL
	T_DEBOUNCE   iec.TIME // default T#10ms
	T_ON_MAX     iec.TIME
	T_DIMM_START iec.TIME // default T#1s
	T_DIMM       iec.TIME // default T#3s
	MIN_ON       iec.BYTE // default 50
	MAX_ON       iec.BYTE // default 255
	RST_OUT      iec.BOOL // default TRUE
	SOFT_DIMM    iec.BOOL // default TRUE
	DBL1_TOG     iec.BOOL
	DBL2_TOG     iec.BOOL
	DBL1_SET     iec.BOOL
	DBL2_SET     iec.BOOL
	DBL1_POS     iec.BYTE
	DBL2_POS     iec.BYTE
	Q            iec.BOOL
	D1           iec.BOOL
	D2           iec.BOOL
	OUT          *iec.BYTE

	t1, t2      timers.TOF
	dc1, dc2    CLICK_MODE
	t3          timers.TON
	dim         engineering.RMP_B_
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (d *DIMM_2) INIT() {
	*d = DIMM_2{
		OUT: d.OUT, VAL: 255,
		T_DEBOUNCE: iec.TIME(10 * time.Millisecond), T_DIMM_START: iec.TIME(time.Second), T_DIMM: iec.TIME(3 * time.Second),
		MIN_ON: 50, MAX_ON: 255, RST_OUT: true, SOFT_DIMM: true,
		initialized: true,
	}
	d.dc1.INIT()
	d.dc2.INIT()
	d.dim.INIT()
}

// Execute runs the block once.
func (d *DIMM_2) Execute(now time.Time) {
	if !d.initialized {
		d.initialized = true
		d.dc1.INIT()
		d.dc2.INIT()
		d.dim.INIT()
	}
	if d.OUT == nil {
		return
	}
	// Debounce I1 and I2.
	d.t1.IN, d.t1.PT = d.I1, d.T_DEBOUNCE
	d.t1.Execute(now)
	d.t2.IN, d.t2.PT = d.I2, d.T_DEBOUNCE
	d.t2.Execute(now)
	d.dc1.IN, d.dc1.T_LONG = d.t1.Q, d.T_DIMM_START
	d.dc1.Execute(now)
	d.dc2.IN, d.dc2.T_LONG = d.t2.Q, d.T_DIMM_START
	d.dc2.Execute(now)

	switch {
	case bool(d.RST):
		if d.RST_OUT {
			*d.OUT = 0
		}
		d.Q, d.D1, d.D2 = false, false, false
	case bool(d.SET):
		*d.OUT = d.VAL
		d.Q = true
	case bool(d.dc1.SINGLE_):
		// a click of I1 turns the dimmer on
		*d.OUT = LIMIT(max(d.MIN_ON, 1), *d.OUT, d.MAX_ON)
		d.Q = true
	case bool(d.dc2.SINGLE_):
		// a click of I2 turns it off
		d.Q = false
	case bool(d.dc1.TP_LONG):
		// holding I1 dims up
		if !d.Q {
			*d.OUT = SEL(d.SOFT_DIMM, LIMIT(max(d.MIN_ON, 1), *d.OUT, d.MAX_ON), 1)
		}
		d.Q = true
		d.dim.DIR = true
	case bool(d.dc2.TP_LONG):
		// holding I2 dims down
		d.dim.DIR = false
	}

	// The double clicks.
	if !d.DBL1_TOG {
		d.D1 = false
	}
	if d.dc1.DOUBLE {
		if d.DBL1_SET {
			*d.OUT = d.DBL1_POS
			d.Q = true
		} else {
			d.D1 = !d.D1
		}
	}
	if !d.DBL2_TOG {
		d.D2 = false
	}
	if d.dc2.DOUBLE {
		if d.DBL2_SET {
			*d.OUT = d.DBL2_POS
			d.Q = true
		} else {
			d.D2 = !d.D2
		}
	}

	// Dim OUT up or down while a button is held.
	d.dim.DIR = d.dc2.LONG
	d.dim.E = d.dc1.LONG || d.dc2.LONG
	d.dim.TR = d.T_DIMM
	d.dim.RMP = d.OUT
	d.dim.Execute(now)

	// OUT at 0 turns the dimmer off.
	if *d.OUT == 0 {
		d.Q = false
	}

	// The maximum time on.
	if d.T_ON_MAX > 0 {
		d.t3.IN, d.t3.PT = d.Q, d.T_ON_MAX
		d.t3.Execute(now)
		d.Q = d.Q != d.t3.Q
	}
}

// DIMM_I is a dimmer with one push button IN, reconfigured and debounced by
// a SW_RECONFIG: a click turns it on or off, and holding IN longer than
// T_DIMM_START dims OUT, over T_DIMM from 0 to 255, reversing the direction
// with each long press. SET sets OUT to VAL, RST turns the dimmer off (and
// OUT to 0 if RST_OUT). A double click toggles DBL. Q goes off after
// T_ON_MAX if that is not 0.
type DIMM_I struct {
	SET          iec.BOOL
	VAL          iec.BYTE // default 255
	IN           iec.BOOL
	RST          iec.BOOL
	T_DEBOUNCE   iec.TIME // default T#10ms
	T_RECONFIG   iec.TIME // default T#10s
	T_ON_MAX     iec.TIME
	T_DIMM_START iec.TIME // default T#1s
	T_DIMM       iec.TIME // default T#3s
	MIN_ON       iec.BYTE // default 50
	MAX_ON       iec.BYTE // default 255
	SOFT_DIMM    iec.BOOL // default TRUE
	DBL_TOGGLE   iec.BOOL
	RST_OUT      iec.BOOL
	Q            iec.BOOL
	DBL          iec.BOOL
	OUT          *iec.BYTE

	config      SW_RECONFIG
	decode      CLICK_MODE
	t3          timers.TON
	dim         engineering.RMP_B_
	dir         iec.BOOL
	initialized bool
}

// INIT resets the block and sets its inputs to their initial values.
func (d *DIMM_I) INIT() {
	*d = DIMM_I{
		OUT: d.OUT, VAL: 255,
		T_DEBOUNCE: iec.TIME(10 * time.Millisecond), T_RECONFIG: iec.TIME(10 * time.Second),
		T_DIMM_START: iec.TIME(time.Second), T_DIMM: iec.TIME(3 * time.Second),
		MIN_ON: 50, MAX_ON: 255, SOFT_DIMM: true,
		initialized: true,
	}
	d.decode.INIT()
	d.dim.INIT()
}

// Execute runs the block once.
func (d *DIMM_I) Execute(now time.Time) {
	if !d.initialized {
		d.initialized = true
		d.decode.INIT()
		d.dim.INIT()
	}
	if d.OUT == nil {
		return
	}
	// Reconfigure and debounce the input.
	d.config.IN, d.config.TD, d.config.TR = d.IN, d.T_DEBOUNCE, d.T_RECONFIG
	d.config.Execute(now)
	d.decode.IN, d.decode.T_LONG = d.config.Q, d.T_DIMM_START
	d.decode.Execute(now)

	// The direction is set reversed, as the next long press reverses it.
	switch {
	case bool(d.RST):
		if d.RST_OUT {
			*d.OUT = 0
		}
		d.Q = false
		d.dir = *d.OUT > 127
	case bool(d.SET):
		*d.OUT = d.VAL
		d.Q = true
		d.dir = *d.OUT > 127
	case bool(d.decode.SINGLE_):
		// a click toggles Q
		d.Q = !d.Q
		if d.Q {
			*d.OUT = LIMIT(max(d.MIN_ON, 1), *d.OUT, d.MAX_ON)
		}
		d.dir = *d.OUT > 127
	case bool(d.decode.TP_LONG):
		if !d.Q {
			if d.SOFT_DIMM {
				*d.OUT = 1
				d.dir = true
			} else {
				*d.OUT = LIMIT(max(d.MIN_ON, 1), *d.OUT, d.MAX_ON)
				d.dir = *d.OUT < 127
			}
			d.Q = true
		} else {
			// each long press reverses the direction
			d.dir = !d.dir
		}
	}

	// The double click.
	if !d.DBL_TOGGLE {
		d.DBL = false
	}
	if d.decode.DOUBLE {
		d.DBL = !d.DBL
	}

	// Dim OUT up or down while the button is held.
	d.dim.DIR = d.dir
	d.dim.E = d.decode.LONG && d.Q
	d.dim.TR = d.T_DIMM
	d.dim.RMP = d.OUT
	d.dim.Execute(now)

	// The limits reverse the direction.
	if *d.OUT == 0 {
		d.dir = true
	} else if *d.OUT == 255 {
		d.dir = false
	}

	// The maximum time on.
	if d.T_ON_MAX > 0 {
		d.t3.IN, d.t3.PT = d.Q, d.T_ON_MAX
		d.t3.Execute(now)
		d.Q = d.Q != d.t3.Q
	}
}

// F_LAMP drives a fluorescent lamp: LAMP is 255 while SWITCH is on for the
// first T_NO_DIMM hours of the lamp's life, DIMM after. It counts the
// ONTIME in seconds and the CYCLES, which RST clears, and STATUS is 120
// once ONTIME passes T_MAINTENANCE hours.
type F_LAMP struct {
	SWITCH        iec.BOOL
	DIMM          iec.BYTE // default 255
	RST           iec.BOOL
	LAMP          iec.BYTE
	STATUS        iec.BYTE
	ONTIME        *iec.UDINT
	CYCLES        *iec.UDINT
	T_NO_DIMM     iec.UINT // default 100
	T_MAINTENANCE iec.UINT // default 15000

	runtime engineering.ONTIME
}

// INIT resets the block and sets its inputs to their initial values.
func (f *F_LAMP) INIT() {
	*f = F_LAMP{ONTIME: f.ONTIME, CYCLES: f.CYCLES, DIMM: 255, T_NO_DIMM: 100, T_MAINTENANCE: 15000}
}

// Execute runs the block once.
func (f *F_LAMP) Execute(now time.Time) {
	if f.ONTIME == nil || f.CYCLES == nil {
		return
	}
	f.runtime.IN = f.SWITCH
	f.runtime.SECONDS, f.runtime.CYCLES = f.ONTIME, f.CYCLES
	f.runtime.Execute(now)

	if f.RST {
		*f.ONTIME = 0
		*f.CYCLES = 0
	}
	switch {
	case bool(!f.SWITCH):
		f.LAMP = 0
		f.STATUS = 110
	case *f.ONTIME < iec.UDINT(f.T_NO_DIMM)*3600:
		f.LAMP = 255
		f.STATUS = 111
	default:
		f.LAMP = f.DIMM
		f.STATUS = 112
	}
	if *f.ONTIME >= iec.UDINT(f.T_MAINTENANCE)*3600 && f.T_MAINTENANCE > 0 {
		f.STATUS = 120
	}
}
