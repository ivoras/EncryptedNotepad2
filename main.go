package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

const (
	AppName    = "Encrypted Notepad 2"
	AppVersion = "1.0.0"
)

// AppState holds the application state
type AppState struct {
	currentFile string
	password    string
	modified    bool
	textWidget  *TextWidget
	statusLabel *TLabelWidget
	leftLabel   *TLabelWidget
}

var app AppState

func main() {
	// Set window title and size
	App.WmTitle(AppName)
	ActivateTheme("azure light")
	App.Configure(Width("80c"), Height("50c"))

	// Create menu bar
	createMenuBar()

	// Create main frame
	mainFrame := TFrame(Padding("2"))
	Grid(mainFrame, Row(0), Column(0), Sticky("nsew"))

	// Create text editor with scrollbar
	createTextEditor(mainFrame)

	// Create status bar
	createStatusBar()

	// Configure grid weights for resizing
	GridColumnConfigure(App, 0, Weight(1))
	GridRowConfigure(App, 0, Weight(1))

	// Set up text change tracking
	setupTextChangeTracking()

	// Update status bar initially
	updateStatusBar()

	// Set up window close handler
	WmProtocol(App, "WM_DELETE_WINDOW", Command(handleExit))

	App.Wait()
}

func createMenuBar() {
	menubar := Menu()

	// File menu
	fileMenu := menubar.Menu()
	fileMenu.AddCommand(Lbl("New"), Underline(0), Accelerator("Ctrl+N"), Command(handleNew))
	fileMenu.AddCommand(Lbl("Open..."), Underline(0), Accelerator("Ctrl+O"), Command(handleOpen))
	fileMenu.AddSeparator()
	fileMenu.AddCommand(Lbl("Save"), Underline(0), Accelerator("Ctrl+S"), Command(handleSave))
	fileMenu.AddCommand(Lbl("Save As..."), Underline(5), Accelerator("Ctrl+Shift+S"), Command(handleSaveAs))
	fileMenu.AddSeparator()
	fileMenu.AddCommand(Lbl("Exit"), Underline(1), Accelerator("Alt+F4"), Command(handleExit))
	menubar.AddCascade(Lbl("File"), Underline(0), Mnu(fileMenu))

	// Edit menu
	editMenu := menubar.Menu()
	editMenu.AddCommand(Lbl("Cut"), Underline(2), Accelerator("Ctrl+X"), Command(handleCut))
	editMenu.AddCommand(Lbl("Copy"), Underline(0), Accelerator("Ctrl+C"), Command(handleCopy))
	editMenu.AddCommand(Lbl("Paste"), Underline(0), Accelerator("Ctrl+V"), Command(handlePaste))
	editMenu.AddSeparator()
	editMenu.AddCommand(Lbl("Select All"), Underline(7), Accelerator("Ctrl+A"), Command(handleSelectAll))
	menubar.AddCascade(Lbl("Edit"), Underline(0), Mnu(editMenu))

	// Help menu
	helpMenu := menubar.Menu()
	helpMenu.AddCommand(Lbl("About..."), Underline(0), Command(handleAbout))
	menubar.AddCascade(Lbl("Help"), Underline(0), Mnu(helpMenu))

	App.Configure(Mnu(menubar))

	// Keyboard bindings
	Bind(App, "<Control-n>", Command(func(e *Event) { handleNew() }))
	Bind(App, "<Control-o>", Command(func(e *Event) { handleOpen() }))
	Bind(App, "<Control-s>", Command(func(e *Event) { handleSave() }))
	Bind(App, "<Control-Shift-S>", Command(func(e *Event) { handleSaveAs() }))
	Bind(App, "<Control-q>", Command(func(e *Event) { handleExit() }))
}

func createTextEditor(parent *TFrameWidget) {
	// Create frame for text and scrollbar
	textFrame := parent.TFrame()
	Grid(textFrame, Row(0), Column(0), Sticky("nsew"))
	GridColumnConfigure(parent, 0, Weight(1))
	GridRowConfigure(parent, 0, Weight(1))

	// Create text widget
	app.textWidget = textFrame.Text(
		Wrap("word"),
		Undo(true),
		Font("TkFixedFont"),
		Width(80),
		Height(25),
	)
	Grid(app.textWidget, Row(0), Column(0), Sticky("nsew"))

	// Create vertical scrollbar connected to text widget
	vscroll := textFrame.TScrollbar(
		Orient("vertical"),
		Command(func(e *Event) { e.Yview(app.textWidget) }),
	)
	Grid(vscroll, Row(0), Column(1), Sticky("ns"))

	// Create horizontal scrollbar connected to text widget
	hscroll := textFrame.TScrollbar(
		Orient("horizontal"),
		Command(func(e *Event) { e.Xview(app.textWidget) }),
	)
	Grid(hscroll, Row(1), Column(0), Sticky("ew"))

	// Configure grid weights for text frame
	GridColumnConfigure(textFrame, 0, Weight(1))
	GridRowConfigure(textFrame, 0, Weight(1))
}

func createStatusBar() {
	statusFrame := TFrame(Padding("2 2 2 2"))
	Grid(statusFrame, Row(1), Column(0), Sticky("ew"))

	// Left side - modified indicator and filename
	app.leftLabel = statusFrame.TLabel(Txt("Ready"))
	Grid(app.leftLabel, Row(0), Column(0), Sticky("w"))

	// Right side - line numbers
	app.statusLabel = statusFrame.TLabel(Txt("Ln 1, Col 1 | Lines: 1"))
	Grid(app.statusLabel, Row(0), Column(1), Sticky("e"))

	GridColumnConfigure(statusFrame, 0, Weight(1))
	GridColumnConfigure(statusFrame, 1, Weight(0))
}

func setupTextChangeTracking() {
	// Bind to key release events to track modifications and update status
	Bind(app.textWidget, "<KeyRelease>", Command(func(e *Event) {
		if !app.modified {
			app.modified = true
			updateWindowTitle()
			updateLeftStatus()
		}
		updateStatusBar()
	}))

	// Bind to mouse clicks to update cursor position in status bar
	Bind(app.textWidget, "<ButtonRelease-1>", Command(func(e *Event) {
		updateStatusBar()
	}))
}

func updateStatusBar() {
	if app.textWidget == nil || app.statusLabel == nil {
		return
	}

	// Get current cursor position
	index := app.textWidget.Index("insert")
	parts := strings.Split(index, ".")
	line := "1"
	col := "1"
	if len(parts) >= 2 {
		line = parts[0]
		// Column is 0-based in Tk, display as 1-based
		colNum := 0
		fmt.Sscanf(parts[1], "%d", &colNum)
		col = fmt.Sprintf("%d", colNum+1)
	}

	// Count total lines
	endIndex := app.textWidget.Index("end-1c")
	endParts := strings.Split(endIndex, ".")
	totalLines := "1"
	if len(endParts) >= 1 {
		totalLines = endParts[0]
	}

	// Update status label
	statusText := fmt.Sprintf("Ln %s, Col %s | Lines: %s", line, col, totalLines)
	app.statusLabel.Configure(Txt(statusText))
}

func updateLeftStatus() {
	if app.leftLabel == nil {
		return
	}

	var status string
	if app.modified {
		status = "● Modified"
	} else {
		status = "Ready"
	}

	if app.currentFile != "" {
		status = status + " - " + filepath.Base(app.currentFile)
	}

	app.leftLabel.Configure(Txt(status))
}

func updateWindowTitle() {
	title := AppName
	if app.currentFile != "" {
		title = filepath.Base(app.currentFile) + " - " + AppName
	}
	if app.modified {
		title = "*" + title
	}
	App.WmTitle(title)
}

// File operations

func handleNew() {
	if app.modified {
		if !confirmDiscard() {
			return
		}
	}

	app.textWidget.Delete("1.0", "end")
	app.currentFile = ""
	app.password = ""
	app.modified = false
	updateWindowTitle()
	updateLeftStatus()
	updateStatusBar()
}

func handleOpen() {
	if app.modified {
		if !confirmDiscard() {
			return
		}
	}

	files := GetOpenFile(
		Filetypes("{{Encrypted Files} {.asc}} {{All Files} {*}}"),
		Title("Open Encrypted File"),
	)
	if len(files) == 0 || files[0] == "" {
		return
	}
	filename := files[0]

	// Read the encrypted file
	data, err := os.ReadFile(filename)
	if err != nil {
		MessageBox(Icon("error"), Title("Error"), Msg(fmt.Sprintf("Failed to read file: %v", err)))
		return
	}

	// Ask for password
	password := askPassword("Enter Password", "Enter the password to decrypt the file:", false)
	if password == "" {
		return
	}

	// Decrypt the content
	plaintext, err := DecryptText(string(data), password)
	if err != nil {
		MessageBox(Icon("error"), Title("Decryption Error"), Msg(fmt.Sprintf("Failed to decrypt file: %v\n\nMake sure you entered the correct password.", err)))
		return
	}

	// Set the text
	app.textWidget.Delete("1.0", "end")
	app.textWidget.Insert("1.0", plaintext)
	app.currentFile = filename
	app.password = password
	app.modified = false
	updateWindowTitle()
	updateLeftStatus()
	updateStatusBar()
}

func handleSave() {
	if app.currentFile == "" {
		handleSaveAs()
		return
	}

	// If we don't have a password yet, ask for one
	if app.password == "" {
		password := askPassword("Set Password", "Enter a password to encrypt the file:", true)
		if password == "" {
			return
		}
		app.password = password
	}

	saveFile(app.currentFile, app.password)
}

func handleSaveAs() {
	filename := GetSaveFile(
		Filetypes("{{Encrypted Files} {.asc}} {{All Files} {*}}"),
		Title("Save Encrypted File As"),
		Defaultextension(".asc"),
	)
	if filename == "" {
		return
	}

	// Ensure .asc extension
	if !strings.HasSuffix(strings.ToLower(filename), ".asc") {
		filename += ".asc"
	}

	// Always ask for a new password on Save As
	password := askPassword("Set Password", "Enter a password to encrypt the file:", true)
	if password == "" {
		return
	}

	app.currentFile = filename
	app.password = password
	saveFile(filename, password)
}

func saveFile(filename, password string) {
	// Get the text content
	contentParts := app.textWidget.Get("1.0", "end-1c")
	content := ""
	if len(contentParts) > 0 {
		content = contentParts[0]
	}

	// Encrypt the content
	encrypted, err := EncryptText(content, password)
	if err != nil {
		MessageBox(Icon("error"), Title("Encryption Error"), Msg(fmt.Sprintf("Failed to encrypt: %v", err)))
		return
	}

	// Write to file
	err = os.WriteFile(filename, []byte(encrypted), 0644)
	if err != nil {
		MessageBox(Icon("error"), Title("Error"), Msg(fmt.Sprintf("Failed to write file: %v", err)))
		return
	}

	app.modified = false
	updateWindowTitle()
	updateLeftStatus()
}

// Password dialog
// Note: Using Text widgets instead of TEntry because TEntry doesn't expose a Get method
// and Variable option doesn't work with TEntry (needs Textvariable)

func askPassword(title, message string, confirm bool) string {
	// Create a toplevel dialog
	dialog := Toplevel()
	dialog.WmTitle(title)
	WmTransient(dialog, App)
	dialog.SetResizable(false, false)

	// Center the dialog relative to main window
	WmGeometry(dialog.Window, "+400+300")

	var result string

	// Main frame
	frame := dialog.TFrame(Padding("20"))
	Grid(frame, Row(0), Column(0), Sticky("nsew"))

	// Message label
	msgLabel := frame.TLabel(Txt(message))
	Grid(msgLabel, Row(0), Column(0), Columnspan(2), Sticky("w"), Pady("0 10"))

	// Note about password visibility
	noteLabel := frame.TLabel(Txt("(Password will be visible while typing)"))
	Grid(noteLabel, Row(1), Column(0), Columnspan(2), Sticky("w"), Pady("0 5"))

	// Password entry using Text widget (single line)
	pwdLabel := frame.TLabel(Txt("Password:"))
	Grid(pwdLabel, Row(2), Column(0), Sticky("e"), Padx("0 10"))
	passwordText := frame.Text(Width(30), Height(1), Wrap("none"))
	Grid(passwordText, Row(2), Column(1), Sticky("w"))

	// Confirm password entry (if needed)
	var confirmText *TextWidget
	if confirm {
		confirmLabel := frame.TLabel(Txt("Confirm Password:"))
		Grid(confirmLabel, Row(3), Column(0), Sticky("e"), Padx("0 10"), Pady("5 0"))
		confirmText = frame.Text(Width(30), Height(1), Wrap("none"))
		Grid(confirmText, Row(3), Column(1), Sticky("w"), Pady("5 0"))
	}

	// Button frame
	btnFrame := frame.TFrame()
	Grid(btnFrame, Row(4), Column(0), Columnspan(2), Pady("15 0"))

	okPressed := false

	onOK := func() {
		pwdParts := passwordText.Get("1.0", "end-1c")
		pwd := ""
		if len(pwdParts) > 0 {
			pwd = pwdParts[0]
		}
		if pwd == "" {
			MessageBox(Icon("warning"), Title("Warning"), Msg("Password cannot be empty."))
			return
		}
		if confirm && confirmText != nil {
			confirmParts := confirmText.Get("1.0", "end-1c")
			confirmPwd := ""
			if len(confirmParts) > 0 {
				confirmPwd = confirmParts[0]
			}
			if pwd != confirmPwd {
				MessageBox(Icon("warning"), Title("Warning"), Msg("Passwords do not match."))
				return
			}
		}
		result = pwd
		okPressed = true
		Destroy(dialog)
	}

	okBtn := btnFrame.TButton(Txt("OK"), Width(10), Command(onOK))
	Grid(okBtn, Row(0), Column(0), Padx("0 5"))

	cancelBtn := btnFrame.TButton(Txt("Cancel"), Width(10), Command(func() {
		Destroy(dialog)
	}))
	Grid(cancelBtn, Row(0), Column(1))

	// Bind Enter key to OK
	Bind(dialog, "<Return>", Command(func(e *Event) {
		onOK()
	}))

	// Bind Escape key to Cancel
	Bind(dialog, "<Escape>", Command(func(e *Event) {
		Destroy(dialog)
	}))

	// Focus the password entry
	Focus(passwordText)

	// Make dialog modal
	Grab(dialog)
	dialog.Wait()
	GrabRelease(dialog)

	if okPressed {
		return result
	}
	return ""
}

// Confirm discard changes dialog

func confirmDiscard() bool {
	response := MessageBox(
		Icon("question"),
		Title("Unsaved Changes"),
		Msg("You have unsaved changes. Do you want to discard them?"),
		Type("yesno"),
	)
	return response == "yes"
}

// Edit operations

func handleCut() {
	app.textWidget.Cut()
	app.modified = true
	updateWindowTitle()
	updateLeftStatus()
}

func handleCopy() {
	app.textWidget.Copy()
}

func handlePaste() {
	app.textWidget.Paste()
	app.modified = true
	updateWindowTitle()
	updateLeftStatus()
}

func handleSelectAll() {
	app.textWidget.SelectAll()
}

// About dialog

func handleAbout() {
	aboutText := fmt.Sprintf(`%s
Version %s

A secure text editor that encrypts files using OpenPGP.

Files are encrypted with AES-256 and stored in
the standard OpenPGP ASCII-armored format (.asc),
compatible with other OpenPGP tools.

© 2024 - Licensed under open source terms.`, AppName, AppVersion)

	MessageBox(
		Icon("info"),
		Title("About "+AppName),
		Msg(aboutText),
	)
}

// Exit handler

func handleExit() {
	if app.modified {
		response := MessageBox(
			Icon("question"),
			Title("Unsaved Changes"),
			Msg("You have unsaved changes. Do you want to save before exiting?"),
			Type("yesnocancel"),
		)
		switch response {
		case "yes":
			handleSave()
			if app.modified {
				// Save was cancelled or failed
				return
			}
		case "cancel":
			return
		}
	}
	Destroy(App)
}
