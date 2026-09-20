package film_models

import (
	"fmt"
)

type FilmInfo struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func (filmInfo FilmInfo) String() string {
	return fmt.Sprintf(
		"%d - %s - %d - %s\n",
		filmInfo.ID,
		filmInfo.Title,
		filmInfo.Year,
		filmInfo.Director,
	)
}
