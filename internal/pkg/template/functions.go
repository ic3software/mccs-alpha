package template

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func arrToSting(stringArray []string) string {
	return strings.Join(stringArray, ",")
}

func tagsToString(tags []*types.TagField) string {
	var sb strings.Builder
	l := len(tags) - 1
	for i, t := range tags {
		sb.WriteString(t.Name)
		if i != l {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func tagsToSearchString(tags []*types.TagField) string {
	var sb strings.Builder
	l := len(tags) - 1
	for i, t := range tags {
		sb.WriteString(t.Name)
		if i != l {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func add(number int, inc int) int {
	return number + inc
}

func minus(number int, inc int) int {
	return number - inc
}

func n(start int, end int) []int {
	numbers := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		numbers = append(numbers, i)
	}
	return numbers
}

func idToString(id primitive.ObjectID) string {
	return id.Hex()
}

func formatTime(t time.Time) string {
	tt := t.UTC()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d UTC",
		tt.Year(), tt.Month(), tt.Day(),
		tt.Hour(), tt.Minute(), tt.Second())
}

func formatAccountBalance(balance float64) string {
	return fmt.Sprintf("%.2f", balance)
}

func formatTransactionID(id string) string {
	return id[0:8]
}

func shouldDisplayTime(t time.Time) bool {
	return !t.IsZero()
}

func includesID(list []primitive.ObjectID, target primitive.ObjectID) bool {
	for _, id := range list {
		if id == target {
			return true
		}
	}
	return false
}

func timeNow() string {
	return time.Now().Format("2006-01-02")
}

func daysBefore(days int) string {
	return time.Now().AddDate(0, 0, -days).Format("2006-01-02")
}

func sortAdminTags(tags []*types.AdminTag) []*types.AdminTag {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})
	return tags
}

func containPrefix(arr []string, prefix string) bool {
	for _, s := range arr {
		if strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(prefix)) {
			return true
		}
	}
	return false
}
