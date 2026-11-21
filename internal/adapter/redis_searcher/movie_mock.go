package redis_searcher

import (
	"context"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
)

var testMovies = []domain.Movie{
	domain.Movie{
		ID:       3,
		Title:    "Kings of the Sun",
		Year:     1963,
		Director: "J. Lee Thompson",
	},
	domain.Movie{
		ID:       4,
		Title:    "Bandits in Milan",
		Year:     1968,
		Director: "Carlo Lizzani",
	},
}

type MovieSearcherMock struct{}

func (searcher MovieSearcherMock) IndexMovie(ctx context.Context, movie domain.Movie) error {
	return nil
}

func (searcher MovieSearcherMock) Search(ctx context.Context, words []string) ([]int, error) {
	return []int{3, 4}, nil
}

func (searcher MovieSearcherMock) Clear(ctx context.Context) error {
	return nil
}
func (searcher MovieSearcherMock) GetCached(ctx context.Context, query string) (movies []domain.Movie, ok bool, err error) {
	return testMovies, true, nil
}

func (searcher MovieSearcherMock) SetCached(ctx context.Context, query string, movies []domain.Movie) error {
	return nil
}

var _ usecase.MovieSearcher = MovieSearcherMock{}
