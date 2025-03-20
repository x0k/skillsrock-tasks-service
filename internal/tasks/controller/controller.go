package tasks_controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
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
	PruneOverdueTasks(ctx context.Context) *shared.ServiceError
}

type Controller struct {
	log          *logger.Logger
	tasksService TasksService
}

func New(
	router fiber.Router,
	log *logger.Logger,
	tasksService TasksService,
) *Controller {
	c := &Controller{log, tasksService}
	router.Get("/", c.findTasks)
	router.Post("/", c.createTask)
	router.Put("/:id", c.updateTaskById)
	router.Delete("/:id", c.removeTaskById)
	router.Post("/import", c.importTasks)
	router.Get("/export", c.exportTasks)
	return c
}
