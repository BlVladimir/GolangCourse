package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/BlVladimir/GolangCourse/internal/adapter/postgres/sqlc"
	"github.com/BlVladimir/GolangCourse/internal/domain"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Repository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func (repository *Repository) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	row, err := repository.queries.GetMovie(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Movie{}, domain.ErrNotFound
		}
		return domain.Movie{}, fmt.Errorf("postgres: get movie %d: %w", id, err)
	}
	return toDomain(row), nil
}

func (repository *Repository) GetMoviesRange(ctx context.Context, from, to int) ([]domain.Movie, error) {
	rows, err := repository.queries.GetMoviesRange(ctx, sqlc.GetMoviesRangeParams{
		FromID: int64(from),
		ToID:   int64(to),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: get movies [%d..%d]: %w", from, to, err)
	}
	movies := make([]domain.Movie, len(rows))
	for i, row := range rows {
		movies[i] = toDomain(row)
	}
	return movies, nil
}

func (repository *Repository) GetMoviesByIDs(ctx context.Context, ids []int) ([]domain.Movie, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	ids64 := make([]int64, len(ids))
	for i, id := range ids {
		ids64[i] = int64(id)
	}
	rows, err := repository.queries.GetMoviesByIDs(ctx, ids64)
	if err != nil {
		return nil, fmt.Errorf("postgres: get movies by ids: %W", err)
	}

	movies := make([]domain.Movie, len(rows))
	for i, row := range rows {
		movies[i] = toDomain(row)
	}
	return movies, nil
}

func (repository *Repository) LastMovie(ctx context.Context) (domain.Movie, error) {
	row, err := repository.queries.GetLastMovie(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Movie{}, domain.ErrNotFound
		}
		return domain.Movie{}, fmt.Errorf("postgres: last movie: %w", err)
	}
	return toDomain(row), nil
}

func (repository *Repository) ExistIDs(ctx context.Context) (map[int]struct{}, error) {
	ids, err := repository.queries.ExistIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres: exist ids: %w", err)
	}
	result := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		result[int(id)] = struct{}{}
	}
	return result, nil
}

func (repository *Repository) UpsertMovie(ctx context.Context, movie domain.Movie) error {
	err := repository.queries.UpsertMovie(ctx, sqlc.UpsertMovieParams{
		ID:       int64(movie.ID),
		Title:    movie.Title,
		Year:     int32(movie.Year),
		Director: movie.Director,
	})
	if err != nil {
		return fmt.Errorf("postgres: upsert movie %d: %w", movie.ID, err)
	}
	return nil
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

const migrationsDir = "migrations"

func (r *Repository) Migrate(ctx context.Context) error {
	sqlDB := stdlib.OpenDBFromPool(r.pool)
	defer sqlDB.Close()

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return err
	}
	return goose.UpContext(ctx, sqlDB, migrationsDir)
}

func toDomain(movie sqlc.Movie) domain.Movie {
	return domain.Movie{
		ID:       int(movie.ID),
		Title:    movie.Title,
		Year:     int(movie.Year),
		Director: movie.Director,
	}
}

var _ usecase.Repository = (*Repository)(nil)
