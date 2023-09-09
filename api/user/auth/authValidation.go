package auth

import (
	"errors"
	"postman/amzn/api/user"
)

func ValidateUserRequest(req apiUser) error {
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if !user.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	return nil
}