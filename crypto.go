package main

import (
	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// EncryptText encrypts plaintext using OpenPGP symmetric encryption with a password.
// Returns ASCII-armored encrypted text compatible with standard OpenPGP tools.
func EncryptText(plaintext, password string) (string, error) {
	pgp := crypto.PGP()

	encHandle, err := pgp.Encryption().
		Password([]byte(password)).
		Compress().
		New()
	if err != nil {
		return "", err
	}

	pgpMessage, err := encHandle.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}

	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		return "", err
	}

	return string(armored), nil
}

// DecryptText decrypts ASCII-armored OpenPGP encrypted text using a password.
// Returns the plaintext content.
func DecryptText(encryptedText, password string) (string, error) {
	pgp := crypto.PGP()

	decHandle, err := pgp.Decryption().
		Password([]byte(password)).
		New()
	if err != nil {
		return "", err
	}

	pgpMessage, err := crypto.NewPGPMessageFromArmored(encryptedText)
	if err != nil {
		return "", err
	}

	decrypted, err := decHandle.Decrypt(pgpMessage.Bytes(), crypto.Auto)
	if err != nil {
		return "", err
	}

	return string(decrypted.Bytes()), nil
}
