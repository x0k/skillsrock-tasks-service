package tasks_controller

import (
	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
)

func (t *Controller) removeTaskById(c *fiber.Ctx) error {
	taskId, err := t.taskId(c, c.Params("id"))
	if err != nil {
		return err
	}
	if err := t.tasksService.RemoveTaskById(c.Context(), taskId); err != nil {
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}
