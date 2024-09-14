package handlers

import (
	"app/api/logic/facades"
	"app/api/schemas"
	"app/db"
	"app/db/models"
	"app/utils/http-utils"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Response struct {
	Status     string                    `json:"status"`
	StatusCode string                    `json:"status_code"`
	Data       schemas.JWTResponseSchema `json:"data"`
}

func RegisterHandler(log *slog.Logger, s *db.Storage, jwtSecretKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user schemas.UserRegisterSchema
		validate := validator.New()

		if err := http_utils.DecodeBody(r, &user); err != nil {
			http.Error(w, "Ошибка в валидации JSON: "+err.Error(), http.StatusTeapot)
			return
		}

		if err := validate.Struct(user); err != nil {
			http.Error(w, "Ошибка в теле запроса: "+err.Error(), http.StatusTeapot)
			return
		}

		token, err := facades.UserRegister(user, s, jwtSecretKey)

		if err != nil {
			http.Error(w, "Ошибка создания пользователя: "+err.Error(), http.StatusTeapot)
			return
		}

		if err := http_utils.EncodeBody(w, token); err != nil {
			http.Error(w, "Ошибка в формировании ответа: "+err.Error(), http.StatusTeapot)
		}
	}
}

func LoginHandler(log *slog.Logger, s *db.Storage, jwtSecretKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		validate := validator.New()

		if err := http_utils.DecodeBody(r, &user); err != nil {
			http.Error(w, "Ошибка в валидации JSON: "+err.Error(), http.StatusTeapot)
			return
		}
		if err := validate.Struct(user); err != nil {
			http.Error(w, "Ошибка в теле запроса: "+err.Error(), http.StatusTeapot)
		}
		token, err := facades.UserLogin(user, s, jwtSecretKey)
		if err != nil {
			http.Error(w, "Ошибка в регистрации пользователя: "+err.Error(), http.StatusTeapot)
			return
		}
		if err := http_utils.EncodeBody(w, token); err != nil {
			http.Error(w, "Ошибка в формировании ответа: "+err.Error(), http.StatusTeapot)
		}
	}
}

func MeHandler(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
