package usecase

import "context"

func syncOneMovie(ctx context.Context, provider MovieProvider, repository Repository, movieID int, existIDs map[int]struct{}) error {
	if _, ok := existIDs[movieID]; ok {
		return nil
	}

	movie, err := provider.GetMovie(ctx, movieID)
	if err != nil {
		return err
	}

	err = repository.UpsertMovie(ctx, movie)
	if err != nil {
		return err
	}

	return nil
}
