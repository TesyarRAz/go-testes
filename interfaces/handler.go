package interfaces

import "github.com/gofiber/fiber/v2"

func CustomMessage(c *fiber.Ctx, msg string, status int) error {
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
		"data":    nil,
	})
}

func SimpleMessage(c *fiber.Ctx, msg string) error {
	return c.JSON(fiber.Map{
		"message": msg,
		"data":    nil,
	})
}

func CustomMessageAndData(c *fiber.Ctx, msg string, status int, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func SimpleMessageAndData(c *fiber.Ctx, msg string, data interface{}) error {
	return c.JSON(fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func ErrorMessage(c *fiber.Ctx, err interface{}, status int) error {
	return c.JSON(fiber.Map{
		"message": nil,
		"errors":  err,
	})
}
