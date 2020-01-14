package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{
			"",
			time.Time{},
		},
		{
			"error time",
			time.Time{},
		},
		{
			"18 January 2019",
			time.Date(2019, time.January, 18, 00, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := ParseTime(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{
			time.Date(2019, time.January, 1, 00, 00, 0, 0, time.UTC),
			"2019-01-01 00:00:00 UTC",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := FormatTime(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
