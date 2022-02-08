package main

import (
	tui "github.com/GrandOichii/go-termui"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// create the window
	w, err := tui.CreateWindow("WordChoice tester")
	checkErr(err)
	wc, err := tui.NewWordChoice([]string{"${red}Red ${normal}sus", "option", "${cyan}Blue ${normal}sus :)"}, tui.AlignCenter, 1, 1)
	checkErr(err)
	b, err := tui.NewButton("[press me]", 3, 1, func() error {
		m := wc.GetSelected().ToString()
		tui.MessageBox(w, "You picked "+m, []string{}, "normal")
		return nil
	}, tui.KeyEnter)
	checkErr(err)
	tui.Link(wc, b)
	w.AddElement(wc)
	w.AddElement(b)
	w.Focus(wc)
	w.Start()
}
