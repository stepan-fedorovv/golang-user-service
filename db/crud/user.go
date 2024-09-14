package crud

import (
	"app/api/schemas"
	"app/db"
	"app/db/models"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(body schemas.UserRegisterSchema, s *db.Storage) (schemas.UserResponseSchema, error) {
	var user schemas.UserResponseSchema
	stmt := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username`
	bytes, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return schemas.UserResponseSchema{}, err
	}
	err = s.DB.QueryRow(context.Background(), stmt, body.Username, string(bytes)).Scan(
		&user.ID,
		&user.Username,
	)
	if err != nil {
		return schemas.UserResponseSchema{}, err
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

func SelectUserById(id int, s *db.Storage) (models.User, error) {
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
