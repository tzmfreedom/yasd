package cli

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
)

func encrypt(plain []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	cipherText := make([]byte, aes.BlockSize+len(plain))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(cipherText[aes.BlockSize:], plain)

	sEnc := base64.StdEncoding.EncodeToString(cipherText)
	return sEnc, nil
}

func decrypt(b64EncodedCipherText string, key []byte) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(b64EncodedCipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCBCDecrypter(block, cipherText[:aes.BlockSize])
	decrypter.CryptBlocks(cipherText[aes.BlockSize:], cipherText)

	return string(cipherText[aes.BlockSize:]), nil
}

func encryptCredential(cfg *config) (string, error) {
	file, err := os.Open(cfg.EncyptionKeyPath)
	if err != nil {
		return "", err
	}
	key := make([]byte, 32)
	if _, err = file.Read(key); err != nil {
		return "", err
	}
	encryptedPassword, err := encrypt([]byte(cfg.Password), key)
	if err != nil {
		return "", err
	}
	return encryptedPassword, nil
}

func decryptCredential(cfg *config) (string, error) {
	file, err := os.Open(cfg.EncyptionKeyPath)
	if err != nil {
		return "", err
	}
	key := make([]byte, 32)
	_, err = file.Read(key)
	if err != nil {
		return "", err
	}
	password, err := decrypt(cfg.Password, key)
	if err != nil {
		return "", err
	}
	return password, nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func generateEncryptionKey(path string) error {
	key, err := generateKey()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, key, 0600)
	return err
}
