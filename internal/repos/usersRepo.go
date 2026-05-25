package repos

import (
	"context"
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
		return 0, err
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
		return models.User{}, err
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
	return err
}
