package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("Window with button")
	// create the button
	button, _ := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		tui.Beep()
		return nil
	}, tui.KeyEnter)
	// add the button to the window
	w.AddElement(button)
	// set focus on the button
	w.Focus(button)
	// start the window
	w.Start()
}
