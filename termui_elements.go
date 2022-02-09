package termui

import (
	"fmt"
	"math"
	"strconv"

	nc "github.com/rthornton128/goncurses"
)

const (
	colorStep     = 1
	startingColor = 10
)

// A standard label
type Label struct {
	data    *UIElementData
	cctText *CCTMessage
}

// Creates a new label
func NewLabel(text string, y, x int) (*Label, error) {
	result := Label{}
	var err error
	result.cctText, err = ToCCTMessage(text)
	if err != nil {
		return nil, err
	}
	result.data = createUIED(y, x)
	return &result, nil
}

// Draws the label
func (l Label) Draw(win *nc.Window) error {
	l.cctText.Draw(win, l.data.yPos, l.data.xPos)
	// put(pWin, data.yPos, data.xPos, l.text, attr)
	return nil
}

// Doesn't do anything
func (l Label) HandleKey(key nc.Key) error {
	return nil
}

// Sets the text of the label
func (l *Label) SetText(text string) error {
	var err error
	l.cctText, err = ToCCTMessage(text)
	return err
}

// Returns the element data of the label
func (l Label) GetElementData() *UIElementData {
	return l.data
}

// Returns the height of the label
func (l Label) Height() int {
	return 1
}

// Returns the width of the label
func (l Label) Width() int {
	return l.cctText.Length()
}

// A clickable button
type Button struct {
	click    func() error
	clickKey nc.Key
	data     *UIElementData
	cctText  *CCTMessage
}

// Creates a new button
func NewButton(text string, y, x int, click func() error, clickKey nc.Key) (*Button, error) {
	result := Button{}
	err := result.SetText(text)
	if err != nil {
		return nil, err
	}
	result.click = click
	result.clickKey = clickKey
	result.data = createUIED(y, x)
	return &result, nil
}

// Draws the button
func (b Button) Draw(win *nc.Window) error {
	attr := nc.A_NORMAL
	if b.data.focused {
		attr = hightlightKey
	}
	b.cctText.Draw(win, b.data.yPos, b.data.xPos, attr)
	// put(pWin, data.yPos, data.xPos, b.text, attr)
	return nil
}

// Sets the text of the button
func (b *Button) SetText(text string) error {
	var err error
	b.cctText, err = ToCCTMessage(text)
	return err
}

// On ENTER or mouse click calls click
func (b Button) HandleKey(key nc.Key) error {
	if key == b.clickKey {
		return b.click()
	}
	// if key == nc.KEY_MOUSE {
	// 	return b.click()
	// }
	return nil
}

// Returns the element data of the button
func (b Button) GetElementData() *UIElementData {
	return b.data
}

// Returns the height of the button
func (b Button) Height() int {
	return 1
}

// Returns the width of the button
func (b Button) Width() int {
	return b.cctText.Length()
}

// A pie chart element
type PieChart struct {
	values []int
	rads   []float64
	total  int
	height int
	width  int
	bcolor string
	colors []nc.Char
	data   *UIElementData
}

// Creates a color pie chart element.
func NewPieChart(win *Window, y, x, height, width int, values []int, colorPairs []string, borderColor string) (*PieChart, error) {
	var err error
	result := PieChart{}
	result.height = height
	result.width = width
	result.total = 0
	result.data = createUIED(y, x)

	result.bcolor = borderColor
	result.SetValues(values)
	if len(colorPairs) == 0 {
		// create custom colors
		result.colors = make([]nc.Char, 0, len(values))
		for i := startingColor; i < len(values)*colorStep+startingColor; i += colorStep {
			pair := strconv.Itoa(i) + "-normal"
			color, err := parseColorPair(pair)
			if err != nil {
				return nil, err
			}
			result.colors = append(result.colors, color)
		}
	} else {
		err = result.setColors(colorPairs)
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}

// Returns the element data of the pie chart
func (p PieChart) GetElementData() *UIElementData {
	return p.data
}

// Draws the pie chart
func (p PieChart) Draw(win *nc.Window) error {
	var err error
	yPos := p.data.yPos
	xPos := p.data.xPos
	centerY := p.height/2 + yPos
	centerX := p.width/2 + xPos
	radius := minInt(p.height/2, p.width/2) - 1
	DrawBox(win, yPos, xPos, p.height, p.width, p.bcolor)
	win.MovePrintf(yPos, xPos, "%v", p.values)
	for i := 0; i < p.height; i++ {
		for j := 0; j < p.width; j++ {
			y := yPos + i
			x := xPos + j
			distance := math.Sqrt(math.Pow(float64(centerY-y), 2) + math.Pow(float64(centerX-x)/2, 2))
			if distance < float64(radius) {
				win.MoveAddChar(y, x, 'a')
				// continue
				if p.total == 0 {
					win.MoveAddChar(y, x, nc.ACS_BLOCK)
					continue
				}
				top := (y - centerY) * 2
				bottom := (x - centerX)
				rad := math.Atan2(float64(top), float64(bottom))
				ri := 0
				for i, rr := range p.rads {
					if rad <= rr {
						ri = i
						break
					}
				}
				err = win.AttrOn(p.colors[ri])
				if err != nil {
					return err
				}
				win.MoveAddChar(y, x, nc.ACS_BLOCK)
				err = win.AttrOff(p.colors[ri])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Doens't do anything
func (p PieChart) HandleKey(key nc.Key) error {
	return nil
}

// Returns the height of the pie chart
func (p PieChart) Height() int {
	return p.height
}

// Returns the width of the pie chart
func (p PieChart) Width() int {
	return p.width
}

// Set the values of the pie chart
func (p *PieChart) SetValues(values []int) {
	p.total = sumInt(values...)
	if p.total == 0 {
		return
	}
	p.values = make([]int, len(values))
	p.values[0] = values[0]
	for i := 1; i < len(values); i++ {
		p.values[i] = values[i] + p.values[i-1]
	}
	p.rads = make([]float64, 0, len(p.values))
	for _, v := range p.values {
		p.rads = append(p.rads, float64(v)*math.Pi*2/float64(p.total)-math.Pi)
	}
}

// Sets the colors of the pie chart sectors
func (p *PieChart) setColors(colorPairs []string) error {
	if len(p.values) != len(colorPairs) {
		return fmt.Errorf("termui - amount of colors and values has to be the same for PieChart (v: %v, c: %v)", len(p.values), len(colorPairs))
	}
	p.colors = make([]nc.Char, 0, len(colorPairs))
	for _, colorPair := range colorPairs {
		color, err := parseColorPair(colorPair)
		if err != nil {
			return err
		}
		p.colors = append(p.colors, color)
	}
	return nil
}

// A word choice element
type WordChoice struct {
	wct    *WordChoiceTemplate
	data   *UIElementData
	IncKey nc.Key
	DecKey nc.Key
}

// Creates a word choice element
func NewWordChoice(options []string, alignment Alignment, y, x int) (*WordChoice, error) {
	var err error
	result := WordChoice{}
	result.IncKey = KeyRight
	result.DecKey = KeyLeft
	result.data = createUIED(y, x)
	result.wct, err = CreateWordChoiceTemplate(options, alignment)
	return &result, err
}

// Returns the currently selected option
func (w WordChoice) GetSelected() *CCTMessage {
	return w.wct.GetSelected()
}

func (w WordChoice) Draw(win *nc.Window) error {
	return w.wct.Draw(win, w.data.yPos, w.data.xPos, w.data.focused)
}

// Returns the element data of the element
func (w WordChoice) GetElementData() *UIElementData {
	return w.data
}

// Toggles between the options
func (w WordChoice) HandleKey(key nc.Key) error {
	switch key {
	case KeyRight:
		w.wct.FocusNext()
	case KeyLeft:
		w.wct.FocusPrev()
	}
	return nil
}

// Returns 1
func (w WordChoice) Height() int {
	return 1
}

// Returns the length of the longest option + 2
func (w WordChoice) Width() int {
	return w.wct.maxLen
}

// A line edit element
type LineEdit struct {
	let  *LineEditTemplate
	data *UIElementData
}

// Creates a new line edit element
func NewLineEdit(text string, maxLength int, y, x int) (*LineEdit, error) {
	result := LineEdit{}
	result.data = createUIED(y, x)
	result.let = CreateLineEditTemplate(text, maxLength)
	return &result, nil
}

// Sets the text of the element
func (l *LineEdit) SetText(text string) error {
	return l.let.SetText(text)
}

// Returns the entered text
func (l LineEdit) GetText() string {
	return l.let.content
}

// Returns the element data of the element
func (l LineEdit) GetElementData() *UIElementData {
	return l.data
}

// Draws the element
func (l LineEdit) Draw(win *nc.Window) error {
	return l.let.Draw(win, l.data.yPos, l.data.xPos, l.data.focused)
}

// On left/right moves the cursor.
// On letters and some other characters enters them.
// On backspace removes the current character.
func (l LineEdit) HandleKey(key nc.Key) error {
	switch key {
	case KeyLeft:
		l.let.MoveCursorLeft()
	case KeyRight:
		l.let.MoveCursorRight()
	case KeyBackspace:
		l.let.DeleteSelected()
	default:
		l.let.AddCh(rune(key))
	}
	return nil
}

// Returns 1
func (l LineEdit) Height() int {
	return 1
}

// Returns the maxLength
func (l LineEdit) Width() int {
	return l.let.maxLen
}

// A list element
type List struct {
	data          *UIElementData
	lt            *ListTemplate
	bcolor        string
	maxWidth      int
	click         func(choice, cursor int, option *CCTMessage) error
	scrollUpKey   nc.Key
	scrollDownKey nc.Key
	clickKey      nc.Key
}

// Creates a list element
func NewList(options []string, maxDisplayAmount int, optionClick func(choice, cursor int, option *CCTMessage) error, y, x int, borderColorPair string) (*List, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("termui - can;t create List with no options")
	}
	result := List{}
	result.lt = createListTemplate([]*CCTMessage{}, maxDisplayAmount)
	result.data = createUIED(y, x)
	result.bcolor = borderColorPair
	result.click = optionClick
	result.scrollUpKey = '<'
	result.scrollDownKey = '>'
	result.clickKey = KeyEnter
	result.SetOptions(options)
	return &result, nil
}

// Adds an option to the template
func (l *List) AddOption(option string) error {
	cct, err := ToCCTMessage(option)
	if err != nil {
		return err
	}
	l.lt.AddOption(cct)
	return nil
}

// Sets the options of the template
func (l *List) SetOptions(options []string) error {
	cct, err := GetCCTs(options)
	if err != nil {
		return err
	}
	l.lt.SetOptions(cct)
	l.maxWidth = l.lt.options[0].Length()
	for i, o := range l.lt.options {
		if i == 0 {
			continue
		}
		l.maxWidth = maxInt(l.maxWidth, o.Length())
	}
	return nil
}

// Returns the element data of the element
func (l List) GetElementData() *UIElementData {
	return l.data
}

// Draws the scroller
func (l List) drawScroller(win *nc.Window) error {
	if len(l.lt.options) > l.lt.maxDisplayAmount {
		y := l.data.yPos
		x := l.data.xPos
		height := l.Height()
		width := l.Width()
		// draw the arrows
		if l.lt.pageN != 0 {
			win.MoveAddChar(1+y, l.Width()-2+x, nc.ACS_UARROW)
		}
		if l.lt.pageN != len(l.lt.options)-l.lt.maxDisplayAmount {
			win.MoveAddChar(height-2+y, width-2+x, nc.ACS_DARROW)
		}
		// draw the line
		scrollerL := height - 4
		for i := 0; i < scrollerL; i++ {
			win.MoveAddChar(2+y+i, width-2+x, nc.ACS_VLINE)
		}
		// draw the scroller
		sbHeight := l.lt.maxDisplayAmount*scrollerL/len(l.lt.options) + 1
		sbOffset := l.lt.pageN * scrollerL / len(l.lt.options)
		// MessageBox(&Window{win: win}, l.bcolor, []string{}, "normal")
		colorPair := ReverseColorPair(l.bcolor)
		// MessageBox(&Window{win: win}, colorPair, []string{}, "normal")
		color, err := parseColorPair(colorPair)
		if err != nil {
			return err
		}
		win.AttrOn(color)
		for i := 0; i < sbHeight; i++ {
			win.MoveAddChar(2+y+i+sbOffset, width-2+x, ' ')
		}
		win.AttrOff(color)
	}
	return nil
}

// Draws the list
func (l List) Draw(win *nc.Window) error {
	var err error
	DrawBox(win, l.data.yPos, l.data.xPos, l.Height(), l.Width(), l.bcolor)
	err = l.drawScroller(win)
	if err != nil {
		return err
	}
	return l.lt.Draw(win, l.data.yPos+1, l.data.xPos+1, l.data.focused)
}

// On scroll keys scrolls the list
func (l List) HandleKey(key nc.Key) error {
	switch key {
	case l.scrollDownKey:
		l.lt.ScrollDown()
	case l.scrollUpKey:
		l.lt.ScrollUp()
	case l.clickKey:
		if len(l.lt.options) > 0 {
			return l.click(l.lt.choice, l.lt.cursor, l.lt.options[l.lt.choice])
		}
	}
	return nil
}

// Returns the height of the element
func (l List) Height() int {
	return l.lt.maxDisplayAmount + 2
}

// Returns the width of the list
func (l List) Width() int {
	return l.maxWidth + 4
}
