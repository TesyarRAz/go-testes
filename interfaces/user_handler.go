package interfaces

import (
	"net/http"

	"github.com/TesyarRAz/testes/application"
	"github.com/TesyarRAz/testes/domain/entity"
	"github.com/TesyarRAz/testes/infrastructure/service"
	"github.com/TesyarRAz/testes/infrastructure/validator"
	"github.com/gofiber/fiber/v2"
)

type Users struct {
	ui application.UserAppInterface
	ai service.AuthInterface
	ti service.TokenInterface
}

func NewUsers(ui application.UserAppInterface, ai service.AuthInterface, ti service.TokenInterface) *Users {
	return &Users{ui, ai, ti}
}

func (s *Users) Save(c *fiber.Ctx) error {
	var user entity.User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	validateMsg := validator.ValidateUser(&user, "register")
	if len(validateMsg) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(validateMsg)
	}

	if err := s.ui.Save(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(user.PublicUser())
}

func (s *Users) Get(c *fiber.Ctx) error {
	users, err := s.ui.Get()
	publicUsers := make([]entity.PublicUser, len(users))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	for i, user := range users {
		publicUsers[i] = *user.PublicUser()
	}

	return c.JSON(publicUsers)
}
