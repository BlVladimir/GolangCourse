package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	const from, to int = 1, 10
	const workersCount int = 10
	const timeout int = 5

	filmIDs := make(chan int)
	loadedFilms := make(chan LoadedFilm)

	var wg sync.WaitGroup

	go func() {
		defer close(filmIDs)
		for filmID := from; filmID <= to; filmID++ {
			filmIDs <- filmID
		}
	}()

	for workerNumber := 0; workerNumber < workersCount; workerNumber++ {
		wg.Add(1)
		go filmWorker(filmIDs, loadedFilms, &wg)
	}

	go func() {
		wg.Wait()
		close(loadedFilms)
	}()

	for loadedFilm := range loadedFilms {
		if loadedFilm.Err == nil {
			printFilmInfo(loadedFilm.FilmInfo)
		} else {
			fmt.Printf("%d - error: %s", loadedFilm.ID, loadedFilm.Err)
		}
	}
}

func getFilmData(filmID int) (FilmInfo, error) {
	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", filmID)

	responce, err := http.Get(url)
	if err != nil {
		return FilmInfo{}, err
	}
	defer responce.Body.Close()

	body, err := io.ReadAll(responce.Body)
	if err != nil {
		return FilmInfo{}, err
	}

	var filmData FilmInfo
	err = json.Unmarshal(body, &filmData)
	if err != nil {
		return FilmInfo{}, err
	}

	return filmData, nil
}

func filmWorker(filmIDs <-chan int, loadedFilms chan<- LoadedFilm, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range filmIDs {
		filmData, err := getFilmData(id)
		film := LoadedFilm{ID: id, FilmInfo: filmData, Err: err}
		loadedFilms <- film
	}
}

type LoadedFilm struct {
	ID       int
	FilmInfo FilmInfo
	Err      error
}

type FilmInfo struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func printFilmInfo(filmInfo FilmInfo) {
	fmt.Printf(
		"%d - %s - %d - %s\n",
		filmInfo.ID,
		filmInfo.Title,
		filmInfo.Year,
		filmInfo.Director,
	)
}
