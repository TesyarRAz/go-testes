package interfaces

import (
	"github.com/TesyarRAz/testes/application"
	"github.com/TesyarRAz/testes/infrastructure/service"
	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	ui application.UserAppInterface
	ai service.AuthInterface
	ti service.TokenInterface
}

func NewAuth(ui application.UserAppInterface, ai service.AuthInterface, ti service.TokenInterface) *Auth {
	return &Auth{ui, ai, ti}
}

func (a *Auth) Login(c *fiber.Ctx) error {
	return nil
}
