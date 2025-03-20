package tasks_controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func (t *Controller) ImportTasks(c *fiber.Ctx) error {
	var dto []TaskDTO
	if err := c.BodyParser(&dto); err != nil {
		t.log.Debug(c.Context(), "failed to decode body")
		return err
	}
	if err := validator_adapter.ValidateStruct(&dto); err != nil {
		t.log.Debug(c.Context(), "invalid tasks dto")
		return err
	}
	tasks := make([]tasks.Task, len(dto))
	var err error
	for i, item := range dto {
		if tasks[i], err = taskFromDTO(item); err != nil {
			t.log.Debug(c.Context(), "failed to construct task from dto", slog.Any("task", item))
			return fiber_adapter.BadRequest(err)
		}
	}
	if err := t.tasksService.ImportTasks(c.Context(), tasks); err != nil {
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}
