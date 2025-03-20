package tasks_controller

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	logger_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/logger"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func (t *Controller) importTasks(c *fiber.Ctx) error {
	var dto []TaskDTO
	if err := c.BodyParser(&dto); err != nil {
		t.log.Debug(c.Context(), "failed to decode body")
		return err
	}
	if err := validator_adapter.ValidateArray(dto); err != nil {
		t.log.Debug(c.Context(), "invalid tasks dto")
		return err
	}
	tasksList := make([]tasks.Task, len(dto))
	var err error
	for i, item := range dto {
		if tasksList[i], err = taskFromDTO(item); err != nil {
			t.log.Debug(c.Context(), "failed to construct task from dto", slog.Any("task", item))
			return fiber_adapter.BadRequest(err)
		}
	}
	if err := t.tasksService.ImportTasks(c.Context(), tasksList); err != nil {
		logger_adapter.LogServiceError(t.log, c, err)
		if errors.Is(err.Err, tasks.ErrTaskIdsConflict) {
			return fiber.ErrConflict
		}
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}
