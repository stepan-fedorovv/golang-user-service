package settings

import (
	"app/db"
	"app/db/crud"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func JwtAuthMiddleware(jwtSecretKey []byte, s *db.Storage) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			splitToken := strings.Split(bearerToken, " ")
			if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			token := splitToken[1]
			claims := jwt.MapClaims{}
			parsedToken, err := jwt.ParseWithClaims(
				token,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					return jwtSecretKey, nil
				})
			if err != nil || !parsedToken.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			user, err := crud.SelectUserById(claims["sub"].(float64), s)
			if &user.Username == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
