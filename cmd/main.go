package main

import (
	_ "jennyfood/docs"
	"jennyfood/internal/auth"
	"jennyfood/internal/config"
	database "jennyfood/internal/db"
	"jennyfood/internal/srv"
	"jennyfood/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

// @title			J3nnyFoo JSON Project
// @version		1.0
// @description	Small project what accepts jsons and storing it in the postgreSQL
// @host			localhost:9090

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {

	cfg := config.InitConfig()
	postgres, err := database.ConnectToPGSQL(*cfg)
	cache := storage.NewCache()
	auth := auth.AuthInstance{AuthCfg: cfg.AuthCfg, Cache: cache, Pgdb: postgres.Pg}
	if err != nil {
		log.Fatalf("%v", err)
	}
	pgHandler := database.InitializeHandler(postgres, cache, auth)

	log.Fatal(
		srv.RunGinServer(
			gin.Default(),
			pgHandler,
			auth,
		),
	)
	defer postgres.Pg.Close()
}
