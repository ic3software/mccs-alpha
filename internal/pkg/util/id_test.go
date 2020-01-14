package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToIDStrings(t *testing.T) {
	idString1 := "5d5516a0f613a4f874b1bf1d"
	idString2 := "5d5516a0f613a4f874b1bf1e"
	idString3 := "5d5516a0f613a4f874b1bf1f"
	idString4 := "5d5516a0f613a4f874b1bf20"
	id1, _ := primitive.ObjectIDFromHex(idString1)
	id2, _ := primitive.ObjectIDFromHex(idString2)
	id3, _ := primitive.ObjectIDFromHex(idString3)
	id4, _ := primitive.ObjectIDFromHex(idString4)

	tests := []struct {
		input    []primitive.ObjectID
		expected []string
	}{
		{
			[]primitive.ObjectID{id1, id2, id3, id4},
			[]string{idString1, idString2, idString3, idString4},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := ToIDStrings(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
