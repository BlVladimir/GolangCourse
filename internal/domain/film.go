package domain

import (
	"fmt"
)

// All information about single film.
type Film struct {
	ID       int
	Title    string
	Year     int
	Director string
}

func (film Film) String() string {
	return fmt.Sprintf(
		"%d - %s - %d - %s\n",
		film.ID,
		film.Title,
		film.Year,
		film.Director,
	)
}
