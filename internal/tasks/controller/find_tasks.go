package tasks_controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func (t *Controller) findTasks(c *fiber.Ctx) error {
	var filter tasks.TasksFilter
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
