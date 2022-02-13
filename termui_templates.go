package termui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	nc "github.com/rthornton128/goncurses"
)

type Alignment int

const (
	progressBarUnit = '#'

	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
)

// List template. Use for drawing lists
type ListTemplate struct {
	options          []DrawableAsLine
	maxDisplayAmount int
	cursor           int
	choice           int
	pageN            int
}

// Creates a list template
func CreateListTemplate(options []DrawableAsLine, maxDisplayAmount int) *ListTemplate {
	result := ListTemplate{}
	result.options = options
	result.maxDisplayAmount = maxDisplayAmount
	result.cursor = 0
	result.choice = 0
	result.pageN = 0
	return &result
}

// Draws the list tamplate
func (l ListTemplate) Draw(win *nc.Window, y, x int, focusSelected bool) error {
	for i := 0; i < minInt(l.maxDisplayAmount, len(l.options)); i++ {
		attr := nc.A_NORMAL
		if i == l.cursor && focusSelected {
			attr = nc.A_REVERSE
		}
		l.options[i+l.pageN].Draw(win, y+i, x, attr)
		// put(win, y+i, x, options[i+pageN], attr)
	}
	return nil
}

// Sets the options
func (l *ListTemplate) SetOptions(options []DrawableAsLine) {
	if len(l.options) > len(options) {
		l.cursor = 0
		l.choice = 0
		l.pageN = 0
	}
	l.options = options
}

// Adds an option
func (l *ListTemplate) AddOption(option DrawableAsLine) {
	l.options = append(l.options, option)
}

// Moves the cursor of the list template up
func (l *ListTemplate) ScrollUp() {
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
func (l *ListTemplate) ScrollDown() {
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
type LineEditTemplate struct {
	content string
	blank   string
	cursor  int
	maxLen  int
}

// Creates the line edit template
func CreateLineEditTemplate(text string, maxLen int) *LineEditTemplate {
	result := LineEditTemplate{}
	result.cursor = 0
	result.content = text
	result.blank = strings.Repeat("_", maxLen)
	result.maxLen = maxLen
	return &result
}

// Moves the cursor to the left
func (l *LineEditTemplate) MoveCursorLeft() {
	l.cursor--
	if l.cursor == 0 {
		l.cursor = 0
	}
}

// Moves the cursor to the right
func (l *LineEditTemplate) MoveCursorRight() {
	l.cursor++
	if l.cursor > len(l.content) {
		l.cursor = len(l.content)
	}
}

// Adds the character to the cursor location
func (l *LineEditTemplate) AddCh(ch rune) {
	if l.cursor < l.maxLen && isValidLineEditCh(ch) {
		l.content = l.content[:l.cursor] + string(ch) + l.content[l.cursor:]
		l.MoveCursorRight()
	}
}

// Draws the line edit template
func (l LineEditTemplate) Draw(win *nc.Window, yPos, xPos int, focused bool) error {
	win.MovePrintf(yPos, xPos, l.blank)
	win.MovePrintf(yPos, xPos, l.content)
	if focused && l.cursor < l.maxLen {
		win.Move(yPos, xPos+l.cursor)
		win.AttrOn(focusedAttribute)
		win.Print(" ")
		win.AttrOff(focusedAttribute)
	}
	return nil
}

// Removes the element at the cursor
func (l *LineEditTemplate) DeleteSelected() {
	if l.cursor == 0 {
		return
	}
	l.content = l.content[:l.cursor-1] + l.content[l.cursor:]
	l.MoveCursorLeft()
}

// Sets the text of the template
func (l *LineEditTemplate) SetText(text string) error {
	if len(text) > l.maxLen {
		return fmt.Errorf("termui - can't set lineEditTemplate text to %v - maxLen is %v", text, l.maxLen)
	}
	l.content = text
	l.cursor = len(l.content)
	return nil
}

// Word choice template use for prompting user to pick a word from options
type WordChoiceTemplate struct {
	options []*CCTMessage
	choice  int
	maxLen  int
	acolor  nc.Char
	al      Alignment
}

// Creates a word choice template
func CreateWordChoiceTemplate(options []string, alignment Alignment, arrowColor string) (*WordChoiceTemplate, error) {
	result := WordChoiceTemplate{}
	var err error
	if len(options) == 0 {
		return nil, errors.New("termui - can't create a WordChoice with no options")
	}
	result.options, err = GetCCTs(options)
	if err != nil {
		return nil, err
	}
	result.choice = 0
	result.maxLen = result.options[0].Length()
	for i, o := range result.options {
		if i == 0 {
			continue
		}
		result.maxLen = maxInt(result.maxLen, o.Length())
	}
	result.al = alignment
	result.acolor, err = parseColorPair(arrowColor)
	return &result, err
}

// Returns the currently focused option
func (w WordChoiceTemplate) GetSelected() *CCTMessage {
	return w.options[w.choice]
}

// Focuses on the next option
func (w *WordChoiceTemplate) FocusNext() {
	w.choice++
	if w.choice == len(w.options) {
		w.choice = 0
	}
}

// Focuses on the previous option
func (w *WordChoiceTemplate) FocusPrev() {
	w.choice--
	if w.choice < 0 {
		w.choice = len(w.options) - 1
	}
}

// Draws the word choice template
func (w WordChoiceTemplate) Draw(win *nc.Window, y, x int, focused bool) error {
	if focused {
		win.AttrOn(focusedAttribute)
	}
	win.AttrOn(w.acolor)
	win.MoveAddChar(y, x, '<')
	win.MoveAddChar(y, x+w.maxLen+1, '>')
	win.AttrOff(w.acolor)
	if focused {
		win.AttrOff(focusedAttribute)
	}
	option := w.options[w.choice]
	xl := x + 1
	switch w.al {
	case AlignCenter:
		xl += (w.maxLen - option.Length()) / 2
	case AlignRight:
		xl += (w.maxLen - option.Length())
	}
	option.Draw(win, y, xl)
	return nil
}

// A progress bar template
type ProgressBarTemplate struct {
	barLength int
	max       int
	maxs      string
	current   int
	clears    string
	si        bool
	bcolor    nc.Char
	icolor    nc.Char
}

// Creates a progress bar template
func CreateProgressBarTemplate(barLength, max int, showInfo bool, barColor string, infoColor string) (*ProgressBarTemplate, error) {
	result := ProgressBarTemplate{}
	result.barLength = barLength
	result.max = max
	result.current = 0
	var err error
	result.bcolor, err = parseColorPair(barColor)
	if err != nil {
		return nil, err
	}
	result.icolor, err = parseColorPair(infoColor)
	if err != nil {
		return nil, err
	}
	result.clears = "[" + strings.Repeat(" ", barLength) + "]"
	if showInfo {
		ispace := strings.Repeat(" ", len(strconv.Itoa(max)))
		result.clears += " (" + ispace + "/" + ispace + ")"
	}
	result.si = showInfo
	result.maxs = strconv.Itoa(max)
	return &result, nil
}

// Sets the current value of the template
func (p *ProgressBarTemplate) Set(value int) {
	p.current = value
	if value > p.max {
		p.current = p.max
	}
}

// Draws the template
func (p ProgressBarTemplate) Draw(win *nc.Window, y, x int) error {
	win.AttrOn(p.icolor)
	// draw the white space
	win.MovePrint(y, x, p.clears)
	if p.si {
		// draw the info
		win.MovePrint(y, x+p.barLength+4, strconv.Itoa(p.current))
		win.MovePrint(y, len(p.clears)-len(p.maxs)+1, strconv.Itoa(p.max))
	}
	win.AttrOff(p.icolor)
	// draw the bar
	l := p.current * p.barLength / p.max
	s := strings.Repeat(string(progressBarUnit), l)
	put(win, y, x+1, s, p.bcolor)
	return nil

}
