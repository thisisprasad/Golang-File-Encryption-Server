package GUI

import (
	"fmt"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func LaunchWindow() {
	fmt.Println("Launching window...")

	app := app.New()
	w := app.NewWindow("Hello Window World!")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne GUI!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()
}
