package usecase

import (
	"context"
	// "sync"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type MovieUseCase struct {
	workersCount int
	provider     MovieProvider
	repository   Repository
}

func (useCase *MovieUseCase) LoadMovies(ctx context.Context, from int, to int) ([]domain.Movie, error) {
	return useCase.repository.GetMoviesRange(ctx, from, to)
	// requestValidationErr := requestValidate(from, to)
	// if requestValidationErr != nil {
	// 	return LoadReport{}, requestValidationErr
	// }

	// movieIDs := make(chan int)
	// loadedMovies := make(chan MovieResult)

	// var wg sync.WaitGroup

	// go func() {
	// 	defer close(movieIDs)
	// 	for movieID := from; movieID <= to; movieID++ {
	// 		select {
	// 		case movieIDs <- movieID:
	// 		case <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

	// for workerNumber := 0; workerNumber < useCase.workersCount; workerNumber++ {
	// 	wg.Add(1)
	// 	go readMovieWorker(ctx, movieIDs, loadedMovies, useCase.provider, &wg)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(loadedMovies)
	// }()

	// var movieResults []MovieResult
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return LoadReport{}, ctx.Err()
	// 	case movieResult, ok := <-loadedMovies:
	// 		if !ok {
	// 			return LoadReport{MovieResults: movieResults}, nil
	// 		}
	// 		movieResults = append(movieResults, movieResult)
	// 	}
	// }
}

func (useCase *MovieUseCase) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	return useCase.repository.LastMovie(ctx)
	// return findLastMovie(ctx, useCase.provider)
}

func (useCase *MovieUseCase) SyncDatabase(ctx context.Context) error {
	existIDs, err := useCase.repository.ExistIDs(ctx)
	if err != nil {
		return err
	}

	lastMovie, err := findLastMovie(ctx, useCase.provider)
	if err != nil {
		return err
	}

	for id := 1; id <= lastMovie.ID; id++ {
		err := syncOneMovie(ctx, useCase.provider, useCase.repository, id, existIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewFilmService(config UseCaseConfig, provider MovieProvider, database Repository) *MovieUseCase {
	return &MovieUseCase{
		workersCount: config.WorkersCount,
		provider:     provider,
		repository:   database,
	}
}
