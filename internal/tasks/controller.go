package tasks

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	validator_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/validator"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
)

type TasksService interface {
	CreateTask(ctx context.Context, params TaskParams) *shared.ServiceError
	FindTasks(ctx context.Context, filter TasksFilter) ([]Task, *shared.ServiceError)
	UpdateTaskById(ctx context.Context, id TaskId, params TaskParams) *shared.ServiceError
	RemoveTaskById(ctx context.Context, id TaskId) *shared.ServiceError
	ExportTasks(ctx context.Context) ([]Task, *shared.ServiceError)
	ImportTasks(ctx context.Context, tasks []Task) *shared.ServiceError
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

type CreateTaskDTO struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status" validate:"required"`
	Priority    string  `json:"priority" validate:"required"`
	DueDate     string  `json:"due_date" validate:"required"`
}

func (t *Controller) CreateTask(c *fiber.Ctx) error {
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

type TaskDTO struct {
	Id          string  `json:"id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status" validate:"required"`
	Priority    string  `json:"priority" validate:"required"`
	DueDate     string  `json:"due_date" validate:"required"`
	CreatedAt   string  `json:"created_at" validate:"required"`
	UpdatedAt   string  `json:"updated_at" validate:"required"`
}

func taskToDTO(task Task) TaskDTO {
	return TaskDTO{
		Id:          uuid.UUID(task.Id).String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status.String(),
		Priority:    task.Priority.String(),
		DueDate:     task.DueDate.Format(time.DateOnly),
		CreatedAt:   task.CreatedAt.Format(time.DateOnly),
		UpdatedAt:   task.UpdatedAt.Format(time.DateOnly),
	}
}

func taskFromDTO(dto TaskDTO) (Task, error) {
	task := Task{
		Title:       dto.Title,
		Description: dto.Description,
	}
	var err error
	if task.Id, err = NewTaskId(dto.Id); err != nil {
		return task, err
	}
	if task.Status, err = NewStatus(dto.Status); err != nil {
		return task, err
	}
	if task.Priority, err = NewPriority(dto.Priority); err != nil {
		return task, err
	}
	if task.DueDate, err = time.Parse(time.DateOnly, dto.DueDate); err != nil {
		return task, err
	}
	if task.CreatedAt, err = time.Parse(time.DateOnly, dto.CreatedAt); err != nil {
		return task, err
	}
	if task.UpdatedAt, err = time.Parse(time.DateOnly, dto.UpdatedAt); err != nil {
		return task, err
	}
	return task, nil
}

func (t *Controller) FindTasks(c *fiber.Ctx) error {
	filter := TasksFilter{}
	title := c.Query("title")
	if title != "" {
		filter.Title = &title
	}
	status := c.Query("status")
	if status != "" {
		if err := t.setStatus(c, filter.Status, status); err != nil {
			return err
		}
	}
	priority := c.Query("priority")
	if priority != "" {
		if err := t.setPriority(c, filter.Priority, priority); err != nil {
			return err
		}
	}
	dueBefore := c.Query("due_before")
	if dueBefore != "" {
		if err := t.setDate(c, filter.DueBefore, dueBefore); err != nil {
			return err
		}
	}
	dueAfter := c.Query("due_after")
	if dueAfter != "" {
		if err := t.setDate(c, filter.DueAfter, dueAfter); err != nil {
			return err
		}
	}
	tasks, err := t.tasksService.FindTasks(c.Context(), filter)
	if err != nil {
		t.log.Debug(c.Context(), "failed to find tasks", slog.Any("filter", filter))
		return fiber_adapter.ServiceError(err)
	}
	tasksDto := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		tasksDto[i] = taskToDTO(t)
	}
	return c.JSON(tasksDto)
}

func (t *Controller) UpdateTaskById(c *fiber.Ctx) error {
	var taskId TaskId
	if err := t.setTaskId(c, &taskId, c.Params("id")); err != nil {
		return err
	}
	params, err := t.taskParams(c)
	if err != nil {
		return err
	}
	if err := t.tasksService.UpdateTaskById(c.Context(), taskId, params); err != nil {
		t.log.Debug(c.Context(), "failed to update task", slog.Any("task_id", taskId), slog.Any("task_params", params))
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (t *Controller) RemoveTaskById(c *fiber.Ctx) error {
	var taskId TaskId
	if err := t.setTaskId(c, &taskId, c.Params("id")); err != nil {
		return err
	}
	if err := t.tasksService.RemoveTaskById(c.Context(), taskId); err != nil {
		return fiber_adapter.ServiceError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

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
	tasks := make([]Task, len(dto))
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

func (t *Controller) taskParams(c *fiber.Ctx) (TaskParams, error) {
	var dto CreateTaskDTO
	if err := c.BodyParser(&dto); err != nil {
		t.log.Debug(c.Context(), "failed to decode body")
		return TaskParams{}, err
	}
	if err := validator_adapter.ValidateStruct(&dto); err != nil {
		t.log.Debug(c.Context(), "invalid create task dto struct")
		return TaskParams{}, fiber_adapter.BadRequest(err)
	}
	params := TaskParams{
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

func (t *Controller) setTaskId(c *fiber.Ctx, out *TaskId, value string) error {
	taskId, err := NewTaskId(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid task id value", slog.String("task_id", value))
		return fiber_adapter.BadRequest(err)
	}
	*out = taskId
	return nil
}

func (t *Controller) setStatus(c *fiber.Ctx, out *Status, value string) error {
	status, err := NewStatus(value)
	if err != nil {
		t.log.Debug(c.Context(), "invalid status value", slog.String("status", value))
		return fiber_adapter.BadRequest(err)
	}
	*out = status
	return nil
}

func (t *Controller) setPriority(c *fiber.Ctx, out *Priority, value string) error {
	priority, err := NewPriority(value)
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
