package tasks_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	logger_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/logger"
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
		fmt.Println(filter.Status)
		if s, err := t.status(c, status); err != nil {
			return err
		} else {
			filter.Status = &s
		}
	}
	priority := c.Query("priority")
	if priority != "" {
		if p, err := t.priority(c, priority); err != nil {
			return err
		} else {
			filter.Priority = &p
		}
	}
	dueBefore := c.Query("due_before")
	if dueBefore != "" {
		if d, err := t.date(c, dueBefore); err != nil {
			return err
		} else {
			filter.DueBefore = &d
		}
	}
	dueAfter := c.Query("due_after")
	if dueAfter != "" {
		if d, err := t.date(c, dueAfter); err != nil {
			return err
		} else {
			filter.DueAfter = &d
		}
	}
	tasks, err := t.tasksService.FindTasks(c.Context(), filter)
	if err != nil {
		logger_adapter.LogServiceError(t.log, c, err)
		return fiber_adapter.ServiceError(err)
	}
	tasksDto := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		tasksDto[i] = taskToDTO(t)
	}
	return c.JSON(tasksDto)
}
