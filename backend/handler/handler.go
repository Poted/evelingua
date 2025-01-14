package handler

import (
	db "evelinqua/database"
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

func HttpHandler() {

	app := fiber.New()

	app.Get("/", HelloHandler)
	app.Get("/elasticsearch", ElasticSearchHandler)
	app.Post("/addword", AddWordHandler)
	app.Get("/getword", GetWordHandler)

	log.Println("Starting server on http://localhost:4000")

	err := app.Listen(":4000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func HelloHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func ElasticSearchHandler(c *fiber.Ctx) error {
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

	res, err := db.ES.Index("words", req)
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

	res, err := db.ES.Search(
		db.ES.Search.WithIndex("words"),
		db.ES.Search.WithBody(query),
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
