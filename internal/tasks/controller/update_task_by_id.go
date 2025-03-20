package tasks_controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
)

func (t *Controller) updateTaskById(c *fiber.Ctx) error {
	taskId, err := t.taskId(c, c.Params("id"))
	if err != nil {
		return err
	}
	params, err := t.taskParams(c)
	if err != nil {
		return err
	}
	if err := t.tasksService.UpdateTaskById(c.Context(), taskId, params); err != nil {
		t.log.Debug(
			c.Context(),
			"failed to update task",
			slog.Any("task_id", taskId),
			slog.Any("task_params", params),
		)
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}
