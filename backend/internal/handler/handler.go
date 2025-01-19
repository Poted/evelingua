package handler

import (
	"evelinqua/es"
	"evelinqua/helpers/colors"
	"evelinqua/internal/repository"
	"evelinqua/internal/service"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
	app     *fiber.App
	Version string
	routerGroup
}

type routerGroup struct {
	public fiber.Router
	auth   fiber.Router
}

func NewServer() *Server {

	server := &Server{
		app:     fiber.New(),
		Version: "1.0.0",
	}

	server.public = server.app.Group("/v1")
	server.auth = server.public.Group("/auth", JWTMiddleware)

	server.AddRoutes(
		NewAuthHandler(),
		NewWordHandler(
			service.NewWordService(
				repository.NewWordRepository(es.Client(), "1"),
			)),
	)

	return server
}

func (s *Server) Start() {

	colors.LogInColors(colors.White, "Starting server on http://localhost:4000")
	if err := s.app.Listen(":4000"); err != nil {
		colors.ErrInColors("Failed to start server: %v", err)
	}

}

type RouteConfigurator interface {
	SetupRoutes(router fiber.Router)
	SetupAuthRoutes(router fiber.Router)
}

func (s *Server) AddRoutes(handler ...RouteConfigurator) {

	for _, rc := range handler {
		rc.SetupRoutes(s.public)
		rc.SetupAuthRoutes(s.auth)
	}

}

func HelloHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func GenerateJWT(userID string, isAdmin bool) (string, error) {

	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func JWTMiddleware(c *fiber.Ctx) error {

	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid token",
		})
	}

	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid authorization format",
		})
	}
	tokenString = parts[1]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		colors.ErrInColors("invalid token", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Pass user data to the next middleware
	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_id", claims["user_id"])
	c.Locals("is_admin", claims["is_admin"])

	return c.Next()
}
