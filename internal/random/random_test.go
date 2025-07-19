package random

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenAddress(t *testing.T) {
	tests := []struct {
		name string
		size int
		r    *rand.Rand
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 100",
			size: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := GenAddress(tt.size)
			str2 := GenAddress(tt.size)

			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			// Check that two generated strings are different
			// This is not an absolute guarantee that the function works correctly,
			// but this is a good heuristic for a simple random generator.
			assert.NotEqual(t, str1, str2)
		})
	}
}
