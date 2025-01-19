package handler

import (
	"github.com/gofiber/fiber/v2"
)

type authHandler struct{}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}

func (h *authHandler) SetupRoutes(router fiber.Router) {
	router.Get("/login", h.Login)
	router.Post("/register", h.Register)
	router.Post("/logout", h.Logout)
	router.Get("/hello", HelloHandler)
}

func (h *authHandler) SetupAdminRoutes(router fiber.Router) {
	router.Get("/login-admin", HelloHandler)
}

func (h *authHandler) Login(c *fiber.Ctx) error {

	token, err := GenerateJWT("user123", true)
	if err != nil {
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
