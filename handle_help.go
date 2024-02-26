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

	infoLines := container.NewVBox(
		versionText,
		widget.NewLabel(fmt.Sprintf("Copyright 2024-%d Ivan Voras <ivoras@gmail.com>", time.Now().Year())),
		widget.NewLabel(fmt.Sprintf("Original source repo: https://github.com/ivoras/EncryptedNotepad2")),
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
