package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	const from int = 1
	const to int = 10
	const workers int = 10
	const timeout int = 5

	responce, error := http.Get("https://homeworksite.site/1/info.0.json")

	if error != nil {
		fmt.Println("Error response:", error)
	}

	defer responce.Body.Close()

	body, error := io.ReadAll(responce.Body)

	if error != nil {
		fmt.Println("Error reading:", error)
	}

	fmt.Println("Status:", responce.StatusCode)
	printFilmInfo(body)
}

type MovieData struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func printFilmInfo(infoBytes []byte) {
	var movieData MovieData

	error := json.Unmarshal(infoBytes, &movieData)
	if error != nil {
		fmt.Println("Error parcing JSON:", error)
		return
	}

	fmt.Printf(
		"%d - %s - %d - %s\n",
		movieData.ID,
		movieData.Title,
		movieData.Year,
		movieData.Director,
	)
}
