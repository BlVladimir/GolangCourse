package usecase

import (
	"context"
	"sync"
)

// func readFilm(ctx context.Context, baseURL string, filmID int, client *http.Client, timeout time.Duration) (domain.Film, error) {

// 	requestCtx, cancle := context.WithTimeout(ctx, timeout)
// 	defer cancle()

// 	url := fmt.Sprintf("%s/%d/info.0.json", baseURL, filmID)
// 	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
// 	if err != nil {
// 		return domain.Film{}, err
// 	}

// 	response, err := client.Do(request)
// 	if err != nil {
// 		return domain.Film{}, err
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode == http.StatusNotFound {
// 		return domain.Film{}, domain.ErrNotFound
// 	} else if response.StatusCode != http.StatusOK {
// 		return domain.Film{}, fmt.Errorf("status %d", response.StatusCode)
// 	}

// 	var filmData domain.Film
// 	err = json.NewDecoder(response.Body).Decode(&filmData)
// 	if err != nil {
// 		return domain.Film{}, err
// 	}

// 	return filmData, nil
// }

func readFilmWorker(
	ctx context.Context,
	filmIDs <-chan int,
	loadedFilms chan<- FilmResult,
	provider FilmProvider,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case id, ok := <-filmIDs:
			if !ok {
				return
			}

			filmData, err := provider.GetMovie(ctx, id)
			if ctx.Err() != nil {
				return
			}

			film := FilmResult{ID: id, FilmInfo: filmData, Err: err}
			select {
			case loadedFilms <- film:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
