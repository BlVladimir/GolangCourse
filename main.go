package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_service"
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

	service := film_service.NewFilmService("https://homeworksite.site", *workersCount, *timeout)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	filmReport, err := service.LoadMovies(ctx, *from, *to)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	for _, loadedFilm := range filmReport.FilmResults {
		fmt.Println(loadedFilm.String())
	}

	fmt.Println(service.GetLastMovie(ctx))
}
