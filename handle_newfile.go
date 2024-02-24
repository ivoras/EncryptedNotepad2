package main

import "fyne.io/fyne/v2/dialog"

func (win *ENWindow) handleNewFile() {
	if win.isChanged {
		dialog.ShowConfirm("Save document?",
			"There are unsaved changes in the document. Do you wish to save the document?",
			func(b bool) {
				if b {
					win.Reset()
				}
			},
			win.win)
	} else {
		win.Reset()
	}
}
