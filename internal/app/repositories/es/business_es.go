package es

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/pagination"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type business struct {
	c     *elastic.Client
	index string
}

var Business = &business{}

func (es *business) Register(client *elastic.Client) {
	es.c = client
	es.index = "businesses"
}

func (es *business) Create(id primitive.ObjectID, data *types.BusinessData) error {
	body := types.BusinessESRecord{
		BusinessID:      id.Hex(),
		BusinessName:    data.BusinessName,
		Offers:          data.Offers,
		Wants:           data.Wants,
		LocationCity:    data.LocationCity,
		LocationCountry: data.LocationCountry,
		Status:          constant.Business.Pending,
		AdminTags:       data.AdminTags,
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

func matchStatuses(q *elastic.BoolQuery, c *types.SearchCriteria) {
	if len(c.Statuses) != 0 {
		qq := elastic.NewBoolQuery()
		for _, status := range c.Statuses {
			qq.Should(elastic.NewMatchQuery("status", status))
		}
		q.Must(qq)
	}
}

func matchTags(q *elastic.BoolQuery, c *types.SearchCriteria) {
	// "Tag Added After" will associate with "tags".
	if c.TagType == constant.OFFERS && len(c.Tags) != 0 {
		qq := elastic.NewBoolQuery()
		// weighted is used to make sure the tags are shown in order.
		weighted := 2.0
		for _, o := range c.Tags {
			qq.Should(newFuzzyWildcardTimeQueryForTag("offers", o.Name, c.CreatedOnOrAfter).
				Boost(weighted))
			weighted *= 0.9
		}
		// Must match one of the "Should" queries.
		q.Must(qq)
	} else if c.TagType == constant.WANTS && len(c.Tags) != 0 {
		qq := elastic.NewBoolQuery()
		// weighted is used to make sure the tags are shown in order.
		weighted := 2.0
		for _, w := range c.Tags {
			qq.Should(newFuzzyWildcardTimeQueryForTag("wants", w.Name, c.CreatedOnOrAfter).
				Boost(weighted))
			weighted *= 0.9
		}
		// Must match one of the "Should" queries.
		q.Must(qq)
	}
}

func (es *business) Find(c *types.SearchCriteria, page int64) ([]string, int, int, error) {
	if page < 0 || page == 0 {
		return nil, 0, 0, e.New(e.InvalidPageNumber, "find business failed")
	}

	var ids []string
	size := viper.GetInt("page_size")
	from := viper.GetInt("page_size") * (int(page) - 1)

	q := elastic.NewBoolQuery()

	q.Should(elastic.NewMatchQuery("status", constant.Trading.Accepted))

	if c.ShowUserFavoritesOnly {
		idQuery := elastic.NewIdsQuery().Ids(util.ToIDStrings(c.FavoriteBusinesses)...)
		q.Must(idQuery)
	}

	matchStatuses(q, c)

	if c.BusinessName != "" {
		q.Must(newFuzzyWildcardQuery("businessName", c.BusinessName))
	}
	if c.LocationCountry != "" {
		q.Must(elastic.NewMatchQuery("locationCountry", c.LocationCountry))
	}
	if c.LocationCity != "" {
		q.Must(newFuzzyWildcardQuery("locationCity", c.LocationCity))
	}

	if c.AdminTag != "" {
		q.Must(elastic.NewMatchQuery("adminTags", c.AdminTag))
	}

	matchTags(q, c)

	res, err := es.c.Search().
		Index(es.index).
		From(from).
		Size(size).
		Query(q).
		Do(context.Background())

	if err != nil {
		return nil, 0, 0, e.Wrap(err, "BusinessES Find failed")
	}

	for _, hit := range res.Hits.Hits {
		var record types.BusinessESRecord
		err := json.Unmarshal(hit.Source, &record)
		if err != nil {
			return nil, 0, 0, e.Wrap(err, "BusinessES Find failed")
		}
		ids = append(ids, record.BusinessID)
	}

	numberOfResults := res.Hits.TotalHits.Value
	totalPages := pagination.Pages(numberOfResults, viper.GetInt64("page_size"))

	return ids, int(numberOfResults), totalPages, nil
}

func (es *business) UpdateBusiness(id primitive.ObjectID, data *types.BusinessData) error {
	params := map[string]interface{}{
		"businessName":    data.BusinessName,
		"locationCity":    data.LocationCity,
		"locationCountry": data.LocationCountry,
		"offersAdded":     helper.ToTagFields(data.OffersAdded),
		"wantsAdded":      helper.ToTagFields(data.WantsAdded),
		"offersRemoved":   data.OffersRemoved,
		"wantsRemoved":    data.WantsRemoved,
	}
	if data.Status != "" {
		params["status"] = data.Status
	}
	params["adminTags"] = data.AdminTags

	script := elastic.
		NewScript(`
			ctx._source.businessName = params.businessName;
			ctx._source.locationCity = params.locationCity;
			ctx._source.locationCountry = params.locationCountry;

			if (params.status !== null) {
				ctx._source.status = params.status;
			}

			if (params.adminTags !== null) {
				if (params.adminTags.length !== 0) {
					ctx._source.adminTags = params.adminTags;
				} else {
					ctx._source.adminTags = [];
				}
			}

			if (params.offersRemoved !== null && params.offersRemoved.length !== 0) {
				for (int i = 0; i < ctx._source.offers.length; i++) {
					if (params.offersRemoved.contains(ctx._source.offers[i].name)) {
						ctx._source.offers.remove(i);
						i--
					}
				}
			}
			if (params.wantsRemoved !== null && params.wantsRemoved.length !== 0) {
				for (int i = 0; i < ctx._source.wants.length; i++) {
					if (params.wantsRemoved.contains(ctx._source.wants[i].name)) {
						ctx._source.wants.remove(i);
						i--
					}
				}
			}

			if (params.offersAdded !== null && params.offersAdded.length !== 0) {
				for (int i = 0; i < params.offersAdded.length; i++) {
					ctx._source.offers.add(params.offersAdded[i]);
				}
			}
			if (params.wantsAdded !== null && params.wantsAdded.length !== 0) {
				for (int i = 0; i < params.wantsAdded.length; i++) {
					ctx._source.wants.add(params.wantsAdded[i]);
				}
			}
		`).
		Params(params)

	_, err := es.c.Update().
		Index(es.index).
		Id(id.Hex()).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *business) UpdateTradingInfo(id primitive.ObjectID, data *types.TradingRegisterData) error {
	doc := map[string]interface{}{
		"businessName":    data.BusinessName,
		"locationCity":    data.LocationCity,
		"locationCountry": data.LocationCountry,
		"status":          constant.Trading.Pending,
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

func (es *business) UpdateAllTagsCreatedAt(id primitive.ObjectID, t time.Time) error {
	params := map[string]interface{}{
		"createdAt": t,
	}

	script := elastic.
		NewScript(`
			for (int i = 0; i < ctx._source.offers.length; i++) {
				ctx._source.offers[i].createdAt = params.createdAt
			}
			for (int i = 0; i < ctx._source.wants.length; i++) {
				ctx._source.wants[i].createdAt = params.createdAt
			}
		`).
		Params(params)

	_, err := es.c.Update().
		Index(es.index).
		Id(id.Hex()).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *business) RenameTag(old string, new string) error {
	query := elastic.NewBoolQuery()
	query.Should(elastic.NewMatchQuery("offers.name", old))
	query.Should(elastic.NewMatchQuery("wants.name", old))
	script := elastic.
		NewScript(`
			for (int i = 0; i < ctx._source.offers.length; i++) {
				if (ctx._source.offers[i].name == params.old) {
					ctx._source.offers[i].name = params.new
				}
			}
			for (int i = 0; i < ctx._source.wants.length; i++) {
				if (ctx._source.wants[i].name == params.old) {
					ctx._source.wants[i].name = params.new
				}
			}
		`).
		Params(map[string]interface{}{"new": new, "old": old})
	_, err := es.c.UpdateByQuery(es.index).
		Query(query).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *business) RenameAdminTag(old string, new string) error {
	query := elastic.NewMatchQuery("adminTags", old)
	script := elastic.
		NewScript(`
			if (ctx._source.adminTags.contains(params.old)) {
				ctx._source.adminTags.remove(ctx._source.adminTags.indexOf(params.old));
				ctx._source.adminTags.add(params.new);
			}
		`).
		Params(map[string]interface{}{"new": new, "old": old})
	_, err := es.c.UpdateByQuery(es.index).
		Query(query).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *business) Delete(id string) error {
	_, err := es.c.Delete().
		Index(es.index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *business) DeleteTag(name string) error {
	query := elastic.NewBoolQuery()
	query.Should(elastic.NewMatchQuery("offers.name", name))
	query.Should(elastic.NewMatchQuery("wants.name", name))
	script := elastic.
		NewScript(`
			for (int i = 0; i < ctx._source.offers.length; i++) {
				if (ctx._source.offers[i].name == params.name) {
					ctx._source.offers.remove(i);
					break;
				}
			}
			for (int i = 0; i < ctx._source.wants.length; i++) {
				if (ctx._source.wants[i].name == params.name) {
					ctx._source.wants.remove(i);
					break;
				}
			}
		`).
		Params(map[string]interface{}{"name": name})
	_, err := es.c.UpdateByQuery(es.index).
		Query(query).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (es *business) DeleteAdminTags(name string) error {
	query := elastic.NewMatchQuery("adminTags", name)
	script := elastic.
		NewScript(`
			if (ctx._source.adminTags.contains(params.name)) {
				ctx._source.adminTags.remove(ctx._source.adminTags.indexOf(params.name));
			}
		`).
		Params(map[string]interface{}{"name": name})
	_, err := es.c.UpdateByQuery(es.index).
		Query(query).
		Script(script).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
