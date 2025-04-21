# golang_base_authenticated

## 🧑‍💻 `UserController` - Quản lý người dùng

Controller này quản lý tất cả các API liên quan đến người dùng trong hệ thống.

---

### 📁 File: `controller/user_controller.go`

---

## ✅ **Endpoints**

### `POST /api/v1/register` – Đăng ký người dùng mới

**Request body:**

```json
{
  "name": "John Doe",
  "email": "johndoe@example.com",
  "password": "securepassword"
}
```

**Response:**

```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "johndoe@example.com",
    "created_at": "2025-04-21T12:00:00Z"
  }
}
```

---

### `POST /api/v1/login` – Đăng nhập và nhận JWT

**Request body:**

```json
{
  "email": "johndoe@example.com",
  "password": "securepassword"
}
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### `GET /api/v1/profile` – Lấy thông tin người dùng (yêu cầu JWT)

**Headers:**

```
Authorization: Bearer <JWT Token>
```

**Response:**

```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "johndoe@example.com",
  "created_at": "2025-04-21T12:00:00Z"
}
```

---

## 🔒 Middleware

- Route `/profile` sử dụng middleware JWT để kiểm tra xác thực.
- `userID` được lấy từ claim `sub` trong JWT.

---

## ⚙️ Phụ thuộc

- `service.UserService`
- `github.com/golang-jwt/jwt/v5`
- `github.com/gofiber/fiber/v2`

---

## 🚧 TODO

- [ ] Thêm xác thực email.
- [ ] Mã hóa mật khẩu (hiện đang để plaintext).
- [ ] Thêm Swagger/OpenAPI để auto-gen docs.
- [ ] Thêm unit tests.

