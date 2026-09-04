package validation

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(
	`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[A-Za-z]{2,}$`,
)

func ValidateEmail(email string, required bool) (string, error) {

	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		if required {
			return "", errors.New("email is required")
		}
		return "", nil
	}

	if len(email) > 100 {
		return "", errors.New("email must not exceed 100 characters")
	}

	if !emailRegex.MatchString(email) {
		return "", errors.New("invalid email address")
	}

	return email, nil
}
