package resource

import (
	"github.com/TesyarRAz/testes/domain/entity"
	"github.com/gofiber/fiber/v2"
)

func LoginResource(user entity.User, ts entity.TokenDetails) fiber.Map {
	return fiber.Map{
		"id":            user.ID,
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
	}
}

func RefreshResource(ts entity.TokenDetails) fiber.Map {
	return fiber.Map{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
}
