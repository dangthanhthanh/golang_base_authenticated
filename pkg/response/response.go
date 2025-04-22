package response

import "github.com/gofiber/fiber/v2"

// SuccessResponse returns a standard success structure
func SuccessResponse(message string, data interface{}) fiber.Map {
	return fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	}
}

// ErrorResponse returns a standard error structure
func ErrorResponse(message string, status int) error {
	return fiber.NewError(status, message)
}

// CustomResponse allows for full manual control if needed
func CustomResponse(success bool, message string, data interface{}) fiber.Map {
	return fiber.Map{
		"success": success,
		"message": message,
		"data":    data,
	}
}
