package utils

import "math/rand"

func GenerateRandomSecretPhrase() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	secretPhrase := make([]byte, 32)
	for i := range secretPhrase {
		secretPhrase[i] = charset[rand.Intn(len(charset))]
	}

	return string(secretPhrase)
}
