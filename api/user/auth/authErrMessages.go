package auth

func AuthErrMessage(errorMessage string) string {
	baseStr := "pq: duplicate key value violates unique constraint"
	if errorMessage == baseStr+" \"users_username_key\"" {
		return "username already exists"
	}
	if errorMessage == baseStr+" \"users_email_key\"" {
		return "email already exists"
	}
	if errorMessage == baseStr+" \"users_phonenum_key\"" {
		return "phone number already exists"
	}
	return ""
}