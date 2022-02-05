package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("EnterString tester")
	// create the button
	button, _ := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		result, _ := tui.EnterString(w, "", "Enter your ${cyan}name", 20)
		tui.MessageBox(w, result, []string{})
		return nil
	}, tui.KeyEnter)
	// add the button to the window
	w.AddElement(button)
	// focus on the button
	w.Focus(button)
	// start the window
	w.Start()
}
