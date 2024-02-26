package main

import "fyne.io/fyne/v2/dialog"

func (ed *EditorWindow) handleNewFile() {
	if ed.isChanged {
		dialog.ShowConfirm("Save document?",
			"There are unsaved changes in the document. Do you wish to save the document?",
			func(b bool) {
				if b {
					// TODO: Save file
					ed.Reset()
				}
			},
			ed.win)
	} else {
		ed.Reset()
	}
}
