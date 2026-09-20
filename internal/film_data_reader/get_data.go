package filmdatareader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
)

func fetchFilmData(filmID int) (film_models.FilmInfo, error) {
	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", filmID)

	responce, err := http.Get(url)
	if err != nil {
		return film_models.FilmInfo{}, err
	}
	defer responce.Body.Close()

	body, err := io.ReadAll(responce.Body)
	if err != nil {
		return film_models.FilmInfo{}, err
	}

	var filmData film_models.FilmInfo
	err = json.Unmarshal(body, &filmData)
	if err != nil {
		return film_models.FilmInfo{}, err
	}

	return filmData, nil
}
