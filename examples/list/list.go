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
	options := []string{}
	optionI := 0
	add := func() error {
		options = append(options, fmt.Sprintf("value - *%v*", optionI))
		optionI++
		if list != nil {
			list.SetOptions(options)
		}
		return nil
	}
	add()
	// create the list
	list, _ = tui.NewList(options, 10, func(choice, cursor int, option *tui.CCTMessage) error {
		tui.MessageBox(w, option.ToString(), []string{}, "normal")
		return nil
	}, 0, 0, "magenta")
	// create the button
	button, _ := tui.NewButton("[click me]", 12, 2, add, tui.KeyEnter)
	// add the elements
	menu.AddElement(list)
	menu.AddElement(button)
	// link the elements
	tui.Link(button, list)
	// focus on the list
	menu.(*tui.NormalMenu).Focus(list)
	// start the window
	w.Start()
}
