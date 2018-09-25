package cli

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func encrypt(plain []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedPlain := padPKCS7(plain)
	cipherText := make([]byte, aes.BlockSize+len(paddedPlain))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(cipherText[aes.BlockSize:], paddedPlain)

	sEnc := base64.StdEncoding.EncodeToString(cipherText)
	return sEnc, nil
}

func decrypt(b64EncodedCipherText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(b64EncodedCipherText)
	if err != nil {
		return "", err
	}

	plain := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, cipherText[:aes.BlockSize])
	decrypter.CryptBlocks(plain, cipherText[aes.BlockSize:])
	padSize := int(plain[len(plain)-1])

	return string(plain[:len(plain)-padSize]), nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func createEncryptionKeyFile(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	key, err := generateKey()
	if err != nil {
		return err
	}
	b64key := base64.StdEncoding.EncodeToString(key)
	err = ioutil.WriteFile(path, []byte(b64key), 0600)
	return err
}

func padPKCS7(data []byte) []byte {
	padSize := 0
	if len(data)%aes.BlockSize != 0 {
		padSize = aes.BlockSize - len(data)%aes.BlockSize
	}
	appendChars := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(data, appendChars...)
}
