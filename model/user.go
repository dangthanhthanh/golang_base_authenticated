package model

import (
	"time"
)

// User mô tả cấu trúc dữ liệu người dùng
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"` // Mật khẩu sẽ được mã hóa trước khi lưu vào DB
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"` // Thời gian tạo tài khoản
}
