package cmonitor

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

//Start service
func Start(c *gin.Context) {
	var err error

	service := c.DefaultQuery("service", "")

	if len(service) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Service name empty"})
		return
	}

	cmd := exec.Command("systemctl", "start", service)
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Service started"})
}

//Stop service
func Stop(c *gin.Context) {
	var err error

	service := c.DefaultQuery("service", "")

	if len(service) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Service name empty"})
		return
	}

	cmd := exec.Command("systemctl", "stop", service)
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Service stopped"})
}

//Restart service
func Restart(c *gin.Context) {
	var err error

	service := c.DefaultQuery("service", "")

	if len(service) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Service name empty"})
		return
	}

	cmd := exec.Command("systemctl", "restart", service)
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Service reloaded"})
}
