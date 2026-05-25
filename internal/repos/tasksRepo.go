package repos

import "github.com/jackc/pgx/v5/pgxpool"

type TasksRepo struct {
	db *pgxpool.Pool
}

func NewTasksRepo(db *pgxpool.Pool) *TasksRepo {
	return &TasksRepo{db: db}
}
