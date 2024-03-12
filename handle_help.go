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

	label1 := widget.NewLabel(fmt.Sprintf("(C) %d Ivan Voras <ivoras@gmail.com>", time.Now().Year()))
	label1.Wrapping = fyne.TextWrapWord

	label2 := widget.NewRichTextFromMarkdown("Source: [https://github.com/ivoras/EncryptedNotepad2](https://github.com/ivoras/EncryptedNotepad2)")
	//label2 := widget.NewLabel("Original source repo: https://github.com/ivoras/EncryptedNotepad2")
	label2.Wrapping = fyne.TextWrapWord

	infoLines := container.NewVBox(
		label1,
		label2,
	)

	readmeEdit := widget.NewRichTextFromMarkdown(readmeFile)
	readmeEdit.Wrapping = fyne.TextWrapWord
	readmeEdit.Scroll = container.ScrollVerticalOnly

	mainLayout := container.NewBorder(infoLines, nil, nil, nil, readmeEdit)

	dlg := dialog.NewCustom(fmt.Sprintf("About %s v%s", APP_NAME, APP_VERSION), "Ok", mainLayout, ed.win)
	dlg.Resize(fyne.NewSize(640, 480))
	dlg.Show()
}
