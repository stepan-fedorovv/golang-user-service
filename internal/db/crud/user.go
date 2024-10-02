package crud

import (
	"app/internal/api/schemas"
	"app/internal/db"
	"app/internal/db/models"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(body schemas.CreateUserSchema, s *db.Storage) (models.User, error) {
	var user models.User
	stmt := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, password`
	bytes, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return models.User{}, err
	}
	err = s.DB.QueryRow(context.Background(), stmt, body.Username, string(bytes)).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func SelectUserByUsername(username string, s *db.Storage) (models.User, error) {
	var user models.User
	stmt := `
			SELECT 
			    "id",
			    "username",
			    "password",
			    "email",
			    "name",
			    "surname",
			    "created_at",
			    "updated_at"
			FROM users WHERE username = $1`
	err := s.DB.QueryRow(context.Background(), stmt, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.As(err, &sql.ErrNoRows) {
		return models.User{}, nil
	} else if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func SelectUserById(id float64, s *db.Storage) (models.User, error) {
	var user models.User
	stmt := `
			SELECT 
			    "id",
			    "username",
			    "password",
			    "email",
			    "name",
			    "surname",
			    "created_at",
			    "updated_at"
			FROM users WHERE id = $1`
	err := s.DB.QueryRow(context.Background(), stmt, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.As(err, &sql.ErrNoRows) {
		return models.User{}, nil
	} else if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func PartialUpdateUser(user schemas.UserResponseSchema, s *db.Storage, id int) (schemas.UserResponseSchema, error) {
	stmt := `
			UPDATE users SET username = COALESCE($1, username),
			                 email = COALESCE($2, email),
			                 name = COALESCE($3, name),
			                 surname = COALESCE($4, surname)
			WHERE id = $5 RETURNING id, username, email, name, surname, created_at, updated_at`
	err := s.DB.QueryRow(
		context.Background(),
		stmt,
		user.Username,
		user.Email,
		user.Name,
		user.Surname,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Surname, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return schemas.UserResponseSchema{}, err
	}
	return user, nil
}
