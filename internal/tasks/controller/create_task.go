package tasks_controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
)

type CreateTaskDTO struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status" validate:"required"`
	Priority    string  `json:"priority" validate:"required"`
	DueDate     string  `json:"due_date" validate:"required"`
}

func (t *Controller) createTask(c *fiber.Ctx) error {
	params, err := t.taskParams(c)
	if err != nil {
		return err
	}
	if err := t.tasksService.CreateTask(c.Context(), params); err != nil {
		t.log.Debug(c.Context(), "failed to create task", slog.Any("task_params", params))
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}
