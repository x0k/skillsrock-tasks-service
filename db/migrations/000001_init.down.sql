DROP INDEX IF EXISTS idx_task_due_date;
DROP INDEX IF EXISTS idx_task_priority;
DROP INDEX IF EXISTS idx_task_status;
DROP INDEX IF EXISTS idx_task_title;

DROP EXTENSION IF EXISTS pg_trgm;

DROP TABLE IF EXISTS task;

DROP TYPE IF EXISTS task_priority;
DROP TYPE IF EXISTS task_status;

DROP TABLE IF EXISTS "user";