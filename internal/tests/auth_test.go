package tests

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/x0k/skillrock-tasks-service/internal/auth"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
)

func newAuthServer(t *testing.T) *httptest.Server {
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
	app := fiber.New()
	auth.NewController(
		app,
		log,
		auth.NewService(
			log,
			[]byte("secret"),
			time.Hour,
			auth.NewRepo(
				log,
				db.New(pool),
			),
		),
	)
	return httptest.NewServer(adaptor.FiberApp(app))
}

func TestRegistrationAndLogin(t *testing.T) {
	server := newAuthServer(t)
	defer server.Close()

	credentials := auth.Credentials{
		Login:    "login",
		Password: "password",
	}

	e := httpexpect.Default(t, server.URL)
	e.POST("/register").
		WithJSON(credentials).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().Value("access_token").String()

	e.POST("/login").
		WithJSON(credentials).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("access_token").String()
}

func TestDoubleRegistration(t *testing.T) {
	server := newAuthServer(t)
	defer server.Close()

	credentials := auth.Credentials{
		Login:    "login",
		Password: "password",
	}

	e := httpexpect.Default(t, server.URL)
	e.POST("/register").
		WithJSON(credentials).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().Value("access_token").String()

	e.POST("/register").
		WithJSON(credentials).
		Expect().
		Status(http.StatusConflict)
}

func TestUnregisteredUser(t *testing.T) {
	server := newAuthServer(t)
	defer server.Close()

	credentials := auth.Credentials{
		Login:    "login",
		Password: "password",
	}

	e := httpexpect.Default(t, server.URL)
	e.POST("/login").
		WithJSON(credentials).
		Expect().
		Status(http.StatusUnauthorized)
}

func TestInvalidPassword(t *testing.T) {
	server := newAuthServer(t)
	defer server.Close()

	credentials := auth.Credentials{
		Login:    "login",
		Password: "password",
	}

	e := httpexpect.Default(t, server.URL)
	e.POST("/register").
		WithJSON(credentials).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().Value("access_token").String()

	credentials.Password = "passsword"

	e.POST("/login").
		WithJSON(credentials).
		Expect().
		Status(http.StatusUnauthorized)
}
