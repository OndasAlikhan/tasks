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
	"time"
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

	go worker(ctx, usersRepo)

	<-ctx.Done()
	cleanUp(db)
	slog.Info("shutting down")
}

func cleanUp(db *pgxpool.Pool) {
	db.Close()
}

func worker(ctx context.Context, usersRepo *repos.UsersRepo) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

		case <-ctx.Done():
			return
		}
	}
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

	// creating task
	task := models.Task{
		UserID:      resultUser.ID,
		Title:       "Wash dishes",
		Description: "You need to wash dishes now",
		Status:      "ready_to_review",
	}
	taskID, err := tasksRepo.Insert(ctx, task)
	if err != nil {
		return err
	}
	task, err = tasksRepo.Find(ctx, taskID)
	if err != nil {
		return err
	}
	slog.Info("created task", "id", taskID, "task", task)

	// created reviewed task
	task2 := models.Task{
		UserID:      resultUser.ID,
		Title:       "Fold clothes",
		Description: "You need to fold clothes now",
		Status:      "review",
	}
	taskID2, err := tasksRepo.Insert(ctx, task2)
	if err != nil {
		return err
	}
	task2, err = tasksRepo.Find(ctx, taskID2)
	if err != nil {
		return err
	}
	slog.Info("created reviewed task", "id", taskID2, "task", task2)

	// creating comment
	comment := models.Comment{
		TaskID: taskID,
		Text:   "dishes should be cleaned with Fairy",
	}
	commentID, err := commentsRepo.Insert(ctx, comment)
	if err != nil {
		return err
	}
	comment, err = commentsRepo.Find(ctx, commentID)
	slog.Info("created comment", "id", commentID, "comment", comment)

	// locking and updating
	err = tasksRepo.LockUpdateStatus(ctx)
	if err != nil {
		return err
	}
	slog.Info("locked and updated tasks in transaction")

	return nil
}
