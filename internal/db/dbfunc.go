package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type jsonPlaceholder struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Type  string  `json:"type"`
	Owner string  `json:"owner"`
}

func ConnectToPGSQL(connectionSTR string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connectionSTR)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

func SentPGSQL(connectionSTR string, data []byte) {
	var jdata jsonPlaceholder
	json.Unmarshal(data, &jdata)
	db, err := sql.Open("postgres", connectionSTR)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	result, err := db.Exec("INSERT INTO data (Name, price, type, owner) VALUES ($1, $2, $3, $4)", jdata.Name, jdata.Price, jdata.Type, jdata.Owner)
	if err != nil {
		fmt.Printf("Error: %v", err)
		println(result)
	}
	db.Close()
}

func GetFromPGSQL(connectionSTR string) ([]byte, error) {
	var data []byte
	db, err := sql.Open("postgres", connectionSTR)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	var items, item = []jsonPlaceholder{}, jsonPlaceholder{}
	count, _ := getRows(db)
	for i := 0; i < count; i++ {
		err = db.QueryRow("SELECT name, price, type, owner FROM data WHERE id = $1", i+1).Scan(&item.Name, &item.Price, &item.Type, &item.Owner)
		if err != nil {
			fmt.Printf("ERROR: %v", err)
		}
		items = append(items, item)
	}
	data, err = json.Marshal(items)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	db.Close()
	return data, nil
}

func getRows(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM data").Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
