package filmdatareader

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func worker(ctx context.Context, filmIDs <-chan int, loadedFilms chan<- film_models.LoadInfo, client *http.Client, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range filmIDs {
		filmData, err := fetch(ctx, id, client, timeout)
		film := film_models.LoadInfo{ID: id, FilmInfo: filmData, Err: err}
		loadedFilms <- film
	}
}
