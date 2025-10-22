CREATE SCHEMA IF NOT EXISTS gocourse;

CREATE TABLE IF NOT EXISTS gocourse.todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    createdAt TIMESTAMP DEFAULT NOW(),
    updatedAt TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_todos_created_at ON gocourse.todos(createdAt);
CREATE INDEX IF NOT EXISTS idx_todos_completed ON gocourse.todos(completed);

