package validation

import (
	"errors"
	"regexp"
	"strings"
)

var digitsOnlyRegex = regexp.MustCompile(`^[0-9]+$`)

func ValidatePhone(phone string) (string, error) {

	phone = strings.TrimSpace(phone)

	if !digitsOnlyRegex.MatchString(phone) {
		return "", errors.New("phone number can contain only digits")
	}

	if len(phone) != 10 {
		return "", errors.New("phone number must contain exactly 10 digits")
	}

	return phone, nil
}
