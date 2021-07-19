package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/controllers/cmonitor"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//InitRoutes init routes
func InitRoutes(router *gin.Engine) *gin.Engine {
	router.GET("/api/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	routes := router.Group("/api/")
	routes.Use(MainMiddleware())
	{
		routes.POST("/monitor", cmonitor.Post)

		routes.GET("/service/start", cmonitor.Start)
		routes.GET("/service/stop", cmonitor.Stop)
		routes.GET("/service/restart", cmonitor.Restart)
	}

	return router
}

//MainMiddleware open connect to monitoring database
func MainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var DB *gorm.DB

		DB, err = database.OpenDatabase(settings.DatabaseMonitorHost, settings.DatabaseMonitorName, settings.DatabaseMonitorUser, settings.DatabaseMonitorPassword, "UTC")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}
		defer DB.Close()

		c.Set("DB", DB)

		c.Next()
	}
}
