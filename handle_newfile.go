package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func (ed *EditorWindow) clickedNewFile() {
	if ed.isChanged {
		dialog.ShowConfirm("Save document?",
			"There are unsaved changes in the document. Do you wish to save the document?",
			func(saveFile bool) {
				if saveFile {
					if ed.fileName != "" && ed.oldPassword != "" {
						// Just save the file with the existing filename and password
						ed.saveWithExistingFileAndPassword()
						ed.Reset()
					} else {
						// Need to ask for filename and password, then
						fileSave := ed.newSaveFileDialog(ed.saveFileAndReset)
						fileSave.Show()
					}
				} else {
					ed.Reset()
				}
			},
			ed.win)
	} else {
		ed.Reset()
	}
}

func (ed *EditorWindow) saveFileAndReset(fwc fyne.URIWriteCloser, err error) {
	ed.handleSaveFileCallbackGeneric(fwc, err, func() {
		ed.Reset()
	})
}
