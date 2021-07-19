package cmonitor

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

//Restart перезагружаем прогамму
func Restart(c *gin.Context) {
	var err error

	cmd := exec.Command("systemctl", "restart", "tuvis")
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Служба не перезагружена, так как в случае перезагрузки это сообщение не могло вернуться с сервера"})
}
