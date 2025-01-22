package handler

import (
	"encoding/json"
	"evelinqua/es"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/gofiber/fiber/v2"
)

type categoryHandler struct {
	prefix string
}

func NewCategoryHandler(prefix string) *categoryHandler {
	return &categoryHandler{prefix: prefix}
}

func (h *categoryHandler) SetupRoutes(router *fiber.App) {
	router.Post("/category", h.AddCategoryHandler)
	router.Get("/category/words", h.GetWordsInCategoryHandler)
	router.Post("/category/check-translation", h.CheckTranslationHandler)
}

func (h *categoryHandler) SetupAuthRoutes(router *fiber.App) {
	router.Post("/category", h.AddCategoryHandler)
	router.Get("/category/words", h.GetWordsInCategoryHandler)
	router.Post("/category/check-translation", h.CheckTranslationHandler)
}

func (h *categoryHandler) AddCategoryHandler(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	var payload struct {
		Category string `json:"category"`
		Words    []struct {
			Word        string `json:"word"`
			Language    string `json:"language"`
			Translation string `json:"translation"`
		} `json:"words"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid payload")
	}

	doc := map[string]interface{}{
		"userId":   userID,
		"category": payload.Category,
		"words":    payload.Words,
	}

	res, err := es.Client().Index("categories", esutil.NewJSONReader(doc))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to add category")
	}
	defer res.Body.Close()

	return c.Status(fiber.StatusCreated).SendString("Category added successfully")
}

func (h *categoryHandler) GetWordsInCategoryHandler(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)
	category := c.Query("category")
	if category == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Category is required")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match": map[string]interface{}{"userId": userID}},
					{"match": map[string]interface{}{"category": category}},
				},
			},
		},
	}

	res, err := es.Client().Search(
		es.Client().Search.WithIndex("categories"),
		es.Client().Search.WithBody(esutil.NewJSONReader(query)),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch words")
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse response")
	}

	return c.JSON(result)
}

func (h *categoryHandler) CheckTranslationHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var payload struct {
		Category    string `json:"category"`
		Word        string `json:"word"`
		Translation string `json:"translation"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid payload")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match": map[string]interface{}{"userId": userID}},
					{"match": map[string]interface{}{"category": payload.Category}},
					{"nested": map[string]interface{}{
						"path": "words",
						"query": map[string]interface{}{
							"bool": map[string]interface{}{
								"must": []map[string]interface{}{
									{"match": map[string]interface{}{"words.word": payload.Word}},
									{"match": map[string]interface{}{"words.translation": payload.Translation}},
								},
							},
						},
					}},
				},
			},
		},
	}

	res, err := es.Client().Search(
		es.Client().Search.WithIndex("categories"),
		es.Client().Search.WithBody(esutil.NewJSONReader(query)),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to validate translation")
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse response")
	}

	if len(result["hits"].(map[string]interface{})["hits"].([]interface{})) > 0 {
		return c.SendString("Correct translation!")
	}

	return c.SendString("Incorrect translation.")
}
