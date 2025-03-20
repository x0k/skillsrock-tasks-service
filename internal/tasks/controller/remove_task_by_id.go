package tasks_controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	logger_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/logger"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func (t *Controller) removeTaskById(c *fiber.Ctx) error {
	taskId, err := t.taskId(c, c.Params("id"))
	if err != nil {
		return err
	}
	if err := t.tasksService.RemoveTaskById(c.Context(), taskId); err != nil {
		logger_adapter.LogServiceError(t.log, c, err)
		if errors.Is(err.Err, tasks.ErrTaskNotFound) {
			return fiber.ErrNotFound
		}
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
