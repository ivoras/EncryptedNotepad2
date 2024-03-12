VERSION=0.2
BUILDNO=4

all: dist dist/Icon.png dist/linux_x64/EncryptedNotepad2.tar.xz dist/windows/EncryptedNotepad2.exe dist/android/EncryptedNotepd2.aab
	true

dist:
	mkdir -p dist

dist/linux_x64:
	mkdir -p dist/linux_x64

dist/linux_x64/EncryptedNotepad2.tar.xz: *.go
	fyne package -os linux -icon Icon.png -tags osusergo,netgo -release && mkdir -p dist/linux_x64 && mv EncryptedNotepad2.tar.xz dist/linux_x64/

dist/windows:
	mkdir -p dist/windows

dist/windows/EncryptedNotepad2.exe: *.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc fyne package -os windows -icon Icon.png -release -appID com.encryptednotepad2 && mkdir -p dist/windows && mv EncryptedNotepad2.exe dist/windows/

dist/android:
	mkdir -p dist/android

dist/android/EncryptedNotepd2.aab: *.go
	rm -f EncryptedNotepad2.aab ; fyne release -os android/arm64 --appID com.encryptednotepad2 --icon Icon.png --appVersion $(VERSION) --appBuild $(BUILDNO) --keyStore en2.keystore --keyName en2 && mkdir -p dist/android && mv EncryptedNotepad2.aab dist/android/

dist/Icon.png: Icon.png
	cp Icon.png dist/

