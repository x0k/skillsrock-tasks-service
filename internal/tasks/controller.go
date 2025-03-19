package tasks

import (
	"context"
	"iter"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
)

type TasksService interface {
	CreateTask(ctx context.Context, dto CreateTaskDTO) *shared.ServiceError
	FindTasks(ctx context.Context, dto FindTasksDTO) ([]Task, *shared.ServiceError)
	UpdateTask(ctx context.Context, id TaskId, dto CreateTaskDTO) *shared.ServiceError
	RemoveTask(ctx context.Context, id TaskId) *shared.ServiceError
	ExportTasks(ctx context.Context) iter.Seq2[Task, error]
	ImportTasks(ctx context.Context, tasks iter.Seq2[TaskDTO, error]) *shared.ServiceError
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

func (t *Controller) CreateTask(c *fiber.Ctx) error {
	var dto CreateTaskDTO
	if err := c.BodyParser(&dto); err != nil {
		t.log.Debug(c.Context(), "failed to decode body")
		return err
	}
	if err := validator_adapter.ValidateStruct(&dto); err != nil {
		t.log.Debug(c.Context(), "invalid create task dto struct")
		return fiber_adapter.BadRequest(err)
	}
	if err := t.tasksService.CreateTask(c.Context(), dto); err != nil {
		t.log.Debug(c.Context(), "failed to create task", slog.Any("task_dto", dto))
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (t *Controller) FindTasks(c *fiber.Ctx) error {
	dto := FindTasksDTO{
		Title:     c.Query("title"),
		Status:    c.Query("status"),
		Priority:  c.Query("priority"),
		DueBefore: c.Query("due_before"),
		DueAfter:  c.Query("due_after"),
	}
	tasks, err := t.tasksService.FindTasks(c.Context(), dto)
	if err != nil {
		t.log.Debug(c.Context(), "failed to find tasks", slog.Any("filter", dto))
	}

}
