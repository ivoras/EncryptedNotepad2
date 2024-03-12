package main

import (
	_ "embed"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	fwex "github.com/matwachich/fynex-widgets"
)

//go:embed README.md
var readmeFile string

func (ed *EditorWindow) handleHelp() {
	versionText := widget.NewLabel(fmt.Sprintf("%s v%s", APP_NAME, APP_VERSION))
	versionText.TextStyle.Bold = true

	label1 := widget.NewLabel(fmt.Sprintf("Copyright %d Ivan Voras <ivoras@gmail.com>", time.Now().Year()))
	label1.Wrapping = fyne.TextWrapBreak
	label2 := widget.NewLabel("Original source repo: https://github.com/ivoras/EncryptedNotepad2")
	label2.Wrapping = fyne.TextWrapBreak

	infoLines := container.NewVBox(
		versionText,
		label1,
		label2,
	)
	readmeEdit := fwex.NewEntryEx(10)
	readmeEdit.SetText(readmeFile)
	readmeEdit.TextStyle.Monospace = true
	readmeEdit.SetMinRowsVisible(10)
	readmeEdit.Wrapping = fyne.TextWrapWord
	readmeEdit.SetReadOnly(true)

	mainLayout := container.NewBorder(infoLines, nil, nil, nil, readmeEdit)

	dlg := dialog.NewCustom(fmt.Sprintf("About %s", APP_NAME), "Ok", mainLayout, ed.win)
	dlg.Show()
}
