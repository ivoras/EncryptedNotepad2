VERSION=0.5
BUILDNO=8

all: dist dist/README.html dist/screenshot.png dist/Icon.png dist/linux_x64 dist/linux_x64/EncryptedNotepad.tgz
	true

clean:
	rm -rf dist

dist:
	mkdir -p dist

dist/linux_x64:
	mkdir -p dist/linux_x64

dist/Icon.png: Icon.png
	cp Icon.png dist/

dist/README.html: README.md
	pandoc --verbose -f markdown -t html5 --standalone --css=pandoc.css --metadata title="" README.md -o dist/README.html
	cp pandoc.css dist/

dist/screenshot.png: screenshot.png
	cp screenshot.png dist/

dist/linux_x64/EncryptedNotepad.tgz: EncryptedNotepad Linux/install_linux.sh Linux/EncryptedNotepad.desktop Icon.png
	mkdir -p build/linux_x64
	cp EncryptedNotepad build/linux_x64/
	cp Linux/install_linux.sh build/linux_x64/
	cp Linux/EncryptedNotepad.desktop build/linux_x64/
	cp Icon.png build/linux_x64/
	cd build/linux_x64 && tar -czvf ../EncryptedNotepad_Linux.tgz * && cd ../..
	mv build/linux_x64/EncryptedNotepad_Linux.tgz dist/
	rm -rf build/linux_x64
