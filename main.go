package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.encryptednotepad2")

	win := newMainWindow(myApp)

	win.win.ShowAndRun()
}
