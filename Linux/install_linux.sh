#!/bin/sh -e

mkdir -p ~/bin
cp EncryptedNotepad ~/bin/
echo "Encrypted Notepad installed to ~/bin"
chmod a+x ~/bin/EncryptedNotepad

mkdir -p ~/.local/share/applications
cp EncryptedNotepad.desktop ~/.local/share/applications/
echo "Desktop entry installed to ~/.local/share/applications/EncryptedNotepad.desktop"

cp Icon.png ~/.local/share/icons/EncryptedNotepad_icon.png
echo "Icon installed to ~/.local/share/icons/EncryptedNotepad_icon.png"

echo "Done"
