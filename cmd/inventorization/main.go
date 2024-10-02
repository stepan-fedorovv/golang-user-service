package main

import (
	"app/internal/db"
	settings2 "app/internal/settings"
	"app/utils/logger"
	"errors"
	"net/http"
	"os"
)

func main() {
	cfg := settings2.MustLoad()
	log := settings2.SetupLogger(cfg.Env)
	conn, err := settings2.Connect(cfg)
	if err != nil {
		log.Error("Error connection to ldap: " + err.Error())
		os.Exit(1)
	}
	log.Info("Starting application")
	log.Debug("Debug mode is active now")
	storage, err := db.New(cfg.StoragePath)
	_ = storage
	if err != nil {
		log.Error(err.Error())
		log.Error("Error opening storage")
		os.Exit(1)
	}
	router := settings2.Router(storage, log, cfg, conn)
	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to start server", logger.Err(err))
	}
}
