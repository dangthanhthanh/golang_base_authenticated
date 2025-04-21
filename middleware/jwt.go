// File: middleware/jwt.go
package middleware

import (
	"base-app/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	user, err := utils.ParseJWT(tokenStr)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	c.Locals("user", user)
	return c.Next()
}
