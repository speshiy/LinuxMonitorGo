package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/speshiy/LinuxMonitorGo/service/controllers"
)

//InitRoutes инициализирует пути
func InitRoutes(router *gin.Engine) *gin.Engine {

	g1 := router.Group("/api/service", gin.BasicAuth(gin.Accounts{
		"migrate": "843g43r-2kfp=2-342kfds3",
	}))
	{
		g1.GET("/migrate", controllers.Migrate)
	}

	return router
}
