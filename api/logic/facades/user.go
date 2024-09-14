package facades

import (
	"app/api/schemas"
	"app/db"
	"app/db/crud"
	"app/db/models"
	jwtUtils "app/utils/jwt-utils"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func UserRegister(
	body schemas.UserRegisterSchema, s *db.Storage, jwtSecretKey []byte,
) (schemas.JWTResponseSchema, error) {
	if _, err := crud.SelectUserByUsername(body.Username, s); err != nil {
		return schemas.JWTResponseSchema{}, errors.New("user Already Exists")
	}
	user, err := crud.CreateUser(body, s)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	payload := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token, err := jwtUtils.CreateToken(payload, jwtSecretKey)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	return token, nil
}

func UserLogin(
	body models.User, s *db.Storage, jwtSecretKey []byte,
) (schemas.JWTResponseSchema, error) {
	user, err := crud.SelectUserByUsername(*body.Username, s)
	if user.Username == nil {
		return schemas.JWTResponseSchema{}, errors.New("user Not found")
	}
	payload := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token, err := jwtUtils.CreateToken(payload, jwtSecretKey)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*body.Password)); err != nil {
		return schemas.JWTResponseSchema{}, errors.New("password incorrect")
	}
	return token, nil
}
