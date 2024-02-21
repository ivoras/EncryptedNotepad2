package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func (win *ENWindow) handleOpenFile() {
	dialog.ShowFileOpen(win.handleOpenFileCallback, win.win)
}

func (win *ENWindow) handleOpenFileCallback(frc fyne.URIReadCloser, err error) {

}
