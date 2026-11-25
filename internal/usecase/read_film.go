package usecase

import (
	"context"
	"sync"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

func readFilmWorker(
	ctx context.Context,
	filmIDs <-chan int,
	loadedFilms chan<- domain.Film,
	provider FilmProvider,
	wg *sync.WaitGroup,
) error {
	defer wg.Done()

	for {
		select {
		case id, ok := <-filmIDs:
			if !ok {
				return nil
			}

			film, err := provider.GetMovie(ctx, id)
			if err != nil {
				return err
			}

			select {
			case loadedFilms <- film:
			case <-ctx.Done():
				return nil
			}

		case <-ctx.Done():
			return nil
		}
	}
}
