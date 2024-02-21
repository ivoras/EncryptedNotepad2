package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	win := newMainWindow(myApp)

	win.win.ShowAndRun()
}
