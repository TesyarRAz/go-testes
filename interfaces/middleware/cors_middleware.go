package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Append("Access-Control-Allow-Origin", "*")
		c.Append("Access-Control-Allow-Credentials", "true")
		c.Append("Access-Control-Allow-Headers", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With")
		c.Append("Access-Control-Allow-Methods", "POST", "GET", "PUT", "PATCH", "DELETE")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(http.StatusMethodNotAllowed)
		}

		return c.Next()
	}
}
