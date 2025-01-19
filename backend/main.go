package main

import (
	"evelinqua/es"
	"evelinqua/helpers/colors"
	"evelinqua/internal/handler"
	"flag"
	"fmt"

	"github.com/Poted/getenv"
)

func main() {

	// Load the environment variables
	err := getenv.LoadEnv(".env")
	if err != nil {
		colors.ErrInColors("Error loading environment variables: ", err)
	}

	runElasticSearch := flag.Bool("es", true, "Run the elasticsearch connection (default: true)")
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			App(*runElasticSearch)

		}
	}()

	App(*runElasticSearch)
}

func App(runElasticSearch bool) {

	// Start the listener for managing app while running
	// go listener.Listen()

	// Start the Elastic Search connection
	if runElasticSearch {
		es.ElasticSearchConnection()
	}

	// Start the HTTP handler
	handler.NewServer().Start()

}
