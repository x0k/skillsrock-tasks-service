-- name: UserById :one
SELECT * FROM user WHERE login = $1;

-- name: InsertUser :exec
INSERT INTO user (login, password_hash) VALUES ($1, $2);

-- name: InsertTask :exec
INSERT INTO task (title, description, status, priority, due_date);

-- name: UpdateTask :exec
UPDATE task SET
  title = $2,
  description = $3,
  status = $4,
  priority = $5,
  due_date = $6,
  updated_at = NOW()
WHERE
  task.id = $1;

-- name: DeleteTask :exec
DELETE FROM task WHERE task.id = $1;
