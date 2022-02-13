package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("Window with label")
	// extract the menu
	menu := w.GetMenu()
	// create the label
	tui.NewLabel(menu, 0, 0, "I am a label (Press escape to quit)")
	// start the window
	w.Start()
}
