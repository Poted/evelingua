package main

import (
	"evelinqua/es"
	"evelinqua/handler"
	"evelinqua/listener"
	"flag"
	"fmt"
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

	wait := make(chan bool)

	defer close(wait)

	go listener.Listen()

	if runElasticSearch {
		es.ElasticSearchConnection()
	}

	handler.HttpHandler()

	<-wait
}
