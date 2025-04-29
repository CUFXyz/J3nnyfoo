package initProj

import (
	"log"

	"github.com/joho/godotenv"
)

// Еще лучше:
// - создать директорию config, где будет лежать глобальный конфиг для
// всех сущностей в проекте, и перенести эту логику парсинга туда,
// а то имхо пакет init выглядит странно. С config более структурированно
func LoadEnvVar() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(
			"Error loading .env file: %v",
			err,
		)
	}
}
