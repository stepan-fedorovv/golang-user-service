package schemas

import "time"

type UserResponseSchema struct {
	ID        int       `json:"id"`
	Username  *string   `json:"username"`
	Email     *string   `json:"email"`
	Name      *string   `json:"name"`
	Surname   *string   `json:"surname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserSchema struct {
	Username    string `json:"username" validate:"required,min=3,max=255"`
	Password    string `json:"password" validate:"required,min=8,max=255"`
	FirstName   string `json:"first_name" validate:"required,min=3,max=255"`
	Description string `json:"description"`
}

type JWTResponseSchema struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type UserLoginSchema struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshSchema struct {
	Refresh string `json:"refresh"`
}
