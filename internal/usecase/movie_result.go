package usecase

import (
	"fmt"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

// Result of loading single movie.
type MovieResult struct {
	ID        int
	MovieInfo domain.Movie
	Err       error
}

func (movieResult MovieResult) String() string {
	if movieResult.Err == nil {
		return movieResult.MovieInfo.String()
	} else {
		return fmt.Sprintf("%d - error: %s", movieResult.ID, movieResult.Err)
	}
}
