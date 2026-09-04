package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var usernameRegex = regexp.MustCompile(`^[a-z0-9]+$`)

func ValidateUsername(username string) (string, error) {
	username = strings.TrimSpace(username)

	if username == "" {
		return "", errors.New("username is required")
	}

	length := utf8.RuneCountInString(username)

	if length < 3 || length > 20 {
		return "", errors.New(
			"username must be between 3 and 20 characters",
		)
	}

	if !usernameRegex.MatchString(username) {
		return "", errors.New(
			"username can contain only lowercase letters and numbers",
		)
	}

	return username, nil
}
