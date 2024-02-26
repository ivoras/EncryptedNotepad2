package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (win *ENWindow) handleOpenFile() {
	fileOpen := dialog.NewFileOpen(win.handleOpenFileCallback, win.win)

	fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".asc"}))

	lastDir := win.app.Preferences().StringWithFallback(PREF_LAST_DIR, "")
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

func (win *ENWindow) handleOpenFileCallback(frc fyne.URIReadCloser, err error) {
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
			win.app.Preferences().SetString(PREF_LAST_DIR, strings.TrimPrefix(dir, "file://"))
			win.doOpenFile(openURI, passwordEntry.Text)
		},
		win.win,
	)
	passwordEntry.OnSubmitted = func(s string) {
		form.Submit()
	}
	form.Resize(fyne.NewSize(350, 170))
	form.Show()
	time.Sleep(100 * time.Millisecond)
	win.win.Canvas().Focus(passwordEntry)
}

func (win *ENWindow) doOpenFile(fileName, password string) {
	fileName = strings.TrimPrefix(fileName, "file://")
	//fmt.Println("Opening", fileName)
	f, err := os.Open(fileName)
	if err != nil {
		dialog.ShowError(err, win.win)
		return
	}
	defer f.Close()
	bytesMsg, err := io.ReadAll(f)
	if err != nil {
		dialog.ShowError(err, win.win)
		return
	}
	pgpMsg, err := crypto.NewPGPMessageFromArmored(string(bytesMsg))
	if err != nil {
		dialog.ShowError(err, win.win)
		return
	}
	msg, err := crypto.DecryptMessageWithPassword(pgpMsg, []byte(password))
	if err != nil {
		dialog.ShowError(err, win.win)
		return
	}

	if win.isChanged {
		dialog.ShowConfirm("Save document?",
			"There are unsaved changes in the document. Do you wish to save the document?",
			func(b bool) {
				if b {
					win.Reset()
					win.entry.SetText(msg.GetString())
					win.statusLabel.SetText(fileName)
				}
			},
			win.win)
	} else {
		win.Reset()
		win.entry.SetText(msg.GetString())
		win.statusLabel.SetText(fileName)
	}

}
