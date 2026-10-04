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

// Package telnet is the telnet user interface of OSCAT NETWORK: a screen
// of 24 rows of 80 columns (US_TN_SCREEN) that the TN_SC_* blocks draw on,
// input elements, a menu bar and popup menus that read the keys of the
// client, and TN_SEND_ROWS, which sends the rows that changed as ANSI
// escape sequences. TN_FRAMEWORK runs it all on a TCP port.
//
// An attribute is a byte: the foreground color in bits 4 to 6, the
// background color in bits 0 to 2, bit 3 dark colors and bit 7 blinking.
// The line characters are those of code page 437.
package telnet

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/royaljelly/iec"
)

// The size of the screen.
const (
	rows    = 24
	columns = 80
	cells   = rows * columns
)

// TN_SC_SHADOW_ATTR returns the attribute IBY_ATTR in the dark colors on
// black, for a shadow.
func TN_SC_SHADOW_ATTR(IBY_ATTR iec.BYTE) iec.BYTE {
	return IBY_ATTR&0b1111_0000 | 0b0000_1000
}

// TN_SC_XY_ERROR is TRUE when the position X, Y is not on the screen.
func TN_SC_XY_ERROR(X, Y iec.INT) iec.BOOL {
	return Y < 0 || X < 0 || Y > rows-1 || X > columns-1
}

// TN_SC_XY2_ERROR is TRUE when the position X1, Y1 or X2, Y2 is not on the
// screen.
func TN_SC_XY2_ERROR(X1, Y1, X2, Y2 iec.INT) iec.BOOL {
	return TN_SC_XY_ERROR(X1, Y1) || TN_SC_XY_ERROR(X2, Y2)
}

// offset returns the offset of the position x, y in the screen.
func offset(x, y iec.INT) iec.INT { return y*columns + x }

// TN_SC_WRITE_ATTR sets the attribute of the position IIN_X, IIN_Y to
// IBY_ATTR.
type TN_SC_WRITE_ATTR struct {
	IIN_Y         iec.INT
	IIN_X         iec.INT
	IBY_ATTR      iec.BYTE
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_WRITE_ATTR) INIT() { *b = TN_SC_WRITE_ATTR{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_WRITE_ATTR) Execute(now time.Time) {
	if b.XUS_TN_SCREEN == nil || TN_SC_XY_ERROR(b.IIN_X, b.IIN_Y) {
		return
	}
	b.XUS_TN_SCREEN.BYA_COLOR[offset(b.IIN_X, b.IIN_Y)] = b.IBY_ATTR
}

// TN_SC_WRITE_CHAR sets the character of the position IIN_X, IIN_Y to
// IBY_CHAR.
type TN_SC_WRITE_CHAR struct {
	IIN_Y         iec.INT
	IIN_X         iec.INT
	IBY_CHAR      iec.BYTE
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_WRITE_CHAR) INIT() { *b = TN_SC_WRITE_CHAR{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_WRITE_CHAR) Execute(now time.Time) {
	if b.XUS_TN_SCREEN == nil || TN_SC_XY_ERROR(b.IIN_X, b.IIN_Y) {
		return
	}
	b.XUS_TN_SCREEN.BYA_CHAR[offset(b.IIN_X, b.IIN_Y)] = b.IBY_CHAR
}

// TN_SC_READ_ATTR reads the attribute of the position IIN_X, IIN_Y into
// OBY_ATTR.
type TN_SC_READ_ATTR struct {
	IIN_Y         iec.INT
	IIN_X         iec.INT
	OBY_ATTR      iec.BYTE
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_READ_ATTR) INIT() { *b = TN_SC_READ_ATTR{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_READ_ATTR) Execute(now time.Time) {
	if b.XUS_TN_SCREEN == nil || TN_SC_XY_ERROR(b.IIN_X, b.IIN_Y) {
		return
	}
	b.OBY_ATTR = b.XUS_TN_SCREEN.BYA_COLOR[offset(b.IIN_X, b.IIN_Y)]
}

// TN_SC_READ_CHAR reads the character of the position IIN_X, IIN_Y into
// OBY_CHAR.
type TN_SC_READ_CHAR struct {
	IIN_Y         iec.INT
	IIN_X         iec.INT
	OBY_CHAR      iec.BYTE
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_READ_CHAR) INIT() { *b = TN_SC_READ_CHAR{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_READ_CHAR) Execute(now time.Time) {
	if b.XUS_TN_SCREEN == nil || TN_SC_XY_ERROR(b.IIN_X, b.IIN_Y) {
		return
	}
	b.OBY_CHAR = b.XUS_TN_SCREEN.BYA_CHAR[offset(b.IIN_X, b.IIN_Y)]
}

// TN_SC_WRITE writes IST_STRING at the position IIN_X, IIN_Y, in the
// attribute IBY_ATTR or, if 0, in those there, if it ends on the screen,
// and sets the end of the string IN_EOS_OFFSET of the screen past it.
type TN_SC_WRITE struct {
	IIN_Y         iec.INT
	IIN_X         iec.INT
	IBY_ATTR      iec.BYTE
	IST_STRING    iec.STRING
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_WRITE) INIT() { *b = TN_SC_WRITE{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_WRITE) Execute(now time.Time) {
	sc := b.XUS_TN_SCREEN
	n := LEN(b.IST_STRING)
	// OSCAT checks the end only: an empty string at column 0 is not written.
	if sc == nil || TN_SC_XY_ERROR(b.IIN_X+n-1, b.IIN_Y) {
		return
	}
	sc.BYA_LINE_UPDATE[b.IIN_Y] = true
	o := offset(b.IIN_X, b.IIN_Y)
	for i := iec.INT(1); i <= n; i++ {
		if o >= 0 { // OSCAT writes before the screen from a negative column
			if b.IBY_ATTR != 0 {
				sc.BYA_COLOR[o] = b.IBY_ATTR
			}
			sc.BYA_CHAR[o] = str.CODE(b.IST_STRING, i)
		}
		o++
	}
	sc.IN_EOS_OFFSET = o
}

// TN_SC_WRITE_EOS writes IST_STRING, in the attribute IBY_ATTR, at the end
// of the string written last.
type TN_SC_WRITE_EOS struct {
	IBY_ATTR      iec.BYTE
	IST_STRING    iec.STRING
	XUS_TN_SCREEN *network.US_TN_SCREEN

	write TN_SC_WRITE
}

// INIT resets the block.
func (b *TN_SC_WRITE_EOS) INIT() { *b = TN_SC_WRITE_EOS{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_WRITE_EOS) Execute(now time.Time) {
	if b.XUS_TN_SCREEN == nil {
		return
	}
	w := &b.write
	w.IIN_Y, w.IIN_X = b.XUS_TN_SCREEN.IN_EOS_OFFSET/columns, b.XUS_TN_SCREEN.IN_EOS_OFFSET%columns
	w.IBY_ATTR, w.IST_STRING, w.XUS_TN_SCREEN = b.IBY_ATTR, b.IST_STRING, b.XUS_TN_SCREEN
	w.Execute(now)
}

// TN_SC_WRITE_C writes IST_STRING, cut or filled with spaces to IIN_LENGTH
// by FIX with the mode IIN_OPTION (2 centers it), at IIN_X, IIN_Y.
type TN_SC_WRITE_C struct {
	IIN_Y         iec.INT
	IIN_X         iec.INT
	IBY_ATTR      iec.BYTE
	IST_STRING    iec.STRING
	IIN_LENGTH    iec.INT
	IIN_OPTION    iec.INT
	XUS_TN_SCREEN *network.US_TN_SCREEN

	write TN_SC_WRITE
}

// INIT resets the block.
func (b *TN_SC_WRITE_C) INIT() { *b = TN_SC_WRITE_C{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_WRITE_C) Execute(now time.Time) {
	w := &b.write
	w.IIN_Y, w.IIN_X, w.IBY_ATTR, w.XUS_TN_SCREEN = b.IIN_Y, b.IIN_X, b.IBY_ATTR, b.XUS_TN_SCREEN
	w.IST_STRING = str.FIX(b.IST_STRING, b.IIN_LENGTH, ' ', b.IIN_OPTION)
	w.Execute(now)
}

// TN_SC_FILL fills the rectangle IIN_X1, IIN_Y1 to IIN_X2, IIN_Y2 with the
// character IBY_CHAR in the attribute IBY_ATTR.
type TN_SC_FILL struct {
	IIN_Y1        iec.INT
	IIN_X1        iec.INT
	IIN_Y2        iec.INT
	IIN_X2        iec.INT
	IBY_CHAR      iec.BYTE
	IBY_ATTR      iec.BYTE
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_FILL) INIT() { *b = TN_SC_FILL{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_FILL) Execute(now time.Time) {
	sc := b.XUS_TN_SCREEN
	if sc == nil || TN_SC_XY2_ERROR(b.IIN_X1, b.IIN_Y1, b.IIN_X2, b.IIN_Y2) {
		return
	}
	for y := b.IIN_Y1; y <= b.IIN_Y2; y++ {
		sc.BYA_LINE_UPDATE[y] = true
		for x := b.IIN_X1; x <= b.IIN_X2; x++ {
			sc.BYA_CHAR[offset(x, y)] = b.IBY_CHAR
			sc.BYA_COLOR[offset(x, y)] = b.IBY_ATTR
		}
	}
}

// TN_SC_AREA_SAVE saves the rectangle IIN_X1, IIN_Y1 to IIN_X2, IIN_Y2 of
// the screen, its characters and attributes, in its BYA_BACKUP, for
// TN_SC_AREA_RESTORE.
type TN_SC_AREA_SAVE struct {
	IIN_Y1        iec.INT
	IIN_X1        iec.INT
	IIN_Y2        iec.INT
	IIN_X2        iec.INT
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_AREA_SAVE) INIT() { *b = TN_SC_AREA_SAVE{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_AREA_SAVE) Execute(now time.Time) {
	sc := b.XUS_TN_SCREEN
	if sc == nil || TN_SC_XY2_ERROR(b.IIN_X1, b.IIN_Y1, b.IIN_X2, b.IIN_Y2) {
		return
	}
	bk := &sc.BYA_BACKUP
	bk[0] = 1
	bk[1], bk[2], bk[3], bk[4] = iec.BYTE(b.IIN_X1), iec.BYTE(b.IIN_Y1), iec.BYTE(b.IIN_X2), iec.BYTE(b.IIN_Y2)
	o := 5
	for y := b.IIN_Y1; y <= b.IIN_Y2; y++ {
		for x := b.IIN_X1; x <= b.IIN_X2 && o+1 < cells; x++ {
			// OSCAT writes past the backup for an area of over 957 cells.
			bk[o], bk[o+1] = sc.BYA_CHAR[offset(x, y)], sc.BYA_COLOR[offset(x, y)]
			o += 2
		}
	}
}

// TN_SC_AREA_RESTORE restores the area TN_SC_AREA_SAVE saved, once.
type TN_SC_AREA_RESTORE struct {
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_AREA_RESTORE) INIT() {}

// Execute runs the block once.
func (b *TN_SC_AREA_RESTORE) Execute(now time.Time) {
	sc := b.XUS_TN_SCREEN
	if sc == nil || sc.BYA_BACKUP[0] == 0 {
		return
	}
	bk := &sc.BYA_BACKUP
	x1, y1, x2, y2 := iec.INT(bk[1]), iec.INT(bk[2]), iec.INT(bk[3]), iec.INT(bk[4])
	if TN_SC_XY2_ERROR(x1, y1, x2, y2) {
		return
	}
	o := 5
	for y := y1; y <= y2; y++ {
		sc.BYA_LINE_UPDATE[y] = true
		for x := x1; x <= x2 && o+1 < cells; x++ {
			sc.BYA_CHAR[offset(x, y)], sc.BYA_COLOR[offset(x, y)] = bk[o], bk[o+1]
			o += 2
		}
	}
	bk[0] = 0 // used
}

// TN_SC_ADD_SHADOW draws the shadow of the rectangle IIN_X1, IIN_Y1 to
// IIN_X2, IIN_Y2: its right column and its bottom row, in the dark colors
// on black if IIN_OPTION is 0, else as black spaces.
type TN_SC_ADD_SHADOW struct {
	IIN_Y1        iec.INT
	IIN_X1        iec.INT
	IIN_Y2        iec.INT
	IIN_X2        iec.INT
	IIN_OPTION    iec.INT
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_ADD_SHADOW) INIT() { *b = TN_SC_ADD_SHADOW{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_ADD_SHADOW) Execute(now time.Time) {
	sc := b.XUS_TN_SCREEN
	if sc == nil || TN_SC_XY2_ERROR(b.IIN_X1, b.IIN_Y1, b.IIN_X2, b.IIN_Y2) {
		return
	}
	shade := func(o iec.INT) {
		if b.IIN_OPTION == 0 {
			sc.BYA_COLOR[o] = TN_SC_SHADOW_ATTR(sc.BYA_COLOR[o])
		} else {
			sc.BYA_COLOR[o], sc.BYA_CHAR[o] = 0, ' '
		}
	}
	for y := b.IIN_Y1; y <= b.IIN_Y2; y++ {
		shade(offset(b.IIN_X2, y))
		sc.BYA_LINE_UPDATE[y] = true
	}
	for x := b.IIN_X1; x <= b.IIN_X2; x++ {
		shade(offset(x, b.IIN_Y2))
	}
}

// lineChars are the characters of a line: the joins left, right, top and
// bottom, the cross, and the vertical and horizontal lines.
type lineChars struct{ left, right, top, bottom, cross, vertical, horizontal iec.BYTE }

// TN_SC_LINE draws a vertical line from IIN_X1, IIN_Y1 to IIN_Y2, or else
// a horizontal one to IIN_X2, in the attribute IBY_ATTR: IBY_BORDER 1 a
// single line, 2 a double one, else of the character IBY_BORDER. Where it
// crosses a line of the other direction it draws the join.
type TN_SC_LINE struct {
	IIN_X1        iec.INT
	IIN_Y1        iec.INT
	IIN_X2        iec.INT
	IIN_Y2        iec.INT
	IBY_ATTR      iec.BYTE
	IBY_BORDER    iec.BYTE
	XUS_TN_SCREEN *network.US_TN_SCREEN
}

// INIT resets the block.
func (b *TN_SC_LINE) INIT() { *b = TN_SC_LINE{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_LINE) Execute(now time.Time) {
	sc := b.XUS_TN_SCREEN
	if sc == nil || TN_SC_XY2_ERROR(b.IIN_X1, b.IIN_Y1, b.IIN_X2, b.IIN_Y2) {
		return
	}
	var c lineChars
	switch b.IBY_BORDER {
	case 1:
		c = lineChars{0xC3, 0xB4, 0xC2, 0xC1, 0xC5, 0xB3, 0xC4}
	case 2:
		c = lineChars{0xCC, 0xB9, 0xCB, 0xCA, 0xCE, 0xBA, 0xCD}
	default:
		k := b.IBY_BORDER
		c = lineChars{k, k, k, k, k, k, k}
	}
	put := func(o iec.INT, ch iec.BYTE) { sc.BYA_CHAR[o], sc.BYA_COLOR[o] = ch, b.IBY_ATTR }
	switch {
	case b.IIN_X1 == b.IIN_X2:
		for y := b.IIN_Y1; y <= b.IIN_Y2; y++ {
			o, ch := offset(b.IIN_X1, y), c.vertical
			if sc.BYA_CHAR[o] == c.horizontal {
				switch y {
				case b.IIN_Y1:
					ch = c.top
				case b.IIN_Y2:
					ch = c.bottom
				default:
					ch = c.cross
				}
			}
			sc.BYA_LINE_UPDATE[y] = true
			put(o, ch)
		}
	case b.IIN_Y1 == b.IIN_Y2:
		sc.BYA_LINE_UPDATE[b.IIN_Y1] = true
		for x := b.IIN_X1; x <= b.IIN_X2; x++ {
			o, ch := offset(x, b.IIN_Y1), c.horizontal
			if sc.BYA_CHAR[o] == c.vertical {
				switch x {
				case b.IIN_X1:
					ch = c.left
				case b.IIN_X2:
					ch = c.right
				default:
					ch = c.cross
				}
			}
			put(o, ch)
		}
	}
}

// TN_SC_BOX draws the box IIN_X1, IIN_Y1 to IIN_X2, IIN_Y2 in the attribute
// IBY_ATTR, filled with IBY_FILL if not 0, with the border IIN_BORDER: 0
// none, 1 single, 2 double, else of spaces.
type TN_SC_BOX struct {
	IIN_Y1        iec.INT
	IIN_X1        iec.INT
	IIN_Y2        iec.INT
	IIN_X2        iec.INT
	IBY_FILL      iec.BYTE
	IBY_ATTR      iec.BYTE
	IIN_BORDER    iec.INT
	XUS_TN_SCREEN *network.US_TN_SCREEN

	fill TN_SC_FILL
}

// INIT resets the block.
func (b *TN_SC_BOX) INIT() { *b = TN_SC_BOX{XUS_TN_SCREEN: b.XUS_TN_SCREEN} }

// Execute runs the block once.
func (b *TN_SC_BOX) Execute(now time.Time) {
	sc := b.XUS_TN_SCREEN
	if sc == nil || TN_SC_XY2_ERROR(b.IIN_X1, b.IIN_Y1, b.IIN_X2, b.IIN_Y2) {
		return
	}
	// The corners: top left, top right, bottom left, bottom right; the
	// sides: top, bottom, left, right.
	var corner [4]iec.BYTE
	var side [4]iec.BYTE
	inset := iec.INT(1)
	switch b.IIN_BORDER {
	case 0:
		inset = 0
	case 1:
		corner, side = [4]iec.BYTE{218, 191, 192, 217}, [4]iec.BYTE{196, 196, 179, 179}
	case 2:
		corner, side = [4]iec.BYTE{201, 187, 200, 188}, [4]iec.BYTE{205, 205, 186, 186}
	default:
		corner, side = [4]iec.BYTE{' ', ' ', ' ', ' '}, [4]iec.BYTE{' ', ' ', ' ', ' '}
	}
	if b.IBY_FILL > 0 { // the inside
		f := &b.fill
		f.IIN_Y1, f.IIN_X1, f.IIN_Y2, f.IIN_X2 = b.IIN_Y1+inset, b.IIN_X1+inset, b.IIN_Y2-inset, b.IIN_X2-inset
		f.IBY_CHAR, f.IBY_ATTR, f.XUS_TN_SCREEN = b.IBY_FILL, b.IBY_ATTR, sc
		f.Execute(now)
	}
	if b.IIN_BORDER == 0 {
		return
	}
	put := func(x, y iec.INT, ch iec.BYTE) {
		sc.BYA_CHAR[offset(x, y)], sc.BYA_COLOR[offset(x, y)] = ch, b.IBY_ATTR
	}
	for y := b.IIN_Y1; y <= b.IIN_Y2; y++ {
		sc.BYA_LINE_UPDATE[y] = true
		put(b.IIN_X1, y, side[2])
		put(b.IIN_X2, y, side[3])
	}
	for x := b.IIN_X1; x <= b.IIN_X2; x++ {
		put(x, b.IIN_Y1, side[0])
		put(x, b.IIN_Y2, side[1])
	}
	put(b.IIN_X1, b.IIN_Y1, corner[0])
	put(b.IIN_X2, b.IIN_Y1, corner[1])
	put(b.IIN_X1, b.IIN_Y2, corner[2])
	put(b.IIN_X2, b.IIN_Y2, corner[3])
}
