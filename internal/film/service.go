package film

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/model"
)

type FilmService struct {
	client       *http.Client
	baseURL      string
	workersCount int
	timeout      time.Duration
}

func (service *FilmService) LoadMovies(ctx context.Context, from int, to int) (model.LoadReport, error) {
	requestValidationErr := requestValidate(from, to)
	if requestValidationErr != nil {
		return model.LoadReport{}, requestValidationErr
	}

	filmIDs := make(chan int)
	loadedFilms := make(chan model.FilmResult)

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

	for workerNumber := 0; workerNumber < service.workersCount; workerNumber++ {
		wg.Add(1)
		go readFilmWorker(ctx, service.baseURL, filmIDs, loadedFilms, service.client, service.timeout, &wg)
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	var filmResults []model.FilmResult
	for {
		select {
		case <-ctx.Done():
			return model.LoadReport{}, ctx.Err() // клиент ушёл — не дочитываем
		case filmResult, ok := <-loadedFilms:
			if !ok {
				return model.LoadReport{FilmResults: filmResults}, nil
			}
			filmResults = append(filmResults, filmResult)
		}
	}
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
