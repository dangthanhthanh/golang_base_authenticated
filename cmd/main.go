package main

import (
	"log"

	"base-app/config"
	"base-app/controller"
	"base-app/pkg/db"
	"base-app/pkg/redis"
	"base-app/repository"
	"base-app/router"
	"base-app/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load config từ .env
	cfg := config.LoadConfig()

	// Kết nối PostgreSQL và Redis
	db.Connect(cfg)
	db.Migrate()
	redis.Connect(cfg)

	// Khởi tạo tầng repository, service, controller
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo, redis.RDB, cfg)
	userController := controller.NewUserController(userService)

	// Khởi tạo Fiber app
	app := fiber.New()

	// Cấu hình routes
	// router.LogRoutes(app, userController)
	router.SetupRoutes(app, userController)

	// xuat router.
	for _, route := range app.GetRoutes() {
		log.Printf("📍 %-6s %s", route.Method, route.Path)
	}

	// Chạy server
	log.Printf("🚀 Server is running on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
