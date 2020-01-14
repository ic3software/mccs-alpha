package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDecimalValid(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"10",
			true,
		},
		{
			"10.1",
			true,
		},
		{
			"10.12",
			true,
		},
		{
			"10.123",
			false,
		},
		{
			"10.123.13",
			false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := IsDecimalValid(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
