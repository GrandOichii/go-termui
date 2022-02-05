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
	w, err := tui.CreateWindow("Window with button")
	checkErr(err)

	button, err := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		tui.Beep()
		return nil
	}, tui.KeyEnter)
	checkErr(err)
	w.Focus(button)
	w.AddElement(button)

	w.Start()
}
