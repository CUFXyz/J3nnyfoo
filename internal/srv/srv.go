package srv

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type jsonPlaceholder struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Type  string  `json:"type"`
	Owner string  `json:"owner"`
}

var placeholder1 = jsonPlaceholder{
	Id:    1,
	Name:  "AirMax",
	Price: 30,
	Type:  "Boots",
	Owner: "Nike",
}

var placeholders = []jsonPlaceholder{placeholder1}

func deleteByID(r http.ResponseWriter, req *http.Request) {
	var data jsonPlaceholder
	switch req.Method {
	case http.MethodPost:
		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatalf("%v", err)
		}
		json.Unmarshal(bytes, &data)
		placeholders = removebyField(placeholders, data.Id)
	default:
		r.WriteHeader(http.StatusBadRequest)
		r.Write([]byte("Unsupported request method"))
	}
}

func sent(r http.ResponseWriter, req *http.Request) {
	var data jsonPlaceholder
	switch req.Method {
	case http.MethodPost:
		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatalf("%v", err)
		}
		json.Unmarshal(bytes, &data)
		placeholders = append(placeholders, data)
	default:
		r.WriteHeader(http.StatusBadRequest)
		r.Write([]byte("Unsupported request method"))
	}
}

func index(r http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if body, err := json.Marshal(placeholders); err != nil {
			log.Fatalf("%v", err)
		} else {
			r.Header().Set("Content-Type", "application/json")
			r.Write(body)
		}
	default:
		r.WriteHeader(http.StatusBadRequest)
		r.Write([]byte("Unsupported request method"))
	}
}

func DefaultSetupServer(mux *http.ServeMux) *http.Server {

	mux.HandleFunc("/", index)
	mux.HandleFunc("/sent", sent)
	mux.HandleFunc("/delete", deleteByID)

	return &http.Server{
		Handler: mux,
		Addr:    ":9090",
	}
}

func RunServer(srv *http.Server) {
	http.ListenAndServe(srv.Addr, srv.Handler)
}
