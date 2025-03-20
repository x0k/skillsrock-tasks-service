package tasks_controller

import (
	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
)

func (t *Controller) ExportTasks(c *fiber.Ctx) error {
	tasks, err := t.tasksService.ExportTasks(c.Context())
	if err != nil {
		t.log.Debug(c.Context(), "failed to export tasks")
		return fiber_adapter.ServiceError(err)
	}
	tasksDto := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		tasksDto[i] = taskToDTO(t)
	}
	return c.JSON(tasksDto)
}
