package settings

import (
	"app/api/handlers"
	"app/db"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func Router(storage *db.Storage, log *slog.Logger, cfg *Config) chi.Router {
	var jwtSecretKey = []byte(cfg.SecretKey)
	router := chi.NewRouter()
	router.Mount("/api", API(storage, log, jwtSecretKey))
	return router
}

func API(storage *db.Storage, log *slog.Logger, jwtSecretKey []byte) chi.Router {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Post("/register/", handlers.RegisterHandler(log, storage, jwtSecretKey))
		r.Post("/login/", handlers.LoginHandler(log, storage, jwtSecretKey))
	})
	return r
}
