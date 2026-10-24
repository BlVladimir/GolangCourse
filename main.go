package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	filmdatareader "github.com/BlVladimir/GolangCourse/internal/film_data_reader"
)

func main() {
	from := flag.Int("from", 0, "id first film (required arg)")
	to := flag.Int("to", 0, "id last film (required arg)")
	workersCount := flag.Int("workers", 10, "count of workers")
	timeout := flag.Duration("timeout", 5*time.Second, "timeout of one request")

	flag.Parse()

	set := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { set[f.Name] = true })

	if !set["from"] {
		fmt.Fprintln(os.Stderr, "error: --from is required")
		flag.Usage()
		os.Exit(1)
	}
	if !set["to"] {
		fmt.Fprintln(os.Stderr, "error: --to is required")
		flag.Usage()
		os.Exit(1)
	}

	validationErr := validate(*from, *to, *workersCount, *timeout)
	if validationErr != nil {
		fmt.Println(validationErr)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	client := &http.Client{}

	defer stop()
	for _, loadedFilm := range filmdatareader.Run(ctx, *from, *to, *workersCount, client, *timeout).FilmResults {
		fmt.Println(loadedFilm.String())
	}
}

func validate(from int, to int, workers int, timeout time.Duration) error {

	if from <= 0 {
		return fmt.Errorf("--from must be > 0, got %d", from)
	}
	if to <= 0 {
		return fmt.Errorf("--to must be > 0, got %d", to)
	}
	if from > to {
		return fmt.Errorf("--from (%d) must not be greater than --to (%d)", from, to)
	}

	if workers <= 0 {
		return fmt.Errorf("--workers must be > 0, got %d", workers)
	}
	if timeout <= 0 {
		return fmt.Errorf("--timeout must be > 0, got %s", timeout)
	}

	return nil
}
