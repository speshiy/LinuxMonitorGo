package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/speshiy/LinuxMonitorGo/service/controllers"
)

//InitRoutes init routes
func InitRoutes(router *gin.Engine) *gin.Engine {

	g1 := router.Group("/api/service", gin.BasicAuth(gin.Accounts{
		"migrate": "<your-password>",
	}))
	{
		g1.GET("/migrate", controllers.Migrate)
	}

	return router
}
