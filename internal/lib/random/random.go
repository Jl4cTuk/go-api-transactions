package random

import (
	"database/sql"

	"math/rand"
)

type Storage struct {
	db *sql.DB
}

var letterRunes = []rune("1234567890abcdef")

// GenAddress generates a random wallet address.
func GenAddress(length int, r *rand.Rand) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
