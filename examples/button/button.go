package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("Window with button")
	// get the menu of the window
	menu := w.GetMenu()
	// create the button
	button, _ := tui.NewButton(menu, 0, 0, "Press enter to click me!", func() error {
		tui.Beep()
		return nil
	}, tui.KeyEnter)
	// set focus on the button
	menu.Focus(button)
	// start the window
	w.Start()
}
