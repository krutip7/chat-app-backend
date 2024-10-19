package models

import "time"

type User struct {
	Id        string    `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"_" db:"password"`
	CreatedAt time.Time `json:"_" db:"created_at"`
	UpdatedAt time.Time `json:"_" db:"updated_at"`
}
