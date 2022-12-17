package cmonitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	procmeminfo "github.com/guillermo/go.procmeminfo"
	"github.com/shirou/gopsutil/disk"
	"github.com/speshiy/LinuxMonitorGo/_monitoring/models/mmonitor"
	"github.com/speshiy/LinuxMonitorGo/database"
	"github.com/speshiy/LinuxMonitorGo/settings"
)

//GetMonitors return all records
func GetMonitors(c *gin.Context) {
	var err error
	var monitors []mmonitor.Monitor

	DBMonitor, err := database.ReturnConnectionToDBMonitoring()
	if err != nil {
		return
	}
	defer DBMonitor.Close()

	err = mmonitor.GetMonitors(DBMonitor, &monitors)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": monitors})
}

//PostRAM into database directly
func PostRAM() {
	if settings.Server == "" {
		return
	}

	DBMonitor, err := database.ReturnConnectionToDBMonitoring()
	if err != nil {
		return
	}
	defer DBMonitor.Close()

	var monitor mmonitor.Monitor

	meminfo := &procmeminfo.MemInfo{}
	meminfo.Update()

	monitor.Server = settings.Server
	err = monitor.GetByServer(DBMonitor)
	if err != nil {
		return
	}

	monitor.IP = settings.ServerIP
	monitor.Server = settings.Server
	monitor.IsBackupServer = settings.IsBackupServer
	monitor.IsDatabaseServer = settings.IsDatabaseServer

	monitor.RAM = meminfo.Total()
	monitor.RAMUsed = meminfo.Used()

	swapTotal := (*meminfo)["SwapTotal"]
	swapFree := (*meminfo)["SwapFree"]
	monitor.Swap = swapTotal
	monitor.SwapUsed = swapTotal - swapFree

	var total uint64
	var totalUsed uint64

	parts, _ := disk.Partitions(false)
	for _, p := range parts {
		device := p.Mountpoint
		s, _ := disk.Usage(device)

		if s.Total == 0 {
			continue
		}

		total += s.Total
		totalUsed += s.Used
	}

	monitor.Disk = total / 1024 / 1024 / 1024
	monitor.DiskUsed = totalUsed / 1024 / 1024 / 1024

	if monitor.ID == 0 {
		err = monitor.Post(DBMonitor)
	} else {
		err = monitor.Put(DBMonitor)
	}

	if err != nil {
		return
	}

}

//Post into monitor info from another server
func Post(c *gin.Context) {
	var err error
	var monitor mmonitor.Monitor
	var monitorInDB mmonitor.Monitor

	DBMonitor, err := database.ReturnConnectionToDBMonitoring()
	if err != nil {
		return
	}
	defer DBMonitor.Close()

	if err := c.ShouldBindJSON(&monitor); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	if monitor.Server == "" {
		return
	}

	if monitor.IP == "" {
		return
	}

	monitorInDB.Server = monitor.Server
	err = monitorInDB.GetByServer(DBMonitor)
	if err != nil {
		return
	}

	if monitorInDB.ID == 0 {
		err = monitor.Post(DBMonitor)
	} else {
		monitor.ID = monitorInDB.ID
		err = monitor.Put(DBMonitor)
	}
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, nil)
}
