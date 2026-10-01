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

// Package basic and the packages below it are a Go port of the OSCAT BASIC
// library, version 3.35 (documents/oscat_basic_335.st). The port follows the
// conventions of royaljelly, the runtime beedance's transpiler targets, so
// the library can be called from transpiled IEC 61131-3 programs.
//
// # Names
//
// Every POU keeps its OSCAT name. Go does not export a name that starts with
// an underscore, so a name such as _BUFFER_CLEAR moves the underscore to the
// end: BUFFER_CLEAR_. Inputs, outputs and structure members keep their OSCAT
// names in upper case.
//
// # Types
//
// IEC types are royaljelly's: BOOL is iec.BOOL, REAL is iec.REAL (32 bit),
// TIME is iec.TIME and so on. A POINTER to a buffer or array, which OSCAT
// passes with its size, is a Go slice followed by the size.
//
// # Functions
//
// A FUNCTION is a Go function that takes its VAR_INPUTs in their declared
// order and returns the function's value. A VAR_IN_OUT is a pointer.
//
// # Function blocks
//
// A FUNCTION_BLOCK is a struct. Its VAR_INPUTs and VAR_OUTPUTs are exported
// fields, a VAR_IN_OUT is a pointer field, and its internal state is not
// exported. Every function block has two methods:
//
//	INIT()                 // resets the block and sets the inputs' initial values
//	Execute(now time.Time) // runs the block once, at the scan time now
//
// INIT sets every variable of the block to its initial value, as a new
// instance has them; only the VAR_IN_OUT pointers stay bound.
//
// A block never reads the wall clock; time is the now it is given. INIT sets
// the initial values OSCAT declares for the inputs, so a block whose inputs
// have initial values, such as FT_AVG's N := 32, needs INIT before it first
// runs. The initial values of outputs and internal state need no INIT: a
// block sets them the first time it runs.
//
// # Strings
//
// OSCAT strings hold single-byte characters (ISO 8859-1). An iec.STRING
// holds the same characters as UTF-8: positions and lengths count
// characters, and a character's code is its Unicode code point, which for
// ISO 8859-1 is the same number.
package basic
