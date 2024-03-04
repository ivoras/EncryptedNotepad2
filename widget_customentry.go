package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type shortcutableEntry struct {
	widget.Entry
	shortcuts map[desktop.CustomShortcut]func(ks fyne.Shortcut)
}

// NewMultiLineEntry creates a new entry that allows multiple lines
func NewMultiLineShortcutableEntry() *shortcutableEntry {
	e := &shortcutableEntry{
		Entry:     widget.Entry{MultiLine: true, Wrapping: fyne.TextTruncate},
		shortcuts: map[desktop.CustomShortcut]func(ks fyne.Shortcut){},
	}
	e.ExtendBaseWidget(e)
	return e
}

func (se *shortcutableEntry) AddShortcut(s desktop.CustomShortcut, f func(ks fyne.Shortcut)) {
	se.shortcuts[s] = f
}

func (se *shortcutableEntry) TypedShortcut(s fyne.Shortcut) {
	cs, ok := s.(*desktop.CustomShortcut)
	if !ok {
		se.Entry.TypedShortcut(s)
		return
	}
	if handler, found := se.shortcuts[*cs]; found {
		handler(cs)
	}
}
