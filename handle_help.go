package main

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (win *ENWindow) handleHelp() {
	dlg := dialog.NewCustom("About", "Ok", widget.NewLabel("Lorem Ipsum"), win.win)
	dlg.Show()
}
