package helper

import (
	"regexp"
	"strings"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/types"
)

var (
	specialCharRe *regexp.Regexp
	multiDashRe   *regexp.Regexp
	ltDashRe      *regexp.Regexp
	adminTagRe    *regexp.Regexp
)

func init() {
	specialCharRe = regexp.MustCompile("(&quot;)|([^a-zA-Z-]+)")
	multiDashRe = regexp.MustCompile("-+")
	ltDashRe = regexp.MustCompile("(^-+)|(-+$)")
	adminTagRe = regexp.MustCompile("[0-9]|(&quot;)|([^a-zA-Z ]+)")
}

// GetTags transforms tags from the user inputs into a standard format.
// dog walking -> dog-walking (one word)
func GetTags(words string) []*types.TagField {
	splitFn := func(c rune) bool {
		return c == ','
	}
	tagArray := strings.FieldsFunc(strings.ToLower(words), splitFn)

	encountered := map[string]bool{}
	tags := make([]*types.TagField, 0, len(tagArray))
	for _, tag := range tagArray {
		tag = strings.Replace(tag, " ", "-", -1)
		tag = specialCharRe.ReplaceAllString(tag, "")
		tag = multiDashRe.ReplaceAllString(tag, "-")
		tag = ltDashRe.ReplaceAllString(tag, "")
		if len(tag) == 0 {
			continue
		}
		// remove duplicates
		if !encountered[tag] {
			tags = append(tags, &types.TagField{
				Name:      tag,
				CreatedAt: time.Now(),
			})
			encountered[tag] = true
		}
	}
	return tags
}

// ToSearchTags transforms tags from user inputs into searching tags.
// dog walking -> dog, walking (two words)
func ToSearchTags(words string) []*types.TagField {
	splitFn := func(c rune) bool {
		return c == ',' || c == ' '
	}
	tagArray := strings.FieldsFunc(strings.ToLower(words), splitFn)

	encountered := map[string]bool{}
	tags := make([]*types.TagField, 0, len(tagArray))
	for _, tag := range tagArray {
		tag = strings.Replace(tag, " ", "-", -1)
		tag = specialCharRe.ReplaceAllString(tag, "")
		tag = multiDashRe.ReplaceAllString(tag, "-")
		tag = ltDashRe.ReplaceAllString(tag, "")
		if len(tag) == 0 {
			continue
		}
		// remove duplicates
		if !encountered[tag] {
			tags = append(tags, &types.TagField{
				Name:      tag,
				CreatedAt: time.Now(),
			})
			encountered[tag] = true
		}
	}
	return tags
}

func getAdminTags(words string) []string {
	splitFn := func(c rune) bool {
		return c == ','
	}
	tags := make([]string, 0, 8)
	encountered := map[string]bool{}
	for _, tag := range strings.FieldsFunc(words, splitFn) {
		tag = adminTagRe.ReplaceAllString(tag, "")
		// remove duplicates
		if !encountered[tag] {
			tags = append(tags, tag)
			encountered[tag] = true
		}
	}
	return tags
}

func FormatAdminTag(tag string) string {
	return adminTagRe.ReplaceAllString(tag, "")
}

// TagDifference finds out the new added tags.
func TagDifference(tags, oldTags []*types.TagField) ([]string, []string) {
	encountered := map[string]int{}
	added := []string{}
	removed := []string{}
	for _, tag := range oldTags {
		if _, ok := encountered[tag.Name]; !ok {
			encountered[tag.Name]++
		}
	}
	for _, tag := range tags {
		encountered[tag.Name]--
	}
	for name, flag := range encountered {
		if flag == -1 {
			added = append(added, name)
		}
		if flag == 1 {
			removed = append(removed, name)
		}
	}
	return added, removed
}

func SameTags(new, old []*types.TagField) bool {
	added, removed := TagDifference(new, old)
	if len(added)+len(removed) > 0 {
		return false
	}
	return true
}

// ToTagFields converts tags into TagFields.
func ToTagFields(tags []string) []*types.TagField {
	tagFields := make([]*types.TagField, 0, len(tags))
	for _, tagName := range tags {
		tagField := &types.TagField{
			Name:      tagName,
			CreatedAt: time.Now(),
		}
		tagFields = append(tagFields, tagField)
	}
	return tagFields
}

// GetTagNames gets tag name from TagField.
func GetTagNames(tags []*types.TagField) []string {
	names := make([]string, 0, len(tags))
	for _, t := range tags {
		names = append(names, t.Name)
	}
	return names
}

// GetAdminTagNames gets admin tag name from AdminTagField.
func GetAdminTagNames(tags []*types.AdminTag) []string {
	names := make([]string, 0, len(tags))
	for _, t := range tags {
		names = append(names, t.Name)
	}
	return names
}
