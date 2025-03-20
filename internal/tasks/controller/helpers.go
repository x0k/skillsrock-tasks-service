package tasks_controller

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func (t *Controller) taskParams(c *fiber.Ctx) (tasks.TaskParams, error) {
	var dto CreateTaskDTO
	if err := c.BodyParser(&dto); err != nil {
		t.log.Debug(c.Context(), "failed to decode body")
		return tasks.TaskParams{}, err
	}
	if err := validator_adapter.ValidateStruct(&dto); err != nil {
		t.log.Debug(c.Context(), "invalid create task dto struct")
		return tasks.TaskParams{}, fiber_adapter.BadRequest(err)
	}
	params := tasks.TaskParams{
		Title:       dto.Title,
		Description: dto.Description,
	}
	var err error
	if params.Status, err = t.status(c, dto.Status); err != nil {
		return params, err
	}
	if params.Priority, err = t.priority(c, dto.Priority); err != nil {
		return params, err
	}
	if params.DueDate, err = t.date(c, dto.DueDate); err != nil {
		return params, err
	}
	return params, nil
}

func (t *Controller) taskId(c *fiber.Ctx, value string) (tasks.TaskId, error) {
	taskId, err := tasks.ParseTaskId(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid task id value", slog.String("task_id", value))
		return taskId, fiber_adapter.BadRequest(err)
	}
	return taskId, nil
}

func (t *Controller) status(c *fiber.Ctx, value string) (tasks.Status, error) {
	status, err := tasks.ParseStatus(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid status value", slog.String("status", value))
		return status, fiber_adapter.BadRequest(err)
	}
	return status, nil
}

func (t *Controller) priority(c *fiber.Ctx, value string) (tasks.Priority, error) {
	priority, err := tasks.ParsePriority(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid priority value", slog.String("priority", value))
		return priority, fiber_adapter.BadRequest(err)
	}
	return priority, nil
}

func (t *Controller) date(c *fiber.Ctx, value string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid date value", slog.String("date", value))
		return date, fiber_adapter.BadRequest(err)
	}
	return date, nil
}
