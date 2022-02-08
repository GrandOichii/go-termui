package termui

import (
	"C"

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

type hasElementData interface {
	GetElementData() *UIElementData
}

type Drawable interface {
	hasElementData
	Draw(win *nc.Window) error
}

// A UI element
type UIElement interface {
	Drawable

	HandleKey(key nc.Key) error
	Height() int
	Width() int
}

// Sets the location of the element
func SetYX(element hasElementData, y, x int) {
	data := element.GetElementData()
	data.yPos = y + yOffset
	data.xPos = x + xOffset
}

// Sets the next element for target
func SetNext(target hasElementData, element UIElement) {
	target.GetElementData().next = element
}

// Sets the previous element for target
func SetPrev(target hasElementData, element UIElement) {
	target.GetElementData().prev = element
}

// Sets the key for selecting the next element
func SetNextKey(element hasElementData, key nc.Key) {
	element.GetElementData().nextKey = key
}

// Sets the key for selecting the prev element
func SetPrevKey(element hasElementData, key nc.Key) {
	element.GetElementData().prevKey = key
}

// Toggles the visibility of the element
func ToggleVisibility(element hasElementData, value bool) {
	element.GetElementData().Visible = value
}

// Element data. Describes the location, visibility and several keys of the element
type UIElementData struct {
	yPos, xPos       int
	focused          bool
	Visible          bool
	next, prev       UIElement
	nextKey, prevKey nc.Key
}

// Creates the element data
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

// Standard window
type Window struct {
	Title string

	height      int
	width       int
	focusedElID int
	running     bool
	borderColor string
	cctTitle    *CCTMessage
	elements    []UIElement
	win         *nc.Window
}

// Sets the border color of the window
func (w *Window) SetBorderColor(borderColor string) {
	w.borderColor = borderColor
}

// Returns the height and width of the window
func (w Window) GetMaxYX() (int, int) {
	return w.height, w.width
}

// Draws the window
func (w Window) Draw() error {
	w.win.Erase()
	var err error
	err = DrawBorders(w.win, w.borderColor)
	w.cctTitle.Draw(w.win, 0, 1)
	if err != nil {
		return err
	}
	pWin := w.win
	for _, el := range w.elements {
		elData := el.GetElementData()
		if elData.Visible {
			err = el.Draw(pWin)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Retunrs GetChar result
func (w Window) GetKey() nc.Key {
	return w.win.GetChar()
}

// Returns the goncurses window
func (w Window) GetWin() *nc.Window {
	return w.win
}

// Exits the window
func (w *Window) Exit() {
	w.running = false
	nc.End()
}

// Returns the element that is located at the point
func (w Window) elementAt(y, x int) UIElement {
	for _, el := range w.elements {
		elData := el.GetElementData()
		if y >= elData.yPos && y <= elData.yPos+el.Height() && x >= elData.xPos && x <= elData.xPos+el.Width() {
			return el
		}
	}
	return nil
}

// If esc is pressed, exits the application.
// If mouse is clicked, focuses on the clicked element. If element is already focused, calls the HandleKey method in element.
// Otherwise calls the HandleKey method in the focused element
func (w *Window) HandleKey(key nc.Key) error {
	if key == nc.KEY_ESC {
		w.Exit()
		return nil
	}
	// if key == nc.KEY_MOUSE {
	// 	md := nc.GetMouse()
	// 	element := w.elementAt(md.Y, md.X)
	// 	elData := element.GetElementData()
	// 	if element == nil {
	// 		return nil
	// 	}
	// 	if elData.focused {
	// 		// MessageBox(w, fmt.Sprintf("%v", elData), []string{})
	// 		element.HandleKey(key)
	// 	} else {
	// 		w.Focus(element)
	// 	}
	// 	return nil
	// }
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
			case elData.prevKey:
				// focus on the elData.prev
				if elData.prev == nil {
					continue
				}
				elData.focused = false
				elData.prev.GetElementData().focused = true
			default:
				el.HandleKey(key)
			}
			break
		}
	}
	return nil
}

// Returns the elements of the window
func (w Window) GetElements() []UIElement {
	return w.elements
}

// Basic goncurses configuration
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

// Starts the window
func (w *Window) Start() error {
	var err error
	w.running = true
	if err != nil {
		return err
	}
	w.config()
	defer nc.Cursor(1)
	defer w.Exit()
	var key nc.Key
	for w.running {
		// draw
		err = w.Draw()
		if err != nil {
			return err
		}
		// handle key
		key = w.GetKey()
		err = w.HandleKey(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// Adds the element to the window
func (w *Window) AddElement(element UIElement) {
	w.elements = append(w.elements, element)
}

// Unfocuses all the elements in the window, then focuses the element
func (w *Window) Focus(element hasElementData) {
	w.unfocusAll()
	element.GetElementData().focused = true
}

// Unfocuses all the elements in the window
func (w *Window) unfocusAll() {
	for _, el := range w.elements {
		el.GetElementData().focused = false
	}
}

// Creates new window (should only be called once)
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
	result.borderColor = "normal"
	return &result, nil
}

// Links all the elements
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

// Calls the goncurses Flash method
func Flash() {
	nc.Flash()
}

// Calls the goncurses Beep method
func Beep() {
	nc.Beep()
}
