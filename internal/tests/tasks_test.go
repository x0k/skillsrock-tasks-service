package tests

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
	tasks_controller "github.com/x0k/skillrock-tasks-service/internal/tasks/controller"
)

const insertTasks = `
INSERT INTO task
  (id, title, description, status, priority, due_date, created_at, updated_at)
VALUES
	('11111111-1111-1111-1111-111111111111', 'Fix login bug',        'Investigate and fix login issue for users.', 'pending',     'high',   '2025-02-02', '2025-02-01', '2025-02-02'),
	('22222222-2222-2222-2222-222222222222', 'Refactor API',         NULL,                                         'in_progress', 'medium', '2025-02-03', '2025-02-02', '2025-02-03'),
  ('33333333-3333-3333-3333-333333333333', 'Write tests',          'Increase test coverage for task module.',    'pending',     'low',    '2025-02-04', '2025-02-03', '2025-02-04'),
  ('44444444-4444-4444-4444-444444444444', 'Update documentation', 'Document new API endpoints.',                'done',        'low',    '2025-02-05', '2025-02-04', '2025-02-05'),
  ('55555555-5555-5555-5555-555555555555', 'Deploy new release',   NULL,                                         'in_progress', 'high',   '2025-02-06', '2025-02-05', '2025-02-06');
`

func newTasksServer(t *testing.T) (*httptest.Server, *tasks_controller.Controller) {
	var buf bytes.Buffer
	log := logger.New(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(buf.String())
		}
	})
	pool := setupPgxPool(t, log.Logger)
	execSql(t, pool, insertTasks)
	app := fiber.New()
	c := tasks_controller.New(
		app,
		log,
		tasks.NewService(
			log,
			tasks.NewRepo(
				log,
				pool,
				db.New(pool),
			),
		),
	)
	return httptest.NewServer(adaptor.FiberApp(app)), c
}

func TestFindTasks(t *testing.T) {
	server, _ := newTasksServer(t)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	e.GET("/").Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(5)

	e.GET("/").WithQuery("status", "pending").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(2)

	e.GET("/").WithQuery("priority", "high").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(2)

	e.GET("/").WithQuery("due_before", "2025-02-04").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(2)

	e.GET("/").WithQuery("due_after", "2025-02-03").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(3)

	e.GET("/").WithQuery("title", "re").
		WithQuery("status", "in_progress").
		WithQuery("due_after", "2025-02-02").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(2)
}

func TestCreateTask(t *testing.T) {
	server, _ := newTasksServer(t)
	defer server.Close()

	now := time.Now()
	dueDate := now.Add(24 * time.Hour).Format(time.DateOnly)

	e := httpexpect.Default(t, server.URL)
	e.POST("/").WithJSON(map[string]string{
		"title":    "foo",
		"status":   "pending",
		"priority": "low",
		"due_date": dueDate,
	}).Expect().Status(http.StatusCreated)

	e.GET("/").WithQuery("title", "foo").
		WithQuery("status", "pending").
		WithQuery("priority", "low").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(1)
}

func TestUpdateTask(t *testing.T) {
	server, _ := newTasksServer(t)
	defer server.Close()

	dueDate := time.Date(2025, 04, 02, 0, 0, 0, 0, time.Local).Format(time.DateOnly)

	e := httpexpect.Default(t, server.URL)
	e.PUT("/11111111-1111-1111-1111-111111111111").WithJSON(map[string]string{
		"title":    "foo",
		"status":   "pending",
		"priority": "low",
		"due_date": dueDate,
	}).Expect().Status(http.StatusNoContent)

	e.GET("/").WithQuery("title", "foo").
		WithQuery("status", "pending").
		WithQuery("priority", "low").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(1)

	e.PUT("/11111111-1111-1111-1111-111111111112").WithJSON(map[string]string{
		"title":    "foo",
		"status":   "pending",
		"priority": "low",
		"due_date": dueDate,
	}).Expect().Status(http.StatusNotFound)

	e.PUT("/44444444-4444-4444-4444-444444444444").WithJSON(map[string]string{
		"title":    "foo",
		"status":   "pending",
		"priority": "low",
		"due_date": dueDate,
	}).Expect().Status(http.StatusNotFound)
}

func TestDeleteTask(t *testing.T) {
	server, _ := newTasksServer(t)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	e.DELETE("/11111111-1111-1111-1111-111111111111").Expect().
		Status(http.StatusNoContent)

	e.DELETE("/11111111-1111-1111-1111-111111111111").Expect().
		Status(http.StatusNotFound)
}

func TestExportTasks(t *testing.T) {
	server, _ := newTasksServer(t)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	e.GET("/export").Expect().
		Status(http.StatusOK).
		JSON().Array().Length().IsEqual(5)
}

func TestImportTasks(t *testing.T) {
	server, _ := newTasksServer(t)
	defer server.Close()

	dto := make([]tasks_controller.TaskDTO, 100)
	now := time.Now()
	nowDate := now.Format(time.RFC3339)
	dueDate := now.Add(time.Hour).Format(time.DateOnly)
	for i := range 100 {
		dto[i] = tasks_controller.TaskDTO{
			Id:        tasks.NewTaskId().String(),
			Title:     strconv.Itoa(i),
			Status:    tasks.Pending.String(),
			Priority:  tasks.Low.String(),
			DueDate:   dueDate,
			CreatedAt: nowDate,
			UpdatedAt: nowDate,
		}
	}

	e := httpexpect.Default(t, server.URL)
	e.POST("/import").WithJSON(dto).
		Expect().Status(http.StatusCreated)

	e.GET("/").Expect().JSON().
		Array().Length().IsEqual(105)

	e.POST("/import").WithJSON(dto).
		Expect().Status(http.StatusConflict)
}

func TestPruneOverdueTasks(t *testing.T) {
	server, c := newTasksServer(t)
	defer server.Close()

	c.PruneOverdueTasks(t.Context())

	e := httpexpect.Default(t, server.URL)
	e.GET("/").Expect().JSON().
		Array().Length().IsEqual(1)

	dto := make([]tasks_controller.TaskDTO, 10)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	dayDuration := 24 * time.Hour
	for i := range 10 {
		date := today.Add(-time.Duration(i) * dayDuration)
		dto[i] = tasks_controller.TaskDTO{
			Id:        tasks.NewTaskId().String(),
			Title:     strconv.Itoa(i),
			Status:    tasks.Pending.String(),
			Priority:  tasks.Low.String(),
			DueDate:   date.Add(time.Hour).Format(time.DateOnly),
			CreatedAt: date.Format(time.RFC3339),
			UpdatedAt: date.Format(time.RFC3339),
		}
	}
	e.POST("/import").WithJSON(dto).
		Expect().Status(http.StatusCreated)

	c.PruneOverdueTasks(t.Context())

	e.GET("/").Expect().JSON().
		Array().Length().IsEqual(9)
}
