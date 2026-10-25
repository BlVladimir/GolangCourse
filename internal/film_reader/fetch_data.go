package film_reader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

var errNotFound = errors.New("film not found")

func fetch(ctx context.Context, baseURL string, filmID int, client *http.Client, timeout time.Duration) (film_models.Film, error) {

	requestCtx, cancle := context.WithTimeout(ctx, timeout)
	defer cancle()

	url := fmt.Sprintf("%s/%d/info.0.json", baseURL, filmID)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return film_models.Film{}, err
	}

	response, err := client.Do(request)
	if err != nil {
		return film_models.Film{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return film_models.Film{}, errNotFound
	} else if response.StatusCode != http.StatusOK {
		return film_models.Film{}, fmt.Errorf("status %d", response.StatusCode)
	}

	var filmData film_models.Film
	err = json.NewDecoder(response.Body).Decode(&filmData)
	if err != nil {
		return film_models.Film{}, err
	}

	return filmData, nil
}
