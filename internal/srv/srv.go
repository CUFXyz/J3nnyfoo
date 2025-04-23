package srv

import (
	"io"
	_ "jennyfood/docs"
	database "jennyfood/internal/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const connectString = "" // Gonna be deleted soon

type Status struct {
	Status string `json:"PGstatus"`
}

// @Summary		Status of PostgreSQL
// @Description	Returns status of PostgreSQL in json
// @Produce		json
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/dbstatus [get]
func dbStatus(c *gin.Context) { // Endpoint to get info about access to the PG
	var dbStatus Status
	if db, err := database.ConnectToPGSQL(connectString); err != nil {
		dbStatus.Status = "DOWN"
		c.JSON(http.StatusNotFound, dbStatus)
		defer db.Close()
	}
	dbStatus.Status = "UP"
	c.JSON(http.StatusOK, dbStatus)
}

// @Summary		Sending data to PostgreSQL
// @Description	Sending JSON to service and saving in PostgreSQL
// @Accept			json
// @Param			name body string false "name of the product"
// @Param			price body number false "price of the product"
// @Param			type body string false "type of the product"
// @Param			owner body string false "brand or name"
// @Success		200
// @Failure		400
// @Failure		404
// @Router			/sent [post]
func sent(c *gin.Context) {
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

func RunGinServer(engine *gin.Engine) {

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:9090/swagger/doc.json")))

	engine.POST("/sent", sent)
	engine.GET("/dbstatus", dbStatus)
	engine.GET("/data", index)
	engine.Run(":9090")
}
