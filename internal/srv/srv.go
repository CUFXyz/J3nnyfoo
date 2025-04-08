package srv

import (
	"encoding/json"
	"log"
	"net/http"
)

type jsonPlaceholder struct {
	Price float32 `json:"price"`
	Type  string  `json:"type"`
	Owner string  `json:"owner"`
}

var placeholder1 = jsonPlaceholder{
	Price: 30,
	Type:  "Boots",
	Owner: "Nike",
}

var placeholders = []jsonPlaceholder{placeholder1}

func sent(r http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var data jsonPlaceholder
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&data); err != nil {
			log.Fatalf("%v", err)
		}
		placeholders = append(placeholders, data)
	}
}

func index(r http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {
		if body, err := json.Marshal(placeholders); err != nil {
			log.Fatalf("%v", err)
		} else {
			r.Header().Set("Content-Type", "application/json")
			r.Write(body)
		}
	}
}

func DefaultSetupServer(mux *http.ServeMux) *http.Server {

	mux.HandleFunc("/", index)
	mux.HandleFunc("/sent", sent)

	return &http.Server{
		Handler: mux,
		Addr:    ":9090",
	}
}

func RunServer(srv *http.Server) {
	http.ListenAndServe(srv.Addr, srv.Handler)
}
