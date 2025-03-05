package ui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Just and example for testing
func RunApp() {
	a := app.New()
	w := a.NewWindow("Mr Verse")

	hello := widget.NewLabel("Hello!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome!")
		}),
	))

	w.ShowAndRun()
}

