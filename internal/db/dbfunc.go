package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jennyfood/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

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

// Func sending data to the product table
func SentProductPGSQL(connectionSTR string, data []byte) {
	var jdata models.JsonPlaceholder
	json.Unmarshal(data, &jdata)
	db, err := sql.Open("postgres", connectionSTR)
	if err != nil {
		fmt.Printf("Error due opening postgres\n")
	}

	result, err := db.Exec("INSERT INTO data (Name, price, type, owner) VALUES ($1, $2, $3, $4)", jdata.Name, jdata.Price, jdata.Type, jdata.Owner)
	if err != nil {
		fmt.Printf("Error due execute insert to the data table\n")
		println(result)
	}
	db.Close()
}

// Func sending user data to the users table
func SendRegDataPGSQL(connectionSTR string, data []byte) {
	var usr models.User
	json.Unmarshal(data, &usr)
	db, err := sql.Open("postgres", connectionSTR)
	if err != nil {
		fmt.Printf("Error due opening postgres\n")
	}
	result, err := db.Exec("INSERT INTO users (uid, email, password, role, token) VALUES ($1, $2, $3, $4, $5)", usr.Uid, usr.Email, usr.Password, usr.Role, usr.Token)
	if err != nil {
		fmt.Printf("Error due execute insert to the data table\n")
	}
	rws, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Can't get affected rows\n")
	}
	fmt.Printf("%v rows affected", rws)

}

func GetUserFromPGSQL(connectionSTR string, user models.RegisterData) (*models.RegisterData, error) {
	db, err := sql.Open("postgres", connectionSTR)
	if err != nil {
		return &models.RegisterData{}, fmt.Errorf("error due opening postgre")
	}
	row, err := db.Query("SELECT email, password FROM users")
	if err != nil {
		return &models.RegisterData{}, fmt.Errorf("error due Quering to postgre")
	}
	err = row.Scan(&user.Email, &user.Password)
	if err != nil {
		return &models.RegisterData{}, fmt.Errorf("error due Scanning data to struct")
	}
	return &user, nil
}

func GetFromPGSQL(connectionSTR string) ([]byte, error) {
	var data []byte
	db, err := sql.Open("postgres", connectionSTR)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	items, item := []models.JsonPlaceholder{}, models.JsonPlaceholder{}
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
