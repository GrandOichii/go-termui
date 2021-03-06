package termui

import (
	"fmt"
	"strings"

	nc "github.com/rthornton128/goncurses"
)

type DDBChoiceType int

const (
	focusedAttribute = nc.A_REVERSE

	SingleElement DDBChoiceType = iota
	MultipleElements
)

var (
	allowedRanges = [][2]rune{
		{'a', 'z'},
		{'A', 'Z'},
		{'0', '9'},
	}
	allowedSingles = []rune{
		'=',
		'"',
		' ',
	}
)

// Checks whether the character can be added to the line edit template
//
// Returns true if the character was added
func isValidLineEditCh(ch rune) bool {
	for _, ar := range allowedRanges {
		if ch >= ar[0] && ch <= ar[1] {
			return true
		}
	}
	for _, as := range allowedSingles {
		if ch == as {
			return true
		}
	}
	return false
}

// Returns the max element
func MaxInt(a ...int) int {
	result := a[0]
	for _, v := range a {
		if v > result {
			result = v
		}
	}
	return result
}

// Returns the min element
func MinInt(a ...int) int {
	result := a[0]
	for _, v := range a {
		if v < result {
			result = v
		}
	}
	return result
}

// Returns the sum of all the elements
func SumInt(a ...int) int {
	result := 0
	for _, i := range a {
		result += i
	}
	return result
}

// More convinient way to add to window with multiple attributes
func Put(win *nc.Window, y, x int, line string, attrs ...nc.Char) {
	for _, attr := range attrs {
		win.AttrOn(attr)
		defer win.AttrOff(attr)
	}
	win.MovePrint(y, x, line)
}

func ReverseColorPair(colorPair string) string {
	split := strings.Split(colorPair, "-")
	if len(split) == 1 {
		return "normal-" + colorPair
	}
	result := ""
	for _, s := range split {
		result = "-" + s + result
	}
	return result[1:]
}

// Draws a box
func DrawBox(win *nc.Window, y, x, height, width int, borderColor string) error {
	bcolor, err := ParseColorPair(borderColor)
	if err != nil {
		return err
	}
	win.AttrOn(bcolor)
	win.MoveAddChar(y, x, nc.ACS_ULCORNER)
	win.MoveAddChar(y+height-1, x, nc.ACS_LLCORNER)
	win.MoveAddChar(y, x+width-1, nc.ACS_URCORNER)
	win.MoveAddChar(y+height-1, x+width-1, nc.ACS_LRCORNER)
	for i := 1; i < height-1; i++ {
		win.MoveAddChar(y+i, x, nc.ACS_VLINE)
		win.MoveAddChar(y+i, x+width-1, nc.ACS_VLINE)
	}
	for i := 1; i < width-1; i++ {
		win.MoveAddChar(y, x+i, nc.ACS_HLINE)
		win.MoveAddChar(y+height-1, x+i, nc.ACS_HLINE)
	}
	win.AttrOff(bcolor)
	return nil
}

// Draws the borders of the window
func DrawBorders(win *nc.Window, colorPair string) error {
	var err error
	color, err := ParseColorPair(colorPair)
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
	wwidth := MaxInt(choicesLen, cctMessage.Length()+4)
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
		Put(win, wheight-3, 1, whiteSpace)
		for i, choice := range cctChoices {
			// nc.Flash()
			sl := choice.Length()
			if i == choiceID {
				Put(win, wheight-3, pos-2, "["+strings.Repeat(" ", sl)+"]")
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
	moptions := make([]DrawableAsLine, 0, len(cctOptions))
	for _, o := range cctOptions {
		moptions = append(moptions, o)
	}
	lt := CreateListTemplate(moptions, maxDisplayAmount)
	whiteSpace := strings.Repeat(" ", width-2)
	bc, err := ParseColorPair(borderColor)
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
			Put(win, i, 1, whiteSpace)
		}
		// draw
		lt.Draw(win, 1, 1, true)
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
			lt.ScrollUp()
		case nc.KEY_DOWN:
			lt.ScrollDown()
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
	let := CreateLineEditTemplate(text, maxLength)
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
