package analytics

import (
	"context"
	"errors"
	"time"

	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

type AnalyticsRepo interface {
	SaveReport(ctx context.Context, report Report) error
	Report(ctx context.Context) (Report, error)
}

type TasksRepo interface {
	TasksCountByStatus(ctx context.Context) (map[tasks.Status]int64, error)
	AverageCompletionTime(ctx context.Context) (float64, error)
	CompleteAndOverdueTasks(ctx context.Context, duration time.Duration) (int64, int64, error)
}

type Service struct {
	log            *logger.Logger
	tasksRepo      TasksRepo
	analyticsRepo  AnalyticsRepo
	reportDuration time.Duration
}

func NewService(
	log *logger.Logger,
	tasksRepo TasksRepo,
	analyticsRepo AnalyticsRepo,
) *Service {
	return &Service{log, tasksRepo, analyticsRepo, 7 * 24 * time.Hour}
}

func (s *Service) GenerateReport(ctx context.Context) *shared.ServiceError {
	tasksCountByStatus, err := s.tasksRepo.TasksCountByStatus(ctx)
	if err != nil {
		return shared.NewUnexpectedError(err, "failed to count tasks by status")
	}
	averageCompletionTime, err := s.tasksRepo.AverageCompletionTime(ctx)
	if err != nil {
		return shared.NewUnexpectedError(err, "failed to calculate average completion time")
	}
	completeTasks, overdueTasks, err := s.tasksRepo.CompleteAndOverdueTasks(ctx, s.reportDuration)
	if err != nil {
		return shared.NewUnexpectedError(err, "failed to count complete and overdue tasks")
	}
	if err := s.analyticsRepo.SaveReport(ctx, Report{
		TasksCountByStatus:        tasksCountByStatus,
		AverageTaskCompletionTime: averageCompletionTime,
		AmountOfCompletedTasks:    completeTasks,
		AmountOfOverdueTasks:      overdueTasks,
	}); err != nil {
		return shared.NewServiceError(err, "failed to save report")
	}
	return nil
}

func (s *Service) Report(ctx context.Context) (Report, *shared.ServiceError) {
	r, err := s.analyticsRepo.Report(ctx)
	if errors.Is(err, ErrReportNotFound) {
		return r, shared.NewServiceError(err, "the report isn't ready yet")
	}
	if err != nil {
		return r, shared.NewUnexpectedError(err, "failed to load report")
	}
	return r, nil
}
