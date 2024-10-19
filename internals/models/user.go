package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string    `json:"id,omitempty" db:"id"`
	FirstName string    `json:"firstName,omitempty" db:"first_name"`
	LastName  string    `json:"lastName,omitempty" db:"last_name"`
	Email     string    `json:"email,omitempty" db:"email"`
	Username  string    `json:"username,omitempty" db:"username"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

func (user *User) VerifyPassword(password string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return err == nil, err
}
