package srv

import (
	"io"
	"log"
	"net/http"

	_ "jennyfood/docs"

	database "jennyfood/internal/db"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

const connectString = "host=localhost port=5454 dbname=postgres user=postgres password=Verynice1qwe sslmode=disable" // Gonna be deleted soon

type Status struct {
	Status string `json:"pg_status"`
}

// @Summary		Status of PostgreSQL
// @Description	Returns status of PostgreSQL in json
// @Produce		json
// @Success		200 {object} Status
// @Failure		400 {object} Status
// @Failure		404
// @Router			/dbstatus [get]
func dbStatus(c *gin.Context) { // Endpoint to get info about access to the PG
	var dbStatus Status

	switch database.ConnectToPGSQL(connectString) {
	case nil:
		dbStatus.Status = "UP"
		c.JSON(
			http.StatusOK,
			dbStatus,
		)
		return
	default:
		dbStatus.Status = "DOWN"
		c.JSON(
			http.StatusNotFound,
			dbStatus,
		)
	}
}

// @Summary		Sending data to PostgreSQL
// @Description	Sending JSON to service and saving in PostgreSQL
// @Accept			json
// @Param			name body database.JsonPlaceholder true "Actual data to store in db"
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/send [post]
func send(c *gin.Context) {
	bytes, err := io.ReadAll(c.Request.Body) // Writing the body of request
	if err != nil {
		log.Fatalf("%v", err) // Server killing itself after that error
	}
	database.SentPGSQL(connectString, bytes)
}

// @Summary		Returns whole data what stored in PostgreSQL
// @Description	Returns an array of JSONS with data
// @Produce		json
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/data [get]
func index(c *gin.Context) {
	if data, err := database.GetFromPGSQL(connectString); err == nil { // Getting info from PG
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
	engine.POST("/send", send)
	engine.GET("/dbstatus", dbStatus)
	engine.GET("/data", index)
	return engine.Run(":9090")
}
