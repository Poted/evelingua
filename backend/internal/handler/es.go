package handler

import (
	"evelinqua/es"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Word struct {
	Word        string `json:"word"`
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

func ElasticSearchHandler(app *fiber.App) {

	app.Get("/elasticsearch", CheckConnectionHandler)
	app.Post("/addword", AddWordHandler)
	app.Get("/getword", GetWordHandler)

}

func CheckConnectionHandler(c *fiber.Ctx) error {
	resp, err := http.Get("http://localhost:9200/")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to connect to Elasticsearch")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response from Elasticsearch")
	}

	return c.JSON(string(body))
}

func AddWordHandler(c *fiber.Ctx) error {
	var word Word
	if err := c.BodyParser(&word); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	req := strings.NewReader(fmt.Sprintf(`{
        "word": "%s",
        "language": "%s",
        "translation": "%s"
    }`, word.Word, word.Language, word.Translation))

	res, err := es.Client().Index("words", req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to add word to Elasticsearch")
	}
	defer res.Body.Close()

	if res.IsError() {
		return c.Status(fiber.StatusInternalServerError).SendString("Error response from Elasticsearch")
	}

	return c.Status(fiber.StatusCreated).SendString("Word added successfully")
}

func GetWordHandler(c *fiber.Ctx) error {
	word := c.Query("word")
	if word == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing word parameter")
	}

	query := strings.NewReader(fmt.Sprintf(`{
        "query": {
            "match": {
                "word": "%s"
            }
        }
    }`, word))

	res, err := es.Client().Search(
		es.Client().Search.WithIndex("words"),
		es.Client().Search.WithBody(query),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to search for word in Elasticsearch")
	}
	defer res.Body.Close()

	if res.IsError() {
		return c.Status(fiber.StatusInternalServerError).SendString("Error response from Elasticsearch")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response from Elasticsearch")
	}

	fmt.Printf("string(body): %v\n", string(body))

	return c.JSON(string(body))
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

	res, err := es.Client().Indices.Create(indexName, es.Client().Indices.Create.WithBody(req))
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
