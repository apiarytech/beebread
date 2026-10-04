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

package telnet

import (
	"time"

	. "github.com/apiarytech/beebread/basic"
	"github.com/apiarytech/beebread/basic/math"
	str "github.com/apiarytech/beebread/basic/string"
	"github.com/apiarytech/beebread/network"
	"github.com/apiarytech/beebread/network/encoding"
	"github.com/apiarytech/royaljelly/iec"
)

// The keys, as TN_RECEIVE reads them.
const (
	keyBackspace iec.BYTE = 8
	keyTab       iec.BYTE = 9
	keyReturn    iec.BYTE = 13
	keyEscape    iec.BYTE = 27
	keyUp        iec.BYTE = 65
	keyDown      iec.BYTE = 66
	keyRight     iec.BYTE = 67
	keyLeft      iec.BYTE = 68
)

// The separators of the menu texts: '#' between the items, '%' between the
// menus.
const (
	sepItem iec.BYTE = '#'
	sepMenu iec.BYTE = '%'
)

// title writes the title of the element d with w: left of it if its
// offsets are 0, else at them.
func title(now time.Time, w *TN_SC_WRITE, d *network.US_TN_INPUT_CONTROL_DATA, sc *network.US_TN_SCREEN) {
	if d.IN_TITLE_X_OFFSET == 0 && d.IN_TITLE_Y_OFFSET == 0 {
		w.IIN_Y, w.IIN_X = d.IN_Y, d.IN_X-LEN(d.ST_TITLE_STRING)
	} else {
		w.IIN_Y, w.IIN_X = d.IN_Y+d.IN_TITLE_Y_OFFSET, d.IN_X+d.IN_TITLE_X_OFFSET
	}
	w.IBY_ATTR, w.IST_STRING, w.XUS_TN_SCREEN = d.BY_TITLE_ATTR, d.ST_TITLE_STRING, sc
	w.Execute(now)
}

// attr returns the attribute of the element d: with or without the focus.
func attr(d *network.US_TN_INPUT_CONTROL_DATA) iec.BYTE {
	if d.BO_FOCUS {
		return d.BY_ATTR_MF
	}
	return d.BY_ATTR_OF
}

// showSelected writes the item IN_SELECTED of the items of the element d,
// fitted to its mask by FIX with its IN_INPUT_OPTION, into ST_INPUT_STRING
// and at the element.
func showSelected(now time.Time, w *TN_SC_WRITE, d *network.US_TN_INPUT_CONTROL_DATA, sc *network.US_TN_SCREEN) {
	s := encoding.ELEMENT_GET(sepItem, d.IN_SELECTED-1, &d.ST_INPUT_DATA)
	s = str.FIX(s, LEN(d.ST_INPUT_MASK), ' ', d.IN_INPUT_OPTION)
	d.ST_INPUT_STRING = s
	w.IIN_Y, w.IIN_X, w.IBY_ATTR, w.IST_STRING, w.XUS_TN_SCREEN = d.IN_Y, d.IN_X, attr(d), s, sc
	w.Execute(now)
}

// TN_INPUT_CONTROL runs the input elements XUS_TN_INPUT_CONTROL on the
// screen XUS_TN_SCREEN while BO_ENABLE and the menu bar is closed: tab and
// the cursor keys up and down move the focus, the element with the focus
// gets the keys, and its tooltip is shown. The elements are of the types
// IN_TYPE 1 TN_INPUT_EDIT_LINE, 2 TN_INPUT_SELECT_TEXT and 3
// TN_INPUT_SELECT_POPUP; the popup elements share one popup.
type TN_INPUT_CONTROL struct {
	XUS_TN_SCREEN        *network.US_TN_SCREEN
	XUS_TN_INPUT_CONTROL *network.US_TN_INPUT_CONTROL

	selectPopup   TN_INPUT_SELECT_POPUP
	selectText    TN_INPUT_SELECT_TEXT
	writeC        TN_SC_WRITE_C
	editLine      TN_INPUT_EDIT_LINE
	scrollOffset  iec.INT
	toolTipUpdate iec.BOOL
}

// INIT resets the block.
func (b *TN_INPUT_CONTROL) INIT() {
	*b = TN_INPUT_CONTROL{XUS_TN_SCREEN: b.XUS_TN_SCREEN, XUS_TN_INPUT_CONTROL: b.XUS_TN_INPUT_CONTROL}
}

// Execute runs the block once.
func (b *TN_INPUT_CONTROL) Execute(now time.Time) {
	sc, ic := b.XUS_TN_SCREEN, b.XUS_TN_INPUT_CONTROL
	if sc == nil || ic == nil || !ic.BO_ENABLE || sc.BO_MENUE_BAR_DIALOG {
		return
	}
	data := func(i iec.INT) *network.US_TN_INPUT_CONTROL_DATA { return &ic.USA_TN_INPUT_CONTROL_DATA[i-1] }

	if ic.BO_RESET_FOKUS { // the focus to the first element
		ic.BO_RESET_FOKUS = false
		for i := iec.INT(1); i <= ic.IN_COUNT; i++ {
			data(i).BO_FOCUS = false
		}
		if ic.IN_COUNT >= 1 {
			data(1).BO_FOCUS = true
			ic.IN_FOCUS_AT = 1
		}
		ic.BO_UPDATE_ALL = true
	}

	if !sc.BO_MODAL_DIALOG {
		switch sc.BY_INPUT_EXTEN_CODE {
		case keyTab, keyDown:
			b.scrollOffset = 1
		case keyUp:
			b.scrollOffset = -1
		}
	}
	if b.scrollOffset != 0 && ic.IN_COUNT < 1 {
		b.scrollOffset = 0 // OSCAT moves the focus of element 0, outside its array
	}
	if b.scrollOffset != 0 { // move the focus
		ic.IN_FOCUS_AT = LIMIT(1, ic.IN_FOCUS_AT, ic.IN_COUNT)
		next := ic.IN_FOCUS_AT + b.scrollOffset
		if next > ic.IN_COUNT {
			next = 1
		}
		if next < 1 {
			next = ic.IN_COUNT
		}
		old := data(ic.IN_FOCUS_AT)
		old.BO_FOCUS, old.BO_UPDATE_INPUT = false, true
		d := data(next)
		d.BO_FOCUS, d.BO_UPDATE_INPUT = true, true
		ic.IN_FOCUS_AT = next
		b.toolTipUpdate = true
		b.scrollOffset = 0
	}

	if ic.BO_UPDATE_ALL {
		for i := iec.INT(1); i <= ic.IN_COUNT; i++ {
			data(i).BO_UPDATE_ALL = true
		}
		ic.BO_UPDATE_ALL = false
		b.toolTipUpdate = true
	}

	for i := iec.INT(1); i <= ic.IN_COUNT; i++ {
		d := data(i)
		if d.IN_TYPE <= 0 {
			continue
		}
		if d.BO_FOCUS {
			d.BY_INPUT_ASCII_CODE, d.BY_INPUT_EXTEN_CODE, d.BO_INPUT_ASCII_ISNUM = sc.BY_INPUT_ASCII_CODE, sc.BY_INPUT_EXTEN_CODE, sc.BO_INPUT_ASCII_ISNUM
			if b.toolTipUpdate {
				w := &b.writeC
				w.IIN_Y, w.IIN_X, w.IBY_ATTR = ic.IN_TOOLTIP_Y, ic.IN_TOOLTIP_X, ic.BY_TOOLTIP_ATTR
				w.IIN_OPTION, w.IIN_LENGTH, w.IST_STRING, w.XUS_TN_SCREEN = 2, ic.IN_TOOLTIP_SIZE, d.ST_INPUT_TOOLTIP, sc
				w.Execute(now)
				b.toolTipUpdate = false
			}
		}
		switch d.IN_TYPE {
		case 1:
			b.editLine.XUS_TN_INPUT_CONTROL_DATA, b.editLine.XUS_TN_SCREEN = d, sc
			b.editLine.Execute(now)
		case 2:
			b.selectText.XUS_TN_INPUT_CONTROL_DATA, b.selectText.XUS_TN_SCREEN = d, sc
			b.selectText.Execute(now)
		case 3:
			b.selectPopup.XUS_TN_INPUT_CONTROL_DATA, b.selectPopup.XUS_TN_SCREEN = d, sc
			b.selectPopup.Execute(now)
		}
		if d.BO_FOCUS {
			sc.IN_CURSOR_X, sc.IN_CURSOR_Y = d.IN_CURSOR_X, d.IN_CURSOR_Y
		}
	}
}

// TN_INPUT_EDIT_LINE is an input line: its mask ST_INPUT_MASK has a space
// where a character may be entered, and the line ST_INPUT_DATA starts as
// the mask. The cursor keys move over the spaces of the mask, backspace
// clears, and return puts the line, trimmed, in ST_INPUT_STRING, sets
// BO_INPUT_ENTERED and starts again. BO_INPUT_ONLY_NUM takes digits only,
// and BO_INPUT_HIDDEN shows '*' for the characters entered.
type TN_INPUT_EDIT_LINE struct {
	XUS_TN_SCREEN             *network.US_TN_SCREEN
	XUS_TN_INPUT_CONTROL_DATA *network.US_TN_INPUT_CONTROL_DATA

	write TN_SC_WRITE
}

// INIT resets the block.
func (b *TN_INPUT_EDIT_LINE) INIT() {
	*b = TN_INPUT_EDIT_LINE{XUS_TN_SCREEN: b.XUS_TN_SCREEN, XUS_TN_INPUT_CONTROL_DATA: b.XUS_TN_INPUT_CONTROL_DATA}
}

// Execute runs the block once.
func (b *TN_INPUT_EDIT_LINE) Execute(now time.Time) {
	sc, d := b.XUS_TN_SCREEN, b.XUS_TN_INPUT_CONTROL_DATA
	if sc == nil || d == nil {
		return
	}
	isSpace := func(s iec.STRING, i iec.INT) bool { return MID(s, 1, i) == " " }
	move := iec.INT(0)
	if d.BO_FOCUS {
		if d.BY_INPUT_ASCII_CODE > 0 && (!d.BO_INPUT_ONLY_NUM || d.BO_INPUT_ASCII_ISNUM) {
			d.ST_INPUT_DATA = REPLACE(d.ST_INPUT_DATA, str.CHR_TO_STRING(d.BY_INPUT_ASCII_CODE), 1, d.IN_CURSOR_POS)
			move = 1
			d.BO_UPDATE_INPUT = true
		}
		switch d.BY_INPUT_EXTEN_CODE {
		case keyRight:
			move = 1
		case keyLeft:
			move = -1
		case keyBackspace:
			d.ST_INPUT_DATA = REPLACE(d.ST_INPUT_DATA, " ", 1, d.IN_CURSOR_POS)
			d.BO_UPDATE_INPUT = true
			move = -1
		case keyReturn:
			if LEN(d.ST_INPUT_MASK) > 0 {
				d.ST_INPUT_STRING = str.TRIME(d.ST_INPUT_DATA)
				d.BO_INPUT_ENTERED = true
				d.ST_INPUT_DATA = d.ST_INPUT_MASK // the next entry
				d.IN_CURSOR_POS = 1
				d.BO_UPDATE_INPUT = true
				if !isSpace(d.ST_INPUT_MASK, d.IN_CURSOR_POS) {
					move = 1 // to the first space of the mask
				}
			}
		}
	}

	if move > 0 || d.IN_CURSOR_POS == 0 { // to the next space of the mask
		for i := d.IN_CURSOR_POS + 1; i <= LEN(d.ST_INPUT_MASK); i++ {
			if isSpace(d.ST_INPUT_MASK, i) {
				d.IN_CURSOR_POS = i
				d.BO_UPDATE_INPUT = true
				break
			}
		}
	}
	if move < 0 { // to the space before
		for i := d.IN_CURSOR_POS - 1; i >= 1; i-- {
			if isSpace(d.ST_INPUT_MASK, i) {
				d.IN_CURSOR_POS = i
				d.BO_UPDATE_INPUT = true
				break
			}
		}
	}
	if d.IN_CURSOR_POS > LEN(d.ST_INPUT_DATA) {
		d.IN_CURSOR_POS = LEN(d.ST_INPUT_DATA)
		d.BO_UPDATE_INPUT = true
	}
	if d.IN_CURSOR_POS < 1 {
		d.IN_CURSOR_POS = 1
		d.BO_UPDATE_INPUT = true
	}

	if !d.BO_UPDATE_INPUT && !d.BO_UPDATE_ALL {
		return
	}
	if d.BO_UPDATE_ALL {
		title(now, &b.write, d, sc)
	}
	s := d.ST_INPUT_DATA
	if d.BO_INPUT_HIDDEN {
		s = ""
		if LEN(d.ST_INPUT_DATA) == LEN(d.ST_INPUT_MASK) {
			for i := iec.INT(1); i <= LEN(d.ST_INPUT_DATA); i++ {
				c := MID(d.ST_INPUT_DATA, 1, i)
				if isSpace(d.ST_INPUT_MASK, i) && c != " " {
					c = "*" // an entered character
				}
				s = CONCAT(s, c)
			}
		}
	}
	w := &b.write
	w.IIN_Y, w.IIN_X, w.IBY_ATTR, w.IST_STRING, w.XUS_TN_SCREEN = d.IN_Y, d.IN_X, attr(d), s, sc
	w.Execute(now)
	d.IN_CURSOR_Y = d.IN_Y
	d.IN_CURSOR_X = d.IN_X + d.IN_CURSOR_POS - 1
	d.BO_UPDATE_ALL, d.BO_UPDATE_INPUT = false, false
}

// TN_INPUT_SELECT_TEXT selects one of the items, separated by '#', of
// ST_INPUT_DATA: space shows the next, and return sets BO_INPUT_ENTERED.
// IN_SELECTED is the item, from 1, and ST_INPUT_STRING its text.
type TN_INPUT_SELECT_TEXT struct {
	XUS_TN_SCREEN             *network.US_TN_SCREEN
	XUS_TN_INPUT_CONTROL_DATA *network.US_TN_INPUT_CONTROL_DATA

	write TN_SC_WRITE
}

// INIT resets the block.
func (b *TN_INPUT_SELECT_TEXT) INIT() {
	*b = TN_INPUT_SELECT_TEXT{XUS_TN_SCREEN: b.XUS_TN_SCREEN, XUS_TN_INPUT_CONTROL_DATA: b.XUS_TN_INPUT_CONTROL_DATA}
}

// Execute runs the block once.
func (b *TN_INPUT_SELECT_TEXT) Execute(now time.Time) {
	sc, d := b.XUS_TN_SCREEN, b.XUS_TN_INPUT_CONTROL_DATA
	if sc == nil || d == nil {
		return
	}
	if d.BO_FOCUS {
		if d.BY_INPUT_ASCII_CODE == ' ' {
			d.IN_SELECTED++
			d.BO_UPDATE_INPUT = true
		}
		if d.BY_INPUT_EXTEN_CODE == keyReturn {
			d.BO_INPUT_ENTERED = true
		}
	}
	if !d.BO_UPDATE_INPUT && !d.BO_UPDATE_ALL {
		return
	}
	if d.BO_UPDATE_ALL {
		title(now, &b.write, d, sc)
	}
	if n := encoding.ELEMENT_COUNT(sepItem, &d.ST_INPUT_DATA); n > 0 {
		if d.IN_SELECTED < 1 || d.IN_SELECTED > n {
			d.IN_SELECTED = 1
		}
		showSelected(now, &b.write, d, sc)
	}
	d.IN_CURSOR_Y, d.IN_CURSOR_X = d.IN_Y, d.IN_X
	d.BO_UPDATE_ALL, d.BO_UPDATE_INPUT = false, false
}

// TN_INPUT_SELECT_POPUP selects one of the items, separated by '#', of
// ST_INPUT_DATA from a popup menu that return opens; see
// TN_INPUT_SELECT_TEXT.
type TN_INPUT_SELECT_POPUP struct {
	XUS_TN_SCREEN             *network.US_TN_SCREEN
	XUS_TN_INPUT_CONTROL_DATA *network.US_TN_INPUT_CONTROL_DATA

	popupData network.US_TN_MENU_POPUP
	popup     TN_INPUT_MENU_POPUP
	write     TN_SC_WRITE
}

// INIT resets the block.
func (b *TN_INPUT_SELECT_POPUP) INIT() {
	*b = TN_INPUT_SELECT_POPUP{XUS_TN_SCREEN: b.XUS_TN_SCREEN, XUS_TN_INPUT_CONTROL_DATA: b.XUS_TN_INPUT_CONTROL_DATA}
}

// Execute runs the block once.
func (b *TN_INPUT_SELECT_POPUP) Execute(now time.Time) {
	sc, d := b.XUS_TN_SCREEN, b.XUS_TN_INPUT_CONTROL_DATA
	if sc == nil || d == nil {
		return
	}
	p := &b.popupData
	if d.BO_FOCUS {
		key := d.BY_INPUT_EXTEN_CODE
		if !p.BO_ACTIV {
			p.BY_INPUT_EXTEN_CODE = 0
			if key == keyReturn { // open the popup
				p.IN_Y, p.IN_X = d.IN_Y, d.IN_X
				p.BY_ATTR_MF, p.BY_ATTR_OF = d.BY_ATTR_MF, d.BY_ATTR_OF
				p.ST_MENU_TEXT = d.ST_INPUT_DATA
				p.BO_CREATE = true
				p.IN_CUR_ITEM = 0
			}
			if p.IN_CUR_ITEM > 0 { // an item selected
				d.IN_SELECTED = p.IN_CUR_ITEM
				d.BO_INPUT_ENTERED = true
				d.BO_UPDATE_INPUT = true
				p.IN_CUR_ITEM = 0
			}
		} else {
			p.BY_INPUT_EXTEN_CODE = key
		}
		b.popup.XUS_TN_MENU_POPUP, b.popup.XUS_TN_SCREEN = p, sc
		b.popup.Execute(now)
	}
	if !d.BO_UPDATE_INPUT && !d.BO_UPDATE_ALL {
		return
	}
	if d.BO_UPDATE_ALL {
		title(now, &b.write, d, sc)
	}
	if n := encoding.ELEMENT_COUNT(sepItem, &d.ST_INPUT_DATA); n > 0 {
		d.IN_SELECTED = LIMIT(1, d.IN_SELECTED, n)
		showSelected(now, &b.write, d, sc)
	}
	d.IN_CURSOR_X, d.IN_CURSOR_Y = d.IN_X, d.IN_Y
	d.BO_UPDATE_ALL, d.BO_UPDATE_INPUT = false, false
}

// TN_INPUT_MENU_POPUP is a popup menu of the items, separated by '#', of
// ST_MENU_TEXT, an item '-' a line: BO_CREATE opens it at IN_X, IN_Y, kept
// on the screen, with a border and a shadow, saving the area below. The
// cursor keys up and down select the item IN_CUR_ITEM, return closes it,
// and escape and the keys left and right close it with IN_CUR_ITEM 0.
// BO_DESTROY closes it.
type TN_INPUT_MENU_POPUP struct {
	XUS_TN_MENU_POPUP *network.US_TN_MENU_POPUP
	XUS_TN_SCREEN     *network.US_TN_SCREEN

	areaSave     TN_SC_AREA_SAVE
	areaRestore  TN_SC_AREA_RESTORE
	line         TN_SC_LINE
	shadow       TN_SC_ADD_SHADOW
	box          TN_SC_BOX
	writeC       TN_SC_WRITE_C
	scrollOffset iec.INT
}

// INIT resets the block.
func (b *TN_INPUT_MENU_POPUP) INIT() {
	*b = TN_INPUT_MENU_POPUP{XUS_TN_MENU_POPUP: b.XUS_TN_MENU_POPUP, XUS_TN_SCREEN: b.XUS_TN_SCREEN}
}

// Execute runs the block once.
func (b *TN_INPUT_MENU_POPUP) Execute(now time.Time) {
	m, sc := b.XUS_TN_MENU_POPUP, b.XUS_TN_SCREEN
	if m == nil || sc == nil {
		return
	}
	item := func(i iec.INT) iec.STRING { return encoding.ELEMENT_GET(sepItem, i-1, &m.ST_MENU_TEXT) }

	if !m.BO_ACTIV && m.BO_CREATE { // open
		m.BO_ACTIV, m.BO_CREATE, m.BO_DESTROY, m.BO_UPDATE = true, false, false, true
		m.IN_MENU_E_COUNT = encoding.ELEMENT_COUNT(sepItem, &m.ST_MENU_TEXT)
		m.IN_ROWS = m.IN_MENU_E_COUNT
		m.IN_COLS = 0
		for i := iec.INT(1); i <= m.IN_ROWS; i++ {
			m.IN_COLS = max(LEN(item(i)), m.IN_COLS)
		}
		if m.IN_X+m.IN_COLS > 75 { // keep it on the screen
			m.IN_X = 79 - m.IN_COLS - 4
		}
		if m.IN_Y+m.IN_ROWS > 21 {
			m.IN_Y = 23 - m.IN_ROWS - 2
		}
		sc.BO_MODAL_DIALOG = true

		a := &b.areaSave
		a.IIN_X1, a.IIN_Y1, a.IIN_X2, a.IIN_Y2, a.XUS_TN_SCREEN = m.IN_X, m.IN_Y, m.IN_X+m.IN_COLS+4, m.IN_Y+m.IN_ROWS+2, sc
		a.Execute(now)
		x := &b.box
		x.IIN_X1, x.IIN_Y1, x.IIN_X2, x.IIN_Y2 = m.IN_X, m.IN_Y, m.IN_X+m.IN_COLS+3, m.IN_Y+m.IN_ROWS+1
		x.IBY_FILL, x.IIN_BORDER, x.IBY_ATTR, x.XUS_TN_SCREEN = ' ', 1, m.BY_ATTR_OF, sc
		x.Execute(now)
		s := &b.shadow
		s.IIN_X1, s.IIN_Y1, s.IIN_X2, s.IIN_Y2 = m.IN_X+1, m.IN_Y+1, m.IN_X+m.IN_COLS+4, m.IN_Y+m.IN_ROWS+2
		s.IIN_OPTION, s.XUS_TN_SCREEN = 0, sc
		s.Execute(now)

		row := m.IN_Y
		for i := iec.INT(1); i <= m.IN_MENU_E_COUNT; i++ {
			row++
			if item(i) == "-" { // a line
				l := &b.line
				l.IIN_Y1, l.IIN_X1, l.IIN_Y2, l.IIN_X2 = row, m.IN_X, row, m.IN_X+m.IN_COLS+3
				l.IBY_ATTR, l.IBY_BORDER, l.XUS_TN_SCREEN = m.BY_ATTR_OF, 1, sc
				l.Execute(now)
			}
		}
	}

	if m.BO_ACTIV {
		switch m.BY_INPUT_EXTEN_CODE {
		case keyUp:
			b.scrollOffset = -1
			m.BO_UPDATE = true
		case keyDown:
			b.scrollOffset = 1
			m.BO_UPDATE = true
		case keyReturn:
			m.BO_DESTROY = true
		case keyEscape, keyRight, keyLeft:
			m.BO_DESTROY = true
			m.IN_CUR_ITEM = 0
		}
	}

	if m.BO_ACTIV && m.BO_UPDATE {
		m.BO_UPDATE = false
		i := m.IN_CUR_ITEM
		if i < 1 || i > m.IN_MENU_E_COUNT {
			b.scrollOffset = 1
			i = m.IN_MENU_E_COUNT
		}
		if item(i) == "-" {
			b.scrollOffset = 1
		}
		if b.scrollOffset != 0 { // to the next item that is not a line
			for range m.IN_MENU_E_COUNT {
				i = math.INC2(i, b.scrollOffset, 1, m.IN_MENU_E_COUNT)
				if item(i) != "-" {
					break
				}
			}
		}
		m.IN_CUR_ITEM = i
		b.scrollOffset = 0

		row := m.IN_Y
		for n := iec.INT(1); n <= m.IN_MENU_E_COUNT; n++ {
			row++
			w := &b.writeC
			w.IST_STRING = item(n)
			if w.IST_STRING != "-" {
				w.IIN_Y, w.IIN_X, w.IIN_LENGTH, w.IIN_OPTION, w.XUS_TN_SCREEN = row, m.IN_X+1, m.IN_COLS+2, 2, sc
				w.IBY_ATTR = SEL(n == m.IN_CUR_ITEM, m.BY_ATTR_OF, m.BY_ATTR_MF)
				w.Execute(now)
			}
		}
	}

	if m.BO_ACTIV && m.BO_DESTROY { // close
		m.BO_DESTROY, m.BO_ACTIV = false, false
		b.areaRestore.XUS_TN_SCREEN = sc
		b.areaRestore.Execute(now)
		sc.BO_MODAL_DIALOG = false
	}
}

// TN_INPUT_MENU_BAR is a menu bar at IN_X, IN_Y: ST_MENU_TEXT is the bar,
// its items separated by '#', then, separated by '%', the popup menu of
// each item. BO_CREATE shows it, escape opens and closes it, the keys left
// and right move along the bar and up and down in its popup, and return
// closes it with IN_MENU_SELECTED the item times 10 plus that of its
// popup. BO_DESTROY closes and removes it.
type TN_INPUT_MENU_BAR struct {
	XUS_TN_MENU   *network.US_TN_MENU
	XUS_TN_SCREEN *network.US_TN_SCREEN

	popup         TN_INPUT_MENU_POPUP
	popupData     network.US_TN_MENU_POPUP
	writeC        TN_SC_WRITE_C
	menuBar       iec.STRING
	scrollOffset  iec.INT
	resetPosition iec.BOOL
}

// INIT resets the block.
func (b *TN_INPUT_MENU_BAR) INIT() {
	*b = TN_INPUT_MENU_BAR{XUS_TN_MENU: b.XUS_TN_MENU, XUS_TN_SCREEN: b.XUS_TN_SCREEN}
}

// Execute runs the block once.
func (b *TN_INPUT_MENU_BAR) Execute(now time.Time) {
	m, sc, p := b.XUS_TN_MENU, b.XUS_TN_SCREEN, &b.popupData
	if m == nil || sc == nil {
		return
	}
	if m.IN_STATE == 0 && m.BO_CREATE { // show the bar
		m.IN_STATE = 1
		m.BO_CREATE = false
		m.IN_MENU_SELECTED = 0
		m.BO_UPDATE = true
		b.resetPosition = true
		b.menuBar = encoding.ELEMENT_GET(sepMenu, 0, &m.ST_MENU_TEXT)
		m.IN_MENU_E_COUNT = encoding.ELEMENT_COUNT(sepItem, &b.menuBar)
	}

	key := sc.BY_INPUT_EXTEN_CODE
	p.BY_INPUT_EXTEN_CODE = 0
	if key == keyUp || key == keyDown { // the popup takes up and down only
		p.BY_INPUT_EXTEN_CODE = key
	}
	// close closes the menu, with the item selected.
	close := func(selected iec.INT) {
		p.BO_DESTROY = true
		m.BO_UPDATE = true
		m.IN_STATE = 1
		b.resetPosition = true
		sc.BO_MENUE_BAR_DIALOG = false
		m.IN_MENU_SELECTED = selected
	}

	if key == keyEscape {
		if m.IN_STATE == 1 && !sc.BO_MODAL_DIALOG { // open
			m.IN_STATE = 2
			m.BO_UPDATE = true
			m.IN_MENU_SELECTED = 0
			sc.BO_MENUE_BAR_DIALOG = true
		} else if m.IN_STATE == 2 {
			close(m.IN_MENU_SELECTED)
		}
	}
	if m.IN_STATE == 2 {
		if key == keyReturn {
			close(m.IN_CUR_MENU_ITEM*10 + m.IN_CUR_SUB_ITEM)
			sc.BY_INPUT_EXTEN_CODE = 0
		} else if m.BO_DESTROY {
			close(0)
			sc.BY_INPUT_EXTEN_CODE = 0
		}
		switch key {
		case keyRight:
			b.scrollOffset = 1
			m.BO_UPDATE = true
			sc.BY_INPUT_EXTEN_CODE = 0
		case keyLeft:
			b.scrollOffset = -1
			m.BO_UPDATE = true
			sc.BY_INPUT_EXTEN_CODE = 0
		}
	}

	if b.resetPosition {
		b.resetPosition = false
		m.IN_X_SM_OLD, m.IN_Y_SM_OLD, m.IN_X_SM_NEW, m.IN_Y_SM_NEW = -1, -1, -1, -1
	}

	if m.IN_STATE > 0 && m.BO_UPDATE { // draw the bar
		m.BO_UPDATE = false
		row, col := m.IN_Y, m.IN_X
		if b.scrollOffset != 0 {
			m.IN_CUR_MENU_ITEM = math.INC2(m.IN_CUR_MENU_ITEM, b.scrollOffset, 1, m.IN_MENU_E_COUNT)
			b.scrollOffset = 0
		} else if m.IN_CUR_MENU_ITEM == 0 {
			m.IN_CUR_MENU_ITEM = 1
		}
		for n := iec.INT(1); n <= m.IN_MENU_E_COUNT; n++ {
			a := m.BY_ATTR_OF
			if n == m.IN_CUR_MENU_ITEM && m.IN_STATE > 1 {
				a = m.BY_ATTR_MF
				m.IN_Y_SM_NEW, m.IN_X_SM_NEW = row+1, col // its popup
			}
			t := encoding.ELEMENT_GET(sepItem, n-1, &b.menuBar)
			w := &b.writeC
			w.IIN_Y, w.IIN_X, w.IBY_ATTR, w.IIN_LENGTH, w.IIN_OPTION, w.IST_STRING, w.XUS_TN_SCREEN = row, col, a, LEN(t)+2, 2, t, sc
			w.Execute(now)
			col += LEN(t) + 2
		}
	}

	b.popup.XUS_TN_MENU_POPUP, b.popup.XUS_TN_SCREEN = p, sc
	if m.IN_X_SM_NEW != m.IN_X_SM_OLD { // another popup
		if p.BO_ACTIV {
			p.BO_DESTROY = true
			b.popup.Execute(now)
		}
		if !p.BO_ACTIV {
			p.IN_Y, p.IN_X = m.IN_Y_SM_NEW, m.IN_X_SM_NEW
			p.BY_ATTR_MF, p.BY_ATTR_OF = m.BY_ATTR_MF, m.BY_ATTR_OF
			p.ST_MENU_TEXT = encoding.ELEMENT_GET(sepMenu, m.IN_CUR_MENU_ITEM, &m.ST_MENU_TEXT)
			p.BO_CREATE = true
			m.IN_X_SM_OLD, m.IN_Y_SM_OLD = m.IN_X_SM_NEW, m.IN_Y_SM_NEW
		}
	}
	b.popup.Execute(now)
	m.IN_CUR_SUB_ITEM = p.IN_CUR_ITEM

	if m.IN_STATE == 1 && m.BO_DESTROY { // remove the bar
		m.IN_STATE = 0
		m.BO_DESTROY = false
		m.IN_CUR_MENU_ITEM = 0
		m.IN_CUR_SUB_ITEM = 0
	}
}
