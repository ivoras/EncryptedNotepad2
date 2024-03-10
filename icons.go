package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icons/new.svg
var iconNewBytes []byte

//go:embed icons/open.svg
var iconOpenBytes []byte

//go:embed icons/save.svg
var iconSaveBytes []byte

//go:embed icons/save-as.svg
var iconSaveAsBytes []byte

//go:embed icons/help.svg
var iconHelpBytes []byte

var iconMap map[string]fyne.Resource = map[string]fyne.Resource{
	"new":     fyne.NewStaticResource("new", iconNewBytes),
	"open":    fyne.NewStaticResource("open", iconOpenBytes),
	"save":    fyne.NewStaticResource("save", iconSaveBytes),
	"save-as": fyne.NewStaticResource("save-as", iconSaveAsBytes),
	"help":    fyne.NewStaticResource("help", iconHelpBytes),
}
