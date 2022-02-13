package main

import (
	"fmt"

	tui "github.com/GrandOichii/go-termui"
)

func main() {
	var list *tui.List
	// create the window
	w, _ := tui.CreateWindow("Window with label")
	// extract the menu
	menu := w.GetMenu()
	// create the options slice
	options := []tui.DrawableAsLine{}
	optionI := 0
	add := func() error {
		m, _ := tui.ToCCTMessage(fmt.Sprintf("value - *%v*", optionI))
		options = append(options, m)
		optionI++
		if list != nil {
			list.SetOptions(options)
		}
		return nil
	}
	// add the first element
	add()
	// create the list
	list, _ = tui.NewList(menu, 0, 0, options, 10, func(choice, cursor int, option tui.DrawableAsLine) error {
		tui.MessageBox(w, option.(*tui.CCTMessage).ToString(), []string{}, "normal")
		return nil
	}, "magenta")
	// create the button
	button, _ := tui.NewButton(menu, 12, 2, "[click me]", add, tui.KeyEnter)
	// link the elements
	tui.Link(button, list)
	// focus on the list
	menu.(*tui.NormalMenu).Focus(list)
	// start the window
	w.Start()
}
