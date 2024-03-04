package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (ed *EditorWindow) newSaveFileDialog(callBack func(fyne.URIWriteCloser, error)) (fileSave *dialog.FileDialog) {
	fileSave = dialog.NewFileSave(callBack, ed.win)

	fileSave.SetFilter(storage.NewExtensionFileFilter(recognizedFileExtensions))

	lastDir := ed.app.Preferences().StringWithFallback(PREF_LAST_DIR, "")
	if lastDir != "" {
		fileLister, err := storage.ListerForURI(storage.NewFileURI(lastDir))
		if err != nil {
			fmt.Println(err)
		} else {
			fileSave.SetLocation(fileLister)
		}
	}

	return
}

func (ed *EditorWindow) clickedSaveFile() {
	if ed.oldPassword == "" || ed.fileName == "" {
		fileSave := ed.newSaveFileDialog(ed.handleSaveFileCallback)
		fileSave.Show()
	} else {
		ed.saveWithExistingFileAndPassword()
	}
}

func (ed *EditorWindow) saveWithExistingFileAndPassword() {
	fwc, err := storage.Writer(storage.NewFileURI(ed.fileName))
	if err != nil {
		fmt.Println("Cannot create Writer on", ed.fileName)
		dialog.ShowError(err, ed.win)
		return
	}
	ed.saveEditorToWriterWithPassword(fwc, ed.oldPassword, nil)
}

func (ed *EditorWindow) handleSaveFileCallback(fwc fyne.URIWriteCloser, err error) {
	ed.handleSaveFileCallbackGeneric(fwc, err, nil)
}

func (ed *EditorWindow) handleSaveFileCallbackGeneric(fwc fyne.URIWriteCloser, err error, callback func()) {
	if err != nil {
		fmt.Println(err)
		dialog.ShowError(err, ed.win)
		return
	}
	if fwc == nil || fwc.URI() == nil {
		return
	}

	// Ask for password
	passwordEntry1 := widget.NewPasswordEntry()
	passwordEntry1.SetPlaceHolder("Password")
	passwordEntry2 := widget.NewPasswordEntry()
	passwordEntry2.SetPlaceHolder("Confirm password")

	form := dialog.NewForm(
		"File password",
		"Ok",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Enter password", Widget: passwordEntry1},
			{Text: "Confirm password", Widget: passwordEntry2},
		},
		func(b bool) {
			if !b {
				return
			}
			if passwordEntry1.Text != passwordEntry2.Text {
				dialog.ShowError(fmt.Errorf("The passwords do not match!"), ed.win)
				return
			}
			saveURI := fwc.URI().String()
			dir := saveURI[0:strings.LastIndex(saveURI, "/")]
			ed.app.Preferences().SetString(PREF_LAST_DIR, strings.TrimPrefix(dir, "file://"))
			ed.saveEditorToWriterWithPassword(fwc, passwordEntry1.Text, callback)
		},
		ed.win,
	)
	passwordEntry1.OnSubmitted = func(s string) {
		ed.win.Canvas().Focus(passwordEntry2)
	}
	passwordEntry2.OnSubmitted = func(s string) {
		form.Submit()
	}
	form.Resize(fyne.NewSize(350, 210))
	form.Show()
	ed.win.Canvas().Focus(passwordEntry1)
}

func (ed *EditorWindow) saveEditorToWriterWithPassword(fwc fyne.URIWriteCloser, password string, callBack func()) {
	defer fwc.Close()

	pgpMsg, err := crypto.EncryptMessageWithPassword(crypto.NewPlainMessageFromString(ed.entry.Text), []byte(password))
	if err != nil {
		dialog.ShowError(err, ed.win)
		return
	}
	ed.oldPassword = password

	aMsg, err := pgpMsg.GetArmored()
	if err != nil {
		dialog.ShowError(err, ed.win)
		return
	}

	_, err = fwc.Write([]byte(aMsg))
	if err != nil {
		dialog.ShowError(err, ed.win)
		return
	}

	ed.isChanged = false
	ed.fileName = strings.TrimPrefix(fwc.URI().String(), "file://")
	ed.statusLabel.SetText(ed.fileName)
	ed.OnCursorChanged()

	if callBack != nil {
		callBack()
	}
}
