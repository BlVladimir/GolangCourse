package redis_searcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/domain"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
	"github.com/redis/go-redis/v9"
)

const (
	indexKeyPrefix = "index:"
	cacheKeyPrefix = "cache:q="
	cacheTimeLive  = time.Hour
)

type MoviesSearcher struct {
	client   *redis.Client
	timeLive time.Duration
}

func NewMoviesSearcher(client *redis.Client) *MoviesSearcher {
	return &MoviesSearcher{
		client:   client,
		timeLive: cacheTimeLive,
	}
}

func (searcher MoviesSearcher) IndexMovie(ctx context.Context, movie domain.Movie) error {
	text := movie.Title + " " + movie.Director
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return nil
	}

	pipe := searcher.client.Pipeline()
	for _, word := range words {
		pipe.SAdd(ctx, indexKeyPrefix+word, movie.ID)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis index movie %d: %w", movie.ID, err)
	}

	return nil
}

func (searcher MoviesSearcher) Search(ctx context.Context, words []string) ([]int, error) {
	if len(words) == 0 {
		return nil, nil
	}

	keys := make([]string, len(words))
	for i, w := range words {
		keys[i] = indexKeyPrefix + strings.ToLower(w)
	}

	result, err := searcher.client.SInter(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis sinter: %w", err)
	}

	ids := make([]int, 0, len(result))
	for _, v := range result {
		var id int
		if _, err := fmt.Sscanf(v, "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (searcher MoviesSearcher) Clear(ctx context.Context) error {
	patterns := []string{indexKeyPrefix + "*", cacheKeyPrefix + "*"}
	for _, pattern := range patterns {
		iter := searcher.client.Scan(ctx, 0, pattern, 1000).Iterator()
		for iter.Next(ctx) {
			if err := searcher.client.Del(ctx, iter.Val()).Err(); err != nil {
				return fmt.Errorf("redis clear key %s: %w", iter.Val(), err)
			}
		}
		if err := iter.Err(); err != nil {
			return fmt.Errorf("redis scan %s: %w", pattern, err)
		}
	}
	return nil
}
func (searcher MoviesSearcher) GetCached(ctx context.Context, query string) ([]domain.Movie, bool, error) {
	key := cacheKeyPrefix + query
	data, err := searcher.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("redis get cache: %w", err)
	}

	var movies []domain.Movie
	if err := json.Unmarshal(data, &movies); err != nil {
		return nil, false, fmt.Errorf("redis unmarshal cache: %w", err)
	}
	return movies, true, nil
}

func (searcher MoviesSearcher) SetCached(ctx context.Context, query string, movies []domain.Movie) error {
	data, err := json.Marshal(movies)
	if err != nil {
		return fmt.Errorf("redis marshal cache: %w", err)
	}

	key := cacheKeyPrefix + query
	if err := searcher.client.Set(ctx, key, data, searcher.timeLive).Err(); err != nil {
		return fmt.Errorf("redis set cache: %w", err)
	}
	return nil
}

var _ usecase.MovieSearcher = MoviesSearcher{}
