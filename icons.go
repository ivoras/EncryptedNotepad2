package main

import (
	"bytes"
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

var iconMapLight map[string]fyne.Resource = map[string]fyne.Resource{
	"new":     fyne.NewStaticResource("new", iconNewBytes),
	"open":    fyne.NewStaticResource("open", iconOpenBytes),
	"save":    fyne.NewStaticResource("save", iconSaveBytes),
	"save-as": fyne.NewStaticResource("save-as", iconSaveAsBytes),
	"help":    fyne.NewStaticResource("help", iconHelpBytes),
}

var iconMapDark map[string]fyne.Resource = map[string]fyne.Resource{
	"new":     fyne.NewStaticResource("new", lightenSVG(iconNewBytes)),
	"open":    fyne.NewStaticResource("open", lightenSVG(iconOpenBytes)),
	"save":    fyne.NewStaticResource("save", lightenSVG(iconSaveBytes)),
	"save-as": fyne.NewStaticResource("save-as", lightenSVG(iconSaveAsBytes)),
	"help":    fyne.NewStaticResource("help", lightenSVG(iconHelpBytes)),
}

// As a very, very special case, this will replace a single HTML color in a byte string
// with a very light gray.
func lightenSVG(inSVG []byte) (outSVG []byte) {
	outSVG = make([]byte, len(inSVG))
	copy(outSVG, inSVG)
	p := bytes.IndexByte(outSVG, '#')
	if p == -1 {
		return
	}
	outSVG[p+1] = 'd'
	outSVG[p+2] = '0'
	outSVG[p+3] = 'd'
	outSVG[p+4] = '0'
	outSVG[p+5] = 'd'
	outSVG[p+6] = '0'
	return
}
