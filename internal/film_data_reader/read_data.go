package filmdatareader

import (
	"sync"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func Run(from int, to int, workersCount int) <-chan film_models.LoadInfo {
	filmIDs := make(chan int)
	loadedFilms := make(chan film_models.LoadInfo)

	var wg sync.WaitGroup

	go func() {
		defer close(filmIDs)
		for filmID := from; filmID <= to; filmID++ {
			filmIDs <- filmID
		}
	}()

	for workerNumber := 0; workerNumber < workersCount; workerNumber++ {
		wg.Add(1)
		go filmWorker(filmIDs, loadedFilms, &wg)
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	return loadedFilms
}
