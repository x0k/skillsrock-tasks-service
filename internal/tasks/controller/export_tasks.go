package tasks_controller

import (
	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	logger_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/logger"
)

func (t *Controller) exportTasks(c *fiber.Ctx) error {
	tasks, err := t.tasksService.ExportTasks(c.Context())
	if err != nil {
		logger_adapter.LogServiceError(t.log, c, err)
		return fiber_adapter.ServiceError(err)
	}
	tasksDto := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		tasksDto[i] = taskToDTO(t)
	}
	return c.JSON(tasksDto)
}
