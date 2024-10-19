package models

import "time"

type User struct {
	Id        string    `json:"id,omitempty" db:"id"`
	FirstName string    `json:"firstName,omitempty" db:"first_name"`
	LastName  string    `json:"lastName,omitempty" db:"last_name"`
	Email     string    `json:"email,omitempty" db:"email"`
	Username  string    `json:"username,omitempty" db:"username"`
	Password  string    `json:"_" db:"password"`
	CreatedAt time.Time `json:"_" db:"created_at"`
	UpdatedAt time.Time `json:"_" db:"updated_at"`
}
