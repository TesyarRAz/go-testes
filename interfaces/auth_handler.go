package interfaces

import (
	"net/http"

	"github.com/TesyarRAz/testes/application"
	"github.com/TesyarRAz/testes/domain/entity"
	"github.com/TesyarRAz/testes/infrastructure/service"
	"github.com/TesyarRAz/testes/infrastructure/validator"
	"github.com/TesyarRAz/testes/interfaces/resource"
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
	var (
		user         entity.User
		loggedUser   *entity.User
		tokenDetails *entity.TokenDetails

		err error
	)

	if err = c.BodyParser(&user); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}
	if msg := validator.ValidateUser(&user, "login"); len(msg) > 0 {
		return ErrorMessage(c, msg, http.StatusUnprocessableEntity)
	}
	if loggedUser, err = a.ui.FindUserByUsernameAndPassword(&user); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}
	if tokenDetails, err = a.ti.CreateToken(loggedUser.ID); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}
	if err = a.ai.CreateAuth(c.Context(), loggedUser.ID, tokenDetails); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}

	return c.JSON(
		resource.LoginResource(*loggedUser, *tokenDetails),
	)
}

func (a *Auth) Register(c *fiber.Ctx) error {
	var (
		user entity.User

		err error
	)

	if err = c.BodyParser(&user); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}
	if msg := validator.ValidateUser(&user, "register"); len(msg) > 0 {
		return ErrorMessage(c, msg, http.StatusUnprocessableEntity)
	}
	if err = a.ui.Save(&user); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}

	return SimpleMessage(c, "Success Register")
}

func (a *Auth) Logout(c *fiber.Ctx) error {
	metadata, err := a.ti.ExtractAccessTokenMetadata(c.GetReqHeaders())
	if err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnauthorized)
	}

	if err := a.ai.DeleteTokens(c.Context(), metadata); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnauthorized)
	}

	return SimpleMessage(c, "Success Logout")
}

func (a *Auth) Refresh(c *fiber.Ctx) error {
	var (
		mapToken = map[string]string{}
		metadata *entity.RefreshDetails

		tokenDetails *entity.TokenDetails

		err error
	)

	if err = c.BodyParser(&mapToken); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}

	if metadata, err = a.ti.ExtractRefreshTokenMetadata(mapToken); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnauthorized)
	}

	if err := a.ai.DeleteRefresh(c.Context(), metadata); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnauthorized)
	}

	if tokenDetails, err = a.ti.CreateToken(metadata.UserId); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}
	if err = a.ai.CreateAuth(c.Context(), metadata.UserId, tokenDetails); err != nil {
		return CustomMessage(c, err.Error(), http.StatusUnprocessableEntity)
	}

	return c.JSON(
		resource.RefreshResource(*tokenDetails),
	)
}
