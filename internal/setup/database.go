package setup

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"time"
)

func SetupDB(ctx context.Context, dbURL string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		slog.Error("unable to parse db config", "err", err)
		os.Exit(1)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("unable to create connection pool", "err", err)
		os.Exit(1)
	}
	return dbpool
}
