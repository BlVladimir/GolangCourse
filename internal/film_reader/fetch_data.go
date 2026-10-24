package film_reader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func fetch(ctx context.Context, filmID int, client *http.Client, timeout time.Duration) (film_models.Film, error) {

	requestCtx, cancle := context.WithTimeout(ctx, timeout)
	defer cancle()

	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", filmID)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return film_models.Film{}, err
	}

	response, err := client.Do(request)
	if err != nil {
		return film_models.Film{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return film_models.Film{}, fmt.Errorf("status %d", response.StatusCode)
	}

	var filmData film_models.Film
	err = json.NewDecoder(response.Body).Decode(&filmData)
	if err != nil {
		return film_models.Film{}, err
	}

	return filmData, nil
}
