package srv

import (
	_ "jennyfood/docs"

	"jennyfood/internal/auth"
	database "jennyfood/internal/db"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func RunGinServer(engine *gin.Engine, pgHandler *database.Handler, auth auth.AuthInstance) error {

	engine.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL(
				"http://localhost:9090/swagger/doc.json"),
		),
	)

	// Group of endpoints for logged users
	userGroup := engine.Group("/usr")
	userGroup.Use(auth.AuthHandler)
	userGroup.POST("/send", pgHandler.Send)
	userGroup.GET("/usr/:email", pgHandler.ReadingCache)
	userGroup.GET("/data", pgHandler.Index)
	//

	// For losers
	engine.POST("/register", pgHandler.RegisterUser)
	engine.POST("/login", pgHandler.LoginUser)
	engine.GET("/dbstatus", pgHandler.DbStatus)
	return engine.Run(":9090")
}
