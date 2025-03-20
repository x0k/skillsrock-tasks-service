package tasks_controller

import (
	"context"
	"log/slog"

	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
)

func (c *Controller) PruneOverdueTasks(ctx context.Context) {
	if err := c.tasksService.PruneOverdueTasks(ctx); err != nil {
		c.log.Error(
			ctx,
			"failed to prune overdue tasks",
			slog.String("message", err.Msg),
			sl.Err(err.Err),
		)
	}
}
