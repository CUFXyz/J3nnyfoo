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
// @version		0.1
// @description	Small project what accepts jsons and storing it in the postgreSQL
// @host			localhost:9090
func main() {

	cfg := config.InitConfig()
	cache := storage.NewCache()
	auth := auth.AuthInstance{AuthCfg: cfg.AuthCfg, Cache: cache}
	postgres, err := database.ConnectToPGSQL(*cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	pgHandler := database.InitializeHandler(postgres, auth.AuthCfg, cache)
	log.Fatal(
		srv.RunGinServer(
			gin.Default(),
			pgHandler,
			auth,
		),
	)
	defer postgres.Pg.Close()
}
