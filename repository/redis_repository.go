package repository

import (
	"base-app/model"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisRepository định nghĩa các hàm thao tác với Redis
type RedisRepository interface {
	// Auth
	SetAccessToken(ctx context.Context, token, userID, role string, ttl time.Duration) error
	IsTokenValid(ctx context.Context, token string) bool
	RevokeToken(ctx context.Context, token string) error

	// Refresh Token
	SetRefreshToken(ctx context.Context, refreshToken, userID string, ttl time.Duration) error
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error)

	// Rate Limiting
	IncrementRate(ctx context.Context, key string, ttl time.Duration) (int64, error)

	// Role & Permission
	AddPermissionToRole(ctx context.Context, role, permission string) error
	RoleHasPermission(ctx context.Context, role, permission string) (bool, error)

	// Session Management
	AddTokenToUser(ctx context.Context, userID, token string) error
	RemoveTokenFromUser(ctx context.Context, userID, token string) error
	GetAllUserTokens(ctx context.Context, userID string) ([]string, error)
	RevokeAllUserTokens(ctx context.Context, userID string) error

	// User Profile
	SetUserProfile(ctx context.Context, userID, email string, ttl time.Duration) error
	SetUserProfileFull(ctx context.Context, user *model.User, ttl time.Duration) error

	GetUserProfile(ctx context.Context, userID string) (*model.User, error)

	// list user mail
	AddUserEmailToList(ctx context.Context, email string) error
	GetUserEmails(ctx context.Context, start, stop int64) ([]string, error)

	// Xóa dữ liệu người dùng khỏi Redis
	DeletedAllDataUserAccount(ctx context.Context, key string) error
}

// redisRepo là implement chính
type redisRepo struct {
	client *redis.Client
}

// NewRedisRepository khởi tạo repository
func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepo{client: client}
}

// ======================= AUTH =======================

func (r *redisRepo) SetAccessToken(ctx context.Context, token, userID, role string, ttl time.Duration) error {
	key := "auth:token:" + token
	value := fmt.Sprintf("%s|%s", userID, role)
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisRepo) IsTokenValid(ctx context.Context, token string) bool {
	key := "auth:token:" + token
	exists, _ := r.client.Exists(ctx, key).Result()
	return exists == 1
}

func (r *redisRepo) RevokeToken(ctx context.Context, token string) error {
	key := "auth:token:" + token
	return r.client.Del(ctx, key).Err()
}

// ======================= REFRESH TOKEN =======================

func (r *redisRepo) SetRefreshToken(ctx context.Context, refreshToken, userID string, ttl time.Duration) error {
	key := "auth:refresh:" + refreshToken
	return r.client.Set(ctx, key, userID, ttl).Err()
}

func (r *redisRepo) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := "auth:refresh:" + refreshToken
	return r.client.Get(ctx, key).Result()
}

// ======================= RATE LIMITING =======================

func (r *redisRepo) IncrementRate(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		r.client.Expire(ctx, key, ttl)
	}
	return count, nil
}

// ======================= ROLE - PERMISSION =======================

func (r *redisRepo) AddPermissionToRole(ctx context.Context, role, permission string) error {
	key := "user:role:" + role
	return r.client.SAdd(ctx, key, permission).Err()
}

func (r *redisRepo) RoleHasPermission(ctx context.Context, role, permission string) (bool, error) {
	key := "user:role:" + role
	return r.client.SIsMember(ctx, key, permission).Result()
}

// ======================= SESSION =======================

func (r *redisRepo) AddTokenToUser(ctx context.Context, userID, token string) error {
	key := "auth:user:" + userID + ":tokens"
	return r.client.SAdd(ctx, key, token).Err()
}

func (r *redisRepo) RemoveTokenFromUser(ctx context.Context, userID, token string) error {
	key := "auth:user:" + userID + ":tokens"
	return r.client.SRem(ctx, key, token).Err()
}

func (r *redisRepo) GetAllUserTokens(ctx context.Context, userID string) ([]string, error) {
	key := "auth:user:" + userID + ":tokens"
	return r.client.SMembers(ctx, key).Result()
}

func (r *redisRepo) RevokeAllUserTokens(ctx context.Context, userID string) error {
	tokens, err := r.GetAllUserTokens(ctx, userID)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		r.RevokeToken(ctx, token)
	}
	key := "auth:user:" + userID + ":tokens"
	return r.client.Del(ctx, key).Err()
}

// ======================= USER PROFILE =======================

// Lưu email đơn giản
func (r *redisRepo) SetUserProfile(ctx context.Context, userID, email string, ttl time.Duration) error {
	key := "user:profile:" + userID
	return r.client.Set(ctx, key, email, ttl).Err()
}

// Lưu đầy đủ profile dưới dạng Hash
func (r *redisRepo) SetUserProfileFull(ctx context.Context, user *model.User, ttl time.Duration) error {
	key := "user:profile:" + user.ID
	data := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	}
	if err := r.client.HMSet(ctx, key, data).Err(); err != nil {
		return err
	}
	return r.client.Expire(ctx, key, ttl).Err()
}

// GetUserProfile - Lấy thông tin người dùng từ Redis
func (r *redisRepo) GetUserProfile(ctx context.Context, userID string) (*model.User, error) {
	key := "user:profile:" + userID
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err // Redis lỗi hoặc không có key
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("user not found in cache") // Không tìm thấy trong cache
	}

	// Chuyển đổi dữ liệu Redis (map) thành đối tượng model.User
	user := &model.User{
		ID:    userID,
		Name:  data["name"],
		Email: data["email"],
		Role:  data["role"],
	}

	return user, nil
}

// AddUserEmailToList thêm email vào danh sách người dùng
func (r *redisRepo) AddUserEmailToList(ctx context.Context, email string) error {
	key := "user:list_email"
	return r.client.LPush(ctx, key, email).Err()
}

// GetUserEmails lấy danh sách email người dùng từ Redis
func (r *redisRepo) GetUserEmails(ctx context.Context, start, stop int64) ([]string, error) {
	key := "user:list_email"
	return r.client.LRange(ctx, key, start, stop).Result()
}

// DeletedUserAccount - Xóa tất cả dữ liệu liên quan đến người dùng khỏi Redis
func (r *redisRepo) DeletedAllDataUserAccount(ctx context.Context, userID string) error {
	// Xóa thông tin profile người dùng
	err := r.client.Del(ctx, "user:profile:"+userID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user profile from Redis: %v", err)
	}

	// Xóa email người dùng (nếu có lưu cache email)
	err = r.client.Del(ctx, "user:email:"+userID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user email from Redis: %v", err)
	}

	// Xóa danh sách token của người dùng
	err = r.client.Del(ctx, "auth:user:"+userID+":tokens").Err()
	if err != nil {
		return fmt.Errorf("failed to delete user tokens from Redis: %v", err)
	}

	// Xóa các khóa liên quan đến session, nếu có
	err = r.client.Del(ctx, "auth:user:"+userID+":sessions").Err()
	if err != nil {
		return fmt.Errorf("failed to delete user sessions from Redis: %v", err)
	}

	// Xóa các khóa liên quan đến refresh token (nếu có)
	err = r.client.Del(ctx, "auth:refresh:"+userID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user refresh token from Redis: %v", err)
	}

	// Xóa tất cả các quyền và role nếu có
	err = r.client.Del(ctx, "user:role:"+userID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user role from Redis: %v", err)
	}

	// Nếu có các dữ liệu khác cần xóa, thêm vào đây

	return nil
}
