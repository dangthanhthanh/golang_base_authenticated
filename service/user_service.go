package service

import (
	"base-app/config"
	"base-app/model"
	"base-app/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  repository.UserRepository
	redis *redis.Client
	cfg   config.Config
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client, cfg config.Config) *UserService {
	return &UserService{
		repo:  repo,
		redis: redisClient,
		cfg:   cfg,
	}
}

// Register - Đăng ký người dùng mới
func (s *UserService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Mã hóa mật khẩu người dùng
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Tạo đối tượng người dùng mới
	newUser := &model.User{
		ID:        uuid.New().String(), // UUID mới
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword), // Lưu mật khẩu đã mã hóa
		Role:      "user",                 // Vai trò mặc định
		CreatedAt: time.Now(),             // Thời gian tạo
	}

	if err := s.repo.Create(newUser); err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	// Cache user vào Redis
	s.redis.Set(ctx, "user:"+newUser.ID, newUser.Email, time.Minute*10)

	return newUser, nil
}

// Login - Xác thực người dùng và sinh JWT
func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if user.Password != password {
		return "", errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	return token, nil
}

// GetUserProfile - Lấy thông tin người dùng theo ID
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*model.User, error) {
	// Check Redis trước
	if email, err := s.redis.Get(ctx, "user:"+userID).Result(); err == nil {
		// Tuỳ mục tiêu ông chủ có thể trả về tạm từ Redis hoặc tiếp tục gọi DB
		fmt.Printf("Cache hit for user %s: %s\n", userID, email)
	}

	// id, err := strconv.ParseInt(userID, 10, 64)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid user ID format")
	// }

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return user, nil
}

// generateJWT - Tạo access token JWT
func (s *UserService) generateJWT(userID string) (string, error) {
	secretKey := []byte(s.cfg.JWTSecret)
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}
	return signedToken, nil
}
