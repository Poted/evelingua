package main

import (
	db "evelinqua/database"
	"evelinqua/handler"
)

func main() {

	db.GetElasticSearchClient()

	handler.HttpHandler()

}
