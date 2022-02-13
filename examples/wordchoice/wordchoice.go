package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("WordChoice tester")
	// extract the menu
	menu := w.GetMenu()
	// create the word choice element
	wc, _ := tui.NewWordChoice(menu, 1, 1, []string{"${red}Red ${normal}sus", "option", "${cyan}Blue ${normal}sus :)"}, tui.AlignCenter, "normal")
	// create the button
	b, _ := tui.NewButton(menu, 3, 1, "[press me]", func() error {
		m := wc.GetSelected().ToString()
		tui.MessageBox(w, "You picked "+m, []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	// link the elements
	tui.Link(wc, b)
	// focus on the word choice
	menu.Focus(wc)
	// start the window
	w.Start()
}
