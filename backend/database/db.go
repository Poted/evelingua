package db

import (
	"log"
	"strings"

	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client
var once sync.Once

func GetElasticSearchClient() {

	once.Do(func() {
		var err error
		newConnection, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{"http://elasticsearch:9200"},
		})
		// newConnection, err := elasticsearch.NewDefaultClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		ES = newConnection
	})

	res, err := ES.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

}

func CreateIndex(indexName string) {

	req := strings.NewReader(`{
        "mappings": {
            "properties": {
                "word": { "type": "text" },
                "language": { "type": "text" },
                "translation": { "type": "text" }
            }
        }
    }`)

	res, err := ES.Indices.Create(indexName, ES.Indices.Create.WithBody(req))
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error response from Elasticsearch: %s", res.String())
	} else {
		log.Printf("Index %s created successfully", indexName)
	}
}
