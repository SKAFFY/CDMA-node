package walsh

import (
	some_stuff "CDMA-telecom-lab1/internal/some-stuff"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextPowerOf2(t *testing.T) {

	tests := [...]struct {
		name     string
		argument int64
		expected int64
	}{
		{
			name:     "got 0 expect 1",
			argument: 0,
			expected: 1,
		},
		{
			name:     "got 1 expect 1",
			argument: 1,
			expected: 1,
		},
		{
			name:     "got 2 expect 2",
			argument: 2,
			expected: 2,
		},
		{
			name:     "got 3 expect 4",
			argument: 3,
			expected: 4,
		},
		{
			name:     "got 4 expect 4",
			argument: 4,
			expected: 4,
		},
		{
			name:     "got 5 expect 8",
			argument: 5,
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := some_stuff.NextPowOf2(tt.argument)

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestNewWalshTable(t *testing.T) {

	tests := [...]struct {
		name     string
		argument int64
		expected [][]int64
	}{
		{
			name:     "got 1 expect 1",
			argument: 1,
			expected: [][]int64{{1}},
		},
		{
			name:     "got 2 expect 2",
			argument: 2,
			expected: H2,
		},
		{
			name:     "got 3 expect 3",
			argument: 3,
			expected: H4[0:3],
		},
		{
			name:     "got 4 expect 4",
			argument: 4,
			expected: H4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWalshTable(tt.argument)

			assert.Equal(t, tt.expected, got)
		})
	}
}
