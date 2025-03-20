package analytics

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	fiber_adapter "github.com/x0k/skillrock-tasks-service/internal/adapters/fiber"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
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
		a.log.Debug(ctx, "failed to generate report")
	}
}

type ReportDTO struct {
	PendingTasksCount           int64  `json:"PendingTasksCount"`
	InProgressTasksCount        int64  `json:"InProgressTasksCount"`
	DoneTasksCount              int64  `json:"DoneTasksCount"`
	AverageCompletionTimeInDays string `json:"AverageCompletionTimeInDays"`
	AmountOfCompletedTasks      int64  `json:"AmountOfCompletedTasks"`
	AmountOfOverdueTasks        int64  `json:"AmountOfOverdueTasks"`
}

func reportToDTO(r Report) ReportDTO {
	return ReportDTO{
		PendingTasksCount:           r.TasksCountByStatus[tasks.Pending],
		InProgressTasksCount:        r.TasksCountByStatus[tasks.InProgress],
		DoneTasksCount:              r.TasksCountByStatus[tasks.Done],
		AverageCompletionTimeInDays: fmt.Sprintf("%.2f", r.AverageTaskCompletionTime),
		AmountOfCompletedTasks:      r.AmountOfCompletedTasks,
		AmountOfOverdueTasks:        r.AmountOfCompletedTasks,
	}
}

func (a *Controller) report(c *fiber.Ctx) error {
	r, err := a.analyticsService.Report(c.Context())
	if err != nil {
		a.log.Debug(c.Context(), "failed to load report")
		return fiber_adapter.ServiceError(err)
	}
	return c.JSON(reportToDTO(r))
}
