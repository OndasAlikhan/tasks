-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    id         SERIAL PRIMARY KEY,
    task_id    INT NOT NULL REFERENCES tasks(id),
    text       TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS comments;
