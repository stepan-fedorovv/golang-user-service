package settings

import (
	"app/api/handlers"
	"app/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-ldap/ldap"
	"log/slog"
)

func Router(storage *db.Storage, log *slog.Logger, cfg *Config, conn *ldap.Conn) chi.Router {
	var jwtSecretKey = []byte(cfg.SecretKey)
	router := chi.NewRouter()
	router.Use(JSONMiddleware)
	router.Mount("/api", API(storage, log, jwtSecretKey, conn))
	return router
}

func API(storage *db.Storage, log *slog.Logger, jwtSecretKey []byte, conn *ldap.Conn) chi.Router {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Post("/login/", handlers.LoginHandler(log, storage, jwtSecretKey, conn))
		r.Post("/refresh/", handlers.RefreshHandler(log, storage, jwtSecretKey))
		r.With(JwtAuthMiddleware(jwtSecretKey, storage)).Get(
			"/me/", handlers.MeHandler(log),
		)
		r.With(JwtAuthMiddleware(jwtSecretKey, storage)).Patch(
			"/{id}/", handlers.PartialUpdateHandler(log, storage),
		)
	})
	return r
}
