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

type listTemplate struct {
	win              *nc.Window
	options          []*CCTMessage
	maxDisplayAmount int
	cursor           int
	choice           int
	pageN            int
}

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

func (l listTemplate) draw(y, x int, focusSelected bool) error {
	return drawList(l.win, y, x, l.options, l.maxDisplayAmount, l.cursor, l.pageN, focusSelected)
}

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

type lineEditTemplate struct {
	content string
	blank   string
	cursor  int
	maxLen  int
}

func createLineEditTemplate(text string, maxLen int) *lineEditTemplate {
	result := lineEditTemplate{}
	result.cursor = 0
	result.content = text
	result.blank = strings.Repeat(" ", maxLen)
	result.maxLen = maxLen
	return &result
}

func (l *lineEditTemplate) MoveCursorLeft() {
	l.cursor--
	if l.cursor == 0 {
		l.cursor = 0
	}
}

func (l *lineEditTemplate) MoveCursorRight() {
	l.cursor++
	if l.cursor > len(l.content) {
		l.cursor = len(l.content)
	}
}

func (l *lineEditTemplate) AddCh(ch rune) {
	if l.cursor < l.maxLen && isValidLineEditCh(ch) {
		l.content = l.content[:l.cursor] + string(ch) + l.content[l.cursor:]
		l.MoveCursorRight()
	}
}

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

func (l *lineEditTemplate) DeleteSelected() {
	if l.cursor == 0 {
		return
	}
	l.content = l.content[:l.cursor-1] + l.content[l.cursor:]
	l.MoveCursorLeft()
}

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

func maxInt(a ...int) int {
	result := a[0]
	for _, v := range a {
		if v > result {
			result = v
		}
	}
	return result
}

func minInt(a ...int) int {
	result := a[0]
	for _, v := range a {
		if v < result {
			result = v
		}
	}
	return result
}

func put(win *nc.Window, y, x int, line string, attrs ...nc.Char) {
	for _, attr := range attrs {
		win.AttrOn(attr)
		defer win.AttrOff(attr)
	}
	win.MovePrint(y, x, line)
}

func DrawBorders(win *nc.Window) error {
	return win.Border(nc.ACS_VLINE, nc.ACS_VLINE, nc.ACS_HLINE, nc.ACS_HLINE, nc.ACS_ULCORNER, nc.ACS_URCORNER, nc.ACS_LLCORNER, nc.ACS_LRCORNER)
}

func MessageBox(parent *Window, message string, choices []string) (string, error) {
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
	DrawBorders(win)
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

func DropDownBox(options []string, maxDisplayAmount, y, x int, choiceType DDBChoiceType) ([]int, error) {
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
	DrawBorders(win)
	lt := createListTemplate(win, cctOptions, maxDisplayAmount)
	whiteSpace := strings.Repeat(" ", width-2)
	for {
		// clear lines
		win.MoveAddChar(1, width-1, nc.ACS_VLINE)
		win.MoveAddChar(height-2, width-1, nc.ACS_VLINE)
		for i := 1; i < height-1; i++ {
			put(win, i, 1, whiteSpace)
		}
		// draw
		lt.draw(1, 1, true)
		if len(options) > maxDisplayAmount {
			if lt.pageN != 0 {
				win.MoveAddChar(1, width-1, nc.ACS_UARROW)
			}
			if lt.pageN != len(options)-maxDisplayAmount {
				win.MoveAddChar(height-2, width-1, nc.ACS_DARROW)
			}
		}
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

func EnterString(parent *Window, text string, prompt string, maxLength int) (string, error) {
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
	DrawBorders(w)
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
