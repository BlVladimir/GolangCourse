package filmstand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
)

type Client struct {
	baseURL string
	client  *http.Client
	timeout time.Duration
}

func (client *Client) GetMovie(ctx context.Context, id int) (domain.Film, error) {
	requestCtx, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()

	url := fmt.Sprintf("%s/%d/info.0.json", client.baseURL, id)
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Film{}, err
	}

	response, err := client.client.Do(request)
	if err != nil {
		return domain.Film{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return domain.Film{}, domain.ErrNotFound
	} else if response.StatusCode != http.StatusOK {
		return domain.Film{}, fmt.Errorf("status %d", response.StatusCode)
	}

	var filmData standMovie
	err = json.NewDecoder(response.Body).Decode(&filmData)
	if err != nil {
		return domain.Film{}, err
	}

	return filmData.toDomain(), nil
}

func NewClient(config Config) *Client {
	client := &http.Client{}

	return &Client{
		client:  client,
		baseURL: config.MovieBaseURL,
		timeout: config.Timeout,
	}
}

var _ usecase.FilmProvider = (*Client)(nil)
