package main

import (
	"fmt"

	tui "github.com/GrandOichii/go-termui"
)

func main() {
	var list *tui.List
	// create the window
	w, err := tui.CreateWindow("Window with label")
	if err != nil {
		panic(err)
	}
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
	list, err = tui.NewList(options, 10, func(choice, cursor int, option *tui.CCTMessage) error {
		tui.MessageBox(w, option.ToString(), []string{}, "normal")
		return nil
	}, 0, 0, "magenta")
	if err != nil {
		panic(err)
	}
	// create the button
	button, _ := tui.NewButton("[click me]", 12, 2, add, tui.KeyEnter)
	// add the elements
	w.AddElement(list)
	w.AddElement(button)
	// link the elements
	tui.Link(button, list)
	// focus on the list
	w.Focus(list)
	// start the window
	err = w.Start()
	if err != nil {
		panic(err)
	}
}
