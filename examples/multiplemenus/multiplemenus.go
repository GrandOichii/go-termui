package main

import tui "github.com/GrandOichii/go-termui"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	w, err := tui.CreateWindow("Menu 1")
	checkErr(err)
	firstMenu := w.GetMenu()
	secondMenu, err := tui.NewNormalMenu("${red-cyan}Menu 2")
	checkErr(err)
	b1, err := tui.NewButton("[click me]", 1, 1, func() error {
		w.SetMenu(secondMenu)
		return nil
	}, tui.KeyEnter)
	checkErr(err)
	b2, err := tui.NewButton("${red-cyan}[go back]", 5, 5, func() error {
		w.SetMenu(firstMenu)
		return nil
	}, tui.KeyEnter)
	checkErr(err)
	firstMenu.AddElement(b1)
	secondMenu.AddElement(b2)
	firstMenu.Focus(b1)
	secondMenu.Focus(b2)
	secondMenu.SetBorderColor("red-cyan")
	w.Start()
}
