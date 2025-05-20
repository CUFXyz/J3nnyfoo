package database

import (
	"encoding/json"
	"fmt"
	"io"
	"jennyfood/internal/auth"
	"jennyfood/internal/storage"
	"jennyfood/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response after loggining user into our service
type TokenResponse struct {
	Token string `json:"token"`
}

// Handler contains db instance and cache instance
type Handler struct {
	db         *Postgres
	storage    *storage.Cache
	authSystem auth.AuthInstance
}

// Constructor for our handler struct
func InitializeHandler(p *Postgres, storage *storage.Cache, auth auth.AuthInstance) *Handler {
	return &Handler{
		db:         p,
		storage:    storage,
		authSystem: auth,
	}
}

// A little constructor to sending registerdata to postgresql
func (p *Handler) SetupRegisterData(userd models.RegisterData) (models.User, error) {
	// variable for id
	var newusid int64
	// Just trying to get number of rows to setup id
	// I gonna kill myself with that bullshit later
	usid, err := p.db.Pg.Exec("SELECT COUNT(*) FROM users")
	if err != nil {
		return models.User{}, fmt.Errorf("error due setuping register data")
	}
	newusid, _ = usid.RowsAffected()
	// Preparing and returning struct to sending
	return models.User{
		Uid: int(newusid),
		UserData: models.UserData{
			Email:    userd.Email,
			Password: string(p.authSystem.CryptPassword(userd.Password)),
			Role:     "user",
		},
	}, nil
}

// @Summary		Status of PostgreSQL
// @Description	Returns status of PostgreSQL in json
// @Produce		json
// @Success		200	{object}	models.Status
// @Failure		400	{object}	models.Status
// @Failure		404
// @Router			/dbstatus [get]
func (p *Handler) DbStatus(c *gin.Context) { // Endpoint to get info about access to the PG
	var dbStatus models.Status

	switch p.db.Pg {
	case nil: // If we can't connect to our DB
		dbStatus.CurStatus = "DOWN"
		c.JSON(
			http.StatusNotFound,
			dbStatus,
		)
	default: // We can connect and do things with our DB
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
// @Param			name	body	models.JsonPlaceholder	true	"Actual data to store in db"
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/send [post]
func (p *Handler) Send(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body) // Writing the body of request
	if err != nil {
		log.Fatalf("%v", err) // Server killing itself after that error
	}
	p.db.SendProductPGSQL(bytes)
}

// @Summary	Registrate user to get new features
// @Accept		json
// @Param		name	body	models.RegisterData	true	"Register user in this service"
// @Success	200
// @Failure	400
// @Failure	404
// @Router		/register [post]
func (p *Handler) RegisterUser(c *gin.Context) {
	var user models.RegisterData

	//Reading data from request
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			"Error to read your request",
		)
	}
	json.Unmarshal(bytes, &user)

	// Small constructor to preparing models.user to sending into DB
	userd, err := p.SetupRegisterData(user)
	if err != nil {
		fmt.Printf("%v", err)
	}
	RegDataToSend, err := json.Marshal(userd)
	if err != nil {
		fmt.Printf("%v", err)
	}

	//Sending prepared models.user to the DB
	p.db.SendRegDataPGSQL(RegDataToSend)
}

// @Summary	loggining user into service and creating token
// @Accept		json
// @Param		name	body	models.RegisterData	true	"registerdata"
// @Success	200
// @Failure	400
// @Failure	404
// @Router		/login [post]
func (p *Handler) LoginUser(c *gin.Context) {
	var (
		logindata models.RegisterData
	)

	// Reading data from request body
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

	// Getting user data from DB, logindata must be email+password
	// If something wrong - returns badgtw and error
	newuser, err := p.db.GetUserFromPGSQL(logindata)
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(
			http.StatusBadGateway,
			"Error while working with db",
		)
		return
	}

	// Checking passwords hashed and !hashed
	// If hashed != !hashed returns badgtw
	// else going further to generating token
	err = p.authSystem.AuthFunc([]byte(newuser.Password), []byte(logindata.Password))
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(
			http.StatusBadGateway,
			fmt.Sprintf("%v", err),
		)
		return
	}

	// Generating token for logged user
	token := p.authSystem.GenerateToken(logindata)

	//Writing token to the cache storage
	p.storage.WriteCache(token, logindata.Email)

	// Making Response to user
	response := TokenResponse{
		Token: token,
	}

	//Setting token to store in the user browser and responding to him
	c.SetCookie("token", response.Token, 3600, "/login", "localhost", true, true)
	c.JSON(
		http.StatusOK,
		response,
	)
}

func (p *Handler) ReadingCache(c *gin.Context) {
	email := c.Param("email")
	var emaildb string
	query := "SELECT email FROM users WHERE email = $1"
	err := p.db.Pg.QueryRow(query, &email).Scan(&emaildb)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("%v | %v", error.Error(err), "pg skill issue"),
		)
		c.Abort()
		return
	}
	result, err := p.storage.GetValue(email)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("%v | %v", error.Error(err), result),
		)
		return
	}
	if email == emaildb {
		c.JSON(
			http.StatusOK,
			result,
		)
		return
	} else {
		c.JSON(
			http.StatusBadRequest,
			"User not found",
		)
		c.Abort()
		return
	}
}
