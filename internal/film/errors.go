package film

import "errors"

var (
	ErrInvalidRange = errors.New("invalid range")

	ErrNotFound = errors.New("film not found")

	ErrUpstream = errors.New("upstream unavailable")
)
