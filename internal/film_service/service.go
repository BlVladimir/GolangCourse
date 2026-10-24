package film_service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/film_models"
	"github.com/BlVladimir/GolangCourse/internal/film_reader"
)

type FilmService struct {
	client       *http.Client
	baseURL      string
	workersCount int
	timeout      time.Duration
}

func (service *FilmService) LoadMovies(ctx context.Context, from int, to int) film_models.LoadReport {
	requestValidationErr := requestValidate(from, to)
	if requestValidationErr != nil {
		fmt.Println(requestValidationErr)
		os.Exit(1)
	}
	return film_reader.Run(ctx, from, to, service.workersCount, service.client, service.timeout)
}

func NewFilmService(baseURL string, workersCount int, timeout time.Duration) *FilmService {
	configValidationErr := configValidate(workersCount, timeout)
	if configValidationErr != nil {
		fmt.Println(configValidationErr)
		os.Exit(1)
	}

	client := &http.Client{}

	return &FilmService{
		client:       client,
		baseURL:      baseURL,
		workersCount: workersCount,
		timeout:      timeout,
	}
}
