package tests

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/x0k/skillrock-tasks-service/internal/analytics"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
)

func newAnalyticsServer(t *testing.T) (*httptest.Server, *analytics.Controller) {
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
	red := setupRedisClient(t, log.Logger)
	app := fiber.New()
	tasksRepo := tasks.NewRepo(
		log,
		pool,
		db.New(pool),
	)
	c := analytics.NewController(
		app,
		log,
		analytics.NewService(
			log,
			tasksRepo,
			analytics.NewRepo(
				log,
				red,
			),
		),
	)
	return httptest.NewServer(adaptor.FiberApp(app)), c
}

func TestReport(t *testing.T) {
	server, c := newAnalyticsServer(t)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)
	e.GET("/").Expect().Status(http.StatusNotFound)

	c.GenerateReport(t.Context())

	var actual analytics.ReportDTO
	e.GET("/").Expect().Status(http.StatusOK).JSON().
		Object().Decode(&actual)

	expected := analytics.ReportDTO{
		PendingTasksCount:           2,
		InProgressTasksCount:        2,
		DoneTasksCount:              1,
		AverageCompletionTimeInDays: "1.00",
		AmountOfCompletedTasks:      0,
		AmountOfOverdueTasks:        0,
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}
