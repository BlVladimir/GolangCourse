package filmstand

import "github.com/BlVladimir/GolangCourse/internal/domain"

type standMovie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func (stand standMovie) toDomain() domain.Film {
	return domain.Film{ID: stand.ID, Title: stand.Title, Year: stand.Year, Director: stand.Director}
}
