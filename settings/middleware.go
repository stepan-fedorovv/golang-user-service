package settings

import (
	"app/db/crud"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type UserClaims struct {
	jwt.RegisteredClaims
	userId int `json:"user_id"`
}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userClaims UserClaims
		bearerToken := r.Header.Get("Authorization")
		token := strings.Split(bearerToken, " ")[1]
		if &token == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		_, err := jwt.ParseWithClaims(
			token,
			&userClaims,
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecretKey, nil
			})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		user, err := crud.SelectUserById(userClaims.userId, s)
		if &user.Username == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
