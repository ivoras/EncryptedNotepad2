package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icons/new.png
var iconNewBytes []byte

//go:embed icons/open.png
var iconOpenBytes []byte

//go:embed icons/save.png
var iconSaveBytes []byte

//go:embed icons/save-as.png
var iconSaveAsBytes []byte

//go:embed icons/help.png
var iconHelpBytes []byte

var iconMap map[string]fyne.Resource = map[string]fyne.Resource{
	"new":     fyne.NewStaticResource("new", iconNewBytes),
	"open":    fyne.NewStaticResource("open", iconOpenBytes),
	"save":    fyne.NewStaticResource("save", iconSaveBytes),
	"save-as": fyne.NewStaticResource("save-as", iconSaveAsBytes),
	"help":    fyne.NewStaticResource("help", iconHelpBytes),
}
