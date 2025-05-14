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
)

func (p *Handler) SetupRegisterData(userd models.RegisterData) (models.User, error) {
	var newusid int64
	usid, err := p.db.Pg.Exec("SELECT COUNT(*) FROM users")
	if err != nil {
		return models.User{}, fmt.Errorf("error due setuping register data")
	}
	newusid, _ = usid.RowsAffected()
	return models.User{
		Uid: int(newusid),
		UserData: models.UserData{
			Email:    userd.Email,
			Password: string(auth.CryptPassword(userd.Password)),
			Role:     "user",
		},
	}, nil
}

type TokenResponse struct {
	Token string `json:"token"`
}

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
// @Success		200	{object}	models.Status
// @Failure		400	{object}	models.Status
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
	p.db.SentProductPGSQL(bytes)
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
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			"Error to read your request",
		)
	}
	json.Unmarshal(bytes, &user)
	userd, err := p.SetupRegisterData(user)
	if err != nil {
		fmt.Printf("%v", err)
	}
	RegDataToSend, err := json.Marshal(userd)
	if err != nil {
		fmt.Printf("%v", err)
	}
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

	newuser, err := p.db.GetUserFromPGSQL(logindata)
	fmt.Println(newuser.Password)
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(
			http.StatusBadGateway,
			"Error while working with db",
		)
		return
	}
	err = auth.AuthFunc([]byte(newuser.Password), []byte(logindata.Password))
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(
			http.StatusBadGateway,
			fmt.Sprintf("%v", err),
		)
		return
	}

	token := auth.GenerateToken(logindata)
	response := TokenResponse{
		Token: token,
	}
	c.SetCookie("token", response.Token, 3600, "/login", "localhost", true, true)
	c.JSON(
		http.StatusOK,
		response,
	)

}
