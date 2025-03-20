package logger_adapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
)

func LogServiceError(log *logger.Logger, c *fiber.Ctx, err *shared.ServiceError) {
	if err.Expected {
		log.Debug(c.Context(), err.Msg, sl.Err(err.Err))
	} else {
		log.Error(c.Context(), err.Msg, sl.Err(err.Err))
	}
}
