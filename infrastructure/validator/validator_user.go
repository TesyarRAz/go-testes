package validator

import (
	"strings"

	"github.com/TesyarRAz/testes/domain/entity"
)

const (
	usernameRequired = "username is required"
	passwordRequired = "password is required"
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
	}

	return msg
}
