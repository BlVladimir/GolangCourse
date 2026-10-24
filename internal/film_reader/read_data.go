package film_reader

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func Run(ctx context.Context, from int, to int, workersCount int, client *http.Client, timeout time.Duration) film_models.LoadReport {
	filmIDs := make(chan int)
	loadedFilms := make(chan film_models.FilmResult)

	var wg sync.WaitGroup

	go func() {
		defer close(filmIDs)
		for filmID := from; filmID <= to; filmID++ {
			filmIDs <- filmID
		}
	}()

	for workerNumber := 0; workerNumber < workersCount; workerNumber++ {
		wg.Add(1)
		go worker(ctx, filmIDs, loadedFilms, client, timeout, &wg)
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	var filmResults []film_models.FilmResult
	for filmResult := range loadedFilms {
		filmResults = append(filmResults, filmResult)
	}

	return film_models.LoadReport{FilmResults: filmResults}
}
