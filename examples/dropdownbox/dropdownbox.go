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
	w, err := tui.CreateWindow("DropDownBox tester")
	checkErr(err)

	button, err := tui.NewButton("Press enter to click me!", 0, 0, func() error {
		ddbOptions := []string{"Choose ${green}me", "or ${red}me", "${magenta-white}probably ${normal}me"}
		result, _ := tui.DropDownBox(ddbOptions, 5, 1, 25, tui.SingleElement)
		if len(result) == 0 {
			// user didn't choose anything
			return nil
		}
		tui.MessageBox(w, ddbOptions[result[0]], []string{})
		return nil
	}, tui.KeyEnter)
	checkErr(err)
	w.Focus(button)
	w.AddElement(button)

	w.Start()
}
