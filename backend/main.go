package main

import (
	"evelinqua/es"
	"evelinqua/handler"
	"evelinqua/listener"
	"flag"
	"fmt"

	"github.com/Poted/getenv"
)

func main() {

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

	// This channel is preventing from returning while restarting handler
	wait := make(chan bool)
	defer close(wait)

	// Load the environment variables
	getenv.LoadEnv(".env")

	// Start the listener for managing app while running
	go listener.Listen()

	// Start the Elastic Search connection
	if runElasticSearch {
		es.ElasticSearchConnection()
	}

	// Start the HTTP handler
	handler.HttpHandler()

	<-wait
}
