package util

import (
	"regexp"
)

func ValidatePhoneFormat(in string) (isMatch bool) {
	//pattern := `^\+62(?:\d{8,15}|-\d{8,15})$` // ini untuk limit 8-16 karakter
	pattern := `^\+62(?:\d{10,13}|-\d{10,13})$`
	re := regexp.MustCompile(pattern)
	isMatch = re.MatchString(in)
	return
}

func ValidateFullNameFormat(in string) (isMatch bool) {
	pattern := `^.{3,60}$`
	re := regexp.MustCompile(pattern)
	isMatch = re.MatchString(in)
	return
}

func ValidatePasswordFormat(in string) (isMatch bool) {
	pattern := `^.{6,64}$`
	re := regexp.MustCompile(pattern)
	isMatch = re.MatchString(in)
	return
}
