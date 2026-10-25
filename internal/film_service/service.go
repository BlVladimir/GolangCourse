package film_service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
	"github.com/BlVladimir/GolangCourse/internal/film_reader"
)

type FilmService struct {
	client       *http.Client
	baseURL      string
	workersCount int
	timeout      time.Duration
}

func (service *FilmService) LoadMovies(ctx context.Context, from int, to int) (film_models.LoadReport, error) {
	requestValidationErr := requestValidate(from, to)
	if requestValidationErr != nil {
		return film_models.LoadReport{}, requestValidationErr
	}

	filmIDs := make(chan int)
	loadedFilms := make(chan film_models.FilmResult)

	var wg sync.WaitGroup

	go func() {
		defer close(filmIDs)
		for filmID := from; filmID <= to; filmID++ {
			filmIDs <- filmID
		}
	}()

	for workerNumber := 0; workerNumber < service.workersCount; workerNumber++ {
		wg.Add(1)
		go film_reader.ReadFilm(ctx, service.baseURL, filmIDs, loadedFilms, service.client, service.timeout, &wg)
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	var filmResults []film_models.FilmResult
	for filmResult := range loadedFilms {
		filmResults = append(filmResults, filmResult)
	}

	return film_models.LoadReport{FilmResults: filmResults}, nil
}

func (service *FilmService) GetLastMovie(ctx context.Context) (film_models.Film, error) {
	return film_reader.GetLastFilm(ctx, service.baseURL, service.client, service.timeout)
}

func NewFilmService(config film_models.Config) *FilmService {
	client := &http.Client{}

	return &FilmService{
		client:       client,
		baseURL:      config.MovieBaseURL,
		workersCount: config.Workers,
		timeout:      config.Timeout,
	}
}
