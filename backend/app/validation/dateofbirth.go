package validation

import (
	"errors"
	"strings"
	"time"
)

func ValidateDateOfBirth(
	value string,
) (*time.Time, error) {
	value = strings.TrimSpace(value)

	// Allows the user to remove an optional date of birth.
	if value == "" {
		return nil, nil
	}

	dateOfBirth, err := time.Parse(
		"02-01-2006",
		value,
	)
	if err != nil {
		return nil, errors.New(
			"date of birth must use DD-MM-YYYY format",
		)
	}

	if dateOfBirth.After(time.Now().UTC()) {
		return nil, errors.New(
			"date of birth cannot be in the future",
		)
	}

	return &dateOfBirth, nil
}
