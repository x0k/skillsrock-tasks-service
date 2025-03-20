package tasks_controller

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

type TasksService interface {
	CreateTask(ctx context.Context, params tasks.TaskParams) *shared.ServiceError
	FindTasks(ctx context.Context, filter tasks.TasksFilter) ([]tasks.Task, *shared.ServiceError)
	UpdateTaskById(ctx context.Context, id tasks.TaskId, params tasks.TaskParams) *shared.ServiceError
	RemoveTaskById(ctx context.Context, id tasks.TaskId) *shared.ServiceError
	ExportTasks(ctx context.Context) ([]tasks.Task, *shared.ServiceError)
	ImportTasks(ctx context.Context, tasks []tasks.Task) *shared.ServiceError
}

type Controller struct {
	log          *logger.Logger
	tasksService TasksService
}

func NewController(
	log *logger.Logger,
	tasksService TasksService,
) *Controller {
	return &Controller{log, tasksService}
}

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
	if err := t.setStatus(c, &params.Status, dto.Status); err != nil {
		return params, err
	}
	if err := t.setPriority(c, &params.Priority, dto.Priority); err != nil {
		return params, err
	}
	if err := t.setDate(c, &params.DueDate, dto.DueDate); err != nil {
		return params, err
	}
	return params, nil
}

func (t *Controller) setTaskId(c *fiber.Ctx, out *tasks.TaskId, value string) error {
	taskId, err := tasks.NewTaskId(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid task id value", slog.String("task_id", value))
		return fiber_adapter.BadRequest(err)
	}
	*out = taskId
	return nil
}

func (t *Controller) setStatus(c *fiber.Ctx, out *tasks.Status, value string) error {
	status, err := tasks.NewStatus(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid status value", slog.String("status", value))
		return fiber_adapter.BadRequest(err)
	}
	*out = status
	return nil
}

func (t *Controller) setPriority(c *fiber.Ctx, out *tasks.Priority, value string) error {
	priority, err := tasks.NewPriority(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid priority value", slog.String("priority", value))
		return fiber_adapter.BadRequest(err)
	}
	*out = priority
	return nil
}

func (t *Controller) setDate(c *fiber.Ctx, out *time.Time, value string) error {
	date, err := time.Parse(time.DateOnly, value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid date value", slog.String("date", value))
		return fiber_adapter.BadRequest(err)
	}
	*out = date
	return nil
}
