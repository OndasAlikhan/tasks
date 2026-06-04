package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"tasks/internal/models"
)

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{db: db}
}

func (u *UsersRepo) Insert(ctx context.Context, user models.User) (int, error) {
	sql := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := u.db.QueryRow(ctx, sql, user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error in UsersRepo.Insert: %w", err)
	}
	return id, nil
}

func (u *UsersRepo) Find(ctx context.Context, id int) (models.User, error) {
	sql := `
		select id, name, email, created_at from users where id = $1
	`

	var user models.User
	err := u.db.QueryRow(ctx, sql, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("error in UsersRepo.Find: %w", err)
	}

	return user, nil
}

func (u *UsersRepo) Update(ctx context.Context, user models.User) error {
	sql := `
		UPDATE users
		SET name = $1,
		    email = $2
		WHERE id = $3
	`

	_, err := u.db.Exec(ctx, sql, user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("error in UsersRepo.Update: %w", err)
	}
	return nil
}

func (u *UsersRepo) ShowUsersWithTheirTasksComments(ctx context.Context) ([]models.UserAggregated, error) {
	sql := `
		select u.id, u.name, t.id, t.title, t.description, t.status, t.updated_at, t.created_at, c.id, c.text, c.created_at  
		from users u
		join tasks t on u.id = t.user_id
		left join comments c on t.id = c.task_id
	`

	rows, err := u.db.Query(ctx, sql)
	if err != nil {
		return []models.UserAggregated{}, fmt.Errorf("error query UsersRepo.ShowUsersWithTheirTasksComments: %w", err)
	}

	defer rows.Close()

	users := make([]models.UserAggregated, 0)
	for rows.Next() {
		var (
			commentID        *int
			commentText      *string
			commentCreatedAt *time.Time
		)
		user := models.UserAggregated{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Task.ID,
			&user.Task.Title,
			&user.Task.Description,
			&user.Task.Status,
			&user.Task.UpdatedAt,
			&user.Task.CreatedAt,
			&commentID,
			&commentText,
			&commentCreatedAt); err != nil {
			return []models.UserAggregated{}, fmt.Errorf("error scan UsersRepo.ShowUsersWithTheirTasksComments: %w", err)
		}
		if commentID != nil {
			user.Task.Comment = models.Comment{
				ID:        *commentID,
				Text:      *commentText,
				CreatedAt: *commentCreatedAt,
			}
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return []models.UserAggregated{}, fmt.Errorf("error rows UsersRepo.ShowUsersWithTheirTasksComments: %w", err)
	}
	return users, nil
}
