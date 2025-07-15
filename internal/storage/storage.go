package storage

import "errors"

var (
	ErrInvalidWallet     = errors.New("wrong wallet adress")
	ErrInsufficientFunds = errors.New("low balance")
	ErrWalletNotFound    = errors.New("wallet not found")
)
