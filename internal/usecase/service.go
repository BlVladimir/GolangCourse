package usecase

import (
	"context"

	// "sync"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"golang.org/x/sync/errgroup"
)

type MovieUseCase struct {
	workersCount int
	provider     MovieProvider
	repository   Repository
}

func (useCase *MovieUseCase) LoadMovies(ctx context.Context, from int, to int) ([]domain.Movie, error) {
	err := requestValidate(from, to)
	if err != nil {
		return nil, err
	}
	return useCase.repository.GetMoviesRange(ctx, from, to)
}

func (useCase *MovieUseCase) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	return useCase.repository.LastMovie(ctx)
	// return findLastMovie(ctx, useCase.provider)
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

func NewFilmService(config UseCaseConfig, provider MovieProvider, database Repository) *MovieUseCase {
	return &MovieUseCase{
		workersCount: config.WorkersCount,
		provider:     provider,
		repository:   database,
	}
}
