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
	wct    *wordChoiceTemplate
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
	result.wct, err = createWordChoiceTemplate(options, alignment)
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
	let  *lineEditTemplate
	data *UIElementData
}

// Creates a new line edit element
func NewLineEdit(text string, maxLength int, y, x int) (*LineEdit, error) {
	result := LineEdit{}
	result.data = createUIED(y, x)
	result.let = createLineEditTemplate(text, maxLength)
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
