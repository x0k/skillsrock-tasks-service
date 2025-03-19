package users

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
)

func NewAuth(
	app fiber.Router,
	log *logger.Logger,
	queries *db.Queries,
	secret []byte,
	tokenLifetime time.Duration,
) error {
	repo := newRepo(
		log.With(sl.Component("users_repository")),
		queries,
	)

	service := newService(
		log.With(sl.Component("users_service")),
		secret,
		tokenLifetime,
		repo,
	)

	controller := newAuthController(
		log.With(sl.Component("auth_controller")),
		service,
	)

	app.Post("/register", controller.Register)
	app.Post("/login", controller.Login)

	return nil
}
