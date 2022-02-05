package termui

import (
	"fmt"
	"strings"

	nc "github.com/rthornton128/goncurses"
)

type DDBChoiceType int

const (
	lineEditFocusedAttribute = nc.A_REVERSE

	SingleElement DDBChoiceType = iota
	MultipleElements
)

// List template. Use for drawing lists
type listTemplate struct {
	win              *nc.Window
	options          []*CCTMessage
	maxDisplayAmount int
	cursor           int
	choice           int
	pageN            int
}

// Creates a list template
func createListTemplate(win *nc.Window, options []*CCTMessage, maxDisplayAmount int) *listTemplate {
	result := listTemplate{}
	result.win = win
	result.options = options
	result.maxDisplayAmount = maxDisplayAmount
	result.cursor = 0
	result.choice = 0
	result.pageN = 0
	return &result
}

// Draws the list tamplate
func (l listTemplate) draw(y, x int, focusSelected bool) error {
	return drawList(l.win, y, x, l.options, l.maxDisplayAmount, l.cursor, l.pageN, focusSelected)
}

// Moves the cursor of the list template up
func (l *listTemplate) scrollUp() {
	l.choice--
	l.cursor--
	if l.cursor < 0 {
		if len(l.options) > l.maxDisplayAmount {
			if l.pageN == 0 {
				l.cursor = l.maxDisplayAmount - 1
				l.choice = len(l.options) - 1
				l.pageN = len(l.options) - l.maxDisplayAmount
			} else {
				l.pageN--
				l.cursor++
			}
		} else {
			l.cursor = len(l.options) - 1
			l.choice = l.cursor
		}
	}
}

// Moves the cursor of the list tamplate down
func (l *listTemplate) scrollDown() {
	l.choice++
	l.cursor++
	if len(l.options) > l.maxDisplayAmount {
		if l.cursor >= l.maxDisplayAmount {
			l.cursor--
			l.pageN++
			if l.choice == len(l.options) {
				l.choice = 0
				l.cursor = 0
				l.pageN = 0
			}
		}
	} else {
		if l.cursor >= len(l.options) {
			l.cursor = 0
			l.choice = 0
		}
	}
}

// Line edit template. Use for drawing and interacting with writable lines
type lineEditTemplate struct {
	content string
	blank   string
	cursor  int
	maxLen  int
}

// Creates the line edit template
func createLineEditTemplate(text string, maxLen int) *lineEditTemplate {
	result := lineEditTemplate{}
	result.cursor = 0
	result.content = text
	result.blank = strings.Repeat(" ", maxLen)
	result.maxLen = maxLen
	return &result
}

// Moves the cursor to the left
func (l *lineEditTemplate) MoveCursorLeft() {
	l.cursor--
	if l.cursor == 0 {
		l.cursor = 0
	}
}

// Moves the cursor to the right
func (l *lineEditTemplate) MoveCursorRight() {
	l.cursor++
	if l.cursor > len(l.content) {
		l.cursor = len(l.content)
	}
}

// Adds the character to the cursor location
func (l *lineEditTemplate) AddCh(ch rune) {
	if l.cursor < l.maxLen && isValidLineEditCh(ch) {
		l.content = l.content[:l.cursor] + string(ch) + l.content[l.cursor:]
		l.MoveCursorRight()
	}
}

// Draws the line edit template
func (l lineEditTemplate) Draw(win *nc.Window, yPos, xPos int, focused bool) error {
	win.MovePrintf(yPos, xPos, l.blank)
	win.MovePrintf(yPos, xPos, l.content)
	if l.cursor < l.maxLen {
		win.Move(yPos, xPos+l.cursor)
		win.AttrOn(lineEditFocusedAttribute)
		win.Print(" ")
		win.AttrOff(lineEditFocusedAttribute)
	}
	return nil
}

// Removes the element at the cursor
func (l *lineEditTemplate) DeleteSelected() {
	if l.cursor == 0 {
		return
	}
	l.content = l.content[:l.cursor-1] + l.content[l.cursor:]
	l.MoveCursorLeft()
}

// Checks whether the character can be added to the line edit template
func isValidLineEditCh(ch rune) bool {
	if ch >= 'a' && ch <= 'z' {
		return true
	}
	if ch >= 'A' && ch <= 'Z' {
		return true
	}
	if ch >= '0' && ch <= '9' {
		return true
	}
	return ch == ' '
}

// Draws the list
func drawList(win *nc.Window, y, x int, options []*CCTMessage, maxDisplayAmount, cursor, pageN int, focusSelected bool) error {
	for i := 0; i < minInt(maxDisplayAmount, len(options)); i++ {
		attr := nc.A_NORMAL
		if i == cursor && focusSelected {
			attr = nc.A_REVERSE
		}
		options[i+pageN].Draw(win, y+i, x, attr)
		// put(win, y+i, x, options[i+pageN], attr)
	}
	return nil
}

// Returns the max element
func maxInt(a ...int) int {
	result := a[0]
	for _, v := range a {
		if v > result {
			result = v
		}
	}
	return result
}

// Returns the min element
func minInt(a ...int) int {
	result := a[0]
	for _, v := range a {
		if v < result {
			result = v
		}
	}
	return result
}

// More convinient way to add to window with multiple attributes
func put(win *nc.Window, y, x int, line string, attrs ...nc.Char) {
	for _, attr := range attrs {
		win.AttrOn(attr)
		defer win.AttrOff(attr)
	}
	win.MovePrint(y, x, line)
}

// Draws the borders of the window
func DrawBorders(win *nc.Window, colorPair string) error {
	var err error
	color, err := parseColors(colorPair)
	if err != nil {
		return err
	}
	win.AttrOn(color)
	err = win.Border(nc.ACS_VLINE, nc.ACS_VLINE, nc.ACS_HLINE, nc.ACS_HLINE, nc.ACS_ULCORNER, nc.ACS_URCORNER, nc.ACS_LLCORNER, nc.ACS_LRCORNER)
	if err != nil {
		return err
	}
	win.AttrOff(color)
	return nil
}

// Displays a message box
// Choices can't be more than 3 elements
// If choices is empty, it becomes {"Ok"}
// Returns the picked element
func MessageBox(parent *Window, message string, choices []string, borderColor string) (string, error) {
	if len(choices) == 0 {
		choices = []string{"Ok"}
	}
	hasCancel := false
	for _, choice := range choices {
		if choice == "Cancel" {
			hasCancel = true
		}
	}
	parentWin := parent.win
	height, width := parentWin.MaxYX()
	if len(choices) > 3 {
		return "", fmt.Errorf("termui - %v can't be choices for MessageBox", choices)
	}
	cctChoices, err := GetCCTs(choices)
	if err != nil {
		return "", err
	}
	cctMessage, err := ToCCTMessage(message)
	if err != nil {
		return "", err
	}
	choiceID := 0
	done := false
	// maxWidth := width - 2
	choicesLen := (len(cctChoices) + 1) * 2
	for _, ch := range cctChoices {
		choicesLen += ch.Length()
	}
	wwidth := maxInt(choicesLen, cctMessage.Length()+4)
	wheight := 7
	ypos := (height - wheight) / 2
	xpos := (width - wwidth) / 2
	win, err := nc.NewWindow(wheight, wwidth, ypos, xpos)
	if err != nil {
		return "", err
	}
	err = win.Keypad(true)
	if err != nil {
		return "", err
	}
	defer win.Clear()
	DrawBorders(win, borderColor)
	cctMessage.Draw(win, 2, 2)
	// put(win, 2, 2, message)
	whiteSpace := strings.Repeat(" ", wwidth-2)

	for !done {
		// draw
		pos := 3
		put(win, wheight-3, 1, whiteSpace)
		for i, choice := range cctChoices {
			// nc.Flash()
			sl := choice.Length()
			if i == choiceID {
				put(win, wheight-3, pos-2, "["+strings.Repeat(" ", sl)+"]")
			}
			choice.Draw(win, wheight-3, pos-1)
			// put(win, wheight-3, pos-1, s)
			pos += sl + 2
		}
		// key handling
		key := win.GetChar()
		if key == nc.KEY_LEFT {
			choiceID--
			if choiceID < 0 {
				choiceID = len(choices) - 1
			}
		}
		if key == nc.KEY_RIGHT {
			choiceID++
			if choiceID >= len(choices) {
				choiceID = 0
			}
		}
		if key == 10 { // nc.KEY_ENTER doesn't work
			done = true
		}
		if key == nc.KEY_ESC && hasCancel {
			return "Cancel", nil
		}
	}
	return choices[choiceID], nil
}

// Displays a drop down box
// Returns the indicies of the picked options
func DropDownBox(options []string, maxDisplayAmount, y, x int, choiceType DDBChoiceType, borderColor string) ([]int, error) {
	if len(options) == 0 {
		return nil, nil
	}
	// height := minInt(len(options), maxDisplayAmount) + 2
	height := maxDisplayAmount + 2
	cctOptions, err := GetCCTs(options)
	if err != nil {
		return nil, err
	}
	width := cctOptions[0].Length()
	for i, line := range cctOptions {
		if i == 0 {
			continue
		}
		l := line.Length()
		if l > width {
			width = l
		}
	}
	width += 3
	win, err := nc.NewWindow(height, width, y, x)
	if err != nil {
		return nil, err
	}
	defer win.Clear()
	win.Keypad(true)
	DrawBorders(win, borderColor)
	lt := createListTemplate(win, cctOptions, maxDisplayAmount)
	whiteSpace := strings.Repeat(" ", width-2)
	bc, err := parseColors(borderColor)
	if err != nil {
		return nil, err
	}
	for {
		// clear lines
		win.AttrOn(bc)
		win.MoveAddChar(1, width-1, nc.ACS_VLINE)
		win.MoveAddChar(height-2, width-1, nc.ACS_VLINE)
		win.AttrOff(bc)
		for i := 1; i < height-1; i++ {
			put(win, i, 1, whiteSpace)
		}
		// draw
		lt.draw(1, 1, true)
		win.AttrOn(bc)
		if len(options) > maxDisplayAmount {
			if lt.pageN != 0 {
				win.MoveAddChar(1, width-1, nc.ACS_UARROW)
			}
			if lt.pageN != len(options)-maxDisplayAmount {
				win.MoveAddChar(height-2, width-1, nc.ACS_DARROW)
			}
		}
		win.AttrOff(bc)
		// handle key
		key := win.GetChar()
		switch key {
		case nc.KEY_ESC:
			return nil, nil
		case nc.KEY_UP:
			lt.scrollUp()
		case nc.KEY_DOWN:
			lt.scrollDown()
		case 10:
			if lt.choice == -1 {
				break
			}
			return []int{lt.choice}, nil
		}
	}
}

// Displays a box where the user will have to enter a string
// Returns the entered string
func EnterString(parent *Window, text string, prompt string, maxLength int, borderColor string) (string, error) {
	pheight, pwidth := parent.win.MaxYX()
	cctprompt, err := ToCCTMessage(prompt)
	if err != nil {
		return "", err
	}
	height := 5
	width := 2 + cctprompt.Length() + 2 + maxLength + 2
	y := (pheight - height) / 2
	x := (pwidth - width) / 2
	w, err := nc.NewWindow(height, width, y, x)
	if err != nil {
		return "", err
	}
	defer w.Clear()
	w.Keypad(true)
	DrawBorders(w, borderColor)
	y = 2
	x = cctprompt.Length() + 4
	cctprompt.Draw(w, y, 2)
	w.MovePrint(y, x-2, ": ")
	let := createLineEditTemplate(text, maxLength)
l:
	for {
		let.Draw(w, y, x, true)
		key := w.GetChar()
		switch key {
		case KeyEnter:
			break l
		case KeyLeft:
			let.MoveCursorLeft()
		case KeyRight:
			let.MoveCursorRight()
		case KeyBackspace:
			let.DeleteSelected()
		default:
			let.AddCh(rune(key))
		}
	}
	return let.content, nil
}
