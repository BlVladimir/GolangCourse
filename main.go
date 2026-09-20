package main

import (
	"fmt"

	filmdatareader "github.com/BlVladimir/GolangCourse/internal/film_data_reader"
)

func main() {
	const from, to int = 1, 10
	const workersCount int = 10
	const timeout int = 5

	for loadedFilm := range filmdatareader.Run(from, to, workersCount) {
		fmt.Println(loadedFilm.String())
	}
}
