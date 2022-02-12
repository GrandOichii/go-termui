package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func main() {
	// create the window
	w, _ := tui.CreateWindow("WordChoice tester")
	// extract the menu
	menu := w.GetMenu()
	wc, _ := tui.NewWordChoice([]string{"${red}Red ${normal}sus", "option", "${cyan}Blue ${normal}sus :)"}, tui.AlignCenter, 1, 1)
	b, _ := tui.NewButton("[press me]", 3, 1, func() error {
		m := wc.GetSelected().ToString()
		tui.MessageBox(w, "You picked "+m, []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	// link the elements
	tui.Link(wc, b)
	// add the elements
	menu.AddElement(wc)
	menu.AddElement(b)
	menu.Focus(wc)
	w.Start()
}
