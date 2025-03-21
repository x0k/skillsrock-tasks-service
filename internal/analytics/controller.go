package analytics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	logger_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

type AnalyticsService interface {
	GenerateReport(ctx context.Context) *shared.ServiceError
	Report(ctx context.Context) (Report, *shared.ServiceError)
}

type Controller struct {
	log              *logger.Logger
	analyticsService AnalyticsService
}

func NewController(
	router fiber.Router,
	log *logger.Logger,
	analyticsService AnalyticsService,
) *Controller {
	c := &Controller{log, analyticsService}
	router.Get("/", c.report)
	return c
}

func (a *Controller) GenerateReport(ctx context.Context) {
	if err := a.analyticsService.GenerateReport(ctx); err != nil {
		a.log.Error(
			ctx,
			"failed to generate report",
			slog.String("message", err.Msg),
			sl.Err(err.Err),
		)
	}
}

type ReportDTO struct {
	PendingTasksCount           int64  `json:"pending_tasks_count"`
	InProgressTasksCount        int64  `json:"in_progress_tasks_count"`
	DoneTasksCount              int64  `json:"done_tasks_count"`
	AverageCompletionTimeInDays string `json:"average_completion_time_in_days"`
	AmountOfCompletedTasks      int64  `json:"amount_of_completed_tasks"`
	AmountOfOverdueTasks        int64  `json:"amount_of_overdue_tasks"`
}

const dayInSeconds = 24 * 60 * 60

func reportToDTO(r Report) ReportDTO {
	return ReportDTO{
		PendingTasksCount:           r.TasksCountByStatus[tasks.Pending],
		InProgressTasksCount:        r.TasksCountByStatus[tasks.InProgress],
		DoneTasksCount:              r.TasksCountByStatus[tasks.Done],
		AverageCompletionTimeInDays: fmt.Sprintf("%.2f", r.AverageTaskCompletionTime/dayInSeconds),
		AmountOfCompletedTasks:      r.AmountOfCompletedTasks,
		AmountOfOverdueTasks:        r.AmountOfCompletedTasks,
	}
}

func (a *Controller) report(c *fiber.Ctx) error {
	r, err := a.analyticsService.Report(c.Context())
	if err != nil {
		logger_adapter.LogServiceError(a.log, c, err)
		if errors.Is(err.Err, ErrReportNotFound) {
			return fiber.ErrNotFound
		}
		return fiber_adapter.ServiceError(err)
	}
	return c.JSON(reportToDTO(r))
}
