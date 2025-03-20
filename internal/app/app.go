package app

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/x0k/skillrock-tasks-service/internal/analytics"
	"github.com/x0k/skillrock-tasks-service/internal/auth"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
	"github.com/x0k/skillrock-tasks-service/internal/lib/migrator"
	"github.com/x0k/skillrock-tasks-service/internal/tasks"
	tasks_controller "github.com/x0k/skillrock-tasks-service/internal/tasks/controller"
)

func Run(ctx context.Context, cfg *Config, log *logger.Logger) error {
	m := migrator.New(
		log.Logger.With(sl.Component("migrator")),
		cfg.Postgres.ConnectionURI,
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

	app := fiber.New()

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

	tasksRepo := tasks.NewRepo(
		log.With(sl.Component("tasks_repo")),
		pgxPool,
		queries,
	)
	tasksGroup := app.Group("/tasks")
	tasks_controller.New(
		tasksGroup,
		log.With(sl.Component("tasks_controller")),
		tasks.NewService(
			log.With(sl.Component("tasks_service")),
			tasksRepo,
		),
	)

	analyticsGroup := app.Group("/analytics")
	analytics.NewController(
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

	return app.Listen(cfg.Server.Address)
}
