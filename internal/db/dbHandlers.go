package database

import (
	"encoding/json"
	"fmt"
	"io"
	"jennyfood/internal/auth"
	"jennyfood/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	db *Postgres
}

func InitializeHandler(p *Postgres) *Handler {
	return &Handler{
		db: p,
	}
}

// @Summary		Status of PostgreSQL
// @Description	Returns status of PostgreSQL in json
// @Produce		json
// @Success		200 {object} models.Status
// @Failure		400 {object} models.Status
// @Failure		404
// @Router			/dbstatus [get]
func (p *Handler) DbStatus(c *gin.Context) { // Endpoint to get info about access to the PG
	var dbStatus models.Status

	switch p.db.Pg {
	case nil:
		dbStatus.CurStatus = "DOWN"
		c.JSON(
			http.StatusNotFound,
			dbStatus,
		)
	default:
		dbStatus.CurStatus = "UP"
		c.JSON(
			http.StatusOK,
			dbStatus,
		)
		return
	}
}

// @Summary		Returns whole data what stored in PostgreSQL
// @Description	Returns an array of JSONS with data
// @Produce		json
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/data [get]
func (p *Handler) Index(c *gin.Context) {
	if data, err := p.db.GetFromPGSQL(); err == nil { // Getting info from PG
		c.Data(http.StatusOK, "application/json", data)
	}
}

// @Summary		Sending data to PostgreSQL
// @Description	Sending JSON to service and saving in PostgreSQL
// @Accept			json
// @Param			name body models.JsonPlaceholder true "Actual data to store in db"
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/send [post]
func (p *Handler) Send(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body) // Writing the body of request
	if err != nil {
		log.Fatalf("%v", err) // Server killing itself after that error
	}
	p.db.SentProductPGSQL(bytes)
}

// @Summary Registrate user to get new features
// @Accept json
// @Param name body models.RegisterData true "Register user in this service"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /register [post]
func (p *Handler) RegisterUser(c *gin.Context) {
	var user models.RegisterData
	var userd models.User
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			"Error to read your request",
		)
	}
	json.Unmarshal(bytes, &user)
	userd.Uid = uuid.New()
	userd.Email = user.Email
	userd.Password = string(auth.CryptPassword(user.Password))
	userd.Role = "user"
	userd.Token = ""
	RegDataToSend, err := json.Marshal(userd)
	if err != nil {
		fmt.Printf("%v", err)
	}
	p.db.SendRegDataPGSQL(RegDataToSend)
}

func (p *Handler) LoginUser(c *gin.Context) {
	var (
		logindata models.RegisterData
		dbinfo    models.RegisterData
	)

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			"Error to read your request",
		)
		return
	}
	err = json.Unmarshal(bytes, &logindata)
	if err != nil {
		fmt.Printf("Error due unmarshal json")
		return
	}

	newuser, err := p.db.GetUserFromPGSQL(dbinfo)
	if err != nil {
		fmt.Printf("Error due getting user from db")
		return
	}
	result, err := auth.AuthFunc([]byte(newuser.Password), []byte(logindata.Password))
	if err != nil {
		fmt.Printf("Error due auth user")
		return
	}

	if !result {
		c.JSON(
			http.StatusForbidden,
			"Email or Password is not correct",
		)
	}

}
