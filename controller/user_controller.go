package controller

import (
	"base-app/pkg/response"
	service "base-app/service"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	service *service.UserService
}

// NewUserController tạo một controller mới
func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

// Register là endpoint để đăng ký người dùng mới
func (uc *UserController) Register(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON body
	if err := c.BodyParser(&input); err != nil {
		return response.ErrorResponse("Invalid request body", fiber.StatusBadRequest)
	}

	// Gọi service để đăng ký user
	user, err := uc.service.Register(c.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		return response.ErrorResponse(err.Error(), fiber.StatusConflict)
	}

	// Trả về user đã tạo
	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse("User registered successfully", fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}))
}

// Login là endpoint để đăng nhập và lấy JWT
func (uc *UserController) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON body
	if err := c.BodyParser(&input); err != nil {
		return response.ErrorResponse("Invalid request body", fiber.StatusBadRequest)
	}

	// Gọi service để login
	token, err := uc.service.Login(context.Background(), input.Email, input.Password)
	if err != nil {
		return response.ErrorResponse(err.Error(), fiber.StatusUnauthorized)
	}

	// Trả về token nếu thành công
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Login successful", fiber.Map{
		"token": token,
	}))
}

// GetProfile là endpoint lấy thông tin người dùng từ JWT
func (uc *UserController) GetProfile(c *fiber.Ctx) error {
	userToken := c.Locals("user")
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return response.ErrorResponse("Invalid token", fiber.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return response.ErrorResponse("Invalid token claims", fiber.StatusUnauthorized)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return response.ErrorResponse("Invalid token subject", fiber.StatusUnauthorized)
	}

	// Gọi service để lấy thông tin user
	user, err := uc.service.GetUserProfile(c.Context(), userID)
	if err != nil {
		return response.ErrorResponse(err.Error(), fiber.StatusNotFound)
	}

	// Trả về thông tin user
	return c.JSON(response.SuccessResponse("User profile fetched successfully", fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}))
}

// UpdateProfile là endpoint để cập nhật thông tin người dùng
func (uc *UserController) UpdateProfile(c *fiber.Ctx) error {
	userToken := c.Locals("user")
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return response.ErrorResponse("Invalid token", fiber.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return response.ErrorResponse("Invalid token claims", fiber.StatusUnauthorized)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return response.ErrorResponse("Invalid token subject", fiber.StatusUnauthorized)
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return response.ErrorResponse("Invalid request body", fiber.StatusBadRequest)
	}

	// Gọi service để cập nhật thông tin user
	user, err := uc.service.UpdateUserProfile(c.Context(), userID, input.Name, input.Email)
	if err != nil {
		return response.ErrorResponse(err.Error(), fiber.StatusNotFound)
	}

	// Trả về thông tin user đã được cập nhật
	return c.JSON(response.SuccessResponse("User profile updated successfully", fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}))
}

// ChangePassword là endpoint để thay đổi mật khẩu người dùng
func (uc *UserController) ChangePassword(c *fiber.Ctx) error {
	userToken := c.Locals("user")
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return response.ErrorResponse("Invalid token", fiber.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return response.ErrorResponse("Invalid token claims", fiber.StatusUnauthorized)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return response.ErrorResponse("Invalid token subject", fiber.StatusUnauthorized)
	}

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return response.ErrorResponse("Invalid request body", fiber.StatusBadRequest)
	}

	// Gọi service để thay đổi mật khẩu
	err := uc.service.ChangePassword(c.Context(), userID, input.OldPassword, input.NewPassword)
	if err != nil {
		return response.ErrorResponse(err.Error(), fiber.StatusBadRequest)
	}

	return c.JSON(response.SuccessResponse("Password updated successfully", nil))
}

// DeleteAccount là endpoint để xóa tài khoản người dùng
func (uc *UserController) DeleteAccount(c *fiber.Ctx) error {
	userToken := c.Locals("user")
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return response.ErrorResponse("Invalid token", fiber.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return response.ErrorResponse("Invalid token claims", fiber.StatusUnauthorized)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return response.ErrorResponse("Invalid token subject", fiber.StatusUnauthorized)
	}

	// Gọi service để xóa tài khoản user
	err := uc.service.ForceDeletedUserAccount(c.Context(), userID)
	if err != nil {
		return response.ErrorResponse(err.Error(), fiber.StatusNotFound)
	}

	return c.JSON(response.SuccessResponse("User account deleted successfully", nil))
}
