package srv

import (
	"encoding/json"
	"fmt"
	"io"
	database "jennyfood/internal/db"
	"log"
	"net/http"
)

type jsonPlaceholder struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Type  string  `json:"type"`
	Owner string  `json:"owner"`
}

var placeholder1 = jsonPlaceholder{
	Name:  "AirMax",
	Price: 30,
	Type:  "Boots",
	Owner: "Nike",
}

const connectString = "host=localhost port=5454 dbname=postgres user=postgres password=Verynice1qwe sslmode=disable"

var placeholders = []jsonPlaceholder{placeholder1}

func dbStatus(r http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		db, err := database.ConnectToPGSQL(connectString)
		if err != nil {
			fmt.Printf("%v", err)
			r.Write([]byte(fmt.Sprintf("%v", err)))
			db.Close()
		} else {
			r.Write([]byte("STATUS: OK"))
			db.Close()
		}
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

func DBsent(r http.ResponseWriter, req *http.Request) {
	//var data jsonPlaceholder
	switch req.Method {
	case http.MethodPost:
		bytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatalf("%v", err)
		}
		database.SentPGSQL(connectString, bytes)

	default:
		r.WriteHeader(http.StatusBadRequest)
		r.Write([]byte("Unsupported request method"))
	}
}

func index(r http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if data, err := database.GetFromPGSQL(connectString); err != nil {
			panic(err)
		} else {
			r.Header().Set("Content-Type", "application/json")
			r.Write(data)
		}
	default:
		r.WriteHeader(http.StatusBadRequest)
		r.Write([]byte("Unsupported request method"))
	}
}

func DefaultSetupServer(mux *http.ServeMux) *http.Server {

	mux.HandleFunc("/", index)
	mux.HandleFunc("/sent", sent)
	mux.HandleFunc("/dbsent", DBsent)
	mux.HandleFunc("/dbstatus", dbStatus)

	return &http.Server{
		Handler: mux,
		Addr:    ":9090",
	}
}

func RunServer(srv *http.Server) {
	http.ListenAndServe(srv.Addr, srv.Handler)
}
