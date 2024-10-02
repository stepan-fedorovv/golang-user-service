package jwt_utils

import (
	"app/internal/api/schemas"
	"app/internal/db/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(user models.User, jwtSecretKey []byte) (schemas.JWTResponseSchema, error) {
	accessPayload := jwt.MapClaims{
		"sub":  user.ID,
		"type": "access",
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	at, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	refreshPayload := jwt.MapClaims{
		"sub":  user.ID,
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshPayload)
	rt, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	return schemas.JWTResponseSchema{
		Access:  at,
		Refresh: rt,
	}, nil
}
