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

package engineering

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// A function's VAR_INPUT CONSTANTs are parameters after its inputs, in
// OSCAT's order; their initial values are in the comments, for a caller
// that leaves them out.

// AIN converts the lowest bits bits of an A/D converter's value to a REAL
// from low to high, negative if the bit sign is set. sign is 255, for none,
// and high 10.0 by default.
func AIN(in iec.DWORD, bits, sign iec.BYTE, low, high iec.REAL) iec.REAL {
	sx := sign < 32 && SHR(in, sign)&1 != 0
	temp1 := SHR(iec.DWORD(0xFFFFFFFF), 32-int(bits))
	temp2 := in & temp1
	out := (high-low)*iec.REAL(temp2)/iec.REAL(temp1) + low
	if sx {
		return -out
	}
	return out
}

// AIN1 converts the bits BIT_0..BIT_N of an A/D converter's value, from
// CODE_MIN..CODE_MAX, to OUT from OUT_MIN to OUT_MAX, negative if the bit
// SIGN_BIT is set. The bit ERROR_BIT or the value ERROR_CODE is an error,
// giving ERROR_OUTPUT, and the bit OVERFLOW_BIT, the value OVERFLOW_CODE or
// a value outside CODE_MIN..CODE_MAX is an overflow, giving
// OVERFLOW_OUTPUT. A bit number of 255 is none.
type AIN1 struct {
	IN               iec.DWORD
	SIGN_BIT         iec.INT // default 255
	ERROR_BIT        iec.INT // default 255
	ERROR_CODE_EN    iec.BOOL
	ERROR_CODE       iec.DWORD
	OVERFLOW_BIT     iec.INT // default 255
	OVERFLOW_CODE_EN iec.BOOL
	OVERFLOW_CODE    iec.DWORD
	BIT_0            iec.INT
	BIT_N            iec.INT // default 31
	OUT_MIN          iec.REAL
	OUT_MAX          iec.REAL // default 10.0
	CODE_MIN         iec.DWORD
	CODE_MAX         iec.DWORD // default 16#FFFFFFFF
	ERROR_OUTPUT     iec.REAL
	OVERFLOW_OUTPUT  iec.REAL // default 10.0
	OUT              iec.REAL
	SIGN             iec.BOOL
	ERROR            iec.BOOL
	OVERFLOW         iec.BOOL
}

// INIT resets the block and sets its inputs to their initial values.
func (a *AIN1) INIT() {
	*a = AIN1{SIGN_BIT: 255, ERROR_BIT: 255, OVERFLOW_BIT: 255, BIT_N: 31, OUT_MAX: 10,
		CODE_MAX: 0xFFFFFFFF, OVERFLOW_OUTPUT: 10}
}

// Execute runs the block once.
func (a *AIN1) Execute(now time.Time) {
	a.ERROR = iec.BOOL(SHR(a.IN, a.ERROR_BIT)&1 == 1 || bool(a.ERROR_CODE_EN) && a.ERROR_CODE == a.IN)
	if a.ERROR {
		a.OUT = a.ERROR_OUTPUT
		return
	}
	tb := SHR(SHL(a.IN, 31-a.BIT_N), 31-a.BIT_N+a.BIT_0)
	a.OVERFLOW = iec.BOOL(SHR(a.IN, a.OVERFLOW_BIT)&1 == 1 || bool(a.OVERFLOW_CODE_EN) && a.OVERFLOW_CODE == a.IN ||
		tb < a.CODE_MIN || tb > a.CODE_MAX)
	if a.OVERFLOW {
		a.OUT = a.OVERFLOW_OUTPUT
		return
	}
	a.SIGN = SHR(a.IN, a.SIGN_BIT)&1 == 1
	a.OUT = iec.REAL(tb-a.CODE_MIN)*(a.OUT_MAX-a.OUT_MIN)/iec.REAL(a.CODE_MAX-a.CODE_MIN) + a.OUT_MIN
	if a.SIGN {
		a.OUT = -a.OUT
	}
}

// AOUT converts in, from low to high, to a value of bits bits for a D/A
// converter, with the bit sign set for a negative value. sign is 255, for
// none, and high 10.0 by default.
func AOUT(in iec.REAL, bits, sign iec.BYTE, low, high iec.REAL) iec.DWORD {
	var sx iec.BOOL
	in2 := in
	if sign < 32 {
		sx = math.SIGN_R(in)
		in2 = ABS(in)
	}
	in2 = LIMIT(low, in2, high)
	out := REAL_TO_DWORD((in2 - low) / (high - low) * iec.REAL(SHL(iec.DWORD(1), bits)-1))
	if sx {
		out |= SHL(iec.DWORD(1), sign)
	}
	return out
}

// AOUT1 converts in, from low to high, to the bits bit0..bitN of a value
// for a D/A converter, with the bit sign set for a negative value. bitN is
// 31, sign 255, for none, and high 10.0 by default.
func AOUT1(in iec.REAL, bit0, bitN, sign iec.INT, low, high iec.REAL) iec.DWORD {
	var sx iec.BOOL
	in2 := in
	if sign < 32 {
		sx = math.SIGN_R(in)
		in2 = ABS(in)
	}
	in2 = LIMIT(low, in2, high)
	out := SHL(REAL_TO_DWORD((in2-low)/(high-low)*iec.REAL(SHL(iec.DWORD(1), bitN-bit0+1)-1)), bit0)
	if sx {
		out |= SHL(iec.DWORD(1), sign)
	}
	return out
}

// BYTE_TO_RANGE converts a byte to a REAL from low to high.
func BYTE_TO_RANGE(x iec.BYTE, low, high iec.REAL) iec.REAL {
	return (high-low)*iec.REAL(x)/255.0 + low
}

// DELAY delays IN by N scans, 0..32. RST loads the delay with IN.
type DELAY struct {
	IN  iec.REAL
	N   iec.INT
	RST iec.BOOL
	OUT iec.REAL

	buf  [32]iec.REAL
	i    iec.INT
	init iec.BOOL
}

// INIT resets the block.
func (d *DELAY) INIT() { *d = DELAY{} }

// Execute runs the block once.
func (d *DELAY) Execute(now time.Time) {
	stop := LIMIT(0, d.N, 32) - 1
	switch {
	case bool(d.RST || !d.init):
		d.init = true
		for i := iec.INT(0); i <= stop; i++ {
			d.buf[i] = d.IN
		}
		d.OUT = d.IN
		d.i = 0
	case stop < 0:
		d.OUT = d.IN
	default:
		d.OUT = d.buf[d.i]
		d.buf[d.i] = d.IN
		d.i = math.INC1(d.i, d.N)
	}
}

// DELAY_4 delays IN by 1 to 4 scans on OUT1..OUT4.
type DELAY_4 struct {
	IN                     iec.REAL
	OUT1, OUT2, OUT3, OUT4 iec.REAL

	temp iec.REAL
}

// INIT resets the block.
func (d *DELAY_4) INIT() { *d = DELAY_4{} }

// Execute runs the block once.
func (d *DELAY_4) Execute(now time.Time) {
	d.OUT4, d.OUT3, d.OUT2, d.OUT1, d.temp = d.OUT3, d.OUT2, d.OUT1, d.temp, d.IN
}

// FADE fades Y from IN1 to IN2 over TF while F is true, and back while it
// is false. RST sets Y to IN1 or IN2 at once.
type FADE struct {
	IN1, IN2 iec.REAL
	F        iec.BOOL
	TF       iec.TIME
	RST      iec.BOOL
	Y        iec.REAL

	rmx         RMP_W
	initialized bool
}

// INIT resets the block.
func (f *FADE) INIT() {
	*f = FADE{initialized: true}
	f.rmx.INIT()
}

// Execute runs the block once.
func (f *FADE) Execute(now time.Time) {
	if !f.initialized {
		f.initialized = true
		f.rmx.INIT()
	}
	f.rmx.RST, f.rmx.SET, f.rmx.PT, f.rmx.UP = f.RST && !f.F, f.RST && f.F, f.TF, f.F
	f.rmx.Execute(now)
	f.Y = (f.IN2-f.IN1)/65535.0*iec.REAL(f.rmx.OUT) + f.IN1
}

// FILTER_DW is a low pass filter with the time T for DWORD values.
type FILTER_DW struct {
	X iec.DWORD
	T iec.TIME
	Y iec.DWORD

	last iec.DWORD
	init iec.BOOL
	yi   iec.REAL
}

// INIT resets the block.
func (f *FILTER_DW) INIT() { *f = FILTER_DW{} }

// Execute runs the block once.
func (f *FILTER_DW) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !f.init || f.T == 0 {
		f.init = true
		f.yi = iec.REAL(f.X)
	} else {
		f.yi = f.yi + (iec.REAL(f.X)-iec.REAL(f.Y))*iec.REAL(tx-f.last)/TIME_TO_REAL(f.T)
	}
	f.last = tx
	f.Y = REAL_TO_DWORD(f.yi)
}

// FILTER_I is a low pass filter with the time T for INT values.
type FILTER_I struct {
	X iec.INT
	T iec.TIME
	Y iec.INT

	yi   iec.DINT
	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block.
func (f *FILTER_I) INIT() { *f = FILTER_I{} }

// Execute runs the block once.
func (f *FILTER_I) Execute(now time.Time) {
	tx := PLC_MS(now)
	if !f.init || f.T == 0 {
		f.init = true
		f.yi = iec.DINT(f.X) * 1000
	} else {
		f.yi = f.yi + iec.DINT(f.X-f.Y)*iec.DINT(tx-f.last)*1000/iec.DINT(ms(f.T))
	}
	f.last = tx
	f.Y = iec.INT(f.yi / 1000)
}

// FILTER_MAV_DW is a moving average over the last N values of X, up to 32,
// for DWORD values.
type FILTER_MAV_DW struct {
	X   iec.DWORD
	N   iec.UINT
	RST iec.BOOL
	Y   iec.DWORD

	init   iec.BOOL
	buffer [32]iec.DWORD
	i      iec.INT
}

// INIT resets the block.
func (f *FILTER_MAV_DW) INIT() { *f = FILTER_MAV_DW{} }

// Execute runs the block once.
func (f *FILTER_MAV_DW) Execute(now time.Time) {
	f.N = min(f.N, 32)
	if !f.init || f.RST || f.N == 0 {
		f.init = true
		for i := 0; i < int(f.N); i++ {
			f.buffer[i] = f.X
		}
		f.Y = f.X
		return
	}
	f.i = math.INC1(f.i, iec.INT(f.N))
	f.Y = f.Y + (f.X-f.buffer[f.i])/iec.DWORD(f.N)
	f.buffer[f.i] = f.X
}

// FILTER_MAV_W is a moving average over the last N values of X, up to 32,
// for WORD values.
//
// OSCAT's buffer is ARRAY[1..32] but the block uses positions 0..N-1, and it
// starts its sum from the previous Y; the port keeps both.
type FILTER_MAV_W struct {
	X   iec.WORD
	N   iec.UINT
	RST iec.BOOL
	Y   iec.WORD

	init   iec.BOOL
	buffer [33]iec.WORD
	i      iec.INT
	sum    iec.DWORD
}

// INIT resets the block.
func (f *FILTER_MAV_W) INIT() { *f = FILTER_MAV_W{} }

// Execute runs the block once.
func (f *FILTER_MAV_W) Execute(now time.Time) {
	f.N = min(f.N, 32)
	if !f.init || f.RST || f.N == 0 {
		f.init = true
		for i := 1; i <= int(f.N)-1; i++ {
			f.buffer[i] = f.X
		}
		f.sum = iec.DWORD(f.Y) * iec.DWORD(f.N)
		f.Y = f.X
		return
	}
	f.i = math.INC1(f.i, iec.INT(f.N))
	f.sum = f.sum + iec.DWORD(f.X) - iec.DWORD(f.buffer[f.i])
	f.Y = iec.WORD(f.sum / iec.DWORD(f.N))
	f.buffer[f.i] = f.X
}

// FILTER_W is a low pass filter with the time T for WORD values.
type FILTER_W struct {
	X iec.WORD
	T iec.TIME
	Y iec.WORD

	last iec.DWORD
	init iec.BOOL
}

// INIT resets the block.
func (f *FILTER_W) INIT() { *f = FILTER_W{} }

// Execute runs the block once.
func (f *FILTER_W) Execute(now time.Time) {
	tx := PLC_MS(now)
	switch {
	case !bool(f.init) || f.T == 0:
		f.init = true
		f.last = tx
		f.Y = f.X
	case f.Y == f.X:
		f.last = tx
	default:
		tmp := iec.DWORD(f.X-f.Y) * (tx - f.last) / ms(f.T)
		if tmp != 0 {
			f.Y = iec.WORD(int32(f.Y) + int32(tmp))
			f.last = tx
		}
	}
}

// FILTER_WAV is a moving average over the last 16 values of X, weighted by
// W: the newest value by W[0].
type FILTER_WAV struct {
	X   iec.REAL
	W   [16]iec.REAL
	RST iec.BOOL
	Y   iec.REAL

	init   iec.BOOL
	buffer [16]iec.REAL
	i      iec.INT
}

// INIT resets the block.
func (f *FILTER_WAV) INIT() { *f = FILTER_WAV{} }

// Execute runs the block once.
func (f *FILTER_WAV) Execute(now time.Time) {
	if !f.init || f.RST {
		f.init = true
		for i := range f.buffer {
			f.buffer[i] = f.X
		}
		f.i = 15
	} else {
		f.i = math.INC1(f.i, 16)
		f.buffer[f.i] = f.X
	}
	f.Y = 0
	for n := 0; n < 16; n++ {
		f.Y = f.buffer[f.i]*f.W[n] + f.Y
		f.i = math.DEC1(f.i, 16)
	}
}

// MIX mixes a and b: (1 - m) * a + m * b.
func MIX(a, b, m iec.REAL) iec.REAL {
	return (1.0-m)*a + m*b
}

// MUX_R2 returns in0 if a is false and in1 if it is true.
func MUX_R2(in0, in1 iec.REAL, a iec.BOOL) iec.REAL {
	return SEL(a, in0, in1)
}

// MUX_R4 returns the input in0..in3 the address a1, a0 selects.
func MUX_R4(in0, in1, in2, in3 iec.REAL, a0, a1 iec.BOOL) iec.REAL {
	if a1 {
		return SEL(a0, in2, in3)
	}
	return SEL(a0, in0, in1)
}

// OFFSET returns x, or def if d is true, plus each offset whose input
// o1..o4 is true.
func OFFSET(x iec.REAL, o1, o2, o3, o4, d iec.BOOL, offset1, offset2, offset3, offset4, def iec.REAL) iec.REAL {
	out := SEL(d, x, def)
	for _, o := range []struct {
		on iec.BOOL
		v  iec.REAL
	}{{o1, offset1}, {o2, offset2}, {o3, offset3}, {o4, offset4}} {
		if o.on {
			out += o.v
		}
	}
	return out
}

// OFFSET2 returns x, or def if d is true, plus the offset of the highest
// input o1..o4 that is true.
func OFFSET2(x iec.REAL, o1, o2, o3, o4, d iec.BOOL, offset1, offset2, offset3, offset4, def iec.REAL) iec.REAL {
	out := SEL(d, x, def)
	switch {
	case bool(o4):
		out += offset4
	case bool(o3):
		out += offset3
	case bool(o2):
		out += offset2
	case bool(o1):
		out += offset1
	}
	return out
}

// OVERRIDE returns the input x1..x3 of the greatest absolute value whose
// enable e1..e3 is true, or 0.
func OVERRIDE(x1, x2, x3 iec.REAL, e1, e2, e3 iec.BOOL) iec.REAL {
	var out iec.REAL
	if e1 {
		out = x1
	}
	if e2 && ABS(x2) > ABS(out) {
		out = x2
	}
	if e3 && ABS(x3) > ABS(out) {
		out = x3
	}
	return out
}

// RANGE_TO_BYTE converts x from low to high to a byte.
func RANGE_TO_BYTE(x, low, high iec.REAL) iec.BYTE {
	return iec.BYTE(TRUNC((LIMIT(low, x, high) - low) * 255.0 / (high - low)))
}

// RANGE_TO_WORD converts x from low to high to a word.
func RANGE_TO_WORD(x, low, high iec.REAL) iec.WORD {
	return iec.WORD(TRUNC((LIMIT(low, x, high) - low) * 65535.0 / (high - low)))
}

// SCALE returns x * k + o, limited to mn..mx.
func SCALE(x, k, o, mx, mn iec.REAL) iec.REAL {
	return LIMIT(mn, x*k+o, mx)
}

// SCALE_B scales a byte x from iLo..iHi to oLo..oHi, as OSCAT does, which
// neither subtracts iLo nor adds oLo.
func SCALE_B(x, iLo, iHi iec.BYTE, oLo, oHi iec.REAL) iec.REAL {
	if iHi == iLo {
		return oLo
	}
	return (oHi - oLo) / iec.REAL(iHi-iLo) * iec.REAL(LIMIT(iLo, x, iHi))
}

// SCALE_B2 scales 2 bytes, each 0..255 to its range inN_min..inN_max, sums
// them, and returns the sum * k + o. The maxima are 1000.0 by default.
func SCALE_B2(in1, in2 iec.BYTE, k, o, in1Min, in1Max, in2Min, in2Max iec.REAL) iec.REAL {
	return (((in1Max-in1Min)*iec.REAL(in1)+(in2Max-in2Min)*iec.REAL(in2))*0.003921569+in1Min+in2Min)*k + o
}

// SCALE_B4 scales 4 bytes; see SCALE_B2.
func SCALE_B4(in1, in2, in3, in4 iec.BYTE, k, o, in1Min, in1Max, in2Min, in2Max, in3Min, in3Max, in4Min, in4Max iec.REAL) iec.REAL {
	return (((in1Max-in1Min)*iec.REAL(in1)+(in2Max-in2Min)*iec.REAL(in2)+(in3Max-in3Min)*iec.REAL(in3)+
		(in4Max-in4Min)*iec.REAL(in4))*0.003921569+in1Min+in2Min+in3Min+in4Min)*k + o
}

// SCALE_B8 scales 8 bytes; see SCALE_B2.
func SCALE_B8(in1, in2, in3, in4, in5, in6, in7, in8 iec.BYTE, k, o,
	in1Min, in1Max, in2Min, in2Max, in3Min, in3Max, in4Min, in4Max,
	in5Min, in5Max, in6Min, in6Max, in7Min, in7Max, in8Min, in8Max iec.REAL) iec.REAL {
	return (((in1Max-in1Min)*iec.REAL(in1)+(in2Max-in2Min)*iec.REAL(in2)+(in3Max-in3Min)*iec.REAL(in3)+
		(in4Max-in4Min)*iec.REAL(in4)+(in5Max-in5Min)*iec.REAL(in5)+(in6Max-in6Min)*iec.REAL(in6)+
		(in7Max-in7Min)*iec.REAL(in7)+(in8Max-in8Min)*iec.REAL(in8))*0.003921569+
		in1Min+in2Min+in3Min+in4Min+in5Min+in6Min+in7Min+in8Min)*k + o
}

// SCALE_D scales a DWORD x from iLo..iHi, to which it is limited, to
// oLo..oHi.
func SCALE_D(x, iLo, iHi iec.DWORD, oLo, oHi iec.REAL) iec.REAL {
	if iHi == iLo {
		return oLo
	}
	return (oHi-oLo)/iec.REAL(iHi-iLo)*iec.REAL(LIMIT(iLo, x, iHi)-iLo) + oLo
}

// SCALE_R scales a REAL x from iLo..iHi, to which it is limited, to
// oLo..oHi.
func SCALE_R(x, iLo, iHi, oLo, oHi iec.REAL) iec.REAL {
	if iLo == iHi {
		return oLo
	}
	return (oHi-oLo)/(iHi-iLo)*(LIMIT(iLo, x, iHi)-iLo) + oLo
}

// SCALE_X2 sums, for 2 inputs, inN_max if the input is true and inN_min if
// it is false, and returns the sum * k + o. The maxima are 1000.0 by
// default.
func SCALE_X2(in1, in2 iec.BOOL, k, o, in1Min, in1Max, in2Min, in2Max iec.REAL) iec.REAL {
	return (SEL(in1, in1Min, in1Max)+SEL(in2, in2Min, in2Max))*k + o
}

// SCALE_X4 is SCALE_X2 for 4 inputs.
func SCALE_X4(in1, in2, in3, in4 iec.BOOL, k, o, in1Min, in1Max, in2Min, in2Max, in3Min, in3Max, in4Min, in4Max iec.REAL) iec.REAL {
	return (SEL(in1, in1Min, in1Max)+SEL(in2, in2Min, in2Max)+SEL(in3, in3Min, in3Max)+SEL(in4, in4Min, in4Max))*k + o
}

// SCALE_X8 is SCALE_X2 for 8 inputs.
func SCALE_X8(in1, in2, in3, in4, in5, in6, in7, in8 iec.BOOL, k, o,
	in1Min, in1Max, in2Min, in2Max, in3Min, in3Max, in4Min, in4Max,
	in5Min, in5Max, in6Min, in6Max, in7Min, in7Max, in8Min, in8Max iec.REAL) iec.REAL {
	return (SEL(in1, in1Min, in1Max)+SEL(in2, in2Min, in2Max)+SEL(in3, in3Min, in3Max)+SEL(in4, in4Min, in4Max)+
		SEL(in5, in5Min, in5Max)+SEL(in6, in6Min, in6Max)+SEL(in7, in7Min, in7Max)+SEL(in8, in8Min, in8Max))*k + o
}

// SEL2_OF_3 averages 3 redundant signals that are within D of each other.
// A signal that is not is left out and its number shown on W; if no two
// agree, E is true, W is 4 and Y keeps its value.
type SEL2_OF_3 struct {
	IN1, IN2, IN3, D iec.REAL
	Y                iec.REAL
	W                iec.INT
	E                iec.BOOL
}

// INIT resets the block.
func (s *SEL2_OF_3) INIT() { *s = SEL2_OF_3{} }

// Execute runs the block once.
func (s *SEL2_OF_3) Execute(now time.Time) {
	d12 := ABS(s.IN1-s.IN2) <= s.D
	d23 := ABS(s.IN2-s.IN3) <= s.D
	d31 := ABS(s.IN3-s.IN1) <= s.D
	switch {
	case d12 && d23 || d12 && d31 || d23 && d31:
		s.Y, s.E, s.W = (s.IN1+s.IN2+s.IN3)*0.333333333333, false, 0
	case d12:
		s.Y, s.E, s.W = (s.IN1+s.IN2)*0.5, false, 3
	case d23:
		s.Y, s.E, s.W = (s.IN2+s.IN3)*0.5, false, 1
	case d31:
		s.Y, s.E, s.W = (s.IN3+s.IN1)*0.5, false, 2
	default:
		s.E, s.W = true, 4
	}
}

// SEL2_OF_3B is the majority of 3 redundant binary signals. W is true when
// they have disagreed for TD.
type SEL2_OF_3B struct {
	IN1, IN2, IN3 iec.BOOL
	TD            iec.TIME
	Q             iec.BOOL
	W             iec.BOOL

	tdel timers.TON
}

// INIT resets the block.
func (s *SEL2_OF_3B) INIT() { *s = SEL2_OF_3B{} }

// Execute runs the block once.
func (s *SEL2_OF_3B) Execute(now time.Time) {
	s.Q = s.IN1 && s.IN2 || s.IN1 && s.IN3 || s.IN2 && s.IN3
	s.tdel.IN = s.IN1 != s.IN2 || s.IN1 != s.IN3 || s.IN2 != s.IN3
	s.tdel.PT = s.TD
	s.tdel.Execute(now)
	s.W = s.tdel.Q
}

// SH samples IN on a rising edge of CLK and holds it in OUT. TRIG is true
// for one scan when it samples.
type SH struct {
	IN   iec.REAL
	CLK  iec.BOOL
	OUT  iec.REAL
	TRIG iec.BOOL

	edge iec.BOOL
}

// INIT resets the block.
func (s *SH) INIT() { *s = SH{} }

// Execute runs the block once.
func (s *SH) Execute(now time.Time) {
	if s.CLK && !s.edge {
		s.OUT = s.IN
		s.TRIG = true
	} else {
		s.TRIG = false
	}
	s.edge = s.CLK
}

// SH_1 samples IN every PT and holds it in OUT. TRIG is true for one scan
// when it samples.
type SH_1 struct {
	IN   iec.REAL
	PT   iec.TIME
	OUT  iec.REAL
	TRIG iec.BOOL

	last iec.DWORD
}

// INIT resets the block.
func (s *SH_1) INIT() { *s = SH_1{} }

// Execute runs the block once.
func (s *SH_1) Execute(now time.Time) {
	tx := PLC_MS(now)
	if tx-s.last >= ms(s.PT) {
		s.last = tx
		s.OUT = s.IN
		s.TRIG = true
	} else {
		s.TRIG = false
	}
}

// SH_2 samples IN every PT and holds it in OUT, and gives the average, the
// lowest and the highest of the last N samples, up to 16, leaving out DISC
// samples: the lowest, then the highest, and so on.
type SH_2 struct {
	IN   iec.REAL
	PT   iec.TIME
	N    iec.INT // default 16
	DISC iec.INT
	OUT  iec.REAL
	TRIG iec.BOOL
	AVG  iec.REAL
	HIGH iec.REAL
	LOW  iec.REAL

	buf2 [16]iec.REAL
	last iec.DWORD
}

// INIT resets the block and sets N to its initial value.
func (s *SH_2) INIT() { *s = SH_2{N: 16} }

// Execute runs the block once.
func (s *SH_2) Execute(now time.Time) {
	tx := PLC_MS(now)
	d2 := s.DISC >> 1
	if tx-s.last < ms(s.PT) {
		s.TRIG = false
		return
	}
	s.last = tx
	s.TRIG = true
	m := LIMIT(1, s.N, 16)
	for i := m - 1; i >= 1; i-- {
		s.buf2[i] = s.buf2[i-1]
	}
	s.buf2[0] = s.IN
	s.OUT = s.IN
	buf := s.buf2
	for start := iec.INT(0); start <= m-2; start++ {
		for i := start + 1; i <= m-1; i++ {
			if buf[start] > buf[i] {
				buf[start], buf[i] = buf[i], buf[start]
			}
		}
	}
	stop := m - 1 - d2
	start := d2
	if !math.EVEN(iec.DINT(s.DISC)) {
		start++
	}
	s.AVG = 0
	for i := start; i <= stop; i++ {
		s.AVG += buf[LIMIT(0, i, 15)]
	}
	s.AVG /= iec.REAL(stop - start + 1)
	s.LOW = buf[LIMIT(0, start, 15)]
	s.HIGH = buf[LIMIT(0, stop, 15)]
}

// SH_T follows IN with OUT while E is true and holds it while E is false.
type SH_T struct {
	IN  iec.REAL
	E   iec.BOOL
	OUT iec.REAL
}

// INIT resets the block.
func (s *SH_T) INIT() { *s = SH_T{} }

// Execute runs the block once.
func (s *SH_T) Execute(now time.Time) {
	if s.E {
		s.OUT = s.IN
	}
}

// STAIR rounds x to steps of d, or returns x if d is not above 0.
func STAIR(x, d iec.REAL) iec.REAL {
	if d > 0.0 {
		return iec.REAL(REAL_TO_DINT(x/d)) * d
	}
	return x
}

// STAIR2 follows X in steps of D, with D as a hysteresis.
type STAIR2 struct {
	X, D iec.REAL
	Y    iec.REAL
}

// INIT resets the block.
func (s *STAIR2) INIT() { *s = STAIR2{} }

// Execute runs the block once.
func (s *STAIR2) Execute(now time.Time) {
	if s.D > 0.0 {
		if s.X >= s.Y+s.D || s.X <= s.Y-s.D {
			s.Y = iec.REAL(math.FLOOR(s.X/s.D)) * s.D
		}
	} else {
		s.Y = s.X
	}
}

// TREND shows how X changes: TU for one scan when it rises, TD when it
// falls, Q when it changes, and D the change.
type TREND struct {
	X      iec.REAL
	Q      iec.BOOL
	TU, TD iec.BOOL
	D      iec.REAL

	lastX iec.REAL
}

// INIT resets the block.
func (t *TREND) INIT() { *t = TREND{} }

// Execute runs the block once.
func (t *TREND) Execute(now time.Time) {
	t.TU = t.X > t.lastX
	t.TD = t.X < t.lastX
	t.Q = t.TU || t.TD
	t.D = t.X - t.lastX
	t.lastX = t.X
}

// TREND_DW shows how X changes: TU for one scan when it rises, TD when it
// falls, Q true after a rise and false after a fall, and D the change.
type TREND_DW struct {
	X      iec.DWORD
	Q      iec.BOOL
	TU, TD iec.BOOL
	D      iec.DWORD

	lastX iec.DWORD
}

// INIT resets the block.
func (t *TREND_DW) INIT() { *t = TREND_DW{} }

// Execute runs the block once.
func (t *TREND_DW) Execute(now time.Time) {
	switch {
	case t.X > t.lastX:
		t.TU, t.TD, t.D, t.Q = true, false, t.X-t.lastX, true
	case t.X < t.lastX:
		t.TD, t.TU, t.D, t.Q = true, false, t.lastX-t.X, false
	default:
		t.TU, t.TD, t.D = false, false, 0
	}
	t.lastX = t.X
}

// WORD_TO_RANGE converts a word to a REAL from low to high.
func WORD_TO_RANGE(x iec.WORD, low, high iec.REAL) iec.REAL {
	return (high-low)*iec.REAL(x)*0.00001525902189669640 + low
}
