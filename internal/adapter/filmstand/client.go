package filmstand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/config"
	"github.com/BlVladimir/GolangCourse/internal/domain"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
)

type Client struct {
	baseURL string
	client  *http.Client
	timeout time.Duration
}

func (client *Client) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	requestCtx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	url := fmt.Sprintf("%s/%d/info.0.json", client.baseURL, id)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Movie{}, err
	}

	response, err := client.client.Do(request)
	if err != nil {
		return domain.Movie{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return domain.Movie{}, domain.ErrNotFound
	} else if response.StatusCode != http.StatusOK {
		return domain.Movie{}, fmt.Errorf("status %d", response.StatusCode)
	}

	var filmData standMovie
	err = json.NewDecoder(response.Body).Decode(&filmData)
	if err != nil {
		return domain.Movie{}, err
	}

	return filmData.toDomain(), nil
}

func NewClient(cnf config.Config) *Client {
	client := &http.Client{}

	return &Client{
		client:  client,
		baseURL: cnf.MovieBaseURL,
		timeout: cnf.Timeout,
	}
}

var _ usecase.MovieProvider = (*Client)(nil)
