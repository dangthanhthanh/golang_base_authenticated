// File: utils/jwt.go
package utils

import (
	"base-app/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user model.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenStr
}

func ParseJWT(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
