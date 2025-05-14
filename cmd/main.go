package main

import (
	_ "jennyfood/docs"
	"jennyfood/internal/config"
	database "jennyfood/internal/db"
	"jennyfood/internal/srv"
	"log"

	"github.com/gin-gonic/gin"
)

// @title			J3nnyFoo JSON Project
// @version		0.1
// @description	Small project what accepts jsons and storing it in the postgreSQL
// @host			localhost:9090
func main() {

	cfg := config.InitConfig()
	postgres, err := database.ConnectToPGSQL(*cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	pgHandler := database.InitializeHandler(postgres)
	log.Fatal(
		srv.RunGinServer(
			gin.Default(),
			pgHandler,
		),
	)
	defer postgres.Pg.Close()
}
