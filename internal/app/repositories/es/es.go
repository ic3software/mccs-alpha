package es

import (
	"log"
	"time"

	"github.com/ic3network/mccs-alpha/global"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

var client *elastic.Client

func init() {
	global.Init()
	client = New()
	registerCollections(client)
}

func registerCollections(client *elastic.Client) {
	Business.Register(client)
	User.Register(client)
	Tag.Register(client)
}

// New returns an initialized ES instance.
func New() *elastic.Client {
	var client *elastic.Client
	var err error

	for {
		client, err = elastic.NewClient(
			elastic.SetURL(viper.GetString("es.url")),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Printf("ElasticSearch connection error: %+v \n", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	checkIndexes(client)
	return client
}

// Client is for seed/restore data
func Client() *elastic.Client {
	return client
}
