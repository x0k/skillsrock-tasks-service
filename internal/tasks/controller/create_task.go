package tasks_controller

import (
	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	logger_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/logger"
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
		logger_adapter.LogServiceError(t.log, c, err)
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}
