package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
)

type TasksRepo interface {
	SaveTask(ctx context.Context, task Task) error
	FindTasks(ctx context.Context, filter TasksFilter) ([]Task, error)
	UpdateTaskById(ctx context.Context, id TaskId, params TaskParams) error
	RemoveTaskById(ctx context.Context, id TaskId) error
	SaveTasks(ctx context.Context, tasks []Task) error
	AllTasks(ctx context.Context) ([]Task, error)
	RemoveOverdueTasksWithDueDateBefore(ctx context.Context, date time.Time) error
}

type Service struct {
	log           *logger.Logger
	tasksRepo     TasksRepo
	pruneDuration time.Duration
}

func NewService(
	log *logger.Logger,
	repo TasksRepo,
) *Service {
	return &Service{log, repo, 7 * 24 * time.Hour}
}

func (s *Service) CreateTask(ctx context.Context, params TaskParams) *shared.ServiceError {
	now := time.Now()
	task, err := NewTask(
		NewTaskId(),
		params.Title,
		params.Description,
		params.Status,
		params.Priority,
		params.DueDate,
		now,
		now,
	)
	if err != nil {
		return shared.NewServiceError(err, "failed to create task")
	}
	if err := s.tasksRepo.SaveTask(ctx, task); err != nil {
		return shared.NewUnexpectedError(err, "failed to save task")
	}
	return nil
}

func (s *Service) FindTasks(ctx context.Context, filter TasksFilter) ([]Task, *shared.ServiceError) {
	tasks, err := s.tasksRepo.FindTasks(ctx, filter)
	if err != nil {
		return tasks, shared.NewUnexpectedError(err, "failed to filter tasks")
	}
	return tasks, nil
}

func (s *Service) UpdateTaskById(ctx context.Context, id TaskId, params TaskParams) *shared.ServiceError {
	err := s.tasksRepo.UpdateTaskById(ctx, id, params)
	if errors.Is(err, ErrTaskNotFound) {
		return shared.NewServiceError(err, fmt.Sprintf("task with id %q not found", id.String()))
	}
	if errors.Is(err, ErrTaskIsAlreadyDone) {
		return shared.NewServiceError(err, "the task to be updated has already been completed")
	}
	if err != nil {
		return shared.NewUnexpectedError(err, "failed to update task")
	}
	return nil
}

func (s *Service) RemoveTaskById(ctx context.Context, id TaskId) *shared.ServiceError {
	err := s.tasksRepo.RemoveTaskById(ctx, id)
	if errors.Is(err, ErrTaskNotFound) {
		return shared.NewServiceError(err, fmt.Sprintf("task with id %q not found", id.String()))
	}
	if err != nil {
		return shared.NewUnexpectedError(err, "failed to remove task")
	}
	return nil
}

func (s *Service) ExportTasks(ctx context.Context) ([]Task, *shared.ServiceError) {
	if tasks, err := s.tasksRepo.AllTasks(ctx); err != nil {
		return tasks, shared.NewUnexpectedError(err, "failed to load tasks")
	} else {
		return tasks, nil
	}
}

func (s *Service) ImportTasks(ctx context.Context, tasks []Task) *shared.ServiceError {
	if err := s.tasksRepo.SaveTasks(ctx, tasks); err != nil {
		return shared.NewUnexpectedError(err, "failed to save tasks")
	}
	return nil
}

func (s *Service) PruneOverdueTasks(ctx context.Context) *shared.ServiceError {
	date := time.Now().Add(-s.pruneDuration)
	err := s.tasksRepo.RemoveOverdueTasksWithDueDateBefore(ctx, date)
	if err != nil {
		return shared.NewUnexpectedError(err, "failed to remove overdue tasks")
	}
	return nil
}
