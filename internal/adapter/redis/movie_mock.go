package redis

import (
	"context"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
)

type MovieSearcherMock struct{}

func (searcher MovieSearcherMock) IndexMovie(ctx context.Context, movie domain.Movie) error {
	return nil
}

func (searcher MovieSearcherMock) Search(ctx context.Context, words []string) ([]int, error) {
	return []int{}, nil
}

func (searcher MovieSearcherMock) Clear(ctx context.Context) error {
	return nil
}
func (searcher MovieSearcherMock) GetCached(ctx context.Context, query string) (movies []domain.Movie, ok bool, err error) {
	return []domain.Movie{}, true, nil
}

func (searcher MovieSearcherMock) SetCached(ctx context.Context, query string, movies []domain.Movie) error {
	return nil
}

var _ usecase.MovieSearcher = MovieSearcherMock{}
