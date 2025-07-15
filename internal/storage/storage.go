package storage

import "errors"

var (
	ErrWalletNotFound = errors.New("no rows in result set")
	ErrURLExists      = errors.New("alias exists")
)
