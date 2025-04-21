package router

import (
	"base-app/controller"

	"github.com/gofiber/fiber/v2"
)

func LogRoutes(app *fiber.App, userController *controller.UserController) {
	// Group API theo version
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Nhóm user routes
	user := v1.Group("/user")
	user.Post("/register", userController.Register)
	user.Post("/login", userController.Login)
	user.Get("/profile", userController.GetProfile) // cần middleware Auth nếu có

	// Route show tất cả routes
	app.Get("/routes", func(c *fiber.Ctx) error {
		routes := app.GetRoutes()
		result := make([]fiber.Map, 0, len(routes))

		for _, route := range routes {
			result = append(result, fiber.Map{
				"method": route.Method,
				"path":   route.Path,
				"name":   route.Name,
			})
		}

		return c.JSON(fiber.Map{
			"routes": result,
		})
	})
}
