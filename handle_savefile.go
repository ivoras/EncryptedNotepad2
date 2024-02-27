package main

import (
	"fmt"
	"strings"
	"time"

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

func (ed *EditorWindow) handleSaveFile() {
	fileSave := ed.newSaveFileDialog(ed.handleSaveFileCallback)
	fileSave.Show()
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
			ed.doSaveFile(fwc, passwordEntry1.Text, callback)
		},
		ed.win,
	)
	passwordEntry1.OnSubmitted = func(s string) {
		ed.win.Canvas().Focus(passwordEntry2)
	}
	passwordEntry2.OnSubmitted = func(s string) {
		form.Submit()
	}
	form.Resize(fyne.NewSize(350, 250))
	form.Show()
	time.Sleep(100 * time.Millisecond)
	ed.win.Canvas().Focus(passwordEntry1)
}

func (ed *EditorWindow) doSaveFile(fwc fyne.URIWriteCloser, password string, callBack func()) {
	defer fwc.Close()

	pgpMsg, err := crypto.EncryptMessageWithPassword(crypto.NewPlainMessageFromString(ed.entry.Text), []byte(password))
	if err != nil {
		dialog.ShowError(err, ed.win)
		return
	}

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
	ed.OnCursorChanged()

	if callBack != nil {
		callBack()
	}
}
