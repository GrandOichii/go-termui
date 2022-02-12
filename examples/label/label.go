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
	label, _ := tui.NewLabel("I am a label (Press escape to quit)", 0, 0)
	// add the label to the window
	menu.AddElement(label)
	// start the window
	w.Start()
}
