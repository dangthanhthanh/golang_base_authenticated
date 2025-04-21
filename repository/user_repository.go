package repository

import (
	"base-app/model"
	"database/sql"
	"errors"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id int64) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
	return err
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *userRepository) FindByID(id int64) (*model.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}
