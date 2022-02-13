package main

import (
	tui "github.com/GrandOichii/go-termui"
	"github.com/rthornton128/goncurses"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	w, err := tui.CreateWindow("ProgressBar tester")
	checkErr(err)
	menu := w.GetMenu()
	// tui.NewProgressBar(menu, 1, 1, 10, 100, true, "normal", "normal")
	pb, err := tui.NewProgressBar(menu, 1, 1, 10, 100, true, "red", "cyan")
	checkErr(err)
	count := 0
	go func() {
		for {
			goncurses.Nap(100)
			count++
			pb.Set(count)
			w.GetMenu().Draw()
			if count == 100 {
				break
			}
		}
	}()
	w.Start()
}
