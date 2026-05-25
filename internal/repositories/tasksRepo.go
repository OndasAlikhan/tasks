package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type TasksRepo struct {
	db *pgxpool.Pool
}
