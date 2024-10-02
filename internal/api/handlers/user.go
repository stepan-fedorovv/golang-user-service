package handlers

import (
	"app/internal/api/logic/facades"
	"app/internal/api/schemas"
	"app/internal/db"
	"app/internal/db/models"
	"app/utils/http-utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-ldap/ldap"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Status     string                    `json:"status"`
	StatusCode string                    `json:"status_code"`
	Data       schemas.JWTResponseSchema `json:"data"`
}

func LoginHandler(log *slog.Logger, s *db.Storage, jwtSecretKey []byte, conn *ldap.Conn, baseDN string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		validate := validator.New()
		if err := http_utils.DecodeBody(r, &user); err != nil {
			http.Error(w, "Ошибка в валидации JSON: "+err.Error(), http.StatusTeapot)
			return
		}
		if err := validate.Struct(user); err != nil {
			http.Error(w, "Ошибка в теле запроса: "+err.Error(), http.StatusTeapot)
			return
		}
		user, err := facades.LDAPAuth(conn, user, s, baseDN)
		if err != nil {
			http.Error(w, "Ошибка в авторизации пользователя: "+err.Error(), http.StatusTeapot)
		}
		token, err := facades.UserLogin(user, s, jwtSecretKey)
		if err != nil {
			http.Error(w, "Ошибка в регистрации пользователя: "+err.Error(), http.StatusTeapot)
			return
		}
		if err := http_utils.EncodeBody(w, token); err != nil {
			http.Error(w, "Ошибка в формировании ответа: "+err.Error(), http.StatusTeapot)
			return
		}
	}
}

func RefreshHandler(log *slog.Logger, s *db.Storage, jwtSecretKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user schemas.RefreshSchema
		validate := validator.New()

		if err := http_utils.DecodeBody(r, &user); err != nil {
			http.Error(w, "Ошибка в валидации JSON: "+err.Error(), http.StatusTeapot)
			return
		}
		if err := validate.Struct(user); err != nil {
			http.Error(w, "Ошибка в теле запроса: "+err.Error(), http.StatusTeapot)
			return
		}
		tokens, err := facades.Refresh(user, s, jwtSecretKey)
		if err != nil {
			http.Error(w, "Ошибка в обновлении токена: "+err.Error(), http.StatusTeapot)
			return
		}
		if err := http_utils.EncodeBody(w, tokens); err != nil {
			http.Error(w, "Ошибка в формировании ответа: "+err.Error(), http.StatusTeapot)
			return
		}

	}
}

func MeHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(models.User)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userResponse := facades.ConvertToResponse(user)
		if err := http_utils.EncodeBody(w, userResponse); err != nil {
			http.Error(w, "Ошибка в формировании ответа: "+err.Error(), http.StatusTeapot)
		}
	}
}

func PartialUpdateHandler(log *slog.Logger, s *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user schemas.UserResponseSchema
		if err := http_utils.DecodeBody(r, &user); err != nil {
			http.Error(w, "Ошибка в валидации JSON: "+err.Error(), http.StatusTeapot)
			return
		}

		urlParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(urlParam)
		if err != nil {
			http.Error(w, "Некорректный id"+err.Error(), http.StatusTeapot)
			return
		}
		user, err = facades.UserPartialUpdate(user, s, id)
		if err != nil {
			http.Error(w, "Ошибка в обновление пользователя: "+err.Error(), http.StatusTeapot)
			return
		}
		if err := http_utils.EncodeBody(w, user); err != nil {
			http.Error(w, "Ошибка в формировании ответа: "+err.Error(), http.StatusTeapot)
		}
	}
}
