package film_models

import (
	"fmt"
)

// All information about single film.
type Film struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
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
