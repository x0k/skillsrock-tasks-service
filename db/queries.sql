-- name: UserById :one
SELECT * FROM "user" WHERE login = $1;

-- name: InsertUser :exec
INSERT INTO "user" (login, password_hash) VALUES ($1, $2);

-- name: AllTasks :many
SELECT * FROM task;

-- name: InsertTask :exec
INSERT INTO task
  (id, title, description, status, priority, due_date, created_at, updated_at)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateTask :execrows
UPDATE task SET
  title = $2,
  description = $3,
  status = $4,
  priority = $5,
  due_date = $6,
  updated_at = CURRENT_DATE
WHERE
  task.id = $1 AND task.status != 'done';

-- name: DeleteTask :execrows
DELETE FROM task WHERE task.id = $1;

-- name: DeleteOverdueTasks :exec
DELETE FROM task WHERE status != 'done' and due_date < $1;

-- name: CountTasksByStatus :many
SELECT count(*) AS tasks_count, status FROM task GROUP BY status;

-- name: AverageTaskCompletionTime :one
SELECT
  AVG(EXTRACT(EPOCH FROM (updated_at - created_at))) AS average_completion_time
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
