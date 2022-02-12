package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("DropDownBox tester")
	// extract the menu
	menu := w.GetMenu()
	// create the button
	button, _ := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		ddbOptions := []string{"Choose ${green}me", "or ${red}me", "${magenta-white}proably ${normal}me"}
		result, _ := tui.DropDownBox(ddbOptions, 2, 1, 25, tui.SingleElement, "cyan-gray")
		if len(result) == 0 {
			// user didn't choose anything
			return nil
		}
		// display the choice
		tui.MessageBox(w, ddbOptions[result[0]], []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	// add the button to the window
	menu.AddElement(button)
	// focus on the button
	menu.Focus(button)
	// start the window
	w.Start()
}
