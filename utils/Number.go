package utils

import "regexp"

func CheckDigits(number string) bool {
	pattern := "^[0-9]+$"
	re := regexp.MustCompile(pattern)
	if re.MatchString(number) {
		return true
	} else {
		return false
	}
}
