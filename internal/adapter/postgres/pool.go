package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/adapter/postgres/sqlc"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(ctx context.Context, dsn string) (*Repository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}

	return &Repository{db: db, queries: sqlc.New(db)}, nil
}
