package usecase

import (
	"context"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type Repository interface {
	GetMovie(ctx context.Context, id int) (domain.Movie, error)
	GetMoviesByIDs(ctx context.Context, ids []int) ([]domain.Movie, error)
	LastMovie(ctx context.Context) (domain.Movie, error)
	ExistIDs(ctx context.Context) (map[int]struct{}, error)
	UpsertMovie(ctx context.Context, movie domain.Movie) error
}
