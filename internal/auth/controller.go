package auth

import (
	"context"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
)

type UsersService interface {
	Register(ctx context.Context, username string, password string) (string, *shared.ServiceError)
	Login(ctx context.Context, username string, password string) (string, *shared.ServiceError)
}

type Controller struct {
	log         *logger.Logger
	authService UsersService
}

func NewController(log *logger.Logger, service UsersService) *Controller {
	return &Controller{log, service}
}

type Credentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
}

func (ac *Controller) Register(c *fiber.Ctx) error {
	return ac.handleAccess(c, func(credentials Credentials) (string, *shared.ServiceError) {
		return ac.authService.Register(c.Context(), credentials.Login, credentials.Password)
	})
}

func (ac *Controller) Login(c *fiber.Ctx) error {
	return ac.handleAccess(c, func(credentials Credentials) (string, *shared.ServiceError) {
		return ac.authService.Login(c.Context(), credentials.Login, credentials.Password)
	})
}

func (ac *Controller) handleAccess(c *fiber.Ctx, handle func(Credentials) (string, *shared.ServiceError)) error {
	var credentials Credentials
	if err := c.BodyParser(&credentials); err != nil {
		ac.log.Debug(c.Context(), "failed to decode credentials")
		return err
	}
	if err := validator_adapter.ValidateStruct(&credentials); err != nil {
		ac.log.Debug(c.Context(), "invalid credentials struct")
		return fiber_adapter.BadRequest(err)
	}
	accessToken, sErr := handle(credentials)
	if sErr != nil {
		ac.log.Debug(c.Context(), sErr.Msg, sl.Err(sErr.Err))
		return fiber_adapter.ServiceError(sErr)
	}
	return c.JSON(Tokens{
		AccessToken: accessToken,
	})
}
