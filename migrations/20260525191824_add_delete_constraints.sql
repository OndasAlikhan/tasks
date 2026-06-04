-- +goose Up
ALTER TABLE tasks
    DROP CONSTRAINT tasks_user_id_fkey,
    ADD CONSTRAINT tasks_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE comments
    DROP CONSTRAINT comments_task_id_fkey,
    ADD CONSTRAINT comments_task_id_fkey
        FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE tasks
    DROP CONSTRAINT tasks_user_id_fkey,
    ADD CONSTRAINT tasks_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE comments
    DROP CONSTRAINT comments_task_id_fkey,
    ADD CONSTRAINT comments_task_id_fkey
        FOREIGN KEY (task_id) REFERENCES tasks(id);
