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

package logic

import (
	. "beebread/basic/math"
)

// Fifo16 is a 16-element DWORD FIFO memory.
type Fifo16 struct {
	Dout  uint32
	Empty bool
	Full  bool

	// internal state
	fifo [16]uint32
	pr   int // read pointer
	pw   int // write pointer
}

// Update executes the FIFO logic for one cycle.
func (f *Fifo16) Update(din uint32, e, rd, wd, rst bool) {
	if rst {
		f.pw = f.pr
		f.Full = false
		f.Empty = true
		f.Dout = 0
		return
	}

	if !e {
		return
	}

	// A read and a write can happen in the same cycle.
	if !f.Empty && rd {
		f.Dout = f.fifo[f.pr]
		f.pr = Inc1(f.pr, 16)
		f.Empty = (f.pr == f.pw)
		f.Full = false
	}

	if !f.Full && wd {
		f.fifo[f.pw] = din
		f.pw = Inc1(f.pw, 16)
		f.Full = (f.pw == f.pr)
		f.Empty = false
	}
}

// Fifo32 is a 32-element DWORD FIFO memory.
type Fifo32 struct {
	Dout  uint32
	Empty bool
	Full  bool

	// internal state
	fifo [32]uint32
	pr   int // read pointer
	pw   int // write pointer
}

// Update executes the FIFO logic for one cycle.
func (f *Fifo32) Update(din uint32, e, rd, wd, rst bool) {
	if rst {
		f.pw = f.pr
		f.Full = false
		f.Empty = true
		f.Dout = 0
		return
	}

	if !e {
		return
	}

	if !f.Empty && rd {
		f.Dout = f.fifo[f.pr]
		f.pr = Inc1(f.pr, 32)
		f.Empty = (f.pr == f.pw)
		f.Full = false
	}

	if !f.Full && wd {
		f.fifo[f.pw] = din
		f.pw = Inc1(f.pw, 32)
		f.Full = (f.pw == f.pr)
		f.Empty = false
	}
}

// Stack16 is a 16-element DWORD LIFO (Last-In, First-Out) stack memory.
type Stack16 struct {
	Dout  uint32
	Empty bool
	Full  bool

	// internal state
	stack [16]uint32
	pt    int // stack pointer
}

// Update executes the stack logic for one cycle.
func (s *Stack16) Update(din uint32, e, rd, wd, rst bool) {
	if rst {
		s.pt = 0
		s.Empty = true
		s.Full = false
		s.Dout = 0
		return
	}

	if !e {
		return
	}

	// A read and a write can happen in the same cycle.
	if !s.Empty && rd {
		s.pt--
		s.Dout = s.stack[s.pt]
		s.Empty = (s.pt == 0)
		s.Full = false
	}

	if !s.Full && wd {
		s.stack[s.pt] = din
		s.pt++
		s.Full = (s.pt >= 16)
		s.Empty = false
	}
}

// Stack32 is a 32-element DWORD LIFO (Last-In, First-Out) stack memory.
type Stack32 struct {
	Dout  uint32
	Empty bool
	Full  bool

	// internal state
	stack [32]uint32
	pt    int // stack pointer
}

// Update executes the stack logic for one cycle.
func (s *Stack32) Update(din uint32, e, rd, wd, rst bool) {
	if rst {
		s.pt = 0
		s.Empty = true
		s.Full = false
		s.Dout = 0
		return
	}

	if !e {
		return
	}

	// A read and a write can happen in the same cycle.
	if !s.Empty && rd {
		s.pt--
		s.Dout = s.stack[s.pt]
		s.Empty = (s.pt == 0)
		s.Full = false
	}

	if !s.Full && wd {
		s.stack[s.pt] = din
		s.pt++
		s.Full = (s.pt >= 32)
		s.Empty = false
	}
}
