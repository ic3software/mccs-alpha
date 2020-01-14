package es

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tag struct {
	c     *elastic.Client
	index string
}

var Tag = &tag{}

func (es *tag) Register(client *elastic.Client) {
	es.c = client
	es.index = "tags"
}

func (es *tag) Create(id primitive.ObjectID, name string) error {
	body := types.TagESRecord{
		TagID: id.Hex(),
		Name:  name,
	}
	_, err := es.c.Index().
		Index(es.index).
		Id(id.Hex()).
		BodyJson(body).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *tag) UpdateOffer(id string, name string) error {
	exists, err := es.c.Exists().Index(es.index).Id(id).Do(context.TODO())
	if err != nil {
		return e.Wrap(err, "TagES UpdateOffer failed")
	}
	if !exists {
		body := types.TagESRecord{
			TagID:        id,
			Name:         name,
			OfferAddedAt: time.Now(),
		}
		_, err = es.c.Index().
			Index(es.index).
			Id(id).
			BodyJson(body).
			Do(context.Background())
		if err != nil {
			return e.Wrap(err, "TagES UpdateOffer failed")
		}
		return nil
	}

	params := map[string]interface{}{
		"offerAddedAt": time.Now(),
	}
	script := elastic.
		NewScript(`
			ctx._source.offerAddedAt = params.offerAddedAt;
		`).
		Params(params)

	_, err = es.c.Update().
		Index(es.index).
		Id(id).
		Script(script).
		Do(context.Background())
	if err != nil {
		return e.Wrap(err, "TagES UpdateOffer failed")
	}
	return nil
}

func (es *tag) UpdateWant(id string, name string) error {
	exists, err := es.c.Exists().Index(es.index).Id(id).Do(context.TODO())
	if err != nil {
		return e.Wrap(err, "TagES UpdateWant failed")
	}
	if !exists {
		body := types.TagESRecord{
			TagID:       id,
			Name:        name,
			WantAddedAt: time.Now(),
		}
		_, err = es.c.Index().
			Index(es.index).
			Id(id).
			BodyJson(body).
			Do(context.Background())
		if err != nil {
			return e.Wrap(err, "TagES UpdateWant failed")
		}
		return nil
	}

	params := map[string]interface{}{
		"wantAddedAt": time.Now(),
	}
	script := elastic.
		NewScript(`
			ctx._source.wantAddedAt = params.wantAddedAt;
		`).
		Params(params)

	_, err = es.c.Update().
		Index(es.index).
		Id(id).
		Script(script).
		Do(context.Background())
	if err != nil {
		return e.Wrap(err, "TagES UpdateWant failed")
	}
	return nil
}

func (es *tag) Rename(t *types.Tag) error {
	params := map[string]interface{}{
		"name": t.Name,
	}
	script := elastic.
		NewScript(`
			ctx._source.name = params.name;
		`).
		Params(params)

	_, err := es.c.Update().
		Index(es.index).
		Id(t.ID.Hex()).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *tag) DeleteByID(id string) error {
	_, err := es.c.Delete().
		Index(es.index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// MatchOffer matches wants for the given offer.
func (es *tag) MatchOffer(offer string, lastLoginDate time.Time) ([]string, error) {
	q := newTagQuery(offer, lastLoginDate, "wantAddedAt")
	res, err := es.c.Search().
		Index(es.index).
		Query(q).
		Do(context.Background())

	if err != nil {
		return nil, e.Wrap(err, "TagES MatchOffer failed")
	}

	matchTags := make([]string, 0, 8)
	for _, hit := range res.Hits.Hits {
		var record types.TagESRecord
		err := json.Unmarshal(hit.Source, &record)
		if err != nil {
			return nil, e.Wrap(err, "TagES MatchOffer failed")
		}
		matchTags = append(matchTags, record.Name)
	}

	return matchTags, nil
}

// MatchWant matches offers for the given want.
func (es *tag) MatchWant(want string, lastLoginDate time.Time) ([]string, error) {
	q := newTagQuery(want, lastLoginDate, "offerAddedAt")
	res, err := es.c.Search().
		Index(es.index).
		Query(q).
		Do(context.Background())

	if err != nil {
		return nil, e.Wrap(err, "TagES MatchWant failed")
	}

	matchTags := make([]string, 0, 8)
	for _, hit := range res.Hits.Hits {
		var record types.TagESRecord
		err := json.Unmarshal(hit.Source, &record)
		if err != nil {
			return nil, e.Wrap(err, "TagES MatchWant failed")
		}
		matchTags = append(matchTags, record.Name)
	}

	return matchTags, nil
}
