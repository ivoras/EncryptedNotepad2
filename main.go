package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()

	win := createMainWindow(myApp)

	win.ShowAndRun()
}
