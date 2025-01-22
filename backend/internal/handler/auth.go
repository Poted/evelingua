package handler

import (
	"evelinqua/helpers/colors"
	"evelinqua/internal/service"

	"github.com/gofiber/fiber/v2"
)

// region setup

type authHandler struct {
	prefix  string
	service *service.AuthService
}

func NewAuthHandler(prefix string, service *service.AuthService) *authHandler {
	return &authHandler{
		prefix:  prefix,
		service: service,
	}
}

func (h *authHandler) SetupRoutes(router fiber.Router) {
	router.Get("/login", h.Login)
	router.Post("/register", h.Register)
	router.Post("/logout", h.Logout)
	router.Get("/hello", HelloHandler)
}

func (h *authHandler) SetupAuthRoutes(router fiber.Router) {
	router.Get("/after-login", HelloHandler)
}

// endregion setup
// region methods

func (h *authHandler) Login(c *fiber.Ctx) error {

	// userID, isAdmin
	token, err := GenerateJWT("user123", true)
	if err != nil {
		colors.ErrInColors("failed to create jwt-token", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}
	return c.JSON(fiber.Map{"token": token})
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	return c.SendString("Registered")
}

func (h *authHandler) Logout(c *fiber.Ctx) error {
	return c.SendString("Logged out")
}

// endregion methods
