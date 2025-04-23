package main

import (
	_ "jennyfood/docs"
	"jennyfood/internal/srv"

	"github.com/gin-gonic/gin"
)

//	@title			J3nnyFoo JSON Project
//	@version		0.1
//	@description	Small project what accepts jsons and storing it in the postgreSQL
//	@host			localhost:9090
func main() {
	srv.RunGinServer(gin.Default())
}
