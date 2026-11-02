package usecase

import (
	"context"
	"errors"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

func findLastFilm(ctx context.Context, provider FilmProvider) (domain.Film, error) {
	var err error = nil
	var lowID int = 1
	var highID int = 100000
	for lowID <= highID {
		var midID int = (lowID + highID) / 2
		_, err = provider.GetMovie(ctx, midID)
		if errors.Is(err, domain.ErrNotFound) {
			highID = midID - 1
		} else if err != nil {
			return domain.Film{}, err
		} else {
			lowID = midID + 1
		}
	}
	return provider.GetMovie(ctx, lowID-1)
}
