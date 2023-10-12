package es

import (
	"context"
	"encoding/json"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/ic3network/mccs-alpha/internal/pkg/pagination"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	c     *elastic.Client
	index string
}

var User = &user{}

func (es *user) Register(client *elastic.Client) {
	es.c = client
	es.index = "users"
}

// Create creates an UserESRecord in Elasticsearch.
func (es *user) Create(u *types.User) error {
	body := types.UserESRecord{
		UserID:    u.ID.Hex(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	_, err := es.c.Index().
		Index(es.index).
		Id(u.ID.Hex()).
		BodyJson(body).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// Find finds users from Elasticsearch.
func (es *user) Find(u *types.User, page int64) ([]string, int, int, error) {
	if page < 0 || page == 0 {
		return nil, 0, 0, e.New(e.InvalidPageNumber, "find user failed")
	}

	var ids []string
	size := viper.GetInt("page_size")
	from := viper.GetInt("page_size") * (int(page) - 1)

	q := elastic.NewBoolQuery()

	if u.LastName != "" {
		q.Must(newFuzzyWildcardQuery("lastName", u.LastName))
	}
	if u.Email != "" {
		q.Must(newFuzzyWildcardQuery("email", u.Email))
	}

	res, err := es.c.Search().
		Index(es.index).
		From(from).
		Size(size).
		Query(q).
		Do(context.Background())

	if err != nil {
		return nil, 0, 0, e.Wrap(err, "find user failed")
	}

	for _, hit := range res.Hits.Hits {
		var record types.UserESRecord
		err := json.Unmarshal(hit.Source, &record)
		if err != nil {
			return nil, 0, 0, e.Wrap(err, "find user failed")
		}
		ids = append(ids, record.UserID)
	}

	numberOfResults := res.Hits.TotalHits.Value
	totalPages := pagination.Pages(numberOfResults, viper.GetInt64("page_size"))

	return ids, int(numberOfResults), totalPages, nil
}

func (es *user) Update(u *types.User) error {
	doc := map[string]interface{}{
		"email":     u.Email,
		"firstName": u.FirstName,
		"lastName":  u.LastName,
	}

	_, err := es.c.Update().
		Index(es.index).
		Id(u.ID.Hex()).
		Doc(doc).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *user) UpdateTradingInfo(
	id primitive.ObjectID,
	data *types.TradingRegisterData,
) error {
	doc := map[string]interface{}{
		"firstName": data.FirstName,
		"lastName":  data.LastName,
	}
	_, err := es.c.Update().
		Index(es.index).
		Id(id.Hex()).
		Doc(doc).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *user) Delete(id string) error {
	_, err := es.c.Delete().
		Index(es.index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
