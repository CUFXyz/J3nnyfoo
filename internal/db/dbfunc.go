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

func ConnectToPGSQL(cfg config.ConfigGlobal) *Postgres {

	db, err := sqlx.Connect("postgres", cfg.Pgcfg.Constr)
	if err != nil {
		fmt.Printf("failed to connect to db: %v", err)
		return &Postgres{
			Pg:  nil,
			cfg: cfg.Pgcfg.Constr,
		}
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	return &Postgres{
		Pg:  db,
		cfg: cfg.Pgcfg.Constr,
	}
}

// Func sending data to the product table
func (p *Postgres) SentProductPGSQL(data []byte) {
	var jdata models.JsonPlaceholder
	json.Unmarshal(data, &jdata)
	result, err := p.Pg.Exec("INSERT INTO data (Name, price, type, owner) VALUES ($1, $2, $3, $4)", jdata.Name, jdata.Price, jdata.Type, jdata.Owner)
	if err != nil {
		fmt.Printf("Error due execute insert to the data table\n")
		println(result)
	}
}

// Func sending user data to the users table
func (p *Postgres) SendRegDataPGSQL(data []byte) {
	var usr models.User
	json.Unmarshal(data, &usr)
	result, err := p.Pg.Exec("INSERT INTO users (uid, email, password, role, token) VALUES ($1, $2, $3, $4, $5)", usr.Uid, usr.Email, usr.Password, usr.Role, usr.Token)
	if err != nil {
		fmt.Printf("Error due execute insert to the data table\n")
	}
	rws, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Can't get affected rows\n")
	}
	fmt.Printf("%v rows affected", rws)
}

func (p *Postgres) GetUserFromPGSQL(user models.RegisterData) (*models.RegisterData, error) {
	usr := &models.RegisterData{}
	err := p.Pg.QueryRow("SELECT email, password FROM users WHERE email LIKE $1", &user.Email).Scan(&usr.Email, &usr.Password)
	if err != nil {
		return &models.RegisterData{}, err
	}
	return usr, nil
}

func (p *Postgres) GetFromPGSQL() ([]byte, error) {
	var data []byte
	items, item := []models.JsonPlaceholder{}, models.JsonPlaceholder{}
	rows, err := p.Pg.Query("SELECT name, price, type, owner FROM data")
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
