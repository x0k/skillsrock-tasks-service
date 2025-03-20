package tests

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tRedis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
	"github.com/x0k/skillrock-tasks-service/internal/lib/migrator"

	// migration tools
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func setupPgxPool(t testing.TB, log *slog.Logger) *pgxpool.Pool {
	pgContainer, err := postgres.Run(
		t.Context(),
		"postgres:17.4-alpine3.21",
		postgres.WithDatabase("tasks"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("admin"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	testcontainers.CleanupContainer(t, pgContainer)

	uri, err := pgContainer.ConnectionString(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	mg := migrator.New(
		log,
		strings.Replace(uri, "postgres://", "pgx5://", 1),
		"file://../../db/migrations",
	)
	if err := mg.Migrate(t.Context()); err != nil {
		t.Fatal(err)
	}
	conn, err := pgxpool.New(t.Context(), uri)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(conn.Close)
	return conn
}

func setupRedisClient(t testing.TB, log *slog.Logger) *redis.Client {
	redisContainer, err := tRedis.Run(
		t.Context(),
		"redis:7.4.2-alpine3.21",
	)
	if err != nil {
		t.Fatal(err)
	}
	testcontainers.CleanupContainer(t, redisContainer)

	uri, err := redisContainer.ConnectionString(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	opts, err := redis.ParseURL(uri)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(opts)
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			log.LogAttrs(t.Context(), slog.LevelError, "failed to close redis client", sl.Err(err))
		}
	})
	return client
}
