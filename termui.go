package termui

import (
	"C"
	"fmt"

	nc "github.com/rthornton128/goncurses"
)

const (
	yOffset       = 1
	xOffset       = 1
	hightlightKey = nc.A_REVERSE

	KeyEnter = 10
	KeyLeft  = nc.KEY_LEFT
	KeyRight = nc.KEY_RIGHT
)

type Drawable interface {
	Draw(win *nc.Window, yPos, xPos int, focused bool) error
}

type hasElementData interface {
	GetElementData() *UIElementData
}

type UIElement interface {
	hasElementData
	Drawable

	HandleKey(key nc.Key) error
	Height() int
	Width() int
}

func SetYX(element hasElementData, y, x int) {
	data := element.GetElementData()
	data.yPos = y + yOffset
	data.xPos = x + xOffset
}

func SetNext(target hasElementData, element UIElement) {
	target.GetElementData().next = element
}

func SetPrev(target hasElementData, element UIElement) {
	target.GetElementData().prev = element
}

type UIElementData struct {
	yPos, xPos       int
	focused          bool
	Visible          bool
	next, prev       UIElement
	nextKey, prevKey nc.Key
}

func createUIED(y, x int) *UIElementData {
	result := UIElementData{}
	result.yPos = y + yOffset
	result.xPos = x + xOffset
	result.prev = nil
	result.next = nil
	result.prevKey = nc.KEY_UP
	result.nextKey = nc.KEY_DOWN
	result.Visible = true
	return &result
}

type Label struct {
	data    *UIElementData
	cctText *CCTMessage
}

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

func (l Label) Draw(win *nc.Window, yPos, xPos int, focused bool) error {
	attr := nc.A_NORMAL
	if focused {
		attr = hightlightKey
	}
	l.cctText.Draw(win, yPos, xPos, attr)
	// put(pWin, data.yPos, data.xPos, l.text, attr)
	return nil
}

func (l Label) HandleKey(key nc.Key) error {
	return nil
}

func (l *Label) SetText(text string) error {
	var err error
	l.cctText, err = ToCCTMessage(text)
	return err
}

func (l Label) GetElementData() *UIElementData {
	return l.data
}

func (l Label) Height() int {
	return 1
}

func (l Label) Width() int {
	return l.cctText.Length()
}

type Button struct {
	click    func() error
	clickKey nc.Key
	data     *UIElementData
	cctText  *CCTMessage
}

func NewButton(text string, y, x int, click func() error, clickKey nc.Key) (*Button, error) {
	result := Button{}
	var err error
	result.cctText, err = ToCCTMessage(text)
	if err != nil {
		return nil, err
	}
	result.click = click
	result.clickKey = clickKey
	result.data = createUIED(y, x)
	return &result, nil
}

func (b Button) Draw(win *nc.Window, yPos, xPos int, focused bool) error {
	attr := nc.A_NORMAL
	if focused {
		attr = hightlightKey
	}
	b.cctText.Draw(win, yPos, xPos, attr)
	// put(pWin, data.yPos, data.xPos, b.text, attr)
	return nil
}

func (b *Button) SetText(text string) error {
	var err error
	b.cctText, err = ToCCTMessage(text)
	return err
}

func (b Button) HandleKey(key nc.Key) error {
	if key == b.clickKey || key == nc.KEY_MOUSE {
		return b.click()
	}
	return nil
}

func (b Button) GetElementData() *UIElementData {
	return b.data
}

func (b Button) Height() int {
	return 1
}

func (b Button) Width() int {
	return b.cctText.Length()
}

type Window struct {
	Title string

	height      int
	width       int
	focusedElID int
	running     bool
	cctTitle    *CCTMessage
	elements    []UIElement
	win         *nc.Window
}

func (w Window) GetMaxYX() (int, int) {
	return w.height, w.width
}

func (w Window) Draw() error {
	w.win.Erase()
	var err error
	err = DrawBorders(w.win)
	w.cctTitle.Draw(w.win, 0, 1)
	if err != nil {
		return err
	}
	pWin := w.win
	for _, el := range w.elements {
		elData := el.GetElementData()
		if elData.Visible {
			err = el.Draw(pWin, elData.yPos, elData.xPos, elData.focused)
			if err != nil {
				return err
			}
		}
	}
	w.win.Refresh()
	return nil
}

func (w Window) GetKey() nc.Key {
	return w.win.GetChar()
}

func (w Window) GetWin() *nc.Window {
	return w.win
}

func (w *Window) Exit() {
	w.running = false
	nc.End()
}

func (w Window) elementAt(y, x int) UIElement {
	for _, el := range w.elements {
		elData := el.GetElementData()
		if y >= elData.yPos && y <= elData.yPos+el.Height() && x >= elData.xPos && x <= elData.xPos+el.Width() {
			return el
		}
	}
	return nil
}

func (w *Window) HandleKey(key nc.Key) error {
	if key == nc.KEY_ESC {
		w.Exit()
		return nil
	}
	if key == nc.KEY_MOUSE {
		md := nc.GetMouse()
		element := w.elementAt(md.Y, md.X)
		elData := element.GetElementData()
		if element == nil {
			return nil
		}
		if elData.focused {
			// MessageBox(w, fmt.Sprintf("%v", elData), []string{})
			element.HandleKey(key)
		} else {
			w.Focus(element)
		}
		return nil
	}
	for _, el := range w.elements {
		elData := el.GetElementData()
		if elData.focused {
			switch key {
			case elData.nextKey:
				// focus on the elData.next
				if elData.next == nil {
					continue
				}
				elData.focused = false
				elData.next.GetElementData().focused = true
				_ = 1
			case elData.prevKey:
				// focus on the elData.prev
				if elData.prev == nil {
					continue
				}
				elData.focused = false
				elData.prev.GetElementData().focused = true
				_ = 1
			default:
				el.HandleKey(key)
			}
			break
		}
	}
	return nil
}

func (w Window) GetElements() []UIElement {
	return w.elements
}

func (w *Window) config() {
	// remove the delay from pressing the escape key
	nc.SetEscDelay(0)
	w.win.Keypad(true)

	nc.Raw(true)
	nc.Echo(false)
	nc.Cursor(0)
	nc.CBreak(true)
	nc.MouseInterval(50)

	nc.MouseMask(nc.M_B1_PRESSED, nil) // only detect left mouse clicks
}

func (w *Window) Start() error {
	var err error
	w.running = true
	if err != nil {
		return err
	}
	w.config()
	defer nc.Cursor(1)
	var key nc.Key
	for w.running {
		// draw
		w.Draw()
		// handle key
		key = w.GetKey()
		w.HandleKey(key)
		// clear screen
	}
	return nil
}

func (w *Window) AddElement(element UIElement) {
	w.elements = append(w.elements, element)
}

func (w *Window) PrintData() {
	for _, el := range w.elements {
		fmt.Println(el, el.GetElementData())
	}
}

func (w *Window) Focus(element hasElementData) {
	w.unfocusAll()
	element.GetElementData().focused = true
}

func (w *Window) unfocusAll() {
	for _, el := range w.elements {
		el.GetElementData().focused = false
	}
}

func CreateWindow(title string) (*Window, error) {
	var err error
	result := Window{}
	result.win, err = nc.Init()
	if err != nil {
		return nil, err
	}
	initColors()
	result.focusedElID = 0
	result.Title = title
	result.cctTitle, err = ToCCTMessage(title)
	if err != nil {
		return nil, err
	}
	result.running = false
	result.elements = []UIElement{}
	return &result, nil
}

func Link(elements ...UIElement) {
	if len(elements) == 0 {
		return
	}
	lastElI := len(elements) - 1
	firstEl := elements[0]
	lastEl := elements[lastElI]
	SetNext(lastEl, firstEl)
	SetPrev(firstEl, lastEl)
	for i, el := range elements {
		data := el.GetElementData()
		if i != 0 {
			data.prev = elements[i-1]
		}
		if i != lastElI {
			data.next = elements[i+1]
		}
	}
}

func Flash() {
	nc.Flash()
}

func Beep() {
	nc.Beep()
}
