package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/OR-Sasaki/cat-backend/controllers"

	_ "github.com/OR-Sasaki/cat-backend/docs"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		controllers.UsersRouter(api)
		controllers.SeriesRouter(api)
		controllers.OutfitsRouter(api)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
