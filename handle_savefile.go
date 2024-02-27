package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	_ "github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (ed *EditorWindow) handleSaveFile() {
	fileSave := dialog.NewFileSave(ed.handleSaveFileCallback, ed.win)

	fileSave.Show()
}

func (ed *EditorWindow) handleSaveFileCallback(url fyne.URIWriteCloser, err error) {

}
