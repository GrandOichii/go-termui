package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("EnterString tester")
	// extract the menu
	menu := w.GetMenu()
	// create the button
	button, _ := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		result, _ := tui.EnterString(w, "", "Enter your ${cyan}name", 20, "23")
		tui.MessageBox(w, result, []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	// add the button to the window
	menu.AddElement(button)
	// focus on the button
	menu.Focus(button)
	// start the window
	w.Start()
}
