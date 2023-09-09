package profile

import (
	"errors"
	"postman/amzn/api/user"
)

func ValidateUpdateRequest(req returnUser) error {
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if !user.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	return nil
}