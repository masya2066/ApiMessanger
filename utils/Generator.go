package utils

import "math/rand"

func GenerateId() string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyz"

	id := ""

	for i := 0; i < 32; i++ {
		id += string(charset[rand.Intn(len(charset))])
		if i == 6 || i == 14 || i == 22 {
			id += "-"
		}
	}

	return id
}

func GenerateVerify() int {
	code := rand.Intn(10000)
	code = code % 10000

	return code
}
