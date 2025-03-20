package auth

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	logger_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/logger"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
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

func NewController(
	router fiber.Router,
	log *logger.Logger,
	service UsersService,
) *Controller {
	c := &Controller{log, service}
	router.Post("/register", c.register)
	router.Post("/login", c.login)
	return c
}

type Credentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
}

func (ac *Controller) register(c *fiber.Ctx) error {
	var credentials Credentials
	if err := c.BodyParser(&credentials); err != nil {
		ac.log.Debug(c.Context(), "failed to decode credentials")
		return err
	}
	if err := validator_adapter.ValidateStruct(&credentials); err != nil {
		ac.log.Debug(c.Context(), "invalid credentials struct")
		return fiber_adapter.BadRequest(err)
	}
	accessToken, sErr := ac.authService.Register(c.Context(), credentials.Login, credentials.Password)
	if sErr != nil {
		logger_adapter.LogServiceError(ac.log, c, sErr)
		if errors.Is(sErr.Err, ErrLoginIsTaken) {
			return fiber_adapter.SpecificServiceError(sErr, fiber.StatusConflict)
		}
		return fiber_adapter.ServiceError(sErr)
	}
	return c.Status(fiber.StatusCreated).JSON(Tokens{
		AccessToken: accessToken,
	})
}

func (ac *Controller) login(c *fiber.Ctx) error {
	var credentials Credentials
	if err := c.BodyParser(&credentials); err != nil {
		ac.log.Debug(c.Context(), "failed to decode credentials")
		return err
	}
	if err := validator_adapter.ValidateStruct(&credentials); err != nil {
		ac.log.Debug(c.Context(), "invalid credentials struct")
		return fiber_adapter.BadRequest(err)
	}
	accessToken, sErr := ac.authService.Login(c.Context(), credentials.Login, credentials.Password)
	if sErr != nil {
		logger_adapter.LogServiceError(ac.log, c, sErr)
		if errors.Is(sErr.Err, ErrUserNotFound) || errors.Is(sErr.Err, ErrPasswordsMismatch) {
			return fiber_adapter.SpecificServiceError(sErr, fiber.StatusUnauthorized)
		}
		return fiber_adapter.ServiceError(sErr)
	}
	return c.JSON(Tokens{
		AccessToken: accessToken,
	})
}
