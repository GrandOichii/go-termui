package main

import (
	"log"

	tui "github.com/GrandOichii/go-termui"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	w, err := tui.CreateWindow("${red}Multiple ${normal}elements")
	checkErr(err)

	label1, err := tui.NewLabel("Funny button", 1, 1)
	checkErr(err)
	button1, err := tui.NewButton("Click", 2, 1, func() error {
		tui.MessageBox(w, "Funny button clicked", []string{})
		return nil
	}, tui.KeyEnter)
	checkErr(err)
	label2, err := tui.NewLabel("Exit button", 1, 21)
	checkErr(err)
	button2, err := tui.NewButton("Click", 2, 21, func() error {
		w.Exit()
		return nil
	}, tui.KeyEnter)
	checkErr(err)

	tui.Link(button1, button2)

	w.AddElement(label1)
	w.AddElement(button1)
	w.AddElement(label2)
	w.AddElement(button2)

	w.Focus(button1)

	w.Start()
}
