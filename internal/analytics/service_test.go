package analytics_test

import (
	"bytes"
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/x0k/skillrock-tasks-service/internal/analytics"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

type serviceMocks struct {
	analyticsRepo *analytics.MockAnalyticsRepo
	tasksRepo     *analytics.MockTasksRepo
}

func newTestService(t *testing.T, setup func(serviceMocks)) *analytics.Service {
	var buf bytes.Buffer
	log := logger.New(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	analyticRepo := analytics.NewMockAnalyticsRepo(t)
	tasksRepo := analytics.NewMockTasksRepo(t)
	if setup != nil {
		setup(serviceMocks{
			analyticsRepo: analyticRepo,
			tasksRepo:     tasksRepo,
		})
	}
	return analytics.NewService(
		log,
		tasksRepo,
		analyticRepo,
	)
}

func TestServiceGenerateReport(t *testing.T) {
	tasksCountByStatus := map[tasks.Status]int64{
		tasks.Pending:    1,
		tasks.InProgress: 2,
		tasks.Done:       3,
	}
	averageCompletionTime := float64(1.5)
	completedTasksCount := int64(5)
	overdueTasksCount := int64(6)
	report := analytics.Report{
		TasksCountByStatus:        tasksCountByStatus,
		AverageTaskCompletionTime: averageCompletionTime,
		AmountOfCompletedTasks:    completedTasksCount,
		AmountOfOverdueTasks:      overdueTasksCount,
	}
	unexpectedErr := errors.New("unexpected error")
	cases := []struct {
		name    string
		service *analytics.Service
		err     *shared.ServiceError
	}{
		{
			name: "happy path",
			service: newTestService(t, func(sm serviceMocks) {
				sm.tasksRepo.EXPECT().TasksCountByStatus(mock.Anything).Return(tasksCountByStatus, nil)
				sm.tasksRepo.EXPECT().AverageCompletionTime(mock.Anything).Return(averageCompletionTime, nil)
				sm.tasksRepo.EXPECT().
					CountCompletedAndOverdueTasks(mock.Anything, mock.AnythingOfType("time.Time")).
					Return(completedTasksCount, overdueTasksCount, nil)
				sm.analyticsRepo.EXPECT().SaveReport(mock.Anything, report).Return(nil)
			}),
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(sm serviceMocks) {
				sm.tasksRepo.EXPECT().TasksCountByStatus(mock.Anything).Return(nil, unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.service.GenerateReport(t.Context()); err != nil {
				if c.err == nil ||
					!errors.Is(err.Err, c.err.Err) ||
					err.Expected != c.err.Expected ||
					(c.err.Msg != "" && err.Msg != c.err.Msg) {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
		})
	}
}

func TestServiceReport(t *testing.T) {
	report := analytics.Report{
		TasksCountByStatus:        nil,
		AverageTaskCompletionTime: 1,
		AmountOfCompletedTasks:    2,
		AmountOfOverdueTasks:      3,
	}
	unexpectedErr := errors.New("unexpected error")
	cases := []struct {
		name    string
		service *analytics.Service
		report  analytics.Report
		err     *shared.ServiceError
	}{
		{
			name: "happy path",
			service: newTestService(t, func(sm serviceMocks) {
				sm.analyticsRepo.EXPECT().Report(mock.Anything).Return(report, nil)
			}),
			report: report,
		},
		{
			name: "report not found",
			service: newTestService(t, func(sm serviceMocks) {
				sm.analyticsRepo.EXPECT().Report(mock.Anything).
					Return(report, analytics.ErrReportNotFound)
			}),
			err: shared.NewServiceError(analytics.ErrReportNotFound, ""),
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(sm serviceMocks) {
				sm.analyticsRepo.EXPECT().Report(mock.Anything).
					Return(report, unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			report, err := c.service.Report(t.Context())
			if err != nil {
				if c.err == nil ||
					!errors.Is(err.Err, c.err.Err) ||
					err.Expected != c.err.Expected ||
					(c.err.Msg != "" && err.Msg != c.err.Msg) {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !reflect.DeepEqual(c.report, report) {
				t.Fatalf("expected %v, but got %v", c.report, report)
			}
		})
	}
}
