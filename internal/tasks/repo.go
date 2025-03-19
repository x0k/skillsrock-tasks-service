package tasks

import (
	"context"
	"strconv"
	"strings"

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
		CreatedAt: pgtype.Date{
			Time:  task.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Date{
			Time:  task.UpdatedAt,
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
		push(pgtype.Date{
			Time:  t.CreatedAt,
			Valid: true,
		})
		q.WriteByte(',')
		push(pgtype.Date{
			Time:  t.UpdatedAt,
			Valid: true,
		})
		q.WriteByte(')')
	}
	q.WriteByte(';')
	_, err := r.pool.Exec(ctx, q.String(), args...)
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
			push(*f.Title)
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
			q.WriteString("due_date >= ")
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
		var i db.Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.Priority,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, Task{
			Id:          i.ID.Bytes,
			Title:       i.Title,
			Description: r.descriptionFromPg(i.Description),
			Status:      Status(i.Status),
			Priority:    Priority(i.Priority),
			DueDate:     i.DueDate.Time,
			CreatedAt:   i.CreatedAt.Time,
			UpdatedAt:   i.UpdatedAt.Time,
		})
	}
	return items, rows.Err()
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
