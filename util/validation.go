package util

import (
	"regexp"
)

func ValidatePhoneForRegister(in string) (isMatch bool, err error) {
	pattern := `^\+62(?:\d{8,15}|-\d{8,15})$`
	re := regexp.MustCompile(pattern)
	isMatch = re.MatchString(in)
	return
}

func ValidateNameForRegister(in string) (err error) {

	return
}
