package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	filmdatareader "github.com/BlVladimir/GolangCourse/internal/film_data_reader"
)

func main() {
	const from, to int = 1, 10
	const workersCount int = 10
	const timeout time.Duration = 5 * time.Second

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	client := &http.Client{}

	defer stop()
	for loadedFilm := range filmdatareader.Run(ctx, from, to, workersCount, client, timeout) {
		fmt.Println(loadedFilm.String())
	}
}
