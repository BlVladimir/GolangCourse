package usecase

import (
	"context"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type MovieSercher interface {
	IndexMovie(ctx context.Context, movie domain.Movie) error
	Search(ctx context.Context, words []string) ([]int, error)
	Clear(ctx context.Context) error
	GetCached(ctx context.Context, query string) (movies []domain.Movie, ok bool, err error)
	SetCached(ctx context.Context, query string, movies []domain.Movie) error
}
