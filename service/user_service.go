package service

import (
	"base-app/config"
	"base-app/model"
	"base-app/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  repository.UserRepository
	redis repository.RedisRepository
	cfg   config.Config
}

func NewUserService(repo repository.UserRepository, redisRepo repository.RedisRepository, cfg config.Config) *UserService {
	return &UserService{
		repo:  repo,
		redis: redisRepo,
		cfg:   cfg,
	}
}

// Register - Đăng ký người dùng mới
func (s *UserService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	// Check email đã tồn tại trong DB
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash mật khẩu
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %v", err)
	}

	// Tạo user mới trong DB
	newUser, err := s.repo.Create(name, email, string(hashPassword))
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	// Lưu email vào Redis (danh sách email người dùng)
	err = s.redis.AddUserEmailToList(ctx, email)
	if err != nil {
		// Cảnh báo nếu Redis gặp lỗi nhưng không làm gián đoạn quá trình
		fmt.Printf("warning: failed to store email in Redis list: %v\n", err)
	}

	// Lưu thông tin user vào Redis (cache user profile)
	err = s.redis.SetUserProfile(ctx, newUser.ID, newUser.Email, 10*time.Minute) // 10 phút cache
	if err != nil {
		// Cảnh báo nếu Redis gặp lỗi nhưng không làm gián đoạn quá trình
		fmt.Printf("warning: failed to cache user profile in Redis: %v\n", err)
	}

	return newUser, nil
}

// Login - Xác thực người dùng và sinh JWT
func (s *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	// Kiểm tra user trong PostgreSQL
	user, err := s.repo.FindByEmail(email) // PostgreSQL
	if err != nil {
		return "", errors.New("user not found")
	}

	// Kiểm tra mật khẩu
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Sinh JWT token
	token, err := s.generateJWT(user.ID, user.Role, time.Hour*1)
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	// Lưu token vào Redis (để xác thực nhanh chóng)
	if err := s.redis.SetAccessToken(ctx, token, user.ID, user.Role, time.Hour*1); err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %v", err)
	}

	// Lưu token vào danh sách của user (phục vụ logout all)
	if err := s.redis.AddTokenToUser(ctx, user.ID, token); err != nil {
		// Tùy chiến lược: có thể log lại hoặc fail luôn
		fmt.Printf("warning: failed to add token to user's list in Redis: %v\n", err)
	}

	// Optional: Lưu email vào Redis danh sách người dùng (nếu chưa có)
	if err := s.redis.AddUserEmailToList(ctx, email); err != nil {
		// Không làm gián đoạn quá trình đăng nhập nếu Redis có lỗi
		fmt.Printf("warning: failed to store email in Redis list: %v\n", err)
	}

	// Lưu thông tin người dùng vào Redis (cache profile)
	if err := s.redis.SetUserProfileFull(ctx, user, time.Hour*24); err != nil {
		// Log cảnh báo nhưng không làm gián đoạn quá trình đăng nhập
		fmt.Printf("warning: failed to cache user profile in Redis: %v\n", err)
	}

	return token, nil
}

// GetUserProfile - Lấy thông tin người dùng theo ID
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*model.User, error) {
	// Kiểm tra Redis trước để tìm thông tin người dùng
	cachedUser, err := s.redis.GetUserProfile(ctx, userID)
	if err == nil && cachedUser != nil {
		// Nếu tìm thấy trong Redis, trả về dữ liệu cache
		fmt.Printf("Cache hit for user %s\n", userID)
		return cachedUser, nil
	}

	// Nếu không có trong Redis, lấy từ cơ sở dữ liệu
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Lưu thông tin người dùng vào Redis để cache cho lần sau
	err = s.redis.SetUserProfileFull(ctx, user, time.Hour*24) // Lưu cache trong 24 giờ
	if err != nil {
		// Không làm gián đoạn quá trình nếu Redis gặp lỗi
		fmt.Printf("warning: failed to cache user profile in Redis: %v\n", err)
	}

	return user, nil
}

// generateJWT - Tạo access token JWT
func (s *UserService) generateJWT(userID string, role string, ttl time.Duration) (string, error) {
	// Key bí mật từ cấu hình
	secretKey := []byte(s.cfg.JWTSecret)

	// Khai báo claims
	claims := jwt.MapClaims{
		"sub":  userID,                     // user ID (subject)
		"role": role,                       // role của người dùng
		"iat":  time.Now().Unix(),          // Thời gian tạo token
		"exp":  time.Now().Add(ttl).Unix(), // Thời gian hết hạn (ttl được truyền vào)
	}

	// Tạo JWT token với phương thức HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token và trả về
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}
	return signedToken, nil
}

// UpdateUserProfile - Cập nhật thông tin người dùng
func (s *UserService) UpdateUserProfile(ctx context.Context, userID, name, email string) (*model.User, error) {
	// Kiểm tra email mới có bị trùng với người dùng khác không
	existingUser, err := s.repo.FindByEmail(email)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		// Trường hợp email đã tồn tại và không phải của người dùng hiện tại
		return nil, errors.New("email already registered")
	}

	// Cập nhật thông tin người dùng trong DB
	user, err := s.repo.Update(userID, name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Cập nhật lại thông tin trong Redis để đồng bộ
	err = s.redis.SetUserProfileFull(ctx, user, time.Hour*24) // Cập nhật cache trong 24h
	if err != nil {
		// Không làm gián đoạn quá trình nếu Redis gặp lỗi
		fmt.Printf("warning: failed to update user profile in Redis: %v\n", err)
	}

	return user, nil
}

// ChangePassword - Thay đổi mật khẩu người dùng
func (s *UserService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Lấy thông tin người dùng từ DB
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Kiểm tra mật khẩu cũ
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("incorrect old password")
	}

	// Mã hóa mật khẩu mới
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Cập nhật mật khẩu mới vào DB
	err = s.repo.UpdatePassword(userID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	// Cập nhật lại thông tin mật khẩu trong Redis nếu có lưu
	// (Nếu có lưu thông tin người dùng trong Redis, cần cập nhật lại hoặc xóa cache nếu cần)
	err = s.redis.RemoveTokenFromUser(ctx, userID, "") // Xóa tất cả các token của user nếu cần
	if err != nil {
		// Log cảnh báo nếu có lỗi khi xóa token trong Redis
		fmt.Printf("warning: failed to remove user tokens from Redis: %v\n", err)
	}

	// // Nếu có lưu thông tin người dùng trong Redis (cache), có thể cập nhật lại hoặc xóa cache
	// err = s.redis.DeletedUserAccount(ctx, "user:profile:"+userID) // Xóa cache của user
	// khon lu password trong redis nen khong can xoas nos
	if err != nil {
		// Log cảnh báo nếu có lỗi khi xóa cache
		fmt.Printf("warning: failed to remove user profile from Redis: %v\n", err)
	}

	return nil
}

// DeleteUserAccount - Xóa tài khoản người dùng
func (s *UserService) ForceDeletedUserAccount(ctx context.Context, userID string) error {
	err := s.repo.Delete(userID)
	if err != nil {
		return fmt.Errorf("failed to delete user account from database: %v", err)
	}

	err = s.redis.DeletedAllDataUserAccount(ctx, userID)
	if err != nil {
		// Log cảnh báo nếu có lỗi trong việc xóa cache, nhưng không ngừng thực hiện
		fmt.Printf("warning: failed to delete user data from Redis: %v\n", err)
	}
	// ...
	return nil
}
