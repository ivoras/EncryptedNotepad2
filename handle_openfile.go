package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (ed *EditorWindow) clickedOpenFile() {
	fileOpen := dialog.NewFileOpen(ed.handleOpenFileCallback, ed.win)

	fileOpen.SetFilter(storage.NewExtensionFileFilter(recognizedFileExtensions))

	lastDir := ed.app.Preferences().StringWithFallback(PREF_LAST_DIR, "")
	if lastDir != "" {
		fileLister, err := storage.ListerForURI(storage.NewFileURI(lastDir))
		if err != nil {
			fmt.Println(err)
		} else {
			fileOpen.SetLocation(fileLister)
		}
	}

	fileOpen.Show()
}

func (ed *EditorWindow) handleOpenFileCallback(frc fyne.URIReadCloser, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	if frc == nil || frc.URI() == nil {
		return
	}
	// Ask for password
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	form := dialog.NewForm(
		"File password",
		"Ok",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Enter password", Widget: passwordEntry},
		},
		func(b bool) {
			if !b {
				return
			}
			openURI := frc.URI().String()
			dir := openURI[0:strings.LastIndex(openURI, "/")]
			ed.app.Preferences().SetString(PREF_LAST_DIR, strings.TrimPrefix(dir, "file://"))
			ed.doOpenFile(frc, passwordEntry.Text)
		},
		ed.win,
	)
	passwordEntry.OnSubmitted = func(s string) {
		form.Submit()
	}
	form.Resize(fyne.NewSize(350, 170))
	form.Show()
	time.Sleep(100 * time.Millisecond)
	ed.win.Canvas().Focus(passwordEntry)
}

func (ed *EditorWindow) doOpenFile(frc fyne.URIReadCloser, password string) {
	defer frc.Close()
	bytesMsg, err := io.ReadAll(frc)
	if err != nil {
		dialog.ShowError(err, ed.win)
		return
	}
	pgpMsg, err := crypto.NewPGPMessageFromArmored(string(bytesMsg))
	if err != nil {
		dialog.ShowError(err, ed.win)
		return
	}
	msg, err := crypto.DecryptMessageWithPassword(pgpMsg, []byte(password))
	if err != nil {
		dialog.ShowError(err, ed.win)
		return
	}
	fileName := strings.TrimPrefix(frc.URI().String(), "file://")

	if ed.isChanged {
		dialog.ShowConfirm("Save document?",
			"There are unsaved changes. Do you wish to save the document?",
			func(saveFile bool) {
				if saveFile {
					if ed.fileName != "" && ed.oldPassword != "" {
						// Just save the file with the existing filename and password
						ed.saveWithExistingFileAndPassword()
						ed.setEditorFile(fileName, msg.GetString())
						ed.oldPassword = password // new password
					} else {
						// Need to ask for filename and password,
						// then save the old file,
						// then load the new file data into the editor.
						fileSave := ed.newSaveFileDialog(func(fwc fyne.URIWriteCloser, err error) {
							ed.handleSaveFileCallbackGeneric(fwc, err, func() {
								ed.setEditorFile(fileName, msg.GetString())
								ed.oldPassword = password
							})
						})
						fileSave.Show()
					}
				}
			},
			ed.win)
	} else {
		ed.setEditorFile(fileName, msg.GetString())
		ed.oldPassword = password
	}

}

func (ed *EditorWindow) setEditorFile(fileName, text string) {
	ed.Reset()
	ed.entry.SetText(text)
	ed.fileName = fileName
	ed.statusLabel.SetText(fileName)
	ed.isChanged = false
	ed.OnCursorChanged()
}
