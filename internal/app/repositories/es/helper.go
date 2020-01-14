package es

import (
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

func newFuzzyWildcardQuery(name, text string) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()
	q.Should(elastic.NewMatchQuery(name, text).Fuzziness("auto").Boost(3))
	q.Should(elastic.NewWildcardQuery(name, strings.ToLower(text)+"*").Boost(2))
	q.Should(elastic.NewRegexpQuery(name, ".*"+strings.ToLower(text)+".*"))
	return q
}

// Should match one of the three queries. (MatchQuery, WildcardQuery, RegexpQuery)
func newFuzzyWildcardTimeQueryForTag(tagField, tagName string, createdOnOrAfter time.Time) *elastic.NestedQuery {
	q := elastic.NewBoolQuery()

	qq := elastic.NewBoolQuery()
	qq.Must((elastic.NewMatchQuery(tagField+".name", tagName).Fuzziness("auto").Boost(2)))
	qq.Must(elastic.NewRangeQuery(tagField + ".createdAt").Gte(createdOnOrAfter))
	q.Should(qq)

	qq = elastic.NewBoolQuery()
	qq.Must(elastic.NewWildcardQuery(tagField+".name", strings.ToLower(tagName)+"*").Boost(1.5))
	qq.Must(elastic.NewRangeQuery(tagField + ".createdAt").Gte(createdOnOrAfter))
	q.Should(qq)

	qq = elastic.NewBoolQuery()
	qq.Must(elastic.NewRegexpQuery(tagField+".name", ".*"+strings.ToLower(tagName)+".*"))
	qq.Must(elastic.NewRangeQuery(tagField + ".createdAt").Gte(createdOnOrAfter))
	q.Should(qq)

	nestedQ := elastic.NewNestedQuery(tagField, q)

	return nestedQ
}

// Should match one of the three queries. (MatchQuery, WildcardQuery, RegexpQuery)
func newTagQuery(tag string, lastLoginDate time.Time, timeField string) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()

	// The default value for both offerAddedAt and wantAddedAt is 0001-01-01T00:00:00.000+0000.
	// If the user never login before and his lastLoginDate will be 0001-01-01T00:00:00.000+0000.
	// And we will match the user's own tags.
	// Added this filter to solve the problem.
	q.MustNot(elastic.NewRangeQuery(timeField).Lte(time.Time{}))

	qq := elastic.NewBoolQuery()
	qq.Must(elastic.NewMatchQuery("name", tag).Fuzziness("auto"))
	qq.Must(elastic.NewRangeQuery(timeField).Gte(lastLoginDate))
	q.Should(qq)

	qq = elastic.NewBoolQuery()
	qq.Must(elastic.NewWildcardQuery("name", strings.ToLower(tag)+"*"))
	qq.Must(elastic.NewRangeQuery(timeField).Gte(lastLoginDate))
	q.Should(qq)

	qq = elastic.NewBoolQuery()
	qq.Must(elastic.NewRegexpQuery("name", ".*"+strings.ToLower(tag)+".*"))
	qq.Must(elastic.NewRangeQuery(timeField).Gte(lastLoginDate))
	q.Should(qq)

	return q
}
