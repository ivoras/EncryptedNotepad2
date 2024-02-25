package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func (win *ENWindow) handleOpenFile() {
	fileOpen := dialog.NewFileOpen(win.handleOpenFileCallback, win.win)

	fileOpen.SetFilter(storage.NewExtensionFileFilter([]string{".asc"}))

	lastDir := win.app.Preferences().StringWithFallback(PREF_LAST_DIR, "")
	if lastDir != "" {
		fmt.Println("lastDir:", lastDir)
		fileLister, _ := storage.ListerForURI(storage.NewFileURI(lastDir))
		fileOpen.SetLocation(fileLister)
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
			win.app.Preferences().SetString(PREF_LAST_DIR, dir)
			win.doOpenFile(openURI, passwordEntry.Text)
		},
		win.win,
	)
	passwordEntry.OnSubmitted = func(s string) {
		form.Submit()
	}
	form.Resize(fyne.NewSize(350, 170))
	form.Show()
}

func (win *ENWindow) doOpenFile(fileName, password string) {
	fileName = strings.TrimPrefix(fileName, "file://")
	fmt.Println("Opening", fileName)
}
