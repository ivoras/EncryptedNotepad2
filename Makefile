
all: dist dist/Icon.png dist/linux_x64 dist/linux_x64/EncryptedNotepad2.tar.xz dist/windows dist/windows/EncryptedNotepad2.exe
	true

dist:
	mkdir -p dist

dist/linux_x64:
	mkdir -p dist/linux_x64

dist/linux_x64/EncryptedNotepad2.tar.xz: *.go
	fyne package -os linux -icon Icon.png -tags osusergo,netgo -release && mv EncryptedNotepad2.tar.xz dist/linux_x64/

dist/windows:
	mkdir -p dist/windows

dist/windows/EncryptedNotepad2.exe: *.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc fyne package -os windows -icon Icon.png -release -appID com.encryptednotepad2  && mv EncryptedNotepad2.exe dist/windows/

dist/Icon.png: Icon.png
	cp Icon.png dist/


