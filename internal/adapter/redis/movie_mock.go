package redis

import (
	"context"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
)

type MovieSercherMock struct{}

func (searcher MovieSercherMock) IndexMovie(ctx context.Context, movie domain.Movie) error {
	return nil
}

func (searcher MovieSercherMock) Search(ctx context.Context, words []string) ([]int, error) {
	return []int{}, nil
}

func (searcher MovieSercherMock) Clear(ctx context.Context) error {
	return nil
}
func (searcher MovieSercherMock) GetCached(ctx context.Context, query string) (movies []domain.Movie, ok bool, err error) {
	return []domain.Movie{}, true, nil
}

func (searcher MovieSercherMock) SetCached(ctx context.Context, query string, movies []domain.Movie) error {
	return nil
}

var _ usecase.MovieSercher = MovieSercherMock{}
