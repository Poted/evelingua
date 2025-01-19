package handler

import (
	"evelinqua/internal/repository"
	"evelinqua/internal/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// region setup

type wordHandler struct {
	prefix  string
	service *service.WordService
}

func NewWordHandler(prefix string, service *service.WordService) *wordHandler {
	return &wordHandler{
		prefix: prefix,
	}
}

func (h *wordHandler) SetupRoutes(router fiber.Router) {
	wg := router.Group(h.prefix)
	wg.Get("/add", h.AddWord)
}

func (h *wordHandler) SetupAuthRoutes(router fiber.Router) {
	wg := router.Group(h.prefix)
	wg.Get("/word", HelloHandler)
}

// endregion setup
// region methods

func (h *wordHandler) AddWord(c *fiber.Ctx) error {

	fmt.Printf("\"D::DD\": %v\n", "D::DD")

	var word repository.Word
	if err := c.BodyParser(&word); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.service.AddWord(word); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add word"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Word added successfully"})
}

func (h *wordHandler) SearchWords(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter is required"})
	}

	words, err := h.service.SearchWords(query, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to search words"})
	}

	return c.JSON(words)
}

// endregion methods
