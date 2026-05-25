package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tasks/internal/models"
	"tasks/internal/repos"
	"tasks/internal/setup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbURL := os.Getenv("DATABASE_URL")
	db := setup.SetupDB(ctx, dbURL)
	slog.Info("connected to db")

	usersRepo := repos.NewUsersRepo(db)
	tasksRepo := repos.NewTasksRepo(db)
	commentsRepo := repos.NewCommentsRepo(db)

	_, err := db.Exec(context.Background(), "TRUNCATE TABLE comments, tasks, users RESTART IDENTITY CASCADE")
	if err != nil {
		slog.Error("failed to truncate tables", "err", err)
		os.Exit(1)
	}

	_ = businessLogic(usersRepo, tasksRepo, commentsRepo)

	<-ctx.Done()
	cleanUp(db)
	slog.Info("shutting down")
}

func cleanUp(db *pgxpool.Pool) {
	db.Close()
}

func businessLogic(usersRepo *repos.UsersRepo, tasksRepo *repos.TasksRepo, commentsRepo *repos.CommentsRepo) (err error) {
	defer func() {
		if err != nil {
			slog.Error("error", "err", err)
		}
	}()

	ctx := context.Background()

	// create new user
	user := models.User{
		Name:  "Alikhan",
		Email: "alikhan@gmail.com",
	}
	userID, err := usersRepo.Insert(ctx, user)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("created user with ID %d", userID))

	resultUser, err := usersRepo.Find(ctx, userID)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("found user by ID %d", userID), "user", resultUser)

	// update user
	resultUser.Email = "SOMEWORD@gmail.com"
	err = usersRepo.Update(ctx, resultUser)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("updated user %d's email", resultUser.ID))

	resultUser, err = usersRepo.Find(ctx, resultUser.ID)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("found user by ID %d", resultUser.ID), "user", resultUser)

	return nil
}
