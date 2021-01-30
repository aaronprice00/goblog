package formaterror

import (
	"errors"
	"strings"
)

// FormatError changes error response for testing
func FormatError(err string) error {
	if strings.Contains(err, "username") {
		return errors.New("Username Already Taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email Already Used")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title Already Used")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")
}
