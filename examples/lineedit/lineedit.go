package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("LineEdit tester")
	// extract the menu
	menu := w.GetMenu()
	// create the label
	tui.NewLabel(menu, 1, 1, "Your name:")
	// create the line edit
	lineedit, _ := tui.NewLineEdit(menu, 1, 12, "", 20, "normal")
	// create the button
	button, _ := tui.NewButton(menu, 2, 12, "[click me]", func() error {
		tui.MessageBox(w, "Your name is ${red}"+lineedit.GetText(), []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	// link the elements
	tui.Link(lineedit, button)
	// focus on the button
	menu.Focus(lineedit)
	// start the window
	w.Start()
}
