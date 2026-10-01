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

package logic

import (
	"time"

	"github.com/apiarytech/beebread/basic/math"
	"github.com/apiarytech/royaljelly/iec"
)

// fifo is a first in first out memory of n DWORDs, with the inputs and
// outputs of FIFO_16 and FIFO_32.
type fifo struct {
	DIN   iec.DWORD
	E     iec.BOOL // default TRUE
	RD    iec.BOOL
	WD    iec.BOOL
	RST   iec.BOOL
	DOUT  iec.DWORD
	EMPTY iec.BOOL // default TRUE
	FULL  iec.BOOL

	pr, pw      iec.INT
	initialized bool
}

func (f *fifo) defaults() {
	f.initialized = true
	f.EMPTY = true
}

func (f *fifo) run(buf []iec.DWORD) {
	n := iec.INT(len(buf))
	if f.RST {
		f.pw = f.pr
		f.FULL = false
		f.EMPTY = true
		f.DOUT = 0
	} else if f.E {
		if !f.EMPTY && f.RD {
			f.DOUT = buf[f.pr]
			f.pr = math.INC1(f.pr, n)
			f.EMPTY = f.pr == f.pw
			f.FULL = false
		}
		if !f.FULL && f.WD {
			buf[f.pw] = f.DIN
			f.pw = math.INC1(f.pw, n)
			f.FULL = f.pw == f.pr
			f.EMPTY = false
		}
	}
}

// FIFO_16 is a first in first out memory of 16 DWORDs. While E is true, RD
// reads the oldest value to DOUT and WD writes DIN, once each scan.
type FIFO_16 struct {
	fifo
	buf [16]iec.DWORD
}

// INIT resets the block and sets E to its initial value.
func (f *FIFO_16) INIT() {
	*f = FIFO_16{}
	f.defaults()
	f.E = true
}

// Execute runs the block once.
func (f *FIFO_16) Execute(now time.Time) {
	if !f.initialized {
		f.defaults()
	}
	f.run(f.buf[:])
}

// FIFO_32 is a first in first out memory of 32 DWORDs. While E is true, RD
// reads the oldest value to DOUT and WD writes DIN, once each scan.
type FIFO_32 struct {
	fifo
	buf [32]iec.DWORD
}

// INIT resets the block and sets E to its initial value.
func (f *FIFO_32) INIT() {
	*f = FIFO_32{}
	f.defaults()
	f.E = true
}

// Execute runs the block once.
func (f *FIFO_32) Execute(now time.Time) {
	if !f.initialized {
		f.defaults()
	}
	f.run(f.buf[:])
}

// stack is a last in first out memory of n DWORDs, with the inputs and
// outputs of STACK_16 and STACK_32.
type stack struct {
	DIN   iec.DWORD
	E     iec.BOOL // default TRUE
	RD    iec.BOOL
	WD    iec.BOOL
	RST   iec.BOOL
	DOUT  iec.DWORD
	EMPTY iec.BOOL // default TRUE
	FULL  iec.BOOL

	pt          iec.INT
	initialized bool
}

func (s *stack) defaults() {
	s.initialized = true
	s.EMPTY = true
}

func (s *stack) run(buf []iec.DWORD) {
	n := iec.INT(len(buf)) - 1
	if s.RST {
		s.pt = 0
		s.EMPTY = true
		s.FULL = false
		s.DOUT = 0
	} else if s.E {
		if !s.EMPTY && s.RD {
			s.pt--
			s.DOUT = buf[s.pt]
			s.EMPTY = s.pt == 0
			s.FULL = false
		}
		if !s.FULL && s.WD {
			buf[s.pt] = s.DIN
			s.pt++
			s.FULL = s.pt > n
			s.EMPTY = false
		}
	}
}

// STACK_16 is a last in first out memory of 16 DWORDs. While E is true, RD
// reads the newest value to DOUT and WD writes DIN, once each scan.
type STACK_16 struct {
	stack
	buf [16]iec.DWORD
}

// INIT resets the block and sets E to its initial value.
func (s *STACK_16) INIT() {
	*s = STACK_16{}
	s.defaults()
	s.E = true
}

// Execute runs the block once.
func (s *STACK_16) Execute(now time.Time) {
	if !s.initialized {
		s.defaults()
	}
	s.run(s.buf[:])
}

// STACK_32 is a last in first out memory of 32 DWORDs. While E is true, RD
// reads the newest value to DOUT and WD writes DIN, once each scan.
type STACK_32 struct {
	stack
	buf [32]iec.DWORD
}

// INIT resets the block and sets E to its initial value.
func (s *STACK_32) INIT() {
	*s = STACK_32{}
	s.defaults()
	s.E = true
}

// Execute runs the block once.
func (s *STACK_32) Execute(now time.Time) {
	if !s.initialized {
		s.defaults()
	}
	s.run(s.buf[:])
}
