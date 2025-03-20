package tasks_test

import (
	"bytes"
	"errors"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func newTestService(t *testing.T, setup func(repo *tasks.MockTasksRepo)) *tasks.Service {
	var buf bytes.Buffer
	log := logger.New(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	repo := tasks.NewMockTasksRepo(t)
	if setup != nil {
		setup(repo)
	}
	return tasks.NewService(
		log,
		repo,
	)
}

func TestServiceCreateTask(t *testing.T) {
	title := "title"
	dueDate := time.Now().Add(time.Hour)
	params := tasks.TaskParams{
		Title:    title,
		DueDate:  dueDate,
		Status:   tasks.Pending,
		Priority: tasks.Low,
	}
	unexpectedErr := errors.New("unexpected err")
	cases := []struct {
		name    string
		service *tasks.Service
		params  tasks.TaskParams
		err     *shared.ServiceError
	}{
		{
			name: "valid params",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				paramsMatcher := mock.MatchedBy(func(t tasks.Task) bool {
					return t.Title == title && t.DueDate.Equal(dueDate) &&
						t.Status == tasks.Pending && t.Priority == tasks.Low
				})
				repo.EXPECT().SaveTask(mock.Anything, paramsMatcher).Return(nil)
			}),
			params: params,
		},
		{
			name:    "invalid due date",
			service: newTestService(t, nil),
			params: tasks.TaskParams{
				Title:    title,
				DueDate:  time.Now(),
				Status:   tasks.Pending,
				Priority: tasks.Low,
			},
			err: shared.NewServiceError(tasks.ErrInvalidDueDate, ""),
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().SaveTask(mock.Anything, mock.Anything).Return(unexpectedErr)
			}),
			params: params,
			err:    shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.service.CreateTask(t.Context(), c.params); err != nil {
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

func TestServiceFindTask(t *testing.T) {
	var filter tasks.TasksFilter
	unexpectedErr := errors.New("unexpected error")
	now := time.Now()
	task, tErr := tasks.NewTask(
		tasks.NewTaskId(),
		"title",
		nil,
		tasks.Pending,
		tasks.Low,
		now.Add(time.Hour),
		now,
		now,
	)
	if tErr != nil {
		t.Fatal("failed to prepare task")
	}
	cases := []struct {
		name    string
		service *tasks.Service
		filter  tasks.TasksFilter
		tasks   []tasks.Task
		err     *shared.ServiceError
	}{
		{
			name: "empty filter",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().FindTasks(mock.Anything, filter).Return([]tasks.Task{task}, nil)
			}),
			tasks: []tasks.Task{task},
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().FindTasks(mock.Anything, mock.Anything).Return(nil, unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tasks, err := c.service.FindTasks(t.Context(), c.filter)
			if err != nil {
				if c.err == nil ||
					!errors.Is(err.Err, c.err.Err) ||
					err.Expected != c.err.Expected ||
					(c.err.Msg != "" && err.Msg != c.err.Msg) {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !reflect.DeepEqual(c.tasks, tasks) {
				t.Fatalf("expected tasks %v, but got %v", c.tasks, tasks)
			}
		})
	}
}

func TestServiceUpdateTaskById(t *testing.T) {
	taskId := tasks.NewTaskId()
	params := tasks.TaskParams{
		Title:    "new title",
		Status:   tasks.Pending,
		Priority: tasks.High,
		DueDate:  time.Now().Add(time.Hour),
	}
	unexpectedErr := errors.New("unexpected error")
	cases := []struct {
		name    string
		service *tasks.Service
		taskId  tasks.TaskId
		params  tasks.TaskParams
		err     *shared.ServiceError
	}{
		{
			name: "happy path",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().UpdateTaskById(mock.Anything, taskId, params).Return(nil)
			}),
			taskId: taskId,
			params: params,
		},
		{
			name: "task not found",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().UpdateTaskById(mock.Anything, mock.Anything, mock.Anything).Return(tasks.ErrTaskNotFound)
			}),
			err: shared.NewServiceError(tasks.ErrTaskNotFound, ""),
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().UpdateTaskById(mock.Anything, mock.Anything, mock.Anything).Return(unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.service.UpdateTaskById(t.Context(), c.taskId, c.params); err != nil {
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

func TestRemoveTaskById(t *testing.T) {
	taskId := tasks.NewTaskId()
	unexpectedErr := errors.New("unexpected err")
	cases := []struct {
		name    string
		service *tasks.Service
		taskId  tasks.TaskId
		err     *shared.ServiceError
	}{
		{
			name: "happy path",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().RemoveTaskById(mock.Anything, taskId).Return(nil)
			}),
			taskId: taskId,
		},
		{
			name: "task not found",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().RemoveTaskById(mock.Anything, mock.Anything).Return(tasks.ErrTaskNotFound)
			}),
			err: shared.NewServiceError(tasks.ErrTaskNotFound, ""),
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().RemoveTaskById(mock.Anything, mock.Anything).Return(unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.service.RemoveTaskById(t.Context(), c.taskId); err != nil {
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

func TestServiceExportTasks(t *testing.T) {
	now := time.Now()
	task, tErr := tasks.NewTask(
		tasks.NewTaskId(),
		"title",
		nil,
		tasks.Pending,
		tasks.Low,
		now.Add(time.Hour),
		now,
		now,
	)
	if tErr != nil {
		t.Fatal("failed to prepare task")
	}
	unexpectedErr := errors.New("unexpected err")
	cases := []struct {
		name    string
		service *tasks.Service
		tasks   []tasks.Task
		err     *shared.ServiceError
	}{
		{
			name: "happy path",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().FindTasks(mock.Anything, tasks.TasksFilter{}).Return([]tasks.Task{task}, nil)
			}),
			tasks: []tasks.Task{task},
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().FindTasks(mock.Anything, mock.Anything).Return(nil, unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tasks, err := c.service.ExportTasks(t.Context())
			if err != nil {
				if c.err == nil ||
					!errors.Is(err.Err, c.err.Err) ||
					err.Expected != c.err.Expected ||
					(c.err.Msg != "" && err.Msg != c.err.Msg) {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !reflect.DeepEqual(c.tasks, tasks) {
				t.Fatalf("expected tasks %v, but got %v", c.tasks, tasks)
			}
		})
	}
}

func TestServiceImportTasks(t *testing.T) {
	now := time.Now()
	task, tErr := tasks.NewTask(
		tasks.NewTaskId(),
		"title",
		nil,
		tasks.Pending,
		tasks.Low,
		now.Add(time.Hour),
		now,
		now,
	)
	if tErr != nil {
		t.Fatal("failed to prepare task")
	}
	ts := []tasks.Task{task}
	unexpectedErr := errors.New("unexpected error")
	cases := []struct {
		name    string
		service *tasks.Service
		tasks   []tasks.Task
		err     *shared.ServiceError
	}{
		{
			name: "happy path",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().SaveTasks(mock.Anything, ts).Return(nil)
			}),
			tasks: ts,
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().SaveTasks(mock.Anything, mock.Anything).Return(unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.service.ImportTasks(t.Context(), c.tasks); err != nil {
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

func TestServicePruneOverdueTasks(t *testing.T) {
	unexpectedErr := errors.New("unexpected error")
	cases := []struct {
		name    string
		service *tasks.Service
		err     *shared.ServiceError
	}{
		{
			name: "happy path",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().
					RemoveOverdueTasksWithDueDateBefore(mock.Anything, mock.AnythingOfType("time.Time")).
					Return(nil)
			}),
		},
		{
			name: "unexpected error",
			service: newTestService(t, func(repo *tasks.MockTasksRepo) {
				repo.EXPECT().
					RemoveOverdueTasksWithDueDateBefore(mock.Anything, mock.Anything).
					Return(unexpectedErr)
			}),
			err: shared.NewUnexpectedError(unexpectedErr, ""),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.service.PruneOverdueTasks(t.Context()); err != nil {
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
