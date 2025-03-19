-- name: UserById :one
SELECT * FROM "user" WHERE login = $1;

-- name: InsertUser :exec
INSERT INTO "user" (login, password_hash) VALUES ($1, $2);

-- name: InsertTask :exec
INSERT INTO task (title, description, status, priority, due_date)
VALUES ($1, $2, $3, $4, $5);

-- TODO: Check affected rows
-- name: UpdateTask :exec
UPDATE task SET
  title = $2,
  description = $3,
  status = $4,
  priority = $5,
  due_date = $6,
  updated_at = CURRENT_DATE
WHERE
  task.id = $1 AND task.status != 'done';

-- name: DeleteTask :exec
DELETE FROM task WHERE task.id = $1;

-- name: DeleteOverdueTasks :exec
DELETE FROM task WHERE status != 'done' and due_date < $1;

-- name: CountTasksByStatus :many
SELECT count(*) AS tasks_count, status FROM task GROUP BY status;

-- name: AverageTaskCompletionTime :one
SELECT
    AVG(updated_at - completed_at) AS average_completion_time
FROM
    task
WHERE
    task.status = 'done';

-- name: CountCompletedAndOverdueTasks :one
WITH last_week_task AS (
  SELECT *
  FROM task
  WHERE updated_at >= $1
)
SELECT
  (SELECT count(*) FROM last_week_task WHERE status = 'done') AS completed_count,
  (SELECT count(*) FROM last_week_task WHERE status != 'done' AND due_date < CURRENT_DATE) AS overdue_count;
