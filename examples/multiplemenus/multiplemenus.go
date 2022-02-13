package main

import tui "github.com/GrandOichii/go-termui"

func main() {
	w, _ := tui.CreateWindow("Menu 1")
	firstMenu := w.GetMenu()
	secondMenu, _ := tui.NewNormalMenu("${red-cyan}Menu 2")
	b1, _ := tui.NewButton(firstMenu, 1, 1, "[click me]", func() error {
		w.SetMenu(secondMenu)
		return nil
	}, tui.KeyEnter)
	b2, _ := tui.NewButton(secondMenu, 5, 5, "${red-cyan}[go back]", func() error {
		w.SetMenu(firstMenu)
		return nil
	}, tui.KeyEnter)
	firstMenu.Focus(b1)
	secondMenu.Focus(b2)
	secondMenu.SetBorderColor("red-cyan")
	w.Start()
}
