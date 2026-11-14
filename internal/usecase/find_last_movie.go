package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

func findLastMovie(ctx context.Context, provider MovieProvider) (domain.Movie, error) {
	var err error = nil
	var lowID int = 1
	var highID int = 100000
	for lowID <= highID {
		var midID int = (lowID + highID) / 2
		_, err = provider.GetMovie(ctx, midID)
		if errors.Is(err, domain.ErrNotFound) {
			highID = midID - 1
		} else if err != nil {
			return domain.Movie{}, err
		} else {
			lowID = midID + 1
		}
	}
	if lowID == 1 {
		return domain.Movie{}, fmt.Errorf("stand has no films: %w", domain.ErrUpstream)
	}
	return provider.GetMovie(ctx, lowID-1)
}
