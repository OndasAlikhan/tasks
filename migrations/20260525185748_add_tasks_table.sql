-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(id),
    title       VARCHAR(255) NOT NULL,
    status      VARCHAR(50) NOT NULL,
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS tasks;
