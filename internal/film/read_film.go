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

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return model.Film{}, fmt.Errorf("film %d: %w", filmID, ErrNotFound)
	default:
		return model.Film{}, fmt.Errorf("%w: film %d: status %d", ErrUpstream, filmID, response.StatusCode)
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
	loadedFilms chan<- model.Film,
	client *http.Client,
	timeout time.Duration,
	wg *sync.WaitGroup,
) error {
	defer wg.Done()

	for {
		select {
		case id, ok := <-filmIDs:
			if !ok {
				return nil
			}

			film, err := readFilm(ctx, baseURL, id, client, timeout)
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
