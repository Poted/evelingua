package es

import (
	"evelinqua/helpers/colors"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	es   *elasticsearch.Client
	once sync.Once
)

func Client() *elasticsearch.Client {
	return es
}

func ElasticSearchConnection() {

	defer func() {
		if r := recover(); r != nil {
			colors.FuncInColors(colors.Yellow, func() {
				log.Default().Printf("Recovered from panic while connecting to Elastic Search server: %s\n\n", r)
				debug.PrintStack()
				println("\n")
			})

		}
	}()

	colors.LogInColors(colors.Cyan, "Connecting to Elasticsearch")

	once.Do(func() {
		var err error
		address := fmt.Sprintf("http://%s:%s", os.Getenv("ELASTICSEARCH_HOST"), os.Getenv("ELASTICSEARCH_PORT"))
		newConnection, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{address},
		})
		if err != nil {
			colors.ErrInColors("error creating the client: " + err.Error())
		}

		es = newConnection
	})

	if es == nil {
		colors.ErrInColors("elasticsearch connection failed")
		return
	}

	res, err := es.Info()
	if err != nil {
		colors.ErrInColors("Error getting test response from elastic search: " + err.Error())
	}
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()

	if res != nil {
		colors.LogInColors(colors.Green, "Elasticsearch connection successful")
	}

}
