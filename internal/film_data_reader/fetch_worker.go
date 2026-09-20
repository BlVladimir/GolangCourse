package filmdatareader

import (
	"sync"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func filmWorker(filmIDs <-chan int, loadedFilms chan<- film_models.LoadInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range filmIDs {
		filmData, err := fetchFilmData(id)
		film := film_models.LoadInfo{ID: id, FilmInfo: filmData, Err: err}
		loadedFilms <- film
	}
}
