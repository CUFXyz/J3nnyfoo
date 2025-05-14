package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PgConfig struct {
	Constr string `env:"CONSTR"`
}

type AuthConfig struct {
	Secret string `env:"SECRET"`
}

type ConfigGlobal struct {
	Pgcfg   PgConfig
	AuthCfg AuthConfig
}

func InitConfig() *ConfigGlobal {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(
			"Error loading .env file: %v",
			err,
		)
	}

	constr := os.Getenv("CONSTR")
	secret := os.Getenv("SECRET")

	return &ConfigGlobal{
		Pgcfg: PgConfig{
			Constr: constr,
		},
		AuthCfg: AuthConfig{
			Secret: secret,
		},
	}
}
