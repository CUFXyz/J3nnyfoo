package main

import (
	_ "jennyfood/docs"
	initProj "jennyfood/internal/init"
	"jennyfood/internal/srv"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	initProj.LoadEnvVar()
}

// @title			J3nnyFoo JSON Project
// @version		0.1
// @description	Small project what accepts jsons and storing it in the postgreSQL
// @host			localhost:9090
func main() {
	log.Fatal(
		srv.RunGinServer(
			gin.Default(),
		),
	)
}
