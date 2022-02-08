package main

import (
	"fmt"

	tui "github.com/GrandOichii/go-termui"
)

var w *tui.Window

func main() {
	w, _ = tui.CreateWindow("${red}Pie chart test")
	w.SetBorderColor("cyan")

	values := []int{1, 1, 1, 1}

	pcheight := 21
	pcwidth := 41
	pcy := 0
	pcx := 0
	piechart, _ := tui.NewPieChart(w, pcy, pcx, pcheight, pcwidth, values, []string{}, "cyan")
	// tui.MessageBox(w, fmt.Sprintf("*%v*", values), []string{}, "normal")
	buttons := make([]tui.UIElement, len(values))
	buttonText := "[Increase value %v (%v)]"
	for i := 0; i < len(values); i++ {
		val := i
		var b *tui.Button
		b, _ = tui.NewButton(fmt.Sprintf(buttonText, i+1, values[val]), pcy+1+i*2, pcx+pcwidth+1, func() error {
			values[val]++
			piechart.SetValues(values)
			b.SetText(fmt.Sprintf(buttonText, val+1, values[val]))
			return nil
		}, tui.KeyEnter)
		buttons[i] = b
		w.AddElement(b)
	}
	tui.Link(buttons...)
	w.Focus(buttons[0])
	w.AddElement(piechart)

	w.Start()
}
