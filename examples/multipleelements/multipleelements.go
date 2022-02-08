package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("${red}Multiple ${normal}elements")
	// create the first label
	label1, _ := tui.NewLabel("Funny button", 1, 1)
	// create the first button
	button1, _ := tui.NewButton("Click", 2, 1, func() error {
		tui.MessageBox(w, "Funny button clicked", []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	tui.SetNextKey(button1, tui.KeyRight)
	tui.SetPrevKey(button1, tui.KeyLeft)
	// create the second label
	label2, _ := tui.NewLabel("Exit button", 1, 21)
	// create the second button
	button2, _ := tui.NewButton("Click", 2, 21, func() error {
		w.Exit()
		return nil
	}, tui.KeyEnter)
	tui.SetNextKey(button2, tui.KeyRight)
	tui.SetPrevKey(button2, tui.KeyLeft)

	// add all the elements
	w.AddElement(label1)
	w.AddElement(button1)
	w.AddElement(label2)
	w.AddElement(button2)

	// link up all the buttons
	tui.Link(button1, button2)

	// focus on the first button
	w.Focus(button1)

	// start the window
	w.Start()
}
