package domain

import (
	"fmt"
)

// All information about single movie.
type Movie struct {
	ID       int
	Title    string
	Year     int
	Director string
}

func (movie Movie) String() string {
	return fmt.Sprintf(
		"%d - %s - %d - %s\n",
		movie.ID,
		movie.Title,
		movie.Year,
		movie.Director,
	)
}
