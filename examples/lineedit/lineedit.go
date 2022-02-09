package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("LineEdit tester")
	// create the label
	label, _ := tui.NewLabel("Your name:", 1, 1)
	// create the line edit

	lineedit, _ := tui.NewLineEdit("", 20, 1, 12)

	// create the button
	button, _ := tui.NewButton("[click me]", 2, 12, func() error {
		tui.MessageBox(w, "Your name is ${red}"+lineedit.GetText(), []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	// add the elements
	w.AddElement(button)
	w.AddElement(lineedit)
	w.AddElement(label)
	// link the elements
	tui.Link(lineedit, button)
	// focus on the button
	w.Focus(lineedit)
	// start the window
	w.Start()
}