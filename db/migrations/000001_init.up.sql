CREATE TABLE
  "user" (
    login VARCHAR(255) PRIMARY KEY,
    password_hash BYTEA NOT NULL
  );

CREATE TYPE task_status AS ENUM ('pending', 'in_progress', 'done');

CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high');

CREATE TABLE
  task (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status NOT NULL,
    priority task_priority NOT NULL,
    due_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
  );

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_task_title ON task USING gin (title gin_trgm_ops);
CREATE INDEX idx_task_status ON task (status);
CREATE INDEX idx_task_priority ON task (priority);
CREATE INDEX idx_task_due_date ON task (due_date);
