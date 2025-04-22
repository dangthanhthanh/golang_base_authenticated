package model

import (
	"time"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User mô tả cấu trúc dữ liệu người dùng
type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`             // Sử dụng kiểu uint cho ID, PostgreSQL sẽ tự động sinh giá trị (SERIAL)
	Name      string    `gorm:"not null" json:"name"`             // Đảm bảo Name không được null
	Email     string    `gorm:"not null;unique" json:"email"`     // Đảm bảo Email không được null và là duy nhất
	Password  string    `gorm:"not null" json:"password"`         // Mật khẩu cần thiết và phải được mã hóa trước khi lưu
	Role      string    `gorm:"not null" json:"role"`             // Vai trò của người dùng, có thể có giá trị như "admin", "user", v.v.
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Tự động gán thời gian tạo
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"Updated_at"` // Tự động gán thời gian autoUpdateTime
}
