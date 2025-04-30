package srv

import (
	_ "jennyfood/docs"

	database "jennyfood/internal/db"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func RunGinServer(engine *gin.Engine, pgHandler *database.Handler) error {

	engine.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL(
				"http://localhost:9090/swagger/doc.json"),
		),
	)
	engine.POST("/register", pgHandler.RegisterUser)
	engine.POST("/send", pgHandler.Send)
	engine.POST("/login", pgHandler.LoginUser)
	engine.GET("/dbstatus", pgHandler.DbStatus)
	engine.GET("/data", pgHandler.Index)
	return engine.Run(":9090")
}
