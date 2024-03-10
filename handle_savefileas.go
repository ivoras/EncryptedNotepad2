package main

func (ed *EditorWindow) clickedSaveFileAs() {
	fileSave := ed.newSaveFileDialog(ed.handleSaveFileCallback)
	fileSave.Show()
}
