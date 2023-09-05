package users

import (
	"errors"
	"regexp"
)

func isValidEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

func ValidateUserRequest(req apiUser) error {
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	return nil
}