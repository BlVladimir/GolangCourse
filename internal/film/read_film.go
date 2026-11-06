package film

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/model"
)

func readFilm(ctx context.Context, baseURL string, filmID int, client *http.Client, timeout time.Duration) (model.Film, error) {

	requestCtx, cancle := context.WithTimeout(ctx, timeout)
	defer cancle()

	url := fmt.Sprintf("%s/%d/info.0.json", baseURL, filmID)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return model.Film{}, err
	}

	response, err := client.Do(request)
	if err != nil {
		return model.Film{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return model.Film{}, ErrNotFound
	} else if response.StatusCode != http.StatusOK {
		return model.Film{}, fmt.Errorf("status %d", response.StatusCode)
	}

	var filmData model.Film
	err = json.NewDecoder(response.Body).Decode(&filmData)
	if err != nil {
		return model.Film{}, err
	}

	return filmData, nil
}

func readFilmWorker(
	ctx context.Context,
	baseURL string,
	filmIDs <-chan int,
	loadedFilms chan<- model.FilmResult,
	client *http.Client,
	timeout time.Duration,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case id, ok := <-filmIDs:
			if !ok {
				return
			}

			filmData, err := readFilm(ctx, baseURL, id, client, timeout)
			if ctx.Err() != nil {
				return
			}

			film := model.FilmResult{ID: id, FilmInfo: filmData, Err: err}
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
