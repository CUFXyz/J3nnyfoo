package main

import (
	_ "jennyfood/docs"
	initProj "jennyfood/internal/init"
	"jennyfood/internal/srv"
	"log"

	"github.com/gin-gonic/gin"
)

// Не уверен что оно тут нужно в init()
// Можно просто перенести в main и все
func init() {
	initProj.LoadEnvVar()
}

// @title			J3nnyFoo JSON Project
// @version		0.1
// @description	Small project what accepts jsons and storing it in the postgreSQL
// @host			localhost:9090
func main() {
	// initProj.LoadEnvVar() можно и сюда в целом

	log.Fatal(
		srv.RunGinServer(
			gin.Default(),
		),
	)
}
