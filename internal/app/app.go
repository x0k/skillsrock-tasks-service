package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/x0k/skillrock-tasks-service/internal/analytics"
	"github.com/x0k/skillrock-tasks-service/internal/auth"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
	"github.com/x0k/skillrock-tasks-service/internal/lib/migrator"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
	tasks_controller "github.com/x0k/skillrock-tasks-service/internal/tasks/controller"

	// migration tools
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(ctx context.Context, cfg *Config, log *logger.Logger) error {
	m := migrator.New(
		log.Logger.With(sl.Component("migrator")),
		strings.Replace(cfg.Postgres.ConnectionURI, "postgres://", "pgx5://", 1),
		cfg.Postgres.MigrationsURI,
	)
	if err := m.Migrate(ctx); err != nil {
		return fmt.Errorf("migrator: %w", err)
	}

	pgxPool, err := pgxpool.New(ctx, cfg.Postgres.ConnectionURI)
	if err != nil {
		return fmt.Errorf("pgx pool: %w", err)
	}
	defer pgxPool.Close()
	queries := db.New(pgxPool)

	redisOpts, err := redis.ParseURL(cfg.Redis.ConnectionURI)
	if err != nil {
		return fmt.Errorf("parse redis url: %w", err)
	}
	redisClient := redis.NewClient(redisOpts)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error(ctx, "failed to close redis client", sl.Err(err))
		}
	}()

	app := fiber.New()

	app.Use(slogfiber.New(log.Logger))
	app.Use(recover.New())

	auth.NewController(
		app.Group("/auth"),
		log.With(sl.Component("auth_controller")),
		auth.NewService(
			log.With(sl.Component("auth_service")),
			[]byte(cfg.Auth.Secret),
			cfg.Auth.TokenLifetime,
			auth.NewRepo(
				log.With(sl.Component("auth_repo")),
				queries,
			),
		),
	)

	authMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Auth.Secret)},
	})

	tasksRepo := tasks.NewRepo(
		log.With(sl.Component("tasks_repo")),
		pgxPool,
		queries,
	)
	tasksGroup := app.Group("/tasks").Use(authMiddleware)
	tasksController := tasks_controller.New(
		tasksGroup,
		log.With(sl.Component("tasks_controller")),
		tasks.NewService(
			log.With(sl.Component("tasks_service")),
			tasksRepo,
		),
	)

	analyticsGroup := app.Group("/analytics").Use(authMiddleware)
	analyticsController := analytics.NewController(
		analyticsGroup,
		log.With(sl.Component("analytics_controller")),
		analytics.NewService(
			log.With(sl.Component("analytics_service")),
			tasksRepo,
			analytics.NewRepo(
				log.With(sl.Component("analytics_repo")),
				redisClient,
			),
		),
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				tasksController.PruneOverdueTasks(ctx)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				analyticsController.GenerateReport(ctx)
			}
		}
	}()

	context.AfterFunc(ctx, func() {
		if err := app.Shutdown(); err != nil {
			log.Error(ctx, "shutdown failed", sl.Err(err))
		}
	})

	if cfg.Metrics.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mux := http.NewServeMux()
			reg := prometheus.NewRegistry()
			reg.MustRegister(
				collectors.NewGoCollector(),
				collectors.NewProcessCollector(
					collectors.ProcessCollectorOpts{},
				),
			)
			mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
			srv := &http.Server{
				Addr:    cfg.Metrics.Address,
				Handler: mux,
			}
			context.AfterFunc(ctx, func() {
				if err := srv.Shutdown(ctx); err != nil {
					log.Error(ctx, "failed to shutdown metrics server", sl.Err(err))
				}
			})
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error(ctx, "failed to stop metrics server", sl.Err(err))
			}
		}()
	}

	err = app.Listen(cfg.Server.Address)
	wg.Wait()
	return err
}
