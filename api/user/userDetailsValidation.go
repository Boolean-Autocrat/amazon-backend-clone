package user

import (
	"errors"
	"regexp"
)

func RegexMatch(regex string, str string) bool {
	re := regexp.MustCompile(regex)
	return re.MatchString(str)
}

func IsValidEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

func IsValidPhoneNum(phoneNum string) bool {
	phoneNumPattern := `^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`
	re := regexp.MustCompile(phoneNumPattern)
	return re.MatchString(phoneNum)
}

func IsValidUsername(username string) error {
	if len(username) < 4 {
		return errors.New("username must be at least 4 characters long")
	}
	usernameRegex := `^[a-zA-Z0-9_]*$`
	if !RegexMatch(usernameRegex, username) {
		return errors.New("username must only contain alphanumeric characters and underscores")
	}
	return nil
}