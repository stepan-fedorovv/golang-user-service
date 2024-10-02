package facades

import (
	"app/internal/api/schemas"
	"app/internal/db"
	"app/internal/db/crud"
	"app/internal/db/models"
	jwtUtils "app/utils/jwt-utils"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func UserLogin(
	body models.User, s *db.Storage, jwtSecretKey []byte,
) (schemas.JWTResponseSchema, error) {
	user, err := crud.SelectUserByUsername(*body.Username, s)
	if user.Username == nil {
		return schemas.JWTResponseSchema{}, errors.New("user Not found")
	}
	if *user.Password != *body.Password {
		return schemas.JWTResponseSchema{}, errors.New("password incorrect")
	}
	token, err := jwtUtils.CreateToken(user, jwtSecretKey)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	return token, nil
}

func Refresh(body schemas.RefreshSchema, s *db.Storage, jwtSecretKey []byte) (schemas.JWTResponseSchema, error) {
	var user models.User
	refreshToken := body.Refresh
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect signing method")
		}
		return jwtSecretKey, nil
	})
	if err != nil || !token.Valid {
		return schemas.JWTResponseSchema{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return schemas.JWTResponseSchema{}, errors.New("invalid token claims")
	}
	if claims["type"] != "refresh" {
		return schemas.JWTResponseSchema{}, errors.New("invalid token type")
	}
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return schemas.JWTResponseSchema{}, errors.New("invalid expiration time")
	}
	expTime := time.Unix(int64(expFloat), 0)
	if expTime.Before(time.Now()) {
		return schemas.JWTResponseSchema{}, errors.New("refresh token expired")
	}
	user, err = crud.SelectUserById(claims["sub"].(float64), s)
	if &user.Username == nil {
		return schemas.JWTResponseSchema{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}
	tokens, err := jwtUtils.CreateToken(user, jwtSecretKey)
	if err != nil {
		return schemas.JWTResponseSchema{}, err
	}
	return tokens, nil
}

func ConvertToResponse(user models.User) schemas.UserResponseSchema {
	return schemas.UserResponseSchema{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserPartialUpdate(user schemas.UserResponseSchema, s *db.Storage, id int) (schemas.UserResponseSchema, error) {
	updatedUser, err := crud.PartialUpdateUser(user, s, id)
	if err != nil {
		return schemas.UserResponseSchema{}, err
	}
	return updatedUser, nil
}

func LDAPAuth(conn *ldap.Conn, user models.User, s *db.Storage, baseDN string) (models.User, error) {
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(sAMAccountName=%s)", *user.Username),
		[]string{"dn", "cn", "sAMAccountName", "givenName", "description"},
		nil,
	)
	searchResp, err := conn.Search(searchRequest)
	if err != nil {
		return models.User{}, err
	}
	entries := searchResp.Entries
	if len(entries) != 1 {
		return models.User{}, errors.New("найдено больше одной записи")
	}
	entry := entries[0]
	existUser, _ := crud.SelectUserByUsername(*user.Username, s)
	if existUser.Username == nil {
		createdUser, err := crud.CreateUser(schemas.CreateUserSchema{
			Username: entry.GetAttributeValue("sAMAccountName"),
			Password: *user.Password,
		}, s)
		if err != nil {
			return models.User{}, err
		}
		return createdUser, nil
	}
	return existUser, nil
}
