// File: pkg/db/migrate.go
package db

import (
	"base-app/model"
)

func Migrate() {
	DB.AutoMigrate(&model.User{}) // có thể thêm nhiều model khác ở đây
}
