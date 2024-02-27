package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func (ed *EditorWindow) handleNewFile() {
	if ed.isChanged {
		dialog.ShowConfirm("Save document?",
			"There are unsaved changes in the document. Do you wish to save the document?",
			func(b bool) {
				if b {
					if ed.fileName == "" {
						fileSave := ed.newSaveFileDialog(ed.saveFileAndReset)
						fileSave.Show()
					} else {
						fwc, err := storage.Writer(storage.NewFileURI(ed.fileName))
						if err != nil {
							fmt.Println("Cannot create Writer on", ed.fileName)
							dialog.ShowError(err, ed.win)
							return
						}
						ed.saveFileAndReset(fwc, nil)
					}
				}
			},
			ed.win)
	} else {
		ed.Reset()
	}
}

func (ed *EditorWindow) saveFileAndReset(fwc fyne.URIWriteCloser, err error) {
	ed.handleSaveFileCallbackGeneric(fwc, nil, func() {
		ed.Reset()
	})
}
