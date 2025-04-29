package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type GlobalConfig struct {
	Postgres PostgresConfig
	Auth     AuthConfig
}

func NewGlobalConfig() *GlobalConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(
			"Error loading .env file: %v",
			err,
		)
	}
	// эту хуйню можно не делать, почитай про этот пакет ниже
	// github.com/caarlos0/env/v11
	constr := os.Getenv("CONSTR")
	secret := os.Getenv("SECRET")

	return &GlobalConfig{
		Postgres: PostgresConfig{
			Constr: constr,
		},
		Auth: AuthConfig{
			Secret: secret,
		},
	}
}

type PostgresConfig struct {
	Constr string `env:"CONSTR"`
}

type AuthConfig struct {
	Secret string `env:"SECRET"`
}
