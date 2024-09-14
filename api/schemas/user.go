package schemas

import "time"

type UserResponseSchema struct {
	ID        int       `json:"ID"`
	Username  *string   `json:"username"`
	Email     *string   `json:"email"`
	Name      *string   `json:"name"`
	Surname   *string   `json:"surname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegisterSchema struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type JWTResponseSchema struct {
	Access string `json:"access"`
}

type UserLoginSchema struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
