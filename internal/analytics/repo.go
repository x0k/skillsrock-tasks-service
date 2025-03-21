package analytics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
)

type Repo struct {
	log   *logger.Logger
	redis *redis.Client
}

func NewRepo(
	log *logger.Logger,
	redis *redis.Client,
) *Repo {
	return &Repo{log, redis}
}

func (r *Repo) SaveReport(ctx context.Context, report Report) error {
	bytes, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	if err := r.redis.Set(ctx, "report", bytes, 0).Err(); err != nil {
		return fmt.Errorf("failed to persist report %w", err)
	}
	return nil
}

func (r *Repo) Report(ctx context.Context) (Report, error) {
	var report Report
	val, err := r.redis.Get(ctx, "report").Result()
	if errors.Is(err, redis.Nil) {
		return report, ErrReportNotFound
	}
	if err != nil {
		return report, fmt.Errorf("failed to retrieve report %w", err)
	}
	if err := json.Unmarshal([]byte(val), &report); err != nil {
		return report, fmt.Errorf("failed to unmarshal report %w", err)
	}
	return report, nil
}
