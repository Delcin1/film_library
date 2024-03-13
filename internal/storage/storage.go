package storage

import "errors"

var (
	ErrMovieNotFound = errors.New("movie not found")
	ErrURLExists     = errors.New("url exists")
)
