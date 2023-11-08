package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func EncryptMessage(message string, key string) (string, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	if len([]byte("9137650284736294")) != aes.BlockSize {
		return "", fmt.Errorf("IV length must be 16 bytes")
	}

	stream := cipher.NewCFBEncrypter(block, []byte("9137650284736294"))
	ciphertext := make([]byte, len(message))
	stream.XORKeyStream(ciphertext, []byte(message))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptMessage(message string, key string) (string, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	if len([]byte("9137650284736294")) != aes.BlockSize {
		return "", fmt.Errorf("IV length must be 16 bytes")
	}

	stream := cipher.NewCFBDecrypter(block, []byte("9137650284736294"))
	ciphertext, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	decryptedText := make([]byte, len(ciphertext))
	stream.XORKeyStream(decryptedText, ciphertext)

	return string(decryptedText), nil
}
