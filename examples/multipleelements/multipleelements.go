package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("${red}Multiple ${normal}elements")
	// extract the menu
	menu := w.GetMenu()
	// create the first label
	tui.NewLabel(menu, 1, 1, "Funny button")
	// create the first button
	button1, _ := tui.NewButton(menu, 2, 1, "Click", func() error {
		tui.MessageBox(w, "Funny button clicked", []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	tui.SetNextKey(button1, tui.KeyRight)
	tui.SetPrevKey(button1, tui.KeyLeft)
	// create the second label
	tui.NewLabel(menu, 1, 21, "Exit button")
	// create the second button
	button2, _ := tui.NewButton(menu, 2, 21, "Click", func() error {
		w.Exit()
		return nil
	}, tui.KeyEnter)
	tui.SetNextKey(button2, tui.KeyRight)
	tui.SetPrevKey(button2, tui.KeyLeft)
	// link up all the buttons
	tui.Link(button1, button2)
	// focus on the first button
	menu.Focus(button1)
	// start the window
	w.Start()
}
