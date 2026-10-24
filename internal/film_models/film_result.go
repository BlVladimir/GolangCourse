package film_models

import (
	"fmt"
)

// Result of loading single film.
type FilmResult struct {
	ID       int
	FilmInfo Film
	Err      error
}

func (filmResult FilmResult) String() string {
	if filmResult.Err == nil {
		return filmResult.FilmInfo.String()
	} else {
		return fmt.Sprintf("%d - error: %s", filmResult.ID, filmResult.Err)
	}
}
