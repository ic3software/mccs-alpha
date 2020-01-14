package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToIDStrings converts Object IDs into strings.
func ToIDStrings(objIDs []primitive.ObjectID) []string {
	ids := make([]string, 0, len(objIDs))
	for _, objID := range objIDs {
		ids = append(ids, objID.Hex())
	}
	return ids
}
