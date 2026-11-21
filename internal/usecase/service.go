package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"golang.org/x/sync/errgroup"
)

type MovieUseCase struct {
	workersCount int
	provider     MovieProvider
	repository   Repository
	searcher     MovieSearcher
}

func (useCase *MovieUseCase) LoadMovies(ctx context.Context, from int, to int) ([]domain.Movie, error) {
	err := requestValidate(from, to)
	if err != nil {
		return nil, err
	}

	idsRange := make([]int, to+1-from)
	for i := 0; i < to+1-from; i++ {
		idsRange[i] = i + from
	}

	return useCase.repository.GetMoviesByIDs(ctx, idsRange)
}

func (useCase *MovieUseCase) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	return useCase.repository.LastMovie(ctx)
}

func (useCase *MovieUseCase) SearchMovie(ctx context.Context, query string) ([]domain.Movie, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return nil, domain.ErrInvalidSearchQuery
	}

	if movies, ok, err := useCase.searcher.GetCached(ctx, query); err != nil {
		return nil, fmt.Errorf("ceshe get: %w", err)
	} else if ok {
		return movies, nil
	}

	words := strings.Fields(query)
	ids, err := useCase.searcher.Search(ctx, words)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	if len(ids) == 0 {
		err := useCase.searcher.SetCached(ctx, query, []domain.Movie{})
		if err != nil {
			return nil, fmt.Errorf("cache set: %w", err)
		}
		return nil, nil
	}

	movies, err := useCase.repository.GetMoviesByIDs(ctx, ids)

	err = useCase.searcher.SetCached(ctx, query, movies)
	if err != nil {
		return nil, fmt.Errorf("cache set: %w", err)
	}
	return movies, nil

}

func (useCase *MovieUseCase) SyncDatabase(ctx context.Context) error {
	lastMovie, err := findLastMovie(ctx, useCase.provider)
	if err != nil {
		return err
	}

	existIDs, err := useCase.repository.ExistIDs(ctx)
	if err != nil {
		return err
	}

	movieIDs := make(chan int)

	group, groupCTX := errgroup.WithContext(ctx)

	group.Go(func() error {
		defer close(movieIDs)
		for movieID := 1; movieID <= lastMovie.ID; movieID++ {
			if _, ok := existIDs[movieID]; ok {
				continue
			}
			select {
			case movieIDs <- movieID:
			case <-ctx.Done():
				return groupCTX.Err()
			}
		}
		return nil
	})

	for i := 0; i < useCase.workersCount; i++ {
		group.Go(func() error { return syncMovieWorker(groupCTX, movieIDs, useCase.provider, useCase.repository) })
	}

	return group.Wait()
}

func NewFilmService(config UseCaseConfig, provider MovieProvider, database Repository, movieSercher MovieSearcher) *MovieUseCase {
	return &MovieUseCase{
		workersCount: config.WorkersCount,
		provider:     provider,
		repository:   database,
		searcher:     movieSercher,
	}
}
