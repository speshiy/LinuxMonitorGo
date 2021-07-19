package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/LinuxMonitorGo/_main/models/mbackup"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//Migrate all bases
func Migrate(c *gin.Context) {
	MigrateMain(c)
	c.String(http.StatusOK, "Migration command done.")
}

//MigrateMain migrate all user models
func MigrateMain(c *gin.Context) {
	db, err := database.OpenDatabase("backup_tuvis", "root", settings.DBRP, "UTC")
	if err != nil {
		if c != nil {
			c.String(http.StatusOK, "Connection to DB in service FAILED. %s", err.Error())
		}
		log.Println(err.Error())
		return
	}
	defer db.Close()

	db.AutoMigrate(
		&mbackup.Backup{},
	)
	log.Println("Models in DB backup created")

	log.Println("Migration done.")
}
