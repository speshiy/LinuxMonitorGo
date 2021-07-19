package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/speshiy/LinuxMonitorGo/_main/controllers/cbackup"
	"github.com/speshiy/LinuxMonitorGo/_main/controllers/cpurge"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/controllers/cmonitor"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//InitRoutes инициализирует пути
func InitRoutes(router *gin.Engine) *gin.Engine {
	router.GET("/api/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	//**************************************************
	//Static FILES ROUTES
	//**************************************************
	router.StaticFS("/var/LinuxMonitorGo/mysql/dumps", http.Dir("/var/LinuxMonitorGo/mysql/dumps"))
	//**************************************************

	routes := router.Group("/api/")
	routes.Use(MainMiddleware())
	{
		routes.GET("/mysql/backups", cbackup.GetBackupsByDBName)
		routes.POST("/mysql/backup", cbackup.BackupByDBName)
		routes.POST("/mysql/backup/restore", cbackup.BackupRestore)
		routes.PUT("/mysql/database/move/from/server/to/server", cbackup.DatabaseMoveFromServerToServer)
		routes.DELETE("/mysql/database/purge", cpurge.PurgeUnusedDatabases)

		// Monitor
		routes.GET("/runtime/restart", cmonitor.Restart)
	}

	return router
}

//MainMiddleware открывает коннект к главной базе
func MainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var DB *gorm.DB

		DB, err = database.OpenDatabase("backup", "root", settings.DBRP, "UTC")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
			return
		}
		defer DB.Close()

		c.Set("DB", DB)

		c.Next()
	}
}
