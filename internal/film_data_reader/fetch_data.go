package filmdatareader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func fetch(ctx context.Context, filmID int, client *http.Client, timeout time.Duration) (film_models.FilmInfo, error) {

	requestCtx, cancle := context.WithTimeout(ctx, timeout)
	defer cancle()

	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", filmID)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return film_models.FilmInfo{}, err
	}

	response, err := client.Do(request)
	if err != nil {
		return film_models.FilmInfo{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return film_models.FilmInfo{}, fmt.Errorf("status %d", response.StatusCode)
	}

	var filmData film_models.FilmInfo
	err = json.NewDecoder(response.Body).Decode(&filmData)
	if err != nil {
		return film_models.FilmInfo{}, err
	}

	return filmData, nil
}
