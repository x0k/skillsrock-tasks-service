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
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    updated_at DATE NOT NULL DEFAULT CURRENT_DATE
  );