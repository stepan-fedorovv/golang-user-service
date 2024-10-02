package models

import "time"

type User struct {
	ID        int       `json:"ID"`
	Username  *string   `json:"username"`
	Password  *string   `json:"password"`
	Email     *string   `json:"email"`
	Name      *string   `json:"name"`
	Surname   *string   `json:"surname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
