package usecase

import (
	"context"
	"sync"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"golang.org/x/sync/errgroup"
)

type FilmUseCase struct {
	workersCount int
	provider     FilmProvider
}

func (useCase *FilmUseCase) LoadMovies(ctx context.Context, from int, to int) ([]domain.Film, error) {
	requestValidationErr := requestValidate(from, to)
	if requestValidationErr != nil {
		return nil, requestValidationErr
	}

	filmIDs := make(chan int)
	loadedFilms := make(chan domain.Film)

	group, groupCTX := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer close(filmIDs)
		for filmID := from; filmID <= to; filmID++ {
			select {
			case filmIDs <- filmID:
			case <-groupCTX.Done():
				return groupCTX.Err()
			}
		}
		return nil
	})

	var wg sync.WaitGroup
	for workerNumber := 0; workerNumber < useCase.workersCount; workerNumber++ {
		wg.Add(1)
		group.Go(func() error {
			return readFilmWorker(groupCTX, filmIDs, loadedFilms, useCase.provider, &wg)
		})
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	var films []domain.Film
	for film := range loadedFilms {
		films = append(films, film)
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}
	return films, nil
}

func (useCase *FilmUseCase) GetLastMovie(ctx context.Context) (domain.Film, error) {
	return findLastFilm(ctx, useCase.provider)
}

func NewFilmService(config UseCaseConfig, provider FilmProvider) *FilmUseCase {
	return &FilmUseCase{
		workersCount: config.WorkersCount,
		provider:     provider,
	}
}
