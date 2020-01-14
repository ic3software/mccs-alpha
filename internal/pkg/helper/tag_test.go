package helper

import (
	"testing"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"should format user input tags",
			"egg,apple2,apple2,apple3,james@,TRE~E",
			[]string{"egg", "apple", "james", "tree"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := []string{}
			for _, v := range GetTags(tt.input) {
				actual = append(actual, v.Name)
			}
			assert.ElementsMatch(t, actual, tt.expected)
		})
	}
}

func TestGetAdminTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"should split user input tags",
			"egg,apple,apple,pineapple,car,WALL",
			[]string{"egg", "apple", "pineapple", "car", "WALL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getAdminTags(tt.input)
			assert.ElementsMatch(t, actual, tt.expected)
		})
	}
}

func TestTagDifference(t *testing.T) {
	tests := []struct {
		name    string
		tags    []*types.TagField
		oldTags []*types.TagField
		added   []string
		removed []string
	}{
		{
			"should create tags for a new business",
			[]*types.TagField{
				&types.TagField{Name: "newTag1"},
				&types.TagField{Name: "newTag2"},
			},
			nil,
			[]string{"newTag1", "newTag2"},
			[]string{},
		},
		{
			"should update tags of a existed business",
			[]*types.TagField{
				&types.TagField{Name: "oldTag1"},
				&types.TagField{Name: "oldTag2"},
				&types.TagField{Name: "newTag1"},
				&types.TagField{Name: "newTag2"},
			},
			[]*types.TagField{
				&types.TagField{Name: "oldTag1"},
				&types.TagField{Name: "oldTag2"},
				&types.TagField{Name: "oldTag3"},
			},
			[]string{"newTag1", "newTag2"},
			[]string{"oldTag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			added, removed := TagDifference(tt.tags, tt.oldTags)
			assert.ElementsMatch(t, added, tt.added)
			assert.ElementsMatch(t, removed, tt.removed)
		})
	}
}
func TestToTagFields(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"should convert tags into tag fields",
			[]string{"egg", "apple", "james", "tree"},
			[]string{"egg", "apple", "james", "tree"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := []string{}
			for _, v := range ToTagFields(tt.input) {
				actual = append(actual, v.Name)
			}
			assert.ElementsMatch(t, actual, tt.expected)
		})
	}
}

func TestGetTagNames(t *testing.T) {
	tests := []struct {
		name     string
		input    []*types.TagField
		expected []string
	}{
		{
			"should convert tags into tag fields",
			[]*types.TagField{
				&types.TagField{Name: "egg"},
				&types.TagField{Name: "apple"},
				&types.TagField{Name: "water"},
				&types.TagField{Name: "oil"},
			},
			[]string{"egg", "apple", "water", "oil"},
		},
	}

	for _, tt := range tests {
		actual := GetTagNames(tt.input)
		assert.ElementsMatch(t, actual, tt.expected)
	}
}
