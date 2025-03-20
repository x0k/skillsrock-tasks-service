package tasks_controller

import (
	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func (t *Controller) removeTaskById(c *fiber.Ctx) error {
	var taskId tasks.TaskId
	if err := t.setTaskId(c, &taskId, c.Params("id")); err != nil {
		return err
	}
	if err := t.tasksService.RemoveTaskById(c.Context(), taskId); err != nil {
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}
