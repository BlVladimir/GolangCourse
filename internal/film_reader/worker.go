package film_reader

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func ReadFilm(ctx context.Context, baseURL string, filmIDs <-chan int, loadedFilms chan<- film_models.FilmResult, client *http.Client, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range filmIDs {
		filmData, err := fetch(ctx, baseURL, id, client, timeout)

		if errors.Is(err, context.Canceled) {
			return
		}
		film := film_models.FilmResult{ID: id, FilmInfo: filmData, Err: err}
		loadedFilms <- film
	}
}
