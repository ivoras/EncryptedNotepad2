package main

import (
	_ "embed"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

//go:embed ABOUT.md
var readmeFile string

func (ed *EditorWindow) handleHelp() {

	label1 := widget.NewLabel(fmt.Sprintf("Copyright %d Ivan Voras <ivoras@gmail.com>", time.Now().Year()))
	label1.Wrapping = fyne.TextWrapWord

	label2 := widget.NewRichTextFromMarkdown("Original source repo: [https://github.com/ivoras/EncryptedNotepad2](https://github.com/ivoras/EncryptedNotepad2)")
	//label2 := widget.NewLabel("Original source repo: https://github.com/ivoras/EncryptedNotepad2")
	label2.Wrapping = fyne.TextWrapWord

	infoLines := container.NewVBox(
		label1,
		label2,
	)
	/*
		readmeEdit := fwex.NewEntryEx(10)
		readmeEdit.SetText(readmeFile)
		readmeEdit.TextStyle.Monospace = true
		readmeEdit.SetMinRowsVisible(10)
		readmeEdit.Wrapping = fyne.TextWrapWord
		readmeEdit.SetReadOnly(true)
	*/
	readmeEdit := widget.NewRichTextFromMarkdown(readmeFile)
	readmeEdit.Wrapping = fyne.TextWrapWord
	readmeEdit.Scroll = container.ScrollVerticalOnly

	mainLayout := container.NewBorder(infoLines, nil, nil, nil, readmeEdit)

	dlg := dialog.NewCustom(fmt.Sprintf("About %s v%s", APP_NAME, APP_VERSION), "Ok", mainLayout, ed.win)
	dlg.Resize(fyne.NewSize(640, 480))
	dlg.Show()
}
