package jwt_utils

import (
	"app/api/schemas"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(payload jwt.MapClaims, jwtSecretKey []byte) (schemas.JWTResponseSchema, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	return schemas.JWTResponseSchema{
		Access: t,
	}, nil
}
