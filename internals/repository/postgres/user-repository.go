package postgres

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/krutip7/chat-app-server/internals/models"
)

type UserRepository struct {
	DB *sqlx.DB
}

func (userRepo *UserRepository) GetUserById(userId int) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	var users []models.User

	query := `SELECT id, first_name, last_name, email FROM users WHERE id=$1`

	err := userRepo.DB.SelectContext(ctx, &users, query, userId)
	if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (userRepo *UserRepository) GetUserByEmail(email string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	var users []models.User

	query := `SELECT id, first_name, last_name, email, password FROM users WHERE email=$1`

	err := userRepo.DB.SelectContext(ctx, &users, query, email)
	if err != nil {
		return nil, err
	} else if len(users) < 1 {
		return nil, errors.New("user not found")
	}

	return &users[0], nil
}

func (userRepo *UserRepository) GetAllUsers() ([]models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	var users []models.User

	query := `SELECT id, first_name, last_name, email, username FROM users`

	err := userRepo.DB.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}
