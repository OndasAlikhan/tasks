package repos

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tasks/internal/models"
)

type TasksRepo struct {
	db *pgxpool.Pool
}

func NewTasksRepo(db *pgxpool.Pool) *TasksRepo {
	return &TasksRepo{db: db}
}

func (t *TasksRepo) Insert(ctx context.Context, task models.Task) (int, error) {
	sql := `
		INSERT INTO tasks (user_id, title, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := t.db.QueryRow(ctx, sql, task.UserID, task.Title, task.Description, task.Status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (t *TasksRepo) Find(ctx context.Context, id int) (models.Task, error) {
	sql := `
		SELECT id, user_id, title, description, status, updated_at, created_at
		FROM tasks
		WHERE id = $1
	`

	var task models.Task
	err := t.db.QueryRow(ctx, sql, id).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description,
		&task.Status, &task.UpdatedAt, &task.CreatedAt,
	)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (t *TasksRepo) Update(ctx context.Context, task models.Task) error {
	sql := `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    status = $3,
		    updated_at = NOW()
		WHERE id = $4
	`

	_, err := t.db.Exec(ctx, sql, task.Title, task.Description, task.Status, task.ID)
	return err
}

func (t *TasksRepo) LockUpdateStatus(ctx context.Context) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	selectSQL := `
		SELECT id FROM tasks
		WHERE status = 'ready_to_review'
		FOR UPDATE
	`
	rows, err := tx.Query(ctx, selectSQL)
	if err != nil {
		return err
	}
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	rows.Close()

	updateSQL := `
		UPDATE tasks
		SET status = 'review',
		    updated_at = NOW()
		WHERE id = ANY($1)
	`
	_, err = tx.Exec(ctx, updateSQL, ids)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
