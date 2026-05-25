-- +goose Up
CREATE INDEX idx_tasks_status ON tasks(status);

-- +goose Down
DROP INDEX idx_tasks_status;
