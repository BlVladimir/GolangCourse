package film

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/model"
	"golang.org/x/sync/errgroup"
)

type FilmService struct {
	client       *http.Client
	baseURL      string
	workersCount int
	timeout      time.Duration
}

func (service *FilmService) LoadMovies(ctx context.Context, from int, to int) ([]model.Film, error) {
	requestValidationErr := requestValidate(from, to)
	if requestValidationErr != nil {
		return nil, requestValidationErr
	}

	filmIDs := make(chan int)
	loadedFilms := make(chan model.Film)

	// var wg sync.WaitGroup
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
	for workerNumber := 0; workerNumber < service.workersCount; workerNumber++ {
		wg.Add(1)
		group.Go(func() error {
			return readFilmWorker(groupCTX, service.baseURL, filmIDs, loadedFilms, service.client, service.timeout, &wg)
		})
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	var films []model.Film
	for film := range loadedFilms {
		films = append(films, film)
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}
	return films, nil
}

func (service *FilmService) GetLastMovie(ctx context.Context) (model.Film, error) {
	return findLastFilm(ctx, service.baseURL, service.client, service.timeout)
}

func NewFilmService(config Config) *FilmService {
	client := &http.Client{}

	return &FilmService{
		client:       client,
		baseURL:      config.MovieBaseURL,
		workersCount: config.Workers,
		timeout:      config.Timeout,
	}
}
