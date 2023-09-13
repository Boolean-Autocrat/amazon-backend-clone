package auth

import (
	"errors"
	"regexp"
)

func RegexMatch(regex string, str string) bool {
	re := regexp.MustCompile(regex)
	return re.MatchString(str)
}

func IsPasswordValid(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	uppercase_regex := `[A-Z]`
	lowercase_regex := `[a-z]`
	number_regex := `[0-9]`
	special_regex := `[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`
	if !RegexMatch(uppercase_regex, password) {
		return errors.New("password must contain at least 1 uppercase letter")
	}
	if !RegexMatch(lowercase_regex, password) {
		return errors.New("password must contain at least 1 lowercase letter")
	}
	if !RegexMatch(number_regex, password) {
		return errors.New("password must contain at least 1 number")
	}
	if !RegexMatch(special_regex, password) {
		return errors.New("password must contain at least 1 special character")
	}
	return nil
}