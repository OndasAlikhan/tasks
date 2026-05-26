package repos

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tasks/internal/models"
)

type CommentsRepo struct {
	db *pgxpool.Pool
}

func NewCommentsRepo(db *pgxpool.Pool) *CommentsRepo {
	return &CommentsRepo{db: db}
}

func (c *CommentsRepo) Insert(ctx context.Context, comment models.Comment) (int, error) {
	sql := `
		INSERT INTO comments (task_id, text)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := c.db.QueryRow(ctx, sql, comment.TaskID, comment.Text).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CommentsRepo) Find(ctx context.Context, id int) (models.Comment, error) {
	sql := `
		SELECT id, task_id, text, created_at
		FROM comments
		WHERE id = $1
	`

	var comment models.Comment
	err := c.db.QueryRow(ctx, sql, id).Scan(
		&comment.ID, &comment.TaskID, &comment.Text, &comment.CreatedAt,
	)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func (c *CommentsRepo) Update(ctx context.Context, comment models.Comment) error {
	sql := `
		UPDATE comments
		SET text = $1
		WHERE id = $2
	`

	_, err := c.db.Exec(ctx, sql, comment.Text, comment.ID)
	return err
}
