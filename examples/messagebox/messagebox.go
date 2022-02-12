package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("MessageBox tester")
	// extract the menu
	menu := w.GetMenu()
	// create the button
	button, _ := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		// let the user pick from the two options
		result, _ := tui.MessageBox(w, "Red of blue?", []string{"${red}Red", "${blue}Blue"}, "cyan")
		// show the user the picked option
		tui.MessageBox(w, "You chose "+result, []string{}, "red-white")
		return nil
	}, tui.KeyEnter)
	// add the button to the window
	menu.AddElement(button)
	// focus on the button
	menu.Focus(button)
	// start the window
	w.Start()
}
