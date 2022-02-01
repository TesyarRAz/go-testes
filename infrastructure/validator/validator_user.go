package validator

import (
	"strings"

	"github.com/TesyarRAz/testes/domain/entity"
)

const (
	firstNameRequired = "first name is required"
	lastNameRequired  = "last name is required"
	usernameRequired  = "username is required"
	passwordRequired  = "password is required"
)

func ValidateUser(user *entity.User, action string) map[string]string {
	msg := make(map[string]string)

	switch strings.ToLower(action) {
	case "login":
		if user.Username == "" {
			msg["username_required"] = usernameRequired
		}
		if user.Password == "" {
			msg["password_required"] = passwordRequired
		}
	case "register":
		if user.FirstName == "" {
			msg["first_name_required"] = firstNameRequired
		}
		if user.LastName == "" {
			msg["last_name_required"] = lastNameRequired
		}
		if user.Username == "" {
			msg["username_required"] = usernameRequired
		}
		if user.Password == "" {
			msg["password_required"] = passwordRequired
		}
	}

	return msg
}
