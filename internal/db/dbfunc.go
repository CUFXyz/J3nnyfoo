package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type JsonPlaceholder struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Type  string  `json:"type"`
	Owner string  `json:"owner"`
}

func ConnectToPGSQL(connectionSTR string) error {
	db, err := sqlx.Connect("postgres", connectionSTR)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// No need for passing db instance anywhere else.
	// Just close it here.
	defer db.Close()
	return nil
}

func SentPGSQL(connectionSTR string, data []byte) {
	var jdata JsonPlaceholder
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
		return nil, fmt.Errorf("error: %v", err)
	}

	items, item := []JsonPlaceholder{}, JsonPlaceholder{}
	rows, err := db.Query("SELECT name, price, type, owner FROM data")
	if err != nil {
		return nil, fmt.Errorf("ERROR: %w", err)
	}

	for rows.Next() {
		err = rows.Scan(&item.Name, &item.Price, &item.Type, &item.Owner)
		if err != nil {
			fmt.Printf("ERROR ROWS: %v", err)
			continue
		}
		items = append(items, item)
	}
	defer rows.Close()
	data, err = json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %w", err)
	}
	db.Close()
	return data, nil
}
