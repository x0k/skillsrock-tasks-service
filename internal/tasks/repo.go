package tasks

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/x0k/skillrock-tasks-service/internal/lib/db"
	"github.com/x0k/skillrock-tasks-service/internal/lib/logger"
)

type Repo struct {
	log     *logger.Logger
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewRepo(
	log *logger.Logger,
	pool *pgxpool.Pool,
	queries *db.Queries,
) *Repo {
	return &Repo{log, pool, queries}
}

func (r *Repo) SaveTask(ctx context.Context, task Task) error {
	return r.queries.InsertTask(ctx, db.InsertTaskParams{
		ID: pgtype.UUID{
			Bytes: task.Id,
			Valid: true,
		},
		Title:       task.Title,
		Description: r.descriptionToPg(task.Description),
		Status:      db.TaskStatus(task.Status),
		Priority:    db.TaskPriority(task.Priority),
		DueDate: pgtype.Date{
			Time:  task.DueDate,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  task.CreatedAt.UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  task.UpdatedAt.UTC(),
			Valid: true,
		},
	})
}

func (r *Repo) UpdateTaskById(ctx context.Context, id TaskId, params TaskParams) error {
	rowsAffected, err := r.queries.UpdateTask(ctx, db.UpdateTaskParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		Title:       params.Title,
		Description: r.descriptionToPg(params.Description),
		Status:      db.TaskStatus(params.Status),
		Priority:    db.TaskPriority(params.Priority),
		DueDate: pgtype.Date{
			Time:  params.DueDate,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *Repo) RemoveTaskById(ctx context.Context, id TaskId) error {
	rowsAffected, err := r.queries.DeleteTask(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *Repo) SaveTasks(ctx context.Context, tasks []Task) error {
	if len(tasks) == 0 {
		return nil
	}
	q := strings.Builder{}
	q.WriteString(`INSERT INTO task
(id, title, description, status, priority, due_date, created_at, updated_at)
VALUES `)
	var args []any
	push := func(arg any) {
		args = append(args, arg)
		q.WriteByte('$')
		q.WriteString(strconv.Itoa(len(args)))
	}
	for i, t := range tasks {
		if i > 0 {
			q.WriteString(", ")
		}
		q.WriteByte('(')
		push(pgtype.UUID{
			Bytes: t.Id,
			Valid: true,
		})
		q.WriteByte(',')
		push(t.Title)
		q.WriteByte(',')
		push(r.descriptionToPg(t.Description))
		q.WriteByte(',')
		push(t.Status)
		q.WriteByte(',')
		push(t.Priority)
		q.WriteByte(',')
		push(pgtype.Date{
			Time:  t.DueDate,
			Valid: true,
		})
		q.WriteByte(',')
		push(pgtype.Timestamp{
			Time:  t.CreatedAt.UTC(),
			Valid: true,
		})
		q.WriteByte(',')
		push(pgtype.Timestamp{
			Time:  t.UpdatedAt.UTC(),
			Valid: true,
		})
		q.WriteByte(')')
	}
	q.WriteByte(';')
	_, err := r.pool.Exec(ctx, q.String(), args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrTaskIdsConflict
		}
	}
	return err
}

func (r *Repo) FindTasks(ctx context.Context, f TasksFilter) ([]Task, error) {
	q := strings.Builder{}
	q.WriteString(`SELECT id, title, description, status, priority, due_date, created_at, updated_at FROM task`)
	var args []any
	push := func(arg any) {
		args = append(args, arg)
		q.WriteByte('$')
		q.WriteString(strconv.Itoa(len(args)))
	}
	isFirst := true
	prepare := func() {
		if isFirst {
			q.WriteString(" WHERE ")
			isFirst = false
		} else {
			q.WriteString(" AND ")
		}
	}
	if !f.IsEmpty() {
		if f.Title != nil {
			prepare()
			q.WriteString("title ILIKE ")
			push("%" + *f.Title + "%")
		}
		if f.Status != nil {
			prepare()
			q.WriteString("status = ")
			push(*f.Status)
		}
		if f.Priority != nil {
			prepare()
			q.WriteString("priority = ")
			push(*f.Priority)
		}
		if f.DueAfter != nil {
			prepare()
			q.WriteString("due_date > ")
			push(pgtype.Date{
				Time:  *f.DueAfter,
				Valid: true,
			})
		}
		if f.DueBefore != nil {
			prepare()
			q.WriteString("due_date < ")
			push(pgtype.Date{
				Time:  *f.DueBefore,
				Valid: true,
			})
		}
	}
	q.WriteByte(';')
	rows, err := r.pool.Query(ctx, q.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var row db.Task
		if err := rows.Scan(
			&row.ID,
			&row.Title,
			&row.Description,
			&row.Status,
			&row.Priority,
			&row.DueDate,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, err
		}
		task, err := NewTask(
			row.ID.Bytes,
			row.Title,
			r.descriptionFromPg(row.Description),
			Status(row.Status),
			Priority(row.Priority),
			row.DueDate.Time,
			row.CreatedAt.Time,
			row.UpdatedAt.Time,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, task)
	}
	return items, rows.Err()
}

func (r *Repo) AllTasks(ctx context.Context) ([]Task, error) {
	rows, err := r.queries.AllTasks(ctx)
	if err != nil {
		return nil, err
	}
	tasks := make([]Task, len(rows))
	for i, row := range rows {
		if tasks[i], err = NewTask(
			row.ID.Bytes,
			row.Title,
			r.descriptionFromPg(row.Description),
			Status(row.Status),
			Priority(row.Priority),
			row.DueDate.Time,
			row.CreatedAt.Time,
			row.UpdatedAt.Time,
		); err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

func (r *Repo) TasksCountByStatus(ctx context.Context) (map[Status]int64, error) {
	rows, err := r.queries.CountTasksByStatus(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[Status]int64, len(rows))
	for _, r := range rows {
		m[Status(r.Status)] = r.TasksCount
	}
	return m, nil
}

func (r *Repo) AverageCompletionTime(ctx context.Context) (float64, error) {
	return r.queries.AverageTaskCompletionTime(ctx)
}

func (r *Repo) CountCompletedAndOverdueTasks(ctx context.Context, date time.Time) (int64, int64, error) {
	row, err := r.queries.CountCompletedAndOverdueTasks(ctx, pgtype.Timestamp{
		Time:  date.UTC(),
		Valid: true,
	})
	return row.CompletedCount, row.OverdueCount, err
}

func (r *Repo) RemoveOverdueTasksWithDueDateBefore(ctx context.Context, date time.Time) error {
	return r.queries.DeleteOverdueTasks(ctx, pgtype.Date{
		Time:  date,
		Valid: true,
	})
}

func (r *Repo) descriptionToPg(d *string) pgtype.Text {
	var t pgtype.Text
	if d != nil {
		t.String = *d
		t.Valid = true
	}
	return t
}

func (r *Repo) descriptionFromPg(t pgtype.Text) *string {
	if t.Valid {
		return &t.String
	}
	return nil
}
