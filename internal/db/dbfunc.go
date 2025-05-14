package database

import (
	"encoding/json"
	"fmt"
	"jennyfood/internal/config"
	"jennyfood/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Pg  *sqlx.DB
	cfg string
}

func ConnectToPGSQL(cfg config.ConfigGlobal) (*Postgres, error) {

	db, err := sqlx.Connect("postgres", cfg.Pgcfg.Constr)
	if err != nil {
		return &Postgres{
			Pg:  nil,
			cfg: cfg.Pgcfg.Constr,
		}, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	return &Postgres{
		Pg:  db,
		cfg: cfg.Pgcfg.Constr,
	}, nil
}

// Func sending data to the product table
func (p *Postgres) SentProductPGSQL(data []byte) {
	var jdata models.JsonPlaceholder
	query := "INSERT INTO data (Name, price, type, owner) VALUES ($1, $2, $3, $4)"
	json.Unmarshal(data, &jdata)
	result, err := p.Pg.Exec(query, jdata.Name, jdata.Price, jdata.Type, jdata.Owner)
	if err != nil {
		fmt.Printf("Error due execute insert to the data table\n")
		println(result)
	}
}

// Func sending user data to the users table
func (p *Postgres) SendRegDataPGSQL(data []byte) {
	var usr models.User
	query := "INSERT INTO users (userID, email, password, role) VALUES ($1, $2, $3, $4)"
	json.Unmarshal(data, &usr)
	_, err := p.Pg.Exec(query, usr.Uid, usr.Email, usr.Password, usr.Role)
	if err != nil {
		fmt.Printf("Error due execute insert to the data table\n")
	}
}

func (p *Postgres) GetUserFromPGSQL(user models.RegisterData) (*models.RegisterData, error) {
	usr := &models.RegisterData{}
	query := "SELECT email, password FROM users WHERE email LIKE $1"
	err := p.Pg.QueryRow(query, &user.Email).Scan(&usr.Email, &usr.Password)
	if err != nil {
		return &models.RegisterData{}, err
	}
	return usr, nil
}

func (p *Postgres) GetFromPGSQL() ([]byte, error) {
	var data []byte
	items, item := []models.JsonPlaceholder{}, models.JsonPlaceholder{}
	query := "SELECT name, price, type, owner FROM data"
	rows, err := p.Pg.Query(query)
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
	return data, nil
}
