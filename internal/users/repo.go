package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
)

type repo struct {
	log     *logger.Logger
	queries *db.Queries
}

func newRepo(log *logger.Logger, queries *db.Queries) *repo {
	return &repo{log, queries}
}

func (r *repo) SaveUser(ctx context.Context, user *user) error {
	if err := r.queries.InsertUser(ctx, db.InsertUserParams{
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrLoginIsTaken
		}
		return err
	}
	return nil
}

func (r *repo) UserByLogin(ctx context.Context, login string) (*user, error) {
	u, err := r.queries.UserById(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return newUser(login, u.PasswordHash), nil
}
