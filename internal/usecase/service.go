package usecase

import (
	"context"
	"sync"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type FilmUseCase struct {
	workersCount int
	provider     FilmProvider
}

func (useCase *FilmUseCase) LoadMovies(ctx context.Context, from int, to int) (LoadReport, error) {
	requestValidationErr := requestValidate(from, to)
	if requestValidationErr != nil {
		return LoadReport{}, requestValidationErr
	}

	filmIDs := make(chan int)
	loadedFilms := make(chan FilmResult)

	var wg sync.WaitGroup

	go func() {
		defer close(filmIDs)
		for filmID := from; filmID <= to; filmID++ {
			select {
			case filmIDs <- filmID:
			case <-ctx.Done():
				return
			}
		}
	}()

	for workerNumber := 0; workerNumber < useCase.workersCount; workerNumber++ {
		wg.Add(1)
		go readFilmWorker(ctx, filmIDs, loadedFilms, useCase.provider, &wg)
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	var filmResults []FilmResult
	for {
		select {
		case <-ctx.Done():
			return LoadReport{}, ctx.Err() // клиент ушёл — не дочитываем
		case filmResult, ok := <-loadedFilms:
			if !ok {
				return LoadReport{FilmResults: filmResults}, nil
			}
			filmResults = append(filmResults, filmResult)
		}
	}
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
