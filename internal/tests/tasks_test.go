package tests

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

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
	('11111111-1111-1111-1111-111111111111', 'Fix login bug',        'Investigate and fix login issue for users.', 'pending',     'high',   '2025-04-02', '2025-04-01', '2025-04-01'),
	('22222222-2222-2222-2222-222222222222', 'Refactor API',         NULL,                                         'in_progress', 'medium', '2025-04-03', '2025-04-02', '2025-04-02'),
  ('33333333-3333-3333-3333-333333333333', 'Write tests',          'Increase test coverage for task module.',    'pending',     'low',    '2025-04-04', '2025-04-03', '2025-04-03'),
  ('44444444-4444-4444-4444-444444444444', 'Update documentation', 'Document new API endpoints.',                'done',        'low',    '2025-04-05', '2025-04-04', '2025-04-04'),
  ('55555555-5555-5555-5555-555555555555', 'Deploy new release',   NULL,                                         'in_progress', 'high',   '2025-04-06', '2025-04-05', '2025-04-05');
`

func newTasksServer(t *testing.T) *httptest.Server {
	var buf bytes.Buffer
	log := logger.New(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(buf)
		}
	})
	pool := setupPgxPool(t, log.Logger)
	execSql(t, pool, insertTasks)
	app := fiber.New()
	tasks_controller.New(
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
	return httptest.NewServer(adaptor.FiberApp(app))
}

func TestFindTasks(t *testing.T) {
	server := newTasksServer(t)
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

	e.GET("/").WithQuery("due_before", "2025-04-04").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(2)

	e.GET("/").WithQuery("due_after", "2025-04-03").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(3)

	e.GET("/").WithQuery("title", "re").
		WithQuery("status", "in_progress").
		WithQuery("due_after", "2025-04-02").
		Expect().Status(http.StatusOK).
		JSON().Array().Length().IsEqual(2)
}
