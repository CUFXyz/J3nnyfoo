package srv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type jsonPlaceholder struct {
	Price float32 `json:"price"`
	Type  string  `json:"type"`
	Owner string  `json:"owner"`
}

func index(r http.ResponseWriter, req *http.Request) {
	placeholder1 := jsonPlaceholder{
		Price: 30,
		Type:  "Boots",
		Owner: "Nike",
	}
	if req.Method == http.MethodGet {
		if body, err := json.Marshal(placeholder1); err != nil {
			log.Fatalf("%v", err)
		} else {
			r.Header().Set("Content-Type", "application/json")
			r.Write(body)
		}
	} else if req.Method == http.MethodPost {
		decoder := json.NewDecoder(req.Body)
		var test jsonPlaceholder
		if err := decoder.Decode(&test); err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Println(test)
	}

}

func DefaultSetupServer(mux *http.ServeMux) *http.Server {

	mux.HandleFunc("/", index)

	return &http.Server{
		Handler: mux,
		Addr:    ":9090",
	}
}

func RunServer(srv *http.Server) {
	http.ListenAndServe(srv.Addr, srv.Handler)
}
