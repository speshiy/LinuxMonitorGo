package cbackup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/LinuxMonitorGo/_main/models/mbackup"
)

//GetBackupsByDBName return all backups
func GetBackupsByDBName(c *gin.Context) {
	var err error
	var backups []mbackup.Backup

	dbName := c.Query("db")

	err = mbackup.GetBackupsByDBName(c, dbName, &backups, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": backups})
}
