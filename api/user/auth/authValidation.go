package auth

import (
	"errors"
	"postman/amzn/api/user"
)

func ValidateUserRequest(req apiUser) error {
	if IsPasswordValid(req.Password) != nil {
		return IsPasswordValid(req.Password)
	}
	if user.IsValidUsername(req.Username) != nil {
		return user.IsValidUsername(req.Username)
	}
	if !user.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if !user.IsValidPhoneNum(req.PhoneNum) {
		return errors.New("invalid phone number format")
	}

	return nil
}