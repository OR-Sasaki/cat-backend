package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/OR-Sasaki/cat-backend/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/users/register", controllers.UserRegister)
	}

	return router
}
