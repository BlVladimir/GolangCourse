package film_reader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func GetLastFilm(ctx context.Context, baseURL string, client *http.Client, timeout time.Duration) (film_models.Film, error) {
	var lastFilmID int = 1
	var err error = nil
	for err == nil {
		lastFilmID *= 2
		_, err = fetch(ctx, baseURL, lastFilmID, client, timeout)
		if !errors.Is(err, errNotFound) && err != nil {
			return film_models.Film{}, err
		}
	}

	fmt.Println("Last film ID:", lastFilmID-1)

	var lowID int = 1
	var highID int = lastFilmID - 1
	for lowID <= highID {
		var midID int = (lowID + highID) / 2
		_, err = fetch(ctx, baseURL, midID, client, timeout)
		if errors.Is(err, errNotFound) {
			highID = midID - 1
		} else if err != nil {
			return film_models.Film{}, err
		} else {
			lowID = midID + 1
		}
	}
	fmt.Println("Low ID:", lowID-1)
	return fetch(ctx, baseURL, lowID-1, client, timeout)
}
