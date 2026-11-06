package film

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/model"
)

func findLastFilm(ctx context.Context, baseURL string, client *http.Client, timeout time.Duration) (model.Film, error) {
	var err error = nil
	var lowID int = 1
	var highID int = 100000
	for lowID <= highID {
		var midID int = (lowID + highID) / 2
		_, err = readFilm(ctx, baseURL, midID, client, timeout)
		if errors.Is(err, ErrNotFound) {
			highID = midID - 1
		} else if err != nil {
			return model.Film{}, err
		} else {
			lowID = midID + 1
		}
	}
	return readFilm(ctx, baseURL, lowID-1, client, timeout)
}
