package random

import (
	"database/sql"
	"time"

	"math/rand"
)

type Storage struct {
	db *sql.DB
}

var letterRunes = []rune("1234567890abcdef")
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenAddress generates a random wallet address.
func GenAddress(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
