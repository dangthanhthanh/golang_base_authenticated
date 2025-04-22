package repository

import (
	"base-app/model"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository là interface cho các thao tác với người dùng
type UserRepository interface {
	Create(name string, email string, hashPassword string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Update(userID string, name string, email string) (*model.User, error)
	UpdatePassword(userID, password string) error
	Delete(userID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(name string, email string, hashPassword string) (*model.User, error) {
	newUser := &model.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  hashPassword,
		Role:      model.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.Create(newUser).Error; err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}

	return newUser, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}

func (r *userRepository) Update(userID string, name string, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	user.Name = name
	user.Email = email
	user.UpdatedAt = time.Now()

	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdatePassword(userID string, newPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", newPassword).Error
}

func (r *userRepository) Delete(userID string) error {
	return r.db.Where("id = ?", userID).Delete(&model.User{}).Error
}
