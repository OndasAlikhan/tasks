package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbURL := os.Getenv("DATABASE_URL")
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("unable to create connection pool", "err", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	slog.Info("connected to db")

	<-ctx.Done()
	slog.Info("shutting down")
}
