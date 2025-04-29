package srv

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "jennyfood/docs"
	"jennyfood/models"

	"jennyfood/internal/auth"
	database "jennyfood/internal/db"

	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	var logindata, dbinfo models.RegisterData
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

	newuser, err := database.GetUserFromPGSQL(os.Getenv("CONSTR"), dbinfo)
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

// @Summary Registrate user to get new features
// @Accept json
// @Param name body models.RegisterData true "Register user in this service"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /register [post]
func RegisterUser(c *gin.Context) {
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
	database.SendRegDataPGSQL(os.Getenv("CONSTR"), RegDataToSend)
}

// @Summary		Status of PostgreSQL
// @Description	Returns status of PostgreSQL in json
// @Produce		json
// @Success		200 {object} models.Status
// @Failure		400 {object} models.Status
// @Failure		404
// @Router			/dbstatus [get]
func dbStatus(c *gin.Context) { // Endpoint to get info about access to the PG
	var dbStatus models.Status

	switch database.ConnectToPGSQL(os.Getenv("CONSTR")) {
	case nil:
		dbStatus.CurStatus = "UP"
		c.JSON(
			http.StatusOK,
			dbStatus,
		)
		return
	default:
		dbStatus.CurStatus = "DOWN"
		c.JSON(
			http.StatusNotFound,
			dbStatus,
		)
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
func send(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body) // Writing the body of request
	if err != nil {
		log.Fatalf("%v", err) // Server killing itself after that error
	}
	database.SentProductPGSQL(os.Getenv("CONSTR"), bytes)
}

// @Summary		Returns whole data what stored in PostgreSQL
// @Description	Returns an array of JSONS with data
// @Produce		json
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/data [get]
func index(c *gin.Context) {
	if data, err := database.GetFromPGSQL(os.Getenv("CONSTR")); err == nil { // Getting info from PG
		c.Data(http.StatusOK, "application/json", data)
	}
}

func RunGinServer(engine *gin.Engine) error {
	engine.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL(
				"http://localhost:9090/swagger/doc.json"),
		),
	)
	engine.POST("/register", RegisterUser)
	engine.POST("/send", send)
	engine.GET("/dbstatus", dbStatus)
	engine.GET("/data", index)
	return engine.Run(":9090")
}
