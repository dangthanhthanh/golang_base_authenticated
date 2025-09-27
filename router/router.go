package router

import (
	"base-app/config"
	"base-app/controller"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/gofiber/jwt/v3"
)

func SetupRoutes(app *fiber.App, userController *controller.UserController) {
	// Group API
	api := app.Group("/api/v1")

	// Hello endpoint - public, no authentication required
	api.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Golang Base Authenticated API!",
			"status":  "success",
		})
	})

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", userController.Register)
	auth.Post("/login", userController.Login)

	// User routes - require JWT
	user := api.Group("/user")

	// JWT middleware
	user.Use(jwt.New(jwt.Config{
		SigningKey:   []byte(config.LoadConfig().JWTSecret),
		ErrorHandler: jwtErrorHandler,
	}))

	user.Get("/profile", userController.GetProfile)
}

// jwtErrorHandler - xử lý lỗi nếu token không hợp lệ
func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized or invalid token",
	})
}
