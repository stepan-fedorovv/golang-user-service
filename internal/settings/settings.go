package settings

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"log/slog"
	"os"
	"time"
)

const (
	envDev        = "development"
	envProduction = "production"
)

type Config struct {
	Env         string `yaml:"env" env:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env:"storage_path" env-required:"true"`
	SecretKey   string `yaml:"secret_key" env:"secret_key" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	LDAP        `yaml:"ldap"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"idle_timeout" env-default:"60s"`
}

type LDAP struct {
	BaseDN       string `yaml:"base_dn" env:"base_dn" env-required:"true"`
	BindDN       string `yaml:"bind_dn" env:"bind_dn" env-required:"true"`
	BindPassword string `yaml:"bind_password" env:"bind_password" env-required:"true"`
	LdapHost     string `yaml:"ldap_host" env:"ldap_host" env-required:"true"`
	LdapPort     int    `yaml:"ldap_port" env:"ldap_port" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}

func SetupLogger(env string) *slog.Logger {
	var initLog *slog.Logger
	switch env {
	case envDev:
		initLog = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProduction:
		initLog = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	initLog.With(
		slog.String("env", env),
	)

	return initLog
}

func Connect(cfg *Config) (*ldap.Conn, error) {

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", cfg.LdapHost, cfg.LdapPort))
	if err != nil {
		return nil, err
	}
	if err := conn.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
		return nil, err
	}
	return conn, nil
}
