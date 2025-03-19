package app

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger/sl"
	"github.com/x0k/skillrock-tasks-service/internal/lib/migrator"
	"github.com/x0k/skillrock-tasks-service/internal/users"
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

	app := fiber.New()

	auth := app.Group("/auth")
	if err := users.NewAuth(
		auth,
		log.With(sl.Module("auth")),
		queries,
		[]byte(cfg.Auth.Secret),
		cfg.Auth.TokenLifetime,
	); err != nil {
		return fmt.Errorf("failed to instantiate auth module: %w", err)
	}

	return app.Listen(cfg.Server.Address)
}
