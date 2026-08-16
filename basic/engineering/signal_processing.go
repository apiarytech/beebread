/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package engineering

import (
	"math"
	"sort"
	"time"

	"beebread/basic/logic"
	beeMath "beebread/basic/math"
)

// AIN converts signals from A/D converters to a real value.
func AIN(in uint32, bits, sign byte, low, high float64) float64 {
	var sx bool
	if sign < 32 {
		sx = (in>>sign)&1 == 1
	}

	mask := (uint32(1) << bits) - 1
	val := in & mask

	var result float64
	if mask > 0 {
		result = (high-low)*float64(val)/float64(mask) + low
	} else {
		result = low
	}

	if sx {
		return -result
	}
	return result
}

// AIN1 converts signals from A/D converters to a real value with error handling.
type AIN1 struct {
	Out      float64
	Sign     bool
	Error    bool
	Overflow bool
}

// Update executes the conversion logic.
func (a *AIN1) Update(in, errorCode, overflowCode, codeMin, codeMax uint32, signBit, errorBit, overflowBit, bit0, bitN int, errCodeEn, ovfCodeEn bool, outMin, outMax, errorOut, overflowOut float64) {
	a.Error = (logic.BIT_OF_DWORD(in, uint(errorBit))) || (errCodeEn && errorCode == in)
	if a.Error {
		a.Out = errorOut
		return
	}

	// Strip off the data input
	tb := (in << (31 - uint(bitN))) >> (31 - uint(bitN) + uint(bit0))

	a.Overflow = (logic.BIT_OF_DWORD(in, uint(overflowBit))) || (ovfCodeEn && overflowCode == in) || (tb < codeMin || tb > codeMax)
	if a.Overflow {
		a.Out = overflowOut
		return
	}

	a.Sign = logic.BIT_OF_DWORD(in, uint(signBit))

	// Convert in to out
	if codeMax-codeMin > 0 {
		a.Out = (float64(tb-codeMin)*(outMax-outMin)/float64(codeMax-codeMin) + outMin)
	} else {
		a.Out = outMin
	}

	if a.Sign {
		a.Out *= -1.0
	}
}

// AOUT converts a real value for a D/A converter.
func AOUT(in, low, high float64, bits, sign byte) uint32 {
	var sx bool
	var in2 = in

	if sign < 32 {
		sx = in < 0.0
		in2 = math.Abs(in)
	}

	in2 = beeMath.LIMIT(low, in2, high)

	var result uint32
	if high-low != 0.0 {
		mask := (uint32(1) << bits) - 1
		result = uint32((in2 - low) / (high - low) * float64(mask))
	}

	if sx {
		result |= (1 << sign)
	}
	return result
}

// AOUT1 converts a real value for a D/A converter with bit field placement.
func AOUT1(in, low, high float64, bit0, bitN, sign int) uint32 {
	var sx bool
	var in2 = in

	if sign < 32 {
		sx = in < 0.0
		in2 = math.Abs(in)
	}

	in2 = beeMath.LIMIT(low, in2, high)

	var result uint32
	if high-low != 0.0 {
		mask := (uint32(1) << (bitN - bit0 + 1)) - 1
		result = uint32((in2-low)/(high-low)*float64(mask)) << uint(bit0)
	}

	if sx {
		result |= (1 << uint(sign))
	}
	return result
}

// BYTE_TO_RANGE converts a byte into a real value between low and high.
func BYTE_TO_RANGE(x byte, low, high float64) float64 {
	return (high-low)*float64(x)/255.0 + low
}

// DELAY delays input values by N program cycles.
type DELAY struct {
	Out float64

	// internal state
	buf  []float64
	i    int
	init bool
}

// Update executes the delay logic.
func (d *DELAY) Update(in float64, n int, rst bool) {
	n = int(beeMath.LIMIT(0, float64(n), 32))

	if rst || !d.init || len(d.buf) != n {
		d.init = true
		d.buf = make([]float64, n)
		for j := 0; j < n; j++ {
			d.buf[j] = in
		}
		d.Out = in
		d.i = 0
	} else if n > 0 {
		d.Out = d.buf[d.i]
		d.buf[d.i] = in
		d.i = beeMath.INC1(d.i, n)
	} else {
		d.Out = in
	}
}

// DELAY_4 delays input values by 4 program cycles.
type DELAY_4 struct {
	Out1, Out2, Out3, Out4 float64
	temp                   float64
}

// Update executes the delay logic.
func (d *DELAY_4) Update(in float64) {
	d.Out4 = d.Out3
	d.Out3 = d.Out2
	d.Out2 = d.Out1
	d.Out1 = d.temp
	d.temp = in
}

// FILTER_DW is a low-pass filter for a DWORD value.
type FILTER_DW struct {
	Y    uint32
	yi   float64
	last int64
	init bool
}

// Update executes the filter logic.
func (f *FILTER_DW) Update(x uint32, t time.Duration) {
	tx := logic.T_PLC_US() / 1000 // T_PLC_MS

	if !f.init || t == 0 {
		f.init = true
		f.yi = float64(x)
	} else {
		f.yi += (float64(x) - float64(f.Y)) * float64(tx-f.last) / t.Seconds() / 1000.0
	}
	f.last = tx
	f.Y = uint32(f.yi)
}

// FILTER_I is a low-pass filter for an INT value.
type FILTER_I struct {
	Y    int
	yi   int64
	last int64
	init bool
}

// Update executes the filter logic.
func (f *FILTER_I) Update(x int, t time.Duration) {
	tx := logic.T_PLC_US() / 1000 // T_PLC_MS

	if !f.init || t == 0 {
		f.init = true
		f.yi = int64(x) * 1000
	} else if t > 0 {
		f.yi += (int64(x-f.Y) * (tx - f.last) * 1000) / t.Milliseconds()
	}
	f.last = tx
	f.Y = int(f.yi / 1000)
}

// FILTER_MAV_DW is a moving average filter for DWORD data.
type FILTER_MAV_DW struct {
	Y      uint32
	init   bool
	buffer [32]uint32
	i      int
}

// Update executes the filter logic.
func (f *FILTER_MAV_DW) Update(x uint32, n int, rst bool) {
	n = int(beeMath.LIMIT(0, float64(n), 32))

	if !f.init || rst || n == 0 {
		f.init = true
		for j := 0; j < n; j++ {
			f.buffer[j] = x
		}
		f.Y = x
		f.i = 0
	} else {
		f.i = beeMath.INC1(f.i, n)
		f.Y = f.Y + (x-f.buffer[f.i])/uint32(n)
		f.buffer[f.i] = x
	}
}

// FILTER_MAV_W is a moving average filter for WORD data.
type FILTER_MAV_W struct {
	Y      uint16
	init   bool
	buffer [32]uint16
	i      int
	sum    uint32
}

// Update executes the filter logic.
func (f *FILTER_MAV_W) Update(x uint16, n int, rst bool) {
	n = int(beeMath.LIMIT(0, float64(n), 32))

	if !f.init || rst || n == 0 {
		f.init = true
		for j := 0; j < n; j++ {
			f.buffer[j] = x
		}
		f.Y = x
		f.sum = uint32(x) * uint32(n)
		f.i = 0
	} else {
		f.i = beeMath.INC1(f.i, n)
		f.sum = f.sum + uint32(x) - uint32(f.buffer[f.i])
		f.Y = uint16(f.sum / uint32(n))
		f.buffer[f.i] = x
	}
}

// FILTER_W is a low-pass filter for a WORD value.
type FILTER_W struct {
	Y    uint16
	last int64
	init bool
}

// Update executes the filter logic.
func (f *FILTER_W) Update(x uint16, t time.Duration) {
	tx := logic.T_PLC_US() / 1000 // T_PLC_MS

	if !f.init || t == 0 {
		f.init = true
		f.last = tx
		f.Y = x
	} else if f.Y == x {
		f.last = tx
	} else if t > 0 {
		tmp := (int64(x) - int64(f.Y)) * (tx - f.last) / t.Milliseconds()
		if tmp != 0 {
			f.Y = uint16(int64(f.Y) + tmp)
			f.last = tx
		}
	}
}

// FILTER_WAV is a weighted moving average filter.
type FILTER_WAV struct {
	Y      float64
	init   bool
	buffer [16]float64
	i      int
}

// Update executes the filter logic.
func (f *FILTER_WAV) Update(x float64, w [16]float64, rst bool) {
	if !f.init || rst {
		f.init = true
		for j := 0; j < 16; j++ {
			f.buffer[j] = x
		}
		f.i = 15
		f.Y = x
	} else {
		f.i = beeMath.INC1(f.i, 16)
		f.buffer[f.i] = x
	}

	f.Y = 0.0
	idx := f.i
	for n := 0; n < 16; n++ {
		f.Y += f.buffer[idx] * w[n]
		idx = (idx - 1 + 16) % 16 // DEC1
	}
}

// MIX is an analog mixer: Y = (1-M)*A + M*B.
func MIX(a, b, m float64) float64 {
	return (1.0-m)*a + m*b
}

// MUX_R2 is a 2-to-1 analog multiplexer.
func MUX_R2(in0, in1 float64, a bool) float64 {
	if a {
		return in1
	}
	return in0
}

// MUX_R4 is a 4-to-1 analog multiplexer.
func MUX_R4(in0, in1, in2, in3 float64, a0, a1 bool) float64 {
	if a1 {
		return MUX_R2(in2, in3, a0)
	}
	return MUX_R2(in0, in1, a0)
}

// OFFSET adds multiple offsets to an input signal.
func OFFSET(x, o1, o2, o3, o4, def float64, b1, b2, b3, b4, d bool) float64 {
	var res float64
	if d {
		res = def
	} else {
		res = x
	}
	if b1 {
		res += o1
	}
	if b2 {
		res += o2
	}
	if b3 {
		res += o3
	}
	if b4 {
		res += o4
	}
	return res
}

// OFFSET2 adds a prioritized offset to an input signal.
func OFFSET2(x, o1, o2, o3, o4, def float64, b1, b2, b3, b4, d bool) float64 {
	var res float64
	if d {
		res = def
	} else {
		res = x
	}

	if b4 {
		return res + o4
	} else if b3 {
		return res + o3
	} else if b2 {
		return res + o2
	} else if b1 {
		return res + o1
	}
	return res
}

// OVERRIDE selects an input based on enabled flags, prioritizing the one with the largest absolute value.
func OVERRIDE(x1, x2, x3 float64, e1, e2, e3 bool) float64 {
	var res float64
	if e1 {
		res = x1
	}
	if e2 && math.Abs(x2) > math.Abs(res) {
		res = x2
	}
	if e3 && math.Abs(x3) > math.Abs(res) {
		res = x3
	}
	return res
}

// RANGE_TO_BYTE converts a real value between low and high into a byte.
func RANGE_TO_BYTE(x, low, high float64) byte {
	if high == low {
		return 0
	}
	val := (beeMath.LIMIT(low, x, high) - low) * 255.0 / (high - low)
	return byte(math.Trunc(val))
}

// RANGE_TO_WORD converts a real value between low and high into a word.
func RANGE_TO_WORD(x, low, high float64) uint16 {
	if high == low {
		return 0
	}
	val := (beeMath.LIMIT(low, x, high) - low) * 65535.0 / (high - low)
	return uint16(math.Trunc(val))
}

// SCALE scales and limits an input signal. Y = (X*K + O) limited by MN and MX.
func SCALE(x, k, o, mx, mn float64) float64 {
	return beeMath.LIMIT(mn, x*k+o, mx)
}

// SCALE_B scales a byte input to a real output range.
func SCALE_B(x, iLo, iHi byte, oLo, oHi float64) float64 {
	if iHi == iLo {
		return oLo
	}
	val := beeMath.LIMIT_B(iLo, x, iHi)
	return (oHi-oLo)/float64(iHi-iLo)*float64(val) + oLo
}

// SCALE_B2 scales and sums two byte inputs.
func SCALE_B2(in1, in2 byte, k, o, in1Min, in1Max, in2Min, in2Max float64) float64 {
	val1 := (in1Max-in1Min)*float64(in1) + (in2Max-in2Min)*float64(in2)
	return (val1*0.003921569+in1Min+in2Min)*k + o
}

// SCALE_B4 scales and sums four byte inputs.
func SCALE_B4(in1, in2, in3, in4 byte, k, o, in1Min, in1Max, in2Min, in2Max, in3Min, in3Max, in4Min, in4Max float64) float64 {
	val1 := (in1Max-in1Min)*float64(in1) + (in2Max-in2Min)*float64(in2)
	val2 := (in3Max-in3Min)*float64(in3) + (in4Max-in4Min)*float64(in4)
	return ((val1+val2)*0.003921569+in1Min+in2Min+in3Min+in4Min)*k + o
}

// SCALE_B8 scales and sums eight byte inputs.
func SCALE_B8(in1, in2, in3, in4, in5, in6, in7, in8 byte, k, o float64, ranges [8][2]float64) float64 {
	var sum float64
	var minSum float64
	for i := 0; i < 8; i++ {
		sum += (ranges[i][1] - ranges[i][0]) * float64([]byte{in1, in2, in3, in4, in5, in6, in7, in8}[i])
		minSum += ranges[i][0]
	}
	return (sum*0.003921569+minSum)*k + o
}

// SCALE_D scales a DWORD input to a real output range.
func SCALE_D(x, iLo, iHi uint32, oLo, oHi float64) float64 {
	if iHi == iLo {
		return oLo
	}
	val := beeMath.LIMIT_DW(iLo, x, iHi)
	return (oHi-oLo)/float64(iHi-iLo)*float64(val-iLo) + oLo
}

// SCALE_R scales a REAL input to a real output range.
func SCALE_R(x, iLo, iHi, oLo, oHi float64) float64 {
	if iHi == iLo {
		return oLo
	}
	val := beeMath.LIMIT(iLo, x, iHi)
	return (oHi-oLo)/(iHi-iLo)*(val-iLo) + oLo
}

// SCALE_X2 scales and sums two boolean inputs.
func SCALE_X2(in1, in2 bool, k, o, in1Min, in1Max, in2Min, in2Max float64) float64 {
	var v1, v2 float64
	if in1 {
		v1 = in1Max
	} else {
		v1 = in1Min
	}
	if in2 {
		v2 = in2Max
	} else {
		v2 = in2Min
	}
	return (v1+v2)*k + o
}

// SCALE_X4 scales and sums four boolean inputs.
func SCALE_X4(in1, in2, in3, in4 bool, k, o float64, ranges [4][2]float64) float64 {
	var sum float64
	inputs := []bool{in1, in2, in3, in4}
	for i := 0; i < 4; i++ {
		if inputs[i] {
			sum += ranges[i][1]
		} else {
			sum += ranges[i][0]
		}
	}
	return sum*k + o
}

// SCALE_X8 scales and sums eight boolean inputs.
func SCALE_X8(in1, in2, in3, in4, in5, in6, in7, in8 bool, k, o float64, ranges [8][2]float64) float64 {
	var sum float64
	inputs := []bool{in1, in2, in3, in4, in5, in6, in7, in8}
	for i := 0; i < 8; i++ {
		if inputs[i] {
			sum += ranges[i][1]
		} else {
			sum += ranges[i][0]
		}
	}
	return sum*k + o
}

// SEL2_OF_3 selects the average of two out of three inputs that are closest to each other.
type SEL2_OF_3 struct {
	Y float64
	W int
	E bool
}

// Update executes the selection logic.
func (s *SEL2_OF_3) Update(in1, in2, in3, d float64) {
	d12 := math.Abs(in1-in2) <= d
	d23 := math.Abs(in2-in3) <= d
	d31 := math.Abs(in3-in1) <= d

	if (d12 && d23) || (d12 && d31) || (d23 && d31) {
		s.Y = beeMath.MID3(in1, in2, in3)
		s.E = false
		s.W = 0
	} else if d12 {
		s.Y = (in1 + in2) * 0.5
		s.E = false
		s.W = 3
	} else if d23 {
		s.Y = (in2 + in3) * 0.5
		s.E = false
		s.W = 1
	} else if d31 {
		s.Y = (in3 + in1) * 0.5
		s.E = false
		s.W = 2
	} else {
		s.E = true
		s.W = 4
	}
}

// SEL2_OF_3B is a 2-out-of-3 voter for boolean signals.
type SEL2_OF_3B struct {
	Q    bool
	W    bool
	tdel logic.TON
}

// Update executes the voter logic.
func (s *SEL2_OF_3B) Update(in1, in2, in3 bool, td time.Duration) {
	s.Q = (in1 && in2) || (in1 && in3) || (in2 && in3)
	s.tdel.Update((in1 != in2) || (in1 != in3), td)
	s.W = s.tdel.Q
}

// SH is a sample and hold block, triggered by a clock.
type SH struct {
	Out  float64
	Trig bool
	edge bool
}

// Update executes the logic.
func (s *SH) Update(in float64, clk bool) {
	if clk && !s.edge {
		s.Out = in
		s.Trig = true
	} else {
		s.Trig = false
	}
	s.edge = clk
}

// SH_1 is a sample and hold block, triggered by a timer.
type SH_1 struct {
	Out  float64
	Trig bool
	last time.Time
}

// Update executes the logic.
func (s *SH_1) Update(in float64, pt time.Duration) {
	tx := time.Now()
	if tx.Sub(s.last) >= pt {
		s.last = tx
		s.Out = in
		s.Trig = true
	} else {
		s.Trig = false
	}
}

// SH_2 is a sample and hold block with statistics.
type SH_2 struct {
	Out  float64
	Trig bool
	Avg  float64
	High float64
	Low  float64

	buf  [16]float64
	last time.Time
	m    int
}

// Update executes the logic.
func (s *SH_2) Update(in float64, pt time.Duration, n, disc int) {
	tx := time.Now()
	s.Trig = false

	if tx.Sub(s.last) >= pt {
		s.last = tx
		s.Trig = true
		s.m = int(beeMath.LIMIT(1, float64(n), 16))

		// Shift buffer
		copy(s.buf[1:], s.buf[:s.m-1])
		s.buf[0] = in
		s.Out = in

		// Sort a copy for statistics
		sortedBuf := make([]float64, s.m)
		copy(sortedBuf, s.buf[:s.m])
		sort.Float64s(sortedBuf)

		d2 := disc / 2
		start := d2
		stop := s.m - 1 - d2
		if disc%2 != 0 { // odd
			start++
		}

		if start > stop {
			s.Avg, s.Low, s.High = 0.0, 0.0, 0.0
			return
		}

		var sum float64
		for i := start; i <= stop; i++ {
			sum += sortedBuf[i]
		}
		s.Avg = sum / float64(stop-start+1)
		s.Low = sortedBuf[start]
		s.High = sortedBuf[stop]
	}
}

// SH_T is a sample and hold block, transparent when enabled.
type SH_T struct {
	Out float64
}

// Update executes the logic.
func (s *SH_T) Update(in float64, e bool) {
	if e {
		s.Out = in
	}
}

// STAIR converts an analog signal to a staircase-like output.
func STAIR(x, d float64) float64 {
	if d > 0.0 {
		return math.Trunc(x/d) * d
	}
	return x
}

// STAIR2 is a staircase function with hysteresis.
type STAIR2 struct {
	Y float64
}

// Update executes the logic.
func (s *STAIR2) Update(x, d float64) {
	if d > 0.0 {
		if x >= s.Y+d || x <= s.Y-d {
			s.Y = math.Floor(x/d) * d
		}
	} else {
		s.Y = x
	}
}

// TREND analyses the trend of a real input signal.
type TREND struct {
	Q     bool
	TU    bool
	TD    bool
	D     float64
	lastX float64
}

// Update executes the trend analysis logic.
func (t *TREND) Update(x float64) {
	t.TU = x > t.lastX
	t.TD = x < t.lastX
	t.Q = t.TU || t.TD
	t.D = x - t.lastX
	t.lastX = x
}

// WORD_TO_RANGE converts a word into a real value between low and high.
func WORD_TO_RANGE(x uint16, low, high float64) float64 {
	return (high-low)*float64(x)*0.0000152590218966964 + low
}
