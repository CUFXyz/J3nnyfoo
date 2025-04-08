package main

import (
	"jennyfood/internal/srv"
	"net/http"
)

func main() {
	srv.RunServer(srv.DefaultSetupServer(http.NewServeMux()))
}
