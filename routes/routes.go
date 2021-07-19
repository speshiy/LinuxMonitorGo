package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/controllers/cmonitor"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//InitRoutes инициализирует пути
func InitRoutes(router *gin.Engine) *gin.Engine {
	router.GET("/api/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	routes := router.Group("/api/")
	routes.Use(MainMiddleware())
	{
		routes.GET("/runtime/restart", cmonitor.Restart)
	}

	return router
}

//MainMiddleware открывает коннект к главной базе
func MainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var DB *gorm.DB

		DB, err = database.OpenDatabase(settings.DatabaseName, settings.DatabaseMonitorUser, settings.DatabaseMonitorPassword, "UTC")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}
		defer DB.Close()

		c.Set("DB", DB)

		c.Next()
	}
}
