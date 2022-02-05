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
	w, err := tui.CreateWindow("Window with label")
	checkErr(err)

	label, err := tui.NewLabel("I am a label (Press escape to quit)", 0, 0)
	checkErr(err)
	w.AddElement(label)

	w.Start()
}
