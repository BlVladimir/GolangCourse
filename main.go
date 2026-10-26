package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film"
	"github.com/BlVladimir/GolangCourse/internal/http_api"
)

func main() {
	cfg, err := film.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	svc := film.NewFilmService(cfg)

	router := http_api.NewRouter(http_api.New(svc))

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second, // GetLastMovie: ~2·log(N) запросов
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("HTTP listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
