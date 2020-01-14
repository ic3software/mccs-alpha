package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"email@domain.com",
			true,
		},
		{
			"firstname.lastname@domain.com",
			true,
		},
		{
			"email@subdomain.domain.com",
			true,
		},
		{
			"firstname+lastname@domain.com",
			true,
		},
		{
			"plainaddress",
			false,
		},
		{
			"#@%^%#$@#$@#.com",
			false,
		},
		{
			"@domain.com",
			false,
		},
		{
			"あいうえお@domain.com",
			false,
		},
		{
			"jdoe1test@opencredit.network",
			true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := IsValidEmail(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
