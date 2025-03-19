package tasks

import (
	"context"

	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
)

type TasksService interface {
	CreateTask(ctx context.Context, dto CreateTaskDTO) error
	FindTasks(ctx context.Context, dto FindTasksDTO) []task
	UpdateTask(ctx context.Context, id taskId, dto CreateTaskDTO) error
	RemoveTask(ctx context.Context, id taskId) error
}

type controller struct {
	log          *logger.Logger
	tasksService TasksService
}
