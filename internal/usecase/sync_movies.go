package usecase

import (
	"context"
	"fmt"
)

func syncMovieWorker(
	ctx context.Context,
	movieIDs <-chan int,
	provider MovieProvider,
	repository Repository,
) error {
	for {
		select {
		case id, ok := <-movieIDs:
			if !ok {
				return nil
			}

			movie, err := provider.GetMovie(ctx, id)
			if err != nil {
				return fmt.Errorf("sync: get id=%d: %w", id, err)
			}

			if err := repository.UpsertMovie(ctx, movie); err != nil {
				return fmt.Errorf("sync: upsert id=%d: %w", id, err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
