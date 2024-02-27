package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type EditorWindow struct {
	fileName    string
	isChanged   bool
	app         fyne.App
	win         fyne.Window
	statusLabel *widget.Label
	infoLabel   *widget.Label
	entry       *widget.Entry
}

func newMainWindow(app fyne.App) (ed EditorWindow) {
	ed.app = app
	ed.win = app.NewWindow(fmt.Sprintf("%s v%s", APP_NAME, APP_VERSION))

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentIcon(), ed.handleNewFile),
		widget.NewToolbarAction(theme.FolderOpenIcon(), ed.handleOpenFile),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), ed.handleSaveFile),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), ed.handleHelp),
	)

	topToolbar := container.NewHBox(toolbar)

	ed.statusLabel = widget.NewLabel("<new document>")
	ed.infoLabel = widget.NewLabel("000:000")
	bottomStatus := container.NewHBox(
		ed.statusLabel,
		layout.NewSpacer(),
		ed.infoLabel,
	)

	ed.entry = widget.NewMultiLineEntry()
	ed.entry.SetPlaceHolder("Just Because You're Paranoid Doesn't Mean They're Not After You")
	ed.entry.OnCursorChanged = ed.OnCursorChanged
	ed.entry.OnChanged = ed.OnChanged
	middleContent := container.NewMax(ed.entry)

	topLayout := container.NewBorder(topToolbar, bottomStatus, nil, nil, middleContent)

	ed.win.Resize(fyne.NewSize(800, 600))
	ed.win.SetContent(topLayout)

	ed.Reset()
	ed.win.CenterOnScreen()

	return
}

func (ed *EditorWindow) Reset() {
	ed.statusLabel.SetText("<new document>")
	ed.infoLabel.SetText("000:000")
	ed.entry.SetText("")
	ed.fileName = ""
	ed.isChanged = false
	ed.win.Canvas().Focus(ed.entry)
	ed.OnCursorChanged()
}

func (ed *EditorWindow) OnCursorChanged() {
	changeMark := ""
	if ed.isChanged {
		changeMark = "*"
	}
	ed.infoLabel.SetText(fmt.Sprintf("%s %03d:%03d", changeMark, ed.entry.CursorColumn+1, ed.entry.CursorRow+1))
}

func (ed *EditorWindow) OnChanged(s string) {
	ed.isChanged = true
}
