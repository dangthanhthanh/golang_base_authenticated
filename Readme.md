# 📦 Router Package

Gói `router` chịu trách nhiệm định nghĩa và cấu hình các tuyến (routes) API cho ứng dụng. Đây là nơi thiết lập luồng giao tiếp giữa client và các controller tương ứng, đồng thời tích hợp các middleware như xác thực JWT.

---

## 🧭 Cấu trúc Route

Tất cả các route được đặt dưới prefix `/api/v1`.

### 🔐 Auth Routes (`/api/v1/auth`)

| Method | Endpoint     | Mô tả                  |
|--------|--------------|------------------------|
| POST   | /register    | Đăng ký người dùng mới |
| POST   | /login       | Đăng nhập và lấy token |

### 👤 User Routes (`/api/v1/user`)

> ⚠️ **Yêu cầu JWT token hợp lệ**

| Method | Endpoint     | Mô tả                    |
|--------|--------------|--------------------------|
| GET    | /profile     | Lấy thông tin người dùng |

---

## 🛡️ Middleware: JWT

Các route yêu cầu xác thực sẽ sử dụng middleware JWT từ `github.com/gofiber/jwt/v3`. Middleware này sẽ:

- Kiểm tra token từ header `Authorization`
- Xác thực chữ ký token với secret từ config
- Trả lỗi `401 Unauthorized` nếu token không hợp lệ hoặc không tồn tại

### 🔧 Cấu hình JWT

```go
jwt.New(jwt.Config{
    SigningKey:   []byte(config.LoadConfig().JWTSecret),
    ErrorHandler: jwtErrorHandler,
})