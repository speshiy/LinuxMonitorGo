package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/models/mmonitor"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//Migrate created and update tables
func Migrate(c *gin.Context) {
	MigrateDo(c)
	c.String(http.StatusOK, "Migration command done")
}

//MigrateDo migrate all user models
func MigrateDo(c *gin.Context) {
	db, err := database.OpenDatabase(settings.DatabaseMonitorHost, settings.DatabaseMonitorName, settings.DatabaseMonitorUser, settings.DatabaseMonitorPassword, "UTC")
	if err != nil {
		if c != nil {
			c.String(http.StatusOK, "Connection to DB in service FAILED. %s", err.Error())
		}
		log.Println(err.Error())
		return
	}
	defer db.Close()

	db.AutoMigrate(
		&mmonitor.Monitor{},
	)

	log.Println("Migration done.")
}
