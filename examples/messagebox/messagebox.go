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
	w, err := tui.CreateWindow("MessageBox tester")
	checkErr(err)

	button, err := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		result, _ := tui.MessageBox(w, "Red of blue?", []string{"${red}Red", "${blue}Blue"})
		tui.MessageBox(w, "You chose "+result, []string{})
		return nil
	}, tui.KeyEnter)
	checkErr(err)
	w.Focus(button)
	w.AddElement(button)

	w.Start()
}
