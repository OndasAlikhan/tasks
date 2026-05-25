package repos

import "github.com/jackc/pgx/v5/pgxpool"

type CommentsRepo struct {
	db *pgxpool.Pool
}

func NewCommentsRepo(db *pgxpool.Pool) *CommentsRepo {
	return &CommentsRepo{db: db}
}
