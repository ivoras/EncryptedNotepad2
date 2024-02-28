
all: dist dist/Icon.png dist/linux dist/linux/EncryptedNotepad2.tar.xz dist/windows dist/windows/EncryptedNotepad2.exe
	true

dist:
	mkdir -p dist

dist/linux:
	mkdir -p dist/linux

dist/linux/EncryptedNotepad2.tar.xz:
	fyne package -os linux -icon Icon.png -tags osusergo,netgo -release && mv EncryptedNotepad2.tar.xz dist/linux/

dist/windows:
	mkdir -p dist/windows

dist/windows/EncryptedNotepad2.exe:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc fyne package -os windows -icon Icon.png -release 

dist/Icon.png: Icon.png
	cp Icon.png dist/


