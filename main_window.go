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

func newMainWindow(app fyne.App) (win EditorWindow) {
	win.app = app
	win.win = app.NewWindow(fmt.Sprintf("%s v%s", APP_NAME, APP_VERSION))

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentIcon(), win.handleNewFile),
		widget.NewToolbarAction(theme.FolderOpenIcon(), win.handleOpenFile),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), win.handleSaveFile),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), win.handleHelp),
	)

	topToolbar := container.NewHBox(toolbar)

	win.statusLabel = widget.NewLabel("<new document>")
	win.infoLabel = widget.NewLabel("000:000")
	bottomStatus := container.NewHBox(
		win.statusLabel,
		layout.NewSpacer(),
		win.infoLabel,
	)

	win.entry = widget.NewMultiLineEntry()
	win.entry.SetPlaceHolder("Just Because You're Paranoid Doesn't Mean They're Not After You")
	win.entry.OnCursorChanged = win.OnCursorChanged
	win.entry.OnChanged = win.OnChanged
	middleContent := container.NewMax(win.entry)

	topLayout := container.NewBorder(topToolbar, bottomStatus, nil, nil, middleContent)

	win.win.Resize(fyne.NewSize(800, 600))
	win.win.SetContent(topLayout)

	win.Reset()
	win.win.CenterOnScreen()

	win.win.Canvas().Focus(win.entry)

	return
}

func (ed *EditorWindow) Reset() {
	ed.statusLabel.SetText("<new document>")
	ed.infoLabel.SetText("000:000")
	ed.entry.SetText("")
	ed.fileName = ""
	ed.isChanged = false
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
